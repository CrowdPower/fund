package models

import (
	"time"
)

type Deposit struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Amount   int    `json:"amount"`
	Time     string `json:"time"`
}

type DepositArgs struct {
	Oldest    time.Time `query:"oldest"`
	Newest    time.Time `query:"newest"`
	MinAmount int       `query:"minamount"`
	MaxAmount int       `query:"maxamount"`
	Offset    int       `query:"offset"`
	Count     int       `query:"count"`
}
