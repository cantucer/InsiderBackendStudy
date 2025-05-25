package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateTables(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS teams (
			name TEXT PRIMARY KEY,
			strength INT NOT NULL,
			points INT GENERATED ALWAYS AS (3 * won + drawn) STORED,
			played INT GENERATED ALWAYS AS (won + drawn + lost) STORED,
			won INT NOT NULL DEFAULT 0,
			drawn INT NOT NULL DEFAULT 0,
			lost INT NOT NULL DEFAULT 0,
			goalsFor INT NOT NULL DEFAULT 0,
			goalsAgainst INT NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS matches (
			homeTeam TEXT NOT NULL REFERENCES teams(name),
			awayTeam TEXT NOT NULL REFERENCES teams(name),
			isPlayed BOOLEAN NOT NULL DEFAULT FALSE,
			homeGoals INT NOT NULL,
			awayGoals INT NOT NULL,
			week INT NOT NULL,
			result INT GENERATED ALWAYS AS (
				CASE
					WHEN homeGoals > awayGoals THEN 1
					WHEN homeGoals < awayGoals THEN 2
					ELSE 0
				END
			) STORED,
			PRIMARY KEY (homeTeam, awayTeam)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func ClearTables(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), `
		TRUNCATE TABLE teams, matches RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		return err
	}

	return nil
}
