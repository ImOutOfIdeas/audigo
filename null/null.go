package null

import (
	"fmt"
	"github.com/imoutofideas/audigo"
)

type Backend struct{}

func New() *Backend {
	return &Backend{}
}

func (b *Backend) OpenStream(cfg audigo.StreamConfig) (audigo.Stream, error) {
	return &stream{cfg: cfg}, nil
}

func (b *Backend) Close() error {
	return nil
}

type stream struct {
	cfg     audigo.StreamConfig
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
