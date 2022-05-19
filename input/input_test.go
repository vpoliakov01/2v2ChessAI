package input_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	. "github.com/vpoliakov01/2v2ChessAI/input"
)

type TestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestParsePGN() {
	r := s.Require()

	moves, err := ParsePGN(`
1. h2-h3 .. b7-c7 .. g13-g12 .. m8-l8
2. f2-f3 .. b9-c9 .. Qh14-e11 .. Qn7-k10
3. Bi1-h2 .. b5-c5 .. Qe11-h11 .. Qk10-g6
4. Qg1-e3 .. Qa8-c6 .. Qh11-j11 .. Qg6-h7
5. Qe3-j8 .. Qc6-g6 .. Bf14-h12 .. m10-k10
6. d2-d3 .. b11-c11 .. Qj11-d5 .. m6-l6
7. e2-e4 .. b6-c6 .. Qd5-i5 .. Nn10-l9
`)
	r.NoError(err)
	fmt.Println(moves)
}
