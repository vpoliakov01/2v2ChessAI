package ai

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

var (
	MaxEvalLimit = int(1e12)
	ErrGameEnded = errors.New("the game has ended")
	ErrNoMoves   = errors.New("no move can be made in this position")
	cpus         = runtime.NumCPU()

	DefaultDepth        = 12
	DefaultCaptureDepth = 12
	DefaultSpread       = 8
	DefaultSpreadDrop   = 2
	DefaultEvalLimit    = MaxEvalLimit
)

type moveScore struct {
	move  game.Move
	score float64
}

func init() {
	fmt.Printf("Running on %v CPUs (GOMAXPROCS=%v)\n", cpus, runtime.GOMAXPROCS(0))
}

// WithEnableDebug enables debug analytics.
func WithEnableDebug(enableDebug bool) func(*AI) {
	return func(ai *AI) {
		ai.enableDebug = enableDebug
	}
}

// Stop stops evaluation of GetBestMove.
func (ai *AI) Stop() {
	ai.EvalsCount = ai.EvalLimit
}
