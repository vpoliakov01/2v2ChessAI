package game

import "math"

const (
	PiecesAtTheStart = 64
)

var (
	// Strength stores relative base strengths for pieces.
	Strength = map[PieceKind]float64{
		KindPawn:   1.0,
		KindKnight: 2.2,
		KindBishop: 5.0,
		KindRook:   4.5,
		KindQueen:  14.0,
		KindKing:   3.0,
	}
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

// CalculateBonusCoef calculates the overall bonus coef.
func CalculateBonusCoef(moves, movesMin, movesMax int, positionCef float64) float64 {
	return 0.5 + (2*float64(moves-movesMin)/float64(movesMax-movesMin)+positionCef)/3
}
