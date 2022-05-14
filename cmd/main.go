package main

import (
	"fmt"
	"math/rand"

	"github.com/vpoliakov01/2v2ChessAI/game"
)

func main() {
	g := game.New()

	for i := 0; i < 20; i++ {
		moves := g.GetMoves()

		move := moves[rand.Intn(len(moves))]
		g.Play(move)
		g.Board.Draw()
		fmt.Println(move)
	}
}
