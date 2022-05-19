package ai_test

import (
	"fmt"
	"runtime"
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
	engine := New(4)

	startTime := time.Now()

	g := game.New()

	for i := 0; i < 10; i++ {
		// for i := 0; !g.HasEnded(); i++ {
		move, score, err := engine.GetBestMove(g)
		if err != nil {
			if err == ErrGameEnded {
				fmt.Printf("%v: Team %v won!\n", i, g.Score)
			} else {
				fmt.Println(err)
			}
			break
		}

		piece := game.Piece(g.Board.GetPiece(move.From))
		if !g.Board.IsEmpty(move.To) {
			capturedPiece := game.Piece(g.Board.GetPiece(move.To))
			fmt.Printf("%v: %v takes %v after %v\n", i, piece, capturedPiece, move)
		} else {
			fmt.Printf("%v: %v moves %v\n", i, piece, move)
		}

		g.Play(*move)
		g.Board.Draw()
		fmt.Println("Evaluation: ", score)
	}

	fmt.Println(time.Since(startTime))
}

func (s *TestSuite) TestMultithreading() {
	cpus := runtime.NumCPU()
	times := []time.Duration{}

	for i := 1; i <= cpus; i++ {
		runtime.GOMAXPROCS(i)

		startTime := time.Now()
		s.TestGetBestMove()

		times = append(times, time.Since(startTime))
	}

	for i, t := range times {
		fmt.Printf("%v CPU: %v\n", i+1, t)
	}

	s.Require().Less(times[len(times)-1], times[0])
}

func (s *TestSuite) TestEngineDepthsPerformance() {
	moves := 1
	depths := 9

	last := time.Duration(1)
	for depth := 1; depth <= depths; depth++ {
		start := time.Now()
		engine := New(depth)
		g := game.New()

		for i := 0; i < moves; i++ {
			move, _, err := engine.GetBestMove(g)
			if err != nil {
				fmt.Println(err)
				break
			}
			g.Play(*move)
		}

		t := time.Since(start)
		fmt.Printf("Depth: %v\tTime: %v\tPer move: %v\tRatio: %.3fx\n", depth, t, t/time.Duration(moves), float64(t)/float64(last))
		last = t
	}
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

		g.Board.PlacePiece(piece, game.Square{rank, file})
	}

	engine := New(2)
	g.Board.Draw()

	for i := 0; i < 30; i++ {
		move, _, err := engine.GetBestMove(g)
		if err != nil {
			if err == ErrGameEnded {
				fmt.Printf("%v: Team %v won!\n", i, g.Score)
			} else {
				fmt.Println(err)
			}
			break
		}

		fmt.Println(move)

		if !g.Board.IsEmpty(move.To) {
			capturedPiece := game.Piece(g.Board.GetPiece(move.To))
			opponent := capturedPiece.Player()
			piece := game.Piece(g.Board.GetPiece(move.From))
			player := piece.Player()
			fmt.Printf("%v: P%v's %v takes P%v's %v after %v\n", i, player, piece, opponent, capturedPiece, move)
		}

		g.Play(*move)
		g.Board.Draw()
	}
}
