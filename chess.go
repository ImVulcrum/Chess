package main

import (
	"chess"
	. "gfxw"
)

func main() {
	var w_x, w_y uint16 = 904, 904
	var a uint16 = calc_a(w_x, w_y)
	var white_is_current_player bool = true

	initialize(w_x, w_y, a)

	var pieces [16]chess.Piece

	pieces[0] = &chess.Rook{
		ChessObject: chess.ChessObject{
			Positioning: chess.Positioning{
				Position: [2]uint16{5, 4},
			},
			White: true,
		},
	}

	pieces[1] = &chess.Bishop{
		ChessObject: chess.ChessObject{
			Positioning: chess.Positioning{
				Position: [2]uint16{5, 5},
			},
			White: false,
		},
	}

	pieces[2] = &chess.Pawn{
		ChessObject: chess.ChessObject{
			Positioning: chess.Positioning{
				Position: [2]uint16{1, 6},
			},
			White: false,
		},
	}

	// fmt.Println(pieces[0].Give_Pos())

	// pieces[0].Move_To([2]uint16{3, 4})

	// fmt.Println(pieces[0].Give_Pos())

	// legal := pieces[1].Give_Legal_Moves()
	// fmt.Println(legal[0].Position[1])

	draw_pieces(pieces, w_x, w_y, a)

	for { //gameloop
		button, status, m_x, m_y := MausLesen1()

		if status == 1 && button == 1 {
			var current_field [2]uint16 = calc_field(w_x, w_y, m_x, m_y)
			var current_piece chess.Piece

			for i := 0; i < len(pieces); i++ {
				if pieces[i] != nil {
					if current_field == pieces[i].Give_Pos() {
						current_piece = pieces[i]
					}
				}
			}

			if current_piece != nil && current_piece.Is_White_Piece() == white_is_current_player { //wenn die maus ein piece angeklickt hat, welches dem aktuellen spieler gehört
				current_piece.Calc_Moves(pieces)
				var current_legal_moves [][2]uint16 = current_piece.Give_Legal_Moves()
				for k := 0; k < len(current_legal_moves); k++ {
					highlight(a, current_legal_moves[k])
				}
				// fmt.Println(current_piece.Give_Pos())
			}
		}
	}
}

func initialize(w_x, w_y, a uint16) {
	Fenster(w_x, w_y)
	Fenstertitel("Chess")
	Stiftfarbe(0, 255, 0)
	Vollrechteck(0, 0, w_x, w_y)
	draw_board(a)
}

func calc_a(w_x, w_y uint16) uint16 {
	var a uint16
	if w_x < w_y {
		a = w_x / 8
	} else {
		a = w_y / 8
	}
	return a
}

func draw_pieces(pieces [16]chess.Piece, w_x, w_y, a uint16) {
	for i := 0; i < len(pieces); i++ {
		if pieces[i] != nil {
			chess.Draw(pieces[i], w_x, w_y, a)
		}
	}
}

func highlight(a uint16, pos [2]uint16) {

	var cord_x uint16 = a * pos[0]
	var cord_y uint16 = a * pos[1]

	Transparenz(170)
	Stiftfarbe(0, 255, 0)
	Vollrechteck(cord_x, cord_y, a, a)
	Transparenz(0)
}

func draw_board(a uint16) {
	var f_x uint16 = 0
	var f_y uint16 = 0
	for i := 0; i <= 7; i++ {
		for k := 0; k <= 7; k++ {
			if k%2 == 0 {
				if i%2 == 0 {
					Stiftfarbe(240, 217, 181)
				} else {
					Stiftfarbe(181, 136, 99)
				}

			} else {
				if i%2 == 1 {
					Stiftfarbe(240, 217, 181)
				} else {
					Stiftfarbe(181, 136, 99)
				}
			}

			Vollrechteck(f_x, f_y, a, a)
			f_x = f_x + a
		}
		f_x = 0
		f_y = f_y + a
	}
}

func calc_field(w_x, w_y, m_x, m_y uint16) [2]uint16 {
	var a uint16
	if w_x < w_y {
		a = w_x / 8
	} else {
		a = w_y / 8
	}

	var current_field [2]uint16

	current_field[0] = m_x / a
	current_field[1] = m_y / a
	return current_field
}
