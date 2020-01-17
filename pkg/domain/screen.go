package domain

type Screen struct {
	TileImages            [][]TileImage
	SpriteImages          []SpriteImage
	DisableSpriteMask     bool
	DisableBackgroundMask bool
}
