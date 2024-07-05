// test675 : stdlib5.kpkg

package kpkg

// go get "golang.org/x/crypto/sha3"

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"stdlib5/kdb"
	"stdlib5/kio"
	"stdlib5/ksc"
	"stdlib5/ksign"
	"stdlib5/kzip"
	"time"

	"golang.org/x/crypto/sha3"
)

// init toolbox, osnum 0:any 1:windows 2:linuxmint
func Initkpkg(osnum int) toolbox {
	var out toolbox
	out.Osnum = osnum
	return out
}

type toolbox struct {
	Public   string  // RSA public key, empty if not using
	Private  string  // RSA private key, empty if not using
	Name     string  // package name
	Version  float64 // package version
	Text     string  // package explanation
	Rel_date string  // package release date
	Dwn_date string  // package download date
	Osnum    int     // any : 0,  windows : 1, linux mint : 2
}

// pack package by osnum/dirpaths, gen file at respath
func (tbox *toolbox) Pack(osnums []int, dirpaths []string, respath string) error {
	// setting required : public, private, name, version, text
	if len(osnums) != len(dirpaths) {
		return errors.New("invalid package order")
	}
	tbox.Rel_date = time.Now().Format("20060102")

	// chunk0
	ts := ""
	for i, r := range osnums {
		ts = ts + fmt.Sprintf("pkg%d = %d\n", i, r)
	}
	w0 := kdb.Initkdb()
	w0.Read(ts + "name = 0\nversion = 0\ntext = 0\nrelease = 0")
	w0.Fix("name", tbox.Name)
	w0.Fix("version", tbox.Version)
	w0.Fix("text", tbox.Text)
	w0.Fix("release", tbox.Rel_date)
	ts, err := w0.Write()
	if err != nil {
		return err
	}
	chunk0 := []byte(ts) // chunk0 : basic info

	// init temp dir, pack with kzip
	dirinit(false)
	for i, r := range dirpaths {
		err = kzip.Dozip([]string{r}, "webp", fmt.Sprintf("./temp675/%d.webp", i))
		if err != nil {
			return err
		}
	}

	var chunk1 []byte
	if tbox.Private != "" {
		w1 := kdb.Initkdb()
		ts = ""
		for i := range osnums {
			ts = ts + fmt.Sprintf("pkg%d = 0\n", i)
		}
		ts = ts + "public = 0"
		w1.Read(ts)

		for i := range osnums {
			f, _ := kio.Open(fmt.Sprintf("./temp675/%d.webp", i), "r")
			data, _ := kio.Read(f, -1)
			f.Close()
			h := sha3.New512()
			h.Write(append(append(make([]byte, 0), chunk0...), data...))
			hvalue := h.Sum(nil)
			enc, err := ksign.Sign(tbox.Private, hvalue)
			if err != nil {
				return errors.New("invalid signkey")
			}
			w1.Fix(fmt.Sprintf("pkg%d", i), enc)
		}
		w1.Fix("public", tbox.Public)
		ts, err = w1.Write()
		if err != nil {
			return err
		}
		chunk1 = []byte(ts)
	}

	// make basic structure
	w2 := ksc.Initksc()
	w2.Path = respath
	w2.Prehead = genwebp()
	w2.Subtype = []byte("KPKG")
	w2.Reserved = append(ksc.Crc32hash(chunk0), ksc.Crc32hash(chunk1)...)
	w2.Writef()
	w2.Linkf(chunk0)
	w2.Linkf(chunk1)

	// add kzip files
	for i := range osnums {
		err = w2.Addf(fmt.Sprintf("./temp675/%d.webp", i))
		if err != nil {
			return err
		}
	}
	dirinit(false)
	return nil
}

// unpack package by os, update internal data, gen folder at temp675/temp/, returns dir path
func (tbox *toolbox) Unpack(path string) (string, error) {
	// setting required : osnum
	dirinit(false)
	w0 := ksc.Initksc()
	w0.Predetect = true
	w0.Path = path
	w0.Readf()
	if !kio.Bequal(w0.Subtype, []byte("KPKG")) {
		return "", errors.New("invalid package")
	}
	f, _ := kio.Open(path, "r")
	f.Seek(int64(w0.Chunkpos[0]+8), 0)
	chunk0 := make([]byte, w0.Chunksize[0])
	f.Read(chunk0) // basic data
	f.Seek(int64(w0.Chunkpos[1]+8), 0)
	chunk1 := make([]byte, w0.Chunksize[1])
	f.Read(chunk1) // sign data
	f.Close()
	if !kio.Bequal(w0.Reserved, append(ksc.Crc32hash(chunk0), ksc.Crc32hash(chunk1)...)) {
		return "", errors.New("invalid package")
	}

	c0data := kdb.Initkdb() // get info chunk
	c0data.Read(string(chunk0))
	c1valid := true
	if len(chunk1) == 0 {
		c1valid = false
	}
	c1data := kdb.Initkdb()
	if c1valid {
		c1data.Read(string(chunk1))
	}

	tv, _ := c0data.Get("name")
	tbox.Name = tv.Dat6
	tv, _ = c0data.Get("version")
	tbox.Version = tv.Dat3
	tv, _ = c0data.Get("text")
	tbox.Text = tv.Dat6
	tv, _ = c0data.Get("release")
	tbox.Rel_date = tv.Dat6
	tbox.Dwn_date = time.Now().Format("20060102")

	// get os numbers
	osnums := make([]int, 0)
	ti := 0
	tv, ta := c0data.Get(fmt.Sprintf("pkg%d", ti))
	for len(ta) != 0 {
		osnums = append(osnums, tv.Dat2)
		ti = ti + 1
		tv, ta = c0data.Get(fmt.Sprintf("pkg%d", ti))
	}
	if !slices.Contains(osnums, 0) && !slices.Contains(osnums, tbox.Osnum) {
		return "", errors.New("OSnotsupportPKG")
	}
	var pos int
	if slices.Contains(osnums, 0) {
		pos = slices.Index(osnums, 0)
	} else {
		pos = slices.Index(osnums, tbox.Osnum)
	}

	// extract chunk corresponds osnum
	f, _ = kio.Open(path, "r")
	t, _ := kio.Open("./temp675/temp.webp", "w")
	f.Seek(int64(w0.Chunkpos[pos+2]+8), 0)
	size := w0.Chunksize[pos+2]
	tb := make([]byte, 10485760)
	for i := 0; i < size/10485760; i++ {
		f.Read(tb)
		t.Write(tb)
	}
	if size%10485760 != 0 {
		tb = make([]byte, size%10485760)
		f.Read(tb)
		t.Write(tb)
	}
	t.Close()
	f.Close()

	// check sign
	tbox.Public = ""
	if c1valid {
		f, _ := kio.Open("./temp675/temp.webp", "r")
		tb, _ = kio.Read(f, -1)
		f.Close()
		h := sha3.New512()
		h.Write(append(append(make([]byte, 0), chunk0...), tb...))
		hvalue := h.Sum(nil)

		tv, _ = c1data.Get(fmt.Sprintf("pkg%d", pos))
		enc := tv.Dat5
		tv, _ = c1data.Get("public")
		tbox.Public = tv.Dat6
		flag, _ := ksign.Verify(tbox.Public, enc, hvalue)
		if !flag {
			return "", errors.New("invalid sign")
		}
	}

	// unpack package
	os.Mkdir("./temp675/temp", os.ModePerm)
	err := kzip.Unzip("./temp675/temp.webp", "./temp675/temp", true)
	if err != nil {
		return "", err
	}
	tl, _ := os.ReadDir("./temp675/temp")
	pkgpath := kio.Abs("./temp675/temp/" + tl[0].Name())
	f, _ = kio.Open(pkgpath+"_ST5_VERSION.txt", "w")
	w1 := kdb.Initkdb()
	w1.Read("name = 0\nversion = 0\ntext = 0\nrelease = 0\ndownload = 0")
	w1.Fix("name", tbox.Name)
	w1.Fix("version", tbox.Version)
	w1.Fix("text", tbox.Text)
	w1.Fix("release", tbox.Rel_date)
	w1.Fix("download", tbox.Dwn_date)
	ts, _ := w1.Write()
	kio.Write(f, []byte(ts))
	f.Close()
	return pkgpath, nil
}

// init dir True : 674, False : 675
func dirinit(mode bool) {
	var name string
	if mode {
		name = "./temp674"
	} else {
		name = "./temp675"
	}
	os.RemoveAll(name)
	os.Mkdir(name, os.ModePerm)
}

func genwebp() []byte {
	var temp []byte
	temp = append(temp, 82, 73, 70, 70, 150, 2, 0, 0, 87, 69, 66, 80, 86, 80, 56, 32, 138, 2, 0, 0, 80, 15, 0, 157, 1, 42, 64, 0, 64, 0, 62, 109, 46, 147, 70, 164, 34, 161, 161, 36)
	temp = append(temp, 14, 216, 128, 13, 137, 106, 0, 192, 212, 100, 65, 71, 178, 125, 128, 126, 153, 244, 209, 239, 196, 84, 222, 246, 218, 190, 170, 246, 212, 121, 128, 253, 110, 253, 128, 247, 124, 254, 229, 234)
	temp = append(temp, 3, 161, 155, 214, 3, 208, 3, 246, 51, 211, 47, 217, 91, 252, 197, 125, 31, 122, 72, 96, 204, 53, 143, 3, 208, 135, 59, 63, 72, 17, 33, 194, 145, 53, 162, 81, 212, 224, 57, 6)
	temp = append(temp, 180, 93, 173, 106, 25, 75, 124, 3, 104, 190, 248, 100, 49, 225, 97, 58, 255, 249, 133, 14, 32, 0, 90, 139, 1, 139, 242, 129, 37, 48, 96, 0, 68, 70, 76, 87, 76, 6, 51, 85)
	temp = append(temp, 226, 51, 135, 239, 74, 226, 254, 12, 253, 205, 4, 178, 34, 239, 192, 181, 127, 240, 167, 11, 54, 203, 176, 136, 139, 227, 9, 20, 33, 31, 118, 224, 92, 217, 143, 209, 255, 253, 108, 251)
	temp = append(temp, 221, 217, 146, 175, 173, 209, 154, 83, 119, 15, 223, 255, 191, 179, 121, 203, 48, 42, 190, 188, 223, 78, 53, 8, 169, 226, 128, 205, 204, 140, 214, 15, 58, 145, 27, 10, 239, 133, 85, 7)
	temp = append(temp, 8, 77, 108, 103, 88, 126, 169, 201, 160, 185, 97, 247, 17, 104, 223, 92, 128, 166, 25, 254, 134, 53, 154, 222, 79, 117, 238, 109, 161, 205, 159, 35, 121, 182, 45, 16, 5, 29, 20, 152)
	temp = append(temp, 199, 6, 252, 51, 89, 36, 121, 19, 37, 96, 4, 76, 24, 187, 34, 197, 164, 175, 134, 226, 136, 3, 90, 233, 164, 192, 214, 125, 134, 69, 163, 246, 57, 52, 103, 158, 206, 126, 114, 40)
	temp = append(temp, 81, 192, 125, 160, 120, 23, 203, 41, 197, 93, 244, 91, 6, 175, 90, 174, 228, 74, 152, 221, 229, 145, 202, 210, 107, 123, 20, 45, 229, 162, 175, 93, 163, 144, 255, 105, 239, 7, 62, 12)
	temp = append(temp, 149, 211, 35, 108, 239, 69, 240, 48, 52, 72, 204, 175, 41, 25, 18, 63, 109, 97, 137, 34, 204, 180, 34, 52, 43, 95, 219, 193, 246, 230, 29, 23, 115, 253, 56, 132, 233, 252, 71, 208)
	temp = append(temp, 140, 2, 253, 48, 202, 101, 200, 182, 250, 33, 162, 243, 14, 31, 105, 60, 64, 127, 205, 158, 73, 97, 250, 189, 29, 38, 251, 98, 232, 135, 63, 43, 140, 150, 114, 42, 159, 195, 213, 46)
	temp = append(temp, 234, 64, 177, 159, 50, 180, 138, 32, 254, 234, 245, 245, 83, 210, 220, 225, 233, 152, 232, 79, 210, 204, 42, 159, 194, 85, 68, 193, 96, 6, 241, 89, 12, 46, 186, 190, 124, 95, 56, 163)
	temp = append(temp, 121, 118, 172, 252, 217, 25, 15, 199, 216, 255, 250, 181, 35, 92, 45, 229, 174, 151, 188, 26, 227, 210, 205, 1, 121, 236, 188, 9, 142, 198, 109, 141, 200, 43, 148, 188, 128, 19, 184, 127)
	temp = append(temp, 27, 210, 222, 1, 39, 203, 194, 253, 223, 197, 206, 18, 140, 130, 149, 51, 55, 86, 138, 18, 240, 205, 152, 108, 24, 226, 208, 157, 12, 103, 249, 131, 73, 185, 179, 32, 202, 232, 3, 86)
	temp = append(temp, 108, 33, 100, 175, 37, 210, 201, 56, 0, 101, 10, 193, 144, 126, 143, 107, 38, 239, 101, 252, 127, 48, 251, 202, 222, 197, 177, 6, 35, 116, 184, 182, 128, 237, 115, 13, 85, 22, 121, 86)
	temp = append(temp, 147, 159, 100, 43, 152, 164, 160, 202, 62, 22, 205, 126, 153, 242, 211, 49, 94, 37, 74, 255, 86, 39, 35, 56, 137, 146, 216, 98, 4, 52, 55, 25, 21, 123, 117, 19, 38, 253, 101, 254)
	temp = append(temp, 57, 236, 172, 119, 17, 90, 36, 187, 105, 61, 127, 36, 127, 168, 238, 144, 199, 34, 162, 56, 252, 119, 217, 78, 207, 225, 131, 64, 0, 0)
	temp = append(temp, make([]byte, 1024-len(temp))...)
	return temp
}
