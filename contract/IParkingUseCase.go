package contract

import "github.com/khafidprayoga/parking-app/internal/types"

type IParkingUseCase interface {
	EnterArea(request types.CarDTO) (areaId int, err error)
	LeaveArea(request types.CarDTO) (exitedCar types.Car, err error)
}
