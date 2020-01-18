package impl

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// Bus ...
type Bus struct {
	wram       []byte
	wramMirror []byte
	io         []byte
	exrom      []byte
	exram      []byte
	programROM *domain.PRGROM

	ppu  domain.PPU
	cpu  domain.CPU
	pad1 domain.Pad
	pad2 domain.Pad

	charactorROM *domain.CHRROM

	vram *domain.VRAM

	pad1ReadCount int
	pad2ReadCount int

	pad1WriteBuf byte
	pad2WriteBuf byte

	setupped bool
}

// NewBus ...
func NewBus() domain.Bus {
	return &Bus{
		wram:  make([]byte, 0x0800),
		io:    make([]byte, 0x0020),
		exrom: make([]byte, 0x1FE0),
		exram: make([]byte, 0x2000),

		pad1ReadCount: 0,
		pad2ReadCount: 0,

		setupped: false,
	}
}

// Setup ...
func (b *Bus) Setup(rom *domain.ROM, ppu domain.PPU, cpu domain.CPU, vram *domain.VRAM, pad1 domain.Pad, pad2 domain.Pad) {
	b.programROM = rom.Prgrom
	b.charactorROM = rom.Chrrom
	b.ppu = ppu
	b.cpu = cpu
	b.vram = vram
	b.pad1 = pad1
	b.pad2 = pad2

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
			log.Warn("Bus.readByCPU[addr=%#v][%v] => %#v", addr, target, err)
		} else {
			log.Trace("Bus.readByCPU[addr=%#v][%v] => %#v", addr, target, data)
		}
	}()

	if b == nil {
		err = xerrors.Errorf("failed to readByCPU, bus is nil")
		return data, err
	}
	if !b.setupped {
		err = xerrors.Errorf("failed to readByCPU, bus setup is not completed")
		return data, err
	}

	// 0x0000～0x07FF	0x0800	WRAM
	if addr >= 0x0000 && addr <= 0x07FF {
		target = "WRAM"
		data = b.wram[addr]
		return data, err
	}

	// 0x0800～0x0FFF	-	WRAMのミラー1
	if addr >= 0x0800 && addr <= 0x0FFF {
		target = "WRAM Mirror 1"
		data = b.wram[addr-0x0800]
		return data, err
	}

	// 0x1000～0x17FF	-	WRAMのミラー2
	if addr >= 0x1000 && addr <= 0x17FF {
		target = "WRAM Mirror 2"
		data = b.wram[addr-0x1000]
		return data, err
	}

	// 0x1800～0x1FFF	-	WRAMのミラー3
	if addr >= 0x1800 && addr <= 0x1FFF {
		target = "WRAM Mirror 3"
		data = b.wram[addr-0x1800]
		return data, err
	}

	// 0x2000～0x2007	0x0008	PPU レジスタ
	if addr >= 0x2000 && addr <= 0x2007 {
		target = "PPU Register"
		if data, err = b.ppu.ReadRegisters(addr); err != nil {
			err = xerrors.Errorf(": %w", err)
		}
		return data, err
	}

	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	if addr >= 0x2008 && addr <= 0x3FFF {
		target = "PPU Register Mirror"
		if data, err = b.ppu.ReadRegisters(addr); err != nil {
			err = xerrors.Errorf(": %w", err)
		}
		return data, err
	}

	// 0x4016 PAD1
	if addr == 0x4016 {
		var pressed bool
		switch b.pad1ReadCount {
		case 0:
			pressed = b.pad1.IsPressed(domain.ButtonTypeA)
		case 1:
			pressed = b.pad1.IsPressed(domain.ButtonTypeB)
		case 2:
			pressed = b.pad1.IsPressed(domain.ButtonTypeSelect)
		case 3:
			pressed = b.pad1.IsPressed(domain.ButtonTypeStart)
		case 4:
			pressed = b.pad1.IsPressed(domain.ButtonTypeUp)
		case 5:
			pressed = b.pad1.IsPressed(domain.ButtonTypeDown)
		case 6:
			pressed = b.pad1.IsPressed(domain.ButtonTypeLeft)
		case 7:
			pressed = b.pad1.IsPressed(domain.ButtonTypeRight)
		}

		if pressed {
			data = 1
		} else {
			data = 0
		}

		if b.pad1ReadCount < 8 {
			b.pad1ReadCount = b.pad1ReadCount + 1
		}
		return data, nil
	}
	// 0x4017 PAD2
	if addr == 0x4017 {
		var pressed bool
		switch b.pad2ReadCount {
		case 0:
			pressed = b.pad2.IsPressed(domain.ButtonTypeA)
		case 1:
			pressed = b.pad2.IsPressed(domain.ButtonTypeB)
		case 2:
			pressed = b.pad2.IsPressed(domain.ButtonTypeSelect)
		case 3:
			pressed = b.pad2.IsPressed(domain.ButtonTypeStart)
		case 4:
			pressed = b.pad2.IsPressed(domain.ButtonTypeUp)
		case 5:
			pressed = b.pad2.IsPressed(domain.ButtonTypeDown)
		case 6:
			pressed = b.pad2.IsPressed(domain.ButtonTypeLeft)
		case 7:
			pressed = b.pad2.IsPressed(domain.ButtonTypeRight)
		}

		if pressed {
			data = 1
		} else {
			data = 0
		}

		if b.pad2ReadCount < 8 {
			b.pad2ReadCount = b.pad2ReadCount + 1
		}
		return data, nil
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

	return 0, xerrors.Errorf("failed read, addr out of range; addr: %#v", addr)
}

// WriteByCPU ...
func (b *Bus) WriteByCPU(addr domain.Address, data byte) error {
	var err error
	var target string
	log.Trace("Bus.writeByCPU[addr=%#v] (<=%#v) ...", addr, data)
	defer func() {
		if err != nil {
			log.Warn("Bus.writeByCPU[addr=%#v][%v] => %#v", addr, target, err)
		} else {
			log.Trace("Bus.writeByCPU[addr=%#v][%v] <= %#v", addr, target, data)
		}
	}()

	if b == nil {
		err = xerrors.Errorf("failed to writeByCPU, bus is nil")
		return err
	}
	if !b.setupped {
		err = xerrors.Errorf("failed to writeByCPU, bus setup is not completed")
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
		err = b.ppu.WriteRegisters(addr, data)
		return err
	}

	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	if addr >= 0x2008 && addr <= 0x3FFF {
		target = "PPU Register Mirror"
		err = b.ppu.WriteRegisters(addr, data)
		return err
	}

	// 0x4014 OAMDMA
	if addr == 0x4014 {
		target = "OAMDMA"
		err = b.ppu.WriteRegisters(addr, data)
		return err
	}

	// 0x4016 PAD1
	if addr == 0x4016 {
		if b.pad1WriteBuf == 0x01 && data == 0x00 {
			b.pad1ReadCount = 0
		}
		b.pad1WriteBuf = data
		return nil
	}

	// 0x4017 PAD2
	if addr == 0x4017 {
		if b.pad2WriteBuf == 0x01 && data == 0x00 {
			b.pad2ReadCount = 0
		}
		b.pad2WriteBuf = data
		return nil
	}

	// 0x4000～0x401F	0x0020	APU I/O、PAD
	if addr >= 0x4000 && addr <= 0x401F {
		target = "APU I/O, PAD"
		b.io[addr-0x4000] = data
		return nil
	}

	// 0x4020～0x5FFF	0x1FE0	拡張ROM
	if addr >= 0x4000 && addr <= 0x401F {
		return xerrors.Errorf("failed write, cannot write EX ROM; addr: %#v", addr)
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
		return xerrors.Errorf("failed write, cannot write PRG-ROM; addr: %#v", addr)
	}

	return xerrors.Errorf("failed write, addr out of range; addr: %#v", addr)
}

// ReadByPPU ...
func (b *Bus) ReadByPPU(addr domain.Address) (data byte, err error) {
	var target string
	log.Trace("Bus.readByPPU[addr=%#v] ...", addr)
	defer func() {
		if err != nil {
			log.Warn("Bus.readByPPU[addr=%#v][%v] => %#v", addr, target, err)
		} else {
			log.Trace("Bus.readByPPU[addr=%#v][%v] => %#v", addr, target, data)
		}
	}()

	if b == nil {
		err = xerrors.Errorf("failed to readByPPU, bus is nil")
		return data, err
	}
	if !b.setupped {
		err = xerrors.Errorf("failed to readByPPU, bus setup is not completed")
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
		data = b.vram.NameTable0[addrTmp-0x2000]
		target = "NameTable0"
		return
	}

	// 0x23C0～0x23FF	0x0040	属性テーブル0
	if addrTmp >= 0x23C0 && addrTmp <= 0x23FF {
		data = b.vram.AttributeTable0[addrTmp-0x23C0]
		target = "AttributeTable0"
		return
	}

	// 0x2400～0x27BF	0x03C0	ネームテーブル1
	if addrTmp >= 0x2400 && addrTmp <= 0x03C0 {
		data = b.vram.NameTable1[addrTmp-0x2400]
		target = "NameTable1"
		return
	}

	// 0x27C0～0x27FF	0x0040	属性テーブル1
	if addrTmp >= 0x27C0 && addrTmp <= 0x27FF {
		data = b.vram.AttributeTable1[addrTmp-0x27C0]
		target = "AttributeTable1"
		return
	}

	// 0x2800～0x2BBF	0x03C0	ネームテーブル2
	if addrTmp >= 0x2800 && addrTmp <= 0x2BBF {
		data = b.vram.NameTable2[addrTmp-0x2800]
		target = "NameTable2"
		return
	}

	// 0x2BC0～0x2BFF	0x0040	属性テーブル2
	if addrTmp >= 0x2BC0 && addrTmp <= 0x0040 {
		data = b.vram.AttributeTable2[addrTmp-0x2BC0]
		target = "AttributeTable2"
		return
	}

	// 0x2C00～0x2FBF	0x03C0	ネームテーブル3
	if addrTmp >= 0x2C00 && addrTmp <= 0x2FBF {
		data = b.vram.NameTable3[addrTmp-0x2C00]
		target = "NameTable3"
		return
	}

	// 0x2FC0～0x2FFF	0x0040	属性テーブル3
	if addrTmp >= 0x2FC0 && addrTmp <= 0x2FFF {
		data = b.vram.AttributeTable3[addrTmp-0x2FC0]
		target = "AttributeTable3"
		return
	}

	// 0x3F00～0x3F0F	0x0010	バックグラウンドパレット
	if addrTmp >= 0x3F00 && addrTmp <= 0x3F0F {
		pIdx := (addrTmp - 0x3F00) / 4
		bitIdx := (addrTmp - 0x3F00) % 4
		data = b.vram.BackgroundPalette[pIdx][bitIdx]
		target = "BackgroundPalette"
		return
	}
	// 0x3F10～0x3F1F	0x0010	スプライトパレット
	if addrTmp >= 0x3F10 && addrTmp <= 0x3F1F {
		pIdx := (addrTmp - 0x3F10) / 4
		bitIdx := (addrTmp - 0x3F10) % 4
		data = b.vram.SpritePalette[pIdx][bitIdx]
		target = "SpritePalette"
		return
	}

	err = xerrors.Errorf("failed to read by PPU; addr: %v", addr)
	return
}

// WriteByPPU ...
func (b *Bus) WriteByPPU(addr domain.Address, data byte) (err error) {
	var target string
	log.Trace("Bus.writeByPPU[addr=%#v] (<=%#v)...", addr, data)
	defer func() {
		if err != nil {
			log.Warn("Bus.writeByPPU[addr=%#v][%v] => %#v", addr, target, err)
		} else {
			log.Trace("Bus.writeByPPU[addr=%#v][%v] <= %#v", addr, target, data)
		}
	}()

	if b == nil {
		err = xerrors.Errorf("failed to writeByPPU, bus is nil")
		return err
	}
	if !b.setupped {
		err = xerrors.Errorf("failed to writeByPPU, bus setup is not completed")
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
		err = xerrors.Errorf("failed write, PatternTable0(CHR-ROM) is read only; addr: %#v", addr)
		target = "PatternTable0"
		return
	}

	// 0x1000～0x1FFF	0x1000	パターンテーブル1
	if addrTmp >= 0x1000 && addrTmp <= 0x1FFF {
		err = xerrors.Errorf("failed write, PatternTable1(CHR-ROM) is read only; addr: %#v", addr)
		target = "PatternTable1"
		return
	}

	// 0x2000～0x23BF	0x03C0	ネームテーブル0
	if addrTmp >= 0x2000 && addrTmp <= 0x23BF {
		b.vram.NameTable0[addrTmp-0x2000] = data
		target = "NameTable0"
		return
	}

	// 0x23C0～0x23FF	0x0040	属性テーブル0
	if addrTmp >= 0x23C0 && addrTmp <= 0x23FF {
		b.vram.AttributeTable0[addrTmp-0x23C0] = data
		target = "AttributeTable0"
		return
	}

	// 0x2400～0x27BF	0x03C0	ネームテーブル1
	if addrTmp >= 0x2400 && addrTmp <= 0x27BF {
		b.vram.NameTable1[addrTmp-0x2400] = data
		target = "NameTable1"
		return
	}

	// 0x27C0～0x27FF	0x0040	属性テーブル1
	if addrTmp >= 0x27C0 && addrTmp <= 0x27FF {
		b.vram.AttributeTable1[addrTmp-0x27C0] = data
		target = "AttributeTable1"
		return
	}

	// 0x2800～0x2BBF	0x03C0	ネームテーブル2
	if addrTmp >= 0x2800 && addrTmp <= 0x2BBF {
		b.vram.NameTable2[addrTmp-0x2800] = data
		target = "NameTable2"
		return
	}

	// 0x2BC0～0x2BFF	0x0040	属性テーブル2
	if addrTmp >= 0x2BC0 && addrTmp <= 0x2BFF {
		b.vram.AttributeTable2[addrTmp-0x2BC0] = data
		target = "AttributeTable2"
		return
	}

	// 0x2C00～0x2FBF	0x03C0	ネームテーブル3
	if addrTmp >= 0x2C00 && addrTmp <= 0x2FBF {
		b.vram.NameTable3[addrTmp-0x2C00] = data
		target = "NameTable3"
		return
	}

	// 0x2FC0～0x2FFF	0x0040	属性テーブル3
	if addrTmp >= 0x2FC0 && addrTmp <= 0x2FFF {
		b.vram.AttributeTable3[addrTmp-0x2FC0] = data
		target = "AttributeTable3"
		return
	}

	// 0x3F00～0x3F0F	0x0010	バックグラウンドパレット
	if addrTmp >= 0x3F00 && addrTmp <= 0x3F0F {
		pIdx := (addrTmp - 0x3F00) / 4
		bitIdx := (addrTmp - 0x3F00) % 4
		b.vram.BackgroundPalette[pIdx][bitIdx] = data

		if bitIdx == 0 {
			b.vram.SpritePalette[pIdx][bitIdx] = data
		}

		target = "BackgroundPalette"
		return
	}
	// 0x3F10～0x3F1F	0x0010	スプライトパレット
	if addrTmp >= 0x3F10 && addrTmp <= 0x3F1F {
		pIdx := (addrTmp - 0x3F10) / 4
		bitIdx := (addrTmp - 0x3F10) % 4
		b.vram.SpritePalette[pIdx][bitIdx] = data

		if bitIdx == 0 {
			b.vram.BackgroundPalette[pIdx][bitIdx] = data
		}

		target = "SpritePalette"
		return
	}

	err = xerrors.Errorf("failed to read by PPU; addr: %v", addr)
	return
}

// GetTileNo ...
func (b *Bus) GetTileNo(nameTblIdx uint8, p domain.NameTablePoint) (no uint8, err error) {
	log.Trace("Bus.GetTileNo[%#v] ...", p)
	defer func() {
		if err != nil {
			log.Warn("Bus.GetTileNo[%#v] => %#v", p, err)
		} else {
			log.Trace("Bus.GetTileNo[%#v] => %#v", p, no)
		}
	}()

	if err = p.Validate(); err != nil {
		return
	}

	switch nameTblIdx {
	case 0:
		no = b.vram.NameTable0[p.ToIndex()]
	case 1:
		no = b.vram.NameTable1[p.ToIndex()]
	case 2:
		no = b.vram.NameTable2[p.ToIndex()]
	case 3:
		no = b.vram.NameTable3[p.ToIndex()]
	}

	return
}

// GetTilePattern ...
func (b *Bus) GetTilePattern(patternTblIdx, no uint8) *domain.TilePattern {
	return b.charactorROM.GetTilePattern(patternTblIdx, no)
}

// GetAttribute ...
func (b *Bus) GetAttribute(tableIndex uint8, p domain.NameTablePoint) (attribute byte, err error) {
	log.Trace("Bus.GetAttribute[%v][%#v] ...", tableIndex, p)
	defer func() {
		if err != nil {
			log.Warn("Bus.GetAttribute[%v][%#v] => %#v", tableIndex, p, err)
		} else {
			log.Trace("Bus.GetAttribute[%v][%#v] => %#v", tableIndex, p, attribute)
		}
	}()

	err = p.Validate()
	if err != nil {
		return
	}

	var tbl []byte
	switch tableIndex {
	case 0:
		tbl = b.vram.AttributeTable0
	case 1:
		tbl = b.vram.AttributeTable1
	case 2:
		tbl = b.vram.AttributeTable2
	case 3:
		tbl = b.vram.AttributeTable3
	}

	attribute = tbl[p.ToAttributeTableIndex()]
	return
}

// GetPaletteNo ...
func (b *Bus) GetPaletteNo(p domain.NameTablePoint, attribute byte) (no uint8, err error) {
	log.Trace("Bus.GetPaletteNo[%#v][%#v] ...", p, attribute)
	defer func() {
		if err != nil {
			log.Warn("Bus.GetPaletteNo[%#v][%#v] => %#v", p, attribute, err)
		} else {
			log.Trace("Bus.GetPaletteNo[%#v][%#v] => %#v", p, attribute, no)
		}
	}()

	attributeIndex := p.ToAttributeIndex()
	no = (attribute & (0x03 << attributeIndex)) >> attributeIndex

	return
}

// GetPalette ...
func (b *Bus) GetPalette(no uint8) *domain.Palette {
	index := no & 0x03
	if (no & 0x0C) == 0 {
		return &b.vram.BackgroundPalette[index]
	}
	return &b.vram.SpritePalette[index]
}

// SendNMI ...
func (b *Bus) SendNMI(active bool) {
	b.cpu.ReceiveNMI(active)
}
