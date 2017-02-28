package resources

import (
	"fmt"
	"io"
)

//go:generate sh -c "mockery -name='Resolver' -case=underscore"

type InvalidIdentifier struct {
	Identifier string
}

// Error implements interface "error".
func (err InvalidIdentifier) Error() string {
	return fmt.Sprintf("Invalid identifier %q", err.Identifier)
}

// Resolver is an interface for something that can resolve a resource
// to a byte stream by its identifier.
type Resolver interface {
	GetResource(identifier string) (io.ReadSeeker, error)
}
