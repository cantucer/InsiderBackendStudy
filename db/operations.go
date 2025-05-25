package db

import (
	"context"
	"insiderbackendstudy/types"

	"github.com/jackc/pgx/v5"
)

func AddTeam(conn *pgx.Conn, name string, strength int) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO teams (name, strength)
		VALUES ($1, $2)
		ON CONFLICT (name) DO NOTHING;
	`, name, strength)
	if err != nil {
		return err
	}
	return nil
}

func AddMatch(conn *pgx.Conn, homeTeam, awayTeam string, week int) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO matches (homeTeam, awayTeam, week, homeGoals, awayGoals, isPlayed)
		VALUES ($1, $2, $3, 0, 0, FALSE)
		ON CONFLICT (homeTeam, awayTeam) DO NOTHING;
	`, homeTeam, awayTeam, week)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMatchResult(conn *pgx.Conn, homeTeam, awayTeam string, homeGoals, awayGoals int) error {
	_, err := conn.Exec(context.Background(), `
		UPDATE matches
		SET homeGoals = $1, awayGoals = $2, isPlayed = TRUE
		WHERE homeTeam = $3 AND awayTeam = $4;
	`, homeGoals, awayGoals, homeTeam, awayTeam)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTeamStats(conn *pgx.Conn, teamName string, goalsFor, goalsAgainst int, winChange, drawChange, loseChange int) error {
	_, err := conn.Exec(context.Background(), `
		UPDATE teams
		SET goalsFor = goalsFor + $1,
		    goalsAgainst = goalsAgainst + $2,
		    won = won + $3,
		    drawn = drawn + $4,
		    lost = lost + $5
		WHERE name = $6;
	`, goalsFor, goalsAgainst, winChange, drawChange, loseChange, teamName)
	if err != nil {
		return err
	}
	return nil
}

func GetTeams(conn *pgx.Conn) ([]types.Team, error) {
	rows, err := conn.Query(context.Background(), `
		SELECT name, strength, points, played, won, drawn, lost, goalsFor, goalsAgainst
		FROM teams
		ORDER BY points DESC, goalsFor - goalsAgainst DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []types.Team
	for rows.Next() {
		var team types.Team
		if err := rows.Scan(&team.Name, &team.Strength, &team.Points, &team.Played,
			&team.Won, &team.Drawn, &team.Lost, &team.GoalsFor, &team.GoalsAgainst); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func GetAllMatches(conn *pgx.Conn) ([]types.Match, error) {
	rows, err := conn.Query(context.Background(), `
		SELECT homeTeam, awayTeam, isPlayed, homeGoals, awayGoals, week
		FROM matches
		ORDER BY week;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []types.Match
	for rows.Next() {
		var match types.Match
		if err := rows.Scan(&match.HomeTeam.Name, &match.AwayTeam.Name, &match.IsPlayed,
			&match.HomeGoals, &match.AwayGoals, &match.Week); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func GetAllUnplayedMatches(conn *pgx.Conn) ([]types.Match, error) {
	rows, err := conn.Query(context.Background(), `
		SELECT homeTeam, awayTeam, isPlayed, homeGoals, awayGoals, week
		FROM matches
		WHERE isPlayed = FALSE
		ORDER BY week;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []types.Match
	for rows.Next() {
		var match types.Match
		if err := rows.Scan(&match.HomeTeam.Name, &match.AwayTeam.Name, &match.IsPlayed,
			&match.HomeGoals, &match.AwayGoals, &match.Week); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func GetMatches(conn *pgx.Conn, week int) ([]types.Match, error) {
	rows, err := conn.Query(context.Background(), `
		SELECT homeTeam, awayTeam, isPlayed, homeGoals, awayGoals
		FROM matches
		WHERE week = $1
		ORDER BY homeTeam, awayTeam;
	`, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []types.Match
	for rows.Next() {
		var match types.Match
		if err := rows.Scan(&match.HomeTeam.Name, &match.AwayTeam.Name, &match.IsPlayed,
			&match.HomeGoals, &match.AwayGoals); err != nil {
			return nil, err
		}
		match.Week = week
		matches = append(matches, match)
	}

	return matches, nil
}

func GetUnplayedMatches(conn *pgx.Conn, week int) ([]types.Match, error) {
	rows, err := conn.Query(context.Background(), `
		SELECT homeTeam, awayTeam, isPlayed, homeGoals, awayGoals
		FROM matches
		WHERE week = $1 AND isPlayed = FALSE
		ORDER BY homeTeam, awayTeam;
	`, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []types.Match
	for rows.Next() {
		var match types.Match
		if err := rows.Scan(&match.HomeTeam.Name, &match.AwayTeam.Name, &match.IsPlayed,
			&match.HomeGoals, &match.AwayGoals); err != nil {
			return nil, err
		}
		match.Week = week
		matches = append(matches, match)
	}

	return matches, nil
}
