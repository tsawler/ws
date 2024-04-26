package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

const (
	TextMessage = 1
	JSONMessage = 2
)

// Sockets is the main type for this library.
type Sockets struct {
	ClientChan chan Payload
	Clients    map[WebSocketConnection]string
	ErrorChan  chan error
}

// New is a factory function to return a new *Sockets object.
func New() *Sockets {
	return &Sockets{
		ClientChan: make(chan Payload),                   // The channel we send ws payloads (from client) to.
		Clients:    make(map[WebSocketConnection]string), // A map of all connected clients.
		ErrorChan:  make(chan error),                     // A channel where errors (or nil) is sent.
	}
}

// upgradeConnection is the upgraded connection needed for ws.
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var clientsMutex sync.Mutex

// Payload defines the data we receive from the client.
type Payload struct {
	MessageType int                 `json:"message_type"`
	Message     string              `json:"message"`
	Conn        WebSocketConnection `json:"-"`
}

// JSONResponse defines the JSON we send back to client
type JSONResponse struct {
	Message     string              `json:"message"`
	CurrentConn WebSocketConnection `json:"-"`
}

// WebSocketConnection is a simple wrapper which holds a websocket connection.
type WebSocketConnection struct {
	*websocket.Conn
}

// SocketEndPoint handles websocket connections.
func (s *Sockets) SocketEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		s.ErrorChan <- err
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
