package domain

// VRAM ...
type VRAM struct {
	NameTable0      []byte
	AttributeTable0 []byte
	NameTable1      []byte
	AttributeTable1 []byte
	NameTable2      []byte
	AttributeTable2 []byte
	NameTable3      []byte
	AttributeTable3 []byte

	BackgroundPalette []Palette
	SpritePalette     []Palette
}

// NewVRAM ...
func NewVRAM() *VRAM {
	bp := []Palette{}
	sp := []Palette{}
	for i := 0; i < 4; i++ {
		newBP := NewPalette()
		newSP := NewPalette()
		bp = append(bp, *newBP)
		sp = append(sp, *newSP)
	}

	return &VRAM{
		NameTable0:        make([]byte, 0x03C0),
		AttributeTable0:   make([]byte, 0x0040),
		NameTable1:        make([]byte, 0x03C0),
		AttributeTable1:   make([]byte, 0x0040),
		NameTable2:        make([]byte, 0x03C0),
		AttributeTable2:   make([]byte, 0x0040),
		NameTable3:        make([]byte, 0x03C0),
		AttributeTable3:   make([]byte, 0x0040),
		BackgroundPalette: bp,
		SpritePalette:     sp,
	}
}
