package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"time"
	"unsafe"
)

// go build -buildmode=c-shared -o ex.dll

// (c_char_p, length int) -> []byte
func recv(arr *C.char, length C.int) []byte {
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
func send(arr []byte) *C.char {
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

// func recv(arr *C.char, length C.int) []byte {...}
// func send(arr []byte) *C.char {...}

//export freeptr
func freeptr(arr *C.char) {
	C.free(unsafe.Pointer(arr))
}

//export work
func work(parm0 *C.char, parm1 C.int) *C.char {
	time.Sleep(2 * time.Second)
	toproc := recv(parm0, parm1)

	time.Sleep(2 * time.Second)
	for i, r := range toproc {
		toproc[i] = r + 16
	}

	time.Sleep(2 * time.Second)
	buf := send(toproc)

	time.Sleep(2 * time.Second)
	return buf
}

func main() {}

/*
import (
	"fmt"
	"runtime/debug"
	"stdlib5/kio"
	"stdlib5/ksc"
	"time"
)

func read(path string) []byte {
	f, _ := kio.Open(path, "r")
	data, _ := kio.Read(f, -1)
	f.Close()
	return data
}

func checksum(data []byte) []byte {
	return ksc.Crc32hash(data)
}

// func read(path string) []byte {...}
// func checksum(data []byte) []byte {...}

func main() {
	for i := 0; i < 4; i++ {
		fmt.Println(checksum(read("./big.bin")))
		debug.FreeOSMemory()
		time.Sleep(2 * time.Second)
	}
}
*/
