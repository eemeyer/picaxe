package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//go:generate sh -c "mockery -name='ResourceResolver' -case=underscore"

type InvalidIdentifier struct {
	Identifier string
}

// Error implements interface "error".
func (err InvalidIdentifier) Error() string {
	return fmt.Sprintf("Invalid identifier %q", err.Identifier)
}

// ResourceResolver is an interface for something that can resolve a resource
// to a byte stream by its identifier.
type ResourceResolver interface {
	GetResource(identifier string) (io.ReadSeeker, error)
}

type httpResourceResolver struct {
	client *http.Client
}

// NewHTTPResourceResolver returns a resource resolver that downloads HTTP and HTTPS
// URLs.
func NewHTTPResourceResolver(client *http.Client) ResourceResolver {
	return &httpResourceResolver{
		client: client,
	}
}

// GetResource implements interface ResourceResolver.
func (h httpResourceResolver) GetResource(identifier string) (io.ReadSeeker, error) {
	u, ok := h.parseIdentifier(identifier)
	if !ok {
		return nil, InvalidIdentifier{Identifier: identifier}
	}

	resp, err := h.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Could not retrieve %s: Failed with HTTP status %d: %s",
			u, resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(&io.LimitedReader{
		R: resp.Body,
		N: maxBodyLength})
	if err != nil {
		return nil, fmt.Errorf("Error retrieving %s: %s", u, err)
	}
	if len(body) == maxBodyLength {
		return nil, fmt.Errorf("Rejecting %s: Body too large (%d)", u, len(body))
	}

	return bytes.NewReader(body), nil
}

func (h httpResourceResolver) parseIdentifier(identifier string) (*url.URL, bool) {
	u, err := url.Parse(strings.TrimSpace(identifier))
	if err != nil {
		return nil, false
	}

	if u.Scheme == "http" || u.Scheme == "https" {
		return u, true
	}

	return nil, false
}

// HTTPResourceResolver is the default HTTP resolver.
var HTTPResourceResolver = NewHTTPResourceResolver(&http.Client{
	Timeout: time.Duration(10 * time.Second),
})

const maxBodyLength = 10 * 1024 * 1024
