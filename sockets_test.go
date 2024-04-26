package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	s := New()
	if reflect.TypeOf(s).String() != "*ws.Sockets" {
		t.Errorf("wrong type; expected %s but got %s", "*ws.Sockets", reflect.TypeOf(s).String())
	}
	if reflect.TypeOf(s.ClientChan).String() != "chan ws.Payload" {
		t.Errorf("wrong type; expected %s but got %s", "chan ws.Payload", reflect.TypeOf(s.ClientChan).String())
	}
	if reflect.TypeOf(s.Clients).String() != "map[ws.WebSocketConnection]string" {
		t.Errorf("wrong type; expected %s but got %s", "map[ws.WebSocketConnection]string", reflect.TypeOf(s.Clients).String())
	}
}

func TestWebSocketConnection(t *testing.T) {
	// Create test server.
	s := httptest.NewServer(http.HandlerFunc(testWS.SocketEndPoint))
	defer s.Close()

	// Convert http:// to ws://.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server.
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	// Wait for our connected client to show up as a map entry in testWs.Clients.
	time.Sleep(10 * time.Millisecond)

	// Create a payload.
	payload := JSONResponse{
		Message: "hi",
	}

	// Broadcast it.
	testWS.BroadcastJSONToAll(payload)

	// Read response.
	messageType, b, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if messageType != websocket.TextMessage {
		t.Errorf("wrong message type; expected 1 but got %d", messageType)
	}

	if !strings.Contains(string(b), `"message":"hi"`) {
		t.Errorf("wrong response; expected %s but got %s", `{"message":"hi"}`, string(b))
	}
}

func Test_ListenToWsChannel(t *testing.T) {
	// Create test server.
	s := httptest.NewServer(http.HandlerFunc(testWS.SocketEndPoint))
	defer s.Close()

	// Convert http//: to ws://.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server.
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	// fire off ListenToWsChannel.
	go testWS.ListenToWsChannel()

	payload := Payload{
		MessageType: JSONMessage,
		Message:     "Hello",
		Conn:        WebSocketConnection{ws},
	}
	testWS.ClientChan <- payload

	_, b, err := ws.ReadMessage()
	if err != nil {
		t.Error("failed to read")
	}

	if !strings.Contains(string(b), "Hello") {
		if err != nil {
			t.Error("response JSON does not have correct text")
		}
	}

	payload = Payload{
		MessageType: TextMessage,
		Message:     "Hello",
		Conn:        WebSocketConnection{ws},
	}
	testWS.ClientChan <- payload

	_, b, err = ws.ReadMessage()
	if err != nil {
		t.Error("failed to read")
	}

	if !strings.Contains(string(b), "Hello") {
		if err != nil {
			t.Error("response JSON does not have correct text")
		}
	}

}

func Test_listenForWS(t *testing.T) {
	// Create test server.
	s := httptest.NewServer(http.HandlerFunc(testWS.SocketEndPoint))
	defer s.Close()

	// Convert http//: to ws://.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server.
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	go testWS.listenForWS(&WebSocketConnection{ws})
	payload := Payload{
		Message: "",
	}
	err = ws.WriteJSON(payload)
	if err != nil {
		t.Fatalf("%v", err)
	}
}
