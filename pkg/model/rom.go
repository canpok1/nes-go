package model

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/canpok1/nes-go/pkg/log"
)

// INESHeader ...
// 仕様: https://wiki.nesdev.com/w/index.php/INES
type INESHeader struct {
	PRGROMSize uint8 // 4: Size of PRG ROM in 16 KB units
	CHRROMSize uint8 // 5: Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
}

// Sprite ...
type Sprite []byte

// toColorMap
func (s Sprite) toColorMap() [][]byte {
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

	log.Debug("Sprite.toColorMap => %#v", colorMap)

	return colorMap
}

// ToSpriteImage ...
func (s Sprite) ToSpriteImage(p *Palette) *SpriteImage {
	r := make([][]byte, SpriteHeight)
	g := make([][]byte, SpriteHeight)
	b := make([][]byte, SpriteHeight)

	colorMap := s.toColorMap()
	for y, line := range colorMap {
		r[y] = make([]byte, SpriteWidth)
		g[y] = make([]byte, SpriteWidth)
		b[y] = make([]byte, SpriteWidth)
		for x, paletteNo := range line {
			cIndex := (*p)[paletteNo]
			c := colors[cIndex]
			r[y][x] = c[0]
			g[y][x] = c[1]
			b[y][x] = c[2]
		}
	}
	return &SpriteImage{
		R: r,
		G: g,
		B: b,
	}
}

// SpriteImage ...
type SpriteImage struct {
	R [][]byte
	G [][]byte
	B [][]byte
}

// NewSpriteImage ...
func NewSpriteImage() *SpriteImage {
	r := make([][]byte, SpriteHeight)
	g := make([][]byte, SpriteHeight)
	b := make([][]byte, SpriteHeight)

	for y := 0; y < SpriteHeight; y++ {
		r[y] = make([]byte, SpriteWidth)
		g[y] = make([]byte, SpriteWidth)
		b[y] = make([]byte, SpriteWidth)
	}

	return &SpriteImage{
		R: r,
		G: g,
		B: b,
	}
}

// PRGROM ...
type PRGROM []byte

// CHRROM ...
type CHRROM []byte

// GetSprite ...
func (c *CHRROM) GetSprite(no uint8) *Sprite {
	begin := uint16(no) * 0x0010
	end := begin + 0x000F + 1
	s := Sprite((*c)[begin:end])
	return &s
}

// ROM ...
type ROM struct {
	Header *INESHeader
	Prgrom *PRGROM
	Chrrom *CHRROM
}

func readFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("failed to open rom\nromPath: %#v\nerr: %w", p, err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to open rom\nromPath: %#v\nerr: %w", p, err)
	}
	return b, nil
}

func parseINESHeader(rom []byte) (*INESHeader, error) {
	if rom == nil {
		return nil, fmt.Errorf("failed to parse, rom is nil")
	}
	if len(rom) < 6 {
		return nil, fmt.Errorf("failed to parse, rom is too short")
	}

	prg := uint8(rom[4])
	chr := uint8(rom[5])

	return &INESHeader{
		PRGROMSize: prg,
		CHRROMSize: chr,
	}, nil
}

func parseROM(rom []byte) (*ROM, error) {
	h, err := parseINESHeader(rom)
	if err != nil {
		return nil, err
	}

	log.Debug("rom header: %#v", h)

	begin := 0x0010
	prgromEnd := 0x0010 + int(h.PRGROMSize)*0x4000
	chrromEnd := prgromEnd + int(h.CHRROMSize)*0x2000

	log.Debug("prg-rom byte index: %#v-%#v", begin, (prgromEnd - 1))
	log.Debug("chr-rom byte index: %#v-%#v", prgromEnd, (chrromEnd - 1))

	p := PRGROM(rom[begin:prgromEnd])
	c := CHRROM(rom[prgromEnd:chrromEnd])
	return &ROM{
		Header: h,
		Prgrom: &p,
		Chrrom: &c,
	}, nil
}

// FetchROM ...
func FetchROM(romPath string) (*ROM, error) {
	log.Info("fetch[rom]: %v", romPath)
	f, err := readFile(romPath)
	if err != nil {
		return nil, fmt.Errorf("failed fetch rom; %w", err)
	}

	r, err := parseROM(f)
	if err != nil {
		return nil, fmt.Errorf("failed fetch rom; %w", err)
	}

	return r, nil
}
