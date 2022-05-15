package game

type Player int

type Team int // 0, 1.
type Side int // 1, -1.

// IsTeamMate returns true if p and other are on the same team (including p == other).
func (p Player) IsTeamMate(other Player) bool {
	return (p^other)&1 == 0 // Last bit must match.
}

func (p Player) GetTeam() Team {
	return Team(p & 1) // Last bit.
}

func (t Team) Opposite() Team {
	return t ^ 1
}

func (t Team) Side() Side {
	return Side(1 - 2*t)
}

func (s Side) Team() Team {
	return Team((1 - s) / 2)
}
