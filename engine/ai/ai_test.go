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
	games  []*game.Game
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupTest() {
	s.engine = New(3, 6, 0)

	pgns := []string{
		`
1. f2-f3 b6-c6 g13-g12 m8-l8`, // mate in 1
		`
1. f2-f3 b7-c7 d13-d12 m7-l7
2. g1-b6`, // free queen
		`
1. k2-k4 b7-d7 i13-i12 m6-k6
2. f2-f4 a8-b7 g13-g12 m8-l8
3. e1-f3 a10-c9 e14-f12 m10-l10
4. g2-g4 b11-d11 k13-k12 m7-l7`, // Complex real
		`
1. h2-h3 b7-c7 i13-i12 m8-l8`, // Mate in 6, basic beginning trap
		`
1. h2-h3 b9-c9 i13-i12 m8-l8
2. g1-j4 a8-d11 e13-e12 m5-l5
3. e2-e3 d11-g8 h14-k11 n7-l9`, // 4 queens in the middle, bishops ready
		`
1. j2-j3 b5-c5 j14-i12 n5-l6
2. e2-e3 a6-f1 e13-e12 m7-k7
3. g1-f1 a5-c4 j13-j12 n10-l9
4. f1-c4 b7-c7 h13-h12 m5-l5`, // 6/10 engine game
	}
	for _, pgn := range pgns {
		game, err := game.LoadPGN(pgn)
		s.Require().NoError(err)
		s.games = append(s.games, game.Game)
	}
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
		{4, 8},
		// {5, 5},
		// {5, 6},
		// {5, 7},
		// {6, 8},
	}

	last := time.Duration(0)
	totalStart := time.Now()
	for _, g := range s.games {
		for _, d := range depths {
			start := time.Now()
			engine := New(d.depth, d.captureDepth, 0)

			for i := 0; i < moves; i++ {
				_, _, err := engine.GetBestMove(g)
				if err != nil {
					fmt.Println(err)
					break
				}
			}

			t := time.Since(start)
			if last == time.Duration(0) {
				last = t
			}
			fmt.Printf(
				"Depth: %v/%v\t t: %.2fs\tt/m: %.2fs\tr: %.2fx\te: %v\tt/e: %.2fµs\n",
				d.depth,
				d.captureDepth,
				t.Seconds(),
				t.Seconds()/float64(moves),
				float64(t)/float64(last),
				engine.EvalsCount,
				float64(t.Microseconds())/float64(engine.EvalsCount),
			)
			last = t
			engine.PrintBestMoveIndexes()
		}
	}
	fmt.Printf("Total time: %.2fs\n", time.Since(totalStart).Seconds())
}

func (s *TestSuite) TestBestMoveIndexes() {
	engine := New(2, 2, 0)
	g := s.games[2].Copy()

	_, _, err := engine.GetBestMove(g)
	s.Require().NoError(err)
	engine.PrintBestMoveIndexes()
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

		g := s.games[2].Copy()

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
		name      string
		gameIndex int
		move      string
	}

	testCases := []testCase{
		{
			"mate",
			0,
			"g1-a7",
		},
		{
			"free queen",
			1,
			"a7-b6",
		},
	}
	failures := 0
	for _, tc := range testCases {
		g := s.games[tc.gameIndex].Copy()
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

func (s *TestSuite) TestGetBestMove() {
	engine := s.engine

	startTime := time.Now()

	g := s.games[2].Copy()

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
