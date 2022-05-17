package game

// Game represents a state of the game.
type Game struct {
	ActivePlayer Player
	Board        *Board
	Score        Team // Red/Yellow win: 1, Blue/Green win: -1.
}

// New creates a new Game.
func New() *Game {
	game := Game{
		ActivePlayer: 0,
		Board:        NewBoard(),
		Score:        0,
	}

	game.Board.SetStartingPosition()

	return &game
}

// GetMoves returns all moves for the active player.
// TODO: return map[Square][]Square instead for specific piece moves lookup?
func (g *Game) GetMoves() []Move {
	moves := []Move{}

	for _, s := range g.Board.PieceSquares[g.ActivePlayer].Elements() {
		square := s.(Square)
		piece := Piece(g.Board.GetPiece(square)).GamePiece()
		moves = append(moves, piece.GetMoves(g.Board, square)...)
	}

	return moves
}

// Play plays a move in the game.
func (g *Game) Play(move Move) {
	if !g.Board.IsEmpty(move.To) {
		capturedPiece := Piece(g.Board.GetPiece(move.To))
		if capturedPiece.Kind() == KindKing {
			g.Score = g.ActivePlayer.Team()
		}
	}

	g.Board.Move(move)
	g.ActivePlayer = (g.ActivePlayer + 1) % 4
}

// HasKing checks if the player still has a king.
func (g *Game) HasKing(player Player) bool {
	for _, sq := range g.Board.PieceSquares[player].Elements() {
		square := sq.(Square)
		piece := Piece(g.Board.GetPiece(square))
		if piece.Kind() == KindKing {
			return true
		}
	}
	return false
}

// HasEnded returns whether the game has ended.
func (g *Game) HasEnded() bool {
	return g.Score != 0
}

// Copy returns a deep copy of the game.
func (g *Game) Copy() *Game {
	newGame := *g
	newGame.Board = g.Board.Copy()
	return &newGame
}
