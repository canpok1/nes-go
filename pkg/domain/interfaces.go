package domain

// CPU ...
type CPU interface {
	SetBus(Bus)
	Run() (int, error)
	String() string
	ReceiveNMI()
}

// PPU ...
type PPU interface {
	SetBus(Bus)
	ReadRegisters(Address) (byte, error)
	WriteRegisters(Address, byte) error
	Run(int) ([][]SpriteImage, error)
	String() string
}

// Bus ...
type Bus interface {
	Setup(ROM, PPU)
	ReadByCPU(Address) (byte, error)
	WriteByCPU(Address, byte) error
	ReadByPPU(Address) (byte, error)
	WriteByPPU(Address, byte) error
	GetSpriteNo(uint8, NameTablePoint) (uint8, error)
	GetSprite(uint8) *Sprite
	GetPaletteNo(NameTablePoint) (uint8, error)
	GetBackgroundPalette(uint8) *Palette
	GetSpritePalette(uint8) *Palette
	SendNMI()
}

// Renderer ...
type Renderer interface {
	Render([][]SpriteImage) error
}
