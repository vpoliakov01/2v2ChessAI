package game

type Player int

// IsTeamMate returns true if p and other are on the same team (including p == other).
func (p Player) IsTeamMate(other Player) bool {
	return (p^other)&1 == 0 // Last bit must match.
}

func (p Player) GetTeam() int {
	return int(p) & 1
}
