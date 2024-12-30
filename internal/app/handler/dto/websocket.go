package dto

type ChatMessage struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Message   string `json:"message"`
}

type GameResultMessage struct {
	Name   string      `json:"name"`
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
}
