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
	case MessageTypeGetAvailableMoves:
		c.processGetAvailableMoves()
	case MessageTypePlayerMove:
		move, err := CastData[PGNMove](msg.Data)
		if err != nil {
			log.Printf("Error casting move: %v", err)
			return
		}
		c.processPlayerMove(move)
	case MessageTypeSaveGame:
		c.processSaveGame()
	case MessageTypeLoadGame:
		c.processLoadGame(msg.Data.(string))
	case MessageTypeNewGame:
		c.processNewGame()
	case MessageTypeSetCurrentMove:
		c.processSetCurrentMove(int(msg.Data.(float64)))
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func (c *Connection) processGetAvailableMoves() {
	gameMoves := c.gs.GetMoves().Flatten()
	moves := make([]PGNMove, len(gameMoves))
	for i, gameMove := range gameMoves {
		moves[i] = PGNMoveFromGameMove(gameMove)
	}
	c.SendMessage(MessageTypeAvailableMoves, moves)
}

func (c *Connection) processPlayerMove(move PGNMove) {
	game := c.gs

	gameMove := GameMoveFromPGN(move)
	if err := game.ValidateMove(&gameMove); err != nil {
		c.SendMessage(MessageTypeInvalidMove, err.Error())
		return
	}
	game.Play(gameMove)
	game.Board.Draw()

	c.proceedUntilPlayerMove()
}

func (c *Connection) processSaveGame() {
	c.SendMessage(MessageTypeSaveGameResponse, SaveGameResponse{
		PGN: c.gs.PGN(),
	})
}

func (c *Connection) processLoadGame(data string) {
	game, err := g.LoadPGN(data)
	if err != nil {
		log.Printf("Error loading game: %v", err)
		return
	}
	c.gs = game

	c.SendMessage(MessageTypeLoadGameResponse, LoadGameResponse{
		PastMoves:   PGNMovesFromGameMoves(c.gs.PastMoves),
		CurrentMove: c.gs.CurrentMove,
	})
	c.proceedUntilPlayerMove()
}

func (c *Connection) processNewGame() {
	c.gs = g.NewGameSession()
	c.SendMessage(MessageTypeLoadGameResponse, LoadGameResponse{
		PastMoves:   PGNMovesFromGameMoves(c.gs.PastMoves),
		CurrentMove: c.gs.CurrentMove,
	})
	c.proceedUntilPlayerMove()
}

func (c *Connection) processSetCurrentMove(moveIndex int) {
	err := c.gs.SetCurrentMove(moveIndex)
	if err != nil {
		log.Printf("Error setting current move: %v", err)
		return
	}
	c.SendMessage(MessageTypeLoadGameResponse, LoadGameResponse{
		PastMoves:   PGNMovesFromGameMoves(c.gs.PastMoves),
		CurrentMove: c.gs.CurrentMove,
	})
	c.proceedUntilPlayerMove()
}

func (c *Connection) proceedUntilPlayerMove() {
	game := c.gs

	for !slices.Contains(c.cfg.HumanPlayers, game.ActivePlayer) {
		now := time.Now()
		bestMove, score, err := c.engine.GetBestMove(game.Game)
		elapsed := time.Since(now)
		if err != nil {
			log.Printf("Error getting best move: %v", err)
			return
		}
		c.SendMessage(MessageTypeEngineMove, BestMoveResponse{
			Move:        PGNMoveFromGameMove(*bestMove),
			Score:       math.Round(score*float64(game.ActivePlayer.Team())*100) / 100,
			Time:        math.Round(elapsed.Seconds()*100) / 100,
			Evaluations: c.engine.EvalsCount,
		})
		game.Play(*bestMove)
		game.Board.Draw()
	}

	c.processGetAvailableMoves()
}
