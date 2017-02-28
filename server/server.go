package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/eemeyer/chi"
	"github.com/eemeyer/chi/middleware"
	"github.com/t11e/picaxe/iiif"
	"github.com/t11e/picaxe/resources"
)

type ServerOptions struct {
	ResourceResolver resources.Resolver
	Processor        iiif.Processor
}

type Server struct {
	ServerOptions
}

func NewServer(opts ServerOptions) *Server {
	return &Server{ServerOptions: opts}
}

func (s *Server) Run(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Unable to listen on address %q: %s", address, err.Error())
	}
	log.Printf("Listening on %s", address)
	return http.Serve(listener, s.Handler())
}

func (s *Server) Handler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.CloseNotify)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(QueryETagMatcher)
	r.Get("/api/picaxe/ping", s.handlePing)
	r.Get("/api/picaxe/v1/iiif/*", s.handleImage)
	return r
}

func (s *Server) handlePing(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("picaxe"))
}

func (s *Server) handleImage(w http.ResponseWriter, req *http.Request) {
	spec := chi.URLParam(req, "*")

	buf := bytes.NewBuffer(make([]byte, 0, 1024*50))

	err := s.Processor.Process(spec, s.ResourceResolver, buf)
	if err != nil {
		if invalid, ok := err.(resources.InvalidIdentifier); ok {
			respondWithError(w, http.StatusBadRequest, "invalid identifier %q", invalid.Identifier)
			return
		}

		if invalid, ok := err.(iiif.InvalidSpec); ok {
			respondWithError(w, http.StatusBadRequest, "invalid request: %s", invalid)
			return
		}

		log.Printf("Error processing %q: %s", spec, err)
		respondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-type", "image/png")
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 365*24*60*60))
	w.Header().Set("ETag", buildETagFromRequest(req))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, buf)
}

func respondWithError(w http.ResponseWriter, statusCode int, format string, args ...interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-type", "text/plain")
	w.Write([]byte(fmt.Sprintf(format, args...)))
}
