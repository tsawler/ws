<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/tsawler/persist/blob/main/LICENSE.md)
<a href="https://pkg.go.dev/github.com/tsawler/ws"><img src="https://img.shields.io/badge/godoc-reference-%23007d9c.svg"></a>
![Tests](https://github.com/tsawler/ws/actions/workflows/tests.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tsawler/ws)](https://goreportcard.com/report/github.com/tsawler/ws)

# ws

Package ws is a simple library that makes it easy to use WebSockets (specifically 
[Gorilla Websockets](https://github.com/gorilla/websocket)) in your Go application.

## Installation
Install it the usual way:

~~~
go get github.com/tsawler/ws
~~~

## Usage
Create a variable of type `ws.Sockets` by calling the `ws.New()` function:

~~~go
ws := sockets.New()
~~~

The `ws.Sockets` type has five fields:

~~~go
// Sockets is the main type for this library.
type Sockets struct {
    ClientChan      chan Payload                // The channel which receives message payloads.
    Clients         map[WebSocketConnection]any // A map of all connected clients.
    ErrorChan       chan error                  // A channel which receives errors.
    ReadBufferSize  int                         // I/O read buffer size in bytes. Defaults to 1024.
    WriteBufferSize int                         // I/O write buffer size in bytes. Defaults to 1024.
}
~~~

The `ws.Sockets` type has the following exposed methods:

~~~go
SocketEndPoint(w http.ResponseWriter, r *http.Request)  // A handler for the websocket endpoint.
ListenToWsChannel()                                     // A goroutine that listens to the SocketsChan and pushes data to broadcast function.
BroadcastTextToAll(payload JSONResponse)                // Pushes textual data to all connected clients.
BroadcastJSONToAll(payload JSONResponse)                // Pushes JSON data to all connected clients.
~~~

1. `SocketEndPoint`: A handler used to listen for (and upgrade) http(s) connections to ws(s) connections. 
This is what client side javascript will connect to.
2. `ListenToWsChannel`: Run this concurrently as a goroutine. It listens to the `ClientChan` 
(type `chan Payload`) in the `Sockets` type and sends client payloads to the appropriate broadcast method.
3. `BroadcastTextToAll`: sends a textual message to all connected clients.
4. `BroadcastJSONToAll`: sends a message in JSON format to all connected clients.

To *push data* over websockets from the client to the server, JSON must be able to be marshalled into the 
`ws.Payload` type:

~~~go
// Payload defines the data we receive from the client.
type Payload struct {
    MessageType int                 `json:"message_type"`   // ws.TextMessage or 1 - text message; ws.JSONMessage or 2: JSON message.
    Message     string              `json:"message"`        // The message.
    Data        any                 `json:"data,omitempty"` // A field for custom structured data.
    Conn        WebSocketConnection `json:"-"`              // Useful when you want to send a message to everyone except the originator.
}
~~~

Obviously, custom data types can simply be put in the `Data` field.

Data that comes back from the server to the client must conform to the `ws.JSONResponse` type:

~~~go
// JSONResponse defines the JSON we send back to client
type JSONResponse struct {
	Message     string              `json:"message"`
	Data        any                 `json:"data,omitempty"`
	CurrentConn WebSocketConnection `json:"-"`
}
~~~

Again, custom data of any type can be put into the `Data` field of this type.

## Sample app
A working web application can be [found here](https://github.com/tsawler/ws-sample-app).
