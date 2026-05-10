package game

type Queen Piece

var _ PieceType = (*Queen)(nil)

// GetMoves returns the moves the piece can make.
func (p Queen) GetMoves(board *Board, from Square) []Square {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	return GetDirectionalMoves(board, from, dirs)
}

// GetStrength returns an estimate of the piece's strength.
func (p Queen) GetStrength(board *Board, square Square, player Player) float64 {
	return StrengthPrecomputed[KindQueen][square.Rank][square.File]
}
