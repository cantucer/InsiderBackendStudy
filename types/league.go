package types

type League struct {
	teams []Team `json:"teams"`
	matches []Match `json:"matches"`
}