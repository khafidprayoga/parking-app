package server

import (
	"fmt"
	"strconv"

	"github.com/khafidprayoga/parking-app/internal/types"
)

func (srv *ParkingAppServer) HandleIncomingMsg(msg types.Socket) (response string, err error) {
	switch msg.Command {
	case types.CmdCreateStore:
		parkingCap, errCv := strconv.Atoi(msg.Data.(string))
		if errCv != nil {
			err = fmt.Errorf("failed to convert string to int at %s actions", msg.Command)
			return
		}

		errOpen := srv.service.OpenParkingArea(parkingCap)
		if err != nil {
			err = fmt.Errorf("failed to open parking area: %s", errOpen.Error())
			return
		}

		response = fmt.Sprintf("success initalize parking lot with %v capacity", parkingCap)
		return
	case types.CmdPark:
		incomingCarData := types.CarDTO{
			RequestId:    msg.XRequestId,
			PoliceNumber: msg.Data.(map[string]any)["police_number"].(string),
		}

		areaId, errParking := srv.service.EnterArea(incomingCarData)
		if errParking != nil {
			err = fmt.Errorf("failed to enter area, %s", errParking.Error())
			return
		}

		response = fmt.Sprintf(
			"successfully parked car. with police number %s and SLOT number id %v",
			incomingCarData.PoliceNumber,
			areaId,
		)
		return
	case types.CmdLeave:
		incomingCarData := types.CarDTO{
			PoliceNumber: msg.Data.(map[string]any)["police_number"].(string),
			Hours:        int(msg.Data.(map[string]any)["hours"].(float64)),
		}

		metadata, errLeave := srv.service.LeaveArea(incomingCarData)
		if errLeave != nil {
			err = fmt.Errorf("failed to exit area with police id %s, %s", incomingCarData.PoliceNumber, errLeave.Error())
			return
		}

		response = fmt.Sprintf(
			"successfully leave car. with police number %s and total hours elapsed  %v on area number %d",
			metadata.PoliceNumber,
			incomingCarData.Hours,
			metadata.AreaNumber,
		)
		return
	case types.CmdStatus:
		dataBytes, errGetStatus := srv.service.Status()
		if errGetStatus != nil {
			err = fmt.Errorf("failed to parking app status %s", errGetStatus.Error())
			return
		}

		response = string(dataBytes)
		return
	}

	return "", nil
}
