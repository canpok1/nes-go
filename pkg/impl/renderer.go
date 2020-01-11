package impl

import (
	"fmt"
	"nes-go/pkg/domain"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Renderer ...
type Renderer struct {
	width    int
	height   int
	scale    float64
	title    string
	imageBuf *ebiten.Image

	lastRenderedTime time.Time
	fps              float64
}

// NewRenderer ...
func NewRenderer(w int, h int, scale float64, title string) (*Renderer, error) {
	imageBuf, err := ebiten.NewImage(domain.ResolutionWidth, domain.ResolutionHeight, ebiten.FilterDefault)
	if err != nil {
		return nil, fmt.Errorf("failed to NewMonitor; err: %w", err)
	}

	return &Renderer{
		width:    w,
		height:   h,
		scale:    scale,
		title:    title,
		imageBuf: imageBuf,
	}, nil
}

// update ...
func (m *Renderer) update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.DrawImage(m.imageBuf, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", m.fps))

	return nil
}

// Run ...
func (m *Renderer) Run() error {
	return ebiten.Run(m.update, m.height, m.width, m.scale, m.title)
}

// Render ...
func (m *Renderer) Render(sis [][]domain.TileImage) error {
	p := toPixels(sis)
	m.imageBuf.ReplacePixels(p)

	m.fps = 1 / time.Since(m.lastRenderedTime).Seconds()
	m.lastRenderedTime = time.Now()

	return nil
}

// toPixels ...
func toPixels(sis [][]domain.TileImage) []byte {
	pixels := make([]byte, 4*domain.ResolutionHeight*domain.ResolutionWidth)

	idx := 0
	for y := 0; y < domain.ResolutionHeight; y++ {
		for x := 0; x < domain.ResolutionWidth; x++ {
			r, g, b, a := getPixel(sis, MonitorX(x), MonitorY(y))

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

func getPixel(sis [][]domain.TileImage, x MonitorX, y MonitorY) (r, g, b, a byte) {
	s := sis[y/domain.SpriteHeight][x/domain.SpriteWidth]

	iy := y % domain.SpriteHeight
	ix := x % domain.SpriteWidth

	return s.R[iy][ix], s.G[iy][ix], s.B[iy][ix], 0xFF
}
