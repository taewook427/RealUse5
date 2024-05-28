// test664 : stdlib5.legsup gen4kaes

package legsup

import (
	"errors"
	"fmt"
	"os"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strings"
	"sync"
)

// gen4 kaes universal reader/writer
type G4io struct {
	IsBin    bool
	IsReader bool

	buffer []byte
	pos    int

	file *os.File
	path string
}

// gen4 kaes setting IOnode binary
func (data *G4io) OpenB(raw []byte, isreader bool) {
	data.IsBin = true
	data.IsReader = isreader
	data.buffer = raw
	data.pos = 0
	data.file = nil
	data.path = ""
}

// gen4 kaes setting IOnode file
func (data *G4io) OpenF(path string, isreader bool) error {
	data.IsBin = false
	data.IsReader = isreader
	data.buffer = nil
	data.pos = 0
	var err error
	if isreader {
		data.file, err = kio.Open(path, "r")
	} else {
		data.file, err = kio.Open(path, "w")
	}
	data.path = path
	return err
}

// gen4 kaes closing IOnode binary
func (data *G4io) CloseB() []byte {
	if data.IsBin {
		temp := data.buffer
		data.pos = 0
		data.buffer = nil
		return temp
	} else {
		return nil
	}
}

// gen4 kaes closing IOnode file
func (data *G4io) CloseF() {
	data.file.Close()
	data.path = ""
}

// gen4 kaes get size (binary/file)
func (data *G4io) Size() int {
	if data.IsBin {
		return len(data.buffer)
	} else {
		return kio.Size(data.path)
	}
}

// gen4 kaes seek (binary/file), works readmode only
func (data *G4io) Seek(pos int) {
	if data.IsReader {
		if data.IsBin {
			data.pos = pos
		} else {
			data.file.Seek(int64(pos), 0)
		}
	}
}

// gen4 kaes read (binary/file), returns nil if writemode
func (data *G4io) Read(size int) []byte {
	if data.IsReader {
		if data.IsBin {
			start := data.pos
			end := data.pos + size
			if end > len(data.buffer) {
				end = len(data.buffer)
			}
			data.pos = end
			return data.buffer[start:end]
		} else {
			temp, _ := kio.Read(data.file, size)
			return temp
		}
	} else {
		return nil
	}
}

// gen4 kaes write (binary/file), does nothing if readmode
func (data *G4io) Write(chunk []byte) {
	if !data.IsReader && chunk != nil {
		if data.IsBin {
			data.buffer = append(data.buffer, chunk...)
		} else {
			kio.Write(data.file, chunk)
		}
	}
}

// gen4 kaes internal calc, pos : 0~31 : enc / 32~63 : dec
func g4gencalc(buf []byte, key [][]byte, iv [][]byte, pos int, wg *sync.WaitGroup) {
	defer wg.Done()
	var pos0 int
	var pos1 int
	if pos < 32 {
		pos0 = 131072 * pos
		pos1 = pos0 + 131072
		copy(buf[pos0:pos1], aes16en(buf[pos0:pos1], key[pos], iv[pos]))
	} else {
		pos = pos - 32
		pos0 = 131072 * pos
		pos1 = pos0 + 131072
		copy(buf[pos0:pos1], aes16de(buf[pos0:pos1], key[pos], iv[pos]))
	}
}

// gen4 kaes generic encrypt
func g4genenc(before *G4io, after *G4io, key [][]byte, iv [][]byte, header []byte) {
	after.Write(header)
	size := before.Size()
	num0 := size / 4194304
	num1 := size % 4194304
	num2 := num1 / 131072
	num3 := num1 % 131072
	var inbuf []byte = nil
	var exbuf []byte = nil
	var tmbuf []byte = nil
	var wg sync.WaitGroup

	if num0 > 0 {
		inbuf = before.Read(4194304)

		for i := 0; i < num0-1; i++ {
			wg.Add(32)
			for j := 0; j < 32; j++ {
				go g4gencalc(inbuf, key, iv, j, &wg)
			}
			after.Write(exbuf)
			tmbuf = before.Read(4194304)
			wg.Wait()
			exbuf = inbuf
			inbuf = tmbuf
		}

		wg.Add(32)
		for i := 0; i < 32; i++ {
			go g4gencalc(inbuf, key, iv, i, &wg)
		}
		after.Write(exbuf)
		tmbuf = before.Read(131072 * num2)
		wg.Wait()
		exbuf = inbuf
		inbuf = tmbuf

	} else {
		inbuf = before.Read(131072 * num2)
	}

	wg.Add(num2)
	for i := 0; i < num2; i++ {
		go g4gencalc(inbuf, key, iv, i, &wg)
	}
	after.Write(exbuf)
	tmbuf = before.Read(num3)
	wg.Wait()
	exbuf = inbuf
	inbuf = tmbuf

	after.Write(exbuf)
	after.Write(aes1en(inbuf, key[num2], iv[num2]))
}

// gen4 kaes generic decrypt
func g4gendec(before *G4io, after *G4io, key [][]byte, iv [][]byte, stpoint int) {
	before.Seek(stpoint)
	size := before.Size() - stpoint
	num0 := size / 4194304
	num1 := size % 4194304
	if num1 == 0 {
		num0 = num0 - 1
		num1 = 4194304
	}
	num2 := num1 / 131072
	num3 := num1 % 131072
	if num3 == 0 {
		num2 = num2 - 1
		num3 = 131072
	}
	var inbuf []byte = nil
	var exbuf []byte = nil
	var tmbuf []byte = nil
	var wg sync.WaitGroup

	if num0 > 0 {
		inbuf = before.Read(4194304)

		for i := 0; i < num0-1; i++ {
			wg.Add(32)
			for j := 0; j < 32; j++ {
				go g4gencalc(inbuf, key, iv, j+32, &wg)
			}
			after.Write(exbuf)
			tmbuf = before.Read(4194304)
			wg.Wait()
			exbuf = inbuf
			inbuf = tmbuf
		}

		wg.Add(32)
		for i := 0; i < 32; i++ {
			go g4gencalc(inbuf, key, iv, i+32, &wg)
		}
		after.Write(exbuf)
		tmbuf = before.Read(131072 * num2)
		wg.Wait()
		exbuf = inbuf
		inbuf = tmbuf

	} else {
		inbuf = before.Read(131072 * num2)
	}

	wg.Add(num2)
	for i := 0; i < num2; i++ {
		go g4gencalc(inbuf, key, iv, i+32, &wg)
	}
	after.Write(exbuf)
	tmbuf = before.Read(num3)
	wg.Wait()
	exbuf = inbuf
	inbuf = tmbuf

	after.Write(exbuf)
	after.Write(aes1de(inbuf, key[num2], iv[num2]))
}

// gen4 kaes all-mode
type G4kaesall struct {
	Hint []byte

	stpoint  int
	salt     []byte
	pwhash   []byte
	ckeydata []byte
	tkeydata []byte
	enctitle []byte
}

// gen4 kaes all-mode encrypt binary
func (tbox *G4kaesall) EnBin(pw []byte, kf []byte, data []byte) ([]byte, error) {
	tbox.salt = genrand(32)
	tb := append(append(append(append(append(make([]byte, 0), pw...), pw...), kf...), kf...), pw...)
	tbox.pwhash = hashscrypt(tb, tbox.salt, 524288, 8, 1, 256)
	tb = append(append(append(append(append(make([]byte, 0), kf...), pw...), pw...), kf...), pw...)
	mkey := hashscrypt(tb, tbox.salt, 16384, 8, 1, 48)
	ckey := genrand(1536)
	tbox.ckeydata = aes16en(ckey, mkey[0:32], mkey[32:48])

	tg0 := make(map[string]G4data)
	var tg1 G4data
	tg1.Set("WHOLE")
	tg0["MODE"] = tg1
	var tg2 G4data
	tg2.Set(tbox.salt)
	tg0["SALT"] = tg2
	var tg3 G4data
	tg3.Set(tbox.pwhash)
	tg0["PWH"] = tg3
	var tg4 G4data
	tg4.Set(tbox.ckeydata)
	tg0["CKDT"] = tg4
	var tg5 G4data
	tg5.Set(tbox.Hint)
	tg0["HINT"] = tg5
	tb = G4pic()
	header := append(tb, make([]byte, 128-len(tb)%128)...)
	tb = []byte(G4DBwrite(tg0))
	header = append(append(append(header, []byte("KAES4")...), kobj.Encode(len(tb), 3)...), tb...)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		keys[i] = ckey[32*i : 32*i+32]
		ivs[i] = ckey[1024+16*i : 1024+16*i+16]
	}

	var inbuf G4io
	inbuf.OpenB(data, true)
	var exbuf G4io
	exbuf.OpenB(make([]byte, 0), false)
	g4genenc(&inbuf, &exbuf, keys, ivs, header)
	return exbuf.CloseB(), nil
}

// gen4 kaes all-mode encrypt file
func (tbox *G4kaesall) EnFile(pw []byte, kf []byte, path string) (string, error) {
	tbox.salt = genrand(32)
	tb := append(append(append(append(append(make([]byte, 0), pw...), pw...), kf...), kf...), pw...)
	tbox.pwhash = hashscrypt(tb, tbox.salt, 524288, 8, 1, 256)
	tb = append(append(append(append(append(make([]byte, 0), kf...), pw...), pw...), kf...), pw...)
	mkey := hashscrypt(tb, tbox.salt, 16384, 8, 1, 48)
	ckey := genrand(1536)
	mkcpy := append(make([]byte, 0), mkey...)
	tbox.ckeydata = aes16en(ckey, mkey[0:32], mkey[32:48])

	path = kio.Abs(path)
	name := []byte(path[strings.LastIndex(path, "/")+1:])
	newpath := fmt.Sprintf("%s/%s.png", path[0:strings.LastIndex(path, "/")], kio.Bprint(genrand(3)))
	tkey := genrand(48)
	tbox.tkeydata = aes16en(tkey, mkcpy[0:32], mkcpy[32:48])
	tbox.enctitle = aes1en(name, tkey[0:32], tkey[32:48])

	tg0 := make(map[string]G4data)
	var tg1 G4data
	tg1.Set("WHOLE")
	tg0["MODE"] = tg1
	var tg2 G4data
	tg2.Set(tbox.salt)
	tg0["SALT"] = tg2
	var tg3 G4data
	tg3.Set(tbox.pwhash)
	tg0["PWH"] = tg3
	var tg4 G4data
	tg4.Set(tbox.ckeydata)
	tg0["CKDT"] = tg4
	var tg5 G4data
	tg5.Set(tbox.Hint)
	tg0["HINT"] = tg5
	var tg6 G4data
	tg6.Set(tbox.tkeydata)
	tg0["TKDT"] = tg6
	var tg7 G4data
	tg7.Set(tbox.enctitle)
	tg0["NMDT"] = tg7
	tb = G4pic()
	header := append(tb, make([]byte, 128-len(tb)%128)...)
	tb = []byte(G4DBwrite(tg0))
	header = append(append(append(header, []byte("KAES4")...), kobj.Encode(len(tb), 3)...), tb...)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		keys[i] = ckey[32*i : 32*i+32]
		ivs[i] = ckey[1024+16*i : 1024+16*i+16]
	}

	var inbuf G4io
	err := inbuf.OpenF(path, true)
	if err == nil {
		defer inbuf.CloseF()
	} else {
		return "", err
	}
	var exbuf G4io
	err = exbuf.OpenF(newpath, false)
	if err == nil {
		defer exbuf.CloseF()
	} else {
		return "", err
	}
	g4genenc(&inbuf, &exbuf, keys, ivs, header)
	return newpath, nil
}

// gen4 kaes all-mode view encbinary
func (tbox *G4kaesall) ViewBin(data []byte) error {
	pos := 0
	for pos+5 < len(data) {
		if kio.Bequal(data[pos:pos+5], []byte("KAES4")) {
			break
		} else {
			pos = pos + 128
		}
	}
	if pos+5 >= len(data) {
		return errors.New("invalid gen4 kaes file")
	}

	pos = pos + 8
	ti := kobj.Decode(data[pos-3 : pos])
	tbox.stpoint = pos + ti

	header := G4DBread(string(data[pos:tbox.stpoint]))
	if header["MODE"].StrV != "WHOLE" {
		return errors.New("invalid encryption mode")
	}
	tbox.salt = header["SALT"].ByteV
	tbox.pwhash = header["PWH"].ByteV
	tbox.ckeydata = header["CKDT"].ByteV
	tbox.Hint = header["HINT"].ByteV

	temp, exist := header["TKDT"]
	if exist {
		tbox.tkeydata = temp.ByteV
	} else {
		tbox.tkeydata = nil
	}
	temp, exist = header["NMDT"]
	if exist {
		tbox.enctitle = temp.ByteV
	} else {
		tbox.enctitle = nil
	}
	return nil
}

// gen4 kaes all-mode view encfile
func (tbox *G4kaesall) ViewFile(path string) error {
	var temp G4io
	err := temp.OpenF(path, true)
	if err == nil {
		defer temp.CloseF()
	} else {
		return err
	}

	pos := 0
	size := temp.Size()
	for pos+5 < size {
		if kio.Bequal(temp.Read(5), []byte("KAES4")) {
			break
		} else {
			pos = pos + 128
			temp.Seek(pos)
		}
	}
	if pos+5 >= size {
		return errors.New("invalid gen4 kaes file")
	}

	pos = pos + 8
	temp.Seek(pos - 3)
	ti := kobj.Decode(temp.Read(3))
	tbox.stpoint = pos + ti

	header := G4DBread(string(temp.Read(ti)))
	if header["MODE"].StrV != "WHOLE" {
		return errors.New("invalid encryption mode")
	}
	tbox.salt = header["SALT"].ByteV
	tbox.pwhash = header["PWH"].ByteV
	tbox.ckeydata = header["CKDT"].ByteV
	tbox.Hint = header["HINT"].ByteV

	tg, exist := header["TKDT"]
	if exist {
		tbox.tkeydata = tg.ByteV
	} else {
		tbox.tkeydata = nil
	}
	tg, exist = header["NMDT"]
	if exist {
		tbox.enctitle = tg.ByteV
	} else {
		tbox.enctitle = nil
	}
	return nil
}

// gen4 kaes all-mode decrypt binary
func (tbox *G4kaesall) DeBin(pw []byte, kf []byte, data []byte) ([]byte, error) {
	if len(tbox.salt) == 0 {
		return nil, errors.New("should done View() first")
	}
	tb := append(append(append(append(append(make([]byte, 0), pw...), pw...), kf...), kf...), pw...)
	pwhcmp := hashscrypt(tb, tbox.salt, 524288, 8, 1, 256)
	tb = append(append(append(append(append(make([]byte, 0), kf...), pw...), pw...), kf...), pw...)
	mkey := hashscrypt(tb, tbox.salt, 16384, 8, 1, 48)
	if !kio.Bequal(tbox.pwhash, pwhcmp) {
		return nil, errors.New("InvalidPW")
	}
	ckey := aes16de(tbox.ckeydata, mkey[0:32], mkey[32:48])

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		keys[i] = ckey[32*i : 32*i+32]
		ivs[i] = ckey[1024+16*i : 1024+16*i+16]
	}

	var inbuf G4io
	inbuf.OpenB(data, true)
	var exbuf G4io
	exbuf.OpenB(make([]byte, 0), false)
	g4gendec(&inbuf, &exbuf, keys, ivs, tbox.stpoint)
	return exbuf.CloseB(), nil
}

// gen4 kaes all-mode decrypt file
func (tbox *G4kaesall) DeFile(pw []byte, kf []byte, path string) (string, error) {
	if len(tbox.salt) == 0 {
		return "", errors.New("should done View() first")
	}
	tb := append(append(append(append(append(make([]byte, 0), pw...), pw...), kf...), kf...), pw...)
	pwhcmp := hashscrypt(tb, tbox.salt, 524288, 8, 1, 256)
	tb = append(append(append(append(append(make([]byte, 0), kf...), pw...), pw...), kf...), pw...)
	mkey := hashscrypt(tb, tbox.salt, 16384, 8, 1, 48)
	if !kio.Bequal(tbox.pwhash, pwhcmp) {
		return "", errors.New("InvalidPW")
	}
	mkcpy := append(make([]byte, 0), mkey...)
	ckey := aes16de(tbox.ckeydata, mkey[0:32], mkey[32:48])

	tkey := aes16de(tbox.tkeydata, mkcpy[0:32], mkcpy[32:48])
	name := string(aes1de(tbox.enctitle, tkey[0:32], tkey[32:48]))
	path = kio.Abs(path)
	newpath := path[0:strings.LastIndex(path, "/")+1] + name

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		keys[i] = ckey[32*i : 32*i+32]
		ivs[i] = ckey[1024+16*i : 1024+16*i+16]
	}

	var inbuf G4io
	err := inbuf.OpenF(path, true)
	if err == nil {
		defer inbuf.CloseF()
	} else {
		return "", err
	}
	var exbuf G4io
	err = exbuf.OpenF(newpath, false)
	if err == nil {
		defer exbuf.CloseF()
	} else {
		return "", err
	}
	g4gendec(&inbuf, &exbuf, keys, ivs, tbox.stpoint)
	return newpath, nil
}

// gen4 kaes func-mode
type G4kaesfunc struct {
	// should be initialized bin/file, !! g4func does not close it !!
	Inbuf G4io // reader
	Exbuf G4io // writer
}

// gen4 kaes func-mode encryption
func (tbox *G4kaesfunc) Encrypt(mkey []byte) error {
	if len(mkey) != 48 {
		return errors.New("mkey should be 48B")
	}
	ckey := genrand(1536)
	ckeydata := aes16en(ckey, mkey[0:32], mkey[32:48])

	tg0 := make(map[string]G4data)
	var tg1 G4data
	tg1.Set("FUNC")
	tg0["MODE"] = tg1
	var tg2 G4data
	tg2.Set(ckeydata)
	tg0["CKDT"] = tg2
	tb := []byte(G4DBwrite(tg0))
	header := append(append([]byte("KAES4"), kobj.Encode(len(tb), 3)...), tb...)

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		keys[i] = ckey[32*i : 32*i+32]
		ivs[i] = ckey[1024+16*i : 1024+16*i+16]
	}

	g4genenc(&tbox.Inbuf, &tbox.Exbuf, keys, ivs, header)
	return nil
}

// gen4 kaes func-mode decryption
func (tbox *G4kaesfunc) Decrypt(mkey []byte) error {
	pos := 0
	size := tbox.Inbuf.Size()
	for pos+5 < size {
		if kio.Bequal(tbox.Inbuf.Read(5), []byte("KAES4")) {
			break
		} else {
			pos = pos + 128
			tbox.Inbuf.Seek(pos)
		}
	}
	if pos+5 >= size {
		return errors.New("invalid gen4 kaes file")
	}

	pos = pos + 8
	tbox.Inbuf.Seek(pos - 3)
	ti := kobj.Decode(tbox.Inbuf.Read(3))
	header := G4DBread(string(tbox.Inbuf.Read(ti)))
	stpoint := pos + ti
	tbox.Inbuf.Seek(0)

	if len(mkey) != 48 {
		return errors.New("mkey should be 48B")
	}
	if header["MODE"].StrV != "FUNC" {
		return errors.New("invalid encryption mode")
	}
	ckeydata := header["CKDT"].ByteV
	ckey := aes16de(ckeydata, mkey[0:32], mkey[32:48])

	keys := make([][]byte, 32)
	ivs := make([][]byte, 32)
	for i := 0; i < 32; i++ {
		keys[i] = ckey[32*i : 32*i+32]
		ivs[i] = ckey[1024+16*i : 1024+16*i+16]
	}

	g4gendec(&tbox.Inbuf, &tbox.Exbuf, keys, ivs, stpoint)
	return nil
}
