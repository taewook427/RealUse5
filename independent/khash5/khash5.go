// test708 : independent.khash

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"stdlib5/cliget"
	"stdlib5/kdb"
	"stdlib5/kio"
	"stdlib5/kobj"
	"stdlib5/ksc"
	"stdlib5/ksign"
	"strconv"
	"strings"
	"time"
)

// get cliget.pathsel
func getpath() *cliget.PathSel {
	names := []string{"Local User", "Desktop"}
	paths := make([]string, 2)
	paths[0], _ = os.UserHomeDir()
	tp0 := filepath.Join(paths[0], "Desktop")
	tp1 := filepath.Join(paths[0], "desktop")
	tp2 := filepath.Join(paths[0], "바탕화면")

	var err error
	if _, err = os.Stat(tp0); err == nil {
		paths[1] = tp0
	} else if _, err = os.Stat(tp1); err == nil {
		paths[1] = tp1
	} else if _, err = os.Stat(tp2); err == nil {
		paths[1] = tp2
	} else {
		paths[1] = filepath.Join(paths[0], "DESKTOP")
	}

	var out cliget.PathSel
	out.Init(names, paths)
	return &out
}

// write data to path
func wrdata(path string, data string) {
	f, _ := kio.Open(path, "w")
	defer f.Close()
	kio.Write(f, []byte(data))
}

// public.txt : { name@S, date@S, public@S, phash@C }
// private.txt : { name@S, date@S, strength@I, private@S, public@S, phash@C }
// sign.txt : { name@S, strength@I, public@S, phash@C, explain@S, date@S, fhash@C, enc@C }

// khash file/folder
func func1() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< Select Path To Hash >", "target file", "target folder", "file mode"}, []string{"bool", "file", "folder", "bool"}, []string{"T", "*", "NR", "*"}, *getpath(), nil)
	sel.GetOpt()

	tgt := ""
	if sel.ByteRes[3][0] == 0 { // file hash
		tgt = sel.StrRes[1]
	} else { // folder hash
		tgt = sel.StrRes[2]
	}
	hv := ksign.Khash(tgt)
	fmt.Printf("\ntgt : %s\nhash : %s\n", tgt, kio.Bprint(hv))
	return err
}

// sign file/folder
func func2() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< Select Path & Key >", "target file", "target folder", "file mode", "private key"}, []string{"bool", "file", "folder", "bool", "keyfile"}, []string{"T", "*", "NR", "*", "*"}, *getpath(), nil)
	sel.GetOpt()

	tgt := ""
	if sel.ByteRes[3][0] == 0 { // file sign
		tgt = sel.StrRes[1]
	} else { // folder sign
		tgt = sel.StrRes[2]
	}
	worker := kdb.Initkdb()
	err = worker.Read(string(sel.ByteRes[4]))
	if err != nil {
		return err
	}

	tv, _ := worker.Get("name")
	pri_name := tv.Dat6
	tv, _ = worker.Get("strength")
	pri_strength := tv.Dat2
	tv, _ = worker.Get("private")
	pri_private := tv.Dat6
	tv, _ = worker.Get("public")
	pri_public := tv.Dat6
	tv, _ = worker.Get("phash")
	pri_phash := tv.Dat5
	if !kio.Bequal(pri_phash, ksc.Crc32hash([]byte(pri_public))) {
		return errors.New("private key error : phash not match")
	}

	sign_explain := kio.Input(fmt.Sprintf("Sign Info\nkey name : %s, strength : %d, phash : %s\nSign Explanation >>> ", pri_name, pri_strength, kio.Bprint(pri_phash)))
	sign_date := time.Now().Local().Format("2006.01.02;15:04:05")
	sign_fhash := ksign.Khash(tgt)
	var sign_enc []byte
	sign_enc, err = ksign.Sign(pri_private, sign_fhash)
	if err != nil {
		return err
	}

	worker = kdb.Initkdb()
	worker.Read("name = 0\nstrength = 0\npublic = 0\nphash = 0\nexplain = 0\ndate = 0\nfhash = 0\nenc = 0")
	worker.Fix("name", pri_name)
	worker.Fix("strength", pri_strength)
	worker.Fix("public", pri_public)
	worker.Fix("phash", pri_phash)
	worker.Fix("explain", sign_explain)
	worker.Fix("date", sign_date)
	worker.Fix("fhash", sign_fhash)
	worker.Fix("enc", sign_enc)
	data, _ := worker.Write()

	if tgt[len(tgt)-1] == '/' {
		tgt = tgt[0:strings.LastIndex(tgt[0:len(tgt)-1], "/")+1] + "sign.txt"
	} else {
		tgt = tgt[0:strings.LastIndex(tgt, "/")+1] + "sign.txt"
	}
	wrdata(tgt, data)
	fmt.Println("package sign complete!")
	return err
}

// verify sign of file/folder
func func3() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< Select Path & Sign >", "target file", "target folder", "file mode", "sign text"}, []string{"bool", "file", "folder", "bool", "file"}, []string{"T", "*", "NR", "*", "txt"}, *getpath(), nil)
	sel.GetOpt()

	tgt := ""
	if sel.ByteRes[3][0] == 0 { // file sign
		tgt = sel.StrRes[1]
	} else { // folder sign
		tgt = sel.StrRes[2]
	}
	f, _ := kio.Open(sel.StrRes[4], "r")
	data, _ := kio.Read(f, -1)
	f.Close()
	worker := kdb.Initkdb()
	err = worker.Read(string(data))
	if err != nil {
		return err
	}

	tv, _ := worker.Get("public")
	public := tv.Dat6
	tv, _ = worker.Get("phash")
	phash := tv.Dat5
	tv, _ = worker.Get("fhash")
	fhash := tv.Dat5
	tv, _ = worker.Get("enc")
	enc := tv.Dat5
	data = ksign.Khash(tgt)

	if !kio.Bequal(phash, ksc.Crc32hash([]byte(public))) {
		return errors.New("sign data error : phash not match")
	}
	if !kio.Bequal(fhash, data) {
		return errors.New("sign data error : fhash not match")
	}
	if s, e := ksign.Verify(public, enc, fhash); !s || e != nil {
		return errors.New("verify fail : invalid enc")
	}
	fmt.Println("verifying package complete!")
	return err
}

// generate new sign key
func func4() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< Set Signkey Field >", "signkey name", "strength"}, []string{"bool", "string", "int"}, []string{"T", "1+", "+"}, *getpath(), nil)
	sel.GetOpt()

	name := sel.StrRes[1]
	date := time.Now().Local().Format("2006.01.02;15:04:05")
	strength, _ := strconv.Atoi(sel.StrRes[2])
	if !slices.Contains([]int{1024, 2048, 4096, 8192, 16384}, strength) {
		return errors.New("invalid strength : should be 2**N shape (ex 2048)")
	}
	public, private, err := ksign.Genkey(strength)
	if err != nil {
		return err
	}
	phash := ksc.Crc32hash([]byte(public))

	worker := kdb.Initkdb()
	worker.Read("name = 0\ndate = 0\npublic = 0\nphash = 0")
	worker.Fix("name", name)
	worker.Fix("date", date)
	worker.Fix("public", public)
	worker.Fix("phash", phash)
	public_out, _ := worker.Write()

	worker = kdb.Initkdb()
	worker.Read("name = 0\ndate = 0\nstrength = 0\nprivate = 0\npublic = 0\nphash = 0")
	worker.Fix("name", name)
	worker.Fix("date", date)
	worker.Fix("strength", strength)
	worker.Fix("private", private)
	worker.Fix("public", public)
	worker.Fix("phash", phash)
	private_out, _ := worker.Write()

	wrdata("./public.txt", public_out)
	wrdata("./private.txt", private_out)
	fmt.Println("signkey generation complete!")
	return err
}

func main() {
	kobj.Repath()
	fmt.Println("===== RealUse5 independent KHASH5 =====")
	flag := true
	for flag {
		fmt.Printf("\n%18s   (0) %14s   (1) %14s\n(2) %14s   (3) %14s   (4) %14s\n", "< Mode Selection >", "Exit", "Get Hash Value", "Sign Package", "Verify Sign", "Generate Key")
		mode := kio.Input(">>> ")
		switch mode {
		case "0", "0\r", "0\n", "0\r\n":
			flag = false
		case "1":
			if err := func1(); err != nil {
				fmt.Println(err)
			}
		case "2":
			if err := func2(); err != nil {
				fmt.Println(err)
			}
		case "3":
			if err := func3(); err != nil {
				fmt.Println(err)
			}
		case "4":
			if err := func4(); err != nil {
				fmt.Println(err)
			}
		}
	}
	kio.Input("press ENTER to exit... ")
}
