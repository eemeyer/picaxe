package resources

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

type httpResolver struct {
	client *http.Client
}

// NewHTTPResolver returns a resource resolver that downloads HTTP and HTTPS
// URLs.
func NewHTTPResolver(client *http.Client) Resolver {
	return &httpResolver{
		client: client,
	}
}

// GetResource implements interface Resolver.
func (h httpResolver) GetResource(identifier string) (io.ReadSeeker, error) {
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

func (h httpResolver) parseIdentifier(identifier string) (*url.URL, bool) {
	u, err := url.Parse(strings.TrimSpace(identifier))
	if err != nil {
		return nil, false
	}

	if u.Scheme == "http" || u.Scheme == "https" {
		return u, true
	}

	return nil, false
}

// HTTPResolver is the default HTTP resolver.
var HTTPResolver = NewHTTPResolver(&http.Client{
	Timeout: time.Duration(10 * time.Second),
})

const maxBodyLength = 10 * 1024 * 1024
