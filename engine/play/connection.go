package play

import (
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/vpoliakov01/2v2ChessAI/engine/ai"
	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

// Config is the config for the game.
type Config struct {
	Depth        int           `json:"depth"`
	CaptureDepth int           `json:"captureDepth"`
	HumanPlayers []game.Player `json:"humanPlayers"`
	MoveLimit    int           `json:"moveLimit"`  // Number of moves to play before stopping.
	EvalLimit    int           `json:"evalLimit"`  // Max number of evaluations to perform per move.
	Evaluation   bool          `json:"evaluation"` // Whether to display the evaluation of the position.
	Load         string        `json:"load"`       // PGN file to load.
}

type Connection struct {
	conn   *websocket.Conn
	cfg    *Config
	gs     *game.GameSession
	engine *ai.AI

	pauseEngine chan struct{}
	writeLock   sync.Mutex
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
