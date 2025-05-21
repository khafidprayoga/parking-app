package main

import (
	"encoding/json"
	"fmt"
	"github.com/khafidprayoga/parking-app/internal/types"
	"log"
	"net"
	"os"
	"strconv"
)

// main to send input from file to running backend services
func main() {

	defaultMsg := fmt.Sprintf(
		"Parking App Service CLI:\n"+
			"\nExample: `parking-app create_parking_lot 12`\n\n"+
			"available commands:\n"+
			"\t%s {lotCapacity:int} => for initialize parking lot size\n"+
			"\t%s {carNumber:string} => parking a car\n"+
			"\t%s {carNumber:string} {hours:int}  => for a car to exit parking area\n"+
			"\t%s => view status of the parking area app service\n"+
			"\thelp  => show this message",
		types.CmdCreateStore,
		types.CmdPark,
		types.CmdLeave,
		types.CmdStatus)
	if len(os.Args) == 1 {
		fmt.Println(defaultMsg)
		return
	}

	command := os.Args[1]
	args := os.Args[2:]
	if command != types.CmdStatus && len(os.Args) < 3 {
		log.Fatalln(defaultMsg)
	}

	switch command {
	case types.CmdCreateStore:
		parkingLotCap := args[0]
		if errSendReq := sendRequest(command, parkingLotCap); errSendReq != nil {
			log.Fatal(errSendReq)
		}

		log.Printf("CLIENT:Creating parking with capacity of %v lot", parkingLotCap)
	case types.CmdPark:
		if errSendReq := sendRequest(command, types.CarDTO{
			PoliceNumber: args[0],
		}); errSendReq != nil {
			log.Fatal(errSendReq)
		}
	case types.CmdLeave:
		durationInHours, errParseDur := strconv.Atoi(args[1])
		if errParseDur != nil {
			log.Fatal(errParseDur)
		}

		if errSendReq := sendRequest(command, types.CarDTO{
			PoliceNumber: args[0],
			Hours:        durationInHours,
		}); errSendReq != nil {
			log.Fatal(errSendReq)
		}
	case types.CmdStatus:
		if errSendReq := sendRequest(command, nil); errSendReq != nil {
			log.Fatal(errSendReq)
		}
	default:
		log.Fatalln(defaultMsg)
	}
}

func sendRequest(command string, data any) error {
	conn, errDial := net.Dial("tcp", "localhost:8080")
	if errDial != nil {
		return fmt.Errorf("cannot connect to server: %v", errDial)
	}

	reqBytes, errMarshall := json.Marshal(types.Socket{
		Command: command,
		Data:    data,
	})
	if errMarshall != nil {
		return fmt.Errorf("cannot marshal json: %v", errMarshall)
	}

	_, errSend := conn.Write(reqBytes)
	if errSend != nil {
		return fmt.Errorf("cannot send request: %v", errSend)
	}

	return nil
}
