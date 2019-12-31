package model

import (
	"fmt"
)

// INESHeader ...
// 仕様: https://wiki.nesdev.com/w/index.php/INES
type INESHeader struct {
	PRGROMSize uint8 // 4: Size of PRG ROM in 16 KB units
	CHRROMSize uint8 // 5: Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
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
