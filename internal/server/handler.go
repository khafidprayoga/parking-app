package server

import (
	"github.com/khafidprayoga/parking-app/internal/types"
	"log"
	"strconv"
)

func (p *ParkingAppServer) HandleIncomingMsg(msg types.Socket) {
	switch msg.Command {
	case types.CmdCreateStore:
		parkingCap, errCv := strconv.Atoi(msg.Data.(string))
		if errCv != nil {
			log.Printf("failed to convert string to int at %s actions", msg.Command)
			return
		}

		if p.service.LotCapacity > 0 {
			log.Printf("failed, already initalize the parking lot capacity")
			return
		}

		p.service.LotCapacity = parkingCap
		p.service.Store = make([]*types.Car, parkingCap)
	case types.CmdPark:
	case types.CmdLeave:
	case types.CmdStatus:
	}
}
