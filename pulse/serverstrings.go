// Locates and parses pulse server strings. Populates serverInfo structs for use
// with connection attempts via net.Dial

package pulse

import (
	"errors"
	"net"
	"os"
	"path"
	"strings"
)

type serverInfo struct {
	protocol  string
	address   string
	localname string
}

// Default TCP port for hostname only command string
const defaultPulseAudioTCPPort = "4713"

// Return a list of servers([]serverInfo) parsed from a server string
// These are located like so:
//  - user provided serverString argument
//  - if empty, PULSE_SERVER environment variable
//  - if empty or isn't set, the linux default (unix:[XDG_RUNTIME_DIR]/pulse/native)
func getServers(serverString string) ([]serverInfo, error) {
	var servers []serverInfo
	// Attempt to locate a server string
	if serverString != "" {
		servers = parseServerString(serverString)
	} else if serverEnv, ok := os.LookupEnv("PULSE_SERVER"); ok && serverEnv != "" {
		servers = parseServerString(serverEnv)
	} else {
		servers = defaultServerString()
	}
	if len(servers) == 0 {
		return []serverInfo{}, errors.New("pulseaudio: no valid server")
	}

	return servers, nil
}

// Get the default server string if one can't be found otherwise
func defaultServerString() []serverInfo {
	return []serverInfo{{
		protocol: "unix",
		address:  path.Join(os.Getenv("XDG_RUNTIME_DIR"), "pulse/native"),
	}}
}

// Get all servers from a server string and return a slice of serverInfo's
func parseServerString(str string) []serverInfo {
	s := strings.Fields(str)
	var result []serverInfo
	for _, s := range s {
		server, ok := parseOneServerString(s)
		if !ok {
			continue
		}
		result = append(result, server)
	}
	return result
}

// parse a single server from the server string
func parseOneServerString(str string) (serverInfo, bool) {
	var server serverInfo
	if str[0] == '{' {
		end := strings.IndexByte(str, '}')
		server.localname = str[1:end]
		str = str[end+1:]
	}
	switch {
	case len(str) == 0:
		// no server string
		return serverInfo{}, false
	case str[0] == '/':
		// rule #2
		server.protocol = "unix"
		server.address = str
	case strings.HasPrefix(str, "unix:"):
		// rule #2
		server.protocol = "unix"
		server.address = str[5:]
	case strings.HasPrefix(str, "tcp6:"):
		// rule #4
		server.protocol = "tcp6"
		server.address = str[5:]
	case strings.HasPrefix(str, "tcp4:"):
		// rule #3
		server.protocol = "tcp4"
		server.address = str[5:]
	case strings.HasPrefix(str, "tcp:"):
		// rule #3
		server.protocol = "tcp"
		server.address = str[4:]
	default:
		// rule #5
		if _, _, err := net.SplitHostPort(str); err == nil {
			server.protocol = "tcp"
			server.address = str
		} else {
			// Adding a default port, as only providing a hostname is valid
			server.protocol = "tcp"
			server.address = net.JoinHostPort(str, defaultPulseAudioTCPPort)
		}
	}
	return server, true
}
