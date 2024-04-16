// test645 : stdlib5.ksc

package ksc

import (
	"bytes"
	"errors"
	"hash/crc32"
	"stdlib5/kio"
)

type toolbox struct {
	Prehead  []byte // prehead + pad 512nB
	Common   []byte // common sign 4B
	Subtype  []byte // subtype sign 4B
	Reserved []byte // reserved 8B

	Headp int    // mainheader start point (512n)
	Rsize int    // valid data size
	Path  string // file path

	Predetect bool  // read chunk info ahead
	Chunkpos  []int // chunk (8B + nB) pos
	Chunksize []int // chunk data (nB) size
}

// read by file path
func (tbox *toolbox) Readf() error {
	var temp []byte
	tbox.Prehead = make([]byte, 0)
	tbox.Headp = 0
	size := kio.Size(tbox.Path)
	if (tbox.Path == "") || (size == -1) {
		return errors.New("NoSuchFile")
	}

	f, err := kio.Open(tbox.Path, "r")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	for size >= tbox.Headp+4 {
		temp, _ = kio.Read(f, 4)
		if kio.Bequal(temp, tbox.Common) {
			break
		} else {
			tbox.Prehead = append(tbox.Prehead, temp...)
			temp, _ = kio.Read(f, 508)
			tbox.Prehead = append(tbox.Prehead, temp...)
			tbox.Headp = tbox.Headp + 512
		}
	}
	if size < tbox.Headp+16 {
		return errors.New("InvalidKSC5")
	}

	tbox.Subtype, _ = kio.Read(f, 4)
	tbox.Reserved, _ = kio.Read(f, 8)
	tbox.Rsize = size - tbox.Headp

	tbox.Chunkpos = make([]int, 0)
	tbox.Chunksize = make([]int, 0)
	pos := tbox.Headp + 16
	if tbox.Predetect {
		for size >= pos+8 {
			temp, _ = kio.Read(f, 8)
			if kio.Bequal(temp, []byte{255, 255, 255, 255, 255, 255, 255, 255}) {
				break
			} else {
				ti := Decode(temp)
				tbox.Chunkpos = append(tbox.Chunkpos, pos)
				tbox.Chunksize = append(tbox.Chunksize, ti)
				pos = pos + 8 + ti
				te0, te1 := f.Seek(int64(pos), 0)
				if (int(te0) != pos) || (te1 != nil) {
					return te1
				}
			}
		}
	}
	return nil
}

// init & write mainhead by file path
func (tbox *toolbox) Writef() error {
	tbox.Headp = len(tbox.Prehead)
	tbox.Rsize = 0
	tbox.Chunkpos = make([]int, 0)
	tbox.Chunksize = make([]int, 0)
	if (tbox.Headp%512 != 0) || (len(tbox.Common) != 4) || (len(tbox.Subtype) != 4) || (len(tbox.Reserved) != 8) {
		return errors.New("InvalidKSC5")
	}
	if tbox.Path == "" {
		return errors.New("NoSuchFile")
	}

	f, err := kio.Open(tbox.Path, "w")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}
	temp := bytes.Join([][]byte{tbox.Prehead, tbox.Common, tbox.Subtype, tbox.Reserved}, make([]byte, 0))
	_, err = kio.Write(f, temp)
	return err
}

// write one chunk by content of file, empty path : end (8x FF)
func (tbox *toolbox) Addf(path string) error {
	f, err := kio.Open(tbox.Path, "a")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}

	if path == "" {
		_, err = kio.Write(f, []byte{255, 255, 255, 255, 255, 255, 255, 255})
		return err
	} else {
		t, err := kio.Open(path, "r")
		if err == nil {
			defer t.Close()
		} else {
			return err
		}

		size := kio.Size(path)
		_, err = kio.Write(f, Encode(size, 8))
		if err != nil {
			return err
		}

		for i := 0; i < size/10485760; i++ {
			temp, err := kio.Read(t, 10485760)
			if err != nil {
				return err
			}
			_, err = kio.Write(f, temp)
			if err != nil {
				return err
			}
		}
		temp, err := kio.Read(t, size%10485760)
		if err != nil {
			return err
		}
		_, err = kio.Write(f, temp)
		return err
	}
}

// write one chunk by binary data
func (tbox *toolbox) Linkf(data []byte) error {
	f, err := kio.Open(tbox.Path, "a")
	if err == nil {
		defer f.Close()
	} else {
		return err
	}

	_, err = kio.Write(f, Encode(len(data), 8))
	if err != nil {
		return err
	}
	_, err = kio.Write(f, data)
	return err
}

// read by binary data
func (tbox *toolbox) Readb(data []byte) error {
	var temp []byte
	tbox.Prehead = make([]byte, 0)
	tbox.Headp = 0
	size := len(data)

	for size >= tbox.Headp+4 {
		temp = data[tbox.Headp : tbox.Headp+4]
		if kio.Bequal(temp, tbox.Common) {
			break
		} else {
			tbox.Prehead = append(tbox.Prehead, data[tbox.Headp:tbox.Headp+512]...)
			tbox.Headp = tbox.Headp + 512
		}
	}
	if size < tbox.Headp+16 {
		return errors.New("InvalidKSC5")
	}

	tbox.Subtype = make([]byte, 4)
	copy(tbox.Subtype, data[tbox.Headp+4:tbox.Headp+8])
	tbox.Reserved = make([]byte, 8)
	copy(tbox.Reserved, data[tbox.Headp+8:tbox.Headp+16])
	tbox.Rsize = size - tbox.Headp

	tbox.Chunkpos = make([]int, 0)
	tbox.Chunksize = make([]int, 0)
	pos := tbox.Headp + 16

	if tbox.Predetect {
		for size >= pos+8 {
			temp = data[pos : pos+8]
			if kio.Bequal(temp, []byte{255, 255, 255, 255, 255, 255, 255, 255}) {
				break
			} else {
				ti := Decode(temp)
				tbox.Chunkpos = append(tbox.Chunkpos, pos)
				tbox.Chunksize = append(tbox.Chunksize, ti)
				pos = pos + 8 + ti
			}
		}
	}
	return nil
}

// init & write mainhead by binary data
func (tbox *toolbox) Writeb() ([]byte, error) {
	tbox.Headp = len(tbox.Prehead)
	tbox.Rsize = 0
	tbox.Chunkpos = make([]int, 0)
	tbox.Chunksize = make([]int, 0)
	if (tbox.Headp%512 != 0) || (len(tbox.Common) != 4) || (len(tbox.Subtype) != 4) || (len(tbox.Reserved) != 8) {
		return nil, errors.New("InvalidKSC5")
	}
	return bytes.Join([][]byte{tbox.Prehead, tbox.Common, tbox.Subtype, tbox.Reserved}, make([]byte, 0)), nil
}

// write one chunk by content of file, empty path : end (8x FF)
func (tbox *toolbox) Addb(stream []byte, path string) ([]byte, error) {
	var out []byte
	if path == "" {
		out = []byte{255, 255, 255, 255, 255, 255, 255, 255}
	} else {
		size := kio.Size(path)
		out = Encode(size, 8)
		t, err := kio.Open(path, "r")
		if err == nil {
			defer t.Close()
		} else {
			return nil, err
		}

		t0, t1 := kio.Read(t, -1)
		if t1 != nil {
			return nil, t1
		}
		out = append(out, t0...)
	}
	return append(stream, out...), nil
}

// write one chunk by binary data
func (tbox *toolbox) Linkb(stream []byte, data []byte) []byte {
	return bytes.Join([][]byte{stream, Encode(len(data), 8), data}, make([]byte, 0))
}

// Init and return toolbox
func Initksc() toolbox {
	var out toolbox
	out.Prehead = Webpbase()
	out.Common = []byte("KSC5")
	out.Subtype = make([]byte, 4)
	out.Reserved = make([]byte, 8)
	out.Headp = 512
	out.Rsize = 0
	out.Path = ""
	out.Predetect = false
	out.Chunkpos = make([]int, 0)
	out.Chunksize = make([]int, 0)
	return out
}

// little endian encoding, I -> B
func Encode(num int, length int) []byte {
	temp := make([]byte, length)
	for i := 0; i < length; i++ {
		temp[i] = byte(num % 256)
		num = num / 256
	}
	return temp
}

// little endian decoding, B -> I
func Decode(data []byte) int {
	temp := 0
	for i, r := range data {
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

// crc32
func Crc32hash(data []byte) []byte {
	temp := int(crc32.ChecksumIEEE(data))
	return Encode(temp, 4)
}

// basic KSC5 prehead webp data
func Webpbase() []byte {
	var temp []byte
	temp = append(temp, 82, 73, 70, 70, 188, 1, 0, 0, 87, 69, 66, 80, 86, 80, 56, 32, 176, 1, 0, 0, 80, 9, 0, 157, 1, 42, 64, 0, 64, 0, 62, 105, 44, 144, 69, 164, 34, 161, 154, 250)
	temp = append(temp, 52, 204, 64, 6, 132, 179, 128, 103, 44, 209, 255, 250, 122, 113, 150, 247, 40, 54, 130, 136, 14, 71, 191, 152, 186, 117, 8, 46, 237, 207, 106, 253, 13, 243, 237, 33, 194, 46, 252, 192)
	temp = append(temp, 196, 135, 35, 241, 247, 91, 154, 132, 230, 20, 63, 91, 193, 27, 117, 80, 185, 174, 228, 148, 182, 73, 192, 0, 236, 143, 194, 185, 18, 220, 156, 42, 84, 77, 31, 219, 232, 240, 221, 0)
	temp = append(temp, 195, 202, 255, 72, 196, 75, 249, 80, 118, 132, 114, 104, 223, 212, 237, 201, 156, 224, 120, 245, 244, 97, 110, 99, 82, 216, 21, 124, 229, 35, 213, 59, 151, 9, 103, 137, 106, 139, 227, 24)
	temp = append(temp, 176, 138, 164, 65, 87, 138, 151, 107, 200, 6, 44, 208, 214, 206, 230, 187, 66, 111, 140, 1, 95, 165, 111, 52, 207, 190, 129, 146, 51, 63, 156, 170, 37, 22, 41, 83, 39, 161, 72, 216)
	temp = append(temp, 142, 88, 195, 253, 1, 52, 82, 54, 108, 37, 252, 192, 104, 237, 44, 13, 149, 154, 138, 211, 137, 144, 247, 42, 162, 59, 176, 198, 212, 245, 25, 81, 52, 188, 114, 47, 45, 184, 70, 50)
	temp = append(temp, 139, 89, 51, 8, 177, 118, 55, 119, 99, 65, 165, 14, 238, 85, 125, 115, 46, 91, 142, 80, 153, 65, 254, 150, 242, 52, 54, 58, 218, 208, 217, 23, 137, 19, 160, 155, 171, 249, 72, 189)
	temp = append(temp, 202, 248, 54, 51, 14, 152, 217, 232, 149, 148, 42, 176, 146, 66, 1, 71, 224, 83, 167, 253, 101, 56, 106, 101, 107, 57, 60, 108, 202, 45, 81, 50, 51, 205, 104, 92, 223, 35, 174, 88)
	temp = append(temp, 79, 14, 125, 187, 219, 152, 86, 98, 196, 97, 148, 42, 221, 81, 205, 121, 235, 56, 153, 74, 74, 107, 166, 45, 123, 197, 186, 187, 146, 154, 244, 125, 183, 220, 87, 198, 253, 233, 102, 42)
	temp = append(temp, 143, 248, 4, 243, 40, 51, 12, 181, 18, 139, 57, 208, 170, 206, 174, 172, 144, 119, 45, 227, 242, 176, 35, 33, 167, 79, 38, 251, 81, 69, 158, 39, 195, 176, 106, 106, 11, 211, 253, 60)
	temp = append(temp, 49, 101, 5, 250, 33, 100, 160, 54, 30, 93, 139, 212, 233, 215, 225, 89, 114, 47, 12, 28, 71, 62, 28, 191, 224, 21, 192, 87, 162, 48, 71, 198, 77, 2, 189, 191, 111, 182, 252, 18)
	temp = append(temp, 16, 76, 59, 222, 234, 221, 123, 0, 0, 0, 0, 0)

	temp = append(temp, make([]byte, 512-len(temp))...)
	return temp
}
