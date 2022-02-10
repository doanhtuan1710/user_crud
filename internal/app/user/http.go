package user

import (
	"net/http"
	"user_crud/internal/pkg/api"

	"github.com/go-chi/chi/v5"
)

type HTTPServer struct {
	*http.Server
}

func NewHTTPServer(userHandler *api.UserHandler) (httpServer *HTTPServer) {

	// Create multiplexing
	mux := chi.NewRouter()

	// User mux
	userMux := mux.Group(nil)
	userMux.Mount("/api/v1/user", userHandler.Route())

	// Assign mux to httpServer
	server := &http.Server{Handler: mux}
	httpServer = &HTTPServer{server}

	return
}
