package pulse

import (
	"net"
)

type Client struct {
	connection net.Conn

}

// Attempt connection to each serverInfo struct
// Upon success get cookie and send auth request
// Recieve and validate AuthReply
// Profit!
func connect() {

}
