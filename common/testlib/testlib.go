// test715 : common.testlib

package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"stdlib5/kobj"
	"strings"
	"unsafe"
)

// (c_char_p, length int) -> []byte
func Recv(arr *C.char, length C.int) []byte {
	return C.GoBytes(unsafe.Pointer(arr), length)
}

// []byte -> cptr(nB data)
func Send(arr []byte) *C.char {
	return (*C.char)(C.CBytes(arr))
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

//export func0
func func0(arr0 *C.char) {
	// free : C char arr
	Free(arr0)
}

//export func1
func func1(arr0 *C.char) *C.char {
	// get 8B uint64
	tgt := kobj.Decode(Recv(arr0, 8))
	prime := make([]string, 0)
	if tgt < 2 {
		prime = append(prime, fmt.Sprint(tgt))
	} else {
		div := 2
		for div < tgt {
			if tgt%div == 0 {
				prime = append(prime, fmt.Sprint(div))
				tgt = tgt / div
			} else {
				div = div + 1
			}
		}
		prime = append(prime, fmt.Sprint(div))
	}
	return Sendauto([]byte(strings.Join(prime, ", ")))
}

func main() {
}
