package game

// GetDirectionalMoves returns a list of possible moves in the specified directions.
// (Used for queens, rooks, and bishops.)
func GetDirectionalMoves(board *Board, from Square, vectors [][]int) []Square {
	toSquares := []Square{}

	for _, vector := range vectors {
		for dist := 1; ; dist++ {
			to := from.Add(dist*vector[0], dist*vector[1])

			if !to.IsValid() {
				break
			} else if board.IsEmpty(to) {
				toSquares = append(toSquares, to)
				continue
			} else if !Piece(board.GetPiece(to)).Player().IsTeamMate(Piece(board.GetPiece(from)).Player()) {
				toSquares = append(toSquares, to)
			}
			break
		}
	}

	return toSquares
}

// GetEnumeratedMoves returns a list of possible produced by adding the specified vectors.
// (Used for kings and knights.)
func GetEnumeratedMoves(board *Board, from Square, vectors [][]int) []Square {
	toSquares := []Square{}

	for _, vector := range vectors {
		to := from.Add(vector[0], vector[1])

		if !to.IsValid() {
			continue
		} else if board.IsEmpty(to) || !Piece(board.GetPiece(to)).Player().IsTeamMate(Piece(board.GetPiece(from)).Player()) {
			toSquares = append(toSquares, to)
		}
	}

	return toSquares
}
