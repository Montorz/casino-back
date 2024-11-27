package dto

type AuthResponse struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}
