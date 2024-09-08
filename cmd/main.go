package main

import (
	"github.com/Binozo/GoAlsa/pkg/alsa"
)

func main() {
	device, err := alsa.NewPlaybackDevice("hw:0,0", alsa.Config{
		Channels:   2,
		Format:     alsa.FormatS16LE,
		SampleRate: 48000,
	})
	if err != nil {
		panic(err)
	}
	defer device.Close()
}
