package sockets

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Sockets struct{}

// New is a factory function to return a new *Sockets object.
func New() *Sockets {
	return &Sockets{}
}

var wsChan = make(chan WsPayload)                  // The channel we send ws payloads (from client) to.
var Clients = make(map[WebSocketConnection]string) // A map of all connected Clients.

// UpgradeConnection is the upgraded connection needed for ws.
var UpgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WsPayload defines the data we receive from the client
type WsPayload struct {
	Message string              `json:"message"`
	Conn    WebSocketConnection `json:"-"`
}

// WsJsonResponse defines the json we send back to client
type WsJsonResponse struct {
	Message     string              `json:"message"`
	CurrentConn WebSocketConnection `json:"-"`
}

// WebSocketConnection holds the websocket connection
type WebSocketConnection struct {
	*websocket.Conn
}

// SocketEndPoint handles websocket connections
func (s *Sockets) SocketEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println(fmt.Sprintf("Client Connected from %s", r.RemoteAddr))

	// Create a WebSocketConnection object with the client's connection.
	conn := WebSocketConnection{Conn: ws}
	// Add the client to the map of connected Clients.
	Clients[conn] = ""

	// Start listening for this client.
	go ListenForWS(&conn)
}

// ListenForWS is the goroutine that listens for communication from Clients.
func ListenForWS(conn *WebSocketConnection) {
	// If this dies, just restart it.
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	// payload is the variable we read a payload into.
	var payload WsPayload

	// This loop will run forever, waiting for something to come
	// in on a websocket.
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			// Send the incoming payload to wsChan.
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

// ListenToWsChannel listens to the wsChan and pushes data to broadcast function.
func (s *Sockets) ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChan
		response.Message = e.Message
		BroadcastToAll(response)
	}
}

// BroadcastToAll broadcasts data to all connected Clients.
func BroadcastToAll(response WsJsonResponse) {
	for client := range Clients {
		// broadcast to every connected client
		err := client.WriteJSON(response)
		if err != nil {
			_ = client.Close()
			delete(Clients, client)
		}
	}
}
