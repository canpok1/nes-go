package domain

type Screen struct {
	TileImages           [][]TileImage
	SpriteImages         []SpriteImage
	EnableSpriteMask     bool
	EnableBackgroundMask bool
}
