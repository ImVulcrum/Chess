package buttons

import (
	"time"

	gfx "../gfxw"
)

type Button struct {
	X            uint16
	Y            uint16
	Length       uint16
	Height       uint16
	Name         string
	R            uint8
	G            uint8
	B            uint8
	R_Label      uint8
	G_Label      uint8
	B_Label      uint8
	Label_Offset uint16
	font_size    int
}

func New(x uint16, y uint16, length uint16, height uint16, name string, r, g, b, r_label, g_label, b_label uint8, label_offset uint16, font_size int) *Button {
	var button *Button = new(Button)
	(*button).X = x
	(*button).Y = y
	(*button).Length = length
	(*button).Height = height
	(*button).Name = name
	(*button).R = r
	(*button).G = g
	(*button).B = b
	(*button).R_Label = r_label
	(*button).G_Label = g_label
	(*button).B_Label = b_label
	(*button).Label_Offset = label_offset
	(*button).font_size = font_size
	return button
}

func (b *Button) Draw() {
	gfx.SetzeFont("./resources/fonts/firamono.ttf", b.font_size)
	gfx.Stiftfarbe(b.R, b.G, b.B)
	gfx.Vollrechteck((*b).X, (*b).Y, (*b).Length, (*b).Height)
	gfx.Stiftfarbe(b.R_Label, b.G_Label, b.B_Label)
	gfx.SchreibeFont((*b).X+(*b).Label_Offset, (*b).Y+(*b).Height/10, (*b).Name)
}

func (b *Button) Is_Clicked(x, y uint16) bool {
	if x >= b.X && x <= b.X+b.Length && y >= b.Y && y <= b.Y+b.Height {
		gfx.SetzeFont("./resources/fonts/firamono.ttf", b.font_size)
		gfx.Stiftfarbe(0, 0, 0)
		gfx.Transparenz(120)
		gfx.Vollrechteck((*b).X, (*b).Y, (*b).Length, (*b).Height)
		time.Sleep(time.Duration(100) * time.Millisecond)
		gfx.Stiftfarbe(b.R, b.G, b.B)
		gfx.Transparenz(0)
		gfx.Vollrechteck((*b).X, (*b).Y, (*b).Length, (*b).Height)
		gfx.Stiftfarbe(b.R_Label, b.G_Label, b.B_Label)
		gfx.SchreibeFont((*b).X+(*b).Label_Offset, (*b).Y+(*b).Height/10, (*b).Name)
		return true
	}
	return false
}
