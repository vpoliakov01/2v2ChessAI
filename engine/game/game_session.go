package game

import "fmt"

// GameSession represents a game with additional metadata.
type GameSession struct {
	*Game
	CurrentMove int
	PastMoves   []Move
}

// NewGameSession creates a new GameSession.
// The abstraction is useful for keeping track of game data without
// convoluting the engine logic with game metadata.
func NewGameSession() *GameSession {
	return &GameSession{
		Game:        New(),
		CurrentMove: -1,
		PastMoves:   []Move{},
	}
}

// Play plays a move in the game session.
func (g *GameSession) Play(move Move) {
	g.PastMoves = g.PastMoves[:g.CurrentMove+1]
	g.Game.Play(move)
	g.PastMoves = append(g.PastMoves, move)
	g.CurrentMove++
}

// SetCurrentMove sets the current move index.
func (g *GameSession) SetCurrentMove(moveIndex int) error {
	if moveIndex < 0 || moveIndex > len(g.PastMoves) {
		return fmt.Errorf("move index out of range")
	}
	g.CurrentMove = moveIndex

	g.Game = New()
	for i := 0; i <= moveIndex; i++ {
		g.Game.Play(g.PastMoves[i])
	}

	return nil
}
