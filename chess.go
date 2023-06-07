package main

import (
	"fmt"
	. "gfx2"
	"pieces"
)

func main() {
	var w_x, w_y uint16 = 800, 800
	fmt.Println("starting game...")
	Fenster(w_x, w_y)
	Fenstertitel("Chess")
	Stiftfarbe(0, 255, 0)
	Vollrechteck(0, 0, w_x, w_y)
	draw_board(w_x, w_y)

	pieces := make([]pieces.Piece, 8)
	ROkk := &pieces.Knight{}
	pieces[0] = &pieces.Rook{}

	for { //gameloop
		button, status, m_x, m_y := MausLesen1()
		fmt.Println("button:", button, "status", status)

		if status == 1 && button == 1 {
			fmt.Println(calc_field(w_x, w_y, m_x, m_y))
		}

	}
	// TastaturLesen1()
}

func draw_board(w_x, w_y uint16) {
	var a uint16
	if w_x < w_y {
		a = w_x / 8
	} else {
		a = w_y / 8
	}
	var f_x uint16 = 0
	var f_y uint16 = 0
	for i := 0; i <= 7; i++ {
		for k := 0; k <= 7; k++ {
			if k%2 == 0 {
				if i%2 == 0 {
					Stiftfarbe(255, 255, 255)
				} else {
					Stiftfarbe(0, 0, 0)
				}

			} else {
				if i%2 == 1 {
					Stiftfarbe(255, 255, 255)
				} else {
					Stiftfarbe(0, 0, 0)
				}
			}

			Vollrechteck(f_x, f_y, a, a)
			f_x = f_x + a
		}
		f_x = 0
		f_y = f_y + a
	}
}

func calc_field(w_x, w_y, m_x, m_y uint16) (x, y uint16) {
	var a uint16
	if w_x < w_y {
		a = w_x / 8
	} else {
		a = w_y / 8
	}

	var field_x uint16
	var field_y uint16

	field_x = m_x / a
	field_y = m_y / a
	return field_x, field_y
}
