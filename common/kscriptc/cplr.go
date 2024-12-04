// test722 : common.kscriptc compiler

package main

import (
	"errors"
	"fmt"
	"stdlib5/kdb"
	"stdlib5/kio"
	"stdlib5/kscript"
	"strconv"
	"strings"
	"time"
)

// compiler worker
type langc struct {
	option    [5]bool   // -token -ast -asm -optconst -optasm
	data      [7]string // source, info, iconpath, signkey, token, ast, asm
	parser    kscript.Parser
	compiler  kscript.Compiler
	assmebler kscript.Assembler
}

// init langc
func (tbox *langc) init() {
	tbox.option = [5]bool{false, false, false, false, false}
	tbox.data = [7]string{"", "", "", "", "", "", ""}
	tbox.parser.Init()
	tbox.parser.Type_Function = make([]string, 0)
	tbox.compiler.Init()
	tbox.compiler.OuterNum = make(map[string]int)
	tbox.compiler.OuterParms = make(map[string]int)
	tbox.assmebler.SetKey("", "", "")
	tbox.assmebler.Info = ""
	tbox.assmebler.ABIf = 0
}

// add kspkg
func (tbox *langc) addpkg(path string) error {
	var temp kspkg
	err := temp.read(path)
	if err != nil {
		return err
	}
	tbox.assmebler.ABIf = tbox.assmebler.ABIf + temp.pkgcode
	for i, r := range temp.funcname {
		tbox.parser.Type_Function = append(tbox.parser.Type_Function, r)
		if _, ext := tbox.compiler.OuterNum[r]; ext {
			return fmt.Errorf("double define : %s", r)
		}
		tbox.compiler.OuterNum[r] = temp.intrnum[i]
		tbox.compiler.OuterParms[r] = temp.parmsnum[i]
	}
	return nil
}

// compile txt
func (tbox *langc) compile() (exe []byte, tpass float64, err error) {
	tstamp := time.Now().UnixMicro()
	err = nil
	tpass = 0.0
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
		tpass = float64(time.Now().UnixMicro()-tstamp) / 1000000
	}()

	err = tbox.parser.Split(tbox.data[0])
	if err != nil {
		return nil, tpass, err
	}
	err = tbox.parser.Parse()
	if err != nil {
		return nil, tpass, err
	}
	token, err := tbox.parser.Structify()
	if tbox.option[0] {
		temp := make([]string, 0)
		for _, r := range token {
			temp = append(temp, r.Write(0))
		}
		tbox.data[4] = strings.Join(temp, "\n")
	}
	if err != nil {
		return nil, tpass, err
	}
	fns, pgr, err := tbox.parser.GenAST(token)
	if tbox.option[1] {
		temp := make([]string, 0)
		for _, r := range fns {
			temp = append(temp, r.Write(0))
		}
		temp = append(temp, pgr.Write(0))
		tbox.data[5] = strings.Join(temp, "\n")
	}
	if err != nil {
		return nil, tpass, err
	}

	tbox.compiler.OptConst = tbox.option[3]
	tbox.compiler.OptAsm = tbox.option[4]
	asm, err := tbox.compiler.Compile(&pgr, fns)
	if tbox.option[2] {
		tbox.data[6] = asm
	}
	if err != nil {
		return nil, tpass, err
	}

	if tbox.data[3] == "" {
		tbox.assmebler.SetKey(tbox.data[2], "", "")
	} else {
		worker := kdb.Initkdb()
		err = worker.Read(tbox.data[3])
		if err != nil {
			return nil, tpass, err
		}
		tv, _ := worker.Get("private")
		private := tv.Dat6
		tv, _ = worker.Get("public")
		public := tv.Dat6
		tbox.assmebler.SetKey(tbox.data[2], public, private)
	}
	tbox.assmebler.Info = tbox.data[1]
	exe, err = tbox.assmebler.GenExe(asm)
	if err != nil {
		return nil, tpass, err
	}
	return exe, tpass, err
}

// kscript extension
type kspkg struct {
	name     string   // package name
	pkgcode  int      // package code (2**N)
	funcname []string // function name
	parmsnum []int    // parms num
	intrnum  []int    // interupt num
}

// read kscript extension txt
func (tbox *kspkg) read(path string) error {
	tbox.funcname = make([]string, 0)
	tbox.parmsnum = make([]int, 0)
	tbox.intrnum = make([]int, 0)
	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	d, err := kio.Read(f, -1)
	if err != nil {
		return err
	}
	raw := string(d)
	idx := strings.Index(raw, ">")
	tbox.name = raw[strings.Index(raw, "<")+1 : idx]
	idx = strings.Index(raw, "]")
	temp, err := strconv.Atoi(raw[strings.Index(raw, "[")+1 : idx])
	if err == nil {
		tbox.pkgcode = temp
	} else {
		return err
	}
	raw = raw[idx+1:]

	for strings.Contains(raw, "[") {
		switch raw[0] {
		case '{':
			idx = strings.Index(raw, "}")
			tbox.funcname = append(tbox.funcname, raw[1:idx])
			raw = raw[idx+1:]
		case '(':
			idx = strings.Index(raw, ")")
			temp := strings.ReplaceAll(raw[1:idx], " ", "")
			if len(temp) == 0 {
				tbox.parmsnum = append(tbox.parmsnum, 0)
			} else {
				tbox.parmsnum = append(tbox.parmsnum, strings.Count(temp, ",")+1)
			}
			raw = raw[idx+1:]
		case '[':
			idx = strings.Index(raw, "]")
			temp, err = strconv.Atoi(raw[1:idx])
			if err == nil {
				tbox.intrnum = append(tbox.intrnum, temp)
			} else {
				return err
			}
			raw = raw[idx+1:]
		case '#':
			raw = raw[strings.Index(raw, "\n")+1:]
		default:
			raw = raw[1:]
		}
	}

	if len(tbox.funcname) != len(tbox.parmsnum) || len(tbox.funcname) != len(tbox.intrnum) {
		return errors.New("invalid kscript extension")
	}
	return nil
}
