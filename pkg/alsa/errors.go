package alsa

/*
#cgo pkg-config: alsa
#include <alsa/asoundlib.h>
*/
import "C"
import "errors"

var (
	// ErrOverrun signals an overrun error
	ErrOverrun = errors.New("overrun")
	// ErrUnderrun signals an underrun error
	ErrUnderrun    = errors.New("underrun")
	ErrOpenError   = errors.New("could not open alsa device")
	ErrParamsError = errors.New("could not set params")
	ErrReadError   = errors.New("could not read")
	ErrWriteError  = errors.New("could not write")
)

func GetErrorMessage(errorcode int) string {
	return C.GoString(C.snd_strerror(C.int(errorcode)))
}
