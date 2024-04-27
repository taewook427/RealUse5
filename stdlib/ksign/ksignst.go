// test639 : stdlib5.ksign st

package ksign

// go get "golang.org/x/crypto/sha3"

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"sort"
	"stdlib5/kio"

	"golang.org/x/crypto/sha3"
)

// check if path is folder
func isdir(path string) bool {
	temp, err := os.Stat(path)
	if err == nil {
		return temp.IsDir()
	} else {
		return false
	}
}

// file pash (abspath)
func hashf(path string, ret chan []byte) {
	var v []byte
	defer func() {
		if err := recover(); err != nil {
			v = nil
		}
		ret <- v
	}()
	f, _ := kio.Open(path, "r")
	defer f.Close()
	temp, _ := kio.Read(f, -1)
	if len(temp) == 0 {
		v = make([]byte, 64)
	} else {
		h := sha3.New512()
		h.Write(temp)
		v = h.Sum(nil)
		temp = nil
		h = nil
	}
}

// folder hash (abspath/)
func hashd(path string, ret chan []byte) {
	var v []byte
	defer func() {
		if err := recover(); err != nil {
			v = nil
		}
		ret <- v
	}()
	fs, _ := os.ReadDir(path)
	nmb := make([][]byte, len(fs))
	for i, r := range fs {
		nms := path + r.Name()
		if isdir(nms) && nms[len(nms)-1] != '/' {
			nms = nms + "/"
		}
		nmb[i] = []byte(nms)
	}
	sort.Slice(nmb, func(i, j int) bool { return bytes.Compare(nmb[i], nmb[j]) < 0 })
	if len(nmb) == 0 {
		ret <- make([]byte, 64)
	} else {
		mem := make([]chan []byte, len(nmb))
		for i, r := range nmb {
			mem[i] = make(chan []byte, 1)
			nms := string(r)
			if nms[len(nms)-1] == '/' {
				go hashd(nms, mem[i])
			} else {
				go hashf(nms, mem[i])
			}
		}
		temp := make([]byte, 0, 64*len(nmb))
		for i := 0; i < len(nmb); i++ {
			temp = append(temp, <-mem[i]...)
		}
		if len(temp) == 64*len(nmb) {
			h := sha3.New512()
			h.Write(temp)
			v = h.Sum(nil)
		} else {
			v = nil
		}
	}
}

// get folder info (abspath/)
func infod(path string, ret chan int) {
	var size, folder, file int
	defer func() {
		if err := recover(); err != nil {
			size = 0
			folder = 0
			file = 0
		}
		ret <- size
		ret <- file
		ret <- folder
	}()
	fs, _ := os.ReadDir(path)
	nms := make([]string, len(fs))
	for i, r := range fs {
		temp := path + r.Name()
		if isdir(temp) && temp[len(temp)-1] != '/' {
			temp = temp + "/"
		}
		nms[i] = temp
	}
	mem := make([]chan int, 0)
	size = 0
	file = 0
	folder = 1
	for _, r := range nms {
		if r[len(r)-1] == '/' {
			com := make(chan int, 3)
			go infod(r, com)
			mem = append(mem, com)
		} else {
			size = size + kio.Size(r)
			file = file + 1
		}
	}
	for _, r := range mem {
		size = size + <-r
		file = file + <-r
		folder = folder + <-r
	}
}

// file/folder -> 64B khash value (general path)
func Khash(path string) []byte { // runtime.GC(), debug.FreeOSMemory()
	path = kio.Abs(path)
	ret := make(chan []byte, 1)
	if path[len(path)-1] == '/' {
		go hashd(path, ret)
	} else {
		go hashf(path, ret)
	}
	return <-ret
}

// get folder/file num, size info (general path) (size, file, folder)
func Kinfo(path string) (int, int, int) {
	path = kio.Abs(path)
	if path[len(path)-1] == '/' {
		ret := make(chan int, 3)
		go infod(path, ret)
		size := <-ret
		file := <-ret
		folder := <-ret
		return size, file, folder
	} else {
		return kio.Size(path), 1, 0
	}
}

// PEM formating
func fmtpem(key []byte, keyType string) string {
	block := &pem.Block{
		Type:  keyType,
		Bytes: key,
	}

	pemKey := pem.EncodeToMemory(block)
	return string(pemKey)
}

// gen N bit public, private key 2048/4096/8192
func Genkey(n int) (string, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, n)
	if err != nil {
		return "", "", err
	}
	pubASN1, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return "", "", err
	}
	priASN1 := x509.MarshalPKCS1PrivateKey(key)
	privateKey := fmtpem(priASN1, "PRIVATE KEY")
	publicKey := fmtpem(pubASN1, "PUBLIC KEY")
	return publicKey, privateKey, nil
}

// private S + plain nB -> enc B
func Sign(private string, plain []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(private))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	hashed := sha3.Sum512(plain)
	signature, err := rsa.SignPSS(rand.Reader, key, crypto.SHA3_512, hashed[:], &rsa.PSSOptions{SaltLength: 0, Hash: crypto.SHA3_512})
	return signature, err
}

// public S + enc B + plain nB -> T/F (True is ok)
func Verify(public string, enc []byte, plain []byte) (bool, error) {
	block, _ := pem.Decode([]byte(public))
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	rsaKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("not an RSA public key")
	}
	hashed := sha3.Sum512([]byte(plain))
	err = rsa.VerifyPSS(rsaKey, crypto.SHA3_512, hashed[:], enc, &rsa.PSSOptions{SaltLength: 0, Hash: crypto.SHA3_512})
	if err != nil {
		return false, nil
	}
	return true, nil
}

// format : PKIX (public), PKCS1 (private), PEM, PSS
