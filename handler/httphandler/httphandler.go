package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ubiqueworks/go-solid-tutorial/service"
)

const (
	// ContentType is a constant for the HTTP header Content-Type
	ContentType = "Content-Type"
)

type httpError struct {
	Error string `json:"error"`
}

func New(s service.Service, bindAddr string, bindPort int) *HttpHandler {
	h := &HttpHandler{
		addr:    fmt.Sprintf("%s:%d", bindAddr, bindPort),
		service: s,
	}
	return h
}

type HttpHandler struct {
	addr    string
	service service.Service
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

// func (h HttpHandler) handleUsersGet(w http.ResponseWriter, r *http.Request) {
// 	userId := chi.URLParam(r, "userId")
// 	if _, err := uuid.Parse(userId); err != nil {
// 		_ = sendError(w, 400, err)
// 		return
// 	}
//
// }

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