package impl

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"github.com/hajimehoshi/ebiten"
)

// Pad ...
type Pad struct {
	mapping map[domain.ButtonType]ebiten.Key
	states  map[domain.ButtonType]bool
}

// NewPad ...
func NewPad(mapping map[domain.ButtonType]ebiten.Key) domain.Pad {
	return &Pad{
		mapping: mapping,
		states:  map[domain.ButtonType]bool{},
	}
}

// Load ...
func (p *Pad) Load() error {
	for _, b := range domain.ButtonList {
		p.states[b] = p.isPressed(b)
	}
	return nil
}

// IsPressed ...
func (p *Pad) IsPressed(key domain.ButtonType) (pressed bool) {
	defer func() {
		if pressed {
			log.Debug("Pad.IsPressed[%v] => %v", key, pressed)
		} else {
			log.Debug("Pad.IsPressed[%v] => %v", key, pressed)
		}
	}()

	var ok bool
	pressed, ok = p.states[key]
	if !ok {
		pressed = false
	}

	return
}

// isPressed ...
func (p *Pad) isPressed(nesKey domain.ButtonType) (pressed bool) {
	pcKey, ok := p.mapping[nesKey]
	if !ok {
		pressed = false
	} else {
		pressed = ebiten.IsKeyPressed(pcKey)
	}

	return
}
