package types

type Team struct {
	name         string `json:"name"`
	strength     int    `json:"strength"`
	played       int    `json:"played"`
	won          int    `json:"won"`
	drawn        int    `json:"drawn"`
	lost         int    `json:"lost"`
	goalsFor     int    `json:"goalsFor"`
	goalsAgainst int    `json:"goalsAgainst"`
}