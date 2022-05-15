package ai_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	. "github.com/vpoliakov01/2v2ChessAI/ai"
	"github.com/vpoliakov01/2v2ChessAI/game"
	"github.com/vpoliakov01/2v2ChessAI/set"
)

type TestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestNegamax() {
	g := game.New()
	startTime := time.Now()

	for i := 0; i < 100; i++ {
		move, err := GetBestMove(g)
		if err != nil {
			break
		}

		if !g.Board.IsEmpty(move.To) {
			capturedPiece := game.Piece(g.Board.Get(move.To))
			opponent := capturedPiece.GetPlayer()
			piece := game.Piece(g.Board.Get(move.From))
			player := piece.GetPlayer()
			fmt.Printf("%v: P%v's %v takes P%v's %v after %v\n", i, player, piece, opponent, capturedPiece, move)
			g.Board.Draw()
		}

		g.Play(*move)
		// g.Board.Draw()
		// fmt.Println(bestMove)
		// fmt.Println(EvaluateCurrent(g))
	}

	fmt.Println(time.Since(startTime))
	g.Board.Draw()
	fmt.Println(EvaluateCurrent(g))
}

func (s *TestSuite) TestGetMoves() {
	g := game.New()
	startTime := time.Now()

	for i := 0; i < 1000; i++ {
		if g.HasEnded() {
			break
		}

		moves := g.GetMoves()

		if (len(moves)) == 0 {
			break
		}

		move := moves[rand.Intn(len(moves))]
		g.Play(move)

		fmt.Println(i, move)
		fmt.Println(EvaluateCurrent(g))
	}

	fmt.Println(time.Since(startTime))
	g.Board.Draw()
}

func (s *TestSuite) TestMisc() {
	r := s.Require()
	st := set.New()

	for i := 0; i < 10; i++ {
		a := i
		st.Add(&a)
	}

	r.Equal(10, st.Size())
}
