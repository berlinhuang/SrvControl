package util

import (
	"bytes"
	"sync"
)

//// 缓冲区 线程池
//var byteBufferPool = &sync.Pool{
//	New: func() interface{} { //对象池没有对象的时候 调用Get会调用New来获取对象
//		return &bytes.Buffer{}
//	},
//}
//
//// 获取缓冲区
//func AcquireByteBuffer() *bytes.Buffer {
//	//get顺序 local private ->local shared -> remote shared -> victim
//	return byteBufferPool.Get().(*bytes.Buffer)
//}
//
////
//func ReleaseByteBuffer(b *bytes.Buffer) {
//	if b != nil {
//		b.Reset()
//		byteBufferPool.Put(b)
//	}
//}

var bufferPool sync.Pool

func AcquireByteBuffer() *bytes.Buffer {
	v := bufferPool.Get()
	if v == nil {
		return bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
	}
	buf := v.(*bytes.Buffer)
	return buf
}

func ReleaseByteBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}
