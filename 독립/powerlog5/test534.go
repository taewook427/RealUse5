package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

func unixt() int64 {
	t := time.Now()
	// 원하는 날짜 및 시간 형식을 정의합니다.
	datf := "2006-01-02_15:04:05" // Go의 특별한 날짜 형식

	unixt := t.Format(datf)
	// 유닉스 시간을 계산합니다.
	tint, _ := time.Parse(datf, unixt)
	tstamp := tint.Unix()
	return tstamp
}

func mktime() string {
	t := time.Now()

	// 원하는 날짜 및 시간 형식을 정의합니다.
	datf := "2006-01-02_15:04:05" // Go의 특별한 날짜 형식

	// 로컬 날짜 및 유닉스 시간을 포맷팅합니다.
	local := t.Format(datf)

	// 결과를 출력합니다.
	result := local + "/" + fmt.Sprintf("%d", unixt()) + "#POWERED\n"
	return result
}

func mknew(path string) {
	f, _ := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	defer f.Close()
	f.Write([]byte("2006-01-02_15:04:05/1136181845#POWERED"))
	fmt.Println("새 기록 텍스트 파일이 생성되었습니다")
}

func wrnow(path string) {
	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	temp := mktime()
	f.Write([]byte(temp))
	fmt.Println("현재 로그가 저장되었습니다\n" + temp)
}

func main() {
	fmt.Println("KOS2023 - gen5 powerlog : POWER ON")

	usr, _ := user.Current()
	path := usr.HomeDir + "/Desktop/powerlog.txt"
	path = strings.ReplaceAll(path, "\\", "/")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		mknew(path)
	}

	dat, _ := ioutil.ReadFile(path)
	logs := strings.Split(fmt.Sprintf("%s", dat), "\n")
	if logs[len(logs)-1] == "" {
		logs = logs[0 : len(logs)-1]
	}
	dates := make([]int, len(logs))
	for i, r := range logs {
		dates[i], _ = strconv.Atoi(r[strings.Index(r, "/")+1 : strings.Index(r, "#")])
	}

	max := 157680000 // 5년
	current := int(unixt())
	if current-(dates[0]) < max {
		fmt.Println("가장 오래된 기록이 5년 이내입니다.")
	} else {
		num := 0
		f, _ := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
		for i, r := range dates {
			if current-r > max {
				num = num + 1
			} else {
				f.Write([]byte(logs[i]))
			}
		}
		f.Close()
		fmt.Printf("삭제된 기록 : %d 개\n보존된 기록 : %d 개\n", num, len(dates)-num)
	}

	wrnow(path)
	time.Sleep(5 * time.Second)
}
