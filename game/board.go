package game

import (
	"github.com/vpoliakov01/2v2ChessAI/set"
)

const (
	BoardSize  = 14
	CornerSize = 3 // 2v2 chess board has corners (3 x 3) cut out.
)

type Board struct {
	Grid         [BoardSize][BoardSize]Piece
	PieceSquares map[Player]*set.Set
}

func NewBoard() *Board {
	b := Board{
		Grid:         [BoardSize][BoardSize]Piece{},
		PieceSquares: map[Player]*set.Set{},
	}

	for player := 0; player < 4; player++ {
		b.PieceSquares[Player(player)] = set.New()
	}

	for rank := 0; rank < BoardSize; rank++ {
		b.Grid[rank] = [BoardSize]Piece{}

		for file := 0; file < BoardSize; file++ {
			if !IsSquareValid(rank, file) {
				b.Grid[rank][file] = Piece(InactiveSquare)
			}
		}
	}

	return &b
}

func (b *Board) Get(s Square) Piece {
	return b.Grid[s.Rank][s.File]
}

func (b *Board) IsEmpty(s Square) bool {
	return b.Grid[s.Rank][s.File] == Piece(EmptySquare)
}

func (b *Board) Clear() {
	newBoard := NewBoard()
	*b = *newBoard
}

func (b *Board) PlacePiece(piece Piece, square Square) {
	b.Grid[square.Rank][square.File] = piece
	b.PieceSquares[piece.GetPlayer()].Add(square)
}

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
				b.PlacePiece(NewPiece(player, kind), NewSquare(rank, file))
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

	b.Grid[move.To.Rank][move.To.File] = b.Grid[move.From.Rank][move.From.File]
	b.Grid[move.From.Rank][move.From.File] = Piece(EmptySquare)
}
