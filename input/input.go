package input

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/vpoliakov01/2v2ChessAI/game"
)

var moveRegex = regexp.MustCompile("([a-n])([0-9][0-9]?){1,2}([a-n])([0-9][0-9]?)")

func ReadMove() (*game.Move, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the move in format e2e4: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	input = strings.ToLower(input)
	matches := moveRegex.FindStringSubmatch(input)

	fmt.Printf("Parsed %v\n", matches)
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
