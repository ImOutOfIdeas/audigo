package pulse

import (
	"net"

	"github.com/ImOutOfIdeas/audigo/internal"
)

type pulseBackend struct {
	version    uint32   // Client/Server protocol version
	connection net.Conn // Connection to PulseAudio server
	sequenceId uint32   // Id of packet incremented each use
}

type pulseStream struct {

}

func New() (internal.Backend, error) {
	b := pulseBackend{}

	conn, version, err := connect()
	if err != nil {
		return &pulseBackend{}, err
	}
	b.connection = conn
	b.version = version

	println("auth reply ok")
	println("version: ", b.version)

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
