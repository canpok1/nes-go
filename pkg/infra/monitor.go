package infra

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Monitor ...
type Monitor struct {
	width  int
	height int
	scale  float64
	title  string
}

// NewMonitor ...
func NewMonitor(w int, h int, scale float64, title string) *Monitor {
	return &Monitor{
		width:  w,
		height: h,
		scale:  scale,
		title:  title,
	}
}

// update ...
func (m *Monitor) update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	ebitenutil.DebugPrint(screen, "sample monitor")
	return nil
}

// Run ...
func (m *Monitor) Run() error {
	return ebiten.Run(m.update, m.height, m.width, m.scale, m.title)
}
