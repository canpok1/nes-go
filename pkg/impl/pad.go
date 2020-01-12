package impl

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"github.com/hajimehoshi/ebiten"
)

// Pad ...
type Pad struct {
	mapping map[domain.ButtonType]ebiten.Key
}

// NewPad ...
func NewPad(mapping map[domain.ButtonType]ebiten.Key) domain.Pad {
	return &Pad{
		mapping: mapping,
	}
}

// IsPressed ...
func (p *Pad) IsPressed(nesKey domain.ButtonType) (pressed bool) {
	defer func() {
		if pressed {
			log.Debug("Pad.IsPressed[%v] => %v", nesKey, pressed)
		} else {
			log.Debug("Pad.IsPressed[%v] => %v", nesKey, pressed)
		}
	}()

	pcKey, ok := p.mapping[nesKey]
	if !ok {
		pressed = false
	} else {
		pressed = ebiten.IsKeyPressed(pcKey)
	}

	return
}
