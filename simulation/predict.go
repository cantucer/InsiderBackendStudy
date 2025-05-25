package simulation

import (
	"insiderbackendstudy/types"
	"math"
	"math/rand"
)

func PredictMatch(match types.Match) types.Match {
	// My super realistic prediction algorithm.
	// The position of ball is divided into 5 positions.
	// [ 0, 1, 2, 3, 4, 5, 6, 7, 8]
	// 0 is goal inside home's goal, 4 is starting position, 8 is goal inside away's goal. So left side is home, right side is away.
	// Others are intermediate positions.
	// Movement of ball is determined by the strength of teams + 5 advantage +- 15 luck. But minimum strength is 10.
	// Game is played for 90 iterations.

	homeStrength := math.Max(10.0, float64(match.HomeTeam.Strength)+5.0+(30*rand.Float64()-15.0))
	awayStrength := math.Max(10.0, float64(match.AwayTeam.Strength)+(30*rand.Float64()-15.0))

	totalStrength := homeStrength + awayStrength
	chanceOfMovementToRight := homeStrength / totalStrength

	randChoice := func() int {
		if rand.Float64() < chanceOfMovementToRight {
			return 1 // Move to the right
		}
		return -1 // Move to the left
	}

	ballPos := 4
	homeScore := 0
	awayScore := 0

	for range 90 {
		ballPos += randChoice()

		if ballPos == 0 {
			awayScore++
			ballPos = 4 // Reset ball position after a goal
		} else if ballPos == 9 {
			homeScore++
			ballPos = 4 // Reset ball position after a goal
		}
	}

	return types.Match{
		HomeTeam:  match.HomeTeam,
		AwayTeam:  match.AwayTeam,
		IsPlayed:  true,
		HomeGoals: homeScore,
		AwayGoals: awayScore,
		Week:      match.Week,
	}
}
