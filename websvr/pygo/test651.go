package main

// go mod init example.com
// go build -buildmode=c-shared -o mydll.dll myso.so
// func Free() {}, func main() {}

/*
#include <stdlib.h>
*/
import "C"
import (
	"runtime/debug"
	"stdlib5/kobj"
	"time"
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

func longterm(logger *mblock) {
	logger.proc = 0.0
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		logger.proc = float64(i+1) / 10
	}
	logger.proc = 2.0
}

//export exf0
func exf0(p0 *C.char) {
	// 메모리 해제 함수
	Free(p0)
}

//export exf1
func exf1(p0 *C.char, p1 C.int) *C.char {
	// 메모리 사용량 검사 함수
	time.Sleep(3 * time.Second)
	i0 := Recv(p0, p1)
	time.Sleep(3 * time.Second)
	for i, r := range i0 {
		i0[i] = r + 1
	}
	time.Sleep(3 * time.Second)
	o0 := Send(i0)
	time.Sleep(3 * time.Second)
	i0 = nil
	debug.FreeOSMemory()
	time.Sleep(3 * time.Second)
	return o0
}

//export exf2
func exf2(p0 *C.char, p1 C.int) *C.char {
	// 메모리 복사시간 검사 함수
	i0 := Recv(p0, p1)
	return Send(i0)
}

//export exf3
func exf3(p0 *C.char, p1 C.int) {
	// 전역 메모리 저장 함수
	gmem.chunk = Recv(p0, p1)
}

//export exf4
func exf4() *C.char {
	// 전역 메모리 읽기 함수
	return Sendauto(gmem.chunk)
}

//export exf5
func exf5() {
	// 고루틴 할당 함수
	go longterm(&gmem)
}

//export exf6
func exf6() C.float {
	// 고루틴 진행도 체크 함수
	return C.float(gmem.proc)
}

type mblock struct {
	chunk []byte
	proc  float64
}

var gmem mblock

func main() {
	// dll 에선 실행 안되는듯
	gmem.chunk = make([]byte, 0)
	gmem.proc = -1.0
}
