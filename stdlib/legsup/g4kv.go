// test665 : stdlib5.legsup gen4kvault

package legsup

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strconv"
	"strings"
)

// gen4 kv4 internal file
type g4file struct {
	name string
	fptr []byte
}

// set internal data by file info (N$&name/fptr)
func (tbox *g4file) set(data []byte) {
	ti := 0
	for 48 <= data[ti] && data[ti] < 58 {
		ti = ti + 1
	}
	data = data[ti+1:]
	ti = bytes.Index(data, []byte("/"))
	tbox.name = string(data[0:ti])
	tbox.fptr = data[ti+1:]
}

// gen file info by internal data
func (tbox *g4file) gen(depth int) []byte {
	temp := []byte(fmt.Sprintf("%d$%s/", depth, tbox.name))
	return append(temp, tbox.fptr...)
}

// gen4 kv4 internal dir
type g4dir struct {
	name    string
	subdir  []g4dir
	subfile []g4file
}

// set internal data by file info (N#name)
func (tbox *g4dir) set(depth []int, data [][]byte) {
	tbox.subdir = make([]g4dir, 0)
	tbox.subfile = make([]g4file, 0)
	ti := 0
	for 48 <= data[0][ti] && data[0][ti] < 58 {
		ti = ti + 1
	}
	tbox.name = string(data[0][ti+1:])

	cur := 1
	for cur < len(depth) {
		ti := 0
		for 48 <= data[cur][ti] && data[cur][ti] < 58 {
			ti = ti + 1
		}
		if data[cur][ti] == 35 {
			ti = cur + 1
			var temp g4dir
			if len(depth) == ti {
				temp.set(depth[cur:], data[cur:])
			} else {
				for depth[cur] < depth[ti] {
					ti = ti + 1
					if len(depth) == ti {
						break
					}
				}
				temp.set(depth[cur:ti], data[cur:ti])
			}
			tbox.subdir = append(tbox.subdir, temp)
			cur = ti
		} else {
			var temp g4file
			temp.set(data[cur])
			tbox.subfile = append(tbox.subfile, temp)
			cur = cur + 1
		}
	}
}

// gen file info by internal data
func (tbox *g4dir) gen(depth int) []byte {
	temp := make([][]byte, 1)
	temp[0] = []byte(fmt.Sprintf("%d#%s", depth, tbox.name))
	for _, r := range tbox.subfile {
		temp = append(temp, r.gen(depth+1))
	}
	for _, r := range tbox.subdir {
		temp = append(temp, r.gen(depth+1))
	}
	return bytes.Join(temp, []byte("\n"))
}

// gen4 kv4 internal worker
type g4work struct {
	path  *string            // cluster path (~/)
	fkeys *map[string][]byte // kio.Bprint fptr-fkey
	fnums *map[string]int    // kio.Bprint fptr-chunknum
}

// get fptr from chunknum
func (tbox *g4work) getfptr(num int) {
	fs, _ := os.ReadDir(fmt.Sprintf("%s%d/", *tbox.path, num))
	for _, r := range fs {
		name := r.Name()
		if len(name) == 20 {
			if name[16:20] == ".kv4" {
				(*tbox.fnums)[kio.Bprint(kb64de(name[0:16]))] = num
			}
		}
	}
}

// make header file by mainhead, filesys, filekey
func (tbox *g4work) hpush(encmh []byte, encfs []byte, encfk []byte) {
	f, _ := kio.Open(fmt.Sprintf("%sheader.kv4", *tbox.path), "w")
	defer f.Close()
	kio.Write(f, []byte("KV4H"))
	kio.Write(f, kobj.Encode(len(encmh), 8))
	kio.Write(f, encmh)
	kio.Write(f, kobj.Encode(len(encfs), 8))
	kio.Write(f, encfs)
	kio.Write(f, kobj.Encode(len(encfk), 8))
	kio.Write(f, encfk)
}

// get mainhead, filesys, filekey from cluster head
func (tbox *g4work) hpop() ([]byte, []byte, []byte) {
	f, _ := kio.Open(fmt.Sprintf("%sheader.kv4", *tbox.path), "r")
	defer f.Close()
	tb, _ := kio.Read(f, 4)
	ti := 0
	if !kio.Bequal(tb, []byte("KV4H")) {
		return nil, nil, nil
	}
	tb, _ = kio.Read(f, 8)
	ti = kobj.Decode(tb)
	encmh, _ := kio.Read(f, ti)
	tb, _ = kio.Read(f, 8)
	ti = kobj.Decode(tb)
	encfs, _ := kio.Read(f, ti)
	tb, _ = kio.Read(f, 8)
	ti = kobj.Decode(tb)
	encfk, _ := kio.Read(f, ti)
	return encmh, encfs, encfk
}

// make encfile by path, fptr, fkey
func (tbox *g4work) fpush(path string, fptr []byte, fkey []byte) error {
	maxnum := 500 // arbitrarly fixed maxnum
	temp := 0
	fs, err := os.ReadDir(fmt.Sprintf("%s%d/", *tbox.path, temp))
	for len(fs) >= maxnum && err == nil {
		temp = temp + 1
		fs, err = os.ReadDir(fmt.Sprintf("%s%d/", *tbox.path, temp))
	}
	if err != nil {
		os.Mkdir(fmt.Sprintf("%s%d/", *tbox.path, temp), os.ModePerm)
	}
	(*tbox.fnums)[kio.Bprint(fptr)] = temp
	(*tbox.fkeys)[kio.Bprint(fptr)] = fkey
	var wk G4kaesfunc
	wk.Inbuf.OpenF(path, true)
	defer wk.Inbuf.CloseF()
	wk.Exbuf.OpenF(fmt.Sprintf("%s%d/%s.kv4", *tbox.path, temp, kb64en(fptr)), false)
	defer wk.Exbuf.CloseF()
	fkcpy := append(make([]byte, 0), fkey...)
	err = wk.Encrypt(fkcpy)
	return err
}

// get decfile by newpath, fptr, fkey
func (tbox *g4work) fpop(newpath string, fptr []byte, fkey []byte) error {
	path := fmt.Sprintf("%s%d/%s.kv4", *tbox.path, (*tbox.fnums)[kio.Bprint(fptr)], kb64en(fptr))
	var wk G4kaesfunc
	wk.Inbuf.OpenF(path, true)
	defer wk.Inbuf.CloseF()
	wk.Exbuf.OpenF(newpath, false)
	defer wk.Exbuf.CloseF()
	err := wk.Decrypt(fkey)
	return err
}

// make filesys, filekey
func (tbox *g4work) mkhead(binf *g4dir, mainf *g4dir) ([]byte, []byte) {
	fs := append(append(append(binf.gen(0), 10), mainf.gen(0)...), 10)
	temp := make([]string, 0)
	tbs := make([][]byte, 2*len(*tbox.fkeys))
	for i := range *tbox.fkeys {
		temp = append(temp, i)
	}
	sort.Strings(temp)
	for i, r := range temp {
		tbs[2*i], _ = kio.Bread(r)
		tbs[2*i+1] = (*tbox.fkeys)[r]
	}
	fk := bytes.Join(tbs, nil)
	return fs, fk
}

// interpret filesys, filekey data
func (tbox *g4work) rdhead(fs []byte, fk []byte) (*g4dir, *g4dir) {
	temp := bytes.Split(fs, []byte("\n"))
	if len(temp[len(temp)-1]) == 0 {
		temp = temp[0 : len(temp)-1]
	}
	depth := make([]int, len(temp))
	var ti int
	for i, r := range temp {
		ti = 0
		for 48 <= r[ti] && r[ti] < 58 {
			ti = ti + 1
		}
		ti, err := strconv.Atoi(string(r[0:ti]))
		if err != nil {
			ti = 1
		}
		depth[i] = ti
	}
	ti = 1
	for depth[ti] > 0 {
		ti = ti + 1
	}
	var binf g4dir
	binf.set(depth[0:ti], temp[0:ti])
	var mainf g4dir
	mainf.set(depth[ti:], temp[ti:])

	for i := 0; i < len(fk)/60; i++ {
		ti = 60 * i
		(*tbox.fkeys)[kio.Bprint(fk[ti:ti+12])] = fk[ti+12 : ti+60]
	}
	return &binf, &mainf
}

// make g4dir with folder, path (~/)
func (tbox *g4work) encdir(path string) (*g4dir, error) {
	path = path[0 : len(path)-1]
	var out g4dir
	out.name = path[strings.LastIndex(path, "/")+1:]
	out.subdir = make([]g4dir, 0)
	out.subfile = make([]g4file, 0)

	fs, err := os.ReadDir(path + "/")
	if err != nil {
		return &out, err
	}
	for _, r := range fs {
		name := r.Name()
		if name[len(name)-1] == '/' {
			name = name[0 : len(name)-1]
		}
		if r.IsDir() {
			td, err := tbox.encdir(path + "/" + name + "/")
			if err != nil {
				return &out, err
			}
			out.subdir = append(out.subdir, *td)
		} else {
			regen := true
			var fptr []byte
			for regen {
				fptr = genrand(12)
				_, regen = (*tbox.fkeys)[kio.Bprint(fptr)]
				if bytes.Contains(fptr, []byte("\n")) {
					regen = true
				}
			}
			fkey := genrand(48)
			err = tbox.fpush(path+"/"+name, fptr, fkey)
			if err != nil {
				return &out, err
			}
			var tf g4file
			tf.name = name
			tf.fptr = fptr
			out.subfile = append(out.subfile, tf)
		}
	}
	return &out, nil
}

// make folder with g4dir, root (~/)
func (tbox *g4work) decdir(dir *g4dir, root string) error {
	path := root + dir.name + "/"
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}
	for _, r := range dir.subdir {
		err = tbox.decdir(&r, path)
		if err != nil {
			return err
		}
	}
	for _, r := range dir.subfile {
		err = tbox.fpop(path+r.name, r.fptr, (*tbox.fkeys)[kio.Bprint(r.fptr)])
		if err != nil {
			return err
		}
	}
	return nil
}

// gen4 kv4
type g4kv4 struct {
	Path string // cluster path (~/)
	Hint []byte // cluster hint

	pwhash   []byte
	salt     []byte
	akeydata []byte
	tkeydata []byte
	encfs    []byte
	encfk    []byte
	worker   g4work

	bins  g4dir             // bin/ dir
	mains g4dir             // main/ dir
	fkeys map[string][]byte // kio.Bprint fptr-fkey
	fnums map[string]int    // kio.Bprint fptr-chunknum
}

// make mainhead by pw, kf, hint, akey 48B, tkey 48B
func (tbox *g4kv4) mkhead(pw []byte, kf []byte, akey []byte, tkey []byte) []byte {
	tbox.salt = genrand(32)
	tb := append(append(append(append(append(make([]byte, 0), pw...), pw...), kf...), kf...), pw...)
	tbox.pwhash = hashscrypt(tb, tbox.salt, 524288, 8, 1, 256)
	tb = append(append(append(append(append(make([]byte, 0), kf...), pw...), pw...), kf...), pw...)
	mkey := hashscrypt(tb, tbox.salt, 16384, 8, 1, 48)
	mkcpy := append(make([]byte, 0), mkey...)
	tbox.akeydata = aes16en(akey, mkey[0:32], mkey[32:48])
	tbox.tkeydata = aes16en(tkey, mkcpy[0:32], mkcpy[32:48])

	tg0 := make(map[string]G4data)
	var tg1 G4data
	tg1.Set("KV4")
	tg0["MODE"] = tg1
	var tg2 G4data
	tg2.Set(tbox.salt)
	tg0["SALT"] = tg2
	var tg3 G4data
	tg3.Set(tbox.pwhash)
	tg0["PWH"] = tg3
	var tg4 G4data
	tg4.Set(tbox.akeydata)
	tg0["AKDT"] = tg4
	var tg5 G4data
	tg5.Set(tbox.tkeydata)
	tg0["TKDT"] = tg5
	var tg6 G4data
	tg6.Set(tbox.Hint)
	tg0["HINT"] = tg6
	return []byte(G4DBwrite(tg0))
}

// read mainhead, gen akdt, tkdt, hint, salt
func (tbox *g4kv4) rdhead(data []byte) error {
	header := G4DBread(string(data))
	if header["MODE"].StrV != "KV4" {
		return errors.New("invalid encryption mode")
	}
	tbox.salt = header["SALT"].ByteV
	tbox.pwhash = header["PWH"].ByteV
	tbox.akeydata = header["AKDT"].ByteV
	tbox.tkeydata = header["TKDT"].ByteV
	tbox.Hint = header["HINT"].ByteV
	return nil
}

// gen4 kv4 view cluster
func (tbox *g4kv4) View() error {
	encmh, encfs, encfk := tbox.worker.hpop()
	if len(encmh) == 0 {
		return errors.New("invalid cluster")
	}
	tbox.encfs = encfs
	tbox.encfk = encfk
	fs, err := os.ReadDir(tbox.Path)
	if err != nil {
		return err
	}
	for _, r := range fs {
		ti, err := strconv.Atoi(strings.Replace(strings.Replace(r.Name(), "/", "", -1), "\\", "", -1))
		if err == nil {
			tbox.worker.getfptr(ti)
		}
	}
	return tbox.rdhead(encmh)
}

// gen4 kv4 read cluster, gen bin/, main/ in newpath
func (tbox *g4kv4) Read(pw []byte, kf []byte, newpath string) error {
	if len(tbox.salt) == 0 {
		return errors.New("should done View() first")
	}
	newpath = kio.Abs(newpath)
	tb := append(append(append(append(append(make([]byte, 0), pw...), pw...), kf...), kf...), pw...)
	pwhcmp := hashscrypt(tb, tbox.salt, 524288, 8, 1, 256)
	tb = append(append(append(append(append(make([]byte, 0), kf...), pw...), pw...), kf...), pw...)
	mkey := hashscrypt(tb, tbox.salt, 16384, 8, 1, 48)
	mkcpy := append(make([]byte, 0), mkey...)
	if !kio.Bequal(pwhcmp, tbox.pwhash) {
		return errors.New("InvalidPW")
	}

	akey := aes16de(tbox.akeydata, mkey[0:32], mkey[32:48])
	tkey := aes16de(tbox.tkeydata, mkcpy[0:32], mkcpy[32:48])
	var wk G4kaesfunc
	wk.Inbuf.OpenB(tbox.encfs, true)
	wk.Exbuf.OpenB(make([]byte, 0), false)
	err := wk.Decrypt(tkey)
	if err != nil {
		return err
	}
	fs := wk.Exbuf.CloseB()
	wk.Inbuf.OpenB(tbox.encfk, true)
	wk.Exbuf.OpenB(make([]byte, 0), false)
	err = wk.Decrypt(akey)
	if err != nil {
		return err
	}
	fk := wk.Exbuf.CloseB()

	binf, mainf := tbox.worker.rdhead(fs, fk)
	tbox.bins = *binf
	tbox.mains = *mainf
	err = tbox.worker.decdir(&tbox.bins, newpath)
	if err != nil {
		return err
	}
	return tbox.worker.decdir(&tbox.mains, newpath)
}

// gen4 kv4 write cluster
func (tbox *g4kv4) Write(pw []byte, kf []byte, tgtpath string) error {
	tbox.bins.name = "bin"
	tbox.bins.subdir = make([]g4dir, 0)
	tbox.bins.subfile = make([]g4file, 0)
	tbox.mains.name = "main"
	tbox.mains.subdir = make([]g4dir, 0)
	tbox.mains.subfile = make([]g4file, 0)
	tgtpath = kio.Abs(tgtpath)
	temp, err := tbox.worker.encdir(tgtpath)
	if err != nil {
		return err
	}
	tbox.mains.subdir = append(tbox.mains.subdir, *temp)

	fs, fk := tbox.worker.mkhead(&tbox.bins, &tbox.mains)
	akey := genrand(48)
	tkey := genrand(48)
	header := tbox.mkhead(pw, kf, akey, tkey)

	var wk G4kaesfunc
	wk.Inbuf.OpenB(fs, true)
	wk.Exbuf.OpenB(make([]byte, 0), false)
	err = wk.Encrypt(tkey)
	if err != nil {
		return err
	}
	tbox.encfs = wk.Exbuf.CloseB()
	wk.Inbuf.OpenB(fk, true)
	wk.Exbuf.OpenB(make([]byte, 0), false)
	err = wk.Encrypt(akey)
	if err != nil {
		return err
	}
	tbox.encfk = wk.Exbuf.CloseB()
	tbox.worker.hpush(header, tbox.encfs, tbox.encfk)
	return nil
}

// gen4 kv4 init g4kv4 struct, cluster path
func InitKV4(path string) *g4kv4 {
	var out g4kv4
	out.Path = kio.Abs(path)
	out.fnums = make(map[string]int)
	out.fkeys = make(map[string][]byte)
	out.worker.path = &out.Path
	out.worker.fnums = &out.fnums
	out.worker.fkeys = &out.fkeys
	return &out
}
