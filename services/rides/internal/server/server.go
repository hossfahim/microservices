package server

import (
	"net/http"
	"rides/internal/database"
)

type Server struct {
	db              *database.Database
	usersServiceURL string
}

func NewServer(db *database.Database, usersServiceURL string) *Server {
	return &Server{db: db, usersServiceURL: usersServiceURL}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /rides", s.createRide)
	mux.HandleFunc("GET /rides/{id}", s.getRide)
	mux.HandleFunc("PATCH /rides/{id}/status", s.updateRideStatus)

	mux.ServeHTTP(w, r)
}
