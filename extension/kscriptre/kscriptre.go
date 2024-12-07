// test725 : extension.kscriptre

package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"stdlib5/kdb"
	"stdlib5/kio"
	"stdlib5/kobj"
	"stdlib5/ksc"
	"stdlib5/kscript"
	kscript_lib "stdlib5/kscript/lib"

	"strconv"
)

// runtime struct
type langre struct {
	option      [4]bool         // safemem, runone, errhlt, chksign
	sign_public []string        // trusted public key
	sign_phash  []string        // trusted key phash
	sign_name   []string        // trusted key name
	abi_name    []string        // abi name
	abi_code    []int           // abi code
	abi_using   []bool          // abi using
	runtime     kscript.KVM     // KVM runtime
	lib         kscript_lib.Lib // KVM lib
}

// init struct
func (tbox *langre) init() {
	tbox.option = [4]bool{true, false, true, true}
	tbox.sign_public = make([]string, 0)
	tbox.sign_phash = make([]string, 0)
	tbox.sign_name = make([]string, 0)
	tbox.abi_name = make([]string, 0)
	tbox.abi_code = make([]int, 0)
	tbox.abi_using = make([]bool, 0)
	tbox.runtime.Init()
	tbox.lib.Init(false, false, false, false)
}

// load program
func (tbox *langre) load(path string) error {
	info, abif, public, err := tbox.runtime.View(path)
	if err != nil {
		return err
	}
	fmt.Printf("[runtime load] info : %s, abi : %d\n", info, abif)
	if tbox.option[3] {
		if public == "" {
			return errors.New("no sign exists")
		} else {
			flag := true
			for i, r := range tbox.sign_public {
				if r == public {
					flag = false
					fmt.Printf("[runtime load] sign : %s (%s)\n", tbox.sign_name[i], tbox.sign_phash[i])
					break
				}
			}
			if flag {
				return fmt.Errorf("not trusted sign (%s)", kio.Bprint(ksc.Crc32hash([]byte(public))))
			}
		}
	}

	tbox.runtime.SafeMem = tbox.option[0]
	tbox.runtime.RunOne = tbox.option[1]
	tbox.runtime.ErrHlt = tbox.option[2]
	err = tbox.runtime.Load(tbox.option[3])
	if err != nil {
		return err
	}
	codes := frag(abif)
	for _, r := range codes {
		if !slices.Contains(tbox.abi_code, r) {
			return fmt.Errorf("unsupported abi : %d", r)
		}
	}

	f, _ := kio.Open("../../_ST5_CONFIG.txt", "r")
	d, _ := kio.Read(f, -1)
	f.Close()
	worker := kdb.Initkdb()
	worker.Read(string(d))
	tv, _ := worker.Get("dev.os")
	tbox.lib.Init(tbox.abi_using[slices.Index(tbox.abi_code, 2)], tbox.abi_using[slices.Index(tbox.abi_code, 4)], tbox.abi_using[slices.Index(tbox.abi_code, 8)], tv.Dat6 == "windows")
	tv, _ = worker.Get("path.desktop")
	tbox.lib.P_desktop = tv.Dat6
	tv, _ = worker.Get("path.local")
	tbox.lib.P_local = tv.Dat6
	tbox.lib.P_starter = kio.Abs("../../")
	tbox.lib.P_base = kio.Abs("./")
	fmt.Printf("[runtime load] success : %s\n", path)
	return nil
}

// run until exit code
func (tbox *langre) run() error {
	defer tbox.lib.Exit()
	fmt.Print("[runtime msg] program start\n\n")
	usetlib := tbox.abi_using[slices.Index(tbox.abi_code, 1)]
	for {
		intr := tbox.runtime.Run()
		if intr >= 32 { // kscript_lib
			vout, eout := tbox.lib.Run(tbox.runtime.CallMem, intr)
			if eout != nil {
				if tbox.runtime.ErrHlt {
					return eout
				} else {
					fmt.Printf("[runtime error] %s\n", eout)
				}
			}
			tbox.runtime.SetRet(vout)

		} else if intr >= 16 { // testfunc
			if usetlib {
				tbox.runtime.SetRet(kscript.TestIO(intr, tbox.runtime.CallMem))
			} else if tbox.runtime.ErrHlt {
				return errors.New("not supported func")
			} else {
				fmt.Println("[runtime error] not supported func")
			}

		} else if intr >= 0 { // kscript
			if intr == 1 {
				return nil
			}
			if intr == 2 {
				if tbox.runtime.ErrHlt {
					return errors.New(tbox.runtime.ErrMsg)
				} else {
					fmt.Println("[runtime error] " + tbox.runtime.ErrMsg)
				}
			}

		} else { // critical
			return fmt.Errorf("critical : %s", tbox.runtime.ErrMsg)
		}
	}
}

// divide abif
func frag(abif int) []int {
	codes := make([]int, 0)
	num := 1
	for abif > 0 {
		if abif&1 == 1 {
			codes = append(codes, num)
		}
		abif = abif >> 1
		num = num << 1
	}
	return codes
}

func main() {
	kobj.Repath()
	var rt langre
	rt.init()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[runtime error] critical : %s\n", err)
		}
		kio.Input("\n[runtime msg] program end : press ENTER to exit... ")
	}()

	f, _ := kio.Open("./settings.txt", "r")
	d, _ := kio.Read(f, -1)
	f.Close()
	worker := kdb.Initkdb()
	worker.Read(string(d))
	tv, _ := worker.Get("safemem")
	rt.option[0] = tv.Dat1
	tv, _ = worker.Get("runone")
	rt.option[1] = tv.Dat1
	tv, _ = worker.Get("errhlt")
	rt.option[2] = tv.Dat1
	tv, _ = worker.Get("chksign")
	rt.option[3] = tv.Dat1
	tv, _ = worker.Get("maxstk")
	if tv.Dat2 > 0 {
		rt.runtime.MaxStk = tv.Dat2
	}
	fmt.Printf("[runtime msg] safe memory : %t, run one op : %t, halt when error : %t, check sign : %t\n", rt.option[0], rt.option[1], rt.option[2], rt.option[3])
	num := 0
	for {
		if _, ext := worker.Name[fmt.Sprintf("%d.name", num)]; ext {
			num = num + 1
		} else {
			break
		}
	}
	fmt.Print("[runtime msg] ")
	for i := 0; i < num; i++ {
		tv, _ = worker.Get(fmt.Sprintf("%d.name", i))
		rt.abi_name = append(rt.abi_name, tv.Dat6)
		tv, _ = worker.Get(fmt.Sprintf("%d.code", i))
		rt.abi_code = append(rt.abi_code, tv.Dat2)
		tv, _ = worker.Get(fmt.Sprintf("%d.using", i))
		rt.abi_using = append(rt.abi_using, tv.Dat1)
		fmt.Printf("%s : %t, ", rt.abi_name[i], rt.abi_using[i])
	}
	fmt.Print("\n")

	f, _ = kio.Open("../../_ST5_SIGN.txt", "r")
	d, _ = kio.Read(f, -1)
	f.Close()
	worker = kdb.Initkdb()
	worker.Read(string(d))
	num = 0
	for {
		if _, ext := worker.Name[fmt.Sprintf("%d.name", num)]; ext {
			num = num + 1
		} else {
			break
		}
	}
	for i := 0; i < num; i++ {
		tv, _ = worker.Get(fmt.Sprintf("%d.name", i))
		rt.sign_name = append(rt.sign_name, tv.Dat6)
		tv, _ = worker.Get(fmt.Sprintf("%d.phash", i))
		rt.sign_phash = append(rt.sign_phash, kio.Bprint(tv.Dat5))
		tv, _ = worker.Get(fmt.Sprintf("%d.public", i))
		rt.sign_public = append(rt.sign_public, tv.Dat6)
	}

	pgrs := make([]string, 0)
	fs, _ := os.ReadDir("./_ST5_DATA/")
	for _, r := range fs {
		pgrs = append(pgrs, r.Name())
	}
	for i, r := range pgrs {
		fmt.Printf("[%03d] %s   ", i, r)
		if (i+1)%3 == 0 {
			fmt.Print("\n")
		}
	}
	num, err := strconv.Atoi(kio.Input("\nSelect program number : "))
	if err != nil {
		fmt.Printf("[runtime load] fail : %s\n", err)
		return
	}

	if err = rt.load("./_ST5_DATA/" + pgrs[num]); err != nil {
		fmt.Printf("[runtime load] fail : %s\n", err)
		return
	}
	if err = rt.run(); err != nil {
		fmt.Printf("[runtime error] %s\n", err)
		return
	}
}
