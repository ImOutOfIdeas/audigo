package audigo

// Sample encoding
type SampleFormat uint8

const (
	Float32LE SampleFormat = iota
	Int16LE
	Int32LE
)

type Backend interface {
	OpenStream(StreamConfig) (Stream, error)
	Close() error
}

type Stream interface {
	Start() error
	Stop() error
	Close() error
	Write([]float32) (int, error)
}

// Callers requested stream parameters
type StreamConfig struct {
	Channels   int
	SampleRate float64
	BufferSize int
	Format 	   SampleFormat
}

// Backend provided stream parameters
// (may differ from config)
type StreamInfo struct {
	Channels   int
	SampleRate float64
	BufferSize int
	Latency    float64
}
