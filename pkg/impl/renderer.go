package impl

import (
	"fmt"
	"nes-go/pkg/domain"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/xerrors"
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

	enableDebugPrint bool
}

// NewRenderer ...
func NewRenderer(scale float64, title string, enableDebugPrint bool) (domain.Renderer, error) {
	imageBuf, err := ebiten.NewImage(domain.ResolutionWidth, domain.ResolutionHeight, ebiten.FilterDefault)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewMonitor; err: %w", err)
	}

	return &Renderer{
		width:            domain.ResolutionWidth,
		height:           domain.ResolutionHeight,
		scale:            scale,
		title:            title,
		imageBuf:         imageBuf,
		enableDebugPrint: enableDebugPrint,
	}, nil
}

// update ...
func (m *Renderer) update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.DrawImage(m.imageBuf, nil)
	if m.enableDebugPrint {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", m.fps))
	}

	return nil
}

// Run ...
func (m *Renderer) Run() error {
	return ebiten.Run(m.update, m.height, m.width, m.scale, m.title)
}

// Render ...
func (m *Renderer) Render(s *domain.Screen) error {
	p := toPixels(s)
	m.imageBuf.ReplacePixels(p)

	m.fps = 1 / time.Since(m.lastRenderedTime).Seconds()
	m.lastRenderedTime = time.Now()

	return nil
}

// toPixels ...
func toPixels(s *domain.Screen) []byte {
	pixels := make([]byte, 4*domain.ResolutionHeight*domain.ResolutionWidth)

	if s.Images != nil {
		idx := 0
		for _, line := range s.Images {
			for _, pixel := range line {
				pixels[idx] = pixel.R
				idx++

				pixels[idx] = pixel.G
				idx++

				pixels[idx] = pixel.B
				idx++

				pixels[idx] = pixel.A
				idx++
			}
		}
		return pixels
	}

	// タイルを描画
	idx := 0
	for y := uint16(0); y < domain.ResolutionHeight; y++ {
		for x := uint16(0); x <= domain.ResolutionWidth-1; x++ {
			r, g, b, a := getPixel(s.TileImages, domain.MonitorX(x), domain.MonitorY(y))

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

	// スプライトを描画
	for _, sImage := range s.SpriteImages {
		if sImage.Y <= 0 {
			// 1行目は表示しない
			continue
		}

		baseIdx := int(sImage.Y)*domain.ResolutionWidth*4 + int(sImage.X)*4

		for offsetY := 0; offsetY < int(sImage.H); offsetY++ {
			y := int(sImage.Y) + offsetY
			if y < 0 || y > (domain.ResolutionHeight-1) {
				// スプライトの一部が画面外のときは描画しないためスキップ
				continue
			}

			for offsetX := 0; offsetX < int(sImage.W); offsetX++ {
				x := int(sImage.X) + offsetX
				if x < 0 || x > (domain.ResolutionWidth-1) {
					// スプライトの一部が画面外のときは描画しないためスキップ
					continue
				}
				if !s.DisableSpriteMask && x < 8 {
					continue
				}

				idx := baseIdx + (offsetY * domain.ResolutionWidth * 4) + (offsetX * 4)

				var r, g, b, a byte
				if sImage.EnableFlipHorizontal && sImage.EnableFlipVertical {
					// 左右上下反転
					iY := (domain.SpriteHeight - 1) - offsetY
					iX := (domain.SpriteWidth - 1) - offsetX
					r = sImage.R[iY][iX]
					g = sImage.G[iY][iX]
					b = sImage.B[iY][iX]
					a = sImage.A[iY][iX]
				} else if sImage.EnableFlipHorizontal && !sImage.EnableFlipVertical {
					// 左右反転
					iY := offsetY
					iX := (domain.SpriteWidth - 1) - offsetX
					r = sImage.R[iY][iX]
					g = sImage.G[iY][iX]
					b = sImage.B[iY][iX]
					a = sImage.A[iY][iX]
				} else if !sImage.EnableFlipHorizontal && sImage.EnableFlipVertical {
					// 上下反転
					iY := (domain.SpriteHeight - 1) - offsetY
					iX := offsetX
					r = sImage.R[iY][iX]
					g = sImage.G[iY][iX]
					b = sImage.B[iY][iX]
					a = sImage.A[iY][iX]
				} else {
					// 反転なし
					iY := offsetY
					iX := offsetX
					r = sImage.R[iY][iX]
					g = sImage.G[iY][iX]
					b = sImage.B[iY][iX]
					a = sImage.A[iY][iX]
				}

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
