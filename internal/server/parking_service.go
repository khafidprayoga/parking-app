package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/khafidprayoga/parking-app/internal/types"
)

type ParkingServiceImpl struct {
	lotCapacity int          `json:"lot_capacity"`
	store       []*types.Car `json:"store"`
	revenue     float64      `json:"revenue"`
}

func (p *ParkingServiceImpl) Status() (_ []byte, err error) {
	state := types.AppStatus{
		Revenue:            p.revenue,
		LotParkingCapacity: p.lotCapacity,
		TxCount:            100_0,
		CarList:            p.store,
	}

	dataBytes, errMarshall := json.Marshal(state)
	if errMarshall != nil {
		err = fmt.Errorf("failed to marshall parking data")
		return
	}

	return dataBytes, nil
}

func (p *ParkingServiceImpl) OpenParkingArea(parkingCap int) (err error) {

	if p.lotCapacity > 0 {
		err = fmt.Errorf("failed, already initalize the parking lot capacity")
		return
	}

	p.lotCapacity = parkingCap
	p.store = make([]*types.Car, parkingCap)

	return
}

func (p *ParkingServiceImpl) EnterArea(request types.CarDTO) (areaId int, err error) {

	// validate if  car number not already exist on the parking area
	for _, car := range p.store {
		if car != nil {
			if strings.EqualFold(car.PoliceNumber, request.PoliceNumber) {
				err = fmt.Errorf("failed, already parked the parking lot capacity")
				return
			}
		}
	}

	for index, car := range p.store {
		// allocating nearest parking lot from the door gateway
		if car == nil {
			id := index + 1
			p.store[index] = &types.Car{
				Id:           request.RequestId,
				AreaNumber:   id,
				PoliceNumber: request.PoliceNumber,
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

	for i, car := range p.store {
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
	p.revenue = p.revenue + carDetail.Cost
	p.store[carIndex] = nil

	exitedCar = carDetail
	return
}

func NewParkingService() *ParkingServiceImpl {
	return &ParkingServiceImpl{}
}

func (p *ParkingServiceImpl) pay() {
	// todo

}
