package simulation

import (
	"insiderbackendstudy/types"
)

func CreateFixture(teams []types.Team) [][]types.Match {
	weekCount := (len(teams) - 1) * 2 // Total weeks in a round-robin tournament.

	// Assuming team number is even.
	matchCountPerWeek := len(teams) / 2

	matches := make([]types.Match, weekCount*matchCountPerWeek)
	weeks := make([][]types.Match, weekCount)

	for week := range weekCount / 2 {
		// Using the circle rotation method.
		fixedTeam := teams[0]

		for match := range matchCountPerWeek {
			homeTeam := teams[match]
			awayTeam := teams[len(teams)-1-match]

			matches[week*matchCountPerWeek+match] = types.Match{
				HomeTeam:  homeTeam,
				AwayTeam:  awayTeam,
				IsPlayed:  false,
				HomeGoals: 0,
				AwayGoals: 0,
				Week:      week + 1,
			}
		}

		weeks[week] = matches[week*matchCountPerWeek : (week+1)*matchCountPerWeek]

		// Rotate teams for next week.
		rotatingTeams := teams[1:]
		lastTeam := rotatingTeams[len(rotatingTeams)-1:]
		rotatedTeams := append(lastTeam, rotatingTeams[:len(rotatingTeams)-1]...)

		teams = append([]types.Team{fixedTeam}, rotatedTeams...)
	}

	// Copying from the first half of the fixture.

	for week := weekCount / 2; week < weekCount; week++ {

		for match := range matchCountPerWeek {
			homeTeam := matches[(weekCount-week-1)*matchCountPerWeek+match].HomeTeam
			awayTeam := matches[(weekCount-week-1)*matchCountPerWeek+match].AwayTeam
			matches[week*matchCountPerWeek+match] = types.Match{
				HomeTeam:  awayTeam,
				AwayTeam:  homeTeam,
				IsPlayed:  false,
				HomeGoals: 0,
				AwayGoals: 0,
				Week:      week + 1,
			}
		}

		weeks[week] = matches[week*matchCountPerWeek : (week+1)*matchCountPerWeek]
	}

	return weeks
}
