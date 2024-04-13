package entity

import "time"

type DoorSwitchState struct {
	IsOpen    bool      `json:"is-open"`
	Timestamp time.Time `json:"timestamp"`
	DeviceID  string    `json:"deviceID"`
}
