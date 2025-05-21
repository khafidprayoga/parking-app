package contract

import "github.com/khafidprayoga/parking-app/internal/types"

type IParkingUseCase interface {
	EnterArea(policeNumber string) (areaId int, err error)
	LeaveArea(request types.CarDTO) (exitedCar types.Car, err error)
}
