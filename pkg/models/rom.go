package model

import (
	"fmt"
	"io/ioutil"
	"os"
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

func FetchROM(romPath string) (*ROM, error) {
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
