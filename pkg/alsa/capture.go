package alsa

import (
	"errors"
	"github.com/Binozo/GoAlsa/internal"
	"reflect"
	"unsafe"
)

/*
#cgo pkg-config: alsa
#include <alsa/asoundlib.h>
#include <stdint.h>
*/
import "C"
import (
	"fmt"
)

type CaptureDevice struct {
	defaultDevice
}

func NewCaptureDevice(deviceName string, audioConfig Config) (*CaptureDevice, error) {
	internalDevice, err := newDevice(deviceName, audioConfig, true, BufferParams{})
	if err != nil {
		return nil, err
	}
	return &CaptureDevice{
		*internalDevice,
	}, nil
}

func (c *CaptureDevice) Read(buffer []float32) (read int, err error) {
	bufVal := reflect.ValueOf(buffer)
	bufferLen := bufVal.Len()
	targetBuf := bufVal.Slice(0, bufferLen)

	frames := bufferLen / c.AudioConfig.Channels
	bufPtr := unsafe.Pointer(targetBuf.Index(0).Addr().Pointer())

	readResult := C.snd_pcm_readi(c.pcmDevice, bufPtr, C.snd_pcm_uframes_t(frames))
	if readResult == -C.EPIPE {
		C.snd_pcm_prepare(c.pcmDevice)
		return 0, ErrOverrun
	} else if int(readResult) < 0 {
		return 0, errors.Join(ErrReadError, fmt.Errorf("could not read: %d (%s)", int(readResult), GetErrorMessage(int(readResult))))
	}
	return int(readResult) * c.AudioConfig.Channels, nil
}

func (c *CaptureDevice) ReadBytes(buffer []byte) (read int, err error) {
	buf := make([]float32, len(buffer)/4)
	bufRead, err := c.Read(buf)
	if err != nil {
		return 0, err
	}
	buffer = internal.ConvertFloatsToBytes(buf[0:bufRead])
	return len(buffer), nil
}
