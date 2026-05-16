package ai

import (
	"fmt"
)

type AvgAcc struct {
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

// recordBestMove updates per-depth move-ordering analytics. Safe to call
// from multiple goroutines; only runs when debug analytics are enabled.
func (ai *AI) recordBestMove(data BestMoveData) {
	acc := &ai.BestMoves[data.Depth]
	acc.Count++
	acc.IndexSum += data.MoveIndex
	acc.MaxIndex = max(acc.MaxIndex, data.MoveIndex)
	acc.TotalMoves += data.TotalMoves
	acc.ScoreDelta += data.ScoreDelta
}

func (ai *AI) PrintBestMoveIndexes() {
	fmt.Println("        dep  best max  moves    ratio    Δscore    total")
	for i := range ai.BestMoves {
		if i == 0 {
			continue
		}
		if ai.BestMoves[i].Count == 0 {
			break
		}

		avgIndex := 1 + float64(ai.BestMoves[i].IndexSum)/float64(ai.BestMoves[i].Count)
		maxIndex := 1 + ai.BestMoves[i].MaxIndex
		moves := float64(ai.BestMoves[i].TotalMoves) / float64(ai.BestMoves[i].Count)
		scoreDelta := ai.BestMoves[i].ScoreDelta / float64(ai.BestMoves[i].Count)

		fmt.Printf(
			"\t %2v:%5.2f (%2d) /%4.0f  = %4.1f%%   %7.2f  %7v\n",
			i,
			avgIndex,
			maxIndex,
			moves,
			avgIndex/moves*100,
			scoreDelta,
			ai.BestMoves[i].Count,
		)

		if i == ai.Depth {
			fmt.Println("\t -----------------------------------------------")
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
