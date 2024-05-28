// test659 : stdlib5.legsup gen3kaes

package legsup

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strconv"
	"strings"
	"sync"
)

// gen3 kaes generate 64B pwhash
func g3pwhash(salt []byte, kf []byte, pw []byte) []byte {
	temp := make([]byte, 0)
	var tb []byte
	temp = append(temp, pw...)
	for i := 0; i < 100; i++ {
		temp = hash512(append(temp, kf...))
		for j := 0; j < 5000; j++ {
			tb = make([]byte, 0)
			tb = append(tb, salt...)
			tb = append(tb, temp...)
			temp = hash512(tb)
			tb = make([]byte, 0)
			tb = append(tb, temp...)
			tb = append(tb, salt...)
			temp = hash512(tb)
		}
	}
	return temp
}

// gen3 kaes generate 64B masterkey
func g3mkey(salt []byte, kf []byte, pw []byte) []byte {
	nsalt := make([]byte, len(salt))
	nkf := make([]byte, len(kf))
	npw := make([]byte, len(pw))
	for i, r := range salt {
		nsalt[len(nsalt)-i-1] = r
	}
	for i, r := range kf {
		nkf[len(nkf)-i-1] = r
	}
	for i, r := range pw {
		npw[len(npw)-i-1] = r
	}

	temp := make([]byte, 0)
	var tb []byte
	temp = append(temp, npw...)
	for i := 0; i < 10; i++ {
		temp = hash512(append(temp, nkf...))
		for j := 0; j < 5000; j++ {
			tb = make([]byte, 0)
			tb = append(tb, nsalt...)
			tb = append(tb, temp...)
			temp = hash512(tb)
			tb = make([]byte, 0)
			tb = append(tb, temp...)
			tb = append(tb, nsalt...)
			temp = hash512(tb)
		}
	}
	return temp
}

// gen3 kaes 256B mainkey -> 32B * num ckey
func g3expkey(mainkey []byte, num int) [][]byte {
	out := make([][]byte, num)
	divsize := len(mainkey) / 4
	a := mainkey[0:divsize]
	b := mainkey[divsize : divsize*2]
	c := mainkey[divsize*2 : divsize*3]
	d := mainkey[divsize*3 : divsize*4]

	for i := 0; i < num; i++ {
		pre := []byte(fmt.Sprintf("%d", num-i))
		sub := make([]byte, 0)
		switch i % 4 {
		case 0:
			pre = append(pre, b...)
			sub = append(append(append(sub, c...), a...), d...)
		case 1:
			pre = append(pre, d...)
			sub = append(append(append(sub, a...), c...), b...)
		case 2: // !!! original py code was WRONG !!!
			//pre = append(pre, a...)
			//sub = append(append(append(sub, b...), d...), c...)
			pre = append(pre, c...)
			sub = append(append(append(sub, d...), b...), a...)
		case 3:
			pre = append(pre, c...)
			sub = append(append(append(sub, d...), b...), a...)
		}

		temp := sub
		for j := 0; j < 10; j++ {
			temp = append(temp, bytes.Repeat([]byte{216}, int(temp[0]))...)
			for k := 0; k < 10000; k++ {
				temp = hash512(append(pre, temp...))
			}
		}
		out[i] = temp[0:32]
	}
	return out
}

// simple mainhead string convertion
func g3hconv(data string, pw string) string {
	key := make([]byte, len(data))
	tb := []byte(pw)
	for i := 0; i < len(key); i++ {
		key[i] = tb[i%len(tb)]
	}
	if strings.Contains(data, "[") {
		temp := []byte(data)
		for i, r := range temp {
			temp[i] = r + key[i]
		}
		out := make([]string, len(temp))
		for i, r := range temp {
			out[i] = fmt.Sprintf("%d", r)
		}
		return strings.Join(out, " ")
	} else {
		data = strings.Replace(data, "\n", " ", -1)
		temp := strings.Split(data, " ")
		out := make([]byte, 0)
		for _, r := range temp {
			if r != "" {
				ti, err := strconv.Atoi(r)
				if err == nil {
					out = append(out, byte(ti))
				} else {
					return ""
				}
			}
		}
		for i, r := range out {
			out[i] = r - key[i]
		}
		return string(out)
	}
}

// gen3 kaes internal worker
func g3calc(inbuf []byte, exbuf []byte, key [][]byte, iv [][]byte, num int, chunk int, isenc bool, wg *sync.WaitGroup) {
	defer wg.Done()
	ti0 := chunk * num
	ti1 := ti0 + chunk
	if isenc {
		copy(exbuf[ti0:ti1], aes16en(inbuf[ti0:ti1], key[num], iv[num]))
	} else {
		copy(exbuf[ti0:ti1], aes16de(inbuf[ti0:ti1], key[num], iv[num]))
	}
}

// gen3 kaes internal encryption
func g3enc(path string, tgt string, key [][]byte, iv [][]byte, head []byte, core int, chunk int) error {
	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	t, err := kio.Open(tgt, "w")
	if err == nil {
		defer t.Close()
	} else {
		return err
	}
	_, err = kio.Write(t, head)
	if err != nil {
		return nil
	}

	num0 := kio.Size(path)
	num1 := num0 / chunk
	num2 := num0 % chunk
	inbuf := make([]byte, core*chunk)
	exbuf := make([]byte, core*chunk)
	var wg sync.WaitGroup

	for i := 0; i < num1/32; i++ {
		f.Read(inbuf)
		wg.Add(32)
		for j := 0; j < 32; j++ {
			go g3calc(inbuf, exbuf, key, iv, j, chunk, true, &wg)
		}
		wg.Wait()
		t.Write(exbuf)
	}

	if num1%32 != 0 {
		inbuf = make([]byte, (num1%32)*chunk)
		exbuf = make([]byte, (num1%32)*chunk)
		f.Read(inbuf)
		wg.Add(num1 % 32)
		for i := 0; i < num1%32; i++ {
			go g3calc(inbuf, exbuf, key, iv, i, chunk, true, &wg)
		}
		wg.Wait()
		t.Write(exbuf)
	}

	inbuf, _ = kio.Read(f, num2)
	exbuf = aes1en(inbuf, key[num1%core], iv[num1%core])
	kio.Write(t, exbuf)
	return nil
}

// gen3 kaes internal decryption
func g3dec(path string, tgt string, key [][]byte, iv [][]byte, stpoint int, core int, chunk int) error {
	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	t, err := kio.Open(tgt, "w")
	if err == nil {
		defer t.Close()
	} else {
		return err
	}
	f.Seek(int64(stpoint), 0)

	num0 := kio.Size(path) - stpoint
	num1 := num0 / chunk
	num2 := num0 % chunk
	if num2 == 0 {
		num1 = num1 - 1
		num2 = chunk
	}
	inbuf := make([]byte, core*chunk)
	exbuf := make([]byte, core*chunk)
	var wg sync.WaitGroup

	for i := 0; i < num1/32; i++ {
		f.Read(inbuf)
		wg.Add(32)
		for j := 0; j < 32; j++ {
			go g3calc(inbuf, exbuf, key, iv, j, chunk, false, &wg)
		}
		wg.Wait()
		t.Write(exbuf)
	}

	if num1%32 != 0 {
		inbuf = make([]byte, (num1%32)*chunk)
		exbuf = make([]byte, (num1%32)*chunk)
		f.Read(inbuf)
		wg.Add(num1 % 32)
		for i := 0; i < num1%32; i++ {
			go g3calc(inbuf, exbuf, key, iv, i, chunk, false, &wg)
		}
		wg.Wait()
		t.Write(exbuf)
	}

	inbuf, _ = kio.Read(f, num2)
	exbuf = aes1de(inbuf, key[num1%core], iv[num1%core])
	kio.Write(t, exbuf)
	return nil
}

// gen3 kaes all-mode
type G3kaesall struct {
	Prehead  []byte // png 1024nB / 0B
	Metadata []byte // 18B meta size
	Mainhead []byte // header
	Subhead  []byte // reserved subheader

	Hidename bool   // hide original name
	Hint     string // hint string
	Respath  string // conversion result path

	core      int    // multithreading core
	chunksize int    // size of one chunk
	salt      []byte // salt bytes
	pwhash    []byte // pwhash bytes
	ckeydata  []byte // encrypted content key
	tkeydata  []byte // encrypted title key
	enctitle  []byte // encrypted file name
}

// make metadata, mainheader by internal value
func (tbox *G3kaesall) mkhead() error {
	var db G3kdb
	var err error
	if tbox.Hidename {
		err = db.Read("[meta]{[core]{0}[chunk]{0}[keydt]{\"\"}[salt]{\"\"}[pwhash]{\"\"}[hint]{\"\"}[tkeydt]{\"\"}[title]{\"\"}}")
	} else {
		err = db.Read("[meta]{[core]{0}[chunk]{0}[keydt]{\"\"}[salt]{\"\"}[pwhash]{\"\"}[hint]{\"\"}}")
	}
	if err != nil {
		return err
	}
	temp := db.Locate("meta#core")
	temp.Data.IntV = tbox.core
	temp = db.Locate("meta#chunk")
	temp.Data.IntV = tbox.chunksize
	temp = db.Locate("meta#keydt")
	temp.Data.StrV = kb64en(tbox.ckeydata)
	temp = db.Locate("meta#salt")
	temp.Data.StrV = kb64en(tbox.salt)
	temp = db.Locate("meta#pwhash")
	temp.Data.StrV = kb64en(tbox.pwhash)
	temp = db.Locate("meta#hint")
	temp.Data.StrV = tbox.Hint
	if tbox.Hidename {
		temp = db.Locate("meta#tkeydt")
		temp.Data.StrV = kb64en(tbox.tkeydata)
		temp = db.Locate("meta#title")
		temp.Data.StrV = kb64en(tbox.enctitle)
	}
	db.Zipexp = true
	db.Zipstr = true
	ts := g3hconv(db.Write(), "Hello, world!")
	if ts == "" {
		return errors.New("invalid header conversion")
	}
	tbox.Mainhead = []byte(ts)
	tbox.Metadata = append([]byte("KES3"), 0, 0)
	tbox.Metadata = append(tbox.Metadata, kobj.Encode(len(tbox.Mainhead), 4)...)
	tbox.Metadata = append(tbox.Metadata, kobj.Encode(len(tbox.Subhead), 4)...)
	tb := make([]byte, 0)
	tb = append(append(tb, tbox.Mainhead...), tbox.Subhead...)
	tb = hash32(tb)
	for i := 3; i >= 0; i-- {
		tbox.Metadata = append(tbox.Metadata, tb[i])
	}
	return nil
}

// make internal data by mainhead
func (tbox *G3kaesall) rdhead() error {
	ts := g3hconv(string(tbox.Mainhead), "Hello, world!")
	var db G3kdb
	err := db.Read(ts)
	if err != nil {
		return nil
	}

	temp := db.Locate("meta#core")
	if temp == nil {
		return errors.New("no data named meta#core")
	} else {
		tbox.core = temp.Data.IntV
	}
	temp = db.Locate("meta#chunk")
	if temp == nil {
		return errors.New("no data named meta#chunk")
	} else {
		tbox.chunksize = temp.Data.IntV
	}
	temp = db.Locate("meta#keydt")
	if temp == nil {
		return errors.New("no data named meta#keydt")
	} else {
		tbox.ckeydata = kb64de(temp.Data.StrV)
	}
	temp = db.Locate("meta#salt")
	if temp == nil {
		return errors.New("no data named meta#salt")
	} else {
		tbox.salt = kb64de(temp.Data.StrV)
	}
	temp = db.Locate("meta#pwhash")
	if temp == nil {
		return errors.New("no data named meta#pwhash")
	} else {
		tbox.pwhash = kb64de(temp.Data.StrV)
	}
	temp = db.Locate("meta#hint")
	if temp == nil {
		return errors.New("no data named meta#hint")
	} else {
		tbox.Hint = temp.Data.StrV
	}

	temp = db.Locate("meta#tkeydt")
	if temp == nil {
		tbox.Hidename = false
	} else {
		tbox.Hidename = true
		tbox.tkeydata = kb64de(temp.Data.StrV)
	}
	temp = db.Locate("meta#title")
	if temp == nil {
		tbox.Hidename = false
	} else {
		tbox.enctitle = kb64de(temp.Data.StrV)
	}
	return nil
}

// gen3 kaes all-mode init, for base setting (32, 128k), set (core, chunk) to 0
func (tbox *G3kaesall) Init(core int, chunk int) {
	tbox.Prehead = G3pic()
	tbox.Prehead = append(tbox.Prehead, make([]byte, 1024-len(tbox.Prehead)%1024)...)
	tbox.Metadata = nil
	tbox.Mainhead = nil
	tbox.Subhead = nil
	tbox.Hidename = true
	tbox.Hint = ""
	tbox.Respath = ""
	if core > 0 {
		tbox.core = core
	} else {
		tbox.core = 32
	}
	if chunk > 0 {
		tbox.chunksize = chunk
	} else {
		tbox.chunksize = 131072
	}
}

// gen3 kaes all-mode encrypt
func (tbox *G3kaesall) Encrypt(path string, pw string, kf []byte) error {
	path = kio.Abs(path)
	if tbox.Hidename {
		ti := kobj.Decode(genrand(3))%900000 + 100000
		tbox.Respath = path[0:strings.LastIndex(path, "/")+1] + fmt.Sprintf("%d", ti)
	} else {
		tbox.Respath = path
	}
	if tbox.Prehead == nil {
		tbox.Respath = tbox.Respath + ".k"
	} else {
		tbox.Respath = tbox.Respath + ".png"
	}
	tbox.salt = genrand(64)
	tbox.pwhash = g3pwhash(tbox.salt, kf, []byte(pw))
	mkey := g3mkey(tbox.salt, kf, []byte(pw))
	ckey := genrand(256)
	tkey := genrand(32)
	iv := hash128(tbox.salt)
	name := path[strings.LastIndex(path, "/")+1:]

	tbox.ckeydata = aes1en(ckey, mkey[0:32], iv)
	iv = hash128(tbox.salt)
	tbox.tkeydata = aes1en(tkey, mkey[32:64], iv)
	iv = hash128(tbox.salt)
	tbox.enctitle = aes1en([]byte(name), tkey, iv)
	iv = hash128(tbox.salt)
	err := tbox.mkhead()
	if err != nil {
		return err
	}

	keys := g3expkey(ckey, tbox.core)
	ivs := make([][]byte, tbox.core)
	for i := 0; i < tbox.core; i++ {
		ivs[i] = make([]byte, 0)
		ivs[i] = append(ivs[i], iv...)
	}
	heads := append(tbox.Prehead, tbox.Metadata...)
	heads = append(heads, tbox.Mainhead...)
	heads = append(heads, tbox.Subhead...)
	err = g3enc(path, tbox.Respath, keys, ivs, heads, tbox.core, tbox.chunksize)
	return err
}

// gen3 kaes all-mode view
func (tbox *G3kaesall) View(path string) error {
	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}

	var temp []byte
	tbox.Prehead = nil
	temp, err = kio.Read(f, 4)
	if err != nil {
		return err
	}
	for !kio.Bequal(temp, []byte("KES3")) {
		tbox.Prehead = append(tbox.Prehead, temp...)
		temp, err = kio.Read(f, 1020)
		if err != nil {
			return err
		}
		if len(temp) == 1020 {
			tbox.Prehead = append(tbox.Prehead, temp...)
		} else {
			return errors.New("invalid gen3 kaes file")
		}
		temp, err = kio.Read(f, 4)
		if err != nil {
			return err
		}
	}

	tbox.Metadata = temp
	temp, err = kio.Read(f, 14)
	if err == nil {
		tbox.Metadata = append(tbox.Metadata, temp...)
	} else {
		return err
	}
	tbox.Mainhead, err = kio.Read(f, kobj.Decode(tbox.Metadata[6:10]))
	if err != nil {
		return err
	}
	tbox.Subhead, err = kio.Read(f, kobj.Decode(tbox.Metadata[10:14]))
	if err != nil {
		return err
	}

	temp = make([]byte, 0)
	temp = append(temp, tbox.Mainhead...)
	temp = append(temp, tbox.Subhead...)
	temp = hash32(temp)
	for i := 3; i >= 0; i-- {
		if temp[i] != tbox.Metadata[17-i] {
			return errors.New("wrong CRC32 value")
		}
	}
	return tbox.rdhead()
}

// gen3 kaes all-mode decrypt
func (tbox *G3kaesall) Decrypt(path string, pw string, kf []byte) error {
	if tbox.Mainhead == nil {
		return errors.New("should done View() first")
	}
	path = kio.Abs(path)
	if !kio.Bequal(tbox.pwhash, g3pwhash(tbox.salt, kf, []byte(pw))) {
		return errors.New("InvalidPW")
	}

	iv := hash128(tbox.salt)
	mkey := g3mkey(tbox.salt, kf, []byte(pw))
	ckey := aes1de(tbox.ckeydata, mkey[0:32], iv)
	iv = hash128(tbox.salt)
	var tkey []byte
	if tbox.Hidename {
		tkey = aes1de(tbox.tkeydata, mkey[32:64], iv)
		iv = hash128(tbox.salt)
	}
	keys := g3expkey(ckey, tbox.core)
	ivs := make([][]byte, tbox.core)
	for i := 0; i < tbox.core; i++ {
		ivs[i] = make([]byte, 0)
		ivs[i] = append(ivs[i], iv...)
	}

	var newpath string
	if tbox.Hidename {
		newpath = path[0:strings.LastIndex(path, "/")+1] + string(aes1de(tbox.enctitle, tkey, iv))
	} else {
		newpath = path[0:strings.LastIndex(path, ".")]
	}
	ti := len(tbox.Prehead) + len(tbox.Metadata) + len(tbox.Mainhead) + len(tbox.Subhead)
	return g3dec(path, newpath, keys, ivs, ti, tbox.core, tbox.chunksize)
}

// gen3 kaes func-mode
type G3kaesfunc struct {
	// not using Prehead/Subhead, enc setting fixed to (32, 131072)
	Metadata []byte // 18B meta size
	Mainhead []byte // header

	core      int    // multithreading core
	chunksize int    // size of one chunk
	iv        []byte // iv bytes
	ckeydata  []byte // encrypted content key
}

// make metadata, mainheader by internal value
func (tbox *G3kaesfunc) mkhead() error {
	var db G3kdb
	err := db.Read("[meta]{[core]{0}[chunk]{0}[iv]{\"\"}[ckeydt]{\"\"}}")
	if err != nil {
		return err
	}
	temp := db.Locate("meta#core")
	temp.Data.IntV = tbox.core
	temp = db.Locate("meta#chunk")
	temp.Data.IntV = tbox.chunksize
	temp = db.Locate("meta#iv")
	temp.Data.StrV = kb64en(tbox.iv)
	temp = db.Locate("meta#ckeydt")
	temp.Data.StrV = kb64en(tbox.ckeydata)

	db.Zipexp = true
	db.Zipstr = true
	ts := g3hconv(db.Write(), "Hello, world!")
	if ts == "" {
		return errors.New("invalid header conversion")
	}
	tbox.Mainhead = []byte(ts)
	tbox.Metadata = append([]byte("KES3"), 0, 0)
	tbox.Metadata = append(tbox.Metadata, kobj.Encode(len(tbox.Mainhead), 4)...)
	tbox.Metadata = append(tbox.Metadata, 0, 0, 0, 0)
	tb := make([]byte, 0)
	tb = append(tb, tbox.Mainhead...)
	tb = hash32(tb)
	for i := 3; i >= 0; i-- {
		tbox.Metadata = append(tbox.Metadata, tb[i])
	}
	return nil
}

// make internal data by mainhead
func (tbox *G3kaesfunc) rdhead() error {
	ts := g3hconv(string(tbox.Mainhead), "Hello, world!")
	var db G3kdb
	err := db.Read(ts)
	if err != nil {
		return nil
	}

	temp := db.Locate("meta#core")
	if temp == nil {
		return errors.New("no data named meta#core")
	} else {
		tbox.core = temp.Data.IntV
	}
	temp = db.Locate("meta#chunk")
	if temp == nil {
		return errors.New("no data named meta#chunk")
	} else {
		tbox.chunksize = temp.Data.IntV
	}
	temp = db.Locate("meta#ckeydt")
	if temp == nil {
		return errors.New("no data named meta#ckeydt")
	} else {
		tbox.ckeydata = kb64de(temp.Data.StrV)
	}
	temp = db.Locate("meta#iv")
	if temp == nil {
		return errors.New("no data named meta#iv")
	} else {
		tbox.iv = kb64de(temp.Data.StrV)
	}
	return nil
}

// gen3 kaes func-mode encrypt
func (tbox *G3kaesfunc) Encrypt(before string, after string, akey []byte) error {
	tbox.core = 32
	tbox.chunksize = 131072
	ckey := genrand(256)
	tbox.iv = genrand(16)

	if len(akey) != 32 {
		return errors.New("akey should be 32B")
	}
	temp := make([]byte, 16)
	copy(temp, tbox.iv)
	tbox.ckeydata = aes1en(ckey, akey, temp)
	err := tbox.mkhead()
	if err != nil {
		return err
	}

	keys := g3expkey(ckey, tbox.core)
	ivs := make([][]byte, tbox.core)
	for i := 0; i < tbox.core; i++ {
		ivs[i] = make([]byte, 0)
		ivs[i] = append(ivs[i], tbox.iv...)
	}
	heads := make([]byte, 0)
	heads = append(heads, tbox.Metadata...)
	heads = append(heads, tbox.Mainhead...)
	return g3enc(before, after, keys, ivs, heads, tbox.core, tbox.chunksize)
}

// gen3 kaes func-mode decrypt
func (tbox *G3kaesfunc) Decrypt(before string, after string, akey []byte) error {
	f, err := kio.Open(before, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}

	temp, err := kio.Read(f, 4)
	if err != nil {
		return err
	}
	if !kio.Bequal(temp, []byte("KES3")) {
		return errors.New("invalid gen3 kaes file")
	}

	tbox.Metadata = temp
	temp, err = kio.Read(f, 14)
	if err != nil {
		return err
	}
	tbox.Metadata = append(tbox.Metadata, temp...)
	tbox.Mainhead, err = kio.Read(f, kobj.Decode(tbox.Metadata[6:10]))
	if err != nil {
		return err
	}

	if kobj.Decode(tbox.Metadata[10:14]) != 0 {
		return errors.New("uncontrollable subhead")
	}
	temp = hash32(tbox.Mainhead)
	for i := 3; i >= 0; i-- {
		if temp[i] != tbox.Metadata[17-i] {
			return errors.New("wrong CRC32 value")
		}
	}
	err = tbox.rdhead()
	if err != nil {
		return err
	}

	if len(akey) != 32 {
		return errors.New("akey should be 32B")
	}
	temp = make([]byte, 16)
	copy(temp, tbox.iv)
	ckey := aes1de(tbox.ckeydata, akey, temp)

	keys := g3expkey(ckey, tbox.core)
	ivs := make([][]byte, tbox.core)
	for i := 0; i < tbox.core; i++ {
		ivs[i] = make([]byte, 0)
		ivs[i] = append(ivs[i], tbox.iv...)
	}
	return g3dec(before, after, keys, ivs, len(tbox.Metadata)+len(tbox.Mainhead), tbox.core, tbox.chunksize)
}

// gen3 kaes kv3 (simple vault)
type G3kv3 struct {
	Hint   string // vault pw hint
	salt   []byte // salt bytes
	pwhash []byte // pwhash bytes
	akeydt []byte // encrypted akey
}

// gen3 kaes kv3 enc/dec convert folder (~/)
func (tbox *G3kv3) calc(path string, akey []byte, isenc bool) error {
	infos, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	var engine G3kaesfunc
	var temp string
	for _, r := range infos {
		if r.IsDir() {
			temp = path + r.Name()
			if temp[len(temp)-1] != '/' {
				temp = temp + "/"
			}
			err = tbox.calc(temp, akey, isenc)
		} else {
			temp = path + r.Name()
			if isenc {
				err = engine.Encrypt(temp, temp+".kv3", akey)
			} else {
				err = engine.Decrypt(temp, temp[0:strings.LastIndex(temp, ".")], akey)
			}
		}
		if err != nil {
			return err
		}
		os.Remove(temp)
	}
	return nil
}

// gen3 kaes kv3 encrypt folder/file
func (tbox *G3kv3) Encrypt(pw string, kf []byte, path string) error {
	tbox.salt = genrand(32)
	akey := genrand(32)
	tbox.pwhash = g3pwhash(tbox.salt, kf, []byte(pw))
	mkey := g3mkey(tbox.salt, kf, []byte(pw))
	tbox.akeydt = aes1en(akey, mkey[0:32], mkey[32:48])

	var db G3kdb
	err := db.Read("[meta]{[hint]{\"\"}[salt]{\"\"}[pwhash]{\"\"}[akeydt]{\"\"}}")
	if err != nil {
		return err
	}
	temp := db.Locate("meta#hint")
	temp.Data.StrV = tbox.Hint
	temp = db.Locate("meta#salt")
	temp.Data.StrV = kb64en(tbox.salt)
	temp = db.Locate("meta#pwhash")
	temp.Data.StrV = kb64en(tbox.pwhash)
	temp = db.Locate("meta#akeydt")
	temp.Data.StrV = kb64en(tbox.akeydt)
	db.Zipexp = false
	db.Zipstr = true
	header := []byte(db.Write())
	var hpath string

	path = kio.Abs(path)
	if path[len(path)-1] == '/' {
		err = tbox.calc(path, akey, true)
		hpath = path[0:len(path)-1] + ".ench.txt"
	} else {
		var engine G3kaesfunc
		err = engine.Encrypt(path, path+".kv3", akey)
		os.Remove(path)
		hpath = path + ".ench.txt"
	}
	if err != nil {
		return err
	}

	f, err := kio.Open(hpath, "w")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	_, err = kio.Write(f, header)
	return err
}

// gen3 kaes kv3 view ench.txt metadata
func (tbox *G3kv3) View(path string) error {
	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	temp, err := kio.Read(f, -1)
	if err != nil {
		return err
	}

	var db G3kdb
	err = db.Read(string(temp))
	if err != nil {
		return err
	}

	tgt := db.Locate("meta#hint")
	if tgt == nil {
		return errors.New("no data named meta#hint")
	} else {
		tbox.Hint = tgt.Data.StrV
	}
	tgt = db.Locate("meta#salt")
	if tgt == nil {
		return errors.New("no data named meta#salt")
	} else {
		tbox.salt = kb64de(tgt.Data.StrV)
	}
	tgt = db.Locate("meta#pwhash")
	if tgt == nil {
		return errors.New("no data named meta#pwhash")
	} else {
		tbox.pwhash = kb64de(tgt.Data.StrV)
	}
	tgt = db.Locate("meta#akeydt")
	if tgt == nil {
		return errors.New("no data named meta#akeydt")
	} else {
		tbox.akeydt = kb64de(tgt.Data.StrV)
	}
	return nil
}

// gen3 kaes kv3 decrypt folder/file
func (tbox *G3kv3) Decrypt(pw string, kf []byte, path string) error {
	if tbox.pwhash == nil {
		return errors.New("should done View() first")
	}
	path = kio.Abs(path)
	if !kio.Bequal(tbox.pwhash, g3pwhash(tbox.salt, kf, []byte(pw))) {
		return errors.New("InvalidPW")
	}
	mkey := g3mkey(tbox.salt, kf, []byte(pw))
	akey := aes1de(tbox.akeydt, mkey[0:32], mkey[32:48])

	var err error
	if path[len(path)-1] == '/' {
		err = tbox.calc(path, akey, false)
		os.Remove(path[0:len(path)-1] + ".ench.txt")
	} else {
		var engine G3kaesfunc
		err = engine.Decrypt(path, path[0:strings.LastIndex(path, ".")], akey)
		os.Remove(path)
		os.Remove(path[0:strings.LastIndex(path, ".")] + ".ench.txt")
	}
	return err
}
