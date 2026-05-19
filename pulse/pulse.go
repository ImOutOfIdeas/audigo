package pulse

import (
    "fmt"
    "net"
    "os"

    "github.com/imoutofideas/audigo"
)

type backend struct {
    conn net.Conn
    seq  uint32
}

func New() (audigo.Backend, error) {
    conn, err := net.Dial("unix", socket_path())
    if err != nil {
        return nil, fmt.Errorf("pulse: connect: %w", err)
    }
    b := &backend{conn: conn}
    if err := b.handshake(); err != nil {
        conn.Close()
        return nil, err
    }
    return b, nil
}

func (b *backend) Close() error {
    return b.conn.Close()
}

func (b *backend) OpenStream(cfg audigo.StreamConfig) (audigo.Stream, error) {
    return nil, fmt.Errorf("pulse: OpenStream not yet implemented")
}

func (b *backend) handshake() error {
    cookie, err := read_cookie()
    if err != nil {
        return err
    }

    b.seq++
    if err := WritePacket(b.conn, ControlChannel, BuildAuth(b.seq, cookie)); err != nil {
        return fmt.Errorf("pulse: send auth: %w", err)
    }
    _, payload, err := ReadPacket(b.conn)
    if err != nil {
        return fmt.Errorf("pulse: recv auth reply: %w", err)
    }
    if err := expect_reply(payload); err != nil {
        return fmt.Errorf("pulse: auth: %w", err)
    }

    b.seq++
    if err := WritePacket(b.conn, ControlChannel, BuildSetClientName(b.seq, "audigo")); err != nil {
        return fmt.Errorf("pulse: send client name: %w", err)
    }
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
    return os.ReadFile(config_dir + "/pulse/cookie")
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
