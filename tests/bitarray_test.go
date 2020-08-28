package test

import (
	"SrvControl/utils/bitarray"
	"fmt"
	"testing"
)

func TestBitarray(t *testing.T) {
	var b bitarray.BitArray
	b.InitBitArray(10)
	b.FromInt(5)

	fmt.Println("b的十進制:", b.ToInt())

	for i := 0; i < 10; i++ {
		fmt.Println(b.GetAt(i))
	}

	b.SetAt(1, true)
	fmt.Println("b.SetAt(1,true)后b的十進制:", b.ToInt())
	b.SetAt(10, true)
	fmt.Println("array len:", b.GetArrayLen())
	fmt.Println("aaaaaaaaaa")
}

func TestCreateBitArray(t *testing.T) {
	b := bitarray.CreateBitArray(10, 5)
	for i := 0; i < 10; i++ {
		fmt.Println(b.GetAt(i))
	}
}
