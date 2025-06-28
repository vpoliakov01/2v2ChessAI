package io

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"

	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

func Save(g *game.Game) (string, error) {
	bytes, err := g.JSON()
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(bytes)
	file := fmt.Sprintf("%x.save", hash[0:2])

	err = ioutil.WriteFile(file, bytes, 0644)
	if err != nil {
		return "", err
	}

	return file, nil
}
