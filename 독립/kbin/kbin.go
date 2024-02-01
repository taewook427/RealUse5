package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"stdlib5/cliget"
	"strconv"
	"strings"
	"time"
)

// test597 : kbin (go)

func input(msg string) string {
	fmt.Print(msg)
	temp, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if len(temp) == 0 {
		return temp
	} else if len(temp) == 1 {
		if temp == "\n" || temp == "\r" {
			return ""
		} else {
			return temp
		}
	} else {
		if temp[len(temp)-1] == '\n' {
			temp = temp[0 : len(temp)-1]
		}
		if temp[len(temp)-1] == '\r' {
			temp = temp[0 : len(temp)-1]
		}
		return temp
	}
}

// 실제 오프셋, 데이터 512B, 추가메세지 S
func bprint(offset int, data []byte, msg string) {
	buffer := make([]string, 34)
	sreg := make([]string, 36)
	nums := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
	conv := make([]string, 256)
	for i := 0; i < 256; i++ {
		if i < 32 {
			conv[i] = "."
		} else if i < 127 {
			conv[i] = string(byte(i))
		} else {
			ts := string(byte(i))
			if len(ts) == 1 {
				conv[i] = ts
			} else {
				conv[i] = "."
			}
		}
	}
	sreg[0] = fmt.Sprintf("%010d", offset/1024)
	if offset%1024 == 0 {
		sreg[1] = "A"
	} else {
		sreg[1] = "B"
	}
	sreg[2] = "   "
	for i, r := range nums {
		sreg[i+3] = fmt.Sprintf("%02d ", i)
		sreg[i+20] = r
	}
	sreg[19] = "  "
	buffer[0] = strings.Join(sreg, "")

	for i := 0; i < 32; i++ {
		sreg = make([]string, 35)
		sreg[0] = fmt.Sprintf("       %03d ", 16*i)
		sreg[1] = "   "
		for j := 0; j < 16; j++ {
			if 16*i+j < len(data) {
				r := data[16*i+j]
				sreg[j+2] = nums[r/16] + nums[r%16] + " "
				sreg[j+19] = conv[r]
			} else {
				sreg[j+2] = "   "
				sreg[j+19] = " "
			}
		}
		sreg[18] = "  "
		buffer[i+1] = strings.Join(sreg, "")
	}

	buffer[33] = msg
	fmt.Println("")
	fmt.Println(strings.Join(buffer, "\n"))
}

// 바로가기 경로 생성
func mktp(ostype string) [][]string {
	out := [][]string{{"", ""}, {"", ""}}
	if ostype == "windows" {
		out[0][0] = "사용자 폴더"
		out[0][1] = "바탕화면"
		dp, _ := os.UserHomeDir()
		ddp := filepath.Join(dp, "Desktop")
		out[1][0] = dp
		out[1][1] = ddp
	} else {
		out[0][0] = "사용자 폴더"
		out[0][1] = "바탕화면"
		dp, _ := os.UserHomeDir()
		ddp := filepath.Join(dp, "Desktop")
		out[1][0] = dp
		out[1][1] = ddp
	}
	return out
}

// 읽기
func read(path string, offset int, size int) []byte {
	f, err := os.Open(path)
	defer f.Close()
	if err == nil {
		f.Seek(int64(offset), 0)
		temp := make([]byte, size)
		f.Read(temp)
		return temp
	} else {
		return make([]byte, size)
	}
}

// 명령 인터프리터
func interpret() (ret []string) {
	defer func() {
		if err := recover(); err != nil {
			ret = nil
		}
	}()
	temp := strings.Split(input(">>> "), " ")
	if len(temp) == 0 {
		return nil
	} else {
		switch temp[0] {
		case "x":
			return []string{"x"}
		case "p":
			return []string{"p"}
		case "n":
			return []string{"n"}
		case "open":
			return []string{"open"}
		case "tp":
			ti, err := strconv.Atoi(temp[1])
			if err == nil {
				if ti < 0 {
					return []string{"tp", "0"}
				} else {
					return []string{"tp", temp[1]}
				}
			} else {
				return nil
			}
		case "str":
			rega := strings.Split(temp[1], ".")
			regb := strings.Split(temp[2], ".")
			ia, err := strconv.Atoi(rega[0])
			if err != nil {
				return nil
			}
			ib, err := strconv.Atoi(rega[1])
			if err != nil {
				return nil
			}
			ic, err := strconv.Atoi(regb[0])
			if err != nil {
				return nil
			}
			id, err := strconv.Atoi(regb[1])
			if err != nil {
				return nil
			}
			if ib < 0 {
				ib = 0
			}
			if id < 0 {
				id = 0
			}
			if ib > 511 {
				ib = 511
			}
			if id > 511 {
				id = 511
			}
			return []string{"str", fmt.Sprint(ia), fmt.Sprint(ib), fmt.Sprint(ic), fmt.Sprint(id)}
		default:
			return nil
		}
	}
}

func main() {
	var args []string
	for _, r := range os.Args {
		args = append(args, r)
	}

	var path string
	if len(args) < 2 {
		fmt.Println("읽을 파일을 선택하세요.")
		time.Sleep(700 * time.Millisecond)
		finder := cliget.Init()
		tps := mktp("windows")
		finder.Tpname = tps[0]
		finder.Tppath = tps[1]
		path = finder.Getfile("./")
	} else {
		path = args[1]
	}

	f, _ := os.Stat(path)
	size := int(f.Size())
	var buffer []byte
	offset := 0
	flag := true
	for flag {

		readable := size - offset
		if readable > 512 {
			buffer = read(path, offset, 512)
		} else if readable > 0 {
			buffer = read(path, offset, readable)
		} else {
			buffer = make([]byte, 0)
		}
		bprint(offset, buffer, "종료(x)   이전(p)   다음(n)   열기(open)   이동(tp N)   문자열화(str N.N N.N)")

		order := interpret()
		if order != nil {
			switch order[0] {
			case "x":
				flag = false
			case "p":
				if offset != 0 {
					offset = offset - 512
				}
			case "n":
				if readable > 512 {
					offset = offset + 512
				}
			case "open":
				flag = false
				ct := exec.Command("./hxdraw.exe", path) // win linux 명령어 차이
				ct.Output()
			case "tp":
				ti, _ := strconv.Atoi(order[1])
				if ti*1024 > size {
					ti = size / 1024
				}
				offset = ti * 1024
			case "str":
				ia, _ := strconv.Atoi(order[1])
				ib, _ := strconv.Atoi(order[2])
				ic, _ := strconv.Atoi(order[3])
				id, _ := strconv.Atoi(order[4])
				start := offset + 512*ia + ib
				end := offset + 512*ic + id
				if 0 <= start && end <= size && start < end && end-start < 16384 {
					fmt.Println("========== UTF-8 str start ==========")
					fmt.Println(string(read(path, start, end-start)))
					fmt.Println("========== UTF-8 str end ==========")
				}
			default:
				flag = true
			}
		}
	}
}
