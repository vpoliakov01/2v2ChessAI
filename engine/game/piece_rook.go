package game

type Rook Piece

var _ PieceType = (*Rook)(nil)

// GetMoves returns the moves the piece can make.
func (p Rook) GetMoves(board *Board, from Square) []Square {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

	return GetDirectionalMoves(board, from, dirs)
}

// GetStrength returns an estimate of the piece's strength.
func (p Rook) GetStrength(board *Board, square Square, player Player) float64 {
	if player.Team() == 1 {
		return StrengthPrecomputed[KindRook][square.Rank][square.File]
	} else {
		return StrengthPrecomputed[KindRook][square.File][square.Rank] // Flip for Blue/Green perspective.
	}
}
