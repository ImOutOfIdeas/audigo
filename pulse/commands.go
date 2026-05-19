// Adapted from github.com/jfreymuth/pulse (MIT License, Copyright Johann Freymuth)
package pulse

import (
	"fmt"
)

// Command opcodes sent as the first u32 in every tagstruct payload
// Only the commands needed for connect and playback are defined
const (
	// Sent by client
	CommandAuth                = 0x000B
	CommandSetClientName       = 0x0013
	CommandCreatePlaybackStream = 0x0003

	// Sent by server
	CommandReply   = 0x0000 // successful reply to a command
	CommandError   = 0x0001 // error reply
	CommandRequest = 0x0002 // server asking for more audio data
)

// Auth protocol version 28
const ProtoVersion = 0x1C

// Undefined is a placeholder for default option
const Undefined = 0xFFFFFFFF

// Error codes returned in CommandError replies
const (
	ErrAccess        = 1
	ErrCommand       = 2
	ErrInvalid       = 3
	ErrExist         = 4
	ErrNoEntity      = 5
	ErrConnectionRefused = 6
	ErrProtocol      = 7
	ErrTimeout       = 8
	ErrAuthKey       = 9
	ErrInternal      = 10
	ErrConnectionTerminated = 11
	ErrKilled        = 12
	ErrInvalidServer = 13
	ErrNotSupported  = 14
	ErrUnknown       = 15
	ErrNoExtension   = 16
	ErrObsolete      = 17
	ErrNotImplemented = 18
	ErrForked        = 19
	ErrIO            = 20
	ErrBusy          = 21
)

// BuildAuth builds the AUTH command payload
// cookie is the 256-byte contents of ~/.config/pulse/cookie
// seqTag is the sequence counter, start at 1 and increment each command
func BuildAuth(seqTag uint32, cookie []byte) []byte {
	w := &TagWriter{}
	w.WriteU32(CommandAuth)
	w.WriteU32(seqTag)
	w.WriteU32(ProtoVersion)
	w.WriteArbitrary(cookie)
	return w.Bytes()
}

// BuildSetClientName builds the SET_CLIENT_NAME command payload
// name is shown in pactl output to identify your application
func BuildSetClientName(seqTag uint32, name string) []byte {
	w := &TagWriter{}
	w.WriteU32(CommandSetClientName)
	w.WriteU32(seqTag)
	// PA expects a proplist: key/value pairs terminated by empty string
	w.WriteString("application.name")
	w.WriteString(name)
	w.WriteString("") // terminate proplist
	return w.Bytes()
}

// PlaybackStreamConfig holds the parameters for a new playback stream
type PlaybackStreamConfig struct {
	SampleSpec  SampleSpec
	SinkName    string // empty = default sink
	BufferSize  uint32 // bytes, 0xFFFFFFFF = server decides
	SeqTag      uint32
}

// BuildCreatePlaybackStream builds the CREATE_PLAYBACK_STREAM command payload
func BuildCreatePlaybackStream(cfg PlaybackStreamConfig) []byte {
	w := &TagWriter{}
	w.WriteU32(CommandCreatePlaybackStream)
	w.WriteU32(cfg.SeqTag)
	w.WriteSampleSpec(cfg.SampleSpec)

	// channel map: let server decide (write Undefined channels)
	// for stereo write a basic left/right map
	w.WriteU8(cfg.SampleSpec.Channels)
	if cfg.SampleSpec.Channels == 2 {
		w.WriteU8(1) // left
		w.WriteU8(2) // right
	} else {
		for i := uint8(0); i < cfg.SampleSpec.Channels; i++ {
			w.WriteU8(8) // mono
		}
	}

	// sink index (Undefined uses default or matches by name)
	w.WriteU32(Undefined)
	// sink name (empty uses default)
	w.WriteString(cfg.SinkName)

	// buffer size in bytes (Undefined lets server decides)
	bufSize := cfg.BufferSize
	if bufSize == 0 {
		bufSize = Undefined
	}
	w.WriteU32(bufSize)

	// corked: false (start playing immediately)
	w.WriteBool(false)

	// tlength, prebuf, minreq (Undefined, let server decide)
	w.WriteU32(Undefined) // tlength
	w.WriteU32(Undefined) // prebuf
	w.WriteU32(Undefined) // minreq

	// sync group id
	w.WriteU32(Undefined)

	// volume: write PA_VOLUME_NORM to each channel
	w.WriteU32(uint32(cfg.SampleSpec.Channels))
	for i := uint8(0); i < cfg.SampleSpec.Channels; i++ {
		w.WriteU32(0x10000) // PA_VOLUME_NORM = 0x10000
	}

	return w.Bytes()
}

// ParseReply reads the command and sequence tag from a server reply payload
// Returns the command (CommandReply or CommandError), the sequence tag,
// and a reader positioned after those two fields for further parsing
func ParseReply(payload []byte) (cmd uint32, seqTag uint32, r *TagReader, err error) {
	r = NewTagReader(payload)
	cmd, err = r.ReadU32()
	if err != nil {
		return
	}
	seqTag, err = r.ReadU32()
	return
}

// ParseRequest reads a CommandRequest payload.
// Returns the stream index and the number of bytes the server wants.
func ParseRequest(payload []byte) (streamIndex uint32, nbytes uint32, err error) {
	r := NewTagReader(payload)
	cmd, err := r.ReadU32()
	if err != nil {
		return
	}
	if cmd != CommandRequest {
		err = fmt.Errorf("proto: expected REQUEST command, got %d", cmd)
		return
	}
	// REQUEST has no seqTag — just index and nbytes
	streamIndex, err = r.ReadU32()
	if err != nil {
		return
	}
	nbytes, err = r.ReadU32()
	return
}
