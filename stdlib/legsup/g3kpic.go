// test660 : stdlib5.legsup gen3kpic

package legsup

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"slices"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strings"
	"sync"

	"golang.org/x/image/bmp"
)

// gen3 kpic
type G3kpic struct {
	moldsize [2]int // n * m mold pic size
	moldhead []byte // mold bmp header
	moldbody []byte // mold bmp content

	Pcover bool // cover data with pic
}

// gen3 kpic encode picture (path : ~.bmp)
func (tbox *G3kpic) enc(path string, data []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	f, _ := kio.Open(path, "w")
	defer func() {
		f.Close()
		pconv(path, true)
	}()

	if tbox.Pcover {
		temp := make([]byte, 0)
		temp = append(temp, tbox.moldbody...)
		for i, r := range data {
			temp[2*i] = 16*(temp[2*i]/16) + r/16
			temp[2*i+1] = 16*(temp[2*i+1]/16) + r%16
		}
		kio.Write(f, tbox.moldhead)
		kio.Write(f, temp)
	} else {
		kio.Write(f, tbox.moldhead)
		kio.Write(f, data)
	}
}

// gen3 kpic decode picture (path : ~.png)
func (tbox *G3kpic) dec(path string, data *[]byte, wg *sync.WaitGroup) {
	defer wg.Done()
	pconv(path, false)
	f, _ := kio.Open(path[0:len(path)-3]+"bmp", "r")
	defer func() {
		f.Close()
		os.Remove(path[0:len(path)-3] + "bmp")
	}()
	bmpv, _ := kio.Read(f, -1)
	temp := bmpv[kobj.Decode(bmpv[10:14]):]

	if tbox.Pcover {
		tb := make([]byte, len(temp)/2)
		for i := 0; i < len(tb); i++ {
			tb[i] = 16*(temp[2*i]%16) + temp[2*i+1]%16
		}
		*data = tb
	} else {
		*data = temp
	}
}

// gen3 kpic init, set empty string to use basic pic, size -1/4n
func (tbox *G3kpic) Init(path string, row int, col int) error {
	tbox.Pcover = true
	var img image.Image
	var err error
	if path == "" {
		img, err = png.Decode(bytes.NewReader(G3mold()))
	} else {
		f, _ := kio.Open(path, "r")
		defer f.Close()
		img, err = png.Decode(f)
	}
	if err != nil {
		return err
	}

	if row < 0 || col < 0 {
		bound := img.Bounds()
		row = bound.Max.X
		col = bound.Max.Y
	} else {
		img = *presize(&img, row, col)
	}
	if row%4 != 0 || col%4 != 0 {
		return errors.New("mold should be 4N*4M size")
	}
	tbox.moldsize[0] = row
	tbox.moldsize[1] = col

	temp := bytes.NewBuffer(make([]byte, 0))
	bmp.Encode(temp, img)
	bmpv := temp.Bytes()
	hsize := kobj.Decode(bmpv[10:14])
	csize := kobj.Decode(bmpv[28:30])
	if csize != 24 {
		return errors.New("color should be 24 bit")
	}
	tbox.moldhead = bmpv[0:hsize]
	tbox.moldbody = bmpv[hsize:]
	return nil
}

// gen3 kpic detect pic names (name, num, err)
func (tbox *G3kpic) Detect(path string) (string, int, error) {
	path = kio.Abs(path)
	if path[len(path)-1] != '/' {
		return "", 0, errors.New("path should be folder")
	}
	chars := make([]string, 26)
	for i := 0; i < 26; i++ {
		chars[i] = string(rune(65 + i))
	}
	nums := make([]string, 10)
	for i := 0; i < 10; i++ {
		nums[i] = string(rune(48 + i))
	}
	infos, err := os.ReadDir(path)
	if err != nil {
		return "", 0, err
	}

	temp := make([]string, 0)
	for _, r := range infos {
		name := r.Name()
		if strings.Contains(name, ".") && len(name) > 8 {
			fr0 := name[0:strings.LastIndex(name, ".")]
			fr1 := name[strings.LastIndex(name, ".")+1:]
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
			if flag && fr1 == "png" {
				temp = append(temp, name)
			}
		}
	}

	if len(temp) == 0 {
		return "", 0, nil
	} else {
		name := temp[0][0:4]
		num := 0
		for slices.Contains(temp, fmt.Sprintf("%s%d.png", name, num)) {
			num = num + 1
		}
		return name, num, nil
	}
}

// gen3 kpic file pack, generate pic at exdir, (name, num)
func (tbox *G3kpic) Pack(tgt string, exdir string) (string, int) {
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
	temp := genrand(4)
	name := chars0[temp[0]%26] + chars0[temp[1]%26] + chars0[temp[2]%26]
	if tbox.Pcover {
		name = name + chars2[temp[0]%21]
	} else {
		name = name + chars1[temp[0]%5]
	}

	exdir = kio.Abs(exdir)
	num0 := kio.Size(tgt)
	num1 := tbox.moldsize[0] * tbox.moldsize[1] * 3
	if tbox.Pcover {
		num1 = num1 / 2
	}
	num2 := num0 / num1
	var wg sync.WaitGroup
	f, _ := kio.Open(tgt, "r")
	defer f.Close()
	count := 0
	var buf []byte

	for i := 0; i < num2/32; i++ {
		buf, _ = kio.Read(f, 32*num1)
		wg.Add(32)
		for j := 0; j < 32; j++ {
			go tbox.enc(fmt.Sprintf("%s%s%d.bmp", exdir, name, count), buf[j*num1:j*num1+num1], &wg)
			count = count + 1
		}
		wg.Wait()
	}

	if num2%32 != 0 {
		buf, _ = kio.Read(f, (num2%32)*num1)
		wg.Add(num2 % 32)
		for i := 0; i < num2%32; i++ {
			go tbox.enc(fmt.Sprintf("%s%s%d.bmp", exdir, name, count), buf[i*num1:i*num1+num1], &wg)
			count = count + 1
		}
		wg.Wait()
	}

	if num0%num1 != 0 {
		buf, _ = kio.Read(f, num0%num1)
		buf = append(buf, make([]byte, num1-num0%num1)...)
		wg.Add(1)
		go tbox.enc(fmt.Sprintf("%s%s%d.bmp", exdir, name, count), buf, &wg)
		count = count + 1
		wg.Wait()
	}
	return name, count
}

// gen3 kpic picture unpack, pic at tgtdir, generate file at path
func (tbox *G3kpic) Unpack(path string, tgtdir string, name string, num int) {
	tgtdir = kio.Abs(tgtdir)
	if name[3] == 'A' || name[3] == 'E' || name[3] == 'I' || name[3] == 'O' || name[3] == 'U' {
		tbox.Pcover = false
	} else {
		tbox.Pcover = true
	}

	f, _ := kio.Open(path, "w")
	defer f.Close()
	var wg sync.WaitGroup
	count := 0
	var buf [][]byte

	for i := 0; i < num/32; i++ {
		buf = make([][]byte, 32)
		wg.Add(32)
		for j := 0; j < 32; j++ {
			go tbox.dec(fmt.Sprintf("%s%s%d.png", tgtdir, name, count), &buf[j], &wg)
			count = count + 1
		}
		wg.Wait()
		for j := 0; j < 32; j++ {
			kio.Write(f, buf[j])
		}
	}

	if num%32 != 0 {
		buf = make([][]byte, num%32)
		wg.Add(num % 32)
		for i := 0; i < num%32; i++ {
			go tbox.dec(fmt.Sprintf("%s%s%d.png", tgtdir, name, count), &buf[i], &wg)
			count = count + 1
		}
		wg.Wait()
		for i := 0; i < num%32; i++ {
			kio.Write(f, buf[i])
		}
	}
}
