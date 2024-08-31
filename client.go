package gosocketioclient

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
    conn *Connection
    // Add a map to store event handlers
    handlers map[string]func(interface{})
}

// NewClient creates a new Socket.IO client with the provided server URL.
func NewClient(url string) (*Client, error) {
    // Ensure the URL includes the required Socket.IO query parameters
    if !strings.Contains(url, "EIO=") {
        if strings.Contains(url, "?") {
            url += "&EIO=4&transport=websocket"
        } else {
            url += "?EIO=4&transport=websocket"
        }
    }

    // Establish the WebSocket connection
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Socket.IO server: %w", err)
    }

    // Wrap the WebSocket connection in your Connection struct
    connection := &Connection{
        Conn: conn,
    }

    // Return a new Client instance with the Connection
    return &Client{
        conn: connection,
    }, nil
}

func (c *Client) On(event string, handler func(interface{})) {
    c.handlers[event] = handler
}

func (c *Client) Emit(event string, data interface{}) error {
    // Logic to encode the event and send it via WebSocket
    return c.conn.Conn.WriteJSON(map[string]interface{}{
        "event": event,
        "data":  data,
    })
}

func (c *Client) Listen() {
    for {
        var message map[string]interface{}
        err := c.conn.Conn.ReadJSON(&message)
        if err != nil {
            log.Println("Error reading message:", err)
            return
        }

        event := message["event"].(string)
        data := message["data"]

        if handler, ok := c.handlers[event]; ok {
            handler(data)
        }
    }
}
