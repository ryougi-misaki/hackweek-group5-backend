package models

import "time"

type ChatHistory struct {
	ID        int       `json:"id" form:"id" `
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	From      int       `json:"from" form:"from"`
	To        int       `json:"to" form:"to"`
	Msg       string    `json:"msg" form:"msg"`
}

