package boot

import (
	"errors"
	"fmt"
	"github.com/khafidprayoga/parking-app/internal/backend"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/khafidprayoga/parking-app/internal/server"
)

func StartApp() {
	// init socket
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error listening on port :8080 with reason %v", err)
	}

	log.Printf("Parking App Server %s is listening on port :8080\n", AppConfig.AppVersion)

	uc := backend.NewParkingService()
	service := server.CreateAppServer(uc)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// running server main thread for backend service in background
	go func() {
		for {
			// handling incoming request
			conn, errAcc := listener.Accept()
			if errAcc != nil {
				errMsg := errors.Unwrap(errAcc)
				if strings.Contains(errMsg.Error(), "close") {
					// exit on closed server (reject request)
					log.Println("server is going to shutdown, exit request listener")
					return
				}

				log.Printf("error accepting connection: %v", errAcc)
				continue
			}

			// set req-res timeout
			conn.SetDeadline(time.Now().Add(AppConfig.ConnLifetime))

			// emit data to service
			go emit(conn, service)
		}
	}()

	// watch shutdown signal
	<-quit
	errCloseTcp := listener.Close()
	if errCloseTcp != nil {
		log.Println(errCloseTcp)
	}

	log.Println("shutting down app in a few seconds")
	ticker := time.NewTicker(1 * time.Second)

	for x := 0; x < 3; x++ {
		<-ticker.C
		fmt.Printf(".")
	}
	ticker.Stop()
	os.Exit(0)
}
