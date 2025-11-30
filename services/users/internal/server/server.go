package server

import (
	"net/http"
	"users/internal/database"
)

type Server struct {
	db *database.Database
}

func NewServer(db *database.Database) *Server {
	return &Server{db: db}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /drivers", s.createDriver)
	mux.HandleFunc("GET /drivers", s.getDrivers)
	mux.HandleFunc("PATCH /drivers/{id}/status", s.setStatus)

	mux.HandleFunc("POST /passengers", s.createPassenger)
	mux.HandleFunc("GET /passengers", s.getPassengers)
	mux.HandleFunc("GET /passengers/{id}", s.getPassenger)
	mux.HandleFunc("PUT /passengers/{id}", s.updatePassenger)
	mux.HandleFunc("DELETE /passengers/{id}", s.deletePassenger)

	mux.ServeHTTP(w, r)
}
