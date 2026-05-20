//go:build darwin

package audigo

import (
	"errors"
	"github.com/ImOutOfIdeas/audigo/internal"
)

func DefaultBackend() (internal.Backend, error) {
	return nil, errors.New("audigo: darwin is not supported yet")
}
