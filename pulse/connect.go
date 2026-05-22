package pulse

import (
	"fmt"
	"net"
	"os"
)

func connect() (net.Conn, uint32, error) {
	servers, err := getServers("")
	if err != nil {
		return nil, 0, err
	}

	localname, err := os.Hostname()
	if err != nil {
		return nil, 0, err
	}

	cookiePath := os.Getenv("HOME") + "/.config/pulse/cookie"
	if path, ok := os.LookupEnv("PULSE_COOKIE"); ok {
		cookiePath = path
	}

	cookie, err := os.ReadFile(cookiePath)
	if os.IsNotExist(err) {
		cookie = make([]byte, 256)
	} else if err != nil {
		return nil, 0, fmt.Errorf("Failed to read cookie file: %w", err)
	}

	auth_packet := make_auth_packet(cookie)

	var lastErr error
	for _, server := range servers {
		// Skip connections with different localnames
		if server.localname != "" && localname != server.localname {
			continue
		}

		conn, err := net.Dial(server.protocol, server.address)
		if err != nil {
			lastErr = err
			continue
		}

		_, err = conn.Write(auth_packet)
		if err != nil {
			lastErr = err
			continue
		}

		version, err := read_auth_reply(conn)
		if err != nil {
			lastErr = err
			continue
		}

		return conn, version, nil
	}

	return nil, 0, lastErr
}
