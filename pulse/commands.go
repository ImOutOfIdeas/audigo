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

package pulse
