package ai

import "github.com/vpoliakov01/2v2ChessAI/engine/game"

const (
	MovesUpperBound = 256
)

// buffer holds per-depth reusable storage for one search worker to avoid repeated allocations.
type buffer struct {
	moves               [][]game.Move // candidate moves at each depth
	moveEvals           [][]moveScore // evaluation of each candidate move
	moveIndexesToSearch [][]int       // indexes (into moves/moveEvals) selected for deeper search
	continuation        [][]game.Move // principal variation (predicted line of best play) from each depth down
}

// init populates buffers for searches up to maxDepth. Safe to call repeatedly;
// it's a no-op once the buffers are already at least maxDepth.
func (buff *buffer) init(maxDepth int) {
	if len(buff.moves) >= maxDepth {
		return
	}

	buff.moves = make([][]game.Move, maxDepth)
	buff.moveEvals = make([][]moveScore, maxDepth)
	buff.moveIndexesToSearch = make([][]int, maxDepth)
	buff.continuation = make([][]game.Move, maxDepth)

	for d := range buff.continuation {
		buff.continuation[d] = make([]game.Move, 0, maxDepth)
	}

	for i := range buff.moves {
		buff.moves[i] = make([]game.Move, 0, MovesUpperBound)
		buff.moveEvals[i] = make([]moveScore, 0, MovesUpperBound)
		buff.moveIndexesToSearch[i] = make([]int, 0, MovesUpperBound)
	}
}

// initBuffers lazily allocates one buffer per CPU, sized for the current configuration.
func (ai *AI) initBuffers() {
	if len(ai.buffers) < cpus {
		ai.buffers = make([]buffer, cpus)
	}

	for i := range ai.buffers { // Per each cpu.
		ai.buffers[i].init(ai.CaptureDepth + 2)
	}
}
