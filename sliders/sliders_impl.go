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
	s.x_box_cord = uint16(math.Round(float64((*s).Value*float32((*s).Length)/(*s).max_value + float32((*s).X))))
	s.Y_box = (*s).Y
	s.Name = name
	s.Int = use_int
	gfx.Stiftfarbe(88, 88, 88)
	gfx.Vollrechteck((*s).X, (*s).Y, (*s).Length+(*s).thickness, (*s).Height)
	if (*s).Int == false {
		gfx.SchreibeFont((*s).X+(*s).Length+(*s).thickness+20, (*s).Y, (*s).Name+": "+strconv.FormatFloat(float64((*s).Value), 'f', -1, 32))
	} else {
		gfx.SchreibeFont((*s).X+(*s).Length+(*s).thickness+20, (*s).Y, (*s).Name+": "+strconv.Itoa(int(math.Round(float64((*s).Value)))))
	}
	gfx.Stiftfarbe(195, 195, 195)
	gfx.Vollrechteck((*s).x_box_cord, (*s).Y_box, (*s).thickness, (*s).Height)

	return s
}

func (s *Slider) Redraw(ms_x uint16) *Slider {
	if ms_x > (*s).Length+(*s).X {
		ms_x = (*s).Length + (*s).X
	} else if ms_x < (*s).X {
		ms_x = (*s).X
	}
	gfx.UpdateAus()
	gfx.Stiftfarbe(0, 0, 0)
	if (*s).Int == false {
		gfx.SchreibeFont((*s).X+(*s).Length+(*s).thickness+20, (*s).Y, (*s).Name+": "+strconv.FormatFloat(float64((*s).Value), 'f', -1, 32))
	} else {
		gfx.SchreibeFont((*s).X+(*s).Length+(*s).thickness+20, (*s).Y, (*s).Name+": "+strconv.Itoa(int(math.Round(float64((*s).Value)))))
	}
	(*s).x_box_cord = ms_x
	(*s).Value = (float32((*s).x_box_cord)*(*s).max_value - (*s).max_value*float32((*s).X)) / float32((*s).Length)
	gfx.Stiftfarbe(88, 88, 88)
	gfx.Vollrechteck((*s).X, (*s).Y, (*s).Length+(*s).thickness, (*s).Height)
	if (*s).Int == false {
		gfx.SchreibeFont((*s).X+(*s).Length+(*s).thickness+20, (*s).Y, (*s).Name+": "+strconv.FormatFloat(float64((*s).Value), 'f', -1, 32))
	} else {
		gfx.SchreibeFont((*s).X+(*s).Length+(*s).thickness+20, (*s).Y, (*s).Name+": "+strconv.Itoa(int(math.Round(float64((*s).Value)))))
	}
	// gfx.SchreibeFont((*s).X + (*s).Length + (*s).thickness + 20, (*s).Y, (*s).Name + ": " + strconv.FormatFloat(float64((*s).Value), 'f', -1, 32))
	gfx.Stiftfarbe(195, 195, 195)
	gfx.Vollrechteck((*s).x_box, (*s).Y_box, (*s).thickness, (*s).Height)
	gfx.UpdateAn()
	return s
}

// x_box = Value * Length / max_value + X
// Value = x_box * max_value / Length - X
