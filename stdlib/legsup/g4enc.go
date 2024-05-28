// test663 : stdlib5.legsup gen4enc

package legsup

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"os"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strings"
	"sync"

	"golang.org/x/crypto/sha3"
)

// sha3-256 hash
func sha3256(data []byte) []byte {
	temp := sha3.New256()
	temp.Write(data)
	return temp.Sum(nil)
}

// aes-128 en, pad T : nB, F : 16nB
func aes128en(data []byte, key []byte, iv []byte, pad bool) []byte {
	if pad {
		plen := byte(16 - len(data)%16)
		var i byte
		for i = 0; i < plen; i++ {
			data = append(data, plen)
		}
	}
	if len(data)%16 != 0 || len(key) != 16 || len(iv) != 16 {
		return nil
	}
	temp := make([]byte, len(data))
	block, _ := aes.NewCipher(key)
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(temp, data)
	copy(iv, temp[len(temp)-16:])
	return temp
}

// aes-128 de, pad T : nB, F : 16nB
func aes128de(data []byte, key []byte, iv []byte, pad bool) []byte {
	if len(data)%16 != 0 || len(key) != 16 || len(iv) != 16 {
		return nil
	}
	temp := make([]byte, len(data))
	block, _ := aes.NewCipher(key)
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(temp, data)
	copy(iv, data[len(data)-16:])
	if pad {
		plen := temp[len(temp)-1]
		temp = temp[0 : len(temp)-int(plen)]
	}
	return temp
}

// gen4 enc copy file to fptr
func g4fcopy(f *os.File, path string) error {
	size := kio.Size(path)
	t, err := kio.Open(path, "r")
	if err == nil {
		defer t.Close()
	} else {
		return err
	}
	var temp []byte
	for i := 0; i < size/10485760; i++ {
		temp, err = kio.Read(t, 10485760)
		if err != nil {
			return err
		}
		_, err = kio.Write(f, temp)
		if err != nil {
			return err
		}
	}
	if size%10485760 != 0 {
		temp, err = kio.Read(t, size%10485760)
		if err != nil {
			return err
		}
		_, err = kio.Write(f, temp)
		if err != nil {
			return err
		}
	}
	return nil
}

// gen4 enc copy file from fptr
func g4fmake(f *os.File, path string, size int) error {
	t, err := kio.Open(path, "w")
	if err == nil {
		defer t.Close()
	} else {
		return err
	}
	var temp []byte
	for i := 0; i < size/10485760; i++ {
		temp, err = kio.Read(f, 10485760)
		if err != nil {
			return err
		}
		_, err = kio.Write(t, temp)
		if err != nil {
			return err
		}
	}
	if size%10485760 != 0 {
		temp, err = kio.Read(f, size%10485760)
		if err != nil {
			return err
		}
		_, err = kio.Write(t, temp)
		if err != nil {
			return err
		}
	}
	return nil
}

// gen4 enc zip files to ./tempkaesl
func g4dozip(files []string) error {
	f, err := kio.Open("./tempkaesl", "w")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	kio.Write(f, kobj.Encode(len(files), 2))
	for _, r := range files {
		r = kio.Abs(r)
		name := []byte(r[strings.LastIndex(r, "/")+1:])
		kio.Write(f, kobj.Encode(len(name), 2))
		kio.Write(f, name)
		kio.Write(f, kobj.Encode(kio.Size(r), 8))
		err = g4fcopy(f, r)
		if err != nil {
			return err
		}
	}
	return nil
}

// gen4 enc unzip files from ./tempkaesl
func g4unzip(path string) error {
	path = kio.Abs(path)
	if path[len(path)-1] != '/' {
		return errors.New("path should be folder")
	}
	f, err := kio.Open("./tempkaesl", "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	temp, err := kio.Read(f, 2)
	if err != nil {
		return err
	}
	num := kobj.Decode(temp)
	for i := 0; i < num; i++ {
		temp, err = kio.Read(f, 2)
		if err != nil {
			return err
		}
		temp, err = kio.Read(f, kobj.Decode(temp))
		if err != nil {
			return err
		}
		name := string(temp)
		temp, err = kio.Read(f, 8)
		if err != nil {
			return err
		}
		err = g4fmake(f, path+name, kobj.Decode(temp))
		if err != nil {
			return err
		}
	}
	return err
}

// gen4 enc key expand, 128B ckey -> 32x 16B keys
func g4expkey(ckey []byte) [][]byte {
	if len(ckey) != 128 {
		return nil
	}
	out := make([][]byte, 32)
	var pre []byte
	var sub []byte

	for i := 0; i < 16; i++ {
		ti := (7 * i) % 16
		if ti > 8 {
			pre = ckey[8*ti-64 : 8*ti]
			sub = append(ckey[8*ti:], ckey[0:8*ti-64]...)
		} else {
			pre = append(ckey[8*ti+64:], ckey[0:8*ti]...)
			sub = ckey[8*ti : 8*ti+64]
		}
		if len(pre) != 64 || len(sub) != 64 {
			return nil
		}

		temp := append(make([]byte, 0), sub...)
		for j := 0; j < 10000; j++ {
			temp = sha3256(append(append(make([]byte, 0), pre...), temp...))
		}
		out[i] = temp[0:16]
		out[i+16] = temp[16:32]
	}
	return out
}

// gen4 enc generate (pwhash, mkey) by pw, salt
func g4pwhmkey(pw []byte, salt []byte) ([]byte, []byte) {
	pwhash := append(make([]byte, 0), pw...)
	for i := 0; i < 100000; i++ {
		pwhash = sha3256(append(append(make([]byte, 0), salt...), pwhash...))
	}
	mkey := append(make([]byte, 0), pw...)
	for i := 0; i < 10000; i++ {
		mkey = sha3256(append(append(make([]byte, 0), mkey...), salt...))
	}
	return pwhash, mkey
}

// gen4 enc aes128 internal calc
func g4calc(chunk []byte, key []byte, iv []byte, isenc bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if isenc {
		copy(chunk, aes128en(chunk, key, iv, false))
	} else {
		copy(chunk, aes128de(chunk, key, iv, false))
	}
}

// gen4 enc file encrypt
func g4enc(before string, after string, key [][]byte, iv [][]byte, header []byte) error {
	f, err := kio.Open(before, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	t, err := kio.Open(after, "w")
	if err == nil {
		defer t.Close()
	} else {
		return err
	}
	_, err = kio.Write(t, header)
	if err != nil {
		return nil
	}

	size := kio.Size(before)
	num0 := size / 131072
	num1 := size % 131072
	var inbuf []byte = nil
	var exbuf []byte = nil
	var tmbuf []byte = nil
	var wg sync.WaitGroup

	if num0 > 0 {
		inbuf, _ = kio.Read(f, 131072)
		for i := 0; i < num0-1; i++ {
			wg.Add(1)
			go g4calc(inbuf, key[i%32], iv[i%32], true, &wg)
			kio.Write(t, exbuf)
			tmbuf, _ = kio.Read(f, 131072)
			wg.Wait()
			exbuf = inbuf
			inbuf = tmbuf
		}
		wg.Add(1)
		go g4calc(inbuf, key[(num0-1)%32], iv[(num0-1)%32], true, &wg)
		kio.Write(t, exbuf)
		tmbuf, _ = kio.Read(f, num1)
		wg.Wait()
		exbuf = inbuf
		inbuf = tmbuf
	} else {
		inbuf, _ = kio.Read(f, num1)
	}

	kio.Write(t, exbuf)
	kio.Write(t, aes128en(inbuf, key[num0%32], iv[num0%32], true))
	return nil
}

// gen4 enc file decrypt
func g4dec(before string, after string, key [][]byte, iv [][]byte, stpoint int) error {
	f, err := kio.Open(before, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	t, err := kio.Open(after, "w")
	if err == nil {
		defer t.Close()
	} else {
		return err
	}
	_, err = f.Seek(int64(stpoint), 0)
	if err != nil {
		return err
	}

	size := kio.Size(before) - stpoint
	num0 := size / 131072
	num1 := size % 131072
	if num1 == 0 {
		num0 = num0 - 1
		num1 = 131072
	}
	var inbuf []byte = nil
	var exbuf []byte = nil
	var tmbuf []byte = nil
	var wg sync.WaitGroup

	if num0 > 0 {
		inbuf, _ = kio.Read(f, 131072)
		for i := 0; i < num0-1; i++ {
			wg.Add(1)
			go g4calc(inbuf, key[i%32], iv[i%32], false, &wg)
			kio.Write(t, exbuf)
			tmbuf, _ = kio.Read(f, 131072)
			wg.Wait()
			exbuf = inbuf
			inbuf = tmbuf
		}
		wg.Add(1)
		go g4calc(inbuf, key[(num0-1)%32], iv[(num0-1)%32], false, &wg)
		kio.Write(t, exbuf)
		tmbuf, _ = kio.Read(f, num1)
		wg.Wait()
		exbuf = inbuf
		inbuf = tmbuf
	} else {
		inbuf, _ = kio.Read(f, num1)
	}

	kio.Write(t, exbuf)
	kio.Write(t, aes128de(inbuf, key[num0%32], iv[num0%32], true))
	return nil
}

// gen4 enc
type G4enc struct {
	Hint   string // pw hint
	salt   []byte
	pwhash []byte
	ckeydt []byte
	iv     []byte
}

// gen4 enc encrypt file
func (tbox *G4enc) Encrypt(files []string, pw []byte) (string, error) {
	tbox.salt = genrand(32)
	tbox.iv = genrand(16)
	ckey := genrand(128)
	var mkey []byte
	tbox.pwhash, mkey = g4pwhmkey(pw, tbox.salt)
	tbox.ckeydt = aes128en(ckey, mkey[16:32], mkey[0:16], false)
	keys := g4expkey(ckey)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = append(make([]byte, 0), tbox.iv...)
	}

	path := kio.Abs(files[0])
	newpath := fmt.Sprintf("%s/kaesl%d.ote", path[0:strings.LastIndex(path, "/")], kobj.Decode(genrand(2))%9000+1000)
	err := g4dozip(files)
	if err == nil {
		defer os.Remove("./tempkaesl")
	} else {
		return "", err
	}

	hint := []byte(tbox.Hint)
	header := append(append([]byte("OTE1"), kobj.Encode(len(hint), 2)...), hint...)
	header = append(append(append(append(header, tbox.salt...), tbox.pwhash...), tbox.ckeydt...), tbox.iv...)
	err = g4enc("./tempkaesl", newpath, keys, ivs, header)
	if err == nil {
		return newpath, nil
	} else {
		return "", err
	}
}

// gen4 enc view encfile
func (tbox *G4enc) View(path string) error {
	path = kio.Abs(path)
	if path[strings.LastIndex(path, "."):] != ".ote" {
		return errors.New("invalid extension")
	}
	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	temp, err := kio.Read(f, 4)
	if err == nil {
		if !kio.Bequal(temp, []byte("OTE1")) {
			return errors.New("InvalidFile")
		}
	} else {
		return err
	}
	temp, err = kio.Read(f, 2)
	if err != nil {
		return err
	}
	temp, err = kio.Read(f, kobj.Decode(temp))
	if err != nil {
		return err
	}
	tbox.Hint = string(temp)
	tbox.salt, err = kio.Read(f, 32)
	if err != nil {
		return err
	}
	tbox.pwhash, err = kio.Read(f, 32)
	if err != nil {
		return err
	}
	tbox.ckeydt, err = kio.Read(f, 128)
	if err != nil {
		return err
	}
	tbox.iv, err = kio.Read(f, 16)
	return err
}

// gen4 enc decrypt file
func (tbox *G4enc) Decrypt(path string, pw []byte) error {
	if len(tbox.salt) == 0 {
		return errors.New("should done View() first")
	}
	temp, mkey := g4pwhmkey(pw, tbox.salt)
	if !kio.Bequal(temp, tbox.pwhash) {
		return errors.New("InvalidPW")
	}
	ckey := aes128de(tbox.ckeydt, mkey[16:32], mkey[0:16], false)
	keys := g4expkey(ckey)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = append(make([]byte, 0), tbox.iv...)
	}

	err := g4dec(path, "./tempkaesl", keys, ivs, len([]byte(tbox.Hint))+214)
	if err == nil {
		defer os.Remove("./tempkaesl")
	} else {
		return err
	}
	path = kio.Abs(path)
	return g4unzip(path[0 : strings.LastIndex(path, "/")+1])
}
