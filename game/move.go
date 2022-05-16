package game

import "fmt"

type Move struct {
	From Square
	To   Square
}

func (m Move) String() string {
	return fmt.Sprintf("%v->%v", m.From, m.To)
}
