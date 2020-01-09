package model

import (
	"fmt"
	"image"

	"nes-go/pkg/log"
)

// PPURegisters ...
type PPURegisters struct {
	ppuctrl   byte // 0x2000	PPUCTRL	W	コントロールレジスタ1	割り込みなどPPUの設定
	ppumask   byte // 0x2001	PPUMASK	W	コントロールレジスタ2	背景イネーブルなどのPPU設定
	ppustatus byte // 0x2002	PPUSTATUS	R	PPUステータス	PPUのステータス
	oamaddr   byte // 0x2003	OAMADDR	W	スプライトメモリデータ	書き込むスプライト領域のアドレス
	oamdata   byte // 0x2004	OAMDATA	RW	デシマルモード	スプライト領域のデータ
	ppuscroll byte // 0x2005	PPUSCROLL	W	背景スクロールオフセット	背景スクロール値
	ppuaddr   byte // 0x2006	PPUADDR	W	PPUメモリアドレス	書き込むPPUメモリ領域のアドレス
}

// NewPPURegisters ...
func NewPPURegisters() *PPURegisters {
	return &PPURegisters{
		ppuctrl:   0,
		ppumask:   0,
		ppustatus: 0,
		oamaddr:   0,
		oamdata:   0,
		ppuscroll: 0,
		ppuaddr:   0,
	}
}

// String ...
func (r *PPURegisters) String() string {
	return fmt.Sprintf(
		"{PPUCTRL:%#v, PPUMASK:%#v, PPUSTATUS:%#v, OAMADDR:%#v, OAMDATA:%v, PPUSCROLL:%#v, PPUADDR:%#v}",
		r.ppuctrl,
		r.ppumask,
		r.ppustatus,
		r.oamaddr,
		r.oamdata,
		r.ppuscroll,
		r.ppuaddr,
	)
}

// PPU ...
type PPU struct {
	registers *PPURegisters
	bus       *Bus

	ppuaddrWriteCount uint8   // PPUADDRへの書き込み回数（0→1→2→1→2→...と遷移）
	ppuaddrBuf        Address // 組み立て中のPPUADDR
	ppuaddrFull       Address // 組み立て済のPPUADDR

	spriteImages [][]SpriteImage

	drawingPoint *image.Point
}

// NewPPU ...
func NewPPU() (*PPU, error) {
	sizeY := ResolutionHeight / SpriteHeight
	sizeX := ResolutionWidth / SpriteWidth
	spriteImages := make([][]SpriteImage, sizeY)
	for y := 0; y < sizeY; y++ {
		spriteImages[y] = make([]SpriteImage, sizeX)
		for x := 0; x < sizeX; x++ {
			spriteImages[y][x] = *NewSpriteImage()
		}
	}

	return &PPU{
		registers:         NewPPURegisters(),
		ppuaddrWriteCount: 0,
		ppuaddrBuf:        0,
		ppuaddrFull:       0,
		spriteImages:      spriteImages,
		drawingPoint:      &image.Point{0, 0},
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
	if (p.registers.ppuctrl & 0x04) == 0 {
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
func (p *PPU) ReadRegisters(addr Address) (byte, error) {
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
		data = p.registers.ppustatus
	case 3:
		target = "OAMADDR"
		err = fmt.Errorf("failed to read, PPURegister[OAMADDR] is write only; addr: %#v", addr)
	case 4:
		target = "OAMDATA"
		data = p.registers.oamdata
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
func (p *PPU) WriteRegisters(addr Address, data byte) error {
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
		p.registers.ppuctrl = data
		target = "PPUCTRL"
	case 1:
		p.registers.ppumask = data
		target = "PPUMASK"
	case 2:
		err = fmt.Errorf("failed to write, PPURegister[PPUSTATUS] is read only; addr: %#v", addr)
		target = "PPUSTATUS"
	case 3:
		p.registers.oamaddr = data
		target = "OAMADDR"
	case 4:
		p.registers.oamdata = data
		target = "OAMDATA"
	case 5:
		p.registers.ppuscroll = data
		target = "PPUSCROLL"
	case 6:
		p.registers.ppuaddr = data
		switch p.ppuaddrWriteCount {
		case 0, 2:
			p.ppuaddrBuf = Address(p.registers.ppuaddr) << 8
			p.ppuaddrWriteCount = 1
			target = fmt.Sprintf("PPUADDR(for high 8 bits(ppuaddr:%#v))", p.ppuaddrBuf)
		case 1:
			p.ppuaddrBuf = p.ppuaddrBuf + Address(p.registers.ppuaddr)
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
func (p *PPU) Run(cycle int) (si [][]SpriteImage, err error) {
	for i := 0; i < cycle; i++ {
		if si, err = p.Run1Cycle(); err != nil {
			return
		}
	}
	return
}

// Run1Cycle ...
func (p *PPU) Run1Cycle() ([][]SpriteImage, error) {
	log.Trace("PPU.Run[(x,y)=%v] ...", p.drawingPoint.String())

	defer p.updateDrawingPoint()

	if p.drawingPoint.X >= ResolutionWidth {
		return nil, nil
	}
	if p.drawingPoint.Y >= ResolutionHeight {
		return nil, nil
	}

	// 8ライン単位で書き込む
	shouldDrawline := (p.drawingPoint.X == ResolutionWidth-1) && (p.drawingPoint.Y%8 == 7)
	if shouldDrawline {
		y := p.drawingPoint.Y / SpriteHeight
		for x := 0; x < 0x20; x++ {
			np := NameTablePoint{X: uint8(x), Y: uint8(y)}
			spriteNo, err := p.bus.GetSpriteNo(np)
			if err != nil {
				return nil, err
			}

			paletteNo, err := p.bus.GetPaletteNo(np)
			if err != nil {
				return nil, err
			}

			sprite := p.bus.GetSprite(spriteNo)
			palette := p.bus.GetBackgroundPalette(paletteNo)

			// 書き込む
			si := sprite.ToSpriteImage(palette)
			p.spriteImages[y][x] = *si
		}
	}

	// 1ライン分の書き込み完了直後以外は描画しない
	if p.drawingPoint.X != ResolutionWidth-1 {
		return nil, nil
	}

	// TODO この判定を入れると急激に遅くなるためコメントアウト
	// 1画面分の書き込み完了直後以外は描画しない
	// if p.drawingPoint.Y != ResolutionHeight-1 {
	// 	return nil, nil
	// }

	return p.spriteImages, nil
}
