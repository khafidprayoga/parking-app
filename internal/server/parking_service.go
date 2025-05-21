package server

import "github.com/khafidprayoga/parking-app/internal/types"

type ParkingServiceImpl struct {
	LotCapacity int64
	Store       *[]types.Car
}

func (p ParkingServiceImpl) EnterArea(policeNumber string) error {
	//TODO implement me
	panic("implement me")
}

func (p ParkingServiceImpl) LeaveArea(policeNumber string) error {
	//TODO implement me
	panic("implement me")
}

func NewParkingService() *ParkingServiceImpl {
	// todo init store db data
	return &ParkingServiceImpl{}
}

func (p ParkingServiceImpl) pay() {

}
