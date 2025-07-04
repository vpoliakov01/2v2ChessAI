package play

import (
	"strconv"
	"strings"

	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

type MessageType string

const (
	MessageTypeGetAvailableMoves MessageType = "getAvailableMoves"
	MessageTypeAvailableMoves    MessageType = "availableMoves"
	MessageTypePlayerMove        MessageType = "playerMove"
	MessageTypeEngineMove        MessageType = "engineMove"
	MessageTypeInvalidMove       MessageType = "invalidMove"
	MessageTypeSaveGame          MessageType = "saveGame"
	MessageTypeSaveGameResponse  MessageType = "saveGameResponse"
	MessageTypeLoadGame          MessageType = "loadGame"
	MessageTypeLoadGameResponse  MessageType = "loadGameResponse"
	MessageTypeNewGame           MessageType = "newGame"
	MessageTypeSetCurrentMove    MessageType = "setCurrentMove"
)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

type PGNMove string

type BestMoveResponse struct {
	Move        PGNMove `json:"move"`
	Score       float64 `json:"score"`
	Time        float64 `json:"time"`
	Evaluations int     `json:"evaluations"`
}

type SaveGameResponse struct {
	PGN string `json:"pgn"`
}

type LoadGameResponse struct {
	PastMoves   []PGNMove `json:"pastMoves"`
	CurrentMove int       `json:"currentMove"`
}

func PGNMoveFromGameMove(gameMove game.Move) PGNMove {
	return PGNMove(gameMove.From.String() + "-" + gameMove.To.String())
}

func PGNMovesFromGameMoves(gameMoves []game.Move) []PGNMove {
	moves := make([]PGNMove, len(gameMoves))
	for i, gameMove := range gameMoves {
		moves[i] = PGNMoveFromGameMove(gameMove)
	}
	return moves
}

func GameMoveFromPGN(pgn PGNMove) game.Move {
	pos := strings.Split(string(pgn), "-")
	return game.Move{From: SquareFromPGN(pos[0]), To: SquareFromPGN(pos[1])}
}

// SquareFromPGN returns a square from a pgn string.
func SquareFromPGN(pgn string) game.Square {
	rank, err := strconv.Atoi(string(pgn[1:]))
	if err != nil {
		return game.Square{}
	}
	return game.Square{Rank: rank - 1, File: int(pgn[0] - 'a')}
}
