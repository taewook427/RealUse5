package main

// test596 : khash
// linux에서는 scan이 정상동작하므로 line 282, 368, 417, 428, 498을 지워야함.
import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"stdlib5/cliget"
	"stdlib5/kaes"
	"stdlib5/kdb"
	"stdlib5/ksign"
	"strings"
	"time"
)

// 지연출력
func printlate(msg string, late float64) {
	slpt := time.Duration(late * 200)
	toptr := [5]string{msg, " .", " .", " .", " "}
	for _, r := range toptr {
		fmt.Print(r)
		time.Sleep(slpt * time.Millisecond)
	}
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

// 바이트 출력 (헥스와 base64로)
func bprint(data []byte) (string, string) {
	temp := make([]string, len(data))
	num := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i, r := range data {
		temp[i] = num[r/16] + num[r%16]
	}
	out0 := strings.Join(temp, "")
	out1 := base64.StdEncoding.EncodeToString([]byte(data))
	return out0, string(out1)
}

// 데이터 읽기
func rd(path string) []byte {
	data, _ := ioutil.ReadFile(path)
	return data
}

// 데이터 쓰기
func wr(path string, data []byte) {
	f, _ := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	defer f.Close()
	f.Write(data)
}

// 바이트 값 비교
func bequal(a []byte, b []byte) bool {
	if len(a) == len(b) {
		for i, r := range a {
			if b[i] != r {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

// func 1
func do1() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Critical ERROR : %s\n", err)
		}
	}()

	fmt.Println("khash 해싱할 파일을 선택하세요.")
	time.Sleep(700 * time.Millisecond)
	finder := cliget.Init()
	tps := mktp("windows")
	finder.Tpname = tps[0]
	finder.Tppath = tps[1]
	tgt := finder.Getfile("./")
	htbox := ksign.Khash(tgt)
	oa, ob := bprint(htbox)
	fmt.Printf("Hex : %s\nBase64 : %s\n", oa, ob)
}

// func 2
func do2() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Critical ERROR : %s\n", err)
		}
	}()

	fmt.Println("khash 해싱할 폴더를 선택하세요.")
	time.Sleep(700 * time.Millisecond)
	finder := cliget.Init()
	tps := mktp("windows")
	finder.Tpname = tps[0]
	finder.Tppath = tps[1]
	tgt := finder.Getfolder("./")
	htbox := ksign.Khash(tgt)
	oa, ob := bprint(htbox)
	fmt.Printf("Hex : %s\nBase64 : %s\n", oa, ob)
}

// func 3
func do3() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Critical ERROR : %s\n", err)
		}
	}()

	fmt.Println("검증할 인증서(sign.txt)를 선택하세요.")
	time.Sleep(700 * time.Millisecond)
	finder := cliget.Init()
	tps := mktp("windows")
	finder.Tpname = tps[0]
	finder.Tppath = tps[1]
	tgt := finder.Getfile("./")

	kdbtbox := kdb.Init()
	kdbtbox.Readfile(tgt)
	ts := "rsa.name"
	rsaname := kdbtbox.Getdata(&ts).Dat6
	ts = "rsa.date"
	rsadate := kdbtbox.Getdata(&ts).Dat6
	ts = "rsa.strength"
	strength := kdbtbox.Getdata(&ts).Dat2
	ts = "rsa.public"
	public := kdbtbox.Getdata(&ts).Dat6
	ts = "sign.explain"
	moreinfo := kdbtbox.Getdata(&ts).Dat6
	ts = "sign.date"
	signdate := kdbtbox.Getdata(&ts).Dat6
	ts = "sign.hash"
	hv := kdbtbox.Getdata(&ts).Dat5
	ts = "sign.enc"
	enc := kdbtbox.Getdata(&ts).Dat5

	fmt.Println("검증할 파일을 선택하세요.")
	time.Sleep(700 * time.Millisecond)
	tgt = finder.Getfile("./")
	tgthv := ksign.Khash(tgt)
	fmt.Printf("인증서 : %s, 인증서 생성일 : %s, 인증서 강도(bit) : %d\n서명 정보 : %s, 서명일 : %s\n", rsaname, rsadate, strength, moreinfo, signdate)

	if bequal(hv, tgthv) {
		if ksign.Verify(public, enc, ksign.Fm(rsaname, hv)) {
			fmt.Println("!!! 인증 성공 !!!   파일의 내용이 유효합니다.")
		} else {
			fmt.Println("!!! 인증 실패 !!!   인증서의 RSA 암호를 풀 수 없습니다.")
		}
	} else {
		fmt.Println("!!! 인증 실패 !!!   인증서와 파일의 해시값이 다릅니다.")
	}
}

// func 4
func do4() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Critical ERROR : %s\n", err)
		}
	}()

	fmt.Println("검증할 인증서(sign.txt)를 선택하세요.")
	time.Sleep(700 * time.Millisecond)
	finder := cliget.Init()
	tps := mktp("windows")
	finder.Tpname = tps[0]
	finder.Tppath = tps[1]
	tgt := finder.Getfile("./")

	kdbtbox := kdb.Init()
	kdbtbox.Readfile(tgt)
	ts := "rsa.name"
	rsaname := kdbtbox.Getdata(&ts).Dat6
	ts = "rsa.date"
	rsadate := kdbtbox.Getdata(&ts).Dat6
	ts = "rsa.strength"
	strength := kdbtbox.Getdata(&ts).Dat2
	ts = "rsa.public"
	public := kdbtbox.Getdata(&ts).Dat6
	ts = "sign.explain"
	moreinfo := kdbtbox.Getdata(&ts).Dat6
	ts = "sign.date"
	signdate := kdbtbox.Getdata(&ts).Dat6
	ts = "sign.hash"
	hv := kdbtbox.Getdata(&ts).Dat5
	ts = "sign.enc"
	enc := kdbtbox.Getdata(&ts).Dat5

	fmt.Println("검증할 폴더를 선택하세요.")
	time.Sleep(700 * time.Millisecond)
	tgt = finder.Getfolder("./")
	tgthv := ksign.Khash(tgt)
	fmt.Printf("인증서 : %s, 인증서 생성일 : %s, 인증서 강도(bit) : %d\n서명 정보 : %s, 서명일 : %s\n", rsaname, rsadate, strength, moreinfo, signdate)

	if bequal(hv, tgthv) {
		if ksign.Verify(public, enc, ksign.Fm(rsaname, hv)) {
			fmt.Println("!!! 인증 성공 !!!   폴더의 내용이 유효합니다.")
		} else {
			fmt.Println("!!! 인증 실패 !!!   인증서의 RSA 암호를 풀 수 없습니다.")
		}
	} else {
		fmt.Println("!!! 인증 실패 !!!   인증서와 파일의 해시값이 다릅니다.")
	}
}

// func 5
func do5() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Critical ERROR : %s\n", err)
		}
	}()

	fmt.Println("개인키 파일(private.webp)을 선택하세요.")
	time.Sleep(700 * time.Millisecond)
	finder := cliget.Init()
	tps := mktp("windows")
	finder.Tpname = tps[0]
	finder.Tppath = tps[1]
	tgt := finder.Getfile("./")

	tb := rd(tgt)
	katbox := kaes.Init0()
	hint, msg, stp := katbox.View(tb)
	if msg != "KHASH RSA private file" {
		fmt.Println("KHASH 개인키 파일이 아닙니다.")

	} else {
		fmt.Printf("hint : %s\n개인키 파일 비밀번호 : ", string(hint))
		pw, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		pw = strings.Replace(pw, "\n", "", -1)
		pw = strings.Replace(pw, "\r", "", -1)
		pri := string(katbox.De([]byte(pw), kaes.Genkf("NoSuchPath"), tb, stp))
		kdbtbox := kdb.Init()
		kdbtbox.Readstr(&pri)
		ts := "name"
		name := kdbtbox.Getdata(&ts).Dat6
		ts = "date"
		date := kdbtbox.Getdata(&ts).Dat6
		ts = "strength"
		strength := kdbtbox.Getdata(&ts).Dat2
		ts = "private"
		private := kdbtbox.Getdata(&ts).Dat6
		ts = "public"
		public := kdbtbox.Getdata(&ts).Dat6
		fmt.Printf("< INFO >   이름 : %s, 생성일 : %s, 강도(bit) : %d\n", name, date, strength)

		fmt.Println("서명할 파일을 선택하세요.")
		time.Sleep(700 * time.Millisecond)
		tgt = finder.Getfile("./")
		fmt.Print("이 서명의 추가 정보를 입력하세요. (제작자, 용도 등)\n>>> ")
		var moreinfo string
		fmt.Scan(&moreinfo)
		bufio.NewReader(os.Stdin).ReadString('\n') // Scan 함수로 인해 \n이 버퍼에 남아있으니 비운다.

		curt := time.Now()
		go printlate("sign in progress", 2.5)
		hv := ksign.Khash(tgt)
		plain := ksign.Fm(name, hv)
		enc := ksign.Sign(private, plain)
		ts = "rsa.name = 0\nrsa.date = 0\nrsa.strength = 0\nrsa.public = 0\nsign.explain = 0\nsign.date = 0\nsign.hash = 0\nsign.enc = 0\n"
		kdbtbox = kdb.Init()
		kdbtbox.Readstr(&ts)
		ts = "rsa.name"
		kdbtbox.Fixdata(&ts, name)
		ts = "rsa.date"
		kdbtbox.Fixdata(&ts, date)
		ts = "rsa.strength"
		kdbtbox.Fixdata(&ts, strength)
		ts = "rsa.public"
		kdbtbox.Fixdata(&ts, public)
		ts = "sign.explain"
		kdbtbox.Fixdata(&ts, moreinfo)
		ts = "sign.date"
		kdbtbox.Fixdata(&ts, curt.Format("2006-01-02_15:04:05"))
		ts = "sign.hash"
		kdbtbox.Fixdata(&ts, hv)
		ts = "sign.enc"
		kdbtbox.Fixdata(&ts, enc)
		sv := *kdbtbox.Writestr()
		wr("./sign.txt", []byte(sv))

		passt := int(time.Since(curt).Milliseconds())
		if passt < 2600 {
			time.Sleep(time.Duration(2600-passt) * time.Millisecond)
		}
		fmt.Println("complete!")
		fmt.Println("< WARNING >   전자서명이 완료되었습니다.\nsign.txt를 서명 대상 파일과 같이 배포하세요.")
	}
}

// func 6
func do6() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Critical ERROR : %s\n", err)
		}
	}()

	fmt.Println("개인키 파일(private.webp)을 선택하세요.")
	time.Sleep(700 * time.Millisecond)
	finder := cliget.Init()
	tps := mktp("windows")
	finder.Tpname = tps[0]
	finder.Tppath = tps[1]
	tgt := finder.Getfile("./")

	tb := rd(tgt)
	katbox := kaes.Init0()
	hint, msg, stp := katbox.View(tb)
	if msg != "KHASH RSA private file" {
		fmt.Println("KHASH 개인키 파일이 아닙니다.")

	} else {
		fmt.Printf("hint : %s\n개인키 파일 비밀번호 : ", string(hint))
		pw, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		pw = strings.Replace(pw, "\n", "", -1)
		pw = strings.Replace(pw, "\r", "", -1)
		pri := string(katbox.De([]byte(pw), kaes.Genkf("NoSuchPath"), tb, stp))
		kdbtbox := kdb.Init()
		kdbtbox.Readstr(&pri)
		ts := "name"
		name := kdbtbox.Getdata(&ts).Dat6
		ts = "date"
		date := kdbtbox.Getdata(&ts).Dat6
		ts = "strength"
		strength := kdbtbox.Getdata(&ts).Dat2
		ts = "private"
		private := kdbtbox.Getdata(&ts).Dat6
		ts = "public"
		public := kdbtbox.Getdata(&ts).Dat6
		fmt.Printf("< INFO >   이름 : %s, 생성일 : %s, 강도(bit) : %d\n", name, date, strength)

		fmt.Println("서명할 폴더를 선택하세요.")
		time.Sleep(700 * time.Millisecond)
		tgt = finder.Getfolder("./")
		fmt.Print("이 서명의 추가 정보를 입력하세요. (제작자, 용도 등)\n>>> ")
		var moreinfo string
		fmt.Scan(&moreinfo)
		bufio.NewReader(os.Stdin).ReadString('\n') // Scan 함수로 인해 \n이 버퍼에 남아있으니 비운다.

		curt := time.Now()
		go printlate("sign in progress", 2.5)
		hv := ksign.Khash(tgt)
		plain := ksign.Fm(name, hv)
		enc := ksign.Sign(private, plain)
		ts = "rsa.name = 0\nrsa.date = 0\nrsa.strength = 0\nrsa.public = 0\nsign.explain = 0\nsign.date = 0\nsign.hash = 0\nsign.enc = 0\n"
		kdbtbox = kdb.Init()
		kdbtbox.Readstr(&ts)
		ts = "rsa.name"
		kdbtbox.Fixdata(&ts, name)
		ts = "rsa.date"
		kdbtbox.Fixdata(&ts, date)
		ts = "rsa.strength"
		kdbtbox.Fixdata(&ts, strength)
		ts = "rsa.public"
		kdbtbox.Fixdata(&ts, public)
		ts = "sign.explain"
		kdbtbox.Fixdata(&ts, moreinfo)
		ts = "sign.date"
		kdbtbox.Fixdata(&ts, curt.Format("2006-01-02_15:04:05"))
		ts = "sign.hash"
		kdbtbox.Fixdata(&ts, hv)
		ts = "sign.enc"
		kdbtbox.Fixdata(&ts, enc)
		sv := *kdbtbox.Writestr()
		wr("./sign.txt", []byte(sv))

		passt := int(time.Since(curt).Milliseconds())
		if passt < 2600 {
			time.Sleep(time.Duration(2600-passt) * time.Millisecond)
		}
		fmt.Println("complete!")
		fmt.Println("< WARNING >   전자서명이 완료되었습니다.\nsign.txt를 서명 대상 폴더와 같이 배포하세요.")
	}
}

// func 7
func do7() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Critical ERROR : %s\n", err)
		}
	}()

	fmt.Print("RSA 키 강도를 설정하세요. (bit 단위이며, 보통 2048 사용)\n>>> ")
	var strength int
	fmt.Scan(&strength)
	bufio.NewReader(os.Stdin).ReadString('\n') // Scan 함수로 인해 \n이 버퍼에 남아있으니 비운다.
	if strength == 0 {
		fmt.Println("잘못된 입력입니다.")

	} else if strength != (strength & -strength) {
		fmt.Println("키 강도는 2의 거듭제곱이여야 합니다.")

	} else {
		fmt.Print("이 서명의 이름을 입력하세요. (이름은 16바이트보다 길 수 없음)\n>>> ")
		var name string
		fmt.Scan(&name)
		bufio.NewReader(os.Stdin).ReadString('\n') // Scan 함수로 인해 \n이 버퍼에 남아있으니 비운다.
		nmb := []byte(name)
		if len(nmb) > 16 {
			fmt.Printf("이름이 너무 깁니다 : %d 바이트\n", len(nmb))

		} else {
			curt := time.Now()
			go printlate("generating RSA key", 2.5)
			pub, pri := ksign.Genkey(strength / 8)
			kdbtbox := kdb.Init()
			ts := "name = 0\ndate = 0\nstrength = 0\nprivate = 0\npublic = 0\n"
			kdbtbox.Readstr(&ts)
			ts = "name"
			kdbtbox.Fixdata(&ts, name)
			ts = "date"
			kdbtbox.Fixdata(&ts, curt.Format("2006-01-02_15:04:05"))
			ts = "strength"
			kdbtbox.Fixdata(&ts, strength)
			ts = "private"
			kdbtbox.Fixdata(&ts, pri)
			ts = "public"
			kdbtbox.Fixdata(&ts, pub)
			private := *kdbtbox.Writestr()
			kdbtbox = kdb.Init()
			ts = "name = 0\ndate = 0\nstrength = 0\npublic = 0\n"
			kdbtbox.Readstr(&ts)
			ts = "name"
			kdbtbox.Fixdata(&ts, name)
			ts = "date"
			kdbtbox.Fixdata(&ts, curt.Format("2006-01-02_15:04:05"))
			ts = "strength"
			kdbtbox.Fixdata(&ts, strength)
			ts = "public"
			kdbtbox.Fixdata(&ts, pub)
			public := *kdbtbox.Writestr()
			passt := int(time.Since(curt).Milliseconds())
			if passt < 2600 {
				time.Sleep(time.Duration(2600-passt) * time.Millisecond)
			}
			fmt.Println("complete!")

			fmt.Print("개인키 암호화용 비밀번호 입력 : ")
			pw, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			pw = strings.Replace(pw, "\n", "", -1)
			pw = strings.Replace(pw, "\r", "", -1)
			fmt.Print("비밀번호 힌트 입력 : ")
			hint, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			hint = strings.Replace(hint, "\n", "", -1)
			pw = strings.Replace(pw, "\r", "", -1)
			katbox := kaes.Init0()
			katbox.Msg = "KHASH RSA private file"
			enpriv := katbox.En([]byte(pw), kaes.Genkf("NoSuchPath"), []byte(hint), []byte(private))
			wr("./public.txt", []byte(public))
			wr("./private.webp", enpriv)
			fmt.Println("< WARNING >   키 생성이 모두 완료되었습니다.\npublic.txt는 공개키로, 모두에게 공개해야 합니다.")
			fmt.Println("private.webp는 암호화된 개인키로, 공개돼선 안됩니다.\n개인키 암호화 비밀번호를 잊지 마세요!")
		}
	}
}

// main UI (CLI)
func main() {
	end := true
	for end {
		fmt.Println("\n========== 모드 숫자를 입력하세요 ==========")
		fmt.Println("파일 해싱 (1)   폴더 해싱 (2)   파일 서명 검증 (3)   폴더 서명 검증 (4)")
		fmt.Println("파일 서명 (5)   폴더 서명 (6)   RSA서명키 생성 (7)   종료 (그 외 입력)")
		var inp int
		fmt.Print(">>> ")
		fmt.Scan(&inp)
		bufio.NewReader(os.Stdin).ReadString('\n') // Scan 함수로 인해 \n이 버퍼에 남아있으니 비운다.
		switch inp {
		case 1:
			do1()
		case 2:
			do2()
		case 3:
			do3()
		case 4:
			do4()
		case 5:
			do5()
		case 6:
			do6()
		case 7:
			do7()
		default:
			end = false
		}
	}
	printlate("KHASH 프로그램을 종료합니다", 2.0)
}
