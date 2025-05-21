package main

import (
	"log"
	"net"
)

func main() {
	// init socket
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error listening on port :8080 with reason %v", err)
	}
	log.Println("listening on port :8080")

	// running backend in background
	for {
		conn, errAcc := ln.Accept()
		if errAcc != nil {
			log.Printf("error accepting connection: %v", errAcc)
			continue
		}

	}

}
