package pieces

import (
	"fmt"
)

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

func (c *ChessObject) Give_Pos() [2]uint16 {
	return c.Position
}

func (c *ChessObject) Piece_Is_White() bool {
	return c.White
}

func NewPawn(x, y uint16, is_white bool) *Pawn {
	return &Pawn{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
		Has_moved:   -1,
	}
}

func NewKnight(x, y uint16, is_white bool) *Knight {
	return &Knight{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

func NewBishop(x, y uint16, is_white bool) *Bishop {
	return &Bishop{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

func NewRook(x, y uint16, is_white bool) *Rook {
	return &Rook{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

func NewQueen(x, y uint16, is_white bool) *Queen {
	return &Queen{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

func NewKing(x, y uint16, is_white bool) *King {
	return &King{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

func (p *Pawn) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()

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
		p.can_do_enpassant(pieces_a, [2]uint16{p.Position[0] + 1, uint16(int16(p.Position[1]) + direction)}, [2]uint16{p.Position[0] + 1, p.Position[1]}, moves_counter) //enpassant
		p.can_do_enpassant(pieces_a, [2]uint16{p.Position[0] - 1, uint16(int16(p.Position[1]) + direction)}, [2]uint16{p.Position[0] - 1, p.Position[1]}, moves_counter) //enpassant
	}
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
