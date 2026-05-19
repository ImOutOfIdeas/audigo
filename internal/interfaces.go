package internal

// The audio backend of choice
type Backend interface {
	OpenStream(StreamConfig) (Stream, error)
	Close() error
}

// A single audio stream
type Stream interface {
	Start() error
	Stop() error
	Close() error
	Write([]float32) (int, error)
}
