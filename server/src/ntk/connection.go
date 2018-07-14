package ntk

import (
	"controlr/plugin/server/src/dto"

	"golang.org/x/net/websocket"
)

// MessageFunc - websocket handler
type MessageFunc func(message *dto.Packet)

// CloseFunc - connection closer
type CloseFunc func(conn *Connection, err error)

// Connection type
type Connection struct {
	ID        uint64
	ws        *websocket.Conn
	onMessage MessageFunc
	onClose   CloseFunc
}

// NewConnection constructor
func NewConnection(id uint64, ws *websocket.Conn, onMessage MessageFunc, onClose CloseFunc) *Connection {
	return &Connection{
		ID:        id,
		ws:        ws,
		onMessage: onMessage,
		onClose:   onClose,
	}
}

func (c *Connection) Read() (err error) {
	for {
		var packet dto.Packet
		err = websocket.JSON.Receive(c.ws, &packet)
		if err != nil {
			go c.onClose(c, err)
			return
		}

		packet.ID = c.ID
		go c.onMessage(&packet)
	}
}

func (c *Connection) Write(packet *dto.Packet) (err error) {
	err = websocket.JSON.Send(c.ws, packet)
	return
}
