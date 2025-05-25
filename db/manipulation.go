package db

import (
	"context"

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
