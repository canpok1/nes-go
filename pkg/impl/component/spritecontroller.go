package component

import "nes-go/pkg/domain"

import "nes-go/pkg/log"

import "image/color"

const (
	MaxSpriteCount = 8
	SpriteByteSize = 4 //単位はbyte
)

// SpriteController ...
type SpriteController struct {
	oam              []byte
	lastWriteAddress uint8

	oam2             []domain.Sprite
	patternRegisterL []ShiftRegister8
	patternRegisterH []ShiftRegister8
	latches          []byte
	counters         []int16

	bus domain.Bus

	n             uint8 // 評価対象スプライト番号(0-63)
	secondarySize uint8
	fetchedCount  uint8
}

// NewSpriteController ...
func NewSpriteController() *SpriteController {
	c := SpriteController{
		oam:              make([]byte, 256),
		lastWriteAddress: 0,

		oam2:             make([]domain.Sprite, MaxSpriteCount),
		patternRegisterL: make([]ShiftRegister8, MaxSpriteCount),
		patternRegisterH: make([]ShiftRegister8, MaxSpriteCount),
		latches:          make([]byte, MaxSpriteCount),
		counters:         make([]int16, MaxSpriteCount),
	}
	for i := 0; i < 8; i++ {
		c.patternRegisterL[i] = ShiftRegister8{}
		c.patternRegisterH[i] = ShiftRegister8{}
	}
	return &c
}

// SetBus ...
func (s *SpriteController) SetBus(bus domain.Bus) {
	s.bus = bus
}

// WriteOAM ...
func (s *SpriteController) WriteOAM(oamaddr uint8, b byte) {
	s.oam[oamaddr] = b
	s.lastWriteAddress = oamaddr
}

// ClearSecondaryOAM ...
func (s *SpriteController) ClearSecondaryOAM(spriteIdx, byteIdx uint16) {
	switch byteIdx {
	case 0:
		s.oam2[spriteIdx].Y = 0xFF
	case 1:
		s.oam2[spriteIdx].TileIndex = 0xFF
	case 2:
		s.oam2[spriteIdx].Attribute = 0xFF
	case 3:
		s.oam2[spriteIdx].X = 0xFF
	default:
		log.Warn("SpriteController.ClearSecondaryOAM byteIdx out of range byteIdx=%v", byteIdx)
	}
	s.n = 0
	s.secondarySize = 0
	s.fetchedCount = 0
}

// EvaluateSprite ... 対象スプライトをセカンダリOAMにコピー
func (s *SpriteController) EvaluateSprite(scanline uint16) {
	if s.secondarySize >= 8 {
		return
	}

	idx := s.n << 2
	top := uint16(s.oam[idx]) + 1
	btm := top + domain.SpriteHeight - 1

	var y uint16
	if scanline <= 239 {
		y = scanline + 1
	}

	sprite := domain.Sprite{
		Y:         s.oam[idx],
		TileIndex: s.oam[idx+1],
		Attribute: s.oam[idx+2],
		X:         s.oam[idx+3],
	}

	if y >= top && y <= btm && sprite.X < 255 {
		// セカンダリにコピー
		s.oam2[s.secondarySize] = sprite
		log.Info("copy to secondaryOAM; scanline: %v, sprite: %v", scanline, s.oam2[s.secondarySize])

		s.secondarySize++
	}

	s.n = (s.n + 1) & 0x3F
}

// FetchSprite ... セカンダリOAMからシフトレジスタ等へコピー
func (s *SpriteController) FetchSprite(scanline uint16, patternTblIdx uint8) {
	if s.fetchedCount >= 8 {
		return
	}

	idx := uint16(s.fetchedCount)

	if s.fetchedCount >= s.secondarySize {
		s.patternRegisterL[idx].Set(0xFF)
		s.patternRegisterH[idx].Set(0xFF)
		s.latches[idx] = 0xFF
		s.counters[idx] = int16(0xFF)
	} else {
		sprite := s.oam2[idx]
		yOffset := scanline - uint16(sprite.Y)
		pattern := s.bus.GetTilePattern(patternTblIdx, sprite.TileIndex)

		s.patternRegisterL[idx].Set(swapbit((*pattern)[yOffset]))
		s.patternRegisterH[idx].Set(swapbit((*pattern)[yOffset+8]))
		s.latches[idx] = sprite.Attribute
		s.counters[idx] = int16(sprite.X)
	}

	s.fetchedCount++
}

// Shift ...
func (s *SpriteController) Shift() {
	for i, counter := range s.counters {
		if counter > 0 {
			s.counters[i]--
			continue
		}

		s.patternRegisterL[i].Shift()
		s.patternRegisterH[i].Shift()
	}
}

// MakePixel ... ピクセルを生成し属性とともに返す
func (s *SpriteController) MakePixel() (color.RGBA, byte) {
	for i := 0; i < 8; i++ {
		counter := s.counters[i]
		if counter != 0 {
			continue
		}

		patternL := s.patternRegisterL[i].Get() & 0x01
		patternH := s.patternRegisterH[i].Get() & 0x01
		pattern := (patternH << 1) | patternL
		if pattern == 0 {
			continue
		}

		attr := s.latches[i]
		palette := s.bus.GetPalette(attr & 0x0F)

		r, g, b := palette.GetColor(pattern)

		return color.RGBA{R: r, G: g, B: b, A: 0xFF}, attr
	}

	// 透明
	return color.RGBA{R: 0, G: 0, B: 0, A: 0}, 0x00
}
