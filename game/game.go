package game

// Game represents a state of the game.
type Game struct {
	ActivePlayer Player
	Board        *Board
	Score        Side // 1: white win, 0: game is ongoing, -1: black win.
}

// New creates a new Game.
func New() *Game {
	game := Game{
		ActivePlayer: 0,
		Board:        NewBoard(),
		Score:        0,
	}

	game.Board.SetStartingPosition()

	for rank := 0; rank < BoardSize; rank++ {
		for file := 0; file < BoardSize; file++ {
			v := game.Board.Get(NewSquare(rank, file))
			if v == emptySquare || v == inactiveSquare {
				continue
			}
		}
	}

	return &game
}

func (g *Game) GetMoves() []Move {
	moves := []Move{}

	for _, s := range g.Board.PieceSquares[g.ActivePlayer].Elements() {
		square := s.(Square)
		piece := Piece(g.Board.Get(square)).GetGamePiece()
		moves = append(moves, piece.GetMoves(g.Board, square)...)
	}

	return moves
}

func (g *Game) Play(move Move) {
	if !g.Board.IsEmpty(move.To) {
		capturedPiece := Piece(g.Board.Get(move.To))
		if capturedPiece.GetKind() == king {
			g.Score = g.ActivePlayer.GetTeam().Side()
		}
	}

	g.Board.Move(move)
	g.ActivePlayer = (g.ActivePlayer + 1) % 4
}

func (g *Game) HasKing(player Player) bool {
	for _, sq := range g.Board.PieceSquares[player].Elements() {
		square := sq.(Square)
		piece := Piece(g.Board.Get(square))
		if piece.GetKind() == king {
			return true
		}
	}
	return false
}

func (g *Game) HasEnded() bool {
	return g.Score != 0
}

func (g *Game) Copy() *Game {
	newGame := *g
	newGame.Board = g.Board.Copy()
	return &newGame
}
