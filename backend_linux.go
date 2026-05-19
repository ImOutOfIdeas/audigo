//go:build linux

package audigo

import (
    "github.com/ImOutOfIdeas/audigo/internal"
    "github.com/ImOutOfIdeas/audigo/pulse"
)

func DefaultBackend() (internal.Backend, error) {
    return pulse.New()
}
