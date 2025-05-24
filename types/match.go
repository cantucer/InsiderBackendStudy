package types

type Match struct {
	homeTeam   Team `json:"homeTeam"`
	awayTeam   Team `json:"awayTeam"`
	isPlayed   bool `json:"isPlayed"`
	homeGoals  int `json:"homeGoals"`
	awayGoals  int `json:"awayGoals"`

	// 0 if draw, 1 if home wins, 2 if away wins
	result     int `json:"result"`
}