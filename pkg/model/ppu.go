package model

import (
	"fmt"

	"github.com/canpok1/nes-go/pkg/log"
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
	ppudata   byte // 0x2007	PPUDATA	RW	PPUメモリデータ	PPUメモリ領域のデータ
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
		ppudata:   0,
	}
}

// String ...
func (r *PPURegisters) String() string {
	return fmt.Sprintf(
		"{PPUCTRL:%#v, PPUMASK:%#v, PPUSTATUS:%#v, OAMADDR:%#v, OAMDATA:%v, PPUSCROLL:%#v, PPUADDR:%#v, PPUDATA:%#v}",
		r.ppuctrl,
		r.ppumask,
		r.ppustatus,
		r.oamaddr,
		r.oamdata,
		r.ppuscroll,
		r.ppuaddr,
		r.ppudata,
	)
}

// Read ...
func (r *PPURegisters) Read(addr Address) (byte, error) {
	var data byte
	var err error
	defer func() {
		if err != nil {
			log.Warn("PPURegisters.Read[%#v] => %#v", addr, err)
		} else {
			log.Warn("PPURegisters.Read[%#v] => %#v", addr, data)
		}
	}()

	switch addr {
	case 0:
		data = r.ppuctrl
	case 1:
		data = r.ppumask
	case 2:
		data = r.ppustatus
	case 3:
		data = r.oamaddr
	case 4:
		data = r.oamdata
	case 5:
		data = r.ppuscroll
	case 6:
		data = r.ppuaddr
	case 7:
		data = r.ppudata
	default:
		err = fmt.Errorf("failed to read PPURegisters, address is out of range; addr: %#v", addr)
	}

	return data, err
}

// Write ...
func (r *PPURegisters) Write(addr Address, data byte) error {
	var err error
	defer func() {
		if err != nil {
			log.Warn("PPURegisters.Write[%#v] => %#v", addr, err)
		} else {
			log.Debug("PPURegisters.Write[%#v] <= %#v", addr, data)
		}
	}()

	switch addr {
	case 0:
		r.ppuctrl = data
	case 1:
		r.ppumask = data
	case 2:
		r.ppustatus = data
	case 3:
		r.oamaddr = data
	case 4:
		r.oamdata = data
	case 5:
		r.ppuscroll = data
	case 6:
		r.ppuaddr = data
	case 7:
		r.ppudata = data
	default:
		err = fmt.Errorf("failed to write PPURegisters, address is out of range; addr: %#v", addr)
	}

	return err
}

// PPU ...
type PPU struct {
	registers *PPURegisters
	bus       *Bus
}

// NewPPU ...
func NewPPU() *PPU {
	return &PPU{
		registers: NewPPURegisters(),
	}
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

// Run ...
func (p *PPU) Run() error {
	// TODO 実装
	return nil
}
