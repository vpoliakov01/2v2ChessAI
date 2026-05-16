package ai_test

import (
	"fmt"
	"runtime"
	"time"

	. "github.com/vpoliakov01/2v2ChessAI/engine/ai"
)

func (s *TestSuite) TestMultithreading() {
	cpus := runtime.NumCPU()
	times := []time.Duration{}
	moves := 1
	engine := New(12, 12, DefaultSpread, DefaultSpreadDrop, 0)

	fmt.Printf("Testing with %v CPUs\n", cpus)
	for i := 1; i <= cpus; i *= 2 {
		runtime.GOMAXPROCS(i)

		startTime := time.Now()

		g := s.solvedGames[2].Copy()

		for i := 0; i < moves; i++ {
			continuation, _, err := engine.GetBestMove(g)
			if err != nil {
				fmt.Println(err)
				break
			}
			g.Play(continuation[0])
		}

		t := time.Since(startTime)
		times = append(times, t)
		fmt.Printf("%v CPU: %v\n", i, t)

		if i != cpus && i*2 > cpus && cpus > 8 {
			i = cpus / 2
		}
	}

	s.Require().Less(times[len(times)-1], times[0])
}
