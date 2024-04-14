// test632 : stdlib5.kdb

package kdb

import (
	"fmt"
	"stdlib5/kio"
	"strconv"
	"strings"
)

// pyvar, generation only
type kdbvar struct {
	Dat0 string
	Dat1 bool
	Dat2 int
	Dat3 float64
	Dat4 complex128
	Dat5 []byte
	Dat6 string
}

// kdbvar 7type setter
func Set(v interface{}) kdbvar {
	var out kdbvar
	switch x := interface{}(v).(type) {
	case nil:
		out.Dat0 = "None"
	case bool:
		out.Dat0 = "bool"
		out.Dat1 = x
	case int:
		out.Dat0 = "int"
		out.Dat2 = x
	case float64:
		out.Dat0 = "float"
		out.Dat3 = x
	case complex128:
		out.Dat0 = "complex"
		out.Dat4 = x
	case []byte:
		out.Dat0 = "bytes"
		out.Dat5 = x
	case string:
		out.Dat0 = "str"
		out.Dat6 = x
	default:
		out.Dat0 = "None"
	}
	return out
}

// for lite setting files
type toolbox struct {
	Name    map[string]int // (int index)[str name]
	Tp      []byte         // type int[]
	Ptr     []int          // pointer int[]
	Fmem    []float64      // mem float[]
	Cmem    []complex128   // mem complex[]
	Bmem    [][]byte       // mem bytes[]
	working []string       // current names list
}

// Init and return toolbox
func Initkdb() toolbox {
	var out toolbox
	out.Name = make(map[string]int)
	return out
}

// parse string
func (t *toolbox) Read(raw string) error {
	// pass id setter nonstr str sharp end
	t.working = make([]string, 0)
	order := strings.Split(raw, "\n")
	for i, r := range order {
		order[i] = r + "\n"
	}

	for _, i := range order {
		status := "pass"
		var mem []string
		name := ""
		v := ""

		for _, j := range i {
			switch status {

			case "pass":
				if (j != ' ') && (j != '\r') {
					if (j == '\n') || (j == ';') {
						status = "pass"
						mem = make([]string, 0)
						name = ""
						v = ""
					} else if j == '=' {
						return fmt.Errorf("invalid key : %s", i)
					} else {
						status = "id"
						mem = append(mem, string(j))
					}
				}

			case "id":
				if (j == '\n') || (j == ';') {
					status = "pass"
					mem = make([]string, 0)
					name = ""
					v = ""
				} else if (j != ' ') && (j != '\r') {
					if j == '=' {
						status = "setter"
						name = strings.Join(mem, "")
						mem = make([]string, 0)
					} else {
						mem = append(mem, string(j))
					}
				}

			case "setter":
				if (j != ' ') && (j != '\r') {
					if (j == '\n') || (j == ';') {
						return fmt.Errorf("invalid value : %s", i)
					} else if j == '"' {
						status = "str"
						mem = append(mem, string(j))
					} else {
						status = "nonstr"
						mem = append(mem, string(j))
					}
				}

			case "nonstr":
				if (j != ' ') && (j != '\r') {
					if j == '\n' {
						v = strings.Join(mem, "")
						t.add(name, v, "\n")
						status = "pass"
						mem = make([]string, 0)
						name = ""
						v = ""
					} else if j == ';' {
						v = strings.Join(mem, "")
						t.add(name, v, ";")
						status = "pass"
						mem = make([]string, 0)
						name = ""
						v = ""
					} else {
						mem = append(mem, string(j))
					}
				}

			case "str":
				if j == '\n' {
					return fmt.Errorf("invalid value : %s", i)
				} else if j == '#' {
					status = "sharp"
				} else if j == '"' {
					mem = append(mem, string(j))
					v = strings.Join(mem, "")
					status = "end"
				} else {
					mem = append(mem, string(j))
				}

			case "sharp":
				switch j {
				case '#':
					mem = append(mem, "#")
					status = "str"
				case 's':
					mem = append(mem, " ")
					status = "str"
				case 'n':
					mem = append(mem, "\n")
					status = "str"
				case '"':
					mem = append(mem, "\"")
					status = "str"
				default:
					return fmt.Errorf("invalid escaping : %s", i)
				}

			case "end":
				if j == '\n' {
					t.add(name, v, "\n")
					status = "pass"
					mem = make([]string, 0)
					name = ""
					v = ""
				} else if j == ';' {
					t.add(name, v, ";")
					status = "pass"
					mem = make([]string, 0)
					name = ""
					v = ""
				}
			}
		}
	}

	return nil
}

// return current DB str
func (t *toolbox) Write() (string, error) {
	out := make([]string, len(t.Tp))
	for name, num := range t.Name {
		tp := t.Tp[num]
		ptr := t.Ptr[num]
		var end string
		if tp/16 == 0 {
			end = "\n"
		} else {
			end = "; "
		}
		tp = tp % 16
		toadd := name + " = "

		var temp string
		var cnvr error
		switch tp {
		case 0:
			temp, cnvr = conv(nil)
		case 1:
			if ptr == 0 {
				temp, cnvr = conv(false)
			} else {
				temp, cnvr = conv(true)
			}
		case 2:
			temp, cnvr = conv(ptr)
		case 3:
			temp, cnvr = conv(t.Fmem[ptr])
		case 4:
			temp, cnvr = conv(t.Cmem[ptr])
		case 5:
			temp, cnvr = conv(t.Bmem[ptr])
		case 6:
			temp, cnvr = conv(string(t.Bmem[ptr]))
		default:
			return "", fmt.Errorf("invalid type : %d", tp)
		}

		if cnvr != nil {
			return "", cnvr
		}
		toadd = toadd + temp + end
		out[num] = toadd
	}
	stro := strings.Join(out, "")
	return stro, nil
}

// get value, [index, type, ptr] by name
func (t *toolbox) Get(name string) (kdbvar, []int) {
	num, cont := t.Name[strings.Replace(name, "/", ".", -1)]
	var out kdbvar
	if !cont {
		return out, nil
	}
	tp := int(t.Tp[num] % 16)
	ptr := t.Ptr[num]

	switch tp {
	case 0:
		out = Set(nil)
	case 1:
		if ptr == 0 {
			out = Set(false)
		} else {
			out = Set(true)
		}
	case 2:
		out = Set(ptr)
	case 3:
		out = Set(t.Fmem[ptr])
	case 4:
		out = Set(t.Cmem[ptr])
	case 5:
		out = Set(t.Bmem[ptr])
	case 6:
		out = Set(string(t.Bmem[ptr]))
	default:
		return out, nil
	}
	return out, []int{num, tp, ptr}
}

// revice data by name
func (t *toolbox) Fix(name string, v interface{}) error {
	num, extv := t.Name[strings.Replace(name, "/", ".", -1)]
	if !extv {
		return fmt.Errorf("KeyError : %s", name)
	}
	end := (t.Tp[num] / 16) * 16

	switch x := interface{}(v).(type) {
	case nil:
		t.Tp[num] = end + 0
		t.Ptr[num] = 0
	case bool:
		t.Tp[num] = end + 1
		if x {
			t.Ptr[num] = 1
		} else {
			t.Ptr[num] = 0
		}
	case int:
		t.Tp[num] = end + 2
		t.Ptr[num] = x
	case float64:
		t.Tp[num] = end + 3
		t.Ptr[num] = len(t.Fmem)
		t.Fmem = append(t.Fmem, x)
	case complex128:
		t.Tp[num] = end + 4
		t.Ptr[num] = len(t.Cmem)
		t.Cmem = append(t.Cmem, x)
	case []byte:
		t.Tp[num] = end + 5
		t.Ptr[num] = len(t.Bmem)
		t.Bmem = append(t.Bmem, x)
	case string:
		t.Tp[num] = end + 6
		t.Ptr[num] = len(t.Bmem)
		t.Bmem = append(t.Bmem, []byte(x))
	default:
		return fmt.Errorf("invalid type : %s", x)
	}
	return nil
}

// convert data -> kformat str
func conv(v interface{}) (string, error) {
	var out string
	switch x := interface{}(v).(type) {
	case nil:
		out = "None"
	case bool:
		if x {
			out = "True"
		} else {
			out = "False"
		}
	case int:
		out = fmt.Sprint(x)
	case float64:
		out = fmt.Sprintf("%f", x)
		for out[len(out)-1] == '0' {
			out = out[0 : len(out)-1]
		}
		if out[len(out)-1] == '.' {
			out = out + "0"
		}
	case complex128:
		st := fmt.Sprint(x)
		out = st[1 : len(st)-1]
	case []byte:
		out = "'" + kio.Bprint(x) + "'"
	case string:
		out = strings.Replace(x, "#", "##", -1)
		out = strings.Replace(out, " ", "#s", -1)
		out = strings.Replace(out, "\n", "#n", -1)
		out = strings.Replace(out, "\"", "#\"", -1)
		out = "\"" + out + "\""
	default:
		return "", fmt.Errorf("invalid type : %s", x)
	}
	return out, nil
}

// add DB by name, var, end
func (t *toolbox) add(name string, v string, end string) {
	nm := strings.Replace(name, "/", ".", -1)
	num := 0
	for nm[num] == '.' {
		num = num + 1
	}
	t.working = append(t.working[0:num], nm[num:])
	nm = strings.Join(t.working, ".")

	var tp byte
	if end == "\n" {
		tp = 0
	} else {
		tp = 16
	}
	tgtv := v

	if tgtv == "None" {
		t.Name[nm] = len(t.Tp)
		t.Tp = append(t.Tp, tp+0)
		t.Ptr = append(t.Ptr, 0)
	} else if tgtv == "True" {
		t.Name[nm] = len(t.Tp)
		t.Tp = append(t.Tp, tp+1)
		t.Ptr = append(t.Ptr, 1)
	} else if tgtv == "False" {
		t.Name[nm] = len(t.Tp)
		t.Tp = append(t.Tp, tp+1)
		t.Ptr = append(t.Ptr, 0)
	} else if tgtv[0] == '"' {
		t.Name[nm] = len(t.Tp)
		t.Tp = append(t.Tp, tp+6)
		t.Ptr = append(t.Ptr, len(t.Bmem))
		t.Bmem = append(t.Bmem, []byte(tgtv[1:len(tgtv)-1]))
	} else if tgtv[0] == '\'' {
		temp, _ := kio.Bread(tgtv[1 : len(tgtv)-1])
		t.Name[nm] = len(t.Tp)
		t.Tp = append(t.Tp, tp+5)
		t.Ptr = append(t.Ptr, len(t.Bmem))
		t.Bmem = append(t.Bmem, temp)
	} else if tgtv[len(tgtv)-1] == 'i' {
		t.Name[nm] = len(t.Tp)
		t.Tp = append(t.Tp, tp+4)
		t.Ptr = append(t.Ptr, len(t.Cmem))
		cv, _ := strconv.ParseComplex(tgtv, 128)
		t.Cmem = append(t.Cmem, cv)
	} else if strings.Contains(tgtv, ".") {
		t.Name[nm] = len(t.Tp)
		t.Tp = append(t.Tp, tp+3)
		t.Ptr = append(t.Ptr, len(t.Fmem))
		fv, _ := strconv.ParseFloat(tgtv, 64)
		t.Fmem = append(t.Fmem, fv)
	} else {
		t.Name[nm] = len(t.Tp)
		t.Tp = append(t.Tp, tp+2)
		iv, _ := strconv.Atoi(tgtv)
		t.Ptr = append(t.Ptr, iv)
	}
}

// import and set by input names & datas & ends (end : \n ;)
func (t *toolbox) Imp(names []string, datas []kdbvar, ends []string) error {
	for i := 0; i < len(ends); i++ {
		name := (names)[i]
		data := (datas)[i]
		var end byte

		num := len(t.Tp)
		t.Name[name] = num
		if ends[i] == "\n" {
			end = 0
		} else {
			end = 16
		}
		t.Tp = append(t.Tp, 0)
		t.Ptr = append(t.Ptr, 0)

		switch data.Dat0 {
		case "None":
			t.Tp[num] = end + 0
			t.Ptr[num] = 0
		case "bool":
			t.Tp[num] = end + 1
			if data.Dat1 {
				t.Ptr[num] = 1
			} else {
				t.Ptr[num] = 0
			}
		case "int":
			t.Tp[num] = end + 2
			t.Ptr[num] = data.Dat2
		case "float":
			t.Tp[num] = end + 3
			t.Ptr[num] = len(t.Fmem)
			t.Fmem = append(t.Fmem, data.Dat3)
		case "complex":
			t.Tp[num] = end + 4
			t.Ptr[num] = len(t.Cmem)
			t.Cmem = append(t.Cmem, data.Dat4)
		case "bytes":
			t.Tp[num] = end + 5
			t.Ptr[num] = len(t.Bmem)
			t.Bmem = append(t.Bmem, data.Dat5)
		case "str":
			t.Tp[num] = end + 6
			t.Ptr[num] = len(t.Bmem)
			t.Bmem = append(t.Bmem, []byte(data.Dat6))
		default:
			return fmt.Errorf("invalid type : %s", data.Dat0)
		}
	}
	return nil
}

// export to precise fullnames, datas, ends (end : \n ;)
func (t *toolbox) Exp() ([]string, []kdbvar, []string) {
	out0 := make([]string, len(t.Tp))
	out1 := make([]kdbvar, len(t.Tp))
	out2 := make([]string, len(t.Tp))
	for name, num := range t.Name {
		out0[num] = name
		tp := t.Tp[num]
		ptr := t.Ptr[num]
		if tp/16 == 0 {
			out2[num] = "\n"
		} else {
			out2[num] = ";"
		}
		tp = tp % 16

		switch tp {
		case 0:
			out1[num] = Set(nil)
		case 1:
			if ptr == 0 {
				out1[num] = Set(false)
			} else {
				out1[num] = Set(true)
			}
		case 2:
			out1[num] = Set(ptr)
		case 3:
			out1[num] = Set(t.Fmem[ptr])
		case 4:
			out1[num] = Set(t.Cmem[ptr])
		case 5:
			out1[num] = Set(t.Bmem[ptr])
		case 6:
			out1[num] = Set(string(t.Bmem[ptr]))
		default:
			return nil, nil, nil
		}
	}
	return out0, out1, out2
}
