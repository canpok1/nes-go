package infra

import (
	"fmt"
	"nes-go/pkg/domain"
	"nes-go/pkg/model"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Monitor ...
type Monitor struct {
	width    int
	height   int
	scale    float64
	title    string
	imageBuf *ebiten.Image

	lastRenderedTime time.Time
	fps              float64
}

// NewMonitor ...
func NewMonitor(w int, h int, scale float64, title string) (*Monitor, error) {
	imageBuf, err := ebiten.NewImage(domain.ResolutionWidth, domain.ResolutionHeight, ebiten.FilterDefault)
	if err != nil {
		return nil, fmt.Errorf("failed to NewMonitor; err: %w", err)
	}

	return &Monitor{
		width:    w,
		height:   h,
		scale:    scale,
		title:    title,
		imageBuf: imageBuf,
	}, nil
}

// update ...
func (m *Monitor) update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.DrawImage(m.imageBuf, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", m.fps))

	return nil
}

// Run ...
func (m *Monitor) Run() error {
	return ebiten.Run(m.update, m.height, m.width, m.scale, m.title)
}

// Render ...
func (m *Monitor) Render(sis [][]domain.SpriteImage) error {
	p := toPixels(sis)
	m.imageBuf.ReplacePixels(p)

	m.fps = 1 / time.Since(m.lastRenderedTime).Seconds()
	m.lastRenderedTime = time.Now()

	return nil
}

// toPixels ...
func toPixels(sis [][]domain.SpriteImage) []byte {
	pixels := make([]byte, 4*domain.ResolutionHeight*domain.ResolutionWidth)

	idx := 0
	for y := 0; y < domain.ResolutionHeight; y++ {
		for x := 0; x < domain.ResolutionWidth; x++ {
			r, g, b, a := getPixel(sis, model.MonitorX(x), model.MonitorY(y))

			pixels[idx] = r
			idx++

			pixels[idx] = g
			idx++

			pixels[idx] = b
			idx++

			pixels[idx] = a
			idx++
		}
	}

	return pixels
}

func getPixel(sis [][]domain.SpriteImage, x model.MonitorX, y model.MonitorY) (r, g, b, a byte) {
	s := sis[y/domain.SpriteHeight][x/domain.SpriteWidth]

	iy := y % domain.SpriteHeight
	ix := x % domain.SpriteWidth

	return s.R[iy][ix], s.G[iy][ix], s.B[iy][ix], 0xFF
}
