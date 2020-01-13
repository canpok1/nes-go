package domain

// CPU ...
type CPU interface {
	SetBus(Bus)
	Run() (int, error)
	String() string
	ReceiveNMI() error
}

// PPU ...
type PPU interface {
	SetBus(Bus)
	ReadRegisters(Address) (byte, error)
	WriteRegisters(Address, byte) error
	Run(int) (*Screen, error)
	String() string
}

// Bus ...
type Bus interface {
	Setup(*ROM, PPU, CPU, *VRAM, Pad, Pad)
	ReadByCPU(Address) (byte, error)
	WriteByCPU(Address, byte) error
	ReadByPPU(Address) (byte, error)
	WriteByPPU(Address, byte) error
	GetTileNo(uint8, NameTablePoint) (uint8, error)
	GetTilePattern(uint8, uint8) *TilePattern
	GetPaletteNo(NameTablePoint, byte) (uint8, error)
	GetPalette(uint8) *Palette
	GetAttribute(uint8, NameTablePoint) (byte, error)
	SendNMI() error
}

// Renderer ...
type Renderer interface {
	Render(*Screen) error
}
