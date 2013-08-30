package webqq

import (
	"time"
)

type WebQQSingle struct {
	Account  Account
	Uin      Uin
	LastTalk time.Time
}

type WebQQGroup struct {
}

type WebQQCustomPic struct {
}
