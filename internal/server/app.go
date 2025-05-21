package server

type ParkingAppServer struct {
	service ParkingServiceImpl
}

func CreateAppServer(service ParkingServiceImpl) *ParkingAppServer {
	return &ParkingAppServer{
		service: service,
	}
}
