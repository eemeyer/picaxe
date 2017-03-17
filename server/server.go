package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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
	"github.com/t11e/picaxe/resources"
)

type ServerOptions struct {
	ResourceResolver resources.Resolver
	Processor        iiif.Processor
	MaxAge           time.Duration
}

type Server struct {
	ServerOptions
	cacheControlHeader string
}

func NewServer(opts ServerOptions) *Server {
	cacheControlHeader := fmt.Sprintf("public,s-maxage=%0.f", opts.MaxAge.Seconds())
	return &Server{ServerOptions: opts, cacheControlHeader: cacheControlHeader}
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
	r.Get("/api/picaxe/ping", s.handlePing)
	r.Get("/api/picaxe/v1/iiif/*", s.handleImage)
	return r
}

func (s *Server) handlePing(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("picaxe"))
}

func (s *Server) handleImage(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(resources.HTTPHeaderPixace) != "" {
		log.Printf("Request contains loop-detecting header %q, refusing", resources.HTTPHeaderPixace)
		writeError(w, http.StatusForbidden, "loop detected")
		return
	}

	spec := chi.URLParam(r, "*")
	if r.URL.RawQuery != "" {
		spec = spec + "?" + r.URL.RawQuery
	}

	req, err := iiif.ParseSpec(spec)
	if err != nil {
		returnError(w, err)
		return
	}

	etag := buildETagFromRequest(req)
	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024*50))

	var result iiif.Result
	if err := s.Processor.Process(*req, s.ResourceResolver, buf, &result); err != nil {
		returnError(w, err)
		return
	}

	w.Header().Set("Content-type", result.ContentType)
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", s.cacheControlHeader)
	w.WriteHeader(http.StatusOK)
	io.Copy(w, buf)
}

func returnError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case net.Error:
		if e.Timeout() {
			writeError(w, http.StatusServiceUnavailable, "timed out")
			return
		}
	case resources.InvalidIdentifier:
		writeError(w, http.StatusBadRequest, "invalid identifier %q", e.Identifier)
		return
	case iiif.InvalidSpec:
		writeError(w, http.StatusBadRequest, "invalid request: %s", e)
		return
	}

	log.Printf("Error: %s", err)
	writeError(w, http.StatusInternalServerError, "internal error")
}

func writeError(w http.ResponseWriter, statusCode int, format string, args ...interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-type", "text/plain")
	w.Write([]byte(fmt.Sprintf(format, args...)))
}

func buildETagFromRequest(req *iiif.Request) string {
	hasher := sha256.New()
	hasher.Write([]byte(req.String()))
	hasher.Write([]byte(cacheVersion))
	return hex.EncodeToString(hasher.Sum(nil))
}

var cacheVersion = "1" // Increase to bust cache
