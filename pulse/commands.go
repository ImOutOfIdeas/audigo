/*
=== PulseAudio Server Packet Format ===

 [4 bytes] length of entire payload
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

/*
TODO:
	implement a command packet constructor with variadic
	argument for serialized command argument fields

	make_command(opcode, tag?, arguments...)

*/

package pulse

import (
	"encoding/binary"
	"io"
	"net"
	"bytes"
	"fmt"
)

func make_auth_packet(cookie []byte) []byte {
    buf := new(bytes.Buffer)

    // frame header
    binary.Write(buf, binary.BigEndian, uint32(0)) // length placeholder
    binary.Write(buf, binary.BigEndian, uint32(0xFFFFFFFF)) // control channel
    binary.Write(buf, binary.BigEndian, uint32(0))
    binary.Write(buf, binary.BigEndian, uint32(0))
    binary.Write(buf, binary.BigEndian, uint32(0))

    payload_start := buf.Len()

    // PA_COMMAND_AUTH
    binary.Write(buf, binary.BigEndian, byte(tag_u32))
    binary.Write(buf, binary.BigEndian, uint32(command_auth))

    // tag sequence
    binary.Write(buf, binary.BigEndian, byte(tag_u32))
    binary.Write(buf, binary.BigEndian, uint32(1))

    // protocol version, no SHM
    binary.Write(buf, binary.BigEndian, byte(tag_u32))
    binary.Write(buf, binary.BigEndian, uint32(protocol_version))

    // 256-byte zeroed cookie
    binary.Write(buf, binary.BigEndian, byte(tag_arbitrary))
    binary.Write(buf, binary.BigEndian, uint32(cookie_length))
    buf.Write(cookie)

    // fill in payload length
    packet := buf.Bytes()
    binary.BigEndian.PutUint32(packet[0:4], uint32(buf.Len()-payload_start))

    return packet
}

func read_auth_reply(conn net.Conn) (uint32, error) {
    header := make([]byte, 20)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return 0, err
	}

    length := binary.BigEndian.Uint32(header[0:4])
    payload := make([]byte, length)
	_, err = io.ReadFull(conn, payload)
	if err != nil {
		return 0, err
	}

    // first field is the command: 0 = PA_COMMAND_REPLY, 1 = PA_COMMAND_ERROR
    command := binary.BigEndian.Uint32(payload[1:5])
    if command != 0 {
        return 0, fmt.Errorf("auth rejected, command: %d", command)
    }

	server_version := binary.BigEndian.Uint32(payload[6:10]) & 0x0000FFFF // mask off SHM/memfd flags
    //if server_version < b.protocol_version {
        //b.protocol_version = server_version
    //}

    return server_version, nil
}
