package main

import (
	"encoding/json"
	"fmt"
	"github.com/khafidprayoga/parking-app/internal/types"
	"log"
	"net"
	"os"
)

// main to send input from file to running backend services
func main() {
	command := os.Args[1]
	args := os.Args[2:]

	if len(os.Args) < 3 {
		log.Fatalf("incorrect command use:\n\t\t options {args:int}")
	}

	log.Printf("got command: %s with argument %s\n", command, args)

	switch command {
	case types.CmdCreateStore:
		parkingLotCap := args[0]
		if errSendReq := sendRequest(command, parkingLotCap); errSendReq != nil {
			log.Fatal(errSendReq)
		}

		log.Printf("CLIENT:Creating parking with capacity of %v lot", parkingLotCap)
	case types.CmdPark:
	case types.CmdLeave:
	case types.CmdStatus:
	default:
		panic("incorrect command use: options {args:int}")
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
