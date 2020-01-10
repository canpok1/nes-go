package impl

import (
	"fmt"

	"nes-go/pkg/domain"
	"nes-go/pkg/log"
)

// Bus ...
type Bus struct {
	wram       []byte
	wramMirror []byte
	io         []byte
	exrom      []byte
	exram      []byte
	programROM *domain.PRGROM

	ppu *PPU

	charactorROM      *domain.CHRROM
	nameTable0        []byte
	attributeTable0   []byte
	nameTable1        []byte
	attributeTable1   []byte
	nameTable2        []byte
	attributeTable2   []byte
	nameTable3        []byte
	attributeTable3   []byte
	backgroundPalette []domain.Palette
	spritePalette     []domain.Palette

	setupped bool
}

// NewBus ...
func NewBus() *Bus {
	bp := []domain.Palette{}
	sp := []domain.Palette{}
	for i := 0; i < 4; i++ {
		newBP := domain.NewPalette()
		newSP := domain.NewPalette()
		bp = append(bp, *newBP)
		sp = append(sp, *newSP)
	}

	return &Bus{
		wram:  make([]byte, 0x0800),
		io:    make([]byte, 0x0020),
		exrom: make([]byte, 0x1FE0),
		exram: make([]byte, 0x2000),

		nameTable0:      make([]byte, 0x03C0),
		attributeTable0: make([]byte, 0x0040),
		nameTable1:      make([]byte, 0x03C0),
		attributeTable1: make([]byte, 0x0040),
		nameTable2:      make([]byte, 0x03C0),
		attributeTable2: make([]byte, 0x0040),
		nameTable3:      make([]byte, 0x03C0),
		attributeTable3: make([]byte, 0x0040),

		backgroundPalette: bp,
		spritePalette:     sp,

		setupped: false,
	}
}

// Setup ...
func (b *Bus) Setup(rom *domain.ROM, ppu *PPU) {
	b.programROM = rom.Prgrom
	b.charactorROM = rom.Chrrom
	b.ppu = ppu

	b.setupped = true
}

// readByCPU ...
func (b *Bus) ReadByCPU(addr domain.Address) (byte, error) {
	var data byte
	var err error
	var target string
	log.Trace("Bus.readByCPU[addr=%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Warn("Bus.readByCPU[addr=%#v] => %#v", addr, err)
		} else {
			log.Trace("Bus.readByCPU[addr=%#v][%v] => %#v", addr, target, data)
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
	if addr >= 0x8000 && addr <= 0xBFFF {
		target = "PRG-ROM"
		r := *b.programROM
		data = r[addr-0x8000]
		return data, err
	}
	// 0xC000～0xFFFF	0x4000	PRG-ROM
	if addr >= 0xC000 && addr <= 0xFFFF {
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
func (b *Bus) WriteByCPU(addr domain.Address, data byte) error {
	var err error
	var target string
	log.Trace("Bus.writeByCPU[addr=%#v] (<=%#v) ...", addr, data)
	defer func() {
		if err != nil {
			log.Warn("Bus.writeByCPU[addr=%#v] => %#v", addr, err)
		} else {
			log.Trace("Bus.writeByCPU[addr=%#v][%v] <= %#v", addr, target, data)
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
func (b *Bus) readByPPU(addr domain.Address) (data byte, err error) {
	var target string
	log.Trace("Bus.readByPPU[addr=%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Trace("Bus.readByPPU[addr=%#v] => %#v", addr, err)
		} else {
			log.Trace("Bus.readByPPU[addr=%#v][%v] => %#v", addr, target, data)
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

	addrTmp := addr
	// 0x3000～0x3EFF	-	0x2000-0x2EFFのミラー
	if addr >= 0x3000 && addr <= 0x3EFF {
		addrTmp = addr - (0x3000 - 0x2000)
	}
	// 0x3F20～0x3FFF	-	0x3F00-0x3F1Fのミラー
	if addr >= 0x3F20 && addr <= 0x3FFF {
		addrTmp = addr - (0x3F20 - 0x3F00)
	}

	// 0x0000～0x0FFF	0x1000	パターンテーブル0
	if addrTmp >= 0x0000 && addrTmp <= 0x0FFF {
		data = (*b.charactorROM)[addrTmp]
		target = "PatternTable0"
		return
	}

	// 0x1000～0x1FFF	0x1000	パターンテーブル1
	if addrTmp >= 0x1000 && addrTmp <= 0x1FFF {
		data = (*b.charactorROM)[addrTmp]
		target = "PatternTable1"
		return
	}

	// 0x2000～0x23BF	0x03C0	ネームテーブル0
	if addrTmp >= 0x2000 && addrTmp <= 0x23BF {
		data = b.nameTable0[addrTmp-0x2000]
		target = "NameTable0"
		return
	}

	// 0x23C0～0x23FF	0x0040	属性テーブル0
	if addrTmp >= 0x23C0 && addrTmp <= 0x23FF {
		data = b.attributeTable0[addrTmp-0x23C0]
		target = "AttributeTable0"
		return
	}

	// 0x2400～0x27BF	0x03C0	ネームテーブル1
	if addrTmp >= 0x2400 && addrTmp <= 0x03C0 {
		data = b.nameTable1[addrTmp-0x2400]
		target = "NameTable1"
		return
	}

	// 0x27C0～0x27FF	0x0040	属性テーブル1
	if addrTmp >= 0x27C0 && addrTmp <= 0x27FF {
		data = b.attributeTable1[addrTmp-0x27C0]
		target = "AttributeTable1"
		return
	}

	// 0x2800～0x2BBF	0x03C0	ネームテーブル2
	if addrTmp >= 0x2800 && addrTmp <= 0x2BBF {
		data = b.nameTable2[addrTmp-0x2800]
		target = "NameTable2"
		return
	}

	// 0x2BC0～0x2BFF	0x0040	属性テーブル2
	if addrTmp >= 0x2BC0 && addrTmp <= 0x0040 {
		data = b.attributeTable2[addrTmp-0x2BC0]
		target = "AttributeTable2"
		return
	}

	// 0x2C00～0x2FBF	0x03C0	ネームテーブル3
	if addrTmp >= 0x2C00 && addrTmp <= 0x2FBF {
		data = b.nameTable3[addrTmp-0x2C00]
		target = "NameTable3"
		return
	}

	// 0x2FC0～0x2FFF	0x0040	属性テーブル3
	if addrTmp >= 0x2FC0 && addrTmp <= 0x2FFF {
		data = b.attributeTable3[addrTmp-0x2FC0]
		target = "AttributeTable3"
		return
	}

	// 0x3F00～0x3F0F	0x0010	バックグラウンドパレット
	if addrTmp >= 0x3F00 && addrTmp <= 0x3F0F {
		pIdx := (addrTmp - 0x3F00) / 4
		bitIdx := (addrTmp - 0x3F00) % 4
		data = b.backgroundPalette[pIdx][bitIdx]
		target = "BackgroundPalette"
		return
	}
	// 0x3F10～0x3F1F	0x0010	スプライトパレット
	if addrTmp >= 0x3F10 && addrTmp <= 0x3F1F {
		pIdx := (addrTmp - 0x3F10) / 4
		bitIdx := (addrTmp - 0x3F10) % 4
		data = b.spritePalette[pIdx][bitIdx]
		target = "SpritePalette"
		return
	}

	err = fmt.Errorf("failed to read by PPU; addr: %v", addr)
	return
}

// writeByPPU ...
func (b *Bus) writeByPPU(addr domain.Address, data byte) (err error) {
	var target string
	log.Trace("Bus.writeByPPU[addr=%#v] (<=%#v)...", addr, data)
	defer func() {
		if err != nil {
			log.Warn("Bus.writeByPPU[addr=%#v] => %#v", addr, err)
		} else {
			log.Trace("Bus.writeByPPU[addr=%#v][%v] <= %#v", addr, target, data)
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
	addrTmp := addr
	// 0x3000～0x3EFF	-	0x2000-0x2EFFのミラー
	if addr >= 0x3000 && addr <= 0x3EFF {
		addrTmp = addr - (0x3000 - 0x2000)
	}
	// 0x3F20～0x3FFF	-	0x3F00-0x3F1Fのミラー
	if addr >= 0x3F20 && addr <= 0x3FFF {
		addrTmp = addr - (0x3F20 - 0x3F00)
	}

	// 0x0000～0x0FFF	0x1000	パターンテーブル0
	if addrTmp >= 0x0000 && addrTmp <= 0x0FFF {
		err = fmt.Errorf("failed write, PatternTable0(CHR-ROM) is read only; addr: %#v", addr)
		target = "PatternTable0"
		return
	}

	// 0x1000～0x1FFF	0x1000	パターンテーブル1
	if addrTmp >= 0x1000 && addrTmp <= 0x1FFF {
		err = fmt.Errorf("failed write, PatternTable1(CHR-ROM) is read only; addr: %#v", addr)
		target = "PatternTable1"
		return
	}

	// 0x2000～0x23BF	0x03C0	ネームテーブル0
	if addrTmp >= 0x2000 && addrTmp <= 0x23BF {
		b.nameTable0[addrTmp-0x2000] = data
		target = "NameTable0"
		return
	}

	// 0x23C0～0x23FF	0x0040	属性テーブル0
	if addrTmp >= 0x23C0 && addrTmp <= 0x23FF {
		b.attributeTable0[addrTmp-0x23C0] = data
		target = "AttributeTable0"
		return
	}

	// 0x2400～0x27BF	0x03C0	ネームテーブル1
	if addrTmp >= 0x2400 && addrTmp <= 0x03C0 {
		b.nameTable1[addrTmp-0x2400] = data
		target = "NameTable1"
		return
	}

	// 0x27C0～0x27FF	0x0040	属性テーブル1
	if addrTmp >= 0x27C0 && addrTmp <= 0x27FF {
		b.attributeTable1[addrTmp-0x27C0] = data
		target = "AttributeTable1"
		return
	}

	// 0x2800～0x2BBF	0x03C0	ネームテーブル2
	if addrTmp >= 0x2800 && addrTmp <= 0x2BBF {
		b.nameTable2[addrTmp-0x2800] = data
		target = "NameTable2"
		return
	}

	// 0x2BC0～0x2BFF	0x0040	属性テーブル2
	if addrTmp >= 0x2BC0 && addrTmp <= 0x0040 {
		b.attributeTable2[addrTmp-0x2BC0] = data
		target = "AttributeTable2"
		return
	}

	// 0x2C00～0x2FBF	0x03C0	ネームテーブル3
	if addrTmp >= 0x2C00 && addrTmp <= 0x2FBF {
		b.nameTable3[addrTmp-0x2C00] = data
		target = "NameTable3"
		return
	}

	// 0x2FC0～0x2FFF	0x0040	属性テーブル3
	if addrTmp >= 0x2FC0 && addrTmp <= 0x2FFF {
		b.attributeTable3[addrTmp-0x2FC0] = data
		target = "AttributeTable3"
		return
	}

	// 0x3F00～0x3F0F	0x0010	バックグラウンドパレット
	if addrTmp >= 0x3F00 && addrTmp <= 0x3F0F {
		pIdx := (addrTmp - 0x3F00) / 4
		bitIdx := (addrTmp - 0x3F00) % 4
		b.backgroundPalette[pIdx][bitIdx] = data
		target = "BackgroundPalette"
		return
	}
	// 0x3F10～0x3F1F	0x0010	スプライトパレット
	if addrTmp >= 0x3F10 && addrTmp <= 0x3F1F {
		pIdx := (addrTmp - 0x3F10) / 4
		bitIdx := (addrTmp - 0x3F10) % 4
		b.spritePalette[pIdx][bitIdx] = data
		target = "SpritePalette"
		return
	}

	err = fmt.Errorf("failed to read by PPU; addr: %v", addr)
	return
}

// GetSpriteNo ...
func (b *Bus) GetSpriteNo(nameTblIdx uint8, p domain.NameTablePoint) (no uint8, err error) {
	log.Trace("Bus.GetSpriteNo[%#v] ...", p)
	defer func() {
		if err != nil {
			log.Warn("Bus.GetSpriteNo[%#v] => %#v", p, err)
		} else {
			log.Trace("Bus.GetSpriteNo[%#v] => %#v", p, no)
		}
	}()

	if err = p.Validate(); err != nil {
		return
	}

	switch nameTblIdx {
	case 0:
		no = b.nameTable0[p.ToIndex()]
	case 1:
		no = b.nameTable1[p.ToIndex()]
	case 2:
		no = b.nameTable2[p.ToIndex()]
	case 3:
		no = b.nameTable3[p.ToIndex()]
	}

	return
}

// GetSprite ...
func (b *Bus) GetSprite(no uint8) *domain.Sprite {
	return b.charactorROM.GetSprite(no)
}

// GetPaletteNo ...
func (b *Bus) GetPaletteNo(p domain.NameTablePoint) (no uint8, err error) {
	log.Trace("Bus.GetPaletteNo[%#v] ...", p)
	defer func() {
		if err != nil {
			log.Warn("Bus.GetPaletteNo[%#v] => %#v", p, err)
		} else {
			log.Trace("Bus.GetPaletteNo[%#v] => %#v", p, no)
		}
	}()

	err = p.Validate()
	if err != nil {
		return
	}

	attribute := b.attributeTable0[p.ToAttributeTableIndex()]
	attributeIndex := p.ToAttributeIndex()
	no = (attribute & (0x03 << attributeIndex)) >> attributeIndex

	return
}

// GetBackgroundPalette ...
func (b *Bus) GetBackgroundPalette(no uint8) *domain.Palette {
	return &b.backgroundPalette[no]
}

// GetSpritePalette ...
func (b *Bus) GetSpritePalette(no uint8) *domain.Palette {
	return &b.spritePalette[no]
}