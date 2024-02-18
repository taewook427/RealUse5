// test512 : xe5 데이터 내장기 (win & linux)

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func pymode(num int) {
	txt := "class xe5:\n"
	name := ""
	fmt.Printf("내장할 파일 경로 %d개를 계속 입력하세요\n", num)
	for i := 0; i < num; i++ {
		fmt.Print(">>> ")
		fmt.Scanf("%s", &name)
		fmt.Scanf("%s")

		dat, err := ioutil.ReadFile(name) // 파일 전부 읽어 메모리에 올리기
		if err == nil {
			varname := "data" + strconv.Itoa(i) // dataN
			txt = txt + "    " + varname + " = b\"\"\n"
			for j := 0; j < len(dat)/40; j++ {
				txt = txt + "    " + varname + " = " + varname + " + "
				txt = txt + byte2pystr(dat[40*j:40*j+40]) + "\n"
			}
			if len(dat)%40 != 0 {
				txt = txt + "    " + varname + " = " + varname + " + "
				txt = txt + byte2pystr(dat[40*(len(dat)/40):]) + "\n\n"
			}
		} else {
			fmt.Printf("%d번째에서 예외 발생 : %s\n", i, err)
		}
	}
	file, _ := os.Create("xe5.py")
	defer file.Close()
	file.WriteString(txt)
}

func byte2pystr(data []byte) string {
	temp := make([]string, len(data))
	for i, b := range data {
		temp[i] = fmt.Sprintf("\\x%02x", b)
	}
	return "b\"" + strings.Join(temp, "") + "\""
}

func gomode(num int) {
	txt := "package main\n\n"
	name := ""
	fmt.Printf("내장할 파일 경로 %d개를 계속 입력하세요\n", num)
	for i := 0; i < num; i++ {
		fmt.Print(">>> ")
		fmt.Scanf("%s", &name)
		fmt.Scanf("%s")

		dat, err := ioutil.ReadFile(name) // 파일 전부 읽어 메모리에 올리기
		if err == nil {
			funcname := "xe5data" + strconv.Itoa(i)
			txt = txt + "func " + funcname + "() *[]byte {\n    var temp []byte\n"
			for j := 0; j < len(dat)/40; j++ {
				txt = txt + "    temp = append(temp, " + byte2gostr(dat[40*j:40*j+40]) + ")\n"
			}
			if len(dat)%40 != 0 {
				txt = txt + "    temp = append(temp, " + byte2gostr(dat[40*(len(dat)/40):]) + ")\n"
			}
			txt = txt + "    return &temp\n}\n\n"
		} else {
			fmt.Printf("%d번째에서 예외 발생 : %s\n", i, err)
		}
	}
	txt = txt + "func main() {\n}"
	file, _ := os.Create("xe5.go")
	defer file.Close()
	file.WriteString(txt)
}

func byte2gostr(data []byte) string {
	temp := make([]string, len(data))
	for i, b := range data {
		temp[i] = strconv.Itoa(int(b))
	}
	return strings.Join(temp, ", ")
}

func main() {
	fmt.Println("5세대 데이터 내장기 KOSxe5")
	fmt.Println("모드와 데이터 개수를 선택하세요 (0 : py, 1 : go)")
	fmt.Print("(ex. 0 3) >>> ")
	var mode string
	var num int
	fmt.Scanf("%s %d", &mode, &num)
	fmt.Scanf("%s") // enter 받기용, 윈도우만 필요함

	if mode == "0" {
		pymode(num)
	} else if mode == "1" {
		gomode(num)
	} else {
		fmt.Println("올바른 입력으로 다시 시도하세요. 입력 : ", mode, num)
	}
	fmt.Print("press ENTER to exit... ")
	fmt.Scanf("%s")
}
