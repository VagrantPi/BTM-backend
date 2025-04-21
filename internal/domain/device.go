package domain

import "time"

type Device struct {
	DeviceID string    `json:"device_id"`
	Paired   bool      `json:"paired"`
	Display  bool      `json:"display"`
	Created  time.Time `json:"created"`
	Name     string    `json:"name"`
}
