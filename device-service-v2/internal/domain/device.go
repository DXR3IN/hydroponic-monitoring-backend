package models

import "time"

type Device struct {
	ID         string    `json:"id"`
	DeviceName string    `json:"device_name"`
	OwnerID    string    `json:"owner_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
