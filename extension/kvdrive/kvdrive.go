// test739 : extension.kvdrive backend

package main

import (
	"bytes"
	"fmt"
	"stdlib5/kobj"
	"stdlib5/kvault"
	"strings"
	"unsafe"
)

/*
#include <stdlib.h>
*/
import "C"

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

var g4sh kvault.G4FSshell
var g5sh kvault.Shell
var isv4 bool
var cmdbuf []string

//export func0
func func0(arg0 *C.char) {
	// free : C char arr
	Free(arg0)
}

//export func1
func func1(arg0 C.int) {
	// set flag : 0 isv4 true, 1 isv4 false, 2 flagsz true, 3 flagsz false
	switch arg0 {
	case 0:
		isv4 = true
	case 1:
		isv4 = false
	case 2:
		g5sh.FlagSz = true
	case 3:
		g5sh.FlagSz = false
	}
}

//export func2
func func2(arg0 C.int) {
	// clear buf : 0 cmdbuf, 1 strbuf, 2 bytesbuf
	switch arg0 {
	case 0:
		cmdbuf = nil
	case 1:
		g4sh.IOstr = nil
		g5sh.IOstr = nil
	case 2:
		g4sh.IObyte = nil
		g5sh.IObyte = nil
	}
}

//export func3
func func3(arg0 *C.char, arg1 C.int, arg2 C.int) {
	// setter : add data (0 str, 1 bytes, 2 curpath, 3 cmdbuf)
	switch arg2 {
	case 0:
		if isv4 {
			g4sh.IOstr = append(g4sh.IOstr, string(Recv(arg0, arg1)))
		} else {
			g5sh.IOstr = append(g5sh.IOstr, string(Recv(arg0, arg1)))
		}
	case 1:
		if isv4 {
			g4sh.IObyte = append(g4sh.IObyte, Recv(arg0, arg1))
		} else {
			g5sh.IObyte = append(g5sh.IObyte, Recv(arg0, arg1))
		}
	case 2:
		if isv4 {
			g4sh.CurPath = string(Recv(arg0, arg1))
		} else {
			g5sh.CurPath = string(Recv(arg0, arg1))
		}
	case 3:
		cmdbuf = append(cmdbuf, string(Recv(arg0, arg1)))
	}
}

//export func4
func func4(arg0 C.int) *C.char {
	// getter : get data & clear (0 str, 1 bytes, 2 curpath, 3 asyncerr)
	var out []byte
	switch arg0 {
	case 0:
		data := make([][]byte, 0)
		if isv4 {
			for _, r := range g4sh.IOstr {
				data = append(data, []byte(r))
			}
			g4sh.IOstr = nil
		} else {
			for _, r := range g5sh.IOstr {
				data = append(data, []byte(r))
			}
			g5sh.IOstr = nil
		}
		out = kobj.Pack(data)
	case 1:
		if isv4 {
			out = kobj.Pack(g4sh.IObyte)
			g4sh.IObyte = nil
		} else {
			out = kobj.Pack(g5sh.IObyte)
			g5sh.IObyte = nil
		}
	case 2:
		if isv4 {
			out = []byte(g4sh.CurPath)
		} else {
			out = []byte(g5sh.CurPath)
		}
	case 3:
		if isv4 {
			out = []byte(g4sh.AsyncErr)
		} else {
			out = []byte(g5sh.AsyncErr)
		}
	}
	return Sendauto(out)
}

//export func5
func func5(arg0 C.int) C.int {
	// get status flag : 0 FlagWK, 1 FlagRo, 2 Flagsz, 3 isv4 -> (0 T, 1 F)
	tgt := false
	switch arg0 {
	case 0:
		if isv4 {
			tgt = g4sh.FlagWk
		} else {
			tgt = g5sh.FlagWk
		}
	case 1:
		if isv4 {
			tgt = g4sh.FlagRo
		} else {
			tgt = g5sh.FlagRo
		}
	case 2:
		if !isv4 {
			tgt = g5sh.FlagSz
		}
	case 3:
		tgt = isv4
	}
	if tgt {
		return 0
	} else {
		return 1
	}
}

//export func6
func func6(arg0 *C.char, arg1 C.int) *C.char {
	// cmd order, returns error
	cmd := string(Recv(arg0, arg1))
	var err error
	var data string
	if isv4 {
		err = g4sh.Command(cmd, cmdbuf)
	} else {
		err = g5sh.Command(cmd, cmdbuf)
	}
	if err == nil {
		data = ""
	} else {
		data = fmt.Sprint(err)
	}
	return Sendauto([]byte(data))
}

//export func7
func func7() *C.char {
	// get status self : name, time, size, lock, subdir, subfile
	out := make([][]byte, 6)
	if isv4 {
		out[0] = []byte(g4sh.CurInfo.Self_name)
		out[1] = []byte(g4sh.CurInfo.Self_time)
		out[2] = kobj.Encode(g4sh.CurInfo.Self_size, 8)
		if g4sh.CurInfo.Self_locked {
			out[3] = []byte{0}
		} else {
			out[3] = []byte{1}
		}
		out[4] = kobj.Encode(g4sh.CurInfo.Self_subdir, 8)
		out[5] = kobj.Encode(g4sh.CurInfo.Self_subfile, 8)
	} else {
		out[0] = []byte(g5sh.CurPath)
		out[1] = []byte("0000.00.00;00:00:00")
		out[2] = kobj.Encode(0, 8)
		out[3] = []byte{1}
		out[4] = kobj.Encode(g5sh.CurNum[0], 8)
		out[5] = kobj.Encode(g5sh.CurNum[1], 8)
	}
	return Sendauto(kobj.Pack(out))
}

//export func8
func func8() *C.char {
	// get status dir : 0 nameLF, 1 timeLF, 2 size8B, 3 lock1B
	out := make([][]byte, 4)
	data := make([][]byte, 0, 64)
	flag := make([]byte, 0, 64)
	if isv4 {
		for _, r := range g4sh.CurInfo.Dir_size {
			data = append(data, kobj.Encode(r, 8))
		}
		for _, r := range g4sh.CurInfo.Dir_locked {
			if r {
				flag = append(flag, 0)
			} else {
				flag = append(flag, 1)
			}
		}
		out[0] = []byte(strings.Join(g4sh.CurInfo.Dir_name, "\n"))
		out[1] = []byte(strings.Join(g4sh.CurInfo.Dir_time, "\n"))
		out[2] = bytes.Join(data, nil)
		out[3] = flag
	} else {
		for _, r := range g5sh.CurSize[:g5sh.CurNum[0]] {
			data = append(data, kobj.Encode(r, 8))
		}
		for _, r := range g5sh.CurLock[:g5sh.CurNum[0]] {
			if r {
				flag = append(flag, 0)
			} else {
				flag = append(flag, 1)
			}
		}
		out[0] = []byte(strings.Join(g5sh.CurName[:g5sh.CurNum[0]], "\n"))
		out[1] = []byte(strings.Join(g5sh.CurTime[:g5sh.CurNum[0]], "\n"))
		out[2] = bytes.Join(data, nil)
		out[3] = flag
	}
	return Sendauto(kobj.Pack(out))
}

//export func9
func func9() *C.char {
	// get status file : 0 nameLF, 1 timeLF, 2 size8B, 3 fptr8B
	out := make([][]byte, 4)
	data := make([][]byte, 0, 64)
	nums := make([][]byte, 0, 64)
	if isv4 {
		for _, r := range g4sh.CurInfo.File_size {
			data = append(data, kobj.Encode(r, 8))
		}
		for _, r := range g4sh.CurInfo.File_fptr {
			nums = append(nums, kobj.Encode(r, 8))
		}
		out[0] = []byte(strings.Join(g4sh.CurInfo.File_name, "\n"))
		out[1] = []byte(strings.Join(g4sh.CurInfo.File_time, "\n"))
		out[2] = bytes.Join(data, nil)
		out[3] = bytes.Join(nums, nil)
	} else {
		for _, r := range g5sh.CurSize[g5sh.CurNum[0]:] {
			data = append(data, kobj.Encode(r, 8))
		}
		for _, r := range g5sh.CurLock[g5sh.CurNum[0]:] {
			if r {
				nums = append(nums, kobj.Encode(0, 8))
			} else {
				nums = append(nums, kobj.Encode(1, 8))
			}
		}
		out[0] = []byte(strings.Join(g5sh.CurName[g5sh.CurNum[0]:], "\n"))
		out[1] = []byte(strings.Join(g5sh.CurTime[g5sh.CurNum[0]:], "\n"))
		out[2] = bytes.Join(data, nil)
		out[3] = bytes.Join(nums, nil)
	}
	return Sendauto(kobj.Pack(out))
}

func main() {

}
