package play

import (
	"sync"

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

// MessageWriter is the minimal interface Connection needs from a websocket
// connection. It allows tests to substitute a mock writer.
type MessageWriter interface {
	WriteMessage(messageType int, data []byte) error
}

type Connection struct {
	conn   MessageWriter
	cfg    *Config
	gs     *game.GameSession
	engine *ai.AI

	stopEngineMovesCh chan struct{}
	writeLock         sync.Mutex
}

func NewConnection(c MessageWriter, cfg *Config) *Connection {
	return newConnection(c, cfg)
}

func newConnection(c MessageWriter, cfg *Config) *Connection {
	connection := &Connection{
		conn:              c,
		cfg:               cfg,
		gs:                game.SetupBoard(cfg.Load),
		engine:            ai.New(cfg.Depth, cfg.CaptureDepth, cfg.EvalLimit),
		stopEngineMovesCh: make(chan struct{}, 10),
	}

	return connection
}
