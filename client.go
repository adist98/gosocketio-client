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

        // Log the raw message for debugging
        log.Println("Raw message received:", string(message))

        // The first character indicates the packet type
        packetType := message[0] - '0'
        payload := message[1:]

        switch packetType {
        case 0: // Open packet
            log.Println("Received 'open' packet")
            var openPayload map[string]interface{}
            err := json.Unmarshal(payload, &openPayload)
            if err != nil {
                log.Println("Error unmarshaling open packet:", err)
                continue
            }
            log.Println("Open packet data:", openPayload)

        case 40: // Connected to namespace
            log.Println("Connected to namespace")
            // Handle namespace connection, no additional payload expected

        case 42: // Event message
            log.Println("Received 'event' packet")
            var eventPayload []interface{}
            err := json.Unmarshal(payload, &eventPayload)
            if err != nil {
                log.Println("Error unmarshaling event packet:", err)
                continue
            }

            // The event name is the first item in the array
            eventName, ok := eventPayload[0].(string)
            if !ok {
                log.Println("Invalid event name in packet:", eventPayload)
                continue
            }

            // The event data is the second item
            if len(eventPayload) > 1 {
                eventData := eventPayload[1]
                if handler, found := c.handlers[eventName]; found {
                    handler(eventData)
                } else {
                    log.Println("No handler found for event:", eventName)
                }
            }

        default:
            log.Println("Unknown packet type:", packetType)
        }
    }
}
