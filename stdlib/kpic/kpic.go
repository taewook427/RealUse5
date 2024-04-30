// test647 : stdlib5.kpic

package kpic

// go get "golang.org/x/crypto/sha3"
// go get "golang.org/x/exp/slices"
// go get "golang.org/x/image/bmp"
// go get "github.com/nfnt/resize"

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"stdlib5/kio"
	"stdlib5/kobj"
	"stdlib5/picdt"
	"strings"
	"sync"

	"github.com/nfnt/resize"
	"golang.org/x/crypto/sha3"
	"golang.org/x/exp/slices"
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
)

// pack bytes -> pic, init temp dir before "webp" pmode
func packpic(head []byte, cnvbody []byte, influx0 chan []byte, influx1 chan string, zmode int, pmode string, flag *sync.WaitGroup) {
	defer func() {
		recover()
		flag.Done()
	}()
	var div byte
	if zmode == 2 {
		div = 16
	} else {
		div = 4
	}
	for data := range influx0 {
		temp := make([]byte, len(cnvbody))
		copy(temp, cnvbody)
		for i, r := range data {
			offset := i * zmode
			adjust := zmode - 1
			for adjust >= 0 {
				temp[offset+adjust] = temp[offset+adjust] + r%div
				adjust = adjust - 1
				r = r / div
			}
		}

		bmpv := make([]byte, 0, len(head)+len(temp))
		bmpv = append(bmpv, head...)
		bmpv = append(bmpv, temp...)
		temp = nil
		path := <-influx1
		switch pmode {
		case "png":
			img, _ := bmp.Decode(bytes.NewReader(bmpv))
			bmpv = nil
			f, _ := kio.Open(path, "w")
			png.Encode(f, img)
			f.Close()
		case "webp":
			genwebp(bmpv, path)
		default:
			f, _ := kio.Open(path, "w")
			kio.Write(f, bmpv)
			f.Close()
		}
	}
}

// unpack pic -> bytes
func unpackpic(zmode int, pmode string, influx chan string, exflux chan []byte) {
	defer func() {
		recover()
		close(exflux)
	}()
	var div byte
	if zmode == 2 {
		div = 16
	} else {
		div = 4
	}
	for path := range influx {
		raw := readpdt(path, pmode)
		temp := make([]byte, len(raw)/zmode)
		for i, r := range temp {
			r = 0
			offset := zmode * i
			adjust := zmode - 1
			var reg byte = 1
			for adjust >= 0 {
				r = r + reg*(raw[offset+adjust]%div)
				adjust = adjust - 1
				reg = reg * div
			}
			temp[i] = r
		}
		exflux <- temp[8:]
	}
}

// examine pic by internal data, chan returns (name, curnum, maxnum)
func expic(zmode int, pmode string, influx chan string, exflux0 chan string, exflux1 chan int) {
	defer func() {
		recover()
		close(exflux0)
		close(exflux1)
	}()
	var div byte
	if zmode == 2 {
		div = 16
	} else {
		div = 4
	}
	for path := range influx {
		raw := readpdt(path, pmode)
		temp := make([]byte, 8)
		for i, r := range temp {
			offset := zmode * i
			adjust := zmode - 1
			var reg byte = 1
			for adjust >= 0 {
				r = r + reg*(raw[offset+adjust]%div)
				adjust = adjust - 1
				reg = reg * div
			}
			temp[i] = r
		}

		flag := true
		for _, r := range temp[0:4] {
			if r < 65 || r > 64+26 {
				exflux0 <- ""
				exflux1 <- -1
				exflux1 <- -1
				flag = false
				break
			}
		}
		if flag {
			name := string(temp[0:4])
			num0 := kobj.Decode(temp[4:6])
			num1 := kobj.Decode(temp[6:8])
			if num0 < num1 {
				exflux0 <- name
				exflux1 <- num0
				exflux1 <- num1
			} else {
				exflux0 <- ""
				exflux1 <- -1
				exflux1 <- -1
			}
		}
	}
}

// feed file paths
func feedpath(files []string, exflux []chan string, proc *float64) {
	defer func() {
		recover()
		for _, r := range exflux {
			close(r)
		}
	}()
	current := 0
	for i := 0; i < len(files)/32; i++ {
		for j := 0; j < 32; j++ {
			exflux[j] <- files[current]
			*proc = float64(current) / float64(len(files))
			current = current + 1
		}
	}
	for j := 0; j < len(files)%32; j++ {
		exflux[j] <- files[current]
		*proc = float64(current) / float64(len(files))
		current = current + 1
	}
}

// generate pic paths, feeds path influx chan
func genpath(folder string, name string, num int, pmode string, exflux []chan string, proc *float64) {
	defer func() {
		recover()
		for _, r := range exflux {
			close(r)
		}
	}()
	count := 0
	for i := 0; i < num/32; i++ {
		for j := 0; j < 32; j++ {
			exflux[j] <- fmt.Sprintf("%s%s%d.%s", folder, name, count, pmode)
			count = count + 1
			*proc = float64(count) / float64(num)
		}
	}
	for j := 0; j < num%32; j++ {
		exflux[j] <- fmt.Sprintf("%s%s%d.%s", folder, name, count, pmode)
		count = count + 1
		*proc = float64(count) / float64(num)
	}
}

// create webp file
func genwebp(bmpv []byte, path string) {
	img, _ := bmp.Decode(bytes.NewReader(bmpv))
	bmpv = nil
	f, _ := kio.Open(path+".png", "w")
	png.Encode(f, img)
	f.Close()
	cmd := exec.Command("./kpic5cwebp.exe", "-lossless", fmt.Sprintf("%s.png", path), "-o", path) // !! use "./kpic5cwebp" at linux!!
	cmd.Run()
	os.Remove(path + ".png")
}

// check valid cwebp binary
func chkwebp() bool {
	winb := []byte{59, 214, 94, 136, 168, 37, 22, 92, 89, 22, 177, 207, 171, 29, 37, 141, 100, 202, 203, 247, 85, 86, 72, 250, 69, 43, 39, 12, 206, 146, 158, 145, 186, 198, 177, 178, 81, 38, 62, 89, 235, 66, 24, 217, 121, 234, 170, 206, 57, 153, 191, 9, 212, 218, 234, 8, 137, 58, 227, 92, 226, 229, 235, 75}
	linuxb := []byte{25, 212, 222, 120, 125, 124, 185, 80, 245, 75, 152, 52, 246, 173, 25, 226, 206, 213, 230, 59, 137, 67, 81, 141, 173, 118, 63, 127, 43, 166, 172, 62, 62, 19, 146, 151, 188, 23, 46, 61, 142, 13, 108, 11, 189, 137, 162, 145, 138, 26, 185, 242, 117, 25, 103, 81, 165, 24, 104, 163, 235, 3, 206, 194}
	f, err := os.Open("./kpic5cwebp.exe")
	if err == nil {
		defer f.Close()
		temp, _ := kio.Read(f, -1)
		h := sha3.New512()
		h.Write(temp)
		return kio.Bequal(h.Sum(nil), winb)
	} else {
		t, err := os.Open("./kpic5cwebp")
		if err == nil {
			defer t.Close()
			temp, _ := kio.Read(t, -1)
			h := sha3.New512()
			h.Write(temp)
			return kio.Bequal(h.Sum(nil), linuxb)
		} else {
			return false
		}
	}
}

// read pic data, returns bmpv
func readpdt(path string, pmode string) []byte {
	var bmpv []byte
	switch pmode {
	case "png":
		f, _ := kio.Open(path, "r")
		defer f.Close()
		img, _ := png.Decode(f)
		temp := bytes.NewBuffer(make([]byte, 0))
		bmp.Encode(temp, img)
		bmpv = temp.Bytes()
	case "webp":
		f, _ := kio.Open(path, "r")
		defer f.Close()
		img, _ := webp.Decode(f)
		temp := bytes.NewBuffer(make([]byte, 0))
		bmp.Encode(temp, img)
		bmpv = temp.Bytes()
	default:
		f, _ := kio.Open(path, "r")
		defer f.Close()
		bmpv, _ = kio.Read(f, -1)
	}
	hsize := kobj.Decode(bmpv[10:14])
	return bmpv[hsize:]
}

type toolbox struct {
	moldsize [2]int
	moldhead []byte
	moldbody []byte

	Target string // conv target file/folder
	Export string // conv result file/folder
	Style  string // conv style, "webp"/"png"/"bmp"

	Proc float64 // progress, -1 : not started, 0~1 : working, 2 : finished
}

// init & set mold pic, (row, col) : -1/4n
func Initpic(path string, row int, col int) (toolbox, error) {
	var out toolbox
	out.moldsize[0] = 0
	out.moldsize[1] = 0
	out.moldhead = nil
	out.moldbody = nil
	out.Target = ""
	out.Export = ""
	out.Style = "webp"
	out.Proc = -1.0

	var img image.Image
	if path == "" {
		img, _ = webp.Decode(bytes.NewReader(picdt.Kp5webp()))
	} else {
		switch strings.ToLower(path[strings.LastIndex(path, ".")+1:]) {
		case "png":
			f, _ := kio.Open(path, "r")
			defer f.Close()
			img, _ = png.Decode(f)
		case "webp":
			f, _ := kio.Open(path, "r")
			defer f.Close()
			img, _ = webp.Decode(f)
		default:
			f, _ := kio.Open(path, "r")
			defer f.Close()
			img, _ = bmp.Decode(f)
		}
	}
	if row < 0 || col < 0 {
		bound := img.Bounds()
		row = bound.Max.X
		col = bound.Max.Y
	} else {
		img = resize.Resize(uint(row), uint(col), img, resize.Lanczos3)
	}
	if row%4 != 0 || col%4 != 0 {
		return out, errors.New("mold should be 4N*4M size")
	}
	out.moldsize[0] = row
	out.moldsize[1] = col

	temp := bytes.NewBuffer(make([]byte, 0))
	bmp.Encode(temp, img)
	bmpv := temp.Bytes()
	temp = nil
	hsize := kobj.Decode(bmpv[10:14])
	csize := kobj.Decode(bmpv[28:30])
	if csize != 24 {
		return out, errors.New("color should be 24 bit")
	}
	out.moldhead = bmpv[0:hsize]
	out.moldbody = bmpv[hsize:]
	return out, nil
}

// detect kpic info by target path -> (name, num, style)
func (tbox *toolbox) Detect() (string, int, string) {
	tbox.Target = kio.Abs(tbox.Target)
	flist := make([]string, 0)
	fs, _ := os.ReadDir(tbox.Target)
	for _, r := range fs {
		flist = append(flist, r.Name())
	}
	plist := make([]string, 0)
	chars := make([]string, 26)
	for i := 0; i < 26; i++ {
		chars[i] = string(rune(65 + i))
	}
	nums := make([]string, 10)
	for i := 0; i < 10; i++ {
		nums[i] = string(rune(48 + i))
	}

	for _, r := range flist {
		if strings.Contains(r, ".") && len(r) > 8 {
			fr0 := r[0:strings.LastIndex(r, ".")]
			fr1 := r[strings.LastIndex(r, ".")+1:]
			flag := true
			for _, l := range fr0[0:4] {
				if !slices.Contains(chars, string(l)) {
					flag = false
					break
				}
			}
			for _, l := range fr0[4:] {
				if !slices.Contains(nums, string(l)) {
					flag = false
					break
				}
			}
			if flag && slices.Contains([]string{"webp", "png", "bmp"}, fr1) {
				plist = append(plist, r)
			}
		}
	}

	if len(plist) == 0 {
		return "", 0, ""
	} else {
		fr := plist[0]
		name := fr[0:4]
		num := 0
		style := fr[strings.LastIndex(fr, ".")+1:]
		for slices.Contains(plist, fmt.Sprintf("%s%d.%s", name, num, style)) {
			num = num + 1
		}
		return name, num, style
	}
}

// pack file to pic, zmode = 2/4
func (tbox *toolbox) Pack(zmode int) (string, int) {
	defer func() { tbox.Proc = 2.0 }()
	tbox.Proc = -1.0
	os.RemoveAll(tbox.Export) // !!! automatically clear export folder !!!
	os.Mkdir(tbox.Export, os.ModePerm)
	tbox.Export = kio.Abs(tbox.Export)
	fsize := kio.Size(tbox.Target)
	csize := tbox.moldsize[0]*tbox.moldsize[1]*3/zmode - 8
	if tbox.Style == "webp" && !chkwebp() {
		return "", 0 // !! webp encoder error !!
	}

	chars0 := make([]string, 26)
	for i := 0; i < 26; i++ {
		chars0[i] = string(rune(65 + i))
	}
	chars1 := []string{"A", "E", "I", "O", "U"}
	chars2 := make([]string, 0)
	for _, r := range chars0 {
		if !slices.Contains(chars1, r) {
			chars2 = append(chars2, r)
		}
	}
	temp := make([]byte, 4)
	rand.Read(temp)
	name := chars0[temp[1]%26] + chars0[temp[2]%26] + chars0[temp[3]%26]
	if zmode == 2 {
		name = chars2[temp[0]%21] + name
	} else {
		name = chars1[temp[0]%5] + name
	}
	var num0 int // pic num
	var num1 int // added bytes length
	var num2 int // x32 repeating
	var num3 int // 1~32 last cycle
	if fsize%csize == 0 {
		num0 = fsize / csize
		num1 = 0
	} else {
		num0 = fsize/csize + 1
		num1 = csize - fsize%csize
	}
	num2 = (num0 - 1) / 32
	num3 = (num0-1)%32 + 1
	head0 := []byte(name)         // name
	head1 := kobj.Encode(num0, 2) // maxnum header

	tbox.Proc = 0.0
	port0 := make([]chan []byte, 32)
	port1 := make([]chan string, 32)
	for i := 0; i < 32; i++ {
		port0[i] = make(chan []byte, 4)
		port1[i] = make(chan string, 4)
	}
	cnvbody := make([]byte, len(tbox.moldbody))
	for i, r := range tbox.moldbody {
		if zmode == 2 {
			cnvbody[i] = 16 * (r / 16)
		} else {
			cnvbody[i] = 4 * (r / 4)
		}
	}

	var wg sync.WaitGroup
	wg.Add(32)
	go genpath(tbox.Export, name, num0, tbox.Style, port1, &tbox.Proc)
	for j := 0; j < 32; j++ {
		go packpic(tbox.moldhead, cnvbody, port0[j], port1[j], zmode, tbox.Style, &wg)
	}
	f, _ := kio.Open(tbox.Target, "r")
	defer f.Close()
	current := 0
	for i := 0; i < num2; i++ {
		for j := 0; j < 32; j++ {
			data := append(make([]byte, 0), head0...)
			data = append(data, kobj.Encode(current, 2)...)
			data = append(data, head1...)
			tb, _ := kio.Read(f, csize)
			data = append(data, tb...)
			tb = nil
			port0[j] <- data
			current = current + 1
		}
	}

	if num3 != 1 {
		for j := 0; j < num3-1; j++ {
			data := append(make([]byte, 0), head0...)
			data = append(data, kobj.Encode(current, 2)...)
			data = append(data, head1...)
			tb, _ := kio.Read(f, csize)
			data = append(data, tb...)
			tb = nil
			port0[j] <- data
			current = current + 1
		}
	}

	data := append(make([]byte, 0), head0...)
	data = append(data, kobj.Encode(current, 2)...)
	data = append(data, head1...)
	tb, _ := kio.Read(f, csize-num1)
	data = append(data, tb...)
	tb = make([]byte, num1)
	rand.Read(tb)
	data = append(data, tb...)
	port0[num3-1] <- data

	for i := 0; i < 32; i++ {
		close(port0[i])
	}
	wg.Wait()
	return name, num0
}

// unpack pic to file with name, num
func (tbox *toolbox) Unpack(name string, num int) {
	defer func() { tbox.Proc = 2.0 }()
	tbox.Proc = -1.0
	tbox.Target = kio.Abs(tbox.Target)
	var zmode int
	switch name[0] {
	case 'A':
		zmode = 4
	case 'E':
		zmode = 4
	case 'I':
		zmode = 4
	case 'O':
		zmode = 4
	case 'U':
		zmode = 4
	default:
		zmode = 2
	}

	tbox.Proc = 0.0
	port0 := make([]chan string, 32)
	port1 := make([]chan []byte, 32)
	for i := 0; i < 32; i++ {
		port0[i] = make(chan string, 4)
		port1[i] = make(chan []byte, 4)
	}
	go genpath(tbox.Target, name, num, tbox.Style, port0, &tbox.Proc)
	for i := 0; i < 32; i++ {
		go unpackpic(zmode, tbox.Style, port0[i], port1[i])
	}

	f, _ := kio.Open(tbox.Export, "w")
	defer f.Close()
	for data := range port1[31] {
		for i := 0; i < 31; i++ {
			kio.Write(f, <-port1[i])
		}
		kio.Write(f, data)
	}
	for i := 0; i < 31; i++ {
		data, ok := <-port1[i]
		if ok {
			kio.Write(f, data)
		}
	}
}

// restore and change name by internal data
func (tbox *toolbox) Restore(files []string, zmode int) (string, int) {
	defer func() { tbox.Proc = 2.0 }()
	tbox.Proc = -1.0
	for i, r := range files {
		files[i] = kio.Abs(r)
	}

	tbox.Proc = 0.0
	names := make([]string, 0)
	curnums := make([]int, 0)
	maxnums := make([]int, 0)
	port0 := make([]chan string, 32)
	port1 := make([]chan string, 32)
	port2 := make([]chan int, 32)
	for i := 0; i < 32; i++ {
		port0[i] = make(chan string, 4)
		port1[i] = make(chan string, 4)
		port2[i] = make(chan int, 8)
	}
	go feedpath(files, port0, &tbox.Proc)
	for i := 0; i < 32; i++ {
		go expic(zmode, tbox.Style, port0[i], port1[i], port2[i])
	}

	for data := range port1[31] {
		for i := 0; i < 31; i++ {
			names = append(names, <-port1[i])
			curnums = append(curnums, <-port2[i])
			maxnums = append(maxnums, <-port2[i])
		}
		names = append(names, data)
		curnums = append(curnums, <-port2[31])
		maxnums = append(maxnums, <-port2[31])
	}
	for i := 0; i < 31; i++ {
		data, ok := <-port1[i]
		if ok {
			names = append(names, data)
			curnums = append(curnums, <-port2[i])
			maxnums = append(maxnums, <-port2[i])
		}
	}

	name := ""
	maxnum := 0
	for i, r := range files {
		if curnums[i] >= 0 && maxnums[i] > 0 {
			os.Rename(r, fmt.Sprintf("%s%s%d.%s", r[0:strings.LastIndex(r, "/")+1], names[i], curnums[i], tbox.Style))
			name = names[i]
			if maxnum < maxnums[i] {
				maxnum = maxnums[i]
			}
		}
	}
	return name, maxnum
}
