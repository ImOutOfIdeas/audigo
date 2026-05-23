package pulse

import (
	"io"
	"encoding/binary"
	"fmt"
	"net"
)

type packetReader struct {
	buf []byte
	pos int
	err error
}

// Constructor
func newPacketReader(data []byte) *packetReader {
	return &packetReader{buf: data}
}

// Bytes remaining unread in reader buffer
func (r *packetReader) remaining() int {
	return len(r.buf) - r.pos
}

// reads the 20-byte frame header and returns the payload
func read_packet(conn net.Conn) ([]byte, error) {
	header := make([]byte, 20)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(header[0:4])
	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// reads opcode and tag from a command payload
func read_command(payload []byte) (opcode uint32, tag uint32, r *packetReader) {
	r = newPacketReader(payload)
	r.byte()        // 'L' tag for opcode
	opcode = r.uint32()
	r.byte()        // 'L' tag for tag
	tag = r.uint32()
	return opcode, tag, r
}

// === Type Specific Read Methods ===

func (r *packetReader) byte() byte {
	if r.err != nil || r.remaining() < 1 {
		r.err = io.ErrUnexpectedEOF
		return 0
	}
	b := r.buf[r.pos]
	r.pos++
	return b
}

func (r *packetReader) uint32() uint32 {
	if r.err != nil || r.remaining() < 4 {
		r.err = io.ErrUnexpectedEOF
		return 0
	}
	v := binary.BigEndian.Uint32(r.buf[r.pos:])
	r.pos += 4
	return v
}

func (r *packetReader) uint64() uint64 {
	if r.err != nil || r.remaining() < 8 {
		r.err = io.ErrUnexpectedEOF
		return 0
	}
	v := binary.BigEndian.Uint64(r.buf[r.pos:])
	r.pos += 8
	return v
}

func (r *packetReader) string() string {
	if r.err != nil {
		return ""
	}
	start := r.pos
	for r.pos < len(r.buf) {
		if r.buf[r.pos] == 0 {
			s := string(r.buf[start:r.pos])
			r.pos++ // skip null terminator
			return s
		}
		r.pos++
	}
	r.err = io.ErrUnexpectedEOF
	return ""
}

func (r *packetReader) arbitrary() []byte {
	if r.err != nil {
		return nil
	}
	length := r.uint32()
	if r.remaining() < int(length) {
		r.err = io.ErrUnexpectedEOF
		return nil
	}
	v := r.buf[r.pos : r.pos+int(length)]
	r.pos += int(length)
	return v
}

// reads the next tagged value and returns it as any
func (r *packetReader) next() any {
	if r.err != nil {
		return nil
	}
	tag := r.byte()
	switch tag {
	case 'L':
		return r.uint32()
	case 'B':
		return r.byte()
	case 'R':
		return r.uint64()
	case 'r':
		return int64(r.uint64())
	case 't':
		return r.string()
	case 'N':
		return nil
	case 'x':
		return r.arbitrary()
	case '1':
		return true
	case '0':
		return false
	case 'a': // sample spec
		format := r.byte()
		channels := r.byte()
		rate := r.uint32()
		return [3]uint32{uint32(format), uint32(channels), rate}
	case 'm': // channel map
		count := r.byte()
		positions := make([]byte, count)
		for i := range positions {
			positions[i] = r.byte()
		}
		return positions
	case 'v': // cvolume
		count := r.byte()
		volumes := make([]uint32, count)
		for i := range volumes {
			volumes[i] = r.uint32()
		}
		return volumes
	case 'U':
		return r.uint64()
	case 'T':
		seconds := r.uint32()
		microseconds := r.uint32()
		return [2]uint32{seconds, microseconds}
	default:
		r.err = fmt.Errorf("unknown tag: %c (0x%02x)", tag, tag)
		return nil
	}
}
