// test629 : stdlib5.kobj

package kobj

import (
	"os"
	"path/filepath"
	"strings"
)

// re-alloc current path, returns cmd args
func Repath() []string {
	var args []string
	args = append(args, os.Args...)
	path, _ := filepath.Abs(args[0])
	path = strings.Replace(path, "\\", "/", -1)
	os.Chdir(path[0:strings.LastIndex(path, "/")])
	return args
}

// little endian encoding
func Encode(num int, length int) []byte {
	temp := make([]byte, length)
	for i := 0; i < length; i++ {
		temp[i] = byte(num % 256)
		num = num / 256
	}
	return temp
}

// little endian decoding
func Decode(data []byte) int {
	temp := 0
	for i, r := range data {
		if r != 0 {
			exp := 1
			for j := 0; j < i; j++ {
				exp = exp * 256
			}
			temp = temp + int(r)*exp
		}
	}
	return temp
}

// package series of bytes, 1B len + (8B size + nB data) * n
func Pack(series [][]byte) []byte {
	out := make([]byte, 1)
	out[0] = byte(len(series))
	for _, r := range series {
		out = append(out, Encode(len(r), 8)...)
		out = append(out, r...)
	}
	return out
}

// unpack packed B, 1B len + (8B size + nB data) * n
func Unpack(chunk []byte) [][]byte {
	length := int(chunk[0])
	out := make([][]byte, length)
	ptr := 1
	for i := 0; i < length; i++ {
		clen := Decode(chunk[ptr : ptr+8])
		ptr = ptr + 8
		out[i] = chunk[ptr : ptr+clen]
		ptr = ptr + clen
	}
	return out
}

/*
// go mod init example.com
// go build -buildmode=c-shared -o mydll.dll myso.so
// func Free() {}, func main() {}

/ *
#include <stdlib.h>
* /
import "C"
import "unsafe"

// (c_char_p, length int) -> []byte
func Recv(arr *C.char, length C.int) []byte {
	// C char array -> Go slice
	gs := (*[1 << 30]C.char)(unsafe.Pointer(arr))[:length:length]

	// convert to []byte
	bs := make([]byte, len(gs))
	for i := 0; i < len(bs); i++ {
		bs[i] = byte(gs[i])
	}

	return bs
}

// []byte -> cptr(nB data)
func Send(arr []byte) *C.char {
	length := len(arr)
	tb := make([]C.char, length)
	for i, r := range arr {
		tb[i] = C.char(r)
	}

	// make C char array
	na := (*C.char)(C.malloc(C.size_t(length) * C.sizeof_char))
	copy((*[1 << 30]C.char)(unsafe.Pointer(na))[:length:length], tb)

	return na
}

// []byte -> cptr(8B len + nB data)
func Sendauto(arr []byte) *C.char {
	temp := kobj.Encode(len(arr), 8)
	temp = append(temp, arr...)
	return Send(temp)
}

// free cptr, MUST make (export free func) at lib code!!
func Free(arr *C.char) {
	C.free(unsafe.Pointer(arr))
}
*/
