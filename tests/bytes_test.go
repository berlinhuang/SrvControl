package test

import (
	"SrvControl/utils/bytes"
	"fmt"
	"testing"
)

func TestBuffer_SeekPos(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteUTFString("John")
	buf.WriteUInt16(24)
	pos := buf.TellPos()

	buf.WriteBool(true)
	buf.WriteInt32(-1024)

	buf.Seek(pos, 0)

	buf.WriteBool(false)
	buf.WriteInt32(8000)

	for i := 1; i <= 100; i++ {
		buf.WriteFloat64(23.5)
	}

	fmt.Println(
		buf.ReadUTFString(),
		buf.ReadUInt16(),
		buf.ReadBool(),
		buf.ReadInt32())
	for i := 1; i <= 100; i++ {
		fmt.Println(buf.ReadFloat64())
	}
}

func TestBuffer_ReadBool(t *testing.T) {
	var buf bytes.Buffer

	buf.WriteBool(true)
	buf.WriteBool(false)
	fmt.Println(buf.ReadBool(), buf.ReadBool())
}

func TestBuffer_ReadInt8(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteInt8(127)
	buf.WriteInt8(-128)
	fmt.Println(buf.ReadInt8(), buf.ReadInt8())
}

func TestBuffer_ReadInt16(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteInt16(32767)
	buf.WriteInt16(-32768)
	fmt.Println(buf.ReadInt16(), buf.ReadInt16())
}

func TestBuffer_ReadInt32(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteInt32(2147483647)
	buf.WriteInt32(-2147483648)
	fmt.Println(buf.ReadInt32(), buf.ReadInt32())
}

func TestBuffer_ReadInt64(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteInt64(9223372036854775807)
	buf.WriteInt64(-9223372036854775808)
	fmt.Println(buf.ReadInt64(), buf.ReadInt64())
}

func TestBuffer_ReadUInt8(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteUInt8(255)
	fmt.Println(buf.ReadUInt8())
}

func TestBuffer_ReadUInt16(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteUInt16(65535)
	fmt.Println(buf.ReadUInt16())
}

func TestBuffer_ReadUInt32(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteUInt32(4294967295)
	fmt.Println(buf.ReadUInt32())
}

func TestBuffer_ReadUInt64(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteUInt64(18446744073709551615)
	fmt.Println(buf.ReadUInt64())
}

func TestBuffer_ReadFloat32(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteFloat32(34.4)
	fmt.Println(buf.ReadFloat32())
}

func TestBuffer_ReadFloat64(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteFloat64(65.23222)
	fmt.Println(buf.ReadFloat64())
}

func TestBuffer_WriteGBKString(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteGBKString("huang")
	buf.WriteGBKString("bolin")
	fmt.Println(buf.ReadGBKString(), buf.ReadGBKString())
}

func TestBuffer_WriteUTFString(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteUTFString("hello")
	buf.WriteUTFString("world")
	fmt.Println(buf.ReadUTFString(), buf.ReadUTFString())
}

func TestBuffer_ReadString(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("aaaaaaaaaa\n")
	str, _ := buf.ReadString('\n')
	fmt.Println(str)
}

func TestBuffer_WriteString(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("hello\x00world\n")
	str1, _ := buf.ReadString(0)
	str2, _ := buf.ReadString('\n')
	fmt.Println(str1)
	fmt.Println(str2)
}

func TestBuffer_WriteRead(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteUTFString("John")
	buf.WriteUInt16(24)
	buf.WriteBool(true)
	buf.WriteInt32(-1024)

	fmt.Println(
		buf.ReadUTFString(),
		buf.ReadUInt16(),
		buf.ReadBool(),
		buf.ReadInt32())
}

func TestBuffer_Write1(t *testing.T) {
	var buf bytes.Buffer
	bytes := []byte{0}
	buf.SetBuf(bytes)
	fmt.Println(buf.ReadInt16())
}
