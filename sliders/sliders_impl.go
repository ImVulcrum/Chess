package sliders

import (
	gfx "gfxw"
	"math"
	"strconv"
)

type Slider struct {
	x             uint16  //(upper left corner)
	y             uint16  //(upper left corner)
	x_box_cord    uint16  //(upper left corner of box)
	y_box_cord    uint16  //(upper left corner of box)
	length        uint16  //(in pixels)
	height        uint16  //(in pixels)
	thickness     uint16  //(in pixels)
	max_value     float32 //(in numbers)
	default_value float32 //(in numbers)
	value         float32 //(calc: (y_box - y) / lenght * max_value) (in number)
	name          string  //(Name of the Slider as string)
	display_int   bool    //controls if the displayed number should be displayed as an integer
}

func New() *Slider {
	var s *Slider = new(Slider)
	return s
}

func (s *Slider) Draw(x uint16, y uint16, length uint16, height uint16, thickness uint16, max_value float32, default_value float32, name string, use_int bool) *Slider {
	s.value = default_value
	s.x = x
	s.y = y
	s.length = length
	s.height = height
	s.thickness = thickness
	s.max_value = max_value
	s.default_value = default_value
	s.x_box_cord = uint16(math.Round(float64(s.value*float32(s.length)/s.max_value + float32(s.x))))
	s.y_box_cord = s.y
	s.name = name
	s.display_int = use_int
	gfx.Stiftfarbe(88, 88, 88)
	gfx.Vollrechteck(s.x, s.y, s.length+s.thickness, s.height)
	if !s.display_int {
		gfx.SchreibeFont(s.x+s.length+s.thickness+20, s.y, s.name+": "+strconv.FormatFloat(float64(s.value), 'f', -1, 32))
	} else {
		gfx.SchreibeFont(s.x+s.length+s.thickness+20, s.y, s.name+": "+strconv.Itoa(int(math.Round(float64(s.value)))))
	}
	gfx.Stiftfarbe(195, 195, 195)
	gfx.Vollrechteck(s.x_box_cord, s.y_box_cord, s.thickness, s.height)

	return s
}

func (s *Slider) Redraw(ms_x uint16) *Slider {
	if ms_x > s.length+s.x {
		ms_x = s.length + s.x
	} else if ms_x < s.x {
		ms_x = s.x
	}
	gfx.UpdateAus()
	gfx.Stiftfarbe(0, 0, 0)
	if !s.display_int {
		gfx.SchreibeFont(s.x+s.length+s.thickness+20, s.y, s.name+": "+strconv.FormatFloat(float64(s.value), 'f', -1, 32))
	} else {
		gfx.SchreibeFont(s.x+s.length+s.thickness+20, s.y, s.name+": "+strconv.Itoa(int(math.Round(float64(s.value)))))
	}
	s.x_box_cord = ms_x
	s.value = (float32(s.x_box_cord)*s.max_value - s.max_value*float32(s.x)) / float32(s.length)
	gfx.Stiftfarbe(88, 88, 88)
	gfx.Vollrechteck(s.x, s.y, s.length+s.thickness, s.height)
	if !s.display_int {
		gfx.SchreibeFont(s.x+s.length+s.thickness+20, s.y, s.name+": "+strconv.FormatFloat(float64(s.value), 'f', -1, 32))
	} else {
		gfx.SchreibeFont(s.x+s.length+s.thickness+20, s.y, s.name+": "+strconv.Itoa(int(math.Round(float64(s.value)))))
	}
	// gfx.SchreibeFont(s.x + s.length + s.thickness + 20, s.y, s.Name + ": " + strconv.FormatFloat(float64(s.value), 'f', -1, 32))
	gfx.Stiftfarbe(195, 195, 195)
	gfx.Vollrechteck(s.x_box_cord, s.y_box_cord, s.thickness, s.height)
	gfx.UpdateAn()
	return s
}

// x_box = value * length / max_value + x
// value = x_box * max_value / length - x
