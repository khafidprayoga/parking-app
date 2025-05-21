package types

type Socket struct {
	Command string `json:"command"`
	Data    any    `json:"data"`
}

const (
	SocketCallSuccess = "SUCCESS"
	SocketCallError   = "ERROR"
)

type SocketServerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
