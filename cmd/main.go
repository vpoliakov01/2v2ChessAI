package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/vpoliakov01/2v2ChessAI/ai"
	"github.com/vpoliakov01/2v2ChessAI/game"
	"github.com/vpoliakov01/2v2ChessAI/input"
)

type flags struct {
	Depth       int
	Moves       int
	HumanPlayer int // game.Team.
}

var flg flags

func main() {
	flag.IntVar(&flg.Depth, "depth", 4, "depth of the engine")
	flag.IntVar(&flg.Moves, "moves", 0, "the number of moves to play (0 for unlimited)")
	flag.IntVar(&flg.HumanPlayer, "human", 1, "the team controlled by the human players")
	flag.Parse()

	fmt.Println("Flags:")
	fmt.Println("	depth - engine depth")
	fmt.Println("	moves - limit on number of moves to be played (useful for engine vs engine game)")
	fmt.Println("	human - team to be controlled by user input (1: Red/Yellow, -1: Blue/Green, 0: engine vs engine")

	fmt.Printf("Depth: %v\nMoves limit: %v\nHuman team: %v\n", flg.Depth, flg.Moves, flg.HumanPlayer)

	engine := ai.New(flg.Depth)

	startTime := time.Now()

	g := game.New()
	g.Board.Draw()

	for i := 0; !g.HasEnded(); i++ {
		if flg.Moves > 0 && i >= flg.Moves {
			break
		}

		var move *game.Move
		var err error

		if g.ActivePlayer.Team() == game.Team(flg.HumanPlayer) {
			for {
				move, err = input.ReadMove()
				if err != nil {
					fmt.Println(err)
					continue
				}

				if err := g.ValidateMove(move); err != nil {
					fmt.Println(err)
					continue
				}
				break
			}
		} else {
			move, err = engine.GetBestMove(g)
			if err != nil {
				if err == ai.ErrGameEnded {
					fmt.Printf("%v: Team %v won!\n", i, g.Score)
				} else {
					fmt.Println(err)
				}
				break
			}
		}
		fmt.Println(move)

		piece := game.Piece(g.Board.GetPiece(move.From))
		if !g.Board.IsEmpty(move.To) {
			capturedPiece := game.Piece(g.Board.GetPiece(move.To))
			fmt.Printf("%v: %v takes %v after %v\n", i, piece, capturedPiece, move)
		} else {
			fmt.Printf("%v: %v moves %v\n", i, piece, move)
		}

		g.Play(*move)
		g.Board.Draw()
	}

	g.Board.Draw()
	fmt.Println("Evaluation: ", engine.EvaluateCurrent(g))
	fmt.Println(time.Since(startTime))
}
