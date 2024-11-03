package model

type Transaction struct {
	Id          int    `json:"-" db:"id"`
	UserId      int    `json:"-" db:"user_id"`
	Type        string `json:"type" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	CreatedDate string `json:"created_date" db:"created_date"`
}
