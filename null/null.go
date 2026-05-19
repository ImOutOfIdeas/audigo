package null

import (
	"fmt"
	"github.com/ImOutOfIdeas/audigo/internal"
)

type nullBackend struct{}

func New() (internal.Backend, error) {
	return &nullBackend{}, nil
}

func (b *nullBackend) OpenStream(cfg internal.StreamConfig) (internal.Stream, error) {
	return &stream{cfg: cfg}, nil
}

func (b *nullBackend) Close() error {
	return nil
}

type stream struct {
	cfg     internal.StreamConfig
	running bool
}

func (s *stream) Start() error {
	s.running = true
	return nil
}

func (s *stream) Stop() error  {
	s.running = false; return nil
}

func (s *stream) Close() error {
	return nil
}

func (s *stream) Write(samples []float32) (int, error) {
	if !s.running {
		return 0, fmt.Errorf("stream not started")
	}
	return len(samples), nil
}
