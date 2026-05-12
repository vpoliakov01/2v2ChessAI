package play

import (
	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

type MessageType string

const (
	MessageTypeSetSettings         MessageType = "setSettings"
	MessageTypeSetSettingsResponse MessageType = "setSettingsResponse"
	MessageTypeGetAvailableMoves   MessageType = "getAvailableMoves"
	MessageTypeAvailableMoves      MessageType = "availableMoves"
	MessageTypePlayerMove          MessageType = "playerMove"
	MessageTypeEngineMove          MessageType = "engineMove"
	MessageTypeInvalidMove         MessageType = "invalidMove"
	MessageTypeSaveGame            MessageType = "saveGame"
	MessageTypeSaveGameResponse    MessageType = "saveGameResponse"
	MessageTypeLoadGame            MessageType = "loadGame"
	MessageTypeLoadGameResponse    MessageType = "loadGameResponse"
	MessageTypeNewGame             MessageType = "newGame"
	MessageTypeSetCurrentMove      MessageType = "setCurrentMove"
	MessageTypeGameEnded           MessageType = "gameEnded"
	MessageTypeProcessing          MessageType = "processing"
	MessageTypeStoppedProcessing   MessageType = "stoppedProcessing"
)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

type PGNMove string

type BestMoveResponse struct {
	Continuation []PGNMove `json:"continuation"`
	MoveNumber   int       `json:"moveNumber"`
	Score        float64   `json:"score"`
	Time         float64   `json:"time"`
	Evaluations  int       `json:"evaluations"`
}

type SaveGameResponse struct {
	PGN string `json:"pgn"`
}

type LoadGameResponse struct {
	PastMoves   []PGNMove `json:"pastMoves"`
	CurrentMove int       `json:"currentMove"`
}

type GameEndedResponse struct {
	King   string `json:"king"`
	Winner string `json:"winner"`
}

func GameMoveFromPGN(pgn PGNMove) game.Move {
	return game.MoveFromPGN(string(pgn))
}

func PGNMoveFromGameMove(gameMove game.Move) PGNMove {
	return PGNMove(gameMove.String())
}

func PGNMovesFromGameMoves(gameMoves []game.Move) []PGNMove {
	moves := make([]PGNMove, len(gameMoves))
	for i, gameMove := range gameMoves {
		moves[i] = PGNMoveFromGameMove(gameMove)
	}
	return moves
}
