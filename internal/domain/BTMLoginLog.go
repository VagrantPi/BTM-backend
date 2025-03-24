package domain

import "time"

type BTMLoginLog struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	UserName  string    `json:"user_name"`
	IP        string    `json:"ip"`
	Browser   string    `json:"browser"`
	CreatedAt time.Time `json:"created_at"`
}
