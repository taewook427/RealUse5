// test693 : stdlib5.kvault

package kvault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"slices"
	"sort"
	"stdlib5/kaes"
	"stdlib5/kdb"
	"stdlib5/kio"
	"stdlib5/kobj"
	"stdlib5/ksc"
	"strconv"
	"strings"
	"sync"
	"time"

	mrand "math/rand"

	"golang.org/x/crypto/scrypt"
)

// go get "golang.org/x/crypto/sha3"
// go get "golang.org/x/crypto/scrypt"

// ===== section ===== logging

// log with limited size
type logger struct {
	abort    bool // emergency abort flag
	working  bool // process working flag
	readonly bool // readonly status flag

	buffer [50]string // log data (not including \n)
	start  int        // queue start
	size   int        // queue size

	lock sync.Mutex
}

// read log, joined with \n
func (lg *logger) read() string {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	if lg.size > 50 {
		lg.size = 50
	}
	if lg.size == 0 {
		return ""
	}
	temp := ""
	for i := 0; i < lg.size; i++ {
		temp = temp + lg.buffer[(lg.start+i)%50] + "\n"
	}
	return temp[0 : len(temp)-1]
}

// write log
func (lg *logger) write(data string) {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	data = strings.Replace(data, "\n", " ", -1)
	if lg.size < 50 {
		lg.buffer[(lg.start+lg.size)%50] = data
		lg.size = lg.size + 1
	} else {
		lg.size = 50
		lg.buffer[lg.start] = data
		lg.start = (lg.start + 1) % 50
	}
}

// clear log
func (lg *logger) clear() {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	for i := range lg.buffer {
		lg.buffer[i] = ""
	}
	lg.start = 0
	lg.size = 0
}

// ===== section ===== safe IO / encryption, always success work

// actual io / encryption layer
type basework struct {
	log   *logger // common logging system
	sleep int     // sleeping time if error occurs
}

// struct init
func (bw *basework) init(log *logger, sleep int) {
	bw.log = log
	bw.sleep = sleep
}

// file read, pos 0 or +, size -1 or +
func (bw *basework) read(path string, pos int, size int) []byte {
	defer func() {
		if ferr := recover(); ferr != nil {
			bw.log.write(fmt.Sprintf("critical : %s", ferr))
			bw.log.abort = true
		}
	}()

	flag := true
	var f *os.File
	var err error
	var data []byte
	for flag {
		if bw.log.abort { // step 0 : abort check
			bw.log.write("msg : IOread abort")
			flag = false
			break
		}

		f, err = os.Open(path) // step 1 : file open
		if err != nil {
			f.Close()
			bw.log.write(fmt.Sprintf("err : file open fail -%s", err))
			time.Sleep(time.Second * time.Duration(bw.sleep))
			continue
		}

		_, err = f.Seek(int64(pos), 0) // step 2 : position move
		if err != nil {
			f.Close()
			bw.log.write(fmt.Sprintf("err : file seek fail -%s", err))
			time.Sleep(time.Second * time.Duration(bw.sleep))
			continue
		}

		data, err = kio.Read(f, size) // step 3 : data IO
		if err != nil {
			f.Close()
			bw.log.write(fmt.Sprintf("err : file read fail -%s", err))
			time.Sleep(time.Second * time.Duration(bw.sleep))
			continue
		}

		f.Close()
		flag = false
	}
	return data
}

// file write, pos 0 or +
func (bw *basework) write(path string, pos int, data []byte) {
	defer func() {
		if ferr := recover(); ferr != nil {
			bw.log.write(fmt.Sprintf("critical : %s", ferr))
			bw.log.abort = true
		}
	}()

	flag := true
	var f *os.File
	var err error
	for flag {
		if bw.log.abort { // step 0 : abort check
			bw.log.write("msg : IOwrite abort")
			flag = false
			break
		}

		f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666) // step 1 : file open
		if err != nil {
			f.Close()
			bw.log.write(fmt.Sprintf("err : file open fail -%s", err))
			time.Sleep(time.Second * time.Duration(bw.sleep))
			continue
		}

		_, err = f.Seek(int64(pos), 0) // step 2 : position move
		if err != nil {
			f.Close()
			bw.log.write(fmt.Sprintf("err : file seek fail -%s", err))
			time.Sleep(time.Second * time.Duration(bw.sleep))
			continue
		}

		_, err = kio.Write(f, data) // step 3 : data IO
		if err != nil {
			f.Close()
			bw.log.write(fmt.Sprintf("err : file write fail -%s", err))
			time.Sleep(time.Second * time.Duration(bw.sleep))
			continue
		}

		f.Close()
		flag = false
	}
}

// generate / delete folder, !!! generated folder turns empty !!!
func (bw *basework) dirctrl(path string, create bool) {
	defer func() {
		if ferr := recover(); ferr != nil {
			bw.log.write(fmt.Sprintf("critical : %s", ferr))
			bw.log.abort = true
		}
	}()

	if create {
		flag := true
		for flag {
			if bw.log.abort { // step 0 : abort check
				bw.log.write("msg : IOdir creation abort")
				flag = false
				break
			}

			os.RemoveAll(path)                 // step 1 : remove
			err := os.Mkdir(path, os.ModePerm) // step 2 : generate
			if err == nil {
				flag = false
			} else {
				bw.log.write(fmt.Sprintf("err : folder generate fail -%s", err))
				time.Sleep(time.Second * time.Duration(bw.sleep))
			}
		}

	} else {
		flag := true
		for flag {
			if bw.log.abort { // step 0 : abort check
				bw.log.write("msg : IOdir delete abort")
				flag = false
				break
			}

			err := os.RemoveAll(path) // step 1 : remove
			if err == nil {
				flag = false
			} else {
				bw.log.write(fmt.Sprintf("err : folder delete fail -%s", err))
				time.Sleep(time.Second * time.Duration(bw.sleep))
			}
		}
	}
}

// get names of subobject in folder
func (bw *basework) dirsub(path string) []string {
	defer func() {
		if ferr := recover(); ferr != nil {
			bw.log.write(fmt.Sprintf("critical : %s", ferr))
			bw.log.abort = true
		}
	}()

	var out []string
	flag := true
	for flag {
		if bw.log.abort { // step 0 : abort check
			bw.log.write("msg : IOdir sublist abort")
			flag = false
			break
		}

		out = make([]string, 0)
		tgts, err := os.ReadDir(path) // step 1 : open dir
		if err != nil {
			bw.log.write(fmt.Sprintf("err : folder info fetch fail -%s", err))
			time.Sleep(time.Second * time.Duration(bw.sleep))
			continue
		}

		for _, r := range tgts {
			temp := r.Name() // step 2 : get name
			if r.IsDir() {
				temp = strings.Replace(temp, "\\", "/", -1)
				if temp[len(temp)-1] != '/' {
					temp = temp + "/"
				}
			}
			out = append(out, temp)
		}
		flag = false
	}
	return out
}

// simple aes-256 calc, key 48B (iv16 key32) / isenc T:encrypt F:decrypt / ispad T:dopad F:nopad
func (bw *basework) aescalc(data []byte, key [48]byte, isenc bool, ispad bool) []byte {
	if len(data)%16 != 0 {
		if !isenc {
			bw.log.write("err : invalid decrypt data length")
			return nil
		} else if !ispad {
			bw.log.write("err : invalid encrypt data length")
			return nil
		}
	}

	var out []byte
	block, _ := aes.NewCipher(key[16:48])
	if isenc {
		encrypter := cipher.NewCBCEncrypter(block, key[0:16])
		if ispad {
			plen := 16 - (len(data) % 16)
			for i := 0; i < plen; i++ {
				data = append(data, byte(plen))
			}
		}
		out = make([]byte, len(data))
		encrypter.CryptBlocks(out, data)

	} else {
		decrypter := cipher.NewCBCDecrypter(block, key[0:16])
		out = make([]byte, len(data))
		decrypter.CryptBlocks(out, data)
		if ispad {
			plen := int(out[len(out)-1])
			out = out[0 : len(out)-plen]
		}
	}
	return out
}

// generate pwhash 192B, mkey 144B
func (bw *basework) genpm(pw []byte, kf []byte, salt []byte) ([192]byte, [144]byte) {
	tb := append(append(append(append(append(make([]byte, 0), pw...), pw...), kf...), pw...), kf...)
	pwh, _ := scrypt.Key(tb, salt, 524288, 8, 1, 192)
	tb = append(append(append(append(append(make([]byte, 0), kf...), pw...), kf...), kf...), pw...)
	mkey, _ := scrypt.Key(tb, salt, 16384, 8, 1, 144)
	return [192]byte(pwh), [144]byte(mkey)
}

// ===== section ===== KOS Virtual File System

// data writing format : (all B is written in bprint)
// folder - (name)\n(depth 4B->8)(time 8B->16)\n - if locked depth + 1000000000
// file - (name)\n(time 8B->16)(size 8B->16)(fptr 5B->10)\n
// folder -> direct files -> subfolders

// []byte conversion, B - hex
func vconv(data []byte, tohex bool) []byte {
	var out []byte
	if tohex {
		out = make([]byte, len(data)*2)
		hex.Encode(out, data)
	} else {
		out = make([]byte, len(data)/2)
		hex.Decode(out, data)
	}
	return out
}

// generate VFS from data bytes, returns nil if error
func vread(data []byte) *vdir {
	frag := bytes.Split(data, []byte{10})
	if len(frag) < 3 || len(frag)%2 == 0 {
		return nil
	}
	frag = frag[0 : len(frag)-1] // 2N

	isfolder := make([]bool, len(frag)/2) // N
	islocked := make([]bool, len(frag)/2) // N
	depth := make([]int, len(frag)/2)     // N
	curdepth := -1
	for i := 0; i < len(frag)/2; i++ {
		tb := frag[2*i]
		if tb[len(tb)-1] == 47 {
			isfolder[i] = true
			ti := kobj.Decode(vconv(frag[2*i+1][0:8], false))
			if ti < 1000000000 {
				islocked[i] = false
				depth[i] = ti
			} else {
				islocked[i] = true
				depth[i] = ti - 1000000000
			}
			curdepth = depth[i]

		} else {
			isfolder[i] = false
			islocked[i] = false
			depth[i] = curdepth + 1
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var out vdir // ex) .{ dir file file { dir {dir} } {dir file} }
	go out.read_inline(frag, isfolder, islocked, depth, &wg)
	wg.Wait()
	return &out
}

// virtual directory
type vdir struct {
	name     string // folder name (~/)
	time     int    // folder generation time
	islocked bool   // is folder protected

	// ordered by name
	subdir  []vdir  // sub folders
	subfile []vfile // sub files
}

// generate data bytes from KVFS_folder, can return nil
func (vfs *vdir) write(wrlocked bool) []byte {
	recv := make(chan []byte, 1)
	go vfs.write_inline(0, wrlocked, recv)
	return <-recv
}

// count size, folders, files
func (vfs *vdir) count(wrlocked bool) (int, int, int) {
	ret0 := make(chan int, 1)
	ret1 := make(chan int, 1)
	ret2 := make(chan int, 1)
	go vfs.count_inline(ret0, ret1, ret2, wrlocked)
	return <-ret0, <-ret1, <-ret2
}

// subdir, subfile ordering
func (vfs *vdir) sort() {
	sort.Slice(vfs.subdir, func(i int, j int) bool {
		return vfs.subdir[i].name < vfs.subdir[j].name
	})
	sort.Slice(vfs.subfile, func(i int, j int) bool {
		return vfs.subfile[i].name < vfs.subfile[j].name
	})
}

// count size, folders, files (inline)
func (vfs *vdir) count_inline(ret0 chan int, ret1 chan int, ret2 chan int, wrlocked bool) {
	defer func() {
		if err := recover(); err != nil {
			ret0 <- 0
			ret1 <- 0
			ret2 <- 0
		}
	}()
	if vfs.islocked && !wrlocked {
		ret0 <- 0
		ret1 <- 0
		ret2 <- 0
	} else {
		size := 0
		folder := 1
		file := len(vfs.subfile)
		for _, r := range vfs.subfile {
			size = size + r.size
		}

		recv0 := make([]chan int, len(vfs.subdir))
		recv1 := make([]chan int, len(vfs.subdir))
		recv2 := make([]chan int, len(vfs.subdir))
		for i, r := range vfs.subdir {
			recv0[i] = make(chan int, 1)
			recv1[i] = make(chan int, 1)
			recv2[i] = make(chan int, 1)
			go r.count_inline(recv0[i], recv1[i], recv2[i], wrlocked)
		}
		for i := 0; i < len(vfs.subdir); i++ {
			size = size + <-recv0[i]
			folder = folder + <-recv1[i]
			file = file + <-recv2[i]
		}

		ret0 <- size
		ret1 <- folder
		ret2 <- file
	}
}

// set internal data by bytes fragment (frag 2N, isfolder N, islocked N, depth N)
func (vfs *vdir) read_inline(frag [][]byte, isfolder []bool, islocked []bool, depth []int, wg *sync.WaitGroup) {
	defer func() {
		recover()
		wg.Done()
	}()
	if len(frag) == 2*len(depth) && len(isfolder) == len(depth) && len(islocked) == len(depth) {
		// current folder
		vfs.name = string(frag[0])
		vfs.time = kobj.Decode(vconv(frag[1][8:24], false))
		vfs.islocked = islocked[0]
		vfs.subdir = make([]vdir, 0)
		vfs.subfile = make([]vfile, 0)

		curdepth := depth[0]
		curpos := 1
		endpos := len(depth)
		flag := true

		// direct sub files
		for flag {
			if curpos == endpos {
				flag = false
			} else if isfolder[curpos] {
				break
			} else {
				var ftemp vfile
				ftemp.name = string(frag[2*curpos])
				ftemp.time = kobj.Decode(vconv(frag[2*curpos+1][0:16], false))
				ftemp.size = kobj.Decode(vconv(frag[2*curpos+1][16:32], false))
				ftemp.fptr = kobj.Decode(vconv(frag[2*curpos+1][32:42], false))
				vfs.subfile = append(vfs.subfile, ftemp)
				curpos = curpos + 1
			}
		}

		// sub folder, multithread
		mempos := curpos // sub folder start
		curpos = curpos + 1
		toreadpos := make([]int, 0) // mempos, curpos
		for flag {
			if curpos == endpos {
				toreadpos = append(append(toreadpos, mempos), curpos)
				flag = false
			} else if curdepth+1 == depth[curpos] {
				toreadpos = append(append(toreadpos, mempos), curpos)
				mempos = curpos
				curpos = curpos + 1
			} else {
				curpos = curpos + 1
			}
		}

		// multithread read
		vfs.subdir = make([]vdir, len(toreadpos)/2) // !! use make() to set memory !!
		var wgret sync.WaitGroup
		wgret.Add(len(toreadpos) / 2)
		for i := 0; i < len(toreadpos)/2; i++ {
			v0 := toreadpos[2*i]
			v1 := toreadpos[2*i+1]
			go vfs.subdir[i].read_inline(frag[2*v0:2*v1], isfolder[v0:v1], islocked[v0:v1], depth[v0:v1], &wgret)
		}
		wgret.Wait()
	}
}

// generate writing data with given depth, depth should be under 900000000
func (vfs *vdir) write_inline(depth int, wrlocked bool, ret chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			ret <- nil
		}
	}()
	if (depth > 900000000) || (!wrlocked && vfs.islocked) {
		ret <- nil
	} else {

		// current folder data
		mem := make([]byte, 0, 1024*len(vfs.subdir)+1024)
		mem = append(mem, []byte(vfs.name+"\n")...)
		if vfs.islocked {
			mem = append(mem, vconv(kobj.Encode(depth+1000000000, 4), true)...)
		} else {
			mem = append(mem, vconv(kobj.Encode(depth, 4), true)...)
		}
		mem = append(append(mem, vconv(kobj.Encode(vfs.time, 8), true)...), 10)

		// direct sub file data
		for _, r := range vfs.subfile {
			mem = append(mem, []byte(r.name+"\n")...)
			mem = append(append(mem, vconv(kobj.Encode(r.time, 8), true)...), vconv(kobj.Encode(r.size, 8), true)...)
			mem = append(append(mem, vconv(kobj.Encode(r.fptr, 5), true)...), 10)
		}

		// sub folder data, multithread
		recv := make([]chan []byte, len(vfs.subdir))
		for i, r := range vfs.subdir {
			recv[i] = make(chan []byte, 1)
			go r.write_inline(depth+1, wrlocked, recv[i])
		}
		for _, r := range recv {
			mem = append(mem, <-r...)
		}

		ret <- mem
	}
}

// virtual file
type vfile struct {
	name string // file name
	time int    // file generation time
	size int    // file size (encrypted)
	fptr int    // file pointer (5B uint)
}

// ===== section ===== fkey save & load

// fptr 5B int - fkey 48B set
type keynode struct {
	fptr int      // fptr int (5B little endian)
	fkey [48]byte // fkey 48B
}

// fptr - fkey storage
type keymap struct {
	data [4194304][]keynode // hash : fptr % 4194304
}

// search fptr, returns (pos0, pos1, fkey) / (-1, -1, zeros)
func (ks *keymap) seek(fptr int) (int, int, [48]byte) {
	hashed := fptr % 4194304 // hashed pos
	for i, r := range ks.data[hashed] {
		if r.fptr == fptr {
			return hashed, i, r.fkey
		}
	}
	return -1, -1, [48]byte{}
}

// push fkey, errors if fkey exists
func (ks *keymap) push(fptr int, fkey [48]byte) error {
	ext, _, _ := ks.seek(fptr)
	if ext != -1 {
		return errors.New("fkey already exists")
	}
	hashed := fptr % 4194304 // hashed pos
	var temp keynode
	temp.fptr = fptr
	temp.fkey = fkey
	ks.data[hashed] = append(ks.data[hashed], temp)
	return nil
}

// delete fkey, errors if fkey not exists
func (ks *keymap) pop(fptr int) error {
	pos0, pos1, _ := ks.seek(fptr)
	if pos0 == -1 {
		return errors.New("fkey not exists")
	}
	ks.data[pos0] = append(ks.data[pos0][:pos1], ks.data[pos0][pos1+1:]...)
	return nil
}

// read & set data from bytes (5B 48B) * n
func (ks *keymap) read(data []byte) error {
	if len(data)%53 != 0 {
		return errors.New("invalid save data")
	}
	for i := range ks.data {
		ks.data[i] = make([]keynode, 0)
	}
	for i := 0; i < len(data)/53; i++ {
		pos := 53 * i
		fptr := kobj.Decode(data[pos : pos+5])
		fkey := [48]byte(data[pos+5 : pos+53])
		ks.push(fptr, fkey)
	}
	return nil
}

// generate save data bytes (5B 48B) * n
func (ks *keymap) write() []byte {
	length := 0
	for _, r := range ks.data {
		length = length + len(r)
	}
	mem := make([]byte, 0, 53*length+1024)
	for _, r := range ks.data {
		for _, l := range r {
			mem = append(mem, kobj.Encode(l.fptr, 5)...)
			mem = append(mem, l.fkey[:]...)
		}
	}
	return mem
}

// ===== section ===== control block RW layer

// control block IO with buffer (L1 mem - L2 local - L3 remote), mem : 64 * 384KiB, local : 512 * 384KiB
type blockctrl struct {
	rawIO   *basework  // direct IO layer
	fphykey [4096]byte // fphy key 4096B
	local   string     // local storage path (~/)
	remote  string     // remote storage path (~/)

	l1data [64][393216]byte // control block data, 64 * 384KiB
	l1pos  [64]int          // control block id, -1 : not allocated
	l2pos  [512]int         // control block id, -1 : not allocated
}

// init & delete internal data, make local/cache/
func (bc *blockctrl) init(bw *basework, fphykey [4096]byte, local string, remote string) {
	bc.rawIO = bw
	bc.fphykey = fphykey
	bc.local = local
	bc.remote = remote
	bc.rawIO.log.write("msg : control block IO Layer init")

	for i := 0; i < 64; i++ {
		bc.l1data[i] = [393216]byte{}
		bc.l1pos[i] = -1
	}
	for i := 0; i < 512; i++ {
		bc.l2pos[i] = -1
	}
	bc.rawIO.dirctrl(bc.local+"cache/", true)
}

// L1 - (L2L3) complete read
func (bc *blockctrl) read(pos int) []byte {
	empty := -1 // not allocated index
	for i, r := range bc.l1pos {
		if r == pos { // allocated data exists
			return bc.l1data[i][:]
		} else if r == -1 {
			empty = i
		}
	}

	if empty == -1 { // empty index not exists
		empty = pos % 64
		bc.write_l1l2(bc.l1pos[empty], bc.l1data[empty]) // push
	}

	// allocated data not exists, empty index exists
	temp := bc.read_l1l2(pos)
	bc.l1data[empty] = temp
	bc.l1pos[empty] = pos
	return temp[:]
}

// L1 - (L2L3) complete write
func (bc *blockctrl) write(pos int, data []byte) {
	if len(data) != 393216 {
		bc.rawIO.log.write("critical : ctrl block invalid length push")
		bc.rawIO.log.abort = true
	} else {
		empty := -1 // not allocated index
		for i, r := range bc.l1pos {
			if r == pos { // allocated data exists
				bc.l1data[i] = [393216]byte(data)
				return
			} else if r == -1 {
				empty = i
			}
		}

		if empty == -1 { // empty index not exists
			empty = pos % 64
			bc.write_l1l2(bc.l1pos[empty], bc.l1data[empty]) // push
		}

		// allocated data not exists, empty index exists
		bc.l1data[empty] = [393216]byte(data)
		bc.l1pos[empty] = pos
	}
}

// L1 - (L2L3) complete sync, deepflush T : L1-L2-L3 F : L1-L2
func (bc *blockctrl) flush(deepflush bool) {
	bc.rawIO.log.write("msg : ctrl block flush -L1L2")
	for i, r := range bc.l1pos {
		if r != -1 {
			bc.write_l1l2(bc.l1pos[i], bc.l1data[i]) // push
		}
	}

	if deepflush {
		bc.rawIO.log.write("msg : ctrl block flush -L2L3")
		for _, r := range bc.l2pos {
			if r != -1 {
				bc.write_l2l3(r, [393216]byte(bc.rawIO.read(fmt.Sprintf("%scache/%dc.kv5", bc.local, r), 0, 393216))) // push
			}
		}
	}
}

// enc/dec control block 384KiB <-> 384KiB
func (bc *blockctrl) ende_inline(data [393216]byte, pos int, encmode bool) [393216]byte {
	ctrlkey := make([]byte, 0, 48)
	ti := 16 * (pos % 256)
	ctrlkey = append(ctrlkey, bc.fphykey[ti:ti+16]...)
	ti = 16 * ((pos / 256) % 256)
	ctrlkey = append(ctrlkey, bc.fphykey[ti:ti+16]...)
	ti = 16 * ((pos / 65536) % 256)
	ctrlkey = append(ctrlkey, bc.fphykey[ti:ti+16]...)
	return [393216]byte(bc.rawIO.aescalc(data[:], [48]byte(ctrlkey), encmode, false))
}

// L2 - L3 read, cached with readonly care -> encdata
func (bc *blockctrl) read_l2l3(pos int) [393216]byte {
	bc.rawIO.log.write(fmt.Sprintf("msg : remote ctrl block read -%d", pos))
	return [393216]byte(bc.rawIO.read(fmt.Sprintf("%s%d/%dc.kv5", bc.remote, pos/256, pos), 0, 393216))
}

// L2 - L3 write, cached with readonly care <- encdata
func (bc *blockctrl) write_l2l3(pos int, encdata [393216]byte) {
	if !bc.rawIO.log.readonly { // if not readonly
		bc.rawIO.log.write(fmt.Sprintf("msg : remote ctrl block write -%d", pos))
		bc.rawIO.write(fmt.Sprintf("%s%d/%dc.kv5", bc.remote, pos/256, pos), 0, encdata[:])
	}
}

// L1 - L2 read, cached -> plaindata
func (bc *blockctrl) read_l1l2(pos int) [393216]byte {
	bc.rawIO.log.write(fmt.Sprintf("msg : local ctrl block read -%d", pos))
	empty := -1 // not allocated index
	for i, r := range bc.l2pos {
		if r == pos { // allocated data exists
			return bc.ende_inline([393216]byte(bc.rawIO.read(fmt.Sprintf("%scache/%dc.kv5", bc.local, pos), 0, 393216)), pos, false)
		} else if r == -1 {
			empty = i
		}
	}

	if empty == -1 { // empty index not exists
		empty = pos % 512
		bc.write_l2l3(bc.l2pos[empty], [393216]byte(bc.rawIO.read(fmt.Sprintf("%scache/%dc.kv5", bc.local, bc.l2pos[empty]), 0, 393216))) // push
		os.Remove(fmt.Sprintf("%scache/%dc.kv5", bc.local, bc.l2pos[empty]))                                                              // delete
	}

	// allocated data not exists, empty index exists
	temp := bc.read_l2l3(pos)
	bc.rawIO.write(fmt.Sprintf("%scache/%dc.kv5", bc.local, pos), 0, temp[:])
	bc.l2pos[empty] = pos
	return bc.ende_inline(temp, pos, false)
}

// L1 - L2 write, cached <- plaindata
func (bc *blockctrl) write_l1l2(pos int, data [393216]byte) {
	bc.rawIO.log.write(fmt.Sprintf("msg : local ctrl block write -%d", pos))
	empty := -1 // not allocated index
	for i, r := range bc.l2pos {
		if r == pos { // allocated data exists
			temp := bc.ende_inline(data, pos, true)
			bc.rawIO.write(fmt.Sprintf("%scache/%dc.kv5", bc.local, pos), 0, temp[:])
			return
		} else if r == -1 {
			empty = i
		}
	}

	if empty == -1 { // empty index not exists
		empty = pos % 512
		bc.write_l2l3(bc.l2pos[empty], [393216]byte(bc.rawIO.read(fmt.Sprintf("%scache/%dc.kv5", bc.local, bc.l2pos[empty]), 0, 393216))) // push
		os.Remove(fmt.Sprintf("%scache/%dc.kv5", bc.local, bc.l2pos[empty]))                                                              // delete
	}

	// allocated data not exists, empty index exists
	temp := bc.ende_inline(data, pos, true)
	bc.rawIO.write(fmt.Sprintf("%scache/%dc.kv5", bc.local, pos), 0, temp[:])
	bc.l2pos[empty] = pos
}

// ===== section ===== abstract file IO layer

// !!! io target size should be bigger than 0 !!!
type fphyIO struct {
	local     string     // local path (~/)
	remote    string     // remote path (~/)
	chunksize int        // 4096 / 32768 / 262144 (4k, 32k, 256k)
	rawIO     *basework  // direct IO layer
	ctrlIO    *blockctrl // control block IO layer, !!! should flush manually !!!

	tgtsize   int   // tgt size
	comsize   int   // computed size
	emptyfptr int   // first empty fptr
	blocknum  int   // existing block num
	orders    []int // RW orders buffer (fptr int)
}

// generate crypto rand int [n, m)
func (fp *fphyIO) rand_inline(n int, m int) int {
	if n < m {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(m-n)))
		return int(num.Int64()) + n
	} else {
		return n
	}
}

// get block num by remote path
func (fp *fphyIO) bnum_inline() int {
	temp := fp.rawIO.dirsub(fp.remote)
	count := 0
	for slices.Contains(temp, fmt.Sprintf("%d/", count)) {
		count = count + 1
	}
	if count == 0 {
		return 0
	}

	temp = fp.rawIO.dirsub(fmt.Sprintf("%s%d/", fp.remote, count-1))
	count = (count - 1) * 256
	for slices.Contains(temp, fmt.Sprintf("%dd.kv5", count)) {
		count = count + 1
	}
	return count
}

// get first empty fptr, makes new block if no space left
func (fp *fphyIO) fptrnext_inline() int {
	// finds next empty fptr !! not including current fptr !!, returns at least (cur fptr + 1)
	tgt := fp.emptyfptr + 1
	flag := true
	for flag {
		if tgt/65536 < fp.blocknum { // inside block
			if fp.ctrlIO.read(tgt / 65536)[6*(tgt%65536)]%128 == 0 {
				flag = false
			} else {
				tgt = tgt + 1
			}
		} else { // out of block
			pos := tgt / 65536
			if !fp.rawIO.log.readonly {
				fp.rawIO.log.write(fmt.Sprintf("msg : new block generated -%d", pos))
				if pos%256 == 0 {
					fp.rawIO.dirctrl(fmt.Sprintf("%s%d/", fp.remote, pos/256), true)
				}

				temp := kaes.Genrand(393216)
				for i := 0; i < 65536; i++ {
					temp[6*i] = 0
				}
				temp_wr := fp.ctrlIO.ende_inline([393216]byte(temp), pos, true)
				fp.rawIO.write(fmt.Sprintf("%s%d/%dc.kv5", fp.remote, pos/256, pos), 0, temp_wr[:])

				path := fmt.Sprintf("%s%d/%dd.kv5", fp.remote, pos/256, pos)
				size := 2048 * fp.chunksize
				for i := 0; i < 32; i++ {
					fp.rawIO.write(path, size*i, kaes.Genrand(size))
				}

				fp.blocknum = fp.blocknum + 1
			}
			flag = false
		}
	}
	return tgt
}

// gc inline function, returns modified ctrl block / deleted chunk number
func (fp *fphyIO) gc_inline(data []byte, ret chan []byte, res chan int) {
	defer func() {
		if err := recover(); err != nil {
			fp.rawIO.log.write(fmt.Sprintf("critical : %s", err))
			fp.rawIO.log.abort = true
			ret <- data
			res <- 0
		}
	}()

	pos := 0
	count := 0
	for i := 0; i < len(data)/6; i++ {
		pos = 6 * i
		switch {
		case 0 < data[pos] && data[pos] < 128:
			data[pos] = 128
			count = count + 1
		case 128 < data[pos]:
			data[pos] = data[pos] - 128
		}
	}
	ret <- data
	res <- count
}

// init struct, empty fptr & block num : -1 to autoset else manual set
func (fp *fphyIO) init(local string, remote string, csize int, iodrv *basework, fphykey [4096]byte, empty int, block int) {
	iodrv.log.write("msg : PhysicalIO init")
	fp.local = local
	fp.remote = remote
	fp.chunksize = csize

	fp.rawIO = iodrv
	var ctrldrv blockctrl
	ctrldrv.init(iodrv, fphykey, local, remote)
	fp.ctrlIO = &ctrldrv

	fp.tgtsize = 0
	fp.comsize = 0
	if block < 0 {
		fp.blocknum = fp.bnum_inline()
	} else {
		fp.blocknum = block
	}
	if empty < 0 {
		fp.emptyfptr = -1
		fp.emptyfptr = fp.fptrnext_inline()
	} else {
		fp.emptyfptr = empty
	}
	fp.orders = make([]int, 0)
}

// clear internal data
func (fp *fphyIO) clear() {
	fp.tgtsize = 0
	fp.comsize = 0
	fp.orders = make([]int, 0)
}

// add file to FS, returns allocated fptr, -1 means failed, need to deep flush
func (fp *fphyIO) fpush(path string) int {
	// size, computed size, order buffer init
	fp.tgtsize = kio.Size(path)
	fp.comsize = 0
	if fp.tgtsize <= 0 {
		fp.rawIO.log.write("err : invalid fpush size")
		return -1
	}
	temp := make([]int, 0, 1024)
	if fp.rawIO.log.readonly {
		fp.rawIO.log.write("err : readonly enabled")
		return -1
	}

	// make space for chunk writing
	for fp.tgtsize > fp.comsize {
		temp = append(temp, fp.emptyfptr)
		fp.comsize = fp.comsize + fp.chunksize
		fp.emptyfptr = fp.fptrnext_inline()
	}

	// shuffle orders cut by length 256
	memnum := len(temp) / 256
	if len(temp)%256 != 0 {
		memnum = memnum + 1
	}
	memint := make([][]int, memnum)
	for i := 0; i < memnum; i++ {
		st := 256 * i
		ed := st + 256
		if ed > len(temp) {
			ed = len(temp)
		}
		memint[i] = temp[st:ed]
	}
	mrand.Shuffle(memnum, func(i int, j int) {
		memint[i], memint[j] = memint[j], memint[i]
	})
	fp.orders = make([]int, 0, 256*memnum)
	for _, r := range memint {
		fp.orders = append(fp.orders, r...)
	}
	if len(temp) != len(fp.orders) {
		fp.rawIO.log.write("err : orders generation fail")
		return -1
	}

	// write chunk data to data block
	fp.comsize = 0
	buffer := make([]int, 0)
	for _, r := range fp.orders[0 : len(fp.orders)-1] {
		if fp.rawIO.log.abort {
			fp.rawIO.log.write("msg : fpush abort")
			return -1
		}

		if len(buffer) == 0 {
			buffer = append(buffer, r)
		} else if buffer[len(buffer)-1]+1 == r && buffer[len(buffer)-1]/65536 == r/65536 {
			buffer = append(buffer, r)

		} else {
			pos := buffer[0] / 65536
			wrpos := (buffer[0] % 65536) * fp.chunksize
			dtsize := len(buffer) * fp.chunksize
			data := fp.rawIO.read(path, fp.comsize, dtsize)

			fp.rawIO.write(fmt.Sprintf("%s%d/%dd.kv5", fp.remote, pos/256, pos), wrpos, data)
			fp.comsize = fp.comsize + dtsize
			buffer = append(make([]int, 0), r)
		}
	}

	if len(buffer) != 0 {
		pos := buffer[0] / 65536
		wrpos := (buffer[0] % 65536) * fp.chunksize
		dtsize := len(buffer) * fp.chunksize
		data := fp.rawIO.read(path, fp.comsize, dtsize)

		fp.rawIO.write(fmt.Sprintf("%s%d/%dd.kv5", fp.remote, pos/256, pos), wrpos, data)
		fp.comsize = fp.comsize + dtsize
		buffer = nil
	}

	pos := fp.orders[len(fp.orders)-1] / 65536
	wrpos := (fp.orders[len(fp.orders)-1] % 65536) * fp.chunksize
	dtsize := fp.tgtsize - fp.comsize
	data := fp.rawIO.read(path, fp.comsize, dtsize)
	if dtsize != fp.chunksize {
		data = append(data, kaes.Genrand(fp.chunksize-dtsize)...)
	}
	fp.rawIO.write(fmt.Sprintf("%s%d/%dd.kv5", fp.remote, pos/256, pos), wrpos, data)
	fp.comsize = fp.comsize + dtsize
	if fp.tgtsize > fp.comsize {
		fp.rawIO.log.write("err : chunk writing error")
		return -1
	}

	// write control block
	fptrid := byte(fp.rand_inline(1, 128))
	for i, r := range fp.orders[0 : len(fp.orders)-1] {
		if fp.rawIO.log.abort {
			fp.rawIO.log.write("msg : fpush abort")
			return -1
		}

		wrpos := 6 * (r % 65536)
		tb := fp.ctrlIO.read(r / 65536)
		tb[wrpos] = fptrid
		copy(tb[wrpos+1:wrpos+6], kobj.Encode(fp.orders[i+1], 5))
		fp.ctrlIO.write(r/65536, tb)
	}
	wrpos = 6 * (fp.orders[len(fp.orders)-1] % 65536)
	tb := fp.ctrlIO.read(fp.orders[len(fp.orders)-1] / 65536)
	tb[wrpos] = fptrid
	copy(tb[wrpos+1:wrpos+6], kobj.Encode(fp.emptyfptr, 5))
	fp.ctrlIO.write(fp.orders[len(fp.orders)-1]/65536, tb)
	return fp.orders[0]
}

// read file from FS, should do Fchk first
func (fp *fphyIO) fpop(path string) {
	fp.comsize = 0
	buffer := make([]int, 0)
	for _, r := range fp.orders[0 : len(fp.orders)-1] {
		if fp.rawIO.log.abort {
			fp.rawIO.log.write("msg : fpop abort")
			break
		}

		if len(buffer) == 0 {
			buffer = append(buffer, r)
		} else if buffer[len(buffer)-1]+1 == r && buffer[len(buffer)-1]/65536 == r/65536 {
			buffer = append(buffer, r)

		} else {
			pos := buffer[0] / 65536
			wrpos := (buffer[0] % 65536) * fp.chunksize
			dtsize := len(buffer) * fp.chunksize
			data := fp.rawIO.read(fmt.Sprintf("%s%d/%dd.kv5", fp.remote, pos/256, pos), wrpos, dtsize)

			fp.rawIO.write(path, fp.comsize, data)
			fp.comsize = fp.comsize + dtsize
			buffer = append(make([]int, 0), r)
		}
	}

	if len(buffer) != 0 {
		pos := buffer[0] / 65536
		wrpos := (buffer[0] % 65536) * fp.chunksize
		dtsize := len(buffer) * fp.chunksize
		data := fp.rawIO.read(fmt.Sprintf("%s%d/%dd.kv5", fp.remote, pos/256, pos), wrpos, dtsize)

		fp.rawIO.write(path, fp.comsize, data)
		fp.comsize = fp.comsize + dtsize
		buffer = nil
	}

	pos := fp.orders[len(fp.orders)-1] / 65536
	wrpos := (fp.orders[len(fp.orders)-1] % 65536) * fp.chunksize
	dtsize := fp.tgtsize - fp.comsize
	data := fp.rawIO.read(fmt.Sprintf("%s%d/%dd.kv5", fp.remote, pos/256, pos), wrpos, fp.chunksize)[0:dtsize]
	fp.rawIO.write(path, fp.comsize, data)
	fp.comsize = fp.comsize + dtsize
	if fp.tgtsize != fp.comsize {
		fp.rawIO.log.write("critical : Fpop size error")
		fp.rawIO.log.abort = true
	}
}

// delete file in FS, should do Fchk first, need to deep flush
func (fp *fphyIO) fdel() {
	if !fp.rawIO.log.readonly {
		pos := fp.orders[0] / 65536
		wrpos := (fp.orders[0] % 65536) * fp.chunksize
		fp.rawIO.write(fmt.Sprintf("%s%d/%dd.kv5", fp.remote, pos/256, pos), wrpos, kaes.Genrand(fp.chunksize))

		for _, r := range fp.orders {
			if fp.rawIO.log.abort {
				fp.rawIO.log.write("msg : fdel abort")
				break
			}

			pos = r / 65536
			wrpos = 6 * (r % 65536)
			tb := fp.ctrlIO.read(pos)
			tb[wrpos] = 128
			fp.ctrlIO.write(pos, tb)
		}

		ti := slices.Min(fp.orders)
		if ti < fp.emptyfptr {
			fp.emptyfptr = ti
		}
	}
}

// fetch & generate orders, check fptr is valid
// change T : GC mode F : fpop fdel mode, result T : valid fptr chain F : invalid chain
func (fp *fphyIO) fchk(fsize int, fptr int, change bool) bool {
	fp.tgtsize = fsize
	fp.comsize = 0
	fp.orders = make([]int, 0)
	if fp.rawIO.log.readonly && change {
		return false
	}

	// set basic info
	nextpos := fptr
	pos := nextpos / 65536
	wrpos := 6 * (nextpos % 65536)
	if pos >= fp.blocknum {
		return false
	}
	fptrid := fp.ctrlIO.read(pos)[wrpos]
	if fptrid == 0 || fptrid == 128 {
		return false
	}

	// loop until end or error
	for fp.tgtsize > fp.comsize {
		pos = nextpos / 65536
		wrpos = 6 * (nextpos % 65536)
		if pos >= fp.blocknum {
			return false
		}
		tb := fp.ctrlIO.read(pos)
		if tb[wrpos] != fptrid {
			return false
		}
		fp.orders = append(fp.orders, nextpos)
		nextpos = kobj.Decode(tb[wrpos+1 : wrpos+6])
		fp.comsize = fp.comsize + fp.chunksize
		if change && fptrid < 128 {
			tb[wrpos] = fptrid + 128
			fp.ctrlIO.write(pos, tb)
		}
	}
	return true
}

// garbage collect / returns deleted chunk num, rewrite T : rewrite GC (128 -> 0) F : normal GC (new chunk to 128), need to deep flush
func (fp *fphyIO) fgc(rewrite bool) int {
	count := 0
	if !fp.rawIO.log.readonly {

		if rewrite { // rewrite GC (rewrite deleted but remaining chunk)
			for i := 0; i < fp.blocknum; i++ {
				if fp.rawIO.log.abort {
					fp.rawIO.log.write("msg : fpop abort")
					break
				}

				tb := fp.ctrlIO.read(i)
				for j := 0; j < 65536; j++ {
					if tb[6*j] == 128 {
						fp.rawIO.write(fmt.Sprintf("%s%d/%dd.kv5", fp.remote, i/256, i), fp.chunksize*j, kaes.Genrand(fp.chunksize))
						tb[6*j] = 0
						count = count + 1
					}
				}
				fp.ctrlIO.write(i, tb)
				if i%32 == 0 {
					fp.rawIO.log.write(fmt.Sprintf("msg : rewrite GC -%d", i))
				}
			}

		} else { // plain GC (clean unused chunks)
			var data [32][]byte
			var buffer [32]chan []byte
			var result [32]chan int
			for i := 0; i < 32; i++ {
				buffer[i] = make(chan []byte, 1)
				result[i] = make(chan int, 1)
			}
			for i := 0; i < fp.blocknum/32; i++ {
				if fp.rawIO.log.abort {
					fp.rawIO.log.write("msg : fpop abort")
					break
				}

				for j := 0; j < 32; j++ {
					data[j] = fp.ctrlIO.read(32*i + j)
				}
				for j := 0; j < 32; j++ {
					go fp.gc_inline(data[j], buffer[j], result[j])
				}
				for j := 0; j < 32; j++ {
					fp.ctrlIO.write(32*i+j, <-buffer[j])
					count = count + <-result[j]
				}
				fp.rawIO.log.write(fmt.Sprintf("msg : plain GC -%d", 32*i))
			}

			i := fp.blocknum / 32
			for j := 0; j < fp.blocknum%32; j++ {
				data[j] = fp.ctrlIO.read(32*i + j)
			}
			for j := 0; j < fp.blocknum%32; j++ {
				go fp.gc_inline(data[j], buffer[j], result[j])
			}
			for j := 0; j < fp.blocknum%32; j++ {
				fp.ctrlIO.write(32*i+j, <-buffer[j])
				count = count + <-result[j]
			}
		}

		fp.emptyfptr = -1
		fp.emptyfptr = fp.fptrnext_inline()
	}
	return count
}

// flush ctrl block, deepflush T : mem - local - remote F : mem - local
func (fp *fphyIO) flush(deepflush bool) {
	if fp.rawIO.log.readonly {
		fp.ctrlIO.flush(false)
	} else {
		fp.ctrlIO.flush(deepflush)
	}
}

// read header, returns nil if not exists
func (fp *fphyIO) hpop(path string) []byte {
	fp.rawIO.log.write(fmt.Sprintf("msg : reading header -%s", path))
	if _, err := os.Stat(path); err == nil {
		return fp.rawIO.read(path, 0, -1)
	} else {
		return nil
	}
}

// write header, autogen bck file, not work when (isremote T && readonly T)
func (fp *fphyIO) hpush(path string, data []byte, isremote bool) {
	if !isremote || !fp.rawIO.log.readonly {
		fp.rawIO.log.write(fmt.Sprintf("msg : writing header -%s", path))
		if _, err := os.Stat(path); err == nil {
			os.Remove(path + ".bck")
			fp.rawIO.write(path+".bck", 0, fp.rawIO.read(path, 0, -1))
		}
		os.Remove(path)
		fp.rawIO.write(path, 0, data)
	}
}

// ===== section ===== pevfs basic structure

// block A (account) KSC (4 section) [ KV5a, CRC32(section 0), CRC32(ClusterName) ]
// 0 : KDB text
//     -> (salt@B hint@B pwhash@B AccountName@S wrsign@B fsys_enckey@B fkey_enckey@B fphy_enckey@B)
// 1 : encrypted fsys (nB)
// 2 : encrypted fkey (53nB)
// 3 : encrypted (emptypos 8B + fphykey 4096B)

// block B (basic) KDB text
// -> (ClusterName@S chunksize@I wrsign@B blocknum@I)

// 1024B base webp
func basewebp() []byte {
	var temp []byte
	temp = append(temp, 82, 73, 70, 70, 130, 2, 0, 0, 87, 69, 66, 80, 86, 80, 56, 32, 118, 2, 0, 0, 240, 12, 0, 157, 1, 42, 64, 0, 64, 0, 62, 109, 44, 147, 69, 164, 34, 161, 151, 10)
	temp = append(temp, 78, 168, 64, 6, 196, 177, 128, 95, 52, 196, 136, 252, 252, 61, 72, 35, 129, 92, 52, 171, 25, 162, 249, 75, 250, 59, 216, 21, 37, 62, 185, 63, 101, 27, 23, 89, 10, 255, 7, 66)
	temp = append(temp, 205, 168, 178, 131, 219, 96, 84, 217, 174, 77, 236, 11, 30, 150, 203, 37, 71, 232, 222, 52, 73, 13, 233, 18, 67, 181, 234, 251, 59, 99, 213, 124, 138, 28, 95, 107, 253, 185, 26, 137)
	temp = append(temp, 8, 173, 230, 117, 231, 143, 46, 1, 217, 90, 42, 128, 0, 254, 251, 113, 215, 73, 116, 185, 5, 36, 95, 170, 218, 68, 56, 110, 42, 30, 65, 177, 185, 196, 35, 122, 198, 94, 175, 46)
	temp = append(temp, 78, 120, 169, 63, 124, 72, 59, 74, 251, 0, 80, 68, 96, 0, 203, 51, 54, 237, 67, 18, 78, 52, 25, 144, 169, 194, 207, 133, 217, 143, 149, 201, 199, 238, 129, 4, 5, 113, 173, 238)
	temp = append(temp, 237, 190, 13, 234, 136, 193, 111, 0, 174, 203, 236, 219, 29, 169, 108, 142, 92, 247, 253, 173, 199, 13, 181, 182, 69, 175, 174, 147, 52, 244, 175, 67, 194, 192, 223, 119, 208, 199, 216, 251)
	temp = append(temp, 35, 102, 251, 255, 142, 222, 243, 83, 235, 87, 179, 109, 159, 7, 62, 212, 167, 145, 249, 7, 52, 213, 30, 251, 136, 173, 228, 75, 150, 198, 16, 251, 155, 183, 27, 91, 50, 90, 122, 103)
	temp = append(temp, 30, 215, 244, 205, 97, 207, 70, 251, 125, 18, 204, 66, 151, 230, 251, 117, 66, 19, 78, 42, 58, 131, 230, 153, 75, 108, 228, 206, 211, 230, 233, 0, 13, 158, 146, 51, 184, 192, 74, 192)
	temp = append(temp, 70, 26, 156, 177, 238, 35, 132, 45, 168, 232, 204, 160, 6, 177, 59, 231, 234, 198, 4, 61, 127, 204, 83, 137, 206, 45, 220, 150, 75, 237, 35, 124, 63, 65, 233, 13, 4, 229, 242, 21)
	temp = append(temp, 243, 88, 237, 141, 26, 92, 108, 9, 240, 144, 2, 174, 211, 169, 215, 168, 93, 250, 33, 193, 196, 154, 103, 105, 127, 244, 167, 188, 171, 243, 104, 9, 67, 167, 197, 70, 130, 190, 238, 252)
	temp = append(temp, 182, 255, 248, 240, 98, 186, 111, 8, 73, 186, 179, 164, 149, 115, 148, 241, 183, 65, 39, 172, 46, 233, 54, 203, 22, 146, 82, 239, 234, 232, 228, 236, 53, 141, 82, 22, 99, 190, 242, 37)
	temp = append(temp, 7, 83, 54, 124, 227, 168, 52, 9, 135, 48, 218, 92, 63, 177, 114, 132, 132, 86, 98, 89, 65, 57, 110, 241, 83, 215, 40, 33, 250, 132, 114, 123, 181, 142, 250, 155, 195, 9, 19, 30)
	temp = append(temp, 202, 54, 186, 2, 39, 240, 62, 225, 86, 40, 51, 217, 46, 217, 133, 32, 251, 182, 217, 26, 168, 155, 176, 50, 61, 214, 20, 214, 190, 162, 87, 207, 40, 78, 54, 133, 73, 93, 207, 41)
	temp = append(temp, 198, 107, 54, 43, 68, 172, 102, 241, 158, 121, 94, 195, 111, 81, 157, 134, 66, 51, 137, 42, 124, 11, 105, 32, 163, 95, 28, 220, 210, 86, 229, 165, 178, 214, 199, 224, 8, 64, 111, 244)
	temp = append(temp, 96, 89, 189, 184, 141, 46, 114, 231, 221, 151, 196, 34, 233, 115, 237, 39, 63, 255, 30, 61, 18, 97, 152, 245, 189, 123, 184, 1, 249, 172, 213, 184, 109, 106, 221, 153, 104, 24, 124, 76)
	temp = append(temp, 50, 203, 174, 95, 111, 244, 206, 44, 87, 83, 237, 48, 201, 227, 100, 175, 26, 21, 178, 249, 178, 255, 241, 41, 159, 216, 176, 125, 237, 176, 86, 176, 215, 159, 138, 95, 5, 42, 176, 7)
	temp = append(temp, 167, 29, 73, 20, 177, 9, 206, 192, 0, 0)
	temp = append(temp, make([]byte, 1024-len(temp))...)
	return temp
}

// block A : salt hint pwhash account wrsign fsyskdt fkeykdt fphykdt
// block B : cluster chunksize wrsign blocknum

// setting info, caller has response of data manage
type sdata struct {
	desktop string // desktop path (~/)
	local   string // local path (~/)
	remote  string // remote path (~/)

	cluster   string  // cluster name
	chunksize int     // cluster chunksize
	wrsign    [8]byte // cluster writing sign (8B)
	blocknum  int     // cluster block num

	account string    // account name
	hint    []byte    // hint
	salt    [64]byte  // salt (64B)
	pwhash  [192]byte // pwhash (192B)

	mkey   [144]byte   // masterkey (144B)
	keybuf [6][48]byte // key (48B) of (fsys fkey fphy), [0:3] plain [3:6] enc
}

// random fill internal data, ( [desktop, local, remote], [cluster, account], chunksize, [pw, kf, hint] )
func (tbox *sdata) fill(paths []string, names []string, csize int, datas [][]byte) {
	tbox.desktop = paths[0]
	tbox.local = paths[1]
	tbox.remote = paths[2]

	tbox.cluster = names[0]
	tbox.chunksize = csize
	tbox.wrsign = [8]byte(kaes.Genrand(8))
	tbox.blocknum = -1

	var worker basework
	tbox.account = names[1]
	tbox.hint = datas[2]
	tbox.salt = [64]byte(kaes.Genrand(64))
	tbox.pwhash, tbox.mkey = worker.genpm(datas[0], datas[1], tbox.salt[:])

	tbox.keybuf[0] = [48]byte(kaes.Genrand(48))
	tbox.keybuf[1] = [48]byte(kaes.Genrand(48))
	tbox.keybuf[2] = [48]byte(kaes.Genrand(48))
	tbox.keybuf[3] = [48]byte{}
	tbox.keybuf[4] = [48]byte{}
	tbox.keybuf[5] = [48]byte{}
}

// read header block A/B text part, returns if wrsign is equal
func (tbox *sdata) rdhead(blockA string, blockB string) (bool, error) {
	worker0 := kdb.Initkdb()
	err := worker0.Read(blockB)
	if err != nil {
		return false, err
	}
	tmp, _ := worker0.Get("cluster")
	tbox.cluster = tmp.Dat6
	tmp, _ = worker0.Get("chunksize")
	tbox.chunksize = tmp.Dat2
	tmp, _ = worker0.Get("wrsign")
	tbox.wrsign = [8]byte(tmp.Dat5)
	tmp, _ = worker0.Get("blocknum")
	tbox.blocknum = tmp.Dat2

	worker1 := kdb.Initkdb()
	err = worker1.Read(blockA)
	if err != nil {
		return false, err
	}
	tmp, _ = worker1.Get("account")
	tbox.account = tmp.Dat6
	tmp, _ = worker1.Get("hint")
	tbox.hint = []byte(tmp.Dat6)
	tmp, _ = worker1.Get("salt")
	tbox.salt = [64]byte(tmp.Dat5)
	tmp, _ = worker1.Get("pwhash")
	tbox.pwhash = [192]byte(tmp.Dat5)
	tmp, _ = worker1.Get("wrsign")
	cmp := tmp.Dat5

	tmp, _ = worker1.Get("fsyskdt")
	tbox.keybuf[3] = [48]byte(tmp.Dat5)
	tmp, _ = worker1.Get("fkeykdt")
	tbox.keybuf[4] = [48]byte(tmp.Dat5)
	tmp, _ = worker1.Get("fphykdt")
	tbox.keybuf[5] = [48]byte(tmp.Dat5)
	return kio.Bequal(tbox.wrsign[:], cmp), nil
}

// write header block A/B text part, returns (A, B), !! should set keybuf[3:6] first !!
func (tbox *sdata) wrhead() (string, string) {
	worker0 := kdb.Initkdb()
	worker0.Read("cluster = 0\nchunksize = 0\nwrsign = 0\nblocknum = 0")
	worker0.Fix("cluster", tbox.cluster)
	worker0.Fix("chunksize", tbox.chunksize)
	worker0.Fix("wrsign", tbox.wrsign[:])
	worker0.Fix("blocknum", tbox.blocknum)
	blockB, _ := worker0.Write()

	worker1 := kdb.Initkdb()
	worker1.Read("account = 0\nhint = 0\nsalt = 0\npwhash = 0\nwrsign = 0\nfsyskdt = 0\nfkeykdt = 0\nfphykdt = 0")
	worker1.Fix("salt", tbox.salt[:])
	worker1.Fix("hint", string(tbox.hint))
	worker1.Fix("pwhash", tbox.pwhash[:])
	worker1.Fix("account", tbox.account)
	worker1.Fix("wrsign", tbox.wrsign[:])

	worker1.Fix("fsyskdt", tbox.keybuf[3][:])
	worker1.Fix("fkeykdt", tbox.keybuf[4][:])
	worker1.Fix("fphykdt", tbox.keybuf[5][:])
	blockA, _ := worker1.Write()
	return blockA, blockB
}

// checks if src contains sub
func (tbox *sdata) contains(src string, sub string) bool {
	if len(src) > len(sub) {
		return false
	}
	if sub[0:len(src)] == src {
		return true
	} else {
		return false
	}
}

// basic module, caller has response of data manage
type bmod struct {
	enc kaes.Funcmode // kaes funcmode
	log logger        // logger
	drv basework      // integrated io

	fsysmod vdir       // fsys module (supreme folder)
	fkeymod keymap     // fkey module
	fphymod fphyIO     // fphy module
	fphykey [4096]byte // fphy cluster key (4096B)
}

// random fill internal data ( [desktop, local, remote], chunksize )
func (tbox *bmod) fill(paths []string, csize int) {
	tbox.log.clear()
	tbox.log.abort = false
	tbox.log.working = false
	tbox.log.readonly = false
	tbox.drv.log = &tbox.log
	tbox.drv.sleep = 4

	tbox.fphykey = [4096]byte(kaes.Genrand(4096))
	tbox.fphymod.init(paths[1], paths[2], csize, &tbox.drv, tbox.fphykey, -1, -1)
	tbox.fphymod.clear()
	tbox.fkeymod.read(nil)

	var root vdir
	root.name = "/"
	root.time = int(time.Now().Unix())
	root.islocked = false

	var bin vdir
	bin.name = "_BIN/"
	bin.time = root.time
	bin.islocked = false

	bufkey := kaes.Genrand(48)
	tbox.drv.dirctrl(paths[1]+"io/", true)
	tbox.enc.Before.Open(make([]byte, 4096*csize), true)
	tbox.enc.After.Open(paths[1]+"io/temp.bin", false)
	tbox.enc.Encrypt(bufkey)
	tbox.enc.Before.Close()
	tbox.enc.After.Close()

	bufptr := tbox.fphymod.fpush(paths[1] + "io/temp.bin")
	if bufptr < 0 {
		tbox.log.write("critical : push fail -/_BUF")
		tbox.log.abort = true
	}
	tbox.fkeymod.push(bufptr, [48]byte(bufkey))

	var buf vfile
	buf.name = "_BUF"
	buf.time = root.time
	buf.size = kio.Size(paths[1] + "io/temp.bin")
	buf.fptr = bufptr

	root.subdir = []vdir{bin}
	root.subfile = []vfile{buf}
	tbox.fsysmod = root
}

// init modules ( [desktop, local, remote], [chunksize, blocknum, sleep], readonly, [fsys, fkey, fphy(8+4096)] )
func (tbox *bmod) init(paths []string, parms []int, ronly bool, datas [][]byte) error {
	tbox.log.abort = false
	tbox.log.working = false
	tbox.log.readonly = ronly
	tbox.drv.log = &tbox.log
	tbox.drv.sleep = parms[2]

	if len(datas[2]) != 4104 {
		return errors.New("invalid fphy")
	}
	tbox.fphykey = [4096]byte(datas[2][8:])
	tbox.fphymod.init(paths[1], paths[2], parms[0], &tbox.drv, tbox.fphykey, kobj.Decode(datas[2][0:8]), parms[1])
	tbox.fphymod.clear()
	err := tbox.fkeymod.read(datas[1])
	if err != nil {
		return err
	}
	tbox.fsysmod = *vread(datas[0])
	return nil
}

// search with name, domain (~/) + tgt (~/) contains pattern (*?)
func (tbox *bmod) search(domain string, tgt *vdir, pattern string, res chan []string) {
	defer func() {
		if err := recover(); err != nil {
			res <- nil
		}
	}()

	out := make([]string, 0)
	re := regexp.MustCompile(pattern)
	if re.FindString(tgt.name) != "" {
		out = append(out, domain+tgt.name)
	}
	for _, r := range tgt.subfile {
		if re.FindString(r.name) != "" {
			out = append(out, domain+tgt.name+r.name)
		}
	}

	ret := make([]chan []string, len(tgt.subdir))
	for i := range ret {
		ret[i] = make(chan []string, 1)
	}
	for i, r := range tgt.subdir {
		go tbox.search(domain+tgt.name, &r, pattern, ret[i])
	}
	for i := range ret {
		out = append(out, <-ret[i]...)
	}
	res <- out
}

// reset dangerous name ( / -> _S_, \ -> _BS_, : -> _C_, * -> _ST_, ? -> _Q_, " -> _QU_, <>| -> _P0_ _P1_ _P2_, \n -> _N_ )
func (tbox *bmod) reset(tgt *vdir, res chan int) {
	defer func() {
		if err := recover(); err != nil {
			res <- 0
		}
	}()
	cont := func(input string) bool {
		for _, r := range input {
			switch r {
			case '/':
				return true
			case '\\':
				return true
			case ':':
				return true
			case '\n':
				return true
			case '*':
				return true
			case '?':
				return true
			case '|':
				return true
			case '"':
				return true
			case '<':
				return true
			case '>':
				return true
			}
		}
		return false
	}
	conv := func(input string) string {
		input = strings.ReplaceAll(strings.ReplaceAll(input, "/", "_S_"), "\\", "_BS_")
		input = strings.ReplaceAll(strings.ReplaceAll(input, ":", "_C_"), "\n", "_N_")
		input = strings.ReplaceAll(strings.ReplaceAll(input, "*", "_ST_"), "?", "_Q_")
		input = strings.ReplaceAll(strings.ReplaceAll(input, "|", "_P2_"), "\"", "_QU_")
		input = strings.ReplaceAll(strings.ReplaceAll(input, "<", "_P0_"), ">", "_P1_")
		return input
	}

	count := 0
	if cont(tgt.name[:len(tgt.name)-1]) {
		count = count + 1
		tgt.name = conv(tgt.name[:len(tgt.name)-1]) + "/"
	}
	for i, r := range tgt.subfile {
		if cont(r.name) {
			count = count + 1
			tgt.subfile[i].name = conv(r.name)
		}
	}

	ret := make([]chan int, len(tgt.subdir))
	for i := range ret {
		ret[i] = make(chan int, 1)
	}
	for i := range tgt.subdir {
		go tbox.reset(&tgt.subdir[i], ret[i])
	}
	for i := range ret {
		count = count + <-ret[i]
	}

	if count != 0 {
		tgt.sort()
	}
	res <- count
}

// print folder to string with indent & \n, return string not start/end with \n
func (tbox *bmod) print(folder *vdir, wrlock bool, indent int, res chan string) {
	defer func() {
		if err := recover(); err != nil {
			res <- ""
		}
	}()
	if folder.islocked && !wrlock { // should not include
		res <- ""
	} else {
		// current folder
		out := make([]string, 0)
		out = append(out, fmt.Sprintf("%s%s  (Lock %t  Time %d)", strings.Repeat("    ", indent), folder.name, folder.islocked, folder.time))

		// direct subfiles
		for _, r := range folder.subfile {
			out = append(out, fmt.Sprintf("%s%s  (Size %d  Time %d  Fptr %d)", strings.Repeat("    ", indent+1), r.name, r.size, r.time, r.fptr))
		}

		// direct subfolders
		ret := make([]chan string, len(folder.subdir))
		for i := 0; i < len(folder.subdir); i++ {
			ret[i] = make(chan string, 1)
		}
		for i, r := range folder.subdir {
			go tbox.print(&r, wrlock, indent+1, ret[i])
		}
		for i := 0; i < len(folder.subdir); i++ {
			temp := <-ret[i]
			if temp != "" {
				out = append(out, temp)
			}
		}

		res <- strings.Join(out, "\n")
	}
}

// activate kaes, key (48B)
func (tbox *bmod) ende(key [48]byte, isencrypt bool) error {
	defer func() {
		if err := recover(); err != nil {
			tbox.log.write(fmt.Sprintf("critical : %s", err))
			tbox.log.abort = true
		}
	}()
	if tbox.log.abort {
		return errors.New("kaes aborted")
	} else {
		if isencrypt {
			return tbox.enc.Encrypt(key[:])
		} else {
			return tbox.enc.Decrypt(key[:])
		}
	}
}

// get folder by path, returns nil if not exists, names : name frag cut by /
func (tbox *bmod) getdir(domain *vdir, names []string) *vdir {
	if tbox.log.abort {
		return nil
	}
	if len(names) == 0 {
		return domain
	}
	for i, r := range domain.subdir {
		if r.name == names[0]+"/" {
			return tbox.getdir(&domain.subdir[i], names[1:])
		}
	}
	return nil
}

// get fkey data (53nB) under folder
func (tbox *bmod) getkeys(folder *vdir, wrlocked bool, res chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			res <- nil
		}
	}()
	if folder.islocked && !wrlocked { // should not include
		res <- nil
	} else {

		temp := make([]byte, 0) // direct subfiles
		for _, r := range folder.subfile {
			ext, _, data := tbox.fkeymod.seek(r.fptr)
			if ext != -1 {
				temp = append(append(temp, kobj.Encode(r.fptr, 5)...), data[:]...)
			}
		}

		ret := make([]chan []byte, len(folder.subdir)) // subfolders
		for i := 0; i < len(folder.subdir); i++ {
			ret[i] = make(chan []byte, 1)
		}
		for i, r := range folder.subdir {
			go tbox.getkeys(&r, wrlocked, ret[i])
		}
		for i := 0; i < len(folder.subdir); i++ {
			temp = append(temp, <-ret[i]...)
		}
		res <- temp
	}
}

// fpop : file fptr -> path
func (tbox *bmod) fpop(fptr int, fsize int, path string) {
	// it will be rewritten if file exists, require : local/io/
	ext, _, key := tbox.fkeymod.seek(fptr)
	if ext == -1 { // fkey seek fail
		tbox.log.abort = true
		tbox.log.write(fmt.Sprintf("critical : no key -%d", fptr))
	} else {
		if !tbox.fphymod.fchk(fsize, fptr, false) { // fphy broken chain
			tbox.log.abort = true
			tbox.log.write(fmt.Sprintf("critical : broken chain -%d", fptr))
		} else {
			os.Remove(tbox.fphymod.local + "io/temp.bin") // fphy fpop to local/io/temp.bin
			tbox.fphymod.fpop(tbox.fphymod.local + "io/temp.bin")
			tbox.fphymod.clear()

			tbox.enc.Before.Open(tbox.fphymod.local+"io/temp.bin", true) // kaes decrypt
			defer tbox.enc.Before.Close()
			tbox.enc.After.Open(path, false)
			defer tbox.enc.After.Close()
			if err := tbox.ende(key, false); err != nil { // kaes decrypt fail
				tbox.log.abort = true
				tbox.log.write(fmt.Sprintf("critical : decrypt fail -%s", err))
			} else {
				tbox.log.write(fmt.Sprintf("msg : fpop -%s", path))
			}
		}
	}
}

// fdel : file delete, broken chain will not be deleted
func (tbox *bmod) fdel(fptr int, fsize int) {
	if tbox.fphymod.fchk(fsize, fptr, false) { // normal chain -> fphy fdel
		tbox.fphymod.fdel()
		tbox.fphymod.clear()
		tbox.fkeymod.pop(fptr)
		tbox.log.write(fmt.Sprintf("msg : fdel -fptr %d size %d", fptr, fsize))
	} else { // fphy broken chain
		tbox.log.abort = true
		tbox.fphymod.clear()
		tbox.fkeymod.pop(fptr)
		tbox.log.write(fmt.Sprintf("critical : broken chain -%d", fptr))
	}
}

// reset islocked status on self/subdir
func (tbox *bmod) relock(folder *vdir, islocked bool) {
	folder.islocked = islocked
	for i := range folder.subdir {
		tbox.relock(&folder.subdir[i], islocked)
	}
}

// ===== section ===== pevfs meta functions

// generate new vault at remote path (~/), path should be empty folder
func PEVFS_New(remote string, cluster string, chunksize int) error {
	if remote[len(remote)-1] != '/' {
		return errors.New("not dir path")
	}
	subs, err := os.ReadDir(remote)
	if err != nil {
		return err
	}
	if len(subs) != 0 {
		return errors.New("not empty dir")
	}

	var worker PEVFS // new PEVFS
	pw := []byte("0000")
	kf := kaes.Basickey()
	hint := []byte("new : PW 0000 KF bkf")
	os.Mkdir("./temp693a/", os.ModePerm)
	defer os.RemoveAll("./temp693a/")
	os.Mkdir("./temp693b/", os.ModePerm)
	defer os.RemoveAll("./temp693b/")

	// fill data
	worker.data.fill([]string{"./temp693a/", "./temp693b/", remote}, []string{cluster, "root"}, chunksize, [][]byte{pw, kf, hint})
	worker.module.fill([]string{"./temp693a/", "./temp693b/", remote}, chunksize)
	for i := 0; i < 3; i++ {
		temp := worker.data.mkey[48*i : 48*i+48]
		worker.data.keybuf[i+3] = [48]byte(worker.module.drv.aescalc(worker.data.keybuf[i][:], [48]byte(temp), true, false))
	}
	worker.Rootpath = worker.module.fsysmod.name
	worker.Cluster = &worker.data.cluster
	worker.Account = &worker.data.account
	worker.flush_all()
	return nil
}

// boot with cluster, init internal data & local/, returns (PEVFS, hint)
func PEVFS_Boot(desktop string, local string, remote string, blockApath string) (*PEVFS, []byte, error) {
	// prebooting minimal data
	var out PEVFS
	var err error
	worker := ksc.Initksc()
	worker.Predetect = true
	out.module.log.readonly = true
	out.module.drv.log = &out.module.log
	out.module.fphymod.rawIO = &out.module.drv
	out.module.drv.dirctrl(local, true)

	// KSC unpacking block A
	worker.Path = blockApath
	err = worker.Readf()
	if err != nil {
		return nil, nil, err
	}
	if !kio.Bequal(worker.Subtype, []byte("KV5a")) {
		return nil, nil, errors.New("invalid block A")
	}
	blockA := out.module.drv.read(blockApath, worker.Chunkpos[0]+8, worker.Chunksize[0])
	out.databuf[3] = out.module.drv.read(blockApath, worker.Chunkpos[1]+8, worker.Chunksize[1])
	out.databuf[4] = out.module.drv.read(blockApath, worker.Chunkpos[2]+8, worker.Chunksize[2])
	out.databuf[5] = out.module.drv.read(blockApath, worker.Chunkpos[3]+8, worker.Chunksize[3])
	blockB := out.module.fphymod.hpop(remote + "0b.txt")
	if !kio.Bequal(worker.Reserved[0:4], ksc.Crc32hash(blockA)) {
		return nil, nil, errors.New("broken block A")
	}
	if blockB == nil {
		return nil, nil, errors.New("invalid block B")
	}

	// setting data section
	out.data.desktop = desktop
	out.data.local = local
	out.data.remote = remote
	out.hpath = blockApath
	wrsign, err := out.data.rdhead(string(blockA), string(blockB))
	if err != nil {
		return nil, nil, err
	}
	if !kio.Bequal(worker.Reserved[4:8], ksc.Crc32hash([]byte(out.data.cluster))) {
		out.module.log.write("err : cluster name not match")
	}
	if !wrsign {
		out.module.log.write("err : cluster session not match")
	}

	// setting readvar section
	out.Cluster = &out.data.cluster
	out.Account = &out.data.account
	if *out.Account == "root" && remote+"0a.webp" == blockApath {
		out.module.log.readonly = false // can write only with valid account name & path
	} else {
		out.module.log.readonly = true
	}
	out.module.log.write("msg : cluster boot success")
	return &out, out.data.hint, nil
}

// delete local folder, erase internal data
func PEVFS_Exit(obj *PEVFS) {
	// readvar section clear
	defer os.RemoveAll(obj.data.local)
	obj.module.log.abort = true
	obj.module.log.readonly = true
	obj.module.log.working = false
	obj.cachedir = nil
	obj.cachepath = ""
	obj.Curdir = nil
	obj.Curpath = ""

	// data section clear
	obj.data.salt = [64]byte{}
	obj.data.hint = nil
	obj.data.pwhash = [192]byte{}
	obj.data.mkey = [144]byte{}
	for i := 0; i < 6; i++ {
		obj.data.keybuf[i] = [48]byte{}
		obj.databuf[i] = nil
	}

	// module section clear
	obj.module.log.clear()
	obj.module.fsysmod.name = ""
	obj.module.fsysmod.time = -1
	obj.module.fsysmod.islocked = false
	obj.module.fsysmod.subdir = nil
	obj.module.fsysmod.subfile = nil
	obj.module.fkeymod.read(make([]byte, 53))
	obj.module.fphymod.clear()
	obj.module.fphymod.tgtsize = -1
	obj.module.fphymod.comsize = -1
	obj.module.fphymod.emptyfptr = -1
	obj.module.fphymod.blocknum = -1
	obj.module.fphymod.ctrlIO.fphykey = [4096]byte{}
	obj.module.fphymod.ctrlIO = nil
	obj.module.fphykey = [4096]byte{}
}

// regenerate fphykey, rewrite control blocks
func PEVFS_Rebuild(remote string, pw []byte, kf []byte) error {
	// block C reader/writer
	var logdrv logger
	var iodrv basework
	iodrv.log = &logdrv
	iodrv.sleep = 4
	var reader blockctrl
	var writer blockctrl
	newkey := [4096]byte(kaes.Genrand(4096))

	// login to PEVFS
	worker, hint, err := PEVFS_Boot("./", "./temp693c/", remote, remote+"0a.webp")
	if err != nil {
		return err
	}
	err = worker.Login(pw, kf, 4)
	if err != nil {
		return err
	}
	defer PEVFS_Exit(worker)
	if worker.module.log.readonly {
		return errors.New("readonly account")
	}
	reader = *worker.module.fphymod.ctrlIO
	writer.init(&iodrv, newkey, "./temp693c/", remote)

	// rewrite block C & save PEVFS
	for i := 0; i < worker.data.blocknum; i++ {
		blockC := reader.ende_inline(reader.read_l2l3(i), i, false)
		writer.write_l2l3(i, writer.ende_inline(blockC, i, true))
	}
	worker.module.fphykey = newkey
	worker.module.fphymod.ctrlIO.fphykey = newkey
	err = worker.AccReset(pw, kf, hint)
	return err
}

// ===== section ===== pevfs general functions

// personal encrypted virtual file system
type PEVFS struct {
	data   sdata // setting data
	module bmod  // worker module

	hpath   string     // header block A path
	databuf [6][]byte  // data (nB) of (fsys fkey fphy), [0:3] plain [3:6] enc
	lock    sync.Mutex // fsys access lock (one at once)

	cachepath string // empty str if cachedir is nil
	cachedir  *vdir  // cache of accessed folder
	tgtpath   string // empty str if tgtdir is nil
	tgtdir    *vdir  // target of current working

	Curpath string // empty str if Curdir is nil
	Curdir  *vdir  // current folder

	Rootpath string  // path of root folder (*/) -readonly const
	Cluster  *string // cluster name -readonly const
	Account  *string // account name (root : RW, else : R) -readonly const
}

// generate block A/B header by current data/buffer
func (tbox *PEVFS) flush_head() ([]byte, []byte) {
	// require : data, data.keybuf[0:6], databuf[0:3]
	var err error
	for i := 0; i < 3; i++ {
		err = tbox.module.enc.Before.Open(tbox.databuf[i], true)
		if err != nil {
			return nil, nil
		}
		err = tbox.module.enc.After.Open(make([]byte, 0, 10485760), false)
		if err != nil {
			return nil, nil
		}
		err = tbox.module.ende(tbox.data.keybuf[i], true)
		if err != nil {
			return nil, nil
		}
		tbox.module.enc.Before.Close()
		tbox.databuf[i+3] = tbox.module.enc.After.Close()
	}
	tbox.data.blocknum = tbox.module.fphymod.blocknum
	tbox.module.log.write("msg : encryption done -flush_head")

	blockA, blockB := tbox.data.wrhead() // write block A/B
	worker := ksc.Initksc()
	worker.Prehead = basewebp()
	worker.Subtype = []byte("KV5a")
	worker.Reserved = append(ksc.Crc32hash([]byte(blockA)), ksc.Crc32hash([]byte(*tbox.Cluster))...)
	temp, err := worker.Writeb()
	if err != nil {
		return nil, nil
	}
	temp = worker.Linkb(temp, []byte(blockA))
	temp = worker.Linkb(temp, tbox.databuf[3])
	temp = worker.Linkb(temp, tbox.databuf[4])
	temp = worker.Linkb(temp, tbox.databuf[5])
	temp, _ = worker.Addb(temp, "")
	tbox.module.log.write("msg : packaging done -flush_head")
	return temp, []byte(blockB)
}

// update all header(block ABC) & clear databuf, only for RW mode, ! FsysLock !
func (tbox *PEVFS) flush_all() {
	// require : data, data.keybuf[0:6], module
	if !tbox.module.log.abort && !tbox.module.log.readonly {
		tbox.lock.Lock()
		tbox.databuf[0] = tbox.module.fsysmod.write(true) // fsys access point
		tbox.lock.Unlock()
		tbox.databuf[1] = tbox.module.fkeymod.write()
		tbox.databuf[2] = append(kobj.Encode(tbox.module.fphymod.emptyfptr, 8), tbox.module.fphykey[:]...)
		tbox.module.log.write("msg : data buffer filled -flush_all")

		blockA, blockB := tbox.flush_head()
		tbox.module.fphymod.hpush(tbox.data.remote+"0a.webp", blockA, true)
		tbox.module.fphymod.hpush(tbox.data.remote+"0b.txt", blockB, true)
		tbox.module.log.write("msg : header block A/B hpush done -flush_all")

		tbox.module.fphymod.flush(true)
		tbox.module.log.write("msg : control block C flush done -flush_all")
		for i := 0; i < 6; i++ {
			tbox.databuf[i] = nil
		}
	}
}

// check folder contains filename, returns inited subfile index, ! FsysLock !
func (tbox *PEVFS) fpush_check(domain *vdir, name string) int {
	// name exist : fptr deleted, not exist : new generated
	tbox.lock.Lock()
	defer tbox.lock.Unlock()
	pos := -1
	for i, r := range domain.subfile { // fsys access point
		if r.name == name {
			pos = i
			break
		}
	}
	if pos == -1 {
		pos = len(domain.subfile)
		var temp vfile
		temp.name = name
		temp.time = int(time.Now().Unix())
		domain.subfile = append(domain.subfile, temp)
	} else {
		tbox.module.fdel(domain.subfile[pos].fptr, domain.subfile[pos].size)
		domain.subfile[pos].time = int(time.Now().Unix())
		tbox.module.log.write(fmt.Sprintf("msg : fpush file exist -%s", name))
	}
	return pos
}

// import file(full path) to domain/name, call folder.sort after this, ! FsysLock !
func (tbox *PEVFS) fpush_file(path string, domain *vdir, name string) {
	// it will be rewritten if file exists, require : local/io/
	pos := tbox.fpush_check(domain, name)

	key := [48]byte(kaes.Genrand(48)) // kaes encrypt
	os.Remove(tbox.data.local + "io/temp.bin")
	tbox.module.enc.Before.Open(path, true)
	tbox.module.enc.After.Open(tbox.data.local+"io/temp.bin", false)
	err := tbox.module.ende(key, true)
	tbox.module.enc.Before.Close()
	tbox.module.enc.After.Close()

	if err != nil { // kaes fail
		tbox.module.log.abort = true
		tbox.module.log.write(fmt.Sprintf("critical : encrypt fail -%s", err))
	} else {
		fptr := tbox.module.fphymod.fpush(tbox.data.local + "io/temp.bin")
		tbox.module.fphymod.clear()
		if fptr == -1 { // fphy fpush fail
			tbox.module.log.abort = true
			tbox.module.log.write(fmt.Sprintf("critical : fpush fail -%s", path))
		} else {

			tbox.lock.Lock()
			domain.subfile[pos].size = kio.Size(tbox.data.local + "io/temp.bin")
			domain.subfile[pos].fptr = fptr
			tbox.lock.Unlock()
			err = tbox.module.fkeymod.push(fptr, key)

			if err != nil { // fkey push fail
				tbox.module.log.abort = true
				tbox.module.log.write(fmt.Sprintf("critical : keypush fail -%d", fptr))
			} else {
				tbox.module.log.write(fmt.Sprintf("msg : fpush -%s", path))
			}
		}
	}
}

// manage /_BUF & make local/io/ & reset wrsign, works only at RW mode, ! FsysLock !
func (tbox *PEVFS) fpush_buf() {
	if !tbox.module.log.abort && !tbox.module.log.readonly {
		tbox.data.wrsign = [8]byte(kaes.Genrand(8))
		tbox.lock.Lock()
		pos := -1
		var tgt vfile
		for i, r := range tbox.module.fsysmod.subfile { // fsys access point
			if r.name == "_BUF" {
				pos = i
				tgt = r
				break
			}
		}
		if pos != -1 {
			tbox.module.fsysmod.subfile = append(tbox.module.fsysmod.subfile[:pos], tbox.module.fsysmod.subfile[pos+1:]...)
		}
		tbox.lock.Unlock()

		tbox.module.drv.dirctrl(tbox.data.local+"io/", true)
		if pos == -1 {
			os.Remove(tbox.data.local + "io/buffer.bin")
			tbox.module.drv.write(tbox.data.local+"io/buffer.bin", 0, make([]byte, 4096*tbox.data.chunksize))
			tbox.fpush_file(tbox.data.local+"io/buffer.bin", &tbox.module.fsysmod, "_BUF")
			tbox.module.log.write("msg : /_BUF generated")
		} else {
			tbox.module.fdel(tgt.fptr, tgt.size)
			tbox.module.log.write("msg : /_BUF deleted")
		}
	}
}

// check folder contains name, returns subdir/subfile index/-1, ! FsysLock !
func (tbox *PEVFS) fpop_check(domain *vdir, name string) int {
	tbox.lock.Lock()
	defer tbox.lock.Unlock()
	if name[len(name)-1] == '/' {
		for i, r := range domain.subdir { // fsys access point
			if r.name == name {
				return i
			}
		}
	} else {
		for i, r := range domain.subfile { // fsys access point
			if r.name == name {
				return i
			}
		}
	}
	return -1
}

// export file(domain/name) to (folder/name), folder : (~/), ! FsysLock !
func (tbox *PEVFS) fpop_file(folder string, domain *vdir, name string) {
	pos := tbox.fpop_check(domain, name)
	if pos == -1 { // no such file
		tbox.module.log.abort = true
		tbox.module.log.write(fmt.Sprintf("critical : no file -%s", name))
	} else {
		tbox.lock.Lock()
		fptr := domain.subfile[pos].fptr
		size := domain.subfile[pos].size
		tbox.lock.Unlock()
		tbox.module.fpop(fptr, size, folder+name)
	}
}

// find folder by path (~/), returns nil if not exists, ! fsys lock !
func (tbox *PEVFS) dir_navigate(path string) *vdir {
	tbox.lock.Lock()
	defer tbox.lock.Unlock()
	if len(path) == 0 {
		return nil
	}
	if path[len(path)-1] != '/' {
		return nil
	}
	if path == tbox.Rootpath {
		return &tbox.module.fsysmod
	}
	if path == tbox.cachepath {
		return tbox.cachedir
	}

	var res *vdir
	if tbox.data.contains(tbox.cachepath, path) { // cache hit
		res = tbox.module.getdir(tbox.cachedir, strings.Split(path[len(tbox.cachepath):len(path)-1], "/")) // fsys access point
	} else if tbox.data.contains(tbox.Rootpath, path) { // cache miss
		res = tbox.module.getdir(&tbox.module.fsysmod, strings.Split(path[len(tbox.Rootpath):len(path)-1], "/")) // fsys access point
	} else { // invalid path
		res = nil
	}

	if res != nil { // update cache
		tbox.module.log.write(fmt.Sprintf("msg : new cache -%s", path))
		tbox.cachepath = path
		tbox.cachedir = res
	}
	return res
}

// generate vdir by real path(~/), using fkey/fphy but not fsys (non-blocking), ! FsysLock !
func (tbox *PEVFS) dir_import(path string) (*vdir, error) {
	var out vdir // gen new vdir
	out.name = path[strings.LastIndex(path[0:len(path)-1], "/")+1:]
	out.time = int(time.Now().Unix())
	out.islocked = false
	tbox.module.log.write(fmt.Sprintf("msg : vdir gen -%s", path))

	for _, r := range tbox.module.drv.dirsub(path) {
		if r[len(r)-1] == '/' { // push dir
			temp, err := tbox.dir_import(path + r)
			if err == nil {
				out.subdir = append(out.subdir, *temp)
				tbox.module.log.write(fmt.Sprintf("msg : vdir add -%s", r))
			} else {
				return nil, err
			}
		} else { // push file
			tbox.fpush_file(path+r, &out, r)
		}
		if tbox.module.log.abort {
			return nil, errors.New("abort")
		}
	}

	out.sort() // sort & return
	return &out, nil
}

// generate folder/domain/, (folder : ~/), ! FsysLock !
func (tbox *PEVFS) dir_export(folder string, domain *vdir) {
	tbox.lock.Lock()
	vdir_name := domain.name
	if vdir_name == "/" {
		vdir_name = "_/"
	}
	vfile_name := make([]string, len(domain.subfile))
	for i, r := range domain.subfile {
		vfile_name[i] = r.name
	}
	tbox.lock.Unlock()
	tbox.module.log.write(fmt.Sprintf("msg : vdir export -%s", domain.name))

	tbox.module.drv.dirctrl(folder+vdir_name, true) // create folder
	for _, r := range vfile_name {                  // create direct subfile
		if tbox.module.log.abort {
			break
		}
		tbox.fpop_file(folder+vdir_name, domain, r)
	}
	for _, r := range domain.subdir { // create direct subfolder
		if tbox.module.log.abort {
			break
		}
		tbox.dir_export(folder+vdir_name, &r)
	}
}

// get fptr/fsize under folder, ! FsysLock !
func (tbox *PEVFS) dir_getf(folder *vdir) ([]int, []int) {
	fptrs := make([]int, 0)
	fsizes := make([]int, 0)
	tbox.lock.Lock()
	for _, r := range folder.subfile {
		fptrs = append(fptrs, r.fptr)
		fsizes = append(fsizes, r.size)
	}
	tbox.lock.Unlock()
	for _, r := range folder.subdir {
		out0, out1 := tbox.dir_getf(&r)
		fptrs = append(fptrs, out0...)
		fsizes = append(fsizes, out1...)
	}
	return fptrs, fsizes
}

// check files under folder is valid (fkey/fphy), ! FsysLock !
func (tbox *PEVFS) sub_check(folder *vdir) int {
	tbox.lock.Lock() // get subfile info
	fptr := make([]int, len(folder.subfile))
	size := make([]int, len(folder.subfile))
	if strings.Contains(folder.name, "\n") || strings.Contains(folder.name[0:len(folder.name)-1], "/") {
		tbox.module.log.write(fmt.Sprintf("err : invalid name -%s", folder.name))
	}
	for i, r := range folder.subfile {
		fptr[i] = r.fptr
		size[i] = r.size
		if r.name == "" || strings.Contains(r.name, "\n") || strings.Contains(r.name, "/") {
			tbox.module.log.write(fmt.Sprintf("err : invalid name -%s", r.name))
		}
	}
	tbox.lock.Unlock()

	count := 0
	for i, r := range fptr { // count subfile
		ext, _, _ := tbox.module.fkeymod.seek(r)
		if ext == -1 { // fkey seek fail
			tbox.module.log.write(fmt.Sprintf("err : no fkey -fptr %d", r))
			count = count + 1
			continue
		}
		if !tbox.module.fphymod.fchk(size[i], r, false) { // fphy check fail
			tbox.module.log.write(fmt.Sprintf("err : broken chain -fptr %d size %d", r, size[i]))
			count = count + 1
		}
	}

	for _, r := range folder.subdir { // count subdir
		count = count + tbox.sub_check(&r)
	}
	return count
}

// check file & push fkey if valid, ! FsysLock !
func (tbox *PEVFS) sub_rebuild(folder *vdir, fk *keymap) int {
	// !! fgc write option is true !!
	count := 0
	vfile_temp := make([]vfile, 0)

	tbox.lock.Lock()
	for _, r := range folder.subfile {
		ext, _, key := tbox.module.fkeymod.seek(r.fptr)
		if ext == -1 { // fkey seek fail
			tbox.module.log.write(fmt.Sprintf("err : no fkey -fptr %d", r.fptr))
			count = count + 1
			continue
		}
		if !tbox.module.fphymod.fchk(r.size, r.fptr, true) { // fphy check fail
			tbox.module.log.write(fmt.Sprintf("err : broken chain -fptr %d size %d", r.fptr, r.size))
			count = count + 1
		} else { // push to new storage
			vfile_temp = append(vfile_temp, r)
			fk.push(r.fptr, key)
		}
	}
	tbox.lock.Unlock()

	folder.subfile = vfile_temp
	for i := range folder.subdir { // count subdir
		count = count + tbox.sub_rebuild(&folder.subdir[i], fk)
	}
	return count
}

// reset cachedir, ! FsysLock !
func (tbox *PEVFS) sub_recache() {
	tbox.lock.Lock()
	tbox.cachedir = &tbox.module.fsysmod
	tbox.cachepath = tbox.cachedir.name
	tbox.lock.Unlock()
}

// check abort/working flag, check readonly if check_ro is true
func (tbox *PEVFS) sub_auth(check_ro bool) error {
	if tbox.module.log.abort {
		return errors.New("aborted cluster")
	} else if tbox.module.log.working {
		return errors.New("working cluster")
	} else if check_ro && tbox.module.log.readonly {
		return errors.New("readonly cluster")
	} else {
		return nil
	}
}

// reset abort/working flag & cache/cur dir, returns abort/working flag
func (tbox *PEVFS) Abort(reset bool, abort bool, working bool) (bool, bool) {
	if reset {
		tbox.module.log.write(fmt.Sprintf("msg : reset -abort %t working %t", abort, working))
		tbox.module.log.abort = abort
		tbox.module.log.working = working

		tbox.cachedir = &tbox.module.fsysmod
		tbox.cachepath = tbox.cachedir.name
		tbox.Curdir = &tbox.module.fsysmod
		tbox.Curpath = tbox.Curdir.name
		tbox.Rootpath = tbox.module.fsysmod.name
		tbox.Cluster = &tbox.data.cluster
		tbox.Account = &tbox.data.account
	}
	return tbox.module.log.abort, tbox.module.log.working
}

// debug info return : [chunksize, blocknum, first fptr], [wrsign(8B), salt(64B), pwhash(192B), fsyskey(48B), fkeykey(48B), fphykey(48B)]
func (tbox *PEVFS) Debug() ([]int, [][]byte) {
	tbox.data.blocknum = tbox.module.fphymod.blocknum
	out0 := make([]int, 3)
	out1 := make([][]byte, 6)
	out0[0] = tbox.data.chunksize
	out0[1] = tbox.data.blocknum
	out0[2] = tbox.module.fphymod.emptyfptr
	out1[0] = append(make([]byte, 0), tbox.data.wrsign[:]...)
	out1[1] = append(make([]byte, 0), tbox.data.salt[:]...)
	out1[2] = append(make([]byte, 0), tbox.data.pwhash[:]...)
	out1[3] = append(make([]byte, 0), tbox.data.keybuf[0][:]...)
	out1[4] = append(make([]byte, 0), tbox.data.keybuf[1][:]...)
	out1[5] = append(make([]byte, 0), tbox.data.keybuf[2][:]...)
	return out0, out1
}

// returns log data joined with \n, returns empty string if reset
func (tbox *PEVFS) Log(reset bool) string {
	if reset {
		tbox.module.log.clear()
		return ""
	} else {
		return tbox.module.log.read()
	}
}

// login & set module, header will be rewritten if root
func (tbox *PEVFS) Login(pw []byte, kf []byte, sleeptime int) error {
	// pwhash & mkey will remain in data, data.keybuf[0:3] will remain, other buffers will be deleted
	cmp0, cmp1 := tbox.module.drv.genpm(pw, kf, tbox.data.salt[:])
	tbox.data.mkey = cmp1
	if tbox.data.pwhash != cmp0 {
		return errors.New("invalid PWKF")
	}

	var err error // decrypt key -> decrypt data, fill keybuf[0:3] databuf[0:3]
	for i := 0; i < 3; i++ {
		temp := tbox.data.mkey[48*i : 48*i+48]
		tbox.data.keybuf[i] = [48]byte(tbox.module.drv.aescalc(tbox.data.keybuf[i+3][:], [48]byte(temp), false, false))
		err = tbox.module.enc.Before.Open(tbox.databuf[i+3], true)
		if err != nil {
			return err
		}
		err = tbox.module.enc.After.Open(make([]byte, 0, 10485760), false)
		if err != nil {
			return err
		}
		err = tbox.module.ende(tbox.data.keybuf[i], false)
		if err != nil {
			return err
		}
		tbox.module.enc.Before.Close()
		tbox.databuf[i] = tbox.module.enc.After.Close()
	}

	// module init, field link
	tss := []string{tbox.data.desktop, tbox.data.local, tbox.data.remote}
	tis := []int{tbox.data.chunksize, tbox.data.blocknum, sleeptime}
	err = tbox.module.init(tss, tis, tbox.module.log.readonly, tbox.databuf[0:3])
	if err != nil {
		return err
	}
	tbox.cachedir = &tbox.module.fsysmod
	tbox.cachepath = tbox.cachedir.name
	tbox.Curdir = &tbox.module.fsysmod
	tbox.Curpath = tbox.Curdir.name
	tbox.Rootpath = tbox.module.fsysmod.name

	if !tbox.module.log.readonly { // rewrite root
		tbox.module.log.write("msg : header rewrite -root")
		tbox.data.salt = [64]byte(kaes.Genrand(64))
		tbox.data.pwhash, tbox.data.mkey = tbox.module.drv.genpm(pw, kf, tbox.data.salt[:])
		for i := 0; i < 3; i++ {
			temp := tbox.data.mkey[48*i : 48*i+48]
			tbox.data.keybuf[i] = [48]byte(kaes.Genrand(48))
			tbox.data.keybuf[i+3] = [48]byte(tbox.module.drv.aescalc(tbox.data.keybuf[i][:], [48]byte(temp), true, false))
		}
		blockA, blockB := tbox.flush_head()
		tbox.module.fphymod.hpush(tbox.data.remote+"0a.webp", blockA, true)
		tbox.module.fphymod.hpush(tbox.data.remote+"0b.txt", blockB, true)
	}
	for i := 0; i < 6; i++ {
		tbox.databuf[i] = nil
	}
	tbox.data.mkey = [144]byte{}
	tbox.module.log.write("msg : cluster login success")
	return nil
}

// reset account PWKF, block A header will be rewritten
func (tbox *PEVFS) AccReset(pw []byte, kf []byte, hint []byte) error {
	if ferr := tbox.sub_auth(false); ferr != nil {
		return ferr
	} else {
		defer func() { tbox.module.log.working = false }()
		tbox.module.log.working = true
	}

	// fill data section (data.keybuf[0:6])
	tbox.data.hint = hint
	tbox.data.salt = [64]byte(kaes.Genrand(64))
	tbox.data.pwhash, tbox.data.mkey = tbox.module.drv.genpm(pw, kf, tbox.data.salt[:])
	for i := 0; i < 3; i++ {
		temp := tbox.data.mkey[48*i : 48*i+48]
		tbox.data.keybuf[i] = [48]byte(kaes.Genrand(48))
		tbox.data.keybuf[i+3] = [48]byte(tbox.module.drv.aescalc(tbox.data.keybuf[i][:], [48]byte(temp), true, false))
	}
	tbox.module.log.write("msg : encryption done -AccReset")

	// fill buffer (databuf[0:3])
	tbox.databuf[0] = tbox.module.fsysmod.write(true)
	tbox.databuf[1] = tbox.module.fkeymod.write()
	tbox.databuf[2] = append(kobj.Encode(tbox.module.fphymod.emptyfptr, 8), tbox.module.fphykey[:]...)
	blockA, _ := tbox.flush_head()
	tbox.module.fphymod.hpush(tbox.hpath, blockA, false)
	for i := 0; i < 6; i++ {
		tbox.databuf[i] = nil
	}
	tbox.data.mkey = [144]byte{}
	tbox.module.log.write("msg : reset complete -AccReset")
	return nil
}

// extend account (curdir becomes new rootdir), returns new block A header path at desktop
func (tbox *PEVFS) AccExtend(pw []byte, kf []byte, hint []byte, account string, wrlocked bool) (string, error) {
	if ferr := tbox.sub_auth(false); ferr != nil {
		return "", ferr
	} else if tbox.Curdir.islocked && !wrlocked {
		return "", errors.New("locked folder")
	} else if account == "" || account == "root" {
		return "", errors.New("invalid account name")
	} else {
		defer func() { tbox.module.log.working = false }()
		tbox.module.log.working = true
		// save & restore session
		tmp_account := *tbox.Account
		tmp_hint := append(make([]byte, 0), tbox.data.hint...)
		tmp_salt := append(make([]byte, 0), tbox.data.salt[:]...)
		tmp_pwhash := append(make([]byte, 0), tbox.data.pwhash[:]...)
		var tmp_keybuf [6][48]byte
		copy(tmp_keybuf[:], tbox.data.keybuf[:])
		defer func() {
			*tbox.Account = tmp_account
			tbox.data.hint = tmp_hint
			tbox.data.account = tmp_account
			tbox.data.salt = [64]byte(tmp_salt)
			tbox.data.pwhash = [192]byte(tmp_pwhash)
			tbox.data.keybuf = tmp_keybuf
		}()
	}

	// fill data section (data.keybuf[0:6])
	tbox.data.hint = hint
	tbox.data.account = account
	tbox.Account = &account
	tbox.data.salt = [64]byte(kaes.Genrand(64))
	tbox.data.pwhash, tbox.data.mkey = tbox.module.drv.genpm(pw, kf, tbox.data.salt[:])
	for i := 0; i < 3; i++ {
		temp := tbox.data.mkey[48*i : 48*i+48]
		tbox.data.keybuf[i] = [48]byte(kaes.Genrand(48))
		tbox.data.keybuf[i+3] = [48]byte(tbox.module.drv.aescalc(tbox.data.keybuf[i][:], [48]byte(temp), true, false))
	}
	tbox.module.log.write("msg : encryption done -AccExtend")

	ret := make(chan []byte, 1)
	go tbox.module.getkeys(tbox.Curdir, wrlocked, ret)
	newpath := fmt.Sprintf("%s%da.webp", tbox.data.desktop, tbox.module.fphymod.rand_inline(100, 1000))
	// fill buffer (databuf[0:3])
	tbox.databuf[0] = tbox.Curdir.write(wrlocked)
	tbox.databuf[1] = <-ret
	tbox.databuf[2] = append(kobj.Encode(tbox.module.fphymod.emptyfptr, 8), tbox.module.fphykey[:]...)
	blockA, _ := tbox.flush_head()
	tbox.module.fphymod.hpush(newpath, blockA, false)
	for i := 0; i < 6; i++ {
		tbox.databuf[i] = nil
	}
	tbox.data.mkey = [144]byte{}
	tbox.module.log.write(fmt.Sprintf("msg : extend complete -%s", tbox.Curpath))
	return newpath, nil
}

// search name under Curdir, * : 0+ str, ? : len1 str, %d : int, %s 1+ ascii, %c len1 ascii, %* %? %% : literal
func (tbox *PEVFS) Search(name string) []string {
	temp := ""
	escape := false
	for _, r := range name {
		if escape {
			switch r {
			case '*':
				temp = temp + "\\*"
			case '?':
				temp = temp + "\\?"
			case '%':
				temp = temp + "\\%"
			case 'd':
				temp = temp + "\\d+"
			case 's':
				temp = temp + "\\w+"
			case 'c':
				temp = temp + "\\w"
			default:
				temp = temp + "\\s"
			}
			escape = false
		} else {
			switch r {
			case '*':
				temp = temp + ".*"
			case '?':
				temp = temp + "."
			case '[':
				temp = temp + "\\["
			case ']':
				temp = temp + "\\]"
			case '%':
				escape = true
			default:
				temp = temp + "[" + string(r) + "]"
			}
		}
	}

	tbox.module.log.write(fmt.Sprintf("msg : search regex -%s", temp))
	path := ""
	if strings.Count(tbox.Curpath, "/") > 1 {
		path = tbox.Curpath[0 : strings.LastIndex(tbox.Curpath[0:len(tbox.Curpath)-1], "/")+1]
	}
	tbox.lock.Lock()
	defer tbox.lock.Unlock()
	ret := make(chan []string, 1)
	go tbox.module.search(path, tbox.Curdir, temp, ret) // fsys access point
	out := <-ret
	tbox.module.log.write(fmt.Sprintf("msg : search complete -%d", len(out)))
	return out
}

// print Curdir & subobjs into string
func (tbox *PEVFS) Print(wrlocked bool) string {
	tbox.lock.Lock()
	defer tbox.lock.Unlock()
	ret := make(chan string, 1)
	go tbox.module.print(tbox.Curdir, wrlocked, 0, ret)
	out := <-ret
	tbox.module.log.write(fmt.Sprintf("msg : print complete -%d", len(out)))
	return out
}

// set Curdir to input path (~/), Curdir will not change if path not exist, returns T if find success
func (tbox *PEVFS) Teleport(path string) bool {
	temp := tbox.dir_navigate(path)
	if temp == nil {
		return false
	} else {
		tbox.Curdir = temp
		tbox.Curpath = path
		tbox.module.log.write(fmt.Sprintf("msg : moved to -%s", path))
		return true
	}
}

// get name/islocked of subdir/subfile under Curdir
func (tbox *PEVFS) NavName() ([]string, []bool) {
	tbox.lock.Lock()
	defer tbox.lock.Unlock()
	out0 := make([]string, 0) // fsys access point
	out1 := make([]bool, 0)
	for _, r := range tbox.Curdir.subdir {
		out0 = append(out0, r.name)
		out1 = append(out1, r.islocked)
	}
	for _, r := range tbox.Curdir.subfile {
		out0 = append(out0, r.name)
		out1 = append(out1, false)
	}
	return out0, out1
}

// Navigate_Name with info (takes more time), index 0 obj is Curdir itself
func (tbox *PEVFS) NavInfo(wrlocked bool) ([]string, []int, []int, []int, []int) {
	// names, times, size, fptr/islocked(0T 1F), [lowerdir, lowerfile]
	tbox.lock.Lock()
	defer tbox.lock.Unlock()
	names := []string{tbox.Curdir.name} // fsys access point
	times := []int{tbox.Curdir.time}
	sizes := []int{0} // size will be added
	var fptrs []int
	if tbox.Curdir.islocked {
		fptrs = []int{0}
	} else {
		fptrs = []int{1}
	}
	var lower []int
	if tbox.Curdir.islocked && !wrlocked {
		lower = []int{0, 0}
		return names, times, sizes, fptrs, lower
	} else {
		lower = []int{1, len(tbox.Curdir.subfile)} // data will be added
	}

	for _, r := range tbox.Curdir.subdir { // fsys access point
		names = append(names, r.name)
		times = append(times, r.time)
		if r.islocked {
			fptrs = append(fptrs, 0)
		} else {
			fptrs = append(fptrs, 1)
		}
		ta, tb, tc := r.count(wrlocked)
		sizes[0] = sizes[0] + ta
		sizes = append(sizes, ta)
		lower[0] = lower[0] + tb
		lower[1] = lower[1] + tc
	}

	for _, r := range tbox.Curdir.subfile { // fsys access point
		names = append(names, r.name)
		times = append(times, r.time)
		sizes = append(sizes, r.size)
		fptrs = append(fptrs, r.fptr)
	}
	return names, times, sizes, fptrs, lower
}

// import binary data under Curdir(->tgtdir), file will be replaced if path exists
func (tbox *PEVFS) ImBin(name string, data []byte) error {
	if ferr := tbox.sub_auth(true); ferr != nil {
		return ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	tbox.fpush_buf() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")
	pos := tbox.fpush_check(tbox.tgtdir, name)

	key := [48]byte(kaes.Genrand(48)) // kaes encrypt
	os.Remove(tbox.data.local + "io/temp.bin")
	tbox.module.enc.Before.Open(data, true)
	tbox.module.enc.After.Open(tbox.data.local+"io/temp.bin", false)
	err := tbox.module.ende(key, true)
	tbox.module.enc.Before.Close()
	tbox.module.enc.After.Close()

	if err != nil { // kaes fail
		tbox.module.log.abort = true
		tbox.module.log.write(fmt.Sprintf("critical : encrypt fail -%s", err))
		return err
	} else {
		fptr := tbox.module.fphymod.fpush(tbox.data.local + "io/temp.bin")
		if fptr == -1 { // fphy fpush fail
			tbox.module.log.abort = true
			tbox.module.log.write(fmt.Sprintf("critical : fpush fail -%s", name))
			return errors.New("fphy fpush error")
		} else {

			tbox.lock.Lock()
			tbox.tgtdir.subfile[pos].size = kio.Size(tbox.data.local + "io/temp.bin")
			tbox.tgtdir.subfile[pos].fptr = fptr
			tbox.lock.Unlock()

			err = tbox.module.fkeymod.push(fptr, key)
			if err != nil { // fkey push fail
				tbox.module.log.abort = true
				tbox.module.log.write(fmt.Sprintf("critical : fkey push fail -%d", fptr))
				return errors.New("fkey push error")
			} else {
				tbox.module.log.write(fmt.Sprintf("msg : fpush -%s", name))
			}
		}
	}

	// sort & header push
	tbox.tgtdir.sort()
	tbox.flush_all()
	tbox.module.log.write(fmt.Sprintf("msg : imbin complete -%d", len(data)))
	return nil
}

// import files under Curdir(->tgtdir), file will be replaced if path exists
func (tbox *PEVFS) ImFiles(paths []string) error {
	if ferr := tbox.sub_auth(true); ferr != nil {
		return ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	tbox.fpush_buf() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	for _, r := range paths { // push files
		r = kio.Abs(r)
		if r[len(r)-1] == '/' {
			return errors.New("invalid file path")
		}
		tbox.fpush_file(r, tbox.tgtdir, r[strings.LastIndex(r, "/")+1:])
		if tbox.module.log.abort {
			return errors.New("abort")
		}
	}

	// sort & header push
	tbox.tgtdir.sort()
	tbox.flush_all()
	tbox.module.log.write(fmt.Sprintf("msg : imfiles complete -%d", len(paths)))
	return nil
}

// import folder under Curdir(->tgtdir), cannot import folder with same name
func (tbox *PEVFS) ImDir(path string) error {
	if ferr := tbox.sub_auth(true); ferr != nil {
		return ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	tbox.fpush_buf() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	path = kio.Abs(path) // checking name
	if path[len(path)-1] != '/' {
		return errors.New("invalid folder path")
	}
	name := path[strings.LastIndex(path[0:len(path)-1], "/")+1:]
	if tbox.fpop_check(tbox.tgtdir, name) != -1 {
		return errors.New("existing folder")
	}

	// generating temp vdir
	temp, err := tbox.dir_import(path)
	if err != nil {
		tbox.module.log.write(fmt.Sprintf("err : ImDir error -%s", err))
		return err
	}

	// sort & header push
	tbox.lock.Lock()
	tbox.tgtdir.subdir = append(tbox.tgtdir.subdir, *temp)
	tbox.tgtdir.sort()
	tbox.lock.Unlock()
	tbox.flush_all()
	tbox.module.log.write(fmt.Sprintf("msg : imdir complete -%s", path))
	return nil
}

// export binary data under Curdir(->tgtdir), find by name
func (tbox *PEVFS) ExBin(name string) ([]byte, error) {
	if ferr := tbox.sub_auth(false); ferr != nil {
		return nil, ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	tbox.module.drv.dirctrl(tbox.data.local+"io/", true) // make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	pos := tbox.fpop_check(tbox.tgtdir, name)
	if pos == -1 { // no such file
		tbox.module.log.abort = true
		tbox.module.log.write(fmt.Sprintf("critical : no file -%s", name))
		return nil, errors.New("no such file")
	}
	tbox.lock.Lock()
	fptr := tbox.tgtdir.subfile[pos].fptr
	size := tbox.tgtdir.subfile[pos].size
	tbox.lock.Unlock()

	ext, _, key := tbox.module.fkeymod.seek(fptr)
	if ext == -1 { // fkey seek fail
		tbox.module.log.abort = true
		tbox.module.log.write(fmt.Sprintf("critical : no key -%d", fptr))
		return nil, errors.New("no key")
	} else {
		if !tbox.module.fphymod.fchk(size, fptr, false) { // fphy broken chain
			tbox.module.log.abort = true
			tbox.module.log.write(fmt.Sprintf("critical : broken chain -%d", fptr))
			return nil, errors.New("broken chain")
		} else {
			os.Remove(tbox.module.fphymod.local + "io/temp.bin") // fphy fpop to local/io/temp.bin
			tbox.module.fphymod.fpop(tbox.data.local + "io/temp.bin")
			tbox.module.fphymod.clear()

			tbox.module.enc.Before.Open(tbox.data.local+"io/temp.bin", true) // kaes decrypt
			tbox.module.enc.After.Open(make([]byte, 0, 10485760), false)
			if err := tbox.module.ende(key, false); err != nil { // kaes decrypt fail
				tbox.module.log.abort = true
				tbox.module.enc.Before.Close()
				tbox.module.enc.After.Close()
				tbox.module.log.write(fmt.Sprintf("critical : decrypt fail -%s", err))
				return nil, err

			} else {
				tbox.module.log.write("msg : fpop -Binary")
				tbox.module.enc.Before.Close()
				temp := tbox.module.enc.After.Close()
				tbox.module.log.write(fmt.Sprintf("msg : exbin complete -%d", len(temp)))
				return temp, nil
			}
		}
	}
}

// export files under Curdir(->tgtdir), find by name, !! generate desktop/kv5export/ !!
func (tbox *PEVFS) ExFiles(names []string) error {
	if ferr := tbox.sub_auth(false); ferr != nil {
		return ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	tbox.module.drv.dirctrl(tbox.data.local+"io/", true) // make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")
	tbox.module.drv.dirctrl(tbox.data.desktop+"kv5export/", true) // make desktop/kv5export/

	for _, r := range names { // pop files
		tbox.fpop_file(tbox.data.desktop+"kv5export/", tbox.tgtdir, r)
		if tbox.module.log.abort {
			return errors.New("abort")
		}
	}
	tbox.module.log.write(fmt.Sprintf("msg : exfiles complete -%d", len(names)))
	return nil
}

// export path folder, !! generate desktop/kv5export/ !!
func (tbox *PEVFS) ExDir(path string) error {
	if ferr := tbox.sub_auth(false); ferr != nil {
		return ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.dir_navigate(path)
		tbox.tgtpath = path
	}
	if tbox.tgtdir == nil {
		return errors.New("invalid path")
	}
	tbox.module.drv.dirctrl(tbox.data.local+"io/", true) // make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")
	tbox.module.drv.dirctrl(tbox.data.desktop+"kv5export/", true) // make desktop/kv5export/

	tbox.dir_export(tbox.data.desktop+"kv5export/", tbox.tgtdir)
	if tbox.module.log.abort {
		return errors.New("abort")
	} else {
		tbox.module.log.write(fmt.Sprintf("msg : exdir complete -%s", path))
		return nil
	}
}

// delete paths, cannot delete / /_BIN/, ! reset cachedir !
func (tbox *PEVFS) Delete(paths []string) error {
	if ferr := tbox.sub_auth(true); ferr != nil {
		return ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	tbox.fpush_buf() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")
	defer tbox.sub_recache() // reset cachedir

	for _, r := range paths { // delete by folder/file path
		switch {
		case r == "/": // rootdir
			return errors.New("del fail rootdir")
		case r == "/_BIN/": // bindir
			return errors.New("del fail bindir")

		case r[len(r)-1] == '/': // tgt is folder
			upper := tbox.dir_navigate(r[0 : strings.LastIndex(r[0:len(r)-1], "/")+1])
			folder := tbox.dir_navigate(r)
			if folder == nil {
				return errors.New("invalid path")
			} else {
				fptr, size := tbox.dir_getf(folder)
				name := r[strings.LastIndex(r[0:len(r)-1], "/")+1:]
				pos := tbox.fpop_check(upper, name)
				tbox.sub_recache()
				upper.subdir = append(upper.subdir[:pos], upper.subdir[pos+1:]...)
				for i, r := range fptr {
					tbox.module.fdel(r, size[i])
				}
			}

		default: // tgt is file
			folder := tbox.dir_navigate(r[0 : strings.LastIndex(r, "/")+1])
			if folder == nil {
				return errors.New("invalid path")
			} else {
				pos := tbox.fpop_check(folder, r[strings.LastIndex(r, "/")+1:])
				if pos == -1 {
					return errors.New("invalid path")
				} else {
					tbox.lock.Lock()
					fptr := folder.subfile[pos].fptr
					size := folder.subfile[pos].size
					folder.subfile = append(folder.subfile[:pos], folder.subfile[pos+1:]...)
					tbox.lock.Unlock()
					tbox.module.fdel(fptr, size)
				}
			}
		}

		if tbox.module.log.abort {
			return errors.New("abort")
		}
	}

	// header push
	tbox.flush_all()
	tbox.module.log.write(fmt.Sprintf("msg : deleted objects -%d", len(paths)))
	return nil
}

// move dir/files in Curdir(->tgtdir) to dst folder, return error if hierarchy problem is detected, ! reset cachedir !
func (tbox *PEVFS) Move(names []string, dst string) error {
	// cannot move /_BIN/ or overlapping names (except move to /_BIN/)
	if ferr := tbox.sub_auth(true); ferr != nil {
		return ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	tbox.fpush_buf() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")
	defer tbox.sub_recache() // reset cachedir

	// get dst folder
	dst_dir := tbox.dir_navigate(dst)
	if dst_dir == nil {
		return errors.New("invalid dst path")
	}

	// check hierarchy & name
	for _, r := range names {
		if tbox.data.contains(tbox.tgtpath+r, dst) {
			return errors.New("invalid src path")
		}
		if tbox.fpop_check(tbox.tgtdir, r) == -1 {
			return errors.New("invalid src path")
		}
		if dst != "/_BIN/" && tbox.fpop_check(dst_dir, r) != -1 {
			return errors.New("overlapping name")
		}
		if r == "_BIN/" && tbox.tgtpath == "/" {
			return errors.New("cannot move bindir")
		}
	}

	// move one by one
	for _, r := range names {
		pos := tbox.fpop_check(tbox.tgtdir, r)
		if r[len(r)-1] == '/' { // move tgt is folder
			dst_dir.subdir = append(dst_dir.subdir, tbox.tgtdir.subdir[pos])
			tbox.tgtdir.subdir = append(tbox.tgtdir.subdir[:pos], tbox.tgtdir.subdir[pos+1:]...)
			tbox.sub_recache()
		} else { // move tgt is file
			dst_dir.subfile = append(dst_dir.subfile, tbox.tgtdir.subfile[pos])
			tbox.tgtdir.subfile = append(tbox.tgtdir.subfile[:pos], tbox.tgtdir.subfile[pos+1:]...)
		}
	}

	// sort & header push
	dst_dir.sort()
	tbox.flush_all()
	tbox.module.log.write(fmt.Sprintf("msg : moved objects -%d", len(names)))
	return nil
}

// rename dir/files in Curdir(->tgtdir), ! reset cachedir !
func (tbox *PEVFS) Rename(before []string, after []string) error {
	// cannot rename /_BIN/ or overlapping names
	if ferr := tbox.sub_auth(true); ferr != nil {
		return ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	tbox.fpush_buf() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")
	defer tbox.sub_recache() // reset cachedir

	// check name
	for i, r := range before {
		l := after[i]
		if tbox.fpop_check(tbox.tgtdir, r) == -1 {
			return errors.New("invalid before path")
		}
		if tbox.fpop_check(tbox.tgtdir, l) != -1 {
			return errors.New("overlapping name")
		}
		if len(l) == 0 {
			return errors.New("invalid after path")
		}

		if r[len(r)-1] == '/' {
			if l[len(l)-1] != '/' {
				return errors.New("invalid after path")
			}
			if strings.Contains(l[0:len(l)-1], "/") {
				return errors.New("invalid after path")
			}
		} else {
			if strings.Contains(l, "/") {
				return errors.New("invalid after path")
			}
		}
		if r == "_BIN/" && tbox.tgtpath == "/" {
			return errors.New("cannot rename bindir")
		}
	}

	// rename one by one
	for i, r := range before {
		l := after[i]
		pos := tbox.fpop_check(tbox.tgtdir, r)
		if r[len(r)-1] == '/' { // move tgt is folder
			tbox.tgtdir.subdir[pos].name = l
			tbox.sub_recache()
		} else { // move tgt is file
			tbox.tgtdir.subfile[pos].name = l
		}
	}

	// sort & header push
	tbox.tgtdir.sort()
	tbox.flush_all()
	tbox.module.log.write(fmt.Sprintf("msg : renamed objects -%d", len(before)))
	return nil
}

// generate new folder at Curdir(->tgtdir), cannot make empty or overlapping names
func (tbox *PEVFS) DirNew(names []string) error {
	if ferr := tbox.sub_auth(true); ferr != nil {
		return ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	tbox.fpush_buf() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	// check name
	for i, r := range names {
		if r[len(r)-1] != '/' {
			r = r + "/"
			names[i] = r
		}
		if tbox.fpop_check(tbox.tgtdir, r) != -1 {
			return errors.New("overlapping name")
		}
		if strings.Contains(r[0:len(r)-1], "/") || len(r) < 2 {
			return errors.New("invalid path")
		}
	}

	// add folders
	for _, r := range names {
		var temp vdir
		temp.name = r
		temp.time = int(time.Now().Unix())
		temp.islocked = false
		tbox.tgtdir.subdir = append(tbox.tgtdir.subdir, temp)
	}

	// sort & header push
	tbox.tgtdir.sort()
	tbox.flush_all()
	tbox.module.log.write(fmt.Sprintf("msg : new folders -%d", len(names)))
	return nil
}

// change lock status of path folder to islocked, changes all subdir if sub T
func (tbox *PEVFS) DirLock(path string, islocked bool, sub bool) error {
	// cannot change lock status of / /_BIN/
	if ferr := tbox.sub_auth(true); ferr != nil {
		return ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	tbox.fpush_buf() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	// check validity
	if path == "/" || path == "/_BIN/" {
		return errors.New("invalid path")
	}
	tbox.tgtdir = tbox.dir_navigate(path)
	tbox.tgtpath = path
	if tbox.tgtdir == nil {
		return errors.New("invalid path")
	}

	tbox.lock.Lock()
	if sub { // reset lock
		tbox.module.relock(tbox.tgtdir, islocked)
	} else {
		tbox.tgtdir.islocked = islocked
	}
	tbox.lock.Unlock()

	// header push
	tbox.flush_all()
	tbox.module.log.write(fmt.Sprintf("msg : set lock -%t %s", islocked, path))
	return nil
}

// check validity under Curdir(-> tgtdir), returns invalid file num
func (tbox *PEVFS) CluCheck() (int, error) {
	if ferr := tbox.sub_auth(true); ferr != nil {
		return 0, ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}

	temp := tbox.sub_check(tbox.tgtdir)
	tbox.module.log.write(fmt.Sprintf("msg : cluster check done -errfile %d", temp))
	return temp, nil
}

// rename : check/fix name, rewrite : reset unallocated chunk, rebuild : check/fix file system
func (tbox *PEVFS) CluRestore(rename bool, rewrite bool, rebuild bool) (int, error) {
	// !! all works are about entire cluster !!, returns wrong file num
	if ferr := tbox.sub_auth(true); ferr != nil {
		return 0, ferr
	} else {
		defer func() {
			tbox.module.log.working = false
			tbox.tgtdir = nil
			tbox.tgtpath = ""
		}()
		tbox.module.log.working = true
		tbox.tgtdir = tbox.Curdir
		tbox.tgtpath = tbox.Curpath
	}
	temp := 0
	defer tbox.sub_recache() // reset cachedir

	if rename { // rename invalid name
		tbox.module.log.write("msg : start rename -cluster restore")
		res := make(chan int, 1)
		go tbox.module.reset(&tbox.module.fsysmod, res)
		temp = temp + <-res
	}

	if rewrite { // rewrite 128 chunk to 0
		tbox.module.log.write("msg : start rewrite -cluster restore")
		temp = temp + tbox.module.fphymod.fgc(true)
	}

	if rebuild { // clear & save only valid file
		tbox.module.log.write("msg : start rebuild -cluster restore")
		var fkey_temp keymap
		fkey_temp.read(nil)
		temp = temp + tbox.sub_rebuild(&tbox.module.fsysmod, &fkey_temp)
		tbox.module.fkeymod = fkey_temp
		tbox.module.log.write("msg : start fGC -cluster restore")
		temp = temp + tbox.module.fphymod.fgc(false)
	}

	// header push
	tbox.flush_all()
	tbox.module.log.write(fmt.Sprintf("msg : cluster restore done -%d", temp))
	return temp, nil
}

// ===== section ===== pevfs commander shell

// wrapper of PEVFS, never panics
type Shell struct {
	// worker modules
	InSys    *PEVFS   // pevfs worker session
	AsyncErr string   // error result of async work
	IOstr    []string // IO string buffer (manual access)
	IObyte   [][]byte // IO []byte buffer (manual access)

	// control flags
	FlagWk bool // isworking flag (readonly)
	FlagRo bool // readonly cluster flag (readonly)
	FlagSz bool // update time&size of CurDir flag

	// info of CurDir (readonly)
	CurPath string   // current session fullpath (~/)
	CurNum  [2]int   // number of subdir&subfile
	CurName []string // names of subdir&subfile
	CurLock []bool   // islocked flag, false if file
	CurTime []string // formatted time of subdir&subfile
	CurSize []int    // sizes of subdir&subfile
}

// init shell with FlagSz
func (sh *Shell) init(flagsz bool) {
	sh.AsyncErr = ""
	sh.IOstr = nil
	sh.IObyte = nil

	sh.FlagWk = false
	sh.FlagRo = true
	sh.FlagSz = flagsz

	sh.CurPath = ""
	sh.CurNum = [2]int{0, 0}
	sh.CurName = nil
	sh.CurLock = nil
	sh.CurTime = nil
	sh.CurSize = nil
}

// update info of CurDir with wrlocked flag, TP to rootdir if CurDir not exists
func (sh *Shell) update() {
	// init info & check session
	sh.CurNum = [2]int{0, 0}
	sh.CurName = nil
	sh.CurLock = nil
	sh.CurTime = nil
	sh.CurSize = nil
	if curtemp := sh.InSys.dir_navigate(sh.CurPath); curtemp == nil {
		sh.InSys.Teleport(sh.InSys.Rootpath)
		sh.CurPath = sh.InSys.Rootpath
	} else {
		sh.InSys.Teleport(sh.CurPath)
	}

	if sh.FlagSz { // update size&time
		t0, t1, t2, t3, _ := sh.InSys.NavInfo(true) // nm t sz fptrlock ldlf
		for i, r := range t0[1:] {
			if r[len(r)-1] == '/' {
				sh.CurNum[0] = sh.CurNum[0] + 1
				sh.CurLock = append(sh.CurLock, t3[i+1] == 0)
			} else {
				sh.CurNum[1] = sh.CurNum[1] + 1
				sh.CurLock = append(sh.CurLock, false)
			}
			sh.CurName = append(sh.CurName, r)
			sh.CurTime = append(sh.CurTime, time.Unix(int64(t1[i+1]), 0).Local().Format("2006.01.02;15:04:05"))
			sh.CurSize = append(sh.CurSize, t2[i+1])
		}

	} else { // fill size&time with zero-value
		sh.CurName, sh.CurLock = sh.InSys.NavName()
		sh.CurTime = make([]string, len(sh.CurName))
		sh.CurSize = make([]int, len(sh.CurName))
		for i, r := range sh.CurName {
			sh.CurTime[i] = "1970.01.01;00:00:00"
			if r[len(r)-1] == '/' {
				sh.CurNum[0] = sh.CurNum[0] + 1
			} else {
				sh.CurNum[1] = sh.CurNum[1] + 1
			}
		}
	}

	if len(sh.CurName) != sh.CurNum[0]+sh.CurNum[1] || len(sh.CurName) != len(sh.CurLock) || len(sh.CurName) != len(sh.CurTime) || len(sh.CurName) != len(sh.CurSize) {
		sh.CurNum = [2]int{0, 0} // check if info unmatch exists
	}
}

// async cluster manipulate, error log will be AsyncErr
func (sh *Shell) asyncwork(worktype int, parms []string) {
	defer func() {
		if ferr := recover(); ferr != nil {
			sh.AsyncErr = fmt.Sprintf("critical : %s", ferr)
		}
		sh.FlagWk = false
	}()
	sh.FlagWk = true

	var err error
	switch worktype {
	case 0: // imbin
		err = sh.InSys.ImBin(parms[0], sh.IObyte[0])
	case 1: // imfile
		err = sh.InSys.ImFiles(parms)
	case 2: // imdir
		err = sh.InSys.ImDir(parms[0])
	case 3: // exbin
		sh.IObyte = make([][]byte, len(parms))
		for i, r := range parms {
			sh.IObyte[i], err = sh.InSys.ExBin(r)
		}
	case 4: // exfile
		err = sh.InSys.ExFiles(parms)
	case 5: // exdir
		err = sh.InSys.ExDir(parms[0])
	case 6: // delete
		err = sh.InSys.Delete(parms)
	case 7: // move
		err = sh.InSys.Move(parms[1:], parms[0])
	case 8: // rename
		err = sh.InSys.Rename(parms[:len(parms)/2], parms[len(parms)/2:])
	case 9: // dirnew
		err = sh.InSys.DirNew(parms)
	case 10: // dirlock
		for i := 0; i < len(parms)/3; i++ {
			err = sh.InSys.DirLock(parms[3*i], parms[3*i+1] == "true", parms[3*i+2] == "true")
		}
	case 11: // check
		var ti int
		ti, err = sh.InSys.CluCheck()
		sh.IOstr[0] = fmt.Sprint(ti)
	case 12: // restore
		var ti int
		ti, err = sh.InSys.CluRestore(parms[0] == "true", parms[1] == "true", parms[2] == "true")
		sh.IOstr[0] = fmt.Sprint(ti)
	}

	if err != nil {
		sh.AsyncErr = fmt.Sprint(err)
	}
	sh.update()
}

// manipulate session with order/option, returns if started successfully
func (sh *Shell) Command(order string, option []string) (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("shell fail : %s", ferr)
		}
		time.Sleep(time.Millisecond * 50)
	}()
	err = nil
	sh.AsyncErr = ""

	switch order {
	case "init": // init Shell -flagsz
		sh.init(option[0] == "true")

	case "new": // PEVFS_New -remote cluster chunksize
		if option[0][len(option[0])-1] != '/' {
			err = errors.New("invalid remote path")
		} else if len(option[1]) == 0 {
			err = errors.New("invalid cluster name")
		} else {
			switch option[2] {
			case "small":
				err = PEVFS_New(option[0], option[1], 4096)
			case "standard":
				err = PEVFS_New(option[0], option[1], 32768)
			case "large":
				err = PEVFS_New(option[0], option[1], 262144)
			default:
				err = PEVFS_New(option[0], option[1], 512)
			}
		}

	case "boot": // PEVFS_Boot -desktop local remote blockA
		if option[0][len(option[0])-1] != '/' {
			err = errors.New("invalid desktop path")
		} else if option[1][len(option[1])-1] != '/' {
			err = errors.New("invalid local path")
		} else if option[2][len(option[2])-1] != '/' {
			err = errors.New("invalid remote path")
		} else {
			sh.IObyte = make([][]byte, 1) // -hint
			sh.InSys, sh.IObyte[0], err = PEVFS_Boot(option[0], option[1], option[2], option[3])
			sh.IOstr = make([]string, 2) // -cluster account
			sh.IOstr[0], sh.IOstr[1] = *sh.InSys.Cluster, *sh.InSys.Account
			sh.FlagRo = sh.InSys.module.log.readonly
		}

	case "exit": // PEVFS_Exit
		PEVFS_Exit(sh.InSys)
		sh.init(false)

	case "rebuild": // PEVFS_Rebuild -remote
		defer func() { sh.IObyte = nil }() // -pw kf
		if option[0][len(option[0])-1] != '/' {
			err = errors.New("invalid remote path")
		} else {
			err = PEVFS_Rebuild(option[0], sh.IObyte[0], sh.IObyte[1])
		}

	case "abort": // abort order/check -reset abort working
		_, sh.FlagWk = sh.InSys.Abort(option[0] == "true", option[1] == "true", option[2] == "true")

	case "debug": // get debug info string -count_locked
		sh.IOstr = make([]string, 4) // -debug *4
		ab, wr := sh.InSys.Abort(false, false, false)
		sh.IOstr[0] = "===== PEVFS public data =====\n"
		sh.IOstr[0] = sh.IOstr[0] + fmt.Sprintf("Cluster : %s\nAccount : %s\n", *sh.InSys.Cluster, *sh.InSys.Account)
		sh.IOstr[0] = sh.IOstr[0] + fmt.Sprintf("RootPath : %s\nCurPath : %s\n", sh.InSys.Rootpath, sh.InSys.Curpath)
		sh.IOstr[0] = sh.IOstr[0] + fmt.Sprintf("Abort : %t\nIsWorking : %t\n", ab, wr)

		t0, t1 := sh.InSys.Debug() // [csize, bnum, fptr], [wrsign, salt, pwhash, fsyskey, fkeykey, fphykey]
		sh.IOstr[1] = "===== PEVFS private data =====\n"
		sh.IOstr[1] = sh.IOstr[1] + fmt.Sprintf("chunksize : %d\nblocknum : %d\nempty fptr : %d\n", t0[0], t0[1], t0[2])
		sh.IOstr[1] = sh.IOstr[1] + fmt.Sprintf("wrsign[8B] : %s\nsalt[64B] : %s\npwhash[192B] : %s\n", kio.Bprint(t1[0]), kio.Bprint(t1[1]), kio.Bprint(t1[2]))
		sh.IOstr[1] = sh.IOstr[1] + fmt.Sprintf("fsys key[48B] : %s\nfkey key[48B] : %s\nfphy key[48B] : %s\n", kio.Bprint(t1[3]), kio.Bprint(t1[4]), kio.Bprint(t1[5]))

		sh.IOstr[2] = "===== Shell public data =====\n"
		sh.IOstr[2] = sh.IOstr[2] + fmt.Sprintf("AsyncErr : %s\nlen IOstr : %d\nlen IObyte : %d\n", sh.AsyncErr, len(sh.IOstr), len(sh.IObyte))
		sh.IOstr[2] = sh.IOstr[2] + fmt.Sprintf("isworking : %t\nreadonly : %t\nviewsize : %t\n", sh.FlagWk, sh.FlagRo, sh.FlagSz)

		t2, t3, t4, t5, t6 := sh.InSys.NavInfo(option[0] == "true") // nm t sz fptrlock ldlf
		sh.IOstr[3] = "===== Shell curdir data =====\n"
		sh.IOstr[3] = sh.IOstr[3] + fmt.Sprintf("path : %s\ntime : %s\n", sh.InSys.Curpath, time.Unix(int64(t3[0]), 0).Local().Format("2006.01.02;15:04:05"))
		sh.IOstr[3] = sh.IOstr[3] + fmt.Sprintf("size : %dB\nlock : %t\nfolder num : %d\nfile num : %d\n", t4[0], t5[0] == 0, t6[0], t6[1])
		for i, r := range t2[1:] {
			if r[len(r)-1] == '/' {
				sh.IOstr[3] = sh.IOstr[3] + fmt.Sprintf("%s (%s, %d B, lock %t)\n", r, time.Unix(int64(t3[i+1]), 0).Local().Format("2006.01.02;15:04:05"), t4[i+1], t5[i+1] == 0)
			} else {
				sh.IOstr[3] = sh.IOstr[3] + fmt.Sprintf("%s (%s, %d B, fptr %d)\n", r, time.Unix(int64(t3[i+1]), 0).Local().Format("2006.01.02;15:04:05"), t4[i+1], t5[i+1])
			}
		}

	case "log": // log read/clear -reset
		sh.IOstr = make([]string, 1) // -log_data
		sh.IOstr[0] = sh.InSys.Log(option[0] == "true")

	case "login": // login to cluster -sleep
		defer func() { sh.IObyte = nil }() // -pw kf
		switch option[0] {
		case "1":
			err = sh.InSys.Login(sh.IObyte[0], sh.IObyte[1], 1)
		case "10":
			err = sh.InSys.Login(sh.IObyte[0], sh.IObyte[1], 10)
		case "30":
			err = sh.InSys.Login(sh.IObyte[0], sh.IObyte[1], 30)
		case "60":
			err = sh.InSys.Login(sh.IObyte[0], sh.IObyte[1], 60)
		default:
			err = sh.InSys.Login(sh.IObyte[0], sh.IObyte[1], 4)
		}
		if err == nil {
			sh.CurPath = sh.InSys.Curpath
			sh.update()
		}

	case "reset": // reset account pwkf
		defer func() { sh.IObyte = nil }() // -pw kf hint
		err = sh.InSys.AccReset(sh.IObyte[0], sh.IObyte[1], sh.IObyte[2])

	case "extend": // extend account -account wrlocked
		defer func() { sh.IObyte = nil }() // -pw kf hint
		if len(option[0]) == 0 || option[0] == "root" {
			err = errors.New("invalid account name")
		} else {
			sh.IOstr = make([]string, 1) // -new blockA path
			sh.IOstr[0], err = sh.InSys.AccExtend(sh.IObyte[0], sh.IObyte[1], sh.IObyte[2], option[0], option[1] == "true")
		}

	case "search": // search for name -name
		sh.IOstr = []string{strings.Join(sh.InSys.Search(option[0]), "\n")} // -result joined with \n

	case "print": // print fsys structure -wrlocked
		sh.IOstr = []string{sh.InSys.Print(option[0] == "true")} // -result

	case "update": // CurDir info update
		sh.update()

	case "navigate": // get subdir names of tgt folder -path
		sh.IOstr = make([]string, 1) // -folder names joined with \n
		temp := sh.InSys.dir_navigate(option[0])
		if temp == nil {
			err = errors.New("invalid path")
		} else if len(temp.subdir) != 0 {
			for _, r := range temp.subdir {
				sh.IOstr[0] = sh.IOstr[0] + r.name + "\n"
			}
			sh.IOstr[0] = sh.IOstr[0][0 : len(sh.IOstr[0])-1]
		}

	case "imbin": // import binary -name
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else if option[0][len(option[0])-1] == '/' {
			err = errors.New("invalid name")
		} else { // IObyte -data
			go sh.asyncwork(0, option)
		}

	case "imfile": // import files -fullpaths ...
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			temp := make([]string, 0)
			for _, r := range option {
				if r[len(r)-1] != '/' {
					temp = append(temp, r)
				}
			}
			go sh.asyncwork(1, temp)
		}

	case "imdir": // import folder -fullpath
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			temp := kio.Abs(option[0])
			if temp[len(temp)-1] != '/' {
				err = errors.New("invalid path")
			} else if slices.Contains(sh.CurName, temp[strings.LastIndex(temp[0:len(temp)-1], "/")+1:]) {
				err = errors.New("overlapping path")
			} else {
				go sh.asyncwork(2, []string{temp})
			}
		}

	case "exbin": // export binary -name *n
		if sh.FlagWk {
			err = errors.New("invalid order")
		} else {
			temp := make([]string, 0)
			for _, r := range option {
				if slices.Contains(sh.CurName, r) {
					temp = append(temp, r)
				}
			}
			go sh.asyncwork(3, temp) // IObytes -bin *n
		}

	case "exfile": // export files -name *n
		if sh.FlagWk {
			err = errors.New("invalid order")
		} else {
			temp := make([]string, 0)
			for _, r := range option {
				if slices.Contains(sh.CurName, r) {
					temp = append(temp, r)
				}
			}
			go sh.asyncwork(4, temp)
		}

	case "exdir": // export folder -fullpath
		if sh.FlagWk {
			err = errors.New("invalid order")
		} else if tgt := sh.InSys.dir_navigate(option[0]); tgt == nil {
			err = errors.New("invalid path")
		} else {
			go sh.asyncwork(5, option)
		}

	case "delete": // delete file/folder -number of name *n
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			temp := make([]string, 0)
			for _, r := range option {
				ti, ierr := strconv.Atoi(r)
				if ierr == nil && ti < len(sh.CurName) {
					path := sh.CurPath + sh.CurName[ti]
					if path != "/_BIN/" && path != "/_BUF" {
						temp = append(temp, path)
					}
				}
			}
			go sh.asyncwork(6, temp)
		}

	case "move": // move file/folder to dst folder -tgtdir, number of name *n
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else if tgt := sh.InSys.dir_navigate(option[0]); tgt == nil {
			err = errors.New("invalid path")
		} else {
			temp := []string{option[0]}
			for _, r := range option[1:] {
				ti, ierr := strconv.Atoi(r)
				if ierr == nil && ti < len(sh.CurName) {
					name := sh.CurName[ti]
					if sh.CurPath != "/" || (name != "_BIN/" && name != "_BUF") {
						temp = append(temp, name)
					}
				}
			}
			go sh.asyncwork(7, temp)
		}

	case "rename": // rename file/folder -number of name *n
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			temp0 := make([]string, 0)
			temp1 := make([]string, 0)
			for i, r := range option {
				ti, ierr := strconv.Atoi(r) // IOstr -newname *n
				if ierr == nil && ti < len(sh.CurName) && !slices.Contains(sh.CurName, sh.IOstr[i]) {
					temp0 = append(temp0, sh.CurName[ti])
					temp1 = append(temp1, sh.IOstr[i])
				}
			}
			go sh.asyncwork(8, append(temp0, temp1...))
		}

	case "dirnew": // make new folder -name *n
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			temp := make([]string, 0)
			for _, r := range option {
				if len(r) > 1 && !slices.Contains(sh.CurName, r) {
					temp = append(temp, r)
				}
			}
			go sh.asyncwork(9, temp)
		}

	case "dirlock": // relock folder -lock, number of name *n (-1 is CurDir)
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			temp := make([]string, 0) // (path lock sub) *n
			for _, r := range option[1:] {
				ti, ierr := strconv.Atoi(r)
				if ierr == nil {
					if ti == -1 {
						temp = append(append(append(temp, sh.CurPath), option[0]), "true")
					} else if ti < sh.CurNum[0] {
						temp = append(append(append(temp, sh.CurPath+sh.CurName[ti]), option[0]), "false")
					}
				}
			}
			go sh.asyncwork(10, temp)
		}

	case "check": // cluster check
		sh.IOstr = []string{"0"} // IOstr -broken files
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			go sh.asyncwork(11, nil)
		}

	case "restore": // cluster restore -mode
		sh.IOstr = []string{"0"} // IOstr -restored files
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			var temp []string
			switch option[0] {
			case "rename":
				temp = []string{"true", "false", "false"}
			case "rewrite":
				temp = []string{"false", "true", "false"}
			case "rebuild":
				temp = []string{"false", "false", "true"}
			default:
				temp = []string{"false", "false", "false"}
			}
			go sh.asyncwork(12, temp)
		}

	default: // unknown order & clear IObuf
		sh.IObyte = nil
		sh.IOstr = nil
		err = errors.New("unknown order")
	}

	return err
}

// ===== section ===== pevfs benchmark test

// cluster generation, basic 4GiB read/write
func Test_Basic(remote string) (float64, float64, float64) {
	fmt.Println("Start : Test_Basic_0")
	tcheck := time.Now()
	PEVFS_New(remote, "test", 32768)
	result0 := float64(time.Since(tcheck).Milliseconds()) / 1000
	fmt.Println("End : Test_Basic_0")
	time.Sleep(3 * time.Second)

	fmt.Println("Start : Test_Basic_1")
	tcheck = time.Now()
	f, _ := kio.Open(remote+"temp.bin", "w")
	for i := 0; i < 8; i++ {
		kio.Write(f, make([]byte, 536870912))
	}
	f.Close()
	result2 := float64(time.Since(tcheck).Milliseconds()) / 1000
	fmt.Println("End : Test_Basic_1")
	time.Sleep(3 * time.Second)

	fmt.Println("Start : Test_Basic_2")
	tcheck = time.Now()
	f, _ = kio.Open(remote+"temp.bin", "r")
	for i := 0; i < 8; i++ {
		kio.Read(f, 536870912)
	}
	f.Close()
	result1 := float64(time.Since(tcheck).Milliseconds()) / 1000
	fmt.Println("End : Test_Basic_2")
	time.Sleep(3 * time.Second)

	f, _ = kio.Open("./io.bin", "w")
	for i := 0; i < 8; i++ {
		kio.Write(f, make([]byte, 536870912))
	}
	f.Close()
	return result0, result1, result2
}

// 4GiB file cluster login/read/write, !! do after Test_Basic !!
func Test_IO(remote string) (float64, float64, float64) {
	fmt.Println("Start : Test_IO_0")
	tcheck := time.Now()
	k, _, _ := PEVFS_Boot("./", "./temp693d/", remote, remote+"0a.webp")
	k.Login([]byte("0000"), kaes.Basickey(), 4)
	defer PEVFS_Exit(k)
	result0 := float64(time.Since(tcheck).Milliseconds()) / 1000
	fmt.Println("End : Test_IO_0")
	time.Sleep(3 * time.Second)

	fmt.Println("Start : Test_IO_1")
	tcheck = time.Now()
	k.ImFiles([]string{"./io.bin"})
	result2 := float64(time.Since(tcheck).Milliseconds()) / 1000
	fmt.Println("End : Test_IO_1")
	time.Sleep(3 * time.Second)

	fmt.Println("Start : Test_IO_2")
	tcheck = time.Now()
	k.ExFiles([]string{"io.bin"})
	result1 := float64(time.Since(tcheck).Milliseconds()) / 1000
	fmt.Println("End : Test_IO_2")
	time.Sleep(3 * time.Second)

	return result0, result1, result2
}

// mk 100 sub (1 lv2 -> 100 lv1 -> 10000 lv0)
func test_sub(pr *vdir, ky *keymap, lv int) {
	if lv == 0 { // add 100 files
		for i := 0; i < 100; i++ {
			var temp vfile
			temp.name = "test_file_low_level"
			temp.time = 1000000000
			temp.size = 2000000000
			temp.fptr = kobj.Decode(kaes.Genrand(4))
			ky.push(temp.fptr, [48]byte(kaes.Genrand(48)))
			pr.subfile = append(pr.subfile, temp)
		}
	} else { // add 100 folders
		pr.subdir = make([]vdir, 100)
		for i := range pr.subdir {
			pr.subdir[i].name = "test_folder_low_level/"
			pr.subdir[i].time = 1000000000
			test_sub(&pr.subdir[i], ky, lv-1)
		}
	}
}

// 100k folder + 10m file generation/read/write
func Test_Multi(remote string) (float64, float64, float64) {
	os.RemoveAll(remote)
	os.Mkdir(remote, os.ModePerm)
	PEVFS_New(remote, "test", 32768)
	k, _, _ := PEVFS_Boot("./", "./temp693d/", remote, remote+"0a.webp")
	k.Login([]byte("0000"), kaes.Basickey(), 4)
	defer PEVFS_Exit(k)

	fmt.Println("Start : Test_Multi_0")
	tcheck := time.Now()
	var t0 vdir
	var t1 keymap
	t0.name = "/"
	t0.time = 1000000000
	t0.islocked = true
	t0.subdir = make([]vdir, 10)
	for i := 0; i < 10; i++ {
		t0.subdir[i].name = "test_folder_low_level/"
		t0.subdir[i].time = 1000000000
		t0.subdir[i].islocked = true
		test_sub(&t0.subdir[i], &t1, 2)
	}
	k.module.fsysmod = t0
	k.module.fkeymod = t1
	result0 := float64(time.Since(tcheck).Milliseconds()) / 1000
	fmt.Println("End : Test_Multi_0")
	fmt.Println(k.module.fsysmod.count(true))
	time.Sleep(3 * time.Second)

	fmt.Println("Start : Test_Multi_1")
	tcheck = time.Now()
	k.flush_all()
	k = nil
	result2 := float64(time.Since(tcheck).Milliseconds()) / 1000
	fmt.Println("End : Test_Multi_1")
	time.Sleep(3 * time.Second)

	fmt.Println("Start : Test_Multi_2")
	tcheck = time.Now()
	k, _, _ = PEVFS_Boot("./", "./temp693d/", remote, remote+"0a.webp")
	k.Login([]byte("0000"), kaes.Basickey(), 4)
	result1 := float64(time.Since(tcheck).Milliseconds()) / 1000
	fmt.Println("End : Test_Multi_2")
	time.Sleep(3 * time.Second)

	return result0, result1, result2
}
