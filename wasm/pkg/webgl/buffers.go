package webgl

import (
	"reflect"
	"unsafe"
)

type BufferData interface {
	Bytes() []byte
}

func float32SliceAsByteSlice(floats []float32) []byte {
	n := 4 * len(floats)

	up := unsafe.Pointer(&(floats[0]))
	pi := (*[1]byte)(up)
	buf := (*pi)[:]
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Len = n
	sh.Cap = n

	return buf
}

type Float32ArrayBuffer []float32

func (b Float32ArrayBuffer) Bytes() []byte {
	return float32SliceAsByteSlice([]float32(b))
}

type ByteArrayBuffer []byte

func (b ByteArrayBuffer) Bytes() []byte {
	return b
}

func byteSliceAsUInt32Slice(bytes []byte) []uint32 {
	n := len(bytes) / 4

	up := unsafe.Pointer(&(bytes[0]))
	pi := (*[1]uint32)(up)
	buf := (*pi)[:]
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Len = n
	sh.Cap = n

	return buf
}

func (b ByteArrayBuffer) UInt32Slice() []uint32 {
	return byteSliceAsUInt32Slice(b)
}

func uint16SliceAsByteSlice(b []uint16) []byte {
	n := 2 * len(b)

	up := unsafe.Pointer(&(b[0]))
	pi := (*[1]byte)(up)
	buf := (*pi)[:]
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Len = n
	sh.Cap = n

	return buf
}

type Uint16ArrayBuffer []uint16

func (b Uint16ArrayBuffer) Bytes() []byte {
	return uint16SliceAsByteSlice([]uint16(b))
}
