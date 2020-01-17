package impl

import "fmt"

// PPURegisters ...
type PPURegisters struct {
	PPUCtrl   *PPUCtrl   // 0x2000	PPUCTRL	W	コントロールレジスタ1	割り込みなどPPUの設定
	PPUMask   *PPUMask   // 0x2001	PPUMASK	W	コントロールレジスタ2	背景イネーブルなどのPPU設定
	PPUStatus *PPUStatus // 0x2002	PPUSTATUS	R	PPUステータス	PPUのステータス
	OAMAddr   byte       // 0x2003	OAMADDR	W	スプライトメモリデータ	書き込むスプライト領域のアドレス
	OAMData   byte       // 0x2004	OAMDATA	RW	デシマルモード	スプライト領域のデータ
	PPUScroll byte       // 0x2005	PPUSCROLL	W	背景スクロールオフセット	背景スクロール値
	PPUAddr   byte       // 0x2006	PPUADDR	W	PPUメモリアドレス	書き込むPPUメモリ領域のアドレス
	OAMDMA    byte       // 0x4014  OAMDMA W
}

// NewPPURegisters ...
func NewPPURegisters() *PPURegisters {
	return &PPURegisters{
		PPUCtrl: &PPUCtrl{
			NMIEnable:                   false,
			SpriteTileSelect:            false,
			BackgroundPatternTableIndex: 0,
			SpritePatternTableIndex:     0,
			VRAMAddressIncrementMode:    0,
			NameTableIndex:              0,
		},
		PPUMask: &PPUMask{
			EmphasizeB:            false,
			EmphasizeG:            false,
			EmphasizeR:            false,
			EnableSprite:          false,
			EnableBackground:      false,
			DisableSpriteMask:     false,
			DisableBackgroundMask: false,
			DisplayType:           0,
		},
		PPUStatus: &PPUStatus{
			VBlankHasStarted: false,
		},
		OAMAddr:   0,
		OAMData:   0,
		PPUScroll: 0,
		PPUAddr:   0,
		OAMDMA:    0,
	}
}

// String ...
func (r PPURegisters) String() string {
	return fmt.Sprintf(
		"{PPUCTRL:%#v, PPUMASK:%#v, PPUSTATUS:%#v, OAMADDR:%#v, OAMDATA:%v, PPUSCROLL:%#v, PPUADDR:%#v, OAMDMA:%#v}",
		r.PPUCtrl,
		r.PPUMask,
		r.PPUStatus,
		r.OAMAddr,
		r.OAMData,
		r.PPUScroll,
		r.PPUAddr,
		r.OAMDMA,
	)
}

// PPUCtrl ...
type PPUCtrl struct {
	NMIEnable                   bool
	SpriteTileSelect            bool
	BackgroundPatternTableIndex uint8
	SpritePatternTableIndex     uint8
	VRAMAddressIncrementMode    uint8
	NameTableIndex              uint8
}

// UpdateAll ...
func (p *PPUCtrl) UpdateAll(b byte) {
	p.NMIEnable = (b & 0x80) == 0x80
	p.SpriteTileSelect = (b & 0x04) == 0x04
	p.BackgroundPatternTableIndex = (b & 0x10) >> 4
	p.SpritePatternTableIndex = (b & 0x08) >> 3
	p.VRAMAddressIncrementMode = (b & 0x04) >> 2
}

// PPUMask ...
type PPUMask struct {
	EmphasizeB            bool
	EmphasizeG            bool
	EmphasizeR            bool
	EnableSprite          bool
	EnableBackground      bool
	DisableSpriteMask     bool
	DisableBackgroundMask bool
	DisplayType           uint8
}

// UpdateAll ...
func (p *PPUMask) UpdateAll(b byte) {
	p.EmphasizeB = (b & 0x80) == 0x80
	p.EmphasizeG = (b & 0x40) == 0x40
	p.EmphasizeR = (b & 0x20) == 0x20
	p.EnableSprite = (b & 0x10) == 0x10
	p.EnableBackground = (b & 0x08) == 0x08
	p.DisableSpriteMask = (b & 0x04) == 0x04
	p.DisableBackgroundMask = (b & 0x02) == 0x02
	p.DisplayType = b & 0x01
}

// PPUStatus ...
type PPUStatus struct {
	VBlankHasStarted bool
}

// ToByte ...
func (p *PPUStatus) ToByte() byte {
	var b byte
	if p.VBlankHasStarted {
		b = b + 0x80
	}
	return b
}
