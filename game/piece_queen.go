package game

type Queen Piece

var _ GamePiece = (*Queen)(nil)

func (p Queen) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	return GetDirectionalMoves(board, from, dirs)
}

func (p Queen) GetStrength(board *Board, square Square, piecesLeft int) float64 {
	coef := 0.5 + GetCenterBonus(square)
	return QueenStrength * coef
}
