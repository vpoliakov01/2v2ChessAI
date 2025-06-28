package game

import (
	"fmt"

	"github.com/vpoliakov01/2v2ChessAI/engine/color"
)

// Draw draws the board. Clunky but does the job.
func (b *Board) Draw() {
	fmt.Print("    ")
	for file := 0; file < BoardSize; file++ {
		fmt.Printf(" %v  ", fmt.Sprintf("%c", int('A')+file))
	}
	fmt.Println()

	for rank := BoardSize - 1; rank >= 0; rank-- {
		fmt.Println(color.Reset, "  +---+---+---+---+---+---+---+---+---+---+---+---+---+---+")

		fmt.Printf("%2v ", rank+1)

		for file := 0; file < BoardSize; file++ {
			fmt.Printf("|%v", b.GetPiece(Square{rank, file}))
		}

		fmt.Printf("| %-2v\n", rank+1)
	}
	fmt.Println(color.Reset, "  +---+---+---+---+---+---+---+---+---+---+---+---+---+---+")

	fmt.Print("    ")
	for file := 0; file < BoardSize; file++ {
		fmt.Printf(" %v  ", fmt.Sprintf("%c", int('A')+file))
	}

	fmt.Println()
}
