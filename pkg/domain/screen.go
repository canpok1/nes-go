package domain

import "image/color"

// Screen ...
type Screen struct {
	TileImages            [][]TileImage
	SpriteImages          []SpriteImage
	DisableSpriteMask     bool
	DisableBackgroundMask bool
	Images                [][]color.RGBA
}
