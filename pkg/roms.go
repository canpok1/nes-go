package rom

import (
	"fmt"
)

// INESHeader ...
// 仕様: https://wiki.nesdev.com/w/index.php/INES
type INESHeader struct {
	PRGROMSize uint8 // 4: Size of PRG ROM in 16 KB units
	CHRROMSize uint8 // 5: Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
}

// PRGROM ...
type PRGROM []byte

// CHRROM ...
type CHRROM []byte

// ROM ...
type ROM struct {
	header *INESHeader
	prgrom *PRGROM
	chrrom *CHRROM
}

func fetchINESHeader(rom []byte) (*INESHeader, error) {
	if rom == nil {
		return nil, fmt.Errorf("failed to fetch, rom is nil")
	}
	if len(rom) < 6 {
		return nil, fmt.Errorf("failed to fetch, rom is too short")
	}

	prg := uint8(rom[4])
	chr := uint8(rom[5])

	return &INESHeader{
		PRGROMSize: prg,
		CHRROMSize: chr,
	}, nil
}

func FetchROM(rom []byte) (*ROM, error) {
	h, err := fetchINESHeader(rom)
	if err != nil {
		return nil, err
	}

	begin := 0x0010
	prgromEnd := 0x0010 + int(h.CHRROMSize)*0x4000
	chrromEnd := prgromEnd + int(h.PRGROMSize)*0x2000

	p := PRGROM(rom[begin:prgromEnd])
	c := CHRROM(rom[prgromEnd:chrromEnd])
	return &ROM{
		header: h,
		prgrom: &p,
		chrrom: &c,
	}, nil
}
