package main

// test615 : ST5adv example common extension
// go mod init example.com
// go build -buildmode=c-shared -o mydll.dll
// go build -buildmode=c-shared -o myso.so

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

//export primef
func primef(num C.int) *C.int {
	tgt := int(num)
	out := make([]int, 1)
	factor := 3
	if tgt < 0 {
		tgt = -tgt
	}

	switch tgt {
	case 0:
		out[0] = 0
	case 1:
		out[0] = 0
	default:
		for tgt != 1 {
			if tgt%2 == 0 {
				out = append(out, 2)
				tgt = tgt / 2
			} else {
				break
			}
		}
		for tgt != 1 {
			if tgt%factor == 0 {
				out = append(out, factor)
				tgt = tgt / factor
			} else {
				factor = factor + 2
			}
		}
		out[0] = len(out) - 1
	}

	clen := len(out)
	tempbi := make([]C.int, clen)
	for i, r := range out {
		tempbi[i] = C.int(r)
	}
	cnew := (*C.int)(C.malloc(C.size_t(clen) * C.sizeof_int))
	copy((*[1 << 30]C.int)(unsafe.Pointer(cnew))[:clen:clen], tempbi)

	return cnew
}

//export freef
func freef(arr *C.int) {
	C.free(unsafe.Pointer(arr))
}

func main() {
}
