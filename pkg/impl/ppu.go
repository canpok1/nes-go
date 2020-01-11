package impl

import (
	"fmt"
	"image"
	"nes-go/pkg/domain"
	"nes-go/pkg/log"
)

// PPU ...
type PPU struct {
	registers *PPURegisters
	bus       *Bus

	ppuaddrWriteCount uint8          // PPUADDRへの書き込み回数（0→1→2→1→2→...と遷移）
	ppuaddrBuf        domain.Address // 組み立て中のPPUADDR
	ppuaddrFull       domain.Address // 組み立て済のPPUADDR

	tileImages   [][]domain.TileImage
	spriteImages []domain.SpriteImage

	drawingPoint *image.Point

	oam *PPUOAM
}

// NewPPU ...
func NewPPU() (*PPU, error) {
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
		spriteImages:      []domain.SpriteImage{},
		drawingPoint:      &image.Point{0, 0},
		oam:               NewPPUOAM(),
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
func (p *PPU) SetBus(b *Bus) {
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
	case 0:
		target = "PPUCTRL"
		err = fmt.Errorf("failed to read, PPURegister[PPUCTRL] is write only; addr: %#v", addr)
	case 1:
		target = "PPUMASK"
		err = fmt.Errorf("failed to read, PPURegister[PPUMASK] is write only; addr: %#v", addr)
	case 2:
		target = "PPUSTATUS"
		data = p.registers.PPUStatus.ToByte()
		p.registers.PPUStatus.VBlankHasStarted = false
	case 3:
		target = "OAMADDR"
		err = fmt.Errorf("failed to read, PPURegister[OAMADDR] is write only; addr: %#v", addr)
	case 4:
		target = "OAMDATA"
		data = p.registers.OAMData
	case 5:
		target = "PPUSCROLL"
		err = fmt.Errorf("failed to read, PPURegister[PPUSCROLL] is write only; addr: %#v", addr)
	case 6:
		target = "PPUADDR"
		err = fmt.Errorf("failed to read, PPURegister[PPUADDR] is write only; addr: %#v", addr)
	case 7:
		target = fmt.Sprintf("PPUDATA(from PPU Memory %#v)", p.ppuaddrFull)
		data, err = p.bus.readByPPU(p.ppuaddrFull)
		p.incrementPPUADDR()
	default:
		target = "-"
		err = fmt.Errorf("failed to read PPURegisters, address is out of range; addr: %#v", addr)
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
	case 0:
		p.registers.PPUCtrl.UpdateAll(data)
		target = "PPUCTRL"
	case 1:
		p.registers.PPUMask.UpdateAll(data)
		target = "PPUMASK"
	case 2:
		err = fmt.Errorf("failed to write, PPURegister[PPUSTATUS] is read only; addr: %#v", addr)
		target = "PPUSTATUS"
	case 3:
		p.registers.OAMAddr = data
		target = "OAMADDR"
	case 4:
		p.oam.Write(p.registers.OAMAddr, data)
		p.registers.OAMAddr = p.registers.OAMAddr + 1
		target = "OAMDATA"
	case 5:
		p.registers.PPUScroll = data
		target = "PPUSCROLL"
	case 6:
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
	case 7:
		target = fmt.Sprintf("PPUDATA(to PPU Memory %#v)", p.ppuaddrFull)
		err = p.bus.writeByPPU(p.ppuaddrFull, data)
		p.incrementPPUADDR()
	default:
		target = "-"
		err = fmt.Errorf("failed to write PPURegisters, address is out of range; addr: %#v", addr)
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
func (p *PPU) Run(cycle int) (tis [][]domain.TileImage, sis []domain.SpriteImage, err error) {
	for i := 0; i < cycle; i++ {
		if tis, sis, err = p.Run1Cycle(); err != nil {
			return
		}
	}
	return
}

// Run1Cycle ...
func (p *PPU) Run1Cycle() ([][]domain.TileImage, []domain.SpriteImage, error) {
	log.Trace("PPU.Run[(x,y)=%v] ...", p.drawingPoint.String())

	defer p.updateDrawingPoint()

	if p.drawingPoint.X == 0 && p.drawingPoint.Y == domain.ResolutionHeight {
		p.registers.PPUStatus.VBlankHasStarted = true
	}
	if p.drawingPoint.X == 0 && p.drawingPoint.Y == 0 {
		p.registers.PPUStatus.VBlankHasStarted = false
	}

	if p.drawingPoint.X == 0 && p.drawingPoint.Y == 0 {
		p.spriteImages = []domain.SpriteImage{}
	}

	if p.registers.PPUCtrl.NMIEnable && p.registers.PPUStatus.VBlankHasStarted {
		if err := p.bus.SendNMI(); err != nil {
			return nil, nil, err
		}
	}

	if p.drawingPoint.X >= domain.ResolutionWidth {
		return nil, nil, nil
	}
	if p.drawingPoint.Y >= domain.ResolutionHeight {
		return nil, nil, nil
	}

	// 8ライン単位で書き込む
	shouldDrawline := (p.drawingPoint.X == domain.ResolutionWidth-1) && (p.drawingPoint.Y%8 == 7)
	if shouldDrawline {
		y := p.drawingPoint.Y / domain.SpriteHeight
		nameTblIdx := p.registers.PPUCtrl.NameTableIndex
		if p.registers.PPUMask.EnableBackground {
			for x := 0; x < 0x20; x++ {
				np := domain.NameTablePoint{X: uint8(x), Y: uint8(y)}
				attribute, err := p.bus.GetAttribute(nameTblIdx, np)
				if err != nil {
					return nil, nil, err
				}

				paletteNo, err := p.bus.GetPaletteNo(np, attribute)
				if err != nil {
					return nil, nil, err
				}

				tileIndex, err := p.bus.GetTileNo(nameTblIdx, np)
				if err != nil {
					return nil, nil, err
				}

				patternTblIdx := p.registers.PPUCtrl.BackgroundPatternTableIndex
				tilePattern := p.bus.GetTilePattern(patternTblIdx, tileIndex)
				palette := p.bus.GetPalette(uint16(paletteNo))

				// 書き込む
				si := tilePattern.ToTileImage(palette)
				p.tileImages[y][x] = *si

			}
		}

		if p.registers.PPUMask.EnableSprite {
			err := p.oam.Each(func(s domain.Sprite) error {
				sy := int(s.Y) / domain.SpriteHeight
				sx := int(s.X) / domain.SpriteWidth
				if y != sy {
					return nil
				}

				np := domain.NameTablePoint{X: uint8(sx), Y: uint8(sy)}
				attribute, err := p.bus.GetAttribute(nameTblIdx, np)
				if err != nil {
					return err
				}

				paletteNo, err := p.bus.GetPaletteNo(np, attribute)
				if err != nil {
					return err
				}

				offset := (uint16(s.Attribute) & 0x03) << 8
				palette := p.bus.GetPalette(uint16(paletteNo) + offset)

				patternTblIdx := p.registers.PPUCtrl.SpritePatternTableIndex
				tilePattern := p.bus.GetTilePattern(patternTblIdx, s.TileIndex)

				// 書き込む
				ti := tilePattern.ToTileImage(palette)
				si := domain.NewSpriteImage(s.X, s.Y, ti)
				p.spriteImages = append(p.spriteImages, *si)

				return nil
			})
			if err != nil {
				return nil, nil, err
			}
		}
	}

	// 1ライン分の書き込み完了直後以外は描画しない
	if p.drawingPoint.X != domain.ResolutionWidth-1 {
		return nil, nil, nil
	}

	// TODO この判定を入れると急激に遅くなるためコメントアウト
	// 1画面分の書き込み完了直後以外は描画しない
	// if p.drawingPoint.Y != ResolutionHeight-1 {
	// 	return nil, nil
	// }

	return p.tileImages, p.spriteImages, nil
}
