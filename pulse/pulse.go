package pulse

import (
	"net"

	"github.com/ImOutOfIdeas/audigo/internal"
)

// Default TCP port for hostname only command string

// Data necassary for holding a pulse server connection
type pulseBackend struct {
	connection net.Conn
	seq        int
}


func New() (internal.Backend, error) {

	return b, nil
}

// === Interface Methods ===

func (b pulseBackend) OpenStream(internal.StreamConfig) (internal.Stream, error) {
	return nil, nil
}

func (b pulseBackend) Close() error {
	return nil
}
