package play

import (
	"encoding/json"
	"fmt"

	"github.com/vpoliakov01/2v2ChessAI/engine/game"
)

// CastData marshals data and unmarshals it into the given type.
func CastData[T any](data interface{}) (T, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return *new(T), fmt.Errorf("error marshalling data: %v", err)
	}
	var unmarshalled T
	err = json.Unmarshal(bytes, &unmarshalled)
	if err != nil {
		return *new(T), fmt.Errorf("error unmarshalling data: %v", err)
	}
	return unmarshalled, nil
}

// AreHumanPlayersEqual checks if two slices of game.Player are equal.
func AreHumanPlayersEqual(a, b []game.Player) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
