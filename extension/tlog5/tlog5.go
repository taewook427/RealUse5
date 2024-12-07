// test729 : extension.tlog5

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"stdlib5/cliget"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strconv"
	"strings"
	"time"
)

// init data folder, returns DB path
func initdb() string {
	if _, err := os.Stat("./_ST5_DATA/"); err != nil {
		os.Mkdir("./_ST5_DATA/", os.ModePerm)
	}
	fs, _ := os.ReadDir("./_ST5_DATA/")
	path := ""
	for _, r := range fs {
		temp := r.Name()
		if len(temp) > 4 && temp[len(temp)-4:] == ".txt" {
			path = "./_ST5_DATA/" + temp
		}
	}
	if path == "" {
		path = "./_ST5_DATA/log729.txt"
		f, _ := kio.Open(path, "w")
		kio.Write(f, []byte(writetxt(0, "log generated")+"\n"))
		f.Close()
	}
	return path
}

// get db
func getdb(path string) []string {
	f, _ := kio.Open(path, "r")
	d, _ := kio.Read(f, -1)
	f.Close()
	d = d[:len(d)-1]
	return strings.Split(string(d), "\n")
}

// read DB string lognum
func readtxt(line string) int {
	if i, err := strconv.Atoi(line[strings.Index(line, "/")+1 : strings.Index(line, "#")]); err == nil {
		return i
	} else {
		return -1
	}
}

// make DB string
func writetxt(num int, info string) string {
	info = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(info, "\n", " "), "/", " "), "#", " ")
	return fmt.Sprintf("%s/%06d#%s", time.Now().Local().Format("2006.01.02;15:04:05"), num, info)
}

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

// logger struct
type logger struct {
	path    string   // logfile path
	data    []string // entire log data
	lastnum int      // last log num
	flag    bool     // working flag
}

// automatic log
func (tbox *logger) func1() {
	info := kio.Input(fmt.Sprintf("%d : ", tbox.lastnum+1))
	info = writetxt(tbox.lastnum+1, info)
	tbox.data = append(tbox.data, info)
	tbox.lastnum = tbox.lastnum + 1
	f, _ := kio.Open(tbox.path, "a")
	kio.Write(f, []byte(info+"\n"))
	f.Close()
	fmt.Println(info)
}

// manual log
func (tbox *logger) func2() {
	temp := kio.Input("Num Info : ")
	var num int
	var info string
	if pos := strings.Index(temp, " "); pos < 0 {
		num = tbox.lastnum + 1
		info = temp
	} else if i, r := strconv.Atoi(temp[0:pos]); r == nil {
		num = i
		info = temp[pos+1:]
	} else {
		num = tbox.lastnum + 1
		info = temp
	}

	info = writetxt(num, info)
	tbox.data = append(tbox.data, info)
	tbox.lastnum = num
	f, _ := kio.Open(tbox.path, "a")
	kio.Write(f, []byte(info+"\n"))
	f.Close()
	fmt.Println(info)
}

// log merge
func (tbox *logger) func3() {
	var sel cliget.OptSel
	sel.Init([]string{"< Select Log >", "main branch", "sub branch"}, []string{"bool", "file", "file"}, []string{"T", "txt", "txt"}, *getpath(), nil)
	sel.GetOpt()
	p_main := sel.StrRes[1]
	p_sub := sel.StrRes[2]
	if p_main == "" || p_sub == "" {
		fmt.Println("[error] invalid log file selection")
		return
	}
	main_pos := 0
	main_db := getdb(p_main)
	main_num := make([]int, len(main_db))
	main_info := make([]string, len(main_db))
	for i, r := range main_db {
		main_num[i] = readtxt(r)
		main_info[i] = r[strings.Index(r, "#")+1:]
	}
	sub_pos := 0
	sub_db := getdb(p_sub)
	sub_num := make([]int, len(sub_db))
	sub_info := make([]string, len(sub_db))
	for i, r := range sub_db {
		sub_num[i] = readtxt(r)
		sub_info[i] = r[strings.Index(r, "#")+1:]
	}

	tbox.data = make([]string, 0) // merge
	for main_pos < len(main_db) && sub_pos < len(sub_db) {
		if main_num[main_pos] < sub_num[sub_pos] { // main < sub
			tbox.data = append(tbox.data, main_db[main_pos]+"\n")
			main_pos = main_pos + 1
		} else if main_num[main_pos] > sub_num[sub_pos] { // main > sub
			tbox.data = append(tbox.data, sub_db[sub_pos]+"\n")
			sub_pos = sub_pos + 1
		} else { // main == sub
			tbox.data = append(tbox.data, main_db[main_pos]+"\n")
			if main_info[main_pos] != sub_info[sub_pos] {
				tbox.data = append(tbox.data, sub_db[sub_pos]+"\n")
			}
			main_pos = main_pos + 1
			sub_pos = sub_pos + 1
		}
	}
	for main_pos < len(main_db) {
		tbox.data = append(tbox.data, main_db[main_pos]+"\n")
		main_pos = main_pos + 1
	}
	for sub_pos < len(sub_db) {
		tbox.data = append(tbox.data, sub_db[sub_pos]+"\n")
		sub_pos = sub_pos + 1
	}

	f, _ := kio.Open(p_main, "w")
	kio.Write(f, []byte(strings.Join(tbox.data, "")))
	f.Close()
	fmt.Printf("merged branch at %s\n", p_main)
}

// log fork
func (tbox *logger) func4() {
	fmt.Println("fork from current branch")
	num, _ := strconv.Atoi(kio.Input("filter number : "))
	temp := make([]string, 0)
	for _, r := range tbox.data {
		if readtxt(r) > num {
			temp = append(temp, r+"\n")
		}
	}
	f, _ := kio.Open("./_ST5_DATA/new_branch.txt", "w")
	kio.Write(f, []byte(strings.Join(temp, "")))
	f.Close()
	fmt.Println("new branch at ./_ST5_DATA/new_branch.txt")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("critical : %s\n", err)
		}
		kio.Input("press ENTER to exit... ")
	}()
	kobj.Repath()
	var worker logger
	worker.path = initdb()
	worker.data = getdb(worker.path)
	for i, r := range worker.data {
		fmt.Printf("[%06d]   %s\n", i, r)
	}
	worker.lastnum = readtxt(worker.data[len(worker.data)-1])
	if worker.lastnum < 0 {
		fmt.Println("[error] invalod log file")
		return
	}

	worker.flag = true
	for worker.flag {
		fmt.Printf("\n%17s   [0] %13s   [1] %13s\n[2] %13s   [3] %13s   [4] %13s\n", "< Select Mode >", "Exit Program", "Automatic Log", "Manual Log", "Merge Log", "Fork Log")
		fmt.Printf("%s [last : %d] ", worker.path, worker.lastnum)
		switch kio.Input(">>> ") {
		case "0":
			worker.flag = false
		case "1":
			worker.func1()
		case "2":
			worker.func2()
		case "3":
			worker.func3()
			worker.flag = false
		case "4":
			worker.func4()
			worker.flag = false
		}
	}
}
