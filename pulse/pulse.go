package pulse

import (
	"math"
	"github.com/ImOutOfIdeas/audigo/internal"

	"github.com/jfreymuth/pulse"
)

type pulseBackend struct {
	client *pulse.Client
}

func New() (internal.Backend, error) {
	var b pulseBackend

	client, err := pulse.NewClient()
	if err != nil {
		return nil, err
	}

	b.client = client

	return &b, nil
}

func (b *pulseBackend) OpenStream(internal.StreamConfig) (internal.Stream, error) {
	stream, err := b.client.NewPlayback(pulse.Float32Reader(synth), pulse.PlaybackLatency(.1))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *pulseBackend) Close() error {
	b.Close()

	return nil
}

var t, phase float32
func synth(out []float32) (int, error) {
	for i := range out {
		if t > 4 {
			return i, pulse.EndOfData
		}
		x := float32(math.Sin(2 * math.Pi * float64(phase)))
		out[i] = x * 0.1
		f := [...]float32{440, 550, 660, 880}[int(2*t)&3]
		phase += f / 44100
		if phase >= 1 {
			phase--
		}
		t += 1. / 44100
	}
	return len(out), nil
}
