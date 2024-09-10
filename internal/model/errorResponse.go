package model

type Error struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}
