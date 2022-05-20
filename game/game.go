package game

import (
	"encoding/json"
	"fmt"
)

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
	g := Game{
		ActivePlayer: 0,
		Board:        NewBoard(),
		Score:        0,
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

// JSON returns json of the game object.
func (g *Game) JSON() ([]byte, error) {
	copy := g.Copy()
	copy.Board.PieceSquares = nil // Exclude pieceSquares due to issues with marshaling map[Square]struct{}.
	return json.Marshal(copy)
}

// LoadJSON returns the game defined by the json.
func LoadJSON(bytes []byte) (*Game, error) {
	g := Game{}

	err := json.Unmarshal(bytes, &g)
	if err != nil {
		return nil, err
	}

	g.Board.SetPieceSquares()

	return &g, nil
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
