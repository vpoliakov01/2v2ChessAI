package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vpoliakov01/2v2ChessAI/engine/ai"
	"github.com/vpoliakov01/2v2ChessAI/engine/game"
	"github.com/vpoliakov01/2v2ChessAI/engine/io"
	"github.com/vpoliakov01/2v2ChessAI/engine/ui"
)

type flags struct {
	Depth       int
	Moves       int
	HumanPlayer int // game.Team.
	Evaluation  bool
	Load        string
	ReactUI     bool
}

var flg flags

func main() {
	// Parse command line flags
	flag.IntVar(&flg.Depth, "depth", 5, "depth of the engine")
	flag.IntVar(&flg.Moves, "moves", 0, "the number of moves to play (0 for unlimited)")
	flag.IntVar(&flg.HumanPlayer, "human", 1, "the team controlled by the human players")
	flag.BoolVar(&flg.Evaluation, "eval", true, "print evalution after every move")
	flag.StringVar(&flg.Load, "load", "", "load pgn notation (no sidelines) to setup the board")
	flag.BoolVar(&flg.ReactUI, "ui", false, "start the React UI")
	flag.Parse()

	flag.CommandLine.Usage()

	fmt.Printf("\nDepth: %v\nMoves limit: %v\nHuman team: %v\nEvaluation: %v\nLoad: %v\n\n", flg.Depth, flg.Moves, flg.HumanPlayer, flg.Evaluation, flg.Load)

	// Get the project root directory for React app
	if flg.ReactUI {
		ex, err := os.Executable()
		if err != nil {
			log.Printf("Failed to get executable path: %v", err)
			// Continue anyway as this only affects the UI
		} else {
			projectRoot := filepath.Dir(filepath.Dir(ex))
			// Start the React app in a goroutine
			go func() {
				if err := ui.StartReactApp(projectRoot); err != nil {
					log.Printf("Failed to start React app: %v", err)
				}
			}()
		}
	}

	// Original game engine logic
	engine := ai.New(flg.Depth)
	startTime := time.Now()

	var g *game.Game

	if flg.Load != "" {
		g, err = io.Load(flg.Load)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	} else {
		g = game.New()
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
				in, err := io.ReadInput()
				if err != nil {
					fmt.Println(err)
					continue
				}

				switch {
				case strings.ToLower(in) == "save":
					file, err := io.Save(g)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Printf("Saved to %v\n", file)
					continue
				default:
					move, err = io.ParseMove(in)
				}

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
