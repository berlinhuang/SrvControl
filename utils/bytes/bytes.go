package bytes

import (
	"encoding/binary"
	"math"
)

const (
	SizeByte   = 1
	SizeShort  = 2 //int16
	SizeInt    = 4 //int32
	SizeLong   = 8 //int64
	SizeFloat  = SizeInt
	SizeDouble = SizeLong
)

func Int16ToBytes(i16 int16) []byte {
	bytes := make([]byte, SizeShort)
	binary.LittleEndian.PutUint16(bytes, uint16(i16))
	return bytes
}

func Int32ToBytes(i32 int32) []byte {
	bytes := make([]byte, SizeInt)
	binary.LittleEndian.PutUint32(bytes, uint32(i32))
	return bytes
}

func Int64ToBytes(i64 int64) []byte {
	bytes := make([]byte, SizeLong)
	binary.LittleEndian.PutUint64(bytes, uint64(i64))
	return bytes
}

func UInt16ToBytes(ui16 uint16) []byte {
	bytes := make([]byte, SizeShort)
	binary.LittleEndian.PutUint16(bytes, ui16)
	return bytes
}

func UInt32ToBytes(ui32 uint32) []byte {
	bytes := make([]byte, SizeInt)
	binary.LittleEndian.PutUint32(bytes, ui32)
	return bytes
}

func UInt64ToBytes(ui64 uint64) []byte {
	bytes := make([]byte, SizeLong)
	binary.LittleEndian.PutUint64(bytes, ui64)
	return bytes
}

// float 2 bytes
func Float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, SizeFloat)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, SizeDouble)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func BytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func StringLenToByte(u32 uint32) ([]byte, int) {
	buf := make([]byte, 4)
	len := 1
	if u32 < 64 {
		buf[0] = (byte)(u32)
	} else {
		buf[0] = (byte)(u32 & 0x3F)
		u32 >>= 6
		for {
			if u32 == 0 {
				break
			}
			buf[len] = (byte)(u32)
			u32 >>= 8
			len++
		}
		buf[0] |= (byte)((len - 1) << 6)
	}
	return buf, len
}

func GetStringLen(buf []byte) uint {
	var u32 uint = 0
	by := int(buf[0] >> 6)
	if by == 0 {
		u32 = uint(buf[0])
	} else {
		u32 = (uint)(buf[0] & 0x3F)
		var i int
		var iBitOff uint = 6
		for i = 0; i < by; i++ {
			u32 += ((uint)(buf[i+1])) << iBitOff
			iBitOff += 8
		}
	}
	return u32
}
