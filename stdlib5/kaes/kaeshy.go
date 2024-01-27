package main

// test580 : kaes hy (go)
// go mod init example.com
// go build -buildmode=c-shared -o mydll.dll
// go build -buildmode=c-shared -o myso.so

/*
#include <stdlib.h>
*/
import "C"
import (
	"crypto/aes"
	"crypto/cipher"
	"os"
	"sync"
	"unsafe"

	"golang.org/x/crypto/scrypt"
)

// exf0 : free ptr
// exf1 : svkey
// exf2 : mkkey
// exf3 : file en
// exf4 : file de

//export exf0
func exf0(arr *C.char) {
	C.free(unsafe.Pointer(arr))
}

//export exf1
func exf1(pwa *C.char, pwl C.int, salta *C.char, saltl C.int, kfa *C.char, kfl C.int) *C.char {
	pw := imarr(pwa, pwl)
	salt := imarr(salta, saltl)
	kf := imarr(kfa, kfl)

	temp := make([]byte, 0)
	temp = append(temp, pw...)
	temp = append(temp, pw...)
	temp = append(temp, kf...)
	temp = append(temp, kf...)
	temp = append(temp, pw...)
	exb := hash(temp, salt, 524288, 8, 1, 256)

	return exarr(exb)
}

//export exf2
func exf2(pwa *C.char, pwl C.int, salta *C.char, saltl C.int, kfa *C.char, kfl C.int) *C.char {
	pw := imarr(pwa, pwl)
	salt := imarr(salta, saltl)
	kf := imarr(kfa, kfl)

	temp := make([]byte, 0)
	temp = append(temp, kf...)
	temp = append(temp, pw...)
	temp = append(temp, pw...)
	temp = append(temp, kf...)
	temp = append(temp, pw...)
	exb := hash(temp, salt, 16384, 8, 1, 48)

	return exarr(exb)
}

//export exf3
func exf3(heada *C.char, headl C.int, kchunk *C.char, beforea *C.char, beforel C.int, aftera *C.char, afterl C.int) {
	header := imarr(heada, headl)
	ckey := imarr(kchunk, 1536)
	before := string(imarr(beforea, beforel))
	after := string(imarr(aftera, afterl))

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = ckey[16*i : 16*i+16]
		keys[i] = ckey[512+32*i : 512+32*i+32]
	}

	inf2(header, keys, ivs, before, after)
}

//export exf4
func exf4(stp C.int, sizet *C.char, kchunk *C.char, beforea *C.char, beforel C.int, aftera *C.char, afterl C.int) {
	stpoint := int(stp)
	size := decode(imarr(sizet, 8))
	ckey := imarr(kchunk, 1536)
	before := string(imarr(beforea, beforel))
	after := string(imarr(aftera, afterl))

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = ckey[16*i : 16*i+16]
		keys[i] = ckey[512+32*i : 512+32*i+32]
	}

	inf3(stpoint, size, keys, ivs, before, after)
}

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

// little endian decoding
func decode(data []byte) int {
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

func hash(pw []byte, salt []byte, round int, word int, core int, length int) []byte { // scrypt hash
	out, _ := scrypt.Key(pw, salt, round, word, core, length)
	return out
}

func pad(data *[]byte) { // 16바이트 패딩
	padlen := 16 - (len(*data) % 16)
	for i := 0; i < padlen; i++ {
		*data = append(*data, byte(padlen))
	}
}

func unpad(data *[]byte) { // 16 바이트 언패딩
	padlen := int((*data)[len(*data)-1])
	*data = (*data)[0 : len(*data)-padlen]
}

func enshort(key []byte, iv []byte, data []byte) []byte { // short encryption no padding, 16B * n
	block, _ := aes.NewCipher(key)
	encrypter := cipher.NewCBCEncrypter(block, iv)
	out := make([]byte, len(data))
	encrypter.CryptBlocks(out, data)
	return out
}

func deshort(key []byte, iv []byte, data []byte) []byte { // short decryption no padding, 16B * n
	block, _ := aes.NewCipher(key)
	decrypter := cipher.NewCBCDecrypter(block, iv)
	out := make([]byte, len(data))
	decrypter.CryptBlocks(out, data)
	return out
}

func enchunk(key [][]byte, iv [][]byte, input []byte, output [][]byte, num int, wg *sync.WaitGroup) { // 16MB 청크에서 512kb 암호화
	defer wg.Done()
	block, _ := aes.NewCipher(key[num])
	encrypter := cipher.NewCBCEncrypter(block, iv[num])
	temp := 524288 * num
	tb := make([]byte, 524288)
	encrypter.CryptBlocks(tb, input[temp:temp+524288])
	output[num] = tb
	copy(iv[num], tb[524272:524288])
}

func dechunk(key [][]byte, iv [][]byte, input []byte, output [][]byte, num int, wg *sync.WaitGroup) { // 16MB 청크에서 512kb 복호화
	defer wg.Done()
	block, _ := aes.NewCipher(key[num])
	decrypter := cipher.NewCBCDecrypter(block, iv[num])
	temp := 524288 * num
	tb := make([]byte, 524288)
	decrypter.CryptBlocks(tb, input[temp:temp+524288])
	output[num] = tb
	copy(iv[num], input[temp+524272:temp+524288])
}

func inf2(header []byte, key [][]byte, iv [][]byte, before string, after string) { // 파일 암호화 내부함수
	fileinfo, _ := os.Stat(before)
	size := int(fileinfo.Size())
	f, _ := os.OpenFile(after, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666) // write as f
	defer f.Close()
	t, _ := os.Open(before) // read as t
	defer t.Close()
	num0 := size / 524288           // 512kb 청크 수
	num1 := size % 524288           // 남는 바이트 수
	buf0 := make([]byte, 16777216)  // 입력 버퍼
	buf0b := make([]byte, 16777216) // 전입력 버퍼
	buf1 := make([][]byte, 32)      // 출력 버퍼
	buf2 := make([][]byte, 32)      // 쓰기 버퍼
	for i := 0; i < 32; i++ {
		buf1[i] = make([]byte, 0)
	}

	var i int = 0
	num2 := num0 / 32
	f.Write(header)
	if num2 > 0 {
		t.Read(buf0b)
	}
	for i = 0; i < num2; i++ {
		copy(buf2, buf1)
		copy(buf0, buf0b)

		var wg sync.WaitGroup
		wg.Add(32)
		for j := 0; j < 32; j++ {
			go enchunk(key, iv, buf0, buf1, j, &wg)
		}

		for j := 0; j < 32; j++ {
			f.Write(buf2[j])
		}
		if i+1 < num2 {
			t.Read(buf0b)
		}
		wg.Wait()
	}

	num2 = num0 % 32
	if num2 == 0 {
		for j := 0; j < 32; j++ {
			f.Write(buf1[j])
		}
	} else {
		copy(buf2, buf1)
		buf0 = make([]byte, 524288*(num2))
		t.Read(buf0)
		var wg sync.WaitGroup
		wg.Add(num2)
		for j := 0; j < num2; j++ {
			go enchunk(key, iv, buf0, buf1, j, &wg)
		}
		for j := 0; j < 32; j++ {
			f.Write(buf2[j])
		}
		wg.Wait()
		for j := 0; j < num2; j++ {
			f.Write(buf1[j])
		}
	}
	buf0 = make([]byte, num1)
	t.Read(buf0)
	pad(&buf0)
	f.Write(enshort(key[num2], iv[num2], buf0))
}

func inf3(stpoint int, size int, key [][]byte, iv [][]byte, before string, after string) { // 파일 복호화 내부함수
	f, _ := os.OpenFile(after, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666) // write as f
	defer f.Close()
	t, _ := os.Open(before) // read as t
	defer t.Close()
	num0 := size / 524288 // 512kb 청크 수
	num1 := size % 524288 // 남는 바이트 수
	if num1 == 0 {
		num1 = 524288
		num0 = num0 - 1
	}
	buf0 := make([]byte, stpoint)   // 입력 버퍼
	buf0b := make([]byte, 16777216) // 전입력 버퍼
	buf1 := make([][]byte, 32)      // 출력 버퍼
	buf2 := make([][]byte, 32)      // 쓰기 버퍼
	for i := 0; i < 32; i++ {
		buf1[i] = make([]byte, 0)
	}
	t.Read(buf0)
	buf0 = make([]byte, 16777216)

	var i int = 0
	num2 := num0 / 32
	if num2 > 0 {
		t.Read(buf0b)
	}
	for i = 0; i < num2; i++ {
		copy(buf2, buf1)
		copy(buf0, buf0b)

		var wg sync.WaitGroup
		wg.Add(32)
		for j := 0; j < 32; j++ {
			go dechunk(key, iv, buf0, buf1, j, &wg)
		}

		for j := 0; j < 32; j++ {
			f.Write(buf2[j])
		}
		if i+1 < num2 {
			t.Read(buf0b)
		}
		wg.Wait()
	}

	num2 = num0 % 32
	if num2 == 0 {
		for j := 0; j < 32; j++ {
			f.Write(buf1[j])
		}
	} else {
		copy(buf2, buf1)
		buf0 = make([]byte, 524288*(num2))
		t.Read(buf0)
		var wg sync.WaitGroup
		wg.Add(num2)
		for j := 0; j < num2; j++ {
			go dechunk(key, iv, buf0, buf1, j, &wg)
		}
		for j := 0; j < 32; j++ {
			f.Write(buf2[j])
		}
		wg.Wait()
		for j := 0; j < num2; j++ {
			f.Write(buf1[j])
		}
	}
	buf0 = make([]byte, num1)
	t.Read(buf0)
	tb := deshort(key[num2], iv[num2], buf0)
	unpad(&tb)
	f.Write(tb)
}

func main() {

}
