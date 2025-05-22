package boot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/khafidprayoga/parking-app/internal/server"
	"github.com/khafidprayoga/parking-app/internal/types"
)

func StartApp() {
	// init socket
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error listening on port :8080 with reason %v", err)
	}

	log.Println("listening on port :8080")

	uc := server.NewParkingService()
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
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))

			// emit data to service
			go func(conn net.Conn) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("panic recovered in conn handler: %v", r)
					}
					conn.Close()
				}()

				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					log.Printf("error reading from connection: %v", err)
					return
				}

				data := types.Socket{}
				if err := json.Unmarshal(buf[:n], &data); err != nil {
					log.Printf("error unmarshalling from connection: %v", err)
					return
				}

				id, e := uuid.Parse(data.XRequestId)

				if e != nil {
					// override invalid uuid or empty uuid
					id = uuid.New()
				}

				log.Printf(
					"Handling Request with Id: %v Command: %v At: %v \n",
					id.String(),
					data.Command,
					time.Now().Format(time.RFC3339),
				)
				resMsg, errProcess := service.HandleIncomingMsg(data)

				response := types.SocketServerResponse{
					Status:  types.SocketCallSuccess,
					Message: resMsg,
				}

				if errProcess != nil {
					response.Status = types.SocketCallError
					response.Message = errProcess.Error()
				}

				resB, errM := json.Marshal(response)
				if errM != nil {
					log.Printf("error reading from connection: %v", err)
					return
				}
				conn.Write(resB)
			}(conn)
		}
	}()

	// watch shutdown signal
	<-quit
	errCloseTcp := listener.Close()
	if errCloseTcp != nil {
		log.Println(errCloseTcp)
	}

	log.Println("shutting down app in 5s")
	ticker := time.NewTicker(1 * time.Second)

	for x := 0; x < 5; x++ {
		<-ticker.C
		fmt.Printf(".")
	}
	ticker.Stop()
	os.Exit(0)
}
