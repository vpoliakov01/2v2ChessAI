package game

import "fmt"

const (
	// Store the piece as ppkkk (last 3 bits specify the kind, 2 bits before them specify the player).
	pieceBitOffset = 3
	pieceBitMask   = 7 // 2^4-1.
)

type PieceKind int

const (
	pawn PieceKind = 1 + iota
	knight
	bishop
	rook
	queen
	king
)

var (
	printMap = map[PieceKind][]string{
		pawn:   {"♟", "♙"},
		knight: {"♞", "♘"},
		bishop: {"♝", "♗"},
		rook:   {"♜", "♖"},
		queen:  {"♛", "♕"},
		king:   {"♚", "♔"},
	}
)

type Piece int

type MovablePiece interface {
	GetMoves(board *Board, from Square) []Move
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

func (p Piece) GetMoves(board *Board, from Square) []Move {
	return nil
}

func (p Piece) String() string {
	mark := " "
	if p.GetPlayer()&2 == 2 {
		mark = "."
	}
	return fmt.Sprintf(" %v%v", printMap[p.GetKind()][p.GetPlayer().GetTeam()], mark)
}

func (p Piece) Cast() MovablePiece {
	switch p.GetKind() {
	case pawn:
		return Pawn(p)
	case knight:
		return Knight(p)
	case bishop:
		return Bishop(p)
	case rook:
		return Rook(p)
	case queen:
		return Queen(p)
	case king:
		return King(p)
	default:
		return p
	}
}
