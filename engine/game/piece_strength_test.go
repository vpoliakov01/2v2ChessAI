package game_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	. "github.com/vpoliakov01/2v2ChessAI/engine/game"
)

type TestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// TestBonuses prints the values for the position bonus maps.
func (s *TestSuite) TestBonuses() {
	bonus := 0.0
	// bonus := 0.5

	testCases := []struct {
		name string
		f    func(Square) float64
	}{
		{
			name: "GetCenterBonus",
			f:    GetCenterBonus,
		},
		{
			name: "GetEdgeBonus",
			f:    GetEdgeBonus,
		},
		{
			name: "GetBalanceBonus",
			f:    GetBalanceBonus,
		},
		{
			name: "GetCenterBonus+GetBalanceBonus",
			f: func(s Square) float64 {
				return (GetCenterBonus(s) + GetBalanceBonus(s)) / 2
			},
		},
		{
			name: "GetEdgeBonus*GetBalanceBonus",
			f: func(s Square) float64 {
				return (GetEdgeBonus(s) + GetBalanceBonus(s)) / 2
			},
		},
		{
			name: "GetDefenseBonus",
			f: func(s Square) float64 {
				return GetDefenseBonus(s, 1)
			},
		},
		{
			name: "GetAttackBonus",
			f: func(s Square) float64 {
				return GetAttackBonus(s, 1)
			},
		},
		{
			name: "GetProximityBonus Bishop",
			f: func(s Square) float64 {
				return GetProximityBonus(s, Square{1, 4})
			},
		},
		{
			name: "GetProximityBonus Rook",
			f: func(s Square) float64 {
				return GetProximityBonus(s, Square{5, 5})
			},
		},
		{
			name: "GetProximityBonus Queen",
			f: func(s Square) float64 {
				return GetProximityBonus(s, Square{6, 6})
			},
		},
	}
	for _, tc := range testCases {
		fmt.Println(tc.name)
		sum := 0.0
		for rank := 0; rank < BoardSize; rank++ {
			for file := 0; file < BoardSize; file++ {
				if isCorner(rank, file) {
					fmt.Printf("     ")
				} else {
					fmt.Printf("%.2f ", tc.f(Square{rank, file})+bonus)
					sum += tc.f(Square{rank, file})
				}
			}
			fmt.Println()
		}
		fmt.Printf("avg: %.2f\n\n", sum/(BoardSize*BoardSize-4*9))
	}
}

func isCorner(rank, file int) bool {
	return (rank < 3 && file < 3) ||
		(rank < 3 && file >= BoardSize-3) ||
		(rank >= BoardSize-3 && file < 3) ||
		(rank >= BoardSize-3 && file >= BoardSize-3)
}
