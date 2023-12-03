package pieces

import (
	"fmt"

	gfx "../gfxw"
	"../path"
)

type chess_object struct { //datentyp chess_object erbt vom datentyp Positioning
	position    [2]uint16
	white       bool
	legal_moves [][3]uint16
	has_moved   int16
}

type Pawn struct { //alle Schachobjekte erben wiederum vom datentyp chess_object
	chess_object
}

type Knight struct {
	chess_object
}

type Bishop struct {
	chess_object
}

type Rook struct {
	chess_object
}

type Queen struct {
	chess_object
}

type King struct {
	chess_object
}

// general functions (chess_object is inherited by every piece type)
func (c *chess_object) Move_To(new_position [2]uint16) {
	c.position = new_position
}
func (c *chess_object) Give_Pos() [2]uint16 {
	return c.position
}

func (c *chess_object) Give_Legal_Moves() [][3]uint16 {
	return c.legal_moves
}
func (c *chess_object) Give_Has_Moved() int16 {
	return c.has_moved
}
func (c *chess_object) Set_Has_Moved(update int16) {
	c.has_moved = update
}
func (c *chess_object) Clear_Legal_Moves() {
	c.legal_moves = nil
}
func (c *chess_object) Append_Legal_Moves(new_legal_move [3]uint16) {
	c.legal_moves = append(c.legal_moves, new_legal_move)
}
func (c *chess_object) Is_White_Piece() bool {
	return c.white
}
func (c *chess_object) Calc_Moves(pieces_a [64]Piece, moves_counter int16) { //not needed but required for the interface --> overwritten by every piece type
	fmt.Println("Error: chess_object is not a valid piece type --> therefore Calc_Moves can't be executed")
}
func (c *chess_object) Give_Piece_Type() string { //not needed but required for the interface --> overwritten by every piece type
	return "Error: chess_object is not a valid piece type --> therefore Give_Piece_Type can't be executed"
}

// deep copy functions need to be written for each pieces individually due to the different data types
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

// give piece type need to be rewritten for each piece. obviously
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

// new functions must be rewritten for each piece type individually as well
func NewPawn(x, y uint16, is_white bool) *Pawn {
	return &Pawn{
		chess_object: chess_object{position: [2]uint16{x, y}, white: is_white, has_moved: -1},
	}
}
func NewKnight(x, y uint16, is_white bool) *Knight {
	return &Knight{
		chess_object: chess_object{position: [2]uint16{x, y}, white: is_white},
	}
}
func NewBishop(x, y uint16, is_white bool) *Bishop {
	return &Bishop{
		chess_object: chess_object{position: [2]uint16{x, y}, white: is_white},
	}
}
func NewRook(x, y uint16, is_white bool) *Rook {
	return &Rook{
		chess_object: chess_object{position: [2]uint16{x, y}, white: is_white},
	}
}
func NewQueen(x, y uint16, is_white bool) *Queen {
	return &Queen{
		chess_object: chess_object{position: [2]uint16{x, y}, white: is_white},
	}
}
func NewKing(x, y uint16, is_white bool) *King {
	return &King{
		chess_object: chess_object{position: [2]uint16{x, y}, white: is_white},
	}
}

// general functions:

func Draw(piece Piece, w_x, w_y, a uint16) {
	Copy_Piece_To_Clipboard(piece, w_x, w_y, a)
}

func Draw_To_Point(piece Piece, w_x, w_y, a, x, y uint16, x_offset, y_offset int16, transparencey uint8, m_x uint16) {
	//create a ghost piece that is spawed at the desired location which is given in pixel cordinates
	if transparencey == 0 { //neeeded for the promotion menu, otherwise only one piece could be displayed at the same time
		gfx.Archivieren()
	}

	if m_x >= 75*a/10 { //makes sure that no ghost piece is drawed to the sidebar area
		x_offset = 0
		x = 70 * a / 10
	}

	gfx.UpdateAus()

	gfx.Restaurieren(0, 0, 8*a, w_y) //needed for the deletion of the ghost piece

	gfx.Archivieren()                           //needed for the deletion of the pieces sprite ()
	Copy_Piece_To_Clipboard(piece, w_x, w_y, a) //copy the desired piece
	gfx.Restaurieren(0, 0, 8*a, w_y)            //needed for the deletion of the pieces sprite

	gfx.Archivieren() //needed for the deletion of the ghost piece ()

	gfx.Transparenz(transparencey)                                                                    //set the transparency
	gfx.Clipboard_einfuegenMitColorKey(uint16(int16(x)+x_offset), uint16(int16(y)+y_offset), 5, 5, 5) //paste the piece to the correct location
	gfx.Transparenz(0)                                                                                //reset transparecy

	gfx.UpdateAn()
}

func Move_Piece_To(piece Piece, new_position [3]uint16, moves_counter int16, pieces_a [64]Piece) ([64]Piece, uint16, string) { //rook and king has moved change
	var promotion uint16 = 64
	var take string = ""

	if king, ok := piece.(*King); ok { //king has castle as a special move
		if new_position[2] == 64 { //normal move
			king.Set_Has_Moved(1)
		} else if new_position[2] <= 63 { //castle
			king.Set_Has_Moved(1)
			pieces_a[new_position[2]].Set_Has_Moved(1)
			rook := pieces_a[new_position[2]] //on pos two of the legal move the index of the rook is saved
			rook.Set_Has_Moved(1)

			if rook.Give_Pos()[0] == 7 { //right castle
				new_position = [3]uint16{king.Give_Pos()[0] + 2, king.Give_Pos()[1], 64}
				rook.Move_To([2]uint16{rook.Give_Pos()[0] - 2, rook.Give_Pos()[1]})
			} else if rook.Give_Pos()[0] == 0 { //left castle
				new_position = [3]uint16{king.Give_Pos()[0] - 2, king.Give_Pos()[1], 64}
				rook.Move_To([2]uint16{rook.Give_Pos()[0] + 3, rook.Give_Pos()[1]})
			} else {
				rook.Set_Has_Moved(0)
				fmt.Println("Panic: Error has occured, Rook for Castle is not on expected position")
			}
		} else if new_position[2] >= 66 && new_position[2] <= 129 { //king take move
			king.Set_Has_Moved(1)
			pieces_a[new_position[2]-66] = nil //delete the coresponding piece
			take = "x"
		} else {
			fmt.Println("panic: Error has occured, king move status is out of range")
		}

	} else if pawn, ok := piece.(*Pawn); ok { //pawn has special moves, so the pawn moves must be considered seperately
		if new_position[2] == 65 { //double_move
			pawn.Set_Has_Moved(moves_counter) //needed for en passant --> to check when the last two pawn move happened
		} else if new_position[2] <= 63 { //en passant
			pieces_a[new_position[2]] = nil
			pawn.Set_Has_Moved(0)
			take = "x"
		} else if new_position[2] == 64 { //normal move
			pawn.Set_Has_Moved(0)
		} else if new_position[2] >= 66 && new_position[2] <= 129 { //take move
			pawn.Set_Has_Moved(0)
			pieces_a[new_position[2]-66] = nil
			take = "x"
		} else {
			fmt.Println("panic: Error has occured, pawn move status is out of range")
		}
		if (new_position[1] == 0 && pawn.Is_White_Piece()) || new_position[1] == 7 && !pawn.Is_White_Piece() { //promotion
			for i := 0; i < len(pieces_a); i++ {
				if pieces_a[i] != nil {
					if pieces_a[i] == pawn { //for the promotion function the index of the panw promoting must be returned
						promotion = uint16(i)
					}
				}
			}
		}

	} else { //normal piece moves
		if new_position[2] == 64 { //normal piece normal move
			piece.Set_Has_Moved(1)
		} else if new_position[2] <= 63 { //normal piece take move indicated by the third element of the legal_moves array being lesser than 64
			piece.Set_Has_Moved(1)
			pieces_a[new_position[2]] = nil
			take = "x"
		} else {
			fmt.Println("panic: Error has occured, normal piece move status is out of range with:", new_position[2])
		}
	}

	piece.Move_To([2]uint16{new_position[0], new_position[1]}) //move the actual piece --> ervery piece has to be moved regardless of the piece type

	return pieces_a, promotion, take
}

func Find_Piece_With_Pos(pieces_a [64]Piece, field [2]uint16) int {
	for i := 0; i < len(pieces_a); i++ { //iterate trough the pieces array to find the piece matching the given field
		if pieces_a[i] != nil && pieces_a[i].Give_Pos() == field {
			return i
		}
	}
	return -1
}

func move_is_in_board(move [2]uint16) bool {
	if move[0] <= 7 && move[1] <= 7 { //check if the given move exceeds the board (lesser than zero must not be tested due to the fact that every value is uint)
		return true
	} else {
		return false
	}
}

func Field_Can_Be_Captured(pieces_that_can_capture_are_white bool, field [2]uint16, pieces_a [64]Piece, moves_counter int16) bool {
	//this function checks if it is even possible for the current player to move to a given field
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

	var direction int16 = 0 //direction in which the pawns are moving
	var last_y uint16       //last possible row of squares
	if p.Is_White_Piece() {
		direction = -1
		last_y = 0
	} else {
		direction = 1
		last_y = 7
	}

	if p.position[1] != last_y {
		//einer move nach vorne
		p.try_to_move(pieces_a, [2]uint16{p.position[0], uint16(int16(p.position[1]) + direction)}, 64)
		//schlagen rechts
		p.try_to_take(pieces_a, [2]uint16{p.position[0] + 1, uint16(int16(p.position[1]) + direction)}, true)
		//schlagen links
		p.try_to_take(pieces_a, [2]uint16{p.position[0] - 1, uint16(int16(p.position[1]) + direction)}, true)
		//zweier move
		if p.position[1] != uint16(int16(last_y)-direction) && p.Give_Has_Moved() == -1 && Find_Piece_With_Pos(pieces_a, [2]uint16{p.Give_Pos()[0], uint16(int16(p.Give_Pos()[1]) + direction)}) == -1 { //überprüft nicht ob etwas davor steht
			p.try_to_move(pieces_a, [2]uint16{p.position[0], uint16(int16(p.position[1]) + direction*2)}, 65)
		}
		//enpassant rechts
		p.calc_enpassant(pieces_a, [2]uint16{p.position[0] + 1, uint16(int16(p.position[1]) + direction)}, [2]uint16{p.position[0] + 1, p.position[1]}, moves_counter)
		//enpassant links
		p.calc_enpassant(pieces_a, [2]uint16{p.position[0] - 1, uint16(int16(p.position[1]) + direction)}, [2]uint16{p.position[0] - 1, p.position[1]}, moves_counter)
	}
}

func (p *Knight) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()
	//calcualte the 8 knight moves
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 2, p.Give_Pos()[1] + 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 2, p.Give_Pos()[1] - 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 2, p.Give_Pos()[1] + 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 2, p.Give_Pos()[1] - 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1] + 2})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1] - 2})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1] + 2})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1] - 2})
}

func (p *Bishop) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()
	p.calc_moves_diagonally(pieces_a)
}

func (p *Rook) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()
	p.calc_moves_vertically_and_horizontally(pieces_a)
}

func (p *Queen) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()
	p.calc_moves_vertically_and_horizontally(pieces_a)
	p.calc_moves_diagonally(pieces_a)
}

func (p *King) Calc_Moves(pieces_a [64]Piece, moves_counter int16) {
	p.Clear_Legal_Moves()

	var right_castle_possible uint16 = 64
	var left_castle_possible uint16 = 64
	var blocking_right_castle bool
	var blocking_left_castle bool

	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1]})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1] + 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] + 1, p.Give_Pos()[1] - 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1]})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1] + 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0] - 1, p.Give_Pos()[1] - 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0], p.Give_Pos()[1] + 1})
	p.Calc_Normal_Move(pieces_a, [2]uint16{p.Give_Pos()[0], p.Give_Pos()[1] - 1})

	//castle
	for i := 0; i < len(pieces_a); i++ { //iterate trough pieces to find the rooks in the same color as the king
		if pieces_a[i] != nil {
			if Rook, ok := pieces_a[i].(*Rook); ok {
				if Rook.Is_White_Piece() == p.Is_White_Piece() {
					if uint16(int16(Rook.Give_Pos()[0])-int16(p.Give_Pos()[0])) == 3 && Rook.Give_Has_Moved() == 0 && p.Give_Has_Moved() == 0 { //check if the rook has the right horizontal distance to the king
						right_castle_possible = uint16(i) //set this var to the index of the corrsponding rook
					} else if uint16(int16(p.Give_Pos()[0])-int16(Rook.Give_Pos()[0])) == 4 && Rook.Give_Has_Moved() == 0 && p.Give_Has_Moved() == 0 { //check if the rook has the right horizontal distance to the king
						left_castle_possible = uint16(i) //set this var to the index of the corrsponding rook
					}
				}
			}
			if pieces_a[i].Give_Pos() == [2]uint16{5, p.Give_Pos()[1]} || pieces_a[i].Give_Pos() == [2]uint16{6, p.Give_Pos()[1]} { //check if there are pieces in the way
				blocking_right_castle = true
			}
			if pieces_a[i].Give_Pos() == [2]uint16{1, p.Give_Pos()[1]} || pieces_a[i].Give_Pos() == [2]uint16{2, p.Give_Pos()[1]} || pieces_a[i].Give_Pos() == [2]uint16{3, p.Give_Pos()[1]} { //check if there are pieces in the way
				blocking_left_castle = true
			}
		}
	}

	if !blocking_right_castle && right_castle_possible != 64 {
		//p.Append_Legal_Moves([3]uint16{7, p.Give_Pos()[1], uint16(right_castle_possible)})
		p.Append_Legal_Moves([3]uint16{6, p.Give_Pos()[1], uint16(right_castle_possible)})
	}
	if !blocking_left_castle && left_castle_possible != 64 {
		//p.Append_Legal_Moves([3]uint16{0, p.Give_Pos()[1], uint16(left_castle_possible)})
		p.Append_Legal_Moves([3]uint16{2, p.Give_Pos()[1], uint16(left_castle_possible)})
	}
}

func (p *King) Is_In_Check(pieces_a [64]Piece, moves_counter int16) bool {
	//this function is only for the king obviously --> checks if the field, the king is standing on can be captured
	var field [2]uint16 = p.Give_Pos()

	if Field_Can_Be_Captured(!p.Is_White_Piece(), field, pieces_a, moves_counter) {
		return true
	} else {
		return false
	}
}

func Copy_Array(pieces_a [64]Piece) [64]Piece { //create a deep copy of the given array (needed cuz the pieces array is an array of pointers)
	var copy_of_pieces_a [64]Piece
	for i, current_piece := range pieces_a {
		if current_piece != nil {
			copy_of_pieces_a[i] = current_piece.DeepCopy(current_piece)
		}
	}
	return copy_of_pieces_a
}

func Calc_Moves_With_Check(pieces_a [64]Piece, moves_counter int16, current_king_index int) ([64]Piece, bool) { //calculates the move of all pieces in the pieces array at once
	//the third step in the moves calculation (Step1: check if the move is mathematically possible, step2: check if the move is possible with no pieces standing in the way, step3: check if the move has no check conflicts)
	var current_legal_moves [][3]uint16
	var checkmate bool = true

	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil && pieces_a[i].Is_White_Piece() == pieces_a[current_king_index].Is_White_Piece() { //calculate the moves for each piece matching the color of the current player

			pieces_a[i].Calc_Moves(pieces_a, moves_counter) //includes step one and two of calculating moves

			current_legal_moves = pieces_a[i].Give_Legal_Moves()

			pieces_a[i].Clear_Legal_Moves() //clear the moves again, cuz they will be readded when matching the given conditions

			for k := 0; k < len(current_legal_moves); k++ { //iterates trough the legal moves for the current piece
				temp_pieces_a := Copy_Array(pieces_a)                                                                       //this creates a deep copy of the pieces array --> resets after every legal move cuz you don't wanna mess up the pieces array
				temp_pieces_a, _, _ = Move_Piece_To(temp_pieces_a[i], current_legal_moves[k], moves_counter, temp_pieces_a) //the current legal move is actually executed in a save copy of the pieces array

				if !temp_pieces_a[current_king_index].(*King).Is_In_Check(temp_pieces_a, moves_counter) { //if the move doesen't result in an instant check of the own king

					//if the current piece is the king of the current player and the current move is a castle, it must be checked if the squares between the king and rook can't be captured
					if i == current_king_index && current_legal_moves[k][2] <= 63 {
						rook := pieces_a[current_legal_moves[k][2]].(*Rook)
						if rook.Give_Pos()[0] == 7 && //check if right castle is blocked
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), [2]uint16{rook.Give_Pos()[0] - 1, rook.Give_Pos()[1]}, pieces_a, moves_counter) &&
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), [2]uint16{rook.Give_Pos()[0] - 2, rook.Give_Pos()[1]}, pieces_a, moves_counter) &&
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), pieces_a[current_king_index].Give_Pos(), pieces_a, moves_counter) {
							pieces_a[i].Append_Legal_Moves(current_legal_moves[k])
							checkmate = false
						} else if rook.Give_Pos()[0] == 0 && //check if left castle is blocked
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), [2]uint16{rook.Give_Pos()[0] + 2, rook.Give_Pos()[1]}, pieces_a, moves_counter) &&
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), [2]uint16{rook.Give_Pos()[0] + 3, rook.Give_Pos()[1]}, pieces_a, moves_counter) &&
							!Field_Can_Be_Captured(!rook.Is_White_Piece(), pieces_a[current_king_index].Give_Pos(), pieces_a, moves_counter) {
							pieces_a[i].Append_Legal_Moves(current_legal_moves[k])
							checkmate = false
						} //else: rochade is not possible because one of the involved squares can be captured
					} else { //if the current move is no castle, and the king is not in check after the move, simply add it to the legal moves
						pieces_a[i].Append_Legal_Moves(current_legal_moves[k])
						checkmate = false
					}
				}
			}
		}
	}
	return pieces_a, checkmate
}

func (p *Pawn) calc_enpassant(pieces_a [64]Piece, field, en_passant_pawn_pos [2]uint16, moves_counter int16) {
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil { //get the pawn the is to be captured via en passant
			if en_passant_pawn, ok := pieces_a[i].(*Pawn); ok && pieces_a[i].Is_White_Piece() != p.Is_White_Piece() && pieces_a[i].Give_Pos() == en_passant_pawn_pos {
				//andersfarbiger pawn rechts neben dem pawn --> en passant rechts
				if en_passant_pawn.Give_Has_Moved() > 0 && en_passant_pawn.Give_Has_Moved()+1 == moves_counter { //check if en passant can be executed (the pawn that is tried to be captured must have made a double step move, one tempo ago)
					p.Append_Legal_Moves([3]uint16{field[0], field[1], uint16(i)}) //if so append
				}
			}
		}
	}
}

func (p *chess_object) try_to_take(pieces_a [64]Piece, field [2]uint16, piece_is_king_or_pawn bool) {
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil { //simply check if there is a piece on the given field
			if pieces_a[i].Is_White_Piece() != p.Is_White_Piece() && pieces_a[i].Give_Pos() == field {
				if piece_is_king_or_pawn { //the third element of kings and pawns works diferently as every other piece, therfore, this condition must be checked
					p.Append_Legal_Moves([3]uint16{field[0], field[1], uint16(i + 66)})
				} else {
					p.Append_Legal_Moves([3]uint16{field[0], field[1], uint16(i)})
				}

			}
		}
	}
}

func (p *chess_object) try_to_move(pieces_a [64]Piece, field [2]uint16, status uint16) {
	var blocking_piece bool
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil { //if there is a piece on the corresponding square
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

func (p *chess_object) check_if_piece_is_blocking(pieces_a [64]Piece, current_pos [2]uint16) bool {
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
		p.Append_Legal_Moves([3]uint16{current_pos[0], current_pos[1], blocking_piece_index}) //the move will be addded as a take move
		var_break = true
	} else if blocking_piece.Is_White_Piece() == p.Is_White_Piece() { //es steht etwas im weg, was aber nicht geschlagen werden kann, daher wird sofort gebreaked
		var_break = true
	} else {
		fmt.Println("fatal: Error in Calculating Moves Method")
	}
	return var_break
}

func (p *chess_object) calc_moves_diagonally(pieces_a [64]Piece) {
	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x < 7 && new_y < 7; { //as long as not out of the board, check diagonal moves (up right)
		new_x++
		new_y++
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x < 7 && new_y != 0; { //as long as not out of the board, check diagonal moves (down right)
		new_x++
		new_y--
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x != 0 && new_y < 7; { //as long as not out of the board, check diagonal moves (up left)
		new_x--
		new_y++
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Give_Pos()[0], p.Give_Pos()[1]; new_x != 0 && new_y != 0; { //as long as not out of the board, check diagonal moves (down left)
		new_x--
		new_y--
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}
}

func (p *chess_object) calc_moves_vertically_and_horizontally(pieces_a [64]Piece) {
	for new_x := p.Give_Pos()[0]; new_x < 7; { //as long as not out of the board, check vertical moves (up)
		new_x++
		var current_pos [2]uint16 = [2]uint16{new_x, p.Give_Pos()[1]}

		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}

	for new_x := p.Give_Pos()[0]; new_x != 0; { //as long as not out of the board, check vertical moves (down)
		new_x--
		var current_pos [2]uint16 = [2]uint16{new_x, p.Give_Pos()[1]}
		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
		if new_x == 0 {
			break
		}
	}

	for new_y := p.Give_Pos()[1]; new_y < 7; { //as long as not out of the board, check horzontal moves (right)
		new_y++
		var current_pos [2]uint16 = [2]uint16{p.Give_Pos()[0], new_y}
		if p.check_if_piece_is_blocking(pieces_a, current_pos) {
			break
		}
	}

	for new_y := p.Give_Pos()[1]; new_y != 0; { //as long as not out of the board, check horzontal moves (left)
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
	gfx.LadeBild(0, 0, (path + "\\resources\\images\\Pieces.bmp")) //load the whole image into the window and cut out the given piece afterwards

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
