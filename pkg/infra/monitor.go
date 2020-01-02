package infra

import (
	"github.com/hajimehoshi/ebiten"
)

// Monitor ...
type Monitor struct {
	width  int
	height int
	scale  float64
	title  string
	pixels []byte
}

// NewMonitor ...
func NewMonitor(w int, h int, scale float64, title string) *Monitor {
	pixels := make([]byte, 4*w*h)

	return &Monitor{
		width:  w,
		height: h,
		scale:  scale,
		title:  title,
		pixels: pixels,
	}
}

// update ...
func (m *Monitor) update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.ReplacePixels(m.pixels)

	return nil
}

// Run ...
func (m *Monitor) Run() error {
	return ebiten.Run(m.update, m.height, m.width, m.scale, m.title)
}

// SetPixels ...
func (m *Monitor) SetPixels(p []byte) error {
	m.pixels = p
	return nil
}
