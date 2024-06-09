// test637 : stdlib5.kaes hy

package main

// go get "golang.org/x/crypto/sha3"
// go get "golang.org/x/crypto/scrypt"

// go mod init example.com
// go build -buildmode=c-shared -o mydll.dll myso.so
// func Free() {}, func main() {}

/*
#include <stdlib.h>
*/
import "C"
import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"stdlib5/kaes"
	"stdlib5/kobj"
	"sync"
	"unsafe"

	"golang.org/x/crypto/scrypt"
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
func func0(arg0 *C.char) {
	// free : C char arr
	Free(arg0)
}

//export func1
func func1() *C.char {
	// kaes5 basickey
	return Sendauto(kaes.Basickey()) // mono sendauto
}

//export func2
func func2(arg0 C.int) C.float {
	// set/get proc, arg0 : 0 set, else get
	if int(arg0) == 0 {
		temp := -1.0
		Proc = &temp
		return C.float(temp)
	} else {
		if Proc == nil {
			return C.float(2.0)
		} else {
			return C.float(*Proc)
		}
	}
}

//export func3
func func3(arg0 *C.char, arg1 C.int, arg2 *C.char, arg3 C.int, arg4 *C.char, arg5 C.int) *C.char {
	// genpm, pw, kf, salt -> 224B known size output (pwh 128B mkey 96B)
	pw := Recv(arg0, arg1)
	kf := Recv(arg2, arg3)
	salt := Recv(arg4, arg5)
	return Send(inf3(pw, kf, salt)) // known size 224B send
}

func inf3(pw []byte, kf []byte, salt []byte) (out []byte) {
	defer func() {
		if err := recover(); err != nil {
			out = make([]byte, 224)
		}
	}()
	tb := append(append(append(append(pw, pw...), kf...), pw...), kf...)
	pwh, _ := scrypt.Key(tb, salt, 524288, 8, 1, 128)
	tb = append(append(append(append(kf, pw...), kf...), kf...), pw...)
	mkey, _ := scrypt.Key(tb, salt, 16384, 8, 1, 96)
	out = append(pwh, mkey...)
	return out
}

//export func4
func func4(arg0 *C.char, arg1 C.int, arg2 *C.char, arg3 C.int, arg4 C.int) (ret *C.char) {
	defer func() {
		if ferr := recover(); ferr != nil {
			ret = Send(nil)
		}
	}()
	// AEScalc, arg0/1 : data, arg2/3 : keys (iv 16B key 32B), arg4 : 0 enc dopad / 1 enc nopad / 2 dec dopad / 3 dec nopad
	data := Recv(arg0, arg1)
	keys := Recv(arg2, arg3)
	mode := int(arg4)
	isenc := true
	ispad := true
	if mode > 1 {
		isenc = false
	}
	if mode%2 == 1 {
		ispad = false
	}

	var out []byte
	block, _ := aes.NewCipher(keys[16:48])
	if isenc {
		encrypter := cipher.NewCBCEncrypter(block, keys[0:16])
		if ispad {
			plen := 16 - (len(data) % 16)
			for i := 0; i < plen; i++ {
				data = append(data, byte(plen))
			}
		}
		out = make([]byte, len(data))
		encrypter.CryptBlocks(out, data)

	} else {
		decrypter := cipher.NewCBCDecrypter(block, keys[0:16])
		out = make([]byte, len(data))
		decrypter.CryptBlocks(out, data)
		if ispad {
			plen := int(out[len(out)-1])
			out = out[0 : len(out)-plen]
		}
	}

	ret = Sendauto(out) // mono sendauto
	return ret
}

//export func5
func func5(arg0 *C.char, arg1 C.int, arg2 *C.char, arg3 C.int, arg4 C.int) (ret *C.char) {
	defer func() {
		if ferr := recover(); ferr != nil {
			ret = Send(make([]byte, 1920))
		}
	}()
	// AESchunk, arg0/1 : ckey 1920B, arg2/3 : data, arg4 : 0 enc / else dec, returns ckey + retdata
	ckey := Recv(arg0, arg1)
	reader := Recv(arg2, arg3)
	writer := make([]byte, 1920+len(reader))
	copy(writer[0:1920], ckey)
	isenc := true
	if int(arg4) != 0 {
		isenc = false
	}

	if len(ckey) == 1920 && len(reader)%524288 == 0 {
		var wg sync.WaitGroup
		wg.Add(40)
		for i := 0; i < 40; i++ {
			go inf5(ckey, reader, writer, isenc, i, &wg)
		}
		wg.Wait()
	}
	ret = Send(writer)
	return ret
}

func inf5(ckey []byte, reader []byte, writer []byte, isenc bool, num int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		recover()
	}()
	rept := len(reader) / 20971520
	ti := len(reader) % 20971520
	if ti/524288 > num {
		rept = rept + 1
	}
	ti = 48 * num
	pos0 := 524288*num - 20971520
	pos1 := pos0 + 1920

	block, _ := aes.NewCipher(ckey[ti+16 : ti+48])
	if isenc {
		encrypter := cipher.NewCBCEncrypter(block, ckey[ti:ti+16])
		for i := 0; i < rept; i++ {
			pos0 = pos0 + 20971520
			pos1 = pos0 + 1920
			encrypter.CryptBlocks(writer[pos1:pos1+524288], reader[pos0:pos0+524288])
		}
		copy(ckey[ti:ti+16], writer[pos1+524272:pos1+524288])
	} else {
		decrypter := cipher.NewCBCDecrypter(block, ckey[ti:ti+16])
		for i := 0; i < rept; i++ {
			pos0 = pos0 + 20971520
			pos1 = pos0 + 1920
			decrypter.CryptBlocks(writer[pos1:pos1+524288], reader[pos0:pos0+524288])
		}
		copy(ckey[ti:ti+16], reader[pos0+524272:pos0+524288])
	}
}

//export func6
func func6(arg0 *C.char, arg1 C.int) (ret *C.char) {
	defer func() {
		if ferr := recover(); ferr != nil {
			ftemp := make([][]byte, 2)
			ftemp[0] = []byte(fmt.Sprintf("critical error : %s", ferr))
			ret = Sendauto(kobj.Pack(ftemp)) // sendauto p2
		}
	}()
	Proc = &Wall.Proc
	// all-mode file encrypt, arg0/1 : p(hint, msg, pub, pri, pw, kf, path, pmode), returns p(err, encpath)
	parms := kobj.Unpack(Recv(arg0, arg1))
	Wall.Hint = string(parms[0])
	Wall.Msg = string(parms[1])
	Wall.Signkey[0] = string(parms[2])
	Wall.Signkey[1] = string(parms[3])
	pw := parms[4]
	kf := parms[5]
	path := string(parms[6])
	pmode := kobj.Decode(parms[7])
	path, err := Wall.EnFile(pw, kf, path, pmode)

	tb := make([][]byte, 2)
	if err != nil {
		tb[0] = []byte(err.Error())
	}
	tb[1] = []byte(path)
	ret = Sendauto(kobj.Pack(tb)) // sendauto p2
	return ret
}

//export func7
func func7(arg0 *C.char, arg1 C.int) (ret *C.char) {
	defer func() {
		if ferr := recover(); ferr != nil {
			ftemp := make([][]byte, 4)
			ftemp[0] = []byte(fmt.Sprintf("critical error : %s", ferr))
			ret = Sendauto(kobj.Pack(ftemp)) // sendauto p4
		}
	}()
	Proc = nil
	// all-mode file view, arg0/1 : path, returns p(err, hint, msg, pub)
	path := string(Recv(arg0, arg1))
	err := Wall.ViewFile(path)

	tb := make([][]byte, 4)
	if err != nil {
		tb[0] = []byte(err.Error())
	}
	tb[1] = []byte(Wall.Hint)
	tb[2] = []byte(Wall.Msg)
	tb[3] = []byte(Wall.Signkey[0])
	ret = Sendauto(kobj.Pack(tb)) // sendauto p4
	return ret
}

//export func8
func func8(arg0 *C.char, arg1 C.int) (ret *C.char) {
	defer func() {
		if ferr := recover(); ferr != nil {
			ftemp := make([][]byte, 2)
			ftemp[0] = []byte(fmt.Sprintf("critical error : %s", ferr))
			ret = Sendauto(kobj.Pack(ftemp)) // sendauto p2
		}
	}()
	Proc = &Wall.Proc
	// all-mode file decrypt, arg0/1 : p(pw, kf, path), returns p(err, decpath)
	parms := kobj.Unpack(Recv(arg0, arg1))
	pw := parms[0]
	kf := parms[1]
	path := string(parms[2])
	path, err := Wall.DeFile(pw, kf, path)

	tb := make([][]byte, 2)
	if err != nil {
		tb[0] = []byte(err.Error())
	}
	tb[1] = []byte(path)
	ret = Sendauto(kobj.Pack(tb)) // sendauto p2
	return ret
}

//export func9
func func9(arg0 *C.char, arg1 C.int) (ret *C.char) {
	defer func() {
		if ferr := recover(); ferr != nil {
			ret = Sendauto([]byte(fmt.Sprintf("critical error : %s", ferr))) // mono sendauto
		}
	}()
	Proc = &Wfunc.Proc
	// func-mode file ende, arg0/1 : p(before, after, akey, mode), returns err
	parms := kobj.Unpack(Recv(arg0, arg1))
	before := string(parms[0])
	after := string(parms[1])
	akey := parms[2]
	mode := kobj.Decode(parms[3]) // mode 0 : enc, else : dec

	Wfunc.Before.Open(before, true)
	defer Wfunc.Before.Close()
	Wfunc.After.Open(after, false)
	defer Wfunc.After.Close()
	var err error
	if mode == 0 {
		err = Wfunc.Encrypt(akey)
	} else {
		err = Wfunc.Decrypt(akey)
	}

	if err == nil {
		ret = Sendauto(nil) // mono sendauto
	} else {
		ret = Sendauto([]byte(err.Error()))
	}
	return ret
}

var Proc *float64
var Wall kaes.Allmode
var Wfunc kaes.Funcmode

func main() {}
