package model

type User struct {
	Default
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
