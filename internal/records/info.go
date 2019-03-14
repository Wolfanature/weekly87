package records

import "github.com/oakmound/weekly87/internal/characters/players"

// RunInfo is a placholder for info that we display after game end.
//
type RunInfo struct {
	Party           []*players.Player
	SectionsCleared int
	EnemiesDefeated int
}