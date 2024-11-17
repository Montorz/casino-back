package dto

type TransactionRequest struct {
	Type   string `json:"type" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

type TransactionResponse struct {
	Id          int    `json:"id" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	CreatedDate string `json:"created_date" binding:"required"`
}
