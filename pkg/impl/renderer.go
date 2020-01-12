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
func (m *Renderer) Render(tis [][]domain.TileImage, sis []domain.SpriteImage) error {
	p := toPixels(tis, sis)
	m.imageBuf.ReplacePixels(p)

	m.fps = 1 / time.Since(m.lastRenderedTime).Seconds()
	m.lastRenderedTime = time.Now()

	return nil
}

// toPixels ...
func toPixels(tImages [][]domain.TileImage, sImages []domain.SpriteImage) []byte {
	pixels := make([]byte, 4*domain.ResolutionHeight*domain.ResolutionWidth)

	idx := 0
	for y := uint16(0); y < domain.ResolutionHeight; y++ {
		for x := uint16(0); x <= domain.ResolutionWidth-1; x++ {
			r, g, b, a := getPixel(tImages, domain.MonitorX(x), domain.MonitorY(y))

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

	for _, sImage := range sImages {
		baseIdx := int(sImage.Y)*domain.ResolutionWidth*4 + int(sImage.X)*4

		for offsetY := 0; offsetY < int(sImage.H); offsetY++ {
			for offsetX := 0; offsetX < int(sImage.W); offsetX++ {
				idx := baseIdx + (offsetY * domain.ResolutionWidth * 4) + (offsetX * 4)

				r := sImage.R[offsetY][offsetX]
				g := sImage.G[offsetY][offsetX]
				b := sImage.B[offsetY][offsetX]
				a := sImage.A[offsetY][offsetX]

				pixels[idx] = r
				pixels[idx+1] = g
				pixels[idx+2] = b
				pixels[idx+3] = a

				if sImage.IsForeground {
					// 上書き
					pixels[idx] = r
					pixels[idx+1] = g
					pixels[idx+2] = b
					pixels[idx+3] = a
				} else {
					// 背景が透明の場合だけ書き込む
				}
			}
		}
	}

	return pixels
}

func getPixel(tis [][]domain.TileImage, x domain.MonitorX, y domain.MonitorY) (r, g, b, a byte) {
	t := tis[y/domain.SpriteHeight][x/domain.SpriteWidth]

	innerY := uint8(y) % domain.SpriteHeight
	innerX := uint8(x) % domain.SpriteWidth

	r = t.R[innerY][innerX]
	g = t.G[innerY][innerX]
	b = t.B[innerY][innerX]
	a = t.A[innerY][innerX]

	return
}
