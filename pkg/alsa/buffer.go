package alsa

/*
#cgo pkg-config: alsa
#include <alsa/asoundlib.h>
#include <stdint.h>
*/
import "C"

// BufferParams specifies the buffer parameters of a device.
type BufferParams struct {
	BufferFrames int
	PeriodFrames int
	Periods      int
}
