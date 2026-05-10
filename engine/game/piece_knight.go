package game

type Knight Piece

var _ PieceType = (*Knight)(nil)

// GetMoves returns the moves the piece can make.
func (p Knight) GetMoves(board *Board, from Square) []Square {
	dirs := [][]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}

	return GetEnumeratedMoves(board, from, dirs)
}

// GetStrength returns an estimate of the piece's strength.
func (p Knight) GetStrength(board *Board, square Square, player Player) float64 {
	return StrengthPrecomputed[KindKnight][square.Rank][square.File]
}
