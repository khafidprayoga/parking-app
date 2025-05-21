package main

import (
	"github.com/khafidprayoga/parking-app/internal/types"
	"log"
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
		log.Printf("Creating parking with capacity of %v lot", args[0])
	case types.CmdPark:
	case types.CmdLeave:
	case types.CmdStatus:
	default:
		panic("incorrect command use: options {args:int}")
	}
}
