package kaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"example.com/kdb"
	"example.com/picdt"
	"golang.org/x/crypto/scrypt"
)

// test579 : kaes st (go)

// ========== ========== ksc5 start ========== ==========

// little endian encoding
func encode(num int, length int) *[]byte {
	temp := make([]byte, length)
	for i := 0; i < length; i++ {
		temp[i] = byte(num % 256)
		num = num / 256
	}
	return &temp
}

// little endian decoding
func decode(data *[]byte) int {
	temp := 0
	for i, r := range *data {
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

// crc32 check
func crc32check(data *[]byte) *[]byte {
	temp := int(crc32.ChecksumIEEE(*data))
	return encode(temp, 4)
}

// find 1024nB + KSC5 pos
func findpos(data *[]byte) int {
	temp := 0
	cmp := []byte("KSC5")
	for len(*data) >= temp+4 {
		d := (*data)[temp : temp+4]
		if bequal(&cmp, &d) {
			return temp
		} else {
			temp = temp + 1024
		}
	}
	return -1
}

// compare two []byte
func bequal(a *[]byte, b *[]byte) bool {
	if len(*a) != len(*b) {
		return false
	}
	for i, r := range *a {
		if r != (*b)[i] {
			return false
		}
	}
	return true
}

// ========== ========== ksc5 end ========== ==========

// ========== ========== kaes5 start ========== ==========

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

func hash(pw []byte, salt []byte, round int, word int, core int, length int) []byte { // scrypt hash
	out, _ := scrypt.Key(pw, salt, round, word, core, length)
	return out
}

func mkkey(pw []byte, salt []byte, kf []byte) []byte { // generate mkey 48B
	temp := make([]byte, 0)
	temp = append(temp, kf...)
	temp = append(temp, pw...)
	temp = append(temp, pw...)
	temp = append(temp, kf...)
	temp = append(temp, pw...)
	return hash(temp, salt, 16384, 8, 1, 48)
}

func svkey(pw []byte, salt []byte, kf []byte) []byte { // generate key storage 256B
	temp := make([]byte, 0)
	temp = append(temp, pw...)
	temp = append(temp, pw...)
	temp = append(temp, kf...)
	temp = append(temp, kf...)
	temp = append(temp, pw...)
	return hash(temp, salt, 524288, 8, 1, 256)
}

func inf0(header []byte, key [][]byte, iv [][]byte, data []byte) []byte { // 바이트 암호화 내부함수
	size := len(data)
	var out bytes.Buffer
	num0 := size / 524288          // 512kb 청크 수
	num1 := size % 524288          // 남는 바이트 수
	buf0 := make([]byte, 16777216) // 입력 버퍼
	buf1 := make([][]byte, 32)     // 출력 버퍼
	buf2 := make([][]byte, 32)     // 쓰기 버퍼
	for i := 0; i < 32; i++ {
		buf1[i] = make([]byte, 0)
	}

	var i int = 0
	num2 := num0 / 32
	for i = 0; i < num2; i++ {
		copy(buf2, buf1)
		temp := 16777216 * i
		buf0 = data[temp : temp+16777216]

		var wg sync.WaitGroup
		wg.Add(32)
		for j := 0; j < 32; j++ {
			go enchunk(key, iv, buf0, buf1, j, &wg)
		}

		for j := 0; j < 32; j++ {
			out.Write(buf2[j])
		}
		wg.Wait()
	}

	num2 = num0 % 32
	if num2 == 0 {
		for j := 0; j < 32; j++ {
			out.Write(buf1[j])
		}
	} else {
		copy(buf2, buf1)
		buf0 = data[16777216*i : 16777216*i+524288*(num0%32)]
		var wg sync.WaitGroup
		wg.Add(num2)
		for j := 0; j < num2; j++ {
			go enchunk(key, iv, buf0, buf1, j, &wg)
		}
		for j := 0; j < 32; j++ {
			out.Write(buf2[j])
		}
		wg.Wait()
		for j := 0; j < num2; j++ {
			out.Write(buf1[j])
		}
	}
	buf0 = data[size-num1:]
	pad(&buf0)
	out.Write(enshort(key[num2], iv[num2], buf0))

	header = append(header, out.Bytes()...)
	return header
}

func inf1(key [][]byte, iv [][]byte, data []byte) []byte { // 바이트 복호화 내부함수
	size := len(data)
	var out bytes.Buffer
	num0 := size / 524288 // 512kb 청크 수
	num1 := size % 524288 // 남는 바이트 수
	if num1 == 0 {
		num1 = 524288
		num0 = num0 - 1
	}
	buf0 := make([]byte, 16777216) // 입력 버퍼
	buf1 := make([][]byte, 32)     // 출력 버퍼
	buf2 := make([][]byte, 32)     // 쓰기 버퍼
	for i := 0; i < 32; i++ {
		buf1[i] = make([]byte, 0)
	}

	var i int = 0
	num2 := num0 / 32
	for i = 0; i < num2; i++ {
		copy(buf2, buf1)
		temp := 16777216 * i
		buf0 = data[temp : temp+16777216]

		var wg sync.WaitGroup
		wg.Add(32)
		for j := 0; j < 32; j++ {
			go dechunk(key, iv, buf0, buf1, j, &wg)
		}

		for j := 0; j < 32; j++ {
			out.Write(buf2[j])
		}
		wg.Wait()
	}

	num2 = num0 % 32
	if num2 == 0 {
		for j := 0; j < 32; j++ {
			out.Write(buf1[j])
		}
	} else {
		copy(buf2, buf1)
		buf0 = data[16777216*i : 16777216*i+524288*(num0%32)]
		var wg sync.WaitGroup
		wg.Add(num2)
		for j := 0; j < num2; j++ {
			go dechunk(key, iv, buf0, buf1, j, &wg)
		}
		for j := 0; j < 32; j++ {
			out.Write(buf2[j])
		}
		wg.Wait()
		for j := 0; j < num2; j++ {
			out.Write(buf1[j])
		}
	}
	buf0 = data[size-num1:]
	tb := deshort(key[num2], iv[num2], buf0)
	unpad(&tb)
	out.Write(tb)

	return out.Bytes()
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

func inf4(size int) int { // 패딩 후 사이즈 반환
	return size + 16 - size%16
}

// compare two slice
func sequal(a []byte, b []byte) bool {
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

// 안전한 난수 nB 반환
func Genrandom(size int) []byte {
	temp := make([]byte, size)
	rand.Read(temp)
	return temp
}

// 키 파일 경로에 따라 키 파일 바이트 생성
func Genkf(path string) []byte {
	temp, err := ioutil.ReadFile(path)
	if err != nil {
		temp = basickey()
	}
	return temp
}

type genbytes struct {
	Vaild bool
	Noerr bool
	Mode  string
	Msg   string
}

// genbytes init
func Init0() genbytes {
	var temp genbytes
	temp.Vaild = true
	temp.Noerr = false
	temp.Mode = "webp"
	temp.Msg = ""
	return temp
}

// pw B, kf B, hint B, data B -> enc B
func (self *genbytes) En(pw []byte, kf []byte, hint []byte, data []byte) []byte {
	salt := Genrandom(32)                            // salt 32B
	pwhash := svkey(pw, salt, kf)                    // pwhash 256B
	mkey := mkkey(pw, salt, kf)                      // master key 48B
	ckey := Genrandom(1536)                          // content key 1536B
	ckeydt := enshort(mkey[16:48], mkey[0:16], ckey) // content key data 1536B
	pw = make([]byte, 64)
	kf = make([]byte, 64)
	mkey = make([]byte, 64)

	mold := "mode = 0\nmsg = 0\nsalt = 0\npwhash = 0\nhint = 0\nckeydt = 0\n"
	kdbtbox := kdb.Init()
	kdbtbox.Readstr(&mold)
	tn := "mode"
	kdbtbox.Fixdata(&tn, "bytes")
	tn = "msg"
	kdbtbox.Fixdata(&tn, self.Msg)
	tn = "salt"
	kdbtbox.Fixdata(&tn, salt)
	tn = "pwhash"
	kdbtbox.Fixdata(&tn, pwhash)
	tn = "hint"
	kdbtbox.Fixdata(&tn, hint)
	tn = "ckeydt"
	kdbtbox.Fixdata(&tn, ckeydt)
	mh := []byte(*kdbtbox.Writestr()) // main header
	mhs := *encode(len(mh), 4)        // main header size

	var fakeh []byte
	switch self.Mode {
	case "png":
		fakeh = *picdt.Ka5png()
		tn := (16384 - len(fakeh)) % 1024
		tb := make([]byte, tn)
		fakeh = append(fakeh, tb...)
	case "webp":
		fakeh = *picdt.Ka5webp()
		tn := (16384 - len(fakeh)) % 1024
		tb := make([]byte, tn)
		fakeh = append(fakeh, tb...)
	default:
		fakeh = make([]byte, 0) // prehead + padding
	}
	commonh := []byte("KSC5")  // common head
	subtypeh := []byte("KAES") // subtype head
	res := *crc32check(&mh)    // reserved
	chunksize := *encode(inf4(len(data)), 8)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = ckey[16*i : 16*i+16]
		keys[i] = ckey[512+32*i : 512+32*i+32]
	}

	hd := append(fakeh, commonh...)
	hd = append(hd, subtypeh...)
	hd = append(hd, res...)
	hd = append(hd, mhs...)
	hd = append(hd, mh...)
	hd = append(hd, chunksize...)
	content := inf0(hd, keys, ivs, data)
	return content
}

// pw B, kf B, data B, stpoint N -> plain B
func (self *genbytes) De(pw []byte, kf []byte, data []byte, stpoint int) []byte {
	tb := data[stpoint : stpoint+4]
	mhs := decode(&tb)
	mh := data[stpoint+4 : stpoint+4+mhs] // main header
	tb = data[stpoint+4+mhs : stpoint+12+mhs]
	chunksize := decode(&tb)
	data = data[stpoint+12+mhs : stpoint+12+mhs+chunksize]

	kdbtbox := kdb.Init()
	ts := string(mh)
	kdbtbox.Readstr(&ts)
	ts = "salt"
	temp := kdbtbox.Getdata(&ts)
	salt := temp.Dat5
	ts = "pwhash"
	temp = kdbtbox.Getdata(&ts)
	pwhash := temp.Dat5
	ts = "ckeydt"
	temp = kdbtbox.Getdata(&ts)
	ckeydt := temp.Dat5
	if !sequal(pwhash, svkey(pw, salt, kf)) {
		panic("Not Valid PWKF")
	}

	mkey := mkkey(pw, salt, kf)                      // master key 48B
	ckey := deshort(mkey[16:48], mkey[0:16], ckeydt) // content key 1536B
	pw = make([]byte, 64)
	kf = make([]byte, 64)
	mkey = make([]byte, 64)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = ckey[16*i : 16*i+16]
		keys[i] = ckey[512+32*i : 512+32*i+32]
	}

	content := inf1(keys, ivs, data)
	return content
}

// data B -> hint B, msg S, stpoint N
func (self *genbytes) View(data []byte) ([]byte, string, int) {
	var pos int
	var tb []byte
	if len(data) > 16384 {
		tb = data[0:16384]
		pos = findpos(&tb)
	} else {
		pos = findpos(&data)
	}
	if pos == -1 {
		panic("Not Valid KSC5")
	}

	subtypeh := data[pos+4 : pos+8] // subtype head
	res := data[pos+8 : pos+12]     // reserved
	tb = data[pos+12 : pos+16]
	mhs := decode(&tb)
	mh := data[pos+16 : pos+16+mhs] // main header

	if !self.Noerr {
		if !sequal(subtypeh, []byte("KAES")) {
			panic("Not Valid KAES")
		}
		if !sequal(res, *crc32check(&mh)) {
			panic("Broken Header")
		}
	}

	kdbtbox := kdb.Init()
	ts := string(mh)
	kdbtbox.Readstr(&ts)
	ts = "hint"
	temp := *kdbtbox.Getdata(&ts)
	hint := temp.Dat5
	ts = "msg"
	temp = *kdbtbox.Getdata(&ts)
	msg := temp.Dat6
	return hint, msg, pos + 12
}

type genfile struct {
	Vaild bool
	Noerr bool
	Mode  string
	Msg   string
}

// genfile init
func Init1() genfile {
	var temp genfile
	temp.Vaild = true
	temp.Noerr = false
	temp.Mode = "webp"
	temp.Msg = ""
	return temp
}

// pw B, kf B, hint B, path str -> new path S
func (self *genfile) En(pw []byte, kf []byte, hint []byte, path string) string {
	path, _ = filepath.Abs(path)
	path = strings.Replace(path, "\\", "/", -1)        // abs path
	fopath := path[0 : strings.LastIndex(path, "/")+1] // folder path
	name := path[strings.LastIndex(path, "/")+1:]      // true name
	var tgt string                                     // write path
	tr := make([]byte, 2)
	rand.Read(tr)
	if self.Mode == "" {
		tgt = fopath + fmt.Sprintf("%03d%03d.ke5", tr[0], tr[1])
	} else {
		tgt = fopath + fmt.Sprintf("%03d%03d.%s", tr[0], tr[1], self.Mode)
	}
	nmb := []byte(name)
	pad(&nmb)

	salt := Genrandom(32)                            // salt 32B
	pwhash := svkey(pw, salt, kf)                    // pwhash 256B
	mkey := mkkey(pw, salt, kf)                      // master key 48B
	tkey := Genrandom(48)                            // title key 48B
	tkeydt := enshort(mkey[16:48], mkey[0:16], tkey) // title key data 48B
	namedt := enshort(tkey[16:48], tkey[0:16], nmb)  // name data
	ckey := Genrandom(1536)                          // content key 1536B
	ckeydt := enshort(mkey[16:48], mkey[0:16], ckey) // content key data 1536B
	pw = make([]byte, 64)
	kf = make([]byte, 64)
	mkey = make([]byte, 64)

	mold := "mode = 0\nmsg = 0\nsalt = 0\npwhash = 0\nhint = 0\ntkeydt = 0\nckeydt = 0\nnamedt = 0\n"
	kdbtbox := kdb.Init()
	kdbtbox.Readstr(&mold)
	tn := "mode"
	kdbtbox.Fixdata(&tn, "file")
	tn = "msg"
	kdbtbox.Fixdata(&tn, self.Msg)
	tn = "salt"
	kdbtbox.Fixdata(&tn, salt)
	tn = "pwhash"
	kdbtbox.Fixdata(&tn, pwhash)
	tn = "hint"
	kdbtbox.Fixdata(&tn, hint)
	tn = "tkeydt"
	kdbtbox.Fixdata(&tn, tkeydt)
	tn = "ckeydt"
	kdbtbox.Fixdata(&tn, ckeydt)
	tn = "namedt"
	kdbtbox.Fixdata(&tn, namedt)
	mh := []byte(*kdbtbox.Writestr()) // main header
	mhs := *encode(len(mh), 4)        // main header size

	var fakeh []byte
	switch self.Mode {
	case "png":
		fakeh = *picdt.Ka5png()
		tn := (16384 - len(fakeh)) % 1024
		tb := make([]byte, tn)
		fakeh = append(fakeh, tb...)
	case "webp":
		fakeh = *picdt.Ka5webp()
		tn := (16384 - len(fakeh)) % 1024
		tb := make([]byte, tn)
		fakeh = append(fakeh, tb...)
	default:
		fakeh = make([]byte, 0) // prehead + padding
	}
	commonh := []byte("KSC5")  // common head
	subtypeh := []byte("KAES") // subtype head
	res := *crc32check(&mh)    // reserved
	f, _ := os.Stat(path)
	chunksize := *encode(inf4(int(f.Size())), 8)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = ckey[16*i : 16*i+16]
		keys[i] = ckey[512+32*i : 512+32*i+32]
	}

	hd := append(fakeh, commonh...)
	hd = append(hd, subtypeh...)
	hd = append(hd, res...)
	hd = append(hd, mhs...)
	hd = append(hd, mh...)
	hd = append(hd, chunksize...)
	inf2(hd, keys, ivs, path, tgt)
	return tgt
}

// pw B, kf B, path str, stpoint N -> new path S
func (self *genfile) De(pw []byte, kf []byte, path string, stpoint int) string {
	path, _ = filepath.Abs(path)
	path = strings.Replace(path, "\\", "/", -1)        // abs path
	fopath := path[0 : strings.LastIndex(path, "/")+1] // folder path
	f, _ := os.Open(path)
	defer f.Close()
	tb := make([]byte, stpoint)
	f.Read(tb)
	tb = make([]byte, 4)
	f.Read(tb)
	mhs := decode(&tb)
	mh := make([]byte, mhs) // main header
	f.Read(mh)
	tb = make([]byte, 8)
	f.Read(tb)
	chunksize := decode(&tb)
	stpoint = stpoint + mhs + 12

	kdbtbox := kdb.Init()
	ts := string(mh)
	kdbtbox.Readstr(&ts)
	ts = "salt"
	temp := kdbtbox.Getdata(&ts)
	salt := temp.Dat5
	ts = "pwhash"
	temp = kdbtbox.Getdata(&ts)
	pwhash := temp.Dat5
	ts = "tkeydt"
	temp = kdbtbox.Getdata(&ts)
	tkeydt := temp.Dat5
	ts = "ckeydt"
	temp = kdbtbox.Getdata(&ts)
	ckeydt := temp.Dat5
	ts = "namedt"
	temp = kdbtbox.Getdata(&ts)
	namedt := temp.Dat5
	if !sequal(pwhash, svkey(pw, salt, kf)) {
		panic("Not Valid PWKF")
	}

	mkey := mkkey(pw, salt, kf)                      // master key 48B
	tkey := deshort(mkey[16:48], mkey[0:16], tkeydt) // title key 48B
	ckey := deshort(mkey[16:48], mkey[0:16], ckeydt) // content key 1536B
	nmb := deshort(tkey[16:48], tkey[0:16], namedt)
	unpad(&nmb)
	name := string(nmb)
	pw = make([]byte, 64)
	kf = make([]byte, 64)
	mkey = make([]byte, 64)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = ckey[16*i : 16*i+16]
		keys[i] = ckey[512+32*i : 512+32*i+32]
	}

	tgt := fopath + name
	inf3(stpoint, chunksize, keys, ivs, path, tgt)
	return tgt
}

// path str -> hint B, msg S, stpoint N
func (self *genfile) View(path string) ([]byte, string, int) {
	f, _ := os.Open(path)
	t, _ := os.Stat(path)
	size := t.Size()
	var tb []byte
	if size > 16384 {
		tb = make([]byte, 16384)
	} else {
		tb = make([]byte, size)
	}
	f.Read(tb)
	pos := findpos(&tb)
	f.Close()
	if pos == -1 {
		panic("Not Valid KSC5")
	}

	f, _ = os.Open(path)
	defer f.Close()
	tb = make([]byte, pos+4)
	f.Read(tb)
	subtypeh := make([]byte, 4) // subtype head
	f.Read(subtypeh)
	res := make([]byte, 4) // reserved
	f.Read(res)
	tb = make([]byte, 4)
	f.Read(tb)
	mhs := decode(&tb)
	mh := make([]byte, mhs) // main header
	f.Read(mh)

	if !self.Noerr {
		if !sequal(subtypeh, []byte("KAES")) {
			panic("Not Valid KAES")
		}
		if !sequal(res, *crc32check(&mh)) {
			panic("Broken Header")
		}
	}

	kdbtbox := kdb.Init()
	ts := string(mh)
	kdbtbox.Readstr(&ts)
	ts = "hint"
	temp := *kdbtbox.Getdata(&ts)
	hint := temp.Dat5
	ts = "msg"
	temp = *kdbtbox.Getdata(&ts)
	msg := temp.Dat6
	return hint, msg, pos + 12
}

type funcbytes struct {
	Vaild bool
}

// funcbytes init
func Init2() funcbytes {
	var temp funcbytes
	temp.Vaild = true
	return temp
}

// key 48B, data B -> enc B
func (self *funcbytes) En(key []byte, data []byte) []byte {
	ckey := Genrandom(1536)                        // content key 1536B
	ckeydt := enshort(key[16:48], key[0:16], ckey) // content key data 1536B
	key = make([]byte, 64)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = ckey[16*i : 16*i+16]
		keys[i] = ckey[512+32*i : 512+32*i+32]
	}

	content := inf0(ckeydt, keys, ivs, data)
	return content
}

// key 48B, data B -> plain B
func (self *funcbytes) De(key []byte, data []byte) []byte {
	ckeydt := data[0:1536]
	ckey := deshort(key[16:48], key[0:16], ckeydt) // content key 1536B
	key = make([]byte, 64)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = ckey[16*i : 16*i+16]
		keys[i] = ckey[512+32*i : 512+32*i+32]
	}

	data = data[1536:]
	content := inf1(keys, ivs, data)
	return content
}

type funcfile struct {
	Vaild bool
}

// funcfile init
func Init3() funcfile {
	var temp funcfile
	temp.Vaild = true
	return temp
}

// key 48B, before -> after
func (self *funcfile) En(key []byte, before string, after string) {
	ckey := Genrandom(1536)                        // content key 1536B
	ckeydt := enshort(key[16:48], key[0:16], ckey) // content key data 1536B
	key = make([]byte, 64)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = ckey[16*i : 16*i+16]
		keys[i] = ckey[512+32*i : 512+32*i+32]
	}

	inf2(ckeydt, keys, ivs, before, after)
}

// key 48B, before -> after
func (self *funcfile) De(key []byte, before string, after string) {
	f, _ := os.Open(before)
	ckeydt := make([]byte, 1536)
	f.Read(ckeydt)
	f.Close()
	ckey := deshort(key[16:48], key[0:16], ckeydt) // content key 1536B
	key = make([]byte, 64)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		ivs[i] = ckey[16*i : 16*i+16]
		keys[i] = ckey[512+32*i : 512+32*i+32]
	}

	fileinfo, _ := os.Stat(before)
	size := int(fileinfo.Size()) - 1536
	inf3(1536, size, keys, ivs, before, after)
}

// ========== ========== kaes5 end ========== ==========

func basickey() []byte {
	var temp []byte
	temp = append(temp, 234, 183, 184, 235, 158, 152, 44, 32, 235, 130, 152, 235, 165, 188, 32, 234, 176, 128, 235, 145, 172, 235, 145, 148, 32, 236, 177, 132, 32, 236, 157, 180, 234, 179, 179, 234, 185, 140, 236, 167)
	temp = append(temp, 128, 32, 236, 158, 172, 235, 176, 140, 235, 138, 148, 32, 236, 151, 172, 236, 160, 149, 32, 235, 179, 180, 235, 131, 136, 235, 139, 136, 63, 13, 10, 235, 172, 180, 235, 132, 136, 236, 160, 184)
	temp = append(temp, 235, 157, 188, 46, 13, 10, 236, 134, 159, 236, 149, 132, 235, 157, 188, 46, 13, 10, 236, 154, 184, 235, 160, 164, 235, 157, 188, 46, 13, 10, 236, 152, 155, 235, 130, 160, 236, 157, 152, 32)
	temp = append(temp, 234, 176, 144, 234, 176, 129, 235, 147, 164, 236, 157, 180, 32, 235, 143, 140, 236, 149, 132, 236, 152, 164, 235, 138, 148, 234, 181, 172, 235, 130, 152, 46, 13, 10, 236, 157, 180, 235, 178, 136)
	temp = append(temp, 236, 151, 148, 32, 236, 138, 164, 236, 138, 164, 235, 161, 156, 236, 157, 152, 32, 237, 158, 152, 236, 156, 188, 235, 161, 156, 32, 235, 130, 152, 235, 165, 188, 32, 235, 167, 137, 236, 157, 132)
	temp = append(temp, 32, 236, 136, 152, 32, 236, 158, 136, 234, 178, 160, 235, 139, 136, 63, 13, 10, 237, 157, 169, 236, 150, 180, 236, 167, 128, 234, 177, 176, 235, 157, 188, 46, 13, 10, 236, 154, 148, 236, 160)
	temp = append(temp, 149, 235, 147, 164, 236, 157, 180, 236, 151, 172, 46, 13, 10, 235, 168, 184, 235, 166, 172, 235, 165, 188, 32, 236, 134, 141, 236, 157, 184, 32, 236, 177, 132, 32, 236, 157, 180, 32, 234, 179)
	temp = append(temp, 179, 236, 151, 144, 236, 132, 156, 32, 236, 157, 180, 235, 159, 176, 32, 235, 139, 185, 235, 143, 140, 237, 149, 156, 32, 236, 167, 147, 236, 157, 132, 32, 235, 152, 144, 32, 235, 139, 164, 236)
	temp = append(temp, 139, 156, 32, 235, 178, 140, 236, 157, 180, 234, 179, 160, 32, 236, 158, 136, 234, 181, 172, 235, 130, 152, 46, 13, 10, 236, 167, 145, 236, 134, 141, 46, 13, 10, 236, 151, 180, 236, 135, 160)
	temp = append(temp, 32, 236, 157, 145, 236, 182, 149, 46, 13, 10, 234, 176, 156, 235, 176, 169, 46, 13, 10, 236, 157, 180, 32, 235, 170, 184, 236, 157, 128, 32, 236, 157, 180, 235, 159, 176, 32, 237, 158, 152)
	temp = append(temp, 235, 143, 132, 32, 236, 147, 184, 32, 236, 136, 152, 32, 236, 158, 136, 234, 181, 172, 235, 130, 152, 46, 32, 237, 157, 165, 235, 175, 184, 235, 161, 173, 236, 167, 128, 235, 167, 140, 32, 235)
	temp = append(temp, 175, 184, 236, 149, 189, 237, 149, 180, 46, 13, 10, 236, 157, 180, 235, 178, 136, 236, 151, 148, 32, 234, 183, 184, 235, 133, 128, 236, 157, 152, 32, 235, 143, 132, 236, 155, 128, 32, 236, 151)
	temp = append(temp, 134, 236, 157, 180, 32, 235, 167, 137, 236, 149, 132, 235, 179, 180, 235, 160, 164, 235, 172, 180, 235, 130, 152, 46, 13, 10, 235, 182, 132, 236, 132, 157, 46, 32, 236, 149, 149, 236, 182, 149)
	temp = append(temp, 46, 32, 236, 160, 132, 234, 176, 156, 46, 13, 10, 235, 130, 180, 32, 236, 149, 158, 236, 151, 144, 32, 236, 132, 156, 236, 167, 128, 32, 235, 167, 144, 234, 177, 176, 235, 157, 188, 46, 13)
	temp = append(temp, 10, 236, 157, 180, 32, 235, 170, 184, 236, 157, 128, 32, 236, 160, 156, 236, 149, 189, 236, 157, 180, 32, 235, 132, 136, 235, 172, 180, 32, 235, 167, 142, 236, 149, 132, 46, 13, 10, 235, 172)
	temp = append(temp, 180, 235, 132, 136, 236, 160, 184, 235, 130, 180, 235, 160, 164, 235, 157, 188, 46, 13, 10, 235, 130, 152, 236, 152, 164, 234, 177, 176, 235, 157, 188, 46, 13, 10, 235, 130, 160, 235, 155, 176)
	temp = append(temp, 234, 177, 176, 235, 157, 188, 46, 13, 10, 236, 157, 180, 234, 179, 179, 236, 157, 152, 32, 237, 158, 152, 236, 157, 132, 32, 235, 141, 148, 32, 235, 168, 188, 236, 160, 128, 32, 236, 149, 140)
	temp = append(temp, 236, 149, 152, 235, 139, 164, 235, 169, 180, 32, 236, 154, 176, 235, 166, 172, 235, 147, 164, 235, 143, 132, 32, 234, 183, 184, 32, 237, 158, 152, 236, 157, 132, 32, 236, 147, 184, 32, 236, 136)
	temp = append(temp, 152, 32, 236, 158, 136, 236, 151, 136, 234, 178, 160, 236, 167, 128, 46, 13, 10, 235, 130, 152, 235, 165, 188, 32, 235, 132, 152, 236, 167, 128, 32, 235, 170, 187, 237, 149, 152, 235, 169, 180)
	temp = append(temp, 32, 234, 178, 176, 234, 181, 173, 32, 235, 152, 144, 32, 235, 139, 164, 236, 139, 156, 32, 235, 168, 184, 235, 166, 172, 236, 151, 144, 32, 235, 176, 159, 237, 158, 144, 32, 235, 191, 144, 236)
	temp = append(temp, 157, 180, 235, 158, 128, 235, 139, 164, 46, 13, 10, 235, 172, 184, 236, 157, 132, 32, 236, 151, 180, 236, 150, 180, 236, 163, 188, 235, 167, 136, 46, 13, 10, 236, 157, 180, 234, 179, 179, 236)
	temp = append(temp, 151, 144, 32, 237, 149, 168, 234, 187, 152, 32, 234, 176, 128, 235, 157, 188, 236, 149, 137, 236, 158, 144, 234, 181, 172, 235, 130, 152, 46, 13, 10, 236, 157, 180, 32, 234, 176, 144, 236, 152)
	temp = append(temp, 165, 236, 151, 144, 236, 132, 156, 32, 235, 130, 152, 234, 176, 132, 235, 139, 164, 32, 237, 149, 152, 235, 141, 148, 235, 157, 188, 235, 143, 132, 32, 237, 152, 188, 236, 158, 144, 236, 132, 156)
	temp = append(temp, 32, 235, 172, 180, 236, 151, 135, 236, 157, 132, 32, 237, 149, 160, 32, 236, 136, 152, 32, 236, 158, 136, 236, 157, 132, 32, 234, 178, 131, 32, 234, 176, 153, 235, 139, 136, 63, 13, 10, 236)
	temp = append(temp, 157, 180, 32, 234, 181, 180, 235, 160, 136, 235, 165, 188, 32, 235, 129, 138, 235, 138, 148, 235, 139, 164, 32, 237, 149, 152, 235, 141, 148, 235, 157, 188, 235, 143, 132, 32, 236, 158, 160, 236)
	temp = append(temp, 139, 156, 235, 191, 144, 236, 157, 180, 236, 167, 128, 46, 13, 10, 236, 152, 133, 236, 150, 180, 236, 167, 128, 235, 138, 148, 234, 181, 172, 235, 130, 152, 46, 13, 10, 236, 158, 160, 236, 157)
	temp = append(temp, 180, 32, 236, 152, 164, 235, 138, 148, 234, 181, 172, 235, 130, 152, 46, 13, 10, 236, 154, 176, 235, 166, 172, 235, 165, 188, 32, 235, 178, 151, 236, 150, 180, 235, 130, 160, 32, 236, 136, 152)
	temp = append(temp, 235, 138, 148, 32, 236, 151, 134, 235, 139, 168, 235, 139, 164, 46, 13, 10, 234, 184, 176, 237, 154, 140, 235, 165, 188, 32, 235, 134, 147, 236, 179, 164, 234, 181, 172, 235, 130, 152, 46, 13)
	temp = append(temp, 10, 234, 181, 189, 237, 158, 136, 236, 167, 128, 32, 235, 170, 187, 237, 150, 136, 234, 181, 172, 235, 130, 152, 46, 13, 10, 237, 152, 188, 236, 158, 144, 236, 132, 156, 235, 138, 148, 32, 237)
	temp = append(temp, 157, 144, 235, 166, 132, 236, 157, 132, 32, 235, 169, 136, 236, 182, 156, 32, 236, 136, 152, 32, 236, 151, 134, 235, 139, 168, 235, 139, 164, 46, 13, 10, 237, 140, 140, 235, 143, 132, 235, 138)
	temp = append(temp, 148, 32, 235, 139, 164, 236, 139, 156, 32, 236, 157, 188, 235, 160, 129, 236, 157, 188, 32, 234, 178, 131, 236, 157, 180, 235, 158, 128, 235, 139, 164, 46, 13, 10, 235, 132, 136, 236, 157, 152)
	temp = append(temp, 32, 235, 175, 184, 236, 136, 153, 237, 149, 168, 236, 157, 180, 235, 139, 164, 46, 13, 10, 234, 177, 176, 235, 140, 128, 237, 149, 156, 32, 237, 157, 144, 235, 166, 132, 236, 157, 132, 32, 236)
	temp = append(temp, 134, 144, 235, 176, 148, 235, 139, 165, 236, 156, 188, 235, 161, 156, 32, 235, 167, 137, 236, 157, 132, 32, 236, 136, 152, 32, 236, 151, 134, 235, 139, 168, 235, 139, 164, 46, 13, 10, 236, 152)
	temp = append(temp, 164, 235, 161, 175, 236, 157, 180, 32, 235, 132, 136, 236, 157, 152, 32, 237, 158, 152, 235, 167, 140, 236, 156, 188, 235, 161, 156, 32, 235, 167, 137, 236, 150, 180, 235, 179, 180, 235, 160, 164)
	temp = append(temp, 235, 172, 180, 235, 130, 152, 46, 13, 10, 237, 155, 140, 235, 165, 173, 237, 149, 152, 234, 181, 172, 235, 130, 152, 46, 13, 10, 234, 183, 184, 235, 158, 152, 44, 32, 236, 157, 180, 32, 236)
	temp = append(temp, 160, 149, 235, 143, 132, 235, 169, 180, 32, 236, 167, 128, 236, 188, 156, 235, 179, 188, 32, 234, 176, 128, 236, 185, 152, 234, 176, 128, 32, 236, 158, 136, 234, 178, 160, 236, 167, 128, 46, 13)
	temp = append(temp, 10, 235, 130, 180, 234, 176, 128, 32, 236, 150, 180, 235, 150, 187, 234, 178, 140, 32, 234, 183, 184, 235, 166, 172, 32, 237, 149, 156, 236, 151, 134, 236, 157, 180, 32, 236, 158, 148, 236, 157)
	temp = append(temp, 184, 237, 149, 180, 236, 167, 136, 32, 236, 136, 152, 32, 236, 158, 136, 236, 151, 136, 235, 138, 148, 236, 167, 128, 32, 236, 149, 140, 235, 160, 164, 236, 164, 132, 234, 185, 140, 63, 13, 10)
	temp = append(temp, 236, 130, 172, 235, 158, 140, 235, 147, 164, 236, 157, 128, 32, 235, 170, 168, 235, 145, 144, 32, 235, 182, 136, 236, 149, 136, 236, 157, 132, 32, 234, 176, 128, 236, 167, 132, 32, 236, 177, 132)
	temp = append(temp, 235, 161, 156, 32, 236, 130, 180, 236, 149, 132, 234, 176, 132, 235, 139, 168, 235, 139, 164, 46, 13, 10, 236, 157, 180, 234, 177, 180, 32, 235, 175, 184, 236, 167, 128, 236, 157, 152, 32, 236)
	temp = append(temp, 152, 129, 236, 151, 173, 236, 157, 132, 32, 235, 167, 158, 235, 139, 165, 235, 156, 168, 235, 166, 180, 32, 235, 149, 140, 32, 235, 138, 144, 235, 129, 188, 235, 138, 148, 32, 235, 139, 185, 236)
	temp = append(temp, 151, 176, 237, 149, 156, 32, 235, 140, 128, 234, 176, 128, 236, 149, 188, 46, 13, 10, 237, 149, 152, 236, 167, 128, 235, 167, 140, 32, 235, 130, 152, 235, 138, 148, 32, 236, 132, 184, 236, 131)
	temp = append(temp, 129, 236, 151, 144, 32, 236, 130, 180, 236, 149, 132, 235, 130, 168, 234, 184, 176, 32, 236, 156, 132, 237, 149, 180, 32, 234, 183, 184, 32, 234, 179, 181, 237, 143, 172, 235, 165, 188, 32, 235)
	temp = append(temp, 176, 155, 236, 149, 132, 235, 147, 164, 236, 157, 180, 236, 167, 128, 32, 236, 149, 138, 236, 157, 128, 32, 236, 177, 132, 32, 236, 138, 164, 236, 138, 164, 235, 161, 156, 32, 235, 168, 185, 236)
	temp = append(temp, 150, 180, 235, 178, 132, 235, 160, 184, 236, 150, 180, 46, 13, 10, 234, 183, 184, 234, 178, 131, 236, 157, 180, 32, 235, 130, 180, 234, 176, 128, 32, 236, 160, 128, 236, 167, 128, 235, 165, 184)
	temp = append(temp, 32, 236, 181, 156, 236, 180, 136, 236, 157, 180, 236, 158, 144, 32, 236, 181, 156, 236, 149, 133, 236, 157, 152, 32, 236, 149, 133, 237, 150, 137, 236, 157, 180, 236, 151, 136, 236, 167, 128, 46)
	temp = append(temp, 13, 10, 237, 155, 132, 237, 154, 140, 237, 149, 152, 236, 167, 132, 32, 236, 149, 138, 236, 149, 132, 46, 32, 236, 131, 157, 236, 161, 180, 236, 157, 132, 32, 236, 156, 132, 237, 149, 156, 32)
	temp = append(temp, 236, 132, 160, 237, 131, 157, 236, 157, 180, 236, 151, 136, 236, 156, 188, 235, 139, 136, 46, 13, 10, 235, 132, 136, 235, 143, 132, 32, 235, 167, 136, 236, 176, 172, 234, 176, 128, 236, 167, 128)
	temp = append(temp, 32, 236, 149, 132, 235, 139, 136, 236, 151, 136, 235, 139, 136, 63, 13, 10, 235, 130, 180, 32, 235, 168, 184, 235, 166, 191, 236, 134, 141, 236, 157, 132, 32, 234, 176, 136, 234, 184, 176, 234)
	temp = append(temp, 176, 136, 234, 184, 176, 32, 235, 182, 132, 237, 149, 180, 236, 139, 156, 237, 130, 164, 235, 141, 152, 32, 235, 132, 136, 236, 157, 152, 32, 237, 145, 156, 236, 160, 149, 236, 157, 132, 32, 235)
	temp = append(temp, 179, 180, 234, 179, 160, 32, 236, 149, 140, 32, 236, 136, 152, 32, 236, 158, 136, 236, 151, 136, 236, 150, 180, 46, 13, 10, 236, 150, 180, 236, 169, 148, 32, 236, 136, 152, 32, 236, 151, 134)
	temp = append(temp, 236, 151, 136, 235, 139, 164, 32, 235, 157, 188, 234, 179, 160, 32, 235, 132, 136, 235, 143, 132, 32, 235, 167, 144, 236, 157, 132, 32, 237, 149, 152, 235, 160, 164, 235, 130, 152, 46, 13, 10)
	temp = append(temp, 236, 157, 180, 236, 160, 156, 235, 138, 148, 32, 235, 139, 164, 236, 139, 156, 32, 236, 158, 138, 234, 179, 160, 32, 236, 158, 136, 235, 141, 152, 32, 234, 179, 181, 237, 143, 172, 235, 165, 188)
	temp = append(temp, 32, 235, 167, 136, 236, 163, 188, 237, 149, 152, 235, 138, 148, 32, 234, 178, 140, 32, 236, 162, 139, 236, 157, 132, 32, 234, 177, 176, 236, 149, 188, 46, 13, 10, 234, 181, 180, 235, 160, 136)
	temp = append(temp, 235, 165, 188, 32, 235, 129, 138, 236, 150, 180, 235, 130, 180, 234, 178, 160, 235, 139, 164, 235, 169, 180, 46, 13, 10, 235, 132, 164, 234, 176, 128, 32, 234, 183, 184, 235, 160, 135, 235, 139)
	temp = append(temp, 164, 235, 169, 180, 32, 234, 183, 184, 235, 159, 176, 32, 234, 178, 131, 236, 157, 180, 234, 178, 160, 236, 167, 128, 46, 13, 10)
	return temp
}
