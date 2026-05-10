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
func (p King) GetStrength(board *Board, square Square, player Player) float64 {
	if player.Team() == 1 {
		return StrengthPrecomputed[KindKing][square.Rank][square.File]
	} else {
		return StrengthPrecomputed[KindKing][square.File][square.Rank] // Flip for Blue/Green perspective.
	}
}
