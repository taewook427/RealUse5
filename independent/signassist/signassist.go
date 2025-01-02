// test743 : independent.signtool

package main

import (
	"fmt"
	"os"
	"os/exec"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"golang.org/x/crypto/sha3"
)

// signtool/
//     programs...
//     keystorage0/
//         ca.cer, ca.pvk, ca.txt

// make new sign
func func0() {
	name := kio.Input("Sign Name : ")
	key := kio.Input("Sign Key : ")
	if name == "" {
		name = "NewSign"
	}
	fdir := name
	if fdir[len(fdir)-1] != '/' {
		fdir = fdir + "/"
	}
	os.Mkdir("./"+fdir, os.ModePerm)
	f, _ := kio.Open("./"+fdir+"ca.txt", "w")
	f.Write([]byte(key + "\n"))
	f.Close()
	clipboard.WriteAll(key)

	fmt.Println("키 문자열이 복사되었습니다. 서명 생성 시 키를 입력하세요.\n서명 완료 후 인증서(ca.cer)를 실행하여 로컬 저장소에 설치해야 합니다.")
	fmt.Println("인증서 더블클릭 -> 인증서 설치 -> 현재 사용자 -> 찾아보기 -> 신뢰할 수 있는 루트 인증기관")
	time.Sleep(time.Second * 3)

	cmd := exec.Command("./makecert", "-n", fmt.Sprintf("CN=%s", name), "-r", "-sv", "./"+fdir+"ca.pvk", "./"+fdir+"ca.cer")
	cmd.Run()
	cmd.Wait()
	clipboard.WriteAll("0000")
	fmt.Printf("sign generated at folder %s with key %s\n", fdir, key)
}

// sign exe
func func1(num int) {
	os.Chdir(dirpath[num])
	clipboard.WriteAll(signkey[num])
	cmd := exec.Command("../signtool", "signwizard")
	cmd.Run()
	cmd.Wait()
	clipboard.WriteAll("0000")
	fmt.Printf("sign complete : %s\n", signname[num])
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("critical : %s\n", err)
		}
		kio.Input("press ENTER to exit... ")
	}()
	kobj.Repath()
	flag := false
	nms := []string{"./Cert2Spc.exe", "./CertMgr.exe", "./makecert.exe", "./signtool.exe"}
	hvs := [][]byte{{165, 159, 235, 70, 183, 159, 52, 43, 147, 23, 25, 66, 13, 142, 166, 190, 156, 149, 106, 10, 220, 1, 67, 37, 221, 34, 131, 181, 133, 160, 70, 73, 88, 194, 235, 0, 219, 85, 106, 206, 226, 55, 46, 148, 209, 63, 198, 60, 1, 5, 119, 35, 107, 138, 231, 225, 124, 10, 74, 83, 65, 16, 28, 247},
		{123, 26, 248, 125, 207, 57, 239, 188, 96, 125, 104, 150, 85, 145, 243, 228, 176, 13, 77, 41, 7, 50, 200, 140, 211, 199, 123, 144, 27, 206, 229, 254, 128, 128, 3, 230, 112, 203, 222, 169, 232, 8, 127, 174, 212, 18, 127, 119, 170, 15, 131, 175, 147, 165, 122, 248, 234, 86, 209, 58, 201, 48, 214, 42},
		{78, 219, 190, 52, 140, 5, 211, 69, 6, 79, 250, 254, 88, 17, 255, 84, 152, 142, 86, 198, 64, 237, 132, 79, 160, 41, 236, 87, 232, 67, 221, 232, 38, 22, 38, 151, 105, 174, 183, 39, 112, 174, 235, 71, 253, 25, 125, 125, 130, 18, 158, 16, 149, 47, 106, 61, 11, 191, 219, 245, 61, 203, 22, 251},
		{198, 99, 87, 65, 196, 28, 248, 21, 56, 67, 28, 134, 14, 139, 148, 211, 219, 141, 72, 18, 224, 102, 244, 40, 112, 148, 246, 222, 241, 253, 76, 19, 231, 73, 203, 23, 104, 123, 190, 204, 59, 36, 116, 151, 123, 109, 2, 219, 26, 204, 227, 88, 77, 116, 51, 242, 228, 90, 98, 38, 53, 189, 109, 241}}
	for i, r := range nms {
		if !checkbin(r, hvs[i]) {
			flag = true
		}
	}

	if flag {
		fmt.Println("hash check fail")
	} else {
		readsign()
		fmt.Printf("%d signs available\n[00] %10s [01] %10s\n", len(signname), "Exit", "Make Sign")
		for i, r := range signname {
			fmt.Printf("[%02d] Do Sign : %s\n", i+2, r)
		}
		num, _ := strconv.Atoi(kio.Input(">>> "))
		if num == 0 {
			fmt.Println("program exit")
		} else if num == 1 {
			func0()
		} else if num-2 < len(signname) {
			func1(num - 2)
		} else {
			fmt.Println("unknown mode")
		}
	}
}

var signname []string
var signkey []string
var dirpath []string

func checkbin(path string, hv []byte) bool {
	f, _ := kio.Open(path, "r")
	defer f.Close()
	d, _ := kio.Read(f, -1)
	h := sha3.New512()
	h.Write(d)
	return kio.Bequal(h.Sum(nil), hv)
}

func readsign() {
	fs, _ := os.ReadDir("./")
	for _, r := range fs {
		if r.IsDir() {
			nm := r.Name()
			path := kio.Abs("./" + nm)
			if kio.Size(path+"ca.cer") > 0 && kio.Size(path+"ca.pvk") > 0 && kio.Size(path+"ca.txt") > 0 {
				f, _ := kio.Open(path+"ca.txt", "r")
				d, _ := kio.Read(f, -1)
				f.Close()
				signname = append(signname, nm)
				signkey = append(signkey, strings.ReplaceAll(strings.ReplaceAll(string(d), "\n", ""), "\r", ""))
				dirpath = append(dirpath, path)
			}
		}
	}
}
