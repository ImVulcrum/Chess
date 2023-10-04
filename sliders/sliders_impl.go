package sliders

import (
	"math"
	"strconv"

	gfx "../gfxw"
)

type slid struct {
	x               uint16 //(upper left corner)
	y               uint16 //(upper left corner)
	x_box_cord      uint16 //(upper left corner of box)
	y_box_cord      uint16 //(upper left corner of box)
	length          uint16 //(in pixels)
	height          uint16 //(in pixels)
	thickness       uint16 //(in pixels)
	min_value       float32
	max_value       float32 //(in numbers)
	default_value   float32 //(in numbers)
	value           float32 //(calc: (y_box - y) / lenght * max_value) (in number)
	name            string  //(Name of the Slider as string)
	display_int     bool    //controls if the displayed number should be displayed as an integer
	bg_color_slider [3]uint8
	fg_color_slider [3]uint8
	bg_color_window [3]uint8
	active          bool
}

func New(x uint16, y uint16, length uint16, height uint16, thickness uint16, min_value float32, max_value float32, default_value float32, name string, use_int bool, bg_color_s [3]uint8, fg_color_s [3]uint8, bg_color_w [3]uint8) *slid {
	var s *slid = new(slid)
	s.default_value = default_value
	s.value = s.default_value
	s.x = x
	s.y = y
	s.length = length
	s.height = height
	s.thickness = thickness
	s.min_value = min_value
	s.max_value = max_value - s.min_value

	if s.min_value > s.default_value {
		panic("The default value of slider: '" + name + " 'is higher than allowed")
	}

	//value times lenght divided by the max value equals the relative postion of the slider indicator, add the postion of the bar to get the real position
	s.x_box_cord = uint16(math.Round(float64((s.value-s.min_value)*float32(s.length)/s.max_value + float32(s.x))))
	s.y_box_cord = s.y
	s.name = name
	s.display_int = use_int
	s.bg_color_slider = bg_color_s
	s.fg_color_slider = fg_color_s
	s.bg_color_window = bg_color_w
	s.active = true
	return s
}

func (s *slid) Draw() {
	s.draw(false)
}

func (s *slid) Get_Value() float32 {
	if s.display_int {
		return float32(math.Round(float64(s.value)))
	} else {
		return s.value
	}
}

func (s *slid) draw(delete bool) {
	gfx.SetzeFont("./resources/fonts/unispace.ttf", int(s.height))

	if !delete {
		//draw the slider
		gfx.Stiftfarbe(s.bg_color_slider[0], s.bg_color_slider[1], s.bg_color_slider[2])
		gfx.Vollrechteck(s.x, s.y, s.length+s.thickness, s.height)
		//draw the slider indicator
		gfx.Stiftfarbe(s.fg_color_slider[0], s.fg_color_slider[1], s.fg_color_slider[2])
		gfx.Vollrechteck(s.x_box_cord, s.y_box_cord, s.thickness, s.height)
	} else {
		gfx.Stiftfarbe(s.bg_color_window[0], s.bg_color_window[1], s.bg_color_window[2])
	}
	if !s.display_int { //if the value should be diplayed as a float, the float must be rounded
		gfx.SchreibeFont(s.x+s.length+3*s.thickness, s.y-s.height/7, s.name+": "+strconv.FormatFloat(float64(s.value), 'f', -1, 32))
	} else { //otherwise the float is just concerted to an integer
		gfx.SchreibeFont(s.x+s.length+3*s.thickness, s.y-s.height/7, s.name+": "+strconv.Itoa(int(math.Round(float64(s.value)))))
	}
}

func (s *slid) Is_Clicked(m_x, m_y uint16) bool {
	if m_x >= s.x && m_x <= s.x+s.length+s.thickness && m_y >= s.y && m_y <= s.y+s.height {
		return true
	}
	return false
}

func (s *slid) Activate() {
	s.active = true
	s.Draw()
}

func (s *slid) Deactivate() {
	s.active = false
	gfx.Stiftfarbe(s.bg_color_window[0], s.bg_color_window[1], s.bg_color_window[2])
	gfx.Vollrechteck(s.x, s.y, s.length+s.thickness, s.height)
	s.draw(true)
}

func (s *slid) If_Clicked_Draw(m_x, m_y uint16) {
	if s.active && s.Is_Clicked(m_x, m_y) {
		s.Redraw(m_x)
		for {
			button, status, m_x, _ := gfx.MausLesen1()
			if button == 1 && status == 0 {
				s.Redraw(m_x)
			} else {
				break
			}
		}
	}
}

func (s *slid) Redraw(x uint16) {
	gfx.UpdateAus()

	if x > s.x+s.length { //if the cord of the mouse is greater than the end of the slider, it's setted to the end
		x = s.x + s.length
	} else if x < s.x { //if the cord of the mouse is smaller than the beginning of the slider, it's setted to the beginning
		x = s.x
	}

	//überschreiben der ursprünglichen font
	s.draw(true)

	s.x_box_cord = x
	s.value = (float32(s.x_box_cord)*s.max_value-s.max_value*float32(s.x))/float32(s.length) + s.min_value //calc the value

	s.draw(false)

	gfx.UpdateAn()
}
