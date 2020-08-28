package bitarray

import (
	"fmt"
	"testing"
)

type TestModule struct {
	a BitArray
}

func (tm TestModule) Always() {
	fmt.Println("Posedge")
}

func TestPosedge(t *testing.T) {
	var ba BitArray
	ba.InitBitArray(3)
	ba.FromInt(3)
	var tm TestModule
	tm.a = ba
	tm.a.AddPosedgeObserver(tm)

	var b BitArray
	b.InitBitArray(3)
	b.FromInt(2)
	tm.a.Add(b)
}
