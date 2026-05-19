// audigo.go
package audigo

import "github.com/ImOutOfIdeas/audigo/internal"

// Re-export internal types so this package is a standalone import
type Backend = internal.Backend
type Stream = internal.Stream
type StreamConfig = internal.StreamConfig
type StreamInfo = internal.StreamInfo
type SampleFormat = internal.SampleFormat
