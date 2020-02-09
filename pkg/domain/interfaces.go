package domain

// CPU ...
type CPU interface {
	SetBus(Bus)
	SetRecorder(*Recorder)
	Run() (int, error)
	String() string
	ReceiveNMI(active bool)
}

// PPU ...
type PPU interface {
	SetBus(Bus)
	SetRecorder(*Recorder)
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
	SendNMI(active bool)
}

// Renderer ...
type Renderer interface {
	Run() error
	Render(*Screen) error
}
