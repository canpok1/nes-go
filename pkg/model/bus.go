package model

import (
	"fmt"

	"github.com/canpok1/nes-go/pkg/log"
)

// Bus ...
type Bus struct {
	wram              []byte
	wramMirror        []byte
	ppuRegister       []byte
	ppuRegisterMirror []byte
	io                []byte
	exrom             []byte
	exram             []byte
	programRom        *PRGROM
}

// NewBus ...
func NewBus(p *PRGROM) *Bus {
	return &Bus{
		wram:        make([]byte, 0x0800),
		ppuRegister: make([]byte, 0x0008),
		io:          make([]byte, 0x0020),
		exrom:       make([]byte, 0x1FE0),
		exram:       make([]byte, 0x2000),
		programRom:  p,
	}
}

// readByCPU ...
func (c *Bus) readByCPU(addr Address) (byte, error) {
	var data byte
	var err error
	var target string
	defer func() {
		if err != nil {
			log.Debug("Bus.readByCPU[addr=%#v] => %#v", addr, err)
		} else {
			log.Debug("Bus.readByCPU[addr=%#v][%v] => %#v", addr, target, data)
		}
	}()

	// 0x0000～0x07FF	0x0800	WRAM
	if addr >= 0x0000 && addr <= 0x07FF {
		target = "WRAM"
		data = c.wram[addr]
		return data, nil
	}

	// 0x0800～0x1FFF	-	WRAMのミラー
	if addr >= 0x0800 && addr <= 0x1FFF {
		target = "WRAM Mirror"
		data = c.wram[addr-0x0800]
		return data, nil
	}

	// 0x2000～0x2007	0x0008	PPU レジスタ
	if addr >= 0x2000 && addr <= 0x2007 {
		target = "PPU Register"
		data = c.ppuRegister[addr-0x2000]
		return data, nil
	}

	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	if addr >= 0x2008 && addr <= 0x3FFF {
		target = "PPU Register Mirror"
		data = c.ppuRegister[addr-0x2008]
		return data, nil
	}

	// 0x4000～0x401F	0x0020	APU I/O、PAD
	if addr >= 0x4000 && addr <= 0x401F {
		target = "APU I/O, PAD"
		data = c.io[addr-0x4000]
		return data, nil
	}

	// 0x4020～0x5FFF	0x1FE0	拡張ROM
	if addr >= 0x4000 && addr <= 0x401F {
		target = "EX ROM"
		data = c.exrom[addr-0x4000]
		return data, nil
	}

	// 0x6000～0x7FFF	0x2000	拡張RAM
	if addr >= 0x6000 && addr <= 0x7FFF {
		target = "EX RAM"
		data = c.exram[addr-0x6000]
		return data, nil
	}

	// 0x8000～0xBFFF	0x4000	PRG-ROM
	// 0xC000～0xFFFF	0x4000	PRG-ROM
	if addr >= 0x8000 && addr <= 0xFFFF {
		target = "PRG-ROM"
		r := *c.programRom
		if len(r) <= 0x4000 {
			data = r[addr-0xC000]
		} else {
			data = r[addr-0x8000]
		}
		return data, nil
	}

	return 0, fmt.Errorf("failed read, addr out of range; addr: %#v", addr)
}

// writeByCPU ...
func (c *Bus) writeByCPU(addr Address, data byte) error {
	var err error
	var target string
	defer func() {
		if err != nil {
			log.Debug("Bus.writeByCPU[addr=%#v] => %#v", addr, err)
		} else {
			log.Debug("Bus.writeByCPU[addr=%#v][%v] <= %#v", addr, target, data)
		}
	}()
	// 0x0000～0x07FF	0x0800	WRAM
	if addr >= 0x0000 && addr <= 0x07FF {
		target = "WRAM"
		c.wram[addr] = data
		return nil
	}

	// 0x0800～0x1FFF	-	WRAMのミラー
	if addr >= 0x0800 && addr <= 0x1FFF {
		target = "WRAM Mirror"
		c.wram[addr-0x0800] = data
		return nil
	}

	// 0x2000～0x2007	0x0008	PPU レジスタ
	if addr >= 0x2000 && addr <= 0x2007 {
		target = "PPU Register"
		c.ppuRegister[addr-0x2000] = data
		return nil
	}

	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	if addr >= 0x2008 && addr <= 0x3FFF {
		target = "PPU Register Mirror"
		c.ppuRegister[addr-0x2008] = data
		return nil
	}

	// 0x4000～0x401F	0x0020	APU I/O、PAD
	if addr >= 0x4000 && addr <= 0x401F {
		target = "APU I/O, PAD"
		c.io[addr-0x4000] = data
		return nil
	}

	// 0x4020～0x5FFF	0x1FE0	拡張ROM
	if addr >= 0x4000 && addr <= 0x401F {
		return fmt.Errorf("failed write, cannot write EX ROM; addr: %#v", addr)
	}

	// 0x6000～0x7FFF	0x2000	拡張RAM
	if addr >= 0x6000 && addr <= 0x7FFF {
		target = "EX RAM"
		c.exram[addr-0x6000] = data
		return nil
	}

	// 0x8000～0xBFFF	0x4000	PRG-ROM
	// 0xC000～0xFFFF	0x4000	PRG-ROM
	if addr >= 0x8000 && addr <= 0xFFFF {
		return fmt.Errorf("failed write, cannot write PRG-ROM; addr: %#v", addr)
	}

	return fmt.Errorf("failed write, addr out of range; addr: %#v", addr)
}
