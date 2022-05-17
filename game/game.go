package game

import "fmt"

// MoveMap represents
type MoveMap map[Square][]Square

// Game represents a state of the game.
type Game struct {
	ActivePlayer Player
	Board        *Board
	Score        Team // Red/Yellow win: 1, Blue/Green win: -1.
	MoveMap      *MoveMap
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
func (g *Game) GetMoves() MoveMap {
	if g.MoveMap != nil {
		return *g.MoveMap
	}

	moves := MoveMap{}

	for _, s := range g.Board.PieceSquares[g.ActivePlayer].Elements() {
		square := s.(Square)
		piece := Piece(g.Board.GetPiece(square)).GamePiece()
		moves[square] = append(moves[square], piece.GetMoves(g.Board, square)...)
	}

	g.MoveMap = &moves
	return moves
}

// Play plays a move in the game.
func (g *Game) Play(move Move) {
	g.MoveMap = nil // After the move is played, MoveMap has to be recalculated.

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

// ValidateMove validates the move.
func (g *Game) ValidateMove(move *Move) error {
	if move == nil || !move.From.IsValid() || !move.To.IsValid() {
		return fmt.Errorf("move %v is invalid", move)
	}

	moves := g.GetMoves()
	for _, to := range moves[move.From] {
		if to == move.To {
			return nil
		}
	}

	return fmt.Errorf("move %v is not available to the player", move)
}

// Flatten transforms MoveMap into []Move.
func (m MoveMap) Flatten() []Move {
	moves := []Move{}

	for from := range m {
		for _, to := range m[from] {
			moves = append(moves, Move{from, to})
		}
	}

	return moves
}
