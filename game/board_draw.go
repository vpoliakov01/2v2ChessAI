package game

import "fmt"

func (b *Board) Draw() {
	for rank := boardSize - 1; rank >= 0; rank-- {
		fmt.Print("   ---------------------------------------------------------\n")

		fmt.Printf("%-3v", rank+1)
		for file := 0; file < boardSize; file++ {
			fmt.Printf("|%v", SquareStringValue(b[rank][file]))
		}

		fmt.Print("|\n")
	}

	fmt.Print("   ---------------------------------------------------------\n")

	fmt.Print("    ")
	for file := 0; file < boardSize; file++ {
		fmt.Printf(" %v  ", string(int('A')+file))
	}

	fmt.Println()
}

func SquareStringValue(v int) string {
	switch v {
	case inactiveSquare:
		return "███"
	case emptySquare:
		return "   "
	default:
		return Piece(v).String()
	}
}
