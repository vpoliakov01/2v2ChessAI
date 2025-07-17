package ai

import (
	"fmt"
	"time"
)

type AvgAcc struct {
	IndexSum   int
	MaxIndex   int
	TotalMoves int
	Count      int
}

type BestMoveIndexes struct {
	Depth      int
	MoveIndex  int
	TotalMoves int
}

func (ai *AI) PrintBestMoveIndexes() {
	time.Sleep(time.Millisecond)
	for i := range ai.BestMoveIndexes {
		if ai.BestMoveIndexes[i].Count == 0 {
			break
		}

		humanIndex := 1 + float64(ai.BestMoveIndexes[i].IndexSum)/float64(ai.BestMoveIndexes[i].Count)
		maxHumanIndex := 1 + ai.BestMoveIndexes[i].MaxIndex
		moves := float64(ai.BestMoveIndexes[i].TotalMoves) / float64(ai.BestMoveIndexes[i].Count)

		if i == ai.Depth {
			fmt.Println("\t----------------------")
		}

		fmt.Printf(
			"\t%v: %5.2f (%v) / %5.2f = %5.2f (total: %v)\n",
			i,
			humanIndex,
			maxHumanIndex,
			moves,
			humanIndex/moves,
			ai.BestMoveIndexes[i].Count,
		)
	}
	fmt.Println()
}
