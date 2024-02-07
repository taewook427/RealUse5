package main

// test609 : kviewer assist
// go mod init example.com
// go build -buildmode=c-shared -o mydll.dll
// go build -buildmode=c-shared -o myso.so

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"stdlib5/ksign"
	"stdlib5/kzip"
	"strings"
	"unsafe"
)

//export ex0
func ex0(arr *C.char, l C.int) *C.char { // 파일폴더 정보 추출
	// 폴더만 작동, 8B 폴더수 8B 파일수 8B 크기
	path := string(imarr(arr, l))
	temp := subfunc0(path)
	res := encode(temp[0], 8)
	res = append(res, encode(temp[1], 8)...)
	res = append(res, encode(temp[2], 8)...)
	return exarr(res)
}

//export ex1
func ex1(arr *C.char, l C.int) *C.char { // khash 계산
	path := string(imarr(arr, l))
	result := exarr(in1(path))
	return result
}

//export ex2
func ex2(arr *C.char, l C.int, m C.int) *C.char { // kzip 제어
	mode := int(m)
	path := string(imarr(arr, l))
	var result []byte
	if mode == 0 { // input folder path -> kzip
		result = []byte(in2(path, true))
	} else { // input kzip path -> folder
		result = []byte(in2(path, false))
	}
	temp := encode(len(result), 4)
	temp = append(temp, result...)
	return exarr(temp)
}

//export ex3
func ex3(arr *C.char, l C.int, m C.int, div C.int) *C.char { // div 제어
	mode := int(m)
	path := string(imarr(arr, l))
	divs := int(div)
	var result []byte
	if mode == 0 { // input file path -> div
		result = []byte(in3(path, true, divs))
	} else { // input div.0 path -> file
		result = []byte(in3(path, false, divs))
	}
	temp := encode(len(result), 4)
	temp = append(temp, result...)
	return exarr(temp)
}

//export ex4
func ex4(arr *C.char) { // c ptr free
	C.free(unsafe.Pointer(arr))
}

// 폴더 정보 구하기 안정화 함수
func subfunc0(path string) (result []int) {
	defer func() {
		if err := recover(); err != nil {
			result = make([]int, 3)
		}
	}()
	path, _ = filepath.Abs(path)
	path = strings.Replace(path, "\\", "/", -1)
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	ret := make(chan []int)
	go in0(path, ret)
	temp := <-ret
	temp[0] = temp[0] + 1
	return temp
}

// 폴더 정보 구하기 하위함수
func in0(path string, ret chan []int) {
	defer func() {
		if err := recover(); err != nil {
			ret <- make([]int, 3)
		}
	}()
	// path는 폴더 표준절대경로, ret로는 폴더수, 파일수, 용량 전송
	fs, _ := ioutil.ReadDir(path)
	fonum := 0
	finum := 0
	sinum := 0
	folders := make([]string, 0)
	files := make([]string, 0)

	for _, r := range fs {
		if r.IsDir() {
			tdir := path + r.Name()
			if tdir[len(tdir)-1] != '/' {
				tdir = tdir + "/"
			}
			folders = append(folders, tdir)
		} else {
			files = append(files, path+r.Name())
		}
	}

	wait := make([]chan []int, len(folders))
	// r 변수의 잘못된 참조를 막기 위함. go ~는 마지막에 r의 값을 일괄 전송함. 불변 string로 바꿈.
	for i, r := range folders {
		wait[i] = make(chan []int) // chan의 초기화 필수
		go in0(r, wait[i])
	}

	fonum = len(folders)
	finum = len(files)
	for _, r := range files {
		fn, _ := os.Stat(r)
		sinum = sinum + int(fn.Size())
	}
	for _, r := range wait {
		temp := <-r
		fonum = fonum + temp[0]
		finum = finum + temp[1]
		sinum = sinum + temp[2]
	}
	new := []int{fonum, finum, sinum}
	ret <- new
}

// khash 하위함수
func in1(path string) (result []byte) {
	defer func() {
		if err := recover(); err != nil {
			result = make([]byte, 64)
		}
	}()
	result = ksign.Khash(path)
	return result
}

// kzip 하위함수
func in2(path string, mode bool) (result string) {
	defer func() {
		if err := recover(); err != nil {
			result = fmt.Sprint(err)
		}
	}()
	path, _ = filepath.Abs(path)
	path = strings.Replace(path, "\\", "/", -1)
	if path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	tbox := kzip.Init()
	if mode { // pack
		tbox.Folder = path
		tbox.Export = path[0:strings.LastIndex(path, "/")] + "/kzip5_result.webp"
		tbox.Zipfolder("webp")
	} else { // unpack
		tbox.Export = path[0:strings.LastIndex(path, "/")] + "/kzip5_result/"
		tbox.Unzip(path)
	}
	return "converted successfully"
}

// div 하위함수
func in3(path string, mode bool, div int) (result string) {
	defer func() {
		if err := recover(); err != nil {
			result = fmt.Sprint(err)
		}
	}()
	path, _ = filepath.Abs(path)
	path = strings.Replace(path, "\\", "/", -1)
	if path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	if mode { // pack
		fs, _ := os.Stat(path)
		size := int(fs.Size())
		num0 := size / div
		num1 := size % div
		f, _ := os.Open(path)
		defer f.Close()
		for i := 0; i < num0; i++ {
			t, _ := os.OpenFile(fmt.Sprintf("%s.%d", path, i), os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
			temp := make([]byte, div)
			f.Read(temp)
			t.Write(temp)
			t.Close()
		}
		if num1 != 0 {
			t, _ := os.OpenFile(fmt.Sprintf("%s.%d", path, num0), os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
			temp := make([]byte, num1)
			f.Read(temp)
			t.Write(temp)
			t.Close()
		}
	} else { // unpack
		base := path[0 : len(path)-2]
		num := 0
		flag := true
		for flag {
			_, err := os.Stat(fmt.Sprintf("%s.%d", base, num))
			if err == nil {
				num = num + 1
			} else {
				flag = false
			}
		}
		f, _ := os.OpenFile(base, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
		defer f.Close()
		for i := 0; i < num; i++ {
			temp, _ := ioutil.ReadFile(fmt.Sprintf("%s.%d", base, i))
			f.Write(temp)
		}
	}
	return "converted successfully"
}

func imarr(arr *C.char, l C.int) []byte { // 바이트열 가져오기
	// C 스타일의 정수 배열을 Go 슬라이스로 변환
	gs := (*[1 << 30]C.char)(unsafe.Pointer(arr))[:l:l]

	// 바이트 배열로 변환
	bs := make([]byte, len(gs))
	for i := 0; i < len(bs); i++ {
		bs[i] = byte(gs[i])
	}

	return bs
}

func exarr(arr []byte) *C.char { // 바이트열 내보내기
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

func encode(num int, length int) []byte { // 리틀엔디안 인코딩
	temp := make([]byte, length)
	for i := 0; i < length; i++ {
		temp[i] = byte(num % 256)
		num = num / 256
	}
	return temp
}

func decode(data []byte) int { // 리틀엔디안 디코딩
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

func main() {
	fmt.Println(subfunc0("C://Users//427ta//AppData//"))
}
