package model

import "time"

type Device struct {
	DeviceID string
	Paired   bool
	Display  bool
	Created  time.Time
	Name     string
}
