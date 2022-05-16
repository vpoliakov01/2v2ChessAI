package ai_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	. "github.com/vpoliakov01/2v2ChessAI/ai"
	"github.com/vpoliakov01/2v2ChessAI/game"
)

type TestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestGetBestMove() {
	ai := New(3)

	startTime := time.Now()

	g := game.New()

	for i := 0; i < 20; i++ {
		move, err := ai.GetBestMove(g)
		if err != nil {
			if err == ErrGameEnded {
				fmt.Printf("%v: Team %v won!\n", i, g.Score.Team())
			} else {
				fmt.Println(err)
			}
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
		g.Board.Draw()
	}

	g.Board.Draw()
	fmt.Println(ai.EvaluateCurrent(g))
	fmt.Println(time.Since(startTime))
}

func (s *TestSuite) TestGetMoves() {
	ai := New(24)

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
		fmt.Println(ai.EvaluateCurrent(g))
	}

	fmt.Println(time.Since(startTime))
	g.Board.Draw()
}

func (s *TestSuite) TestPosition() {
	pieces := [][]int{
		{int(game.NewPiece(0, game.KindKing)), 13, 10},
		{int(game.NewPiece(0, game.KindPawn)), 13, 9},
		{int(game.NewPiece(0, game.KindPawn)), 12, 10},
		{int(game.NewPiece(0, game.KindPawn)), 12, 9},
		{int(game.NewPiece(1, game.KindKing)), 6, 1},
		{int(game.NewPiece(2, game.KindKing)), 12, 6},
		{int(game.NewPiece(3, game.KindKing)), 8, 13},
		{int(game.NewPiece(2, game.KindQueen)), 9, 12},
		{int(game.NewPiece(0, game.KindQueen)), 10, 13},
	}

	g := game.New()
	g.Board.Clear()

	for i := range pieces {
		piece := game.Piece(pieces[i][0])
		rank := pieces[i][1]
		file := pieces[i][2]

		g.Board.PlacePiece(piece, game.NewSquare(rank, file))
	}

	ai := New(2)
	g.Board.Draw()

	for i := 0; i < 30; i++ {
		move, err := ai.GetBestMove(g)
		fmt.Println(move)
		if err != nil {
			if err == ErrGameEnded {
				fmt.Printf("%v: Team %v won!\n", i, g.Score.Team())
			} else {
				fmt.Println(err)
			}
			break
		}

		if !g.Board.IsEmpty(move.To) {
			capturedPiece := game.Piece(g.Board.Get(move.To))
			opponent := capturedPiece.GetPlayer()
			piece := game.Piece(g.Board.Get(move.From))
			player := piece.GetPlayer()
			fmt.Printf("%v: P%v's %v takes P%v's %v after %v\n", i, player, piece, opponent, capturedPiece, move)
		}

		g.Play(*move)
		g.Board.Draw()
	}
}
