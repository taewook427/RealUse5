package ksign

// test565 : ksign standard (go)
// go mod init example.com

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/crypto/sha3"
)

// 파일 해시, 너무 큰 파일은 안됨. /절대경로
func hashfile(path string, ret chan *[]byte) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic("file read error")
	}
	var out []byte
	if len(dat) == 0 {
		out = make([]byte, 64)
	} else {
		h := sha3.New512()
		h.Write(dat)
		out = h.Sum(nil)
	}
	ret <- &out
}

// 폴더 해시, 너무 내부 폴더가 많으면 안됨. /절대경로(~/)
func hashfolder(path string, ret chan *[]byte) {
	ps, err := ioutil.ReadDir(path)
	if err != nil {
		panic("folder read error")
	}
	names := make([]string, 0)
	for _, i := range ps {
		if i.IsDir() {
			tdir := path + i.Name()
			if tdir[len(tdir)-1] != '/' {
				tdir = tdir + "/"
			}
			names = append(names, tdir)
		} else {
			names = append(names, path+i.Name())
		}
	}
	sort.Strings(names)
	var out []byte
	if len(names) == 0 {
		out = make([]byte, 64)
	} else {
		mem := make([][]byte, len(names))
		wait := make([]chan *[]byte, len(names))
		// r 변수의 잘못된 참조를 막기 위함. go ~는 마지막에 r의 값을 일괄 전송함. 불변 string로 바꿈.
		for i, r := range names {
			wait[i] = make(chan *[]byte)
			if r[len(r)-1] == byte('/') {
				go hashfolder(r, wait[i])
			} else {
				go hashfile(r, wait[i])
			}
		}
		for i, r := range wait {
			mem[i] = *(<-r)
		}
		var temp []byte
		for _, r := range mem {
			temp = append(temp, r...)
		}
		h := sha3.New512()
		h.Write(temp)
		out = h.Sum(nil)
	}
	ret <- &out
}

// file/folder path str -> 64B hash []byte
func Khash(path string) []byte {
	path, _ = filepath.Abs(path)
	path = strings.Replace(path, "\\", "/", -1)
	file, _ := os.Open(path)
	fileinfo, _ := file.Stat()
	ret := make(chan *[]byte)
	var out []byte
	if fileinfo.IsDir() {
		if path[len(path)-1] != '/' {
			path = path + "/"
		}
		go hashfolder(path, ret)
	} else {
		go hashfile(path, ret)
	}
	out = *(<-ret)
	return out
}

// PEM 형식화
func fmtpem(key []byte, keyType string) string {
	block := &pem.Block{
		Type:  keyType,
		Bytes: key,
	}

	pemKey := pem.EncodeToMemory(block)
	return string(pemKey)
}

// gen N byte public, private key (N * 8 bit) -> 2048 : 256, 4096 : 512
func Genkey(n int) (string, string) {
	key, err := rsa.GenerateKey(rand.Reader, 8*n)
	if err != nil {
		panic("keygen fail")
	}
	privateKey := fmtpem(x509.MarshalPKCS1PrivateKey(key), "RSA PRIVATE KEY")
	publicKey := fmtpem(x509.MarshalPKCS1PublicKey(&key.PublicKey), "PUBLIC KEY")
	return publicKey, privateKey
}

// private S + plain 80B -> enc B
func Sign(private string, plain []byte) []byte {
	if len(plain) != 80 {
		panic("PlainV should be 80B")
	}
	block, _ := pem.Decode([]byte(private))
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	hashed := sha3.Sum512([]byte(plain))
	signature, _ := rsa.SignPSS(rand.Reader, key, crypto.SHA3_512, hashed[:], nil)
	return signature
}

// public S + enc B + plain 80B -> T/F (True is ok)
func Verify(public string, enc []byte, plain []byte) bool {
	if len(plain) != 80 {
		panic("PlainV should be 80B")
	}
	block, _ := pem.Decode([]byte(public))
	key, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	hashed := sha3.Sum512([]byte(plain))
	err := rsa.VerifyPSS(key, crypto.SHA3_512, hashed[:], enc, nil)
	if err != nil {
		return false
	}
	return true
}

// name S + hashed 64B -> plain 80B
func Fm(name string, hashed []byte) []byte {
	nb := []byte(name + "                ")
	nb = nb[0:16]
	if len(hashed) != 64 {
		panic("HashV should be 64B")
	}
	return append(nb, hashed...)
}
