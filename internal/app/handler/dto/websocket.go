package dto

type ChatMessage struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Message   string `json:"message"`
}
