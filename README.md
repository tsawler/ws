<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/tsawler/persist/blob/main/LICENSE.md)
<a href="https://pkg.go.dev/github.com/tsawler/ws"><img src="https://img.shields.io/badge/godoc-reference-%23007d9c.svg"></a>
![Tests](https://github.com/tsawler/ws/actions/workflows/tests.yml/badge.svg)

# ws

Package ws is a simple library that makes it easy to use WebSockets (specifically 
[Gorilla Websockets](https://github.com/gorilla/websocket)) in your Go application.

## Installation
Install it the usual way:

~~~
go get github.com/tsawler/ws
~~~

## Usage
Create a variable of type ws.Sockets by calling the `ws.New()` function:

~~~go
ws := sockets.New()
~~~

The ws.Sockets type has three fields:

~~~go
// Sockets is the main type for this library.
type Sockets struct {
    ClientChan chan Payload // The channel that receives messages.
    Clients    map[WebSocketConnection]string // A map of connected clients.
    ErrorChan  chan error // A channel to send errors to.
}
~~~

The ws.Sockets type has the following exposed methods:

~~~go
SocketEndPoint(w http.ResponseWriter, r *http.Request) // A handler for the websocket endpoint.
ListenToWsChannel() // A goroutine that listens to the SocketsChan and pushes data to broadcast function.
BroadcastTextToAll(payload JSONResponse) // Pushes textual data to all connected clients.
BroadcastJSONToAll(payload JSONResponse) // Pushes JSON data to all connected clients.
~~~

1. `SocketEndPoint`: You'll need a handler to listen for (and upgrade) http(s) connections to ws(s) connections. 
This is what client side javascript will connect to.
2. `ListenToWsChannel`: Run this concurrently() as a goroutine. It listens to the Clients field
in the `Sockets` type and sends client payloads to the appropriate broadcast method.
3. `BroadcastTextToAll`: sends a textual message to all connected clients.
4. `BroadcastJSONToAll`: sends a message in JSON format to all connected clients.



## Sample app
A working web application can be [found here](https://github.com/tsawler/ws-sample-app).
