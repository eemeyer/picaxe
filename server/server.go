package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/eemeyer/chi"
	"github.com/eemeyer/chi/middleware"
	"github.com/t11e/picaxe/iiif"
)

type ServerOptions struct {
	ResourceResolver ResourceResolver
	ProcessorFactory iiif.ProcessorFactory
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
	r.Get("/api/picaxe/v1/:identifier/:region/:size/:rotation/:rest", s.handleImage)
	return r
}

func (s *Server) handlePing(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("picaxe"))
}

var paramNames = []string{
	"region",
	"size",
	"rotation",
	"quality",
	"format",
}

func (s *Server) handleImage(w http.ResponseWriter, req *http.Request) {
	params := iiif.Params{}
	for _, key := range paramNames {
		params[key] = chi.URLParam(req, key)
	}

	if rest := chi.URLParam(req, "rest"); rest != "" {
		parts := strings.SplitN(chi.URLParam(req, "rest"), ".", 2)
		if len(parts) == 2 {
			params["quality"] = parts[0]
			params["format"] = parts[1]
		} else {
			// Defer to the processor to make sense of it
			params["quality"] = rest
		}
	}

	processor, err := s.ProcessorFactory.NewProcessor(params)
	if err != nil {
		if invalid, ok := err.(iiif.InvalidRequest); ok {
			respondWithError(w, http.StatusBadRequest, "invalid request: %s", invalid)
			return
		}
		log.Printf("Error getting processor: %s", err)
		respondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	identifier := chi.URLParam(req, "identifier")
	if identifier == "" {
		respondWithError(w, http.StatusBadRequest, "identifier not specified")
		return
	}

	reader, err := s.ResourceResolver.GetResource(identifier)
	if err != nil {
		if invalid, ok := err.(InvalidIdentifier); ok {
			respondWithError(w, http.StatusBadRequest, "invalid identifier %q", invalid.Identifier)
			return
		}
		respondWithError(w, http.StatusInternalServerError,
			"internal error retrieving resource %q", identifier)
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024*50))
	err = processor.Process(reader, buf)
	if err != nil {
		log.Printf("Error processing %q: %s", identifier, err)
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
