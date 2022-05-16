package game

import "fmt"

type Square struct {
	Rank int
	File int
}

// New creates a new Square.
func NewSquare(rank, file int) Square {
	return Square{Rank: rank, File: file}
}

func (s *Square) Add(rank, file int) Square {
	return NewSquare(s.Rank+rank, s.File+file)
}

func (s *Square) IsValid() bool {
	return IsSquareValid(s.Rank, s.File)
}

func (s Square) String() string {
	return fmt.Sprintf("%v%v", fmt.Sprintf("%c", int('A')+s.File), s.Rank+1)
}

func IsSquareValid(rank, file int) bool {
	return !((file < CornerSize || file >= BoardSize-CornerSize) && (rank < CornerSize || rank >= BoardSize-CornerSize)) &&
		(file >= 0 && file < BoardSize && rank >= 0 && rank < BoardSize)
}
