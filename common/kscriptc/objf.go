// test724 : common.kscriptc object

package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"stdlib5/kobj"
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
func func1(opt0 C.int, opt1 C.int, opt2 C.int, opt3 C.int, opt4 C.int) {
	// init compiler, set option (0 T, 1 F)
	cplr.init()
	cplr.option = [5]bool{opt0 == 0, opt1 == 0, opt2 == 0, opt3 == 0, opt4 == 0}
	exe = nil
	tpass = 0.0
}

//export func2
func func2(arr0 *C.char, len0 C.int, pos0 C.int) {
	// set data ( langc.data[0:4] )
	data := string(Recv(arr0, len0))
	switch pos0 {
	case 0:
		cplr.data[0] = data
	case 1:
		cplr.data[1] = data
	case 2:
		cplr.data[2] = data
	case 3:
		cplr.data[3] = data
	}
}

//export func3
func func3(pos0 C.int) *C.char {
	// get data ( langc.data[4:7], exe[-1] )
	var data []byte
	switch pos0 {
	case 4:
		data = []byte(cplr.data[4])
	case 5:
		data = []byte(cplr.data[5])
	case 6:
		data = []byte(cplr.data[6])
	case -1:
		data = exe
	}
	return Sendauto(data)
}

//export func4
func func4() C.float {
	// get tpass
	return C.float(tpass)
}

//export func5
func func5(arr0 *C.char, len0 C.int) *C.char {
	// addpkg
	data := string(Recv(arr0, len0))
	err := cplr.addpkg(data)
	if err == nil {
		data = ""
	} else {
		data = fmt.Sprintf("[header error] %s", err)
	}
	return Sendauto([]byte(data))
}

//export func6
func func6() *C.char {
	// compile
	var err error
	var data string
	exe, tpass, err = cplr.compile()
	if err == nil {
		data = ""
	} else {
		data = fmt.Sprintf("[compile error] %s", err)
	}
	return Sendauto([]byte(data))
}

var cplr langc
var exe []byte
var tpass float64

func main() {
}
