package dto

type BetResponse struct {
	BetAmount int `json:"bet_amount" binding:"required"`
}
