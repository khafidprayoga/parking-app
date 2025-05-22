package types

type Socket struct {
	Command    string `json:"command"`
	Data       any    `json:"data"`
	XRequestId string `json:"x_request_id"`
}

const (
	SocketCallSuccess = "OK"
	SocketCallError   = "ERROR"
)

type SocketServerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
