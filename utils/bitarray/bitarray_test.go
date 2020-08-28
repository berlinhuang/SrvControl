package bitarray

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// require的函数会直接导致case结束
// assert虽然也标记为case失败，但case不会退出，而是继续往下执行。
func TestInitialize(t *testing.T) {
	bit := Bit{false}
	assert.Equal(t, false, bit.value)
}

func TestInitBitArray(t *testing.T) {
	var ba BitArray
	ba.InitBitArray(3)
	assert.Equal(t, false, ba.bits[0].value)
	assert.Equal(t, false, ba.bits[1].value)
	assert.Equal(t, false, ba.bits[2].value)
}

func TestGet(t *testing.T) {
	var ba BitArray
	ba.InitBitArray(3)
	ba.Set(3)
	assert.Equal(t, true, ba.Get(0).bits[0].value)
	assert.Equal(t, true, ba.Get(1).bits[0].value)
	assert.Equal(t, false, ba.Get(2).bits[0].value)
}

func TestToInt(t *testing.T) {
	var ba BitArray
	ba.InitBitArray(3)
	ba.Set(5)
	assert.Equal(t, 5, ba.ToInt())
	//オーバーフロー
	ba.Set(10)
	assert.Equal(t, 2, ba.ToInt())
}

func TestAdd(t *testing.T) {
	var a BitArray
	a.InitBitArray(3)
	a.Set(2)
	var b BitArray
	b.InitBitArray(3)
	b.Set(3)
	assert.Equal(t, 5, a.Add(b).ToInt())
}

func TestASub(t *testing.T) {
	var a BitArray
	a.InitBitArray(3)
	a.Set(3)
	var b BitArray
	b.InitBitArray(3)
	b.Set(2)
	assert.Equal(t, 1, a.Sub(b).ToInt())
	//オーバーフロー
	assert.Equal(t, 7, b.Sub(a).ToInt())
}

func TestMul(t *testing.T) {
	var a BitArray
	a.InitBitArray(3)
	a.Set(2)
	var b BitArray
	b.InitBitArray(3)
	b.Set(2)
	assert.Equal(t, 4, a.Mul(b).ToInt())
}

func TestAssign(t *testing.T) {
	var a BitArray
	a.InitBitArray(3)
	a.Set(2)
	var b BitArray
	b.InitBitArray(3)
	b.Assign(a)
	assert.Equal(t, 2, b.ToInt())
}
