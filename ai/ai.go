package ai

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"sort"

	"github.com/vpoliakov01/2v2ChessAI/game"
)

var (
	ErrGameEnded = errors.New("the game has ended")
	ErrNoMoves   = errors.New("no move can be made in this position")
)

type gameScore struct {
	game  *game.Game
	score float64
}

type playerStrength struct {
	player   game.Player
	strength float64
}

type AI struct {
	Depth int
	Cache map[uint64]float64 // Stores scores of evalutated positions.
}

func init() {
	cpus := runtime.NumCPU()
	fmt.Printf("Running on %v CPUs\n", cpus)
	runtime.GOMAXPROCS(cpus)
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

	moveEvalEstimates := map[game.Move]gameScore{}

	for i := range moves {
		gameCopy := g.Copy()
		gameCopy.Play(moves[i])
		moveEvalEstimates[moves[i]] = gameScore{gameCopy, ai.EvaluateCurrent(gameCopy)}
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
	playerStrengths := map[game.Player]float64{}
	piecesLeft := 0

	for player := range g.Board.PieceSquares {
		piecesLeft += g.Board.PieceSquares[player].Size()
	}

	c := make(chan playerStrength, piecesLeft)

	for player := range g.Board.PieceSquares {
		for _, sq := range g.Board.PieceSquares[player].Elements() {
			go func(player game.Player, square game.Square) {
				piece := game.Piece(g.Board.Get(square)).GetGamePiece()
				c <- playerStrength{player, piece.GetStrength(g.Board, square, piecesLeft)}
			}(player, sq.(game.Square))
		}
	}

	for ; piecesLeft > 0; piecesLeft-- {
		ps := <-c
		playerStrengths[ps.player] += ps.strength
	}

	// Account for the advantage of having a balanced pieces distribution between teammates.
	for player := range g.Board.PieceSquares {
		if playerStrengths[player] > 0 {
			playerStrengths[player] = math.Pow(playerStrengths[player], 0.8)
		}
	}

	score := float64(g.ActivePlayer.GetTeam()) * (playerStrengths[0] + playerStrengths[2] - playerStrengths[1] - playerStrengths[3])

	return score
}
