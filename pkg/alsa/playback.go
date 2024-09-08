package alsa

/*
#cgo pkg-config: alsa
#include <alsa/asoundlib.h>
#include <stdint.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"github.com/Binozo/GoAlsa/internal"
	"reflect"
	"unsafe"
)

type PlaybackDevice struct {
	defaultDevice
}

func NewPlaybackDevice(deviceName string, audioConfig Config) (*PlaybackDevice, error) {
	internalDevice, err := newDevice(deviceName, audioConfig, false, BufferParams{})
	if err != nil {
		return nil, err
	}
	return &PlaybackDevice{
		*internalDevice,
	}, nil
}

func (p *PlaybackDevice) Write(buffer []float32) (samples int, err error) {
	bufVal := reflect.ValueOf(buffer)
	bufferLen := bufVal.Len()
	targetBuf := bufVal.Slice(0, bufferLen)

	frames := C.snd_pcm_uframes_t(bufferLen / p.AudioConfig.Channels)
	bufPtr := unsafe.Pointer(targetBuf.Index(0).Addr().Pointer())

	writeResult := C.snd_pcm_writei(p.pcmDevice, bufPtr, frames)
	if writeResult == -C.EPIPE {
		C.snd_pcm_prepare(p.pcmDevice)
		return 0, ErrUnderrun
	} else if writeResult < 0 {
		return 0, errors.Join(ErrWriteError, fmt.Errorf("could not write: %d (%s)", int(writeResult), GetErrorMessage(writeResult)))
	}

	return int(writeResult) * p.AudioConfig.Channels, nil
}

func (p *PlaybackDevice) WriteBytes(buffer []byte) (samples int, err error) {
	buf := internal.ConvertBytesToFloats(buffer)
	written, err := p.Write(buf)
	return written * 4, err
}
