// test742 : extension.mdm5

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"stdlib5/cliget"
	"stdlib5/kaes"
	"stdlib5/kdb"
	"stdlib5/kio"
	"stdlib5/kobj"
	"stdlib5/kpic"
	"stdlib5/kpkg"
	"stdlib5/ksc"
	"stdlib5/kvault"
	"stdlib5/kzip"
	"stdlib5/legsup"
	"stdlib5/picdt"
	"strconv"
	"strings"
)

// "go run ./" "go build -ldflags="-s -w" -trimpath ./"

// filediv pack
func func1() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< FileDiv Pack >", "target file", "div size (B)"}, []string{"bool", "file", "int"}, []string{"T", "*", "+"}, *getpath(), nil)
	sel.StrRes[2] = "25165824"
	sel.GetOpt()

	path := sel.StrRes[1]
	name := path[strings.LastIndex(path, "/")+1:]
	div, _ := strconv.Atoi(sel.StrRes[2])
	size := kio.Size(path)
	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	num := 0
	var d []byte
	for size > 0 {
		t, err := kio.Open(fmt.Sprintf("./_ST5_DATA/%s.%d", name, num), "w")
		if err == nil {
			if size > div {
				d, _ = kio.Read(f, div)
			} else {
				d, _ = kio.Read(f, size)
			}
			kio.Write(t, d)
			t.Close()
		} else {
			return err
		}
		size = size - div
		num = num + 1
	}
	fmt.Printf("FileDiv Pack Success : %d files\n", num)
	return err
}

// filediv unpack
func func2() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< FileDiv Unpack >", "starting file"}, []string{"bool", "file"}, []string{"T", "0"}, *getpath(), nil)
	sel.GetOpt()

	path := sel.StrRes[1][:strings.LastIndex(sel.StrRes[1], ".")]
	name := path[strings.LastIndex(path, "/")+1:]
	num := 0
	for {
		if kio.Size(fmt.Sprintf("%s.%d", path, num)) == -1 {
			break
		}
		num = num + 1
	}
	f, err := kio.Open(fmt.Sprintf("./_ST5_DATA/%s", name), "w")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	for i := 0; i < num; i++ {
		t, err := kio.Open(fmt.Sprintf("%s.%d", path, i), "r")
		if err == nil {
			d, _ := kio.Read(t, -1)
			kio.Write(f, d)
			t.Close()
		} else {
			return err
		}
	}
	fmt.Printf("FileDiv Unpack Success : %d files\n", num)
	return err
}

// g5ksc view
func func3() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KSC View >", "target file"}, []string{"bool", "file"}, []string{"T", "*"}, *getpath(), nil)
	sel.GetOpt()

	wk := ksc.Initksc()
	wk.Predetect = true
	wk.Path = sel.StrRes[1]
	err = wk.Readf()
	if err != nil {
		return err
	}

	fmt.Printf("%s\nprehead %d\ncommon %s, subtype %s, reserved %s\nhead %d, realsize %d\n", wk.Path, len(wk.Prehead), kio.Bprint(wk.Common), kio.Bprint(wk.Subtype), kio.Bprint(wk.Reserved), wk.Headp, wk.Rsize)
	for i, r := range wk.Chunkpos {
		fmt.Printf("data %d : offset %d, size %d\n", i, r, wk.Chunksize[i]+8)
	}
	fmt.Println("G5KSC View Success")
	return err
}

// g5ksc pack
func func4() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KSC Pack >", "prehead image", "head.common", "head.subtype", "head.reserved", "data 0", "data 1", "data 2", "data 3", "data 4", "data 5", "data 6", "data 7", "add end sign"},
		[]string{"bool", "file", "bytes", "bytes", "bytes", "file", "file", "file", "file", "file", "file", "file", "file", "bool"},
		[]string{"T", "*", "4", "4", "8", "*", "*", "*", "*", "*", "*", "*", "*", "*"}, *getpath(), nil)
	sel.ByteRes = [][]byte{{0}, nil, []byte("KSC5"), {0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0}, nil, nil, nil, nil, nil, nil, nil, nil, {0}}
	sel.GetOpt()

	wk := ksc.Initksc()
	if sel.StrRes[1] != "" {
		f, _ := kio.Open(sel.StrRes[1], "r")
		wk.Prehead, _ = kio.Read(f, -1)
		add := make([]byte, 512-len(wk.Prehead)%512)
		if len(add) != 512 {
			wk.Prehead = append(wk.Prehead, add...)
		}
		f.Close()
	}
	wk.Common = sel.ByteRes[2]
	wk.Subtype = sel.ByteRes[3]
	wk.Reserved = sel.ByteRes[4]
	wk.Path = "./_ST5_DATA/result.bin"

	err = wk.Writef()
	if err != nil {
		return err
	}
	for i := 0; i < 8; i++ {
		if sel.StrRes[i+5] != "" {
			err = wk.Addf(sel.StrRes[i+5])
			if err != nil {
				return err
			}
		}
	}
	if sel.ByteRes[13][0] == 0 {
		wk.Addf("")
	}
	fmt.Println("G5KSC Pack Success (result.bin)")
	return err
}

// g5ksc unpack
func func5() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KSC Unpack >", "target file"}, []string{"bool", "file"}, []string{"T", "*"}, *getpath(), nil)
	sel.GetOpt()

	wk := ksc.Initksc()
	wk.Predetect = true
	wk.Path = sel.StrRes[1]
	err = wk.Readf()
	if err != nil {
		return err
	}

	fmt.Printf("%s\nhead %d, realsize %d\n", wk.Path, wk.Headp, wk.Rsize)
	f, _ := kio.Open("./_ST5_DATA/prehead.webp", "w")
	kio.Write(f, wk.Prehead)
	f.Close()
	f, _ = kio.Open("./_ST5_DATA/head.txt", "w")
	kio.Write(f, []byte(fmt.Sprintf("common = '%s'\nsubtype = '%s'\nreserved = '%s'\n", kio.Bprint(wk.Common), kio.Bprint(wk.Subtype), kio.Bprint(wk.Reserved))))
	f.Close()
	f, _ = kio.Open(wk.Path, "r")
	for i, r := range wk.Chunkpos {
		f.Seek(int64(r+8), 0)
		d, _ := kio.Read(f, wk.Chunksize[i])
		t, _ := kio.Open(fmt.Sprintf("./_ST5_DATA/data%d.bin", i), "w")
		kio.Write(t, d)
		t.Close()
	}
	f.Close()
	fmt.Println("G5KSC Unpack Success")
	return err
}

// g5kzip pack
func func6() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KZIP Pack >", "webp mode", "png mode", "file 0", "file 1", "file 2", "file 3", "file 4", "file 5", "dir 0", "dir 1", "dir 2", "dir 3", "dir 4", "dir 5"},
		[]string{"bool", "bool", "bool", "file", "file", "file", "file", "file", "file", "folder", "folder", "folder", "folder", "folder", "folder"},
		[]string{"T", "*", "*", "*", "*", "*", "*", "*", "*", "NR", "NR", "NR", "NR", "NR", "NR"}, *getpath(), nil)
	sel.GetOpt()

	mode := "bin"
	if sel.ByteRes[1][0] == 0 {
		mode = "webp"
	} else if sel.ByteRes[2][0] == 0 {
		mode = "png"
	}
	paths := make([]string, 0)
	for i := 0; i < 6; i++ {
		if sel.StrRes[i+3] != "" {
			paths = append(paths, sel.StrRes[i+3])
		}
		if sel.StrRes[i+9] != "./" {
			paths = append(paths, sel.StrRes[i+9])
		}
	}
	err = kzip.Dozip(paths, mode, "./_ST5_DATA/result."+mode)
	if err != nil {
		return err
	}
	fmt.Printf("G5KZIP Pack Success (%s)\n", "result."+mode)
	return err
}

// g5kzip unpack
func func7() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KZIP Unpack >", "target file", "check error"}, []string{"bool", "file", "bool"}, []string{"T", "*", "*"}, *getpath(), nil)
	sel.GetOpt()

	err = kzip.Unzip(sel.StrRes[1], "./_ST5_DATA/result/", sel.ByteRes[2][0] == 0)
	if err != nil {
		return err
	}
	fmt.Println("G5KZIP Unpack Success (result/)")
	return err
}

// g5kaes enc all
func func8() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KAES.ALL Encrypt >", "target file", "password", "keyfile", "hint", "message", "sign", "webp mode", "png mode"},
		[]string{"bool", "file", "string", "keyfile", "string", "string", "keyfile", "bool", "bool"},
		[]string{"T", "*", "*", "*", "*", "*", "*", "*", "*"}, *getpath(), kaes.Basickey())
	sel.GetOpt()

	pmode := 2
	if sel.ByteRes[7][0] == 0 {
		pmode = 0
	} else if sel.ByteRes[8][0] == 0 {
		pmode = 1
	}
	g5kaes_all.Hint = sel.StrRes[4]
	g5kaes_all.Msg = sel.StrRes[5]
	g5kaes_all.Signkey = getsign(sel.ByteRes[6])
	var res string
	res, err = g5kaes_all.EnFile([]byte(sel.StrRes[2]), sel.ByteRes[3], sel.StrRes[1], pmode)
	if err != nil {
		return err
	}
	fmt.Printf("G5KAES.ALL Encrypt Success (%s)\n", res)
	return err
}

// g5kaes dec all
func func9() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	path := kio.Input("HintViewer file : ")
	err = g5kaes_all.ViewFile(path)
	if err != nil {
		return err
	}

	var sel cliget.OptSel
	sel.Init([]string{"< G5KAES.ALL Decrypt >", "target file", "password", "keyfile", "hint", "message", "sign"},
		[]string{"bool", "file", "string", "keyfile", "string", "string", "string"}, []string{"T", "*", "*", "*", "0", "0", "0"}, *getpath(), kaes.Basickey())
	sel.StrRes = []string{"", path, "", "", g5kaes_all.Hint, g5kaes_all.Msg, cfg.check(g5kaes_all.Signkey[0])}
	sel.GetOpt()
	var res string
	res, err = g5kaes_all.DeFile([]byte(sel.StrRes[2]), sel.ByteRes[3], sel.StrRes[1])
	if err != nil {
		return err
	}
	fmt.Printf("G5KAES.ALL Decrypt Success (%s)\n", res)
	return err
}

// g5kaes enc func
func func10() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KAES.FUNC Encrypt >", "target file", "key"}, []string{"bool", "file", "bytes"}, []string{"T", "*", "48"}, *getpath(), nil)
	sel.ByteRes[2] = make([]byte, 48)
	sel.GetOpt()

	err = g5kaes_func.Before.Open(sel.StrRes[1], true)
	if err != nil {
		return err
	}
	defer g5kaes_func.Before.Close()
	g5kaes_func.After.Open("./_ST5_DATA/result.bin", false)
	defer g5kaes_func.After.Close()
	err = g5kaes_func.Encrypt(sel.ByteRes[2])
	if err != nil {
		return err
	}
	fmt.Println("G5KAES.FUNC Encrypt Success (result.bin)")
	return err
}

// g5kaes dec func
func func11() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KAES.FUNC Decrypt >", "target file", "key"}, []string{"bool", "file", "bytes"}, []string{"T", "*", "48"}, *getpath(), nil)
	sel.ByteRes[2] = make([]byte, 48)
	sel.GetOpt()
	err = g5kaes_func.Before.Open(sel.StrRes[1], true)
	if err != nil {
		return err
	}
	defer g5kaes_func.Before.Close()
	g5kaes_func.After.Open("./_ST5_DATA/result.bin", false)
	defer g5kaes_func.After.Close()
	err = g5kaes_func.Decrypt(sel.ByteRes[2])
	if err != nil {
		return err
	}
	fmt.Println("G5KAES.FUNC Decrypt Success (result.bin)")
	return err
}

// g5kpkg view
func func12() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KPKG View >", "target file"}, []string{"bool", "file"}, []string{"T", "webp"}, *getpath(), nil)
	sel.GetOpt()

	wk := ksc.Initksc()
	wk.Predetect = true
	wk.Path = sel.StrRes[1]
	err = wk.Readf()
	if err != nil {
		return err
	}

	if !kio.Bequal(wk.Subtype, []byte("KPKG")) {
		return errors.New("invalid package")
	}
	f, _ := kio.Open(wk.Path, "r")
	f.Seek(int64(wk.Chunkpos[0]+8), 0)
	d0, _ := kio.Read(f, wk.Chunksize[0])
	f.Seek(int64(wk.Chunkpos[1]+8), 0)
	d1, _ := kio.Read(f, wk.Chunksize[1])
	f.Close()
	fmt.Printf("%s\n%s\n", d0, d1)
	fmt.Println("G5KPKG View Success")
	return err
}

// g5kpkg pack
func func13() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KPKG Pack >", "package name", "package version", "package text", "package sign", "target 0", "osnum 0", "target 1", "osnum 1"},
		[]string{"bool", "string", "float", "string", "keyfile", "folder", "int", "folder", "int"}, []string{"T", "*", "0+", "*", "*", "NR", "0+", "NR", "0+"}, *getpath(), nil)
	sel.GetOpt()

	wk := kpkg.Initkpkg(0)
	t := getsign(sel.ByteRes[4])
	wk.Public = t[0]
	wk.Private = t[1]
	wk.Name = sel.StrRes[1]
	wk.Version, _ = strconv.ParseFloat(sel.StrRes[2], 64)
	wk.Text = sel.StrRes[3]
	ns := make([]int, 0)
	ps := make([]string, 0)
	if sel.StrRes[5] != "./" {
		t, _ := strconv.Atoi(sel.StrRes[6])
		ns = append(ns, t)
		ps = append(ps, sel.StrRes[5])
	}
	if sel.StrRes[7] != "./" {
		t, _ := strconv.Atoi(sel.StrRes[8])
		ns = append(ns, t)
		ps = append(ps, sel.StrRes[7])
	}
	err = wk.Pack(ns, ps, "./_ST5_DATA/result.webp")
	if err != nil {
		return err
	}
	os.RemoveAll("./temp675/")
	fmt.Println("G5KPKG Pack Success (result.webp)")
	return err
}

// g5kpkg unpack
func func14() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KPKG Unpack >", "target file"}, []string{"bool", "file"}, []string{"T", "webp"}, *getpath(), nil)
	sel.GetOpt()

	wk := ksc.Initksc()
	wk.Predetect = true
	wk.Path = sel.StrRes[1]
	err = wk.Readf()
	if err != nil {
		return err
	}

	f, _ := kio.Open(wk.Path, "r")
	for i, r := range wk.Chunkpos {
		f.Seek(int64(r+8), 0)
		d, _ := kio.Read(f, wk.Chunksize[i])
		nm := fmt.Sprintf("./_ST5_DATA/zip%d.webp", i)
		if i == 0 {
			nm = "./_ST5_DATA/info.txt"
		} else if i == 1 {
			nm = "./_ST5_DATA/sign.txt"
		}
		t, _ := kio.Open(nm, "w")
		kio.Write(t, d)
		t.Close()
	}
	f.Close()
	fmt.Println("G5KPKG Unpack Success")
	return err
}

// g5kpic pack
func func15() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KPIC Pack >", "target file", "style webp", "style png", "high integration", "mode [ltl, std, big]", "picture", "size x", "size y"},
		[]string{"bool", "file", "bool", "bool", "bool", "string", "file", "int", "int"}, []string{"T", "*", "*", "*", "*", "3", "*", "+", "+"}, *getpath(), nil)
	sel.StrRes[5] = "dft"
	sel.GetOpt()

	style := "bmp"
	if sel.ByteRes[2][0] == 0 {
		style = "webp"
	} else if sel.ByteRes[3][0] == 0 {
		style = "png"
	}
	high := (sel.ByteRes[4][0] == 0)
	pic := ""
	switch sel.StrRes[5] {
	case "ltl":
		pic = "./little.webp"
	case "std":
		pic = "./standard.webp"
	case "big":
		pic = "./big.webp"
	}
	if sel.StrRes[6] != "" {
		pic = sel.StrRes[6]
	}
	x, _ := strconv.Atoi(sel.StrRes[7])
	if x <= 0 {
		x = -1
	}
	y, _ := strconv.Atoi(sel.StrRes[8])
	if y <= 0 {
		y = -1
	}

	wk, err := kpic.Initpic(pic, x, y)
	if err != nil {
		return err
	}
	wk.Target = sel.StrRes[1]
	wk.Export = "./_ST5_DATA/"
	wk.Style = style
	var name string
	var num int
	if high {
		name, num = wk.Pack(2)
	} else {
		name, num = wk.Pack(4)
	}
	fmt.Printf("G5KPIC Pack Success : %s.%s * %d\n", name, style, num)
	return err
}

// g5kpic unpack
func func16() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KPIC Unpack >", "dir to detect"}, []string{"bool", "folder"}, []string{"T", "NR"}, *getpath(), nil)
	sel.GetOpt()

	wk, err := kpic.Initpic("", -1, -1)
	if err != nil {
		return err
	}
	wk.Target = sel.StrRes[1]
	wk.Export = "./_ST5_DATA/result.bin"
	name, num, style := wk.Detect()
	wk.Style = style
	wk.Unpack(name, num)
	fmt.Println("G5KPIC Unpack Success (result.bin)")
	return err
}

// g5kpic restore
func func17() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G5KPIC Restore >", "dir to restore", "style [webp, png, bmp]", "high integration"},
		[]string{"bool", "folder", "string", "bool"}, []string{"T", "NR", "*", "*"}, *getpath(), nil)
	sel.StrRes[2] = "webp"
	sel.GetOpt()

	wk, err := kpic.Initpic("", -1, -1)
	if err != nil {
		return err
	}
	wk.Style = sel.StrRes[2]
	files := make([]string, 0)
	fs, _ := os.ReadDir(sel.StrRes[1])
	for _, r := range fs {
		files = append(files, sel.StrRes[1]+r.Name())
	}
	var name string
	var num int
	if sel.ByteRes[3][0] == 0 {
		name, num = wk.Restore(files, 2)
	} else {
		name, num = wk.Restore(files, 4)
	}
	fmt.Printf("G5KPIC Restore Success : %s.%s * %d\n", name, wk.Style, num)
	return err
}

// kscript compile
func func18() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< kscript compile >", "source", "optconst", "optasm", "header 0", "header 1", "header 2", "header 3", "header 4", "header 5"},
		[]string{"bool", "file", "bool", "bool", "file", "file", "file", "file", "file", "file"},
		[]string{"T", "txt", "*", "*", "txt", "txt", "txt", "txt", "txt", "txt"}, *getpath(), nil)
	sel.GetOpt()

	g5kscript.init()
	for i := 0; i < 6; i++ {
		if sel.StrRes[i+4] != "" {
			err = g5kscript.addpkg(sel.StrRes[i+4])
			if err != nil {
				return err
			}
		}
	}
	g5kscript.option = [5]bool{true, true, true, sel.ByteRes[2][0] == 0, sel.ByteRes[3][0] == 0}
	f, _ := kio.Open(sel.StrRes[1], "r")
	d, _ := kio.Read(f, -1)
	f.Close()
	g5kscript.data[0] = string(d)

	_, _, err = g5kscript.compile()
	if err != nil {
		return err
	}
	f, _ = kio.Open("./_ST5_DATA/tkn.txt", "w")
	kio.Write(f, []byte(g5kscript.data[4]))
	f.Close()
	f, _ = kio.Open("./_ST5_DATA/ast.txt", "w")
	kio.Write(f, []byte(g5kscript.data[5]))
	f.Close()
	f, _ = kio.Open("./_ST5_DATA/asm.txt", "w")
	kio.Write(f, []byte(g5kscript.data[6]))
	f.Close()
	fmt.Println("Compile Success : tkn.txt, ast.txt, asm.txt")
	return err
}

// kscript assemble
func func19() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< kscript assemble >", "asm", "info", "icon", "sign"}, []string{"bool", "file", "string", "file", "keyfile"}, []string{"T", "txt", "*", "webp", "*"}, *getpath(), nil)
	sel.GetOpt()

	t := getsign(sel.ByteRes[4])
	g5kscript.assmebler.SetKey(sel.StrRes[3], t[0], t[1])
	g5kscript.assmebler.Info = sel.StrRes[2]
	f, _ := kio.Open(sel.StrRes[1], "r")
	d, _ := kio.Read(f, -1)
	f.Close()
	data, err := g5kscript.assmebler.GenExe(string(d))
	if err != nil {
		return err
	}
	f, _ = kio.Open("./_ST5_DATA/exe.webp", "w")
	kio.Write(f, data)
	f.Close()
	fmt.Println("Assemble Success : exe.webp")
	return err
}

// kscript reversing
func func20() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< kscript reversing >", "target file"}, []string{"bool", "file"}, []string{"T", "webp"}, *getpath(), nil)
	sel.GetOpt()

	f, _ := kio.Open(sel.StrRes[1], "r")
	data, _ := kio.Read(f, -1)
	f.Close()
	wk := ksc.Initksc()
	wk.Predetect = true
	err = wk.Readb(data)
	if err != nil {
		return err
	}
	if !kio.Bequal(wk.Subtype, []byte("KELF")) {
		return errors.New("invalid KELF")
	}

	f, _ = kio.Open("./_ST5_DATA/info.txt", "w")
	kio.Write(f, data[wk.Chunkpos[0]+8:wk.Chunkpos[0]+8+wk.Chunksize[0]])
	f.Close()
	t := ".rodata\n\n" + kscript_revdata(data[wk.Chunkpos[1]+8:wk.Chunkpos[1]+8+wk.Chunksize[1]])
	t = t + "\n.data\n\n" + kscript_revdata(data[wk.Chunkpos[2]+8:wk.Chunkpos[2]+8+wk.Chunksize[2]])
	t = t + "\n.text\n\n" + kscript_revcode(data[wk.Chunkpos[3]+8:wk.Chunkpos[3]+8+wk.Chunksize[3]])
	f, _ = kio.Open("./_ST5_DATA/code.txt", "w")
	kio.Write(f, []byte(t))
	f.Close()
	fmt.Println("Reversing Success : info.txt, code.txt")
	return err
}

// kv4adv mknew
func func21() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< KV4adv MkNew >", "cluster name", "chunksize", "fptr/unit"}, []string{"bool", "string", "int", "int"}, []string{"T", "+", "+", "+"}, *getpath(), nil)
	sel.GetOpt()

	name := sel.StrRes[1]
	if name == "" {
		name = "NewDrive"
	}
	csize := sel.StrRes[2]
	if csize == "0" {
		csize = "134217728"
	}
	maxf := sel.StrRes[3]
	if maxf == "0" {
		maxf = "128"
	}
	path := fmt.Sprintf("./_ST5_DATA/clu%d/", kobj.Decode(kaes.Genrand(2)))
	os.Mkdir(path, os.ModePerm)
	g5shell.g4sh.Command("init", nil)
	err = g5shell.g4sh.Command("new", []string{path, name, csize, maxf})
	if err != nil {
		return err
	}
	fmt.Printf("KV4adv MkNew Success : %s\n", path[12:])
	return err
}

// kv4adv login
func func22() (err error) {
	defer func() {
		g5shell.exit()
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	path := kio.Abs(kio.Input("HintViewer cluster : "))
	g5shell.init(true)
	local := fmt.Sprintf("./local%d/", kobj.Decode(kaes.Genrand(2)))
	err = g5shell.g4sh.Command("boot", []string{cfg.desktop, local, path, path + "0a.webp"})
	if err != nil {
		return err
	}

	var sel cliget.OptSel
	sel.Init([]string{"< KV4adv Login >", "cluster", "password", "keyfile", "hint", "sleep [1, 10, 30, 60]"},
		[]string{"bool", "folder", "string", "keyfile", "string", "int"}, []string{"T", "NR", "*", "*", "0", "+"}, *getpath(), kaes.Basickey())
	sel.StrRes = []string{"", path, "", "", fmt.Sprintf("%s [%s@%s]", string(g5shell.g4sh.IObyte[0]), g5shell.g4sh.IOstr[1], g5shell.g4sh.IOstr[0]), ""}
	sel.GetOpt()
	err = g5shell.repl([]byte(sel.StrRes[2]), sel.ByteRes[3], sel.StrRes[5])
	fmt.Println("KV4adv KVdrive Exit")
	return err
}

// kv4adv rebuild
func func23() (err error) {
	defer func() {
		os.RemoveAll("./temp740/")
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	path := kio.Abs(kio.Input("HintViewer cluster : "))
	g5shell.init(true)
	err = g5shell.g4sh.Command("boot", []string{cfg.desktop, "./temp740/", path, path + "0a.webp"})
	if err != nil {
		return err
	}

	var sel cliget.OptSel
	sel.Init([]string{"< KV4adv Rebuild >", "cluster", "password", "keyfile", "hint"},
		[]string{"bool", "folder", "string", "keyfile", "string"}, []string{"T", "NR", "*", "*", "0"}, *getpath(), kaes.Basickey())
	sel.StrRes = []string{"", path, "", "", fmt.Sprintf("%s [%s@%s]", string(g5shell.g4sh.IObyte[0]), g5shell.g4sh.IOstr[1], g5shell.g4sh.IOstr[0])}
	sel.GetOpt()
	g5shell.g4sh.IObyte = [][]byte{[]byte(sel.StrRes[2]), sel.ByteRes[3]}
	err = g5shell.g4sh.Command("rebuild", []string{path})
	if err != nil {
		return err
	}
	fmt.Println("KV4adv Rebuild Success")
	return err
}

// kv5st mknew
func func24() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< KV5st MkNew >", "cluster name", "csize [sml, std, lrg]"}, []string{"bool", "string", "string"}, []string{"T", "+", "3"}, *getpath(), nil)
	sel.GetOpt()

	name := sel.StrRes[1]
	if name == "" {
		name = "NewDrive"
	}
	csize := "default"
	switch sel.StrRes[2] {
	case "sml":
		csize = "small"
	case "std":
		csize = "standard"
	case "lrg":
		csize = "large"
	}
	path := fmt.Sprintf("./_ST5_DATA/clu%d/", kobj.Decode(kaes.Genrand(2)))
	os.Mkdir(path, os.ModePerm)
	g5shell.g5sh.Command("init", []string{"false"})
	err = g5shell.g5sh.Command("new", []string{path, name, csize})
	if err != nil {
		return err
	}
	fmt.Printf("KV5st MkNew Success : %s\n", path[12:])
	return err
}

// kv5st login
func func25() (err error) {
	defer func() {
		g5shell.exit()
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	path := kio.Abs(kio.Input("HintViewer cluster : "))
	g5shell.init(false)
	local := fmt.Sprintf("./local%d/", kobj.Decode(kaes.Genrand(2)))
	err = g5shell.g5sh.Command("boot", []string{cfg.desktop, local, path, path + "0a.webp"})
	if err != nil {
		return err
	}

	var sel cliget.OptSel
	sel.Init([]string{"< KV5st Login >", "cluster", "password", "keyfile", "hint", "sleep [1, 10, 30, 60]"},
		[]string{"bool", "folder", "string", "keyfile", "string", "int"}, []string{"T", "NR", "*", "*", "0", "+"}, *getpath(), kaes.Basickey())
	sel.StrRes = []string{"", path, "", "", fmt.Sprintf("%s [%s@%s]", string(g5shell.g5sh.IObyte[0]), g5shell.g5sh.IOstr[1], g5shell.g5sh.IOstr[0]), ""}
	sel.GetOpt()
	err = g5shell.repl([]byte(sel.StrRes[2]), sel.ByteRes[3], sel.StrRes[5])
	fmt.Println("KV5st KVdrive Exit")
	return err
}

// kv5st rebuild
func func26() (err error) {
	defer func() {
		os.RemoveAll("./temp740/")
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	path := kio.Abs(kio.Input("HintViewer cluster : "))
	g5shell.init(false)
	err = g5shell.g5sh.Command("boot", []string{cfg.desktop, "./temp740/", path, path + "0a.webp"})
	if err != nil {
		return err
	}

	var sel cliget.OptSel
	sel.Init([]string{"< KV5st Rebuild >", "cluster", "password", "keyfile", "hint"},
		[]string{"bool", "folder", "string", "keyfile", "string"}, []string{"T", "NR", "*", "*", "0"}, *getpath(), kaes.Basickey())
	sel.StrRes = []string{"", path, "", "", fmt.Sprintf("%s [%s@%s]", string(g5shell.g5sh.IObyte[0]), g5shell.g5sh.IOstr[1], g5shell.g5sh.IOstr[0])}
	sel.GetOpt()
	g5shell.g5sh.IObyte = [][]byte{[]byte(sel.StrRes[2]), sel.ByteRes[3]}
	err = g5shell.g5sh.Command("rebuild", []string{path})
	if err != nil {
		return err
	}
	fmt.Println("KV5st Rebuild Success")
	return err
}

// g1kenc encrypt
func func27() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G1KENC Encrypt >", "target file", "password", "hint"}, []string{"bool", "file", "string", "string"}, []string{"T", "*", "*", "*"}, *getpath(), nil)
	sel.GetOpt()

	g1kenc.Init()
	g1kenc.Path = sel.StrRes[1]
	g1kenc.Pw = sel.StrRes[2]
	g1kenc.Hint = sel.StrRes[3]
	err = g1kenc.Encrypt()
	if err != nil {
		return err
	}
	fmt.Println("G1KENC Encrypt Success")
	return err
}

// g1kenc decrypt
func func28() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	g1kenc.Init()
	g1kenc.Path = kio.Input("HintViewer file : ")
	err = g1kenc.View()
	if err != nil {
		return err
	}

	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G1KENC Decrypt >", "target file", "password", "hint"}, []string{"bool", "file", "string", "string"}, []string{"T", "k", "*", "0"}, *getpath(), nil)
	sel.StrRes = []string{"", g1kenc.Path, "", g1kenc.Hint}
	sel.GetOpt()
	g1kenc.Pw = sel.StrRes[2]
	err = g1kenc.Decrypt()
	if err != nil {
		return err
	}
	fmt.Println("G1KENC Decrypt Success")
	return err
}

// g2kenc encrypt
func func29() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G2KENC Encrypt >", "target file", "password", "hint", "hidename"}, []string{"bool", "file", "string", "string", "bool"}, []string{"T", "*", "*", "*", "*"}, *getpath(), nil)
	sel.GetOpt()

	g2kenc.Init()
	g2kenc.Path = sel.StrRes[1]
	g2kenc.Pw = sel.StrRes[2]
	g2kenc.Hint = sel.StrRes[3]
	g2kenc.Hidename = (sel.ByteRes[4][0] == 0)
	path, err := g2kenc.Encrypt()
	if err != nil {
		return err
	}
	fmt.Printf("G2KENC Encrypt Success (%s)\n", path)
	return err
}

// g2kenc decrypt
func func30() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	g2kenc.Init()
	g2kenc.Path = kio.Input("HintViewer file : ")
	err = g2kenc.View()
	if err != nil {
		return err
	}

	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G2KENC Decrypt >", "target file", "password", "hint"}, []string{"bool", "file", "string", "string"}, []string{"T", "k", "*", "0"}, *getpath(), nil)
	sel.StrRes = []string{"", g2kenc.Path, "", g2kenc.Hint}
	sel.GetOpt()
	g2kenc.Pw = sel.StrRes[2]
	err = g2kenc.Decrypt()
	if err != nil {
		return err
	}
	fmt.Println("G2KENC Decrypt Success (./_)")
	return err
}

// g3kzip pack file
func func31() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G3KZIP Pack File >", "prehead picture", "file 0", "file 1", "file 2", "file 3", "winsign (/ -> \\)"},
		[]string{"bool", "file", "file", "file", "file", "file", "bool"}, []string{"T", "png", "*", "*", "*", "*", "*"}, *getpath(), nil)
	sel.ByteRes[6][0] = 1
	sel.GetOpt()

	g3kzip.Init()
	if sel.StrRes[1] != "" {
		f, _ := kio.Open(sel.StrRes[1], "r")
		g3kzip.Prehead, _ = kio.Read(f, -1)
		add := make([]byte, 1024-len(g3kzip.Prehead)%1024)
		if len(add) != 1024 {
			g3kzip.Prehead = append(g3kzip.Prehead, add...)
		}
		f.Close()
	}
	files := make([]string, 0)
	for i := 0; i < 4; i++ {
		if sel.StrRes[i+2] != "" {
			files = append(files, sel.StrRes[i+2])
		}
	}
	g3kzip.Winsign = (sel.ByteRes[6][0] == 0)
	err = g3kzip.Packf(files, "./_ST5_DATA/result.png")
	if err != nil {
		return err
	}
	fmt.Println("G3KZIP Pack Success (result.png)")
	return err
}

// g3kzip pack folder
func func32() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G3KZIP Pack Dir >", "prehead picture", "target dir", "winsign (/ -> \\)"},
		[]string{"bool", "file", "folder", "bool"}, []string{"T", "png", "*", "*"}, *getpath(), nil)
	sel.ByteRes[3][0] = 1
	sel.GetOpt()

	g3kzip.Init()
	if sel.StrRes[1] != "" {
		f, _ := kio.Open(sel.StrRes[1], "r")
		g3kzip.Prehead, _ = kio.Read(f, -1)
		add := make([]byte, 1024-len(g3kzip.Prehead)%1024)
		if len(add) != 1024 {
			g3kzip.Prehead = append(g3kzip.Prehead, add...)
		}
		f.Close()
	}
	g3kzip.Winsign = (sel.ByteRes[3][0] == 0)
	err = g3kzip.Packd(sel.StrRes[2], "./_ST5_DATA/result.png")
	if err != nil {
		return err
	}
	fmt.Println("G3KZIP Pack Success (result.png)")
	return err
}

// g3kzip unpack
func func33() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G3KZIP Unpack >", "target file"}, []string{"bool", "file"}, []string{"T", "png"}, *getpath(), nil)
	sel.GetOpt()

	path := sel.StrRes[1]
	g3kzip.Init()
	err = g3kzip.View(path)
	if err != nil {
		return err
	}
	err = g3kzip.Unpack(path)
	if err != nil {
		return err
	}
	os.Rename("./temp261/", "./_ST5_DATA/temp261/")
	fmt.Println("G3KZIP Unpack Success (temp261/)")
	return err
}

// g3kaes enc all
func func34() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G3KAES.ALL Encrypt >", "target file", "password", "keyfile", "hint", "hide name", "prehead picture", "subhead file", "threads num", "chunksize"},
		[]string{"bool", "file", "string", "keyfile", "string", "bool", "file", "file", "int", "int"}, []string{"T", "*", "*", "*", "*", "*", "png", "*", "+", "+"}, *getpath(), legsup.G3kf())
	sel.GetOpt()

	core, _ := strconv.Atoi(sel.StrRes[8])
	chunk, _ := strconv.Atoi(sel.StrRes[9])
	g3kaes_all.Init(core, chunk)
	if sel.StrRes[6] != "" {
		f, _ := kio.Open(sel.StrRes[6], "r")
		g3kaes_all.Prehead, _ = kio.Read(f, -1)
		add := make([]byte, 1024-len(g3kaes_all.Prehead)%1024)
		if len(add) != 1024 {
			g3kaes_all.Prehead = append(g3kaes_all.Prehead, add...)
		}
		f.Close()
	}
	if sel.StrRes[7] != "" {
		f, _ := kio.Open(sel.StrRes[7], "r")
		g3kaes_all.Subhead, _ = kio.Read(f, -1)
		f.Close()
	}
	g3kaes_all.Hidename = (sel.ByteRes[5][0] == 0)
	g3kaes_all.Hint = sel.StrRes[4]
	err = g3kaes_all.Encrypt(sel.StrRes[1], sel.StrRes[2], sel.ByteRes[3])
	if err != nil {
		return err
	}
	fmt.Printf("G3KAES.ALL Encrypt Success (%s)\n", g3kaes_all.Respath)
	return err
}

// g3kaes dec all
func func35() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	path := kio.Input("HintViewer file : ")
	err = g3kaes_all.View(path)
	if err != nil {
		return err
	}

	var sel cliget.OptSel
	sel.Init([]string{"< G3KAES.ALL Decrypt >", "target file", "password", "keyfile", "hint"}, []string{"bool", "file", "string", "keyfile", "string"}, []string{"T", "png", "*", "*", "0"}, *getpath(), legsup.G3kf())
	sel.StrRes = []string{"", path, "", "", g3kaes_all.Hint}
	sel.GetOpt()
	err = g3kaes_all.Decrypt(sel.StrRes[1], sel.StrRes[2], sel.ByteRes[3])
	if err != nil {
		return err
	}
	if len(g3kaes_all.Subhead) != 0 {
		f, _ := kio.Open("./_ST5_DATA/subhead.bin", "w")
		kio.Write(f, g3kaes_all.Subhead)
		f.Close()
	}
	fmt.Printf("G3KAES.ALL Decrypt Success (%s)\n", g3kaes_all.Respath)
	return err
}

// g3kaes enc func
func func36() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G3KAES.FUNC Encrypt >", "target file", "key"}, []string{"bool", "file", "bytes"}, []string{"T", "*", "32"}, *getpath(), nil)
	sel.ByteRes[2] = make([]byte, 32)
	sel.GetOpt()
	err = g3kaes_func.Encrypt(sel.StrRes[1], "./_ST5_DATA/result.bin", sel.ByteRes[2])
	if err != nil {
		return err
	}
	fmt.Println("G3KAES.FUNC Encrypt Success (result.bin)")
	return err
}

// g3kaes dec func
func func37() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G3KAES.FUNC Decrypt >", "target file", "key"}, []string{"bool", "file", "bytes"}, []string{"T", "*", "32"}, *getpath(), nil)
	sel.ByteRes[2] = make([]byte, 32)
	sel.GetOpt()
	err = g3kaes_func.Decrypt(sel.StrRes[1], "./_ST5_DATA/result.bin", sel.ByteRes[2])
	if err != nil {
		return err
	}
	fmt.Println("G3KAES.FUNC Decrypt Success (result.bin)")
	return err
}

// g3kpic pack
func func38() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G3KPIC Pack >", "target file", "high integration", "picture", "size x", "size y"},
		[]string{"bool", "file", "bool", "file", "int", "int"}, []string{"T", "*", "*", "png", "+", "+"}, *getpath(), nil)
	sel.ByteRes[2][0] = 1
	sel.GetOpt()

	x, _ := strconv.Atoi(sel.StrRes[4])
	if x <= 0 {
		x = -1
	}
	y, _ := strconv.Atoi(sel.StrRes[5])
	if y <= 0 {
		y = -1
	}
	err = g3kpic.Init(sel.StrRes[3], x, y)
	if err != nil {
		return err
	}
	g3kpic.Pcover = (sel.ByteRes[2][0] != 0)
	name, num := g3kpic.Pack(sel.StrRes[1], "./_ST5_DATA/")
	fmt.Printf("G3KPIC Pack Success : %s.png * %d\n", name, num)
	return err
}

// g3kpic unpack
func func39() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G3KPIC Unpack >", "dir to detect"}, []string{"bool", "folder"}, []string{"T", "NR"}, *getpath(), nil)
	sel.GetOpt()

	err = g3kpic.Init("", -1, -1)
	if err != nil {
		return err
	}
	name, num, err := g3kpic.Detect(sel.StrRes[1])
	if err != nil {
		return err
	}
	g3kpic.Unpack("./_ST5_DATA/result.bin", sel.StrRes[1], name, num)
	fmt.Println("G3KPIC Unpack Success (result.bin)")
	return err
}

// g3kv enc
func func40() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G3KV Encrypt >", "target dir", "password", "keyfile", "hint"}, []string{"bool", "folder", "string", "keyfile", "string"}, []string{"T", "NR", "*", "*", "*"}, *getpath(), legsup.G3kf())
	sel.GetOpt()

	g3kv.Hint = sel.StrRes[4]
	err = g3kv.Encrypt(sel.StrRes[2], sel.ByteRes[3], sel.StrRes[1])
	if err != nil {
		return err
	}
	fmt.Printf("G3KV Encrypt Success (%s)\n", sel.StrRes[1])
	return err
}

// g3kv dec
func func41() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	path := kio.Input("HintViewer file : ")
	err = g3kv.View(path)
	if err != nil {
		return err
	}

	var sel cliget.OptSel
	sel.Init([]string{"< G3KV Decrypt >", "target dir", "password", "keyfile", "hint", "ench.txt file"},
		[]string{"bool", "folder", "string", "keyfile", "string", "string"}, []string{"T", "NR", "*", "*", "0", "0"}, *getpath(), legsup.G3kf())
	sel.StrRes = []string{"", "", "", "", g3kv.Hint, path}
	sel.GetOpt()
	err = g3kv.Decrypt(sel.StrRes[2], sel.ByteRes[3], sel.StrRes[1])
	if err != nil {
		return err
	}
	fmt.Printf("G3KV Decrypt Success (%s)\n", sel.StrRes[1])
	return err
}

// g3zip release
func func42() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G3ZIP Release >", "picture", "file 0", "file 1", "file 2", "file 3", "file 4", "file 5"},
		[]string{"bool", "file", "file", "file", "file", "file", "file", "file"}, []string{"T", "*", "*", "*", "*", "*", "*", "*"}, *getpath(), nil)
	sel.GetOpt()

	files := make([]string, 0)
	for i := 0; i < 6; i++ {
		if sel.StrRes[i+2] != "" {
			files = append(files, sel.StrRes[i+2])
		}
	}
	var pic []byte
	mode := ".zip"
	if sel.StrRes[1] == "" {
		if len(files) < 4 {
			pic = picdt.Zr5png()
			mode = ".png"
		} else {
			pic = picdt.Zr5webp()
			mode = ".webp"
		}
	} else {
		f, _ := kio.Open(sel.StrRes[1], "r")
		pic, _ = kio.Read(f, -1)
		f.Close()
		if strings.Contains(sel.StrRes[1], ".") {
			mode = sel.StrRes[1][strings.LastIndex(sel.StrRes[1], "."):]
		}
	}
	err = legsup.G3picre(pic, files, "./_ST5_DATA/result"+mode)
	if err != nil {
		return err
	}
	fmt.Printf("G3ZIP Release Success (%s)\n", "result"+mode)
	return err
}

// g4kenc encrypt
func func43() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G4KENC Encrypt >", "password", "hint", "file 0", "file 1", "file 2", "file 3"},
		[]string{"bool", "string", "string", "file", "file", "file", "file"}, []string{"T", "*", "*", "*", "*", "*", "*"}, *getpath(), nil)
	sel.GetOpt()

	files := make([]string, 0)
	for i := 0; i < 4; i++ {
		if sel.StrRes[i+3] != "" {
			files = append(files, sel.StrRes[i+3])
		}
	}
	g4kenc.Hint = sel.StrRes[2]
	path, err := g4kenc.Encrypt(files, []byte(sel.StrRes[1]))
	if err != nil {
		return err
	}
	fmt.Printf("G4KENC Encrypt Success (%s)\n", path)
	return err
}

// g4kenc decrypt
func func44() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	path := kio.Input("HintViewer file : ")
	err = g4kenc.View(path)
	if err != nil {
		return err
	}

	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G4KENC Decrypt >", "target file", "password", "hint"}, []string{"bool", "file", "string", "string"}, []string{"T", "ote", "*", "0"}, *getpath(), nil)
	sel.StrRes = []string{"", path, "", g4kenc.Hint}
	sel.GetOpt()
	err = g4kenc.Decrypt(sel.StrRes[1], []byte(sel.StrRes[2]))
	if err != nil {
		return err
	}
	fmt.Println("G4KENC Decrypt Success")
	return err
}

// g4kaes enc all
func func45() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G4KAES.ALL Encrypt >", "target file", "password", "keyfile", "hint"}, []string{"bool", "file", "string", "keyfile", "string"}, []string{"T", "*", "*", "*", "*"}, *getpath(), legsup.G4kf())
	sel.GetOpt()

	g4kaes_all.Hint = []byte(sel.StrRes[4])
	path, err := g4kaes_all.EnFile([]byte(sel.StrRes[2]), sel.ByteRes[3], sel.StrRes[1])
	if err != nil {
		return err
	}
	fmt.Printf("G4KAES.ALL Encrypt Success (%s)\n", path)
	return err
}

// g4kaes dec all
func func46() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	path := kio.Input("HintViewer file : ")
	err = g4kaes_all.ViewFile(path)
	if err != nil {
		return err
	}

	var sel cliget.OptSel
	sel.Init([]string{"< G4KAES.ALL Decrypt >", "target file", "password", "keyfile", "hint"}, []string{"bool", "file", "string", "keyfile", "string"}, []string{"T", "png", "*", "*", "0"}, *getpath(), legsup.G4kf())
	sel.StrRes = []string{"", path, "", "", string(g4kaes_all.Hint)}
	sel.GetOpt()
	path, err = g4kaes_all.DeFile([]byte(sel.StrRes[2]), sel.ByteRes[3], sel.StrRes[1])
	if err != nil {
		return err
	}
	fmt.Printf("G4KAES.ALL Decrypt Success (%s)\n", path)
	return err
}

// g4kaes enc func
func func47() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G4KAES.FUNC Encrypt >", "target file", "key"}, []string{"bool", "file", "bytes"}, []string{"T", "*", "48"}, *getpath(), nil)
	sel.ByteRes[2] = make([]byte, 48)
	sel.GetOpt()

	err = g4kaes_func.Inbuf.OpenF(sel.StrRes[1], true)
	if err == nil {
		defer g4kaes_func.Inbuf.CloseF()
	} else {
		return err
	}
	err = g4kaes_func.Exbuf.OpenF("./_ST5_DATA/result.bin", false)
	if err == nil {
		defer g4kaes_func.Exbuf.CloseF()
	} else {
		return err
	}
	err = g4kaes_func.Encrypt(sel.ByteRes[2])
	if err != nil {
		return err
	}
	fmt.Println("G4KAES.FUNC Encrypt Success (result.bin)")
	return err
}

// g4kaes dec func
func func48() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G4KAES.FUNC Decrypt >", "target file", "key"}, []string{"bool", "file", "bytes"}, []string{"T", "*", "48"}, *getpath(), nil)
	sel.ByteRes[2] = make([]byte, 48)
	sel.GetOpt()

	err = g4kaes_func.Inbuf.OpenF(sel.StrRes[1], true)
	if err == nil {
		defer g4kaes_func.Inbuf.CloseF()
	} else {
		return err
	}
	err = g4kaes_func.Exbuf.OpenF("./_ST5_DATA/result.bin", false)
	if err == nil {
		defer g4kaes_func.Exbuf.CloseF()
	} else {
		return err
	}
	err = g4kaes_func.Decrypt(sel.ByteRes[2])
	if err != nil {
		return err
	}
	fmt.Println("G4KAES.FUNC Decrypt Success (result.bin)")
	return err
}

// g4kv enc
func func49() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	var sel cliget.OptSel
	sel.Init([]string{"< G4KV Encrypt >", "target dir", "password", "keyfile", "hint"}, []string{"bool", "folder", "string", "keyfile", "string"}, []string{"T", "NR", "*", "*", "*"}, *getpath(), legsup.G4kf())
	sel.GetOpt()

	clu := fmt.Sprintf("./_ST5_DATA/clu%d/", kobj.Decode(kaes.Genrand(2)))
	os.Mkdir(clu, os.ModePerm)
	wk := legsup.InitKV4(clu)
	wk.Hint = []byte(sel.StrRes[4])
	err = wk.Write([]byte(sel.StrRes[2]), sel.ByteRes[3], sel.StrRes[1])
	if err != nil {
		return err
	}
	fmt.Printf("G4KV Encrypt Success (%s)\n", clu[12:])
	return err
}

// g4kv dec
func func50() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	wk := legsup.InitKV4(kio.Input("HintViewer cluster : "))
	err = wk.View()
	if err != nil {
		return err
	}

	var sel cliget.OptSel
	sel.Init([]string{"< G4KV Decrypt >", "target dir", "password", "keyfile", "hint"}, []string{"bool", "folder", "string", "keyfile", "string"}, []string{"T", "NR", "*", "*", "0"}, *getpath(), legsup.G4kf())
	sel.StrRes = []string{"", wk.Path, "", "", string(wk.Hint)}
	sel.GetOpt()
	err = wk.Read([]byte(sel.StrRes[2]), sel.ByteRes[3], "./_ST5_DATA/")
	if err != nil {
		return err
	}
	fmt.Println("G4KV Decrypt Success (bin/ main/)")
	return err
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("critical : %s\n", err)
		}
		kio.Input("press ENTER to exit... ")
	}()
	kobj.Repath()
	flag := true
	options := []string{"Exit Program", "FileDiv Pack", "FileDiv Unpack", "G5KSC View", "G5KSC Pack", "G5KSC Unpack", "G5KZIP Pack", "G5KZIP Unpack", "G5KAES Encrypt All",
		"G5KAES Decrypt All", "G5KAES Encrypt Func", "G5KAES Decrypt Func", "G5KPKG View", "G5KPKG Pack", "G5KPKG Unpack", "G5KPIC Pack", "G5KPIC Unpack", "G5KPIC Restore",
		"Kscript5 Compile", "Kscript5 Assemble", "Kscript5 Reversing", "KV4adv MkNew", "KV4adv Login", "KV4adv Rebuild", "KV5st MkNew", "KV5st Login", "KV5st Rebuild",
		"G1KENC Encrypt", "G1KENC Decrypt", "G2KENC Encrypt", "G2KENC Decrypt", "G3KZIP Pack File", "G3KZIP Pack Folder", "G3KZIP Unpack", "G3KAES Encrypt All", "G3KAES Decrypt All",
		"G3KAES Encrypt Func", "G3KAES Decrypt Func", "G3KPIC Pack", "G3KPIC Unpack", "KV3st Encrypt", "KV3st Decrypt", "G3ZIP Release", "G4KENC Encrypt", "G4KENC Decrypt",
		"G4KAES Encrypt All", "G4KAES Decrypt All", "G4KAES Encrypt Func", "G4KAES Decrypt Func", "KV4st Encrypt", "KV4st Decrypt"}
	if len(options)%3 == 1 {
		options = append(append(options, "_____"), "_____")
	} else if len(options)%3 == 2 {
		options = append(options, "_____")
	}
	cfg.init()
	for flag {
		fmt.Println("\n==============================   Select Mode   ==============================")
		for i := 0; i < len(options)/3; i++ {
			fmt.Printf("[%02d] %19s   [%02d] %19s   [%02d] %19s\n", 3*i, options[3*i], 3*i+1, options[3*i+1], 3*i+2, options[3*i+2])
		}
		mode, err := strconv.Atoi(kio.Input(">>> "))
		if err != nil {
			fmt.Println("Mode Error : unknown mode")
		} else {
			switch mode {
			case 0:
				flag = false
			case 1:
				if err := func1(); err != nil {
					fmt.Println(err)
				}
			case 2:
				if err := func2(); err != nil {
					fmt.Println(err)
				}
			case 3:
				if err := func3(); err != nil {
					fmt.Println(err)
				}
			case 4:
				if err := func4(); err != nil {
					fmt.Println(err)
				}
			case 5:
				if err := func5(); err != nil {
					fmt.Println(err)
				}
			case 6:
				if err := func6(); err != nil {
					fmt.Println(err)
				}
			case 7:
				if err := func7(); err != nil {
					fmt.Println(err)
				}
			case 8:
				if err := func8(); err != nil {
					fmt.Println(err)
				}
			case 9:
				if err := func9(); err != nil {
					fmt.Println(err)
				}
			case 10:
				if err := func10(); err != nil {
					fmt.Println(err)
				}
			case 11:
				if err := func11(); err != nil {
					fmt.Println(err)
				}
			case 12:
				if err := func12(); err != nil {
					fmt.Println(err)
				}
			case 13:
				if err := func13(); err != nil {
					fmt.Println(err)
				}
			case 14:
				if err := func14(); err != nil {
					fmt.Println(err)
				}
			case 15:
				if err := func15(); err != nil {
					fmt.Println(err)
				}
			case 16:
				if err := func16(); err != nil {
					fmt.Println(err)
				}
			case 17:
				if err := func17(); err != nil {
					fmt.Println(err)
				}
			case 18:
				if err := func18(); err != nil {
					fmt.Println(err)
				}
			case 19:
				if err := func19(); err != nil {
					fmt.Println(err)
				}
			case 20:
				if err := func20(); err != nil {
					fmt.Println(err)
				}
			case 21:
				if err := func21(); err != nil {
					fmt.Println(err)
				}
			case 22:
				if err := func22(); err != nil {
					fmt.Println(err)
				}
			case 23:
				if err := func23(); err != nil {
					fmt.Println(err)
				}
			case 24:
				if err := func24(); err != nil {
					fmt.Println(err)
				}
			case 25:
				if err := func25(); err != nil {
					fmt.Println(err)
				}
			case 26:
				if err := func26(); err != nil {
					fmt.Println(err)
				}
			case 27:
				if err := func27(); err != nil {
					fmt.Println(err)
				}
			case 28:
				if err := func28(); err != nil {
					fmt.Println(err)
				}
			case 29:
				if err := func29(); err != nil {
					fmt.Println(err)
				}
			case 30:
				if err := func30(); err != nil {
					fmt.Println(err)
				}
			case 31:
				if err := func31(); err != nil {
					fmt.Println(err)
				}
			case 32:
				if err := func32(); err != nil {
					fmt.Println(err)
				}
			case 33:
				if err := func33(); err != nil {
					fmt.Println(err)
				}
			case 34:
				if err := func34(); err != nil {
					fmt.Println(err)
				}
			case 35:
				if err := func35(); err != nil {
					fmt.Println(err)
				}
			case 36:
				if err := func36(); err != nil {
					fmt.Println(err)
				}
			case 37:
				if err := func37(); err != nil {
					fmt.Println(err)
				}
			case 38:
				if err := func38(); err != nil {
					fmt.Println(err)
				}
			case 39:
				if err := func39(); err != nil {
					fmt.Println(err)
				}
			case 40:
				if err := func40(); err != nil {
					fmt.Println(err)
				}
			case 41:
				if err := func41(); err != nil {
					fmt.Println(err)
				}
			case 42:
				if err := func42(); err != nil {
					fmt.Println(err)
				}
			case 43:
				if err := func43(); err != nil {
					fmt.Println(err)
				}
			case 44:
				if err := func44(); err != nil {
					fmt.Println(err)
				}
			case 45:
				if err := func45(); err != nil {
					fmt.Println(err)
				}
			case 46:
				if err := func46(); err != nil {
					fmt.Println(err)
				}
			case 47:
				if err := func47(); err != nil {
					fmt.Println(err)
				}
			case 48:
				if err := func48(); err != nil {
					fmt.Println(err)
				}
			case 49:
				if err := func49(); err != nil {
					fmt.Println(err)
				}
			case 50:
				if err := func50(); err != nil {
					fmt.Println(err)
				}
			default:
				fmt.Println("Mode Error : unknown mode")
			}
		}
	}
}

// static data & generic functions
var cfg sdata

var g5kaes_all kaes.Allmode
var g5kaes_func kaes.Funcmode
var g5kscript langc // cplr.go from common.kscriptc
var g5shell kvdrive

var g1kenc legsup.G1enc
var g2kenc legsup.G2enc
var g3kzip legsup.G3kzip
var g3kaes_all legsup.G3kaesall
var g3kaes_func legsup.G3kaesfunc
var g3kpic legsup.G3kpic
var g3kv legsup.G3kv3
var g4kenc legsup.G4enc
var g4kaes_all legsup.G4kaesall
var g4kaes_func legsup.G4kaesfunc

// settings data
type sdata struct {
	desktop string   // desktop path
	pname   []string // public key name
	public  []string // public key data
}

// init sdata
func (tbox *sdata) init() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("No Config : %s\n", err)
		}
	}()
	if _, err := os.Stat("./_ST5_DATA/"); err != nil {
		os.Mkdir("./_ST5_DATA/", os.ModePerm)
	}
	tbox.desktop = "./_ST5_DATA/"
	wk := kdb.Initkdb()
	f, _ := kio.Open("../../_ST5_CONFIG.txt", "r")
	d, _ := kio.Read(f, -1)
	f.Close()
	wk.Read(string(d))
	v, _ := wk.Get("path.desktop")
	tbox.desktop = v.Dat6
	wk = kdb.Initkdb()
	f, _ = kio.Open("../../_ST5_SIGN.txt", "r")
	d, _ = kio.Read(f, -1)
	f.Close()
	wk.Read(string(d))
	num := 0
	tbox.pname = make([]string, 0, 10)
	tbox.public = make([]string, 0, 10)
	for {
		nm0 := fmt.Sprintf("%d.name", num)
		nm1 := fmt.Sprintf("%d.public", num)
		if _, ext := wk.Name[nm0]; ext {
			v, _ = wk.Get(nm0)
			tbox.pname = append(tbox.pname, v.Dat6)
			v, _ = wk.Get(nm1)
			tbox.public = append(tbox.public, v.Dat6)
			num = num + 1
		} else {
			break
		}
	}
}

// check publickey
func (tbox *sdata) check(public string) string {
	if public == "" {
		return "No Sign [00000000] : _"
	}
	for i, r := range tbox.public {
		if r == public {
			return fmt.Sprintf("Valid [%s] : %s", kio.Bprint(ksc.Crc32hash([]byte(r))), tbox.pname[i])
		}
	}
	return fmt.Sprintf("Untrusted [%s] : ???", kio.Bprint(ksc.Crc32hash([]byte(public))))
}

// get cliget.pathsel
func getpath() *cliget.PathSel {
	names := []string{"Current", "Desktop"}
	paths := []string{kio.Abs("./"), cfg.desktop}
	var out cliget.PathSel
	out.Init(names, paths)
	return &out
}

// get public, private
func getsign(data []byte) (ret [2]string) {
	defer func() {
		if ferr := recover(); ferr != nil {
			ret = [2]string{"", ""}
		}
	}()
	wk := kdb.Initkdb()
	if err := wk.Read(string(data)); err != nil {
		return [2]string{"", ""}
	}
	v0, _ := wk.Get("public")
	v1, _ := wk.Get("private")
	return [2]string{v0.Dat6, v1.Dat6}
}

// decode signed int 2B 4B
func kscript_dec(data []byte) int {
	if len(data) == 4 {
		return int(int32(data[0]) | int32(data[1])<<8 | int32(data[2])<<16 | int32(data[3])<<24)
	} else if len(data) == 2 {
		return int(int16(data[0]) | int16(data[1])<<8)
	} else {
		return 0
	}
}

// translate kscript bytecode (data)
func kscript_revdata(code []byte) string {
	out := ""
	pos := 0
	for pos < len(code) {
		switch code[pos] {
		case 78: // none
			out = out + "data none\n"
			pos = pos + 1
		case 66: // bool
			if code[pos+1] == 0 {
				out = out + "data bool t\n"
			} else {
				out = out + "data bool f\n"
			}
			pos = pos + 2
		case 73: // int
			temp := bytes.NewReader(code[pos+1 : pos+9])
			var tgt int64
			binary.Read(temp, binary.LittleEndian, &tgt)
			out = out + fmt.Sprintf("data int %d\n", tgt)
			pos = pos + 9
		case 70: // float
			tgt := math.Float64frombits(binary.LittleEndian.Uint64(code[pos+1 : pos+9]))
			out = out + fmt.Sprintf("data float %f\n", tgt)
			pos = pos + 9
		case 83: // string
			length := kscript_dec(code[pos+1 : pos+5])
			pos = pos + 5
			out = out + fmt.Sprintf("data string %s\n", kio.Bprint(code[pos:pos+length]))
			pos = pos + length
		case 67: // bytes
			length := kscript_dec(code[pos+1 : pos+5])
			pos = pos + 5
			out = out + fmt.Sprintf("data bytes %s\n", kio.Bprint(code[pos:pos+length]))
			pos = pos + length
		default:
			return out + "e703 : decode fail DATA\n"
		}
	}
	return out
}

// translate kscript bytecode (code)
func kscript_revcode(code []byte) string {
	out := ""
	for i := 0; i < len(code)/8; i++ {
		op := code[8*i]
		reg := "m" + string(code[8*i+1])
		i16 := kscript_dec(code[8*i+2 : 8*i+4])
		var i16_v string
		switch i16 {
		case 99:
			i16_v = "const"
		case 103:
			i16_v = "global"
		case 108:
			i16_v = "local"
		default:
			i16_v = "unknown"
		}
		i32 := kscript_dec(code[8*i+4 : 8*i+8])

		switch op {
		case 0x00:
			out = out + "hlt\n"
		case 0x01:
			out = out + fmt.Sprintf("nop ;%d\n", i)
		case 0x02:
			out = out + fmt.Sprintf("label @%d\n", i32)

		case 0x10:
			out = out + fmt.Sprintf("intr $%d $%d\n", i16, i32)
		case 0x11:
			out = out + fmt.Sprintf("call %s $%d @%d\n", reg, i16, i32)
		case 0x12:
			out = out + fmt.Sprintf("ret $%d\n", i16)
		case 0x13:
			out = out + fmt.Sprintf("jmp @%d\n", i32)
		case 0x14:
			out = out + fmt.Sprintf("jmpiff @%d\n", i32)
		case 0x15:
			out = out + fmt.Sprintf("forcond %s %s &%d\n", reg, i16_v, i32)
		case 0x16:
			out = out + fmt.Sprintf("forset %s %s &%d\n", reg, i16_v, i32)

		case 0x20:
			out = out + fmt.Sprintf("load %s %s &%d\n", reg, i16_v, i32)
		case 0x21:
			out = out + fmt.Sprintf("store %s %s &%d\n", reg, i16_v, i32)
		case 0x22:
			out = out + fmt.Sprintf("push %s\n", reg)
		case 0x23:
			out = out + fmt.Sprintf("pop %s\n", reg)
		case 0x24:
			out = out + fmt.Sprintf("pushset %s &%d\n", i16_v, i32)
		case 0x25:
			out = out + fmt.Sprintf("popset %s &%d\n", i16_v, i32)

		case 0x30:
			out = out + "add\n"
		case 0x31:
			out = out + "sub\n"
		case 0x32:
			out = out + "mul\n"
		case 0x33:
			out = out + "div\n"

		case 0x40:
			out = out + "divs\n"
		case 0x41:
			out = out + "divr\n"
		case 0x42:
			out = out + "pow\n"

		case 0x50:
			out = out + "eql\n"
		case 0x51:
			out = out + "eqln\n"
		case 0x52:
			out = out + "sml\n"
		case 0x53:
			out = out + "grt\n"
		case 0x54:
			out = out + "smle\n"
		case 0x55:
			out = out + "grte\n"

		case 0x60:
			out = out + fmt.Sprintf("inc %s &%d\n", i16_v, i32)
		case 0x61:
			out = out + fmt.Sprintf("dec %s &%d\n", i16_v, i32)
		case 0x62:
			out = out + fmt.Sprintf("shm %s &%d\n", i16_v, i32)
		case 0x63:
			out = out + fmt.Sprintf("shd %s &%d\n", i16_v, i32)

		case 0x70:
			out = out + fmt.Sprintf("addi $%d\n", i32)
		case 0x71:
			out = out + fmt.Sprintf("muli $%d\n", i32)
		case 0x72:
			out = out + fmt.Sprintf("addr %s $%d\n", reg, i32)
		case 0x73:
			out = out + fmt.Sprintf("jmpi $%d @%d\n", i16, i32)

		default:
			return out + "e705 : invalid opcode\n"
		}
	}
	return out
}

// kv4adv, kv5st shell
type kvdrive struct {
	g4sh     kvault.G4FSshell
	g5sh     kvault.Shell
	isv4     bool
	viewsize bool
}

// init struct
func (tbox *kvdrive) init(isv4 bool) {
	tbox.isv4 = isv4
	tbox.viewsize = false
	if tbox.isv4 {
		tbox.g4sh.Command("init", nil)
	} else {
		tbox.g5sh.Command("init", nil)
	}
}

// exit struct
func (tbox *kvdrive) exit() {
	if tbox.isv4 {
		tbox.g4sh.Command("exit", nil)
	} else {
		tbox.g5sh.Command("exit", nil)
	}
}

// start repl
func (tbox *kvdrive) repl(pw []byte, kf []byte, sleep string) error {
	if tbox.isv4 {
		tbox.g4sh.IObyte = [][]byte{pw, kf}
		if err := tbox.g4sh.Command("login", []string{sleep}); err != nil {
			return err
		}
	} else {
		tbox.g5sh.IObyte = [][]byte{pw, kf}
		if err := tbox.g5sh.Command("login", []string{sleep}); err != nil {
			return err
		}
	}

	cmd := ""
	var opt []string
	tbox.printpage()
	for {
		cmd, opt = tbox.getcmd()
		if err := tbox.runcmd(cmd, opt); err != nil {
			fmt.Println(err)
		}
		if cmd == "cd" {
			tbox.printpage()
		} else if cmd == "exit" {
			break
		}
	}
	return nil
}

// print page
func (tbox *kvdrive) printpage() {
	if tbox.isv4 {
		if tbox.viewsize {
			tbox.g4sh.Command("update", []string{"true"})
		} else {
			tbox.g4sh.Command("update", []string{"false"})
		}
		fmt.Printf("\nWorking : %t, Readonly : %t, [%s@%s]\n%s\n", tbox.g4sh.FlagWk, tbox.g4sh.FlagRo, *tbox.g4sh.InSys.Account, *tbox.g4sh.InSys.Cluster, tbox.g4sh.CurPath)
		num := 0
		for i, r := range tbox.g4sh.CurInfo.Dir_name {
			fmt.Printf("[%5d] %20s   %s %dB lock %t\n", num, r, tbox.g4sh.CurInfo.Dir_time[i], tbox.g4sh.CurInfo.Dir_size[i], tbox.g4sh.CurInfo.Dir_locked[i])
			num = num + 1
		}
		for i, r := range tbox.g4sh.CurInfo.File_name {
			fmt.Printf("[%5d] %20s   %s %dB fptr %d\n", num, r, tbox.g4sh.CurInfo.File_time[i], tbox.g4sh.CurInfo.File_size[i], tbox.g4sh.CurInfo.File_fptr[i])
			num = num + 1
		}
	} else {
		tbox.g5sh.FlagSz = tbox.viewsize
		tbox.g5sh.Command("update", nil)
		fmt.Printf("\nWorking : %t, Readonly : %t, [%s@%s]\n%s\n", tbox.g5sh.FlagWk, tbox.g5sh.FlagRo, *tbox.g5sh.InSys.Account, *tbox.g5sh.InSys.Cluster, tbox.g5sh.CurPath)
		for i, r := range tbox.g5sh.CurName {
			fmt.Printf("[%5d] %20s   %s %dB lock %t\n", i, r, tbox.g5sh.CurTime[i], tbox.g5sh.CurSize[i], tbox.g5sh.CurLock[i])
		}
	}
}

// get order
func (tbox *kvdrive) getcmd() (string, []string) {
	order := kio.Input(">>> ")
	if strings.Contains(order, " ") {
		pos := strings.Index(order, " ")
		return strings.ToLower(order[:pos]), strings.Split(order[pos+1:], " ")
	} else {
		return strings.ToLower(order), nil
	}
}

// run order
func (tbox *kvdrive) runcmd(order string, option []string) (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil

	switch order {
	case "cd":
		if len(option) > 0 {
			n, e := strconv.Atoi(option[0])
			if e == nil {
				err = tbox.cd(n)
			} else {
				err = e
			}
		}
	case "rm", "del":
		n, e := strconv.Atoi(option[0])
		if e == nil {
			err = tbox.rm(n)
		} else {
			err = e
		}
	case "ren":
		nm := strings.Join(option[1:], " ")
		n, e := strconv.Atoi(option[0])
		if e == nil {
			err = tbox.ren(n, nm)
		} else {
			err = e
		}
	case "mv", "move":
		ns := make([]int, 0)
		for _, r := range option {
			n, e := strconv.Atoi(r)
			if e == nil {
				ns = append(ns, n)
			} else {
				return e
			}
		}
		err = tbox.mv(ns)

	case "md", "mkdir":
		err = tbox.md(strings.Join(option, " "))
	case "touch", "mkfile":
		err = tbox.mf(strings.Join(option, " "))
	case "im", "import":
		err = tbox.im(strings.Join(option, " "))
	case "ex", "export":
		if len(option) == 0 && tbox.isv4 {
			err = tbox.g4sh.Command("exdir", []string{""})
		} else if len(option) == 0 {
			err = tbox.g5sh.Command("exdir", []string{tbox.g5sh.CurPath})
		} else {
			n, e := strconv.Atoi(option[0])
			if e == nil {
				err = tbox.ex(n)
			} else {
				err = e
			}
		}

	case "view":
		if len(option) == 0 && tbox.isv4 {
			fmt.Printf("Status KV4adv : Working %t, ViewSize %t, AsyncErr %s\n", tbox.g4sh.FlagWk, tbox.viewsize, tbox.g4sh.AsyncErr)
		} else if len(option) == 0 {
			fmt.Printf("Status KV5st : Working %t, ViewSize %t, AsyncErr %s\n", tbox.g5sh.FlagWk, tbox.viewsize, tbox.g5sh.AsyncErr)
		} else {
			err = tbox.view(option[0])
		}
	case "reset":
		err = tbox.reset(strings.ToLower(option[0]))
	case "exit":
		err = nil

	default:
		fmt.Println("unknown order : KVdrive supports following commands")
		fmt.Println("cd .num[-1, N]\nrm/del num[N]\nren num[N] name\nmv/move nums[N][...]")
		fmt.Println("md/mkdir name\ntouch/mkfile name\nim/import path\nex/export .num[N]")
		fmt.Println("view .mode[log, debug, print]\nreset mode[viewsize t/f, account]\nexit")
	}
	tbox.g4sh.IOstr = nil
	tbox.g4sh.IObyte = nil
	tbox.g5sh.IOstr = nil
	tbox.g5sh.IObyte = nil
	return err
}

func (tbox *kvdrive) getnm(pos int) string {
	if pos < len(tbox.g4sh.CurInfo.Dir_name) {
		return tbox.g4sh.CurInfo.Dir_name[pos]
	} else {
		pos = pos - len(tbox.g4sh.CurInfo.Dir_name)
		return tbox.g4sh.CurInfo.File_name[pos]
	}
}

func (tbox *kvdrive) cd(pos int) error {
	if tbox.isv4 {
		if pos < 0 {
			if strings.Count(tbox.g4sh.CurPath, "/") > 1 {
				temp := tbox.g4sh.CurPath[:len(tbox.g4sh.CurPath)-1]
				tbox.g4sh.CurPath = temp[:strings.LastIndex(temp, "/")+1]
			}
		} else if pos < len(tbox.g4sh.CurInfo.Dir_name) {
			tbox.g4sh.CurPath = tbox.g4sh.CurPath + tbox.g4sh.CurInfo.Dir_name[pos]
		} else {
			return errors.New("invalid option")
		}
	} else {
		if pos < 0 {
			if strings.Count(tbox.g5sh.CurPath, "/") > 1 {
				temp := tbox.g5sh.CurPath[:len(tbox.g5sh.CurPath)-1]
				tbox.g5sh.CurPath = temp[:strings.LastIndex(temp, "/")+1]
			}
		} else if pos < tbox.g5sh.CurNum[0] {
			tbox.g5sh.CurPath = tbox.g5sh.CurPath + tbox.g5sh.CurName[pos]
		} else {
			return errors.New("invalid option")
		}
	}
	return nil
}

func (tbox *kvdrive) rm(pos int) error {
	if tbox.isv4 {
		return tbox.g4sh.Command("delete", []string{tbox.getnm(pos)})
	} else {
		return tbox.g5sh.Command("delete", []string{fmt.Sprint(pos)})
	}
}

func (tbox *kvdrive) ren(pos int, name string) error {
	if tbox.isv4 {
		if pos < len(tbox.g4sh.CurInfo.Dir_name) && name[len(name)-1] != '/' {
			name = name + "/"
		}
		tbox.g4sh.IOstr = []string{name}
		return tbox.g4sh.Command("rename", []string{tbox.getnm(pos)})
	} else {
		if pos < tbox.g5sh.CurNum[0] && name[len(name)-1] != '/' {
			name = name + "/"
		}
		tbox.g5sh.IOstr = []string{name}
		return tbox.g5sh.Command("rename", []string{fmt.Sprint(pos)})
	}
}

func (tbox *kvdrive) mv(pos []int) error {
	var tgt string
	if tbox.isv4 {
		tgt = tbox.g4sh.CurPath
		for {
			if err := tbox.g4sh.Command("navigate", []string{tgt}); err != nil {
				return err
			}
			fmt.Println(tgt)
			names := strings.Split(tbox.g4sh.IOstr[0], "\n")
			for i, r := range names {
				fmt.Printf("[%5d] %s\n", i, r)
			}
			num, _ := strconv.Atoi(kio.Input("Move Tgt [-N, -1, N] : "))
			if num < -1 {
				break
			} else if num == -1 {
				if strings.Count(tgt, "/") > 1 {
					temp := tgt[:len(tgt)-1]
					tgt = temp[:strings.LastIndex(temp, "/")+1]
				}
			} else if num < len(names) {
				tgt = tgt + names[num]
			}
		}
	} else {
		tgt = tbox.g5sh.CurPath
		for {
			if err := tbox.g5sh.Command("navigate", []string{tgt}); err != nil {
				return err
			}
			fmt.Println(tgt)
			names := strings.Split(tbox.g5sh.IOstr[0], "\n")
			for i, r := range names {
				fmt.Printf("[%5d] %s\n", i, r)
			}
			num, _ := strconv.Atoi(kio.Input("Move Tgt [-N, -1, N] : "))
			if num < -1 {
				break
			} else if num == -1 {
				if strings.Count(tgt, "/") > 1 {
					temp := tgt[:len(tgt)-1]
					tgt = temp[:strings.LastIndex(temp, "/")+1]
				}
			} else if num < len(names) {
				tgt = tgt + names[num]
			}
		}
	}

	p := []string{tgt}
	if tbox.isv4 {
		for _, r := range pos {
			p = append(p, tbox.getnm(r))
		}
		return tbox.g4sh.Command("move", p)
	} else {
		for _, r := range pos {
			p = append(p, fmt.Sprint(r))
		}
		return tbox.g5sh.Command("move", p)
	}
}

func (tbox *kvdrive) md(name string) error {
	if name[len(name)-1] != '/' {
		name = name + "/"
	}
	if tbox.isv4 {
		return tbox.g4sh.Command("dirnew", []string{name})
	} else {
		return tbox.g5sh.Command("dirnew", []string{name})
	}
}

func (tbox *kvdrive) mf(name string) error {
	if tbox.isv4 {
		tbox.g4sh.IObyte = [][]byte{[]byte("Hello, world!")}
		return tbox.g4sh.Command("imbin", []string{name})
	} else {
		tbox.g5sh.IObyte = [][]byte{[]byte("Hello, world!")}
		return tbox.g5sh.Command("imbin", []string{name})
	}
}

func (tbox *kvdrive) im(path string) error {
	path = kio.Abs(path)
	if tbox.isv4 {
		if path[len(path)-1] == '/' {
			return tbox.g4sh.Command("imdir", []string{path})
		} else {
			return tbox.g4sh.Command("imfile", []string{path})
		}
	} else {
		if path[len(path)-1] == '/' {
			return tbox.g5sh.Command("imdir", []string{path})
		} else {
			return tbox.g5sh.Command("imfile", []string{path})
		}
	}
}

func (tbox *kvdrive) ex(pos int) error {
	if tbox.isv4 {
		if pos < len(tbox.g4sh.CurInfo.Dir_name) {
			return tbox.g4sh.Command("exdir", []string{tbox.getnm(pos)})
		} else {
			return tbox.g4sh.Command("exfile", []string{tbox.getnm(pos)})
		}
	} else {
		if pos < tbox.g5sh.CurNum[0] {
			return tbox.g5sh.Command("exdir", []string{tbox.g5sh.CurPath + tbox.g5sh.CurName[pos]})
		} else {
			return tbox.g5sh.Command("exfile", []string{tbox.g5sh.CurName[pos]})
		}
	}
}

func (tbox *kvdrive) view(mode string) error {
	switch mode {
	case "log":
		if tbox.isv4 {
			defer tbox.g4sh.Command("log", []string{"true"})
			err := tbox.g4sh.Command("log", []string{"false"})
			fmt.Println(tbox.g4sh.IOstr[0])
			return err
		} else {
			defer tbox.g5sh.Command("log", []string{"true"})
			err := tbox.g5sh.Command("log", []string{"false"})
			fmt.Println(tbox.g5sh.IOstr[0])
			return err
		}
	case "debug":
		if tbox.isv4 {
			err := tbox.g4sh.Command("debug", []string{"true"})
			fmt.Println(strings.Join(tbox.g4sh.IOstr, "\n"))
			return err
		} else {
			err := tbox.g5sh.Command("debug", []string{"true"})
			fmt.Println(strings.Join(tbox.g5sh.IOstr, "\n"))
			return err
		}
	case "print":
		if tbox.isv4 {
			err := tbox.g4sh.Command("print", []string{"true"})
			fmt.Println(tbox.g4sh.IOstr[0])
			return err
		} else {
			err := tbox.g5sh.Command("print", []string{"true"})
			fmt.Println(tbox.g5sh.IOstr[0])
			return err
		}
	default:
		return errors.New("invalid option")
	}
}

func (tbox *kvdrive) reset(mode string) error {
	switch mode {
	case "t", "true", "on":
		tbox.viewsize = true
		return nil
	case "f", "false", "off":
		tbox.viewsize = false
		return nil
	case "acc", "account", "pw", "pwkf":
		var sel cliget.OptSel
		sel.Init([]string{"< PWKF Reset >", "password", "keyfile", "hint"}, []string{"bool", "string", "keyfile", "string"}, []string{"T", "*", "*", "*"}, *getpath(), kaes.Basickey())
		sel.GetOpt()
		if tbox.isv4 {
			tbox.g4sh.IObyte = [][]byte{[]byte(sel.StrRes[1]), sel.ByteRes[2], []byte(sel.StrRes[3])}
			return tbox.g4sh.Command("reset", nil)
		} else {
			tbox.g5sh.IObyte = [][]byte{[]byte(sel.StrRes[1]), sel.ByteRes[2], []byte(sel.StrRes[3])}
			return tbox.g5sh.Command("reset", nil)
		}
	default:
		return errors.New("invalid option")
	}
}
