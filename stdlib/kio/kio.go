// test627 : stdlib5.kio

package kio

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// wait t, print any
func Print(v interface{}, t float64) {
	time.Sleep(time.Duration(t*1000) * time.Millisecond)
	fmt.Print(v)
}

// ask q, get str input
func Input(q string) string {
	fmt.Print(q)
	temp, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if len(temp) == 0 {
		return ""
	} else if temp[len(temp)-1] == '\n' {
		temp = temp[0 : len(temp)-1]
	}
	if len(temp) == 0 {
		return ""
	} else if temp[len(temp)-1] == '\r' {
		temp = temp[0 : len(temp)-1]
	}
	return temp
}

// compare two []byte
func Bequal(a []byte, b []byte) bool {
	if len(a) == len(b) {
		for i, r := range a {
			if b[i] != r {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

// read stringed []byte ("8a9B" -> B)
func Bread(raw string) ([]byte, error) {
	raw = strings.ToLower(raw)
	if len(raw)%2 != 0 {
		return nil, errors.New("not even length")
	}
	table := "0123456789abcdef"
	out := make([]byte, len(raw)/2)
	var cmp byte
	var ptr byte
	for i := 0; i < len(raw)/2; i++ {
		cmp = raw[2*i]
		ptr = 0
		for cmp != table[ptr] {
			ptr = ptr + 1
			if ptr == 16 {
				return nil, errors.New("invalid character")
			}
		}
		out[i] = 16 * ptr

		cmp = raw[2*i+1]
		ptr = 0
		for cmp != table[ptr] {
			ptr = ptr + 1
			if ptr == 16 {
				return nil, errors.New("invalid character")
			}
		}
		out[i] = out[i] + ptr
	}
	return out, nil
}

// conv []byte to string (B -> "39a2")
func Bprint(raw []byte) string {
	table := "0123456789abcdef"
	out := make([]byte, len(raw)*2)
	for i, r := range raw {
		out[2*i] = table[r/16]
		out[2*i+1] = table[r%16]
	}
	return string(out)
}

// absolute path (folder : */, file : *)
func Abs(path string) string {
	path, _ = filepath.Abs(path)
	path = strings.Replace(path, "\\", "/", -1)
	file, _ := os.Open(path)
	defer file.Close()
	fileinfo, _ := file.Stat()
	if fileinfo.IsDir() {
		if path[len(path)-1] != '/' {
			path = path + "/"
		}
	}
	return path
}

// get file Size (-1 : not Exist)
func Size(path string) int {
	fileinfo, err := os.Stat(path)
	if err == nil {
		return int(fileinfo.Size())
	} else {
		return -1
	}
}

// file io pointer, "r"/"w"/"a"/"x"
func Open(path string, mode string) (*os.File, error) {
	var out *os.File
	var err error

	switch mode {
	case "r":
		out, err = os.Open(path)
	case "w":
		out, err = os.Create(path)
	case "a":
		out, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	case "x":
		out, err = os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	default:
		err = errors.New("invalid mode")
	}

	return out, err
}

// read bytes (-1 is readall)
func Read(f *os.File, size int) ([]byte, error) {
	if size < 0 {
		fileinfo, fserr := f.Stat()
		if fserr == nil {
			size = int(fileinfo.Size())
			tell, rderr := f.Seek(0, 1)
			if rderr == nil {
				size = size - int(tell)
			} else {
				return nil, rderr
			}
		} else {
			return nil, fserr
		}
	}
	success := 0
	tb := make([]byte, size)
	for i := 0; i < size/1073741824; i++ {
		temp, err := f.Read(tb[1073741824*i : 1073741824*i+1073741824])
		success = success + temp
		if err != nil {
			return tb[0:success], nil
		}
	}
	if size%1073741824 != 0 {
		temp, err := f.Read(tb[size-size%1073741824:])
		success = success + temp
		if err != nil {
			return tb[0:success], nil
		}
	}
	return tb, nil
}

// write bytes, returns not written bytes (pointer of original!!)
func Write(f *os.File, data []byte) ([]byte, error) {
	success := 0
	size := len(data)
	for i := 0; i < size/1073741824; i++ {
		temp, err := f.Write(data[1073741824*i : 1073741824*i+1073741824])
		success = success + temp
		if err != nil {
			return data[success:], err
		}
	}
	if size%1073741824 != 0 {
		temp, err := f.Write(data[size-size%1073741824:])
		success = success + temp
		if err != nil {
			return data[success:], err
		}
	}
	return nil, nil
}
