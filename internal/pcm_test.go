package internal

import (
	"bytes"
	"reflect"
	"testing"
)

func TestConvertFloatsToBytes(t *testing.T) {
	conversionData := []float32{0, 1, 2, 3}
	want := []byte{0, 0, 0, 0, 0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64}

	convertedData := ConvertFloatsToBytes(conversionData)
	if bytes.Compare(convertedData, want) != 0 {
		t.Errorf("ConvertFloatsToBytes failed: got %v, want %v", convertedData, want)
	}
}

func TestConvertBytesToFloats(t *testing.T) {
	conversionData := []byte{0, 0, 0, 0, 0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64}
	want := []float32{0, 1, 2, 3}

	convertedData := ConvertBytesToFloats(conversionData)
	if !reflect.DeepEqual(convertedData, want) {
		t.Errorf("ConvertBytesToFloats failed: got %v, want %v", convertedData, want)
	}
}
