package server

import (
	"github.com/trapajim/snapmatch-ai/server/middleware"
	"log"
	"log/slog"
	"net/http"
)

// Middleware is a function type that wraps an http.Handler
// It allows adding custom logic before or after the handler is executed
type Middleware func(http.Handler) http.Handler

// Server struct holds the HTTP server configuration
type Server struct {
	Addr        string
	Mux         *http.ServeMux
	Middlewares []Middleware
}

// NewServer creates a new instance of Server with the given address
func NewServer(addr string, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	return &Server{
		Addr:        addr,
		Mux:         http.NewServeMux(),
		Middlewares: []Middleware{middleware.Logger(slog.Default())},
	}
}

// RegisterMiddleware adds a middleware to the server
func (s *Server) RegisterMiddleware(mw Middleware) {
	s.Middlewares = append(s.Middlewares, mw)
}

// RegisterRoute adds a route and its handler to the server's multiplexer
func (s *Server) RegisterRoute(pattern string, handler http.HandlerFunc) {
	wrappedHandler := http.Handler(handler)
	for i := len(s.Middlewares) - 1; i >= 0; i-- {
		wrappedHandler = s.Middlewares[i](wrappedHandler)
	}
	s.Mux.Handle(pattern, wrappedHandler)
}

// Start begins listening on the configured address
func (s *Server) Start() error {
	log.Printf("Starting server on %s", s.Addr)
	return http.ListenAndServe(s.Addr, s.Mux)
}
