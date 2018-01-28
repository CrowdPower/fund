package models

type Payment struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Amount   int    `json:"amount"`
	Time     string `json:"time"`
	Url      string `json:"url"`
}
