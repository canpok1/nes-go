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
func (m *Renderer) Render(sis [][]domain.Tile) error {
	p := toPixels(sis)
	m.imageBuf.ReplacePixels(p)

	m.fps = 1 / time.Since(m.lastRenderedTime).Seconds()
	m.lastRenderedTime = time.Now()

	return nil
}

// toPixels ...
func toPixels(sis [][]domain.Tile) []byte {
	pixels := make([]byte, 4*domain.ResolutionHeight*domain.ResolutionWidth)

	idx := 0
	for y := 0; y < domain.ResolutionHeight; y++ {
		for x := 0; x < domain.ResolutionWidth; x++ {
			r, g, b, a := getPixel(sis, domain.MonitorX(x), domain.MonitorY(y))

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

func getPixel(sis [][]domain.Tile, x domain.MonitorX, y domain.MonitorY) (r, g, b, a byte) {
	t := sis[y/domain.SpriteHeight][x/domain.SpriteWidth]

	iy := y % domain.SpriteHeight
	ix := x % domain.SpriteWidth

	if t.Sprite == nil && t.Background == nil {
		return
	}

	if t.Sprite != nil && t.Sprite.A[iy][ix] != 0 {
		r = t.Sprite.R[iy][ix]
		g = t.Sprite.G[iy][ix]
		b = t.Sprite.B[iy][ix]
		a = t.Sprite.A[iy][ix]
	}
	if t.Background != nil && t.Background.A[iy][ix] != 0 {
		r = t.Background.R[iy][ix]
		g = t.Background.G[iy][ix]
		b = t.Background.B[iy][ix]
		a = t.Background.A[iy][ix]
	}

	return
}
