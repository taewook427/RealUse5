package cliget

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// test575 : cliget (go)

type toolbox struct {
	Tpname []string
	Tppath []string
}

// toolbox 반환, Tpname/Tppath는 따로 설정
func Init() toolbox {
	var temp toolbox
	return temp
}

// 시작 폴더를 입력해 파일 선택
func (self *toolbox) Getfile(start string) string {
	start, _ = filepath.Abs(start)
	start = strings.Replace(start, "\\", "/", -1)
	if start[len(start)-1] != '/' {
		start = start + "/"
	}
	for i, r := range self.Tppath {
		tpt, _ := filepath.Abs(r)
		tpt = strings.Replace(tpt, "\\", "/", -1)
		if tpt[len(tpt)-1] != '/' {
			tpt = tpt + "/"
		}
		self.Tppath[i] = tpt
	}
	curpath := start
	modeext := "*"
	var folder []string
	var file []string
	target := ""

	output := "바로가기(tp N)  이동(N)  선택(sel N)  확장자(only S)  직접입력(S)" // 메세지 출력부
	page := 0                                                      // 지나간 페이지 수
	tp := "  TP"                                                   // 바로가기 표시부
	for i, r := range self.Tpname {
		tp = tp + "  " + fmt.Sprint(i) + " : " + r
	}
	var cur string  // 현재 폴더 경로
	var mode string // 선택 모드

	for target == "" {
		cur = " CUR  " + curpath
		mode = "MODE : FILE  ext : " + modeext

		tdir, _ := ioutil.ReadDir(curpath)
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

		input := printer(&output, page, &tp, &cur, &mode, &folder, &file)
		page = page + 1

		if strings.Contains(input, " ") {
			cmd := strings.Split(input, " ")
			num, err := strconv.Atoi(cmd[1])
			switch cmd[0] {
			case "tp": // tp N
				if err == nil && num < len(self.Tppath) {
					output = "바로가기로 이동하였습니다."
					curpath = self.Tppath[num]
				} else {
					output = fmt.Sprintf("잘못된 tp 인자 : %s", cmd[1])
				}
			case "sel": // sel N
				if err == nil {
					if num < len(folder) {
						output = "폴더를 선택할 수 없습니다."
					} else if num < len(folder)+len(file) {
						target = curpath + file[num-len(folder)]
						output = fmt.Sprintf("선택됨  %s", target)
					}
				} else {
					output = fmt.Sprintf("잘못된 sel 인자 : %s", cmd[1])
				}
			case "only": // only S
				modeext = strings.ToLower(cmd[1])
				output = "전체선택으로 되돌릴려면 only * 입력."
			default: // cmd cmd
				output = "잘못된 명령."
			}
		} else if len(input) > 0 {
			num, err := strconv.Atoi(input)
			if err == nil { // N
				if num == 0 && strings.Count(curpath, "/") > 1 {
					curpath = curpath[0 : strings.LastIndex(curpath[0:len(curpath)-1], "/")+1]
					output = "상위 폴더로 이동합니다."
				} else if num < len(folder) {
					curpath = curpath + folder[num]
					output = "바로가기(tp N)  이동(N)  선택(sel N)  확장자(only S)  직접입력(S)"
				} else if num < len(folder)+len(file) {
					target = curpath + file[num-len(folder)]
					output = fmt.Sprintf("선택됨  %s", target)
				}
			} else { // S
				if input[0] == '"' && input[len(input)-1] == '"' {
					input = input[1 : len(input)-1]
				}
				_, err := os.Open(input)
				if err == nil {
					target = input
					output = fmt.Sprintf("선택됨  %s", target)
				}
			}
		}
	}
	fmt.Println(output)
	return target
}

// 시작 폴더를 입력해 폴더 선택
func (self *toolbox) Getfolder(start string) string {
	start, _ = filepath.Abs(start)
	start = strings.Replace(start, "\\", "/", -1)
	if start[len(start)-1] != '/' {
		start = start + "/"
	}
	for i, r := range self.Tppath {
		tpt, _ := filepath.Abs(r)
		tpt = strings.Replace(tpt, "\\", "/", -1)
		if tpt[len(tpt)-1] != '/' {
			tpt = tpt + "/"
		}
		self.Tppath[i] = tpt
	}
	curpath := start
	var folder []string
	var file []string
	target := ""

	output := "바로가기(tp N)  이동(N)  선택(sel N)  직접입력(S)" // 메세지 출력부
	page := 0                                         // 지나간 페이지 수
	tp := "  TP"                                      // 바로가기 표시부
	for i, r := range self.Tpname {
		tp = tp + "  " + fmt.Sprint(i) + " : " + r
	}
	var cur string  // 현재 폴더 경로
	var mode string // 선택 모드

	for target == "" {
		cur = " CUR  " + curpath
		mode = "MODE : FOLDER"

		tdir, _ := ioutil.ReadDir(curpath)
		folder = make([]string, 1)
		folder[0] = "../"
		for _, r := range tdir {
			tnm := r.Name()
			if r.IsDir() {
				if tnm[len(tnm)-1] != '/' {
					tnm = tnm + "/"
				}
				folder = append(folder, tnm)
			}
		}

		input := printer(&output, page, &tp, &cur, &mode, &folder, &file)
		page = page + 1

		if strings.Contains(input, " ") {
			cmd := strings.Split(input, " ")
			num, err := strconv.Atoi(cmd[1])
			switch cmd[0] {
			case "tp": // tp N
				if err == nil && num < len(self.Tppath) {
					output = "바로가기로 이동하였습니다."
					curpath = self.Tppath[num]
				} else {
					output = fmt.Sprintf("잘못된 tp 인자 : %s", cmd[1])
				}
			case "sel": // sel N
				if err == nil {
					if num == 0 {
						output = "폴더를 선택할 수 없습니다."
					} else if num < len(folder) {
						target = curpath + folder[num]
						output = fmt.Sprintf("선택됨  %s", target)
					}
				} else {
					output = fmt.Sprintf("잘못된 sel 인자 : %s", cmd[1])
				}
			default: // cmd cmd
				output = "잘못된 명령."
			}
		} else if len(input) > 0 {
			num, err := strconv.Atoi(input)
			if err == nil { // N
				if num == 0 && strings.Count(curpath, "/") > 1 {
					curpath = curpath[0 : strings.LastIndex(curpath[0:len(curpath)-1], "/")+1]
					output = "상위 폴더로 이동합니다."
				} else if num < len(folder) {
					curpath = curpath + folder[num]
					output = "바로가기(tp N)  이동(N)  선택(sel N)  직접입력(S)"
				}
			} else { // S
				if input[0] == '"' && input[len(input)-1] == '"' {
					input = input[1 : len(input)-1]
				}
				_, err := os.ReadDir(input)
				if err == nil {
					target = input
					output = fmt.Sprintf("선택됨  %s", target)
				}
			}
		}
	}
	fmt.Println(output)
	return target
}

// printer
func printer(output *string, page int, tp *string, cur *string, mode *string, folder *[]string, file *[]string) string {
	var input string
	fmt.Print("\033[2J")
	fmt.Println(*output)
	fmt.Println(strings.Repeat("=", 32) + fmt.Sprintf(" PAGE %04d START ", page) + strings.Repeat("=", 31))
	fmt.Println(*tp)
	fmt.Println(*cur)
	fmt.Println(*mode)
	count := 0
	for _, r := range *folder {
		fmt.Printf("%04d  %s\n", count, r)
		count = count + 1
	}
	for _, r := range *file {
		fmt.Printf("%04d  %s\n", count, r)
		count = count + 1
	}
	fmt.Println(strings.Repeat("=", 32) + fmt.Sprintf(" PAGE %04d END ", page) + strings.Repeat("=", 33))
	input = reader(">>> ")
	return input
}

// reader
func reader(msg string) string {
	fmt.Print(msg)
	r := bufio.NewReader(os.Stdin)
	s, _ := r.ReadString('\n')
	if s[len(s)-1] == '\n' {
		s = s[0 : len(s)-1]
	}
	if s[len(s)-1] == '\r' {
		s = s[0 : len(s)-1]
	}
	return s
}
