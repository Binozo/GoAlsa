package internal

import (
	"bytes"
	"encoding/binary"
	"math"
)

func ConvertFloatsToBytes(data []float32) []byte {
	var buf bytes.Buffer
	for i := 0; i < len(data); i++ {
		err := binary.Write(&buf, binary.LittleEndian, data[i])
		if err != nil {
			panic(err)
		}
	}

	return buf.Bytes()
}

func ConvertBytesToFloats(data []byte) []float32 {
	buf := make([]float32, len(data)/4)
	for i := 0; i < len(buf); i++ {
		buf[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[i*4 : i*4+4]))
	}
	return buf
}
