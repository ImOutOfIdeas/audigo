package pulse

// Misc
const (
    cookie_length    = 256
	header_length    = 20

    protocol_version = uint32(35)
    control_channel  = uint32(0xFFFFFFFF)
)

// Datatype tags
const (
    tag_arbitrary  = byte('x')
    tag_bool_false = byte('0')
    tag_bool_true  = byte('1')
    tag_string     = byte('t')
    tag_u8         = byte('B')
	tag_u16		   = byte('S')
    tag_u32        = byte('L')
)

// OP codes (commanly used subset)
const (
    command_error                  = uint32(0)
    command_reply                  = uint32(1)
    command_auth                   = uint32(4)
    command_set_client_name        = uint32(5)
    command_create_playback_stream = uint32(9)
    command_delete_playback_stream = uint32(11)
    command_create_record_stream   = uint32(10)
    command_delete_record_stream   = uint32(12)
    command_get_sink_info          = uint32(19)
    command_get_sink_info_list     = uint32(20)
    command_get_source_info        = uint32(22)
    command_get_source_info_list   = uint32(23)
    command_get_server_info        = uint32(26)
    command_subscribe              = uint32(28)
    command_set_sink_volume        = uint32(35)
    command_set_sink_input_volume  = uint32(36)
    command_set_sink_mute          = uint32(38)
    command_cork_playback_stream   = uint32(43)
    command_flush_playback_stream  = uint32(44)
    command_request                = uint32(59)
)
