// test635 : stdlib5.kaes st

package kaes

// go get "golang.org/x/crypto/sha3"
// go get "golang.org/x/crypto/scrypt"

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"os"
	"stdlib5/kdb"
	"stdlib5/kio"
	"stdlib5/ksc"
	"stdlib5/ksign"
	"stdlib5/picdt"
	"strings"

	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
)

// integrated byte/file io
type SimIO struct {
	Buf  []byte   // Bmode buffer
	File *os.File // Fmode file
	Size int      // raw size of IOtgt, for readonly

	rpos  int // reading position, for Bmode readonly
	dpos  int // encdata stpoint
	dsize int // encdata size
}

// open/init struct, v []byte : Bmode / str : Fmode, isreader T:readonly / F:writeonly
func (tbox *SimIO) Open(v interface{}, isreader bool) error {
	tbox.rpos = 0
	tbox.dpos = 0
	tbox.dsize = 0
	tbox.Size = 0
	switch data := v.(type) {
	case []byte:
		tbox.Buf = data
		tbox.File = nil
		if isreader {
			tbox.Size = len(data)
		}
	case string:
		tbox.Buf = nil
		var err error
		if isreader {
			tbox.File, err = kio.Open(data, "r")
			tbox.Size = kio.Size(data)
		} else {
			tbox.File, err = kio.Open(data, "w")
		}
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid datatype")
	}
	return nil
}

// close internal data, returns nil if Fmode / []byte if Bmode
func (tbox *SimIO) Close() []byte {
	var temp []byte
	if tbox.File == nil {
		temp = tbox.Buf
	} else {
		tbox.File.Close()
		temp = nil
	}
	tbox.Buf = nil
	tbox.File = nil
	return temp
}

// seek from start of IOtgt
func (tbox *SimIO) Seek(pos int) {
	if tbox.File == nil {
		if pos > tbox.Size {
			tbox.rpos = tbox.Size
		} else {
			tbox.rpos = pos
		}
	} else {
		tbox.File.Seek(int64(pos), 0)
		tbox.rpos = pos
	}
}

// read by size, readbyte size can be smaller than size, !!! cannot read 1GiB+ !!!
func (tbox *SimIO) Read(size int) []byte {
	var temp []byte
	var ti int
	if tbox.File == nil {
		ti = tbox.rpos + size
		if ti > tbox.Size {
			temp = tbox.Buf[tbox.rpos:tbox.Size]
			tbox.rpos = tbox.Size
		} else {
			temp = tbox.Buf[tbox.rpos:ti]
			tbox.rpos = ti
		}
	} else {
		temp = make([]byte, size)
		ti, _ = tbox.File.Read(temp)
		if ti == size {
			tbox.rpos = tbox.rpos + size
		} else {
			tbox.rpos = tbox.rpos + ti
			temp = temp[0:ti]
		}
	}
	return temp
}

// write []byte, !!! cannot write 1GiB+ !!!
func (tbox *SimIO) Write(data []byte) {
	if tbox.File == nil {
		tbox.Buf = append(tbox.Buf, data...)
	} else {
		tbox.File.Write(data)
	}
}

// hash SHA3-512
func hash3512(data []byte) []byte {
	h := sha3.New512()
	h.Write(data)
	return h.Sum(nil)
}

// simple aes-256 calc, key 32B / iv 16B / isenc T:encrypt F:decrypt / ispad T:dopad F:nopad
func aescalc(data []byte, key []byte, iv []byte, isenc bool, ispad bool) []byte {
	if len(key) != 32 || len(iv) != 16 {
		return nil
	}
	var out []byte
	block, _ := aes.NewCipher(key)
	if isenc {
		encrypter := cipher.NewCBCEncrypter(block, iv)
		if ispad {
			plen := 16 - (len(data) % 16)
			for i := 0; i < plen; i++ {
				data = append(data, byte(plen))
			}
		}
		out = make([]byte, len(data))
		encrypter.CryptBlocks(out, data)
	} else {
		decrypter := cipher.NewCBCDecrypter(block, iv)
		out = make([]byte, len(data))
		decrypter.CryptBlocks(out, data)
		if ispad {
			plen := int(out[len(out)-1])
			out = out[0 : len(out)-plen]
		}
	}
	return out
}

// generate pwhash 128B, mkey 96B
func genpm(pw []byte, kf []byte, salt []byte) ([]byte, []byte) {
	tb := append(append(append(append(append(make([]byte, 0), pw...), pw...), kf...), pw...), kf...)
	pwh, _ := scrypt.Key(tb, salt, 524288, 8, 1, 128)
	tb = append(append(append(append(append(make([]byte, 0), kf...), pw...), kf...), kf...), pw...)
	mkey, _ := scrypt.Key(tb, salt, 16384, 8, 1, 96)
	return pwh, mkey
}

// generate size when n bytes is pad encrypted
func gensize(n int) int {
	return n + 16 - n%16
}

// generate random word with length 2*size
func genpath(size int) string {
	return kio.Bprint(Genrand(size))
}

// analysis/set SimIO reader by ksc5/kaes5 format, returns enchead/signhead
func ksc5read(reader *SimIO) ([]byte, []byte, error) {
	for reader.rpos+8 < reader.Size {
		if kio.Bequal(reader.Read(4), []byte("KSC5")) {
			break
		} else {
			reader.Seek(reader.rpos + 508)
		}
	}
	if reader.Size < reader.rpos+16 {
		return nil, nil, errors.New("InvalidKSC5")
	}
	if !kio.Bequal(reader.Read(4), []byte("KAES")) {
		return nil, nil, errors.New("InvalidKAES5")
	}
	reserved := reader.Read(8)
	encheader := reader.Read(ksc.Decode(reader.Read(8)))
	signheader := reader.Read(ksc.Decode(reader.Read(8)))
	if !kio.Bequal(reserved, append(ksc.Crc32hash(encheader), ksc.Crc32hash(signheader)...)) {
		return nil, nil, errors.New("InvalidCRC32")
	}
	reader.dsize = ksc.Decode(reader.Read(8))
	reader.dpos = reader.rpos
	return encheader, signheader, nil
}

// part-read, (x40) 20m + 512k + nB
func readflow(reader *SimIO, exstrm []chan []byte, isenc bool, proc *float64, fdata *[]byte, fnum *int) {
	defer func() {
		for _, r := range exstrm {
			close(r)
		}
		recover()
	}()
	*proc = 0.0
	reader.Seek(reader.dpos)
	num0 := reader.dsize / 524288 // 512k N
	num1 := reader.dsize % 524288 // left N
	if !isenc && num1 == 0 {
		num0 = num0 - 1
		num1 = 524288
	}
	for i := 0; i < num0; i++ {
		*proc = float64(i) / float64(num0)
		exstrm[i%40] <- reader.Read(524288)
	}
	*fdata = reader.Read(num1)
	*fnum = num0 % 40
}

// part-calc, (x40), ckey 1920B 40*(iv16 + key32), updates ckey
func calcflow(instrm []chan []byte, exstrm []chan []byte, isenc bool, ckey []byte, pnum int) {
	defer func() {
		close(exstrm[pnum])
		recover()
	}()
	ti := 48 * pnum
	var temp []byte
	block, _ := aes.NewCipher(ckey[ti+16 : ti+48])

	if isenc {
		encrypter := cipher.NewCBCEncrypter(block, ckey[ti:ti+16])
		for r := range instrm[pnum] {
			temp = make([]byte, 524288)
			encrypter.CryptBlocks(temp, r)
			copy(ckey[ti:ti+16], temp[524272:524288])
			exstrm[pnum] <- temp
		}

	} else {
		decrypter := cipher.NewCBCDecrypter(block, ckey[ti:ti+16])
		for r := range instrm[pnum] {
			temp = make([]byte, 524288)
			copy(ckey[ti:ti+16], r[524272:524288])
			decrypter.CryptBlocks(temp, r)
			exstrm[pnum] <- temp
		}
	}
}

// part-write, (x40), wr header / do final calc
func writeflow(writer *SimIO, instrm []chan []byte, isenc bool, header []byte, ckey []byte, fdata *[]byte, fnum *int) {
	defer recover()
	writer.Write(header)
	for r := range instrm[39] {
		for i := 0; i < 39; i++ {
			writer.Write(<-instrm[i])
		}
		writer.Write(r)
	}
	for i := 0; i < 40; i++ {
		temp, ext := <-instrm[i]
		if ext {
			writer.Write(temp)
		}
	}
	ti := (*fnum) * 48
	writer.Write(aescalc(*fdata, ckey[ti+16:ti+48], ckey[ti:ti+16], isenc, true))
}

// gen5 kaes all-mode
type Allmode struct {
	Hint    string    // pwkf hint str
	Msg     string    // program msg str
	Signkey [2]string // (pub, pri) ksign rsa key
	Proc    float64   // -1 : not started, 0~1 : working, 2 : end

	salt     []byte
	pwhash   []byte
	ckeydata []byte
	tkeydata []byte
	encname  []byte

	before SimIO
	after  SimIO
}

// make encheader / must set (Hint, Msg), ckey before call
func (tbox *Allmode) mkhead(pw []byte, kf []byte, ckcpy []byte, oldpath string) ([]byte, error) {
	tkey := Genrand(48)
	tbox.salt = Genrand(40)
	var mkey []byte
	tbox.pwhash, mkey = genpm(pw, kf, tbox.salt)
	tbox.ckeydata = aescalc(ckcpy, mkey[16:48], mkey[0:16], true, false)
	tbox.tkeydata = aescalc(tkey, mkey[64:96], mkey[48:64], true, false)
	tbox.encname = aescalc([]byte(oldpath), tkey[16:48], tkey[0:16], true, true)

	worker := kdb.Initkdb()
	worker.Read("salt = 0\npwhash = 0\nckeydata = 0\ntkeydata = 0\nencname = 0\nhint = 0\nmsg = 0")
	worker.Fix("salt", tbox.salt)
	worker.Fix("pwhash", tbox.pwhash)
	worker.Fix("ckeydata", tbox.ckeydata)
	worker.Fix("tkeydata", tbox.tkeydata)
	worker.Fix("encname", tbox.encname)
	worker.Fix("hint", tbox.Hint)
	worker.Fix("msg", tbox.Msg)
	ts, err := worker.Write()
	if err == nil {
		return []byte(ts), nil // c0 : encheader
	} else {
		return nil, err
	}
}

// general internal encryption / must set (Signkey), SimIO before call
func (tbox *Allmode) encrypt(pw []byte, kf []byte, oldpath string, pmode int) error {
	ckey := Genrand(1920)
	ckcpy := append(make([]byte, 0), ckey...) // for writing header
	var tgtb []byte                           // for final calc
	var tgti int                              // for final calc
	flow0 := make([]chan []byte, 40)          // read -> calc
	flow1 := make([]chan []byte, 40)          // calc -> write
	for i := 0; i < 40; i++ {
		flow0[i] = make(chan []byte, 12)
		flow1[i] = make(chan []byte, 12)
	}
	tbox.before.dpos = 0                 // reader tuning
	tbox.before.dsize = tbox.before.Size // reader tuning
	go readflow(&tbox.before, flow0, true, &tbox.Proc, &tgtb, &tgti)
	for i := 0; i < 40; i++ {
		go calcflow(flow0, flow1, true, ckey, i)
	}

	encheader, err := tbox.mkhead(pw, kf, ckcpy, oldpath) // c0 : encheader
	if err != nil {
		return err
	}
	var signheader []byte
	if len(tbox.Signkey[0]) != 0 && len(tbox.Signkey[1]) != 0 {
		w1 := kdb.Initkdb()
		w1.Read("publickey = 0\nsigndata = 0")
		w1.Fix("publickey", tbox.Signkey[0])
		tb, err := ksign.Sign(tbox.Signkey[1], hash3512(encheader))
		if err != nil {
			return err
		}
		w1.Fix("signdata", tb)
		ts, err := w1.Write()
		if err != nil {
			return err
		}
		signheader = []byte(ts)
	} else {
		signheader = make([]byte, 0) // c1 : signheader
	}

	worker := ksc.Initksc()
	switch pmode {
	case 0:
		tb := picdt.Ka5webp()
		tb = append(tb, make([]byte, (102400-len(tb))%512)...)
		worker.Prehead = tb
	case 1:
		tb := picdt.Ka5png()
		tb = append(tb, make([]byte, (102400-len(tb))%512)...)
		worker.Prehead = tb
	default:
		worker.Prehead = make([]byte, 0)
	}
	worker.Subtype = []byte("KAES")
	worker.Reserved = append(ksc.Crc32hash(encheader), ksc.Crc32hash(signheader)...)
	header, err := worker.Writeb() // real header to write
	if err != nil {
		return err
	}
	header = worker.Linkb(header, encheader)
	header = worker.Linkb(header, signheader)
	header = append(header, ksc.Encode(gensize(tbox.before.Size), 8)...)

	writeflow(&tbox.after, flow1, true, header, ckey, &tgtb, &tgti)
	tbox.after.Write([]byte{255, 255, 255, 255, 255, 255, 255, 255})
	return nil
}

// read encheader / update internal value, check signheader
func (tbox *Allmode) rdhead(encheader []byte, signheader []byte) error {
	tbox.Signkey = [2]string{"", ""}
	if len(signheader) != 0 {
		w0 := kdb.Initkdb()
		w0.Read(string(signheader))
		tv, _ := w0.Get("publickey")
		tbox.Signkey[0] = tv.Dat6
		tbox.Signkey[1] = ""
		tv, _ = w0.Get("signdata")
		vfy, err := ksign.Verify(tbox.Signkey[0], tv.Dat5, hash3512(encheader))
		if err != nil {
			return err
		}
		if !vfy {
			return errors.New("invalidRSAsign")
		}
	}

	w1 := kdb.Initkdb()
	err := w1.Read(string(encheader))
	if err != nil {
		return err
	}
	tv, _ := w1.Get("salt")
	tbox.salt = tv.Dat5
	tv, _ = w1.Get("pwhash")
	tbox.pwhash = tv.Dat5
	tv, _ = w1.Get("ckeydata")
	tbox.ckeydata = tv.Dat5
	tv, _ = w1.Get("tkeydata")
	tbox.tkeydata = tv.Dat5
	tv, _ = w1.Get("encname")
	tbox.encname = tv.Dat5
	tv, _ = w1.Get("hint")
	tbox.Hint = tv.Dat6
	tv, _ = w1.Get("msg")
	tbox.Msg = tv.Dat6
	return nil
}

// partial decryption, returns ckey, decname
func (tbox *Allmode) dechead(pw []byte, kf []byte) ([]byte, string, error) {
	pwhcmp, mkey := genpm(pw, kf, tbox.salt)
	if !kio.Bequal(pwhcmp, tbox.pwhash) {
		return nil, "", errors.New("invalidPWKF")
	}
	ckey := aescalc(tbox.ckeydata, mkey[16:48], mkey[0:16], false, false)
	tkey := aescalc(tbox.tkeydata, mkey[64:96], mkey[48:64], false, false)
	name := string(aescalc(tbox.encname, tkey[16:48], tkey[0:16], false, true))
	return ckey, name, nil
}

// general internal decryption / must set SimIO before call
func (tbox *Allmode) decrypt(ckey []byte) {
	var tgtb []byte                  // for final calc
	var tgti int                     // for final calc
	flow0 := make([]chan []byte, 40) // read -> calc
	flow1 := make([]chan []byte, 40) // calc -> write
	for i := 0; i < 40; i++ {
		flow0[i] = make(chan []byte, 12)
		flow1[i] = make(chan []byte, 12)
	}
	go readflow(&tbox.before, flow0, false, &tbox.Proc, &tgtb, &tgti)
	for i := 0; i < 40; i++ {
		go calcflow(flow0, flow1, false, ckey, i)
	}
	writeflow(&tbox.after, flow1, false, nil, ckey, &tgtb, &tgti)
}

// binary encryption, pwkf bytes / data []byte / pmode 0:webp 1:png 2:none, returns encdata
func (tbox *Allmode) EnBin(pw []byte, kf []byte, data []byte, pmode int) ([]byte, error) {
	defer func() { tbox.Proc = 2.0 }()
	tbox.Proc = -1.0
	tbox.before.Open(data, true)
	tbox.after.Open(make([]byte, 0, len(data)+1048576), false)
	err := tbox.encrypt(pw, kf, "NewData.bin", pmode)
	if err != nil {
		return nil, err
	}
	temp := tbox.after.Close()
	return temp, nil
}

// file encryption, pwkf bytes / path string / pmode 0:webp 1:png 2:none, returns encpath
func (tbox *Allmode) EnFile(pw []byte, kf []byte, path string, pmode int) (string, error) {
	defer func() { tbox.Proc = 2.0 }()
	tbox.Proc = -1.0
	path = kio.Abs(path)
	oldpath := path[strings.LastIndex(path, "/")+1:]
	newpath := path[0:strings.LastIndex(path, "/")+1] + genpath(2)
	switch pmode {
	case 0:
		newpath = newpath + ".webp"
	case 1:
		newpath = newpath + ".png"
	default:
		newpath = newpath + ".k"
	}
	tbox.before.Open(path, true)
	defer tbox.before.Close()
	tbox.after.Open(newpath, false)
	defer tbox.after.Close()
	err := tbox.encrypt(pw, kf, oldpath, pmode)
	if err != nil {
		return "", err
	}
	return newpath, nil
}

// view encbin
func (tbox *Allmode) ViewBin(data []byte) error {
	tbox.Proc = -1.0
	tbox.before.Open(data, true)
	encheader, signheader, err := ksc5read(&tbox.before)
	if err != nil {
		return err
	}
	return tbox.rdhead(encheader, signheader)
}

// view encfile
func (tbox *Allmode) ViewFile(path string) error {
	tbox.Proc = -1.0
	tbox.before.Open(path, true)
	defer tbox.before.Close()
	encheader, signheader, err := ksc5read(&tbox.before)
	if err != nil {
		return err
	}
	return tbox.rdhead(encheader, signheader)
}

// binary decryption, pwkf bytes, returns decbin
func (tbox *Allmode) DeBin(pw []byte, kf []byte, data []byte) ([]byte, error) {
	defer func() { tbox.Proc = 2.0 }()
	tbox.Proc = -1.0
	if tbox.before.dpos == 0 {
		return nil, errors.New("should done ViewBin() first")
	}
	ckey, _, err := tbox.dechead(pw, kf)
	if err != nil {
		return nil, err
	}

	ti0 := tbox.before.dpos
	ti1 := tbox.before.dsize
	tbox.before.Open(data, true)
	tbox.before.dpos = ti0
	tbox.before.dsize = ti1

	tbox.after.Open(make([]byte, 0, len(data)+1048576), false)
	tbox.decrypt(ckey)
	temp := tbox.after.Close()
	return temp, nil
}

// file decryption, pwkf bytes / path string, returns decpath
func (tbox *Allmode) DeFile(pw []byte, kf []byte, path string) (string, error) {
	defer func() { tbox.Proc = 2.0 }()
	tbox.Proc = -1.0
	if tbox.before.dpos == 0 {
		return "", errors.New("should done ViewFile() first")
	}
	ckey, name, err := tbox.dechead(pw, kf)
	if err != nil {
		return "", err
	}

	ti0 := tbox.before.dpos
	ti1 := tbox.before.dsize
	tbox.before.Open(path, true)
	defer tbox.before.Close()
	tbox.before.dpos = ti0
	tbox.before.dsize = ti1

	path = kio.Abs(path)
	newpath := path[0:strings.LastIndex(path, "/")+1] + name
	tbox.after.Open(newpath, false)
	defer tbox.after.Close()
	tbox.decrypt(ckey)
	return newpath, nil
}

// gen5 kaes func-mode
type Funcmode struct {
	Before SimIO   // reader, !! caller should Close() it !!
	After  SimIO   // writer, !! caller should Close() it !!
	Proc   float64 // -1 : not started, 0~1 : working, 2 : end
}

// encrypt with akey 48B
func (tbox *Funcmode) Encrypt(akey []byte) error {
	defer func() { tbox.Proc = 2.0 }()
	tbox.Proc = -1.0
	if len(akey) != 48 {
		return errors.New("invalidAKEY")
	}
	ckey := Genrand(1920)
	ckcpy := append(make([]byte, 0), ckey...) // for writing header
	var tgtb []byte                           // for final calc
	var tgti int                              // for final calc
	flow0 := make([]chan []byte, 40)          // read -> calc
	flow1 := make([]chan []byte, 40)          // calc -> write
	for i := 0; i < 40; i++ {
		flow0[i] = make(chan []byte, 12)
		flow1[i] = make(chan []byte, 12)
	}
	tbox.Before.dpos = 0                 // reader tuning
	tbox.Before.dsize = tbox.Before.Size // reader tuning
	go readflow(&tbox.Before, flow0, true, &tbox.Proc, &tgtb, &tgti)
	for i := 0; i < 40; i++ {
		go calcflow(flow0, flow1, true, ckey, i)
	}
	header := aescalc(ckcpy, akey[16:48], akey[0:16], true, false)
	writeflow(&tbox.After, flow1, true, header, ckey, &tgtb, &tgti)
	return nil
}

// decrypt with akey 48B
func (tbox *Funcmode) Decrypt(akey []byte) error {
	defer func() { tbox.Proc = 2.0 }()
	tbox.Proc = -1.0
	if len(akey) != 48 {
		return errors.New("invalidAKEY")
	}
	ckeydata := tbox.Before.Read(1920)
	tbox.Before.dpos = 1920                     // reader tuning
	tbox.Before.dsize = tbox.Before.Size - 1920 // reader tuning
	ckey := aescalc(ckeydata, akey[16:48], akey[0:16], false, false)
	var tgtb []byte                  // for final calc
	var tgti int                     // for final calc
	flow0 := make([]chan []byte, 40) // read -> calc
	flow1 := make([]chan []byte, 40) // calc -> write
	for i := 0; i < 40; i++ {
		flow0[i] = make(chan []byte, 12)
		flow1[i] = make(chan []byte, 12)
	}
	go readflow(&tbox.Before, flow0, false, &tbox.Proc, &tgtb, &tgti)
	for i := 0; i < 40; i++ {
		go calcflow(flow0, flow1, false, ckey, i)
	}
	writeflow(&tbox.After, flow1, false, nil, ckey, &tgtb, &tgti)
	return nil
}

// generate secure random nB
func Genrand(size int) []byte {
	temp := make([]byte, size)
	rand.Read(temp)
	return temp
}

// returns gen5kaes basic keyfile
func Basickey() []byte {
	var temp []byte
	temp = append(temp, 234, 183, 184, 235, 158, 152, 44, 32, 235, 130, 152, 235, 165, 188, 32, 234, 176, 128, 235, 145, 172, 235, 145, 148, 32, 236, 177, 132, 32, 236, 157, 180, 234, 179, 179, 234, 185, 140, 236, 167)
	temp = append(temp, 128, 32, 236, 158, 172, 235, 176, 140, 235, 138, 148, 32, 236, 151, 172, 236, 160, 149, 32, 235, 179, 180, 235, 131, 136, 235, 139, 136, 63, 13, 10, 235, 172, 180, 235, 132, 136, 236, 160, 184)
	temp = append(temp, 235, 157, 188, 46, 13, 10, 236, 134, 159, 236, 149, 132, 235, 157, 188, 46, 13, 10, 236, 154, 184, 235, 160, 164, 235, 157, 188, 46, 13, 10, 236, 152, 155, 235, 130, 160, 236, 157, 152, 32)
	temp = append(temp, 234, 176, 144, 234, 176, 129, 235, 147, 164, 236, 157, 180, 32, 235, 143, 140, 236, 149, 132, 236, 152, 164, 235, 138, 148, 234, 181, 172, 235, 130, 152, 46, 13, 10, 236, 157, 180, 235, 178, 136)
	temp = append(temp, 236, 151, 148, 32, 236, 138, 164, 236, 138, 164, 235, 161, 156, 236, 157, 152, 32, 237, 158, 152, 236, 156, 188, 235, 161, 156, 32, 235, 130, 152, 235, 165, 188, 32, 235, 167, 137, 236, 157, 132)
	temp = append(temp, 32, 236, 136, 152, 32, 236, 158, 136, 234, 178, 160, 235, 139, 136, 63, 13, 10, 237, 157, 169, 236, 150, 180, 236, 167, 128, 234, 177, 176, 235, 157, 188, 46, 13, 10, 236, 154, 148, 236, 160)
	temp = append(temp, 149, 235, 147, 164, 236, 157, 180, 236, 151, 172, 46, 13, 10, 235, 168, 184, 235, 166, 172, 235, 165, 188, 32, 236, 134, 141, 236, 157, 184, 32, 236, 177, 132, 32, 236, 157, 180, 32, 234, 179)
	temp = append(temp, 179, 236, 151, 144, 236, 132, 156, 32, 236, 157, 180, 235, 159, 176, 32, 235, 139, 185, 235, 143, 140, 237, 149, 156, 32, 236, 167, 147, 236, 157, 132, 32, 235, 152, 144, 32, 235, 139, 164, 236)
	temp = append(temp, 139, 156, 32, 235, 178, 140, 236, 157, 180, 234, 179, 160, 32, 236, 158, 136, 234, 181, 172, 235, 130, 152, 46, 13, 10, 236, 167, 145, 236, 134, 141, 46, 13, 10, 236, 151, 180, 236, 135, 160)
	temp = append(temp, 32, 236, 157, 145, 236, 182, 149, 46, 13, 10, 234, 176, 156, 235, 176, 169, 46, 13, 10, 236, 157, 180, 32, 235, 170, 184, 236, 157, 128, 32, 236, 157, 180, 235, 159, 176, 32, 237, 158, 152)
	temp = append(temp, 235, 143, 132, 32, 236, 147, 184, 32, 236, 136, 152, 32, 236, 158, 136, 234, 181, 172, 235, 130, 152, 46, 32, 237, 157, 165, 235, 175, 184, 235, 161, 173, 236, 167, 128, 235, 167, 140, 32, 235)
	temp = append(temp, 175, 184, 236, 149, 189, 237, 149, 180, 46, 13, 10, 236, 157, 180, 235, 178, 136, 236, 151, 148, 32, 234, 183, 184, 235, 133, 128, 236, 157, 152, 32, 235, 143, 132, 236, 155, 128, 32, 236, 151)
	temp = append(temp, 134, 236, 157, 180, 32, 235, 167, 137, 236, 149, 132, 235, 179, 180, 235, 160, 164, 235, 172, 180, 235, 130, 152, 46, 13, 10, 235, 182, 132, 236, 132, 157, 46, 32, 236, 149, 149, 236, 182, 149)
	temp = append(temp, 46, 32, 236, 160, 132, 234, 176, 156, 46, 13, 10, 235, 130, 180, 32, 236, 149, 158, 236, 151, 144, 32, 236, 132, 156, 236, 167, 128, 32, 235, 167, 144, 234, 177, 176, 235, 157, 188, 46, 13)
	temp = append(temp, 10, 236, 157, 180, 32, 235, 170, 184, 236, 157, 128, 32, 236, 160, 156, 236, 149, 189, 236, 157, 180, 32, 235, 132, 136, 235, 172, 180, 32, 235, 167, 142, 236, 149, 132, 46, 13, 10, 235, 172)
	temp = append(temp, 180, 235, 132, 136, 236, 160, 184, 235, 130, 180, 235, 160, 164, 235, 157, 188, 46, 13, 10, 235, 130, 152, 236, 152, 164, 234, 177, 176, 235, 157, 188, 46, 13, 10, 235, 130, 160, 235, 155, 176)
	temp = append(temp, 234, 177, 176, 235, 157, 188, 46, 13, 10, 236, 157, 180, 234, 179, 179, 236, 157, 152, 32, 237, 158, 152, 236, 157, 132, 32, 235, 141, 148, 32, 235, 168, 188, 236, 160, 128, 32, 236, 149, 140)
	temp = append(temp, 236, 149, 152, 235, 139, 164, 235, 169, 180, 32, 236, 154, 176, 235, 166, 172, 235, 147, 164, 235, 143, 132, 32, 234, 183, 184, 32, 237, 158, 152, 236, 157, 132, 32, 236, 147, 184, 32, 236, 136)
	temp = append(temp, 152, 32, 236, 158, 136, 236, 151, 136, 234, 178, 160, 236, 167, 128, 46, 13, 10, 235, 130, 152, 235, 165, 188, 32, 235, 132, 152, 236, 167, 128, 32, 235, 170, 187, 237, 149, 152, 235, 169, 180)
	temp = append(temp, 32, 234, 178, 176, 234, 181, 173, 32, 235, 152, 144, 32, 235, 139, 164, 236, 139, 156, 32, 235, 168, 184, 235, 166, 172, 236, 151, 144, 32, 235, 176, 159, 237, 158, 144, 32, 235, 191, 144, 236)
	temp = append(temp, 157, 180, 235, 158, 128, 235, 139, 164, 46, 13, 10, 235, 172, 184, 236, 157, 132, 32, 236, 151, 180, 236, 150, 180, 236, 163, 188, 235, 167, 136, 46, 13, 10, 236, 157, 180, 234, 179, 179, 236)
	temp = append(temp, 151, 144, 32, 237, 149, 168, 234, 187, 152, 32, 234, 176, 128, 235, 157, 188, 236, 149, 137, 236, 158, 144, 234, 181, 172, 235, 130, 152, 46, 13, 10, 236, 157, 180, 32, 234, 176, 144, 236, 152)
	temp = append(temp, 165, 236, 151, 144, 236, 132, 156, 32, 235, 130, 152, 234, 176, 132, 235, 139, 164, 32, 237, 149, 152, 235, 141, 148, 235, 157, 188, 235, 143, 132, 32, 237, 152, 188, 236, 158, 144, 236, 132, 156)
	temp = append(temp, 32, 235, 172, 180, 236, 151, 135, 236, 157, 132, 32, 237, 149, 160, 32, 236, 136, 152, 32, 236, 158, 136, 236, 157, 132, 32, 234, 178, 131, 32, 234, 176, 153, 235, 139, 136, 63, 13, 10, 236)
	temp = append(temp, 157, 180, 32, 234, 181, 180, 235, 160, 136, 235, 165, 188, 32, 235, 129, 138, 235, 138, 148, 235, 139, 164, 32, 237, 149, 152, 235, 141, 148, 235, 157, 188, 235, 143, 132, 32, 236, 158, 160, 236)
	temp = append(temp, 139, 156, 235, 191, 144, 236, 157, 180, 236, 167, 128, 46, 13, 10, 236, 152, 133, 236, 150, 180, 236, 167, 128, 235, 138, 148, 234, 181, 172, 235, 130, 152, 46, 13, 10, 236, 158, 160, 236, 157)
	temp = append(temp, 180, 32, 236, 152, 164, 235, 138, 148, 234, 181, 172, 235, 130, 152, 46, 13, 10, 236, 154, 176, 235, 166, 172, 235, 165, 188, 32, 235, 178, 151, 236, 150, 180, 235, 130, 160, 32, 236, 136, 152)
	temp = append(temp, 235, 138, 148, 32, 236, 151, 134, 235, 139, 168, 235, 139, 164, 46, 13, 10, 234, 184, 176, 237, 154, 140, 235, 165, 188, 32, 235, 134, 147, 236, 179, 164, 234, 181, 172, 235, 130, 152, 46, 13)
	temp = append(temp, 10, 234, 181, 189, 237, 158, 136, 236, 167, 128, 32, 235, 170, 187, 237, 150, 136, 234, 181, 172, 235, 130, 152, 46, 13, 10, 237, 152, 188, 236, 158, 144, 236, 132, 156, 235, 138, 148, 32, 237)
	temp = append(temp, 157, 144, 235, 166, 132, 236, 157, 132, 32, 235, 169, 136, 236, 182, 156, 32, 236, 136, 152, 32, 236, 151, 134, 235, 139, 168, 235, 139, 164, 46, 13, 10, 237, 140, 140, 235, 143, 132, 235, 138)
	temp = append(temp, 148, 32, 235, 139, 164, 236, 139, 156, 32, 236, 157, 188, 235, 160, 129, 236, 157, 188, 32, 234, 178, 131, 236, 157, 180, 235, 158, 128, 235, 139, 164, 46, 13, 10, 235, 132, 136, 236, 157, 152)
	temp = append(temp, 32, 235, 175, 184, 236, 136, 153, 237, 149, 168, 236, 157, 180, 235, 139, 164, 46, 13, 10, 234, 177, 176, 235, 140, 128, 237, 149, 156, 32, 237, 157, 144, 235, 166, 132, 236, 157, 132, 32, 236)
	temp = append(temp, 134, 144, 235, 176, 148, 235, 139, 165, 236, 156, 188, 235, 161, 156, 32, 235, 167, 137, 236, 157, 132, 32, 236, 136, 152, 32, 236, 151, 134, 235, 139, 168, 235, 139, 164, 46, 13, 10, 236, 152)
	temp = append(temp, 164, 235, 161, 175, 236, 157, 180, 32, 235, 132, 136, 236, 157, 152, 32, 237, 158, 152, 235, 167, 140, 236, 156, 188, 235, 161, 156, 32, 235, 167, 137, 236, 150, 180, 235, 179, 180, 235, 160, 164)
	temp = append(temp, 235, 172, 180, 235, 130, 152, 46, 13, 10, 237, 155, 140, 235, 165, 173, 237, 149, 152, 234, 181, 172, 235, 130, 152, 46, 13, 10, 234, 183, 184, 235, 158, 152, 44, 32, 236, 157, 180, 32, 236)
	temp = append(temp, 160, 149, 235, 143, 132, 235, 169, 180, 32, 236, 167, 128, 236, 188, 156, 235, 179, 188, 32, 234, 176, 128, 236, 185, 152, 234, 176, 128, 32, 236, 158, 136, 234, 178, 160, 236, 167, 128, 46, 13)
	temp = append(temp, 10, 235, 130, 180, 234, 176, 128, 32, 236, 150, 180, 235, 150, 187, 234, 178, 140, 32, 234, 183, 184, 235, 166, 172, 32, 237, 149, 156, 236, 151, 134, 236, 157, 180, 32, 236, 158, 148, 236, 157)
	temp = append(temp, 184, 237, 149, 180, 236, 167, 136, 32, 236, 136, 152, 32, 236, 158, 136, 236, 151, 136, 235, 138, 148, 236, 167, 128, 32, 236, 149, 140, 235, 160, 164, 236, 164, 132, 234, 185, 140, 63, 13, 10)
	temp = append(temp, 236, 130, 172, 235, 158, 140, 235, 147, 164, 236, 157, 128, 32, 235, 170, 168, 235, 145, 144, 32, 235, 182, 136, 236, 149, 136, 236, 157, 132, 32, 234, 176, 128, 236, 167, 132, 32, 236, 177, 132)
	temp = append(temp, 235, 161, 156, 32, 236, 130, 180, 236, 149, 132, 234, 176, 132, 235, 139, 168, 235, 139, 164, 46, 13, 10, 236, 157, 180, 234, 177, 180, 32, 235, 175, 184, 236, 167, 128, 236, 157, 152, 32, 236)
	temp = append(temp, 152, 129, 236, 151, 173, 236, 157, 132, 32, 235, 167, 158, 235, 139, 165, 235, 156, 168, 235, 166, 180, 32, 235, 149, 140, 32, 235, 138, 144, 235, 129, 188, 235, 138, 148, 32, 235, 139, 185, 236)
	temp = append(temp, 151, 176, 237, 149, 156, 32, 235, 140, 128, 234, 176, 128, 236, 149, 188, 46, 13, 10, 237, 149, 152, 236, 167, 128, 235, 167, 140, 32, 235, 130, 152, 235, 138, 148, 32, 236, 132, 184, 236, 131)
	temp = append(temp, 129, 236, 151, 144, 32, 236, 130, 180, 236, 149, 132, 235, 130, 168, 234, 184, 176, 32, 236, 156, 132, 237, 149, 180, 32, 234, 183, 184, 32, 234, 179, 181, 237, 143, 172, 235, 165, 188, 32, 235)
	temp = append(temp, 176, 155, 236, 149, 132, 235, 147, 164, 236, 157, 180, 236, 167, 128, 32, 236, 149, 138, 236, 157, 128, 32, 236, 177, 132, 32, 236, 138, 164, 236, 138, 164, 235, 161, 156, 32, 235, 168, 185, 236)
	temp = append(temp, 150, 180, 235, 178, 132, 235, 160, 184, 236, 150, 180, 46, 13, 10, 234, 183, 184, 234, 178, 131, 236, 157, 180, 32, 235, 130, 180, 234, 176, 128, 32, 236, 160, 128, 236, 167, 128, 235, 165, 184)
	temp = append(temp, 32, 236, 181, 156, 236, 180, 136, 236, 157, 180, 236, 158, 144, 32, 236, 181, 156, 236, 149, 133, 236, 157, 152, 32, 236, 149, 133, 237, 150, 137, 236, 157, 180, 236, 151, 136, 236, 167, 128, 46)
	temp = append(temp, 13, 10, 237, 155, 132, 237, 154, 140, 237, 149, 152, 236, 167, 132, 32, 236, 149, 138, 236, 149, 132, 46, 32, 236, 131, 157, 236, 161, 180, 236, 157, 132, 32, 236, 156, 132, 237, 149, 156, 32)
	temp = append(temp, 236, 132, 160, 237, 131, 157, 236, 157, 180, 236, 151, 136, 236, 156, 188, 235, 139, 136, 46, 13, 10, 235, 132, 136, 235, 143, 132, 32, 235, 167, 136, 236, 176, 172, 234, 176, 128, 236, 167, 128)
	temp = append(temp, 32, 236, 149, 132, 235, 139, 136, 236, 151, 136, 235, 139, 136, 63, 13, 10, 235, 130, 180, 32, 235, 168, 184, 235, 166, 191, 236, 134, 141, 236, 157, 132, 32, 234, 176, 136, 234, 184, 176, 234)
	temp = append(temp, 176, 136, 234, 184, 176, 32, 235, 182, 132, 237, 149, 180, 236, 139, 156, 237, 130, 164, 235, 141, 152, 32, 235, 132, 136, 236, 157, 152, 32, 237, 145, 156, 236, 160, 149, 236, 157, 132, 32, 235)
	temp = append(temp, 179, 180, 234, 179, 160, 32, 236, 149, 140, 32, 236, 136, 152, 32, 236, 158, 136, 236, 151, 136, 236, 150, 180, 46, 13, 10, 236, 150, 180, 236, 169, 148, 32, 236, 136, 152, 32, 236, 151, 134)
	temp = append(temp, 236, 151, 136, 235, 139, 164, 32, 235, 157, 188, 234, 179, 160, 32, 235, 132, 136, 235, 143, 132, 32, 235, 167, 144, 236, 157, 132, 32, 237, 149, 152, 235, 160, 164, 235, 130, 152, 46, 13, 10)
	temp = append(temp, 236, 157, 180, 236, 160, 156, 235, 138, 148, 32, 235, 139, 164, 236, 139, 156, 32, 236, 158, 138, 234, 179, 160, 32, 236, 158, 136, 235, 141, 152, 32, 234, 179, 181, 237, 143, 172, 235, 165, 188)
	temp = append(temp, 32, 235, 167, 136, 236, 163, 188, 237, 149, 152, 235, 138, 148, 32, 234, 178, 140, 32, 236, 162, 139, 236, 157, 132, 32, 234, 177, 176, 236, 149, 188, 46, 13, 10, 234, 181, 180, 235, 160, 136)
	temp = append(temp, 235, 165, 188, 32, 235, 129, 138, 236, 150, 180, 235, 130, 180, 234, 178, 160, 235, 139, 164, 235, 169, 180, 46, 13, 10, 235, 132, 164, 234, 176, 128, 32, 234, 183, 184, 235, 160, 135, 235, 139)
	temp = append(temp, 164, 235, 169, 180, 32, 234, 183, 184, 235, 159, 176, 32, 234, 178, 131, 236, 157, 180, 234, 178, 160, 236, 167, 128, 46, 13, 10)
	return temp
}
