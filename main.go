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

	fmt.Println()

}

type Data struct {
	Val    uint32   // 4 byte
	Label  [10]byte // 10 byte
	Active bool     // 1 byte
}

func NewDataFromNet(dataBytes [16]byte) *Data {
	data := new(Data)
	data.Val = binary.BigEndian.Uint32(dataBytes[:4])
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
	if isLE {
		d.Val = bits.ReverseBytes32(d.Val)
	}
	return *(*[16]byte)(unsafe.Pointer(&d))
}
