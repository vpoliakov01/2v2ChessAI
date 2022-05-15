package game

type Bishop Piece

var _ GamePiece = (*Bishop)(nil)

func (p Bishop) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	return GetDirectionalMoves(board, from, dirs)
}

func (p Bishop) GetStrength(board *Board, square Square, piecesLeft int) float64 {
	coef := 0.5 + GetBalanceBonus(square)
	return BishopStrength * coef
}
