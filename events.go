package gosocketioclient

import (
    "log"
)

func (c *Client) handleIncomingMessage(message map[string]interface{}) {
    event, ok := message["event"].(string)
    if !ok {
        log.Println("Invalid event received")
        return
    }

    if handler, found := c.handlers[event]; found {
        handler(message["data"])
    }
}
