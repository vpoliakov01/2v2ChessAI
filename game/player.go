package game

type Player int

type Team int // Red/Yellow: 1, Blue/Green: -1.

// IsTeamMate returns true if p and other are on the same team (including p == other).
func (p Player) IsTeamMate(other Player) bool {
	return (p^other)&1 == 0 // Last bit must match.
}

// Team returns 1 for Red/Yellow and -1 for Blue/Green.
func (p Player) Team() Team {
	t := p & 1
	return Team(1 - 2*t)
}

// Opposite returns the opposite team.
func (t Team) Opposite() Team {
	return t * -1
}

// String implements the Stringer interface.
func (t Team) String() string {
	switch t {
	case 1:
		return "Red/Yellow"
	case -1:
		return "Blue/Green"
	default:
		panic("unsupported team")
	}
}
