package entity

import "time"

type DoorState struct {
	IsOpen    bool      `json:"is-open"`
	Timestamp time.Time `json:"timestamp"`
	DeviceID  string    `json:"deviceID"`
}
