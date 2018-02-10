package models

import (
	"time"
)

type Payment struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Amount   int    `json:"amount"`
	Time     string `json:"time"`
	Url      string `json:"url"`
}

type PaymentArgs struct {
	Oldest    time.Time `query:"oldest"`
	Newest    time.Time `query:"newest"`
	MinAmount int       `query:"minamount"`
	MaxAmount int       `query:"maxamount"`
	Url       string    `query:"url"`
	Offset    int       `query:"offset"`
	Count     int       `query:"count"`
}
