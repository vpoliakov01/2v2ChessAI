package game

type King Piece

var _ MovablePiece = (*King)(nil)

func (p King) GetMoves(board *Board, from Square) []Move {
	dirs := [][]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	moves := GetEnumeratedMoves(board, from, dirs)

	// TODO: add castling.

	return moves
}
