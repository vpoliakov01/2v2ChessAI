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
	Evaluation  bool
	LoadPGN     string
}

var flg flags

func main() {
	flag.IntVar(&flg.Depth, "depth", 5, "depth of the engine")
	flag.IntVar(&flg.Moves, "moves", 0, "the number of moves to play (0 for unlimited)")
	flag.IntVar(&flg.HumanPlayer, "human", 0, "the team controlled by the human players")
	flag.BoolVar(&flg.Evaluation, "eval", true, "print evalution after every move")
	flag.StringVar(&flg.LoadPGN, "load", "", "load pgn notation (no sidelines) to setup the board")
	flag.Parse()

	flag.CommandLine.Usage()

	fmt.Printf("\nDepth: %v\nMoves limit: %v\nHuman team: %v\nEvaluation: %v\n\n", flg.Depth, flg.Moves, flg.HumanPlayer, flg.Evaluation)

	engine := ai.New(flg.Depth)

	startTime := time.Now()

	g := game.New()
	if flg.LoadPGN != "" {
		moves, err := input.LoadPGN(flg.LoadPGN)
		if err != nil {
			panic(err)
		}

		for _, move := range moves {
			g.Play(move)
		}
	}

	g.Board.Draw()

	for i := 0; !g.HasEnded(); i++ {
		if flg.Moves > 0 && i >= flg.Moves {
			break
		}

		fmt.Println()

		var move *game.Move
		var score float64
		var err error
		moveStartTime := time.Now()

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
			move, score, err = engine.GetBestMove(g)
			if err != nil {
				if err == ai.ErrGameEnded {
					fmt.Printf("%v: Team %v won!\n", i, g.Score)
				} else {
					fmt.Println(err)
				}
				break
			}

			if flg.Evaluation {
				fmt.Printf("Evaluation: %.3f\n", score*float64(g.ActivePlayer.Team()))
			}
		}

		fmt.Printf("Time used: %.3fs\n", time.Since(moveStartTime).Seconds())
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

	if g.Score != 0 {
		fmt.Printf("Team %v won!\n", g.Score)
	}
	fmt.Printf("Total time: %v\n", time.Since(startTime))
}
