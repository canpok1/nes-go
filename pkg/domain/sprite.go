package domain

// Sprite ...
type Sprite struct {
	Y         byte
	TileIndex byte
	Attribute byte
	X         byte
}

// SpriteImage ...
type SpriteImage struct {
	X uint8
	Y uint8
	R [][]byte
	G [][]byte
	B [][]byte
	A [][]byte
}

// NewSpriteImage ...
func NewSpriteImage(x, y uint8, t *TileImage) *SpriteImage {
	return &SpriteImage{
		X: x,
		Y: y,
		R: t.R,
		G: t.G,
		B: t.B,
		A: t.A,
	}
}
