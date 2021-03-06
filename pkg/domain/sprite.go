package domain

// Sprite ...
type Sprite struct {
	Y         byte // Y座標-1が格納される
	TileIndex byte
	Attribute byte
	X         byte
}

// ContainsY ...
func (s Sprite) ContainsY(y uint16) bool {
	if y < uint16(s.Y+1) {
		return false
	}
	if y >= uint16(s.Y+1)+ResolutionHeight {
		return false
	}

	return true
}

// SpriteImage ...
type SpriteImage struct {
	X                    uint16
	Y                    uint16
	W                    uint16
	H                    uint16
	R                    [][]byte
	G                    [][]byte
	B                    [][]byte
	A                    [][]byte
	IsForeground         bool
	EnableFlipHorizontal bool
	EnableFlipVertical   bool
}

// NewSpriteImage ...
func NewSpriteImage(x, y uint16, t *TileImage, isForeground, enableFlipHorizontal, enableFlipVertical bool) *SpriteImage {
	return &SpriteImage{
		X:                    x,
		Y:                    y,
		W:                    t.W,
		H:                    t.H,
		R:                    t.R,
		G:                    t.G,
		B:                    t.B,
		A:                    t.A,
		IsForeground:         isForeground,
		EnableFlipHorizontal: enableFlipHorizontal,
		EnableFlipVertical:   enableFlipVertical,
	}
}

// ContainsX ...
func (s SpriteImage) ContainsX(x uint16) bool {
	if x < s.X {
		return false
	}
	if x >= s.X+s.W {
		return false
	}

	return true
}

// ContainsY ...
func (s SpriteImage) ContainsY(y uint16) bool {
	if y < s.Y {
		return false
	}
	if y >= s.Y+s.H {
		return false
	}

	return true
}
