package main

import "net/http"

func (server *Server) Ping(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Your connection is good",
	}

	_ = server.writeJSON(w, http.StatusAccepted, payload)
}

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {}
