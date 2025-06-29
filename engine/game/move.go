package game

import "fmt"

// Move stores move coordinates.
type Move struct {
	From Square `json:"from"`
	To   Square `json:"to"`
}

// String implements the Stringer interface.
func (m Move) String() string {
	return fmt.Sprintf("%v->%v", m.From, m.To)
}
