package ai_test

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	. "github.com/vpoliakov01/2v2ChessAI/engine/ai"
)

func (s *TestSuite) TestMultithreading() {
	maxCPUs := runtime.NumCPU()
	moves := 1
	engine := New(12, 12, DefaultSpread, DefaultSpreadDrop, 0)

	fmt.Printf("Testing with %v CPUs\n", maxCPUs)
	for i := 0; true; i++ {
		cpus := int(math.Min(math.Pow(2, float64(i)), float64(maxCPUs)))
		runtime.GOMAXPROCS(cpus)
		startTime := time.Now()

		g := s.GetGame("Free queen (a7-b6)")

		for i := 0; i < moves; i++ {
			continuation, _, err := engine.GetBestMove(g.Game)
			if err != nil {
				fmt.Println(err)
				break
			}
			g.Play(continuation[0])
		}

		fmt.Printf("%2d CPU: %v\n", cpus, time.Since(startTime))

		if cpus >= maxCPUs {
			break
		}
	}
}

func (s *TestSuite) TestMultithreadingGetBestMove() {
	maxCPUs := runtime.GOMAXPROCS(0)

	for i := 0; true; i++ {
		cpus := int(math.Min(math.Pow(2, float64(i)), float64(maxCPUs)))

		// Run the same GetBestMove in multiple goroutines
		var wg sync.WaitGroup
		startTime := time.Now()
		for i := 0; i < cpus; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				engine := New(10, 10, DefaultSpread, DefaultSpreadDrop, 0)
				g := s.GetGame("Free queen (a7-b6)")
				_, _, err := engine.GetBestMove(g.Game)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}
		wg.Wait()
		fmt.Printf("%2d CPUs: %v\n", cpus, time.Since(startTime))

		if cpus >= maxCPUs {
			break
		}
	}
}
