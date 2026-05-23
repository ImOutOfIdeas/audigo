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
  	=== serialized command arguments ===
 [1 byte]  byte marker (specifies type)
 [n bytes] associated data
 		   ...
*/

package pulse

import (
	"encoding/binary"
	"bytes"
)

type packetWriter struct {
	buf bytes.Buffer
}

// Constructor
func newPacketWriter() *packetWriter {
	return &packetWriter{}
}

// === Type Specific Write Methods ===

func (w *packetWriter) uint32(v uint32) *packetWriter {
	w.buf.WriteByte(tag_u32)
	binary.Write(&w.buf, binary.BigEndian, v)
	return w
}

func (w *packetWriter) byte(v byte) *packetWriter {
	w.buf.WriteByte(tag_u8)
	w.buf.WriteByte(v)
	return w
}

func (w *packetWriter) uint64(v uint64) *packetWriter {
	w.buf.WriteByte(tag_u64)
	binary.Write(&w.buf, binary.BigEndian, v)
	return w
}

func (w *packetWriter) int64(v int64) *packetWriter {
	w.buf.WriteByte(tag_s64)
	binary.Write(&w.buf, binary.BigEndian, uint64(v))
	return w
}

func (w *packetWriter) string(v string) *packetWriter {
	w.buf.WriteByte(tag_string)
	w.buf.WriteString(v)
	w.buf.WriteByte(0)
	return w
}

func (w *packetWriter) string_null() *packetWriter {
	w.buf.WriteByte(tag_string_null)
	return w
}

func (w *packetWriter) bool(v bool) *packetWriter {
	if v {
		w.buf.WriteByte(tag_bool_true)
	} else {
		w.buf.WriteByte(tag_bool_false)
	}
	return w
}

func (w *packetWriter) arbitrary(v []byte) *packetWriter {
	w.buf.WriteByte(tag_arbitrary)
	binary.Write(&w.buf, binary.BigEndian, uint32(len(v)))
	w.buf.Write(v)
	return w
}

func (w *packetWriter) sample_spec(format byte, rate uint32, channels byte) *packetWriter {
	w.buf.WriteByte(tag_sample_spec)
	w.buf.WriteByte(format)
	w.buf.WriteByte(channels)
	binary.Write(&w.buf, binary.BigEndian, rate)
	return w
}

func (w *packetWriter) channel_map(positions ...byte) *packetWriter {
	w.buf.WriteByte(tag_channel_map)
	w.buf.WriteByte(byte(len(positions)))
	for _, p := range positions {
		w.buf.WriteByte(p)
	}
	return w
}

func (w *packetWriter) cvolume(volumes ...uint32) *packetWriter {
	w.buf.WriteByte(tag_cvolume)
	w.buf.WriteByte(byte(len(volumes)))
	for _, v := range volumes {
		binary.Write(&w.buf, binary.BigEndian, v)
	}
	return w
}

func (w *packetWriter) prop(key string, value []byte) *packetWriter {
	w.string(key)
	w.uint32(uint32(len(value)))
	w.arbitrary(value)
	return w
}

func (w *packetWriter) prop_list_begin() *packetWriter {
    w.buf.WriteByte(tag_proplist)
    return w
}

func (w *packetWriter) prop_list_end() *packetWriter {
    w.buf.WriteByte(tag_string_null)
    return w
}

// Frame the payload
func (w *packetWriter) frame(opcode uint32, tag uint32) []byte {
	// Construct payload header
	payload := new(bytes.Buffer)
	payload.WriteByte(tag_u32)
	binary.Write(payload, binary.BigEndian, opcode)
	payload.WriteByte(tag_u32)
	binary.Write(payload, binary.BigEndian, tag)
	payload.Write(w.buf.Bytes())

	// Frame the payload within a packet header
	packet := new(bytes.Buffer)
	binary.Write(packet, binary.BigEndian, uint32(payload.Len()))
	binary.Write(packet, binary.BigEndian, control_channel)
	binary.Write(packet, binary.BigEndian, uint32(0))
	binary.Write(packet, binary.BigEndian, uint32(0))
	binary.Write(packet, binary.BigEndian, uint32(0))

	// Append the actual payload to the packet
	packet.Write(payload.Bytes())

	return packet.Bytes()
}
