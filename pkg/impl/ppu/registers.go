package ppu

import "fmt"

// PPURegisters ...
type PPURegisters struct {
	PPUCtrl   PPUCtrl   // 0x2000	PPUCTRL	W	コントロールレジスタ1	割り込みなどPPUの設定
	PPUMask   byte      // 0x2001	PPUMASK	W	コントロールレジスタ2	背景イネーブルなどのPPU設定
	PPUStatus PPUStatus // 0x2002	PPUSTATUS	R	PPUステータス	PPUのステータス
	OAMAddr   byte      // 0x2003	OAMADDR	W	スプライトメモリデータ	書き込むスプライト領域のアドレス
	OAMData   byte      // 0x2004	OAMDATA	RW	デシマルモード	スプライト領域のデータ
	PPUScroll byte      // 0x2005	PPUSCROLL	W	背景スクロールオフセット	背景スクロール値
	PPUAddr   byte      // 0x2006	PPUADDR	W	PPUメモリアドレス	書き込むPPUメモリ領域のアドレス
}

// NewPPURegisters ...
func NewPPURegisters() *PPURegisters {
	return &PPURegisters{
		PPUCtrl: PPUCtrl{
			NMIEnable:        false,
			SpriteTileSelect: false,
			NameTableIndex:   0,
		},
		PPUMask: 0,
		PPUStatus: PPUStatus{
			VBlankHasStarted: false,
		},
		OAMAddr:   0,
		OAMData:   0,
		PPUScroll: 0,
		PPUAddr:   0,
	}
}

// String ...
func (r PPURegisters) String() string {
	return fmt.Sprintf(
		"{PPUCTRL:%#v, PPUMASK:%#v, PPUSTATUS:%#v, OAMADDR:%#v, OAMDATA:%v, PPUSCROLL:%#v, PPUADDR:%#v}",
		r.PPUCtrl.String(),
		r.PPUMask,
		r.PPUStatus.String(),
		r.OAMAddr,
		r.OAMData,
		r.PPUScroll,
		r.PPUAddr,
	)
}

// PPUCtrl ...
type PPUCtrl struct {
	NMIEnable        bool
	SpriteTileSelect bool
	NameTableIndex   uint8
}

// String ...
func (p PPUCtrl) String() string {
	return fmt.Sprintf(
		"{NMIEnable: %v, SpriteTileSelect: %v}",
		p.NMIEnable,
		p.SpriteTileSelect,
	)
}

// UpdateAll ...
func (p PPUCtrl) UpdateAll(b byte) {
	p.SpriteTileSelect = (b % 0x04) == 0x04
	p.NMIEnable = (b % 0x80) == 0x80
}

// PPUStatus ...
type PPUStatus struct {
	VBlankHasStarted bool
}

// ToByte ...
func (p PPUStatus) ToByte() byte {
	var b byte
	if p.VBlankHasStarted {
		b = b + 0x80
	}
	return b
}

// String ...
func (p PPUStatus) String() string {
	return fmt.Sprintf(
		"{VBlank: %v}",
		p.VBlankHasStarted,
	)
}
