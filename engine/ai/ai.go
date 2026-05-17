package ai

import (
	"math"
	"sort"

	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

// AI is the ai engine used for evaluating the position and picking the best move.
type AI struct {
	Depth        int
	CaptureDepth int
	Spread       int
	SpreadDrop   int

	EvalsCount int
	EvalLimit  int

	BestMoves   []AvgAcc
	enableDebug bool

	buffers []buffer // One buffer per CPU.
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
		ai.BestMoves = make([]AvgAcc, 100)
	}

	return ai
}

// GetBestMove returns the predicted continuation up to the search depth.
// The first element of the continuation is the best move itself.
func (ai *AI) GetBestMove(g *game.Game) (continuation []game.Move, score float64, err error) {
	ai.EvalsCount = 0
	if g.HasEnded() {
		return nil, float64(g.Winner), ErrGameEnded
	}

	ai.initBuffers()
	buffer := &ai.buffers[0]

	forcedMateScore := 1002 - float64(ai.Depth)
	score = ai.Negamax(g, buffer, 1, ai.EvaluateCurrent(g), -forcedMateScore, forcedMateScore)

	continuation = make([]game.Move, len(buffer.continuation[1]))
	copy(continuation, buffer.continuation[1])

	if len(continuation) == 0 {
		return nil, 0, ErrNoMoves
	}
	return continuation, score, nil
}

// Negamax (minimax + negation) recursively finds the position
// reached by each side picking their best move.
// Alpha and beta params are used for alpha-beta pruning (skipping evalution
// of branches that are guaranteed not to be picked by any of players).
func (ai *AI) Negamax(g *game.Game, buffer *buffer, depth int, eval, alpha, beta float64) (score float64) {
	buffer.continuation[depth] = buffer.continuation[depth][:0] // Reset the buffer.

	// Check base cases.
	if g.HasEnded() {
		return float64(-1001 + depth)
	}
	if depth > ai.CaptureDepth {
		return eval
	}

	moves := g.GetMoves(buffer.moves[depth][:0])
	buffer.moves[depth] = moves // In case moves gets allocated a larger slice with append.

	moveEvals := buffer.moveEvals[depth][:len(moves)]

	// Score each move with a static evaluation of the position after the move.
	for i := range moves {
		capturedPiece := g.Play(moves[i])
		moveEvals[i] = moveScore{moves[i], -ai.EvaluateCurrent(g)}
		g.UnplayMove(moves[i], capturedPiece)
	}
	buffer.moveEvals[depth] = moveEvals

	// Sort to process "immediately stronger" moves first.
	sort.Slice(moveEvals, func(a, b int) bool {
		return moveEvals[a].score > moveEvals[b].score
	})

	// Filter promising moves to search.
	moveIndexesToSearch := ai.getMoveIndexesToSearch(g, moveEvals, depth, buffer.moveIndexesToSearch[depth][:0])
	buffer.moveIndexesToSearch[depth] = moveIndexesToSearch

	bestMoveIndex := moveIndexesToSearch[0]
	bestScore := -math.MaxFloat64

	for _, i := range moveIndexesToSearch {
		move := moveEvals[i].move
		childEval := -moveEvals[i].score

		capturedPiece := g.Play(move)
		opponentScore := ai.Negamax(g, buffer, depth+1, childEval, -beta, -alpha)
		g.UnplayMove(move, capturedPiece)

		score := -opponentScore
		moveEvals[i].score = score

		if score > bestScore {
			bestScore = score
			bestMoveIndex = i

			// Update the continuation at this depth.
			continuation := buffer.continuation[depth][:0]
			continuation = append(continuation, move)
			buffer.continuation[depth] = append(continuation, buffer.continuation[depth+1]...)
		}

		if bestScore > alpha {
			alpha = bestScore
		}

		if alpha >= beta || ai.EvalsCount >= ai.EvalLimit {
			break
		}
	}

	if ai.enableDebug {
		ai.recordBestMove(BestMoveData{
			Depth:      depth,
			MoveIndex:  bestMoveIndex,
			TotalMoves: len(moves),
			ScoreDelta: moveEvals[bestMoveIndex].score - moveEvals[0].score,
		})
	}

	return bestScore
}

// EvaluateCurrent returns the difference between strengths of the team making the move and the opponent team.
func (ai *AI) EvaluateCurrent(g *game.Game) float64 {
	ai.EvalsCount++
	playerStrengths := [4]float64{}

	if g.HasEnded() {
		return float64(g.ActivePlayer.Team()*g.Winner) * 1000
	}

	// For each piece, run piece strength evaluation.
	for player := range g.Board.PieceSquares {
		for square := range g.Board.PieceSquares[player] {
			piece := g.Board.GetPiece(square)
			playerStrengths[player] += piece.GetStrength(g.Board, square, player)
		}
	}

	redYellowStrength := playerStrengths[0] + playerStrengths[2] - math.Abs(playerStrengths[0]-playerStrengths[2])/3
	blueGreenStrength := playerStrengths[1] + playerStrengths[3] - math.Abs(playerStrengths[1]-playerStrengths[3])/3

	return float64(g.ActivePlayer.Team()) * (redYellowStrength - blueGreenStrength)
}

// getMoveIndexesToSearch appends the indexes of moves worth searching to dst and returns the extended slice.
func (ai *AI) getMoveIndexesToSearch(g *game.Game, moveEvals []moveScore, depth int, dst []int) []int {
	movesLeft := max(ai.Spread-depth/4*ai.SpreadDrop, 1)
	capturesLeft := movesLeft/2 + 1

	for i, moveEval := range moveEvals {
		if movesLeft == 0 {
			return dst
		}

		if !g.Board.GetPiece(moveEval.move.To).IsEmpty() { // Capture
			if capturesLeft == 0 {
				continue
			}
			capturesLeft--
		}

		movesLeft--
		dst = append(dst, i)
	}

	return dst
}
