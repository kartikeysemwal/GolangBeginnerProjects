package main

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (server *Server) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Links"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Mount("/debug/pprof", http.HandlerFunc(pprof.Index))
	mux.Mount("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Mount("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Mount("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Mount("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	mux.Get("/", server.Ping)
	mux.Post("/user", server.CreateUser)
	mux.Get("/user", server.ReadUser)
	mux.Patch("/user", server.UpdateUser)
	mux.Delete("/user", server.DeleteUser)

	return mux
}
