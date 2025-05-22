package types

import "time"

type Car struct {
	Id           string     `json:"id"`
	AreaNumber   int        `json:"area_number"`
	Color        string     `json:"color"`
	PoliceNumber string     `json:"police_number"`
	ParkingAt    time.Time  `json:"parking_at"`
	ExitAt       *time.Time `json:"exit_at"`
	Cost         float64    `json:"cost"`
}

type CarDTO struct {
	RequestId    string `json:"request_id"`
	PoliceNumber string `json:"police_number"`
	Hours        int    `json:"hours,omitempty"`
}
