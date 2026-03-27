package models

import "time"

type SessionClaims struct {
	SessionID string    `json:"sid"`
	UserID    int64     `json:"uid"`
	ExpiresAt time.Time `json:"exp"`
}
