package game

type King Piece

var _ GamePiece = (*King)(nil)

func (p King) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	moves := GetEnumeratedMoves(board, from, dirs)

	// TODO: add castling.

	return moves
}

func (p King) GetStrength(board *Board, square Square, piecesLeft int) float64 {
	moves := len(p.GetMoves(board, square))
	progression := 1 - float64(piecesLeft)/PiecesAtTheStart
	return Strength[KindKing] * CalculateBonusCoef(moves, 2, 30, GetEdgeBonus(square)*(1-progression)+GetBalanceBonus(square)*progression)
}
