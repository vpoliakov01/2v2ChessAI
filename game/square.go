package game

import "fmt"

// Square stores a coordinate on the board.
type Square struct {
	Rank int
	File int
}

// Add adds a vector to the square.
func (s *Square) Add(rank, file int) Square {
	return Square{s.Rank + rank, s.File + file}
}

// String implements the Stringer interface.
func (s Square) String() string {
	return fmt.Sprintf("%v%v", fmt.Sprintf("%c", int('A')+s.File), s.Rank+1)
}

// IsValid checs if the square is on the board.
func (s *Square) IsValid() bool {
	return IsSquareValid(s.Rank, s.File)
}

// IsSquareValid returns whether rank and file are within [0, 13] and outside the excluded corners.
func IsSquareValid(rank, file int) bool {
	return !((file < CornerSize || file >= BoardSize-CornerSize) && (rank < CornerSize || rank >= BoardSize-CornerSize)) &&
		(file >= 0 && file < BoardSize && rank >= 0 && rank < BoardSize)
}
