package entity

import "time"

type DoorSwitchState struct {
	IsOpen    bool      `json:"is_open"`
	Timestamp time.Time `json:"timestamp"`
	DeviceID  string    `json:"device_id"`
}
