package model

// Monitor ...
type Monitor interface {
	Render([]byte) error
}
