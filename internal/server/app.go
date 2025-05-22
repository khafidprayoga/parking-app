package server

import "github.com/khafidprayoga/parking-app/contract"

type ParkingAppServer struct {
	service contract.IParkingUseCase
}

func CreateAppServer(service contract.IParkingUseCase) *ParkingAppServer {
	return &ParkingAppServer{
		service: service,
	}
}
