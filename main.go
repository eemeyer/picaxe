package main

import (
	"bytes"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"golang.org/x/net/context"
)

type contextKey int

const (
	httpClientKey contextKey = iota
)

func main() {
	ctx := context.WithValue(context.Background(), httpClientKey, &http.Client{})

	handler := buildHTTPHandler(ctx, xmux.New())

	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func buildHTTPHandler(ctx context.Context, mux *xmux.Mux) http.Handler {
	mux.GET("/api/picaxe/ping", xhandler.HandlerFuncC(pingHandler))

	// TODO: add middleware to check that request is allowed
	mux.GET("/api/picaxe/scale", xhandler.HandlerFuncC(resizeHandler))

	chain := xhandler.Chain{}
	chain.UseC(xhandler.CloseHandler)
	return xhandler.New(ctx, chain.HandlerC(mux))
}

func pingHandler(_ context.Context, resp http.ResponseWriter, _ *http.Request) {
	resp.WriteHeader(200)
	_, _ = resp.Write([]byte("picaxe"))
}

func resizeHandler(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	src := req.FormValue("src")
	var okay = src != ""
	w, err := strconv.Atoi(req.FormValue("w"))
	if err != nil {
		okay = false
	}
	h, err := strconv.Atoi(req.FormValue("h"))
	if err != nil {
		okay = false
	}

	if !okay {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("src, w, and h are required"))
		return
	}
	httpClient, _ := ctx.Value(httpClientKey).(http.Client)
	imgResp, err := httpClient.Get(src)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Unable to get %s: %v", src, err)))
		return
	}
	defer imgResp.Body.Close()
	sourceImg, err := ioutil.ReadAll(imgResp.Body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Unable to retrieve %s: %v", src, err)))
		return
	}
	reader := bytes.NewReader(sourceImg)
	//TODO : write cache headers - maybe need to buffer ProcessImage output to do this though
	resp.Header().Set("Content-type", "image/png")

	spec := ProcessingSpec{
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
	}

	if err = ProcessImage(reader, resp, spec); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Error processing %s: %v", src, err)))
		return
	}
}
