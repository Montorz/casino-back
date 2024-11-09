package model

type Game struct {
	Id          int     `json:"-" db:"id"`
	UserId      int     `json:"-" db:"user_id"`
	SlotId      int     `json:"-" db:"slot_id"`
	Name        string  `json:"name" db:"name" binding:"required"`
	BetAmount   int     `json:"bet_amount" db:"bet_amount" binding:"required"`
	Coefficient float64 `json:"coefficient" db:"coefficient" binding:"required"`
	WinAmount   int     `json:"win_amount" db:"win_amount"`
	CreatedDate string  `json:"created_date" db:"created_date"`
}
