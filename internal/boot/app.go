package boot

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/khafidprayoga/parking-app/internal/server"
	"github.com/khafidprayoga/parking-app/internal/types"
)

func StartApp() {
	// init socket
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error listening on port :8080 with reason %v", err)
	}

	defer listener.Close()
	log.Println("listening on port :8080")

	uc := server.NewParkingService()
	service := server.CreateAppServer(*uc)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// running server main thread for backend service in background
	go func() {
		for {
			// handling incoming request
			conn, errAcc := listener.Accept()
			if errAcc != nil {
				log.Printf("error accepting connection: %v", errAcc)
				continue
			}

			// set req-res timeout
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))

			// emit data to service
			go func(conn net.Conn) {
				defer conn.Close()

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
				resMsg, errProcess := service.HandleIncomingMsg(data)

				response := types.SocketServerResponse{
					Status:  types.SocketCallSuccess,
					Message: resMsg,
				}

				if errProcess != nil {
					response.Status = types.SocketCallError
					response.Message = errProcess.Error()
				}

				resB, _ := json.Marshal(response)
				conn.Write(resB)
			}(conn)
		}
	}()

	// watch shutdown signal
	<-quit
	log.Println("shutting down app")
}
