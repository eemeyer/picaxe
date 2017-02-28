package resources_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t11e/picaxe/resources"
)

func TestHTTPResolver_valid(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "image/png")
		w.WriteHeader(200)
		w.Write([]byte("hello"))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resolver := resources.NewHTTPResolver(http.DefaultClient)
	r, err := resolver.GetResource(ts.URL + "/foo.png")
	require.NoError(t, err)

	b, err := ioutil.ReadAll(r)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), b)
}

func TestHTTPResolver_fileURL(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	url := "file://" + cwd + "/http_test.go"

	resolver := resources.NewHTTPResolver(http.DefaultClient)
	_, err = resolver.GetResource(url)
	require.Error(t, err)
	assert.IsType(t, resources.InvalidIdentifier{}, err)
	e := err.(resources.InvalidIdentifier)
	assert.Equal(t, url, e.Identifier)
}

func TestHTTPResolver_nonHTTPScheme(t *testing.T) {
	url := "ftp://example.com/"
	resolver := resources.NewHTTPResolver(http.DefaultClient)
	_, err := resolver.GetResource(url)
	require.Error(t, err)
	assert.IsType(t, resources.InvalidIdentifier{}, err)
	e := err.(resources.InvalidIdentifier)
	assert.Equal(t, url, e.Identifier)
}
