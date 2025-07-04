package game

import (
	"fmt"
)

// MoveMap represents available moves from each square.
type MoveMap map[Square][]Square

// Game represents a state of the game.
type Game struct {
	ActivePlayer Player
	Board        *Board
	Winner       Team     // Red/Yellow win: 1, Blue/Green win: -1.
	MoveMap      *MoveMap `json:"-"`
}

// New creates a new Game.
func New() *Game {
	g := Game{
		ActivePlayer: 0,
		Board:        NewBoard(),
		Winner:       0,
	}

	g.Board.SetStartingPosition()

	return &g
}

// GetMoves returns all moves for the active player.
func (g *Game) GetMoves() MoveMap {
	if g.MoveMap != nil {
		return *g.MoveMap
	}

	moves := MoveMap{}

	for square := range g.Board.PieceSquares[g.ActivePlayer] {
		piece := Piece(g.Board.GetPiece(square)).PieceType()
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
			g.Winner = g.ActivePlayer.Team()
		}
	}

	g.Board.Move(move)
	g.ActivePlayer = (g.ActivePlayer + 1) % 4
}

// HasKing checks if the player still has a king.
func (g *Game) HasKing(player Player) bool {
	for square := range g.Board.PieceSquares[player] {
		piece := Piece(g.Board.GetPiece(square))
		if piece.Kind() == KindKing {
			return true
		}
	}
	return false
}

// HasEnded returns whether the game has ended.
func (g *Game) HasEnded() bool {
	return g.Winner != 0
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
