package gosocketioclient

import (
    "log"
)

type Client struct {
    conn *Connection
    // Add a map to store event handlers
    handlers map[string]func(interface{})
}

func NewClient(url string) (*Client, error) {
    connection, err := NewConnection(url)
    if err != nil {
        return nil, err
    }
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
