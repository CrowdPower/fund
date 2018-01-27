package models

type Deposit struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Amount   int    `json:"amount"`
	Time     string `json:"time"`
}
