package boot

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/khafidprayoga/parking-app/internal/server"
	"github.com/khafidprayoga/parking-app/internal/types"
	"log"
	"net"
	"time"
)

func emit(conn net.Conn, service *server.ParkingAppServer) {
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
}
