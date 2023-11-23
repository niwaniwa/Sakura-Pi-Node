package entity

import "time"

type Card struct {
	ID        []byte    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}
