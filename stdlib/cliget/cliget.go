// test630 : stdlib5.cliget

package cliget

import (
	"fmt"
	"os"
	"stdlib5/kio"
	"strconv"
	"strings"
)

type toolbox struct {
	tpname []string
	tppath []string
}

// setting tp names, paths
func Initget(names []string, paths []string) toolbox {
	var out toolbox
	var num int
	if len(names) > len(paths) {
		num = len(paths)
	} else {
		num = len(names)
	}
	out.tpname = make([]string, num)
	out.tppath = make([]string, num)
	for i := 0; i < num; i++ {
		out.tpname[i] = names[i]
		out.tppath[i] = kio.Abs(paths[i])
	}
	return out
}

// print one page (output msg + page + tp list + cur dir + sel mode + folder/files)
func printpage(output string, page int, tp string, cur string, mode string, folder []string, file []string) {
	fmt.Print("\033[2J")
	fmt.Println(output)
	fmt.Println(strings.Repeat("=", 32) + fmt.Sprintf(" PAGE %04d START ", page) + strings.Repeat("=", 31))
	fmt.Println(tp)
	fmt.Println(cur)
	fmt.Println(mode)
	count := 0
	for _, r := range folder {
		fmt.Printf("%04d  %s\n", count, r)
		count = count + 1
	}
	for _, r := range file {
		fmt.Printf("%04d  %s\n", count, r)
		count = count + 1
	}
	fmt.Println(strings.Repeat("=", 32) + fmt.Sprintf(" PAGE %04d END ", page) + strings.Repeat("=", 33))
}

// disassemble order, returns only 6 possible value
func interpret() (out []string) {
	defer func() {
		if err := recover(); err != nil {
			out = []string{"nop", ""}
		}
	}()
	raw := kio.Input(">>> ")
	if raw == "!STOP" {
		return []string{"nop", "stop"}
	}
	// tp N, num N, sel N, only S, str S, nop S
	if raw[0] == '"' {
		raw = raw[1:]
		if raw[len(raw)-1] == '"' {
			raw = raw[0 : len(raw)-1]
		}
	}
	if tnum, terr := strconv.Atoi(raw); terr == nil { // pure number
		return []string{"num", fmt.Sprint(tnum)}
	} else if _, terr := os.Stat(raw); terr == nil { // existing file/dir
		return []string{"str", kio.Abs(raw)}
	} else {
		temp := strings.Split(raw, " ")
		switch temp[0] {
		case "tp":
			tnum, terr := strconv.Atoi(temp[1])
			if terr == nil {
				return []string{"tp", fmt.Sprint(tnum)}
			} else {
				return []string{"nop", ""}
			}
		case "sel":
			tnum, terr := strconv.Atoi(temp[1])
			if terr == nil {
				return []string{"sel", fmt.Sprint(tnum)}
			} else {
				return []string{"nop", ""}
			}
		case "only":
			return []string{"only", temp[1]}
		default:
			return []string{"nop", ""}
		}
	}
}

// get one file with start dir init
func (meta *toolbox) Getfile(start string) string {
	curpath := kio.Abs(start) // current view dir
	modeext := "*"            // selection file type
	var folder []string
	var file []string
	target := ""                                            // selected path
	output := "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)" // msg out
	page := 0                                               // passed pages
	tp := "   .TP"                                          // tp msg
	for i, r := range meta.tpname {
		tp = tp + "  " + fmt.Sprint(i) + " : " + r
	}

	for target == "" {
		tdir, _ := os.ReadDir(curpath)
		folder = make([]string, 1)
		file = make([]string, 0)
		folder[0] = "../"
		for _, r := range tdir {
			tnm := r.Name()
			if r.IsDir() {
				if tnm[len(tnm)-1] != '/' {
					tnm = tnm + "/"
				}
				folder = append(folder, tnm)
			} else {
				if modeext == "*" {
					file = append(file, tnm)
				} else if strings.Contains(tnm, ".") {
					if modeext == strings.ToLower(tnm[strings.LastIndex(tnm, ".")+1:]) {
						file = append(file, tnm)
					}
				}
			}
		}

		printpage(output, page, tp, "  .CUR  "+curpath, fmt.Sprintf(" .MODE  File  (with extension %s)", modeext), folder, file)
		page = page + 1
		cmd := interpret()

		switch cmd[0] {
		case "tp":
			tnum, _ := strconv.Atoi(cmd[1])
			if 0 <= tnum && tnum < len(meta.tpname) {
				curpath = meta.tppath[tnum]
				output = "changed current directory"
			} else {
				output = fmt.Sprintf("!wrong_teleport_number : %d", tnum)
			}
		case "num":
			tnum, _ := strconv.Atoi(cmd[1])
			if tnum == 0 && strings.Count(curpath, "/") > 1 {
				curpath = curpath[0 : strings.LastIndex(curpath[0:len(curpath)-1], "/")+1]
				output = "moved to parent directory"
			} else if 0 < tnum && tnum < len(folder) {
				curpath = curpath + folder[tnum]
				output = "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)"
			} else if 0 < tnum && tnum < len(folder)+len(file) {
				target = curpath + file[tnum-len(folder)]
				output = fmt.Sprintf("!file_selected : %s", target)
			} else {
				output = fmt.Sprintf("!wrong_selection_number : %d", tnum)
			}
		case "sel":
			tnum, _ := strconv.Atoi(cmd[1])
			if 0 <= tnum && tnum < len(folder) {
				output = fmt.Sprintf("!cannot_select_folder : %d", tnum)
			} else if 0 < tnum && tnum < len(folder)+len(file) {
				target = curpath + file[tnum-len(folder)]
				output = fmt.Sprintf("!file_selected : %s", target)
			} else {
				output = fmt.Sprintf("!wrong_selection_number : %d", tnum)
			}
		case "only":
			modeext = strings.ToLower(cmd[1])
			output = "to return to full selection, enter (only *)"
		case "str":
			if cmd[1][len(cmd[1])-1] == '/' {
				output = fmt.Sprintf("!cannot_select_folder : %s", cmd[1])
			} else {
				target = cmd[1]
				output = fmt.Sprintf("!file_selected : %s", target)
			}
		case "nop":
			if cmd[1] == "stop" {
				return ""
			} else {
				output = "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)"
			}
		default:
			output = "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)"
		}
	}
	fmt.Println(output)
	return target
}

// get one foldere with start dir init
func (meta *toolbox) Getfolder(start string) string {
	curpath := kio.Abs(start) // current view dir
	modeext := "*"            // selection filter
	var folder []string
	target := ""                                            // selected path
	output := "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)" // msg out
	page := 0                                               // passed pages
	tp := "   .TP"                                          // tp msg
	for i, r := range meta.tpname {
		tp = tp + "  " + fmt.Sprint(i) + " : " + r
	}

	for target == "" {
		tdir, _ := os.ReadDir(curpath)
		folder = make([]string, 1)
		folder[0] = "../"
		for _, r := range tdir {
			tnm := r.Name()
			if r.IsDir() {
				if tnm[len(tnm)-1] != '/' {
					tnm = tnm + "/"
				}
				if modeext == "*" {
					folder = append(folder, tnm)
				} else if strings.Contains(strings.ToLower(tnm), modeext) {
					folder = append(folder, tnm)
				}
			}
		}

		printpage(output, page, tp, "  .CUR  "+curpath, fmt.Sprintf(" .MODE  Folder  (with filter %s)", modeext), folder, nil)
		page = page + 1
		cmd := interpret()

		switch cmd[0] {
		case "tp":
			tnum, _ := strconv.Atoi(cmd[1])
			if 0 <= tnum && tnum < len(meta.tpname) {
				curpath = meta.tppath[tnum]
				output = "changed current directory"
			} else {
				output = fmt.Sprintf("!wrong_teleport_number : %d", tnum)
			}
		case "num":
			tnum, _ := strconv.Atoi(cmd[1])
			if tnum == 0 && strings.Count(curpath, "/") > 1 {
				curpath = curpath[0 : strings.LastIndex(curpath[0:len(curpath)-1], "/")+1]
				output = "moved to parent directory"
			} else if 0 < tnum && tnum < len(folder) {
				curpath = curpath + folder[tnum]
				output = "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)"
			} else {
				output = fmt.Sprintf("!wrong_selection_number : %d", tnum)
			}
		case "sel":
			tnum, _ := strconv.Atoi(cmd[1])
			if tnum == 0 {
				target = curpath
				output = fmt.Sprintf("!folder_selected : %s", target)
			} else if 0 < tnum && tnum < len(folder) {
				target = curpath + folder[tnum]
				output = fmt.Sprintf("!folder_selected : %s", target)
			} else {
				output = fmt.Sprintf("!wrong_selection_number : %d", tnum)
			}
		case "only":
			modeext = strings.ToLower(cmd[1])
			output = "to return to full selection, enter (only *)"
		case "str":
			if cmd[1][len(cmd[1])-1] == '/' {
				target = cmd[1]
				output = fmt.Sprintf("!folder_selected : %s", target)
			} else {
				output = fmt.Sprintf("!cannot_select_file : %s", cmd[1])
			}
		case "nop":
			if cmd[1] == "stop" {
				return ""
			} else {
				output = "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)"
			}
		default:
			output = "OPTION  (tp N)  (N)  (sel N)  (only S)  (S)"
		}
	}
	fmt.Println(output)
	return target
}
