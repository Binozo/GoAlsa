package alsa

import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

/*
#cgo pkg-config: alsa
#include <alsa/asoundlib.h>
#include <stdint.h>
*/
import "C"

type Device interface {
	GetConfig() Config
	Close()
}

type defaultDevice struct {
	AudioConfig  Config
	pcmDevice    *C.snd_pcm_t
	frames       int
	bufferParams BufferParams
}

func (d *defaultDevice) GetConfig() Config {
	return d.AudioConfig
}

func (d *defaultDevice) Close() {
	if d.pcmDevice != nil {
		C.snd_pcm_drain(d.pcmDevice)
		C.snd_pcm_close(d.pcmDevice)
		d.pcmDevice = nil
	}
	runtime.SetFinalizer(d, nil)
}

func newDevice(deviceName string, audioConfig Config, captureMode bool, bufferParams BufferParams) (*defaultDevice, error) {
	deviceNameC := C.CString(deviceName)
	defer C.free(unsafe.Pointer(deviceNameC))
	var internalDevice *C.snd_pcm_t
	var openResult C.int
	if captureMode {
		openResult = C.snd_pcm_open(&internalDevice, deviceNameC, C.SND_PCM_STREAM_CAPTURE, 0)
	} else {
		openResult = C.snd_pcm_open(&internalDevice, deviceNameC, C.SND_PCM_STREAM_PLAYBACK, 0)
	}
	if openResult < 0 {
		return nil, errors.Join(ErrOpenError, fmt.Errorf("error code: %d", int(openResult)))
	}
	alsaDevice := &defaultDevice{
		AudioConfig: audioConfig,
		pcmDevice:   internalDevice,
	}
	runtime.SetFinalizer(alsaDevice, (*defaultDevice).Close)
	var hwParams *C.snd_pcm_hw_params_t
	if allocResult := C.snd_pcm_hw_params_malloc(&hwParams); allocResult < 0 {
		return nil, errors.Join(ErrParamsError, fmt.Errorf("could not alloc hw params: %d", int(allocResult)))
	}
	defer C.snd_pcm_hw_params_free(hwParams)
	if defaultParamsResult := C.snd_pcm_hw_params_any(alsaDevice.pcmDevice, hwParams); defaultParamsResult < 0 {
		return nil, errors.Join(ErrParamsError, fmt.Errorf("could not set default hw params: %d", int(defaultParamsResult)))
	}
	if setAccessParamsResult := C.snd_pcm_hw_params_set_access(alsaDevice.pcmDevice, hwParams, C.SND_PCM_ACCESS_RW_INTERLEAVED); setAccessParamsResult < 0 {
		return nil, errors.Join(ErrParamsError, fmt.Errorf("could not set access params: %d", int(setAccessParamsResult)))
	}
	if setFormatParamsResult := C.snd_pcm_hw_params_set_format(alsaDevice.pcmDevice, hwParams, C.snd_pcm_format_t(audioConfig.Format)); setFormatParamsResult < 0 {
		return nil, errors.Join(ErrParamsError, fmt.Errorf("could not set format params: %d", int(setFormatParamsResult)))
	}
	if setChannelParamsResult := C.snd_pcm_hw_params_set_channels(alsaDevice.pcmDevice, hwParams, C.uint(audioConfig.Channels)); setChannelParamsResult < 0 {
		return nil, errors.Join(ErrParamsError, fmt.Errorf("could not set channel params: %d", int(setChannelParamsResult)))
	}
	if setSampleRateParamsResult := C.snd_pcm_hw_params_set_rate(alsaDevice.pcmDevice, hwParams, C.uint(audioConfig.SampleRate), 0); setSampleRateParamsResult < 0 {
		return nil, errors.Join(ErrParamsError, fmt.Errorf("could not set sample rate params: %d", int(setSampleRateParamsResult)))
	}

	var bufferSize = C.snd_pcm_uframes_t(bufferParams.BufferFrames)
	if bufferParams.BufferFrames == 0 {
		// Default buffer size: max buffer size
		if bufferSizeMaxResult := C.snd_pcm_hw_params_get_buffer_size_max(hwParams, &bufferSize); bufferSizeMaxResult < 0 {
			return nil, errors.Join(ErrParamsError, fmt.Errorf("could get buffer size: %d", int(bufferSizeMaxResult)))
		}
	}
	if setBufferSizeResult := C.snd_pcm_hw_params_set_buffer_size_near(alsaDevice.pcmDevice, hwParams, &bufferSize); setBufferSizeResult < 0 {
		return nil, errors.Join(ErrParamsError, fmt.Errorf("could not set buffer size: %d", int(setBufferSizeResult)))
	}

	// Default period size: 1/8 of a second
	var periodFrames = C.snd_pcm_uframes_t(audioConfig.SampleRate / 8)
	if bufferParams.PeriodFrames > 0 {
		periodFrames = C.snd_pcm_uframes_t(bufferParams.PeriodFrames)
	} else if bufferParams.Periods > 0 {
		periodFrames = C.snd_pcm_uframes_t(int(bufferSize) / bufferParams.Periods)
	}

	if setPeriodSizeResult := C.snd_pcm_hw_params_set_period_size_near(alsaDevice.pcmDevice, hwParams, &periodFrames, nil); setPeriodSizeResult < 0 {
		return nil, errors.Join(ErrParamsError, fmt.Errorf("could not set period size: %d", int(setPeriodSizeResult)))
	}
	var periods = C.uint(0)
	if getPeriods := C.snd_pcm_hw_params_get_periods(hwParams, &periods, nil); getPeriods < 0 {
		return nil, errors.Join(ErrParamsError, fmt.Errorf("could get periods: %d", int(getPeriods)))
	}
	if setHwParamsResult := C.snd_pcm_hw_params(alsaDevice.pcmDevice, hwParams); setHwParamsResult < 0 {
		return nil, errors.Join(ErrParamsError, fmt.Errorf("could set hw params: %d", int(setHwParamsResult)))
	}
	alsaDevice.frames = int(periodFrames)
	alsaDevice.bufferParams.Periods = int(periods)
	alsaDevice.bufferParams.PeriodFrames = int(periodFrames)
	alsaDevice.bufferParams.BufferFrames = int(bufferSize)

	return alsaDevice, nil
}
