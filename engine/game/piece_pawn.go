package game

type Pawn Piece

var _ PieceType = (*Pawn)(nil)

var (
	pawnMoveDirs = [4][3][2]int{
		{{1, 0}, {1, -1}, {1, 1}},
		{{0, 1}, {-1, 1}, {1, 1}},
		{{-1, 0}, {-1, -1}, {-1, 1}},
		{{0, -1}, {-1, -1}, {1, -1}},
	}
	pawnCaptureDirs = [4][2][2]int{
		{{1, -1}, {1, 1}},
		{{-1, 1}, {1, 1}},
		{{-1, -1}, {-1, 1}},
		{{-1, -1}, {1, -1}},
	}
)

// GetMoves returns the moves the piece can make.
func (p Pawn) GetMoves(board *Board, from Square) []Square {
	player := Piece(board.GetPiece(from)).Player()

	// Move directions depend on which player's pawn it is.
	dirs := pawnMoveDirs[player]
	moves := []Square{}

	// Move forward by 1.
	to := from.Add(dirs[0][0], dirs[0][1])
	if to.IsValid() && board.IsEmpty(to) {
		moves = append(moves, to)

		// Move forward by 2.
		to = from.Add(2*dirs[0][0], 2*dirs[0][1])
		if to.IsValid() && board.IsEmpty(to) {
			// Since pawns capture sideways, they can end up on other players' pawns' starting positions.
			// Therefore, it's not enough to just check if the pawn is on rank 1, file 1, etc.
			if (player == 0 && from.Rank == 1) ||
				(player == 1 && from.File == 1) ||
				(player == 2 && from.Rank == BoardSize-2) ||
				(player == 3 && from.File == BoardSize-2) {
				moves = append(moves, to)
			}
		}
	}

	// Capture.
	for i := 1; i <= 2; i++ {
		dir := dirs[i]
		to := from.Add(dir[0], dir[1])

		if !to.IsValid() {
			continue
		} else if !board.IsEmpty(to) && !Piece(board.GetPiece(to)).Player().IsTeamMate(Piece(board.GetPiece(from)).Player()) {
			moves = append(moves, to)
		}
	}

	// TODO: add en passant and promotions.

	return moves
}

// GetStrength returns an estimate of the piece's strength.
func (p Pawn) GetStrength(board *Board, square Square, player Player) float64 {
	// Check pawn structure.
	dirs := pawnCaptureDirs[player]

	coef := 0.9
	for _, dir := range dirs {
		inFront := square.Add(dir[0], dir[1])
		if !inFront.IsValid() || board.IsEmpty(inFront) {
			continue
		}

		coef += 0.2
	}

	return Strength[KindPawn] * coef
}
