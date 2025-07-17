package game

import "math"

const (
	PiecesAtTheStart = 64
)

// Strength stores relative base strengths for pieces.
var Strength = map[PieceKind]float64{
	KindPawn:   1.0,
	KindKnight: 2.2,
	KindBishop: 5.0,
	KindRook:   4.5,
	KindQueen:  14.0,
	KindKing:   7.0,
}

// GetCenterBonus returns a value between 0 and 1.
// Squares closest to the center of the board produce 1, closest to the edge - 0.
// The coefficients are for scaling the result to (0, 1) range.
func GetCenterBonus(s Square) float64 {
	return 1 - GetEdgeBonus(s)
}

// GetEdgeBonus returns a value between 0 and 1.
// Squares closest to the edge of the board produce 1, closest to the center - 0.
func GetEdgeBonus(s Square) float64 {
	return ((math.Abs(float64(s.Rank)-6.5) + math.Abs(float64(s.File)-6.5)) - 1) / 9
}

// GetBalanceBonus returns a value between 0 and 1.
// Squares equidistant from the center and the edges produce 1.
// The coefficients are for scaling the result to (0, 1) range.
func GetBalanceBonus(s Square) float64 {
	return (math.Pow(GetCenterBonus(s), 2) + math.Pow(GetEdgeBonus(s), 2) - 1.01) / -.5
}

// GetDefenseBonus returns a value between 0 and 1.
// Squares equidistant from the center and the edges produce 1.
// The coefficients are for scaling the result to (0, 1) range.
func GetDefenseBonus(s Square, team Team) float64 {
	switch team {
	case 1:
		return (math.Abs(float64(s.Rank)-6.5) - 0.5) / 6
	case -1:
		return (math.Abs(float64(s.File)-6.5) - 0.5) / 6
	}
	return 0
}

// GetAttackBonus returns a value between 0 and 1.
// Squares equidistant from the center and the edges produce 1.
// The coefficients are for scaling the result to (0, 1) range.
func GetAttackBonus(s Square, team Team) float64 {
	return 1 - GetDefenseBonus(s, team)
}

// GetProximityBonus returns a value between 0 and 1.
// The ideal square is assumed to be placed in all
// symmetric positions around the center of the board
// in regard to the ideal square.
func GetProximityBonus(s Square, idealSquare Square) float64 {
	rankFromCenter := math.Abs(float64(s.Rank) - 6.5)
	fileFromCenter := math.Abs(float64(s.File) - 6.5)

	idealRankFromCenter := math.Abs(float64(idealSquare.Rank) - 6.5)
	idealFileFromCenter := math.Abs(float64(idealSquare.File) - 6.5)

	sFromDiagonal := math.Abs(rankFromCenter - fileFromCenter)
	idealSFromDiagonal := math.Abs(idealRankFromCenter - idealFileFromCenter)

	sFromCenter := rankFromCenter + fileFromCenter
	idealSFromCenter := idealRankFromCenter + idealFileFromCenter

	distFromIdeal := math.Sqrt(math.Pow(sFromCenter-idealSFromCenter, 2) + math.Pow(sFromDiagonal-idealSFromDiagonal, 2))

	return 1 - distFromIdeal/9
}

// CalculateBonusCoef calculates the overall bonus coef. (0.5, 1.5)
func CalculateBonusCoef(moves, movesMin, movesMax int, positionCef float64) float64 {
	return 0.5 + (2*float64(moves-movesMin)/float64(movesMax-movesMin)+positionCef)/2
}
