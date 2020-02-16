package impl

import (
	"fmt"
	"image/color"
	"nes-go/pkg/domain"
	"nes-go/pkg/impl/component"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// PPU2 ...
type PPU2 struct {
	registers         *component.PPURegisters
	internalRegisters *component.PPUInternalRegisters
	bus               domain.Bus

	bgController *component.BackgroundController
	spController *component.SpriteController

	images [][]color.RGBA

	dot      uint16
	scanline uint16

	enableOAMDMA bool

	rendered bool

	recorder *domain.Recorder
}

// NewPPU2 ...
func NewPPU2() domain.PPU {
	images := make([][]color.RGBA, domain.ResolutionHeight)
	for i := range images {
		images[i] = make([]color.RGBA, domain.ResolutionWidth)
	}

	return &PPU2{
		registers:         component.NewPPURegisters(),
		internalRegisters: component.NewPPUInnerRegisters(),
		bgController:      component.NewBackgroundController(),
		spController:      component.NewSpriteController(),
		images:            images,
		dot:               0,
		scanline:          0,
		enableOAMDMA:      false,
		rendered:          false,
		recorder:          &domain.Recorder{},
	}
}

// String ...
func (p *PPU2) String() string {
	return fmt.Sprintf(
		"PPU Info\nregisters: %v",
		p.registers.String(),
	)
}

// SetBus ...
func (p *PPU2) SetBus(b domain.Bus) {
	p.bus = b
	p.bgController.SetBus(b)
	p.spController.SetBus(b)
}

// SetRecorder ...
func (p *PPU2) SetRecorder(r *domain.Recorder) {
	p.recorder = r
}

// incrementPPUADDR
func (p *PPU2) incrementPPUADDR() {
	if p.registers.PPUCtrl.VRAMAddressIncrementMode == 0 {
		p.registers.PPUAddr.Increment(1)
	} else {
		p.registers.PPUAddr.Increment(32)
	}
}

// ReadRegisters ...
func (p *PPU2) ReadRegisters(addr domain.Address) (byte, error) {
	var data byte
	var err error
	var target string
	log.Trace("begin[%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Warn("end[%#v][%#v] => %#v", addr, target, err)
		} else {
			log.Trace("end[%#v][%#v] => %#v", addr, target, data)
		}
	}()

	if addr < 0x2000 && addr > 0x3FFF {
		target = "-"
		err = xerrors.Errorf("address is out of range; addr: %#v", addr)
		return data, err
	}

	switch addr % 8 {
	case 0:
		target = "PPUCTRL"
		data = p.registers.PPUCtrl.ToByte()
	case 1:
		target = "PPUMASK"
		data = p.registers.PPUMask.ToByte()
	case 2:
		target = "PPUSTATUS"
		data = p.registers.PPUStatus.ToByte()
		p.internalRegisters.ClearW()
		p.registers.PPUStatus.VBlankHasStarted = false
	case 3:
		target = "OAMADDR"
		data = p.registers.OAMAddr
	case 4:
		target = "OAMDATA"
		data = p.registers.OAMData
	case 5:
		target = "PPUSCROLL"
		data = p.registers.PPUScroll.ToByte()
	case 6:
		target = "PPUADDR"
		data = p.registers.PPUAddr.ToByte()
	case 7:
		ppuaddr := p.registers.PPUAddr.ToFullAddress()
		target = fmt.Sprintf("PPUDATA(from PPU Memory %#v)", ppuaddr)
		data, err = p.bus.ReadByPPU(ppuaddr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
		}
		p.incrementPPUADDR()
		p.internalRegisters.IncrementV(p.registers.PPUCtrl)
	default:
		target = "-"
		err = xerrors.Errorf("address is out of range; addr: %#v", addr)
	}

	return data, err
}

// WriteRegisters ...
func (p *PPU2) WriteRegisters(addr domain.Address, data byte) error {
	var err error
	var target string
	log.Trace("begin[%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Warn("end[%#v][%#v] => %#v", addr, target, err)
		} else {
			log.Trace("end[%#v][%#v] <= %#v", addr, target, data)
		}
	}()

	if addr < 0x2000 && addr > 0x3FFF && addr != 0x4014 {
		target = "-"
		err = xerrors.Errorf("address is out of range; addr: %#v", addr)
		return err
	}

	if addr == 0x4014 {
		target = "OAMDMA"
		p.registers.OAMDMA = data
		p.enableOAMDMA = true
		return err
	}

	switch addr % 8 {
	case 0:
		p.registers.PPUCtrl.UpdateAll(data)
		p.internalRegisters.UpdateByPPUCtrl(data)
		target = "PPUCTRL"
	case 1:
		p.registers.PPUMask.UpdateAll(data)
		target = "PPUMASK"
	case 2:
		err = xerrors.Errorf("PPURegister[PPUSTATUS] is read only; addr: %#v", addr)
		target = "PPUSTATUS"
	case 3:
		p.registers.OAMAddr = data
		target = "OAMADDR"
	case 4:
		p.spController.WriteOAM(p.registers.OAMAddr, data)
		p.registers.OAMAddr = p.registers.OAMAddr + 1
		target = "OAMDATA"
	case 5:
		p.registers.PPUScroll.Set(data)
		p.internalRegisters.UpdateByPPUScroll(data)
		target = "PPUSCROLL"
	case 6:
		p.registers.PPUAddr.Set(data)
		p.internalRegisters.UpdateByPPUAddr(data)
		target = "PPUADDR"
	case 7:
		ppuaddr := p.registers.PPUAddr.ToFullAddress()
		target = fmt.Sprintf("PPUDATA(to PPU Memory %#v)", ppuaddr)
		err = p.bus.WriteByPPU(ppuaddr, data)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
		}
		p.incrementPPUADDR()
		p.internalRegisters.IncrementV(p.registers.PPUCtrl)
	}

	return err
}

func (p *PPU2) execOAMDMA() error {
	p.enableOAMDMA = false
	readAddrH := uint16(p.registers.OAMDMA) << 8
	for readAddrL := 0; readAddrL <= 0xFF; readAddrL++ {
		readAddr := domain.Address(readAddrH + uint16(readAddrL))

		// CPUのメモリマップにおけるアドレスからデータを読み込む
		readData, err := p.bus.ReadByCPU(readAddr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return err
		}

		p.spController.WriteOAM(uint8(readAddrL), readData)
	}
	return nil
}

// shift ... 各シフトレジスタのデータをシフト
func (p *PPU2) shift() {
	shouldSkip := true
	if p.dot >= 2 && p.dot <= 257 {
		shouldSkip = false
	}
	if p.dot >= 322 && p.dot <= 337 {
		shouldSkip = false
	}

	if shouldSkip {
		log.Trace("PPU[%v,%v]shift skipped", p.dot, p.scanline)
		return
	}

	p.bgController.Shift()

	log.Trace("PPU[%v,%v]shift completed", p.dot, p.scanline)
}

// setNextData ... 各シフトレジスタに次の値をセット
func (p *PPU2) setNextData() {
	// 9,17,25, ..., 257 もしくは 329, 337 のときだけセット
	shouldSkip := false
	if p.dot < 9 {
		shouldSkip = true
	}
	if p.dot > 257 && p.dot < 329 {
		shouldSkip = true
	}
	if p.dot > 337 {
		shouldSkip = true
	}
	if (p.dot % 8) != 1 {
		shouldSkip = true
	}

	if shouldSkip {
		log.Trace("PPU[%v,%v]set next data skipped", p.dot, p.scanline)
		return
	}

	attrIdx := p.internalRegisters.GetAttributeIndex()
	p.bgController.SetupNextData(attrIdx)

	log.Trace("PPU[%v,%v]set next data completed", p.dot, p.scanline)
}

// updatePixel ...
func (p *PPU2) updatePixel() {
	y := p.scanline
	x := p.dot - 1

	shouldSkip := false
	if y >= domain.ResolutionHeight {
		shouldSkip = true
	}
	if x >= domain.ResolutionWidth {
		shouldSkip = true
	}

	if shouldSkip {
		log.Trace("PPU[%v,%v]update pixel skipped", p.dot, p.scanline)
		return
	}

	var bgPixel, spPixel color.RGBA
	var spAttr byte

	if p.registers.PPUMask.EnableBackground {
		bgPixel = p.bgController.MakePixel(p.internalRegisters.GetFineX())
	}
	if p.registers.PPUMask.EnableSprite {
		spPixel, spAttr = p.spController.MakePixel()
	}

	if bgPixel.A == 0 && spPixel.A == 0 {
		p.images[y][x] = color.RGBA{R: 0, G: 0, B: 0, A: 0xFF}
		return
	}
	if bgPixel.A == 0 && spPixel.A != 0 {
		p.images[y][x] = spPixel
		return
	}
	if bgPixel.A != 0 && spPixel.A == 0 {
		p.images[y][x] = bgPixel
		return
	}
	if (spAttr & 0x08) == 0x08 {
		p.images[y][x] = bgPixel
		return
	}
	p.images[y][x] = spPixel

	//log.Trace("PPU[%v,%v]update pixel completed (x,y)=(%v,%v), (r,g,b)=(%v,%v,%v)", p.dot, p.scanline, x, y, pixel.R, pixel.G, pixel.B)
}

func (p *PPU2) incrementHorizontal() error {
	p.internalRegisters.IncrementHorizontal()
	return nil
}

func (p *PPU2) incrementVertical() error {
	p.internalRegisters.IncrementVertical()
	return nil
}

func (p *PPU2) updateHorizontalToLeftEdge() error {
	p.internalRegisters.UpdateHorizontalToLeftEdge()
	return nil
}

func (p *PPU2) updateVerticalToTopEdge() error {
	p.internalRegisters.UpdateVerticalToTopEdge()
	return nil
}

func (p *PPU2) fetchNTByte() error {
	log.Trace("PPU[%v,%v] fetch NT byte ...", p.dot, p.scanline)

	addr := p.internalRegisters.GetTileIndexAddress()

	var err error
	if p.bgController.NextTileIndex, err = p.bus.ReadByPPU(addr); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	log.Trace("PPU[%v,%v] fetch NT byte (addr: %#v) => %#v", p.dot, p.scanline, addr, p.bgController.NextTileIndex)
	return nil
}

func (p *PPU2) fetchATByte() error {
	log.Trace("PPU[%v,%v] fetch AT byte ...", p.dot, p.scanline)

	addr := p.internalRegisters.GetAttributeAddress()

	var err error
	if p.bgController.NextAttributeTable, err = p.bus.ReadByPPU(addr); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	log.Trace("PPU[%v,%v] fetch AT byte (addr: %#v) => %#v", p.dot, p.scanline, addr, p.bgController.NextAttributeTable)
	return nil
}

func (p *PPU2) fetchLowBGTileByte() error {
	log.Trace("PPU[%v,%v] fetch Low BG Tile byte ...", p.dot, p.scanline)

	addr := p.internalRegisters.GetTilePatternLowAddress(p.bgController.NextTileIndex)

	var err error
	if p.bgController.NextTilePatternL, err = p.bus.ReadByPPU(addr); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	log.Trace("PPU[%v,%v] fetch Low BG Tile byte (addr: %#v) => %#v", p.dot, p.scanline, addr, p.bgController.NextTilePatternL)
	return nil
}

func (p *PPU2) fetchHighBGTileByte() error {
	log.Trace("PPU[%v,%v] fetch High BG Tile byte ...", p.dot, p.scanline)

	addr := p.internalRegisters.GetTilePatternHighAddress(p.bgController.NextTileIndex)

	var err error
	if p.bgController.NextTilePatternH, err = p.bus.ReadByPPU(addr); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	log.Trace("PPU[%v,%v] fetch High BG Tile byte (addr: %#v) => %#v", p.dot, p.scanline, addr, p.bgController.NextTilePatternH)
	return nil
}

// setVBlankFlag ...
func (p *PPU2) setVBlankFlag() error {
	p.registers.PPUStatus.VBlankHasStarted = true
	return nil
}

// clearFlags ...
func (p *PPU2) clearFlags() error {
	p.registers.PPUStatus.VBlankHasStarted = false
	p.rendered = false
	return nil
}

// clearSecondaryOAM ...
func (p *PPU2) clearSecondaryOAM() error {
	spriteIdx := (p.dot & 0x1C) >> 2
	byteIdx := p.dot & 0x03
	p.spController.ClearSecondaryOAM(spriteIdx, byteIdx)
	return nil
}

// evaluateSprite ...
func (p *PPU2) evaluateSprite() error {
	p.spController.EvaluateSprite(p.scanline)
	return nil
}

// fetchSprite ...
func (p *PPU2) fetchSprite() error {
	if (p.dot % 8) != 0 {
		return nil
	}

	p.spController.FetchSprite(p.scanline, p.registers.PPUCtrl.SpritePatternTableIndex)
	return nil
}

// updateSpriteController ...
func (p *PPU2) updateSpriteController() error {

	if p.scanline >= 240 && p.scanline <= 260 {
		return nil
	}

	if p.scanline <= 239 && p.dot >= 1 && p.dot <= 256 {
		p.spController.Shift()
	}

	if p.dot >= 1 && p.dot <= 64 {
		if err := p.clearSecondaryOAM(); err != nil {
			return xerrors.Errorf(": %w", err)
		}
		return nil
	}
	if p.dot >= 65 && p.dot <= 256 {
		if err := p.evaluateSprite(); err != nil {
			return xerrors.Errorf(": %w", err)
		}
		return nil
	}
	if p.dot >= 257 && p.dot <= 320 {
		if err := p.fetchSprite(); err != nil {
			return xerrors.Errorf(": %w", err)
		}
		return nil
	}

	return nil
}

// run1Cycle ...
func (p *PPU2) run1Cycle() error {
	log.Trace("PPU[%v,%v] run start", p.dot, p.scanline)
	defer log.Trace("PPU[%v,%v] run end: internal register : %#v", p.dot, p.scanline, *p.internalRegisters)

	p.shift()
	p.setNextData()
	p.updatePixel()

	p.bus.SendNMI(p.registers.PPUCtrl.NMIEnable && p.registers.PPUStatus.VBlankHasStarted)

	// 次の仕様にしたがって更新
	// http://wiki.nesdev.com/w/index.php/PPU_rendering

	if err := p.updateSpriteController(); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	// Visible scanlines (0-239)
	if p.scanline >= 0 && p.scanline <= 239 {
		if p.dot == 0 {
			return nil
		}

		if p.dot >= 258 && p.dot <= 320 {
			return nil
		}

		if p.dot == 256 {
			if err := p.incrementVertical(); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		}
		if p.dot == 257 {
			if err := p.updateHorizontalToLeftEdge(); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		}
		// if p.dot >= 249 && p.dot <= 255 {
		// 	// unused tile fetch
		// 	return nil
		// }

		switch p.dot % 8 {
		case 0:
			if err := p.incrementHorizontal(); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		case 1:
			if err := p.fetchNTByte(); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		case 3:
			if err := p.fetchATByte(); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		case 5:
			if err := p.fetchLowBGTileByte(); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		case 7:
			if err := p.fetchHighBGTileByte(); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		default:
			return nil
		}
	}

	if p.scanline == 241 && p.dot == 1 {
		if err := p.setVBlankFlag(); err != nil {
			return xerrors.Errorf(": %w", err)
		}
		return nil
	}

	// Pre-render line (261)
	if p.scanline == 261 {
		if p.dot == 1 {
			if err := p.clearFlags(); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		}

		if p.dot == 257 {
			if err := p.updateHorizontalToLeftEdge(); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		}
		if p.dot >= 280 && p.dot <= 304 {
			if err := p.updateVerticalToTopEdge(); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		}

		if (p.dot >= 321 && p.dot <= 336) || (p.dot >= 2 && p.dot <= 256) {
			switch p.dot % 8 {
			case 0:
				if err := p.incrementHorizontal(); err != nil {
					return xerrors.Errorf(": %w", err)
				}
				return nil
			case 1:
				if err := p.fetchNTByte(); err != nil {
					return xerrors.Errorf(": %w", err)
				}
				return nil
			case 3:
				if err := p.fetchATByte(); err != nil {
					return xerrors.Errorf(": %w", err)
				}
				return nil
			case 5:
				if err := p.fetchLowBGTileByte(); err != nil {
					return xerrors.Errorf(": %w", err)
				}
				return nil
			case 7:
				if err := p.fetchHighBGTileByte(); err != nil {
					return xerrors.Errorf(": %w", err)
				}
				return nil
			default:
				return nil
			}
		}
		return nil
	}

	return nil
}

// Run ... 指定サイクル数だけ実行
func (p *PPU2) Run(cycle int) (*domain.Screen, error) {
	defer func() {
		p.recorder.Dot = p.dot
		p.recorder.Scanline = p.scanline
	}()

	for i := 0; i < cycle; i++ {
		if p.enableOAMDMA {
			if err := p.execOAMDMA(); err != nil {
				return nil, xerrors.Errorf(": %w", err)
			}
			continue
		}

		if err := p.run1Cycle(); err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		if p.dot < 340 {
			p.dot++
			continue
		}

		p.dot = 0

		if p.scanline < 261 {
			p.scanline++
			continue
		}

		p.scanline = 0
	}

	// post-render line のときだけ1回返す
	if p.scanline == 240 && !p.rendered {
		log.Trace("PPU[%v,%v] return images", p.dot, p.scanline)
		p.rendered = true
		return &domain.Screen{
			TileImages:            nil,
			SpriteImages:          nil,
			DisableSpriteMask:     false,
			DisableBackgroundMask: false,
			Images:                p.images,
		}, nil
	}
	return nil, nil
}

var colors = []color.RGBA{
	{R: 0x80, G: 0x80, B: 0x80}, {R: 0x00, G: 0x3D, B: 0xA6}, {R: 0x00, G: 0x12, B: 0xB0}, {R: 0x44, G: 0x00, B: 0x96},
	{R: 0xA1, G: 0x00, B: 0x5E}, {R: 0xC7, G: 0x00, B: 0x28}, {R: 0xBA, G: 0x06, B: 0x00}, {R: 0x8C, G: 0x17, B: 0x00},
	{R: 0x5C, G: 0x2F, B: 0x00}, {R: 0x10, G: 0x45, B: 0x00}, {R: 0x05, G: 0x4A, B: 0x00}, {R: 0x00, G: 0x47, B: 0x2E},
	{R: 0x00, G: 0x41, B: 0x66}, {R: 0x00, G: 0x00, B: 0x00}, {R: 0x05, G: 0x05, B: 0x05}, {R: 0x05, G: 0x05, B: 0x05},
	{R: 0xC7, G: 0xC7, B: 0xC7}, {R: 0x00, G: 0x77, B: 0xFF}, {R: 0x21, G: 0x55, B: 0xFF}, {R: 0x82, G: 0x37, B: 0xFA},
	{R: 0xEB, G: 0x2F, B: 0xB5}, {R: 0xFF, G: 0x29, B: 0x50}, {R: 0xFF, G: 0x22, B: 0x00}, {R: 0xD6, G: 0x32, B: 0x00},
	{R: 0xC4, G: 0x62, B: 0x00}, {R: 0x35, G: 0x80, B: 0x00}, {R: 0x05, G: 0x8F, B: 0x00}, {R: 0x00, G: 0x8A, B: 0x55},
	{R: 0x00, G: 0x99, B: 0xCC}, {R: 0x21, G: 0x21, B: 0x21}, {R: 0x09, G: 0x09, B: 0x09}, {R: 0x09, G: 0x09, B: 0x09},
	{R: 0xFF, G: 0xFF, B: 0xFF}, {R: 0x0F, G: 0xD7, B: 0xFF}, {R: 0x69, G: 0xA2, B: 0xFF}, {R: 0xD4, G: 0x80, B: 0xFF},
	{R: 0xFF, G: 0x45, B: 0xF3}, {R: 0xFF, G: 0x61, B: 0x8B}, {R: 0xFF, G: 0x88, B: 0x33}, {R: 0xFF, G: 0x9C, B: 0x12},
	{R: 0xFA, G: 0xBC, B: 0x20}, {R: 0x9F, G: 0xE3, B: 0x0E}, {R: 0x2B, G: 0xF0, B: 0x35}, {R: 0x0C, G: 0xF0, B: 0xA4},
	{R: 0x05, G: 0xFB, B: 0xFF}, {R: 0x5E, G: 0x5E, B: 0x5E}, {R: 0x0D, G: 0x0D, B: 0x0D}, {R: 0x0D, G: 0x0D, B: 0x0D},
	{R: 0xFF, G: 0xFF, B: 0xFF}, {R: 0xA6, G: 0xFC, B: 0xFF}, {R: 0xB3, G: 0xEC, B: 0xFF}, {R: 0xDA, G: 0xAB, B: 0xEB},
	{R: 0xFF, G: 0xA8, B: 0xF9}, {R: 0xFF, G: 0xAB, B: 0xB3}, {R: 0xFF, G: 0xD2, B: 0xB0}, {R: 0xFF, G: 0xEF, B: 0xA6},
	{R: 0xFF, G: 0xF7, B: 0x9C}, {R: 0xD7, G: 0xE8, B: 0x95}, {R: 0xA6, G: 0xED, B: 0xAF}, {R: 0xA2, G: 0xF2, B: 0xDA},
	{R: 0x99, G: 0xFF, B: 0xFC}, {R: 0xDD, G: 0xDD, B: 0xDD}, {R: 0x11, G: 0x11, B: 0x11}, {R: 0x11, G: 0x11, B: 0x11},
}

func swapbit(b byte) byte {
	swapped := byte(0)

	for i := 0; i < 8; i++ {
		swapped = swapped << 1
		if (b & 0x01) == 0x01 {
			swapped = swapped | 0x01
		}
		b = b >> 1
	}

	return swapped
}
