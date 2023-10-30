package pieces

import (
	"fmt"

	gfx "../gfxw"
	"../path"
)

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

func (c *ChessObject) Give_Pos() [2]uint16 {
	return c.Position
}

func (c *ChessObject) Give_Piece_Type() string {
	return "Error"
}

func (c *ChessObject) Piece_Is_White() bool {
	return c.White
}

func (c *ChessObject) Give_Legal_Moves() [][3]uint16 {
	return c.Legal_Moves
}

func (c *ChessObject) Give_Has_Moved() int16 {
	return c.Has_moved
}

func (c *ChessObject) Set_Has_Moved(update int16) {
	c.Has_moved = update
}

func (c *ChessObject) Clear_Legal_Moves() {
	c.Legal_Moves = nil
}

func (c *ChessObject) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	//pass
}

func (c *ChessObject) Append_Legal_Moves(new_legal_move [3]uint16) {
	c.Legal_Moves = append(c.Legal_Moves, new_legal_move)
}

func (c *ChessObject) Is_White_Piece() bool {
	return c.White
}

func (p *Pawn) DeepCopy(current_piece Piece) Piece {
	dst := *p
	return &dst
}

func (p *Rook) DeepCopy(current_piece Piece) Piece {
	dst := *p
	return &dst
}

func (p *Knight) DeepCopy(current_piece Piece) Piece {
	dst := *p
	return &dst
}

func (p *King) DeepCopy(current_piece Piece) Piece {
	dst := *p
	return &dst
}

func (p *Queen) DeepCopy(current_piece Piece) Piece {
	dst := *p
	return &dst
}

func (p *Bishop) DeepCopy(current_piece Piece) Piece {
	dst := *p
	return &dst
}

func (c *Pawn) Give_Piece_Type() string {
	return ""
}
func (c *Knight) Give_Piece_Type() string {
	return "N"
}
func (c *Bishop) Give_Piece_Type() string {
	return "B"
}
func (c *Rook) Give_Piece_Type() string {
	return "R"
}
func (c *Queen) Give_Piece_Type() string {
	return "Q"
}
func (c *King) Give_Piece_Type() string {
	return "K"
}

func NewPawn(x, y uint16, is_white bool) *Pawn {
	return &Pawn{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white, Has_moved: -1},
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

func Draw(piece Piece, w_x, w_y, a uint16) {
	Copy_Piece_To_Clipboard(piece, w_x, w_y, a)
}

func Draw_To_Point(piece Piece, w_x, w_y, a, x, y uint16, x_offset, y_offset int16, transparencey uint8, m_x uint16) {

	if transparencey == 0 {
		gfx.Archivieren()
	}

	if m_x >= 75*a/10 {
		x_offset = 0
		x = 70 * a / 10
	}

	gfx.UpdateAus()

	gfx.Restaurieren(0, 0, 8*a, w_y)

	gfx.Archivieren()
	Copy_Piece_To_Clipboard(piece, w_x, w_y, a)
	gfx.Restaurieren(0, 0, 8*a, w_y)

	gfx.Archivieren()

	gfx.Transparenz(transparencey)
	gfx.Clipboard_einfuegenMitColorKey(uint16(int16(x)+x_offset), uint16(int16(y)+y_offset), 5, 5, 5)
	gfx.Transparenz(0)

	gfx.UpdateAn()
}

func Move_Piece_To(piece Piece, new_position [3]uint16, moves_counter int16, pieces_a [64]Piece) ([64]Piece, uint16, string) { //rook and king has moved change
	var promotion uint16 = 64
	var take string = ""

	if king, ok := piece.(*King); ok {
		if new_position[2] == 64 { //normal move
			king.Has_moved = 1
		} else if new_position[2] <= 63 { //castle
			king.Has_moved = 1
			pieces_a[new_position[2]].Set_Has_Moved(1)
			rook := pieces_a[new_position[2]]

			if rook.Give_Pos()[0] == 7 { //right castle
				new_position = [3]uint16{king.Give_Pos()[0] + 2, king.Give_Pos()[1], 64}
				rook.Move_To([2]uint16{rook.Give_Pos()[0] - 2, rook.Give_Pos()[1]})
			} else if rook.Give_Pos()[0] == 0 { //left castle
				new_position = [3]uint16{king.Give_Pos()[0] - 2, king.Give_Pos()[1], 64}
				rook.Move_To([2]uint16{rook.Give_Pos()[0] + 3, rook.Give_Pos()[1]})
			} else {
				fmt.Println("Panic: Error has occured, Rook for Castle is not on expected position")
			}
		} else if new_position[2] >= 66 && new_position[2] <= 129 { //king take move
			king.Has_moved = 0
			pieces_a[new_position[2]-66] = nil
			take = "x"
		} else {
			fmt.Println("panic: Error has occured, king move status is out of range")
		}

	} else if pawn, ok := piece.(*Pawn); ok {
		if new_position[2] == 65 { //double_move
			pawn.Has_moved = moves_counter
		} else if new_position[2] <= 63 { //en passant
			pieces_a[new_position[2]] = nil
			pawn.Has_moved = 0
			take = "x"
		} else if new_position[2] == 64 { //normal move
			pawn.Has_moved = 0
		} else if new_position[2] >= 66 && new_position[2] <= 129 { //take move
			pawn.Has_moved = 0
			pieces_a[new_position[2]-66] = nil
			take = "x"
		} else {
			fmt.Println("panic: Error has occured, pawn move status is out of range")
		}
		if (new_position[1] == 0 && pawn.Is_White_Piece()) || new_position[1] == 7 && !pawn.Is_White_Piece() { //promotion
			for i := 0; i < len(pieces_a); i++ {
				if pieces_a[i] != nil {
					if pieces_a[i] == pawn {
						promotion = uint16(i)
					}
				}
			}
		}

	} else {
		if new_position[2] == 64 { //normal piece normal move
			piece.Set_Has_Moved(1)
		} else if new_position[2] <= 63 { //normal piece take move
			piece.Set_Has_Moved(1)
			pieces_a[new_position[2]] = nil
			take = "x"
		} else {
			fmt.Println(new_position[2])
			fmt.Println("panic: Error has occured, normal piece move status is out of range")
		}
	}

	piece.Move_To([2]uint16{new_position[0], new_position[1]})

	return pieces_a, promotion, take
}

func Find_Piece_With_Pos(pieces_a [64]Piece, field [2]uint16) int {
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil && pieces_a[i].Give_Pos() == field {
			return i
		}
	}
	return -1
}

func move_is_in_board(move [2]uint16) bool {
	if move[0] <= 7 && move[1] <= 7 {
		return true
	} else {
		return false
	}
}

func Field_Can_Be_Captured(pieces_that_can_capture_are_white bool, field [2]uint16, pieces_a [64]Piece, moves_counter int16) bool {
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil && pieces_a[i].Is_White_Piece() == pieces_that_can_capture_are_white {
			pieces_a[i].Calc_Moves(pieces_a, moves_counter)
			current_legal_moves := pieces_a[i].Give_Legal_Moves()
			for k := 0; k < len(current_legal_moves); k++ {
				legal_move := [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]}
				if legal_move == field {
					return true
				}
			}
		}
	}
	return false
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
		//einer move
		p.try_to_move(pieces_a, [2]uint16{p.Position[0], uint16(int16(p.Position[1]) + direction)}, 64)
		//schlagen rechts
		p.try_to_take(pieces_a, [2]uint16{p.Position[0] + 1, uint16(int16(p.Position[1]) + direction)}, true)
		//schlagen links
		p.try_to_take(pieces_a, [2]uint16{p.Position[0] - 1, uint16(int16(p.Position[1]) + direction)}, true)
		//zweier move
		if p.Position[1] != uint16(int16(last_y)-direction) && p.Has_moved == -1 && Find_Piece_With_Pos(pieces_a, [2]uint16{p.Give_Pos()[0], uint16(int16(p.Give_Pos()[1]) + direction)}) == -1 { //überprüft nicht ob etwas davor steht
			p.try_to_move(pieces_a, [2]uint16{p.Position[0], uint16(int16(p.Position[1]) + direction*2)}, 65)
		}
		//enpassant rechts
		p.can_do_enpassant(pieces_a, [2]uint16{p.Position[0] + 1, uint16(int16(p.Position[1]) + direction)}, [2]uint16{p.Position[0] + 1, p.Position[1]}, moves_counter)
		//enpassant links
		p.can_do_enpassant(pieces_a, [2]uint16{p.Position[0] - 1, uint16(int16(p.Position[1]) + direction)}, [2]uint16{p.Position[0] - 1, p.Position[1]}, moves_counter)
	}
}

func (p *Knight) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 2, p.Give_Pos()[1] + 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 2, p.Give_Pos()[1] - 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 2, p.Give_Pos()[1] + 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 2, p.Give_Pos()[1] - 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1] + 2})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1] - 2})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1] + 2})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1] - 2})
}

func (p *Rook) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()
	p.calc_moves_vertically_and_horizontally(pieces_a)
}

func (p *Bishop) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()
	p.calc_moves_diagonally(pieces_a)
}

func (p *Queen) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()
	p.calc_moves_vertically_and_horizontally(pieces_a)
	p.calc_moves_diagonally(pieces_a)
}

func (p *King) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()

	var right_rochade_possible uint16 = 64
	var left_rochade_possible uint16 = 64
	var blocking_right_rochade bool
	var blocking_left_rochade bool

	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1]})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1] + 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1] - 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1]})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1] + 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1] - 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0], p.Give_Pos()[1] + 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0], p.Give_Pos()[1] - 1})

	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if Rook, ok := pieces_a[i].(*Rook); ok {
				if Rook.Is_White_Piece() == p.Is_White_Piece() {
					if uint16(int16(Rook.Give_Pos()[0])-int16(p.Give_Pos()[0])) == 3 && Rook.Has_moved == 0 && p.Has_moved == 0 {
						right_rochade_possible = uint16(i)
					} else if uint16(int16(p.Give_Pos()[0])-int16(Rook.Give_Pos()[0])) == 4 && Rook.Has_moved == 0 && p.Has_moved == 0 {
						left_rochade_possible = uint16(i)
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

	if !blocking_right_rochade && right_rochade_possible != 64 {
		//p.Append_Legal_Moves([3]uint16{7, p.Give_Pos()[1], uint16(right_rochade_possible)})
		p.Append_Legal_Moves([3]uint16{6, p.Give_Pos()[1], uint16(right_rochade_possible)})
	}
	if !blocking_left_rochade && left_rochade_possible != 64 {
		//p.Append_Legal_Moves([3]uint16{0, p.Give_Pos()[1], uint16(left_rochade_possible)})
		p.Append_Legal_Moves([3]uint16{2, p.Give_Pos()[1], uint16(left_rochade_possible)})
	}
}

func (p *King) Is_In_Check(pieces_a [64]Piece, moves_counter int16) bool {
	var field [2]uint16 = p.Give_Pos()

	if Field_Can_Be_Captured(!p.Is_White_Piece(), field, pieces_a, moves_counter) {
		return true
	} else {
		return false
	}
}

func Copy_Array(pieces_a [64]Piece) [64]Piece {
	var copy_of_pieces_a [64]Piece
	for i, current_piece := range pieces_a {
		if current_piece != nil {
			copy_of_pieces_a[i] = current_piece.DeepCopy(current_piece)
		}
	}
	return copy_of_pieces_a
}

func Calc_Moves_With_Check(pieces_a [64]Piece, moves_counter int16, current_king_index int) ([64]Piece, bool) {

	var current_legal_moves [][3]uint16
	var checkmate bool = true

	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil && pieces_a[i].Is_White_Piece() == pieces_a[current_king_index].Is_White_Piece() {

			pieces_a[i].Calc_Moves(pieces_a, moves_counter) //includes step one and two of calculating moves

			current_legal_moves = pieces_a[i].Give_Legal_Moves()

			pieces_a[i].Clear_Legal_Moves()

			// if i == current_king_index {
			// 	fmt.Println(current_legal_moves)
			// }

			for k := 0; k < len(current_legal_moves); k++ { //iterates trough the legal moves for the current piece
				temp_pieces_a := Copy_Array(pieces_a) //this creates a deep copy of the pieces array --> resets after evry legal move
				temp_pieces_a, _, _ = Move_Piece_To(temp_pieces_a[i], current_legal_moves[k], moves_counter, temp_pieces_a)

				//fmt.Println(temp_pieces_a[current_king_index].(*King).Give_Pos())

				if !temp_pieces_a[current_king_index].(*King).Is_In_Check(temp_pieces_a, moves_counter) {

					if i == current_king_index && current_legal_moves[k][2] <= 63 { //rochade move
						rook := pieces_a[current_legal_moves[k][2]].(*Rook)
						if rook.Give_Pos()[0] == 7 &&
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), rook.Give_Pos(), pieces_a, moves_counter) &&
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), [2]uint16{rook.Give_Pos()[0] - 2, rook.Give_Pos()[1]}, pieces_a, moves_counter) &&
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), pieces_a[current_king_index].Give_Pos(), pieces_a, moves_counter) { //right castle
							pieces_a[i].Append_Legal_Moves(current_legal_moves[k])
							checkmate = false
						} else if rook.Give_Pos()[0] == 0 &&
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), rook.Give_Pos(), pieces_a, moves_counter) &&
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), [2]uint16{rook.Give_Pos()[0] + 1, rook.Give_Pos()[1]}, pieces_a, moves_counter) &&
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), pieces_a[current_king_index].Give_Pos(), pieces_a, moves_counter) { //left castle
							pieces_a[i].Append_Legal_Moves(current_legal_moves[k])
							checkmate = false
						} //else: rochade is not possible because one of the involved squares can be captured
					} else {
						pieces_a[i].Append_Legal_Moves(current_legal_moves[k])
						checkmate = false
					}
				}
			}
		}
	}
	return pieces_a, checkmate
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

func (p *ChessObject) try_to_take(pieces_a [64]Piece, field [2]uint16, piece_is_king_or_pawn bool) {
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if pieces_a[i].Is_White_Piece() != p.Is_White_Piece() && pieces_a[i].Give_Pos() == field {
				if piece_is_king_or_pawn {
					p.Append_Legal_Moves([3]uint16{field[0], field[1], uint16(i + 66)})
				} else {
					p.Append_Legal_Moves([3]uint16{field[0], field[1], uint16(i)})
				}

			}
		}
	}
}

func (p *ChessObject) try_to_move(pieces_a [64]Piece, field [2]uint16, status uint16) {
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
		p.Append_Legal_Moves([3]uint16{field[0], field[1], status})
	}
}

func (p *Knight) Calc_Normal_Move(pieces_a [64]Piece, field [2]uint16) {
	if move_is_in_board(field) {
		p.try_to_move(pieces_a, field, 64)
		p.try_to_take(pieces_a, field, false)
	}
}

func (p *King) Calc_Normal_Move(pieces_a [64]Piece, field [2]uint16) {
	if move_is_in_board(field) {
		p.try_to_move(pieces_a, field, 64)
		p.try_to_take(pieces_a, field, true)
	}
}

func (p *ChessObject) check_if_piece_is_blocking(pieces_a [64]Piece, current_pos [2]uint16) bool {
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

func (p *ChessObject) calc_moves_diagonally(pieces_a [64]Piece) {
	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x < 7 && new_y < 7; {
		new_x++
		new_y++
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x < 7 && new_y != 0; {
		new_x++
		new_y--
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x != 0 && new_y < 7; {
		new_x--
		new_y++
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x != 0 && new_y != 0; {
		new_x--
		new_y--
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}
}

func (p *ChessObject) calc_moves_vertically_and_horizontally(pieces_a [64]Piece) {
	for new_x := p.Give_Pos()[0]; new_x < 7; {
		new_x++
		var current_pos [2]uint16 = [2]uint16{new_x, p.Give_Pos()[1]}

		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}

	for new_x := p.Give_Pos()[0]; new_x != 0; {
		new_x--
		var current_pos [2]uint16 = [2]uint16{new_x, p.Give_Pos()[1]}
		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
		if new_x == 0 {
			break
		}

	}

	for new_y := p.Give_Pos()[1]; new_y < 7; {
		new_y++
		var current_pos [2]uint16 = [2]uint16{p.Give_Pos()[0], new_y}
		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}

	for new_y := p.Give_Pos()[1]; new_y != 0; {
		new_y--
		var current_pos [2]uint16 = [2]uint16{p.Give_Pos()[0], new_y}
		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
		if new_y == 0 {
			break
		}

	}
}

func Copy_Piece_To_Clipboard(piece Piece, w_x, w_y, a uint16) {
	path := path.Give_Path()
	gfx.LadeBild(0, 0, (path + "\\resources\\images\\Pieces.bmp"))

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
