# GoAlsa 🎵

Go bindings for [ALSA](https://www.alsa-project.org)

## Installation 🚀

Make sure you have the development headers installed:
#### Debian/Ubuntu
```bash
$ sudo apt install libasound2-dev -y
```
#### Arch
```bash
$ sudo pacman -S alsa-lib -y
```
#### Fedora
```bash
$ sudo dnf install alsa-lib-devel -y
```

#### Now install the go package:
```bash
$ go get -u github.com/Binozo/GoAlsa
```

## Quickstart 💫
```go
package main

import (
	"github.com/Binozo/GoAlsa/pkg/alsa"
)

func main() {
	//hw:<CARD_NR>,<DEVICE_NR>
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
```