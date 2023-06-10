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

func Move_Piece_To(piece Piece, new_position [2]uint16, moves_counter int16, pieces_a [64]Piece) ([64]Piece, bool) { //rook and king has moved change
	var smth_has_been_taken bool = false

	if King, ok := piece.(*King); ok {

		castle_moves := King.Calc_Castle_Moves(pieces_a)

		for k := 0; k < 2; k++ {
			if castle_moves[k] == new_position {
				if k == 0 { //right castlef

					new_position = [2]uint16{King.Give_Pos()[0] + 2, King.Give_Pos()[1]}
					Rook := pieces_a[castle_moves[2][0]]
					Rook.Move_To([2]uint16{Rook.Give_Pos()[0] - 2, Rook.Give_Pos()[1]})
				} else if k == 1 { //left castle
					new_position = [2]uint16{King.Give_Pos()[0] - 2, King.Give_Pos()[1]}
					Rook := pieces_a[castle_moves[2][1]]
					Rook.Move_To([2]uint16{Rook.Give_Pos()[0] + 3, Rook.Give_Pos()[1]})
				}
			}
		}
	}

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

func (p *King) calc_normal_move(pieces_a [64]Piece, field [2]uint16) {
	can_move(p, pieces_a, field)
	can_take(p, pieces_a, field)
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
		can_move(p, pieces_a, [2]uint16{p.Position[0], uint16(int16(p.Position[1]) + direction)})     //einer move
		can_take(p, pieces_a, [2]uint16{p.Position[0] + 1, uint16(int16(p.Position[1]) + direction)}) //schlagen rechts
		can_take(p, pieces_a, [2]uint16{p.Position[0] - 1, uint16(int16(p.Position[1]) + direction)}) //schlagen links
		if p.Position[1] != uint16(int16(last_y)-direction) && p.Has_moved == -1 {
			can_move(p, pieces_a, [2]uint16{p.Position[0], uint16(int16(p.Position[1]) + direction*2)}) //zweier move
		}
		p.can_do_enpassant(pieces_a, [2]uint16{p.Position[0] + 1, uint16(int16(p.Position[1]) + direction)}, [2]uint16{p.Position[0] + 1, p.Position[1]}, moves_counter) //enpassant
		p.can_do_enpassant(pieces_a, [2]uint16{p.Position[0] - 1, uint16(int16(p.Position[1]) + direction)}, [2]uint16{p.Position[0] - 1, p.Position[1]}, moves_counter) //enpassant
	}
}

func (p *Knight) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	fmt.Printf("Moves of Knight")
}

func (p *Rook) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	fmt.Println("this is rook: ", p.Give_Pos())
	p.Clear_Legal_Moves()
	calc_moves_vertically_and_horizontally(p, pieces_a)
}

func (p *Bishop) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()
	calc_moves_diagonally(p, pieces_a)
}

func (p *Queen) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()
	calc_moves_vertically_and_horizontally(p, pieces_a)
	calc_moves_diagonally(p, pieces_a)
}

func (p *King) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()

	if p.Give_Pos()[0] < 7 {
		p.calc_normal_move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1]})
		if p.Give_Pos()[1] < 7 {
			p.calc_normal_move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1] + 1})
		}
		if p.Give_Pos()[1] > 0 {
			p.calc_normal_move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1] - 1})
		}
	}
	if p.Give_Pos()[0] > 0 {
		p.calc_normal_move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1]})
		if p.Give_Pos()[1] < 7 {
			p.calc_normal_move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1] + 1})
		}
		if p.Give_Pos()[1] > 0 {
			p.calc_normal_move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1] - 1})
		}
	}
	if p.Give_Pos()[1] < 7 {
		p.calc_normal_move(pieces_a, [2]uint16{p.Give_Pos()[0], p.Give_Pos()[1] + 1})
	}
	if p.Give_Pos()[1] > 0 {
		p.calc_normal_move(pieces_a, [2]uint16{p.Give_Pos()[0], p.Give_Pos()[1] - 1})
	}

	castle_moves := p.Calc_Castle_Moves(pieces_a)
	for k := 0; k < 2; k++ {
		if castle_moves[k] != [2]uint16{8, 8} {
			p.Append_Legal_Moves(castle_moves[k])
		}
	}
}

func (p *King) Calc_Castle_Moves(pieces_a [64]Piece) [3][2]uint16 {
	castle_moves := [3][2]uint16{{8, 8}, {8, 8}, {64, 64}}
	var blocking_right_rochade bool
	var blocking_left_rochade bool
	var right_rook_has_moved bool = true
	var left_rook_has_moved bool = true

	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if Rook, ok := pieces_a[i].(*Rook); ok {
				if Rook.Is_White_Piece() == p.Is_White_Piece() {
					if uint16(int16(Rook.Give_Pos()[0])-int16(p.Give_Pos()[0])) == 3 && !Rook.Has_moved {
						right_rook_has_moved = false
						castle_moves[2][0] = uint16(i)
					} else if uint16(int16(p.Give_Pos()[0])-int16(Rook.Give_Pos()[0])) == 4 && !Rook.Has_moved {
						left_rook_has_moved = false
						castle_moves[2][1] = uint16(i)
					}
				}
			}
			if pieces_a[i].Give_Pos() == [2]uint16{5, p.Give_Pos()[1]} || pieces_a[i].Give_Pos() == [2]uint16{6, p.Give_Pos()[1]} {
				blocking_right_rochade = true
			}
			if pieces_a[i].Give_Pos() == [2]uint16{1, p.Give_Pos()[1]} || pieces_a[i].Give_Pos() == [2]uint16{2, p.Give_Pos()[1]} || pieces_a[i].Give_Pos() == [2]uint16{3, p.Give_Pos()[1]} {
				blocking_left_rochade = true
			}
		}
	}

	if !blocking_right_rochade && !right_rook_has_moved && !p.Has_moved {
		castle_moves[0] = [2]uint16{7, p.Give_Pos()[1]} //right castle
	}
	if !blocking_left_rochade && !left_rook_has_moved && !p.Has_moved {
		castle_moves[1] = [2]uint16{0, p.Give_Pos()[1]} //left castle
	}

	return castle_moves
}
