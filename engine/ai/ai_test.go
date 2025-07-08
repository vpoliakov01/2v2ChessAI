package ai_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	. "github.com/vpoliakov01/2v2ChessAI/engine/ai"
	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

type TestSuite struct {
	suite.Suite
	engine *AI
	game   *game.Game
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupTest() {
	s.engine = New(3, 6, 0)
	game, err := game.LoadPGN(`
1. h2-h3 b9-c9 i13-i12 m8-l8
2. g1-j4 a8-d11 e13-e12 m5-l5
3. e2-e3 d11-g8 h14-k11 n7-l9
`)
	s.Require().NoError(err)
	s.game = game.Game
}

func (s *TestSuite) TestGetBestMove() {
	engine := s.engine

	startTime := time.Now()

	g := s.game.Copy()

	for i := 0; i < 5; i++ {
		move, score, err := engine.GetBestMove(g)
		if err != nil {
			if err == ErrGameEnded {
				fmt.Printf("%v: Team %v won!\n", i, g.Winner)
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

	fmt.Println("Depth: ", engine.Depth)
	fmt.Println("Capture depth: ", engine.CaptureDepth)
	fmt.Println(time.Since(startTime))
}

func (s *TestSuite) TestEngineDepthsPerformance() {
	moves := 1
	cpus := 1
	runtime.GOMAXPROCS(cpus)

	depths := []struct {
		depth        int
		captureDepth int
	}{
		// {2, 5},
		// {2, 6},
		// {2, 7},
		// {2, 8},
		// {3, 5},
		// {3, 6},
		// {3, 7},
		// {3, 8},
		// {3, 9},
		// {4, 5},
		// {4, 6},
		// {4, 7},
		// {4, 8},
		// {5, 5},
		// {5, 6},
		{5, 7},
	}

	last := time.Duration(1)
	for _, d := range depths {
		start := time.Now()
		engine := New(d.depth, d.captureDepth, 0)
		g := s.game.Copy()

		for i := 0; i < moves; i++ {
			move, _, err := engine.GetBestMove(g)
			if err != nil {
				fmt.Println(err)
				break
			}
			g.Play(*move)
		}

		t := time.Since(start)
		fmt.Printf("Depth: %v/%v\t Time: %.2fs\tPer move: %.2fs\tRatio: %.2fx\tEvaluations: %v\n", d.depth, d.captureDepth, t.Seconds(), t.Seconds()/float64(moves), float64(t)/float64(last), engine.EvalsCount)
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

	engine := New(2, 2, 0)
	g.Board.Draw()

	for i := 0; i < 30; i++ {
		move, _, err := engine.GetBestMove(g)
		if err != nil {
			if err == ErrGameEnded {
				fmt.Printf("%v: Team %v won!\n", i, g.Winner)
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

func (s *TestSuite) TestMultithreading() {
	cpus := runtime.NumCPU()
	times := []time.Duration{}
	moves := 1
	engine := New(2, 4, 0)

	for i := 1; i <= cpus; i++ {
		runtime.GOMAXPROCS(i)

		startTime := time.Now()

		g := s.game.Copy()

		for i := 0; i < moves; i++ {
			move, _, err := engine.GetBestMove(g)
			if err != nil {
				fmt.Println(err)
				break
			}
			g.Play(*move)
		}

		times = append(times, time.Since(startTime))
		fmt.Printf("%v CPU: %v\n", i, times[i-1])
	}

	s.Require().Less(times[len(times)-1], times[0])
}

func (s *TestSuite) TestObviousMoves() {
	runtime.GOMAXPROCS(1)

	type testCase struct {
		name string
		pgn  string
		move string
	}

	testCases := []testCase{
		{
			"free queen",
			`
1. f2-f3 b7-c7 d13-d12 m7-l7
2. g1-b6`,
			"a7-b6",
		},
		{
			"mate",
			`1. f2-f3 b6-c6 g13-g12 m8-l8`,
			"g1-a7",
		},
	}
	failures := 0
	for _, tc := range testCases {
		gs, err := game.LoadPGN(tc.pgn)
		s.Require().NoError(err)
		g := gs.Game
		engine := s.engine

		move, _, err := engine.GetBestMove(g)
		s.Require().NoError(err)
		if move.String() != tc.move {
			fmt.Printf("%v: %v != %v\n", tc.name, move, tc.move)
			failures++
			g.Board.Draw()
		}
	}
	s.Require().Equal(0, failures)
}
