// test655 : stdlib5.legsup gen1enc

package legsup

import (
	"crypto/rand"
	"errors"
	"stdlib5/kio"
)

// get 1024B content key
func g1genkey(data []byte, salt []byte) []byte {
	temp := kio.Bprint(data)
	out := make([]byte, 1024)
	for i := 0; i < 32; i++ {
		tb := []byte(temp[2*i : 2*i+2])
		tb = append(tb, salt...)
		for j, r := range hash256(tb) {
			out[32*i+j] = r
		}
	}
	return out
}

// add/sub bytes, mode : T-enc/F-dec
func g1calc(target []byte, key []byte, mode bool) {
	if mode {
		for i, r := range target {
			target[i] = r + key[i]
		}
	} else {
		for i, r := range target {
			target[i] = r - key[i]
		}
	}
}

// gen1 enc
type G1enc struct {
	Path string // target file path
	Pw   string // encryption password
	Hint string // pw hint (max 324B)

	salt   []byte
	pwhash []byte
}

// gen1 enc init
func (tbox *G1enc) Init() {
	tbox.Path = ""
	tbox.Pw = ""
	tbox.Hint = ""
	tbox.salt = nil
	tbox.pwhash = nil
}

// gen1 enc file encrypt
func (tbox *G1enc) Encrypt() error {
	hint := []byte(tbox.Hint)
	if len(hint) > 324 {
		return errors.New("hint overflow (324B)")
	} else {
		hpad := make([]byte, 324-len(hint))
		for i := range hpad {
			hpad[i] = 32
		}
		hint = append(hint, hpad...)
	}
	tbox.salt = make([]byte, 40)
	rand.Read(tbox.salt)
	for i, r := range tbox.salt {
		tbox.salt[i] = r%94 + 32
	}
	tbox.pwhash = hash256(append(append(make([]byte, 0), tbox.salt...), []byte(tbox.Pw)...))
	key := g1genkey(hash256(append(append(make([]byte, 0), []byte(tbox.Pw)...), tbox.salt...)), tbox.salt)
	size := kio.Size(tbox.Path)

	f, err := kio.Open(tbox.Path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	t, err := kio.Open(tbox.Path+".k", "w")
	if err == nil {
		defer t.Close()
	} else {
		return err
	}
	kio.Write(t, []byte(".kos"))
	kio.Write(t, tbox.salt)
	kio.Write(t, tbox.pwhash)
	kio.Write(t, hint)

	for i := 0; i < size/1024; i++ {
		temp, _ := kio.Read(f, 1024)
		g1calc(temp, key, true)
		kio.Write(t, temp)
	}
	if size%1024 != 0 {
		temp, _ := kio.Read(f, size%1024)
		temp = append(temp, make([]byte, 1024-size%1024)...)
		g1calc(temp, key, true)
		kio.Write(t, temp[0:size%1024])
	}
	return nil
}

// gen1 enc encfile view
func (tbox *G1enc) View() error {
	if tbox.Path[len(tbox.Path)-2:] != ".k" {
		return errors.New("invalid extension")
	}
	f, err := kio.Open(tbox.Path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	temp, err := kio.Read(f, 400)
	if err != nil {
		return err
	}
	if !kio.Bequal(temp[0:4], []byte(".kos")) {
		return errors.New("InvalidFile")
	} else {
		tbox.salt = temp[4:44]
		tbox.pwhash = temp[44:76]
		tbox.Hint = string(temp[76:])
		return nil
	}
}

// gen1 enc file decrypt
func (tbox *G1enc) Decrypt() error {
	if !kio.Bequal(tbox.pwhash, hash256(append(append(make([]byte, 0), tbox.salt...), []byte(tbox.Pw)...))) {
		return errors.New("InvalidPW")
	}
	key := g1genkey(hash256(append(append(make([]byte, 0), []byte(tbox.Pw)...), tbox.salt...)), tbox.salt)
	size := kio.Size(tbox.Path) - 400

	f, err := kio.Open(tbox.Path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	t, err := kio.Open(tbox.Path[0:len(tbox.Path)-2], "w")
	if err == nil {
		defer t.Close()
	} else {
		return err
	}
	f.Seek(400, 0)

	for i := 0; i < size/1024; i++ {
		temp, _ := kio.Read(f, 1024)
		g1calc(temp, key, false)
		kio.Write(t, temp)
	}
	if size%1024 != 0 {
		temp, _ := kio.Read(f, size%1024)
		temp = append(temp, make([]byte, 1024-size%1024)...)
		g1calc(temp, key, false)
		kio.Write(t, temp[0:size%1024])
	}
	return nil
}
