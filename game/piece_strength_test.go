package game_test

import (
	"fmt"

	. "github.com/vpoliakov01/2v2ChessAI/game"
)

// TestBonuses prints the values for the position bonus maps.
func (s *TestSuite) TestBonuses() {
	funcs := []func(Square) float64{
		GetCenterBonus,
		GetEdgeBonus,
		GetBalanceBonus,
	}

	for _, f := range funcs {
		for rank := 0; rank < BoardSize; rank++ {
			for file := 0; file < BoardSize; file++ {
				fmt.Printf("%-2.2f ", f(NewSquare(rank, file)))
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
