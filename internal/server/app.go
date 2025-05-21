package server

type ParkingAppServer struct {
	service ParkingServiceImpl
}

func NewParkingAppServer(service ParkingServiceImpl) *ParkingAppServer {
	return &ParkingAppServer{
		service: service,
	}
}
