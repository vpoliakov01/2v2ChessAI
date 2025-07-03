package game

type Bishop Piece

var _ PieceType = (*Bishop)(nil)

// GetMoves returns the moves the piece can make.
func (p Bishop) GetMoves(board *Board, from Square) []Square {
	dirs := [][]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	return GetDirectionalMoves(board, from, dirs)
}

// GetStrength returns an estimate of the piece's strength.
func (p Bishop) GetStrength(board *Board, numMoves int, square Square, piecesLeft int) float64 {
	return Strength[KindBishop] * CalculateBonusCoef(numMoves, 2, 12, GetBalanceBonus(square))
}
