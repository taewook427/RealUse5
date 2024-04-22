// test641 : stdlib5.ksign hy

package main

// go mod init example.com
// go build -buildmode=c-shared -o mydll.dll myso.so
// func Free() {}, func main() {}

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"stdlib5/kobj"
	"stdlib5/ksign"
	"unsafe"
)

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

//export func0
func func0(arr0 *C.char) {
	// free : C char arr
	Free(arr0)
}

//export func1
func func1(arr0 *C.char, len0 C.int) *C.char {
	// khash : path -> hashv, error (pack/sendauto)
	v0, v1 := inf1(Recv(arr0, len0))
	return Sendauto(kobj.Pack([][]byte{v0, v1}))
}

func inf1(parm0 []byte) (ret0 []byte, ret1 []byte) {
	defer func() {
		err := recover()
		if err != nil {
			ret0 = nil
			ret1 = []byte(fmt.Sprintf("Critical Error : %s", err))
		}
	}()
	hashv := ksign.Khash(string(parm0))
	return hashv, nil
}

//export func2
func func2(arr0 *C.char, len0 C.int) *C.char {
	// kinfo : path -> size, file, folder, error (pack/sendauto)
	v0, v1, v2, v3 := inf2(Recv(arr0, len0))
	return Sendauto(kobj.Pack([][]byte{v0, v1, v2, v3}))
}

func inf2(parm0 []byte) (ret0 []byte, ret1 []byte, ret2 []byte, ret3 []byte) {
	defer func() {
		err := recover()
		if err != nil {
			ret0 = nil
			ret1 = nil
			ret2 = nil
			ret3 = []byte(fmt.Sprintf("Critical Error : %s", err))
		}
	}()
	v0, v1, v2 := ksign.Kinfo(string(parm0))
	return kobj.Encode(v0, 8), kobj.Encode(v1, 8), kobj.Encode(v2, 8), nil
}

//export func3
func func3(val0 C.int) *C.char {
	// genkey : bit n -> public, private, error (pack/sendauto)
	v0, v1, v2 := inf3(int(val0))
	return Sendauto(kobj.Pack([][]byte{v0, v1, v2}))
}

func inf3(parm0 int) (ret0 []byte, ret1 []byte, ret2 []byte) {
	defer func() {
		err := recover()
		if err != nil {
			ret0 = nil
			ret1 = nil
			ret2 = []byte(fmt.Sprintf("Critical Error : %s", err))
		}
	}()
	v0, v1, v2 := ksign.Genkey(parm0)
	v3 := ""
	if v2 != nil {
		v3 = fmt.Sprintf("Error : %s", v2)
	}
	return []byte(v0), []byte(v1), []byte(v3)
}

//export func4
func func4(arr0 *C.char, len0 C.int, arr1 *C.char, len1 C.int) *C.char {
	// sign : private, plain -> enc, error (pack/sendauto)
	v0, v1 := inf4(Recv(arr0, len0), Recv(arr1, len1))
	return Sendauto(kobj.Pack([][]byte{v0, v1}))
}

func inf4(parm0 []byte, parm1 []byte) (ret0 []byte, ret1 []byte) {
	defer func() {
		err := recover()
		if err != nil {
			ret0 = nil
			ret1 = []byte(fmt.Sprintf("Critical Error : %s", err))
		}
	}()
	v0, v1 := ksign.Sign(string(parm0), parm1)
	v2 := ""
	if v1 != nil {
		v2 = fmt.Sprintf("Error : %s", v1)
	}
	return v0, []byte(v2)
}

//export func5
func func5(arr0 *C.char, len0 C.int, arr1 *C.char, len1 C.int, arr2 *C.char, len2 C.int) *C.char {
	// verify : public, enc, plain -> SU/error (pack/sendauto)
	v0 := inf5(Recv(arr0, len0), Recv(arr1, len1), Recv(arr2, len2))
	return Sendauto(kobj.Pack([][]byte{v0}))
}

func inf5(parm0 []byte, parm1 []byte, parm2 []byte) (ret0 []byte) {
	defer func() {
		err := recover()
		if err != nil {
			ret0 = []byte(fmt.Sprintf("Critical Error : %s", err))
		}
	}()
	v0, v1 := ksign.Verify(string(parm0), parm1, parm2)
	if v1 == nil {
		if v0 {
			return []byte("P")
		} else {
			return []byte("F")
		}
	} else {
		return []byte(fmt.Sprintf("Error : %s", v1))
	}
}

func main() {}
