package models

import "time"

type Session struct {
	ID        string
	Username  string
	Subject   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionInfo struct {
	SessionID string
	Subject   string
	Username  string
}
