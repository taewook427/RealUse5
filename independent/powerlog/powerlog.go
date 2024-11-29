// test712 : independent.powerlog

package main

import (
	"fmt"
	"os"
	"stdlib5/kio"
	"strconv"
	"strings"
	"time"
)

func genlog(path string, word string) {
	t := time.Now().Unix()
	temp := time.Unix(t, 0).Local().Format("2006.01.02;15:04:05")
	temp = fmt.Sprintf("%s/%d#%s\n", temp, t, word)
	f, _ := kio.Open(path, "a")
	kio.Write(f, []byte(temp))
	f.Close()
	fmt.Printf("기록이 저장되었습니다 : %s", temp)
}

func getold(path string) int {
	f, _ := kio.Open(path, "r")
	data, _ := kio.Read(f, 37)
	f.Close()
	raw := string(data)
	raw = raw[strings.Index(raw, "/")+1 : strings.Index(raw, "#")]
	out, err := strconv.Atoi(raw)
	if err == nil {
		fmt.Printf("가장 오래된 데이터 : %s\n", time.Unix(int64(out), 0).Local().Format("2006.01.02"))
	} else {
		fmt.Printf("이전 데이터 읽기 중 오류가 발생했습니다 : %s\n", err)
		out = int(time.Now().Unix())
	}
	return out
}

func main() {
	path, _ := os.UserHomeDir()
	path = kio.Abs(path) + "Desktop/powerlog.txt"
	fmt.Print("RealUse5 powerlog : POWER ON\n\n")
	genlog(path, "POWERED")
	if int(time.Now().Unix())-getold(path) > 94608000 { // 3년보다 오래된 기록
		os.Rename(path, path[:len(path)-12]+"powerlog_old.txt")
		fmt.Println("3년 초과 기록으로 powerlog_old.txt가 생성되었습니다.")
	} else {
		fmt.Println("가장 오래된 기록이 3년 이내입니다.")
	}
	time.Sleep(time.Second * 5)
}
