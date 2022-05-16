package game

import (
	"fmt"

	"github.com/vpoliakov01/2v2ChessAI/color"
)

func (b *Board) Draw() {
	for rank := BoardSize - 1; rank >= 0; rank-- {
		fmt.Println(color.Reset, "  +---+---+---+---+---+---+---+---+---+---+---+---+---+---+")

		fmt.Printf("%-3v", rank+1)
		for file := 0; file < BoardSize; file++ {
			fmt.Printf("|%v", b.Get(NewSquare(rank, file)))
		}

		fmt.Println("|")
	}
	fmt.Println(color.Reset, "  +---+---+---+---+---+---+---+---+---+---+---+---+---+---+")

	fmt.Print("    ")
	for file := 0; file < BoardSize; file++ {
		fmt.Printf(" %v  ", fmt.Sprintf("%c", int('A')+file))
	}

	fmt.Println()
}
