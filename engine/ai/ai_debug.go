package ai

import (
	"fmt"
)

// BestmoveDataAvgAcc accumulates analytics on the indexes of best moves at each depth, for debugging and tuning move ordering heuristics.
type BestmoveDataAvgAcc struct {
	IndexSum   int
	MaxIndex   int
	TotalMoves int
	ScoreDelta float64
	Count      int
}

type BestMoveData struct {
	Depth      int
	MoveIndex  int
	TotalMoves int
	ScoreDelta float64
}

func (ai *AI) InitDebug() {
	ai.BestMoves = make([][]BestmoveDataAvgAcc, cpus)
	for i := range ai.BestMoves {
		ai.BestMoves[i] = make([]BestmoveDataAvgAcc, ai.CaptureDepth+1)
	}
}

// recordBestMove updates per-depth move-ordering analytics.
func (ai *AI) recordBestMove(data BestMoveData, cpu int) {
	acc := &ai.BestMoves[cpu][data.Depth]
	acc.Count++
	acc.IndexSum += data.MoveIndex
	acc.MaxIndex = max(acc.MaxIndex, data.MoveIndex)
	acc.TotalMoves += data.TotalMoves
	acc.ScoreDelta += data.ScoreDelta
}

func (ai *AI) PrintBestMoveIndexes(printIndividualCPUs bool, printAllCPUs bool) {
	fmt.Println("        dep  best max  moves    ratio    Δscore    total")

	cpuAcc := make([][]BestmoveDataAvgAcc, cpus)
	for i := range cpuAcc {
		cpuAcc[i] = make([]BestmoveDataAvgAcc, ai.CaptureDepth+1)
	}

	for cpu := range ai.BestMoves {
		hasData := false
		for depth := range ai.BestMoves[cpu] {
			if ai.BestMoves[cpu][depth].Count > 0 {
				hasData = true
				break
			}
		}
		if !hasData {
			continue
		}

		if printIndividualCPUs {
			fmt.Printf("CPU %v:\n", cpu)
		}

		for depth := range ai.BestMoves[cpu] {
			if depth == 0 {
				continue
			}

			acc := ai.BestMoves[cpu][depth]
			if acc.Count == 0 {
				continue
			}

			avgIndex := float64(acc.IndexSum) / float64(acc.Count)
			maxIndex := acc.MaxIndex
			moves := float64(acc.TotalMoves) / float64(acc.Count)
			scoreDelta := acc.ScoreDelta / float64(acc.Count)

			sharedAcc := &cpuAcc[cpu][depth]
			sharedAcc.Count += acc.Count
			sharedAcc.IndexSum += acc.IndexSum
			sharedAcc.MaxIndex = max(sharedAcc.MaxIndex, acc.MaxIndex)
			sharedAcc.TotalMoves += acc.TotalMoves
			sharedAcc.ScoreDelta += acc.ScoreDelta

			if printIndividualCPUs {
				fmt.Printf(
					"\t %2v:%5.2f (%2d) /%4.0f  = %4.1f%%   %7.2f  %7v\n",
					depth,
					avgIndex+1, // Human index
					maxIndex+1,
					moves,
					avgIndex/moves*100,
					scoreDelta,
					acc.Count,
				)

				if depth == ai.Depth {
					fmt.Println("\t -----------------------------------------------")
				}
			}
		}
	}

	if printAllCPUs {
		fmt.Println("All CPUs")
		for depth := 1; depth <= ai.CaptureDepth; depth++ {
			acc := &cpuAcc[0][depth]

			avgIndex := float64(acc.IndexSum) / float64(acc.Count)
			maxIndex := acc.MaxIndex
			moves := float64(acc.TotalMoves) / float64(acc.Count)
			scoreDelta := acc.ScoreDelta / float64(acc.Count)

			fmt.Printf(
				"\t %2v:%5.2f (%2d) /%4.0f  = %4.1f%%   %7.2f  %7v\n",
				depth,
				avgIndex+1, // Human index
				maxIndex+1,
				moves,
				avgIndex/moves*100,
				scoreDelta,
				acc.Count,
			)
			if depth == ai.Depth {
				fmt.Println("\t -----------------------------------------------")
			}
		}
	}

	fmt.Println()
}

func (ai *AI) TotalPossibleEvals() int {
	total := 1
	for depth := 1; depth <= ai.Depth; depth++ {
		total *= ai.Spread - depth/4*ai.SpreadDrop
	}
	return total
}
