package engineio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"log"

	"github.com/gorilla/websocket"
)

type EngineIOClient struct {
    Conn         *websocket.Conn
    SID          string
    PingInterval int
    PingTimeout  int
}

// NewEngineIOClient initializes a new Engine.IO client
func NewEngineIOClient(serverURL string) (*EngineIOClient, error) {
    // Step 1: Initialize the connection using polling
    handshake, err := InitializePollingConnection(serverURL)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize polling connection: %v", err)
    }

    // Step 2: Upgrade the connection to WebSocket
    client, err := UpgradeToWebSocket(serverURL, handshake)
    if err != nil {
        return nil, fmt.Errorf("failed to upgrade to WebSocket: %v", err)
    }

    log.Println("Successfully upgraded to WebSocket")
    return client, nil
}

func InitializePollingConnection(serverURL string) (*HandshakeResponse, error) {
    u, err := url.Parse(serverURL)
    if err != nil {
        return nil, fmt.Errorf("invalid server URL: %v", err)
    }

    q := u.Query()
    q.Set("EIO", "4")
    q.Set("transport", "polling")
    u.RawQuery = q.Encode()

    resp, err := http.Get(u.String())
    if err != nil {
        return nil, fmt.Errorf("failed to send handshake request: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read handshake response: %v", err)
    }

    log.Printf("Raw handshake response: %s", string(body))  // Log the raw response

    var handshakeData string
    if body[0] == '0' {
        colonIndex := bytes.IndexByte(body, ':')
        if colonIndex != -1 {
            handshakeData = string(body[colonIndex+1:])
        } else {
            handshakeData = string(body[1:])
        }
    } else {
        handshakeData = string(body)
    }

    var handshake HandshakeResponse
    err = json.Unmarshal([]byte(handshakeData), &handshake)
    if err != nil {
        return nil, fmt.Errorf("failed to parse handshake response: %v", err)
    }

    return &handshake, nil
}


func UpgradeToWebSocket(serverURL string, handshake *HandshakeResponse) (*EngineIOClient, error) {
    u, err := url.Parse(serverURL)
    if err != nil {
        return nil, fmt.Errorf("invalid server URL: %v", err)
    }

    u.Scheme = "ws"
    q := u.Query()
    q.Set("EIO", "4")
    q.Set("transport", "websocket")
    q.Set("sid", handshake.SID)
    u.RawQuery = q.Encode()

    conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        return nil, fmt.Errorf("failed to upgrade to websocket: %v", err)
    }

    client := &EngineIOClient{
        Conn:         conn,
        SID:          handshake.SID,
        PingInterval: handshake.PingInterval,
        PingTimeout:  handshake.PingTimeout,
    }

    if err := client.performWebSocketUpgrade(); err != nil {
        return nil, fmt.Errorf("websocket upgrade handshake failed: %v", err)
    }

    return client, nil
}

func (c *EngineIOClient) performWebSocketUpgrade() error {
    _, message, err := c.Conn.ReadMessage()
    if err != nil {
        return fmt.Errorf("failed to read open packet after websocket upgrade: %v", err)
    }

    if len(message) == 0 || message[0] != '0' {
        return fmt.Errorf("invalid open packet after websocket upgrade: %s", string(message))
    }

    return nil
}

func (c *EngineIOClient) Listen() {
    go c.sendPingMessages()

    for {
        messageType, message, err := c.Conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }

        if messageType != websocket.TextMessage {
            log.Printf("Received non-text message: %v", messageType)
            continue
        }

        if len(message) == 0 {
            log.Println("Received empty message")
            continue
        }

        packetType := message[0]
        payload := message[1:]

        switch packetType {
        case '0': 
            log.Println("Received open packet")
        case '2': 
            log.Println("Received ping packet, sending pong")
            err = c.Conn.WriteMessage(websocket.TextMessage, []byte("3"))
            if err != nil {
                log.Printf("Error sending pong: %v", err)
                break
            }
        case '3': 
            log.Println("Received pong packet")
        case '4': 
            log.Printf("Received message: %s", string(payload))
        default:
            log.Printf("Received unknown packet type: %c", packetType)
        }
    }

    c.Conn.Close()
}

func (c *EngineIOClient) sendPingMessages() {
    ticker := time.NewTicker(time.Duration(c.PingInterval) * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {
        err := c.Conn.WriteMessage(websocket.TextMessage, []byte("2"))
        if err != nil {
            log.Printf("Error sending ping: %v", err)
            return
        }
        log.Println("Sent ping to server")
    }
}
