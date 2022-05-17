package game

type Rook Piece

var _ GamePiece = (*Rook)(nil)

// GetMoves returns the moves the piece can make.
func (p Rook) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

	return GetDirectionalMoves(board, from, dirs)
}

// GetStrength returns an estimate of the piece's strength.
func (p Rook) GetStrength(board *Board, square Square, piecesLeft int) float64 {
	moves := len(p.GetMoves(board, square))
	progression := 1 - float64(piecesLeft)/PiecesAtTheStart
	return Strength[KindRook] * CalculateBonusCoef(moves, 2, 20, GetEdgeBonus(square)*(1-progression)+progression)
}
