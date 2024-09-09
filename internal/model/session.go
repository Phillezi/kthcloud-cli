package model

import (
	"go-deploy/dto/v2/body"
	"time"
)

type AuthSession struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiryDate   time.Time `json:"expiry_date"`
}

type Session struct {
	User        body.UserRead `json:"user"`
	ApiKey      body.ApiKey   `json:"api_token,omitempty"`
	AuthSession AuthSession   `json:"auth_session"`
}
