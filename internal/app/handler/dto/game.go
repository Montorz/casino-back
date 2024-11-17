package dto

type GameRequest struct {
	Name        string  `json:"name" binding:"required"`
	BetAmount   float64 `json:"bet_amount" binding:"required"`
	Coefficient float64 `json:"coefficient" binding:"required"`
}

type GameResponse struct {
	Id          int     `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	BetAmount   float64 `json:"bet_amount" binding:"required"`
	Coefficient float64 `json:"coefficient" binding:"required"`
	WinAmount   float64 `json:"win_amount" binding:"required"`
	CreatedDate string  `json:"created_date" binding:"required"`
}
