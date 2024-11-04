package model

type Slot struct {
	Id     int    `json:"-" db:"id"`
	Name   string `json:"name"  db:"name" binding:"required"`
	MinBet int    `json:"min_bet" db:"min_bet" binding:"required"`
	MaxBet int    `json:"max_bet" db:"max_bet" binding:"required"`
}
