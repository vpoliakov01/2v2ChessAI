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
	Depth int
}

func init() {
	fmt.Printf("Running on %v CPUs\n", cpus)
	runtime.GOMAXPROCS(cpus) // Should be equal to runtime.NumCPU() by default since go 1.5, but set just in case.
}

// New creates a new AI.
func New(depth int) *AI {
	return &AI{
		Depth: depth,
	}
}

// GetBestMove returns the best move for the active player to play.
func (ai *AI) GetBestMove(g *game.Game) (bestMove *game.Move, score float64, err error) {
	if g.HasEnded() {
		return nil, float64(g.Score), ErrGameEnded
	}

	bestMove, score = ai.Negamax(g, 0, ai.EvaluateCurrent(g), math.Inf(-1), math.Inf(1))
	if bestMove == nil {
		return nil, 0, ErrNoMoves
	}

	return bestMove, score, nil
}

// Negamax (minimax + negation) recursively finds the position to which
// picking the best move by every player leads.
// Alpha and beta params are used for alpha-beta pruning (skipping evalution
// of branches that are guaranteed not to be picked by any of players).
func (ai *AI) Negamax(g *game.Game, depth int, eval, alpha, beta float64) (nextMove *game.Move, score float64) {
	if depth == ai.Depth {
		return nil, eval
	}

	// Instead of calculating checks, just evaluate until king capture.
	// In 2v2 chess king capture is actually possible since teammate A can
	// unblock the path between a teammate's B piece and the opponent's king.
	if !g.HasKing(g.ActivePlayer) {
		// Prefer finishing the game for the winner and prolonging it for the loser.
		return nil, float64(math.MinInt32 + depth)
	}

	moveMap := g.GetMoves() // map[Square][]Square
	moves := moveMap.Flatten()
	if len(moves) == 0 {
		return nil, 0
	}

	// Channel for communicating results of position evaluations.
	c := make(chan moveScore, len(moves))
	moveEvalEstimates := map[game.Move]moveScore{}

	for i := range moves {
		go func(move game.Move) {
			gameCopy := g.Copy()
			gameCopy.Play(move)
			c <- moveScore{move, ai.EvaluateCurrent(gameCopy)} // Although the keys are different, concurrent map writes are not allowed.
		}(moves[i])
	}

	for range moves {
		ms := <-c
		moveEvalEstimates[ms.move] = ms
	}

	// Sort to process "immediately stronger" moves first.
	sort.Slice(moves, func(a, b int) bool {
		return moveEvalEstimates[moves[a]].score < moveEvalEstimates[moves[b]].score
	})

	for i := range moves {
		gameCopy := g.Copy()
		gameCopy.Play(moves[i])
		_, opponentScore := ai.Negamax(gameCopy, depth+1, moveEvalEstimates[moves[i]].score, -beta, -alpha)
		score := -opponentScore

		// If the score is already better than what the opponent could get by playing
		// another move, we can assume that the opponent will not play this move,
		// so we can stop evaluating this branch.
		if score >= beta {
			return &moves[i], beta
		}

		if score > alpha {
			alpha = score
			nextMove = &moves[i]
		}
	}

	return nextMove, alpha
}

// EvaluateCurrent returns the difference between strengths of the team making the move and the other team.
func (ai *AI) EvaluateCurrent(g *game.Game) float64 {
	playerStrengths := map[game.Player]float64{}
	piecesLeft := 0

	for player := range g.Board.PieceSquares {
		piecesLeft += len(g.Board.PieceSquares[player])
	}

	numMoves := len(g.GetMoves().Flatten())

	// For each piece, run piece strength evaluation.
	for player := range g.Board.PieceSquares {
		for square := range g.Board.PieceSquares[player] {
			piece := game.Piece(g.Board.GetPiece(square)).GamePiece()
			playerStrengths[player] += piece.GetStrength(g.Board, numMoves, square, piecesLeft)
		}
	}

	// Account for the advantage of having a balanced pieces distribution between teammates.
	for player := range g.Board.PieceSquares {
		if playerStrengths[player] > 0 {
			playerStrengths[player] = math.Pow(playerStrengths[player], 0.8)
		}
	}

	score := float64(g.ActivePlayer.Team()) * (playerStrengths[0] + playerStrengths[2] - playerStrengths[1] - playerStrengths[3])

	return score
}
