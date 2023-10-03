package textbox

type Box interface {
	Draw()
	Is_Clicked(x, y uint16) bool
	Write()
	Get_Text() string
	If_Clicked_Write(x, y uint16)
}
