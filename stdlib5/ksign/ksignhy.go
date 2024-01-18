package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"example.com/ksign"
)

// test567 : ksign hybrid (go)
// go mod init example.com
// go build -buildmode=c-shared -o mydll.dll myso.so
// main.go ksign/ksign.go

func imarr(arr *C.char, l C.int) []byte {
	// C 스타일의 정수 배열을 Go 슬라이스로 변환
	gs := (*[1 << 30]C.char)(unsafe.Pointer(arr))[:l:l]

	// 바이트 배열로 변환
	bs := make([]byte, len(gs))
	for i := 0; i < len(bs); i++ {
		bs[i] = byte(gs[i])
	}

	return bs
}

func exarr(arr []byte) *C.char {
	l := len(arr)
	tb := make([]C.char, l)
	for i, r := range arr {
		tb[i] = C.char(r)
	}

	// 새로운 C 스타일의 정수 배열 생성
	na := (*C.char)(C.malloc(C.size_t(l) * C.sizeof_char))
	copy((*[1 << 30]C.char)(unsafe.Pointer(na))[:l:l], tb)

	return na
}

//export freehy
func freehy(arr *C.char) {
	C.free(unsafe.Pointer(arr))
}

//export khashhy
func khashhy(path *C.char, l C.int) *C.char {
	i := imarr(path, l)
	out := ksign.Khash(string(i))
	j := exarr(out)
	return j
}

//export genkeyhy
func genkeyhy(n C.int) *C.char {
	pu, pr := ksign.Genkey(int(n))
	public := []byte(pu)
	private := []byte(pr)
	temp := make([]byte, 0)
	temp = append(temp, byte(len(public)%256), byte(len(public)/256))
	temp = append(temp, public...)
	temp = append(temp, byte(len(private)%256), byte(len(private)/256))
	temp = append(temp, private...)
	o := exarr(temp)
	return o
}

//export signhy
func signhy(priv *C.char, l C.int, pln *C.char) *C.char {
	pr := imarr(priv, l)
	plain := imarr(pln, 80)
	private := string(pr)
	sgn := ksign.Sign(private, plain)
	o := exarr(sgn)
	return o
}

//export verifyhy
func verifyhy(arr0 *C.char, l0 C.int, arr1 *C.char, l1 C.int, arr2 *C.char) C.int {
	pu := imarr(arr0, l0)
	enc := imarr(arr1, l1)
	plain := imarr(arr2, 80)
	public := string(pu)
	if ksign.Verify(public, enc, plain) {
		return 1
	} else {
		return 0
	}
}

func main() {}
