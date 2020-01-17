package impl

import (
	"fmt"
	"image"
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// PPU ...
type PPU struct {
	registers *PPURegisters
	bus       domain.Bus

	ppuaddrWriteCount uint8          // PPUADDRへの書き込み回数（0→1→2→1→2→...と遷移）
	ppuaddrBuf        domain.Address // 組み立て中のPPUADDR
	ppuaddrFull       domain.Address // 組み立て済のPPUADDR

	tileImages [][]domain.TileImage

	drawingPoint *image.Point

	oam *PPUOAM

	enableOAMDMA bool
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
		registers:         NewPPURegisters(),
		ppuaddrWriteCount: 0,
		ppuaddrBuf:        0,
		ppuaddrFull:       0,
		tileImages:        tileImages,
		drawingPoint:      &image.Point{0, 0},
		oam:               NewPPUOAM(),
		enableOAMDMA:      false,
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

// incrementPPUADDR
func (p *PPU) incrementPPUADDR() {
	old := p.ppuaddrFull
	if p.registers.PPUCtrl.VRAMAddressIncrementMode == 0 {
		p.ppuaddrFull = p.ppuaddrFull + 1
	} else {
		p.ppuaddrFull = p.ppuaddrFull + 32
	}
	log.Trace("PPURegisters.update[PPUADDR Full] %#v => %#v", old, p.ppuaddrFull)
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
	log.Trace("PPU.ReadRegisters[%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Warn("PPU.ReadRegisters[%#v][%#v] => %#v", addr, target, err)
		} else {
			log.Trace("PPU.ReadRegisters[%#v][%#v] => %#v", addr, target, data)
		}
	}()

	switch addr {
	case 0x2000:
		target = "PPUCTRL"
		err = xerrors.Errorf("failed to read, PPURegister[PPUCTRL] is write only; addr: %#v", addr)
	case 0x2001:
		target = "PPUMASK"
		err = xerrors.Errorf("failed to read, PPURegister[PPUMASK] is write only; addr: %#v", addr)
	case 0x2002:
		target = "PPUSTATUS"
		data = p.registers.PPUStatus.ToByte()
		p.registers.PPUStatus.VBlankHasStarted = false
	case 0x2003:
		target = "OAMADDR"
		err = xerrors.Errorf("failed to read, PPURegister[OAMADDR] is write only; addr: %#v", addr)
	case 0x2004:
		target = "OAMDATA"
		data = p.registers.OAMData
	case 0x2005:
		target = "PPUSCROLL"
		err = xerrors.Errorf("failed to read, PPURegister[PPUSCROLL] is write only; addr: %#v", addr)
	case 0x2006:
		target = "PPUADDR"
		err = xerrors.Errorf("failed to read, PPURegister[PPUADDR] is write only; addr: %#v", addr)
	case 0x2007:
		target = fmt.Sprintf("PPUDATA(from PPU Memory %#v)", p.ppuaddrFull)
		data, err = p.bus.ReadByPPU(p.ppuaddrFull)
		p.incrementPPUADDR()
	default:
		target = "-"
		err = xerrors.Errorf("failed to read PPURegisters, address is out of range; addr: %#v", addr)
	}

	return data, err
}

// WriteRegisters ...
func (p *PPU) WriteRegisters(addr domain.Address, data byte) error {
	var err error
	var target string
	log.Trace("PPU.WriteRegisters[%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Warn("PPU.WriteRegisters[%#v][%#v] => %#v", addr, target, err)
		} else {
			log.Trace("PPU.WriteRegisters[%#v][%#v] <= %#v", addr, target, data)
		}
	}()

	switch addr {
	case 0x2000:
		p.registers.PPUCtrl.UpdateAll(data)
		target = "PPUCTRL"
	case 0x2001:
		p.registers.PPUMask.UpdateAll(data)
		target = "PPUMASK"
	case 0x2002:
		err = xerrors.Errorf("failed to write, PPURegister[PPUSTATUS] is read only; addr: %#v", addr)
		target = "PPUSTATUS"
	case 0x2003:
		p.registers.OAMAddr = data
		target = "OAMADDR"
	case 0x2004:
		p.oam.Write(p.registers.OAMAddr, data)
		p.registers.OAMAddr = p.registers.OAMAddr + 1
		target = "OAMDATA"
	case 0x2005:
		p.registers.PPUScroll = data
		target = "PPUSCROLL"
	case 0x2006:
		p.registers.PPUAddr = data
		switch p.ppuaddrWriteCount {
		case 0, 2:
			p.ppuaddrBuf = domain.Address(p.registers.PPUAddr) << 8
			p.ppuaddrWriteCount = 1
			target = fmt.Sprintf("PPUADDR(for high 8 bits(ppuaddr:%#v))", p.ppuaddrBuf)
		case 1:
			p.ppuaddrBuf = p.ppuaddrBuf + domain.Address(p.registers.PPUAddr)
			p.ppuaddrFull = p.ppuaddrBuf
			p.ppuaddrWriteCount = 2
			target = fmt.Sprintf("PPUADDR(for low 8 bits(ppuaddr:%#v))", p.ppuaddrBuf)
		}
	case 0x2007:
		target = fmt.Sprintf("PPUDATA(to PPU Memory %#v)", p.ppuaddrFull)
		err = p.bus.WriteByPPU(p.ppuaddrFull, data)
		p.incrementPPUADDR()
	case 0x4014:
		target = "OAMDMA"
		p.registers.OAMDMA = data
		p.enableOAMDMA = true
	default:
		target = "-"
		err = xerrors.Errorf("failed to write PPURegisters, address is out of range; addr: %#v", addr)
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
				return nil, err
			}
			continue
		}

		s, err := p.run1Cycle()
		if err != nil {
			return nil, err
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

	if p.registers.PPUCtrl.NMIEnable && p.registers.PPUStatus.VBlankHasStarted {
		if err := p.bus.SendNMI(); err != nil {
			return nil, err
		}
	}

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
					return nil, err
				}

				paletteNo, err := p.bus.GetPaletteNo(np, attribute)
				if err != nil {
					return nil, err
				}

				tileIndex, err := p.bus.GetTileNo(nameTblIdx, np)
				if err != nil {
					return nil, err
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
				return err
			}

			paletteNo, err := p.bus.GetPaletteNo(np, attribute)
			if err != nil {
				return err
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
			return nil, err
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
			return err
		}

		p.oam.Write(uint8(readAddrL), readData)
	}
	return nil
}
