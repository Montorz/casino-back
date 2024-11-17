package model

type Game struct {
	Id          int
	UserId      int `db:"user_id"`
	SlotId      int `db:"slot_id"`
	Name        string
	BetAmount   int `db:"bet_amount"`
	Coefficient float64
	WinAmount   int    `db:"win_amount"`
	CreatedDate string `db:"created_date"`
}
