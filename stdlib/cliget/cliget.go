// test698 : stdlib5.cliget

package cliget

import (
	"fmt"
	"os"
	"stdlib5/kaes"
	"stdlib5/kcom"
	"stdlib5/kio"
	"strconv"
	"strings"
)

// folder / file selector
type PathSel struct {
	tpname []string // tp name (1+ len)
	tppath []string // tp folder path (~/)

	curpath string   // current view path (~/)
	modeext string   // extension selection mode
	folder  []string // subdir names
	file    []string // subfile names
}

// disassemble order, returns only 6 possible value -> (opcode, operand)
func (tbox *PathSel) interpret() (opcode string, operand string) {
	defer func() {
		if err := recover(); err != nil {
			opcode = "nop"
			operand = ""
		}
	}()
	raw := kio.Input(">>> ")
	if raw == "!STOP" {
		return "nop", "stop"
	}
	// tp N, num N, sel N, only S, str S, nop S
	if raw[0] == '"' {
		raw = raw[1:]
		if raw[len(raw)-1] == '"' {
			raw = raw[0 : len(raw)-1]
		}
	}
	if tnum, terr := strconv.Atoi(raw); terr == nil { // pure number
		return "num", fmt.Sprint(tnum)
	} else if _, terr := os.Stat(raw); terr == nil { // existing file/dir
		return "str", kio.Abs(raw)
	} else {
		temp := strings.Split(raw, " ")
		switch temp[0] {
		case "tp":
			tnum, terr := strconv.Atoi(temp[1])
			if terr == nil {
				return "tp", fmt.Sprint(tnum)
			} else {
				return "nop", ""
			}
		case "sel":
			tnum, terr := strconv.Atoi(temp[1])
			if terr == nil {
				return "sel", fmt.Sprint(tnum)
			} else {
				return "nop", ""
			}
		case "only":
			return "only", temp[1]
		default:
			return "nop", ""
		}
	}
}

// print one page (output msg + tp list + cur dir + sel mode + folder/files)
func (tbox *PathSel) printpage(output string) {
	fmt.Print("\033[2J")
	fmt.Println(strings.Repeat("=", 32) + " PAGE START " + strings.Repeat("=", 32))
	defer fmt.Println(strings.Repeat("=", 32) + "  PAGE END  " + strings.Repeat("=", 32))
	fmt.Println("MSG : " + output)
	fmt.Print("TP : ")
	for i, r := range tbox.tpname {
		fmt.Printf("%d %s  ", i, r)
	}
	fmt.Println("\nCURPATH : " + tbox.curpath)
	fmt.Println("MODE : " + tbox.modeext)

	tdir, _ := os.ReadDir(tbox.curpath)
	tbox.folder = make([]string, 1)
	tbox.file = make([]string, 0)
	tbox.folder[0] = "../"
	for _, r := range tdir {
		tnm := r.Name()
		if r.IsDir() {
			if tnm[len(tnm)-1] != '/' {
				tnm = tnm + "/"
			}
			tbox.folder = append(tbox.folder, tnm)
		} else {
			if tbox.modeext == "*" {
				tbox.file = append(tbox.file, tnm)
			} else if strings.Contains(tnm, ".") {
				if tbox.modeext == strings.ToLower(tnm[strings.LastIndex(tnm, ".")+1:]) {
					tbox.file = append(tbox.file, tnm)
				}
			}
		}
	}

	count := 0
	for _, r := range tbox.folder {
		fmt.Printf("%04d  %s\n", count, r)
		count = count + 1
	}
	for _, r := range tbox.file {
		fmt.Printf("%04d  %s\n", count, r)
		count = count + 1
	}
}

// init selector
func (tbox *PathSel) Init(names []string, paths []string) {
	var num int
	if len(names) > len(paths) {
		num = len(paths)
	} else {
		num = len(names)
	}
	tbox.tpname = make([]string, num)
	tbox.tppath = make([]string, num)
	for i := 0; i < num; i++ {
		tbox.tpname[i] = names[i]
		tbox.tppath[i] = kio.Abs(paths[i])
	}
	tbox.curpath = kio.Abs("./")
	tbox.modeext = "*"
	tbox.folder = nil
	tbox.file = nil
}

// get one file with start dir init
func (tbox *PathSel) GetFile(start string) string {
	tbox.curpath = kio.Abs(start)
	target := ""                                            // selected path
	output := "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)" // msg out

	for target == "" {
		tbox.printpage(output)
		opcode, operand := tbox.interpret()

		switch opcode {
		case "tp": // tp curpath
			tnum, _ := strconv.Atoi(operand)
			if 0 <= tnum && tnum < len(tbox.tpname) {
				tbox.curpath = tbox.tppath[tnum]
				output = "changed current directory"
			} else {
				output = fmt.Sprintf("!wrong_teleport_number : %d", tnum)
			}

		case "num": // select or move
			tnum, _ := strconv.Atoi(operand)
			if tnum == 0 && strings.Count(tbox.curpath, "/") > 1 {
				tbox.curpath = tbox.curpath[0 : strings.LastIndex(tbox.curpath[0:len(tbox.curpath)-1], "/")+1]
				output = "moved to parent directory"
			} else if 0 < tnum && tnum < len(tbox.folder) {
				tbox.curpath = tbox.curpath + tbox.folder[tnum]
				output = "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)"
			} else if 0 < tnum && tnum < len(tbox.folder)+len(tbox.file) {
				target = tbox.curpath + tbox.file[tnum-len(tbox.folder)]
				output = fmt.Sprintf("!file_selected : %s", target)
			} else {
				output = fmt.Sprintf("!wrong_selection_number : %d", tnum)
			}

		case "sel": // select file
			tnum, _ := strconv.Atoi(operand)
			if 0 <= tnum && tnum < len(tbox.folder) {
				output = fmt.Sprintf("!cannot_select_folder : %d", tnum)
			} else if 0 < tnum && tnum < len(tbox.folder)+len(tbox.file) {
				target = tbox.curpath + tbox.file[tnum-len(tbox.folder)]
				output = fmt.Sprintf("!file_selected : %s", target)
			} else {
				output = fmt.Sprintf("!wrong_selection_number : %d", tnum)
			}

		case "only": // set selection mode
			tbox.modeext = strings.ToLower(operand)
			output = "to return to full selection, enter (only *)"

		case "str": // direct selection
			if operand[len(operand)-1] == '/' {
				output = fmt.Sprintf("!cannot_select_folder : %s", operand)
			} else {
				target = operand
				output = fmt.Sprintf("!file_selected : %s", target)
			}

		case "nop": // no op
			if operand == "stop" {
				return ""
			} else {
				output = "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)"
			}

		default: // cmd error
			output = "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)"
		}
	}
	fmt.Println("EXIT : " + output)
	return target
}

// get one folder with start dir init
func (tbox *PathSel) GetFolder(start string) string {
	tbox.curpath = kio.Abs(start)
	tbox.modeext = "**"
	target := ""                                  // selected path
	output := "OPTION  (tp N)  (N)  (sel N)  (S)" // msg out

	for target == "" {
		tbox.printpage(output)
		opcode, operand := tbox.interpret()

		switch opcode {
		case "tp": // tp curpath
			tnum, _ := strconv.Atoi(operand)
			if 0 <= tnum && tnum < len(tbox.tpname) {
				tbox.curpath = tbox.tppath[tnum]
				output = "changed current directory"
			} else {
				output = fmt.Sprintf("!wrong_teleport_number : %d", tnum)
			}

		case "num": // move
			tnum, _ := strconv.Atoi(operand)
			if tnum == 0 && strings.Count(tbox.curpath, "/") > 1 {
				tbox.curpath = tbox.curpath[0 : strings.LastIndex(tbox.curpath[0:len(tbox.curpath)-1], "/")+1]
				output = "moved to parent directory"
			} else if 0 < tnum && tnum < len(tbox.folder) {
				tbox.curpath = tbox.curpath + tbox.folder[tnum]
				output = "OPTION  (tp N)  (N)  (sel N)  (S)"
			} else {
				output = fmt.Sprintf("!wrong_selection_number : %d", tnum)
			}

		case "sel": // select folder
			tnum, _ := strconv.Atoi(operand)
			if tnum == 0 {
				target = tbox.curpath
				output = fmt.Sprintf("!folder_selected : %s", target)
			} else if 0 < tnum && tnum < len(tbox.folder) {
				target = tbox.curpath + tbox.folder[tnum]
				output = fmt.Sprintf("!folder_selected : %s", target)
			} else {
				output = fmt.Sprintf("!wrong_selection_number : %d", tnum)
			}

		case "str": // direct selection
			if operand[len(operand)-1] == '/' {
				target = operand
				output = fmt.Sprintf("!folder_selected : %s", target)
			} else {
				output = fmt.Sprintf("!cannot_select_file : %s", operand)
			}

		case "nop": // no op
			if operand == "stop" {
				return ""
			} else {
				output = "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)"
			}

		default: // cmd error
			output = "OPTION  (tp N)  (N)  (sel N)  (S)"
		}
	}
	fmt.Println("EXIT : " + output)
	return target
}

// keyfile selector
type KeySel struct {
	Data []byte // selection result
	Path string // keyfile path (bkf/direct/com)

	explorer PathSel // folder/file selector
	basic    []byte  // basic key file
}

// disassemble order, returns only 5 possible value -> (opcode, operand)
func (tbox *KeySel) interpret() (opcode string, operand string) {
	defer func() {
		if err := recover(); err != nil {
			opcode = "nop"
			operand = ""
		}
	}()
	raw := kio.Input(">>> ")
	if raw == "!STOP" {
		return "nop", "stop"
	}
	// submit, basic, direct, comm S
	if raw == "submit" {
		return "submit", ""
	} else if raw == "basic" {
		return "basic", ""
	} else if raw == "direct" {
		return "direct", ""
	} else {
		return "comm", raw[strings.Index(raw, " ")+1:]
	}
}

// print one page (output msg + Path + DataLen)
func (tbox *KeySel) printpage(output string) {
	fmt.Println("\n\n" + strings.Repeat("=", 32) + " PAGE START " + strings.Repeat("=", 32))
	defer fmt.Println(strings.Repeat("=", 32) + "  PAGE END  " + strings.Repeat("=", 32))
	fmt.Printf("MSG : %s\nPATH : %s\nLENGTH : %d\n", output, tbox.Path, len(tbox.Data))
}

// read & decrypt kcom reciever
func (tbox *KeySel) kcomread(address string) (output string) {
	defer func() {
		if ferr := recover(); ferr != nil {
			tbox.Data = nil
			tbox.Path = "[NONE]"
			output = fmt.Sprintf("!kcom_err : %s", ferr)
		}
	}()

	// unpack address
	port, key, err := kcom.Unpack(address)
	if err != nil {
		return fmt.Sprintf("!kcom_err : %s", err)
	}
	worker0 := kcom.Initcom()
	worker0.Port = port
	worker0.Close = 15

	// recieve 48B kaes key + nB path
	data, err := worker0.Recieve(key)
	if err != nil {
		return fmt.Sprintf("!kcom_err : %s", err)
	}

	// decrypt kaes.funcmode
	var worker1 kaes.Funcmode
	err = worker1.Before.Open(string(data[48:]), true)
	defer worker1.Before.Close()
	if err != nil {
		return fmt.Sprintf("!kcom_err : %s", err)
	}
	worker1.After.Open(make([]byte, 0, 1024), false)
	err = worker1.Decrypt(data[0:48])
	if err != nil {
		return fmt.Sprintf("!kcom_err : %s", err)
	}
	tbox.Data = worker1.After.Close()
	tbox.Path = fmt.Sprintf("[COMM] %s", string(data[48:]))
	return "!kcom_file_selection"
}

// init selector
func (tbox *KeySel) Init(basic []byte, explorer PathSel) {
	tbox.Data = basic
	tbox.Path = "[BASIC]"
	tbox.basic = basic
	tbox.explorer = explorer
}

// get keyfile
func (tbox *KeySel) GetKey() []byte {
	opcode := "nop"
	operand := ""
	output := "OPTION  (submit)  (basic)  (direct)  (comm S)"
	for opcode != "submit" {
		tbox.printpage(output)
		opcode, operand = tbox.interpret()

		switch opcode {
		case "submit": // submit keyfile
			output = fmt.Sprintf("!selected : %s size %d", tbox.Path, len(tbox.Data))

		case "basic": // set basic keyfile
			output = "!basic_keyfile_selected"
			tbox.Data = tbox.basic
			tbox.Path = "[BASIC]"

		case "direct": // direct file selection
			output = "!direct_file_selection"
			tbox.Path = tbox.explorer.GetFile("./")
			f, _ := kio.Open(tbox.Path, "r")
			tbox.Data, _ = kio.Read(f, -1)
			f.Close()
			tbox.Path = "[DIRECT]  " + tbox.Path

		case "comm": // kcom read
			output = tbox.kcomread(operand)

		case "nop": // no op
			if operand == "stop" {
				opcode = "submit"
			} else {
				output = "OPTION  (submit)  (basic)  (direct)  (comm S)"
			}

		default: // cmd error
			output = "OPTION  (submit)  (basic)  (direct)  (comm S)"
		}
	}
	fmt.Println("EXIT : " + output)
	return tbox.Data
}

// multiple option selector
type OptSel struct {
	Name    []string // option names
	StrRes  []string // selected string
	ByteRes [][]byte // selected bytes

	tplimit   []int    // option type
	lenlimit  []string // option limit
	explorer  PathSel  // folder/file selector
	navigator KeySel   // keyfile selector
}

// disassemble order, returns finite option pos -> (option pos, content)
func (tbox *OptSel) interpret() (pos int, content string) {
	defer func() {
		if err := recover(); err != nil {
			pos = -1
			content = ""
		}
	}()
	raw := kio.Input(">>> ")
	if raw == "!STOP" {
		return -1, "stop"
	} else if raw == "submit" {
		return -1, "submit"
	}

	// format : N%sS or submit
	var err error
	if !strings.Contains(raw, " ") {
		raw = raw + " "
	}
	pos = strings.Index(raw, " ")
	content = raw[pos+1:]
	pos, err = strconv.Atoi(raw[0:pos])
	if err != nil {
		pos = -1
		content = ""
	}
	return pos, content
}

// print one page (output msg + [name, type, limit, result] *n)
func (tbox *OptSel) printpage(output string) {
	fmt.Print("\033[2J")
	fmt.Println(strings.Repeat("=", 32) + " PAGE START " + strings.Repeat("=", 32))
	defer fmt.Println(strings.Repeat("=", 32) + "  PAGE END  " + strings.Repeat("=", 32))
	fmt.Println("MSG : " + output)

	for i, r := range tbox.Name {
		switch tbox.tplimit[i] {
		case 0:
			fmt.Printf("%03d  %s  (BOOL, %s)  %t\n", i, r, tbox.lenlimit[i], tbox.ByteRes[i][0] == 0)
		case 1:
			fmt.Printf("%03d  %s  (INT, %s)  %s\n", i, r, tbox.lenlimit[i], tbox.StrRes[i])
		case 2:
			fmt.Printf("%03d  %s  (FLOAT, %s)  %s\n", i, r, tbox.lenlimit[i], tbox.StrRes[i])
		case 3:
			fmt.Printf("%03d  %s  (STRING, %s)  %s\n", i, r, tbox.lenlimit[i], tbox.StrRes[i])
		case 4:
			if len(tbox.ByteRes[i]) < 64 {
				fmt.Printf("%03d  %s  (BYTES, %s)  %s\n", i, r, tbox.lenlimit[i], kio.Bprint(tbox.ByteRes[i]))
			} else {
				fmt.Printf("%03d  %s  (BYTES, %s)  LEN %d\n", i, r, tbox.lenlimit[i], len(tbox.ByteRes[i]))
			}
		case 5:
			fmt.Printf("%03d  %s  (FOLDER, %s)  %s\n", i, r, tbox.lenlimit[i], tbox.StrRes[i])
		case 6:
			fmt.Printf("%03d  %s  (FILE, %s)  %s\n", i, r, tbox.lenlimit[i], tbox.StrRes[i])
		case 7:
			fmt.Printf("%03d  %s  (KEYFILE, %s)  LEN %d\n", i, r, tbox.lenlimit[i], len(tbox.ByteRes[i]))
		}
	}
}

// check if valid input & update value
func (tbox *OptSel) checkinput(pos int, content string) (output string) {
	defer func() {
		if ferr := recover(); ferr != nil {
			output = fmt.Sprintf("!parse_error : %s", ferr)
		}
	}()
	if pos < 0 || len(tbox.Name) <= pos {
		return "!invalid_pos_number"
	}

	var err error
	passed := false
	switch tbox.tplimit[pos] {
	case 0: // bool (* T F)
		var temp bool
		if content == "t" || content == "T" || content == "true" || content == "True" || content == "TRUE" {
			temp = true
		} else {
			temp = false
		}
		switch tbox.lenlimit[pos] {
		case "T":
			passed = temp
		case "F":
			passed = !temp
		default:
			passed = true
		}
		if passed {
			output = "!bool_value_set"
			if temp {
				tbox.ByteRes[pos][0] = 0
			} else {
				tbox.ByteRes[pos][0] = 1
			}
		} else {
			output = "!bool_value_wrong"
		}

	case 1: // int (* 0 + 0+ - 0-)
		var temp int
		temp, err = strconv.Atoi(content)
		if err != nil {
			return "!parse_error : int"
		}
		switch tbox.lenlimit[pos] {
		case "0":
			passed = (temp == 0)
		case "0+":
			passed = (temp >= 0)
		case "0-":
			passed = (temp <= 0)
		case "+":
			passed = (temp > 0)
		case "-":
			passed = (temp < 0)
		default:
			passed = true
		}
		if passed {
			output = "!int_value_set"
			tbox.StrRes[pos] = content
		} else {
			output = "!int_value_wrong"
		}

	case 2: // float (* 0 + 0+ - 0-)
		var temp float64
		temp, err = strconv.ParseFloat(content, 64)
		if err != nil {
			return "!parse_error : float"
		}
		switch tbox.lenlimit[pos] {
		case "0":
			passed = (temp > -0.001) && (temp < 0.001)
		case "0+":
			passed = (temp >= 0.0)
		case "0-":
			passed = (temp <= 0.0)
		case "+":
			passed = (temp > 0.0)
		case "-":
			passed = (temp < 0.0)
		default:
			passed = true
		}
		if passed {
			output = "!float_value_set"
			tbox.StrRes[pos] = content
		} else {
			output = "!float_value_wrong"
		}

	case 3: // string (* N N+ N-)
		limit := tbox.lenlimit[pos]
		length := len(content)
		switch {
		case limit == "*":
			passed = true
		case strings.Contains(limit, "+"):
			num, _ := strconv.Atoi(limit[0 : len(limit)-1])
			passed = (length >= num)
		case strings.Contains(limit, "-"):
			num, _ := strconv.Atoi(limit[0 : len(limit)-1])
			passed = (length <= num)
		default:
			num, _ := strconv.Atoi(limit)
			passed = (length == num)
		}
		if passed {
			output = "!string_value_set"
			tbox.StrRes[pos] = content
		} else {
			output = "!string_value_wrong"
		}

	case 4: // bytes (* N N+ N-)
		var temp []byte
		temp, err = kio.Bread(content)
		if err != nil {
			return "!parse_error : bytes"
		}
		limit := tbox.lenlimit[pos]
		length := len(temp)
		switch {
		case limit == "*":
			passed = true
		case strings.Contains(limit, "+"):
			num, _ := strconv.Atoi(limit[0 : len(limit)-1])
			passed = (length >= num)
		case strings.Contains(limit, "-"):
			num, _ := strconv.Atoi(limit[0 : len(limit)-1])
			passed = (length <= num)
		default:
			num, _ := strconv.Atoi(limit)
			passed = (length == num)
		}
		if passed {
			output = "!bytes_value_set"
			tbox.ByteRes[pos] = temp
		} else {
			output = "!bytes_value_wrong"
		}

	case 5: // folder (* R NR)
		temp := tbox.explorer.GetFolder("./")
		switch tbox.lenlimit[pos] {
		case "R":
			passed = (strings.Count(temp, "/") == 1)
		case "NR":
			passed = (strings.Count(temp, "/") > 1)
		default:
			passed = true
		}
		if passed {
			output = "!folder_value_set"
			tbox.StrRes[pos] = temp
		} else {
			output = "!folder_value_wrong"
		}

	case 6: // file (* ext)
		temp := tbox.explorer.GetFile("./")
		limit := tbox.lenlimit[pos]
		var ext string
		if strings.Contains(temp, ".") {
			ext = strings.ToLower(temp[strings.LastIndex(temp, ".")+1:])
		} else {
			ext = ""
		}
		switch limit {
		case "*":
			passed = true
		default:
			passed = (ext == limit)
		}
		if passed {
			output = "!file_value_set"
			tbox.StrRes[pos] = temp
		} else {
			output = "!file_value_wrong"
		}

	case 7: // keyfile (* N N+ N-)
		temp := tbox.navigator.GetKey()
		limit := tbox.lenlimit[pos]
		length := len(temp)
		switch {
		case limit == "*":
			passed = true
		case strings.Contains(limit, "+"):
			num, _ := strconv.Atoi(limit[0 : len(limit)-1])
			passed = (length >= num)
		case strings.Contains(limit, "-"):
			num, _ := strconv.Atoi(limit[0 : len(limit)-1])
			passed = (length <= num)
		default:
			num, _ := strconv.Atoi(limit)
			passed = (length == num)
		}
		if passed {
			output = "!keyfile_value_set"
			tbox.ByteRes[pos] = temp
		} else {
			output = "!keyfile_value_wrong"
		}

	default:
		output = "!wrong_order"
	}
	return output
}

// init selector, (type : bool int float string bytes folder file keyfile)
func (tbox *OptSel) Init(names []string, types []string, limits []string, explorer PathSel, basic []byte) {
	tbox.explorer = explorer
	tbox.navigator.Init(basic, explorer)
	length := len(names)
	tbox.Name = make([]string, length)
	tbox.StrRes = make([]string, length)
	tbox.ByteRes = make([][]byte, length)
	tbox.tplimit = make([]int, length)
	tbox.lenlimit = make([]string, length)

	for i, r := range names {
		// set name field
		if len(r) == 0 {
			tbox.Name[i] = "_"
		} else {
			tbox.Name[i] = r
		}

		// set type, result field
		switch types[i] {
		case "bool": // result at B
			tbox.tplimit[i] = 0
			tbox.ByteRes[i] = []byte{0} // 0 T 1 F
		case "int": // result at B
			tbox.tplimit[i] = 1
			tbox.StrRes[i] = "0" // str-stored int
		case "float": // result at S
			tbox.tplimit[i] = 2
			tbox.StrRes[i] = "0.0"
		case "string": // result at S
			tbox.tplimit[i] = 3
			tbox.StrRes[i] = ""
		case "bytes": // result at B
			tbox.tplimit[i] = 4
			tbox.ByteRes[i] = nil
		case "folder": // result at S
			tbox.tplimit[i] = 5
			tbox.StrRes[i] = "./"
		case "file": // result at S
			tbox.tplimit[i] = 6
			tbox.StrRes[i] = ""
		case "keyfile": // result at B
			tbox.tplimit[i] = 7
			tbox.ByteRes[i] = basic
		default:
			tbox.tplimit[i] = 0
			tbox.ByteRes[i] = []byte{1}
		}

		// set limit field
		if len(limits[i]) == 0 {
			tbox.lenlimit[i] = "*"
		} else {
			tbox.lenlimit[i] = limits[i]
		}
	}
}

// get option
func (tbox *OptSel) GetOpt() {
	flag := true
	output := "OPTION  (pos content)  (submit)"
	for flag {
		tbox.printpage(output)
		pos, content := tbox.interpret()
		if pos == -1 {
			switch content {
			case "stop":
				flag = false
			case "submit":
				flag = false
			default:
				output = "OPTION  (pos content)  (submit)"
			}
		} else {
			output = tbox.checkinput(pos, content)
		}
	}
	fmt.Printf("EXIT : %d values\n", len(tbox.Name))
}
