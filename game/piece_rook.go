package game

type Rook Piece

var _ MovablePiece = (*Rook)(nil)

func (p Rook) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

	return GetDirectionalMoves(board, from, dirs)
}
