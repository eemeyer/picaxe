package server_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/t11e/picaxe/iiif"
	iiif_mocks "github.com/t11e/picaxe/iiif/mocks"
	"github.com/t11e/picaxe/resources"
	resources_mocks "github.com/t11e/picaxe/resources/mocks"
	"github.com/t11e/picaxe/server"
)

func TestServer_ping(t *testing.T) {
	ts := newTestServer(server.ServerOptions{
		ResourceResolver: &resources_mocks.Resolver{},
		Processor:        &iiif_mocks.Processor{},
	})
	defer ts.Close()

	resp, body := doRequest(t, ts, "/api/picaxe/ping")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "picaxe", body)
}

func TestServer_invalidParams(t *testing.T) {
	ts := newTestServer(server.ServerOptions{
		ResourceResolver: &resources_mocks.Resolver{},
		Processor:        &iiif_mocks.Processor{},
	})
	defer ts.Close()

	resp, body := doRequest(t, ts, "/api/picaxe/v1/iiif/myidentifier/BEEP/max/0/default.png")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `invalid request: Not a valid set of coordinates: BEEP`, body)
}

func TestServer_loopDetection(t *testing.T) {
	ts := newTestServer(server.ServerOptions{
		ResourceResolver: &resources_mocks.Resolver{},
		Processor:        &iiif_mocks.Processor{},
	})
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/api/picaxe/v1/iiif/http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/0/default.png", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add(resources.HTTPHeaderPixace, "1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

type timeoutErr struct{}

func (err timeoutErr) Error() string   { return "game over, man" }
func (err timeoutErr) Timeout() bool   { return true }
func (err timeoutErr) Temporary() bool { return true }

func TestServer_timeouts(t *testing.T) {
	processor := &iiif_mocks.Processor{}
	processor.On("Process",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(timeoutErr{})

	ts := newTestServer(server.ServerOptions{
		ResourceResolver: &resources_mocks.Resolver{},
		Processor:        processor,
	})
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/api/picaxe/v1/iiif/http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/0/default.png", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestServer_iiifHandler(t *testing.T) {
	resolver := &resources_mocks.Resolver{}

	processor := &iiif_mocks.Processor{}
	processor.On("Process", iiif.Request{
		Identifier: "http://i.imgur.com/J1XaOIa.jpg",
		Region: iiif.Region{
			Kind: iiif.RegionKindFull,
		},
		Size: iiif.Size{
			Kind: iiif.SizeKindMax,
		},
		Format: iiif.FormatPNG,
	}, resolver, mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			w := args.Get(2).(io.Writer)
			w.Write([]byte("result")) // Dummy image data

			result := args.Get(3).(*iiif.Result)
			assert.NotNil(t, result)
			result.ContentType = "image/smurf"
		}).Return(nil)

	ts := newTestServer(server.ServerOptions{
		ResourceResolver: resolver,
		Processor:        processor,
		MaxAge:           time.Hour,
	})
	defer ts.Close()

	resp, body := doRequest(t, ts,
		"/api/picaxe/v1/iiif/http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/0/default.png")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "result", body)
	require.Equal(t, "image/smurf", resp.Header.Get("Content-Type"))
	require.Equal(t, "public,s-maxage=3600", resp.Header.Get("Cache-Control"))
	require.Equal(t, "220ecea9df8e805373fd15a35bc29ef1d5c9119e7b139f3e6c271a2a9bdd79be", resp.Header.Get("ETag"))

	processor.AssertNumberOfCalls(t, "Process", 1)
}

func doRequest(t *testing.T, ts *httptest.Server, path string) (*http.Response, string) {
	req, err := http.NewRequest("GET", ts.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

func newTestServer(options server.ServerOptions) *httptest.Server {
	handler := server.NewServer(options).Handler()
	return httptest.NewServer(handler)
}
