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
	return fmt.Sprintf("%v%v", string(int('A')+s.File), s.Rank+1)
}

func IsSquareValid(rank, file int) bool {
	return !((file < cornerSize || file >= boardSize-cornerSize) && (rank < cornerSize || rank >= boardSize-cornerSize)) &&
		(file >= 0 && file < boardSize && rank >= 0 && rank < boardSize)
}
