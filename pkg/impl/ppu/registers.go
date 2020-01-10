package ppu

import "fmt"

// PPURegisters ...
type PPURegisters struct {
	PPUCtrl   byte // 0x2000	PPUCTRL	W	コントロールレジスタ1	割り込みなどPPUの設定
	PPUMask   byte // 0x2001	PPUMASK	W	コントロールレジスタ2	背景イネーブルなどのPPU設定
	PPUStatus byte // 0x2002	PPUSTATUS	R	PPUステータス	PPUのステータス
	OAMAddr   byte // 0x2003	OAMADDR	W	スプライトメモリデータ	書き込むスプライト領域のアドレス
	OAMData   byte // 0x2004	OAMDATA	RW	デシマルモード	スプライト領域のデータ
	PPUScroll byte // 0x2005	PPUSCROLL	W	背景スクロールオフセット	背景スクロール値
	PPUAddr   byte // 0x2006	PPUADDR	W	PPUメモリアドレス	書き込むPPUメモリ領域のアドレス
}

// NewPPURegisters ...
func NewPPURegisters() *PPURegisters {
	return &PPURegisters{
		PPUCtrl:   0,
		PPUMask:   0,
		PPUStatus: 0,
		OAMAddr:   0,
		OAMData:   0,
		PPUScroll: 0,
		PPUAddr:   0,
	}
}

// String ...
func (r *PPURegisters) String() string {
	return fmt.Sprintf(
		"{PPUCTRL:%#v, PPUMASK:%#v, PPUSTATUS:%#v, OAMADDR:%#v, OAMDATA:%v, PPUSCROLL:%#v, PPUADDR:%#v}",
		r.PPUCtrl,
		r.PPUMask,
		r.PPUStatus,
		r.OAMAddr,
		r.OAMData,
		r.PPUScroll,
		r.PPUAddr,
	)
}