// test656 : stdlib5.legsup gen2enc

package legsup

import (
	"crypto/rand"
	"errors"
	"fmt"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strings"
)

// gen2 enc
type G2enc struct {
	Path     string // target file path
	Pw       string // encryption password
	Hint     string // pw hint (max 600B)
	Hidename bool   // name encryption mode

	salt   []byte
	pwhash []byte
	namedt []byte
	nmsize []byte
}

// gen2 enc init
func (tbox *G2enc) Init() {
	tbox.Path = ""
	tbox.Pw = ""
	tbox.Hint = ""
	tbox.Hidename = true
	tbox.salt = nil
	tbox.pwhash = nil
	tbox.namedt = nil
	tbox.nmsize = nil
}

// gen2 enc file encrypt
func (tbox *G2enc) Encrypt() (string, error) {
	hint := []byte(tbox.Hint)
	if len(hint) > 600 {
		return "", errors.New("hint overflow (600B)")
	} else {
		hpad := make([]byte, 600-len(hint))
		for i := range hpad {
			hpad[i] = 32
		}
		hint = append(hint, hpad...)
	}
	tbox.salt = make([]byte, 80)
	rand.Read(tbox.salt)
	for i, r := range tbox.salt {
		tbox.salt[i] = r%94 + 32
	}

	temp := make([]byte, 0)
	temp = append(temp, tbox.salt[0:40]...)
	temp = append(temp, []byte(tbox.Pw)...)
	temp = append(temp, tbox.salt[40:80]...)
	tbox.pwhash = hash512(temp)
	temp = make([]byte, 0)
	temp = append(temp, tbox.salt[0:20]...)
	temp = append(temp, []byte(tbox.Pw)...)
	temp = append(temp, tbox.salt[20:80]...)
	ckey := hash256(temp)
	civ := hash128(ckey)
	temp = make([]byte, 0)
	temp = append(temp, tbox.salt[0:60]...)
	temp = append(temp, []byte(tbox.Pw)...)
	temp = append(temp, tbox.salt[60:80]...)
	tkey := hash256(temp)
	tiv := hash128(tkey)

	path := kio.Abs(tbox.Path)
	oriname := []byte(path[strings.LastIndex(path, "/")+1:])
	tbox.nmsize = make([]byte, 2)
	var nmmode []byte
	if tbox.Hidename {
		tbox.nmsize[1] = byte(len(oriname))
		oriname = append(oriname, make([]byte, 256-len(oriname))...)
		nmmode = []byte("hi")
		tbox.namedt = aes16en(oriname, tkey, tiv)
	} else {
		oriname = make([]byte, 256)
		nmmode = []byte("op")
		tbox.namedt = oriname
	}

	header := append([]byte("kos2"), tbox.salt...)
	header = append(header, tbox.pwhash...)
	header = append(header, hint...)
	header = append(header, tbox.namedt...)
	header = append(header, tbox.nmsize...)
	header = append(header, nmmode...)
	header = append(header, hash128(header)...)

	var newname string
	if tbox.Hidename {
		temp = make([]byte, 3)
		rand.Read(temp)
		newname = fmt.Sprintf("%s/%d.k", path[0:strings.LastIndex(path, "/")], kobj.Decode(temp)%9000+1000)
	} else {
		newname = path + ".k"
	}
	size := kio.Size(path)
	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return "", err
	}
	t, err := kio.Open(newname, "w")
	if err == nil {
		defer t.Close()
	} else {
		return "", err
	}
	_, err = kio.Write(t, header)
	if err != nil {
		return "", err
	}

	for i := 0; i < size/16384; i++ {
		temp, err = kio.Read(f, 16384)
		if err != nil {
			return "", err
		}
		_, err = kio.Write(t, aes16en(temp, ckey, civ))
		if err != nil {
			return "", err
		}
	}
	if size%16384 == 0 {
		temp = make([]byte, 0)
	} else {
		temp, err = kio.Read(f, size%16384)
		if err != nil {
			return "", err
		}
	}
	_, err = kio.Write(t, aes1en(temp, ckey, civ))
	if err != nil {
		return "", err
	}
	return newname, nil
}

// gen2 enc encfile view
func (tbox *G2enc) View() error {
	if tbox.Path[len(tbox.Path)-2:] != ".k" {
		return errors.New("invalid extension")
	}
	f, err := kio.Open(tbox.Path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	temp, err := kio.Read(f, 1024)
	if err != nil {
		return err
	}
	if !kio.Bequal(temp[0:4], []byte("kos2")) {
		return errors.New("InvalidFile")
	} else if !kio.Bequal(temp[1008:1024], hash128(temp[0:1008])) {
		return errors.New("DamagedHeader")
	} else {
		tbox.salt = temp[4:84]
		tbox.pwhash = temp[84:148]
		tbox.Hint = string(temp[148:748])
		tbox.namedt = temp[748:1004]
		tbox.nmsize = temp[1004:1006]
		if kio.Bequal(temp[1006:1008], []byte("op")) {
			tbox.Hidename = false
		} else if kio.Bequal(temp[1006:1008], []byte("hi")) {
			tbox.Hidename = true
		} else {
			return errors.New("DamagedHeader")
		}
	}
	return nil
}

// gen2 enc file decrypt
func (tbox *G2enc) Decrypt() error {
	temp := make([]byte, 0)
	temp = append(temp, tbox.salt[0:40]...)
	temp = append(temp, []byte(tbox.Pw)...)
	temp = append(temp, tbox.salt[40:80]...)
	temp = hash512(temp)
	if !kio.Bequal(tbox.pwhash, temp) {
		return errors.New("InvalidPW")
	}
	temp = make([]byte, 0)
	temp = append(temp, tbox.salt[0:20]...)
	temp = append(temp, []byte(tbox.Pw)...)
	temp = append(temp, tbox.salt[20:80]...)
	ckey := hash256(temp)
	civ := hash128(ckey)
	temp = make([]byte, 0)
	temp = append(temp, tbox.salt[0:60]...)
	temp = append(temp, []byte(tbox.Pw)...)
	temp = append(temp, tbox.salt[60:80]...)
	tkey := hash256(temp)
	tiv := hash128(tkey)

	path := kio.Abs(tbox.Path)
	size := kio.Size(path) - 1024
	var newname string
	if tbox.Hidename {
		newname = string(aes16de(tbox.namedt, tkey, tiv)[0:tbox.nmsize[1]])
	} else {
		newname = path[0 : len(path)-2]
	}
	f, err := kio.Open(path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	t, err := kio.Open(newname, "w")
	if err == nil {
		defer t.Close()
	} else {
		return err
	}
	f.Seek(1024, 0)

	num0 := size / 16384
	num1 := size % 16384
	if num1 == 0 {
		num0 = num0 - 1
		num1 = 16384
	}
	for i := 0; i < num0; i++ {
		temp, err = kio.Read(f, 16384)
		if err != nil {
			return err
		}
		_, err = kio.Write(t, aes16de(temp, ckey, civ))
		if err != nil {
			return err
		}
	}
	temp, err = kio.Read(f, num1)
	if err != nil {
		return err
	}
	temp = aes1de(temp, ckey, civ)
	if len(temp)%16 == 0 {
		if kio.Bequal(temp[len(temp)-16:], make([]byte, 16)) {
			temp = temp[0 : len(temp)-16]
		}
	}
	_, err = kio.Write(t, temp)
	if err != nil {
		return err
	}
	return nil
}
