//go:build linux

package pulse

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/ImOutOfIdeas/audigo/internal"
)

type backend struct {
	conn net.Conn
	seq  uint32
}

func New() (internal.Backend, error) {
	b := &backend{seq: 0}

	println(socket_path())

	socket, err := net.Dial("unix", socket_path())
	if err != nil {
		return nil, fmt.Errorf("pulse: connect: %w", err)
	}
	b.conn = socket

	if err := b.handshake(); err != nil {
		b.conn.Close()
		return nil, err
	}

	return b, nil
}

func (b *backend) Close() error {
	return b.conn.Close()
}

func (b *backend) OpenStream(cfg internal.StreamConfig) (internal.Stream, error) {
	return nil, fmt.Errorf("pulse: OpenStream not yet implemented")
}

func (b *backend) handshake() error {
	cookie, err := read_cookie()
	if err != nil {
		return err
	}

	fmt.Printf("cookie length: %d\n", len(cookie))

	pl := BuildAuth(b.seq, cookie)
	fmt.Printf("auth payload (%d bytes): %x\n", len(pl), pl)

	if err := WritePacket(b.conn, ControlChannel, BuildAuth(b.seq, cookie)); err != nil {
		return fmt.Errorf("pulse: send auth: %w", err)
	}
	b.seq++ // Inc after every command

	_, payload, err := ReadPacket(b.conn)
	if err != nil {
		return fmt.Errorf("pulse: recv auth reply: %w", err)
	}

	if err := expect_reply(payload); err != nil {
		return fmt.Errorf("pulse: auth: %w", err)
	}

	if err := WritePacket(b.conn, ControlChannel, BuildSetClientName(b.seq, "audigo")); err != nil {
		return fmt.Errorf("pulse: send client name: %w", err)
	}
	b.seq++ // Inc after every command

	_, payload, err = ReadPacket(b.conn)
	if err != nil {
		return fmt.Errorf("pulse: recv client name reply: %w", err)
	}

	if err := expect_reply(payload); err != nil {
		return fmt.Errorf("pulse: set client name: %w", err)
	}

	return nil
}

func socket_path() string {
	if p := os.Getenv("PULSE_SERVER"); p != "" {
		// TODO: handle full server string format properly
		if len(p) > 5 && p[:5] == "unix:" {
			return p[5:]
		}
		return p
	}
	return fmt.Sprintf("/run/user/%d/pulse/native", os.Getuid())
}

func read_cookie() ([]byte, error) {
	if p := os.Getenv("PULSE_COOKIE"); p != "" {
		return os.ReadFile(p)
	}
	config_dir := os.Getenv("XDG_CONFIG_HOME")
	if config_dir == "" {
		config_dir = os.Getenv("HOME") + "/.config"
	}
	data, err := os.ReadFile(config_dir + "/pulse/cookie")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// servers with auth-anonymous=1 accept any 256 byte cookie
			return make([]byte, 256), nil
		}
		return nil, fmt.Errorf("pulse: read cookie: %w", err)
	}
	return data, nil
}

func expect_reply(payload []byte) error {
	cmd, _, r, err := ParseReply(payload)
	if err != nil {
		return err
	}
	if cmd == CommandError {
		code, _ := r.ReadU32()
		return fmt.Errorf("server error code %d", code)
	}
	if cmd != CommandReply {
		return fmt.Errorf("unexpected command %d", cmd)
	}
	return nil
}
