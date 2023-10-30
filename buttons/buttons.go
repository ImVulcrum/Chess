package buttons

type Button interface {
	Draw()
	Is_Clicked(x, y uint16) bool
	Give_State() bool
	Deactivate()
	Activate()
	Is_Active() bool
	Switch(re, gr, bl uint8) bool
}
