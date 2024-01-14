package kweb

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// test562 : kweb (go)

// go mod init example.com
// go get github.com/PuerkitoBio/goquery

// http~.html에서 domain 별 txt 반환
func Gettxt(url string, domain string) string {
	// HTTP GET 요청
	resp, err := http.Get(url)
	if err != nil {
		panic("get fail")
	}
	defer resp.Body.Close()

	// HTTP 응답이 200 OK인지 확인
	if resp.StatusCode != http.StatusOK {
		panic("200 fail")
	}

	// HTML 파싱
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic("parse fail")
	}

	// 클래스가 domain인 요소 찾기
	text := doc.Find(fmt.Sprintf("p#%s", domain)).Text()
	if text == "" {
		panic("text find fail")
	}

	return text
}

// http~/ + name + .num 를 path로 바이너리 생성
func Download(url string, name string, num int, path string) {
	// URL이 끝에 슬래시가 없으면 추가
	if url[len(url)-1] != '/' {
		url = url + "/"
	}

	// 파일 생성
	file, err := os.Create(path)
	if err != nil {
		panic("file creation fail")
	}
	defer file.Close()

	// 파일에 데이터 쓰기
	for i := 0; i < num; i++ {
		resp, err := http.Get(fmt.Sprintf("%s%s.%d", url, name, i))
		if err != nil {
			panic("download fail")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			panic("connection fail")
		}

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			panic("writing fail")
		}
	}
}

// edgedriver 홈페이지에서 정보 받아오기, 홈페이지 + 행동정보 + zip경로 -> 버전 정수
func Driver(url string, datas []string, path string) int {
	// HTTP GET 요청
	resp, err := http.Get(url)
	if err != nil {
		panic("get fail")
	}
	defer resp.Body.Close()

	// HTTP 응답이 200 OK인지 확인
	if resp.StatusCode != http.StatusOK {
		panic("200 fail")
	}

	// HTML 파싱
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic("parse fail")
	}

	// stable box 찾기
	var box *goquery.Selection
	doc.Find(datas[0] + "." + datas[1]).Each(func(i int, s *goquery.Selection) {
		text := strings.ToLower(strings.ReplaceAll(s.Text(), " ", ""))
		if strings.Contains(text, "stable") {
			box = s
		}
	})
	if box == nil {
		panic("No Stable WD")
	}

	// x64 button 찾기
	var but *goquery.Selection
	box.Find(datas[2] + "." + datas[3]).Each(func(i int, s *goquery.Selection) {
		text := strings.ToLower(strings.ReplaceAll(s.Text(), " ", ""))
		if strings.Contains(text, "x64") {
			but = s
		}
	})
	if but == nil {
		panic("No x64 WD")
	}

	// wd.zip link, ver int
	link, _ := but.Attr("href")
	ver := 0
	text := box.Find(datas[4] + "." + datas[5]).Text()
	temp := ""
	nums := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	for _, r := range text {
		for _, l := range nums {
			if r == l {
				temp = temp + string(r)
			}
		}
		if r == '.' {
			ver, _ = strconv.Atoi(temp)
			break
		}
	}

	// wd.zip download
	file, err := os.Create(path)
	if err != nil {
		panic("file creation fail")
	}
	defer file.Close()
	resp, err = http.Get(link)
	if err != nil {
		panic("download fail")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("connection fail")
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic("writing fail")
	}

	return ver
}
