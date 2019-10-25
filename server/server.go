package server

import (
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	port  int
	token string
	key   string
	cert  string

	ipWhiteList []string

	server *http.Server
	log    *logger

	cache map[string]bool
	mux   *sync.RWMutex
}

func New(
	port int,
	token,
	key,
	cert string,
	ipWhiteList []string,
	v bool,
) *Server {
	return &Server{
		port:        port,
		token:       token,
		key:         key,
		cert:        cert,
		ipWhiteList: ipWhiteList,
		log:         newLogger(v),
		cache:       map[string]bool{},
		mux:         new(sync.RWMutex),
	}
}

func (s *Server) Run() {
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s,
	}
	s.log.Log(fmt.Sprintf("server running on :%d", s.port))
	if s.cert != "" && s.key != "" {
		if err := s.server.ListenAndServeTLS(s.cert, s.key); err != nil {
			s.log.Fatal(err)
		}
	} else {
		if err := s.server.ListenAndServe(); err != nil {
			s.log.Fatal(err)
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.processRequest(w, r)
}
