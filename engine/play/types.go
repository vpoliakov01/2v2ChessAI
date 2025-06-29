package play

import (
	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

type MessageType string

const (
	MessageTypeGetBoardState MessageType = "getBoardState"
	MessageTypeBoardState    MessageType = "boardState"
	MessageTypeGetMoves      MessageType = "getMoves"
	MessageTypeMoves         MessageType = "moves"
	MessageTypePlayerMove    MessageType = "playerMove"
	MessageTypeEngineMove    MessageType = "engineMove"
	MessageTypeInvalidMove   MessageType = "invalidMove"
)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type Move struct {
	From Position `json:"from"`
	To   Position `json:"to"`
}

type BestMoveResponse struct {
	Move  Move    `json:"move"`
	Score float64 `json:"score"`
}

func (p Position) ToSquare() game.Square {
	return game.Square{Rank: game.BoardSize - 1 - p.Row, File: p.Col}
}

func PositionFromSquare(square game.Square) Position {
	return Position{Row: game.BoardSize - 1 - square.Rank, Col: square.File}
}

func (m Move) ToGameMove() game.Move {
	return game.Move{
		From: m.From.ToSquare(),
		To:   m.To.ToSquare(),
	}
}

func MoveFromGameMove(gameMove game.Move) Move {
	return Move{
		From: PositionFromSquare(gameMove.From),
		To:   PositionFromSquare(gameMove.To),
	}
}
