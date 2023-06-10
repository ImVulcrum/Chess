package main

import (
	gfx "gfxw"

	"./pieces"
)

func main() {
	var w_x, w_y uint16 = 904, 904
	var a uint16 = calc_a(w_x, w_y)
	var white_is_current_player bool = true

	pieces_a := initialize(w_x, w_y, a)

	var moves_counter int16 = 1

	// fmt.Println(pieces[0].Give_Pos())

	// pieces[0].Move_To([2]uint16{3, 4})

	// fmt.Println(pieces[0].Give_Pos())

	// legal := pieces[1].Give_Legal_Moves()
	// fmt.Println(legal[0].Position[1])

	draw_pieces(pieces_a, w_x, w_y, a)

	for { //gameloop

		button, status, m_x, m_y := gfx.MausLesen1()

		if status == 1 && button == 1 {
			// UpdateAus()
			// draw_board(a)
			// draw_pieces(pieces_a, w_x, w_y, a)
			// UpdateAn()
			var current_field [2]uint16 = calc_field(a, m_x, m_y)
			var current_piece pieces.Piece

			for i := 0; i < len(pieces_a); i++ {
				if pieces_a[i] != nil {
					if current_field == pieces_a[i].Give_Pos() {
						current_piece = pieces_a[i]
						break
					}
				}
			}

			if current_piece != nil && current_piece.Is_White_Piece() == white_is_current_player { //wenn die maus ein piece angeklickt hat, welches dem aktuellen spieler gehört
				current_piece.Calc_Moves(pieces_a, moves_counter)
				var current_legal_moves [][3]uint16 = current_piece.Give_Legal_Moves()
				var x_offset int16 = int16(current_piece.Give_Pos()[0]*a) - int16(m_x)
				var y_offset int16 = int16(current_piece.Give_Pos()[1]*a) - int16(m_y)

				gfx.UpdateAus()
				draw_board(a)
				highlight(a, current_piece.Give_Pos(), 255, 50, 0)
				// fmt.Println(current_legal_moves)
				for k := 0; k < len(current_legal_moves); k++ {
					highlight(a, [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]}, 0, 255, 0)
				}
				draw_pieces(pieces_a, w_x, w_y, a)
				gfx.UpdateAn()

				for {
					button, status, m_x, m_y := gfx.MausLesen1()
					if status != -1 && button == 1 {
						gfx.UpdateAus()

						gfx.Restaurieren(0, 0, w_x, w_y)

						gfx.Archivieren()
						pieces.Draw_To_Mouce(current_piece, w_x, w_y, a, m_x, m_y, x_offset, y_offset)
						gfx.Restaurieren(0, 0, w_x, w_y)

						gfx.Archivieren()

						gfx.Transparenz(150)
						gfx.Clipboard_einfuegenMitColorKey(uint16(int16(m_x)+x_offset), uint16(int16(m_y)+y_offset), 5, 5, 5)
						gfx.Transparenz(0)

						gfx.UpdateAn()
					} else {
						new_field := calc_field(a, uint16(int16(m_x)+x_offset+int16(a)/2), uint16(int16(m_y)+y_offset+int16(a)/2))

						if new_field == current_piece.Give_Pos() {
							gfx.Restaurieren(0, 0, w_x, w_y)
							break
						}
						// highlight(a, new_field, 0, 0, 255)
						for k := 0; k < len(current_legal_moves); k++ {
							if new_field == [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]} { //es wurde eine Figur bewegt
								pieces_a = pieces.Move_Piece_To(current_piece, current_legal_moves[k], moves_counter, pieces_a)
								white_is_current_player = change_player(white_is_current_player)
								moves_counter++
								break
							}
						}
						gfx.UpdateAus()
						draw_board(a)
						draw_pieces(pieces_a, w_x, w_y, a)
						gfx.UpdateAn()
						break
					}
				}
			}
		}
	}
}

func change_player(white_is_current_player bool) bool {
	if white_is_current_player {
		white_is_current_player = false
	} else {
		white_is_current_player = true
	}
	return white_is_current_player
}

func initialize(w_x, w_y, a uint16) [64]pieces.Piece {
	gfx.Fenster(w_x, w_y)
	gfx.Fenstertitel("Chess")
	gfx.Stiftfarbe(0, 255, 0)
	gfx.Vollrechteck(0, 0, w_x, w_y)
	draw_board(a)

	var pieces_a [64]pieces.Piece

	pieces_a[0] = pieces.NewRook(0, 0, false)
	pieces_a[1] = pieces.NewKnight(1, 0, false)
	pieces_a[2] = pieces.NewBishop(2, 0, false)
	pieces_a[3] = pieces.NewQueen(3, 0, false)
	pieces_a[4] = pieces.NewKing(4, 0, false)
	pieces_a[5] = pieces.NewBishop(5, 0, false)
	pieces_a[6] = pieces.NewKnight(6, 0, false)
	pieces_a[7] = pieces.NewRook(7, 0, false)

	var i uint16
	for i = 0; i < 8; i++ {
		pieces_a[i+8] = pieces.NewPawn(i, 1, false)
	}
	for i = 0; i < 8; i++ {
		pieces_a[i+16] = pieces.NewPawn(i, 6, true)
	}

	pieces_a[24] = pieces.NewRook(0, 7, true)
	// pieces_a[25] = pieces.NewKnight(1, 7, true)
	// pieces_a[26] = pieces.NewBishop(2, 7, true)
	// pieces_a[27] = pieces.NewQueen(3, 7, true)
	pieces_a[28] = pieces.NewKing(4, 7, true)
	// pieces_a[29] = pieces.NewBishop(5, 7, true)
	// pieces_a[30] = pieces.NewKnight(6, 7, true)
	pieces_a[31] = pieces.NewRook(7, 7, true)

	// pieces_a[32] = pieces.NewKing(3, 3, true)
	return pieces_a
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

func draw_pieces(pieces_a [64]pieces.Piece, w_x, w_y, a uint16) {
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			gfx.Archivieren()
			pieces.Draw(pieces_a[i], w_x, w_y, a)
			gfx.Restaurieren(0, 0, w_x, w_y)

			gfx.Clipboard_einfuegenMitColorKey(pieces_a[i].Give_Pos()[0]*a, pieces_a[i].Give_Pos()[1]*a, 5, 5, 5)
		}
	}
	gfx.Archivieren()
}

func highlight(a uint16, pos [2]uint16, r, g, b uint8) {

	var cord_x uint16 = a * pos[0]
	var cord_y uint16 = a * pos[1]

	gfx.Transparenz(170)
	gfx.Stiftfarbe(r, g, b)
	gfx.Vollrechteck(cord_x, cord_y, a, a)
	gfx.Transparenz(0)
}

func draw_board(a uint16) {
	var f_x uint16 = 0
	var f_y uint16 = 0
	for i := 0; i <= 7; i++ {
		for k := 0; k <= 7; k++ {
			if k%2 == 0 {
				if i%2 == 0 {
					gfx.Stiftfarbe(240, 217, 181)
				} else {
					gfx.Stiftfarbe(181, 136, 99)
				}

			} else {
				if i%2 == 1 {
					gfx.Stiftfarbe(240, 217, 181)
				} else {
					gfx.Stiftfarbe(181, 136, 99)
				}
			}

			gfx.Vollrechteck(f_x, f_y, a, a)
			f_x = f_x + a
		}
		f_x = 0
		f_y = f_y + a
	}
}

func calc_field(a, m_x, m_y uint16) [2]uint16 {
	var current_field [2]uint16

	current_field[0] = m_x / a
	current_field[1] = m_y / a
	return current_field
}
