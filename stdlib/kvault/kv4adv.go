// test703 : stdlib5.kv4adv

package kvault

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"stdlib5/kaes"
	"stdlib5/kio"
	"stdlib5/kobj"
	"stdlib5/ksc"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ===== section ===== near/far file copy

// safe file copy (about cloud IO)
type g4fcopy struct {
	log *logger   // common logger
	drv *basework // common safe IO
}

// safe far file copy
func (fm *g4fcopy) farcopy_inline(src string, dst string) (res error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			res = fmt.Errorf("critical : %s", ferr)
			fm.log.abort = true
		}
	}()

	size := kio.Size(src) // step 0 : get size, file open
	if size < 0 {
		return errors.New("source not exist")
	}
	f, err0 := os.Open(src)
	if err0 == nil {
		defer f.Close()
	} else {
		return err0
	}
	t, err1 := os.Create(dst)
	if err1 == nil {
		defer t.Close()
	} else {
		return err1
	}

	var temp []byte // step 1 : copy 100MiB
	for i := 0; i < size/104857600; i++ {
		temp, err0 = kio.Read(f, 104857600)
		_, err1 = kio.Write(t, temp)
		if err0 != nil {
			return err0
		}
		if err1 != nil {
			return err1
		}
	}
	if size%104857600 != 0 {
		temp, err0 = kio.Read(f, size%104857600)
		_, err1 = kio.Write(t, temp)
		if err0 != nil {
			return err0
		}
		if err1 != nil {
			return err1
		}
	}
	return nil
}

// init struct
func (fm *g4fcopy) init(logger *logger, iodrv *basework) {
	fm.log = logger
	fm.drv = iodrv
}

// get names of file under far folder
func (fm *g4fcopy) farsub(path string) []string {
	return fm.drv.dirsub(path)
}

// far folder dirctrl
func (fm *g4fcopy) fardir(path string, ctrl bool) {
	fm.drv.dirctrl(path, ctrl)
}

// rename file (domain ~/ src -> dst),  !! cloud IO is writed here !!
func (fm *g4fcopy) farname(domain string, src string, dst string) {
	flag := true
	for flag {
		if err := os.Rename(domain+src, domain+dst); err == nil {
			flag = false
		} else {
			fm.log.write(fmt.Sprintf("err : farname fail -%s", err))
			time.Sleep(time.Second * time.Duration(fm.drv.sleep))
		}
	}
}

// delete far file
func (fm *g4fcopy) fardel(path string) {
	os.Remove(path)
}

// remote file copy !! cloud IO is writed here !!
func (fm *g4fcopy) farcopy(src string, dst string) {
	flag := true
	for flag {
		if fm.log.abort {
			fm.log.write("msg : farcopy aborted")
			flag = false
			break
		}
		err := fm.farcopy_inline(src, dst)
		if err == nil {
			flag = false
			fm.log.write("msg : farcopy complete")
		} else {
			fm.log.write(fmt.Sprintf("err : farcopy fail -%s", err))
			time.Sleep(time.Second * time.Duration(fm.drv.sleep))
		}
	}
}

// local near file copy
func (fm *g4fcopy) nearcopy(src string, pos int, size int, dst string) {
	var i int
	for i = 0; i < size/104857600; i++ {
		fm.drv.write(dst, i*104857600, fm.drv.read(src, pos+i*104857600, 104857600))
	}
	i = size / 104857600
	if size%104857600 != 0 {
		fm.drv.write(dst, i*104857600, fm.drv.read(src, pos+i*104857600, size%104857600))
	}
}

// local near file merge
func (fm *g4fcopy) nearmerge(src []string, dst string) {
	pos := 0
	size := 0
	for _, r := range src {
		size = kio.Size(r)
		if size < 0 {
			fm.log.write("err : merge fail -invalid path")
			break
		}

		var i int
		for i = 0; i < size/104857600; i++ {
			fm.drv.write(dst, pos+i*104857600, fm.drv.read(r, i*104857600, 104857600))
		}
		i = size / 104857600
		if size%104857600 != 0 {
			fm.drv.write(dst, pos+i*104857600, fm.drv.read(r, i*104857600, size%104857600))
		}
		pos = pos + size
	}
}

// ===== section ===== physical folder memory

// file (fptr, num) with path
type g4fchunk struct {
	namecode [16]byte  // actual file namecode
	fragnum  int       // file chunk number
	next     *g4fchunk // next file chunk
}

// set namecode, fragnum
func (fm *g4fchunk) set(namecode []byte, fragnum int) {
	if len(namecode) == 16 {
		fm.namecode = [16]byte(namecode)
	} else {
		fm.namecode = [16]byte{}
	}
	fm.fragnum = fragnum
	fm.next = nil
}

// link another node with same fptr
func (fm *g4fchunk) add(node *g4fchunk) {
	if node.fragnum < fm.fragnum { // link forward
		temp0 := node.namecode
		temp1 := node.fragnum
		node.namecode = fm.namecode
		node.fragnum = fm.fragnum
		fm.namecode = temp0
		fm.fragnum = temp1
		node.next = fm.next
		fm.next = node
	} else if fm.next == nil { // link backward (nil)
		fm.next = node
	} else { // link backward
		tgt := fm
		for tgt.next != nil && tgt.fragnum < node.fragnum {
			tgt = tgt.next
		}
		tgt.add(node)
	}
}

// check if fragnum is increaing well
func (fm *g4fchunk) check() bool {
	tgt := fm
	for tgt.next != nil {
		if tgt.fragnum+1 != tgt.next.fragnum {
			return false
		}
		tgt = tgt.next
	}
	return true
}

// folder page data
type g4fpage struct {
	dirnum int                 // folder number
	key    [48]byte            // filename enckey
	table  map[int](*g4fchunk) // table filechunk[fptr]
}

// set dirnum, key
func (fm *g4fpage) set(dirnum int, key [48]byte) {
	fm.dirnum = dirnum
	fm.key = key
	fm.table = make(map[int](*g4fchunk))
}

// add by filename (str32 + .kv4)
func (fm *g4fpage) add(fname string, worker *basework) {
	temp, _ := kio.Bread(fname[0:32])
	temp = worker.aescalc(temp, fm.key, false, false)
	fptr := kobj.Decode(temp[8:13])
	fnum := kobj.Decode(temp[13:16])
	var newchunk g4fchunk
	newchunk.set(temp, fnum)
	if _, ext := fm.table[fptr]; ext { // existing chain
		fm.table[fptr].add(&newchunk)
	} else { // new chain
		fm.table[fptr] = &newchunk
	}
}

// check if page table is valid, returns invalid file names
func (fm *g4fpage) check(worker *basework) []string {
	invalid := make([]string, 0)
	for i, r := range fm.table {
		if r.fragnum != 0 || !r.check() {
			tgt := r
			for tgt != nil {
				invalid = append(invalid, kio.Bprint(worker.aescalc(tgt.namecode[:], fm.key, true, false))+".kv4")
				tgt = tgt.next
			}
			delete(fm.table, i)
		}
	}
	return invalid
}

// ===== section ===== kv4adv IO management

// remote file move, file divider
type g4IOmgr struct {
	log *logger   // logger
	drv *basework // local IO
	fcp *g4fcopy  // file copy

	local   string    // local path (~/)
	remote  string    // remote path (~/)
	fphykey [128]byte // file name enc key

	maxfptr int       // max fptr in 1 folder
	divsize int       // remote chunk size
	curdir  int       // current writing folder
	table   []g4fpage // folder data page
}

// get file name enc key by (dirnum, fphykey)
func (mgr *g4IOmgr) getkey_inline(dirnum int) [48]byte {
	temp := make([]byte, 0, 48)
	pos := 16 * (dirnum % 8)
	temp = append(temp, mgr.fphykey[pos:pos+16]...)
	pos = 16 * ((dirnum / 8) % 8)
	temp = append(temp, mgr.fphykey[pos:pos+16]...)
	pos = 16 * ((dirnum / 64) % 8)
	temp = append(temp, mgr.fphykey[pos:pos+16]...)
	return [48]byte(temp)
}

// load table from remote
func (mgr *g4IOmgr) getfptr_inline(dirnum int, access *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done() // page table init
	names := mgr.fcp.farsub(fmt.Sprintf("%s%d/", mgr.remote, dirnum))
	mgr.table[dirnum].set(dirnum, mgr.getkey_inline(dirnum))

	for _, r := range names { // delete invalid names
		if strings.Contains(r, ".kv4") && len(r) == 36 {
			mgr.table[dirnum].add(r, mgr.drv)
		} else {
			access.Lock()
			mgr.fcp.fardel(fmt.Sprintf("%s%d/%s", mgr.remote, dirnum, r))
			access.Unlock()
			mgr.log.write(fmt.Sprintf("err : delete while loading -%s", r))
		}
	}

	for _, r := range mgr.table[dirnum].check(mgr.drv) { // delete broken chain
		access.Lock()
		mgr.fcp.fardel(fmt.Sprintf("%s%d/%s", mgr.remote, dirnum, r))
		access.Unlock()
		mgr.log.write(fmt.Sprintf("err : delete while loading -%s", r))
	}
}

// find fpage by fptr, -1 if not exists
func (mgr *g4IOmgr) findfptr_inline(fptr int) int {
	for i, r := range mgr.table {
		if _, ext := r.table[fptr]; ext {
			return i
		}
	}
	return -1
}

// update curdir, make new folder if need
func (mgr *g4IOmgr) update_inline() {
	for len(mgr.table[mgr.curdir].table) >= mgr.maxfptr {
		mgr.curdir = mgr.curdir + 1
		if mgr.curdir == len(mgr.table) {
			var newpage g4fpage
			newpage.set(mgr.curdir, mgr.getkey_inline(mgr.curdir))
			mgr.table = append(mgr.table, newpage)
			mgr.fcp.fardir(fmt.Sprintf("%s%d/", mgr.remote, mgr.curdir), true)
			mgr.log.write(fmt.Sprintf("msg : new folder generated -%d", mgr.curdir))
			break
		}
	}
}

// load curdir from remote
func (mgr *g4IOmgr) load_inline() {
	temp := mgr.fcp.farsub(mgr.remote) // get folder num
	curdir := 0
	for slices.Contains(temp, fmt.Sprintf("%d/", curdir)) {
		curdir = curdir + 1
	}

	var access sync.Mutex
	var wait sync.WaitGroup
	mgr.table = make([]g4fpage, curdir)
	wait.Add(curdir)
	for i := 0; i < curdir; i++ { // get fptr data for each folder
		go mgr.getfptr_inline(i, &access, &wait)
	}
	wait.Wait()
}

// fpush by name
func (mgr *g4IOmgr) fpush_inline(src string, size int, num int, fptr int) {
	start := num * mgr.divsize
	end := start + mgr.divsize
	if end > size {
		end = size
	}
	key := mgr.table[mgr.curdir].key
	namecode := append(append(kaes.Genrand(8), kobj.Encode(fptr, 5)...), kobj.Encode(num, 3)...)
	tpath0 := fmt.Sprintf("%sbuffer/%d.bin", mgr.local, num)
	tpath1 := fmt.Sprintf("%s%d/", mgr.remote, mgr.curdir)
	tpath2 := kio.Bprint(mgr.drv.aescalc(namecode, key, true, false)) + ".kv4"

	mgr.fcp.nearcopy(src, start, end-start, tpath0)
	mgr.log.write(fmt.Sprintf("msg : nearcopy -fptr %d, num %d", fptr, num))
	mgr.fcp.farcopy(tpath0, tpath1+tpath2)
	mgr.log.write(fmt.Sprintf("msg : farcopy -fptr %d, num %d", fptr, num))
	mgr.table[mgr.curdir].add(tpath2, mgr.drv)
}

// init struct, set basic modules
func (mgr *g4IOmgr) init(iodrv *basework, maxfptr int, divsize int, local string, remote string, fphykey [128]byte) {
	mgr.log = iodrv.log
	mgr.drv = iodrv
	var temp g4fcopy
	temp.init(mgr.log, mgr.drv)
	mgr.fcp = &temp

	mgr.local = local
	mgr.remote = remote
	mgr.fphykey = fphykey

	mgr.maxfptr = maxfptr
	mgr.divsize = divsize
	mgr.curdir = 0
	mgr.table = make([]g4fpage, 0)

	mgr.drv.dirctrl(local, true)
	mgr.load_inline()
	mgr.update_inline()
}

// file push (local -> remote), returns curdir num
func (mgr *g4IOmgr) fpush(path string, fptr int) int {
	if mgr.log.readonly {
		mgr.log.write("err : cannot fpush -readonly cluster")
		return -1
	}
	if mgr.findfptr_inline(fptr) != -1 {
		mgr.log.write("err : cannot fpush -existing fptr")
		return -1
	}
	size := kio.Size(path)
	if size < 0 {
		mgr.log.write("err : cannot fpush -invalid path")
		return -1
	}
	defer mgr.log.write(fmt.Sprintf("msg : fpush complete -fptr %d", fptr))
	mgr.drv.dirctrl(mgr.local+"buffer/", true)
	defer mgr.drv.dirctrl(mgr.local+"buffer/", false) // check cond, setup buffer

	for i := 0; i < size/mgr.divsize; i++ { // fpush inline (near/far copy)
		mgr.fpush_inline(path, size, i, fptr)
	}
	if size%mgr.divsize != 0 {
		mgr.fpush_inline(path, size, size/mgr.divsize, fptr)
	}

	temp := mgr.curdir // update curdir
	mgr.update_inline()
	return temp
}

// file pop (remote -> local), return deleted files num
func (mgr *g4IOmgr) fpop(fptr int) int {
	if mgr.log.readonly {
		mgr.log.write("err : cannot fpop -readonly cluster")
		return -1
	}
	pos := mgr.findfptr_inline(fptr)
	if pos == -1 {
		mgr.log.write("err : cannot fpop -no fptr")
		return -1
	} // check cond, get folder pos
	defer mgr.log.write(fmt.Sprintf("msg : fpop complete -fptr %d", fptr))

	todel := make([]string, 0) // get file names
	tgt := mgr.table[pos].table[fptr]
	for tgt != nil {
		todel = append(todel, kio.Bprint(mgr.drv.aescalc(tgt.namecode[:], mgr.table[pos].key, true, false))+".kv4")
		tgt = tgt.next
	}
	delete(mgr.table[pos].table, fptr)

	for _, r := range todel { // delete files & update curdir
		mgr.fcp.fardel(fmt.Sprintf("%s%d/%s", mgr.remote, pos, r))
	}
	if pos < mgr.curdir {
		mgr.curdir = pos
		mgr.update_inline()
	}
	return len(todel)
}

// file seek (remote -> local), return number of chunk
func (mgr *g4IOmgr) fseek(path string, fptr int) int {
	pos := mgr.findfptr_inline(fptr)
	if pos == -1 {
		mgr.log.write("err : cannot fseek -no fptr")
		return -1
	}
	defer mgr.log.write(fmt.Sprintf("msg : fseek complete -fptr %d", fptr))
	mgr.drv.dirctrl(mgr.local+"buffer/", true)
	defer mgr.drv.dirctrl(mgr.local+"buffer/", false) // check cond, setup buffer

	tgt := mgr.table[pos].table[fptr]
	count := 0
	names := make([]string, 0)
	for tgt != nil { // copy all chunks
		tpath0 := fmt.Sprintf("%sbuffer/%d.bin", mgr.local, count)
		tpath1 := fmt.Sprintf("%s%d/", mgr.remote, pos)
		tpath2 := kio.Bprint(mgr.drv.aescalc(tgt.namecode[:], mgr.table[pos].key, true, false)) + ".kv4"
		names = append(names, tpath0)
		mgr.fcp.farcopy(tpath1+tpath2, tpath0)
		tgt = tgt.next
		count = count + 1
	}

	mgr.fcp.nearmerge(names, path) // file merge
	return count
}

// header push (local -> remote)
func (mgr *g4IOmgr) hpush(data []byte, path string, isremote bool) error {
	if isremote && mgr.log.readonly {
		return errors.New("cannot hpush to readonly cluster")
	}
	if _, err := os.Stat(path); err == nil {
		os.Remove(path + ".bck")
		mgr.drv.write(path+".bck", 0, mgr.drv.read(path, 0, -1))
		mgr.log.write("msg : header backup complete")
	}
	mgr.drv.write(path, 0, data)
	mgr.log.write(fmt.Sprintf("msg : header writing complete -%s", path))
	return nil
}

// header seek (remote -> local), returns nil if not exists
func (mgr *g4IOmgr) hseek(path string) []byte {
	if _, err := os.Stat(path); err == nil {
		defer mgr.log.write(fmt.Sprintf("msg : header reading complete -%s", path))
		return mgr.drv.read(path, 0, -1)
	} else {
		mgr.log.write("err : header not exists")
		return nil
	}
}

// ===== section ===== kv4adv basic structure

// block A (account) KSC (4 section) [ KV4a, CRC32(section 0), CRC32(ClusterName) ]
// 0 : KDB text
//     -> (salt@B hint@B pwhash@B AccountName@S wrsign@B fsys_enckey@B fkey_enckey@B fphy_enckey@B)
// 1 : encrypted fsys (nB)
// 2 : encrypted fkey (53nB)
// 3 : encrypted (fptrnum 8B + fphykey 128B)

// block B (basic) KDB text
// -> (ClusterName@S chunksize@I wrsign@B blocknum@I)

// block A : salt hint pwhash account wrsign fsyskdt fkeykdt fphykdt
// block B : cluster chunksize wrsign blocknum

// basic module, caller has response of data manage (kv4adv version)
type g4bmod struct {
	bmod
	fptrcount int       // last used fptr
	fphymod   g4IOmgr   // fphy module (kv4adv version)
	fphykey   [128]byte // fphy cluster key (128B)
}

// find next unused fptr
func (tbox *g4bmod) find_inline() int {
	temp := tbox.fptrcount
	for tbox.fphymod.findfptr_inline(temp) != -1 {
		temp++
	}
	return temp
}

// random fill internal data ( [desktop, local, remote], chunksize, blocknum )
func (tbox *g4bmod) fill(paths []string, csize int, bnum int) {
	tbox.log.clear()
	tbox.log.abort = false
	tbox.log.working = false
	tbox.log.readonly = false
	tbox.drv.log = &tbox.log
	tbox.drv.sleep = 4
	tbox.fptrcount = 0

	tbox.fphykey = [128]byte(kaes.Genrand(128))
	tbox.fphymod.init(&tbox.drv, bnum, csize, paths[1], paths[2], tbox.fphykey)
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
	defer tbox.drv.dirctrl(paths[1]+"io/", false)
	tbox.enc.Before.Open(make([]byte, csize/47+43), true)
	tbox.enc.After.Open(paths[1]+"io/temp.bin", false)
	tbox.enc.Encrypt(bufkey)
	tbox.enc.Before.Close()
	tbox.enc.After.Close()

	var buf vfile
	buf.name = "_BUF"
	buf.time = root.time
	buf.size = kio.Size(paths[1] + "io/temp.bin")
	buf.fptr = tbox.fptrcount
	tbox.fkeymod.push(buf.fptr, [48]byte(bufkey))
	if tbox.fphymod.fpush(paths[1]+"io/temp.bin", buf.fptr) < 0 {
		tbox.log.write("critical : push fail -/_BUF")
		tbox.log.abort = true
	}

	root.subdir = []vdir{bin}
	root.subfile = []vfile{buf}
	tbox.fsysmod = root
}

// init modules ( [desktop, local, remote], [chunksize, blocknum, sleep], readonly, [fsys, fkey, fphy(8+128)] )
func (tbox *g4bmod) init(paths []string, parms []int, ronly bool, datas [][]byte) error {
	tbox.log.abort = false
	tbox.log.working = false
	tbox.log.readonly = ronly
	tbox.drv.log = &tbox.log
	tbox.drv.sleep = parms[2]

	if len(datas[2]) != 136 {
		return errors.New("invalid fphy")
	}
	tbox.fptrcount = kobj.Decode(datas[2][0:8])
	tbox.fphykey = [128]byte(datas[2][8:])
	tbox.fphymod.init(&tbox.drv, parms[1], parms[0], paths[1], paths[2], tbox.fphykey)
	err := tbox.fkeymod.read(datas[1])
	if err != nil {
		return err
	}
	tbox.fsysmod = *vread(datas[0])
	return nil
}

// file push (local -> remote), returns nil if error
func (tbox *g4bmod) fpush(data interface{}, name string) *vfile {
	key := [48]byte(kaes.Genrand(48))
	tbox.drv.dirctrl(tbox.fphymod.local+"io/", true)
	defer tbox.drv.dirctrl(tbox.fphymod.local+"io/", false)
	tbox.enc.Before.Open(data, true)
	tbox.enc.After.Open(tbox.fphymod.local+"io/temp.bin", false)
	err := tbox.ende(key, true)
	tbox.enc.Before.Close()
	tbox.enc.After.Close()
	if err != nil {
		tbox.log.write(fmt.Sprintf("err : encrypt fail -%s", err))
		return nil
	}

	num := tbox.find_inline()
	dirnum := tbox.fphymod.fpush(tbox.fphymod.local+"io/temp.bin", num)
	if dirnum < 0 {
		tbox.log.write("err : fpush fail")
		return nil
	}
	err = tbox.fkeymod.push(num, key)
	if err != nil {
		tbox.log.write(fmt.Sprintf("err : key push fail -%s", err))
		return nil
	}

	defer tbox.log.write(fmt.Sprintf("msg : vfile generated -fptr %d, pdir %d", num, dirnum))
	var newfile vfile
	newfile.name = name
	newfile.time = int(time.Now().Unix())
	newfile.size = kio.Size(tbox.fphymod.local + "io/temp.bin")
	newfile.fptr = num
	tbox.fptrcount = num
	return &newfile
}

// file pop (remote -> local)
func (tbox *g4bmod) fpop(fptr int) error {
	num := tbox.fphymod.fpop(fptr)
	if num < 0 {
		return errors.New("fpop fail")
	} else {
		err := tbox.fkeymod.pop(fptr)
		if err != nil {
			return err
		} else {
			if fptr < tbox.fptrcount {
				tbox.fptrcount = fptr
			}
			tbox.log.write(fmt.Sprintf("msg : remote file deleted -fptr %d, chunk %d", fptr, num))
			return nil
		}
	}
}

// file seek (remote -> local), empty string to bin output
func (tbox *g4bmod) fseek(fptr int, path string) ([]byte, error) {
	ext, _, key := tbox.fkeymod.seek(fptr)
	tbox.drv.dirctrl(tbox.fphymod.local+"io/", true)
	defer tbox.drv.dirctrl(tbox.fphymod.local+"io/", false)
	if ext < 0 {
		return nil, errors.New("no fkey")
	}
	num := tbox.fphymod.fseek(tbox.fphymod.local+"io/temp.bin", fptr)
	if num < 0 {
		return nil, errors.New("fseek fail")
	}

	tbox.enc.Before.Open(tbox.fphymod.local+"io/temp.bin", true)
	if path == "" {
		tbox.enc.After.Open(make([]byte, 0, 10485760), false)
	} else {
		tbox.enc.After.Open(path, false)
	}
	err := tbox.ende(key, false)
	tbox.enc.Before.Close()
	out := tbox.enc.After.Close()
	if err != nil {
		return nil, err
	}
	tbox.log.write(fmt.Sprintf("msg : file generated -fptr %d, chunk %d", fptr, num))
	return out, nil
}

// vdir info (for readonly)
type G4Finfo struct {
	Self_name    string
	Self_time    string
	Self_size    int // moreinfo
	Self_locked  bool
	Self_subdir  int // moreinfo
	Self_subfile int // moreinfo

	Dir_name   []string
	Dir_time   []string
	Dir_size   []int // moreinfo
	Dir_locked []bool

	File_name []string
	File_time []string
	File_size []int
	File_fptr []int
}

// init struct (not update moreinfo)
func (tbox *G4Finfo) init(tgt *vdir, wrlocked bool) {
	tbox.Self_name = tgt.name
	tbox.Self_time = time.Unix(int64(tgt.time), 0).Local().Format("2006.01.02;15:04:05")
	tbox.Self_size = -1
	tbox.Self_locked = tgt.islocked
	tbox.Self_subdir = -1
	tbox.Self_subdir = -1

	tbox.Dir_name = make([]string, 0)
	tbox.Dir_time = make([]string, 0)
	tbox.Dir_size = make([]int, 0)
	tbox.Dir_locked = make([]bool, 0)
	for _, r := range tgt.subdir {
		if wrlocked || !r.islocked {
			tbox.Dir_name = append(tbox.Dir_name, r.name)
			tbox.Dir_time = append(tbox.Dir_time, time.Unix(int64(r.time), 0).Local().Format("2006.01.02;15:04:05"))
			tbox.Dir_size = append(tbox.Dir_size, -1)
			tbox.Dir_locked = append(tbox.Dir_locked, r.islocked)
		}
	}

	tbox.File_name = make([]string, 0)
	tbox.File_time = make([]string, 0)
	tbox.File_size = make([]int, 0)
	tbox.File_fptr = make([]int, 0)
	if wrlocked || !tgt.islocked {
		for _, r := range tgt.subfile {
			tbox.File_name = append(tbox.File_name, r.name)
			tbox.File_time = append(tbox.File_time, time.Unix(int64(r.time), 0).Local().Format("2006.01.02;15:04:05"))
			tbox.File_size = append(tbox.File_size, r.size)
			tbox.File_fptr = append(tbox.File_fptr, r.fptr)
		}
	}
}

// ===== section ===== kv4adv meta functions

// generate new vault at remote path (~/), path should be empty folder
func G4FS_New(remote string, cluster string, divsize int, maxfptr int) error {
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

	var worker G4FScore // new G4FScore
	pw := []byte("0000")
	kf := kaes.Basickey()
	hint := []byte("new : PW 0000 KF bkf")
	os.Mkdir("./temp703a/", os.ModePerm)
	defer os.RemoveAll("./temp703a/")
	os.Mkdir("./temp703b/", os.ModePerm)
	defer os.RemoveAll("./temp703b/")

	// fill data
	os.Mkdir(remote+"0/", os.ModePerm)
	worker.data.fill([]string{"./temp703a/", "./temp703b/", remote}, []string{cluster, "root"}, divsize, [][]byte{pw, kf, hint})
	worker.data.blocknum = maxfptr
	worker.module.fill([]string{"./temp703a/", "./temp703b/", remote}, divsize, maxfptr)
	for i := 0; i < 3; i++ {
		temp := worker.data.mkey[48*i : 48*i+48]
		worker.data.keybuf[i+3] = [48]byte(worker.module.drv.aescalc(worker.data.keybuf[i][:], [48]byte(temp), true, false))
	}
	worker.Rootpath = worker.module.fsysmod.name
	worker.Cluster = &worker.data.cluster
	worker.Account = &worker.data.account
	worker.flush_save()
	return nil
}

// boot with cluster, init internal data & local/, returns (G4FScore, hint)
func G4FS_Boot(desktop string, local string, remote string, blockApath string) (*G4FScore, []byte, error) {
	// prebooting minimal data
	var out G4FScore
	var err error
	worker := ksc.Initksc()
	worker.Predetect = true
	out.module.log.readonly = true
	out.module.drv.log = &out.module.log
	out.module.fphymod.log = &out.module.log
	out.module.fphymod.drv = &out.module.drv

	// KSC unpacking block A
	worker.Path = blockApath
	err = worker.Readf()
	if err != nil {
		return nil, nil, err
	}
	if !kio.Bequal(worker.Subtype, []byte("KV4a")) {
		return nil, nil, errors.New("invalid block A")
	}
	blockA := out.module.drv.read(blockApath, worker.Chunkpos[0]+8, worker.Chunksize[0])
	out.databuf[3] = out.module.drv.read(blockApath, worker.Chunkpos[1]+8, worker.Chunksize[1])
	out.databuf[4] = out.module.drv.read(blockApath, worker.Chunkpos[2]+8, worker.Chunksize[2])
	out.databuf[5] = out.module.drv.read(blockApath, worker.Chunkpos[3]+8, worker.Chunksize[3])
	blockB := out.module.fphymod.hseek(remote + "0b.txt")
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
func G4FS_Exit(obj *G4FScore) {
	// readvar section clear
	defer os.RemoveAll(obj.data.local)
	obj.module.log.abort = true
	obj.module.log.readonly = true
	obj.module.log.working = false
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
	obj.module.fphymod.local = ""
	obj.module.fphymod.remote = ""
	obj.module.fphymod.fphykey = [128]byte{}
	obj.module.fphymod.maxfptr = -1
	obj.module.fphymod.divsize = -1
	obj.module.fphymod.curdir = -1
	obj.module.fphymod.table = nil
	obj.module.fphykey = [128]byte{}
	obj.module.fptrcount = -1
}

// regenerate fphykey, rename files
func G4FS_Rebuild(remote string, pw []byte, kf []byte) error {
	// login to G4FScore
	worker, hint, err := G4FS_Boot("./", "./temp703c/", remote, remote+"0a.webp")
	if err != nil {
		return err
	}
	err = worker.Login(pw, kf, 4)
	if err != nil {
		return err
	}
	defer G4FS_Exit(worker)
	if worker.module.log.readonly {
		return errors.New("readonly account")
	}

	// remote rewriter
	newkey := [128]byte(kaes.Genrand(128))
	writer := &worker.module.fphymod
	writer.fphykey = newkey
	worker.module.fphykey = newkey

	for i, r := range writer.table { // enc.kv4 -> n.kv4 -> enc.kv4
		oldname := writer.fcp.farsub(fmt.Sprintf("%s%d/", remote, i))
		newname := make([]string, len(oldname))
		for j, l := range oldname {
			temp, _ := kio.Bread(l[0:32])
			temp = writer.drv.aescalc(temp, r.key, false, false)
			newname[j] = kio.Bprint(writer.drv.aescalc(temp, writer.getkey_inline(i), true, false)) + ".kv4"
		}
		for j, l := range oldname {
			writer.fcp.farname(fmt.Sprintf("%s%d/", remote, i), l, fmt.Sprintf("%d.kv4", j))
		}
		for j, l := range newname {
			writer.fcp.farname(fmt.Sprintf("%s%d/", remote, i), fmt.Sprintf("%d.kv4", j), l)
		}
	}
	return worker.AccReset(pw, kf, hint)
}

// ===== section ===== kv4adv general functions

// kv4adv filesystem core
type G4FScore struct {
	data   sdata  // setting data
	module g4bmod // worker module

	hpath   string     // header block A path
	databuf [6][]byte  // data (nB) of (fsys fkey fphy), [0:3] plain [3:6] enc
	lock    sync.Mutex // fsys access lock (one at once)

	Curpath string // empty str if Curdir is nil
	Curdir  *vdir  // current folder

	Rootpath string  // path of root folder (*/) -readonly const
	Cluster  *string // cluster name -readonly const
	Account  *string // account name (root : RW, else : R) -readonly const
}

// generate block A/B header by current data/buffer, clear databuf
func (tbox *G4FScore) flush_head() ([]byte, []byte) {
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
	tbox.module.log.write("msg : encryption done -flush_head")

	blockA, blockB := tbox.data.wrhead() // write block A/B
	worker := ksc.Initksc()
	worker.Prehead = basewebp()
	worker.Subtype = []byte("KV4a")
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

	for i := 0; i < 6; i++ {
		tbox.databuf[i] = nil
	}
	return temp, []byte(blockB)
}

// save status to hpath, RW mode only, ! FsysLock !
func (tbox *G4FScore) flush_save() {
	// require : data, data.keybuf[0:6], module
	if !tbox.module.log.abort && !tbox.module.log.readonly {
		tbox.lock.Lock()
		tbox.databuf[0] = tbox.module.fsysmod.write(true) // fsys access point
		tbox.lock.Unlock()
		tbox.databuf[1] = tbox.module.fkeymod.write()
		tbox.databuf[2] = append(kobj.Encode(tbox.module.fptrcount, 8), tbox.module.fphykey[:]...)
		tbox.module.log.write("msg : data buffer filled -flush_all")

		blockA, blockB := tbox.flush_head()
		tbox.module.fphymod.hpush(blockA, tbox.data.remote+"0a.webp", true)
		tbox.module.fphymod.hpush(blockB, tbox.data.remote+"0b.txt", true)
		tbox.module.log.write("msg : header block A/B hpush done -flush_all")
	}
}

// save status to hpath, for AccReset, ! FsysLock !
func (tbox *G4FScore) flush_extend() {
	// require : data, data.keybuf[0:6], module
	if !tbox.module.log.abort {
		tbox.lock.Lock()
		tbox.databuf[0] = tbox.module.fsysmod.write(true) // fsys access point
		tbox.lock.Unlock()
		tbox.databuf[1] = tbox.module.fkeymod.write()
		tbox.databuf[2] = append(kobj.Encode(tbox.module.fptrcount, 8), tbox.module.fphykey[:]...)
		tbox.module.log.write("msg : data buffer filled -flush_extend")

		blockA, _ := tbox.flush_head()
		tbox.module.fphymod.hpush(blockA, tbox.hpath, false)
		tbox.module.log.write("msg : header block A hpush done -flush_extend")
	}
}

// check abort/working flag, check readonly if check_ro is true
func (tbox *G4FScore) check_auth(check_ro bool) error {
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

// check folder contains filename, returns inited subfile index, ! FsysLock !
func (tbox *G4FScore) check_import(domain *vdir, name string) int {
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
		tbox.module.fpop(domain.subfile[pos].fptr)
		domain.subfile[pos].time = int(time.Now().Unix())
		tbox.module.log.write(fmt.Sprintf("msg : fpush file exist -%s", name))
	}
	return pos
}

// check folder contains name, returns subdir/subfile index/-1, ! FsysLock !
func (tbox *G4FScore) check_export(domain *vdir, name string) int {
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

// find vdir by name, returns nil if not exists, ! fsys lock !
func (tbox *G4FScore) find_dir(path string, abspath bool) *vdir {
	tbox.lock.Lock()
	defer tbox.lock.Unlock()
	if len(path) == 0 || path[len(path)-1] != '/' {
		tbox.module.log.write(fmt.Sprintf("err : navigate fail -%s (abs %t)", path, abspath))
		return nil
	}
	if abspath && path == tbox.Rootpath {
		return &tbox.module.fsysmod
	}

	var tgt *vdir // fsys access point
	if abspath {  // find under rootdir
		tgt = tbox.module.getdir(&tbox.module.fsysmod, strings.Split(path[len(tbox.Rootpath):len(path)-1], "/"))
	} else { // find under curdir
		tgt = tbox.module.getdir(tbox.Curdir, strings.Split(path[0:len(path)-1], "/"))
	}
	if tgt == nil {
		tbox.module.log.write(fmt.Sprintf("err : navigate fail -%s (abs %t)", path, abspath))
	}
	return tgt
}

// manage /_BUF & make local/io/ & reset wrsign, works only at RW mode, ! FsysLock !
func (tbox *G4FScore) import_buffer() {
	if !tbox.module.log.abort && !tbox.module.log.readonly {
		tbox.data.wrsign = [8]byte(kaes.Genrand(8))
		tbox.module.drv.dirctrl(tbox.data.local+"io/", true)

		pos := -1
		for i, r := range tbox.module.fsysmod.subfile {
			if r.name == "_BUF" {
				pos = i
				break
			}
		}
		if pos == -1 {
			vf := tbox.module.fpush(make([]byte, tbox.data.chunksize/47+43), "_BUF")
			if vf == nil {
				tbox.module.log.write(fmt.Sprintf("err : buffer push fail -size %d", tbox.data.chunksize/47+43))
			} else {
				tbox.module.fsysmod.subfile = append(tbox.module.fsysmod.subfile, *vf)
			}
		} else {
			tbox.delete_file(&tbox.module.fsysmod, "_BUF")
		}
		tbox.module.fsysmod.sort()
	}
}

// generate vdir by real path(~/), using fkey/fphy but not fsys (non-blocking), ! FsysLock !
func (tbox *G4FScore) import_dir(path string) (*vdir, error) {
	var out vdir // gen new vdir
	out.name = path[strings.LastIndex(path[0:len(path)-1], "/")+1:]
	out.time = int(time.Now().Unix())
	out.islocked = false
	tbox.module.log.write(fmt.Sprintf("msg : vdir gen -%s", path))

	for _, r := range tbox.module.drv.dirsub(path) {
		if r[len(r)-1] == '/' { // push dir
			temp, err := tbox.import_dir(path + r)
			if err == nil {
				out.subdir = append(out.subdir, *temp)
				tbox.module.log.write(fmt.Sprintf("msg : vdir add -%s", r))
			} else {
				return nil, err
			}
		} else { // push file
			vf := tbox.module.fpush(path+r, r)
			if vf == nil {
				tbox.module.log.write(fmt.Sprintf("err : vdir add fail -%s", r))
			} else {
				out.subfile = append(out.subfile, *vf)
			}
		}
		if tbox.module.log.abort {
			return nil, errors.New("abort")
		}
	}

	out.sort() // sort & return
	return &out, nil
}

// generate folder/domain/, (folder : ~/), ! FsysLock !
func (tbox *G4FScore) export_dir(folder string, domain *vdir) {
	tbox.lock.Lock()
	vdir_name := domain.name
	if vdir_name == "/" {
		vdir_name = "_/"
	}
	tbox.lock.Unlock()
	tbox.module.log.write(fmt.Sprintf("msg : vdir export -%s", domain.name))

	tbox.module.drv.dirctrl(folder+vdir_name, true) // create folder
	for _, r := range domain.subfile {              // create direct subfile
		if tbox.module.log.abort {
			break
		}
		tbox.module.fseek(r.fptr, folder+vdir_name+r.name)
	}
	for _, r := range domain.subdir { // create direct subfolder
		if tbox.module.log.abort {
			break
		}
		tbox.export_dir(folder+vdir_name, &r)
	}
}

// delete file by name, works only at RW mode, ! FsysLock !
func (tbox *G4FScore) delete_file(domain *vdir, name string) {
	if !tbox.module.log.abort && !tbox.module.log.readonly {
		pos := tbox.check_export(domain, name)
		if pos == -1 {
			tbox.module.log.write(fmt.Sprintf("err : file delete fail -%s not exists", name))
		} else {
			err := tbox.module.fpop(domain.subfile[pos].fptr)
			if err != nil {
				tbox.module.log.write(fmt.Sprintf("err : file delete fail -%s", err))
			} else {
				domain.subfile = append(domain.subfile[:pos], domain.subfile[pos+1:]...)
				tbox.module.log.write(fmt.Sprintf("msg : file deleted -%s", name))
			}
		}
	}
}

// delete folder & sub, works only at RW mode, ! FsysLock !
func (tbox *G4FScore) delete_dir(folder *vdir) {
	for _, r := range folder.subdir {
		tbox.delete_dir(&r)
	}
	for _, r := range folder.subfile {
		if err := tbox.module.fpop(r.fptr); err != nil {
			tbox.module.log.write(fmt.Sprintf("err : file delete fail -%s", err))
		}
	}
}

// check fkey-fphy & returns new keymap
func (tbox *G4FScore) rebuild_fphy() (*keymap, int) {
	var out keymap
	out.read(nil)
	count := 0
	for _, r := range tbox.module.fphymod.table {
		for j := range r.table { // for all fptr in fphy
			if ext, _, key := tbox.module.fkeymod.seek(j); ext < 0 { // fkey seek fail
				tbox.module.fphymod.fpop(j)
				tbox.module.log.write(fmt.Sprintf("err : deleted while rebuild -fptr %d", j))
				count = count + 1
			} else { // fkey seek success
				out.push(j, key)
			}
		}
	}
	return &out, count
}

// check fkey-fsys, ! FsysLock !
func (tbox *G4FScore) rebuild_fsys(folder *vdir, fk *keymap) int {
	vfile_temp := make([]vfile, 0)
	count := 0
	tbox.lock.Lock()
	for _, r := range folder.subfile { // for all fptr in subfile
		if ext, _, key := tbox.module.fkeymod.seek(r.fptr); ext < 0 { // fkey seek fail
			tbox.module.log.write(fmt.Sprintf("err : deleted while rebuild -%s", r.name))
			count = count + 1
		} else { // fkey seek success
			vfile_temp = append(vfile_temp, r)
			fk.push(r.fptr, key)
		}
	}
	tbox.lock.Unlock()
	folder.subfile = vfile_temp
	for i := range folder.subdir {
		count = count + tbox.rebuild_fsys(&folder.subdir[i], fk)
	}
	return count
}

// reset abort/working flag & cache/cur dir, returns abort/working flag
func (tbox *G4FScore) Abort(reset bool, abort bool, working bool) (bool, bool) {
	if reset {
		tbox.module.log.write(fmt.Sprintf("msg : reset -abort %t working %t", abort, working))
		tbox.module.log.abort = abort
		tbox.module.log.working = working

		tbox.Curdir = &tbox.module.fsysmod
		tbox.Curpath = tbox.Curdir.name
		tbox.Rootpath = tbox.module.fsysmod.name
		tbox.Cluster = &tbox.data.cluster
		tbox.Account = &tbox.data.account
	}
	return tbox.module.log.abort, tbox.module.log.working
}

// debug info return : [divsize, maxfptr, last fptr], [wrsign(8B), salt(64B), pwhash(192B), fsyskey(48B), fkeykey(48B), fphykey(48B)]
func (tbox *G4FScore) Debug() ([]int, [][]byte) {
	out0 := make([]int, 3)
	out1 := make([][]byte, 6)
	out0[0] = tbox.data.chunksize
	out0[1] = tbox.data.blocknum
	out0[2] = tbox.module.fptrcount
	out1[0] = append(make([]byte, 0), tbox.data.wrsign[:]...)
	out1[1] = append(make([]byte, 0), tbox.data.salt[:]...)
	out1[2] = append(make([]byte, 0), tbox.data.pwhash[:]...)
	out1[3] = append(make([]byte, 0), tbox.data.keybuf[0][:]...)
	out1[4] = append(make([]byte, 0), tbox.data.keybuf[1][:]...)
	out1[5] = append(make([]byte, 0), tbox.data.keybuf[2][:]...)
	return out0, out1
}

// returns log data joined with \n, returns empty string if reset
func (tbox *G4FScore) Log(reset bool) string {
	if reset {
		tbox.module.log.clear()
		return ""
	} else {
		return tbox.module.log.read()
	}
}

// login & set module, header will be rewritten if root
func (tbox *G4FScore) Login(pw []byte, kf []byte, sleeptime int) error {
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
		tbox.module.fphymod.hpush(blockA, tbox.data.remote+"0a.webp", true)
		tbox.module.fphymod.hpush(blockB, tbox.data.remote+"0b.txt", true)
	}
	for i := 0; i < 6; i++ {
		tbox.databuf[i] = nil
	}
	tbox.data.mkey = [144]byte{}
	tbox.module.log.write("msg : cluster login success")
	return nil
}

// reset account PWKF, block A header will be rewritten
func (tbox *G4FScore) AccReset(pw []byte, kf []byte, hint []byte) error {
	if ferr := tbox.check_auth(false); ferr != nil {
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

	// save current status
	tbox.flush_extend()
	tbox.module.log.write("msg : reset complete -AccReset")
	return nil
}

// extend account (curdir becomes new rootdir), returns new block A header path at desktop
func (tbox *G4FScore) AccExtend(pw []byte, kf []byte, hint []byte, account string, wrlocked bool) (string, error) {
	if ferr := tbox.check_auth(false); ferr != nil {
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
			tbox.data.salt = [64]byte(tmp_salt)
			tbox.data.pwhash = [192]byte(tmp_pwhash)
			tbox.data.keybuf = tmp_keybuf
		}()
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
	tbox.module.log.write("msg : encryption done -AccExtend")

	ret := make(chan []byte, 1)
	go tbox.module.getkeys(tbox.Curdir, wrlocked, ret)
	newpath := fmt.Sprintf("%s%da.webp", tbox.data.desktop, kobj.Decode(kaes.Genrand(4))%900+100)
	// fill buffer (databuf[0:3])
	tbox.databuf[0] = tbox.Curdir.write(wrlocked)
	tbox.databuf[1] = <-ret
	tbox.databuf[2] = append(kobj.Encode(tbox.module.fptrcount, 8), tbox.module.fphykey[:]...)
	blockA, _ := tbox.flush_head()
	tbox.module.fphymod.hpush(blockA, newpath, false)
	for i := 0; i < 6; i++ {
		tbox.databuf[i] = nil
	}
	tbox.data.mkey = [144]byte{}
	tbox.module.log.write(fmt.Sprintf("msg : extend complete -%s", tbox.Curpath))
	return newpath, nil
}

// search name under Curdir, * : 0+ str, ? : len1 str, %d : int, %s 1+ ascii, %c len1 ascii, %* %? %% : literal
func (tbox *G4FScore) Search(name string) []string {
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
func (tbox *G4FScore) Print(wrlocked bool) string {
	tbox.lock.Lock()
	defer tbox.lock.Unlock()
	ret := make(chan string, 1)
	go tbox.module.print(tbox.Curdir, wrlocked, 0, ret)
	out := <-ret
	tbox.module.log.write(fmt.Sprintf("msg : print complete -%d", len(out)))
	return out
}

// set Curdir to input path (~/), Curdir will not change if path not exist, returns T if find success
func (tbox *G4FScore) Navigate(path string, abspath bool) bool {
	if path[len(path)-1] != '/' {
		tbox.module.log.write(fmt.Sprintf("err : navigate fail -%s (abs %t)", path, abspath))
		return false
	}
	res := tbox.find_dir(path, abspath)
	if res == nil {
		return false
	} else if abspath {
		tbox.Curpath = path
	} else {
		tbox.Curpath = tbox.Curpath + path
	}
	tbox.Curdir = res
	tbox.module.log.write(fmt.Sprintf("msg : moved to -%s", tbox.Curpath))
	return true
}

// get vdir info, returns nil if error, set path empty string to read curdir
func (tbox *G4FScore) Info(path string, wrlocked bool, moreinfo bool) *G4Finfo {
	var tgt *vdir
	if path == "" {
		tgt = tbox.Curdir
	} else {
		tgt = tbox.find_dir(path, true)
	}
	if tgt == nil {
		return nil
	}

	tbox.lock.Lock()
	defer tbox.lock.Unlock()
	var out G4Finfo
	out.init(tgt, wrlocked)
	if !moreinfo {
		return &out
	}

	out.Self_size = 0
	out.Self_subdir = 1
	out.Self_subfile = len(out.File_name)
	count := 0
	for _, r := range out.File_size {
		out.Self_size = out.Self_size + r
	}
	for _, r := range tgt.subdir {
		if wrlocked || !r.islocked {
			a, b, c := r.count(wrlocked)
			out.Dir_size[count] = a
			out.Self_size = out.Self_size + a
			out.Self_subdir = out.Self_subdir + b
			out.Self_subfile = out.Self_subfile + c
			count = count + 1
		}
	}
	return &out
}

// import binary data under Curdir(->tgtdir), file will be replaced if path exists
func (tbox *G4FScore) ImBin(name string, data []byte) error {
	if ferr := tbox.check_auth(true); ferr != nil {
		return ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	tbox.import_buffer() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	pos := tbox.check_import(tgtdir, name)
	vf := tbox.module.fpush(data, name)
	if vf == nil {
		tbox.module.log.write(fmt.Sprintf("err : import fail -%s", name))
		return errors.New("fphy fpush error")
	} else {
		tbox.lock.Lock()
		tgtdir.subfile[pos] = *vf
		tbox.lock.Unlock()
		tbox.module.log.write(fmt.Sprintf("msg : import binary -%s", name))
	}
	tgtdir.sort()
	tbox.flush_save()
	return nil
}

// import files under Curdir(->tgtdir), file will be replaced if path exists
func (tbox *G4FScore) ImFiles(paths []string) error {
	if ferr := tbox.check_auth(true); ferr != nil {
		return ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	tbox.import_buffer() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	for _, r := range paths {
		if tbox.module.log.abort {
			return errors.New("abort")
		}
		r = kio.Abs(r)
		if r[len(r)-1] == '/' {
			return errors.New("invalid file path")
		}
		name := r[strings.LastIndex(r, "/")+1:]

		pos := tbox.check_import(tgtdir, name)
		vf := tbox.module.fpush(r, name)
		if vf == nil {
			tbox.module.log.write(fmt.Sprintf("err : import fail -%s", name))
			return errors.New("fphy fpush error")
		} else {
			tbox.lock.Lock()
			tgtdir.subfile[pos] = *vf
			tbox.lock.Unlock()
			tbox.module.log.write(fmt.Sprintf("msg : import file -%s", name))
		}
	}
	tgtdir.sort()
	tbox.flush_save()
	return nil
}

// import folder under Curdir(->tgtdir), cannot import folder with same name
func (tbox *G4FScore) ImDir(path string) error {
	if ferr := tbox.check_auth(true); ferr != nil {
		return ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	tbox.import_buffer() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	path = kio.Abs(path) // checking name
	if path[len(path)-1] != '/' {
		return errors.New("invalid folder path")
	}
	name := path[strings.LastIndex(path[0:len(path)-1], "/")+1:]
	if tbox.check_export(tgtdir, name) != -1 {
		return errors.New("existing folder")
	}

	// generating temp vdir
	temp, err := tbox.import_dir(path)
	if err != nil {
		tbox.module.log.write(fmt.Sprintf("err : import fail -%s", err))
		return err
	}

	// sort & header push
	tbox.lock.Lock()
	tgtdir.subdir = append(tgtdir.subdir, *temp)
	tgtdir.sort()
	tbox.lock.Unlock()
	tbox.module.log.write(fmt.Sprintf("msg : import folder -%s", name))
	tbox.flush_save()
	return nil
}

// export binary data under Curdir(->tgtdir), find by name
func (tbox *G4FScore) ExBin(name string) ([]byte, error) {
	if ferr := tbox.check_auth(false); ferr != nil {
		return nil, ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	tbox.module.drv.dirctrl(tbox.data.local+"io/", true) // make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	pos := tbox.check_export(tgtdir, name)
	if pos < 0 {
		tbox.module.log.write(fmt.Sprintf("err : no file -%s", name))
		return nil, errors.New("no such file")
	}
	data, err := tbox.module.fseek(tgtdir.subfile[pos].fptr, "")
	if err != nil {
		tbox.module.log.write(fmt.Sprintf("err : export fail -%s", err))
		return nil, err
	} else {
		tbox.module.log.write(fmt.Sprintf("msg : export binary -%s", name))
		return data, nil
	}
}

// export files under Curdir(->tgtdir), find by name, !! generate desktop/kv5export/ !!
func (tbox *G4FScore) ExFiles(names []string) error {
	if ferr := tbox.check_auth(false); ferr != nil {
		return ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	tbox.module.drv.dirctrl(tbox.data.local+"io/", true) // make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	tbox.module.drv.dirctrl(tbox.data.desktop+"kv5export/", true) // make desktop/kv5export/
	for _, r := range names {
		if tbox.module.log.abort {
			return errors.New("abort")
		}
		pos := tbox.check_export(tgtdir, r)
		if pos < 0 {
			tbox.module.log.write(fmt.Sprintf("err : no file -%s", r))
			return errors.New("no such file")
		}

		if _, err := tbox.module.fseek(tgtdir.subfile[pos].fptr, tbox.data.desktop+"kv5export/"+r); err != nil {
			tbox.module.log.write(fmt.Sprintf("err : export fail -%s", err))
			return err
		} else {
			tbox.module.log.write(fmt.Sprintf("msg : export file -%s", r))
		}
	}
	return nil
}

// export name folder (empty string to export curdir), !! generate desktop/kv5export/ !!
func (tbox *G4FScore) ExDir(name string) error {
	if ferr := tbox.check_auth(false); ferr != nil {
		return ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	tbox.module.drv.dirctrl(tbox.data.local+"io/", true) // make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	tbox.module.drv.dirctrl(tbox.data.desktop+"kv5export/", true) // make desktop/kv5export/
	if name != "" {
		if name[len(name)-1] != '/' {
			return errors.New("invalid folder name")
		} else {
			pos := tbox.check_export(tgtdir, name)
			if pos < 0 {
				tbox.module.log.write(fmt.Sprintf("err : export fail -%s", name))
				return errors.New("no such folder")
			} else {
				tgtdir = &tgtdir.subdir[pos]
			}
		}
	}

	tbox.export_dir(tbox.data.desktop+"kv5export/", tgtdir)
	if tbox.module.log.abort {
		return errors.New("abort")
	} else {
		tbox.module.log.write(fmt.Sprintf("msg : export folder -%s", name))
		return nil
	}
}

// delete names, cannot delete / /_BIN/
func (tbox *G4FScore) Delete(names []string) error {
	if ferr := tbox.check_auth(true); ferr != nil {
		return ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	tgtpath := tbox.Curpath
	tbox.import_buffer() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	for _, r := range names {
		if tbox.module.log.abort {
			return errors.New("abort")
		}
		if tgtpath == "/" && r == "_BIN/" {
			return errors.New("del fail bindir")
		} else if tgtpath == "/" && r == "_BUF" {
			return errors.New("del fail buffer")
		} else if len(r) == 0 {
			return errors.New("del fail invalid name")
		}

		if r[len(r)-1] == '/' { // delete folder
			pos := tbox.check_export(tgtdir, r)
			if pos < 0 {
				return errors.New("no such folder")
			}
			tbox.module.log.write(fmt.Sprintf("msg : try to delete folder -%s", r))
			tbox.delete_dir(&tgtdir.subdir[pos])
			tbox.lock.Lock()
			tgtdir.subdir = append(tgtdir.subdir[:pos], tgtdir.subdir[pos+1:]...)
			tbox.lock.Unlock()
		} else { // delete file
			tbox.module.log.write(fmt.Sprintf("msg : try to delete file -%s", r))
			tbox.delete_file(tgtdir, r)
		}
	}
	tbox.flush_save()
	return nil
}

// move dir/files in Curdir(->tgtdir) to dst folder, return error if hierarchy problem is detected
func (tbox *G4FScore) Move(names []string, dst string) error {
	// cannot move /_BUF, /_BIN/ or overlapping names (except move to /_BIN/)
	if ferr := tbox.check_auth(true); ferr != nil {
		return ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	srcpath := tbox.Curpath
	tbox.import_buffer() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	dst_dir := tbox.find_dir(dst, true) // get dst folder
	if dst_dir == nil {
		return errors.New("invalid dst path")
	}

	// check hierarchy & name
	for _, r := range names {
		if tbox.data.contains(srcpath+r, dst) {
			return errors.New("invalid src path")
		}
		if tbox.check_export(tgtdir, r) == -1 {
			return errors.New("invalid src path")
		}
		if dst != "/_BIN/" && tbox.check_export(dst_dir, r) != -1 {
			return errors.New("overlapping name")
		}
		if r == "_BIN/" && srcpath == "/" {
			return errors.New("cannot move bindir")
		}
		if r == "_BUF" && srcpath == "/" {
			return errors.New("cannot move buffer")
		}
	}

	// move one by one
	for _, r := range names {
		pos := tbox.check_export(tgtdir, r)
		if r[len(r)-1] == '/' { // move tgt is folder
			dst_dir.subdir = append(dst_dir.subdir, tgtdir.subdir[pos])
			tgtdir.subdir = append(tgtdir.subdir[:pos], tgtdir.subdir[pos+1:]...)
		} else { // move tgt is file
			dst_dir.subfile = append(dst_dir.subfile, tgtdir.subfile[pos])
			tgtdir.subfile = append(tgtdir.subfile[:pos], tgtdir.subfile[pos+1:]...)
		}
	}

	// sort & header push
	dst_dir.sort()
	tbox.flush_save()
	tbox.module.log.write(fmt.Sprintf("msg : moved objects -%d", len(names)))
	return nil
}

// rename dir/files in Curdir(->tgtdir)
func (tbox *G4FScore) Rename(before []string, after []string) error {
	// cannot rename /_BUF, /_BIN/ or overlapping names
	if ferr := tbox.check_auth(true); ferr != nil {
		return ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	srcpath := tbox.Curpath
	tbox.import_buffer() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	// check name
	for i, r := range before {
		l := after[i]
		if tbox.check_export(tgtdir, r) == -1 {
			return errors.New("invalid before path")
		}
		if tbox.check_export(tgtdir, l) != -1 {
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
		if r == "_BIN/" && srcpath == "/" {
			return errors.New("cannot rename bindir")
		}
		if r == "_BUF" && srcpath == "/" {
			return errors.New("cannot rename buffer")
		}
	}

	// rename one by one
	for i, r := range before {
		l := after[i]
		pos := tbox.check_export(tgtdir, r)
		if r[len(r)-1] == '/' { // move tgt is folder
			tgtdir.subdir[pos].name = l
		} else { // move tgt is file
			tgtdir.subfile[pos].name = l
		}
	}

	// sort & header push
	tgtdir.sort()
	tbox.flush_save()
	tbox.module.log.write(fmt.Sprintf("msg : renamed objects -%d", len(before)))
	return nil
}

// generate new folder at Curdir(->tgtdir), cannot make empty or overlapping names
func (tbox *G4FScore) DirNew(names []string) error {
	if ferr := tbox.check_auth(true); ferr != nil {
		return ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	tbox.import_buffer() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	// check name
	for i, r := range names {
		if r[len(r)-1] != '/' {
			r = r + "/"
			names[i] = r
		}
		if tbox.check_export(tgtdir, r) != -1 {
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
		tgtdir.subdir = append(tgtdir.subdir, temp)
	}

	// sort & header push
	tgtdir.sort()
	tbox.flush_save()
	tbox.module.log.write(fmt.Sprintf("msg : new folders -%d", len(names)))
	return nil
}

// change lock status of name folder to islocked, changes all subdir if sub T, empty string to change curdir
func (tbox *G4FScore) DirLock(name string, islocked bool, sub bool) error {
	// cannot change lock status of /, /_BIN/
	if ferr := tbox.check_auth(true); ferr != nil {
		return ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	tgtdir := tbox.Curdir
	tbox.import_buffer() // /_BUF check, make local/io/
	defer os.RemoveAll(tbox.data.local + "io/")

	// check validity
	if name == "" && tbox.Curpath == "/" {
		return errors.New("cannot change rootdir")
	} else if name == "" && tbox.Curpath == "/_BIN/" {
		return errors.New("cannot change bindir")
	} else if tbox.Curpath == "/" && name == "_BIN/" {
		return errors.New("cannot change bindir")
	} else if name != "" {
		if name[len(name)-1] != '/' {
			return errors.New("invalid name")
		}
		if pos := tbox.check_export(tgtdir, name); pos < 0 {
			return errors.New("no such folder")
		} else {
			tgtdir = &tgtdir.subdir[pos]
		}
	}

	tbox.lock.Lock()
	if sub { // reset lock
		tbox.module.relock(tgtdir, islocked)
	} else {
		tgtdir.islocked = islocked
	}
	tbox.lock.Unlock()

	// header push
	tbox.flush_save()
	tbox.module.log.write(fmt.Sprintf("msg : set lock -%t %s", islocked, name))
	return nil
}

// mode T : check/fix name, F : check/fix file system
func (tbox *G4FScore) Restore(mode bool) (int, error) {
	// !! all works are about entire cluster !!, returns wrong file num
	if ferr := tbox.check_auth(true); ferr != nil {
		return 0, ferr
	}
	defer func() { tbox.module.log.working = false }()
	tbox.module.log.working = true
	defer os.RemoveAll(tbox.data.local + "io/")

	errnum := 0
	if mode { // rename
		tbox.module.log.write("msg : start rename -cluster restore")
		res := make(chan int, 1)
		go tbox.module.reset(&tbox.module.fsysmod, res)
		errnum = <-res

	} else { // rebuild
		var temp *keymap
		tbox.module.log.write("msg : start rebuild -cluster restore")
		temp, errnum = tbox.rebuild_fphy()
		tbox.module.fkeymod = *temp

		var newkey keymap
		newkey.read(nil)
		temp = &newkey
		errnum = errnum + tbox.rebuild_fsys(&tbox.module.fsysmod, temp)
		tbox.module.fkeymod = *temp

		tnum := 0
		temp, tnum = tbox.rebuild_fphy()
		tbox.module.fkeymod = *temp
		errnum = errnum + tnum
	}

	// header push
	tbox.flush_save()
	tbox.module.log.write(fmt.Sprintf("msg : cluster restore done -%d", errnum))
	return errnum, nil
}

// ===== section ===== kv4adv commander shell

// wrapper of G4FScore, never panics
type G4FSshell struct {
	// worker modules
	InSys    *G4FScore // kv4adv worker session
	AsyncErr string    // error result of async work
	IOstr    []string  // IO string buffer (manual access)
	IObyte   [][]byte  // IO []byte buffer (manual access)

	// readonly & path values
	FlagWk  bool     // isworking flag
	FlagRo  bool     // readonly cluster flag
	CurPath string   // full path of current session
	CurInfo *G4Finfo // curdir info

	moreinfo bool // moreinfo last option
}

// clear data & reset (not clear InSys)
func (sh *G4FSshell) init() {
	sh.AsyncErr = ""
	sh.IOstr = nil
	sh.IObyte = nil
	sh.FlagWk = false
	sh.FlagRo = true
	sh.CurPath = ""
	sh.CurInfo = nil
	sh.moreinfo = false
}

// update info of CurDir with wrlocked flag, TP to rootdir if CurDir not exists
func (sh *G4FSshell) update() {
	// init info & check session
	sh.CurInfo = nil
	if sh.CurPath != sh.InSys.Curpath {
		if !sh.InSys.Navigate(sh.CurPath, true) {
			sh.InSys.Navigate(sh.InSys.Rootpath, true)
			sh.CurPath = sh.InSys.Rootpath
		}
	}

	// get session info
	sh.CurInfo = sh.InSys.Info("", true, sh.moreinfo)
	if sh.CurInfo == nil {
		var temp0 vfile
		temp0.name = "critical error"
		temp0.time = int(time.Now().Unix())
		var temp1 vdir
		temp1.name = "critical error/"
		temp1.time = temp0.time
		temp1.subfile = []vfile{temp0}
		var temp2 G4Finfo
		temp2.init(&temp1, true)
		sh.CurInfo = &temp2
	}
}

// async cluster manipulate, error log will be AsyncErr
func (sh *G4FSshell) asyncwork(worktype int, parms []string) {
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
	case 11: // restore
		var ti int
		ti, err = sh.InSys.Restore(parms[0] == "true")
		sh.IOstr[0] = fmt.Sprint(ti)
	}

	if err != nil {
		sh.AsyncErr = fmt.Sprint(err)
	}
	sh.update()
}

// manipulate session with order/option, returns if started successfully
func (sh *G4FSshell) Command(order string, option []string) (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("shell fail : %s", ferr)
		}
		time.Sleep(time.Millisecond * 50)
	}()
	err = nil
	sh.AsyncErr = ""

	switch order {
	case "init": // init Shell
		sh.init()

	case "new": // G4FS_New -remote cluster divsize maxfptr
		if option[0][len(option[0])-1] != '/' {
			err = errors.New("invalid remote path")
		} else if len(option[1]) == 0 {
			err = errors.New("invalid cluster name")
		} else {
			divsize, _ := strconv.Atoi(option[2])
			maxfptr, _ := strconv.Atoi(option[3])
			if divsize == 0 {
				divsize = 104857600
			}
			if maxfptr == 0 {
				maxfptr = 256
			}
			err = G4FS_New(option[0], option[1], divsize, maxfptr)
		}

	case "boot": // G4FS_Boot -desktop local remote blockA
		if option[0][len(option[0])-1] != '/' {
			err = errors.New("invalid desktop path")
		} else if option[1][len(option[1])-1] != '/' {
			err = errors.New("invalid local path")
		} else if option[2][len(option[2])-1] != '/' {
			err = errors.New("invalid remote path")
		} else {
			sh.IObyte = make([][]byte, 1) // -hint
			sh.InSys, sh.IObyte[0], err = G4FS_Boot(option[0], option[1], option[2], option[3])
			sh.IOstr = make([]string, 2) // -cluster account
			sh.IOstr[0], sh.IOstr[1] = *sh.InSys.Cluster, *sh.InSys.Account
			sh.FlagRo = sh.InSys.module.log.readonly
		}

	case "exit": // G4FS_Exit
		G4FS_Exit(sh.InSys)
		sh.init()

	case "rebuild": // G4FS_Rebuild -remote
		defer func() { sh.IObyte = nil }() // -pw kf
		if option[0][len(option[0])-1] != '/' {
			err = errors.New("invalid remote path")
		} else {
			err = G4FS_Rebuild(option[0], sh.IObyte[0], sh.IObyte[1])
		}

	case "abort": // abort order/check -reset abort working
		_, sh.FlagWk = sh.InSys.Abort(option[0] == "true", option[1] == "true", option[2] == "true")

	case "debug": // get debug info string -count_locked
		sh.IOstr = make([]string, 4) // -debug *4
		ab, wr := sh.InSys.Abort(false, false, false)
		sh.IOstr[0] = "===== G4FS public data =====\n"
		sh.IOstr[0] = sh.IOstr[0] + fmt.Sprintf("Cluster : %s\nAccount : %s\n", *sh.InSys.Cluster, *sh.InSys.Account)
		sh.IOstr[0] = sh.IOstr[0] + fmt.Sprintf("RootPath : %s\nCurPath : %s\n", sh.InSys.Rootpath, sh.InSys.Curpath)
		sh.IOstr[0] = sh.IOstr[0] + fmt.Sprintf("Abort : %t\nIsWorking : %t\n", ab, wr)

		t0, t1 := sh.InSys.Debug() // [csize, bnum, fptr], [wrsign, salt, pwhash, fsyskey, fkeykey, fphykey]
		sh.IOstr[1] = "===== G4FS private data =====\n"
		sh.IOstr[1] = sh.IOstr[1] + fmt.Sprintf("divsize : %d\nmaxfptr : %d\nlast fptr : %d\n", t0[0], t0[1], t0[2])
		sh.IOstr[1] = sh.IOstr[1] + fmt.Sprintf("wrsign[8B] : %s\nsalt[64B] : %s\npwhash[192B] : %s\n", kio.Bprint(t1[0]), kio.Bprint(t1[1]), kio.Bprint(t1[2]))
		sh.IOstr[1] = sh.IOstr[1] + fmt.Sprintf("fsys key[48B] : %s\nfkey key[48B] : %s\nfphy key[48B] : %s\n", kio.Bprint(t1[3]), kio.Bprint(t1[4]), kio.Bprint(t1[5]))

		sh.IOstr[2] = "===== Shell public data =====\n"
		sh.IOstr[2] = sh.IOstr[2] + fmt.Sprintf("AsyncErr : %s\nlen IOstr : %d\nlen IObyte : %d\n", sh.AsyncErr, len(sh.IOstr), len(sh.IObyte))
		sh.IOstr[2] = sh.IOstr[2] + fmt.Sprintf("isworking : %t\nreadonly : %t\n", sh.FlagWk, sh.FlagRo)

		fs := sh.InSys.Info("", option[0] == "true", true)
		sh.IOstr[3] = "===== Shell curdir data =====\n"
		if fs == nil {
			sh.IOstr[3] = sh.IOstr[3] + "critical error\n"
		} else {
			sh.IOstr[3] = sh.IOstr[3] + fmt.Sprintf("path : %s\ntime : %s\n", sh.InSys.Curpath, fs.Self_time)
			sh.IOstr[3] = sh.IOstr[3] + fmt.Sprintf("size : %dB\nlock : %t\nfolder num : %d\nfile num : %d\n", fs.Self_size, fs.Self_locked, fs.Self_subdir, fs.Self_subfile)
			for i, r := range fs.Dir_name {
				sh.IOstr[3] = sh.IOstr[3] + fmt.Sprintf("%s (%s, %d B, lock %t)\n", r, fs.Dir_time[i], fs.Dir_size[i], fs.Dir_locked[i])
			}
			for i, r := range fs.File_name {
				sh.IOstr[3] = sh.IOstr[3] + fmt.Sprintf("%s (%s, %d B, fptr %d)\n", r, fs.File_time[i], fs.File_size[i], fs.File_fptr[i])
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

	case "update": // CurDir info update -moreinfo
		sh.moreinfo = option[0] == "true"
		sh.update()

	case "navigate": // get subdir names of tgt folder -fullpath
		sh.IOstr = make([]string, 1) // -folder names joined with \n
		temp := sh.InSys.Info(option[0], true, false)
		if temp == nil {
			err = errors.New("invalid path")
		} else {
			sh.IOstr[0] = strings.Join(temp.Dir_name, "\n")
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
			} else if slices.Contains(sh.CurInfo.Dir_name, temp[strings.LastIndex(temp[0:len(temp)-1], "/")+1:]) {
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
				if slices.Contains(sh.CurInfo.File_name, r) {
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
				if slices.Contains(sh.CurInfo.File_name, r) {
					temp = append(temp, r)
				}
			}
			go sh.asyncwork(4, temp)
		}

	case "exdir": // export folder -name(or empty)
		if sh.FlagWk {
			err = errors.New("invalid order")
		} else {
			go sh.asyncwork(5, option)
		}

	case "delete": // delete file/folder -name *n
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			temp := make([]string, 0)
			for _, r := range option {
				if sh.CurPath != "/" || (r != "_BIN/" && r != "_BUF") {
					temp = append(temp, r)
				}
			}
			go sh.asyncwork(6, temp)
		}

	case "move": // move file/folder to dst folder -tgtdir, name *n
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else if tgt := sh.InSys.Info(option[0], false, false); tgt == nil {
			err = errors.New("invalid path")
		} else {
			temp := []string{option[0]}
			for _, r := range option[1:] {
				if sh.CurPath != "/" || (r != "_BIN/" && r != "_BUF") {
					temp = append(temp, r)
				}
			}
			go sh.asyncwork(7, temp)
		}

	case "rename": // rename file/folder -name *n
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			temp0 := make([]string, 0)
			temp1 := make([]string, 0) // IOstr -newname *n
			for i, r := range option {
				if sh.CurPath != "/" || (r != "_BIN/" && r != "_BUF") {
					temp0 = append(temp0, r)
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
				if len(r) > 1 && !slices.Contains(sh.CurInfo.Dir_name, r) {
					temp = append(temp, r)
				}
			}
			go sh.asyncwork(9, temp)
		}

	case "dirlock": // relock folder -lock, name *n (empty string is CurDir)
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			temp := make([]string, 0) // (path lock sub) *n
			for _, r := range option[1:] {
				temp = append(append(temp, r), option[0])
				if r == "" {
					temp = append(temp, "true")
				} else {
					temp = append(temp, "false")
				}
			}
			go sh.asyncwork(10, temp)
		}

	case "restore": // cluster restore -mode
		sh.IOstr = []string{"0"} // IOstr -restored files
		if sh.FlagWk || sh.FlagRo {
			err = errors.New("invalid order")
		} else {
			switch option[0] {
			case "rename":
				go sh.asyncwork(11, []string{"true"})
			case "rebuild":
				go sh.asyncwork(11, []string{"false"})
			default:
				err = errors.New("invalid order")
			}
		}

	default: // unknown order & clear IObuf
		sh.IObyte = nil
		sh.IOstr = nil
		err = errors.New("unknown order")
	}

	return err
}
