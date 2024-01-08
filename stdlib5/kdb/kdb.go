package kdb

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// test557 : kdb5

// pyvar 유사체, 생성만 가능
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
		out.Dat0 = "nah"
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
		out.Dat0 = "nah"
	}
	return out
}

// for lite setting files
type toolbox struct {
	Name map[string]int
	Tp   []byte
	Ptr  []int
	Fmem []float64
	Cmem []complex128
	Bmem [][]byte
}

// Init and return toolbox
func Init() toolbox {
	var out toolbox
	out.Name = make(map[string]int)
	return out
}

// str parsing
func (t *toolbox) Readstr(raw *string) {
	// pass id setter nonstr str sharp end
	var current []string
	order := strings.Split(*raw, "\n")
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
						panic(fmt.Sprintf("invalid key : %s", i))
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
						panic(fmt.Sprintf("invalid value : %s", i))
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
						current = *t.add(&name, &v, "\n", &current)
						status = "pass"
						mem = make([]string, 0)
						name = ""
						v = ""
					} else if j == ';' {
						v = strings.Join(mem, "")
						current = *t.add(&name, &v, ";", &current)
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
					panic(fmt.Sprintf("invalid value : %s", i))
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
					panic(fmt.Sprintf("invalid escaping : %s", i))
				}

			case "end":
				if j == '\n' {
					current = *t.add(&name, &v, "\n", &current)
					status = "pass"
					mem = make([]string, 0)
					name = ""
					v = ""
				} else if j == ';' {
					current = *t.add(&name, &v, ";", &current)
					status = "pass"
					mem = make([]string, 0)
					name = ""
					v = ""
				}
			}
		}
	}
}

// txt file parsing
func (t *toolbox) Readfile(path string) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	sdat := string(dat)
	t.Readstr(&sdat)
}

// return current db str
func (t *toolbox) Writestr() *string {
	var out []string
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

		out = append(out, name, " = ")
		var temp string
		switch tp {
		case 0:
			temp = *kformat(nil)
		case 1:
			if ptr == 0 {
				temp = *kformat(false)
			} else {
				temp = *kformat(true)
			}
		case 2:
			temp = *kformat(ptr)
		case 3:
			temp = *kformat(t.Fmem[ptr])
		case 4:
			temp = *kformat(t.Cmem[ptr])
		case 5:
			temp = *kformat(t.Bmem[ptr])
		case 6:
			temp = *kformat(string(t.Bmem[ptr]))
		default:
			panic(fmt.Sprintf("invalid type : %d", tp))
		}
		out = append(out, temp, end)
	}
	stro := strings.Join(out, "")
	return &stro
}

// write current db to path
func (t *toolbox) Writefile(path string) {
	stro := *t.Writestr()
	err := ioutil.WriteFile(path, []byte(stro), 0644)
	if err != nil {
		panic(err)
	}
}

// current db str sorted by type & ptr
func (t *toolbox) Writestrs() *string {
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
		switch tp {
		case 0:
			temp = *kformat(nil)
		case 1:
			if ptr == 0 {
				temp = *kformat(false)
			} else {
				temp = *kformat(true)
			}
		case 2:
			temp = *kformat(ptr)
		case 3:
			temp = *kformat(t.Fmem[ptr])
		case 4:
			temp = *kformat(t.Cmem[ptr])
		case 5:
			temp = *kformat(t.Bmem[ptr])
		case 6:
			temp = *kformat(string(t.Bmem[ptr]))
		default:
			panic(fmt.Sprintf("invalid type : %d", tp))
		}
		toadd = toadd + temp + end
		out[num] = toadd
	}
	stro := strings.Join(out, "")
	return &stro
}

// write current db to path sorted by type & ptr
func (t *toolbox) Writefiles(path string) {
	stro := *t.Writestrs()
	err := ioutil.WriteFile(path, []byte(stro), 0644)
	if err != nil {
		panic(err)
	}
}

// find index, type, ptr by name (tp : 0~6, *str -> []int)
func (t *toolbox) Getpara(name *string) []int {
	nm := strings.Replace(*name, "/", ".", -1)
	num := t.Name[nm]
	tp := int(t.Tp[num] % 16)
	ptr := t.Ptr[num]
	out := []int{num, tp, ptr}
	return out
}

// find value by tp, ptr (int, int -> *kdbvar)
func (t *toolbox) Getvalue(tp int, ptr int) *kdbvar {
	var out kdbvar
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
		panic(fmt.Sprintf("invalid type : %d", tp))
	}
	return &out
}

// find value by name (*str -> *kdbvar)
func (t *toolbox) Getdata(name *string) *kdbvar {
	temp := t.Getpara(name)
	return t.Getvalue(temp[1], temp[2])
}

// revice data by name
func (t *toolbox) Fixdata(name *string, v interface{}) {
	nm := strings.Replace(*name, "/", ".", -1)
	num := t.Name[nm]
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
		panic(fmt.Sprintf("invalid type : %s", x))
	}
}

// infunc kformat (7type 입력받아 포매팅된 str 반환)
func kformat(v interface{}) *string {
	var out string
	switch x := interface{}(v).(type) {
	case nil:
		out = "nah"
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
		temp := make([]string, len(x))
		for i, r := range x {
			if r > 15 {
				temp[i] = fmt.Sprintf("%x", r)
			} else {
				temp[i] = "0" + fmt.Sprintf("%x", r)
			}
		}
		out = "'" + strings.Join(temp, "") + "'"
	case string:
		out = strings.Replace(x, "#", "##", -1)
		out = strings.Replace(out, " ", "#s", -1)
		out = strings.Replace(out, "\n", "#n", -1)
		out = strings.Replace(out, "\"", "#\"", -1)
		out = "\"" + out + "\""
	default:
		panic(fmt.Sprintf("invalid type : %s", x))
	}
	return &out
}

// infunc add (자신에게 값 추가)
func (t *toolbox) add(name *string, v *string, end string, current *[]string) *[]string {
	nm := strings.Replace(*name, "/", ".", -1)
	num := 0
	for nm[num] == '.' {
		num = num + 1
	}
	cur := *current
	cur = cur[0:num]
	cur = append(cur, nm[num:])
	nm = strings.Join(cur, ".")

	var tp byte
	if end == "\n" {
		tp = 0
	} else {
		tp = 16
	}
	tgtv := *v

	if tgtv == "nah" {
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
		tgtv = strings.ToLower(tgtv[1 : len(tgtv)-1])
		temp, _ := hex.DecodeString(tgtv)
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

	return &cur
}

// import and set by input names & datas & ends (end는 \n 또는 ;)
func (t *toolbox) Imp(names *[]string, datas *[]kdbvar, ends []string) {
	for i := 0; i < len(ends); i++ {
		name := (*names)[i]
		data := (*datas)[i]
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
		case "nah":
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
			panic(fmt.Sprintf("invalid type : %s", data.Dat0))
		}
	}
}

// export to precise fullnames *[]string & datas *[]kdbvar & ends []string (end는 \n 또는 ;)
func (t *toolbox) Exp() (*[]string, *[]kdbvar, []string) {
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
			panic(fmt.Sprintf("invalid type : %d", tp))
		}
	}
	return &out0, &out1, out2
}
