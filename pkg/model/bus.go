package model

import (
	"fmt"

	"github.com/canpok1/nes-go/pkg/log"
)

// Palette ...
type Palette []byte

// NewPalette ...
func NewPalette() *Palette {
	p := Palette(make([]byte, 4))
	return &p
}

// GetColor ...
func (p *Palette) GetColor(no uint8) (uint8, uint8, uint8) {
	index := (*p)[no]
	c := colors[index]
	return c[0], c[1], c[2]
}

// Bus ...
type Bus struct {
	wram       []byte
	wramMirror []byte
	io         []byte
	exrom      []byte
	exram      []byte
	programROM *PRGROM

	ppu *PPU

	charactorROM      *CHRROM
	nameTable0        []byte
	nameTable1        []byte
	nameTable2        []byte
	nameTable3        []byte
	backgroundPalette []Palette
	spritePalette     []Palette

	setupped bool
}

// NewBus ...
func NewBus() *Bus {
	bp := []Palette{}
	sp := []Palette{}
	for i := 0; i < 4; i++ {
		newBP := NewPalette()
		newSP := NewPalette()
		bp := append(bp, *newBP)
		sp := append(sp, *newSP)
	}

	return &Bus{
		wram:  make([]byte, 0x0800),
		io:    make([]byte, 0x0020),
		exrom: make([]byte, 0x1FE0),
		exram: make([]byte, 0x2000),

		backgroundPalette: bp,
		spritePalette:     sp,

		setupped: false,
	}
}

// Setup ...
func (b *Bus) Setup(rom *ROM, ppu *PPU) {
	b.programROM = rom.Prgrom
	b.charactorROM = rom.Chrrom
	b.ppu = ppu

	b.setupped = true
}

// readByCPU ...
func (b *Bus) readByCPU(addr Address) (byte, error) {
	var data byte
	var err error
	var target string
	log.Debug("Bus.readByCPU[addr=%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Debug("Bus.readByCPU[addr=%#v] => %#v", addr, err)
		} else {
			log.Debug("Bus.readByCPU[addr=%#v][%v] => %#v", addr, target, data)
		}
	}()

	if b == nil {
		err = fmt.Errorf("failed to readByCPU, bus is nil")
		return data, err
	}
	if !b.setupped {
		err = fmt.Errorf("failed to readByCPU, bus setup is not completed")
		return data, err
	}

	// 0x0000～0x07FF	0x0800	WRAM
	if addr >= 0x0000 && addr <= 0x07FF {
		target = "WRAM"
		data = b.wram[addr]
		return data, err
	}

	// 0x0800～0x1FFF	-	WRAMのミラー
	if addr >= 0x0800 && addr <= 0x1FFF {
		target = "WRAM Mirror"
		data = b.wram[addr-0x0800]
		return data, err
	}

	// 0x2000～0x2007	0x0008	PPU レジスタ
	if addr >= 0x2000 && addr <= 0x2007 {
		target = "PPU Register"
		data, err = b.ppu.ReadRegisters(addr - 0x2000)
		return data, err
	}

	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	if addr >= 0x2008 && addr <= 0x3FFF {
		target = "PPU Register Mirror"
		data, err = b.ppu.ReadRegisters(addr - 0x2008)
		return data, err
	}

	// 0x4000～0x401F	0x0020	APU I/O、PAD
	if addr >= 0x4000 && addr <= 0x401F {
		target = "APU I/O, PAD"
		data = b.io[addr-0x4000]
		return data, err
	}

	// 0x4020～0x5FFF	0x1FE0	拡張ROM
	if addr >= 0x4000 && addr <= 0x401F {
		target = "EX ROM"
		data = b.exrom[addr-0x4000]
		return data, err
	}

	// 0x6000～0x7FFF	0x2000	拡張RAM
	if addr >= 0x6000 && addr <= 0x7FFF {
		target = "EX RAM"
		data = b.exram[addr-0x6000]
		return data, err
	}

	// 0x8000～0xBFFF	0x4000	PRG-ROM
	// 0xC000～0xFFFF	0x4000	PRG-ROM
	if addr >= 0x8000 && addr <= 0xFFFF {
		target = "PRG-ROM"
		r := *b.programROM
		if len(r) <= 0x4000 {
			data = r[addr-0xC000]
		} else {
			data = r[addr-0x8000]
		}
		return data, err
	}

	return 0, fmt.Errorf("failed read, addr out of range; addr: %#v", addr)
}

// writeByCPU ...
func (b *Bus) writeByCPU(addr Address, data byte) error {
	var err error
	var target string
	log.Debug("Bus.writeByCPU[addr=%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Debug("Bus.writeByCPU[addr=%#v] => %#v", addr, err)
		} else {
			log.Debug("Bus.writeByCPU[addr=%#v][%v] <= %#v", addr, target, data)
		}
	}()

	if b == nil {
		err = fmt.Errorf("failed to writeByCPU, bus is nil")
		return err
	}
	if !b.setupped {
		err = fmt.Errorf("failed to writeByCPU, bus setup is not completed")
		return err
	}

	// 0x0000～0x07FF	0x0800	WRAM
	if addr >= 0x0000 && addr <= 0x07FF {
		target = "WRAM"
		b.wram[addr] = data
		return nil
	}

	// 0x0800～0x1FFF	-	WRAMのミラー
	if addr >= 0x0800 && addr <= 0x1FFF {
		target = "WRAM Mirror"
		b.wram[addr-0x0800] = data
		return nil
	}

	// 0x2000～0x2007	0x0008	PPU レジスタ
	if addr >= 0x2000 && addr <= 0x2007 {
		target = "PPU Register"
		err = b.ppu.WriteRegisters(addr-0x2000, data)
		return err
	}

	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	if addr >= 0x2008 && addr <= 0x3FFF {
		target = "PPU Register Mirror"
		err = b.ppu.WriteRegisters(addr-0x2008, data)
		return err
	}

	// 0x4000～0x401F	0x0020	APU I/O、PAD
	if addr >= 0x4000 && addr <= 0x401F {
		target = "APU I/O, PAD"
		b.io[addr-0x4000] = data
		return nil
	}

	// 0x4020～0x5FFF	0x1FE0	拡張ROM
	if addr >= 0x4000 && addr <= 0x401F {
		return fmt.Errorf("failed write, cannot write EX ROM; addr: %#v", addr)
	}

	// 0x6000～0x7FFF	0x2000	拡張RAM
	if addr >= 0x6000 && addr <= 0x7FFF {
		target = "EX RAM"
		b.exram[addr-0x6000] = data
		return nil
	}

	// 0x8000～0xBFFF	0x4000	PRG-ROM
	// 0xC000～0xFFFF	0x4000	PRG-ROM
	if addr >= 0x8000 && addr <= 0xFFFF {
		return fmt.Errorf("failed write, cannot write PRG-ROM; addr: %#v", addr)
	}

	return fmt.Errorf("failed write, addr out of range; addr: %#v", addr)
}

// readByPPU ...
func (b *Bus) readByPPU(addr Address) (byte, error) {
	var data byte
	var err error
	var target string
	log.Debug("Bus.readByPPU[addr=%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Debug("Bus.readByPPU[addr=%#v] => %#v", addr, err)
		} else {
			log.Debug("Bus.readByPPU[addr=%#v][%v] => %#v", addr, target, data)
		}
	}()

	if b == nil {
		err = fmt.Errorf("failed to readByPPU, bus is nil")
		return data, err
	}
	if !b.setupped {
		err = fmt.Errorf("failed to readByPPU, bus setup is not completed")
		return data, err
	}

	// TODO 実装
	// 0x0000～0x0FFF	0x1000	パターンテーブル0
	// 0x1000～0x1FFF	0x1000	パターンテーブル1
	// 0x2000～0x23BF	0x03C0	ネームテーブル0
	// 0x23C0～0x23FF	0x0040	属性テーブル0
	// 0x2400～0x27BF	0x03C0	ネームテーブル1
	// 0x27C0～0x27FF	0x0040	属性テーブル1
	// 0x2800～0x2BBF	0x03C0	ネームテーブル2
	// 0x2BC0～0x2BFF	0x0040	属性テーブル2
	// 0x2C00～0x2FBF	0x03C0	ネームテーブル3
	// 0x2FC0～0x2FFF	0x0040	属性テーブル3
	// 0x3000～0x3EFF	-	0x2000-0x2EFFのミラー
	// 0x3F00～0x3F0F	0x0010	バックグラウンドパレット
	// 0x3F10～0x3F1F	0x0010	スプライトパレット
	// 0x3F20～0x3FFF	-	0x3F00-0x3F1Fのミラー

	return data, err
}

// writeByPPU ...
func (b *Bus) writeByPPU(addr Address, data byte) error {
	var err error
	var target string
	log.Debug("Bus.writeByPPU[addr=%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Debug("Bus.writeByPPU[addr=%#v] => %#v", addr, err)
		} else {
			log.Debug("Bus.writeByPPU[addr=%#v][%v] <= %#v", addr, target, data)
		}
	}()

	if b == nil {
		err = fmt.Errorf("failed to writeByPPU, bus is nil")
		return err
	}
	if !b.setupped {
		err = fmt.Errorf("failed to writeByPPU, bus setup is not completed")
		return err
	}

	// TODO 実装
	// 0x0000～0x0FFF	0x1000	パターンテーブル0
	// 0x1000～0x1FFF	0x1000	パターンテーブル1
	// 0x2000～0x23BF	0x03C0	ネームテーブル0
	// 0x23C0～0x23FF	0x0040	属性テーブル0
	// 0x2400～0x27BF	0x03C0	ネームテーブル1
	// 0x27C0～0x27FF	0x0040	属性テーブル1
	// 0x2800～0x2BBF	0x03C0	ネームテーブル2
	// 0x2BC0～0x2BFF	0x0040	属性テーブル2
	// 0x2C00～0x2FBF	0x03C0	ネームテーブル3
	// 0x2FC0～0x2FFF	0x0040	属性テーブル3
	// 0x3000～0x3EFF	-	0x2000-0x2EFFのミラー
	// 0x3F00～0x3F0F	0x0010	バックグラウンドパレット
	// 0x3F10～0x3F1F	0x0010	スプライトパレット
	// 0x3F20～0x3FFF	-	0x3F00-0x3F1Fのミラー

	return err
}

// GetSprite ...
func (b *Bus) GetSprite(no uint8) *Sprite {
	return b.charactorROM.GetSprite(no)
}

// GetBackgroundPalette ...
func (b *Bus) GetBackgroundPalette(no uint8) *Palette {
	return &b.backgroundPalette[no]
}

// GetSpritePalette ...
func (b *Bus) GetSpritePalette(no uint8) *Palette {
	return &b.spritePalette[no]
}
