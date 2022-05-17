package game

type Queen Piece

var _ GamePiece = (*Queen)(nil)

// GetMoves returns the moves the piece can make.
func (p Queen) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	return GetDirectionalMoves(board, from, dirs)
}

// GetStrength returns an estimate of the piece's strength.
func (p Queen) GetStrength(board *Board, square Square, piecesLeft int) float64 {
	moves := len(p.GetMoves(board, square))
	progression := 1 - float64(piecesLeft)/PiecesAtTheStart
	return Strength[KindQueen] * CalculateBonusCoef(moves, 2, 30, GetCenterBonus(square)*(1-progression)+GetBalanceBonus(square)*progression)
}
