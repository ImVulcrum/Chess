package pieces

import (
	"fmt"
	gfx "gfxw"
)

type Piece interface {
	Calc_Moves(pieces_a [64]Piece, moves_counter int16)
	Piece_Is_White() bool
	Give_Legal_Moves() [][3]uint16
	Give_Pos() [2]uint16
	Move_To(new_position [2]uint16)
	Is_White_Piece() bool
	Append_Legal_Moves(new_legal_move [3]uint16)
	Clear_Legal_Moves()
	Set_Has_Moved(update int16)
}

type Positioning struct { //datentyp Positioning
	Position [2]uint16
}

type ChessObject struct { //datentyp ChessObject erbt vom datentyp Positioning
	Positioning
	White       bool
	Legal_Moves [][3]uint16
	Has_moved   int16
}

type Pawn struct { //alle Schachobjekte erben wiederum vom datentyp ChessObject

	ChessObject
}

type Knight struct {
	ChessObject
}

type Bishop struct {
	ChessObject
}

type Rook struct {
	ChessObject
}

type Queen struct {
	ChessObject
}

type King struct {
	ChessObject
}

func (c *ChessObject) Move_To(new_position [2]uint16) {
	c.Position = new_position
}

func can_take(p Piece, pieces_a [64]Piece, field [2]uint16) uint16 {
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if pieces_a[i].Is_White_Piece() != p.Is_White_Piece() && pieces_a[i].Give_Pos() == field {
				return uint16(i + 66)
			}
		}
	}
	return 130
}

func can_move(p Piece, pieces_a [64]Piece, field [2]uint16) bool {
	var blocking_piece bool
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if pieces_a[i].Give_Pos() == field {
				blocking_piece = true
				break
			}
		}
	}
	if !blocking_piece {
		return true
	} else {
		return false
	}
}

func (p *Pawn) can_do_enpassant(pieces_a [64]Piece, field, en_passant_pawn_pos [2]uint16, moves_counter int16) {
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if en_passant_pawn, ok := pieces_a[i].(*Pawn); ok && pieces_a[i].Is_White_Piece() != p.Is_White_Piece() && pieces_a[i].Give_Pos() == en_passant_pawn_pos {
				//andersfarbiger pawn rechts neben dem pawn --> en passant rechts
				if en_passant_pawn.Has_moved > 0 && en_passant_pawn.Has_moved+1 == moves_counter {
					p.Append_Legal_Moves([3]uint16{field[0], field[1], uint16(i)})
				}
			}
		}
	}
}

func check_if_piece_is_blocking(p Piece, pieces_a [64]Piece, current_pos [2]uint16) bool {
	var blocking_piece Piece
	var blocking_piece_index uint16
	var var_break bool = false

	for i := 0; i < len(pieces_a) && blocking_piece == nil; i++ {
		if pieces_a[i] != nil {
			if pieces_a[i].Give_Pos() == current_pos {
				blocking_piece = pieces_a[i]
				blocking_piece_index = uint16(i)
			}
		}
	}

	if blocking_piece == nil { //es steht nichts im weg
		p.Append_Legal_Moves([3]uint16{current_pos[0], current_pos[1], 64})
	} else if blocking_piece.Is_White_Piece() != p.Is_White_Piece() { //es steht etwas im weg, was aber geschlagen werden kann, daher wird danach jedoch gebreaked
		p.Append_Legal_Moves([3]uint16{current_pos[0], current_pos[1], blocking_piece_index})
		var_break = true
	} else if blocking_piece.Is_White_Piece() == p.Is_White_Piece() { //es steht etwas im weg, was aber nicht geschlagen werden kann, daher wird sofort gebreaked
		var_break = true
	} else {
		fmt.Println("fatal: Error in Calculating Moves Method")
	}
	return var_break
}

func calc_moves_diagonally(p Piece, pieces_a [64]Piece) {
	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x < 7 && new_y < 7; {
		new_x++
		new_y++
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x < 7 && new_y != 0; {
		new_x++
		new_y--
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x != 0 && new_y < 7; {
		new_x--
		new_y++
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x != 0 && new_y != 0; {
		new_x--
		new_y--
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}
}

func calc_moves_vertically_and_horizontally(p Piece, pieces_a [64]Piece) {
	for new_x := p.Give_Pos()[0]; new_x < 7; {
		new_x++
		var current_pos [2]uint16 = [2]uint16{new_x, p.Give_Pos()[1]}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x := p.Give_Pos()[0]; new_x != 0; {
		new_x--
		var current_pos [2]uint16 = [2]uint16{new_x, p.Give_Pos()[1]}
		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
		if new_x == 0 {
			break
		}

	}

	for new_y := p.Give_Pos()[1]; new_y < 7; {
		new_y++
		var current_pos [2]uint16 = [2]uint16{p.Give_Pos()[0], new_y}
		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_y := p.Give_Pos()[1]; new_y != 0; {
		new_y--
		var current_pos [2]uint16 = [2]uint16{p.Give_Pos()[0], new_y}
		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
		if new_y == 0 {
			break
		}

	}
}

func Copy_Piece_To_Clipboard(piece Piece, w_x, w_y, a uint16) {
	gfx.LadeBild(0, 0, "C:\\Users\\liamw\\Documents\\_Privat\\_Go\\Chess\\Pieces.bmp")

	switch piece.(type) {
	case *Pawn:
		if piece.Is_White_Piece() {
			gfx.Clipboard_kopieren(0, a, a, a)
		} else {
			gfx.Clipboard_kopieren(0, 0, a, a)
		}
	case *Knight:
		if piece.Is_White_Piece() {
			gfx.Clipboard_kopieren(a, a, a, a)
		} else {
			gfx.Clipboard_kopieren(a, 0, a, a)
		}
	case *Bishop:
		if piece.Is_White_Piece() {
			gfx.Clipboard_kopieren(2*a, a, a, a)
		} else {
			gfx.Clipboard_kopieren(2*a, 0, a, a)
		}
	case *Rook:
		if piece.Is_White_Piece() {
			gfx.Clipboard_kopieren(3*a, a, a, a)
		} else {
			gfx.Clipboard_kopieren(3*a, 0, a, a)
		}
	case *Queen:
		if piece.Is_White_Piece() {
			gfx.Clipboard_kopieren(4*a, a, a, a)
		} else {
			gfx.Clipboard_kopieren(4*a, 0, a, a)
		}
	case *King:
		if piece.Is_White_Piece() {
			gfx.Clipboard_kopieren(5*a, a, a, a)
		} else {
			gfx.Clipboard_kopieren(5*a, 0, a, a)
		}
	default:
		fmt.Println("Unknown Piece type")
	}
}
