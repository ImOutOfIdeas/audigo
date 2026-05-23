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
    tag_string       = byte('t')
    tag_string_null  = byte('N')
    tag_u32          = byte('L')
    tag_u8           = byte('B')
    tag_u64          = byte('R')
    tag_s64          = byte('r')
    tag_sample_spec  = byte('a')
    tag_arbitrary    = byte('x')
    tag_bool_true    = byte('1')
    tag_bool_false   = byte('0')
    tag_timeval      = byte('T')
    tag_usec         = byte('U')
    tag_channel_map  = byte('m')
    tag_cvolume      = byte('v')
    tag_proplist     = byte('P')
    tag_volume       = byte('V')
    tag_format_info  = byte('f')
)

// OP codes (useful subset)
const (
    op_error   = uint32(0)
    op_timeout = uint32(1)
    op_reply   = uint32(2)

    op_create_playback_stream = uint32(3)
    op_delete_playback_stream = uint32(4)
    op_create_record_stream   = uint32(5)
    op_delete_record_stream   = uint32(6)

    op_exit            = uint32(7)
    op_auth            = uint32(8)
    op_set_client_name = uint32(9)

	op_get_server_info        = uint32(20)
	op_get_sink_info          = uint32(21)
	op_get_sink_info_list     = uint32(22)
	op_get_source_info        = uint32(23)
	op_get_source_info_list   = uint32(24)

	op_subscribe              = uint32(35)
	op_set_sink_volume        = uint32(36)
	op_set_sink_input_volume  = uint32(37)
	op_set_sink_mute          = uint32(39)

	op_cork_playback_stream   = uint32(41)
	op_flush_playback_stream  = uint32(42)

	op_request                = uint32(61)
)
