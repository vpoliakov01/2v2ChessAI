package game

type Rook Piece

var _ GamePiece = (*Rook)(nil)

func (p Rook) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

	return GetDirectionalMoves(board, from, dirs)
}

func (p Rook) GetStrength(board *Board, square Square, piecesLeft int) float64 {
	progression := 1 - float64(piecesLeft)/PiecesAtTheStart
	coef := 0.5 + GetEdgeBonus(square)*(1-progression) + progression
	return RookStrength * coef
}
