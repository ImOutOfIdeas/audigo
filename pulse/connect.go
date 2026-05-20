// TODO: Become intelligent and actually figure out what in the actual
// fuck is going on here


package pulse

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"sync"
	"time"
)

type pulseClient struct {
	conn      io.ReadWriter
	reader    *bufio.Reader

	// Single mutex protects all state
	mu        sync.Mutex
	nextTag   uint32
	pending   map[uint32]chan []byte  // tag -> response channel

	timeout   time.Duration
}

func NewSimpleClient(conn io.ReadWriter) *pulseClient {
	sc := &pulseClient{
		conn:    conn,
		reader:  bufio.NewReader(conn),
		pending: make(map[uint32]chan []byte),
		timeout: 5 * time.Second,
	}
	go sc.readLoop()
	return sc
}

// Send a command and wait for reply
func (sc *pulseClient) SendCommand(opcode uint32, payload []byte) ([]byte, error) {
	sc.mu.Lock()
	tag := sc.nextTag
	sc.nextTag++

	// Create response channel
	replyChan := make(chan []byte, 1)
	sc.pending[tag] = replyChan
	sc.mu.Unlock()

	// Build message: [opcode][tag][payload]
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, opcode)
	binary.Write(buf, binary.BigEndian, tag)
	buf.Write(payload)

	// Send to server
	if _, err := sc.conn.Write(buf.Bytes()); err != nil {
		sc.mu.Lock()
		delete(sc.pending, tag)
		sc.mu.Unlock()
		return nil, err
	}

	// Wait for reply with timeout
	select {
	case reply := <-replyChan:
		return reply, nil
	case <-time.After(sc.timeout):
		sc.mu.Lock()
		delete(sc.pending, tag)
		sc.mu.Unlock()
		return nil, io.ErrShortBuffer // or your timeout error
	}
}

// Read loop: matches replies to pending requests
func (sc *pulseClient) readLoop() {
	for {
		var opcode, tag uint32
		if err := binary.Read(sc.reader, binary.BigEndian, &opcode); err != nil {
			return
		}
		if err := binary.Read(sc.reader, binary.BigEndian, &tag); err != nil {
			return
		}

		// Read remaining payload (simplified - read fixed size or use length prefix)
		payload := make([]byte, 256)
		n, err := sc.reader.Read(payload)
		if err != nil {
			return
		}
		payload = payload[:n]

		// Deliver to waiting request
		sc.mu.Lock()
		if ch, ok := sc.pending[tag]; ok {
			delete(sc.pending, tag)
			sc.mu.Unlock()
			ch <- payload
		} else {
			sc.mu.Unlock()
		}
	}
}
