package contract

type IParkingUseCase interface {
	EnterArea(policeNumber string) error
	LeaveArea(policeNumber string) error
}
