package gosocketioclient

import (
	"encoding/json"
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
		conn:     connection,
		handlers: make(map[string]func(interface{})),
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
        _, message, err := c.conn.Conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            return
        }

        // Try to unmarshal the message into a generic map
        var parsedMessage map[string]interface{}
        err = json.Unmarshal(message, &parsedMessage)
        if err != nil {
            log.Println("Error unmarshaling message:", err)
            continue
        }

        // Check if the message has an event and data
        event, ok := parsedMessage["event"].(string)
        if !ok {
            log.Println("Received a message without an 'event' field:", parsedMessage)
            continue
        }

        // Call the registered handler if it exists
        if handler, found := c.handlers[event]; found {
            handler(parsedMessage["data"])
        } else {
            log.Println("No handler found for event:", event)
        }
    }
}
