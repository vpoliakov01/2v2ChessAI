package game

import "math"

const (
	PawnStrength   = 1.0
	KnightStrength = 2.0
	BishopStrength = 4.0
	RookStrength   = 4.0
	QueenStrength  = 10.0

	// In 2v2 chess it's possible to lose the king as opponent1 can
	// move away from the opponent2's piece and open a lane of attack.
	KingStrength = 100.0

	PiecesAtTheStart = 64
)

// GetCenterBonus returns a value between 0 and 1.
// Squares closest to the center of the board produce 1, closest to the edge - 0.
// The coefficients are for scaling the result to (0, 1) range.
func GetCenterBonus(s Square) float64 {
	center := (float64(BoardSize) - 1) / 2
	return 1.1 - (math.Pow(math.Pow(float64(s.Rank)-center, 2)+math.Pow(float64(s.File)-center, 2), 0.5) / 7)
}

// GetEdgeBonus returns a value between 0 and 1.
// Squares closest to the edge of the board produce 1, closest to the center - 0.
func GetEdgeBonus(s Square) float64 {
	return 1 - GetCenterBonus(s)
}

// GetBalanceBonus returns a value between 0 and 1.
// Squares equidistant from the center and the edges produce 1.
// The coefficients are for scaling the result to (0, 1) range.
func GetBalanceBonus(s Square) float64 {
	return (math.Pow(GetCenterBonus(s), 2) + math.Pow(GetEdgeBonus(s), 2) - 1) / -.5
}
