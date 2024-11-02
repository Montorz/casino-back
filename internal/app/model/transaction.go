package model

type Transaction struct {
	Id     int    `json:"-" db:"id"`
	UserId int    `json:"user_id"`
	Type   string `json:"type" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}
