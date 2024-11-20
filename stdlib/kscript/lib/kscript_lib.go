// test706 : stdlib5.kscript stdlib/stdio/osfs

package kscript_lib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"stdlib5/kio"
	"stdlib5/kobj"
	"stdlib5/kscript"
	"stdlib5/ksign"
	"strconv"
	"strings"
	"time"
)

// 3-module worker, (stdlib, stdio, osfs), !!! manual set path !!!
type Lib struct {
	u_stdlib bool // lib using flag (abi 2, iter 32~100)
	u_stdio  bool // lib using flag (abi 4, iter 100~200)
	u_osfs   bool // lib using flag (abi 8, iter 200~300)

	myos string                       // working os name
	hmem map[string]([]kscript.Vunit) // global memory map (mem[id])
	fhnd map[string](*os.File)        // file handle (file[fullpath])

	P_desktop string // desktop path (~/)
	P_local   string // local user path (~/)
	P_starter string // st5 path (~/)
	P_base    string // program file position (~)
}

// stdlib func (32~100)
func (tbox *Lib) f_stdlib(parms []kscript.Vunit, icode int) (*kscript.Vunit, error) {
	if !tbox.u_stdlib {
		return nil, errors.New("not supported func")
	}
	var vout kscript.Vunit
	var err error
	switch {
	case icode < 50: // builtin
		switch icode {
		case 32: // type
			vout.Set([]string{"none", "bool", "int", "float", "str", "bytes"}[parms[0].Vtype])

		case 33: // int
			switch parms[0].Vtype {
			case 2:
				vout.Set(parms[0].Vint)
			case 3:
				vout.Set(int(parms[0].Vfloat))
			case 4:
				temp := 0
				temp, err = strconv.Atoi(parms[0].Vstring)
				vout.Set(temp)
			case 5:
				vout.Set(kobj.Decode(parms[0].Vbytes))
			default:
				return nil, errors.New("type error")
			}

		case 34: // float
			switch parms[0].Vtype {
			case 2:
				vout.Set(float64(parms[0].Vint))
			case 3:
				vout.Set(parms[0].Vfloat)
			case 4:
				temp := 0.0
				temp, err = strconv.ParseFloat(parms[0].Vstring, 64)
				vout.Set(temp)
			default:
				return nil, errors.New("type error")
			}

		case 35: // str
			if parms[0].Vtype == 5 {
				vout.Set(string(parms[0].Vbytes))
			} else {
				vout.Set(parms[0].ToString())
			}

		case 36: // bytes
			switch parms[0].Vtype {
			case 2:
				vout.Set(kobj.Encode(parms[0].Vint, 8))
			case 4:
				vout.Set([]byte(parms[0].Vstring))
			case 5:
				vout.Set(parms[0].Vbytes)
			default:
				return nil, errors.New("type error")
			}

		case 37: // hex
			switch parms[0].Vtype {
			case 2:
				vout.Set(strconv.FormatInt(int64(parms[0].Vint), 16))
			case 5:
				vout.Set(parms[0].ToString())
			default:
				return nil, errors.New("type error")
			}

		case 38: // chr
			if !tbox.check(parms, []string{"i"}) {
				return nil, errors.New("type error")
			}
			vout.Set(string(rune(parms[0].Vint)))

		case 39: // ord
			if !tbox.check(parms, []string{"s"}) {
				return nil, errors.New("type error")
			}
			vout.Set(int([]rune(parms[0].Vstring)[0]))

		case 40: // len
			switch parms[0].Vtype {
			case 4:
				vout.Set(len([]rune(parms[0].Vstring)))
			case 5:
				vout.Set(len(parms[0].Vbytes))
			default:
				return nil, errors.New("type error")
			}

		case 41: // slice
			if !tbox.check(parms, []string{"s|c", "n|i", "n|i"}) {
				return nil, errors.New("type error")
			}
			if parms[0].Vtype == 4 { // str[]
				tgt := []rune(parms[0].Vstring)
				start := 0
				end := len(tgt)
				if parms[1].Vtype == 2 {
					if parms[1].Vint < 0 {
						start = len(tgt) + parms[1].Vint
					} else {
						start = parms[1].Vint
					}
				}
				if parms[2].Vtype == 2 {
					if parms[2].Vint < 0 {
						end = len(tgt) + parms[2].Vint
					} else {
						end = parms[2].Vint
					}
				}
				vout.Set(string(tgt[start:end]))
			} else { // bytes[]
				tgt := parms[0].Vbytes
				start := 0
				end := len(tgt)
				if parms[1].Vtype == 2 {
					if parms[1].Vint < 0 {
						start = len(tgt) + parms[1].Vint
					} else {
						start = parms[1].Vint
					}
				}
				if parms[2].Vtype == 2 {
					if parms[2].Vint < 0 {
						end = len(tgt) + parms[2].Vint
					} else {
						end = parms[2].Vint
					}
				}
				vout.Set(tgt[start:end])
			}

		default:
			return nil, errors.New("not supported func")
		}

	case icode < 60: // memory
		switch icode {
		case 50: // m.malloc
			if !tbox.check(parms, []string{"s", "i"}) {
				return nil, errors.New("type error")
			}
			if parms[1].Vint < 0 {
				return nil, errors.New("invalid size")
			}
			if _, ext := tbox.hmem[parms[0].Vstring]; ext {
				err = errors.New("double alloc")
			}
			tbox.hmem[parms[0].Vstring] = make([]kscript.Vunit, parms[1].Vint)

		case 51: // m.realloc
			if !tbox.check(parms, []string{"s", "i"}) {
				return nil, errors.New("type error")
			}
			if parms[1].Vint < 0 {
				return nil, errors.New("invalid size")
			}
			var tgt []kscript.Vunit
			oldsize := 0
			if fs, ext := tbox.hmem[parms[0].Vstring]; ext {
				tgt = fs
				oldsize = len(fs)
			}
			newsize := parms[1].Vint
			if oldsize < newsize {
				tbox.hmem[parms[0].Vstring] = append(tgt, make([]kscript.Vunit, newsize-oldsize)...)
			} else {
				tbox.hmem[parms[0].Vstring] = tgt[0:newsize]
			}

		case 52: // m.free
			nm := parms[0].ToString()
			if _, ext := tbox.hmem[nm]; !ext {
				err = errors.New("double free")
			}
			delete(tbox.hmem, nm)

		case 53: // m.len
			nm := parms[0].ToString()
			if fs, ext := tbox.hmem[nm]; ext {
				vout.Set(len(fs))
			} else {
				vout.Set(-1)
			}

		case 54: // m.set
			if !tbox.check(parms, []string{"s", "i", "a"}) {
				return nil, errors.New("type error")
			}
			nm := parms[0].Vstring
			pos := parms[1].Vint
			if fs, ext := tbox.hmem[nm]; ext {
				if pos < 0 {
					pos = pos + len(fs)
				}
				fs[pos] = parms[2]
			} else {
				err = errors.New("invalid id")
			}

		case 55: // m.get
			if !tbox.check(parms, []string{"s", "i"}) {
				return nil, errors.New("type error")
			}
			nm := parms[0].Vstring
			pos := parms[1].Vint
			if fs, ext := tbox.hmem[nm]; ext {
				if pos < 0 {
					pos = pos + len(fs)
				}
				vout = fs[pos]
			} else {
				err = errors.New("invalid id")
			}

		case 56: // m.split
			if !tbox.check(parms, []string{"s", "s|c", "s|c"}) || parms[1].Vtype != parms[2].Vtype {
				return nil, errors.New("type error")
			}
			var temp []kscript.Vunit
			if parms[1].Vtype == 4 { // split str
				tgt := strings.Split(parms[1].Vstring, parms[2].Vstring)
				temp = make([]kscript.Vunit, len(tgt))
				for i, r := range tgt {
					temp[i].Set(r)
				}
			} else { // split bytes
				tgt := bytes.Split(parms[1].Vbytes, parms[2].Vbytes)
				temp = make([]kscript.Vunit, len(tgt))
				for i, r := range tgt {
					temp[i].Set(r)
				}
			}
			tbox.hmem[parms[0].Vstring] = temp

		case 57: // m.join
			if !tbox.check(parms, []string{"s", "s|c"}) {
				return nil, errors.New("type error")
			}
			nm := parms[0].Vstring
			if fs, ext := tbox.hmem[nm]; ext {
				if parms[1].Vtype == 4 { // join str
					temp := make([]string, 0)
					for _, r := range fs {
						if r.Vtype == 4 {
							temp = append(temp, r.Vstring)
						}
					}
					vout.Set(strings.Join(temp, parms[1].Vstring))
				} else { // join bytes
					temp := make([][]byte, 0)
					for _, r := range fs {
						if r.Vtype == 5 {
							temp = append(temp, r.Vbytes)
						}
					}
					vout.Set(bytes.Join(temp, parms[1].Vbytes))
				}
			} else {
				err = errors.New("invalid id")
			}

		default:
			return nil, errors.New("not supported func")
		}

	case icode < 70: // string
		switch icode {
		case 60: // str.change
			if !tbox.check(parms, []string{"s", "b"}) {
				return nil, errors.New("type error")
			}
			if parms[1].Vbool { // to upper
				vout.Set(strings.ToUpper(parms[0].Vstring))
			} else { // to lower
				vout.Set(strings.ToLower(parms[0].Vstring))
			}

		case 61: // str.find
			if !tbox.check(parms, []string{"s", "s", "b"}) {
				return nil, errors.New("type error")
			}
			pos := -1
			if parms[2].Vbool { // from head
				pos = strings.Index(parms[0].Vstring, parms[1].Vstring)
			} else { // from tail
				pos = strings.LastIndex(parms[0].Vstring, parms[1].Vstring)
			}
			if pos < 0 {
				vout.Set(-1)
			} else {
				vout.Set(len([]rune(parms[0].Vstring[0:pos])))
			}

		case 62: // str.count
			if !tbox.check(parms, []string{"s", "s"}) {
				return nil, errors.New("type error")
			}
			vout.Set(strings.Count(parms[0].Vstring, parms[1].Vstring))

		case 63: // str.replace
			if !tbox.check(parms, []string{"s", "s", "s", "i"}) {
				return nil, errors.New("type error")
			}
			vout.Set(strings.Replace(parms[0].Vstring, parms[1].Vstring, parms[2].Vstring, parms[3].Vint))

		default:
			return nil, errors.New("not supported func")
		}

	case icode < 80: // time
		switch icode {
		case 70: // t.time
			vout.Set(float64(time.Now().UnixMicro()) / 1000000)

		case 71: // t.stamp
			ti := 0
			if parms[0].Vtype == 2 {
				ti = parms[0].Vint
			} else if parms[0].Vtype == 3 {
				ti = int(parms[0].Vfloat)
			} else {
				return nil, errors.New("type error")
			}
			vout.Set(time.Unix(int64(ti), 0).Local().Format("2006.01.02;15:04:05"))

		case 72: // t.stampf
			if !tbox.check(parms, []string{"i|f", "s"}) {
				return nil, errors.New("type error")
			}
			ti := 0
			if parms[0].Vtype == 2 {
				ti = parms[0].Vint
			} else if parms[0].Vtype == 3 {
				ti = int(parms[0].Vfloat)
			} else {
				return nil, errors.New("type error")
			}
			ts := parms[1].Vstring
			ts = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(ts, "%Y", "2006"), "%M", "01"), "%D", "02")
			ts = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(ts, "%h", "15"), "%m", "04"), "%s", "05")
			vout.Set(time.Unix(int64(ti), 0).Local().Format(ts))

		case 73: // t.sleep
			if parms[0].Vtype == 2 {
				time.Sleep(time.Second * time.Duration(parms[0].Vint))
			} else if parms[0].Vtype == 3 {
				time.Sleep(time.Microsecond * time.Duration(parms[0].Vfloat*1000000))
			} else {
				return nil, errors.New("type error")
			}

		default:
			return nil, errors.New("not supported func")
		}

	case icode < 90: // math
		switch icode {
		case 80: // math.const
			if !tbox.check(parms, []string{"s"}) {
				return nil, errors.New("type error")
			}
			switch parms[0].Vstring {
			case "e":
				vout.Set(math.E)
			case "pi":
				vout.Set(math.Pi)
			case "phi":
				vout.Set(math.Phi)
			case "sqrt2":
				vout.Set(math.Sqrt2)
			default:
				err = errors.New("invalid option")
			}

		case 81: // math.conv
			if !tbox.check(parms, []string{"s", "i|f"}) {
				return nil, errors.New("type error")
			}
			switch parms[0].Vstring {
			case "abs":
				if parms[1].Vtype == 2 {
					vout.Set(int(math.Abs(float64(parms[1].Vint))))
				} else {
					vout.Set(math.Abs(parms[1].Vfloat))
				}
			case "up":
				if parms[1].Vtype == 2 {
					vout.Set(parms[1].Vint)
				} else {
					vout.Set(math.Ceil(parms[1].Vfloat))
				}
			case "down":
				if parms[1].Vtype == 2 {
					vout.Set(parms[1].Vint)
				} else {
					vout.Set(math.Floor(parms[1].Vfloat))
				}
			case "round":
				if parms[1].Vtype == 2 {
					vout.Set(parms[1].Vint)
				} else {
					vout.Set(math.Round(parms[1].Vfloat))
				}
			default:
				err = errors.New("invalid option")
			}

		case 82: // math.log
			if !tbox.check(parms, []string{"i|f", "i|f"}) {
				return nil, errors.New("type error")
			}
			base := 0.0
			num := 0.0
			if parms[0].Vtype == 2 {
				base = float64(parms[0].Vint)
			} else {
				base = parms[0].Vfloat
			}
			if parms[1].Vtype == 2 {
				num = float64(parms[1].Vint)
			} else {
				num = parms[1].Vfloat
			}
			vout.Set(math.Log(num) / math.Log(base))

		case 83: // math.trif
			if !tbox.check(parms, []string{"s", "i|f"}) {
				return nil, errors.New("type error")
			}
			num := 0.0
			if parms[1].Vtype == 2 {
				num = float64(parms[1].Vint)
			} else {
				num = parms[1].Vfloat
			}
			switch parms[0].Vstring {
			case "sin":
				vout.Set(math.Sin(num))
			case "cos":
				vout.Set(math.Cos(num))
			case "tan":
				vout.Set(math.Tan(num))
			case "asin":
				vout.Set(math.Asin(num))
			case "acos":
				vout.Set(math.Acos(num))
			case "atan":
				vout.Set(math.Atan(num))
			default:
				err = errors.New("invalid option")
			}

		case 84: // math.random
			vout.Set(rand.Float64())

		case 85: // math.randrange
			if !tbox.check(parms, []string{"i", "i"}) {
				return nil, errors.New("type error")
			}
			st := parms[0].Vint
			ed := parms[1].Vint
			if st < ed {
				vout.Set(rand.Intn(ed-st) + st)
			} else {
				err = errors.New("invalid option")
			}

		default:
			return nil, errors.New("not supported func")
		}

	default:
		return nil, errors.New("not supported func")
	}
	return &vout, err
}

// stdio func (100~200)
func (tbox *Lib) f_stdio(parms []kscript.Vunit, icode int) (*kscript.Vunit, error) {
	if !tbox.u_stdio {
		return nil, errors.New("not supported func")
	}
	var vout kscript.Vunit
	var err error
	switch icode {
	case 100: // io.input
		vout.Set(kio.Input(parms[0].ToString()))

	case 101: // io.print
		if !tbox.check(parms, []string{"a", "s"}) {
			return nil, errors.New("type error")
		}
		fmt.Print(parms[0].ToString() + parms[1].Vstring)

	case 102: // io.println
		fmt.Println(parms[0].ToString())

	case 103: // io.error
		if !tbox.check(parms, []string{"a", "s"}) {
			return nil, errors.New("type error")
		}
		fmt.Fprint(os.Stderr, parms[0].ToString()+parms[1].Vstring)

	case 104: // io.open
		if !tbox.check(parms, []string{"s", "s"}) {
			return nil, errors.New("type error")
		}
		path, _ := filepath.Abs(parms[0].Vstring)
		if _, ext := tbox.fhnd[path]; ext {
			tbox.fhnd[path].Close()
			tbox.fhnd[path], err = kio.Open(path, parms[1].Vstring)
		} else {
			tbox.fhnd[path], err = kio.Open(path, parms[1].Vstring)
		}

	case 105: // io.close
		if !tbox.check(parms, []string{"s"}) {
			return nil, errors.New("type error")
		}
		path, _ := filepath.Abs(parms[0].Vstring)
		if _, ext := tbox.fhnd[path]; ext {
			tbox.fhnd[path].Close()
			delete(tbox.fhnd, path)
		}

	case 106: // io.seek
		if !tbox.check(parms, []string{"s", "i", "i"}) {
			return nil, errors.New("type error")
		}
		path, _ := filepath.Abs(parms[0].Vstring)
		if _, ext := tbox.fhnd[path]; ext {
			var off int64
			if parms[2].Vint == 0 {
				off, err = tbox.fhnd[path].Seek(int64(parms[1].Vint), 0)
			} else if parms[2].Vint == -1 {
				off, err = tbox.fhnd[path].Seek(int64(parms[1].Vint), 2)
			} else {
				off, err = tbox.fhnd[path].Seek(int64(parms[1].Vint), 1)
			}
			vout.Set(int(off))
		} else {
			return nil, errors.New("handle not exists")
		}

	case 107: // io.readline
		if !tbox.check(parms, []string{"s", "i"}) {
			return nil, errors.New("type error")
		}
		path, _ := filepath.Abs(parms[0].Vstring)
		data := ""
		count := parms[1].Vint
		if f, ext := tbox.fhnd[path]; ext {
			pos, _ := f.Seek(0, 1)
			reader := bufio.NewReader(f)
			for count != 0 {
				ts, er := reader.ReadString('\n')
				data = data + ts
				count = count - 1
				if er != nil {
					break
				}
			}
			f.Seek(pos+int64(len(data)), 0)
			vout.Set(data)
		} else {
			return nil, errors.New("handle not exists")
		}

	case 108: // io.read
		if !tbox.check(parms, []string{"s", "i"}) {
			return nil, errors.New("type error")
		}
		path, _ := filepath.Abs(parms[0].Vstring)
		var data []byte
		if f, ext := tbox.fhnd[path]; ext {
			data, err = kio.Read(f, parms[1].Vint)
			vout.Set(data)
		} else {
			return nil, errors.New("handle not exists")
		}

	case 109: // io.write
		if !tbox.check(parms, []string{"s", "s|c"}) {
			return nil, errors.New("type error")
		}
		path, _ := filepath.Abs(parms[0].Vstring)
		if f, ext := tbox.fhnd[path]; ext {
			if parms[1].Vtype == 4 {
				_, err = kio.Write(f, []byte(parms[1].Vstring))
			} else {
				_, err = kio.Write(f, parms[1].Vbytes)
			}
		} else {
			return nil, errors.New("handle not exists")
		}

	default:
		return nil, errors.New("not supported func")
	}
	return &vout, err
}

// osfs func (200~300)
func (tbox *Lib) f_osfs(parms []kscript.Vunit, icode int) (*kscript.Vunit, error) {
	if !tbox.u_osfs {
		return nil, errors.New("not supported func")
	}
	var vout kscript.Vunit
	var err error
	switch icode {
	case 200: // os.name
		vout.Set(tbox.myos)

	case 201: // os.chdir
		if !tbox.check(parms, []string{"s"}) {
			return nil, errors.New("type error")
		}
		err = os.Chdir(parms[0].Vstring)

	case 202: // os.getpath
		path := ""
		switch parms[0].ToString() {
		case "cwd":
			path, err = os.Getwd()
			path = strings.Replace(path, "\\", "/", -1)
			if path[len(path)-1] != '/' {
				path = path + "/"
			}
		case "desktop":
			path = tbox.P_desktop
		case "local":
			path = tbox.P_local
		case "starter":
			path = tbox.P_starter
		case "base":
			path = tbox.P_base
		default:
			return nil, errors.New("invalid option")
		}
		vout.Set(path)

	case 203: // os.exists
		if _, er := os.Stat(parms[0].ToString()); er == nil {
			vout.Set(true)
		} else {
			vout.Set(false)
		}

	case 204: // os.abspath
		path, _ := filepath.Abs(parms[0].ToString())
		if _, er := os.Stat(path); er == nil {
			path = kio.Abs(path)
		}
		vout.Set(path)

	case 205: // os.is
		if !tbox.check(parms, []string{"s", "b"}) {
			return nil, errors.New("type error")
		}
		fs, er := os.Stat(parms[0].Vstring)
		if er != nil {
			vout.Set(false)
		} else if parms[1].Vbool { // check if file
			vout.Set(!fs.IsDir())
		} else { // check if folder
			vout.Set(fs.IsDir())
		}

	case 206: // os.finfo
		if !tbox.check(parms, []string{"s", "b"}) {
			return nil, errors.New("type error")
		}
		fs, er := os.Stat(parms[0].Vstring)
		if er != nil {
			vout.Set(-1)
		} else if parms[1].Vbool { // check size
			sz, _, _ := ksign.Kinfo(parms[0].Vstring)
			vout.Set(sz)
		} else { // check time
			vout.Set(int(fs.ModTime().Unix()))
		}

	case 207: // os.listdir
		if !tbox.check(parms, []string{"s"}) {
			return nil, errors.New("type error")
		}
		if _, er := os.Stat(parms[0].Vstring); er != nil {
			return nil, errors.New("invalid path")
		}
		temp := make([]string, 0)
		tgt, _ := os.ReadDir(parms[0].Vstring)
		for _, r := range tgt {
			nm := r.Name()
			if r.IsDir() && nm[len(nm)-1] != '/' {
				nm = nm + "/"
			}
			temp = append(temp, nm)
		}
		vout.Set(strings.Join(temp, "\n"))

	case 208: // os.mkdir
		err = os.Mkdir(parms[0].ToString(), os.ModePerm)

	case 209: // os.rename
		err = os.Rename(parms[0].ToString(), parms[1].ToString())

	case 210: // os.move
		err = os.Rename(parms[0].ToString(), parms[1].ToString())

	case 211: // os.remove
		err = os.RemoveAll(parms[0].ToString())

	default:
		return nil, errors.New("not supported func")
	}
	return &vout, err
}

// check parms types (nbifsc | a)
func (tbox *Lib) check(parms []kscript.Vunit, tps []string) bool {
	if len(parms) != len(tps) {
		return false
	}
	for i, r := range parms {
		flag := false
		for _, l := range strings.Split(tps[i], "|") {
			switch {
			case l == "a":
				flag = true
			case l == "n" && r.Vtype == 0:
				flag = true
			case l == "b" && r.Vtype == 1:
				flag = true
			case l == "i" && r.Vtype == 2:
				flag = true
			case l == "f" && r.Vtype == 3:
				flag = true
			case l == "s" && r.Vtype == 4:
				flag = true
			case l == "c" && r.Vtype == 5:
				flag = true
			}
			if flag {
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

// set module use, init struct, ! set paths after init !
func (tbox *Lib) Init(u_stdlib bool, u_stdio bool, u_osfs bool, iswin bool) {
	tbox.u_stdlib = u_stdlib
	tbox.u_stdio = u_stdio
	tbox.u_osfs = u_osfs
	if iswin {
		tbox.myos = "windows"
	} else {
		tbox.myos = "linux"
	}
	tbox.hmem = make(map[string][]kscript.Vunit)
	tbox.fhnd = make(map[string]*os.File)
	tbox.P_desktop = "!undefined_desktop"
	tbox.P_local = "!undefined_local"
	tbox.P_starter = "!undefined_starter"
	tbox.P_base = "!undefined_base"
}

// close all file handles
func (tbox *Lib) Exit() {
	for _, r := range tbox.fhnd {
		r.Close()
	}
	tbox.Init(false, false, false, false)
}

// outer func works
func (tbox *Lib) Run(parms []kscript.Vunit, icode int) (vout *kscript.Vunit, eout error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			eout = fmt.Errorf("critical : %s", ferr)
		}
		if vout == nil {
			var temp kscript.Vunit
			vout = &temp
		}
	}()
	if icode < 32 { // not supported
		return nil, errors.New("not supported func")
	} else if icode < 100 { // stdlib
		vout, eout = tbox.f_stdlib(parms, icode)
	} else if icode < 200 { // stdio
		vout, eout = tbox.f_stdio(parms, icode)
	} else if icode < 300 { // osfs
		vout, eout = tbox.f_osfs(parms, icode)
	} else { // not supported
		return nil, errors.New("not supported func")
	}
	return vout, eout
}
