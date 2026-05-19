// Adapted from github.com/jfreymuth/pulse (MIT License, Copyright Johann Freymuth)
package pulse

import (
	"encoding/binary"
	"fmt"
)

// Tag bytes prefix every value in a tagstruct payload.
// Each value is: [1 byte tag][value bytes]
const (
	tagU32        = 'L' // uint32, 4 bytes big-endian
	tagU8         = 'B' // uint8,  1 byte
	tagString     = 't' // null-terminated UTF-8 string
	tagStringNull = 'N' // null string (empty/unset)
	tagArbitrary  = 'x' // raw bytes: [4 byte length][data]
	tagBoolean    = 'b' // boolean: ascii '1' true, '0' false (PA uses different — see below)
	tagBoolTrue   = '1'
	tagBoolFalse  = '0'
	tagSampleSpec = 'a' // format(u8) + rate(u32) + channels(u8)
	tagUsec       = 'U' // microseconds, uint64 8 bytes big-endian
)

// SampleFormat constants match PulseAudio's pa_sample_format_t
const (
	SampleU8        uint8 = 0
	SampleALAW      uint8 = 1
	SampleULAW      uint8 = 2
	SampleS16LE     uint8 = 3
	SampleS16BE     uint8 = 4
	SampleFloat32LE uint8 = 5
	SampleFloat32BE uint8 = 6
	SampleS32LE     uint8 = 7
	SampleS32BE     uint8 = 8
	SampleS24LE     uint8 = 9
)

// SampleSpec describes a PCM audio format.
type SampleSpec struct {
	Format   uint8
	Rate     uint32
	Channels uint8
}

// --- Writer ---

// TagWriter builds a tagstruct payload.
// Call bytes() to get the finished payload to pass to WritePacket.
type TagWriter struct {
	buf []byte
}

func (w *TagWriter) WriteU32(v uint32) {
	w.buf = append(w.buf, tagU32)
	w.buf = binary.BigEndian.AppendUint32(w.buf, v)
}

func (w *TagWriter) WriteU8(v uint8) {
	w.buf = append(w.buf, tagU8, v)
}

func (w *TagWriter) WriteString(s string) {
	if s == "" {
		w.buf = append(w.buf, tagStringNull)
		return
	}
	w.buf = append(w.buf, tagString)
	w.buf = append(w.buf, []byte(s)...)
	w.buf = append(w.buf, 0) // null terminator
}

func (w *TagWriter) WriteArbitrary(b []byte) {
	w.buf = append(w.buf, tagArbitrary)
	w.buf = binary.BigEndian.AppendUint32(w.buf, uint32(len(b)))
	w.buf = append(w.buf, b...)
}

func (w *TagWriter) WriteBool(v bool) {
	if v {
		w.buf = append(w.buf, tagBoolTrue)
	} else {
		w.buf = append(w.buf, tagBoolFalse)
	}
}

func (w *TagWriter) WriteSampleSpec(s SampleSpec) {
	w.buf = append(w.buf, tagSampleSpec)
	w.buf = append(w.buf, s.Format)
	w.buf = binary.BigEndian.AppendUint32(w.buf, s.Rate)
	w.buf = append(w.buf, s.Channels)
}

func (w *TagWriter) WriteUsec(v uint64) {
	w.buf = append(w.buf, tagUsec)
	w.buf = binary.BigEndian.AppendUint64(w.buf, v)
}

// Bytes returns the completed tagstruct payload.
func (w *TagWriter) Bytes() []byte {
	return w.buf
}

// --- Reader ---

// TagReader reads typed values from a tagstruct payload.
type TagReader struct {
	buf []byte
	pos int
}

// NewTagReader wraps a raw payload for reading.
func NewTagReader(payload []byte) *TagReader {
	return &TagReader{buf: payload}
}

func (r *TagReader) peek() (byte, error) {
	if r.pos >= len(r.buf) {
		return 0, fmt.Errorf("proto: tagstruct: unexpected end of data")
	}
	return r.buf[r.pos], nil
}

func (r *TagReader) expectTag(expected byte) error {
	tag, err := r.peek()
	if err != nil {
		return err
	}
	if tag != expected {
		return fmt.Errorf("proto: tagstruct: expected tag %c (%d), got %c (%d)", expected, expected, tag, tag)
	}
	r.pos++
	return nil
}

func (r *TagReader) ReadU32() (uint32, error) {
	if err := r.expectTag(tagU32); err != nil {
		return 0, err
	}
	if r.pos+4 > len(r.buf) {
		return 0, fmt.Errorf("proto: tagstruct: not enough bytes for u32")
	}
	v := binary.BigEndian.Uint32(r.buf[r.pos:])
	r.pos += 4
	return v, nil
}

func (r *TagReader) ReadU8() (uint8, error) {
	if err := r.expectTag(tagU8); err != nil {
		return 0, err
	}
	if r.pos >= len(r.buf) {
		return 0, fmt.Errorf("proto: tagstruct: not enough bytes for u8")
	}
	v := r.buf[r.pos]
	r.pos++
	return v, nil
}

func (r *TagReader) ReadString() (string, error) {
	tag, err := r.peek()
	if err != nil {
		return "", err
	}
	if tag == tagStringNull {
		r.pos++
		return "", nil
	}
	if err := r.expectTag(tagString); err != nil {
		return "", err
	}
	start := r.pos
	for r.pos < len(r.buf) && r.buf[r.pos] != 0 {
		r.pos++
	}
	if r.pos >= len(r.buf) {
		return "", fmt.Errorf("proto: tagstruct: unterminated string")
	}
	s := string(r.buf[start:r.pos])
	r.pos++ // consume null terminator
	return s, nil
}

func (r *TagReader) ReadArbitrary() ([]byte, error) {
	if err := r.expectTag(tagArbitrary); err != nil {
		return nil, err
	}
	if r.pos+4 > len(r.buf) {
		return nil, fmt.Errorf("proto: tagstruct: not enough bytes for arb length")
	}
	length := binary.BigEndian.Uint32(r.buf[r.pos:])
	r.pos += 4
	if r.pos+int(length) > len(r.buf) {
		return nil, fmt.Errorf("proto: tagstruct: not enough bytes for arb data")
	}
	v := make([]byte, length)
	copy(v, r.buf[r.pos:])
	r.pos += int(length)
	return v, nil
}

func (r *TagReader) ReadBool() (bool, error) {
	tag, err := r.peek()
	if err != nil {
		return false, err
	}
	r.pos++
	switch tag {
	case tagBoolTrue:
		return true, nil
	case tagBoolFalse:
		return false, nil
	default:
		return false, fmt.Errorf("proto: tagstruct: expected bool tag, got %c (%d)", tag, tag)
	}
}

func (r *TagReader) ReadSampleSpec() (SampleSpec, error) {
	if err := r.expectTag(tagSampleSpec); err != nil {
		return SampleSpec{}, err
	}
	if r.pos+6 > len(r.buf) {
		return SampleSpec{}, fmt.Errorf("proto: tagstruct: not enough bytes for sample spec")
	}
	spec := SampleSpec{
		Format:   r.buf[r.pos],
		Rate:     binary.BigEndian.Uint32(r.buf[r.pos+1:]),
		Channels: r.buf[r.pos+5],
	}
	r.pos += 6
	return spec, nil
}

func (r *TagReader) ReadUsec() (uint64, error) {
	if err := r.expectTag(tagUsec); err != nil {
		return 0, err
	}
	if r.pos+8 > len(r.buf) {
		return 0, fmt.Errorf("proto: tagstruct: not enough bytes for usec")
	}
	v := binary.BigEndian.Uint64(r.buf[r.pos:])
	r.pos += 8
	return v, nil
}

// Remaining returns how many bytes haven't been read yet.
// Useful for debugging unexpected trailing data.
func (r *TagReader) Remaining() int {
	return len(r.buf) - r.pos
}
