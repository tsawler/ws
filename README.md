# ws

Package ws is a simple library that makes it simple to use WebSockets (specifically 
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

The ws.Sockets type has two fields:

~~~go
// Sockets is the main type for this library.
type Sockets struct {
	ClientChan chan Payload // Data to be handled by this library is sent to this channel.
	Clients    map[WebSocketConnection]string // A map of all connected clients.
}
~~~

The ws.Sockets type has the following exposed methods:

~~~go
SocketEndPoint(w http.ResponseWriter, r *http.Request) // A handler to for the websocket endpoint.
ListenToWsChannel() // A goroutine that listens to the SocketsChan and pushes data to broadcast function
BroadcastTextToAll(payload JSONResponse) // Pushes textual data to all connected clients.
BroadcastJSONToAll(payload JSONResponse) // Pushes JSON data to all connected clients.
~~~

## Sample app
A working web application can be [found here](https://github.com/tsawler/ws-sample-app).
