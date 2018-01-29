package models

type User struct {
	Username          string `json:"username"`
	Password          string `json:"password,omitempty"`
	Email             string `json:"email"`
	Balance           int    `json:"balance"`
	InvalidatedTokens bool   `json:"-"`
}
