package picaxe_test

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/eemeyer/chi"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	picaxe "github.com/t11e/picaxe"
	"github.com/t11e/picaxe/mocks"
)

func TestPing(t *testing.T) {
	ts := httptest.NewServer(newHandler(picaxe.ServerOptions{
		ResourceResolver: &mocks.ResourceResolver{},
	}))
	defer ts.Close()

	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/ping", nil); resp != nil {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "picaxe", body)
	}
}

// Support limited subset of iiiaf `full` and `square` regions combined with size `!w,h`, `max`
func TestRequiredParameters(t *testing.T) {
	resolver := &mocks.ResourceResolver{}
	resolver.On("GetResource", "an-invalid-identifier").
		Return(nil, picaxe.InvalidIdentifier{Identifier: "an-invalid-identifier"})
	resolver.On("GetResource", "erroring-identifier").
		Return(nil, errors.New("fail"))

	ts := httptest.NewServer(newHandler(picaxe.ServerOptions{
		ResourceResolver: resolver,
	}))
	defer ts.Close()

	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/an-invalid-identifier/full/max/0/default.png", nil); resp != nil {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, `invalid identifier "an-invalid-identifier"`, body)
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/erroring-identifier/full/max/0/default.png", nil); resp != nil {
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, `internal error retrieving resource "erroring-identifier"`, body)
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/invalid/max/0/default.png", nil); resp != nil {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid or unsupported region 'invalid'", body)
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/0/0/default.png", nil); resp != nil {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid size '0'", body)
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/max/0/invalid.png", nil); resp != nil {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid or unsupported quality 'invalid'", body)
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/max/0/default.invalid", nil); resp != nil {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid or unsupported format 'invalid'", body)
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/max/0/default.tif", nil); resp != nil {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid or unsupported format 'tif'", body)
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/max/1/default.png", nil); resp != nil {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid or unsupported rotation '1'", body)
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/10,/0/default.png", nil); resp != nil {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid size '10,'", body)
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/,10/0/default.png", nil); resp != nil {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid size ',10'", body)
	}
}

const testImagePath = "/unit-test/image.jpg"

func TestConversion(t *testing.T) {
	testConversion(t, "testimages/baby-duck-10x10.png",
		fmt.Sprintf("/api/picaxe/v1/%s/square/10,10/0/default.png",
			url.QueryEscape(testImagePath)))
	testConversion(t, "testimages/baby-duck-10x10.jpg",
		fmt.Sprintf("/api/picaxe/v1/%s/square/10,10/0/default.jpg",
			url.QueryEscape(testImagePath)))
	testConversion(t, "testimages/baby-duck-10x10.gif",
		fmt.Sprintf("/api/picaxe/v1/%s/square/10,10/0/default.gif",
			url.QueryEscape(testImagePath)))
	testConversion(t, "testimages/baby-duck-full.png",
		fmt.Sprintf("/api/picaxe/v1/%s/full/max/0/default.png",
			url.QueryEscape(testImagePath)))
	testConversion(t, "testimages/baby-duck-full.jpg",
		fmt.Sprintf("/api/picaxe/v1/%s/full/max/0/default.jpg",
			url.QueryEscape(testImagePath)))
	testConversion(t, "testimages/baby-duck-full.gif",
		fmt.Sprintf("/api/picaxe/v1/%s/full/max/0/default.gif",
			url.QueryEscape(testImagePath)))
}

func testConversion(t *testing.T, expectedFile string, requestURL string) {
	imageData, err := ioutil.ReadFile("testimages/baby-duck.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	x, err := exif.Decode(bytes.NewBuffer(imageData))
	if err != nil {
		log.Fatal(err)
	}
	tag, err := x.Get("Software")
	require.NoError(t, err)
	require.Equal(t, `"ACD Systems Digital Imaging"`, tag.String())

	resolver := &mocks.ResourceResolver{}
	resolver.On("GetResource", "/unit-test/image.jpg").Return(bytes.NewReader(imageData), nil)

	ts := httptest.NewServer(newHandler(picaxe.ServerOptions{
		ResourceResolver: resolver,
	}))
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", requestURL, nil)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	expected, err := ioutil.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("cannot read expected image %s", expectedFile)
	}
	assert.Equal(t, string(expected), body, fmt.Sprintf("converted image %s", expectedFile))
	if data, err := exif.Decode(bytes.NewReader([]byte(body))); err == nil {
		t.Errorf("Expect to not be able to extract exif from converted image, but got %#v", data)
	}

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

func newHandler(options picaxe.ServerOptions) *chi.Mux {
	return picaxe.NewServer(options).Handler()
}
