package kzip

// test571 : kzip (go)

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"example.com/picdt"
)

// ========== ========== ksc5 start ========== ==========

// little endian encoding
func encode(num int, length int) *[]byte {
	temp := make([]byte, length)
	for i := 0; i < length; i++ {
		temp[i] = byte(num % 256)
		num = num / 256
	}
	return &temp
}

// little endian decoding
func decode(data *[]byte) int {
	temp := 0
	for i, r := range *data {
		if r != 0 {
			exp := 1
			for j := 0; j < i; j++ {
				exp = exp * 256
			}
			temp = temp + int(r)*exp
		}
	}
	return temp
}

// crc32 check
func crc32check(data *[]byte) *[]byte {
	temp := int(crc32.ChecksumIEEE(*data))
	return encode(temp, 4)
}

// find 1024nB + KSC5 pos
func findpos(data *[]byte) int {
	temp := 0
	cmp := []byte("KSC5")
	for len(*data) >= temp+4 {
		d := (*data)[temp : temp+4]
		if bequal(&cmp, &d) {
			return temp
		} else {
			temp = temp + 1024
		}
	}
	return -1
}

// compare two []byte
func bequal(a *[]byte, b *[]byte) bool {
	if len(*a) != len(*b) {
		return false
	}
	for i, r := range *a {
		if r != (*b)[i] {
			return false
		}
	}
	return true
}

// ========== ========== ksc5 end ========== ==========

// ========== ========== kzip5 start ========== ==========

// file 존재시 삭제
func initfile(path string) {
	os.Remove(path)
}

// folder 존재시 삭제
func initfolder(path string) {
	os.RemoveAll(path)
}

// root : 표준절대경로/, path : 표준절대경로/에 폴더,파일 표준절대경로 추가
func getlist(root *string, path *string, lstf *schunk, lstd *schunk) {
	ft, _ := ioutil.ReadDir(*root + *path)
	temp := make([]string, len(ft))
	for i, r := range ft {
		nt := r.Name()
		if r.Mode().IsDir() {
			if nt[len(nt)-1] == '/' {
				temp[i] = *path + nt
			} else {
				temp[i] = *path + nt + "/"
			}
		} else {
			temp[i] = *path + nt
		}
	}
	for _, r := range temp {
		if r[len(r)-1] == '/' {
			lstd.add(&r)
			getlist(root, &r, lstf, lstd)
		} else {
			lstf.add(&r)
		}
	}
}

// add only str list
type schunk struct {
	data []string
}

func (c *schunk) add(path *string) {
	c.data = append(c.data, *path)
}

// ! 실행 전 파라미터 4개를 맞춰야 함 !
type toolbox struct {
	Noerr  bool     // crc 오류 무시 여부
	Export string   // 결과 출력 위치
	Folder string   // 폴더 절대경로
	File   []string // 파일들 절대경로
}

func Init() toolbox {
	var temp toolbox
	temp.Noerr = false
	temp.Export = "./temp570"
	temp.Folder = ""
	temp.File = make([]string, 0)
	return temp
}

// 파일/폴더 표준절대경로화
func (self *toolbox) Abs() {
	self.Folder, _ = filepath.Abs(self.Folder)
	self.Folder = strings.Replace(self.Folder, "\\", "/", -1)
	if self.Folder[len(self.Folder)-1] != '/' {
		self.Folder = self.Folder + "/"
	}
	temp := make([]string, len(self.File))
	for i, r := range self.File {
		temp[i], _ = filepath.Abs(r)
		temp[i] = strings.Replace(temp[i], "\\", "/", -1)
	}
	self.File = temp
}

// 파일을 패키징, mode = "png"/"webp"/"nah"
func (self *toolbox) Zipfile(mode string) {
	var fakeh []byte // prehead + padding
	switch mode {
	case "png":
		fakeh = *picdt.Kz5png()
		tn := (16384 - len(fakeh)) % 1024
		temp := make([]byte, tn)
		fakeh = append(fakeh, temp...)
	case "webp":
		fakeh = *picdt.Kz5webp()
		tn := (16384 - len(fakeh)) % 1024
		temp := make([]byte, tn)
		fakeh = append(fakeh, temp...)
	default:
		fakeh = make([]byte, 0)
	}
	commonh := []byte("KSC5")  // common head
	subtypeh := []byte("KZIP") // subtype head
	h0 := fmt.Sprintf("folder = 0; file = %d", len(self.File))
	h1 := make([]string, len(self.File))
	for i, r := range self.File {
		ht := strings.Replace(r, "\\", "/", -1)
		h1[i] = ht[strings.LastIndex(ht, "/")+1:]
	}
	mainh := []byte(h0)             // main header
	hsize := *encode(len(mainh), 4) // main head size
	res := *crc32check(&mainh)      // reserved

	initfile(self.Export)
	f, _ := os.OpenFile(self.Export, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666) // write as f
	defer f.Close()
	f.Write(fakeh)
	f.Write(commonh)
	f.Write(subtypeh)
	f.Write(res)
	f.Write(hsize)
	f.Write(mainh)

	for i, r := range self.File {
		nmb := []byte(h1[i])        // name B
		nms := *encode(len(nmb), 8) // name len
		f.Write(nms)
		f.Write(nmb)

		fileinfo, _ := os.Stat(r)
		fsize := int(fileinfo.Size())
		fls := *encode(fsize, 8)
		f.Write(fls)
		t, _ := os.Open(r) // read as t
		defer t.Close()
		temp := make([]byte, 10485760)
		for j := 0; j < fsize/10485760; j++ {
			t.Read(temp)
			f.Write(temp)
		}
		temp = make([]byte, fsize%10485760) // file B
		if len(temp) != 0 {
			t.Read(temp)
			f.Write(temp)
		}
	}
}

// 폴더를 패키징, mode = "png"/"webp"/"nah"
func (self *toolbox) Zipfolder(mode string) {
	var fakeh []byte // prehead + padding
	switch mode {
	case "png":
		fakeh = *picdt.Kz5png()
		tn := (16384 - len(fakeh)) % 1024
		temp := make([]byte, tn)
		fakeh = append(fakeh, temp...)
	case "webp":
		fakeh = *picdt.Kz5webp()
		tn := (16384 - len(fakeh)) % 1024
		temp := make([]byte, tn)
		fakeh = append(fakeh, temp...)
	default:
		fakeh = make([]byte, 0)
	}
	commonh := []byte("KSC5")  // common head
	subtypeh := []byte("KZIP") // subtype head

	root := strings.Replace(self.Folder, "\\", "/", -1)
	if root[len(root)-1] == '/' { // ~/에서 /제거
		root = root[0 : len(root)-1]
	}
	fr0 := root[0 : strings.LastIndex(root, "/")+1]    // root의 상위 폴더 절대경로
	fr1 := root[strings.LastIndex(root, "/")+1:] + "/" // 패키징할 root
	var fl schunk
	var dl schunk
	dl.add(&fr1)
	getlist(&fr0, &fr1, &fl, &dl)
	wrfile := fl.data
	wrfolder := dl.data

	h0 := fmt.Sprintf("folder = %d; file = %d", len(wrfolder), len(wrfile))
	h1 := h0 + "\n" + strings.Join(wrfolder, "\n")
	mainh := []byte(h1)             // main header
	hsize := *encode(len(mainh), 4) // main head size
	res := *crc32check(&mainh)      // reserved

	initfile(self.Export)
	f, _ := os.OpenFile(self.Export, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666) // write as f
	defer f.Close()
	f.Write(fakeh)
	f.Write(commonh)
	f.Write(subtypeh)
	f.Write(res)
	f.Write(hsize)
	f.Write(mainh)

	for _, r := range wrfile {
		nmb := []byte(r)            // name B
		nms := *encode(len(nmb), 8) // name len
		f.Write(nms)
		f.Write(nmb)

		fileinfo, _ := os.Stat(fr0 + r)
		fsize := int(fileinfo.Size())
		fls := *encode(fsize, 8)
		f.Write(fls)
		t, _ := os.Open(fr0 + r) // read as t
		defer t.Close()
		temp := make([]byte, 10485760)
		for j := 0; j < fsize/10485760; j++ {
			t.Read(temp)
			f.Write(temp)
		}
		temp = make([]byte, fsize%10485760) // file B
		if len(temp) != 0 {
			t.Read(temp)
			f.Write(temp)
		}
	}
}

// 패키징 해제, path는 kzip파일 경로
func (self *toolbox) Unzip(path string) {
	fileinfo, _ := os.Stat(path)
	var temp []byte
	if fileinfo.Size() < 16384 {
		temp = make([]byte, fileinfo.Size())
	} else {
		temp = make([]byte, 16384)
	}
	f, _ := os.Open(path)
	f.Read(temp)
	f.Close()
	pos := findpos(&temp)
	if pos == -1 {
		panic("Not Valid KSC5")
	}

	expath, _ := filepath.Abs(self.Export)
	expath = strings.Replace(expath, "\\", "/", -1)
	if expath[len(expath)-1] != '/' {
		expath = expath + "/"
	}
	initfolder(expath)
	os.Mkdir(expath, os.ModePerm)

	f, _ = os.Open(path)
	defer f.Close()
	temp = make([]byte, pos+4)
	f.Read(temp)
	subtypeh := make([]byte, 4) // subtype head
	f.Read(subtypeh)
	res := make([]byte, 4) // reserved
	f.Read(res)
	temp = make([]byte, 4)
	f.Read(temp)
	fsize := decode(&temp)
	mainh := make([]byte, fsize)
	f.Read(mainh) // main header

	if !self.Noerr {
		cmp := []byte("KZIP")
		if !bequal(&subtypeh, &cmp) {
			panic("Not Valid KZIP")
		}
		cmp = *crc32check(&mainh)
		if !bequal(&res, &cmp) {
			panic("Broken Header")
		}
	}

	doc := strings.Split(string(mainh), "\n")
	if doc[len(doc)-1] == "" {
		doc = doc[0 : len(doc)-1]
	}
	infoline := doc[0] // file folder num
	if len(doc) == 1 {
		doc = make([]string, 0)
	} else {
		doc = doc[1:] // folders
	}

	foldernum := 0
	filenum := 0
	for _, r := range strings.Split(infoline, ";") {
		if strings.Contains(r, "folder") {
			r = strings.Replace(r, " ", "", -1)
			foldernum, _ = strconv.Atoi(r[strings.Index(r, "=")+1:])
		} else if strings.Contains(r, "file") {
			r = strings.Replace(r, " ", "", -1)
			filenum, _ = strconv.Atoi(r[strings.Index(r, "=")+1:])
		}
	}
	if foldernum != 0 {
		for _, r := range doc {
			os.Mkdir(expath+r, os.ModePerm)
		}
	}

	for i := 0; i < filenum; i++ {
		temp = make([]byte, 8)
		f.Read(temp)
		nms := decode(&temp)
		temp = make([]byte, nms)
		f.Read(temp)
		nm := string(temp)
		temp = make([]byte, 8)
		f.Read(temp)
		fls := decode(&temp)
		t, _ := os.OpenFile(expath+nm, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
		defer t.Close()
		temp = make([]byte, 10485760)
		for j := 0; j < fls/10485760; j++ {
			f.Read(temp)
			t.Write(temp)
		}
		temp = make([]byte, fls%10485760)
		if len(temp) != 0 {
			f.Read(temp)
			t.Write(temp)
		}
	}
}

// ========== ========== kzip5 end ========== ==========
