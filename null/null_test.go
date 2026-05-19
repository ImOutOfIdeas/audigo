package null_test

import (
    "testing"

    "github.com/ImOutOfIdeas/audigo/internal"
    "github.com/ImOutOfIdeas/audigo/null"
)

func TestStreamLifecycle(t *testing.T) {
    b, err := null.New()
	if err != nil {
		t.Fatal(err)
	}
    defer b.Close()

    s, err := b.OpenStream(internal.StreamConfig{
        Channels:   2,
        SampleRate: 44100,
    })
    if err != nil {
        t.Fatal(err)
    }
    defer s.Close()

    if err := s.Start(); err != nil {
        t.Fatal(err)
    }

    samples := make([]float32, 512)
    n, err := s.Write(samples)
    if err != nil {
        t.Fatal(err)
    }
    if n != len(samples) {
        t.Fatalf("wrote %d, want %d", n, len(samples))
    }

    if err := s.Stop(); err != nil {
        t.Fatal(err)
    }
}
