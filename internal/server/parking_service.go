package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/khafidprayoga/parking-app/internal/types"
)

type ParkingServiceImpl struct {
	LotCapacity int          `json:"lot_capacity"`
	Store       []*types.Car `json:"store"`
	Revenue     float64      `json:"revenue"`
}

func (p *ParkingServiceImpl) EnterArea(policeNumber string) (areaId int, err error) {

	// validate if  car number not already exist on the parking area
	for _, car := range p.Store {
		if car != nil {
			if strings.EqualFold(car.PoliceNumber, policeNumber) {
				err = fmt.Errorf("failed, already parked the parking lot capacity")
				return
			}
		}
	}

	for index, car := range p.Store {
		// allocating nearest parking lot from the door gateway
		if car == nil {
			id := index + 1
			p.Store[index] = &types.Car{
				Id:           id,
				PoliceNumber: policeNumber,
				ParkingAt:    time.Now(),
				ExitAt:       nil,
			}

			areaId = id
			break
		}
	}

	return
}

func (p *ParkingServiceImpl) LeaveArea(req types.CarDTO) (exitedCar types.Car, err error) {
	carDetail := types.Car{}
	carIndex := -1

	for i, car := range p.Store {
		if car != nil {
			if strings.EqualFold(car.PoliceNumber, req.PoliceNumber) {
				carDetail = *car
				carIndex = i
				break
			}
		}
	}

	if carIndex < 0 {
		err = fmt.Errorf("failed, car with police number %s does not exist", req.PoliceNumber)
		return
	}

	start := carDetail.ParkingAt
	end := start.Add(time.Duration(req.Hours) * time.Hour)
	carDetail.ExitAt = &end
	carDetail.Cost = 50.0

	// flush
	p.Revenue = p.Revenue + carDetail.Cost
	p.Store[carIndex] = nil

	exitedCar = carDetail
	return
}

func NewParkingService() *ParkingServiceImpl {
	return &ParkingServiceImpl{}
}

func (p *ParkingServiceImpl) pay() {
	// todo

}
