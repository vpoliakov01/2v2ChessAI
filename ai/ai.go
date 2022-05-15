package ai

import (
	"errors"
	"math"
	"sort"

	. "github.com/vpoliakov01/2v2ChessAI/game" // TODO: remove dot importing.
)

var (
	ErrGameEnded = errors.New("the game has ended")
	ErrNoMoves   = errors.New("no move can be made in this position")

	depth = 4
)

func GetBestMove(g *Game) (*Move, error) {
	if g.HasEnded() {
		return nil, ErrGameEnded
	}

	moves := g.GetMoves()
	if (len(moves)) == 0 {
		return nil, ErrNoMoves
	}

	bestScore := math.Inf(-1)
	var bestMove *Move

	for i := range moves {
		gameCopy := g.Copy()
		gameCopy.Play(moves[i])
		score := -Negamax(gameCopy, depth, math.Inf(1), math.Inf(-1))

		if score > bestScore {
			bestScore = score
			bestMove = &moves[i]
		}
	}

	return bestMove, nil
}

func Negamax(g *Game, depth int, alpha, beta float64) float64 {
	if depth == 0 {
		return EvaluateCurrent(g)
	}

	moves := g.GetMoves()
	// If the active player cannot make a move, treat it as a loss.
	if len(moves) == 0 {
		if g.HasKing(g.ActivePlayer) {
			return 0
		}
		return -KingStrength
	}

	type evaluation struct {
		game  *Game
		score float64
	}

	moveEvalEstimates := map[Move]evaluation{}

	for i := range moves {
		gameCopy := g.Copy()
		gameCopy.Play(moves[i])
		moveEvalEstimates[moves[i]] = evaluation{gameCopy, EvaluateCurrent(gameCopy)}
	}

	sort.Slice(moves, func(a, b int) bool {
		return moveEvalEstimates[moves[b]].score-moveEvalEstimates[moves[a]].score < 0
	})

	for _, move := range moves {
		score := -Negamax(moveEvalEstimates[move].game, depth-1, -beta, -alpha)
		if score >= beta {
			return beta
		}
		alpha = math.Max(score, alpha)
	}

	return alpha
}

func EvaluateCurrent(g *Game) float64 {
	playerStrength := map[Player]float64{}
	piecesLeft := 0

	for player := range g.Board.PieceSquares {
		piecesLeft += g.Board.PieceSquares[player].Size()
	}

	for player := range g.Board.PieceSquares {
		for _, sq := range g.Board.PieceSquares[player].Elements() {
			square := sq.(Square)
			piece := Piece(g.Board.Get(square)).GetGamePiece()
			playerStrength[player] += piece.GetStrength(g.Board, square, piecesLeft)
		}

		playerStrength[player] -= KingStrength
		// Account for the advantage of having a balanced pieces distribution over the teammates.
		playerStrength[player] = math.Pow(playerStrength[player], 0.8)
	}

	side := g.ActivePlayer.GetTeam().Side()

	return float64(side) * (playerStrength[0] + playerStrength[2] - playerStrength[1] - playerStrength[3])
}
