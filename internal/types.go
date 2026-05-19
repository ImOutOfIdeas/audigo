package internal

// Sample encoding
type SampleFormat uint8
const (
	Float32LE SampleFormat = iota
	Int16LE
	Int32LE
)

// Callers requested stream parameters
type StreamConfig struct {
	Channels   int
	SampleRate float64
	BufferSize int
	Format 	   SampleFormat
}

// Backend provided stream parameters
// (may differ from StreamConfig)
type StreamInfo struct {
	Channels   int
	SampleRate float64
	BufferSize int
	Latency    float64
}
