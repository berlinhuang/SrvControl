package bytes

import (
	errors "SrvControl/utils/error"
	"bytes"
	"encoding/binary"
	"io"
	"unicode/utf8"
)

type Buffer struct {
	buf       []byte   // contents are the bytes buf[off : len(buf)]
	off       int      // read at &buf[off], write at &buf[len(buf)]
	bootstrap [64]byte // memory to hold first slice; helps small buffers avoid allocation.
	lastRead  readOp   // last read operation, so that Unread* can work correctly.
	pos       int
}

// The readOp constants describe the last action performed on
// the buffer, so that UnreadRune and UnreadByte can check for
// invalid usage. opReadRuneX constants are chosen such that
// converted to int they correspond to the rune size that was read.
type readOp int

const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid          = 0  // Non-read operation.
	opReadRune1        = 1  // Read rune of size 1.
	opReadRune2        = 2  // Read rune of size 2.
	opReadRune3        = 3  // Read rune of size 3.
	opReadRune4        = 4  // Read rune of size 4.
)

const (
	SeekStart   = 0 // seek relative to the origin of the file
	SeekCurrent = 1 // seek relative to the current offset
	SeekEnd     = 2 // seek relative to the end
)

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrTooLarge = errors.New("bytes.Buffer: too large")

// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like Read, Write, Reset, or Truncate).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
func (b *Buffer) Bytes() []byte { return b.buf[b.off:] }

// String returns the contents of the unread portion of the buffer
// as a string. If the Buffer is a nil pointer, it returns "<nil>".
func (b *Buffer) String() string {
	if b == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(b.buf[b.off:])
}

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (b *Buffer) OffLen() int { return len(b.buf) - b.off }

func (b *Buffer) PosLen() int { return len(b.buf) - b.pos }

func (b *Buffer) PosOffLen() int { return b.pos - b.off }

func (b *Buffer) BufLen() int { return len(b.buf) }

func (b *Buffer) SetPos(pos int) { b.pos = pos }

func (b *Buffer) SetOff(off int) { b.off = off }

func (b *Buffer) GetOff() int { return b.off }

func (b *Buffer) GetBuf() []byte { return b.buf }

func (b *Buffer) SetBuf(bytes []byte) { b.buf = bytes }

func (b *Buffer) doCopy() {
	//b.p
}

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space allocated for the buffer's data.
func (b *Buffer) Cap() int { return cap(b.buf) }

// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
func (b *Buffer) Truncate(n int) {
	b.lastRead = opInvalid
	switch {
	case n < 0 || n > b.OffLen():
		panic("bytes.Buffer: truncation out of range")
	case n == 0:
		// Reuse buffer space.
		b.off = 0
	}
	b.buf = b.buf[0 : b.off+n]
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (b *Buffer) Reset() { b.Truncate(0) }

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
//func (b *Buffer) grow(n int) int {
//	m := b.OffLen()
//	// If buffer is empty, reset to recover space.
//	if m == 0 && b.off != 0 {
//		b.Truncate(0)
//	}
//	if len(b.buf)+n > cap(b.buf) {
//		var buf []byte
//		if b.buf == nil && n <= len(b.bootstrap) {
//			buf = b.bootstrap[0:]
//		} else if m+n <= cap(b.buf)/2 {// b.buf足够多 不用分配
//			// We can slide things down instead of allocating a new
//			// slice. We only need m+n <= cap(b.buf) to slide, but
//			// we instead let capacity get twice as large so we
//			// don't spend all our time copying.
//			copy(b.buf[:], b.buf[b.off:]) //向前移动复用空间，后面off会置0
//			buf = b.buf[:m]
//		} else {
//			// not enough space anywhere
//			buf = makeSlice(2*cap(b.buf) + n)
//			copy(buf, b.buf[b.off:])
//		}
//		b.buf = buf
//		b.off = 0
//	}
//	b.buf = b.buf[0 : b.off+m+n]
//	return b.off + m
//}

func (b *Buffer) grow(n int) int {
	if b.pos+n < len(b.buf) { //有进行seek
		return b.pos
	}
	m := b.OffLen()
	// If buffer is empty, reset to recover space.
	if m == 0 && b.pos != 0 { //buffer 没了，写的位置不在最前面
		b.Truncate(0)
	}
	if b.pos+n > cap(b.buf) {
		var buf []byte
		if b.buf == nil && n <= len(b.bootstrap) {
			buf = b.bootstrap[0:]
		} else if b.PosOffLen()+n <= cap(b.buf)/2 { // b.buf足够多 不用分配
			// We can slide things down instead of allocating a new
			// slice. We only need m+n <= cap(b.buf) to slide, but
			// we instead let capacity get twice as large so we
			// don't spend all our time copying.
			copy(b.buf[:], b.buf[b.off:]) //向前移动复用空间，后面off会置0
			buf = b.buf[:m]
		} else {
			// not enough space anywhere
			buf = makeSlice(2*cap(b.buf) + n)
			copy(buf, b.buf[b.off:])
		}
		b.buf = buf
		b.off = 0
		b.pos = b.PosOffLen()
	}
	if b.off+m+n < b.pos+n {
		b.buf = b.buf[0 : b.pos+n]
	} else {
		b.buf = b.buf[0 : b.off+m+n]
	}
	return b.pos
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer) Grow(n int) {
	if n < 0 {
		panic("bytes.Buffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[0:m]
}

// MinRead is the minimum slice size passed to a Read call by
// Buffer.ReadFrom. As long as the Buffer has at least MinRead bytes beyond
// what is required to hold the contents of r, ReadFrom will not grow the
// underlying buffer.
const MinRead = 512

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with ErrTooLarge.
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.lastRead = opInvalid
	// If buffer is empty, reset to recover space.
	if b.off >= len(b.buf) {
		b.Truncate(0)
	}
	for {
		if free := cap(b.buf) - len(b.buf); free < MinRead {
			// not enough space at end
			newBuf := b.buf
			if b.off+free < MinRead {
				// not enough space using beginning of buffer;
				// double buffer capacity
				newBuf = makeSlice(2*cap(b.buf) + MinRead)
			}
			copy(newBuf, b.buf[b.off:])
			b.buf = newBuf[:len(b.buf)-b.off]
			b.off = 0
		}
		m, e := r.Read(b.buf[len(b.buf):cap(b.buf)])
		b.buf = b.buf[0 : len(b.buf)+m]
		n += int64(m)
		if e == io.EOF {
			break
		}
		if e != nil {
			return n, e
		}
	}
	return n, nil // err is EOF, so return nil explicitly
}

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	return make([]byte, n)
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	b.lastRead = opInvalid
	if b.off < len(b.buf) {
		nBytes := b.OffLen()
		m, e := w.Write(b.buf[b.off:])
		if m > nBytes {
			panic("bytes.Buffer.WriteTo: invalid Write count")
		}
		b.off += m
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	// Buffer is now empty; reset.
	b.Truncate(0)
	return
}

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
func (b *Buffer) Read(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	if b.off >= len(b.buf) {
		// Buffer is empty, reset to recover space.
		b.Truncate(0)
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}
	n = copy(p, b.buf[b.off:])
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (b *Buffer) Next(n int) []byte {
	b.lastRead = opInvalid
	m := b.OffLen()
	if n > m {
		n = m
	}
	data := b.buf[b.off : b.off+n]
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return data
}

// ReadByte reads and returns the next byte from the buffer.
// If no byte is available, it returns error io.EOF.
func (b *Buffer) ReadByte() (byte, error) {
	b.lastRead = opInvalid
	if b.off >= len(b.buf) {
		// Buffer is empty, reset to recover space.
		b.Truncate(0)
		return 0, io.EOF
	}
	c := b.buf[b.off]
	b.off++
	b.lastRead = opRead
	return c, nil
}

// ReadRune reads and returns the next UTF-8-encoded
// Unicode code point from the buffer.
// If no bytes are available, the error returned is io.EOF.
// If the bytes are an erroneous UTF-8 encoding, it
// consumes one byte and returns U+FFFD, 1.
func (b *Buffer) ReadRune() (r rune, size int, err error) {
	b.lastRead = opInvalid
	if b.off >= len(b.buf) {
		// Buffer is empty, reset to recover space.
		b.Truncate(0)
		return 0, 0, io.EOF
	}
	c := b.buf[b.off]
	if c < utf8.RuneSelf {
		b.off++
		b.lastRead = opReadRune1
		return rune(c), 1, nil
	}
	r, n := utf8.DecodeRune(b.buf[b.off:])
	b.off += n
	b.lastRead = readOp(n)
	return r, n, nil
}

// UnreadRune unreads the last rune returned by ReadRune.
// If the most recent read or write operation on the buffer was
// not a ReadRune, UnreadRune returns an error.  (In this regard
// it is stricter than UnreadByte, which will unread the last byte
// from any read operation.)
func (b *Buffer) UnreadRune() error {
	if b.lastRead <= opInvalid {
		return errors.New("bytes.Buffer: UnreadRune: previous operation was not ReadRune")
	}
	if b.off >= int(b.lastRead) {
		b.off -= int(b.lastRead)
	}
	b.lastRead = opInvalid
	return nil
}

// UnreadByte unreads the last byte returned by the most recent
// read operation. If write has happened since the last read, UnreadByte
// returns an error.
func (b *Buffer) UnreadByte() error {
	if b.lastRead == opInvalid {
		return errors.New("bytes.Buffer: UnreadByte: previous operation was not a read")
	}
	b.lastRead = opInvalid
	if b.off > 0 {
		b.off--
	}
	return nil
}

// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
// If ReadBytes encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadBytes returns err != nil if and only if the returned data does not end in
// delim.
func (b *Buffer) ReadBytes(delim byte) (line []byte, err error) {
	slice, err := b.readSlice(delim)
	// return a copy of slice. The buffer's backing array may
	// be overwritten by later calls.
	line = append(line, slice...)
	return
}

// readSlice is like ReadBytes but returns a reference to internal buffer data.
func (b *Buffer) readSlice(delim byte) (line []byte, err error) {
	i := bytes.IndexByte(b.buf[b.off:], delim)
	end := b.off + i + 1
	if i < 0 {
		end = len(b.buf)
		err = io.EOF
	}
	line = b.buf[b.off : end-1] //end  or end-1
	b.off = end
	b.lastRead = opRead
	return line, err
}

// ReadString reads until the first occurrence of delim in the input,
// returning a string containing the data up to and including the delimiter.
// If ReadString encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadString returns err != nil if and only if the returned data does not end
// in delim.
func (b *Buffer) ReadString(delim byte) (line string, err error) {
	slice, err := b.readSlice(delim)
	return string(slice), err
}

// NewBuffer creates and initializes a new Buffer using buf as its initial
// contents. It is intended to prepare a Buffer to read existing data. It
// can also be used to size the internal buffer for writing. To do that,
// buf should have the desired capacity but a length of zero.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewBuffer(buf []byte) *Buffer { return &Buffer{buf: buf} }

// NewBufferString creates and initializes a new Buffer using string s as its
// initial contents. It is intended to prepare a buffer to read an existing
// string.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewBufferString(s string) *Buffer {
	return &Buffer{buf: []byte(s)}
}

// WriteMultiByte
// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (b *Buffer) WriteMultiByte(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	m := b.grow(len(p))
	b.pos += len(p)
	return copy(b.buf[m:], p), nil
}

func (b *Buffer) ReadExternalByte(bytes []byte, offset int, length int) {
	b.lastRead = opInvalid
	copy(bytes[offset:length+offset], b.buf[b.off:b.off+length])
	b.off += length
}

func (b *Buffer) ReadMultiByte(length int) []byte {
	if length <= 0 {
		return []byte{}
	}
	b.lastRead = opInvalid
	bytes := b.buf[b.off : b.off+length]
	b.off += length
	return bytes
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	m := b.grow(len(p))
	b.pos += len(p)
	return copy(b.buf[m:], p), nil
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
func (b *Buffer) WriteString(s string) (n int, err error) {
	b.lastRead = opInvalid
	m := b.grow(len(s))
	b.pos += len(s)
	return copy(b.buf[m:], s), nil
}

// WriteByte appends the byte c to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match bufio.Writer's
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// ErrTooLarge.
func (b *Buffer) WriteByte(c byte) error {
	b.lastRead = opInvalid
	m := b.grow(1)
	b.buf[m] = c
	b.pos += 1
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to the
// buffer, returning its length and an error, which is always nil but is
// included to match bufio.Writer's WriteRune. The buffer is grown as needed;
// if it becomes too large, WriteRune will panic with ErrTooLarge.
func (b *Buffer) WriteRune(r rune) (n int, err error) {
	if r < utf8.RuneSelf {
		b.WriteByte(byte(r))
		return 1, nil
	}
	b.lastRead = opInvalid
	m := b.grow(utf8.UTFMax)
	n = utf8.EncodeRune(b.buf[m:m+utf8.UTFMax], r)
	b.buf = b.buf[:m+n]
	return n, nil
}

func (b *Buffer) WriteInt8(int8 int8) {
	b.WriteByte(byte(int8))
}

func (b *Buffer) WriteUInt8(uint8 uint8) {
	b.WriteByte(byte(uint8))
}

func (b *Buffer) WriteInt16(int16 int16) {
	bytes := Int16ToBytes(int16)
	b.WriteMultiByte(bytes)
}

func (b *Buffer) WriteUInt16(uint16 uint16) {
	bytes := UInt16ToBytes(uint16)
	b.WriteMultiByte(bytes)
}

func (b *Buffer) WriteInt32(int32 int32) {
	//p := b.grow(4)
	//(b.buf)[p] = byte(v >> 24)
	//(b.buf)[p+1] = byte(v >> 16)
	//(b.buf)[p+2] = byte(v >> 8)
	//(b.buf)[p+3] = byte(v)
	bytes := Int32ToBytes(int32)
	b.WriteMultiByte(bytes)
}

func (b *Buffer) WriteUInt32(uint32 uint32) {
	bytes := UInt32ToBytes(uint32)
	b.WriteMultiByte(bytes)
}

func (b *Buffer) WriteInt64(int64 int64) {
	bytes := Int64ToBytes(int64)
	b.WriteMultiByte(bytes)
}

func (b *Buffer) WriteUInt64(uint64 uint64) {
	bytes := UInt64ToBytes(uint64)
	b.WriteMultiByte(bytes)
}

func (b *Buffer) WriteFloat32(f32 float32) {
	bytes := Float32ToBytes(f32)
	b.WriteMultiByte(bytes)
}

func (b *Buffer) WriteFloat64(f64 float64) {
	bytes := Float64ToBytes(f64)
	b.WriteMultiByte(bytes)
}

func (b *Buffer) WriteBool(bool bool) {
	if bool {
		b.WriteByte(byte(1))
	} else {
		b.WriteByte(byte(0))
	}
}

///////////////////////////////////////////////////////////////////////

func (b *Buffer) ReadBool() bool {
	bt, _ := b.ReadByte()
	return bt == byte(1)
}

func (b *Buffer) ReadInt8() int8 {
	bt, _ := b.ReadByte()
	return int8(bt)
}

func (b *Buffer) ReadUInt8() uint8 {
	bt, _ := b.ReadByte()
	return uint8(bt)
}

func (b *Buffer) ReadInt16() int16 {
	bytes := make([]byte, SizeShort)
	b.Read(bytes)
	return int16(binary.LittleEndian.Uint16(bytes))
}

func (b *Buffer) ReadUInt16() uint16 {
	bytes := make([]byte, SizeShort)
	b.Read(bytes)
	return binary.LittleEndian.Uint16(bytes)
}

func (b *Buffer) ReadInt32() int32 {
	bytes := make([]byte, SizeInt)
	b.Read(bytes)
	return int32(binary.LittleEndian.Uint32(bytes))
}

func (b *Buffer) ReadUInt32() uint32 {
	bytes := make([]byte, SizeInt)
	b.Read(bytes)
	return binary.LittleEndian.Uint32(bytes)
}

func (b *Buffer) ReadInt64() int64 {
	bytes := make([]byte, SizeLong)
	b.Read(bytes)
	return int64(binary.LittleEndian.Uint64(bytes))
}

func (b *Buffer) ReadUInt64() uint64 {
	bytes := make([]byte, SizeLong)
	b.Read(bytes)
	return binary.LittleEndian.Uint64(bytes)
}

func (b *Buffer) ReadFloat32() float32 {
	bytes := make([]byte, SizeFloat)
	b.Read(bytes)
	return BytesToFloat32(bytes)
}

func (b *Buffer) ReadFloat64() float64 {
	bytes := make([]byte, SizeDouble)
	b.Read(bytes)
	return BytesToFloat64(bytes)
}

func (b *Buffer) WriteUTFString(s string) {
	b.WriteString(s)
	b.WriteString("\x00")
}

func (b *Buffer) ReadUTFString() string {
	str, _ := b.ReadString(0)
	return str
}

func (b *Buffer) WriteGBKString(s string) {
	bytes := []byte(s)
	len := len(bytes)
	buf, len := StringLenToByte(uint32(len))
	b.WriteBytes(buf, 0, len, 0)
	binary.Write(b, binary.LittleEndian, bytes)
}

func (b *Buffer) ReadGBKString() string {
	buf := make([]byte, 4)
	buf[0], _ = b.ReadByte()
	remain := (uint)(buf[0] >> 6)
	if remain != 0 {
		b.ReadExternalByte(buf, 1, int(remain))
		//length := remain + 1
	} else {
		//length := 1
	}
	readLen := GetStringLen(buf)
	return string(b.ReadMultiByte((int)(readLen)))
}

func (b *Buffer) TellPos() int {
	return b.pos
}

func (b *Buffer) Seek(offset int, whence int) (int, error) {
	if whence < 0 || whence > 2 {
		return 0, errors.ErrParamNotExist
	}
	var pos = 0
	if whence == SeekStart { //文件开头
		pos = 0 + offset
	} else if whence == SeekCurrent { //当前位置
		pos = b.pos + offset
	} else if whence == SeekEnd { //文件尾
		pos = b.OffLen() + offset
	}
	if pos < 0 {
		return 0, errors.ErrOverBuffer
	} else {
		b.pos = pos
		return pos, nil
	}
}

///////////////////////////////////////
func (b *Buffer) WriteBytes(bytes []byte, offset int, length int, bytesLen uint32) {

	//if (length<=0) {
	//	length = (int)(len(bytes))-offset
	//}
	//bytesLen = (uint32)(len(bytes))
	//
	//if ( bytesLen <= 0 ){
	//	bytesLen =(uint32)(len(bytes))
	//}
	//if(offset<0 || uint32(offset) >= bytesLen ){
	//	return
	//}
	//if(length<0 || uint32(length +offset) > bytesLen){
	//	return
	//}
	//b.SetPos(offset)
	b.WriteMultiByte(bytes[offset:length])
	b.SetPos(len(b.buf))
}
