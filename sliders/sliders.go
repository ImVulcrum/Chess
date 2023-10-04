package sliders

type Slider interface {
	Draw()
	Redraw(x uint16)
	Is_Clicked(m_x, m_y uint16) bool
	If_Clicked_Draw(m_x, m_y uint16)
	Get_Value() float32
	Activate()
	Deactivate()
}
