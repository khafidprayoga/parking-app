package types

type Socket struct {
	Command string `json:"command"`
	Data    any    `json:"data"`
}
