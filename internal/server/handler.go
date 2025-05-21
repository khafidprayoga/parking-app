package server

import (
	"encoding/json"
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

		if srv.service.LotCapacity > 0 {
			err = fmt.Errorf("failed, already initalize the parking lot capacity")
			return
		}

		srv.service.LotCapacity = parkingCap
		srv.service.Store = make([]*types.Car, parkingCap)

		response = fmt.Sprintf("success initalize parking lot with %v capacity", parkingCap)
		return
	case types.CmdPark:
		incomingCarData := types.Car{
			PoliceNumber: msg.Data.(map[string]any)["police_number"].(string),
		}

		areaId, errParking := srv.service.EnterArea(incomingCarData.PoliceNumber)
		if errParking != nil {
			err = fmt.Errorf("failed to enter area %s", errParking.Error())
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
			err = fmt.Errorf("failed to exit area with police id %s", incomingCarData.PoliceNumber)
			return
		}

		metaByte, errM := json.Marshal(metadata)
		if errM != nil {
			err = fmt.Errorf("failed to marshal metadata")
			return
		}

		response = string(metaByte)
		return
	case types.CmdStatus:
		dataBytes, errMarshall := json.Marshal(srv.service)
		if errMarshall != nil {
			err = fmt.Errorf("failed to marshall parking data")
			return
		}

		response = string(dataBytes)
		return
	}

	return "", nil
}
