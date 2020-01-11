package domain

import "fmt"

// NameTablePoint ...
type NameTablePoint struct {
	X uint8
	Y uint8
}

// Validate ...
func (p NameTablePoint) Validate() error {
	if p.Y < 0 || p.Y >= 30 {
		return fmt.Errorf("y is out of range; y: %v", p.Y)
	}
	return nil
}

// ToIndex ...
func (p NameTablePoint) ToIndex() uint16 {
	return uint16(p.Y)*32 + uint16(p.X)
}

// ToAttributeTableIndex ...
func (p NameTablePoint) ToAttributeTableIndex() uint16 {
	x := p.X / 4
	y := p.Y / 4
	return uint16(y)*8 + uint16(x)
}

// ToAttributeIndex ...
func (p NameTablePoint) ToAttributeIndex() uint16 {
	x := (p.X % 4) / 2
	y := (p.Y % 4) / 2
	return uint16(y)*2 + uint16(x)
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
