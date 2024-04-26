# ws

Package ws is a simple library that makes it simple to use WebSockets (specifically 
[Gorilla Websockets](https://github.com/gorilla/websocket)) in your Go application.

## Installation
Install it the usual way:

~~~
go get github.com/tsawler/ws
~~~

## Usage
Create a variable of type ws.Sockets by calling the ws.New() function:

ws := sockets.New()

The ws.Sockets type has two fields:

~~~go
// Sockets is the main type for this library.
type Sockets struct {
	ClientChan chan WsPayload // Data to be handled by this library is sent to this channel.
	Clients    map[WebSocketConnection]string // A map of all connected clients.
}
~~~

The ws.Sockets type has the following exposed methods:

~~~go
SocketEndPoint(w http.ResponseWriter, r *http.Request) // A handler to for the websocket endpoint.
ListenToWsChannel() // A goroutine that listens to the SocketsChan and pushes data to broadcast function
BroadcastToAll(response WsJsonResponse) // Pushes data to all connected clients.
~~~

## Sample app
~~~go
package main

import (
	"fmt"
	"github.com/tsawler/page"
	"github.com/tsawler/toolbox"
	"github.com/tsawler/ws"
	"net/http"
	"time"
)

type application struct {
	ws        *ws.Sockets
	render    *page.Render
	eventChan chan string
}

const port = 8080

func main() {
	render := page.New()
	render.UseCache = false

	app := application{
		ws:        ws.New(),
		render:    render,
		eventChan: make(chan string, 100),
	}

	// start websocket functionality
	fmt.Println("Starting websocket functionality...")
	go app.ws.ListenToWsChannel()

	go app.RandomString()

	// start the web server
	fmt.Printf("Starting web server on port %d...\n", port)

	// create http server
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func (app *application) RandomString() {
	var t toolbox.Tools

	for {
		time.Sleep(3 * time.Second)
		for k, _ := range app.ws.Clients {
			payload := ws.WsPayload{
				Message: t.RandomString(5),
				Conn:    k,
			}
			app.ws.ClientChan <- payload
		}
	}
}
~~~
