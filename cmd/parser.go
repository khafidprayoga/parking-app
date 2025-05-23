package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/khafidprayoga/parking-app/internal/extra"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	bootstrap "github.com/khafidprayoga/parking-app/internal/boot"

	"github.com/khafidprayoga/parking-app/internal/types"
)

// main to send input from file to running backend services
func main() {

	defaultMsg := fmt.Sprintf(
		"Parking App Service CLI:\n"+
			"\nExample: `EXAMPLE`\n\n"+
			"available commands:\n"+
			"\t%s => start parking app server socket at :8080\n"+
			"\t%s {lotCapacity:int} => for initialize parking lot size\n"+
			"\t%s {carNumber:string} => parking a car\n"+
			"\t%s {carNumber:string} {hours:int}  => for a car to exit parking area\n"+
			"\t%s => view status of the parking area app service\n"+
			"\thelp  => show this message",
		types.CmdServe,
		types.CmdCreateStore,
		types.CmdPark,
		types.CmdLeave,
		types.CmdStatus)

	if len(os.Args) < 2 {
		defaultMsg = strings.Replace(defaultMsg, "EXAMPLE", fmt.Sprintf("parking-app %s 12", types.CmdCreateStore), -1)
		fmt.Println(defaultMsg)
		return
	}

	command := os.Args[1]
	param := os.Args[2:]

	// on check server state
	if command != types.CmdStatus && command != types.CmdServe && len(os.Args) < 3 {
		defaultMsg = strings.Replace(defaultMsg, "EXAMPLE", fmt.Sprintf("parking-app %s 12", types.CmdCreateStore), -1)
		log.Fatalln(defaultMsg)
	}

	switch command {
	case types.CmdServe:
		log.Println("Starting Parking App Server")
		bootstrap.StartApp()
	case types.CmdCreateStore:
		if len(param) == 0 {
			log.Printf("lot capacity not specified")
			defaultMsg = strings.Replace(defaultMsg, "EXAMPLE", fmt.Sprintf("parking-app %s 12", types.CmdCreateStore), -1)
			log.Println(defaultMsg)
			return
		}

		parkingLotCap := param[0]
		log.Printf("CLIENT:Creating parking with capacity of %v lot", parkingLotCap)

		if errSendReq := sendRequest(command, parkingLotCap); errSendReq != nil {
			log.Fatal(errSendReq)
		}
	case types.CmdPark:
		if len(param) == 0 {
			log.Printf("police number on car not specified")
			defaultMsg = strings.Replace(defaultMsg, "EXAMPLE", fmt.Sprintf("parking-app %s KA-01-HH-270", types.CmdPark), -1)
			log.Println(defaultMsg)
			return
		}

		if errSendReq := sendRequest(command, types.CarDTO{
			PoliceNumber: param[0],
		}); errSendReq != nil {
			log.Fatal(errSendReq)
		}
	case types.CmdLeave:
		if len(param) < 2 {
			log.Printf("hours cost not specified")
			defaultMsg = strings.Replace(defaultMsg, "EXAMPLE", fmt.Sprintf("parking-app %s KA-01-HH-270 2", types.CmdLeave), -1)
			log.Println(defaultMsg)
			return
		}

		durationInHours, errParseDur := strconv.Atoi(param[1])
		if errParseDur != nil {
			log.Fatal(errParseDur)
		}

		if errSendReq := sendRequest(command, types.CarDTO{
			PoliceNumber: param[0],
			Hours:        durationInHours,
		}); errSendReq != nil {
			log.Fatal(errSendReq)
		}
	case types.CmdStatus:
		if errSendReq := sendRequest(command, nil); errSendReq != nil {
			log.Fatal(errSendReq)
		}
	case types.CmdImport:
		if len(param) < 1 {
			log.Printf("command instruction file is not specified")
			defaultMsg = strings.Replace(defaultMsg, "EXAMPLE", fmt.Sprintf("parking-app %s KA-01-HH-270 2", types.CmdLeave), -1)
			log.Println(defaultMsg)
			return
		}

		cmdList, errParseCmd := extra.ParseImportCmd(param[0])
		if errParseCmd != nil {
			log.Fatal(errParseCmd)
		}

		// sequential process
		for _, cmd := range cmdList {
			if errSend := sendRequest(cmd.Command, cmd.Data); errSend != nil {
				log.Fatal(errSend)
			}
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
		Command:    command,
		Data:       data,
		XRequestId: uuid.NewString(),
	})
	if errMarshall != nil {
		return fmt.Errorf("cannot marshal json: %v", errMarshall)
	}

	_, errSend := conn.Write(reqBytes)
	if errSend != nil {
		return fmt.Errorf("cannot send request: %v", errSend)
	}

	// reading response
	buf := make([]byte, 2048)
	size, errRead := conn.Read(buf)
	if errRead != nil {
		return fmt.Errorf("cannot read response: %v", errRead)
	}

	res := types.SocketServerResponse{}
	if err := json.Unmarshal(buf[:size], &res); err != nil {
		return fmt.Errorf("cannot unmarshal response: %v", errSend)
	}
	log.Printf("\nSERVER-STATUS: %s\n"+
		"SERVER-RESPONSE: %s",
		res.Status, res.Message)
	return nil
}
