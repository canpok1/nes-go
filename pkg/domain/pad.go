package domain

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
type Pad interface {
	IsPressed(ButtonType) bool
}
