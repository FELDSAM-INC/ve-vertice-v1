package httpd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	ln   net.Listener
	addr string
	err  chan error

	Handler *Handler
}

// NewService returns a new instance of Service.
func NewService(c httpd.Config) *Service {
	s := &Service{
		addr: c.BindAddress,
		err:  make(chan error),
		Handler: NewHandler(
			"0.9",
		),
	}
	return s
}

// Open starts the service
func (s *Service) Open() error {
	s.Logger.Println("Starting HTTP service")

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.Logger.Println("Listening on HTTP:", listener.Addr().String())
	s.ln = listener

	// Begin listening for requests in a separate goroutine.
	go s.serve()
	return nil
}

// Close closes the underlying listener.
func (s *Service) Close() error {
	if s.ln != nil {
		return s.ln.Close()
	}
	return nil
}

// Err returns a channel for fatal errors that occur on the listener.
func (s *Service) Err() <-chan error { return s.err }

// Addr returns the listener's address. Returns nil if listener is closed.
func (s *Service) Addr() net.Addr {
	if s.ln != nil {
		return s.ln.Addr()
	}
	return nil
}

// serve serves the handler from the listener.
func (s *Service) serve() {
	// The listener was closed so exit
	// See https://github.com/golang/go/issues/4373
	err := http.Serve(s.ln, s.Handler)
	if err != nil && !strings.Contains(err.Error(), "closed") {
		s.err <- fmt.Errorf("listener failed: addr=%s, err=%s", s.Addr(), err)
	}
}
