package pulse

import (
	"fmt"
	"net"
	"os"
)

func connect(b *pulseBackend) error {
	// Get list of servers to attempt connections
	servers, err := getServers("")
	if err != nil {
		return err
	}

	// Get hostname to filter servers
	localname, err := os.Hostname()
	if err != nil {
		return err
	}

	// Get cookie path
	cookiePath := os.Getenv("HOME") + "/.config/pulse/cookie"
	if path, ok := os.LookupEnv("PULSE_COOKIE"); ok {
		cookiePath = path
	}

	// Get cookie from path or empty if not found
	// Pulse may not need a cookie so an empty is fine
	cookie, err := os.ReadFile(cookiePath)
	if os.IsNotExist(err) {
		cookie = make([]byte, 256)
	} else if err != nil {
		return fmt.Errorf("Failed to read cookie file: %w", err)
	}

	var lastErr error
	for _, server := range servers {
		// Skip connections with different localnames
		if server.localname != "" && localname != server.localname {
			continue
		}

		// Attempt connection
		conn, err := net.Dial(server.protocol, server.address)
		if err != nil {
			lastErr = err
			continue
		}

		// Create the auth packet
		auth_packet := newPacketWriter().
			uint32(protocol_version).
			arbitrary(cookie).
			frame(op_auth, b.next())

		// Send the auth packet to the server
		_, err = conn.Write(auth_packet)
		if err != nil {
			conn.Close()
			lastErr = err
			continue
		}

		// Get auth reply from server and validate
		payload, err := read_packet(conn)
		if err != nil {
			lastErr = err
			continue
		}
		opcode, tagId, r := read_command(payload)
		if opcode != op_reply {
			lastErr = fmt.Errorf("Auth validation failed (code: %d)", opcode)
			continue
		}

		// Get server version from reply
		version := r.next().(uint32)
		if protocol_version != version {
			lastErr = fmt.Errorf("Protocol version mismatch. client: %d, server: %d\n", protocol_version, version)
		}

		// Create client name packet
		name_packet := newPacketWriter().
			prop_list_begin().
			prop("application.name", []byte("audigo")).
			prop_list_end().
			frame(op_set_client_name, b.next())

		fmt.Printf("name: %x\n", name_packet)

		// Send the client name packet to the server
		_, err = conn.Write(name_packet)
		if err != nil {
			conn.Close()
			lastErr = err
			continue
		}

		// Get client name server response
		payload, err = read_packet(conn)
		fmt.Printf("name res: %x\n", payload)
		if err != nil {
			lastErr = err
			continue
		}

		// Parse name response
		opcode, tagId, r = read_command(payload)
		fmt.Printf("opcode: %d, tagId: %d\n", opcode, tagId)
		if opcode != op_reply {
			errCode := r.next()
			lastErr = fmt.Errorf("Name setting failed (code: %v)", errCode)
			continue
		}
		// Get client index from response
		clientIndex := r.next().(uint32)

		// Set backend struct fields for the now open connection
		b.connection = conn
		b.version = version
		b.clientIndex = clientIndex

		return nil
	}

	return lastErr
}
