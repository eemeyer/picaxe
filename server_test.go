package main_test

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/eemeyer/chi"
	"github.com/rwcarlsen/goexif/exif"
	main "github.com/t11e/picaxe"
)

func TestPing(t *testing.T) {
	ts := httptest.NewServer(newHandler())
	defer ts.Close()
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/ping", nil); resp != nil {
		assertStatus(t, resp, http.StatusOK)
		assertEqual(t, "picaxe", body, "response body")
	}
}

// Support limited subset of iiiaf `full` and `square` regions combined with size `!w,h`, `max`
func TestRequiredParameters(t *testing.T) {
	ts := httptest.NewServer(newHandler())
	defer ts.Close()

	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/source.png/full/max/0/default.png", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "invalid identifier 'source.png'", body, "response body")
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/file:%2f%2fsource.png/full/max/0/default.png", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "invalid identifier 'file://source.png'", body, "response body")
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/invalid/max/0/default.png", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "invalid or unsupported region 'invalid'", body, "response body")
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/0/0/default.png", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "invalid size '0'", body, "response body")
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/max/0/invalid.png", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "invalid or unsupported quality 'invalid'", body, "response body")
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/max/0/default.invalid", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "invalid or unsupported format 'invalid'", body, "response body")
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/max/0/default.tif", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "invalid or unsupported format 'tif'", body, "response body")
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/max/1/default.png", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "invalid or unsupported rotation '1'", body, "response body")
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/10,/0/default.png", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "invalid size '10,'", body, "response body")
	}
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/http:%2f%2fexample.com%2fsource.png/full/,10/0/default.png", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "invalid size ',10'", body, "response body")
	}
}

func TestConversion(t *testing.T) {
	data, err := ioutil.ReadFile("testimages/baby-duck.jpeg")
	if err != nil {
		t.Fatalf("cannot read image file %v", err)
		return
	}
	x, _ := exif.Decode(bytes.NewBuffer(data))
	if tag, err := x.Get("Software"); err != nil || tag.String() != `"ACD Systems Digital Imaging"` {
		t.Fatalf("Original test image should have exif %v '%v'", err, tag.String())
	}
	testImagePath := "/unit-test/image.jpg"
	r := newHandler()
	r.Get(testImagePath, func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Content-Type", "image/jpeg")
		resp.Write(data)
	})
	ts := httptest.NewServer(r)
	testConversion(t, ts, "testimages/baby-duck-10x10.png", fmt.Sprintf("/api/picaxe/v1/%s/square/10,10/0/default.png", url.QueryEscape(ts.URL+testImagePath)))
	testConversion(t, ts, "testimages/baby-duck-10x10.jpg", fmt.Sprintf("/api/picaxe/v1/%s/square/10,10/0/default.jpg", url.QueryEscape(ts.URL+testImagePath)))
	testConversion(t, ts, "testimages/baby-duck-10x10.gif", fmt.Sprintf("/api/picaxe/v1/%s/square/10,10/0/default.gif", url.QueryEscape(ts.URL+testImagePath)))
	testConversion(t, ts, "testimages/baby-duck-full.png", fmt.Sprintf("/api/picaxe/v1/%s/full/max/0/default.png", url.QueryEscape(ts.URL+testImagePath)))
	testConversion(t, ts, "testimages/baby-duck-full.jpg", fmt.Sprintf("/api/picaxe/v1/%s/full/max/0/default.jpg", url.QueryEscape(ts.URL+testImagePath)))
	testConversion(t, ts, "testimages/baby-duck-full.gif", fmt.Sprintf("/api/picaxe/v1/%s/full/max/0/default.gif", url.QueryEscape(ts.URL+testImagePath)))
}

func testConversion(t *testing.T, ts *httptest.Server, expectedFile string, requestURL string) {
	if resp, body := testRequest(t, ts, "GET", requestURL, nil); resp != nil {
		assertStatus(t, resp, http.StatusOK)
		expected, err := ioutil.ReadFile(expectedFile)
		if err != nil {
			t.Fatalf("cannot read expected image %s", expectedFile)
		}
		assertEqual(t, string(expected), body, fmt.Sprintf("converted image %s", expectedFile))
		if data, err := exif.Decode(bytes.NewReader([]byte(body))); err == nil {
			t.Errorf("Expect to not be able to extract exif from converted image, but got %#v", data)
		}
	}
}

func Test403OfSrc(t *testing.T) {
	testImagePath := "/unit-test/image.jpg"
	r := newHandler()
	r.Get(testImagePath, func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusForbidden)
	})
	ts := httptest.NewServer(r)
	if resp, body := testRequest(t, ts, "GET", fmt.Sprintf("/api/picaxe/v1/%s/full/max/0/default.png", url.QueryEscape(ts.URL+testImagePath)), nil); resp != nil {
		assertStatus(t, resp, http.StatusForbidden)
		assertEqual(t, fmt.Sprintf("Unable to get %s%s. Got 403 Forbidden", ts.URL, testImagePath), body, "converted image")
	}
}

func assertEqual(t *testing.T, expected, actual interface{}, message string) {
	if expected != actual {
		t.Errorf("Expected %s to be '%v' but was '%v'", message, expected, actual)
	}
}

func assertStatus(t *testing.T, resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		t.Errorf("Expected status of %d but was %d", expected, resp.StatusCode)
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
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

func newHandler() *chi.Mux {
	return main.NewServer().Handler()
}
