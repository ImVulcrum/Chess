package textbox

import (
	gfx "../gfxw"
)

type tbox struct {
	x          uint16
	y          uint16
	h          uint16
	w          uint16
	bg_color   [3]uint8 //background color
	fg_color   [3]uint8 //foreground color (text)
	hl_color   [3]uint8 //highlighting color (when typing or if nothing was entered)
	max_entry  int
	text       string
	enter_text string
	used       bool
}

func New(x uint16, y uint16, h uint16, w uint16, bg [3]uint8, fg [3]uint8, hl [3]uint8, max int, enter_text string) *tbox {
	var t *tbox = new(tbox)
	t.x = x
	t.y = y
	t.h = h
	t.w = w
	t.bg_color = bg
	t.fg_color = fg
	t.hl_color = hl
	t.max_entry = max
	t.enter_text = enter_text
	t.used = false

	if len(t.enter_text) > t.max_entry { //check if the the default enter text eceeds the length of the slider
		panic("the original input is longer than it should be")
	}
	t.text = t.enter_text

	return t
}

func (t *tbox) Draw() {
	t.draw(false)
}

func (t *tbox) Was_Used() bool {
	return t.used
}

func (t *tbox) draw(highlight bool) {
	if !highlight {
		gfx.Stiftfarbe(t.bg_color[0], t.bg_color[1], t.bg_color[2])
	} else {
		gfx.Stiftfarbe(t.hl_color[0], t.hl_color[1], t.hl_color[2]) //set color to highlight color
	}

	gfx.Vollrechteck(t.x, t.y, t.w, t.h) //draw the slider box in the desired color

	if !t.used { //if the box wasn't used, the enter text is supposed to appear in the highlight color
		gfx.Stiftfarbe(t.hl_color[0], t.hl_color[1], t.hl_color[2])
	} else {
		gfx.Stiftfarbe(t.fg_color[0], t.fg_color[1], t.fg_color[2])
	}

	gfx.SetzeFont("./resources/fonts/unispace.ttf", int(t.h/11*10))
	gfx.SchreibeFont(t.x+t.h/10, t.y+t.h/15, t.text)
}

func (t *tbox) Is_Clicked(m_x, m_y uint16) bool { //returns true if a click is executed directly on the slider
	if m_x >= t.x && m_x <= t.x+t.w && m_y >= t.y && m_y <= t.y+t.h {
		return true
	}
	return false
}

func (t *tbox) If_Clicked_Write(m_x, m_y uint16) { //directly move on to the write cycle if the slider was clicked
	if t.Is_Clicked(m_x, m_y) {
		t.Write()
	}
}

func (t *tbox) Write() {
	t.draw(true) //draw because the bg color is supposed to change to the highlight color to indicate that the user is in the write cycle of this box

	if !t.used { //if the box wasn't used it is supposed to clear the enter text so that the user can type
		t.text = ""
	}

	for {
		key, pressed, depth := gfx.TastaturLesen1()

		if pressed == 1 {
			var text_before string = t.text

			if key >= 97 && key <= 122 { //character
				if key == 122 { //because the gfx package uses an english querty layout, z and y have to be flipped
					key = 121
				} else if key == 121 {
					key = 122
				}

				if depth == 1 {
					key = key - 32
				}
				t.text = t.text + string(rune(int(key)))
				t.used = true

			} else if (key >= 48 && key <= 57) || key == 32 || key == 46 { //number
				t.text = t.text + string(rune(int(key)))
				t.used = true

			} else if key == 8 { //backspace
				if len(t.text) != 0 {
					t.text = t.text[:len(t.text)-1]
				}
			} else if key == 27 || key == 13 { //escape or enter --> quit
				if t.text == "" {
					t.used = false
				}
				if !t.used { //if the text box wasn't used it is supposed to display the enter message again
					t.text = t.enter_text
				}
				t.draw(false)
				break
			}

			if len(t.text) > t.max_entry { //shorten the text if it's to long
				t.text = t.text[:t.max_entry]
			}

			if t.text != text_before { //something has been written, so the box must be drawed
				t.draw(true)
			}
		}
	}
}

func (t *tbox) Get_Text() string {
	if !t.used {
		return ""
	} else {
		return t.text
	}
}
