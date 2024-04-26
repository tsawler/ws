package ws

import (
	"os"
	"testing"
)

var testWS *Sockets

func TestMain(m *testing.M) {
	testWS = New()
	os.Exit(m.Run())
}
