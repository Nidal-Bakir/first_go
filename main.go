package main

import (
	"encoding/binary"
	"fmt"
	"math/bits"

	// "math/bits"
	"unsafe"
)

var isBE bool
var isLE bool

func init() {
	var x uint16 = 0xFF00
	isBE = (*(*[2]byte)(unsafe.Pointer(&x)))[0] == 0xFF
	isLE = !isBE
}

func main() {
	// netData := []byte{0, 132, 95, 237, 80, 104, 111, 110, 101, 0, 0, 0, 0, 0, 1, 0}

	// bs := make([]byte, 2)
	// binary.LittleEndian.PutUint16(bs, 1)
	// fmt.Println(bs)
	// fmt.Println()

	// data := []byte{0x00, 0x01}
	// fmt.Println(*(*uint16)(unsafe.Pointer(&data)))
	// fmt.Println(binary.BigEndian.Uint16(data))
	// fmt.Println(binary.LittleEndian.Uint16(data))

	// var s uint16 = 300

	// lit := []byte{44, 1}
	big := []byte{1, 44}

	fmt.Println(binary.BigEndian.Uint16(big))
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, 300)
	fmt.Println(bs)

	// bs := make([]byte, 2)
	// bs2 := make([]byte, 2)
	// binary.NativeEndian.PutUint16(bs, s)
	// binary.BigEndian.PutUint16(bs2, s)
	// fmt.Println(bs)
	// fmt.Println(bs2)

	// fmt.Println(*(*[2]byte)(unsafe.Pointer(&bs)))

	// fmt.Println()

	// value := binary.BigEndian.Uint16(data)
	// fmt.Println(value)
	// fmt.Println(*(*[2]byte)(unsafe.Pointer(&value)))

	// fmt.Println()

}

type Data struct {
	Val    uint32   // 4 byte
	Label  [10]byte // 10 byte
	Active bool     // 1 byte
}

func NewDataFromNet(dataBytes [16]byte) *Data {
	data := new(Data)
	if isLE {
		data.Val = binary.BigEndian.Uint32(dataBytes[:4])
	} else {
		// data.Val = binary.ByteOrder.Uint32()
	}
	copy(data.Label[:], dataBytes[4:14])
	data.Active = dataBytes[14] != 0
	return data
}

func NewDataFromNetUsingUnsafe(dataBytes [16]byte) *Data {
	data := (*Data)(unsafe.Pointer(&dataBytes))
	if isLE {
		data.Val = bits.ReverseBytes32(data.Val)
	}
	return data
}

func DataToNetSafe(d Data) [16]byte {
	arr := [16]byte{}

	binary.BigEndian.PutUint32(arr[:4], d.Val)
	copy(arr[4:14], d.Label[:])
	if d.Active {
		arr[14] = 1
	}

	return arr
}

func DataToNetUnSafe(d Data) [16]byte {

	arr := *(*[16]byte)(unsafe.Pointer(&d))

	if isLE {

	}

	return arr

}
