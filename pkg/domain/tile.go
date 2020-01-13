package domain

// TilePattern ...
type TilePattern []byte

// toColorMap
func (s TilePattern) toColorMap() [][]byte {
	colorMap := make([][]byte, 8)
	for y := 0; y < 8; y++ {
		colorMap[y] = make([]byte, 8)
	}

	for i, line := range []byte(s) {
		y := i % 8
		indexShift := i / 8
		for x := 0; x < 8; x++ {
			filterShift := 7 - x
			colorIndex := ((line & (0x01 << filterShift)) >> filterShift) << indexShift
			colorMap[y][x] = colorMap[y][x] + colorIndex
		}
	}

	return colorMap
}

// ToTileImage ...
func (s TilePattern) ToTileImage(p *Palette) *TileImage {
	r := make([][]byte, SpriteHeight)
	g := make([][]byte, SpriteHeight)
	b := make([][]byte, SpriteHeight)
	a := make([][]byte, SpriteHeight)

	colorMap := s.toColorMap()
	for y, line := range colorMap {
		r[y] = make([]byte, SpriteWidth)
		g[y] = make([]byte, SpriteWidth)
		b[y] = make([]byte, SpriteWidth)
		a[y] = make([]byte, SpriteWidth)

		if p == nil {
			continue
		}

		for x, paletteNo := range line {
			cIndex := (*p)[paletteNo]
			c := colors[cIndex]
			r[y][x] = c[0]
			g[y][x] = c[1]
			b[y][x] = c[2]
			a[y][x] = 0xFF
		}
	}
	return &TileImage{
		R: r,
		G: g,
		B: b,
		A: a,
		W: SpriteWidth,
		H: SpriteHeight,
	}
}

// TileImage ...
type TileImage struct {
	W uint16
	H uint16
	R [][]byte
	G [][]byte
	B [][]byte
	A [][]byte
}

// NewTileImage ...
func NewTileImage() *TileImage {
	r := make([][]byte, SpriteHeight)
	g := make([][]byte, SpriteHeight)
	b := make([][]byte, SpriteHeight)
	a := make([][]byte, SpriteHeight)

	for y := 0; y < SpriteHeight; y++ {
		r[y] = make([]byte, SpriteWidth)
		g[y] = make([]byte, SpriteWidth)
		b[y] = make([]byte, SpriteWidth)
		a[y] = make([]byte, SpriteWidth)
	}

	return &TileImage{
		W: SpriteWidth,
		H: SpriteHeight,
		R: r,
		G: g,
		B: b,
		A: a,
	}
}
