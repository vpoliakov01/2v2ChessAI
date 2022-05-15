package game

import (
	"github.com/vpoliakov01/2v2ChessAI/set"
)

const (
	BoardSize      = 14
	cornerSize     = 3 // 2v2 chess board has corners (3 x 3) cut out.
	emptySquare    = 0
	inactiveSquare = 1 << 10 // TODO: change value
)

type Board struct {
	b            [BoardSize][BoardSize]int
	PieceSquares map[Player]*set.Set
}

func NewBoard() *Board {
	b := Board{
		b:            [BoardSize][BoardSize]int{},
		PieceSquares: map[Player]*set.Set{},
	}

	for player := 0; player < 4; player++ {
		b.PieceSquares[Player(player)] = set.New()
	}

	for rank := 0; rank < BoardSize; rank++ {
		b.b[rank] = [BoardSize]int{}

		for file := 0; file < BoardSize; file++ {
			if !IsSquareValid(rank, file) {
				b.b[rank][file] = inactiveSquare
			}
		}
	}

	return &b
}

func (b *Board) Get(s Square) int {
	return b.b[s.Rank][s.File]
}

func (b *Board) IsEmpty(s Square) bool {
	return b.b[s.Rank][s.File] == emptySquare
}

func (b *Board) SetStartingPosition() {
	pieces := [][]PieceKind{
		{pawn, pawn, pawn, pawn, pawn, pawn, pawn, pawn},
		{rook, knight, bishop, queen, king, bishop, knight, rook},
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
				b.b[rank][file] = int(NewPiece(player, kind))
				b.PieceSquares[player].Add(NewSquare(rank, file))
			}
		}
	}
}

func (b *Board) Copy() *Board {
	board := *b
	board.PieceSquares = map[Player]*set.Set{}

	for player := range b.PieceSquares {
		board.PieceSquares[player] = b.PieceSquares[player].Copy()
	}

	return &board
}

func (b *Board) Move(move Move) {
	if !b.IsEmpty(move.To) {
		capturedPiece := Piece(b.Get(move.To))
		opponent := capturedPiece.GetPlayer()

		b.PieceSquares[opponent].Delete(move.To)
	}

	player := Piece(b.Get(move.From)).GetPlayer()

	b.PieceSquares[player].Delete(move.From)
	b.PieceSquares[player].Add(move.To)

	b.b[move.To.Rank][move.To.File] = b.b[move.From.Rank][move.From.File]
	b.b[move.From.Rank][move.From.File] = emptySquare
}
