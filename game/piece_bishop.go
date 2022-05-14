package game

type Bishop Piece

var _ MovablePiece = (*Bishop)(nil)

func (p Bishop) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	return GetDirectionalMoves(board, from, dirs)
}
