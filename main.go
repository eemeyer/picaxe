package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

type contextKey int

const (
	httpClientKey contextKey = iota
)

type Options struct {
	ListenAddress string `short:"l" long:"listen" description:"Listen address." value-name:"ADDRESS"`
}

func main() {
	var options Options
	parser := flags.NewParser(&options, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
			return
		}
		return
	}

	listenAddress := ensureAddressWithPort(options.ListenAddress, 7073)
	fmt.Fprintf(os.Stdout, "listening on %s\n", listenAddress)
	handler := buildHTTPHandler()
	if err := http.ListenAndServe(listenAddress, handler); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
func ensureAddressWithPort(address string, defaultPort int) string {
	if address == "" {
		return fmt.Sprintf(":%d", defaultPort)
	} else if !strings.Contains(address, ":") {
		return fmt.Sprintf("%s:%d", address, defaultPort)
	}
	return address
}

func buildHTTPHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CloseNotify)
	r.Use(middleware.WithValue(httpClientKey, &http.Client{Timeout: time.Duration(10 * time.Second)}))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Get("/api/picaxe/ping", pingHandler)

	// TODO: add middleware to check that request is allowed
	r.Get("/api/picaxe/v1/get", resizeHandler)

	return r
}

func pingHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.WriteHeader(200)
	_, _ = resp.Write([]byte("picaxe"))
}

func resizeHandler(resp http.ResponseWriter, req *http.Request) {
	etag := md5sum(req.URL.Query().Encode())
	if match := req.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			resp.WriteHeader(http.StatusNotModified)
			return
		}
	}
	src, spec, err := makeProcessingSpec(req)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(err.Error()))
		return
	}
	httpClient, _ := req.Context().Value(httpClientKey).(http.Client)
	imgResp, err := httpClient.Get(src)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Unable to get %s: %v", src, err)))
		return
	}
	defer imgResp.Body.Close()
	if http.StatusOK != imgResp.StatusCode {
		resp.WriteHeader(imgResp.StatusCode)
		resp.Write([]byte(fmt.Sprintf("Unable to get %s. Got %s", src, imgResp.Status)))
		return
	}

	sourceImg, err := ioutil.ReadAll(imgResp.Body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Unable to retrieve %s: %v", src, err)))
		return
	}

	reader := bytes.NewReader(sourceImg)
	buffer := bytes.NewBuffer(make([]byte, 0, 1024*50))
	if err = ProcessImage(reader, buffer, *spec); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Error processing %s: %v", src, err)))
		return
	}

	resp.Header().Set("Content-type", "image/png")
	resp.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 365*24*60*60))
	resp.Header().Set("ETag", etag)
	resp.WriteHeader(http.StatusOK)
	resp.Write(buffer.Bytes())
}

func makeProcessingSpec(req *http.Request) (string, *ProcessingSpec, error) {
	query := req.URL.Query()
	src := query.Get("src")
	if src == "" {
		return "", nil, errors.New("src is required")
	}
	w, err := strconv.Atoi(query.Get("w"))
	if err != nil {
		return "", nil, errors.New("w is required")
	}
	h, err := strconv.Atoi(query.Get("h"))
	if err != nil {
		return "", nil, errors.New("h is required")
	}

	return src, &ProcessingSpec{
		Format:               ImageFormatPng,
		Trim:                 TrimModeFuzzy,
		TrimFuzzFactor:       0.5,
		Scale:                ScaleModeCover,
		ScaleWidth:           w,
		ScaleHeight:          h,
		Crop:                 CropModeCenter,
		CropWidth:            w,
		CropHeight:           h,
		NormalizeOrientation: true,
		Quality:              0.9,
	}, nil
}

func md5sum(query string) string {
	hasher := md5.New()
	hasher.Write([]byte(query))
	return hex.EncodeToString(hasher.Sum(nil))
}
