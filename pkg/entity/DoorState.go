package entity

import "time"

type KeyState struct {
	Open      bool      `json:"open"`
	Timestamp time.Time `json:"timestamp"`
	DeviceID  string    `json:"deviceID"`
}

type DoorState struct {
	IsOpen    bool      `json:"is-open"`
	Timestamp time.Time `json:"timestamp"`
	DeviceID  string    `json:"deviceID"`
}
