package ai

import (
	"errors"
	"math"
	"sort"

	"github.com/vpoliakov01/2v2ChessAI/game" // TODO: remove dot importing.
)

var (
	ErrGameEnded = errors.New("the game has ended")
	ErrNoMoves   = errors.New("no move can be made in this position")
)

type evaluation struct {
	game  *game.Game
	score float64
}

type AI struct {
	Depth int
	Cache map[uint64]float64 // Stores scores of evalutated positions.
}

func New(depth int) *AI {
	return &AI{
		Depth: depth,
	}
}

func (ai *AI) GetBestMove(g *game.Game) (*game.Move, error) {
	if g.HasEnded() {
		return nil, ErrGameEnded
	}

	bestMove, _ := ai.Negamax(g, 0, math.Inf(-1), math.Inf(1))
	if bestMove == nil {
		return nil, ErrGameEnded
	}

	return bestMove, nil
}

func (ai *AI) Negamax(g *game.Game, depth int, alpha, beta float64) (nextMove *game.Move, score float64) {
	if depth == ai.Depth {
		return nil, ai.EvaluateCurrent(g)
	}

	if !g.HasKing(g.ActivePlayer) {
		// Prefer finishing the game for the winner and prolonging it for the loser.
		return nil, float64(math.MinInt32 + depth)
	}

	moves := g.GetMoves()
	if len(moves) == 0 {
		return nil, 0
	}

	moveEvalEstimates := map[game.Move]evaluation{}

	for i := range moves {
		gameCopy := g.Copy()
		gameCopy.Play(moves[i])
		moveEvalEstimates[moves[i]] = evaluation{gameCopy, ai.EvaluateCurrent(gameCopy)}
	}

	// Sort to process "immediately stronger" moves first.
	sort.Slice(moves, func(a, b int) bool {
		return moveEvalEstimates[moves[a]].score < moveEvalEstimates[moves[b]].score
	})

	for i := range moves {
		_, opponentScore := ai.Negamax(moveEvalEstimates[moves[i]].game, depth+1, -beta, -alpha)
		score := -opponentScore

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

// EvaluateCurrent returns the difference between the team making the move and the other team.
func (ai *AI) EvaluateCurrent(g *game.Game) float64 {
	playerStrength := map[game.Player]float64{}
	piecesLeft := 0

	for player := range g.Board.PieceSquares {
		piecesLeft += g.Board.PieceSquares[player].Size()
	}

	for player := range g.Board.PieceSquares {
		for _, sq := range g.Board.PieceSquares[player].Elements() {
			square := sq.(game.Square)
			piece := game.Piece(g.Board.Get(square)).GetGamePiece()
			playerStrength[player] += piece.GetStrength(g.Board, square, piecesLeft)
		}

		// Account for the advantage of having a balanced pieces distribution between teammates.
		if playerStrength[player] > 0 {
			playerStrength[player] = math.Pow(playerStrength[player], 0.8)
		}
	}

	side := g.ActivePlayer.GetTeam().Side()
	score := float64(side) * (playerStrength[0] + playerStrength[2] - playerStrength[1] - playerStrength[3])

	return score
}
