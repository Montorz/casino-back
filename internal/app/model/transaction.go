package model

type Transaction struct {
	Id          int
	UserId      int `db:"user_id"`
	Type        string
	Amount      int
	CreatedDate string `db:"created_date"`
}
