package types

import "time"

type Car struct {
	Id           int        `json:"id"`
	Color        string     `json:"color"`
	PoliceNumber string     `json:"police_number"`
	ParkingAt    time.Time  `json:"parking_at"`
	ExitAt       *time.Time `json:"exit_at"`
	Cost         float64    `json:"cost"`
}

type CarDTO struct {
	PoliceNumber string `json:"police_number"`

	Hours int `json:"hours"`
}
