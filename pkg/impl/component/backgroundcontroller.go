package component

import (
	"image/color"
	"nes-go/pkg/domain"
)

// BackgroundController ...
type BackgroundController struct {
	patternRegisterL   *ShiftRegister16bit
	patternRegisterH   *ShiftRegister16bit
	attributeRegisterL *ShiftRegister16bit
	attributeRegisterH *ShiftRegister16bit

	NextTileIndex      byte
	NextAttributeTable byte
	NextTilePatternL   byte
	NextTilePatternH   byte

	bus domain.Bus
}

// NewBackgroundController ...
func NewBackgroundController() *BackgroundController {
	return &BackgroundController{
		patternRegisterL:   &ShiftRegister16bit{},
		patternRegisterH:   &ShiftRegister16bit{},
		attributeRegisterL: &ShiftRegister16bit{},
		attributeRegisterH: &ShiftRegister16bit{},
		NextTileIndex:      0,
		NextAttributeTable: 0,
		NextTilePatternL:   0,
		NextTilePatternH:   0,
	}
}

// SetBus ...
func (b *BackgroundController) SetBus(bus domain.Bus) {
	b.bus = bus
}

// Shift ...
func (b *BackgroundController) Shift() {
	b.patternRegisterL.Shift()
	b.patternRegisterH.Shift()
	b.attributeRegisterL.Shift()
	b.attributeRegisterH.Shift()
}

// SetupNextData ... フェッチしたデータのセットアップ
func (b *BackgroundController) SetupNextData(attrIdx byte) {
	b.patternRegisterL.SetHigh(swapbit(b.NextTilePatternL))
	b.patternRegisterH.SetHigh(swapbit(b.NextTilePatternH))

	attr := (b.NextAttributeTable & (0x03 << attrIdx)) >> attrIdx

	// 1タイルにおける属性（適用するパレット番号）は同じなので全ビットを揃える
	if (attr & 0x01) == 0x01 {
		b.attributeRegisterL.SetHigh(0xFF)
	}
	if (attr & 0x02) == 0x02 {
		b.attributeRegisterH.SetHigh(0xFF)
	}
}

// MakePixel ...
func (b *BackgroundController) MakePixel() color.RGBA {
	attrL := b.attributeRegisterL.GetLow() & 0x01
	attrH := b.attributeRegisterH.GetLow() & 0x01
	attr := (attrH << 1) | attrL

	palette := b.bus.GetPalette(attr)

	patternL := b.patternRegisterL.GetLow() & 0x01
	patternH := b.patternRegisterH.GetLow() & 0x01
	pattern := (patternH << 1) | patternL

	red, green, blue := palette.GetColor(pattern)
	return color.RGBA{R: red, G: green, B: blue, A: 0xFF}
}

func swapbit(b byte) byte {
	swapped := byte(0)

	for i := 0; i < 8; i++ {
		swapped = swapped << 1
		if (b & 0x01) == 0x01 {
			swapped = swapped | 0x01
		}
		b = b >> 1
	}

	return swapped
}
