package game

const (
	BoardSize  = 14
	CornerSize = 3 // 2v2 chess board has corners (3 x 3) cut out.
)

// Board represents the chess board.
type Board struct {
	Grid         [BoardSize][BoardSize]Piece    `json:"grid"`
	PieceSquares map[Player]map[Square]struct{} `json:"-"`
}

// NewBoard creates a new board.
func NewBoard() *Board {
	b := Board{
		Grid:         [BoardSize][BoardSize]Piece{}, // Use an array instead of slice for perf optimization.
		PieceSquares: map[Player]map[Square]struct{}{},
	}

	for player := 0; player < 4; player++ {
		b.PieceSquares[Player(player)] = map[Square]struct{}{}
	}

	for rank := 0; rank < BoardSize; rank++ {
		for file := 0; file < BoardSize; file++ {
			if !IsSquareValid(rank, file) {
				b.Grid[rank][file] = Piece(InactiveSquare)
			}
		}
	}

	return &b
}

// GetPiece returns a piece from the square.
// NOTE: it does not check the square's validity.
func (b *Board) GetPiece(s Square) Piece {
	return b.Grid[s.Rank][s.File]
}

// IsEmpty checks if the square is empty (no piece).
// NOTE: it does not check the square's validity.
func (b *Board) IsEmpty(s Square) bool {
	return b.Grid[s.Rank][s.File] == Piece(EmptySquare)
}

// Clear clears all the pieces of the board.
func (b *Board) Clear() {
	newBoard := NewBoard()
	*b = *newBoard
}

// PlacePiece places a piece onto the board.
func (b *Board) PlacePiece(piece Piece, square Square) {
	b.Grid[square.Rank][square.File] = piece
	b.PieceSquares[piece.Player()][square] = struct{}{}
}

// SetPieceSquares sets PieceSquares.
func (b *Board) SetPieceSquares() {
	b.PieceSquares = map[Player]map[Square]struct{}{}
	for player := 0; player < 4; player++ {
		b.PieceSquares[Player(player)] = map[Square]struct{}{}
	}

	for rank := 0; rank < BoardSize; rank++ {
		for file := 0; file < BoardSize; file++ {
			square := Square{rank, file}

			if square.IsValid() && !b.IsEmpty(square) {
				piece := b.GetPiece(square)
				b.PieceSquares[piece.Player()][square] = struct{}{}
			}
		}
	}
}

// SetStartingPosition sets the pieces for 4 players.
func (b *Board) SetStartingPosition() {
	pieces := [][]PieceKind{
		{KindPawn, KindPawn, KindPawn, KindPawn, KindPawn, KindPawn, KindPawn, KindPawn},
		{KindRook, KindKnight, KindBishop, KindQueen, KindKing, KindBishop, KindKnight, KindRook},
	}

	for row := range pieces {
		for col, kind := range pieces[row] {
			playerPositions := [][]int{
				{1 - row, 3 + col},
				{10 - col, 1 - row},
				{12 + row, 10 - col},
				{3 + col, 12 + row},
			}

			for i := range playerPositions {
				player := Player(i)
				rank := playerPositions[i][0]
				file := playerPositions[i][1]
				b.PlacePiece(NewPiece(player, kind), Square{rank, file})
			}
		}
	}
}

// Copy returns a deep copy of the board.
func (b *Board) Copy() *Board {
	board := *b
	board.PieceSquares = map[Player]map[Square]struct{}{}

	for player := range b.PieceSquares {
		copy := map[Square]struct{}{}

		for square := range b.PieceSquares[player] {
			copy[square] = struct{}{}
		}

		board.PieceSquares[player] = copy
	}

	return &board
}

// Move performs a move of a piece on the board.
func (b *Board) Move(move Move) {
	if !b.IsEmpty(move.To) {
		capturedPiece := Piece(b.GetPiece(move.To))
		opponent := capturedPiece.Player()

		delete(b.PieceSquares[opponent], move.To)
	}

	player := Piece(b.GetPiece(move.From)).Player()

	delete(b.PieceSquares[player], move.From)
	b.PieceSquares[player][move.To] = struct{}{}

	b.Grid[move.To.Rank][move.To.File] = b.Grid[move.From.Rank][move.From.File]
	b.Grid[move.From.Rank][move.From.File] = Piece(EmptySquare)
}

// Unmove undoes a move of a piece on the board.
func (b *Board) Unmove(move Move, capturedPiece Piece) {
	b.Grid[move.From.Rank][move.From.File] = b.Grid[move.To.Rank][move.To.File]
	b.Grid[move.To.Rank][move.To.File] = capturedPiece

	player := Piece(b.GetPiece(move.From)).Player()
	b.PieceSquares[player][move.From] = struct{}{}
	delete(b.PieceSquares[player], move.To)

	if !capturedPiece.IsEmpty() {
		b.PieceSquares[capturedPiece.Player()][move.To] = struct{}{}
	}
}
