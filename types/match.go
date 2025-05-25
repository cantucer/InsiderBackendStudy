package types

type Match struct {
	HomeTeam   Team `json:"homeTeam"`
	AwayTeam   Team `json:"awayTeam"`
	IsPlayed   bool `json:"isPlayed"`
	HomeGoals  int `json:"homeGoals"`
	AwayGoals  int `json:"awayGoals"`
	Week       int `json:"week"`
}