package ws

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"testing"
)

func testRoutes() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/ws", func(w http.ResponseWriter, r *http.Request) {

	})

	return mux
}

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
