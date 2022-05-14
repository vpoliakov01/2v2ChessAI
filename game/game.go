package game

import (
	"github.com/vpoliakov01/2v2ChessAI/set"
)

// Game represents a state of the game.
type Game struct {
	ActivePlayer Player
	Board        *Board
	PieceSquares map[Player]*set.Set
}

// New creates a new Game.
func New() *Game {
	game := Game{
		ActivePlayer: 0,
		Board:        NewBoard(),
		PieceSquares: map[Player]*set.Set{},
	}

	for player := 0; player < 5; player++ {
		game.PieceSquares[Player(player)] = set.New()
	}

	game.Board.SetStartingPosition()

	for rank := 0; rank < boardSize; rank++ {
		for file := 0; file < boardSize; file++ {
			v := game.Board[rank][file]
			if v == emptySquare || v == inactiveSquare {
				continue
			}

			game.PieceSquares[Piece(v).GetPlayer()].Add(NewSquare(rank, file))
		}
	}

	return &game
}

func (g *Game) GetMoves() []Move {
	moves := []Move{}

	for _, s := range g.PieceSquares[g.ActivePlayer].Elements() {
		square := s.(Square)
		piece := Piece(g.Board[square.Rank][square.File]).Cast()
		moves = append(moves, piece.GetMoves(g.Board, square)...)
	}

	return moves
}

func (g *Game) Play(move Move) {
	if !g.Board.IsEmpty(move.To) {
		piece := Piece(g.Board.Get(move.To))
		g.PieceSquares[piece.GetPlayer()].Delete(move.To)
	}

	g.PieceSquares[g.ActivePlayer].Delete(move.From)
	g.PieceSquares[g.ActivePlayer].Add(move.To)

	g.Board.Move(move)

	g.ActivePlayer = (g.ActivePlayer + 1) % 4
}
