package dto

// Packet - websocket communication
type Packet struct {
	ID      uint64      `json:"-"`
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
}
