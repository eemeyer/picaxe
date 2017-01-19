package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pressly/chi"
)

func TestPing(t *testing.T) {
	ts := httptest.NewServer(buildHTTPHandler())
	defer ts.Close()
	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/ping", nil); resp != nil {
		assertStatus(t, resp, http.StatusOK)
		assertEqual(t, "picaxe", body, "response body")
	}
}

func TestRequiredParameters(t *testing.T) {
	ts := httptest.NewServer(buildHTTPHandler())
	defer ts.Close()

	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/get", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "src is required", body, "response body")
	}

	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/get?src=file://testimages/baby-duck.jpeg", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "w is required", body, "response body")
	}

	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/get?src=file://testimages/baby-duck.jpeg&w=10", nil); resp != nil {
		assertStatus(t, resp, http.StatusBadRequest)
		assertEqual(t, "h is required", body, "response body")
	}

	if resp, body := testRequest(t, ts, "GET", "/api/picaxe/v1/get?src=file://testimages/baby-duck.jpeg&w=10&h=10", nil); resp != nil {
		assertStatus(t, resp, http.StatusInternalServerError)
		assertEqual(t, `Unable to get file://testimages/baby-duck.jpeg: Get file://testimages/baby-duck.jpeg: unsupported protocol scheme "file"`, body, "response body")
	}
}

func TestConversion(t *testing.T) {
	r := buildHTTPHandler().(chi.Router)
	testImagePath := "/unit-test/image.jpg"
	r.Get(testImagePath, func(resp http.ResponseWriter, req *http.Request) {
		data, err := ioutil.ReadFile("testimages/baby-duck.jpeg")
		if err != nil {
			t.Fatalf("cannot read image file %v", err)
			return
		}
		resp.Header().Add("Content-Type", "image/jpeg")
		resp.Write(data)
	})
	ts := httptest.NewServer(r)
	if resp, body := testRequest(t, ts, "GET", fmt.Sprintf("/api/picaxe/v1/get?src=%s%s&w=10&h=10", ts.URL, testImagePath), nil); resp != nil {
		assertStatus(t, resp, http.StatusOK)
		expected, err := ioutil.ReadFile("testimages/baby-duck-10x10.png")
		if err != nil {
			t.Fatal("cannot read expected image")
		}
		assertEqual(t, string(expected), body, "converted image")
	}

}

func assertEqual(t *testing.T, expected, actual interface{}, message string) {
	if expected != actual {
		t.Fatalf("Expected %s to be '%v' but was '%v'", message, expected, actual)
	}
}

func assertStatus(t *testing.T, resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		t.Fatalf("Expected status of %d but was %d", expected, resp.StatusCode)
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
