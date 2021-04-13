package httpsrv

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type Server struct {
	seq Seq64Generator

	mu      sync.Mutex
	addr    string
	logger  logrus.FieldLogger
	srv     *http.Server
	started bool
}

// New initializes a new ready-to-use HTTP server.
func New(logger logrus.FieldLogger, seq Seq64Generator, addr string) *Server {
	return &Server{
		seq:    seq,
		mu:     sync.Mutex{},
		addr:   addr,
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return errors.New("already started")
	}

	return s.start()
}

func (s *Server) start() error {
	s.started = true
	defer func() { s.started = false }()

	h := NewGetID(s.logger.WithField("handler", "GetID"), s.seq)

	mux := http.NewServeMux()
	mux.Handle("/id64", h)

	s.srv = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	return s.srv.ListenAndServe()
}
