package pulse

import (
	"net"

	"github.com/ImOutOfIdeas/audigo/internal"
)

// Data necassary for holding a pulse server connection
type pulseBackend struct {
	connection net.Conn
	seq        int
}

type pulseStream struct {

}

func New() (internal.Backend, error) {
	b := pulseBackend{}



	return &b, nil
}

// === Backend Interface Methods ===

func (b *pulseBackend) OpenStream(internal.StreamConfig) (internal.Stream, error) {
	return nil, nil
}

func (b *pulseBackend) Close() error {
	return nil
}

// === Stream Interface Methods ===

func (s *pulseStream) Start() error {
	return nil
}

func (s *pulseStream) Stop() error {
	return nil
}

func (s *pulseStream) Close() error {
	return nil
}

func (s *pulseStream) Write(buf []float32) (int, error) {
	return len(buf), nil
}
