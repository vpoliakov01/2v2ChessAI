package game

type Knight Piece

var _ GamePiece = (*Knight)(nil)

func (p Knight) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {-1, 2}, {-1, -2}}

	return GetEnumeratedMoves(board, from, dirs)
}

func (p Knight) GetStrength(board *Board, square Square, piecesLeft int) float64 {
	coef := 0.5 + GetBalanceBonus(square)
	return KnightStrength * coef
}
