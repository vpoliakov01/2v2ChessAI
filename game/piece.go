package game

import (
	"fmt"

	"github.com/vpoliakov01/2v2ChessAI/color"
)

const (
	// Store the piece as ppkkk (last 3 bits specify the kind, 2 bits before them specify the player).
	pieceBitOffset = 3
	pieceBitMask   = 7 // 2^4-1.
)

type Piece uint8

type PieceKind uint8

const (
	EmptySquare Piece = iota
	InactiveSquare
	KindPawn PieceKind = iota
	KindKnight
	KindBishop
	KindRook
	KindQueen
	KindKing
)

var (
	printMap = map[PieceKind]string{
		KindPawn:   "♟",
		KindKnight: "♞",
		KindBishop: "♝",
		KindRook:   "♜",
		KindQueen:  "♛",
		KindKing:   "♚",
	}
	colorMap = map[Player]color.Color{
		0: color.Red,
		1: color.Blue,
		2: color.Yellow,
		3: color.Green,
	}
)

// GamePiece defines functionality a piece should implement.
type GamePiece interface {
	// GetMoves returns a list of moves the peice could make.
	GetMoves(board *Board, from Square) []Move
	// GetStrength returns an estimate of the piece's strength depending on its position and # of pieces left on the board.
	GetStrength(board *Board, square Square, piecesLeft int) float64
}

// New creates a new Piece.
func NewPiece(player Player, kind PieceKind) Piece {
	return Piece(int(player)<<pieceBitOffset + int(kind))
}

func (p Piece) GetPlayer() Player {
	return Player(p >> pieceBitOffset)
}

func (p Piece) GetKind() PieceKind {
	return PieceKind(p & pieceBitMask)
}

func (p Piece) String() string {
	switch p {
	case InactiveSquare:
		return "███"
	case EmptySquare:
		return "   "
	default:
		return fmt.Sprintf(" %v%v%v ", colorMap[p.GetPlayer()], printMap[p.GetKind()], color.Reset)
	}
}

func (p Piece) GetGamePiece() GamePiece {
	switch p.GetKind() {
	case KindPawn:
		return Pawn(p)
	case KindKnight:
		return Knight(p)
	case KindBishop:
		return Bishop(p)
	case KindRook:
		return Rook(p)
	case KindQueen:
		return Queen(p)
	case KindKing:
		return King(p)
	default:
		panic("unsupported piece")
	}
}
