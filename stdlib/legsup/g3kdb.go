// test657 : stdlib5.legsup gen3kdb

package legsup

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// gen3 kdb data node
type G3data struct {
	Next   *G3data // pointer to next data node (if list)
	Vtype  rune    // data type ('i', 'f', 's', 'n')
	IntV   int     // int data
	FloatV float64 // float data
	StrV   string  // string data
}

// init data node
func (mnode *G3data) Init() {
	mnode.Next = nil
	mnode.Vtype = 'n'
	mnode.IntV = 0
	mnode.FloatV = 0.0
	mnode.StrV = ""
}

// append data node(s) to G3data node
func (mnode *G3data) Append(tgt *G3data) {
	temp := mnode
	for temp.Next != nil {
		temp = temp.Next
	}
	temp.Next = tgt
}

// get length of G3data node
func (mnode *G3data) Length() int {
	temp := mnode
	count := 1
	for temp.Next != nil {
		temp = temp.Next
		count = count + 1
	}
	return count
}

// find G3data node by position
func (mnode *G3data) Locate(pos int) *G3data {
	temp := mnode
	for i := 0; i < pos; i++ {
		temp = temp.Next
	}
	return temp
}

// print series of data separated by comma
func (mnode *G3data) Print(zipstr bool, zipexp bool) string {
	temp := make([]string, 0)
	current := mnode
	for current != nil {
		switch current.Vtype {
		case 'i':
			temp = append(temp, fmt.Sprintf("%d", current.IntV))
		case 'f':
			temp = append(temp, fmt.Sprintf("%.4f", current.FloatV))
		case 's':
			ts := current.StrV
			ts = strings.Replace(ts, "#", "##", -1)
			ts = strings.Replace(ts, "\"", "#\"", -1)
			if zipstr {
				ts = strings.Replace(ts, " ", "#s", -1)
				ts = strings.Replace(ts, "\n", "#n", -1)
			}
			temp = append(temp, fmt.Sprintf("\"%s\"", ts))
		}
		current = current.Next
	}
	if zipexp {
		return strings.Join(temp, ",")
	} else {
		return strings.Join(temp, ", ")
	}
}

// gen3 kdb obj/data node
type G3node struct {
	Name  string   // sign name
	Data  *G3data  // content (nil if objmode)
	Child []G3node // child nodes
}

// read str, make internal tree
func (onode *G3node) Read(frag []string, trait []rune) error {
	// trait : 'g' gramer, 'i' int, 'f' float, 's' string
	pos := make([]int, 0)
	for i, r := range frag {
		if trait[i] == 'g' {
			switch r {
			case "[":
				pos = append(pos, i)
			case "]":
				pos = append(pos, i)
			case "{":
				pos = append(pos, i)
			case "}":
				pos = append(pos, i)
			}
		}
	}
	if len(pos)%4 != 0 {
		return errors.New("InvalidInput")
	}

	name := frag[pos[0]+1]
	if strings.Contains(name, "#") {
		return errors.New("comment") // name including "#" should be ignored
	} else {
		onode.Name = name
	}

	if len(pos) == 4 {
		var node0 G3data
		node0.Init()
		var err error
		for i := pos[2] + 1; i < pos[3]; i++ {
			var node1 G3data
			node1.Init()
			switch trait[i] {
			case 'i':
				node1.IntV, err = strconv.Atoi(frag[i])
				if err == nil {
					node1.Vtype = 'i'
				} else {
					return err
				}
			case 'f':
				node1.FloatV, err = strconv.ParseFloat(frag[i], 64)
				if err == nil {
					node1.Vtype = 'f'
				} else {
					return err
				}
			case 's':
				node1.StrV = frag[i]
				node1.Vtype = 's'
			}
			node0.Append(&node1)
		}

		onode.Data = node0.Next
		onode.Child = nil
		node0.Next = nil

	} else {
		current := 2
		stpos := 0
		endpos := 0
		count := 0
		onode.Data = nil
		onode.Child = make([]G3node, 0)

		for current < len(pos)-1 {
			if frag[pos[current]] == "[" {
				if count == 0 {
					stpos = pos[current]
				}
				count = count + 1
			} else if frag[pos[current]] == "}" {
				count = count - 1
				if count == 0 {
					endpos = pos[current]
					var newnode G3node
					err := newnode.Read(frag[stpos:endpos+1], trait[stpos:endpos+1])
					if err == nil {
						onode.Child = append(onode.Child, newnode)
					} else if err.Error() != "comment" {
						return err
					}
				}
			}
			current = current + 1
		}
	}

	return nil
}

// make output string with space indent
func (onode *G3node) Write(indent int, zipstr bool, zipexp bool) string {
	if onode.Data == nil {
		var temp string
		if zipexp {
			temp = fmt.Sprintf("[%s]{", onode.Name)
		} else {
			temp = fmt.Sprintf("%s[%s] {\n", strings.Repeat(" ", indent), onode.Name)
		}
		for _, r := range onode.Child {
			if zipexp {
				temp = temp + r.Write(0, zipstr, zipexp)
			} else {
				temp = temp + r.Write(indent+4, zipstr, zipexp) + "\n"
			}
		}
		if zipexp {
			temp = temp + "}"
		} else {
			temp = temp + fmt.Sprintf("%s}", strings.Repeat(" ", indent))
		}
		return temp

	} else {
		if zipexp {
			return fmt.Sprintf("[%s]{%s}", onode.Name, onode.Data.Print(zipstr, zipexp))
		} else {
			return fmt.Sprintf("%s[%s] {%s}", strings.Repeat(" ", indent), onode.Name, onode.Data.Print(zipstr, zipexp))
		}
	}
}

// find lower node by name token (returns nil if not exists)
func (onode *G3node) Locate(name string) *G3node {
	for _, r := range onode.Child {
		if r.Name == name {
			return &r
		}
	}
	return nil
}

// revise data
func (onode *G3node) Revise(tgt *G3data) error {
	if tgt.Vtype == 'n' || onode.Data == nil {
		return errors.New("InvalidDatatype")
	}
	onode.Data = tgt
	return nil
}

// gen3 kdb database
type G3kdb struct {
	Zipstr bool    // shorten string expression
	Zipexp bool    // shorten grammer expression
	Node   *G3node // first node
}

// read raw str, make internal parse tree
func (db *G3kdb) Read(raw string) error {
	frag := make([]string, 0)
	trait := make([]rune, 0)
	mem := ""
	status := "background" // "background", "name", "content", "nonstr", "str", "expstr"

	for _, r := range raw {
		switch status {
		case "name":
			switch r {
			case ']':
				frag = append(frag, mem)
				trait = append(trait, 's')
				frag = append(frag, "]")
				trait = append(trait, 'g')
				status = "background"
				mem = ""
			default:
				mem = mem + string(r)
			}
		case "content":
			switch r {
			case '[':
				frag = append(frag, "[")
				trait = append(trait, 'g')
				status = "name"
				mem = ""
			case '}':
				frag = append(frag, "}")
				trait = append(trait, 'g')
				status = "background"
				mem = ""
			case '"':
				status = "str"
				mem = ""
			case ',':
				mem = ""
			case ' ':
				mem = ""
			case '\n':
				mem = ""
			default:
				status = "nonstr"
				mem = string(r)
			}
		case "nonstr":
			switch r {
			case ',':
				if strings.Contains(mem, ".") {
					frag = append(frag, mem)
					trait = append(trait, 'f')
				} else {
					frag = append(frag, mem)
					trait = append(trait, 'i')
				}
				status = "content"
				mem = ""
			case '}':
				if strings.Contains(mem, ".") {
					frag = append(frag, mem)
					trait = append(trait, 'f')
				} else {
					frag = append(frag, mem)
					trait = append(trait, 'i')
				}
				frag = append(frag, "}")
				trait = append(trait, 'g')
				status = "background"
				mem = ""
			default:
				mem = mem + string(r)
			}
		case "str":
			switch r {
			case '"':
				frag = append(frag, mem)
				trait = append(trait, 's')
				status = "content"
				mem = ""
			case '#':
				status = "expstr"
			default:
				mem = mem + string(r)
			}
		case "expstr":
			switch r {
			case '#':
				mem = mem + "#"
			case '"':
				mem = mem + "\""
			case 's':
				mem = mem + " "
			case 'n':
				mem = mem + "\n"
			}
			status = "str"
		default:
			switch r {
			case '[':
				frag = append(frag, "[")
				trait = append(trait, 'g')
				status = "name"
				mem = ""
			case '{':
				frag = append(frag, "{")
				trait = append(trait, 'g')
				status = "content"
				mem = ""
			case '}':
				frag = append(frag, "}")
				trait = append(trait, 'g')
				status = "background"
				mem = ""
			}
		}
	}

	var temp G3node
	err := temp.Read(frag, trait)
	db.Node = &temp
	return err
}

// make output string
func (db *G3kdb) Write() string {
	return db.Node.Write(0, db.Zipstr, db.Zipexp)
}

// find node by fullname, returns nil if not exists
func (db *G3kdb) Locate(name string) *G3node {
	pos := strings.Split(name, "#")
	current := db.Node
	if current.Name != pos[0] {
		return nil
	}

	for _, r := range pos[1:] {
		temp := true
		count := 0
		for temp && count < len(current.Child) {
			if current.Child[count].Name == r {
				current = &current.Child[count]
				temp = false
			}
			count = count + 1
		}
		if temp {
			return nil
		}
	}
	return current
}
