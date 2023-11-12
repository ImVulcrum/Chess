package pieces

type Piece interface {
	Calc_Moves(pieces_a [64]Piece, moves_counter int16)
	Is_White_Piece() bool
	Give_Legal_Moves() [][3]uint16
	Give_Pos() [2]uint16
	Move_To(new_position [2]uint16)
	Append_Legal_Moves(new_legal_move [3]uint16)
	Clear_Legal_Moves()
	Set_Has_Moved(update int16)
	Give_Has_Moved() int16
	DeepCopy(Piece) Piece
	Give_Piece_Type() string
}
