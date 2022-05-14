package game

type Pawn Piece

var _ MovablePiece = (*Pawn)(nil)

func (p Pawn) GetMoves(board *Board, from Square) []Move {
	player := Piece(board.Get(from)).GetPlayer()

	// Directions depend on which player's pawn it is.
	dirs := [][][]int{
		{{1, 0}, {1, -1}, {1, 1}},
		{{0, 1}, {-1, 1}, {1, 1}},
		{{-1, 0}, {-1, -1}, {-1, 1}},
		{{0, -1}, {-1, -1}, {1, -1}},
	}[player]
	moves := []Move{}

	// Move forward by 1.
	to := from.Add(dirs[0][0], dirs[0][1])
	if board.IsEmpty(to) {
		moves = append(moves, Move{from, to})

		// Move forward by 2.
		to = from.Add(2*dirs[0][0], 2*dirs[0][1])
		if board.IsEmpty(to) {
			// Since pawns capture sideways, they can end up on other players' pawns' starting positions.
			// Therefore, it's not enough to just check if the pawn is on rank 1, file 1, etc.
			if (player == 0 && from.Rank == 1) ||
				(player == 1 && from.File == 1) ||
				(player == 2 && from.Rank == boardSize-2) ||
				(player == 3 && from.File == boardSize-2) {
				moves = append(moves, Move{from, to})
			}
		}
	}

	// Capture.
	for i := 1; i <= 2; i++ {
		dir := dirs[i]
		to := from.Add(dir[0], dir[1])

		if !to.IsValid() {
			continue
		} else if !board.IsEmpty(to) && !Piece(board.Get(to)).GetPlayer().IsTeamMate(Piece(board.Get(from)).GetPlayer()) {
			moves = append(moves, Move{from, to})
		}
	}

	// TODO: add en passant.

	return moves
}
