package impl

import (
	"fmt"
	"nes-go/pkg/domain"
)

// Monitor ...
type Monitor interface {
	Render([][]domain.TileImage) error
}

// MonitorX ...
type MonitorX uint8

// Validate ...
func (x MonitorX) Validate() error {
	if x < 0 {
		return fmt.Errorf("x is out of range; x: %v", x)
	}
	return nil
}

// MonitorY ...
type MonitorY uint8

// Validate ...
func (y MonitorY) Validate() error {
	if y < 0 || y >= 240 {
		return fmt.Errorf("y is out of range; y: %v", y)
	}
	return nil
}
