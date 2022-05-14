package game

type Queen Piece

var _ MovablePiece = (*Queen)(nil)

func (p Queen) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	return GetDirectionalMoves(board, from, dirs)
}
