package main

import (
	"fmt"
	"os"

	"github.com/ImOutOfIdeas/audigo"
)

func main() {
	backend, err := audigo.DefaultBackend()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in backend init: %v\n", err)

		return
	}
	defer backend.Close()

	return

	stream, err := backend.OpenStream(audigo.StreamConfig{
		Channels:   2,
		SampleRate: 44100,
		BufferSize: 512,
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer stream.Close()

	fmt.Println("stream opened successfully")

	fmt.Scan()
}
