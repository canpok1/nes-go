package component

import (
	"fmt"
	"nes-go/pkg/domain"
	"nes-go/pkg/log"
)

// PPURegisters ...
type PPURegisters struct {
	PPUCtrl   *PPUCtrl   // 0x2000	PPUCTRL	W	コントロールレジスタ1	割り込みなどPPUの設定
	PPUMask   *PPUMask   // 0x2001	PPUMASK	W	コントロールレジスタ2	背景イネーブルなどのPPU設定
	PPUStatus *PPUStatus // 0x2002	PPUSTATUS	R	PPUステータス	PPUのステータス
	OAMAddr   byte       // 0x2003	OAMADDR	W	スプライトメモリデータ	書き込むスプライト領域のアドレス
	OAMData   byte       // 0x2004	OAMDATA	RW	デシマルモード	スプライト領域のデータ
	PPUScroll *PPUScroll // 0x2005	PPUSCROLL	W	背景スクロールオフセット	背景スクロール値
	PPUAddr   *PPUAddr   // 0x2006	PPUADDR	W	PPUメモリアドレス	書き込むPPUメモリ領域のアドレス
	OAMDMA    byte       // 0x4014  OAMDMA W

	ppuaddrWriteCount uint8          // PPUADDRへの書き込み回数（0→1→2→1→2→...と遷移）
	ppuaddrBuf        domain.Address // 組み立て中のPPUADDR
	ppuaddrFull       domain.Address // 組み立て済のPPUADDR
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
		OAMAddr: 0,
		OAMData: 0,
		PPUScroll: &PPUScroll{
			buf:     nil,
			vOffset: 0,
			hOffset: 0,
		},
		PPUAddr: &PPUAddr{
			writeCount: 0,
			buf:        0,
			full:       0,
		},
		OAMDMA: 0,
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

// PPUScroll ...
type PPUScroll struct {
	buf     *byte
	vOffset byte
	hOffset byte
}

// Set ...
func (s *PPUScroll) Set(data byte) {
	if s.buf == nil {
		s.buf = &data
		return
	}
	s.vOffset = *s.buf
	s.hOffset = data
	s.buf = nil
}

// GetVOffset ...
func (s *PPUScroll) GetVOffset() byte {
	return s.vOffset
}

// GetHOffset ...
func (s *PPUScroll) GetHOffset() byte {
	return s.hOffset
}

// PPUAddr ...
type PPUAddr struct {
	writeCount uint8          // PPUADDRへの書き込み回数（0→1→2→1→2→...と遷移）
	buf        domain.Address // 組み立て中のPPUADDR
	full       domain.Address // 組み立て済のPPUADDR
}

// Set ...
func (a *PPUAddr) Set(data byte) {
	switch a.writeCount {
	case 0, 2:
		a.buf = domain.Address(data) << 8
		a.writeCount = 1
	case 1:
		a.buf = a.buf + domain.Address(data)
		a.full = a.buf
		a.writeCount = 2
	}
}

// Increment ...
func (a *PPUAddr) Increment(v uint16) {
	old := a.full
	a.full = domain.Address(uint16(a.full) + v)
	log.Trace("PPUAddr.update %#v => %#v", old, a.full)
}

// ToFullAddress ...
func (a *PPUAddr) ToFullAddress() domain.Address {
	return a.full
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

// PPUInternalRegisters ... PPUの内部状態レジスタ
// http://wiki.nesdev.com/w/index.php/PPU_scrolling
type PPUInternalRegisters struct {
	v uint16 // Current VRAM address
	t uint16 // Temporary VRAM address (15 bits); can also be thought of as the address of the top left onscreen tile.
	x byte   // Fine X scroll (3 bits)
	w byte   // First or second write toggle (1 bit)

	// vとtの仕様
	// bit 0-5   ネームテーブルのX座標(0-31)
	// bit 6-9   ネームテーブルのY座標(0-31)
	// bit 10-11 ネームテーブルの番号(0-3)
	// bit 12-14 タイル内におけるscanlineのオフセット
}

// NewPPUInnerRegisters ...
func NewPPUInnerRegisters() *PPUInternalRegisters {
	return &PPUInternalRegisters{
		v: 0,
		t: 0,
		x: 0,
		w: 0,
	}
}

// UpdateByPPUCtrl ...
func (p *PPUInternalRegisters) UpdateByPPUCtrl(data byte) {
	// t: ...BA.. ........ = d: ......BA
	ba := uint16(data) & 0x03
	p.t = (p.t & 0xF3FF) | (ba << 10)
}

// ClearW ...
func (p *PPUInternalRegisters) ClearW() {
	p.w = 0
}

// UpdateByPPUScroll ...
func (p *PPUInternalRegisters) UpdateByPPUScroll(data byte) {
	if p.w == 0 {
		// t: ....... ...HGFED = d: HGFED...
		hgfed := (uint16(data) & 0xF8) >> 3
		p.t = (p.t & 0xFFE0) | hgfed

		p.x = data & 0x07
		p.w = 1
		return
	}

	// t: CBA..HG FED..... = d: HGFEDCBA
	cba := uint16(data) & 0x07
	hgfed := (uint16(data) & 0xF8)
	p.t = (p.t & 0x0C1F) | (cba << 12) | (hgfed << 2)

	p.w = 0
}

// UpdateByPPUAddr ...
func (p *PPUInternalRegisters) UpdateByPPUAddr(data byte) {
	if p.w == 0 {
		// t: .FEDCBA ........ = d: ..FEDCBA
		fedcba := uint16(data) & 0x3F
		p.t = (p.t & 0xC0FF) | (fedcba << 8)

		// v: X...... ........ = 0
		p.v = p.v & 0x3FFF
		p.w = 1
		return
	}

	// t: ....... HGFEDCBA = d: HGFEDCBA
	p.t = (p.t & 0xFF00) | uint16(data)

	p.v = p.t
	p.w = 0
}

// UpdateHorizontalToLeftEdge ... 水平方向の初期化
func (p *PPUInternalRegisters) UpdateHorizontalToLeftEdge() {
	// v: ....F.. ...EDCBA = t: ....F.. ...EDCBA
	f := p.t & 0x0400
	edcba := p.t & 0x001F
	p.v = (p.v & 0xFBE0) | f | edcba
}

// UpdateVerticalToTopEdge ... 垂直方向の初期化
func (p *PPUInternalRegisters) UpdateVerticalToTopEdge() {
	// v: IHGF.ED CBA..... = t: IHGF.ED CBA.....
	ihgfedcba := p.t & 0x7BE0
	p.v = (p.v & 0x041F) | ihgfedcba
}

// IncrementHorizontal ... 水平方向のインクリメント(8px単位)
func (p *PPUInternalRegisters) IncrementHorizontal() {
	if (p.v & 0x001F) != 0x001F {
		p.v++
		return
	}

	// ネームテーブルをまたぐ
	p.v = p.v & (^uint16(0x001F)) // x = 0
	p.v = p.v ^ uint16(0x0400)    // ネームテーブル切り替え
	return
}

// IncrementVertical ... 垂直方向のインクリメント(1px単位)
func (p *PPUInternalRegisters) IncrementVertical() {
	if (p.v & 0x7000) != 0x7000 {
		// ネームテーブルのY座標は変わらない
		p.v += 0x1000
		return
	}
	p.v = p.v & (^uint16(0x7000))

	// 以下、ネームテーブルのY座標のインクリメントが必要

	y := (p.v & 0x03E0) >> 5
	if y == domain.NameTableHeight-1 {
		// ネームテーブルをまたぐ
		y = 0
		p.v = p.v ^ uint16(0x0800) // ネームテーブル切り替え
	} else if y == 31 {
		// ネームテーブルが範囲外
		y = 0
		// ネームテーブルは切り替えられない
	} else {
		y++
	}
	p.v = (p.v & (^uint16(0x03E0))) | (y << 5)
}

// IncrementV ...
func (p *PPUInternalRegisters) IncrementV(ctrl *PPUCtrl) {
	if ctrl.VRAMAddressIncrementMode == 0 {
		p.v++
	} else {
		p.v = p.v + 32
	}
}

// GetTileIndexAddress ...
func (p *PPUInternalRegisters) GetTileIndexAddress() domain.Address {
	return domain.Address(domain.NameTableBaseAddress | (p.v & 0x0FFF))
}

// GetAttributeAddress ...
func (p *PPUInternalRegisters) GetAttributeAddress() domain.Address {
	offsetX := (p.v & 0x001F) >> 2 // X座標(=ネームテーブルのX座標の上位3bit分)
	offsetY := (p.v & 0x03E0) >> 4 // Y座標(=ネームテーブルのY座標の上位3bit分)
	offsetNameTtbl := p.v & 0x0C00
	return domain.Address(domain.AttributeTableBaseAddress | offsetNameTtbl | offsetY | offsetX)
}

// GetAttributeIndex ...
func (p *PPUInternalRegisters) GetAttributeIndex() byte {
	offsetX := (p.v & 0x0002) >> 1
	offsetY := (p.v & 0x0040) >> 5
	return byte(offsetY + offsetX)
}

// GetTilePatternLowAddress ...
func (p *PPUInternalRegisters) GetTilePatternLowAddress(tileIndex byte) domain.Address {
	offsetY := (p.v & 0x7000) >> 12
	return domain.Address((uint16(tileIndex) << 4) | offsetY)
}

// GetTilePatternHighAddress ...
func (p *PPUInternalRegisters) GetTilePatternHighAddress(tileIndex byte) domain.Address {
	offsetY := (p.v & 0x7000) >> 12
	return domain.Address((uint16(tileIndex) << 4) | offsetY | 0x0008)
}

// GetFineX ...
func (p *PPUInternalRegisters) GetFineX() byte {
	return p.x
}
