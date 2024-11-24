package dto

type BetResponse struct {
	BetAmount int `json:"bet_amount" binding:"required"`
}

type GameRequest struct {
	Name        string  `json:"name" binding:"required"`
	BetAmount   float64 `json:"bet_amount" binding:"required"`
	Coefficient float64 `json:"coefficient"` // не используем binding:"required" так как может быть равен 0 при проигрыше
}

type GameResponse struct {
	Id          int     `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	BetAmount   float64 `json:"bet_amount" binding:"required"`
	Coefficient float64 `json:"coefficient" binding:"required"`
	WinAmount   float64 `json:"win_amount" binding:"required"`
	CreatedDate string  `json:"created_date" binding:"required"`
}
