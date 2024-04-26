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
		t.Errorf("wrong type; expected %s but got %s", "*ws.Sockets", reflect.TypeOf(s).String())
	}
	if reflect.TypeOf(s.ClientChan).String() != "chan ws.WsPayload" {
		t.Errorf("wrong type; expected %s but got %s", "chan ws.WsPayload", reflect.TypeOf(s.ClientChan).String())
	}
	if reflect.TypeOf(s.Clients).String() != "map[ws.WebSocketConnection]string" {
		t.Errorf("wrong type; expected %s but got %s", "map[ws.WebSocketConnection]string", reflect.TypeOf(s.Clients).String())
	}
}

func TestWebSocketConnection(t *testing.T) {
	testWS := New()

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
	for {
		if len(testWS.Clients) > 0 {
			break
		}
	}

	// Create a payload.
	payload := WsJsonResponse{
		Message: "hi",
	}

	// Broadcast it.
	testWS.BroadcastToAll(payload)

	// Read response.
	messageType, b, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if messageType != 1 {
		t.Errorf("wrong message type; expected 1 but got %d", messageType)
	}

	if !strings.Contains(string(b), `"message":"hi"`) {
		t.Errorf("wrong response; expected %s but got %s", `{"message":"hi"}`, string(b))
	}
}

func Test_ListenToWsChannel(t *testing.T) {
	testWS := New()

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

	// Wait for our connected client to show up as a map entry in testWs.Clients.
	for {
		if len(testWS.Clients) > 0 {
			break
		}
	}

	// fire off
	go testWS.ListenToWsChannel()

	payload := WsPayload{
		Message: "Hello",
		Conn:    WebSocketConnection{ws},
	}
	testWS.ClientChan <- payload

	_, b, err := ws.ReadMessage()
	if err != nil {
		t.Error("failed to read")
	}

	t.Log("Bytes", string(b))
}

func Test_listenForWS(t *testing.T) {
	testWS := New()
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
	payload := WsPayload{
		Message: "",
	}
	err = ws.WriteJSON(payload)
	if err != nil {
		t.Fatalf("%v", err)
	}
}
