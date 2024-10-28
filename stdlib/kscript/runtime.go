// test691 : stdlib5.kscript runtime go

package kscript

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"slices"
	"stdlib5/kdb"
	"stdlib5/kio"
	"stdlib5/ksc"
	"stdlib5/ksign"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
)

// variable unit (part 6)
type Vunit struct {
	Vtype   byte    // 0 : none, 1 : bool, 2 : int, 3 : float, 4 : string, 5 : bytes
	Vbool   bool    // bool
	Vint    int     // int
	Vfloat  float64 // float
	Vstring string  // string
	Vbytes  []byte  // bytes
}

// set Vunit with value
func (vu *Vunit) Set(v interface{}) {
	vu.Vtype = 0
	vu.Vbool = false
	vu.Vint = 0
	vu.Vfloat = 0.0
	vu.Vstring = ""
	vu.Vbytes = nil

	switch v := v.(type) {
	case bool:
		vu.Vtype = 1
		vu.Vbool = v
	case int:
		vu.Vtype = 2
		vu.Vint = v
	case float64:
		vu.Vtype = 3
		vu.Vfloat = v
	case string:
		vu.Vtype = 4
		vu.Vstring = v
	case []byte:
		vu.Vtype = 5
		vu.Vbytes = v
	}
}

// get string value
func (vu *Vunit) ToString() string {
	switch vu.Vtype {
	case 1:
		if vu.Vbool {
			return "True"
		} else {
			return "False"
		}
	case 2:
		return fmt.Sprint(vu.Vint)
	case 3:
		return fmt.Sprintf("%f", vu.Vfloat)
	case 4:
		return vu.Vstring
	case 5:
		return kio.Bprint(vu.Vbytes)
	default:
		return "None"
	}
}

// arithmetic logic calculation : get 2 operand, set itself result

func (vu *Vunit) add(oper0 *Vunit, oper1 *Vunit) error {
	switch oper0.Vtype {
	case 1: // bool
		if oper1.Vtype == 1 { // bool
			vu.Vtype = 1
			vu.Vbool = oper0.Vbool || oper1.Vbool
		} else {
			return errors.New("e600 : cannot add bool, ?")
		}
	case 2: // int
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 2
			vu.Vint = oper0.Vint + oper1.Vint
		case 3: // float
			vu.Vtype = 3
			vu.Vfloat = float64(oper0.Vint) + oper1.Vfloat
		default:
			return errors.New("e601 : cannot add int, ?")
		}
	case 3: // float
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 3
			vu.Vfloat = oper0.Vfloat + float64(oper1.Vint)
		case 3: // float
			vu.Vtype = 3
			vu.Vfloat = oper0.Vfloat + oper1.Vfloat
		default:
			return errors.New("e602 : cannot add float, ?")
		}
	case 4: // string
		if oper1.Vtype == 4 { // string
			vu.Vtype = 4
			vu.Vstring = oper0.Vstring + oper1.Vstring
		} else {
			return errors.New("e603 : cannot add string, ?")
		}
	case 5: // bytes
		if oper1.Vtype == 5 { // bytes
			vu.Vtype = 5
			vu.Vbytes = append(oper0.Vbytes, oper1.Vbytes...)
		} else {
			return errors.New("e604 : cannot add bytes, ?")
		}
	default: // nil
		return errors.New("e605 : cannot add none, ?")
	}
	return nil
}

func (vu *Vunit) sub(oper0 *Vunit, oper1 *Vunit) error {
	switch oper0.Vtype {
	case 2: // int
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 2
			vu.Vint = oper0.Vint - oper1.Vint
		case 3: // float
			vu.Vtype = 3
			vu.Vfloat = float64(oper0.Vint) - oper1.Vfloat
		default:
			return errors.New("e606 : cannot sub int, ?")
		}
	case 3: // float
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 3
			vu.Vfloat = oper0.Vfloat - float64(oper1.Vint)
		case 3: // float
			vu.Vtype = 3
			vu.Vfloat = oper0.Vfloat - oper1.Vfloat
		default:
			return errors.New("e607 : cannot sub float, ?")
		}
	default: // bool, string, bytes, nil
		return errors.New("e608 : cannot sub ?, ?")
	}
	return nil
}

func (vu *Vunit) mul(oper0 *Vunit, oper1 *Vunit) error {
	switch oper0.Vtype {
	case 1: // bool
		if oper1.Vtype == 1 { // bool
			vu.Vtype = 1
			vu.Vbool = oper0.Vbool && oper1.Vbool
		} else {
			return errors.New("e609 : cannot mul bool, ?")
		}
	case 2: // int
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 2
			vu.Vint = oper0.Vint * oper1.Vint
		case 3: // float
			vu.Vtype = 3
			vu.Vfloat = float64(oper0.Vint) * oper1.Vfloat
		case 4: // string
			if oper0.Vint < 0 {
				return errors.New("e610 : cannot mul int_n, string")
			} else {
				vu.Vtype = 4
				vu.Vstring = strings.Repeat(oper1.Vstring, oper0.Vint)
			}
		case 5: // bytes
			if oper0.Vint < 0 {
				return errors.New("e611 : cannot mul int_n, bytes")
			} else {
				vu.Vtype = 5
				vu.Vbytes = slices.Repeat(oper1.Vbytes, oper0.Vint)
			}
		default:
			return errors.New("e612 : cannot mul int, ?")
		}
	case 3: // float
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 3
			vu.Vfloat = oper0.Vfloat * float64(oper1.Vint)
		case 3: // float
			vu.Vtype = 3
			vu.Vfloat = oper0.Vfloat * oper1.Vfloat
		default:
			return errors.New("e613 : cannot mul float, ?")
		}
	case 4: // string
		if oper1.Vtype == 2 { // int
			if oper1.Vint < 0 {
				return errors.New("e614 : cannot mul string, int_n")
			} else {
				vu.Vtype = 4
				vu.Vstring = strings.Repeat(oper0.Vstring, oper1.Vint)
			}
		} else {
			return errors.New("e615 : cannot mul string, ?")
		}
	case 5: // bytes
		if oper1.Vtype == 2 { // int
			if oper1.Vint < 0 {
				return errors.New("e616 : cannot mul bytes, int_n")
			} else {
				vu.Vtype = 5
				vu.Vbytes = slices.Repeat(oper0.Vbytes, oper1.Vint)
			}
		} else {
			return errors.New("e617 : cannot mul bytes, ?")
		}
	default: // nil
		return errors.New("e618 : cannot mul none, ?")
	}
	return nil
}

func (vu *Vunit) div(oper0 *Vunit, oper1 *Vunit) error {
	switch oper0.Vtype {
	case 2: // int
		switch oper1.Vtype {
		case 2: // int
			if oper1.Vint == 0 {
				return errors.New("e619 : cannot div int, int_0")
			} else {
				vu.Vtype = 3
				vu.Vfloat = float64(oper0.Vint) / float64(oper1.Vint)
			}
		case 3: // float
			if oper1.Vfloat == 0.0 {
				return errors.New("e620 : cannot div int, float_0")
			} else {
				vu.Vtype = 3
				vu.Vfloat = float64(oper0.Vint) / oper1.Vfloat
			}
		default:
			return errors.New("e621 : cannot div int, ?")
		}
	case 3: // float
		switch oper1.Vtype {
		case 2: // int
			if oper1.Vint == 0 {
				return errors.New("e622 : cannot div float, int_0")
			} else {
				vu.Vtype = 3
				vu.Vfloat = oper0.Vfloat / float64(oper1.Vint)
			}
		case 3: // float
			if oper1.Vfloat == 0.0 {
				return errors.New("e623 : cannot div float, float_0")
			} else {
				vu.Vtype = 3
				vu.Vfloat = oper0.Vfloat / oper1.Vfloat
			}
		default:
			return errors.New("e624 : cannot div float, ?")
		}
	default: // bool, string, bytes, nil
		return errors.New("e625 : cannot div ?, ?")
	}
	return nil
}

func (vu *Vunit) divs(oper0 *Vunit, oper1 *Vunit) error {
	switch oper0.Vtype {
	case 2: // int
		switch oper1.Vtype {
		case 2: // int
			if oper1.Vint == 0 {
				return errors.New("e626 : cannot divs int, int_0")
			} else if oper0.Vint > 0 && oper1.Vint > 0 {
				vu.Vtype = 2
				vu.Vint = oper0.Vint / oper1.Vint
			} else {
				vu.Vtype = 2
				vu.Vint = int(math.Floor(float64(oper0.Vint) / float64(oper1.Vint)))
			}
		case 3: // float
			if oper1.Vfloat == 0.0 {
				return errors.New("e627 : cannot divs int, float_0")
			} else {
				vu.Vtype = 3
				vu.Vfloat = math.Floor(float64(oper0.Vint) / oper1.Vfloat)
			}
		default:
			return errors.New("e628 : cannot divs int, ?")
		}
	case 3: // float
		switch oper1.Vtype {
		case 2: // int
			if oper1.Vint == 0 {
				return errors.New("e629 : cannot divs float, int_0")
			} else {
				vu.Vtype = 3
				vu.Vfloat = math.Floor(oper0.Vfloat / float64(oper1.Vint))
			}
		case 3: // float
			if oper1.Vfloat == 0.0 {
				return errors.New("e630 : cannot divs float, float_0")
			} else {
				vu.Vtype = 3
				vu.Vfloat = math.Floor(oper0.Vfloat / oper1.Vfloat)
			}
		default:
			return errors.New("e631 : cannot divs float, ?")
		}
	default: // bool, string, bytes, nil
		return errors.New("e632 : cannot divs ?, ?")
	}
	return nil
}

func (vu *Vunit) divr(oper0 *Vunit, oper1 *Vunit) error {
	switch oper0.Vtype {
	case 2: // int
		switch oper1.Vtype {
		case 2: // int
			if oper1.Vint == 0 {
				return errors.New("e633 : cannot divr int, int_0")
			} else if oper0.Vint > 0 && oper1.Vint > 0 {
				vu.Vtype = 2
				vu.Vint = oper0.Vint % oper1.Vint
			} else {
				vu.Vtype = 2
				vu.Vint = int(float64(oper0.Vint) - float64(oper1.Vint)*math.Floor(float64(oper0.Vint)/float64(oper1.Vint)))
			}
		case 3: // float
			if oper1.Vfloat == 0.0 {
				return errors.New("e634 : cannot divr int, float_0")
			} else {
				vu.Vtype = 3
				vu.Vfloat = float64(oper0.Vint) - oper1.Vfloat*math.Floor(float64(oper0.Vint)/oper1.Vfloat)
			}
		default:
			return errors.New("e635 : cannot divr int, ?")
		}
	case 3: // float
		switch oper1.Vtype {
		case 2: // int
			if oper1.Vint == 0 {
				return errors.New("e636 : cannot divr float, int_0")
			} else {
				vu.Vtype = 3
				vu.Vfloat = oper0.Vfloat - float64(oper1.Vint)*math.Floor(oper0.Vfloat/float64(oper1.Vint))
			}
		case 3: // float
			if oper1.Vfloat == 0.0 {
				return errors.New("e637 : cannot divr float, float_0")
			} else {
				vu.Vtype = 3
				vu.Vfloat = oper0.Vfloat - oper1.Vfloat*math.Floor(oper0.Vfloat/oper1.Vfloat)
			}
		default:
			return errors.New("e638 : cannot divr float, ?")
		}
	default: // bool, string, bytes, nil
		return errors.New("e639 : cannot divr ?, ?")
	}
	return nil
}

func (vu *Vunit) pow(oper0 *Vunit, oper1 *Vunit) error {
	switch oper0.Vtype {
	case 2: // int
		switch oper1.Vtype {
		case 2: // int
			if oper1.Vint >= 0 {
				vu.Vtype = 2
				vu.Vint = int(math.Pow(float64(oper0.Vint), float64(oper1.Vint)))
			} else {
				vu.Vtype = 3
				vu.Vfloat = math.Pow(float64(oper0.Vint), float64(oper1.Vint))
			}
		case 3: // float
			vu.Vtype = 3
			vu.Vfloat = math.Pow(float64(oper0.Vint), oper1.Vfloat)
		default:
			return errors.New("e640 : cannot pow int, ?")
		}
	case 3: // float
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 3
			vu.Vfloat = math.Pow(oper0.Vfloat, float64(oper1.Vint))
		case 3: // float
			vu.Vtype = 3
			vu.Vfloat = math.Pow(oper0.Vfloat, oper1.Vfloat)
		default:
			return errors.New("e641 : cannot pow float, ?")
		}
	default: // bool, string, bytes, nil
		return errors.New("e642 : cannot pow ?, ?")
	}
	return nil
}

func (vu *Vunit) eql(oper0 *Vunit, oper1 *Vunit) error {
	vu.Vtype = 1
	switch oper0.Vtype {
	case 1: // bool
		if oper1.Vtype == 1 { // bool
			vu.Vbool = oper0.Vbool == oper1.Vbool
		} else {
			vu.Vbool = false
		}
	case 2: // int
		switch oper1.Vtype {
		case 2: // int
			vu.Vbool = oper0.Vint == oper1.Vint
		case 3: // float
			vu.Vbool = float64(oper0.Vint) == oper1.Vfloat
		default:
			vu.Vbool = false
		}
	case 3: // float
		switch oper1.Vtype {
		case 2: // int
			vu.Vbool = oper0.Vfloat == float64(oper1.Vint)
		case 3: // float
			vu.Vbool = oper0.Vfloat == oper1.Vfloat
		default:
			vu.Vbool = false
		}
	case 4: // string
		if oper1.Vtype == 4 { // string
			vu.Vbool = oper0.Vstring == oper1.Vstring
		} else {
			vu.Vbool = false
		}
	case 5: // bytes
		if oper1.Vtype == 5 { // bytes
			vu.Vbool = bytes.Equal(oper0.Vbytes, oper1.Vbytes)
		} else {
			vu.Vbool = false
		}
	default: // nil
		vu.Vbool = oper1.Vtype == 0
	}
	return nil
}

func (vu *Vunit) sml(oper0 *Vunit, oper1 *Vunit) error {
	switch oper0.Vtype {
	case 2: // int
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 1
			vu.Vbool = oper0.Vint < oper1.Vint
		case 3: // float
			vu.Vtype = 1
			vu.Vbool = float64(oper0.Vint) < oper1.Vfloat
		default:
			return errors.New("e643 : cannot sml int, ?")
		}
	case 3: // float
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 1
			vu.Vbool = oper0.Vfloat < float64(oper1.Vint)
		case 3: // float
			vu.Vtype = 1
			vu.Vbool = oper0.Vfloat < oper1.Vfloat
		default:
			return errors.New("e644 : cannot sml float, ?")
		}
	case 4: // string
		if oper1.Vtype == 4 { // string
			vu.Vtype = 1
			vu.Vbool = oper0.Vstring < oper1.Vstring
		} else {
			return errors.New("e645 : cannot sml string, ?")
		}
	case 5: // bytes
		if oper1.Vtype == 5 { // bytes
			vu.Vtype = 1
			vu.Vbool = bytes.Compare(oper0.Vbytes, oper1.Vbytes) == -1
		} else {
			return errors.New("e646 : cannot sml bytes, ?")
		}
	default: // bool, nil
		return errors.New("e647 : cannot sml ?, ?")
	}
	return nil
}

func (vu *Vunit) smle(oper0 *Vunit, oper1 *Vunit) error {
	switch oper0.Vtype {
	case 2: // int
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 1
			vu.Vbool = oper0.Vint <= oper1.Vint
		case 3: // float
			vu.Vtype = 1
			vu.Vbool = float64(oper0.Vint) <= oper1.Vfloat
		default:
			return errors.New("e648 : cannot smle int, ?")
		}
	case 3: // float
		switch oper1.Vtype {
		case 2: // int
			vu.Vtype = 1
			vu.Vbool = oper0.Vfloat <= float64(oper1.Vint)
		case 3: // float
			vu.Vtype = 1
			vu.Vbool = oper0.Vfloat <= oper1.Vfloat
		default:
			return errors.New("e649 : cannot smle float, ?")
		}
	case 4: // string
		if oper1.Vtype == 4 { // string
			vu.Vtype = 1
			vu.Vbool = oper0.Vstring <= oper1.Vstring
		} else {
			return errors.New("e650 : cannot smle string, ?")
		}
	case 5: // bytes
		if oper1.Vtype == 5 { // bytes
			vu.Vtype = 1
			vu.Vbool = bytes.Compare(oper0.Vbytes, oper1.Vbytes) != 1
		} else {
			return errors.New("e651 : cannot smle bytes, ?")
		}
	default: // bool, nil
		return errors.New("e652 : cannot smle ?, ?")
	}
	return nil
}

// shortcut opcode : get 0 or 1 operand, set itself result

func (vu *Vunit) inc() error {
	switch vu.Vtype {
	case 2: // int
		vu.Vint++
	case 3: // float
		vu.Vfloat = vu.Vfloat + 1.0
	default:
		return errors.New("e653 : cannot add ?, 1")
	}
	return nil
}

func (vu *Vunit) dec() error {
	switch vu.Vtype {
	case 2: // int
		vu.Vint--
	case 3: // float
		vu.Vfloat = vu.Vfloat - 1.0
	default:
		return errors.New("e654 : cannot sub ?, 1")
	}
	return nil
}

func (vu *Vunit) shm() error {
	switch vu.Vtype {
	case 2: // int
		vu.Vint = vu.Vint * 2
	case 3: // float
		vu.Vfloat = vu.Vfloat * 2.0
	case 4: // string
		vu.Vstring = vu.Vstring + vu.Vstring
	case 5: // bytes
		vu.Vbytes = append(vu.Vbytes, vu.Vbytes...)
	default:
		return errors.New("e655 : cannot mul ?, 2")
	}
	return nil
}

func (vu *Vunit) shd() error {
	switch vu.Vtype {
	case 2: // int
		vu.Vtype = 3
		vu.Vfloat = float64(vu.Vint) / 2.0
	case 3: // float
		vu.Vfloat = vu.Vfloat / 2.0
	default:
		return errors.New("e656 : cannot div ?, 2")
	}
	return nil
}

func (vu *Vunit) addi(operand int) error {
	switch vu.Vtype {
	case 2: // int
		vu.Vint = vu.Vint + operand
	case 3: // float
		vu.Vfloat = vu.Vfloat + float64(operand)
	default:
		return errors.New("e657 : cannot add ?, int")
	}
	return nil
}

func (vu *Vunit) muli(operand int) error {
	switch vu.Vtype {
	case 2: // int
		vu.Vint = vu.Vint * operand
	case 3: // float
		vu.Vfloat = vu.Vfloat * float64(operand)
	case 4: // string
		if operand < 0 {
			return errors.New("e658 : cannot mul string, int_n")
		} else {
			vu.Vstring = strings.Repeat(vu.Vstring, operand)
		}
	case 5: // bytes
		if operand < 0 {
			return errors.New("e659 : cannot mul bytes, int_n")
		} else {
			vu.Vbytes = slices.Repeat(vu.Vbytes, operand)
		}
	default:
		return errors.New("e660 : cannot mul ?, int")
	}
	return nil
}

// memory structure

type register struct {
	pc Vunit // program counter
	sp Vunit // stack pointer
	ma Vunit // multi-use A
	mb Vunit // multi-use B
}

type readonly struct {
	op  []byte // opcode
	reg []byte // register
	i16 []int  // int16 operand
	i32 []int  // int32 operand
}

type readwrite struct {
	tmp Vunit   // temp variable
	stk []Vunit // stack memory
}

type loadinfo struct {
	info   string
	abif   int
	sign   []byte
	public string

	rodata []byte
	data   []byte
	text   []byte
}

// kscript virtual machine (part 7)
type KVM struct {
	CallMem []Vunit // outercall args
	SafeMem bool    // safe memory access
	RunOne  bool    // run only 1 cycle
	ErrHlt  bool    // halt if error
	ErrMsg  string  // runtime error msg
	MaxStk  int     // stack max size

	loader loadinfo
	reg    register
	rom    readonly
	ram    readwrite
}

// decode little-endian signed int 16 / 32
func (kvm *KVM) dec(data []byte) int {
	if len(data) == 4 {
		return int(int32(data[0]) | int32(data[1])<<8 | int32(data[2])<<16 | int32(data[3])<<24)
	} else if len(data) == 2 {
		return int(int16(data[0]) | int16(data[1])<<8)
	} else {
		return 0
	}
}

// info, abif, sign, public, rodata, data, text
func (kvm *KVM) readelf(path string) error {
	worker0 := ksc.Initksc()
	worker0.Path = path
	worker0.Predetect = true
	err := worker0.Readf()
	if err != nil {
		return fmt.Errorf("%s at KVM.readelf ksc", err)
	}
	if !kio.Bequal(worker0.Subtype, []byte("KELF")) {
		return errors.New("e700 : not kelf file")
	}
	if len(worker0.Chunkpos) < 4 {
		return errors.New("e701 : invalid chunk num")
	}

	f, _ := kio.Open(path, "r")
	data, _ := kio.Read(f, -1)
	f.Close()
	pt0 := data[worker0.Chunkpos[0]+8 : worker0.Chunkpos[0]+worker0.Chunksize[0]+8]
	pt1 := data[worker0.Chunkpos[1]+8 : worker0.Chunkpos[1]+worker0.Chunksize[1]+8]
	pt2 := data[worker0.Chunkpos[2]+8 : worker0.Chunkpos[2]+worker0.Chunksize[2]+8]
	pt3 := data[worker0.Chunkpos[3]+8 : worker0.Chunkpos[3]+worker0.Chunksize[3]+8]
	if !kio.Bequal(ksc.Crc32hash(pt0), worker0.Reserved[0:4]) {
		return errors.New("e702 : invalid header crc32")
	}

	worker1 := kdb.Initkdb()
	worker1.Read(string(pt0))
	info, _ := worker1.Get("info")
	abif, _ := worker1.Get("abi")
	sign, _ := worker1.Get("sign")
	public, _ := worker1.Get("public")
	kvm.loader.info = info.Dat6
	kvm.loader.abif = abif.Dat2
	kvm.loader.sign = sign.Dat5
	kvm.loader.public = public.Dat6
	kvm.loader.rodata = pt1
	kvm.loader.data = pt2
	kvm.loader.text = pt3
	return nil
}

// add data to ram
func (kvm *KVM) load_data(code []byte) error {
	pos := 0
	for pos < len(code) {
		var vu Vunit
		switch code[pos] {
		case 78: // none
			vu.Set(nil)
			pos = pos + 1
		case 66: // bool
			if code[pos+1] == 0 {
				vu.Set(true)
			} else {
				vu.Set(false)
			}
			pos = pos + 2
		case 73: // int
			temp := bytes.NewReader(code[pos+1 : pos+9])
			var tgt int64
			binary.Read(temp, binary.LittleEndian, &tgt)
			vu.Set(int(tgt))
			pos = pos + 9
		case 70: // float
			vu.Set(math.Float64frombits(binary.LittleEndian.Uint64(code[pos+1 : pos+9])))
			pos = pos + 9
		case 83: // string
			length := kvm.dec(code[pos+1 : pos+5])
			pos = pos + 5
			vu.Set(string(code[pos : pos+length]))
			pos = pos + length
		case 67: // bytes
			length := kvm.dec(code[pos+1 : pos+5])
			pos = pos + 5
			vu.Set(code[pos : pos+length])
			pos = pos + length
		default:
			return errors.New("e703 : decode fail DATA")
		}
		kvm.ram.stk = append(kvm.ram.stk, vu)
	}
	return nil
}

// add data to rom
func (kvm *KVM) load_text(code []byte) error {
	if len(code)%8 != 0 {
		return errors.New("e704 : invalid code length")
	}
	// reg Y +4 N +0, i16 non_n Y +2 N +0, i32 non_n Y +1 N +0
	opcond := map[byte]byte{0: 0, 1: 0,
		16: 2, 17: 7, 18: 2, 19: 1, 20: 1, 21: 6, 22: 6,
		32: 6, 33: 6, 34: 4, 35: 4, 36: 2, 37: 2,
		48: 0, 49: 0, 50: 0, 51: 0,
		64: 0, 65: 0, 66: 0,
		80: 0, 81: 0, 82: 0, 83: 0, 84: 0, 85: 0,
		96: 2, 97: 2, 98: 2, 99: 2,
		112: 0, 113: 0, 114: 4, 115: 3}
	codenum := len(code) / 8
	kvm.rom.op = make([]byte, codenum)
	kvm.rom.reg = make([]byte, codenum)
	kvm.rom.i16 = make([]int, codenum)
	kvm.rom.i32 = make([]int, codenum)

	// load, check validity
	for i := 0; i < codenum; i++ {
		temp := 8 * i
		opcode := code[temp]
		reg := code[temp+1]
		i16 := kvm.dec(code[temp+2 : temp+4])
		i32 := kvm.dec(code[temp+4 : temp+8])

		cond, ext := opcond[opcode]
		if !ext {
			return errors.New("e705 : invalid opcode")
		}
		if cond&0x04 == 4 && reg != 97 && reg != 98 {
			return errors.New("e706 : invalid register")
		}
		if cond&0x02 == 2 && i16 < 0 {
			return errors.New("e707 : invalid int16")
		}
		if cond&0x01 == 1 && i32 < 0 {
			return errors.New("e708 : invalid int32")
		}

		kvm.rom.op[i] = opcode
		kvm.rom.reg[i] = reg
		kvm.rom.i16[i] = i16
		kvm.rom.i32[i] = i32
	}
	return nil
}

// seg + addr -> pos
func (kvm *KVM) getpos(i16 int, i32 int) int {
	if i16 == 108 {
		return kvm.reg.sp.Vint + i32
	} else {
		return i32
	}
}

// for (i, r <- v) calc
func (kvm *KVM) forcalc(reg byte, i16 int, i32 int, iscond bool) (Vunit, error) {
	var out Vunit
	var idx int
	pos := kvm.getpos(i16, i32)
	if reg == 97 {
		idx = kvm.reg.ma.Vint
	} else {
		idx = kvm.reg.mb.Vint
	}
	if kvm.SafeMem && (pos < 0 || pos >= len(kvm.ram.stk)) {
		return out, errors.New("e709 : ram access fail")
	}

	kvm.ram.tmp = kvm.ram.stk[pos]
	switch kvm.ram.tmp.Vtype {
	case 2: // int
		if kvm.ram.tmp.Vint > 0 {
			if iscond {
				if idx < kvm.ram.tmp.Vint {
					out.Set(true)
				} else {
					out.Set(false)
				}
			} else {
				out.Set(idx)
			}
		} else if kvm.ram.tmp.Vint < 0 {
			if iscond {
				if -idx > kvm.ram.tmp.Vint {
					out.Set(true)
				} else {
					out.Set(false)
				}
			} else {
				out.Set(-idx)
			}
		} else {
			if iscond {
				out.Set(false)
			}
		}

	case 4: // string
		if iscond {
			if idx < len([]rune(kvm.ram.tmp.Vstring)) {
				out.Set(true)
			} else {
				out.Set(false)
			}
		} else {
			out.Set(string([]rune(kvm.ram.tmp.Vstring)[idx]))
		}

	case 5: // bytes
		if iscond {
			if idx < len(kvm.ram.tmp.Vbytes) {
				out.Set(true)
			} else {
				out.Set(false)
			}
		} else {
			out.Set([]byte{kvm.ram.tmp.Vbytes[idx]})
		}

	default:
		out.Set(false)
		return out, errors.New("e710 : invalid for type")
	}
	return out, nil
}

// fetch - decode - execute - interupt (0 normal, 1 hlt, 2 err, -1 c_err)
func (kvm *KVM) cycle() int {
	var op byte
	var reg byte
	var i16 int
	var i32 int
	for {
		// fetch
		if kvm.SafeMem && kvm.reg.pc.Vint >= len(kvm.rom.op) {
			kvm.ErrMsg = "e711 : rom access fail"
			return -1
		}
		op, reg, i16, i32 = kvm.rom.op[kvm.reg.pc.Vint], kvm.rom.reg[kvm.reg.pc.Vint], kvm.rom.i16[kvm.reg.pc.Vint], kvm.rom.i32[kvm.reg.pc.Vint]
		kvm.reg.pc.Vint++

		// decode & execute
		switch op >> 4 {
		case 0:
			if op&0x0f == 0 { // hlt
				return 1
			}

		case 1:
			switch op & 0x0f {
			case 0: // intr
				if kvm.SafeMem && i16 > len(kvm.ram.stk) {
					kvm.ErrMsg = "e712 : ram access fail"
					return -1
				}
				kvm.ram.tmp.Vint = len(kvm.ram.stk) - i16
				kvm.CallMem = kvm.ram.stk[kvm.ram.tmp.Vint:]
				kvm.ram.stk = kvm.ram.stk[:kvm.ram.tmp.Vint]
				return i32

			case 1: // call
				if kvm.SafeMem && len(kvm.ram.stk) > kvm.MaxStk-4 {
					kvm.ErrMsg = "e713 : stack overflow"
					return -1
				}
				kvm.ram.tmp.Vint = len(kvm.ram.stk)
				if reg == 97 {
					kvm.ram.stk = append(kvm.ram.stk, kvm.reg.ma)
				} else {
					kvm.ram.stk = append(kvm.ram.stk, kvm.reg.mb)
				}
				kvm.ram.stk = append(append(kvm.ram.stk, kvm.reg.pc, kvm.reg.sp), make([]Vunit, i16+1)...)
				kvm.reg.pc.Vint = i32
				kvm.reg.sp.Vint = kvm.ram.tmp.Vint

			case 2: // ret
				kvm.ram.tmp.Vint = kvm.reg.sp.Vint
				kvm.reg.pc.Vint = kvm.ram.stk[kvm.ram.tmp.Vint+1].Vint
				kvm.reg.sp.Vint = kvm.ram.stk[kvm.ram.tmp.Vint+2].Vint
				kvm.reg.ma = kvm.ram.stk[kvm.ram.tmp.Vint+3]
				if kvm.SafeMem && kvm.ram.tmp.Vint < i16 {
					kvm.ErrMsg = "e714 : ram access fail"
					return -1
				}
				kvm.ram.stk = kvm.ram.stk[:kvm.ram.tmp.Vint-i16]

			case 3: // jmp
				kvm.reg.pc.Vint = i32

			case 4: // jmpiff
				kvm.ram.tmp = kvm.ram.stk[len(kvm.ram.stk)-1]
				kvm.ram.stk = kvm.ram.stk[:len(kvm.ram.stk)-1]
				if kvm.ram.tmp.Vtype == 1 {
					if !kvm.ram.tmp.Vbool {
						kvm.reg.pc.Vint = i32
					}
				} else {
					kvm.reg.pc.Vint = i32
					if kvm.ErrHlt {
						kvm.ErrMsg = "e715 : invalid condition"
						return 2
					}
				}

			case 5: // forcond
				temp, err := kvm.forcalc(reg, i16, i32, true)
				kvm.ram.stk = append(kvm.ram.stk, temp)
				if kvm.ErrHlt && err != nil {
					kvm.ErrMsg = fmt.Sprint(err)
					return 2
				}

			case 6: // forset
				kvm.reg.ma, _ = kvm.forcalc(reg, i16, i32, false)
			}

		case 2:
			switch op & 0x0f {
			case 0: // load
				pos := kvm.getpos(i16, i32)
				if kvm.SafeMem && (pos < 0 || pos >= len(kvm.ram.stk)) {
					kvm.ErrMsg = "e716 : ram access fail"
					return -1
				}
				if reg == 97 {
					kvm.reg.ma = kvm.ram.stk[pos]
				} else {
					kvm.reg.mb = kvm.ram.stk[pos]
				}

			case 1: // store
				pos := kvm.getpos(i16, i32)
				if kvm.SafeMem && (pos < 0 || pos >= len(kvm.ram.stk)) {
					kvm.ErrMsg = "e717 : ram access fail"
					return -1
				}
				if kvm.SafeMem && i16 == 99 {
					kvm.ErrMsg = "e718 : cannot write to const"
					return -1
				}
				if reg == 97 {
					kvm.ram.stk[pos] = kvm.reg.ma
				} else {
					kvm.ram.stk[pos] = kvm.reg.mb
				}

			case 2: // push
				if reg == 97 {
					kvm.ram.stk = append(kvm.ram.stk, kvm.reg.ma)
				} else {
					kvm.ram.stk = append(kvm.ram.stk, kvm.reg.mb)
				}

			case 3: // pop
				kvm.ram.tmp = kvm.ram.stk[len(kvm.ram.stk)-1]
				kvm.ram.stk = kvm.ram.stk[:len(kvm.ram.stk)-1]
				if reg == 97 {
					kvm.reg.ma = kvm.ram.tmp
				} else {
					kvm.reg.mb = kvm.ram.tmp
				}

			case 4: // pushset
				pos := kvm.getpos(i16, i32)
				if kvm.SafeMem && (pos < 0 || pos >= len(kvm.ram.stk)) {
					kvm.ErrMsg = "e719 : ram access fail"
					return -1
				}
				kvm.ram.stk = append(kvm.ram.stk, kvm.ram.stk[pos])

			case 5: // popset
				pos := kvm.getpos(i16, i32)
				if kvm.SafeMem && (pos < 0 || pos >= len(kvm.ram.stk)) {
					kvm.ErrMsg = "e720 : ram access fail"
					return -1
				}
				if kvm.SafeMem && i16 == 99 {
					kvm.ErrMsg = "e721 : cannot write to const"
					return -1
				}
				kvm.ram.stk[pos] = kvm.ram.stk[len(kvm.ram.stk)-1]
				kvm.ram.stk = kvm.ram.stk[:len(kvm.ram.stk)-1]
			}

		case 3:
			var temp Vunit
			var err error
			pos := len(kvm.ram.stk)
			switch op & 0x0f {
			case 0: // add
				err = temp.add(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
			case 1: // sub
				err = temp.sub(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
			case 2: // mul
				err = temp.mul(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
			case 3: // div
				err = temp.div(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
			}
			kvm.ram.stk[pos-2] = temp
			kvm.ram.stk = kvm.ram.stk[:pos-1]
			if kvm.ErrHlt && err != nil {
				kvm.ErrMsg = fmt.Sprint(err)
				return 2
			}

		case 4:
			var temp Vunit
			var err error
			pos := len(kvm.ram.stk)
			switch op & 0x0f {
			case 0: // divs
				err = temp.divs(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
			case 1: // divr
				err = temp.divr(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
			case 2: // pow
				err = temp.pow(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
			}
			kvm.ram.stk[pos-2] = temp
			kvm.ram.stk = kvm.ram.stk[:pos-1]
			if kvm.ErrHlt && err != nil {
				kvm.ErrMsg = fmt.Sprint(err)
				return 2
			}

		case 5:
			var temp Vunit
			var err error
			pos := len(kvm.ram.stk)
			switch op & 0x0f {
			case 0: // eql
				err = temp.eql(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
			case 1: // eqln
				err = temp.eql(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
				temp.Vbool = !temp.Vbool
			case 2: // sml
				err = temp.sml(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
			case 3: // grt
				err = temp.sml(&kvm.ram.stk[pos-1], &kvm.ram.stk[pos-2])
			case 4: // smle
				err = temp.smle(&kvm.ram.stk[pos-2], &kvm.ram.stk[pos-1])
			case 5: // grte
				err = temp.smle(&kvm.ram.stk[pos-1], &kvm.ram.stk[pos-2])
			}
			kvm.ram.stk[pos-2] = temp
			kvm.ram.stk = kvm.ram.stk[:pos-1]
			if kvm.ErrHlt && err != nil {
				kvm.ErrMsg = fmt.Sprint(err)
				return 2
			}

		case 6:
			var err error
			pos := kvm.getpos(i16, i32)
			if kvm.SafeMem && (pos < 0 || pos >= len(kvm.ram.stk)) {
				kvm.ErrMsg = "e722 : ram access fail"
				return -1
			}
			if kvm.SafeMem && i16 == 99 {
				kvm.ErrMsg = "e723 : cannot write to const"
				return -1
			}
			switch op & 0x0f {
			case 0: // inc
				err = kvm.ram.stk[pos].inc()
			case 1: // dec
				err = kvm.ram.stk[pos].dec()
			case 2: // shm
				err = kvm.ram.stk[pos].shm()
			case 3: // shd
				err = kvm.ram.stk[pos].shd()
			}
			if kvm.ErrHlt && err != nil {
				kvm.ErrMsg = fmt.Sprint(err)
				return 2
			}

		case 7:
			var err error
			switch op & 0x0f {
			case 0: // addi
				err = kvm.ram.stk[len(kvm.ram.stk)-1].addi(i32)
				if kvm.ErrHlt && err != nil {
					kvm.ErrMsg = fmt.Sprint(err)
					return 2
				}

			case 1: // muli
				err = kvm.ram.stk[len(kvm.ram.stk)-1].muli(i32)
				if kvm.ErrHlt && err != nil {
					kvm.ErrMsg = fmt.Sprint(err)
					return 2
				}

			case 2: // addr
				if reg == 97 {
					kvm.ram.stk = append(kvm.ram.stk, kvm.reg.ma)
				} else {
					kvm.ram.stk = append(kvm.ram.stk, kvm.reg.mb)
				}
				err = kvm.ram.stk[len(kvm.ram.stk)-1].addi(i32)
				if kvm.ErrHlt && err != nil {
					kvm.ErrMsg = fmt.Sprint(err)
					return 2
				}

			case 3: // jmpi
				kvm.ram.tmp.Vbool = false
				switch i16 {
				case 1:
					err = kvm.ram.tmp.eql(&kvm.reg.ma, &kvm.reg.mb)
				case 2:
					err = kvm.ram.tmp.eql(&kvm.reg.ma, &kvm.reg.mb)
					kvm.ram.tmp.Vbool = !kvm.ram.tmp.Vbool
				case 3:
					err = kvm.ram.tmp.sml(&kvm.reg.ma, &kvm.reg.mb)
				case 4:
					err = kvm.ram.tmp.sml(&kvm.reg.mb, &kvm.reg.ma)
				case 5:
					err = kvm.ram.tmp.smle(&kvm.reg.ma, &kvm.reg.mb)
				case 6:
					err = kvm.ram.tmp.smle(&kvm.reg.mb, &kvm.reg.ma)
				default:
					kvm.ErrMsg = "e724 : decode fail JMPI"
					return -1
				}
				if !kvm.ram.tmp.Vbool {
					kvm.reg.pc.Vint = i32
				}
				if kvm.ErrHlt && err != nil {
					kvm.ErrMsg = fmt.Sprint(err)
					return 2
				}
			}
		}

		// interupt
		if kvm.RunOne {
			return 0
		}
	}
}

// init runtime system
func (kvm *KVM) Init() {
	kvm.CallMem = nil
	kvm.RunOne = false
	kvm.SafeMem = true
	kvm.ErrHlt = true
	kvm.ErrMsg = ""
	kvm.MaxStk = 16777216

	kvm.reg.pc.Set(0)
	kvm.reg.sp.Set(0)
	kvm.reg.ma.Set(nil)
	kvm.reg.mb.Set(nil)

	kvm.rom.op = nil
	kvm.rom.reg = nil
	kvm.rom.i16 = nil
	kvm.rom.i32 = nil

	kvm.ram.stk = make([]Vunit, 0, 1024)
	kvm.ram.tmp.Set(nil)

	kvm.loader.info = ""
	kvm.loader.abif = -1
	kvm.loader.sign = nil
	kvm.loader.public = ""
	kvm.loader.rodata = nil
	kvm.loader.data = nil
	kvm.loader.text = nil
}

// view kelf file, returns info / abif / public
func (kvm *KVM) View(path string) (string, int, string, error) {
	err := kvm.readelf(path)
	return kvm.loader.info, kvm.loader.abif, kvm.loader.public, err
}

// load kelf file (should done View first), check sign
func (kvm *KVM) Load(sign bool) error {
	if kvm.loader.text == nil {
		return errors.New("e725 : should done kvm.View() first")
	}
	if sign && kvm.loader.public != "" {
		hworker := sha3.New512()
		hworker.Write(kvm.loader.rodata)
		hworker.Write(kvm.loader.data)
		hworker.Write(kvm.loader.text)
		hvalue := hworker.Sum(nil)
		sus, err := ksign.Verify(kvm.loader.public, kvm.loader.sign, hvalue)
		if !sus || err != nil {
			return errors.New("e726 : ksign verify fail")
		}
	}
	err := kvm.load_data(kvm.loader.rodata)
	if err != nil {
		return fmt.Errorf("%s at KVM.Load rodata", err)
	}
	err = kvm.load_data(kvm.loader.data)
	if err != nil {
		return fmt.Errorf("%s at KVM.Load data", err)
	}
	err = kvm.load_text(kvm.loader.text)
	if err != nil {
		return fmt.Errorf("%s at KVM.Load text", err)
	}
	return nil
}

// run code, exit with intr code
func (kvm *KVM) Run() (retv int) {
	defer func() {
		if err := recover(); err != nil {
			kvm.ErrMsg = fmt.Sprintf("%s at KVM.Run exit", err)
			retv = -1
		}
	}()
	kvm.CallMem = nil
	kvm.ErrMsg = ""
	retv = kvm.cycle()
	return retv
}

// set ma with outercall return
func (kvm *KVM) SetRet(ma *Vunit) {
	kvm.reg.ma = *ma
}

// test io support (input, print, read, write, time, sleep)
func TestIO(mode int, v []Vunit) *Vunit {
	var out Vunit
	switch mode {
	case 16: // input(v)
		out.Set(kio.Input(v[0].ToString()))
	case 17: // print(v)
		fmt.Print(v[0].ToString())
	case 18: // read(s, i)
		if v[0].Vtype == 4 && v[1].Vtype == 2 {
			f, err := kio.Open(v[0].Vstring, "r")
			if err == nil {
				defer f.Close()
				temp, _ := kio.Read(f, v[1].Vint)
				out.Set(temp)
			} else {
				out.Set([]byte{})
			}
		}
	case 19: // write(s, s|b)
		if v[0].Vtype == 4 && v[1].Vtype == 4 {
			f, _ := kio.Open(v[0].Vstring, "w")
			defer f.Close()
			kio.Write(f, []byte(v[1].Vstring))
		} else if v[0].Vtype == 4 && v[1].Vtype == 5 {
			f, _ := kio.Open(v[0].Vstring, "w")
			defer f.Close()
			kio.Write(f, v[1].Vbytes)
		}
	case 20: // time()
		out.Set(float64(time.Now().UnixMicro()) / 1000000)
	case 21: // sleep(i|f)
		if v[0].Vtype == 2 {
			time.Sleep(time.Second * time.Duration(v[0].Vint))
		} else if v[0].Vtype == 3 {
			time.Sleep(time.Microsecond * time.Duration(v[0].Vfloat*1000000))
		}
	}
	return &out
}
