package main

import (
	"log"
	"net/http"
)

type Server struct {
	listenAddr string
	routes     map[string]http.HandlerFunc
}

func NewServer(listenAddr string) *Server {
	return &Server{listenAddr: listenAddr, routes: make(map[string]http.HandlerFunc)}
}

// register a route with the server
func (s *Server) RegisterRoute(path string, handler http.HandlerFunc) {
	s.routes[path] = handler
}

// initialize server with default routes
func (s *Server) Init() {
	s.RegisterRoute("~default~", defaultHandler)
	s.RegisterRoute("/info", routeInfo)
}

func (s *Server) httpHandler(w http.ResponseWriter, r *http.Request) {

	reqPath := r.URL.Path
	for path := range s.routes {
		handler := s.routes[path]

		if reqPath == path {
			handler(w, r)
			return
		}
	}

	defaultHandler := s.routes["~default~"]
	defaultHandler(w, r)
}

func (s *Server) Run() {
	log.Println("Server started at :", s.listenAddr)

	router := http.NewServeMux()

	router.HandleFunc("/", s.httpHandler)
	s.Init()

	http.ListenAndServe(s.listenAddr, router)
}
