package game

type King Piece

var _ PieceType = (*King)(nil)

// GetMoves returns the moves the piece can make.
func (p King) GetMoves(board *Board, from Square) []Square {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	moves := GetEnumeratedMoves(board, from, dirs)

	// TODO: add castling.

	return moves
}

// GetStrength returns an estimate of the piece's strength.
func (p King) GetStrength(board *Board, numMoves int, square Square, piecesLeft int) float64 {
	progression := 1 - float64(piecesLeft)/PiecesAtTheStart
	return Strength[KindKing] * CalculateBonusCoef(numMoves, 2, 30, GetEdgeBonus(square)*(1-progression)+GetBalanceBonus(square)*progression)
}
