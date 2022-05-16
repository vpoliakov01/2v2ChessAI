package game

// GetDirectionalMoves returns a list of possible moves in the specified directions.
// (Used for queens, rooks, and bishops.)
func GetDirectionalMoves(board *Board, from Square, vectors [][]int) []Move {
	moves := []Move{}

	for _, vector := range vectors {
		for dist := 1; ; dist++ {
			to := from.Add(dist*vector[0], dist*vector[1])

			if !to.IsValid() {
				break
			} else if board.IsEmpty(to) {
				moves = append(moves, Move{from, to})
				continue
			} else if !Piece(board.Get(to)).GetPlayer().IsTeamMate(Piece(board.Get(from)).GetPlayer()) {
				moves = append(moves, Move{from, to})
			}
			break
		}
	}

	return moves
}

// GetEnumeratedMoves returns a list of possible produced by adding the specified vectors.
// (Used for kings and knights.)
func GetEnumeratedMoves(board *Board, from Square, vectors [][]int) []Move {
	moves := []Move{}

	for _, vector := range vectors {
		to := from.Add(vector[0], vector[1])

		if !to.IsValid() {
			continue
		} else if board.IsEmpty(to) || !Piece(board.Get(to)).GetPlayer().IsTeamMate(Piece(board.Get(from)).GetPlayer()) {
			moves = append(moves, Move{from, to})
		}
	}

	return moves
}
