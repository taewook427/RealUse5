// test643 : stdlib5.kzip

package kzip

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"stdlib5/kio"
	"stdlib5/ksc"
	"stdlib5/picdt"
	"strings"
)

// check if path is folder
func isdir(path string) bool {
	temp, err := os.Stat(path)
	if err == nil {
		return temp.IsDir()
	} else {
		return false
	}
}

// get folder/file names from tgts
func getlist(tgts []string) ([]string, []string, []string) {
	for i, r := range tgts {
		temp, _ := filepath.Abs(r)
		temp = strings.Replace(temp, "\\", "/", -1)
		if isdir(temp) && temp[len(temp)-1] != '/' {
			temp = temp + "/"
		}
		tgts[i] = temp
	}

	lstd := make([]string, 0)
	lstf := make([]string, 0)
	absf := make([]string, 0)
	count := 0

	for _, r := range tgts {
		if r[len(r)-1] == '/' {
			temp := ""
			if strings.Count(r, "/") == 1 {
				temp = fmt.Sprintf("LargeVolume%d/", count)
				count = count + 1
			} else {
				temp = r[0 : len(r)-1]
				temp = temp[strings.LastIndex(temp, "/")+1:] + "/"
			}
			getsub(r, temp, &lstd, &lstf, &absf)
		} else {
			lstf = append(lstf, r[strings.LastIndex(r, "/")+1:])
			absf = append(absf, r)
		}
	}
	return lstd, lstf, absf
}

// path : stdpath(~/) of root, prefix : ("root/")
func getsub(path string, prefix string, lstdp *[]string, lstfp *[]string, absfp *[]string) {
	*lstdp = append(*lstdp, prefix)
	fs, _ := os.ReadDir(path)
	temp := make([]string, len(fs))
	for i, r := range fs {
		rr := r.Name()
		if isdir(path+rr) && rr[len(rr)-1] != '/' {
			temp[i] = rr + "/"
		} else {
			temp[i] = rr
		}
	}

	for _, r := range temp {
		if r[len(r)-1] == '/' {
			getsub(path+r, prefix+r, lstdp, lstfp, absfp)
		} else {
			*lstfp = append(*lstfp, prefix+r)
			*absfp = append(*absfp, path+r)
		}
	}
}

// zip. tgts : folder/file path, mode ("webp", "png", ""), path ("" -> "./temp570.webp")
func Dozip(tgts []string, mode string, path string) (ferr error) {
	defer func() {
		errsign := recover()
		if errsign != nil {
			ferr = fmt.Errorf("%s", errsign)
		}
	}()
	lstd, lstf, absf := getlist(tgts)
	header := []byte(strings.Join(lstd, "\n"))

	writer := ksc.Initksc()
	switch mode {
	case "webp":
		writer.Prehead = picdt.Kz5webp()
		writer.Prehead = append(writer.Prehead, make([]byte, (16384-len(writer.Prehead))%512)...)
	case "png":
		writer.Prehead = picdt.Kz5png()
		writer.Prehead = append(writer.Prehead, make([]byte, (16384-len(writer.Prehead))%512)...)
	default:
		writer.Prehead = make([]byte, 0)
	}
	writer.Subtype = []byte("kzip")
	writer.Reserved = append(ksc.Crc32hash(header), 0, 0, 0, 0)
	if path == "" {
		writer.Path = "./temp570.webp"
	} else {
		writer.Path = path
	}
	os.RemoveAll(writer.Path)

	writer.Writef()
	writer.Linkf(header)
	for i, r := range lstf {
		writer.Linkf([]byte(r))
		writer.Addf(absf[i])
	}
	writer.Addf("")
	return nil
}

// unzip. kzip file path, export folder ("" -> "./temp570/"), check subtype/crc error
func Unzip(path string, export string, chkerr bool) error {
	if export == "" {
		export = "./temp570/"
	}
	export = strings.Replace(export, "\\", "/", -1)
	if export[len(export)-1] != '/' {
		export = export + "/"
	}
	var err error
	err = os.RemoveAll(export)
	if err != nil {
		return err
	}
	err = os.Mkdir(export, os.ModePerm)
	if err != nil {
		return err
	}

	reader := ksc.Initksc()
	reader.Path = path
	reader.Predetect = true
	err = reader.Readf()
	if err != nil {
		return err
	}

	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	f.Seek(int64(reader.Chunkpos[0]+8), 0)
	header, _ := kio.Read(f, reader.Chunksize[0])
	if chkerr {
		if !kio.Bequal(reader.Subtype, []byte("kzip")) {
			return errors.New("InvalidKZIP")
		}
		if !kio.Bequal(reader.Reserved[0:4], ksc.Crc32hash(header)) {
			return errors.New("InvalidCRC")
		}
		if len(reader.Chunkpos)%2 != 1 {
			return errors.New("InvalidChunk")
		}
	}

	var lstd []string
	if len(header) == 0 {
		lstd = make([]string, 0)
	} else {
		lstd = strings.Split(string(header), "\n")
	}
	for _, r := range lstd {
		err = os.Mkdir(export+r, os.ModePerm)
		if err != nil {
			return err
		}
	}

	for i := 0; i < (len(reader.Chunkpos)-1)/2; i++ {
		pos := 2*i + 1
		f.Seek(int64(reader.Chunkpos[pos]+8), 0)
		nmb, _ := kio.Read(f, reader.Chunksize[pos])
		name := string(nmb)
		size := reader.Chunksize[pos+1]
		f.Seek(int64(reader.Chunkpos[pos+1]+8), 0)

		t, err := kio.Open(export+name, "w")
		if err != nil {
			return err
		}
		for j := 0; j < size/10485760; j++ {
			temp, _ := kio.Read(f, 10485760)
			kio.Write(t, temp)
		}
		if size%10485760 != 0 {
			temp, _ := kio.Read(f, size%10485760)
			kio.Write(t, temp)
		}
		t.Close()
	}
	return nil
}
