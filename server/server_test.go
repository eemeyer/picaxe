package server_test

import (
	"bytes"
	"fmt"
	"html/template"
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
	"github.com/t11e/picaxe/server"
	server_mocks "github.com/t11e/picaxe/server/mocks"
)

func TestServer_ping(t *testing.T) {
	ts := httptest.NewServer(newHandler(server.ServerOptions{
		ResourceResolver: &server_mocks.ResourceResolver{},
		ProcessorFactory: &iiif_mocks.ProcessorFactory{},
	}))
	defer ts.Close()

	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/ping", nil); resp != nil {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "picaxe", body)
	}
}

func TestServer_invalidParams(t *testing.T) {
	processorFactory := &iiif_mocks.ProcessorFactory{}
	processorFactory.On("NewProcessor", mock.Anything).Return(nil, iiif.InvalidRequest{
		Message: "not valid",
	})

	ts := httptest.NewServer(newHandler(server.ServerOptions{
		ResourceResolver: &server_mocks.ResourceResolver{},
		ProcessorFactory: processorFactory,
	}))
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/some-identifier/full/max/0/default.png", nil)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `invalid request: not valid`, body)
}

func TestServer_paramParsing(t *testing.T) {
	for _, format := range []string{"png", "jpg", "gif", "invalid"} {
		for _, region := range []string{"max", "full", "invalid"} {
			for _, size := range []string{"max", "invalid"} {
				for _, rotation := range []string{"0", "90", "180", "invalid"} {
					for _, quality := range []string{"default", "invalid"} {
						t.Run(fmt.Sprintf("region=%s,size=%s,format=%s,rotation=%s,quality=%s",
							region, size, format, rotation, quality),
							func(t *testing.T) {
								url := fmt.Sprintf("/api/picaxe/v1/my-identifier/%s/%s/%s/%s.%s",
									region, size, rotation, quality, format)

								resolver := &server_mocks.ResourceResolver{}
								resolver.On("GetResource", "my-identifier").Return(bytes.NewReader([]byte{}), nil)

								expectParams := iiif.Params{
									"region":   region,
									"size":     size,
									"rotation": rotation,
									"quality":  quality,
									"format":   format,
								}

								processor := &iiif_mocks.Processor{}
								processor.On("Process", mock.Anything, mock.Anything).Run(
									func(args mock.Arguments) {
										w := args.Get(1).(io.Writer)
										w.Write([]byte("result")) // Dummy image data
									}).Return(nil)

								factory := &iiif_mocks.ProcessorFactory{}
								factory.On("NewProcessor", expectParams).Return(processor, nil)

								ts := httptest.NewServer(newHandler(server.ServerOptions{
									ResourceResolver: resolver,
									ProcessorFactory: factory,
								}))
								defer ts.Close()
								resp, body := testRequest(t, ts, "GET", url, nil)
								require.NotNil(t, resp)
								assert.Equal(t, http.StatusOK, resp.StatusCode)
								assert.Equal(t, `result`, body)
							})
					}
				}
			}
		}
	}
}

func TestServer_resolver(t *testing.T) {
	resolver := &server_mocks.ResourceResolver{}
	resolver.On("GetResource", "my-identifier").Return(bytes.NewReader([]byte("hello")), nil)

	processor := &iiif_mocks.Processor{}
	processor.On("Process", mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			r := args.Get(0).(io.Reader)
			b, err := ioutil.ReadAll(r)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, []byte("hello"), b)

			w := args.Get(1).(io.Writer)
			w.Write([]byte("result")) // Dummy image data
		}).Return(nil)

	factory := &iiif_mocks.ProcessorFactory{}
	factory.On("NewProcessor", mock.Anything).Return(processor, nil)

	ts := httptest.NewServer(newHandler(server.ServerOptions{
		ResourceResolver: resolver,
		ProcessorFactory: factory,
	}))
	defer ts.Close()
	resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/my-identifier/full/max/0/default.png", nil)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `result`, body)

	resolver.AssertNumberOfCalls(t, "GetResource", 1)
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

func testUrl(options interface{}) string {
	urlTemplate, err := template.New("test").Parse("/api/picaxe/v1/{{urlquery .Id}}/{{.Region}}/{{.Size}}/{{.Rotation}}//{{.Quality}}.{{.Format}}")
	if err != nil {
		panic("cannot create URL template")
	}
	out := bytes.NewBuffer(nil)
	if err := urlTemplate.Execute(out, options); err != nil {
		panic(err)
	}
	return out.String()
}

func newHandler(options server.ServerOptions) *chi.Mux {
	return server.NewServer(options).Handler()
}
