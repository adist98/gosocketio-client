package engineio

import (
    "errors"
    "net/url"

    "github.com/gorilla/websocket"
)

// Transport interface defines the methods that a transport should implement.
type Transport interface {
    Connect(serverURL string, queryParams map[string]string) (*websocket.Conn, error)
}

// WebSocketTransport implements the Transport interface for WebSocket.
type WebSocketTransport struct{}

// Connect establishes a WebSocket connection.
func (w *WebSocketTransport) Connect(serverURL string, queryParams map[string]string) (*websocket.Conn, error) {
    u, err := url.Parse(serverURL)
    if err != nil {
        return nil, errors.New("invalid server URL")
    }

    u.Scheme = "ws"
    q := u.Query()
    for key, value := range queryParams {
        q.Set(key, value)
    }
    u.RawQuery = q.Encode()

    conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        return nil, err
    }

    return conn, nil
}

// NewWebSocketTransport creates a new WebSocketTransport instance.
func NewWebSocketTransport() *WebSocketTransport {
    return &WebSocketTransport{}
}
