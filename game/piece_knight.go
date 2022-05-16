package game

type Knight Piece

var _ GamePiece = (*Knight)(nil)

func (p Knight) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {-1, 2}, {-1, -2}}

	return GetEnumeratedMoves(board, from, dirs)
}

func (p Knight) GetStrength(board *Board, square Square, piecesLeft int) float64 {
	moves := len(p.GetMoves(board, square))
	return Strength[KindKnight] * CalculateBonusCoef(moves, 2, 6, GetBalanceBonus(square))
}
