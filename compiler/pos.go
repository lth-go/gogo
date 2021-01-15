package compiler

// Pos interface provies two functions to get/set the position for expression or statement.
type Pos interface {
	Position() Position
	SetPosition(Position)
}

// Position provides interface to store code locations.
type Position struct {
	Line   int
	Column int
}

// PosBase provies commonly implementations for Pos.
type PosBase struct {
	pos Position
}

// Position return the position of the expression or statement.
func (x *PosBase) Position() Position {
	return x.pos
}

// SetPosition is a function to specify position of the expression or statement.
func (x *PosBase) SetPosition(pos Position) {
	x.pos = pos
}
