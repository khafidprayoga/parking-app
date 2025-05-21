package server

import (
	"fmt"
	"github.com/khafidprayoga/parking-app/internal/types"
	"log"
	"strconv"
	"strings"
	"time"
)

func (p *ParkingAppServer) HandleIncomingMsg(msg types.Socket) (err error) {
	switch msg.Command {
	case types.CmdCreateStore:
		parkingCap, errCv := strconv.Atoi(msg.Data.(string))
		if errCv != nil {
			err = fmt.Errorf("failed to convert string to int at %s actions", msg.Command)
			return
		}

		if p.service.LotCapacity > 0 {
			err = fmt.Errorf("failed, already initalize the parking lot capacity")
			return
		}

		p.service.LotCapacity = parkingCap
		p.service.Store = make([]*types.Car, parkingCap)
	case types.CmdPark:
		incomingCarData := types.Car{
			PoliceNumber: msg.Data.(map[string]any)["police_number"].(string),
		}

		// validate if  car number not already exist on the parking area
		for _, car := range p.service.Store {
			if car != nil {
				if strings.EqualFold(car.PoliceNumber, incomingCarData.PoliceNumber) {
					err = fmt.Errorf("failed, already parked the parking lot capacity")
					return
				}
			}
		}

		for index, car := range p.service.Store {
			// allocating nearest parking lot from the door gateway
			if car == nil {
				p.service.Store[index] = &types.Car{
					ID:           index + 1,
					Color:        incomingCarData.Color,
					PoliceNumber: incomingCarData.PoliceNumber,
					ParkingAt:    time.Now(),
					ExitAt:       nil,
				}

				incomingCarData.ID = index
				break
			}

		}

		log.Printf("successfully parked car. with police number %s and SLOT number id %v", incomingCarData.PoliceNumber, incomingCarData.ID)
	case types.CmdLeave:
	case types.CmdStatus:
		fmt.Println(1)
	}

	return nil
}
