package main

import (
	"fmt"
	"image"
	"os"

	gfx "./gfxw"
	"./imaging"
	"./pieces"
)

func main() {
	var w_x, w_y uint16 = 1000, 1000
	var a uint16 = calc_a(w_x, w_y)
	var white_is_current_player bool
	var player_change bool = true
	fmt.Println("start game")
	var current_king_index int
	pieces_a, white_king_index, black_king_index := initialize(w_x, w_y, a)
	var checkmate bool
	var check bool
	var current_piece pieces.Piece
	var piece_index int
	var current_legal_moves [][3]uint16
	var moves_counter int16
	var current_field [2]uint16

	draw_pieces(pieces_a, w_x, w_y, a)

	for { //gameloop

		if player_change {

			player_change = false
			white_is_current_player, current_king_index = change_player(white_is_current_player, white_king_index, black_king_index)
			moves_counter++
			pieces_a, checkmate = pieces.Calc_Moves_With_Check(pieces_a, moves_counter, current_king_index)

			check = pieces_a[current_king_index].(*pieces.King).Is_In_Check(pieces_a, moves_counter)
			Draw_Board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
			fmt.Println("---")
			if checkmate {
				fmt.Println("checkmate")
				gfx.TastaturLesen1()
			}

		}

		button, status, m_x, m_y := gfx.MausLesen1()

		if status == 1 && button == 1 {
			current_field = calc_field(a, m_x, m_y, 0)

			for piece_index = 0; piece_index < len(pieces_a); piece_index++ {
				if pieces_a[piece_index] != nil {
					if current_field == pieces_a[piece_index].Give_Pos() {
						current_piece = pieces_a[piece_index]
						break
					}
				}
			}

			if current_piece != nil && current_piece.Is_White_Piece() == white_is_current_player { //wenn die maus ein piece angeklickt hat, welches dem aktuellen spieler gehört
				//current_piece.Calc_Moves(pieces_a, moves_counter)
				current_legal_moves = current_piece.Give_Legal_Moves()

				var x_offset int16 = int16(current_piece.Give_Pos()[0]*a) - int16(m_x)
				var y_offset int16 = int16(current_piece.Give_Pos()[1]*a) - int16(m_y)
				var promotion uint16 = 0

				Draw_Board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, true, current_king_index, check)

				for {
					button, status, m_x, m_y := gfx.MausLesen1() //hält so lange an, bis die maus bewegt wurde
					if status != -1 && button == 1 {
						//schwebenedes piece wenn taste gehalten wird
						pieces.Draw_To_Point(current_piece, w_x, w_y, a, m_x, m_y, x_offset, y_offset, 50)

					} else { //wenn taste losgelassen wird
						new_field := calc_field(a, uint16(int16(m_x)+x_offset+int16(a)/2), uint16(int16(m_y)+y_offset+int16(a)/2), 0)

						if new_field == current_piece.Give_Pos() { //wenn taste über dem gleichen feld losgelassen wird wie die Figur steht
							gfx.Restaurieren(0, 0, w_x, w_y)
							break
						}
						//überprüfen ob das Feld über dem die Maus losgelassen wurde in den Legal Moves des angeklickten Pieces enthalten ist
						for k := 0; k < len(current_legal_moves); k++ {
							if new_field == [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]} { //wenn das der Fall ist, wird das Piece bewegt

								pieces_a, promotion = pieces.Move_Piece_To(current_piece, current_legal_moves[k], moves_counter, pieces_a, false)
								if promotion != 64 {
									Draw_Board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
									pieces_a = Pawm_Promotion(w_x, w_y, a, piece_index, pieces_a)
								}

								player_change = true
								break
							}
						}
						//entweder wurde ein piece bewegt oder die maus wurde auf einem Feld losgelassen, welches nicht in Legal_Moves enthalten ist
						//in jedem Fall wird das Feld neugezeichnet
						Draw_Board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
						break
					}
				}
			}
		}
	}
}

func Pawm_Promotion(w_x, w_y, a uint16, pawn_index int, pieces_a [64]pieces.Piece) [64]pieces.Piece {
	var queen pieces.Piece = pieces.NewQueen(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var knight pieces.Piece = pieces.NewKnight(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var rook pieces.Piece = pieces.NewRook(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var bishop pieces.Piece = pieces.NewBishop(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())

	gfx.Stiftfarbe(0, 0, 0)
	gfx.Transparenz(70)
	gfx.Vollrechteck(0, 0, a*8, a*8)
	gfx.Transparenz(0)
	gfx.Stiftfarbe(221, 221, 221)
	gfx.Vollrechteck((2 * a), (4*a)-(5*a)/10, 4*a, a)

	pieces.Draw_To_Point(queen, w_x, w_y, a, (2 * a), (4*a)-(5*a)/10, 0, 0, 0)
	pieces.Draw_To_Point(knight, w_x, w_y, a, (3 * a), (4*a)-(5*a)/10, 0, 0, 0)
	pieces.Draw_To_Point(rook, w_x, w_y, a, (4 * a), (4*a)-(5*a)/10, 0, 0, 0)
	pieces.Draw_To_Point(bishop, w_x, w_y, a, (5 * a), (4*a)-(5*a)/10, 0, 0, 0)
	for {
		button, status, m_x, m_y := gfx.MausLesen1()
		if button == 1 && status == 1 && m_x >= 2*a && m_x <= 6*a && m_y >= (4*a)-(5*a)/10 && m_y <= (4*a)+(5*a)/10 {
			x_field_pos := (calc_field(a, m_x, m_y, (5*a)/10))[0]
			if x_field_pos == 2 {
				pieces_a[pawn_index] = queen
			} else if x_field_pos == 3 {
				pieces_a[pawn_index] = knight
			} else if x_field_pos == 4 {
				pieces_a[pawn_index] = rook
			} else if x_field_pos == 5 {
				pieces_a[pawn_index] = bishop
			}
			break
		}

	}
	return pieces_a
}

func Draw_Board(a, w_x, w_y uint16, current_piece pieces.Piece, current_legal_moves [][3]uint16, pieces_a [64]pieces.Piece, highlighting_is_activated bool, current_king_index int, check bool) {
	gfx.UpdateAus()
	draw_background(a)
	if check {
		highlight(a, pieces_a[current_king_index].Give_Pos(), 255, 0, 0)
	}
	if highlighting_is_activated {
		highlight(a, current_piece.Give_Pos(), 0, 50, 255)
		for k := 0; k < len(current_legal_moves); k++ {
			highlight(a, [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]}, 0, 255, 0)
		}
	}
	draw_pieces(pieces_a, w_x, w_y, a)
	gfx.UpdateAn()
}

func change_player(white_is_current_player bool, white_king_index, black_king_index int) (bool, int) {
	var current_king_index int
	if white_is_current_player {
		current_king_index = black_king_index
		white_is_current_player = false
	} else {
		current_king_index = white_king_index
		white_is_current_player = true
	}
	return white_is_current_player, current_king_index
}

func initialize(w_x, w_y, a uint16) ([64]pieces.Piece, int, int) {
	gfx.Fenster(w_x, w_y)
	gfx.Fenstertitel("Chess")
	gfx.Stiftfarbe(221, 221, 221)
	gfx.Vollrechteck(0, 0, w_x, w_y)

	rescale_image(a)
	draw_background(a)

	var pieces_a [64]pieces.Piece
	var white_king_index int = -1
	var black_king_index int = -1

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
	pieces_a[25] = pieces.NewKnight(1, 7, true)
	pieces_a[26] = pieces.NewBishop(2, 7, true)
	pieces_a[27] = pieces.NewQueen(3, 7, true)
	pieces_a[28] = pieces.NewKing(4, 7, true)
	pieces_a[29] = pieces.NewBishop(5, 7, true)
	pieces_a[30] = pieces.NewKnight(6, 7, true)
	pieces_a[31] = pieces.NewRook(7, 7, true)

	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if king, ok := pieces_a[i].(*pieces.King); ok {
				if king.Is_White_Piece() {
					white_king_index = i
				} else if !king.Is_White_Piece() {
					black_king_index = i
				}
			}
		}
	}
	return pieces_a, white_king_index, black_king_index
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

func draw_background(a uint16) {
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

func rescale_image(a uint16) {

	// Open the BMP file
	file, err := os.Open("Pieces_Source.bmp")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the BMP file
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	// Specify the desired width and height
	var width int = 6 * int(a)
	var height int = 2 * int(a)

	// Resize the image to the specified dimensions
	resizedImg := imaging.Resize(img, width, height, imaging.NearestNeighbor)

	// Save the resized image to a new BMP file
	err = imaging.Save(resizedImg, "Pieces.bmp")
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
}

func calc_field(a, m_x, m_y, y_offset uint16) [2]uint16 {
	var current_field [2]uint16

	current_field[0] = m_x / a
	current_field[1] = (m_y + y_offset) / a
	return current_field
}
