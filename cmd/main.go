package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/vpoliakov01/2v2ChessAI/ai"
	"github.com/vpoliakov01/2v2ChessAI/game"
)

type flags struct {
	Depth int
	Moves int
}

var flg flags

func main() {
	flag.IntVar(&flg.Depth, "depth", 4, "depth of the engine")
	flag.IntVar(&flg.Moves, "moves", 100, "the number of moves to play")

	engine := ai.New(flg.Depth)

	startTime := time.Now()

	g := game.New()

	for i := 0; i < flg.Moves; i++ {
		move, err := engine.GetBestMove(g)
		if err != nil {
			if err == ai.ErrGameEnded {
				fmt.Printf("%v: Team %v won!\n", i, g.Score)
			} else {
				fmt.Println(err)
			}
			break
		}

		if !g.Board.IsEmpty(move.To) {
			capturedPiece := game.Piece(g.Board.Get(move.To))
			opponent := capturedPiece.GetPlayer()
			piece := game.Piece(g.Board.Get(move.From))
			player := piece.GetPlayer()
			fmt.Printf("%v: P%v's %v takes P%v's %v after %v\n", i, player, piece, opponent, capturedPiece, move)
			g.Board.Draw()
		}

		g.Play(*move)
		g.Board.Draw()
	}

	g.Board.Draw()
	fmt.Println(engine.EvaluateCurrent(g))
	fmt.Println(time.Since(startTime))
}
