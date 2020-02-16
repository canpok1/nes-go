package impl

import (
	"fmt"
	"image"
	"nes-go/pkg/domain"
	"nes-go/pkg/impl/component"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// PPU ...
type PPU struct {
	registers *component.PPURegisters
	bus       domain.Bus

	tileImages [][]domain.TileImage

	drawingPoint *image.Point

	oam *component.PPUOAM

	enableOAMDMA bool

	recorder *domain.Recorder
}

// NewPPU ...
func NewPPU() (domain.PPU, error) {
	sizeY := domain.ResolutionHeight / domain.SpriteHeight
	sizeX := domain.ResolutionWidth / domain.SpriteWidth
	tileImages := make([][]domain.TileImage, sizeY)
	for y := 0; y < sizeY; y++ {
		tileImages[y] = make([]domain.TileImage, sizeX)
		for x := 0; x < sizeX; x++ {
			tileImages[y][x] = *domain.NewTileImage()
		}
	}

	return &PPU{
		registers:    component.NewPPURegisters(),
		tileImages:   tileImages,
		drawingPoint: &image.Point{0, 0},
		oam:          component.NewPPUOAM(),
		enableOAMDMA: false,
		recorder:     &domain.Recorder{},
	}, nil
}

// String ...
func (p *PPU) String() string {
	return fmt.Sprintf(
		"PPU Info\nregisters: %v",
		p.registers.String(),
	)
}

// SetBus ...
func (p *PPU) SetBus(b domain.Bus) {
	p.bus = b
}

// SetRecorder ...
func (p *PPU) SetRecorder(r *domain.Recorder) {
	p.recorder = r
}

// incrementPPUADDR
func (p *PPU) incrementPPUADDR() {
	if p.registers.PPUCtrl.VRAMAddressIncrementMode == 0 {
		p.registers.PPUAddr.Increment(1)
	} else {
		p.registers.PPUAddr.Increment(32)
	}
}

// flatten ...
func flatten(org [][]byte) []byte {
	flat := []byte{}
	for _, o := range org {
		flat = append(flat, o...)
	}
	return flat
}

// ReadRegisters ...
func (p *PPU) ReadRegisters(addr domain.Address) (byte, error) {
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
		err = xerrors.Errorf("PPURegister[PPUCTRL] is write only; addr: %#v", addr)
	case 1:
		target = "PPUMASK"
		err = xerrors.Errorf("PPURegister[PPUMASK] is write only; addr: %#v", addr)
	case 2:
		target = "PPUSTATUS"
		data = p.registers.PPUStatus.ToByte()
		p.registers.PPUStatus.VBlankHasStarted = false
	case 3:
		target = "OAMADDR"
		err = xerrors.Errorf("PPURegister[OAMADDR] is write only; addr: %#v", addr)
	case 4:
		target = "OAMDATA"
		data = p.registers.OAMData
	case 5:
		target = "PPUSCROLL"
		err = xerrors.Errorf("PPURegister[PPUSCROLL] is write only; addr: %#v", addr)
	case 6:
		target = "PPUADDR"
		err = xerrors.Errorf("PPURegister[PPUADDR] is write only; addr: %#v", addr)
	case 7:
		ppuaddr := p.registers.PPUAddr.ToFullAddress()
		target = fmt.Sprintf("PPUDATA(from PPU Memory %#v)", ppuaddr)
		data, err = p.bus.ReadByPPU(ppuaddr)
		p.incrementPPUADDR()
	default:
		target = "-"
		err = xerrors.Errorf("address is out of range; addr: %#v", addr)
	}

	return data, err
}

// WriteRegisters ...
func (p *PPU) WriteRegisters(addr domain.Address, data byte) error {
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
		p.oam.Write(p.registers.OAMAddr, data)
		p.registers.OAMAddr = p.registers.OAMAddr + 1
		target = "OAMDATA"
	case 5:
		p.registers.PPUScroll.Set(data)
		target = "PPUSCROLL"
	case 6:
		p.registers.PPUAddr.Set(data)
	case 7:
		ppuaddr := p.registers.PPUAddr.ToFullAddress()
		target = fmt.Sprintf("PPUDATA(to PPU Memory %#v)", ppuaddr)
		err = p.bus.WriteByPPU(ppuaddr, data)
		p.incrementPPUADDR()
	}

	return err
}

// updateDrawingPoint ...
func (p *PPU) updateDrawingPoint() {
	// 1ライン描画は341クロック
	// 1画面は262ライン
	p.drawingPoint.X = p.drawingPoint.X + 1
	if p.drawingPoint.X >= 341 {
		p.drawingPoint.X = 0
		p.drawingPoint.Y = p.drawingPoint.Y + 1
		if p.drawingPoint.Y >= 262 {
			p.drawingPoint.Y = 0
		}
	}
}

// Run ...
func (p *PPU) Run(cycle int) (*domain.Screen, error) {
	var screen *domain.Screen
	for i := 0; i < cycle; i++ {
		if p.enableOAMDMA {
			if err := p.execOAMDMA(); err != nil {
				return nil, xerrors.Errorf(": %w", err)
			}
			continue
		}

		s, err := p.run1Cycle()
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		if s != nil {
			screen = s
		}
	}
	return screen, nil
}

// run1Cycle ...
func (p *PPU) run1Cycle() (*domain.Screen, error) {
	defer p.updateDrawingPoint()

	if p.drawingPoint.X == 0 && p.drawingPoint.Y == domain.ResolutionHeight {
		p.registers.PPUStatus.VBlankHasStarted = true
	}
	if p.drawingPoint.X == 0 && p.drawingPoint.Y == 0 {
		p.registers.PPUStatus.VBlankHasStarted = false
	}

	p.bus.SendNMI(p.registers.PPUCtrl.NMIEnable && p.registers.PPUStatus.VBlankHasStarted)

	if p.drawingPoint.X >= domain.ResolutionWidth {
		return nil, nil
	}
	if p.drawingPoint.Y >= domain.ResolutionHeight {
		return nil, nil
	}

	nameTblIdx := p.registers.PPUCtrl.NameTableIndex

	// 8ライン単位で書き込む
	shouldDrawline := (p.drawingPoint.X == domain.ResolutionWidth-1) && (p.drawingPoint.Y%8 == 7)
	if shouldDrawline {
		y := p.drawingPoint.Y / domain.SpriteHeight
		if p.registers.PPUMask.EnableBackground {
			for x := 0; x < 0x20; x++ {
				np := domain.NameTablePoint{X: uint8(x), Y: uint8(y)}
				attribute, err := p.bus.GetAttribute(nameTblIdx, np)
				if err != nil {
					return nil, xerrors.Errorf(": %w", err)
				}

				paletteNo, err := p.bus.GetPaletteNo(np, attribute)
				if err != nil {
					return nil, xerrors.Errorf(": %w", err)
				}

				tileIndex, err := p.bus.GetTileNo(nameTblIdx, np)
				if err != nil {
					return nil, xerrors.Errorf(": %w", err)
				}

				patternTblIdx := p.registers.PPUCtrl.BackgroundPatternTableIndex
				tilePattern := p.bus.GetTilePattern(patternTblIdx, tileIndex)
				palette := p.bus.GetPalette(paletteNo)

				// 書き込む
				si := tilePattern.ToTileImage(palette)
				p.tileImages[y][x] = *si
			}
		}

	}

	// 1画面分の書き込み完了直後以外は描画しない
	if p.drawingPoint.X != domain.ResolutionWidth-1 || p.drawingPoint.Y != domain.ResolutionHeight-1 {
		return nil, nil
	}

	// TODO この判定を入れると急激に遅くなるためコメントアウト
	// 1画面分の書き込み完了直後以外は描画しない
	// if p.drawingPoint.Y != ResolutionHeight-1 {
	// 	return nil, nil
	// }

	// スプライトを画像に変換
	sImages := []domain.SpriteImage{}
	if p.registers.PPUMask.EnableSprite {
		err := p.oam.Each(func(s domain.Sprite) error {
			if (s.Y+1) == 0 || (s.Y+1) > (domain.ResolutionHeight-1) || s.X > (domain.ResolutionWidth-1) {
				// 画面外なのでスキップ
				return nil
			}
			sy := int(s.Y+1) / domain.SpriteHeight
			sx := int(s.X) / domain.SpriteWidth

			np := domain.NameTablePoint{X: uint8(sx), Y: uint8(sy)}
			attribute, err := p.bus.GetAttribute(nameTblIdx, np)
			if err != nil {
				return xerrors.Errorf(": %w", err)
			}

			paletteNo, err := p.bus.GetPaletteNo(np, attribute)
			if err != nil {
				return xerrors.Errorf(": %w", err)
			}

			offset := (s.Attribute & 0x3) << 2
			palette := p.bus.GetPalette(paletteNo + offset)

			patternTblIdx := p.registers.PPUCtrl.SpritePatternTableIndex
			tilePattern := p.bus.GetTilePattern(patternTblIdx, s.TileIndex)

			ti := tilePattern.ToTileImage(palette)
			isForeground := (s.Attribute & 0x20) == 0x00
			enableFlipH := (s.Attribute & 0x40) == 0x40
			enableFlipV := (s.Attribute & 0x80) == 0x80
			si := domain.NewSpriteImage(uint16(s.X), uint16(s.Y+1), ti, isForeground, enableFlipH, enableFlipV)
			sImages = append(sImages, *si)

			return nil
		})
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
	}

	return &domain.Screen{
		TileImages:            p.tileImages,
		SpriteImages:          sImages,
		DisableSpriteMask:     p.registers.PPUMask.DisableSpriteMask,
		DisableBackgroundMask: p.registers.PPUMask.DisableBackgroundMask,
	}, nil
}

func (p *PPU) execOAMDMA() error {
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

		p.oam.Write(uint8(readAddrL), readData)
	}
	return nil
}
