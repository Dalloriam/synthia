# Synthia
[![Go Report Card](https://goreportcard.com/badge/github.com/dalloriam/synthia)](https://goreportcard.com/report/github.com/dalloriam/synthia)
[![GoDoc](https://godoc.org/github.com/dalloriam/synthia?status.svg)](https://godoc.org/github.com/dalloriam/synthia)
[![CircleCI](https://circleci.com/gh/dalloriam/synthia.svg?style=svg)](https://circleci.com/gh/dalloriam/synthia)

Synthia is a standalone (as in _dependency-free_) toolkit for synthesizing audio in Go. The project is still in its
early stages, but should improve at a good pace.

## Quick Start

```go
package main

import (
	"github.com/gordonklaus/portaudio"
	"github.com/dalloriam/synthia/modular"
	"github.com/dalloriam/synthia"
)

const (
	bufferSize = 512 // Size of the audio buffer.
	sampleRate = 44100 // Audio sample rate.
	audioChannelCount = 2 // Stereo.
	mixerChannelCount = 2 // Three oscillators.
)

type AudioBackend struct {
	params portaudio.StreamParameters
}

func (b *AudioBackend) Start(callback func(in []float32, out [][]float32)) error {
	stream, err := portaudio.OpenStream(b.params, callback)
	if err != nil {
		return err
	}

	return stream.Start()
}

func (b *AudioBackend) FrameSize() int {
	return b.params.FramesPerBuffer
}

func newBackend() *AudioBackend {
	// Quick-and-dirty way to initialize portaudio
	if err := portaudio.Initialize(); err != nil {
		panic(err)
	}

	devices, err := portaudio.Devices()
	if err != nil {
		panic(err)
	}
	inDevice, outDevice := devices[0], devices[1]

	params := portaudio.LowLatencyParameters(inDevice, outDevice)

	params.Input.Channels = 1
	params.Output.Channels = audioChannelCount

	params.SampleRate = float64(sampleRate)
	params.FramesPerBuffer = bufferSize

	return &AudioBackend{params}
}

func main() {

	backend := newBackend()

	// Set tempo to 60 bpm
	clock := modular.NewClock()
	clock.Tempo.SetValue(60)

	// Set the sequence to a C Major scale
	seq := modular.NewSequencer([]float64{130.81, 146.83, 164.1, 174.61, 196, 220, 246.94, 261.63})
	seq.Clock = clock

	// Play the sequence in eighth notes
	seq.BeatsPerStep = 0.25

	// Create an oscillator and attach the sequencer to it.
	osc1 := modular.NewOscillator()
	osc1.Frequency.Line = seq

	// Create the synthesizer with two mixer channels and set it to output to our audio backend.
	synth := synthia.New(mixerChannelCount, bufferSize, backend)

	// Map two different waves to the two outputs of our mixer.
	synth.Mixer.Channels[0].Input = osc1.Sine
	synth.Mixer.Channels[1].Input = osc1.Triangle

	// Pan the sine wave to the right channel and the triangle wave to the left channel
	synth.Mixer.Channels[0].Pan.SetValue(-1)
	synth.Mixer.Channels[1].Pan.SetValue(1)

	// Block until terminated
	select{}
}
```