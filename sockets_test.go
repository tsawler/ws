package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func Test_New(t *testing.T) {
	s := New()
	if reflect.TypeOf(s).String() != "*ws.Sockets" {
		t.Error("wrong type; got", reflect.TypeOf(s))
	}
	if reflect.TypeOf(s.ClientChan).String() != "chan ws.WsPayload" {
		t.Error("wrong type; got", reflect.TypeOf(s))
	}
}

func TestWebSocketConnection(t *testing.T) {
	testWS := New()

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(testWS.SocketEndPoint))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	for {
		if len(testWS.Clients) > 0 {
			break
		}
	}

	payload := WsJsonResponse{
		Message: "hi",
	}
	testWS.BroadcastToAll(payload)
}
