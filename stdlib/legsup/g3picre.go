// test661 : stdlib5.legsup gen3picre

package legsup

import (
	"archive/zip"
	"errors"
	"os"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strings"
)

// move files to ./temp365/
func g3fcopy(path string) error {
	if kio.Size(path) > 2147483648 {
		return errors.New("file too big (over 2GiB)")
	}
	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	t, err := kio.Open("./temp365"+path[strings.LastIndex(path, "/"):], "w")
	if err == nil {
		defer t.Close()
	} else {
		return err
	}
	temp, err := kio.Read(f, -1)
	if err != nil {
		return err
	}
	_, err = kio.Write(t, temp)
	return err
}

// read file at temp365, returns nil if error
func g3fread(name string) []byte {
	f, err := kio.Open(name, "r")
	if err == nil {
		defer f.Close()
	} else {
		return nil
	}
	temp, err := kio.Read(f, -1)
	if err == nil {
		return temp
	} else {
		return nil
	}
}

// gen3 picre making zip file, pack files to zip
func g3mkzip(files []string, path string) ([]byte, error) {
	f, err := kio.Open(path+".zip", "w")
	if err != nil {
		return nil, err
	}
	z := zip.NewWriter(f)
	os.RemoveAll("./temp365/") // autoclear ./temp365/
	os.Mkdir("./temp365/", os.ModePerm)
	defer os.RemoveAll("./temp365/")

	names := make([]string, len(files))
	for i, r := range files {
		r = kio.Abs(r)
		if r[len(r)-1] == '/' {
			return nil, errors.New("cannot zip folder")
		}
		names[i] = r[strings.LastIndex(r, "/")+1:]
		err = g3fcopy(r)
		if err != nil {
			return nil, err
		}
	}

	for _, r := range names {
		w, err := z.Create("temp365/" + r)
		if err == nil {
			w.Write(g3fread("./temp365/" + r))
		} else {
			return nil, err
		}
	}

	z.Close()
	f.Close()
	res := g3fread(path + ".zip")
	os.Remove(path + ".zip")
	return res, nil
}

// gen3 picre, zip pic + files
func G3picre(pic []byte, files []string, path string) error {
	zipv, err := g3mkzip(files, path)
	if err != nil {
		return err
	}
	f, err := kio.Open(path, "w")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}

	size := len(zipv)
	addlen := len(pic)
	cheadloc := kobj.Decode(zipv[size-6 : size-2])
	fnum := kobj.Decode(zipv[size-12 : size-10])
	fdata := zipv[0:cheadloc]
	endhead := append(append(zipv[size-22:size-6], kobj.Encode(addlen+cheadloc, 4)...), zipv[size-2:]...)

	centerhead := make([]byte, 0)
	curpos := cheadloc
	for i := 0; i < fnum; i++ {
		seta := zipv[curpos : curpos+28]
		namelen := kobj.Decode(zipv[curpos+28 : curpos+30])
		setb := zipv[curpos+30 : curpos+42]
		stpoint := kobj.Encode(kobj.Decode(zipv[curpos+42:curpos+46])+addlen, 4)
		fname := zipv[curpos+46 : curpos+46+namelen]
		centerhead = append(append(centerhead, seta...), kobj.Encode(namelen, 2)...)
		centerhead = append(append(centerhead, setb...), stpoint...)
		centerhead = append(centerhead, fname...)
		curpos = curpos + 46 + namelen
	}

	kio.Write(f, pic)
	kio.Write(f, fdata)
	kio.Write(f, centerhead)
	kio.Write(f, endhead)
	return nil
}
