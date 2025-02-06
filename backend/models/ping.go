package models

import "time"

type Ping struct {
	ID        int       `json:"id"`
	IPAddress string    `json:"ip_address"`
	PingTime  time.Time `json:"ping_time"`
	IsSuccess bool      `json:"is_success"`
}
