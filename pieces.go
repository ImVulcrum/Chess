package pieces

import "fmt"

type Piece interface {
	Moves(pieces Piece)
	Name()
	Color()
	Legal_Moves() []Position
	Give_Pos() [2]uint16
}

type Position struct {
	position [2]uint16
}

type Legal_Moves struct {
	moves []Position
}

func (p *Position) MoveTo(position [2]uint16) {
	p.position = position
}

func (p *Position) Give_Pos() [2]uint16 {
	return (p.position)
}

func (c *Common_Piece) Give_Legal_Moves() []Position {
	return c.Legal_Moves.moves
}

type Common_Piece struct {
	white bool
	Position
	Legal_Moves
}

type Pawn struct {
	has_moved bool
	Common_Piece
}

type Knight struct {
	Common_Piece
}

type Rook struct {
	has_moved bool
	Common_Piece
}

type Bishop struct {
	Common_Piece
}

type Queen struct {
	Common_Piece
}

type King struct {
	has_moved bool
	Common_Piece
}

func (p Pawn) Moves() {
	fmt.Printf("Moves of Pawn")
}

func (p Knight) Moves() {
	fmt.Printf("Moves of Knight")
}

func (p Rook) Moves() {
	fmt.Printf("Moves of Rook")
}

func (p Bishop) Moves() {
	fmt.Printf("Moves of Bishop")
}

func (p Queen) Moves() {
	fmt.Printf("Moves of Queen")
}

func (p King) Moves() {
	fmt.Printf("Moves of King")
}
