package types

type Team struct {
	Name         string `json:"name"`
	Strength     int    `json:"strength"`
	Points       int    `json:"points"`
	Played       int    `json:"played"`
	Won          int    `json:"won"`
	Drawn        int    `json:"drawn"`
	Lost         int    `json:"lost"`
	GoalsFor     int    `json:"goalsFor"`
	GoalsAgainst int    `json:"goalsAgainst"`
}
