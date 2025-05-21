package types

import "time"

type Car struct {
	Color        string     `json:"color"`
	PoliceNumber string     `json:"police_number"`
	ParkingAt    *time.Time `json:"parking_at"`
	ExitAt       *time.Time `json:"exit_at"`
}
