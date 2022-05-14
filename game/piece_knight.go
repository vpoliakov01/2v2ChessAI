package game

type Knight Piece

var _ MovablePiece = (*Knight)(nil)

func (p Knight) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {-1, 2}, {-1, -2}}

	return GetEnumeratedMoves(board, from, dirs)
}
