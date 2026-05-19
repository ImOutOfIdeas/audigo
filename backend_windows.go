//go:build windows

package audigo

import (
    "errors"
    "github.com/ImOutOfIdeas/audigo/internal"
)

func DefaultBackend() (internal.Backend, error) {
    return nil, errors.New("audigo: windows is not supported yet")
}
