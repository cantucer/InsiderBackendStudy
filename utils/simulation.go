package utils

import (
	"insiderbackendstudy/types"
	"math"
	"math/rand"
)

func SimulateMatch(match types.Match) (int, int) {
	// My super realistic simulation algorithm.
	// The position of ball is divided into 5 positions.
	// [ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
	// 0 is goal inside home's goal, 5 is starting position, 10 is goal inside away's goal. So left side is home, right side is away.
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

	ballPos := 5
	homeScore := 0
	awayScore := 0

	for range 90 {
		ballPos += randChoice()

		if ballPos == 0 {
			awayScore++
			ballPos = 5 // Reset ball position after a goal
		} else if ballPos == 10 {
			homeScore++
			ballPos = 5 // Reset ball position after a goal
		}
	}

	return homeScore, awayScore
}
