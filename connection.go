package gosocketioclient

import (
	"github.com/gorilla/websocket"
)

type Connection struct {
    Conn *websocket.Conn
}

func NewConnection(url string) (*Connection, error) {
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return nil, err
    }
    return &Connection{Conn: conn}, nil
}

func (c *Connection) Close() error {
    return c.Conn.Close()
}