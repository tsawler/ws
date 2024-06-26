// Package ws is a simple library that makes it easy to use WebSockets
// (specifically Gorilla Websockets) in your Go application.

package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

const (
	TextMessage = 1 // text format.
	JSONMessage = 2 // JSON format.
)

// Sockets is the main type for this library.
type Sockets struct {
	ClientChan      chan Payload                // The channel which receives message payloads.
	Clients         map[WebSocketConnection]any // A map of all connected clients.
	ErrorChan       chan error                  // A channel which receives errors.
	ReadBufferSize  int                         // I/O read buffer size in bytes.
	WriteBufferSize int                         // I/O write buffer size in bytes.
}

var (
	ReadBufferSize  = 1024 // Set a sensible default of 1024 bytes for read buffer.
	WriteBufferSize = 1024 // Set a sensible default of 1024 bytes for write buffer.
)

// New is a factory function to return a new *Sockets object.
func New() *Sockets {
	return &Sockets{
		ClientChan:      make(chan Payload),
		Clients:         make(map[WebSocketConnection]any),
		ErrorChan:       make(chan error),
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
	}
}

// upgradeConnection is the upgraded connection needed for websockets.
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  ReadBufferSize,
	WriteBufferSize: WriteBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// clientsMutex is used to lock/unlock the map of connected clients in order to avoid
// race conditions.
var clientsMutex sync.Mutex

// Payload defines the data we receive from the client.
type Payload struct {
	MessageType int                 `json:"message_type"`
	Message     string              `json:"message"`
	Data        any                 `json:"data,omitempty"`
	Conn        WebSocketConnection `json:"-"`
}

// JSONResponse defines the JSON we send back to client
type JSONResponse struct {
	Message     string              `json:"message"`
	Data        any                 `json:"data,omitempty"`
	CurrentConn WebSocketConnection `json:"-"`
}

// WebSocketConnection is a simple wrapper which holds a websocket connection.
type WebSocketConnection struct {
	*websocket.Conn
}

// SocketEndPoint handles websocket connections.
func (s *Sockets) SocketEndPoint(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to the WebSocket protocol.
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		s.ErrorChan <- err
		return
	}

	log.Printf("Client Connected from %s", r.RemoteAddr)

	// Create a WebSocketConnection object with the client's connection.
	conn := WebSocketConnection{Conn: ws}

	// Add the client to the map of connected Clients.
	clientsMutex.Lock()
	s.Clients[conn] = ""
	clientsMutex.Unlock()

	// Start listening for this client.
	go s.listenForWS(&conn)
}

// listenForWS is the goroutine that listens for communication from Clients.
func (s *Sockets) listenForWS(conn *WebSocketConnection) {
	// If this dies, just restart it.
	defer func() {
		recover()
	}()

	// payload is the variable we read a payload into.
	var payload Payload

	// This loop will run forever, waiting for something to come
	// in on a websocket.
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			s.ErrorChan <- err
		} else {
			// Send the incoming payload to SocketsChan.
			payload.Conn = *conn
			s.ClientChan <- payload
		}
	}
}

// ListenToWsChannel listens to the SocketsChan and pushes data to broadcast functions.
func (s *Sockets) ListenToWsChannel() {

	for {
		e := <-s.ClientChan

		switch e.MessageType {
		case TextMessage:
			s.BroadcastTextToAll(e.Message)
		case JSONMessage:
			var response JSONResponse
			response.Message = e.Message
			s.BroadcastJSONToAll(response)
		default:
			s.ErrorChan <- fmt.Errorf("invalid message type %d received", e.MessageType)
		}
	}
}

// BroadcastJSONToAll broadcasts JSON data to all connected Clients.
func (s *Sockets) BroadcastJSONToAll(payload any) {
	clientsMutex.Lock()
	for client := range s.Clients {
		// Broadcast to every connected client.
		err := client.WriteJSON(payload)
		if err != nil {
			// Someone left. Remove them from the map of connected users.
			_ = client.Close()
			delete(s.Clients, client)
		}
	}
	clientsMutex.Unlock()
}

// BroadcastTextToAll broadcasts textual data to all connected Clients.
func (s *Sockets) BroadcastTextToAll(payload string) {
	clientsMutex.Lock()
	for client := range s.Clients {
		// Broadcast to every connected client.
		err := client.WriteMessage(websocket.TextMessage, []byte(payload))
		if err != nil {
			// Someone left. Remove them from the map of connected users.
			_ = client.Close()
			delete(s.Clients, client)
		}
	}
	clientsMutex.Unlock()
}
