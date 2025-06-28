package game_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

type TestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestJSON() {
	r := s.Require()

	g := game.New()
	bytes, err := g.JSON()
	r.NoError(err)

	g2, err := game.LoadJSON(bytes)
	r.NoError(err)

	r.NoError(err)
	r.Equal(g.Board.Grid, g2.Board.Grid)
	g2.Board.Draw()
}
