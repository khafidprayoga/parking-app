package types

import "time"

type Car struct {
	ID           int
	Color        string     `json:"color"`
	PoliceNumber string     `json:"police_number"`
	ParkingAt    time.Time  `json:"parking_at"`
	ExitAt       *time.Time `json:"exit_at"`
}

type CarDTO struct {
	PoliceNumber string `json:"police_number"`
	Hours        int    `json:"hours"`
}
