package resources

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const HTTPHeaderPixace = "X-Picaxe"

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
	u := strings.TrimSpace(identifier)

	if err := h.validateIdentifier(u); err != nil {
		return nil, InvalidIdentifier{
			Message:    err.Error(),
			Identifier: identifier,
		}
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(HTTPHeaderPixace, "1") // Used to prevent loops

	resp, err := h.client.Do(req)
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

func (h httpResolver) validateIdentifier(identifier string) error {
	u, err := url.Parse(identifier)
	if err != nil {
		return errors.New("not a valid URL")
	}

	if !(u.Scheme == "http" || u.Scheme == "https") {
		return fmt.Errorf("not a valid scheme: %q", u.Scheme)
	}

	return nil
}

// HTTPResolver is the default HTTP resolver.
var HTTPResolver = NewHTTPResolver(&http.Client{
	Timeout: time.Duration(10 * time.Second),
})

const maxBodyLength = 10 * 1024 * 1024
