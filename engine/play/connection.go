package play

import (
	"github.com/gofiber/websocket/v2"
	"github.com/vpoliakov01/2v2ChessAI/engine/ai"
	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

// Config is the config for the game.
type Config struct {
	Depth        int
	CaptureDepth int
	HumanPlayers []game.Player
	MoveLimit    int    // Number of moves to play before stopping.
	EvalLimit    int    // Max number of evaluations to perform per move.
	Evaluation   bool   // Whether to display the evaluation of the position.
	Load         string // PGN file to load.
}

type Connection struct {
	conn   *websocket.Conn
	cfg    *Config
	gs     *game.GameSession
	engine *ai.AI

	pauseEngine chan struct{}
}

func NewConnection(c *websocket.Conn, cfg *Config) *Connection {
	connection := &Connection{
		conn:        c,
		cfg:         cfg,
		gs:          game.SetupBoard(cfg.Load),
		engine:      ai.New(cfg.Depth, cfg.CaptureDepth, cfg.EvalLimit),
		pauseEngine: make(chan struct{}),
	}

	return connection
}
