package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/khafidprayoga/parking-app/internal/server"
	"github.com/khafidprayoga/parking-app/internal/types"
)

func main() {
	// init socket
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error listening on port :8080 with reason %v", err)
	}

	defer listener.Close()
	log.Println("listening on port :8080")

	uc := server.NewParkingService()
	service := server.CreateAppServer(uc)

	// running backend in background
	for {
		conn, errAcc := listener.Accept()
		if errAcc != nil {
			log.Printf("error accepting connection: %v", errAcc)
			continue
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("error reading from connection: %v", err)
		}

		data := types.Socket{}
		if err := json.Unmarshal(buf[:n], &data); err != nil {
			log.Printf("error unmarshalling from connection: %v", err)
		}

		// emit data to handler
		go func() {
			if errProcess := service.HandleIncomingMsg(data); errProcess != nil {
				log.Printf("error processing message: %v", errProcess)
			}
		}()
	}

}
