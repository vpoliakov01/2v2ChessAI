package play

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"slices"
	"time"

	g "github.com/vpoliakov01/2v2ChessAI/engine/game"
)

// CastData marshals data and unmarshals it into the given type.
func CastData[T any](data interface{}) (T, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return *new(T), fmt.Errorf("error marshalling data: %v", err)
	}
	var unmarshalled T
	err = json.Unmarshal(bytes, &unmarshalled)
	if err != nil {
		return *new(T), fmt.Errorf("error unmarshalling data: %v", err)
	}
	return unmarshalled, nil
}

func (c *Connection) ProcessMessage(msg *Message) {
	switch msg.Type {
	case MessageTypeGetBoardState:
		boardState, err := c.game.JSON()
		if err != nil {
			log.Fatalf("Error getting board state: %v", err)
		}
		c.SendMessage(MessageTypeBoardState, string(boardState))
	case MessageTypeGetMoves:
		c.processGetMoves()
	case MessageTypePlayerMove:
		move, err := CastData[Move](msg.Data)
		if err != nil {
			log.Printf("Error casting move: %v", err)
			return
		}
		c.processPlayerMove(move)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func (c *Connection) processGetMoves() {
	gameMoves := c.game.GetMoves().Flatten()
	moves := make([]Move, len(gameMoves))
	for i, gameMove := range gameMoves {
		moves[i] = MoveFromGameMove(gameMove)
	}
	c.SendMessage(MessageTypeMoves, moves)
}

func (c *Connection) processPlayerMove(move Move) {
	game := c.game

	gameMove := g.Move{
		From: move.From.ToSquare(),
		To:   move.To.ToSquare(),
	}
	if err := game.ValidateMove(&gameMove); err != nil {
		c.SendMessage(MessageTypeInvalidMove, err.Error())
		return
	}
	game.Play(gameMove)

	c.proceedUntilPlayerMove()
}

func (c *Connection) proceedUntilPlayerMove() {
	game := c.game

	for !slices.Contains(c.cfg.HumanPlayers, game.ActivePlayer) {
		now := time.Now()
		bestMove, score, err := c.engine.GetBestMove(game)
		elapsed := time.Since(now)
		if err != nil {
			log.Printf("Error getting best move: %v", err)
			return
		}
		c.SendMessage(MessageTypeEngineMove, BestMoveResponse{
			Move:        MoveFromGameMove(*bestMove),
			Score:       math.Round(score*float64(game.ActivePlayer.Team())*100) / 100,
			Time:        math.Round(elapsed.Seconds()*100) / 100,
			Evaluations: c.engine.EvalsCount,
		})
		game.Play(*bestMove)
		game.Board.Draw()
	}

	c.processGetMoves()
}
