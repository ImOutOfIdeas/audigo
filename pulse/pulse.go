package pulse

import (
	"net"

	"github.com/ImOutOfIdeas/audigo/internal"
)

var _ internal.Backend = (*pulseBackend)(nil)
var _ internal.Stream = (*pulseStream)(nil)

type pulseBackend struct {
	version    uint32   // Client/Server protocol version
	connection net.Conn // Connection to PulseAudio server
	tagId uint32        // Id of packet incremented each use
	clientIndex uint32 // Index of connection to server
}

type pulseStream struct {

}

func New() (internal.Backend, error) {
	b := pulseBackend{}

	// Connect to pulse server
	err := connect(&b)
	if err != nil {
		return &pulseBackend{}, err
	}

	return &b, nil
}

// === Backend Interface Methods ===

func (b *pulseBackend) OpenStream(internal.StreamConfig) (internal.Stream, error) {

	return nil, nil
}

func (b *pulseBackend) Close() error {
	b.connection.Close()
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

// === Pulse Backend Helpers ===
func (b *pulseBackend) next() uint32 {
    b.tagId++
    return b.tagId
}
