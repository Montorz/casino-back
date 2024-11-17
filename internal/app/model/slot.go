package model

type Slot struct {
	Id     int
	Name   string
	MinBet int `db:"min_bet"`
	MaxBet int `db:"max_bet"`
}
