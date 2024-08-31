package engineio

import (
    "encoding/json"
    "fmt"
)

// Packet types for Engine.IO
const (
    Open    = '0'
    Close   = '1'
    Ping    = '2'
    Pong    = '3'
    Message = '4'
    Upgrade = '5'
    Noop    = '6'
)

// HandshakeResponse represents the server's response to the initial handshake.
type HandshakeResponse struct {
    SID          string   `json:"sid"`
    Upgrades     []string `json:"upgrades"`
    PingInterval int      `json:"pingInterval"`
    PingTimeout  int      `json:"pingTimeout"`
}

// ParsePacket parses a raw Engine.IO packet into its type and payload.
func ParsePacket(packet []byte) (packetType byte, payload []byte, err error) {
    if len(packet) < 1 {
        return 0, nil, fmt.Errorf("invalid packet: too short")
    }

    packetType = packet[0]
    payload = packet[1:]
    return packetType, payload, nil
}

// BuildPacket constructs an Engine.IO packet from a type and payload.
func BuildPacket(packetType byte, payload []byte) []byte {
    return append([]byte{packetType}, payload...)
}

// ParseHandshakeResponse parses the handshake response from the server.
func ParseHandshakeResponse(response []byte) (*HandshakeResponse, error) {
    var handshake HandshakeResponse
    if err := json.Unmarshal(response, &handshake); err != nil {
        return nil, fmt.Errorf("failed to parse handshake response: %v", err)
    }
    return &handshake, nil
}
