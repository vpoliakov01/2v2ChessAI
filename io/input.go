package io

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

// ReadInput reads user io from STDIN.
func ReadInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a command or a move in format e2e4: ")

	in, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.Trim(in, " \n\t"), nil
}

// Load attempts to load a game from json and if it fails,
// it attempts to load it from pgn.
func Load(file string) (*game.Game, error) {
	g, err := LoadJSON(file)
	if err == nil {
		return g, nil
	}

	fmt.Printf("Failed to load json: %v\n", err)
	fmt.Println("Attempting to load pgn")

	g = game.New()
	moves, err := LoadPGN(file)
	if err != nil {
		return nil, err
	}

	for _, move := range moves {
		g.Play(move)
	}

	return g, nil
}

// LoadPGN returns the moves (pgn notation) specified in the file.
func LoadPGN(file string) ([]game.Move, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return ParsePGN(string(bytes))
}

// LoadJSON the game from json.
func LoadJSON(file string) (*game.Game, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return game.LoadJSON(bytes)
}

// ParseMove parses a move from a string.
func ParseMove(m string) (*game.Move, error) {
	matches := moveRegex.FindStringSubmatch(m)

	if len(matches) < 5 {
		return nil, fmt.Errorf("bad io %v", matches)
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
func ParsePGN(io string) ([]game.Move, error) {
	lines := strings.Split(io, "\n")
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
