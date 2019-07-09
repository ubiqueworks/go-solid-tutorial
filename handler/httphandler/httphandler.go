//go:generate mockgen -package httphandler -source=httphandler.go -destination httphandler_mock.go

package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ubiqueworks/go-solid-tutorial/domain"
)

const (
	// ContentType is a constant for the HTTP header Content-Type
	ContentType = "Content-Type"
)

type config struct {
	BindAddr string `env:"BIND_ADDR" envDefault:"0.0.0.0" required:"true"`
	Port     int    `env:"HTTP_PORT" envDefault:"8000" required:"true"`
}

func (c config) Validate() error {
	if c.BindAddr == "" {
		return fmt.Errorf("missing required HTTP_ADDR")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid HTTP_PORT: %d", c.Port)
	}

	return nil
}

type httpError struct {
	Error string `json:"error"`
}

type Service interface {
	ListUsers() ([]domain.User, error)
}

func New(s Service) (*HttpHandler, error) {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	h := &HttpHandler{
		addr:    fmt.Sprintf("%s:%d", cfg.BindAddr, cfg.Port),
		service: s,
	}
	return h, nil
}

type HttpHandler struct {
	addr    string
	service Service
}

func (h HttpHandler) Bootstrap() (shutdownFn func() error, err error) {
	log.Printf("starting HTTP server on %s", h.addr)

	router, err := h.createRouter()
	if err != nil {
		return
	}

	listener, err := net.Listen("tcp", h.addr)
	if err != nil {
		return
	}

	server := &http.Server{
		Addr:    h.addr,
		Handler: router,
	}

	shutdownFn = func() error {
		defer func() {
			log.Println("HTTP server shutdown")
		}()
		stopCtx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelTimeout()

		log.Println("shutting down HTTP server...")
		return server.Shutdown(stopCtx)
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return
}

func (h HttpHandler) createRouter() (chi.Router, error) {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

	r.Get("/users", h.handleUsersQuery)

	return r, nil
}

func (h HttpHandler) handleUsersQuery(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers()
	if err != nil {
		_ = sendError(w, 500, err)
		return
	}
	_ = sendJson(w, 200, users)
}

func sendError(w http.ResponseWriter, statusCode int, err error) error {
	return sendJson(w, statusCode, &httpError{
		Error: err.Error(),
	})
}

func sendJson(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Add(ContentType, "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func sendStatus(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}
