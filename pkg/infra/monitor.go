package infra

import (
	"fmt"
	"image/color"
	"nes-go/pkg/model"

	"github.com/hajimehoshi/ebiten"
)

// Monitor ...
type Monitor struct {
	width    int
	height   int
	scale    float64
	title    string
	imageBuf *ebiten.Image
}

// NewMonitor ...
func NewMonitor(w int, h int, scale float64, title string) (*Monitor, error) {
	imageBuf, err := ebiten.NewImage(model.ResolutionWidth, model.ResolutionHeight, ebiten.FilterDefault)
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

	return nil
}

// Run ...
func (m *Monitor) Run() error {
	return ebiten.Run(m.update, m.height, m.width, m.scale, m.title)
}

// Render ...
func (m *Monitor) Render(ss [][]model.SpriteImage) error {
	for y, simgs := range ss {
		for x, simg := range simgs {
			img, err := makeImage(&simg)
			if err != nil {
				return fmt.Errorf("failed to render; err: %w", err)
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))

			m.imageBuf.DrawImage(img, op)
		}
	}

	return nil
}

// makeImage ...
func makeImage(si *model.SpriteImage) (*ebiten.Image, error) {
	img, err := ebiten.NewImage(model.SpriteWidth, model.SpriteHeight, ebiten.FilterDefault)
	if err != nil {
		return nil, fmt.Errorf("failed to makeImage; err: %w", err)
	}

	for y := 0; y < model.SpriteHeight; y++ {
		for x := 0; x < model.SpriteWidth; x++ {
			img.Set(x, y, color.NRGBA{
				R: si.R[y][x],
				G: si.G[y][x],
				B: si.B[y][x],
				A: 0xFF,
			})
		}
	}

	return img, nil
}
