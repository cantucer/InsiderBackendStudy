package utils

import (
	"insiderbackendstudy/types"
)

func PredictWinningChances(teams []types.Team, unplayedMatches []types.Match) map[string]float64 {
	// Runs simulations 100 times to predict winning chances.
	winCounts := make(map[string]int)
	for _, team := range teams {
		winCounts[team.Name] = 0
	}

	points := make(map[string]int)

	for range 1000000 {
		for _, team := range teams {
			points[team.Name] = team.Points
		}

		// Simulate each unplayed match.
		for _, match := range unplayedMatches {
			homeTeam := match.HomeTeam
			awayTeam := match.AwayTeam
			homeGoals, awayGoals := SimulateMatch(match)
			if homeGoals > awayGoals {
				points[homeTeam.Name] += 3 // Home team wins
			} else if awayGoals > homeGoals {
				points[awayTeam.Name] += 3 // Away team wins
			} else {
				points[homeTeam.Name] += 1 // Draw
				points[awayTeam.Name] += 1 // Draw
			}
		}

		// Count points for each team.
		maxPoints := 0
		winningTeam := ""
		for _, team := range teams {
			if points[team.Name] > maxPoints {
				maxPoints = points[team.Name]
				winningTeam = team.Name
			}
		}

		winCounts[winningTeam]++
	}

	chances := make(map[string]float64)

	for team := range winCounts {
		chances[team] = float64(winCounts[team]) / 10000.0
	}
	return chances
}
