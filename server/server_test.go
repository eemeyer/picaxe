package server_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eemeyer/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/t11e/picaxe/iiif"
	iiif_mocks "github.com/t11e/picaxe/iiif/mocks"
	resources_mocks "github.com/t11e/picaxe/resources/mocks"
	"github.com/t11e/picaxe/server"
)

func TestServer_ping(t *testing.T) {
	ts := httptest.NewServer(newHandler(server.ServerOptions{
		ResourceResolver: &resources_mocks.Resolver{},
		Processor:        &iiif_mocks.Processor{},
	}))
	defer ts.Close()

	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/ping", nil); resp != nil {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "picaxe", body)
	}
}

func TestServer_invalidParams(t *testing.T) {
	processor := &iiif_mocks.Processor{}
	processor.On("Process", "myidentifier/full/max/0/default.png",
		mock.Anything, mock.Anything).Return(iiif.InvalidSpec{
		Message: "not valid",
	})

	ts := httptest.NewServer(newHandler(server.ServerOptions{
		ResourceResolver: &resources_mocks.Resolver{},
		Processor:        processor,
	}))
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/iiif/myidentifier/full/max/0/default.png", nil)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `invalid request: not valid`, body)
}

func TestServer_escaping(t *testing.T) {
	processor := &iiif_mocks.Processor{}
	processor.On("Process", "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/0/default.png",
		mock.Anything, mock.Anything).Return(nil)

	ts := httptest.NewServer(newHandler(server.ServerOptions{
		ResourceResolver: &resources_mocks.Resolver{},
		Processor:        processor,
	}))
	defer ts.Close()

	resp, _ := testRequest(t, ts, "GET",
		"/api/picaxe/v1/iiif/http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/0/default.png", nil)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServer_iiifHandler(t *testing.T) {
	resolver := &resources_mocks.Resolver{}

	processor := &iiif_mocks.Processor{}
	processor.On("Process", "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/0/default.png",
		resolver, mock.Anything).Run(
		func(args mock.Arguments) {
			w := args.Get(2).(io.Writer)
			w.Write([]byte("result")) // Dummy image data
		}).Return(nil)

	ts := httptest.NewServer(newHandler(server.ServerOptions{
		ResourceResolver: resolver,
		Processor:        processor,
	}))
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET",
		"/api/picaxe/v1/iiif/http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/0/default.png", nil)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `result`, body)

	processor.AssertNumberOfCalls(t, "Process", 1)
}

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method, path string,
	body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

func newHandler(options server.ServerOptions) *chi.Mux {
	return server.NewServer(options).Handler()
}
