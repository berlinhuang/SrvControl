// github.com/Yamakatsu63/verilog2Go/src/variable/

package bitarray

import (
	"fmt"
	"log"
	"strconv"
)

type Bit struct {
	value bool
}

type BitArray struct {
	bits []Bit
	pos  PosedgeObserver // positive edge 上升沿觀察者
	neg  NegedgeObserver // negative edge 下降沿觀察者
}

func notify(a BitArray, b int) {
	if a.ToInt() < b && a.pos != nil {
		fmt.Println("notify posedge observer")
		a.NotifyPosedgeObserver()
	} else if a.ToInt() > b && a.neg != nil {
		a.NotifyNegedgeObserver()
	}
}

//InitBitArray 将BitArry初始化
func (ba *BitArray) InitBitArray(length int) {
	ba.bits = make([]Bit, length)
	for i := 0; i < length; i++ {
		bit := Bit{false}
		ba.bits[i] = bit
	}
}

// 创建指定长度bitarray并用十进制value进行初始化
func CreateBitArray(length int, value int) BitArray {
	var result BitArray
	result.InitBitArray(length)
	result.FromInt(value)
	return result
}

// 用十进制初始化bitarray
func (ba *BitArray) FromInt(value int) {
	notify(*ba, value)
	length := len(ba.bits)
	comparison := 1 << length
	for i := 1; i <= length; i++ {
		ba.bits[length-i].value = (((value << i) & comparison) >> length) == 1
	}
}
func (ba *BitArray) initFromTo(from, to int, value int) {
	length := from - to + 1
	comparison := 1 << length
	k := 0
	for i := to; i <= from; i++ {
		b := (((value << (k + 1)) & comparison) >> length) == 1
		idx := from - k
		ba.bits[idx].value = b
		k = k + 1
	}
}

// 用形如"A3F82D"字符串初始化bitarray
func (ba *BitArray) FromString(strValue string) {
	length := len(strValue)
	if strValue == "" || length <= 0 {
		ba.FromInt(0)
		return
	}
	if length%2 != 0 { //基数 前面补0
		strValue = "0" + strValue
		length = length + 1
	}
	if ba.GetArrayLen() < length*8 {
		log.Fatal("BitArray长度不够")
	}
	k := 0
	for i := length - 2; i >= 0; i = i - 2 {
		dec, err := strconv.ParseUint(strValue[i:i+2], 16, 8)
		if err != nil {
			log.Fatalln("ParseUint失败")
		}
		from := ba.GetArrayLen() - 1 - k*8
		to := ba.GetArrayLen() - 1 - k*8 - 7
		ba.initFromTo(from, to, int(dec))
		k = k + 1
	}
}

func (ba *BitArray) ToString() string {

	return ""
}

// 将bitarray转化为十进制并返回
func (ba BitArray) ToInt() int {
	ret := 0
	length := len(ba.bits)
	for i := length - 1; i >= 0; i-- {
		//bitの値がtrueの場合，対応する値を加算する
		//bit的值为真时，加上对应的值
		if ba.bits[i].value {
			ret |= 1 << i
		}
	}
	return ret
}

// 对指定位进行置位操作 ba.Set( 4, true )
func (ba BitArray) SetAt(index int, b bool) bool {
	if index < 0 || index >= len(ba.bits) {
		return false
	}
	ba.bits[index : index+1][0].value = b
	return true
}

// 判定指定位置是否置位 ba.GetAt( 5 )
func (ba BitArray) GetAt(index int) bool {
	if index < 0 || index >= len(ba.bits) {
		log.Fatal("溢出") //os.exit()
	}
	return ba.bits[index : index+1][0].value
}

// 获取bitarray的长度
func (ba BitArray) GetArrayLen() int {
	return len(ba.bits)
}

//Set はBitArrayのBitsに値をセットする
func (ba *BitArray) Set(value int) {
	notify(*ba, value)
	length := len(ba.bits)
	comparison := 1 << length
	for i := 1; i <= length; i++ {
		ba.bits[length-i].value = (((value << i) & comparison) >> length) == 1
	}
}

// Get はindexで指定したBitを持つBitArrayを返す
func (ba BitArray) Get(index int) BitArray {
	var result BitArray
	result.InitBitArray(1)
	//スライスを代入
	result.bits = ba.bits[index : index+1]
	return result
}

// Calc はvalueの値をもつBitArrayを返す
func (BitArray) Calc(value int, length int) BitArray {
	comparison := 1 << length
	var result BitArray
	result.InitBitArray(length)
	for i := 1; i <= length; i++ {
		result.bits[length-i].value = (((value << i) & comparison) >> length) == 1
	}
	return result
}

//Add はポート同士の加算を行う
func (ba BitArray) Add(input BitArray) BitArray {
	a := ba.ToInt()
	b := input.ToInt()
	length := len(ba.bits)
	var result BitArray
	result = result.Calc(a+b, length)
	notify(ba, result.ToInt())
	return result
}

//Sub はポート同士の減算を行う
func (ba BitArray) Sub(input BitArray) BitArray {
	a := ba.ToInt()
	b := input.ToInt()
	length := len(ba.bits)
	var result BitArray
	result = result.Calc(a-b, length)
	notify(ba, result.ToInt())
	return result
}

//Mul はポート同士の乗算を行う
func (ba BitArray) Mul(input BitArray) BitArray {
	a := ba.ToInt()
	b := input.ToInt()
	length := len(ba.bits)
	var result BitArray
	result = result.Calc(a*b, length)
	notify(ba, result.ToInt())
	return result
}

//Bitxor はポート同士の排他的論理和を返す
func (ba BitArray) Bitxor(input BitArray) BitArray {
	a := ba.ToInt()
	b := input.ToInt()
	length := len(ba.bits)
	var result BitArray
	result = result.Calc(a^b, length)
	notify(ba, result.ToInt())
	return result
}

//Bitand はポート同士の論理積を返す
func (ba BitArray) Bitand(input BitArray) BitArray {
	a := ba.ToInt()
	b := input.ToInt()
	length := len(ba.bits)
	var result BitArray
	result = result.Calc(a&b, length)
	notify(ba, result.ToInt())
	return result
}

//Bitor はポート同士の論理和を返す
func (ba BitArray) Bitor(input BitArray) BitArray {
	a := ba.ToInt()
	b := input.ToInt()
	length := len(ba.bits)
	var result BitArray
	result = result.Calc(a|b, length)
	notify(ba, result.ToInt())
	return result
}

// Assign は引数のBitArrayを割り当てる
func (ba *BitArray) Assign(result BitArray) {
	length := len(ba.bits)
	for i := 0; i < length; i++ {
		ba.bits[i].value = result.bits[i].value
	}
}
