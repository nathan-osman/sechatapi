package sechatapi

import (
	"net"
	"net/http"

	"github.com/nathan-osman/go-sechat"
	"github.com/sirupsen/logrus"
)

// Config stores configuration information for the server.
type Config struct {
	Email    string
	Password string
	Token    string
}

// Server provides API methods for interacting with the Stack Exchange chat
// server.
type Server struct {
	stopped chan bool
	log     *logrus.Entry
	conn    *sechat.Conn
	mux     *http.ServeMux
	l       net.Listener
	token   string
}

// New creates a new Server instance with the provided configuration.
func New(cfg *Config) (*Server, error) {
	conn, err := sechat.New(cfg.Email, cfg.Password, 1)
	if err != nil {
		return nil, err
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	var (
		mux = http.NewServeMux()
		srv = http.Server{}
		s   = &Server{
			stopped: make(chan bool),
			log:     logrus.WithField("context", "sechatapi"),
			conn:    conn,
			mux:     mux,
			l:       l,
			token:   cfg.Token,
		}
	)
	srv.Handler = s
	mux.HandleFunc("/send", s.handleSend)
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

// ServeHTTP ensures the token was supplied if required.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Token") != s.token {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	s.mux.ServeHTTP(w, r)
}

// Close shuts down the server.
func (s *Server) Close() {
	s.l.Close()
	<-s.stopped
	s.conn.Close()
}
