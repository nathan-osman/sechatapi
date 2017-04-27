package sechatapi

import (
	"net"
	"net/http"

	"github.com/nathan-osman/go-sechat"
)

// Config stores configuration information for the server.
type Config struct {
	Addr     string
	Email    string
	Password string
}

// Server provides API methods for interacting with the Stack Exchange chat
// server.
type Server struct {
	stopped chan bool
	conn    *sechat.Conn
	mux     *http.ServeMux
	l       net.Listener
}

// New creates a new Server instance with the provided configuration.
func New(cfg *Config) (*Server, error) {
	conn, err := sechat.New(cfg.Email, cfg.Password, 1)
	if err != nil {
		return nil, err
	}
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		mux = http.NewServeMux()
		srv = http.Server{
			Handler: mux,
		}
		s = &Server{
			stopped: make(chan bool),
			conn:    conn,
			mux:     mux,
			l:       l,
		}
	)
	mux.HandleFunc("/auth/login", s.handleLogin)
	go func() {
		defer close(s.stopped)
		srv.Serve(l)
	}()
	return s, nil
}

// Addr retrieves the address of the server.
func (s *Server) Addr() string {
	return s.l.Addr().String()
}

// Close shuts down the server.
func (s *Server) Close() {
	s.l.Close()
	<-s.stopped
	s.conn.Close()
}
