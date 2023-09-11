package main

import (
	"fmt"
	"image"
	"os"
	"time"

	gfx "./gfxw"
	"./imaging"
	"./parser"
	"./pieces"
)

func main() {
	fmt.Println("start game")
	var premoves string = ""
	// `
	// 	[Event "Open NOR-ch"]
	// 	[Site "Oslo NOR"]
	// 	[Date "2001.04.08"]
	// 	[Round "3"]
	// 	[White "Flores,R"]
	// 	[Black "Carlsen,M"]
	// 	[Result "0-1"]
	// 	[WhiteElo ""]
	// 	[BlackElo "2064"]
	// 	[ECO "B76"]

	// 	1.e4 c5 2.Nf3 d6 3.d4 Nf6 4.Nc3 cxd4 5.Nxd4 g6 6.f3 Bg7 7.Be3 O-O 8.Qd2 Nc6
	// 	9.Nb3 Be6 10.Bh6 a5 11.Bxg7 Kxg7 12.g4 Ne5 13.Be2 Nc4 14.Bxc4 Bxc4 15.h4 a4
	// 	16.Nd4 e5 17.Ndb5 d5 18.g5 Nh5 19.exd5 Nf4 20.O-O-O Ra5 21.Na3 Bxd5 22.Nxd5 Rxd5
	// 	23.Qe3 Rxd1+ 24.Rxd1 Qc7 25.Qe4 Qc5 26.Qxb7 Ne2+ 27.Kd2 Qf2 28.Qc7 e4 29.fxe4 Re8
	// 	30.e5 Qd4+ 31.Ke1 Qe4 32.Kd2 Rxe5 33.c4 Qf4+ 34.Kc2 Nd4+ 35.Rxd4 Qxd4 36.Qxe5+ Qxe5
	// 	37.b4 Qe2+  0-1
	// 	`
	var w_x, w_y uint16 = 800, 800
	var duration_of_premove_animation int = 0

	premoves_array := parser.Create_Array_Of_Moves(premoves)
	var a uint16 = calc_a(w_x, w_y)
	var white_is_current_player bool
	var player_change bool = true
	var current_king_index int

	pieces_a, white_king_index, black_king_index := initialize(w_x, w_y, a, false)

	var checkmate bool
	var check bool
	var current_piece pieces.Piece
	var temp_current_piece pieces.Piece
	var piece_index int
	var current_legal_moves [][3]uint16
	var moves_counter int16
	var current_field [2]uint16
	var there_are_no_premoves bool = false
	var restart bool

	var ending_premoves bool = true

	var piece_is_selected uint16 = 64

	draw_pieces(pieces_a, w_x, w_y, a)

	var promotion uint16 = 0
	for { //gameloop

		if player_change {

			restart = false
			player_change = false
			white_is_current_player, current_king_index = change_player(white_is_current_player, white_king_index, black_king_index)
			moves_counter++
			pieces_a, checkmate = pieces.Calc_Moves_With_Check(pieces_a, moves_counter, current_king_index)

			check = pieces_a[current_king_index].(*pieces.King).Is_In_Check(pieces_a, moves_counter)

			if there_are_no_premoves || duration_of_premove_animation != 0 {
				draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
			}

			//fmt.Println("---")
			if checkmate && check {
				fmt.Println("Checkmate")
				game_end_visual(0, a, white_is_current_player)
				gfx.TastaturLesen1()
				pieces_a, white_king_index, black_king_index, moves_counter, check, white_is_current_player, restart, player_change = restart_game(w_x, w_y, a)
			} else if checkmate {
				fmt.Println("Stalemate")
				game_end_visual(1, a, white_is_current_player)
				gfx.TastaturLesen1()
				pieces_a, white_king_index, black_king_index, moves_counter, check, white_is_current_player, restart, player_change = restart_game(w_x, w_y, a)
			}

		}

		if len(premoves_array) > 0 {
			var promotion uint16 = 0
			fmt.Println("premove:", premoves_array[0])
			piece_executing_move, index_of_move, piece_promoting_to := parser.Get_Correct_Move(premoves_array[0], pieces_a, current_king_index)

			pieces_a, promotion = pieces.Move_Piece_To(pieces_a[piece_executing_move], pieces_a[piece_executing_move].Give_Legal_Moves()[index_of_move], moves_counter, pieces_a)

			if promotion != 64 {
				pieces_a = pawn_promotion(w_x, w_y, a, piece_executing_move, pieces_a, piece_promoting_to)
			}
			premoves_array = premoves_array[1:]

			time.Sleep(time.Duration(duration_of_premove_animation) * time.Millisecond)
			player_change = true
		} else if ending_premoves {
			draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
			there_are_no_premoves = true
			ending_premoves = false
			fmt.Println("----- End of Premoves -----")
		}

		if there_are_no_premoves && !restart {
			button, status, m_x, m_y := gfx.MausLesen1()

			if status == 1 && button == 1 {
				fmt.Println("-------------Entering Click--------------")
				current_field = calc_field(a, m_x, m_y, 0)

				temp_current_piece = nil
				for piece_index = 0; piece_index < len(pieces_a); piece_index++ {
					if pieces_a[piece_index] != nil {
						if current_field == pieces_a[piece_index].Give_Pos() {
							temp_current_piece = pieces_a[piece_index]
							break
						}
					}
				}

				if piece_is_selected != 64 {
					fmt.Println("Click: piece is selected")
					//überprüfen ob auf ein Feld in Legal Moves geklickt wurde
					pieces_a, piece_is_selected, player_change, promotion = move_if_current_field_is_in_legal_moves(current_field, pieces_a, promotion, piece_is_selected, a, w_x, w_y, current_king_index, check, moves_counter)
				}

				//auswählen eines pieces
				if temp_current_piece != nil && temp_current_piece.Is_White_Piece() == white_is_current_player { //select
					fmt.Println("Click: selected piece")
					current_piece = temp_current_piece
					current_legal_moves = current_piece.Give_Legal_Moves()
					promotion = 0
					draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, true, current_king_index, check)
					piece_is_selected = uint16(piece_index)
				} else if piece_is_selected != 64 && ((temp_current_piece == nil) || (temp_current_piece != nil && (temp_current_piece.Give_Pos() == current_piece.Give_Pos() || temp_current_piece.Is_White_Piece() != white_is_current_player))) { //deselect
					//sobald ein Piece ausgewählt ist: wenn auf kein Piece geklickt wurde, auf das bereits ausgewählte Piece nocheinmal geklickt wurde oder auf ein gegnerisches Piece geklickt wurde, wird das ausgewählte Piece deselected
					fmt.Println("Click: deselect after click on same piece, the field, or an enemy piece")
					current_piece = nil
					draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
					piece_is_selected = 64
				}
			}
			if status == -1 && button == 1 {
				fmt.Println("Entering Release")
				current_field = calc_field(a, m_x, m_y, 0)

				temp_current_piece = nil
				for piece_index = 0; piece_index < len(pieces_a); piece_index++ {
					if pieces_a[piece_index] != nil {
						if current_field == pieces_a[piece_index].Give_Pos() {
							temp_current_piece = pieces_a[piece_index]
							break
						}
					}
				}

				if piece_is_selected != 64 {
					//überprüfen ob in current legal moves
					fmt.Println("Release: piece is selected")
					pieces_a, piece_is_selected, player_change, promotion = move_if_current_field_is_in_legal_moves(current_field, pieces_a, promotion, piece_is_selected, a, w_x, w_y, current_king_index, check, moves_counter)
				}

				//Sobald ein Piece ausgewählt ist: wenn auf kein Piece geklickt wurde oder auf ein generisches Piece geklickt wurde
				if piece_is_selected != 64 && ((temp_current_piece == nil) || (temp_current_piece != nil && temp_current_piece.Is_White_Piece() != white_is_current_player)) { //deselect
					fmt.Println("Release: deselect after release on the field or enemy piece")
					current_piece = nil
					draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
					piece_is_selected = 64

				}
			}
		}
	}
}

func move_if_current_field_is_in_legal_moves(current_field [2]uint16, pieces_a [64]pieces.Piece, promotion uint16, piece_is_selected uint16, a, w_x, w_y uint16, current_king_index int, check bool, moves_counter int16) ([64]pieces.Piece, uint16, bool, uint16) {
	current_piece := pieces_a[piece_is_selected]
	current_legal_moves := current_piece.Give_Legal_Moves()
	for k := 0; k < len(current_legal_moves); k++ {
		if current_field == [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]} { //wenn das der Fall ist, wird das Piece bewegt
			fmt.Println("Click: moving piece because it was released on a legal move field")
			pieces_a, promotion = pieces.Move_Piece_To(current_piece, current_legal_moves[k], moves_counter, pieces_a)
			if promotion != 64 {
				draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
				pieces_a = pawn_promotion(w_x, w_y, a, int(piece_is_selected), pieces_a, "A")
			}
			piece_is_selected = 64
			promotion = 0
			return pieces_a, 64, true, 0
		}
	}
	return pieces_a, piece_is_selected, false, promotion
}

//hfdsoghfdoljghofdhg
// if status == 1 && button == 1 {
// 	current_field = calc_field(a, m_x, m_y, 0)

// 	for piece_index = 0; piece_index < len(pieces_a); piece_index++ {
// 		if pieces_a[piece_index] != nil {
// 			if current_field == pieces_a[piece_index].Give_Pos() {
// 				current_piece = pieces_a[piece_index]
// 				break
// 			}
// 		}
// 	}

// 	if current_piece != nil && current_piece.Is_White_Piece() == white_is_current_player { //wenn die maus ein piece angeklickt hat, welches dem aktuellen spieler gehört
// 		//current_piece.Calc_Moves(pieces_a, moves_counter)
// 		current_legal_moves = current_piece.Give_Legal_Moves()

// 		var x_offset int16 = int16(current_piece.Give_Pos()[0]*a) - int16(m_x)
// 		var y_offset int16 = int16(current_piece.Give_Pos()[1]*a) - int16(m_y)
// 		var promotion uint16 = 0

// 		draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, true, current_king_index, check)

// 		for {
// 			button, status, m_x, m_y := gfx.MausLesen1() //hält so lange an, bis die maus bewegt wurde
// 			if status != -1 && button == 1 {
// 				//schwebenedes piece wenn taste gehalten wird
// 				pieces.Draw_To_Point(current_piece, w_x, w_y, a, m_x, m_y, x_offset, y_offset, 50)

// 			} else { //wenn taste losgelassen wird
// 				new_field := calc_field(a, uint16(int16(m_x)+x_offset+int16(a)/2), uint16(int16(m_y)+y_offset+int16(a)/2), 0)

// 				if new_field == current_piece.Give_Pos() { //wenn taste über dem gleichen feld losgelassen wird wie die Figur steht
// 					gfx.Restaurieren(0, 0, w_x, w_y)
// 					break
// 				}
// 				//überprüfen ob das Feld über dem die Maus losgelassen wurde in den Legal Moves des angeklickten Pieces enthalten ist
// 				for k := 0; k < len(current_legal_moves); k++ {
// 					if new_field == [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]} { //wenn das der Fall ist, wird das Piece bewegt

// 						pieces_a, promotion = pieces.Move_Piece_To(current_piece, current_legal_moves[k], moves_counter, pieces_a)
// 						if promotion != 64 {
// 							draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
// 							pieces_a = pawn_promotion(w_x, w_y, a, piece_index, pieces_a, "A")
// 						}

// 						player_change = true
// 						break
// 					}
// 				}
// 				//entweder wurde ein piece bewegt oder die maus wurde auf einem Feld losgelassen, welches nicht in Legal_Moves enthalten ist
// 				//in jedem Fall wird das Feld neugezeichnet
// 				draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
// 				break
// 			}
// 		}
// 	}
// }

func restart_game(w_x, w_y, a uint16) ([64]pieces.Piece, int, int, int16, bool, bool, bool, bool) {
	pieces_a, white_king_index, black_king_index := initialize(w_x, w_y, a, true)

	return pieces_a, white_king_index, black_king_index, 0, false, false, true, true
}

func game_end_visual(ending_var uint8, a uint16, white_is_current_player bool) {

	gfx.Stiftfarbe(0, 0, 0)
	gfx.Transparenz(70)
	gfx.Vollrechteck(0, 0, a*8, a*8)
	gfx.Transparenz(50)
	gfx.Stiftfarbe(8, 8, 8)
	gfx.Vollrechteck((a), 3*a, 6*a, 2*a)

	gfx.Stiftfarbe(220, 220, 220)
	gfx.SetzeFont("junegull.ttf", int(5*a/10))

	if ending_var == 0 {
		if !white_is_current_player {
			gfx.SchreibeFont(18*a/10, 31*a/10, "player black has")
		} else {
			gfx.SchreibeFont(17*a/10, 31*a/10, "player white has")
		}
		gfx.Stiftfarbe(136, 8, 8)
		gfx.SetzeFont("punk.ttf", int(a))
		gfx.SchreibeFont(29*a/10, 37*a/10, "lost")
	}

	if ending_var == 1 {
		gfx.SchreibeFont(295*a/100, 31*a/10, "That's a ")

		gfx.Stiftfarbe(0, 143, 230)
		gfx.SetzeFont("punk.ttf", int(a))
		gfx.SchreibeFont(144*a/100, 37*a/10, "stalemate")
	}

	gfx.Transparenz(0)

}

func pawn_promotion(w_x, w_y, a uint16, pawn_index int, pieces_a [64]pieces.Piece, premoved string) [64]pieces.Piece {
	var queen pieces.Piece = pieces.NewQueen(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var knight pieces.Piece = pieces.NewKnight(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var rook pieces.Piece = pieces.NewRook(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var bishop pieces.Piece = pieces.NewBishop(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())

	if premoved == "A" {
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

	} else {
		if premoved == "Q" {
			pieces_a[pawn_index] = queen
		} else if premoved == "N" {
			pieces_a[pawn_index] = knight
		} else if premoved == "R" {
			pieces_a[pawn_index] = rook
		} else if premoved == "B" {
			pieces_a[pawn_index] = bishop
		}
	}
	return pieces_a
}

func draw_board(a, w_x, w_y uint16, current_piece pieces.Piece, current_legal_moves [][3]uint16, pieces_a [64]pieces.Piece, highlighting_is_activated bool, current_king_index int, check bool) {
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

func initialize(w_x, w_y, a uint16, restart bool) ([64]pieces.Piece, int, int) {
	if !restart {
		gfx.Fenster(w_x, w_y)
		gfx.Fenstertitel("Chess")
		rescale_image(a)
	}

	gfx.Stiftfarbe(221, 221, 221)
	gfx.Vollrechteck(0, 0, w_x, w_y)
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
	file, err := os.Open("Pieces_Source_Original.bmp")
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
