package game

// GetDirectionalMoves appends valid moves in the given directions (used for queens, rooks, bishops)
// to dst and returns the extended slice. Each direction is followed until the edge of the board
// or a blocking piece is reached.
func GetDirectionalMoves(board *Board, from Square, vectors [][2]int, dst []Square) []Square {
	fromPlayer := Piece(board.GetPiece(from)).Player()

	for _, vector := range vectors {
		for dist := 1; ; dist++ {
			to := from.Add(dist*vector[0], dist*vector[1])

			if !to.IsValid() {
				break
			} else if board.IsEmpty(to) {
				dst = append(dst, to)
				continue
			} else if !Piece(board.GetPiece(to)).Player().IsTeamMate(fromPlayer) {
				dst = append(dst, to)
			}
			break
		}
	}

	return dst
}

// GetEnumeratedMoves appends valid moves produced by adding each vector to from (used for kings
// and knights) to dst and returns the extended slice.
func GetEnumeratedMoves(board *Board, from Square, vectors [][2]int, dst []Square) []Square {
	fromPlayer := Piece(board.GetPiece(from)).Player()

	for _, vector := range vectors {
		to := from.Add(vector[0], vector[1])

		if !to.IsValid() {
			continue
		} else if board.IsEmpty(to) || !Piece(board.GetPiece(to)).Player().IsTeamMate(fromPlayer) {
			dst = append(dst, to)
		}
	}

	return dst
}
