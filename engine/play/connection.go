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
	MoveLimit    int
	Evaluation   bool
	Load         string
}

type Connection struct {
	conn   *websocket.Conn
	cfg    *Config
	game   *game.Game
	engine *ai.AI
}

func NewConnection(c *websocket.Conn, cfg *Config) *Connection {
	connection := &Connection{
		conn:   c,
		cfg:    cfg,
		game:   SetupBoard(cfg.Load),
		engine: ai.New(cfg.Depth, cfg.CaptureDepth),
	}

	return connection
}
