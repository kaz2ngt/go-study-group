package model

type Request struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

type Response struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}
