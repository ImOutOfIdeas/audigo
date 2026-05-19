// Package proto implements the PulseAudio native binary protocol.
// Adapted from github.com/jfreymuth/pulse (MIT License, Copyright Johann Freymuth)
package pulse

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ControlChannel is used for all command packets.
// Audio data packets use the stream's assigned channel number instead.
const ControlChannel = 0xFFFFFFFF

// packet header is always 20 bytes:
//   [0:4]  payload length
//   [4:8]  channel (0xFFFFFFFF = control)
//   [8:12] offset high (unused, always 0)
//  [12:16] offset low  (unused, always 0)
//  [16:20] flags       (unused, always 0)
const headerSize = 20

// WritePacket writes a framed packet to w.
// channel should be ControlChannel for commands, or the stream's
// data channel number for PCM audio.
func WritePacket(w io.Writer, channel uint32, payload []byte) error {
	hdr := make([]byte, headerSize)
	binary.BigEndian.PutUint32(hdr[0:], uint32(len(payload)))
	binary.BigEndian.PutUint32(hdr[4:], channel)
	// offset hi, offset lo, flags — all zero

	if _, err := w.Write(hdr); err != nil {
		return fmt.Errorf("proto: write header: %w", err)
	}
	if _, err := w.Write(payload); err != nil {
		return fmt.Errorf("proto: write payload: %w", err)
	}
	return nil
}

// ReadPacket reads one framed packet from r.
// Returns the channel number and raw payload bytes.
func ReadPacket(r io.Reader) (channel uint32, payload []byte, err error) {
	hdr := make([]byte, headerSize)
	if _, err = io.ReadFull(r, hdr); err != nil {
		return 0, nil, fmt.Errorf("proto: read header: %w", err)
	}

	length := binary.BigEndian.Uint32(hdr[0:])
	channel = binary.BigEndian.Uint32(hdr[4:])

	payload = make([]byte, length)
	if _, err = io.ReadFull(r, payload); err != nil {
		return 0, nil, fmt.Errorf("proto: read payload: %w", err)
	}
	return channel, payload, nil
}
