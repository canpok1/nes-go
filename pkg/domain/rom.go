package domain

import (
	"io/ioutil"
	"nes-go/pkg/log"
	"os"

	"golang.org/x/xerrors"
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

// GetTilePattern ...
func (c *CHRROM) GetTilePattern(patternTblIdx, no uint8) *TilePattern {
	begin := uint16(no) * 0x0010
	if patternTblIdx == 1 {
		begin = 0x1000 + begin
	}
	end := begin + 0x000F + 1
	s := TilePattern((*c)[begin:end])
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
		return nil, xerrors.Errorf("failed to open rom\nromPath: %#v\nerr: %w", p, err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, xerrors.Errorf("failed to open rom\nromPath: %#v\nerr: %w", p, err)
	}
	return b, nil
}

func parseINESHeader(rom []byte) (*INESHeader, error) {
	if rom == nil {
		return nil, xerrors.New("failed to parse, rom is nil")
	}
	if len(rom) < 6 {
		return nil, xerrors.New("failed to parse, rom is too short")
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

	log.Trace("rom header: %#v", h)

	begin := 0x0010
	prgromEnd := 0x0010 + int(h.PRGROMSize)*0x4000
	chrromEnd := prgromEnd + int(h.CHRROMSize)*0x2000

	log.Trace("prg-rom byte index: %#v-%#v", begin, (prgromEnd - 1))
	log.Trace("chr-rom byte index: %#v-%#v", prgromEnd, (chrromEnd - 1))

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
	log.Trace("fetch[rom]: %v", romPath)
	f, err := readFile(romPath)
	if err != nil {
		return nil, xerrors.Errorf("failed fetch rom: %w", err)
	}

	r, err := parseROM(f)
	if err != nil {
		return nil, xerrors.Errorf("failed fetch rom: %w", err)
	}

	return r, nil
}
