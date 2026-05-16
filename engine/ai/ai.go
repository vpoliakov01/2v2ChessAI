package ai

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"sort"

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

// AI is the ai engine used for evaluating the position and picking the best move.
type AI struct {
	Depth        int
	CaptureDepth int
	Spread       int
	SpreadDrop   int

	EvalsCount int
	EvalLimit  int

	BestMoveIndexes []AvgAcc
	enableDebug     bool
}

func init() {
	fmt.Printf("Running on %v CPUs\n", cpus)
	runtime.GOMAXPROCS(cpus) // Should be equal to runtime.NumCPU() by default since go 1.5, but set just in case.
}

// New creates a new AI.
func New(depth, captureDepth, spread, spreadDrop, evalLimit int, options ...func(*AI)) *AI {
	if evalLimit == 0 {
		evalLimit = MaxEvalLimit
	}

	ai := &AI{
		Depth:        depth,
		CaptureDepth: captureDepth,
		Spread:       spread,
		SpreadDrop:   spreadDrop,
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

// recordBestMoveIndex updates per-depth move-ordering analytics. Safe to call
// from multiple goroutines; only runs when debug analytics are enabled.
func (ai *AI) recordBestMoveIndex(data BestMoveIndexes) {
	if !ai.enableDebug {
		return
	}

	acc := &ai.BestMoveIndexes[data.Depth]
	acc.Count++
	acc.IndexSum += data.MoveIndex
	acc.MaxIndex = max(acc.MaxIndex, data.MoveIndex)
	acc.TotalMoves += data.TotalMoves
}

// Stop stops evaluation of GetBestMove.
func (ai *AI) Stop() {
	ai.EvalsCount = ai.EvalLimit
}

// GetBestMove returns the best move for the active player to play along with
// the predicted continuation (the principal variation) up to the search depth.
// The first element of the continuation is the best move itself.
func (ai *AI) GetBestMove(g *game.Game) (continuation []game.Move, score float64, err error) {
	ai.EvalsCount = 0

	if g.HasEnded() {
		return nil, float64(g.Winner), ErrGameEnded
	}

	// Run Negamax at depth 5 and use the score to set alpha-beta pruning bounds for the main Negamax search at depth ai.Depth.

	forcedMateScore := 1002 - float64(ai.Depth) // No point on trying to improve the score if you are forcing mate.
	continuation, score = ai.Negamax(g, 1, ai.EvaluateCurrent(g), -forcedMateScore, forcedMateScore)
	if len(continuation) == 0 {
		return nil, 0, ErrNoMoves
	}

	return continuation, score, nil
}

// Negamax (minimax + negation) recursively finds the position to which
// picking the best move by every player leads.
// Alpha and beta params are used for alpha-beta pruning (skipping evalution
// of branches that are guaranteed not to be picked by any of players).
func (ai *AI) Negamax(g *game.Game, depth int, eval, alpha, beta float64) (continuation []game.Move, score float64) {
	// Check base cases.
	if g.HasEnded() {
		return nil, float64(-1001 + depth)
	}
	if depth > ai.CaptureDepth {
		return nil, eval
	}

	// Get moves to search.
	moves := g.GetMoves().Flatten()
	moveEvalEstimates := map[game.Move]moveScore{}

	for i := range moves {
		capturedPiece := g.Play(moves[i])

		moveEvalEstimates[moves[i]] = moveScore{moves[i], -ai.EvaluateCurrent(g)}
		g.UnplayMove(moves[i], capturedPiece)
	}

	// Sort to process "immediately stronger" moves first.
	sort.Slice(moves, func(a, b int) bool {
		return moveEvalEstimates[moves[a]].score > moveEvalEstimates[moves[b]].score
	})

	bestContinuation := []game.Move{}
	moveIndexesToSearch := ai.GetMoveIndexesToSearch(g, moves, depth)
	bestMoveIndex := moveIndexesToSearch[0]
	bestScore := -math.MaxFloat64

	for _, i := range moveIndexesToSearch {
		move := moves[i]
		eval := -moveEvalEstimates[move].score

		capturedPiece := g.Play(move)
		opponentContinuation, opponentScore := ai.Negamax(g, depth+1, eval, -beta, -alpha)
		g.UnplayMove(move, capturedPiece)

		score := -opponentScore

		if score > bestScore {
			bestScore = score
			bestMoveIndex = i
			bestContinuation = opponentContinuation
		}

		if bestScore > alpha {
			alpha = bestScore
		}

		if alpha >= beta || ai.EvalsCount >= ai.EvalLimit {
			break
		}
	}

	ai.recordBestMoveIndex(BestMoveIndexes{
		Depth:      depth,
		MoveIndex:  bestMoveIndex,
		TotalMoves: len(moves),
	})
	return append([]game.Move{moves[bestMoveIndex]}, bestContinuation...), bestScore
}

// GetMoveIndexesToSearch returns the indexes of the moves to search.
// TODO: return mix of captures, development moves, and king safety moves.
func (ai *AI) GetMoveIndexesToSearch(g *game.Game, moves []game.Move, depth int) []int {
	indexes := []int{}
	movesLeft := ai.Spread - depth/4*ai.SpreadDrop
	if movesLeft < 1 {
		movesLeft = 1
	}
	capturesLeft := movesLeft/2 + 1

	for i, move := range moves {
		if movesLeft == 0 {
			return indexes
		}

		if !g.Board.GetPiece(move.To).IsEmpty() { // Capture
			if capturesLeft == 0 {
				continue
			}
			capturesLeft--
		}

		movesLeft--
		indexes = append(indexes, i)
	}

	return indexes
}

// EvaluateCurrent returns the difference between strengths of the team making the move and the other team.
func (ai *AI) EvaluateCurrent(g *game.Game) float64 {
	ai.EvalsCount++
	playerStrengths := [4]float64{}

	if g.HasEnded() {
		return float64(g.ActivePlayer.Team()*g.Winner) * 1000
	}

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
