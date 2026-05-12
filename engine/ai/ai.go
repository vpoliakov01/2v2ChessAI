package ai

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

var (
	MaxEvalLimit = int(1e12)
	ErrGameEnded = errors.New("the game has ended")
	ErrNoMoves   = errors.New("no move can be made in this position")
	cpus         = runtime.NumCPU()
)

type moveScore struct {
	move  game.Move
	score float64
}

// AI is the ai engine used for evaluating the position and picking the best move.
type AI struct {
	Depth           int
	CaptureDepth    int
	EvalLimit       int
	BestMoveIndexes []AvgAcc

	evalsCount atomic.Int64

	enableDebug          bool
	bestMoveIndexesMutex sync.Mutex
}

func init() {
	fmt.Printf("Running on %v CPUs\n", cpus)
	runtime.GOMAXPROCS(cpus) // Should be equal to runtime.NumCPU() by default since go 1.5, but set just in case.
}

// New creates a new AI.
func New(depth, captureDepth, evalLimit int, options ...func(*AI)) *AI {
	if evalLimit == 0 {
		evalLimit = MaxEvalLimit
	}

	ai := &AI{
		Depth:        depth,
		CaptureDepth: captureDepth,
		EvalLimit:    evalLimit,
	}
	for _, option := range options {
		option(ai)
	}

	if ai.enableDebug {
		ai.BestMoveIndexes = make([]AvgAcc, 100)
	}

	return ai
}

// WithEnableDebug enables debug analytics.
func WithEnableDebug(enableDebug bool) func(*AI) {
	return func(ai *AI) {
		ai.enableDebug = enableDebug
	}
}

// EvalsCount returns the number of positions evaluated during the current
// or most recent search. Safe to call concurrently.
func (ai *AI) EvalsCount() int {
	return int(ai.evalsCount.Load())
}

// Stop stops GetBestMove by pushing the evaluation counter to the limit so
// the recursion gate in Negamax short-circuits.
func (ai *AI) Stop() {
	ai.evalsCount.Store(int64(ai.EvalLimit))
}

// recordBestMoveIndex updates per-depth move-ordering analytics. Safe to call
// from multiple goroutines; only runs when debug analytics are enabled.
func (ai *AI) recordBestMoveIndex(data BestMoveIndexes) {
	if !ai.enableDebug {
		return
	}

	ai.bestMoveIndexesMutex.Lock()
	defer ai.bestMoveIndexesMutex.Unlock()

	acc := &ai.BestMoveIndexes[data.Depth]
	acc.Count++
	acc.IndexSum += data.MoveIndex
	acc.MaxIndex = max(acc.MaxIndex, data.MoveIndex)
	acc.TotalMoves += data.TotalMoves
}

// GetBestMove returns the best move for the active player to play along with
// the predicted continuation (the principal variation) up to the search depth.
// The first element of the continuation is the best move itself.
func (ai *AI) GetBestMove(g *game.Game) (continuation []game.Move, score float64, err error) {
	ai.evalsCount.Store(0)

	if g.HasEnded() {
		return nil, float64(g.Winner), ErrGameEnded
	}

	forcedMateScore := 1002 - float64(ai.Depth) // No point on trying to improve the score if you are forcing mate.
	continuation, score = ai.Negamax(g, 0, 0, ai.EvaluateCurrent(g), -forcedMateScore, forcedMateScore)
	if len(continuation) == 0 {
		return nil, 0, ErrNoMoves
	}

	return continuation, score, nil
}

// Negamax (minimax + negation) recursively finds the position to which
// picking the best move by every player leads.
// Alpha and beta params are used for alpha-beta pruning (skipping evalution
// of branches that are guaranteed not to be picked by any of players).
func (ai *AI) Negamax(g *game.Game, depth, reduction int, eval, alpha, beta float64) (continuation []game.Move, score float64) {
	// Instead of calculating checks, just evaluate until king capture.
	// In 2v2 chess king capture is actually possible since teammate A can
	// unblock the path between a teammate's B piece and the opponent's king.
	if g.HasEnded() {
		ai.recordBestMoveIndex(BestMoveIndexes{
			Depth:      depth,
			MoveIndex:  0,
			TotalMoves: 1,
		})
		return nil, float64(-1001 + depth)
	}

	if depth >= ai.CaptureDepth-reduction {
		return nil, eval
	}

	moves := g.GetMoves().Flatten()

	if depth > ai.Depth-reduction {
		captureMoves := []game.Move{}
		for _, move := range moves {
			if !g.Board.GetPiece(move.To).IsEmpty() {
				captureMoves = append(captureMoves, move)
			}
		}
		moves = captureMoves
	}
	if len(moves) == 0 {
		return nil, eval
	}

	var bestContinuation []game.Move

	// Channel for communicating results of position evaluations.
	c := make(chan moveScore, len(moves))
	moveEvalEstimates := map[game.Move]moveScore{}

	for i := range moves {
		go func(move game.Move) {
			gameCopy := g.Copy()
			gameCopy.Play(move)
			c <- moveScore{move, -ai.EvaluateCurrent(gameCopy)} // Negate the opponent's position evaluation.
		}(moves[i])
	}
	ai.evalsCount.Add(int64(len(moves)))

	for range moves {
		moveScore := <-c
		moveEvalEstimates[moveScore.move] = moveScore
	}

	// Sort to process "immediately stronger" moves first.
	// Process captures first.
	// Strongest moves are the weakest for the opponent (lowest score).
	sort.Slice(moves, func(a, b int) bool {
		if g.Board.GetPiece(moves[a].To).IsEmpty() == g.Board.GetPiece(moves[b].To).IsEmpty() {
			return moveEvalEstimates[moves[a]].score > moveEvalEstimates[moves[b]].score
		}
		if !g.Board.GetPiece(moves[a].To).IsEmpty() {
			return true
		}
		return false
	})

	nextMoveIndex := 0

	depthReductionIndex := ai.getDepthReductionThreshold(len(moves), depth)

	for i := range moves {
		if ai.evalsCount.Load() >= int64(ai.EvalLimit) {
			break
		}

		gameCopy := g.Copy()
		gameCopy.Play(moves[i])

		eval := -moveEvalEstimates[moves[i]].score
		depthReduction := 0
		if i >= depthReductionIndex {
			depthReduction = 2
		}

		opponentContinuation, opponentScore := ai.Negamax(gameCopy, depth+1, depthReduction, eval, -beta, -alpha)
		score := -opponentScore

		if depthReduction > 0 && score > alpha {
			opponentContinuation, opponentScore = ai.Negamax(gameCopy, depth+1, 0, eval, -beta, -alpha)
			score = -opponentScore
		}

		// If the score is already better than what the opponent could get by playing
		// any other move, we can assume that the opponent will not play this move,
		// so we can stop evaluating this branch.
		if score >= beta {
			ai.recordBestMoveIndex(BestMoveIndexes{
				Depth:      depth,
				MoveIndex:  i,
				TotalMoves: len(moves),
			})
			return append([]game.Move{moves[i]}, opponentContinuation...), beta
		}

		if score > alpha {
			alpha = score
			nextMoveIndex = i
			bestContinuation = opponentContinuation
		}
	}

	ai.recordBestMoveIndex(BestMoveIndexes{
		Depth:      depth,
		MoveIndex:  nextMoveIndex,
		TotalMoves: len(moves),
	})
	return append([]game.Move{moves[nextMoveIndex]}, bestContinuation...), alpha
}

func (ai *AI) getDepthReductionThreshold(numMoves, depth int) int {
	if depth > ai.Depth { // If in the capture depth territory, check all captures.
		return numMoves
	}

	if depth > 4 {
		return numMoves / 2
	}
	if depth > 2 {
		return numMoves
	}

	return numMoves
}

// EvaluateCurrent returns the difference between strengths of the team making the move and the other team.
func (ai *AI) EvaluateCurrent(g *game.Game) float64 {
	playerStrengths := [4]float64{}

	// For each piece, run piece strength evaluation.
	for player := range g.Board.PieceSquares {
		for square := range g.Board.PieceSquares[player] {
			piece := game.Piece(g.Board.GetPiece(square)).PieceType()
			playerStrengths[player] += piece.GetStrength(g.Board, square, player)
		}
	}

	redYellowStrength := playerStrengths[0] + playerStrengths[2] - math.Abs(playerStrengths[0]-playerStrengths[2])/3
	blueGreenStrength := playerStrengths[1] + playerStrengths[3] - math.Abs(playerStrengths[1]-playerStrengths[3])/3

	return float64(g.ActivePlayer.Team()) * (redYellowStrength - blueGreenStrength)
}
