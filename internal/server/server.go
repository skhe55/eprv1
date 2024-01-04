package server

import (
	"context"
	"erpv1/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
)

type Server struct {
	Addr string
	cfg  *config.Config
}

func NewServer(addr string, cfg *config.Config) *Server {
	return &Server{
		Addr: addr,
		cfg:  cfg,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	s.MapHandlers(mux)
	server := &http.Server{
		Addr:    s.Addr,
		Handler: mux,
	}

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodPut,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	server.Handler = cors.Handler(mux)

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server shutdown error: %v", err)
		}
		close(stopped)
	}()

	log.Printf("Starting http server on http://localhost%s", s.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP Server ListenAndServe error: %v", err)
	}

	<-stopped

	log.Printf("\nGraceful shutdown")
	return nil
}
