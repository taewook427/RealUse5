// ./govm.exe txt|program
// -sign -info -o1 -o2 -o3 : runone -> errhlt -> safemem
// -sign -asm -o1 -o2 : optconst -> optasm

package main

// go build -ldflags="-s -w" -trimpath x.go

import (
	"fmt"
	"path"
	"stdlib5/kio"
	"stdlib5/kobj"
	"stdlib5/ksc"
	"stdlib5/kscript"
	"stdlib5/ksign"
	"strings"
	"time"
)

func compile(file string, sign bool, asm bool, opt int) string {
	f, err := kio.Open(file, "r")
	if err != nil {
		return fmt.Sprintf("err : %s", err)
	}
	data, err := kio.Read(f, -1)
	if err != nil {
		return fmt.Sprintf("err : %s", err)
	}
	f.Close()

	t := time.Now()
	var w0 kscript.Parser
	w0.Init()
	err = w0.Split(string(data))
	if err != nil {
		return fmt.Sprintf("err : %s", err)
	}
	err = w0.Parse()
	if err != nil {
		return fmt.Sprintf("err : %s", err)
	}
	tks, err := w0.Structify()
	if err != nil {
		return fmt.Sprintf("err : %s", err)
	}
	fs, pg, err := w0.GenAST(tks)
	if err != nil {
		return fmt.Sprintf("err : %s", err)
	}

	var w1 kscript.Compiler
	w1.Init()
	w1.OptConst = opt > 0
	w1.OptAsm = opt > 1
	kasm, err := w1.Compile(&pg, fs)
	if err != nil {
		return fmt.Sprintf("err : %s", err)
	}

	var w2 kscript.Assembler
	w2.Icon = ksc.Webpbase()
	w2.Info = "테스트용 파일입니다."
	w2.ABIf = -1
	if sign {
		pub, pri, _ := ksign.Genkey(2048)
		w2.SetKey("", pub, pri)
	}
	elf, err := w2.GenExe(kasm)
	if err != nil {
		return fmt.Sprintf("err : %s", err)
	}
	dt := time.Since(t)

	if asm {
		f, _ := kio.Open("asm.txt", "w")
		kio.Write(f, []byte(strings.ToUpper(kasm)))
		f.Close()
	}
	f, _ = kio.Open("elf.webp", "w")
	kio.Write(f, elf)
	f.Close()
	return fmt.Sprintf("complete %s with (time %fs, asm size %d, bin size %d)", file, float64(dt.Microseconds())/1000000, len(kasm), len(elf))
}

func run(file string, sign bool, info bool, opt int) string {
	var w3 kscript.KVM
	w3.Init()
	w3.RunOne = opt < 1
	w3.ErrHlt = opt < 2
	w3.SafeMem = opt < 3

	a, b, c, err := w3.View(file)
	if err != nil {
		return fmt.Sprintf("err : %s", err)
	}
	if info {
		return fmt.Sprintf("info : %s\nabi : %d\npublic : %s", a, b, c)
	}
	err = w3.Load(sign)
	if err != nil {
		return fmt.Sprintf("err : %s", err)
	}

	ix := 0
	for ix == 0 {
		ix = w3.Run()
		if ix > 8 {
			w3.SetRet(kscript.TestIO(ix, w3.CallMem))
			ix = 0
		}
	}
	return fmt.Sprintf("exit code : %d with %s", ix, w3.ErrMsg)
}

func main() {
	order := kobj.Repath()
	file := ""
	sign := false
	info := false
	asm := false
	opt := 0

	for _, r := range order[1:] {
		if r[0] == '-' {
			switch strings.ToLower(r) {
			case "-sign":
				sign = true
			case "-info":
				info = true
			case "-asm":
				asm = true
			case "-o1":
				opt = 1
			case "-o2":
				opt = 2
			case "-o3":
				opt = 3
			}
		} else {
			file = r
		}
	}

	ret := ""
	if strings.ToLower(path.Ext(file)) == ".txt" {
		ret = compile(file, sign, asm, opt)
	} else {
		ret = run(file, sign, info, opt)
	}
	fmt.Println("\n" + ret)
	kio.Input("press ENTER to exit... ")
}
