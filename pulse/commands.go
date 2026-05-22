/*
=== PulseAudio Server Packet Format ===

------------- Header -------------
 [4 bytes] length of payload
 [4 bytes] index (0xFFFFFFFF for command messages)
 [8 bytes] offset (0 for commands)
 [4 bytes] flags (0 for commands)

 ------------ Payload ------------
 [1 byte]  'L' (byte marker)
 [4 bytes] command opcode (big-endian uint32)
 [1 byte]  'L' (byte marker)
 [4 bytes] tag (big-endian uint32, used to match replies)
 [...    ] serialized command arguments
*/

package pulse

import (
	"encoding/binary"
	"io"
	"net"
	"bytes"
	"fmt"
)

/*
TODO:
	implement a command packet constructor with variadic
	argument for serialized command argument fields

	make_command(opcode, tag?, arguments...)

*/

// Creates a serialized command packet
func make_command(opcode uint32, tag uint32, arguments ...any) []byte {
	payloadBuf := new(bytes.Buffer)

	// Write Opcode
	payloadBuf.WriteByte('L')
	binary.Write(payloadBuf, binary.BigEndian, opcode)

	// Write Tag
	payloadBuf.WriteByte('L')
	binary.Write(payloadBuf, binary.BigEndian, tag)

	// Write each tag/argument pair into the payload
	for _, arg := range arguments {
		switch v := arg.(type) {
		case uint32:
			payloadBuf.WriteByte(tag_u32)
			binary.Write(payloadBuf, binary.BigEndian, v)
		case uint16:
			payloadBuf.WriteByte(tag_u16)
			binary.Write(payloadBuf, binary.BigEndian, v)
		case byte:
			payloadBuf.WriteByte('B')
			payloadBuf.WriteByte(v)
		case []byte:
			payloadBuf.WriteByte('x')
			binary.Write(payloadBuf, binary.BigEndian, uint32(len(v)))
			payloadBuf.Write(v)
		case string:
			payloadBuf.WriteByte('Z')
			payloadBuf.WriteString(v)
			payloadBuf.WriteByte(0x00) // Null-terminate string
		default:
			// Unrecognized type: skip or handle error
			continue
		}
	}

	return make_command_packet(payloadBuf.Bytes())
}

// Wraps a payload in a command packet header
func make_command_packet(payload []byte) []byte {
	payload_size := len(payload)

	packetBuf := new(bytes.Buffer)
	packetBuf.Grow(header_length + payload_size)

	// Packet header
    binary.Write(packetBuf, binary.BigEndian, uint32(payload_size)) // length
    binary.Write(packetBuf, binary.BigEndian, control_channel)  	// channel
    binary.Write(packetBuf, binary.BigEndian, uint32(0)) 		    // Offset
    binary.Write(packetBuf, binary.BigEndian, uint32(0))			// Offset
    binary.Write(packetBuf, binary.BigEndian, uint32(0))			// Flags

	// Payload
	packetBuf.Write(payload)

	return packetBuf.Bytes()
}

func read_auth_reply(conn net.Conn) error {
	// Read header to extract payload length
	header := make([]byte, 20)
	if _, err := io.ReadFull(conn, header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}
	length := binary.BigEndian.Uint32(header[0:4])

	// Read the payload
	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return fmt.Errorf("failed to read payload: %w", err)
	}

	// Parse feild by feild with a bytes.Reader
	r := bytes.NewReader(payload)

	// Read Opcode
	if marker, _ := r.ReadByte(); marker != 'L' {
		return fmt.Errorf("invalid opcode marker")
	}
	var command uint32
	binary.Read(r, binary.BigEndian, &command)

	// Read Tag
	if marker, _ := r.ReadByte(); marker != 'L' {
		return fmt.Errorf("invalid tag marker")
	}
	var tag uint32
	binary.Read(r, binary.BigEndian, &tag)

	// Validate results
	if command != command_reply {
		return fmt.Errorf("auth rejected by server (opcode: %d)", command)
	}

	return nil
}
