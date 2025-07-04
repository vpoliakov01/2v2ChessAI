package game

import (
	"fmt"

	"github.com/vpoliakov01/2v2ChessAI/engine/color"
)

const (
	// Store the piece as ppkkk (last 3 bits specify the kind, 2 bits before them specify the player).
	pieceBitOffset = 3
	pieceBitMask   = 7 // 2^4-1.
)

// PieceType defines functionality a piece should implement.
type PieceType interface {
	// GetMoves returns a list of moves the peice could make.
	GetMoves(board *Board, from Square) []Square
	// GetStrength returns an estimate of the piece's strength depending on its position and # of pieces left on the board.
	GetStrength(board *Board, numMoves int, square Square, piecesLeft int) float64
}

type Piece uint8 // Use uint8 to save some space (the board is a [][]Piece).

type PieceKind uint8

const (
	// Set values from 0 to 7.
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

// New creates a new Piece.
func NewPiece(player Player, kind PieceKind) Piece {
	return Piece(int(player)<<pieceBitOffset + int(kind))
}

// Player returns the owner of the piece.
func (p Piece) Player() Player {
	return Player(p >> pieceBitOffset)
}

// Kind returns the kind of the piece.
func (p Piece) Kind() PieceKind {
	return PieceKind(p & pieceBitMask)
}

// IsEmpty returns true if the piece is empty.
func (p Piece) IsEmpty() bool {
	return p == EmptySquare
}

// String implements the Stringer interface.
func (p Piece) String() string {
	switch p {
	case InactiveSquare:
		return "███"
	case EmptySquare:
		return "   "
	default:
		return fmt.Sprintf(" %v%v%v ", colorMap[p.Player()], printMap[p.Kind()], color.Reset)
	}
}

// PieceType returns the corresponding PieceType implementation for the kind of the piece.
func (p Piece) PieceType() PieceType {
	switch p.Kind() {
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
