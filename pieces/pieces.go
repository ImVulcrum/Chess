package pieces

import (
	"fmt"
	gfx "gfxw"
)

type Piece interface {
	Calc_Moves(pieces_a [64]Piece, moves_counter int16)
	Piece_Is_White() bool
	Give_Legal_Moves() [][2]uint16
	Give_Pos() [2]uint16
	Move_To(new_position [2]uint16)
	Is_White_Piece() bool
	Append_Legal_Moves(new_legal_move [2]uint16)
	Clear_Legal_Moves() //kann wahrscheinlich weg
}

func (c *ChessObject) Give_Legal_Moves() [][2]uint16 {
	return c.Legal_Moves
}

func (c *ChessObject) Clear_Legal_Moves() {
	c.Legal_Moves = nil
}

func (c *ChessObject) Append_Legal_Moves(new_legal_move [2]uint16) {
	c.Legal_Moves = append(c.Legal_Moves, new_legal_move)
}

func (c *ChessObject) Is_White_Piece() bool {
	return c.White
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

func Draw(piece Piece, w_x, w_y, a uint16) {
	Copy_Piece_To_Clipboard(piece, w_x, w_y, a)
}

func Draw_To_Mouce(piece Piece, w_x, w_y, a, m_x, m_y uint16, x_offset, y_offset int16) {
	Copy_Piece_To_Clipboard(piece, w_x, w_y, a)
}

func Move_Piece_To(piece Piece, new_position [2]uint16, moves_counter int16, pieces_a [64]Piece) ([64]Piece, bool) {
	var smth_has_been_taken bool = false

	for k := 0; k < len(pieces_a); k++ {
		if pieces_a[k] != nil && pieces_a[k].Give_Pos() == new_position {
			if pieces_a[k].Is_White_Piece() != piece.Is_White_Piece() {
				pieces_a[k] = nil
				smth_has_been_taken = true
			} else {
				fmt.Println("panic: Error has occured, trying to take a piece of same color")
			}
			break
		}
	}
	if pawn, ok := piece.(*Pawn); ok {
		var double_move [2]uint16
		if pawn.Is_White_Piece() {
			double_move = [2]uint16{pawn.Position[0], pawn.Position[1] - 2}
		} else {
			double_move = [2]uint16{pawn.Position[0], pawn.Position[1] + 2}
		}
		if new_position == double_move {
			pawn.Has_moved = moves_counter
		} else {
			pawn.Has_moved = 0
		}
		if index := pawn.can_do_enpassant(pieces_a, new_position, [2]uint16{pawn.Position[0] + 1, pawn.Position[1]}, moves_counter); index != -1 {
			pieces_a[index] = nil
			smth_has_been_taken = true
		}
		if index := pawn.can_do_enpassant(pieces_a, new_position, [2]uint16{pawn.Position[0] - 1, pawn.Position[1]}, moves_counter); index != -1 {
			pieces_a[index] = nil
			smth_has_been_taken = true
		}
	}
	piece.Move_To(new_position)

	return pieces_a, smth_has_been_taken
}

func (c *ChessObject) Move_To(new_position [2]uint16) {
	c.Position = new_position
}

func (c *ChessObject) Give_Pos() [2]uint16 {
	return c.Position
}

func (c *ChessObject) Piece_Is_White() bool {
	return c.White
}

type Positioning struct { //datentyp Positioning
	Position [2]uint16
}

type ChessObject struct { //datentyp ChessObject erbt vom datentyp Positioning
	Positioning
	White       bool
	Legal_Moves [][2]uint16
}

type Pawn struct { //alle Schachobjekte erben wiederum vom datentyp ChessObject
	Has_moved int16
	ChessObject
}

func NewPawn(x, y uint16, is_white bool) *Pawn {
	return &Pawn{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
		Has_moved:   -1,
	}
}

type Knight struct {
	ChessObject
}

func NewKnight(x, y uint16, is_white bool) *Knight {
	return &Knight{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

type Bishop struct {
	ChessObject
}

func NewBishop(x, y uint16, is_white bool) *Bishop {
	return &Bishop{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

type Rook struct {
	Has_moved bool
	ChessObject
}

func NewRook(x, y uint16, is_white bool) *Rook {
	return &Rook{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

type Queen struct {
	ChessObject
}

func NewQueen(x, y uint16, is_white bool) *Queen {
	return &Queen{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

type King struct {
	Has_moved bool
	ChessObject
}

func NewKing(x, y uint16, is_white bool) *King {
	return &King{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

func (p *Pawn) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()

	// var blocking_piece_1 bool
	// var blocking_piece_2 bool
	// new_legal_move_1 := [2]uint16{10, 10}
	// new_legal_move_2 := [2]uint16{10, 10}
	// new_legal_move_3 := [2]uint16{10, 10}
	// new_legal_move_4 := [2]uint16{10, 10}
	// en_passant_right := [2]uint16{p.Position[0] + 1, p.Position[1]}
	// en_passant_left := [2]uint16{p.Position[0] - 1, p.Position[1]}

	var direction int16 = 0
	var last_y uint16
	if p.Is_White_Piece() {
		direction = -1
		last_y = 0
	} else {
		direction = 1
		last_y = 7
	}

	if p.Position[1] != last_y {
		p.can_move(pieces_a, [2]uint16{p.Position[0], uint16(int16(p.Position[1]) + direction)})     //einer move
		p.can_take(pieces_a, [2]uint16{p.Position[0] + 1, uint16(int16(p.Position[1]) + direction)}) //schlagen rechts
		p.can_take(pieces_a, [2]uint16{p.Position[0] - 1, uint16(int16(p.Position[1]) + direction)}) //schlagen links
		if p.Position[1] != uint16(int16(last_y)-direction) && p.Has_moved == -1 {
			p.can_move(pieces_a, [2]uint16{p.Position[0], uint16(int16(p.Position[1]) + direction*2)}) //zweier move
		}
		p.can_do_enpassant(pieces_a, [2]uint16{p.Position[0] + 1, uint16(int16(p.Position[1]) + direction)}, [2]uint16{p.Position[0] + 1, p.Position[1]}, moves_counter)
		p.can_do_enpassant(pieces_a, [2]uint16{p.Position[0] - 1, uint16(int16(p.Position[1]) + direction)}, [2]uint16{p.Position[0] - 1, p.Position[1]}, moves_counter)
	}

	// //definition der hypothetischen felder je nach farbe
	// if p.Is_White_Piece() && p.Position[1] != 0 {
	// 	new_legal_move_1 = [2]uint16{p.Position[0], p.Position[1] - 1} //einer move
	// 	if p.Position[1] > 1 && p.Has_moved == -1 {
	// 		new_legal_move_2 = [2]uint16{p.Position[0], p.Position[1] - 2} //zweier move
	// 	}
	// 	new_legal_move_3 = [2]uint16{p.Position[0] + 1, p.Position[1] - 1} //schlagen rechts
	// 	new_legal_move_4 = [2]uint16{p.Position[0] - 1, p.Position[1] - 1} //schlagen links
	// } else if !p.Is_White_Piece() && p.Position[1] != 7 {
	// 	new_legal_move_1 = [2]uint16{p.Position[0], p.Position[1] + 1} //einer move
	// 	if p.Position[1] < 6 && p.Has_moved == -1 {
	// 		new_legal_move_2 = [2]uint16{p.Position[0], p.Position[1] + 2} //zweier move
	// 	}
	// 	new_legal_move_3 = [2]uint16{p.Position[0] + 1, p.Position[1] + 1} //schlagen rechts
	// 	new_legal_move_4 = [2]uint16{p.Position[0] - 1, p.Position[1] + 1} //schlagen links
	// }

	//überprüfen der hyothetischen felder
	// for i := 0; i < len(pieces_a) && (!blocking_piece_1 || !blocking_piece_2); i++ {
	// 	if pieces_a[i] != nil {
	// 		if pieces_a[i].Give_Pos() == new_legal_move_1 {
	// 			blocking_piece_1 = true
	// 		} else if pieces_a[i].Give_Pos() == new_legal_move_2 {
	// 			blocking_piece_2 = true
	// 		} else if pieces_a[i].Give_Pos() == new_legal_move_3 && pieces_a[i].Is_White_Piece() != p.Is_White_Piece() { //schlagen rechts
	// 			p.Append_Legal_Moves(new_legal_move_3)
	// 		} else if pieces_a[i].Give_Pos() == new_legal_move_4 && pieces_a[i].Is_White_Piece() != p.Is_White_Piece() { //schlagen links
	// 			p.Append_Legal_Moves(new_legal_move_4)
	// 		} else if en_passant_pawn1, ok := pieces_a[i].(*Pawn); ok && pieces_a[i].Is_White_Piece() != p.Is_White_Piece() && pieces_a[i].Give_Pos() == en_passant_right {
	// 			//andersfarbiger pawn rechts neben dem pawn --> en passant rechts
	// 			if en_passant_pawn1.Has_moved > 0 && en_passant_pawn1.Has_moved+1 == moves_counter {
	// 				p.Append_Legal_Moves(new_legal_move_3)
	// 			}
	// 		} else if en_passant_pawn2, ok := pieces_a[i].(*Pawn); ok && pieces_a[i].Is_White_Piece() != p.Is_White_Piece() && pieces_a[i].Give_Pos() == en_passant_left {
	// 			//andersfarbiger pawn links neben dem pawn --> en passant links
	// 			fmt.Println("en passant")
	// 			if en_passant_pawn2.Has_moved > 0 && en_passant_pawn2.Has_moved+1 == moves_counter {
	// 				p.Append_Legal_Moves(new_legal_move_4)
	// 			}
	// 		}
	// 	}
	// }

	// if !blocking_piece_1 && new_legal_move_1 != [2]uint16{10, 10} { //es steht nichts im weg direkt davor einer move
	// 	p.Append_Legal_Moves(new_legal_move_1)
	// }
	// if !blocking_piece_1 && !blocking_piece_2 && new_legal_move_2 != [2]uint16{10, 10} { //es steht nichts im weg direkt davor zweier move
	// 	p.Append_Legal_Moves(new_legal_move_2)
	// }

}

func (p *Pawn) can_take(pieces_a [64]Piece, field [2]uint16) {
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if pieces_a[i].Is_White_Piece() != p.Is_White_Piece() && pieces_a[i].Give_Pos() == field {
				p.Append_Legal_Moves(field)
				break
			}
		}
	}
}

func (p *Pawn) can_move(pieces_a [64]Piece, field [2]uint16) {
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
		p.Append_Legal_Moves(field)
	}
}

func (p *Pawn) can_do_enpassant(pieces_a [64]Piece, field, en_passant_pawn_pos [2]uint16, moves_counter int16) int {
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if en_passant_pawn, ok := pieces_a[i].(*Pawn); ok && pieces_a[i].Is_White_Piece() != p.Is_White_Piece() && pieces_a[i].Give_Pos() == en_passant_pawn_pos {
				//andersfarbiger pawn rechts neben dem pawn --> en passant rechts
				if en_passant_pawn.Has_moved > 0 && en_passant_pawn.Has_moved+1 == moves_counter {
					p.Append_Legal_Moves(field)
					return i
				}
			}
		}
	}
	return -1
}

func (p *Knight) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	fmt.Printf("Moves of Knight")
}

func (p *Rook) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Legal_Moves = nil

	for new_x := p.Position[0]; new_x < 7; {
		new_x++
		var current_pos [2]uint16 = [2]uint16{new_x, p.Position[1]}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x := p.Position[0]; new_x != 0; {
		new_x--
		var current_pos [2]uint16 = [2]uint16{new_x, p.Position[1]}
		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
		if new_x == 0 {
			break
		}

	}

	for new_y := p.Position[1]; new_y < 7; {
		new_y++
		var current_pos [2]uint16 = [2]uint16{p.Position[0], new_y}
		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_y := p.Position[1]; new_y != 0; {
		new_y--
		var current_pos [2]uint16 = [2]uint16{p.Position[0], new_y}
		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
		if new_y == 0 {
			break
		}

	}
}

func (p *Bishop) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Legal_Moves = nil

	for new_x, new_y := p.Position[0], p.Position[1]; new_x < 7 && new_y < 7; {
		new_x++
		new_y++
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Position[0], p.Position[1]; new_x < 7 && new_y != 0; {
		new_x++
		new_y--
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Position[0], p.Position[1]; new_x != 0 && new_y < 7; {
		new_x--
		new_y++
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Position[0], p.Position[1]; new_x != 0 && new_y != 0; {
		new_x--
		new_y--
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}
}

func (p *Queen) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	fmt.Printf("Moves of Queen")
}

func (p *King) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	fmt.Printf("Moves of King")
}

func check_if_piece_is_blocking(p Piece, pieces_a [64]Piece, current_pos [2]uint16) bool {
	var blocking_piece Piece
	var var_break bool = false

	for i := 0; i < len(pieces_a) && blocking_piece == nil; i++ {
		if pieces_a[i] != nil {
			if pieces_a[i].Give_Pos() == current_pos {
				blocking_piece = pieces_a[i]
			}
		}
	}

	if blocking_piece == nil { //es steht nichts im weg
		p.Append_Legal_Moves(current_pos)
	} else if blocking_piece.Is_White_Piece() != p.Is_White_Piece() { //es steht etwas im weg, was aber geschlagen werden kann, daher wird danach jedoch gebreaked
		p.Append_Legal_Moves(current_pos)
		var_break = true
	} else if blocking_piece.Is_White_Piece() == p.Is_White_Piece() { //es steht etwas im weg, was aber nicht geschlagen werden kann, daher wird sofort gebreaked
		var_break = true
	} else {
		fmt.Println("fatal: Error in Calculating Moves Method")
	}
	return var_break
}
