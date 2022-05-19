package input

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/vpoliakov01/2v2ChessAI/game"
)

var moveRegex = regexp.MustCompile("[A-Z]?([a-n])([0-9][0-9]?){1,2}-?([a-n])([0-9][0-9]?)")

// ReadMove reads a move from STDIN.
func ReadMove() (*game.Move, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the move in format e2e4: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return ParseMove(input)
}

// LoadPGN returns the moves (pgn notation) specified in the file.
func LoadPGN(file string) ([]game.Move, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return ParsePGN(string(bytes))
}

// ParseMove parses a move from a string.
func ParseMove(m string) (*game.Move, error) {
	matches := moveRegex.FindStringSubmatch(m)

	if len(matches) < 5 {
		return nil, fmt.Errorf("bad input %v", matches)
	}

	fromFile := int(matches[1][0]) - int('a')
	fromRank, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, err
	}
	fromRank--

	toFile := int(matches[3][0]) - int('a')
	toRank, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, err
	}
	toRank--

	return &game.Move{
		From: game.Square{fromRank, fromFile},
		To:   game.Square{toRank, toFile},
	}, nil
}

// ParsePGN parses pgn from a string.
func ParsePGN(input string) ([]game.Move, error) {
	lines := strings.Split(input, "\n")
	moves := []game.Move{}

	for _, line := range lines {
		if len(line) < 4 {
			continue
		}

		line = line[3:]

		for _, moveStr := range strings.Split(line, " .. ") {
			move, err := ParseMove(moveStr)
			if err != nil {
				return nil, err
			}
			moves = append(moves, *move)
		}

	}
	return moves, nil
}
