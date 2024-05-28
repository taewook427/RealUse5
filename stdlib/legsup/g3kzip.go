// test658 : stdlib5.legsup gen3kzip

package legsup

import (
	"bytes"
	"errors"
	"os"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strings"
)

// gen3 kzip file add
func g3fadd(source string, dist *os.File) {
	size := kio.Size(source)
	f, _ := kio.Open(source, "r")
	defer f.Close()
	for i := 0; i < size/10485760; i++ {
		data, _ := kio.Read(f, 10485760)
		kio.Write(dist, data)
	}
	if size%10485760 != 0 {
		data, _ := kio.Read(f, size%10485760)
		kio.Write(dist, data)
	}
}

// gen3 kzip file subtract
func g3fsub(source *os.File, dist string, size int) {
	f, _ := kio.Open(dist, "w")
	defer f.Close()
	for i := 0; i < size/10485760; i++ {
		temp, _ := kio.Read(source, 10485760)
		kio.Write(f, temp)
	}
	if size%10485760 != 0 {
		temp, _ := kio.Read(source, size%10485760)
		kio.Write(f, temp)
	}
}

// gen3 kzip making folder
func g3mkdir(data string) error {
	var db G3kdb
	err := db.Read(data)
	if err != nil {
		return err
	}
	temp := db.Locate("folders#data")
	if temp == nil {
		return errors.New("structure data not exist")
	}
	tgt := temp.Data
	for tgt != nil {
		os.Mkdir("./"+strings.Replace(tgt.StrV, "\\", "/", -1), os.ModePerm)
		tgt = tgt.Next
	}
	return nil
}

// gen3 kzip get file/folder info
func g3fdinfo(root string, path string, folders *[]string, files *[]string) {
	// root : */, path : *, folders/files : actual names to be written
	*folders = append(*folders, path)
	dirs, _ := os.ReadDir(root + path)
	for _, r := range dirs {
		if r.IsDir() {
			temp := path + "/" + r.Name()
			if temp[len(temp)-1] == '/' || temp[len(temp)-1] == '\\' {
				temp = temp[0 : len(temp)-1]
			}
			g3fdinfo(root, temp, folders, files)
		} else {
			*files = append(*files, path+"/"+r.Name())
		}
	}
}

// gen3 kzip gen subhead, (type 1 + size 8 + name 247)
func g3gensubhead(folders []string, files []string, fsize []int, winsign bool) ([][]byte, []byte) {
	output := make([][]byte, len(files)+1)
	var db G3kdb
	db.Read("[folders] { [data] {\"\"} }")
	db.Zipexp = true
	db.Zipstr = true
	temp := db.Locate("folders#data")
	temp.Data.StrV = folders[0]
	for i := 1; i < len(folders); i++ {
		var addnode G3data
		addnode.Vtype = 's'
		if winsign {
			addnode.StrV = strings.Replace(folders[i], "/", "\\", -1)
		} else {
			addnode.StrV = folders[i]
		}
		temp.Data.Append(&addnode)
	}
	ts := db.Write()
	tb := []byte(ts)
	output[0] = append(append([]byte("S"), kobj.Encode(len(tb), 8)...), []byte(strings.Repeat(" ", 247))...)

	for i, r := range files {
		tb = []byte(r)
		ts = r + strings.Repeat(" ", 247-len(tb))
		if winsign {
			ts = strings.Replace(ts, "/", "\\", -1)
		}
		output[i+1] = append(append([]byte("F"), kobj.Encode(fsize[i], 8)...), []byte(ts)...)
	}
	return output, []byte(db.Write())
}

// gen3 kzip
type G3kzip struct {
	Prehead  []byte   // fake header (png + padding 1024nB)
	Header   []byte   // header 18B
	Chunkpos []int    // positions of chunk start (subhead + data)
	Subhead  [][]byte // subheads
	Winsign  bool     // use backslash to sign folder
}

// gen3 kzip init
func (tbox *G3kzip) Init() {
	tbox.Prehead = G3zip()
	tbox.Prehead = append(tbox.Prehead, make([]byte, 1024-len(tbox.Prehead)%1024)...)
	tbox.Header = nil
	tbox.Chunkpos = nil
	tbox.Subhead = nil
	tbox.Winsign = true
}

// gen3 kzip files pack (tgt : files, path : output file)
func (tbox *G3kzip) Packf(tgt []string, path string) error {
	for i, r := range tgt {
		tgt[i] = kio.Abs(r)
		if tgt[i][len(tgt[i])-1] == '/' {
			return errors.New("cannot zip folder")
		}
	}
	f, err := kio.Open(path, "w")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}

	rawlist0 := make([]string, len(tgt))
	rawlist1 := make([]int, len(tgt))
	for i, r := range tgt {
		rawlist1[i] = kio.Size(r)
		rawlist0[i] = "temp261" + r[strings.LastIndex(r, "/"):]
	}
	sub0, sub1 := g3gensubhead([]string{"temp261"}, rawlist0, rawlist1, tbox.Winsign)
	sub2 := bytes.Join(sub0, nil)
	tbox.Header = append(append([]byte("KTS2"), make([]byte, 2)...), kobj.Encode(len(sub0), 3)...)
	tbox.Header = append(tbox.Header, 1, 8, 247, 0, 0)
	temp := hash32(sub2)
	for i := 3; i >= 0; i-- {
		tbox.Header = append(tbox.Header, temp[i])
	}

	kio.Write(f, tbox.Prehead)
	kio.Write(f, tbox.Header)
	kio.Write(f, sub0[0])
	kio.Write(f, sub1)
	for i, r := range tgt {
		kio.Write(f, sub0[i+1])
		g3fadd(r, f)
	}
	return nil
}

// gen3 kzip folder pack (tgt : folder, path : output file)
func (tbox *G3kzip) Packd(tgt string, path string) error {
	tgt = kio.Abs(tgt)
	if tgt[len(tgt)-1] == '/' {
		tgt = tgt[0 : len(tgt)-1]
	} else {
		return errors.New("cannot zip file")
	}
	f, err := kio.Open(path, "w")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}

	rawlist0 := make([]string, 0)
	rawlist1 := make([]string, 0)
	root := tgt[0 : strings.LastIndex(tgt, "/")+1]
	g3fdinfo(root, tgt[strings.LastIndex(tgt, "/")+1:], &rawlist0, &rawlist1)
	for i, r := range rawlist0 {
		rawlist0[i] = "temp261/" + r
	}
	rawlist2 := make([]int, len(rawlist1))
	rawlist3 := make([]string, len(rawlist1))
	for i, r := range rawlist1 {
		rawlist2[i] = kio.Size(root + r)
		rawlist3[i] = "temp261/" + r
	}

	sub0, sub1 := g3gensubhead(append([]string{"temp261"}, rawlist0...), rawlist3, rawlist2, tbox.Winsign)
	sub2 := bytes.Join(sub0, nil)
	tbox.Header = append(append([]byte("KTS2"), make([]byte, 2)...), kobj.Encode(len(sub0), 3)...)
	tbox.Header = append(tbox.Header, 1, 8, 247, 0, 0)
	temp := hash32(sub2)
	for i := 3; i >= 0; i-- {
		tbox.Header = append(tbox.Header, temp[i])
	}

	kio.Write(f, tbox.Prehead)
	kio.Write(f, tbox.Header)
	kio.Write(f, sub0[0])
	kio.Write(f, sub1)
	for i, r := range rawlist1 {
		kio.Write(f, sub0[i+1])
		g3fadd(root+r, f)
	}
	return nil
}

// gen3 kzip file view, check CRC32
func (tbox *G3kzip) View(tgt string) error {
	f, err := kio.Open(tgt, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}

	var temp []byte
	tbox.Prehead = nil
	temp, err = kio.Read(f, 4)
	if err != nil {
		return err
	}
	for !kio.Bequal(temp, []byte("KTS2")) {
		tbox.Prehead = append(tbox.Prehead, temp...)
		temp, err = kio.Read(f, 1020)
		if err != nil {
			return err
		}
		if len(temp) == 1020 {
			tbox.Prehead = append(tbox.Prehead, temp...)
		} else {
			return errors.New("invalid gen3 kzip file")
		}
		temp, err = kio.Read(f, 4)
		if err != nil {
			return err
		}
	}

	tbox.Header = temp
	temp, err = kio.Read(f, 14)
	if err == nil {
		tbox.Header = append(tbox.Header, temp...)
	} else {
		return err
	}

	num0 := kobj.Decode(tbox.Header[6:9])
	num1 := int(tbox.Header[9])
	num2 := int(tbox.Header[10])
	num3 := kobj.Decode(tbox.Header[11:14])
	current := len(tbox.Prehead) + 18
	tbox.Chunkpos = make([]int, num0)
	tbox.Subhead = make([][]byte, num0)
	for i := 0; i < num0; i++ {
		temp, err = kio.Read(f, num1+num2+num3)
		if err != nil {
			return err
		}
		tbox.Chunkpos[i] = current
		tbox.Subhead[i] = temp
		num4 := kobj.Decode(temp[num1 : num1+num2])
		current = current + num1 + num2 + num3 + num4
		f.Seek(int64(num4), 1)
	}

	temp = hash32(bytes.Join(tbox.Subhead, nil))
	for i := 3; i >= 0; i-- {
		if temp[i] != tbox.Header[17-i] {
			return errors.New("wrong CRC32 value")
		}
	}
	return nil
}

// gen3 kzip unpack (tgt : kzip file, output : ./temp261)
func (tbox *G3kzip) Unpack(tgt string) error {
	if tbox.Header == nil {
		return errors.New("should done View() first")
	}
	os.RemoveAll("./temp261/") // !!! autoclear temp261/ !!!
	f, err := kio.Open(tgt, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}

	num0 := int(tbox.Header[9])
	num1 := int(tbox.Header[10])
	num2 := kobj.Decode(tbox.Header[11:14])
	for i, r := range tbox.Chunkpos {
		f.Seek(int64(r+num0+num1+num2), 0)
		if string(tbox.Subhead[i][0:num0]) == "S" {
			temp, err := kio.Read(f, kobj.Decode(tbox.Subhead[i][num0:num0+num1]))
			if err != nil {
				return err
			}
			err = g3mkdir(string(temp))
			if err != nil {
				return err
			}
		} else {
			name := strings.TrimRight(string(tbox.Subhead[i][num0+num1:num0+num1+num2]), " ")
			size := kobj.Decode(tbox.Subhead[i][num0 : num0+num1])
			g3fsub(f, "./"+strings.Replace(name, "\\", "/", -1), size)
		}
	}
	return nil
}
