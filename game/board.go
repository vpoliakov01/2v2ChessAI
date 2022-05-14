package game

const (
	boardSize      = 14
	cornerSize     = 3 // 2v2 chess board has corners (3 x 3) cut out.
	emptySquare    = 0
	inactiveSquare = 1 << 10 // TODO: change value
)

type Board [boardSize][boardSize]int

func NewBoard() *Board {
	b := Board{}

	for rank := 0; rank < boardSize; rank++ {
		b[rank] = [boardSize]int{}

		for file := 0; file < boardSize; file++ {
			if !IsSquareValid(rank, file) {
				b[rank][file] = inactiveSquare
			}
		}
	}

	return &b
}

func (b *Board) Get(s Square) int {
	return b[s.Rank][s.File]
}

func (b *Board) IsEmpty(s Square) bool {
	return b[s.Rank][s.File] == emptySquare
}

func (b *Board) SetStartingPosition() {
	pieces := [][]PieceKind{
		{pawn, pawn, pawn, pawn, pawn, pawn, pawn, pawn},
		{rook, knight, bishop, queen, king, bishop, knight, rook},
	}

	for row := range pieces {
		for col, kind := range pieces[row] {
			b[1-row][3+col] = int(NewPiece(Player(0), kind))
			b[10-col][1-row] = int(NewPiece(Player(1), kind))
			b[12+row][10-col] = int(NewPiece(Player(2), kind))
			b[3+col][12+row] = int(NewPiece(Player(3), kind))
		}
	}
}

func (b *Board) Copy() Board {
	return *b
}

func (b *Board) Move(m Move) {
	b[m.To.Rank][m.To.File] = b[m.From.Rank][m.From.File]
	b[m.From.Rank][m.From.File] = emptySquare
}
