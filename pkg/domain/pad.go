package domain

import (
	"nes-go/pkg/log"

	"github.com/hajimehoshi/ebiten"
)

// ButtonType ...
type ButtonType string

const (
	ButtonTypeA      ButtonType = ButtonType("A")
	ButtonTypeB      ButtonType = ButtonType("B")
	ButtonTypeSelect ButtonType = ButtonType("SELECT")
	ButtonTypeStart  ButtonType = ButtonType("START")
	ButtonTypeUp     ButtonType = ButtonType("UP")
	ButtonTypeDown   ButtonType = ButtonType("DOWN")
	ButtonTypeLeft   ButtonType = ButtonType("LEFT")
	ButtonTypeRight  ButtonType = ButtonType("RIGHT")
)

// Pad ...
type Pad struct {
	mapping map[ButtonType]ebiten.Key
}

// NewPad ...
func NewPad(mapping map[ButtonType]ebiten.Key) *Pad {
	return &Pad{
		mapping: mapping,
	}
}

// IsPressed ...
func (p *Pad) IsPressed(nesKey ButtonType) (pressed bool) {
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
