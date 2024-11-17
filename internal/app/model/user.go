package model

type User struct {
	Id        int
	Name      string
	Login     string
	Password  string
	Balance   int
	AvatarURL string `db:"avatar_url"`
}
