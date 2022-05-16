package game_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	. "github.com/vpoliakov01/2v2ChessAI/game"
)

type TestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestBoardCopy() {
	r := s.Require()

	b := NewBoard()
	b.SetStartingPosition()
	r.Equal(16, b.PieceSquares[0].Size())

	c := b.Copy()
	r.Equal(16, c.PieceSquares[0].Size())
	c.PieceSquares[0].Clear()
	b.PieceSquares[1].Clear()
	r.Equal(0, c.PieceSquares[0].Size())
	r.Equal(16, b.PieceSquares[0].Size())
	r.Equal(16, c.PieceSquares[1].Size())
	r.Equal(0, b.PieceSquares[1].Size())
}
