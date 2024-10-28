// test690 : stdlib5.kscript compiler

package kscript

// go get "golang.org/x/crypto/sha3"

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
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"
)

// syntax parser (part 1)
type Parser struct {
	Type_Allocator []string // memory allocator
	Type_Control   []string // control syntax
	Type_Operator  []string // operator syntax
	Type_Function  []string // basic function

	Sign_Bracket  []string // bracket with pair
	Sign_Blank    []string // whitespace characters
	Sign_Linefeed []string // newline characters
	Sign_Comma    []string // comma characters
	Sign_Comment  []string // comment characters

	Result_String []string // string values
	Result_Type   []string // type of words
}

// init parser with basic settings
func (ps *Parser) Init() {
	ps.Type_Allocator = []string{"=", "<-"}
	ps.Type_Control = []string{"def", "return", "if", "else", "while", "for"}
	ps.Type_Operator = []string{"+", "-", "*", "/", "//", "%", "**", ">", ">=", "<", "<=", "==", "!="}
	ps.Type_Function = []string{"test.input", "test.print", "test.read", "test.write", "test.time", "test.sleep"}

	ps.Sign_Bracket = []string{"(", ")", "{", "}", "[", "]"}
	ps.Sign_Blank = []string{" ", "\t", "\r"}
	ps.Sign_Linefeed = []string{"\n", ";"}
	ps.Sign_Comma = []string{",", "|"}
	ps.Sign_Comment = []string{"#", "~"}
}

// split raw code into tokens
func (ps *Parser) Split(raw string) (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("%s at Parser.Split exit", ferr)
		}
	}()
	raw = raw + "\n"
	status := "default"         // default syntax comma bytes string escape comment
	buffer := make([]string, 0) // storing characters
	ps.Result_String = make([]string, 0)
	ps.Result_Type = make([]string, 0)

	for _, r := range raw {
		l := string(r)
		switch status {
		case "syntax": // functions, controls, other literals
			if slices.Contains(ps.Sign_Bracket, l) { // bracket
				ps.Result_String = append(ps.Result_String, strings.Join(buffer, ""))
				buffer = make([]string, 0)
				ps.Result_Type = append(ps.Result_Type, "syntax") // push buffer
				if slices.Index(ps.Sign_Bracket, l)%2 == 0 {
					ps.Result_String = append(ps.Result_String, "(")
				} else {
					ps.Result_String = append(ps.Result_String, ")")
				}
				ps.Result_Type = append(ps.Result_Type, "bracket") // push bracket
				status = "default"
			} else if slices.Contains(ps.Sign_Blank, l) { // whitespace
				ps.Result_String = append(ps.Result_String, strings.Join(buffer, ""))
				buffer = make([]string, 0)
				ps.Result_Type = append(ps.Result_Type, "syntax") // push buffer
				status = "default"
			} else if slices.Contains(ps.Sign_Linefeed, l) { // newline
				ps.Result_String = append(ps.Result_String, strings.Join(buffer, ""))
				buffer = make([]string, 0)
				ps.Result_Type = append(ps.Result_Type, "syntax") // push buffer
				ps.Result_String = append(ps.Result_String, "\n")
				ps.Result_Type = append(ps.Result_Type, "newline") // push newline
				status = "default"
			} else if slices.Contains(ps.Sign_Comma, l) { // comma
				ps.Result_String = append(ps.Result_String, strings.Join(buffer, ""))
				buffer = make([]string, 0)
				ps.Result_Type = append(ps.Result_Type, "syntax") // push buffer
				ps.Result_String = append(ps.Result_String, ",")
				ps.Result_Type = append(ps.Result_Type, "comma") // push comma
				status = "default"
			} else { // continous syntax
				buffer = append(buffer, l)
			}

		case "bytes":
			if l == "'" { // end sign
				ps.Result_String = append(ps.Result_String, strings.Join(buffer, ""))
				buffer = make([]string, 0)
				ps.Result_Type = append(ps.Result_Type, "bytes") // push buffer
				status = "default"
			} else if !slices.Contains(ps.Sign_Blank, l) && !slices.Contains(ps.Sign_Linefeed, l) { // continous bytes
				buffer = append(buffer, l)
			}

		case "string":
			if l == "\"" { // end sign
				ps.Result_String = append(ps.Result_String, strings.Join(buffer, ""))
				buffer = make([]string, 0)
				ps.Result_Type = append(ps.Result_Type, "string") // push buffer
				status = "default"
			} else if l == "#" { // escape sign
				status = "escape"
			} else { // continous string
				buffer = append(buffer, l)
			}

		case "escape":
			if l == "#" {
				buffer = append(buffer, "#")
			} else if l == "\"" {
				buffer = append(buffer, "\"")
			} else if l == "n" {
				buffer = append(buffer, "\n")
			} else if l == "s" {
				buffer = append(buffer, " ")
			} else {
				return errors.New("e101 : wrong string escape")
			}
			status = "string"

		case "comment":
			if slices.Contains(ps.Sign_Linefeed, l) { // end sign
				ps.Result_String = append(ps.Result_String, "\n")
				ps.Result_Type = append(ps.Result_Type, "newline") // push newline
				status = "default"
			}

		default:
			if slices.Contains(ps.Sign_Bracket, l) { // bracket
				if slices.Index(ps.Sign_Bracket, l)%2 == 0 {
					ps.Result_String = append(ps.Result_String, "(")
				} else {
					ps.Result_String = append(ps.Result_String, ")")
				}
				ps.Result_Type = append(ps.Result_Type, "bracket") // push bracket
				status = "default"
			} else if slices.Contains(ps.Sign_Blank, l) { // whitespace
				status = "default"
			} else if slices.Contains(ps.Sign_Linefeed, l) { // newline
				ps.Result_String = append(ps.Result_String, "\n")
				ps.Result_Type = append(ps.Result_Type, "newline") // push newline
				status = "default"
			} else if slices.Contains(ps.Sign_Comma, l) { // comma
				ps.Result_String = append(ps.Result_String, ",")
				ps.Result_Type = append(ps.Result_Type, "comma") // push comma
				status = "default"
			} else if slices.Contains(ps.Sign_Comment, l) { // start of comment
				status = "comment"
			} else if l == "'" { // start of bytes
				buffer = make([]string, 0)
				status = "bytes"
			} else if l == "\"" { // start of string
				buffer = make([]string, 0)
				status = "string"
			} else { // start of syntax
				buffer = []string{l}
				status = "syntax"
			}
		}
	}

	count := 0 // bracket count
	for i, r := range ps.Result_Type {
		if r == "bracket" && ps.Result_String[i] == "(" {
			count = count + 1
		} else if r == "bracket" && ps.Result_String[i] == ")" {
			count = count - 1
		} else if count < 0 {
			return errors.New("e102 : wrong bracket sequence")
		}
	}
	if count == 0 {
		return nil
	} else {
		return errors.New("e103 : wrong bracket total")
	}
}

// examine int float name
func (ps *Parser) examine(token string) string {
	signs := []string{"+", "-"}
	numbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	isfloat := false
	for i, r := range token {
		l := string(r)
		switch {
		case l == ".":
			if isfloat {
				return "name"
			} else {
				isfloat = true
			}
		case slices.Contains(signs, l):
			if i != 0 {
				return "name"
			}
		case slices.Contains(numbers, l):
		default:
			return "name"
		}
	}
	if isfloat {
		return "float"
	} else {
		return "int"
	}
}

// parse result, set type (bracket newline allocator operator control basefunc none bool int float string bytes name)
func (ps *Parser) Parse() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("%s at Parser.Parse exit", ferr)
		}
	}()
	for i, r := range ps.Result_String {
		if ps.Result_Type[i] == "syntax" {
			switch {
			case slices.Contains(ps.Type_Allocator, r):
				ps.Result_Type[i] = "allocator"
			case slices.Contains(ps.Type_Operator, r):
				ps.Result_Type[i] = "operator"
			case slices.Contains(ps.Type_Control, r):
				ps.Result_Type[i] = "control"
			case slices.Contains(ps.Type_Function, r):
				ps.Result_Type[i] = "basefunc"
			case r == "None":
				ps.Result_Type[i] = "none"
			case r == "True" || r == "False":
				ps.Result_Type[i] = "bool"
			default:
				ps.Result_Type[i] = ps.examine(r)
			}
		}
	}
	return nil
}

// token : one word, expression token : (word...), token[] : whole program
func (ps *Parser) Structify() (result []Token, err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("%s at Parser.Structify exit", ferr)
		}
	}()
	result = make([]Token, 0)

	pos := 0
	for pos < len(ps.Result_Type) {
		var temp Token
		if ps.Result_Type[pos] == "bracket" {
			mempos := pos
			count := 1
			for count != 0 {
				pos = pos + 1
				if ps.Result_Type[pos] == "bracket" && ps.Result_String[pos] == "(" {
					count = count + 1
				} else if ps.Result_Type[pos] == "bracket" && ps.Result_String[pos] == ")" {
					count = count - 1
				}
			}
			err := temp.Read(ps.Result_String[mempos:pos+1], ps.Result_Type[mempos:pos+1])
			if err != nil {
				return result, err
			}
		} else {
			temp.Value = ps.Result_String[pos]
			temp.Vtype = ps.Result_Type[pos]
		}
		result = append(result, temp)
		pos = pos + 1
	}
	return result, nil
}

// get entire code, generate AST (functions, main flow)
func (ps *Parser) GenAST(token []Token) (functions []AST, program AST, err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("%s at Parser.GenAST exit", ferr)
		}
	}()
	functions = make([]AST, 0)
	program.Value = "flow"
	program.Vtype = "control"
	linecount := 0

	tstack := make([]Token, 0)
	for i, r := range token {
		if r.Vtype == "newline" { // flush or ignore
			linecount = linecount + 1
			if len(tstack) == 0 { // newline after newline
				continue
			}
			if i+1 < len(token) && token[i+1].Vtype == "expression" { // non K&R style bracket
				continue
			}
			if i+1 < len(token) && token[i+1].Vtype == "control" && token[i+1].Value == "else" { // if-else statement
				continue
			}

			// need flush
			var temp AST
			var err error
			if len(tstack) > 2 && tstack[0].Vtype == "control" && tstack[0].Value == "def" { // function define
				err = temp.ReadFunc(tstack)
				functions = append(functions, temp)
			} else if len(tstack) > 1 && tstack[0].Vtype == "control" && tstack[0].Value == "return" { // return statement
				err = errors.New("e104 : return statement at main code")
				program.Sub = append(program.Sub, temp)
			} else if len(tstack) > 2 && tstack[0].Vtype == "name" && tstack[1].Vtype == "allocator" { // assign flow
				err = temp.ReadAssign(tstack)
				program.Sub = append(program.Sub, temp)
			} else if len(tstack) > 1 && tstack[0].Vtype == "name" && tstack[1].Vtype == "expression" { // innercall flow
				err = temp.ReadCall(tstack)
				program.Sub = append(program.Sub, temp)
			} else if len(tstack) > 1 && tstack[0].Vtype == "basefunc" && tstack[1].Vtype == "expression" { // outercall flow
				err = temp.ReadCall(tstack)
				program.Sub = append(program.Sub, temp)
			} else { // control flow
				err = temp.ReadControl(tstack)
				program.Sub = append(program.Sub, temp)
			}

			if err == nil { // count line
				for _, r := range tstack {
					linecount = linecount + r.Count()
				}
				tstack = make([]Token, 0)
			} else {
				return functions, program, fmt.Errorf("%s at Parser.GenAST Line %d", err, linecount)
			}

		} else { // add token
			tstack = append(tstack, r)
		}
	}
	return functions, program, nil
}

// token without bracket (part 2)
type Token struct {
	Value string  // string formed value
	Vtype string  // expression newline allocator operator control basefunc none bool int float string bytes name
	Lower []Token // lower tokens
}

// read tokens with form (a b ...)
func (tk *Token) Read(value []string, vtype []string) error {
	if len(value) != len(vtype) || len(value) < 2 {
		return errors.New("e201 : wrong token length")
	}
	tk.Value = ""
	tk.Vtype = "expression"
	tk.Lower = make([]Token, 0)

	pos := 1
	for pos < len(vtype)-1 {
		var temp Token
		if vtype[pos] == "bracket" {
			mempos := pos
			count := 1
			for count != 0 {
				pos = pos + 1
				if vtype[pos] == "bracket" && value[pos] == "(" {
					count = count + 1
				} else if vtype[pos] == "bracket" && value[pos] == ")" {
					count = count - 1
				}
			}
			err := temp.Read(value[mempos:pos+1], vtype[mempos:pos+1])
			if err != nil {
				return err
			}
		} else {
			temp.Value = value[pos]
			temp.Vtype = vtype[pos]
		}
		tk.Lower = append(tk.Lower, temp)
		pos = pos + 1
	}
	return nil
}

// write token in debug mode
func (tk *Token) Write(indent int) string {
	out := fmt.Sprintf("%s %s (%s)\n", strings.Repeat("____", indent), tk.Value, tk.Vtype)
	for _, r := range tk.Lower {
		out = out + r.Write(indent+1)
	}
	return out
}

// count number of newline in token
func (tk *Token) Count() int {
	out := 0
	if tk.Vtype == "newline" {
		out = 1
	}
	for _, r := range tk.Lower {
		out = out + r.Count()
	}
	return out
}

// AST - literal, allocator, func call, bicalc, structures (part 3)
type AST struct {
	Value string // value, name
	Vtype string // literal(6) name bicalc innercall outercall assign function control
	Sub   []AST  // limited size subtree
}

// write ast in debug mode
func (ast *AST) Write(indent int) string {
	out := fmt.Sprintf("%s__ %s (%s)\n", strings.Repeat("| ", indent), ast.Value, ast.Vtype)
	for _, r := range ast.Sub {
		out = out + r.Write(indent+1)
	}
	return out
}

// read one expression / series of tokens -> one AST value (literal(6) name bicalc innercall outercall)
func (ast *AST) ReadExpr(token []Token) error {
	pre := make([]AST, 0)                                                                                                                                    // convert tokens to pre-AST
	post := make([]AST, 0)                                                                                                                                   // temp postfix notation
	stack_oper := make([]AST, 0)                                                                                                                             // value/operator stack
	superiority := map[string]int{"+": 20, "-": 20, "*": 30, "/": 30, "//": 30, "%": 30, "**": 40, ">": 10, ">=": 10, "<": 10, "<=": 10, "==": 10, "!=": 10} // operator superiority
	if len(token) == 0 {
		return errors.New("e301 : empty expression")
	}

	// breakdown to value & operator
	pretype := "operator"
	pos := 0
	for pos < len(token) {
		switch token[pos].Vtype {
		case "none", "bool", "int", "float", "string", "bytes": // literal
			if pretype == "operator" {
				var temp AST
				temp.Value = token[pos].Value
				temp.Vtype = token[pos].Vtype
				pre = append(pre, temp)
			} else {
				return fmt.Errorf("e302 : literal after %s", pretype)
			}
			pretype = "literal"
			pos = pos + 1

		case "name": // name or innercall
			if pretype == "operator" {
				var temp AST
				if pos+1 < len(token) && token[pos+1].Vtype == "expression" {
					err := temp.ReadCall(token[pos : pos+2])
					pos = pos + 1
					if err != nil {
						return fmt.Errorf("%s at AST.ReadExpr innercall_pre", err)
					}
				} else {
					temp.Value = token[pos].Value
					temp.Vtype = token[pos].Vtype
				}
				pre = append(pre, temp)
			} else {
				return fmt.Errorf("e303 : name after %s", pretype)
			}
			pretype = "name"
			pos = pos + 1

		case "basefunc": // outercall
			if pretype == "operator" {
				var temp AST
				if pos+1 < len(token) && token[pos+1].Vtype == "expression" {
					err := temp.ReadCall(token[pos : pos+2])
					pos = pos + 1
					if err != nil {
						return fmt.Errorf("%s at AST.ReadExpr outercall_pre", err)
					}
				} else {
					return fmt.Errorf("e303 : wrong basefunc call %s", pretype)
				}
				pre = append(pre, temp)
			} else {
				return fmt.Errorf("e304 : basefunc after %s", pretype)
			}
			pretype = "basefunc"
			pos = pos + 1

		case "expression": // expression
			if pretype == "operator" {
				var temp AST
				temp.ReadExpr(token[pos].Lower)
				pre = append(pre, temp)
			} else {
				return fmt.Errorf("e305 : expression after %s", pretype)
			}
			pretype = "expression"
			pos = pos + 1

		case "operator": // operator
			if pretype == "operator" {
				return errors.New("e306 : operator after operator")
			} else {
				var temp AST
				temp.Value = token[pos].Value
				temp.Vtype = token[pos].Vtype
				pre = append(pre, temp)
			}
			pretype = "operator"
			pos = pos + 1

		default:
			return fmt.Errorf("e307 : invalid ReadExpr type %s", token[pos].Vtype)
		}
	}

	// convert infix to postfix
	for _, r := range pre {
		switch r.Vtype {
		case "none", "bool", "int", "float", "string", "bytes", "name", "bicalc", "innercall", "outercall":
			post = append(post, r)
		case "operator":
			for len(stack_oper) != 0 && superiority[stack_oper[len(stack_oper)-1].Value] >= superiority[r.Value] {
				post = append(post, stack_oper[len(stack_oper)-1])
				stack_oper = stack_oper[0 : len(stack_oper)-1]
			}
			stack_oper = append(stack_oper, r)
		}
	}
	for len(stack_oper) != 0 {
		post = append(post, stack_oper[len(stack_oper)-1])
		stack_oper = stack_oper[0 : len(stack_oper)-1]
	}

	// convert postfix to bi-tree
	for _, r := range post {
		switch r.Vtype {
		case "none", "bool", "int", "float", "string", "bytes", "name", "bicalc", "innercall", "outercall":
			stack_oper = append(stack_oper, r)
		case "operator":
			var temp AST
			temp.Value = r.Value
			temp.Vtype = "bicalc"
			temp.Sub = []AST{stack_oper[len(stack_oper)-2], stack_oper[len(stack_oper)-1]}
			stack_oper = append(stack_oper[0:len(stack_oper)-2], temp)
		}
	}
	ast.Value = stack_oper[0].Value
	ast.Vtype = stack_oper[0].Vtype
	ast.Sub = stack_oper[0].Sub
	return nil
}

// read function call tokens -> one AST value (innercall outercall)
func (ast *AST) ReadCall(token []Token) error {
	if len(token) != 2 {
		return fmt.Errorf("e308 : invalid token length %d", len(token))
	}
	if (token[0].Vtype != "name" && token[0].Vtype != "basefunc") || token[1].Vtype != "expression" {
		return fmt.Errorf("e309 : invalid token type %s %s", token[0].Vtype, token[1].Vtype)
	}
	if token[0].Vtype == "name" {
		ast.Vtype = "innercall"
	} else {
		ast.Vtype = "outercall"
	}
	ast.Value = token[0].Value
	ast.Sub = make([]AST, 0)

	// cut tokens by comma
	pos := 0
	for i, r := range token[1].Lower {
		if r.Vtype == "comma" {
			var temp AST
			err := temp.ReadExpr(token[1].Lower[pos:i])
			if err != nil {
				return fmt.Errorf("%s at AST.ReadCall args", err)
			}
			ast.Sub = append(ast.Sub, temp)
			pos = i + 1
		}
	}
	if pos != len(token[1].Lower) {
		var temp AST
		err := temp.ReadExpr(token[1].Lower[pos:])
		if err != nil {
			return fmt.Errorf("%s at AST.ReadCall args", err)
		}
		ast.Sub = append(ast.Sub, temp)
	}
	return nil
}

// read variable assign tokens -> one AST token (assign)
func (ast *AST) ReadAssign(token []Token) error {
	if len(token) < 3 {
		return fmt.Errorf("e310 : invalid token length %d", len(token))
	}
	if token[0].Vtype != "name" || token[1].Vtype != "allocator" {
		return fmt.Errorf("e311 : invalid token type %s %s", token[0].Vtype, token[1].Vtype)
	}
	ast.Value = token[0].Value
	ast.Vtype = "assign"
	ast.Sub = make([]AST, 0)

	pos := 2
	for pos < len(token) && token[pos].Vtype != "newline" {
		pos = pos + 1
	}
	var temp AST
	err := temp.ReadExpr(token[2:pos])
	ast.Sub = append(ast.Sub, temp)
	if err != nil {
		return fmt.Errorf("%s at AST.ReadAssign expr", err)
	} else {
		return nil
	}
}

// read function define tokens -> one AST token (function)
func (ast *AST) ReadFunc(token []Token) error {
	if len(token) != 4 {
		return fmt.Errorf("e312 : invalid token length %d", len(token))
	}
	if token[0].Vtype != "control" || token[1].Vtype != "name" || token[2].Vtype != "expression" || token[3].Vtype != "expression" {
		return fmt.Errorf("e313 : invalid token type %s %s %s %s", token[0].Vtype, token[1].Vtype, token[2].Vtype, token[3].Vtype)
	}
	ast.Value = token[1].Value
	ast.Vtype = "function"
	ast.Sub = make([]AST, 0)

	// cut tokens by comma
	pretype := "comma"
	for _, r := range token[2].Lower {
		switch r.Vtype {
		case "comma":
			if pretype == "comma" {
				return errors.New("e314 : comma after comma")
			} else {
				pretype = "comma"
			}
		case "name":
			if pretype == "name" {
				return errors.New("e315 : name after name")
			} else {
				var temp AST
				temp.Value = r.Value
				temp.Vtype = "name"
				ast.Sub = append(ast.Sub, temp)
				pretype = "name"
			}
		default:
			return fmt.Errorf("e316 : invalid argument type %s", r.Vtype)
		}
	}

	var temp AST
	err := temp.ReadFlow(token[3].Lower)
	ast.Sub = append(ast.Sub, temp)
	if err != nil {
		return fmt.Errorf("%s at AST.ReadFunc line", err)
	} else {
		return nil
	}
}

// read multi lines of tokens (no def inside) -> one AST token (control-flow)
func (ast *AST) ReadFlow(token []Token) error {
	if len(token) == 0 {
		return errors.New("e317 : invalid token length")
	}
	ast.Value = "flow"
	ast.Vtype = "control"
	ast.Sub = make([]AST, 0)
	var addtoken Token
	addtoken.Value = "\n"
	addtoken.Vtype = "newline"
	token = append(token, addtoken)

	tstack := make([]Token, 0)
	for i, r := range token {
		if r.Vtype == "newline" { // flush or ignore
			if len(tstack) == 0 { // newline after newline
				continue
			}
			if i+1 < len(token) && token[i+1].Vtype == "expression" { // non K&R style bracket
				continue
			}
			if i+1 < len(token) && token[i+1].Value == "else" && token[i+1].Vtype == "control" { // if-else statement
				continue
			}

			// need flush
			var temp AST
			var err error
			if len(tstack) > 2 && tstack[0].Vtype == "name" && tstack[1].Vtype == "allocator" { // assign flow
				err = temp.ReadAssign(tstack)
			} else if len(tstack) > 1 && tstack[0].Vtype == "name" && tstack[1].Vtype == "expression" { // innercall flow
				err = temp.ReadCall(tstack)
			} else if len(tstack) > 1 && tstack[0].Vtype == "basefunc" && tstack[1].Vtype == "expression" { // outercall flow
				err = temp.ReadCall(tstack)
			} else { // control flow
				err = temp.ReadControl(tstack)
			}
			if err != nil {
				return fmt.Errorf("%s at AST.ReadFlow line", err)
			}
			tstack = make([]Token, 0)
			ast.Sub = append(ast.Sub, temp)

		} else { // add token
			tstack = append(tstack, r)
		}
	}
	return nil
}

// read control tokens (if while for return) -> one AST token (control)
func (ast *AST) ReadControl(token []Token) error {
	if len(token) < 2 || token[0].Vtype != "control" {
		return fmt.Errorf("e318 : invalid token %d %s", len(token), token[0].Vtype)
	}
	ast.Vtype = "control"
	var err error

	// checking control type
	if len(token) == 3 && token[0].Value == "if" { // if statement
		ast.Value = "if"
		ast.Sub = make([]AST, 2)
		if token[1].Vtype != "expression" || token[2].Vtype != "expression" {
			return errors.New("e319 : invalid if statement")
		}
		err = ast.Sub[0].ReadExpr(token[1].Lower)
		if err != nil {
			return fmt.Errorf("%s at AST.ReadControl if-cond", err)
		}
		err = ast.Sub[1].ReadFlow(token[2].Lower)
		if err != nil {
			return fmt.Errorf("%s at AST.ReadControl if-flow", err)
		}

	} else if len(token) == 5 && token[0].Value == "if" && token[3].Value == "else" { // if else statement
		ast.Value = "if"
		ast.Sub = make([]AST, 3)
		if token[1].Vtype != "expression" || token[2].Vtype != "expression" || token[3].Vtype != "control" || token[4].Vtype != "expression" {
			return errors.New("e320 : invalid if-else statement")
		}
		err = ast.Sub[0].ReadExpr(token[1].Lower)
		if err != nil {
			return fmt.Errorf("%s at AST.ReadControl if-cond", err)
		}
		err = ast.Sub[1].ReadFlow(token[2].Lower)
		if err != nil {
			return fmt.Errorf("%s at AST.ReadControl if-flow", err)
		}
		err = ast.Sub[2].ReadFlow(token[4].Lower)
		if err != nil {
			return fmt.Errorf("%s at AST.ReadControl else-flow", err)
		}

	} else if len(token) == 3 && token[0].Value == "while" { // while statement
		ast.Value = "while"
		ast.Sub = make([]AST, 2)
		if token[1].Vtype != "expression" || token[2].Vtype != "expression" {
			return errors.New("e321 : invalid while statement")
		}
		err = ast.Sub[0].ReadExpr(token[1].Lower)
		if err != nil {
			return fmt.Errorf("%s at AST.ReadControl while-cond", err)
		}
		err = ast.Sub[1].ReadFlow(token[2].Lower)
		if err != nil {
			return fmt.Errorf("%s at AST.ReadControl while-flow", err)
		}

	} else if len(token) == 3 && token[0].Value == "for" { // for statement
		ast.Value = "for"
		ast.Sub = make([]AST, 4)
		if token[1].Vtype != "expression" || token[2].Vtype != "expression" {
			return errors.New("e322 : invalid for statement")
		}
		if len(token[1].Lower) != 5 || token[1].Lower[0].Vtype != "name" || token[1].Lower[1].Vtype != "comma" || token[1].Lower[2].Vtype != "name" || token[1].Lower[3].Vtype != "allocator" {
			return errors.New("e323 : invalid for-assign") // for (i, r <- v)
		}
		if !slices.Contains([]string{"int", "string", "bytes", "name"}, token[1].Lower[4].Vtype) { // v should be int string bytes name
			return fmt.Errorf("e324 : invalid for-assign type %s", token[1].Lower[4].Vtype)
		}

		ast.Sub[0].Value = token[1].Lower[0].Value
		ast.Sub[0].Vtype = token[1].Lower[0].Vtype
		ast.Sub[1].Value = token[1].Lower[2].Value
		ast.Sub[1].Vtype = token[1].Lower[2].Vtype
		ast.Sub[2].Value = token[1].Lower[4].Value
		ast.Sub[2].Vtype = token[1].Lower[4].Vtype
		err = ast.Sub[3].ReadFlow(token[2].Lower)
		if err != nil {
			return fmt.Errorf("%s at AST.ReadControl for-flow", err)
		}

	} else if token[0].Value == "return" { // return statement
		ast.Value = "return"
		ast.Sub = make([]AST, 1)
		err = ast.Sub[0].ReadExpr(token[1:])
		if err != nil {
			return fmt.Errorf("%s at AST.ReadControl return-expr", err)
		}

	} else {
		return fmt.Errorf("e325 : invalid control type %s", token[0].Value)
	}
	return nil
}

// AST -> CTree -> kasm (part 4.0 ~ 4.4)
type Compiler struct {
	OptConst bool // fold literal ! manual init !
	OptAsm   bool // use shortcode ! manual init !

	OuterNum   map[string]int // outercall interupt code (32+) ! manual init !
	OuterParms map[string]int // outercall [funcname]parms num ! manual init !
	InnerNum   map[string]int // innercall jmp num
	InnerParms map[string]int // innercall [funcname]parms num
	InnerLocal map[string]int // innercall [funcname]localvar num (include argnum)

	VarGlobal  map[string]int // global [varname]pos
	VarLocal   map[string]int // local [varname]pos ! need init !
	VarArgs    int            // function args num ! need init !
	CountLocal int            // function localvar counter ! need init !

	CountLabel   int                 // jmp label counter
	CountFor     int                 // fornum counter
	CountLiteral map[string]CLiteral // [type + strf]value, (none true false 0 1 0.0 s"" b'' ...)

	il_tree []CTree  // optimized CTree, [0]mainflow [1+]function
	il_main string   // kasm mainflow (joined with \n)
	il_func []string // kasm function codestr
}

// phase 1 - set literal
func (cp *Compiler) readconst(ct *CTree) {
	if ct.Vtype == "literal" {
		if ct.Ivalue.Vtype == "none" { // none
			ct.Ivalue.Vpos = 0
		} else { // not-none
			db, ext := cp.CountLiteral[ct.Ivalue.Vtype+ct.Value]
			if ext { // update ivalue
				ct.Ivalue.Vpos = db.Vpos
			} else { // need to push
				ct.Ivalue.Vpos = len(cp.CountLiteral)
				cp.CountLiteral[ct.Ivalue.Vtype+ct.Value] = *ct.Ivalue
			}
		}
	} else {
		for i := 0; i < len(ct.Sub); i++ {
			cp.readconst(&ct.Sub[i])
		}
	}
}

// phase 2 - read function
func (cp *Compiler) readfunc(ct *CTree) error {
	if ct.Vtype != "function" {
		return fmt.Errorf("e400 : reading non-function %s", ct.Vtype)
	}
	if len(ct.Sub) == 0 {
		return errors.New("e401 : no flow inside function")
	}
	if _, ext := cp.InnerNum[ct.Value]; ext {
		return fmt.Errorf("e402 : double define %s", ct.Value)
	}

	var temp CLiteral
	temp.Set("funcname", ct.Value, len(cp.CountLiteral))
	ct.Ivalue = &temp
	cp.CountLiteral["funcname"+ct.Value] = temp
	cp.InnerNum[ct.Value] = cp.CountLabel
	cp.CountLabel = cp.CountLabel + 1
	cp.InnerParms[ct.Value] = len(ct.Sub) - 1
	return nil
}

// phase 3 - read global variable
func (cp *Compiler) readglobal(ct *CTree) error {
	var err error
	err = nil
	switch ct.Vtype {
	case "function": // should not defined
		return errors.New("e403 : function define at mainflow")

	case "innercall": // using inner func
		if cp.InnerNum[ct.Value] == 0 {
			return fmt.Errorf("e404 : unknown function %s", ct.Value)
		}
		if cp.InnerParms[ct.Value] != len(ct.Sub) {
			return fmt.Errorf("e405 : invalid args number %d", len(ct.Sub))
		}
		for i := 0; i < len(ct.Sub); i++ {
			err = cp.readglobal(&ct.Sub[i])
			if err != nil {
				return err
			}
		}

	case "outercall": // using outer func
		if cp.OuterNum[ct.Value] == 0 {
			return fmt.Errorf("e406 : unknown function %s", ct.Value)
		}
		if cp.OuterParms[ct.Value] != len(ct.Sub) {
			return fmt.Errorf("e407 : invalid args number %d", len(ct.Sub))
		}
		for i := 0; i < len(ct.Sub); i++ {
			err = cp.readglobal(&ct.Sub[i])
			if err != nil {
				return err
			}
		}

	case "name": // using global var
		if cp.VarGlobal[ct.Value] == 0 {
			return fmt.Errorf("e408 : unknown variable %s", ct.Value)
		}

	case "assign": // assign global var
		err = cp.readglobal(&ct.Sub[0])
		if cp.VarGlobal[ct.Value] == 0 { // new global var
			cp.VarGlobal[ct.Value] = len(cp.CountLiteral) + len(cp.VarGlobal)
		}

	case "control": // control (if while for return flow)
		if ct.Value == "for" {
			if cp.VarGlobal[ct.Sub[0].Value] == 0 { // new global var (count)
				cp.VarGlobal[ct.Sub[0].Value] = len(cp.CountLiteral) + len(cp.VarGlobal)
			}
			if cp.VarGlobal[ct.Sub[1].Value] == 0 { // new global var (value)
				cp.VarGlobal[ct.Sub[1].Value] = len(cp.CountLiteral) + len(cp.VarGlobal)
			}
			err = cp.readglobal(&ct.Sub[2]) // tgt value
			if err != nil {
				return err
			}

			var temp CLiteral
			temp.Vtype = "fornum"
			temp.Vpos = cp.CountFor
			cp.VarGlobal[fmt.Sprintf("for count %d", cp.CountFor)] = len(cp.CountLiteral) + len(cp.VarGlobal)
			cp.VarGlobal[fmt.Sprintf("for value %d", cp.CountFor)] = len(cp.CountLiteral) + len(cp.VarGlobal)
			cp.CountFor = cp.CountFor + 1
			ct.Ivalue = &temp
			err = cp.readglobal(&ct.Sub[3])

		} else {
			for i := 0; i < len(ct.Sub); i++ {
				err = cp.readglobal(&ct.Sub[i])
				if err != nil {
					return err
				}
			}
		}

	default: // literal, bicalc
		for i := 0; i < len(ct.Sub); i++ {
			err = cp.readglobal(&ct.Sub[i])
			if err != nil {
				return err
			}
		}
	}
	return err
}

// phase 4 - read function, set parms/local
func (cp *Compiler) readlocal(ct *CTree) error {
	var err error
	err = nil
	switch ct.Vtype {
	case "function": // setting parms
		argnum := cp.InnerParms[ct.Value]
		for i := 0; i < argnum; i++ {
			cp.VarLocal[ct.Sub[i].Value] = i - argnum
		}
		err = cp.readlocal(&ct.Sub[argnum])

	case "innercall": // using inner func
		if cp.InnerNum[ct.Value] == 0 {
			return fmt.Errorf("e409 : unknown function %s", ct.Value)
		}
		if cp.InnerParms[ct.Value] != len(ct.Sub) {
			return fmt.Errorf("e410 : invalid args number %d", len(ct.Sub))
		}
		for i := 0; i < len(ct.Sub); i++ {
			err = cp.readlocal(&ct.Sub[i])
			if err != nil {
				return err
			}
		}

	case "outercall": // using outer func
		if cp.OuterNum[ct.Value] == 0 {
			return fmt.Errorf("e411 : unknown function %s", ct.Value)
		}
		if cp.OuterParms[ct.Value] != len(ct.Sub) {
			return fmt.Errorf("e412 : invalid args number %d", len(ct.Sub))
		}
		for i := 0; i < len(ct.Sub); i++ {
			err = cp.readlocal(&ct.Sub[i])
			if err != nil {
				return err
			}
		}

	case "name": // check local -> global
		if cp.VarLocal[ct.Value] == 0 && cp.VarGlobal[ct.Value] == 0 {
			return fmt.Errorf("e413 : unknown variable %s", ct.Value)
		}

	case "assign": // assign local -> global -> local
		for i := 0; i < len(ct.Sub); i++ {
			err = cp.readlocal(&ct.Sub[i])
			if err != nil {
				return err
			}
		}
		if cp.VarLocal[ct.Value] == 0 && cp.VarGlobal[ct.Value] == 0 { // new local var
			cp.VarLocal[ct.Value] = cp.CountLocal + 4
			cp.CountLocal = cp.CountLocal + 1
		}

	case "control": // control (if while for return flow)
		if ct.Value == "for" {
			if cp.VarLocal[ct.Sub[0].Value] == 0 && cp.VarGlobal[ct.Sub[0].Value] == 0 { // new local var (count)
				cp.VarLocal[ct.Sub[0].Value] = cp.CountLocal + 4
				cp.CountLocal = cp.CountLocal + 1
			}
			if cp.VarLocal[ct.Sub[1].Value] == 0 && cp.VarGlobal[ct.Sub[1].Value] == 0 { // new local var (value)
				cp.VarLocal[ct.Sub[1].Value] = cp.CountLocal + 4
				cp.CountLocal = cp.CountLocal + 1
			}
			err = cp.readlocal(&ct.Sub[2]) // tgt value
			if err != nil {
				return err
			}

			var temp CLiteral
			temp.Vtype = "fornum"
			temp.Vpos = cp.CountFor
			cp.VarLocal[fmt.Sprintf("for count %d", cp.CountFor)] = cp.CountLocal + 4
			cp.VarLocal[fmt.Sprintf("for value %d", cp.CountFor)] = cp.CountLocal + 5
			cp.CountLocal = cp.CountLocal + 2
			cp.CountFor = cp.CountFor + 1
			ct.Ivalue = &temp
			err = cp.readlocal(&ct.Sub[3])

		} else {
			for i := 0; i < len(ct.Sub); i++ {
				err = cp.readlocal(&ct.Sub[i])
				if err != nil {
					return err
				}
			}
		}

	default: // literal, bicalc
		for i := 0; i < len(ct.Sub); i++ {
			err = cp.readlocal(&ct.Sub[i])
			if err != nil {
				return err
			}
		}
	}
	return err
}

// phase 5 - gen kasm (main / func)
func (cp *Compiler) genkasm(ct *CTree) (asm string, err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("%s at Compiler.genkasm exit", ferr)
		}
	}()
	asm = ""
	err = nil

	switch ct.Vtype {
	case "literal":
		asm = asm + fmt.Sprintf("pushset const &%d\n", ct.Ivalue.Vpos)

	case "name":
		if cp.VarLocal[ct.Value] == 0 { // using global
			asm = asm + fmt.Sprintf("pushset global &%d\n", cp.VarGlobal[ct.Value])
		} else { // using local
			asm = asm + fmt.Sprintf("pushset local &%d\n", cp.VarLocal[ct.Value])
		}

	case "bicalc":
		var t0 string
		var t1 error
		if cp.OptAsm { // opt : +const *const (int32 const)
			if ct.Value == "+" { // F + C, C + F, ? + C, C + ?
				if ct.Sub[0].Vtype == "literal" && ct.Sub[0].Ivalue.Vtype == "int" && -2147483648 < ct.Sub[0].Ivalue.Ivalue.(int) && ct.Sub[0].Ivalue.Ivalue.(int) < 2147483648 {
					if ct.Sub[1].Vtype == "innercall" || ct.Sub[1].Vtype == "outercall" { // push from MA
						t0, t1 = cp.gencall(&ct.Sub[1], false)
						asm = asm + t0 + fmt.Sprintf("addr ma $%d\n", ct.Sub[0].Ivalue.Ivalue.(int))
					} else { // push operand 0
						t0, t1 = cp.genkasm(&ct.Sub[1])
						asm = asm + t0 + fmt.Sprintf("addi $%d\n", ct.Sub[0].Ivalue.Ivalue.(int))
					}
					return asm, t1

				} else if ct.Sub[1].Vtype == "literal" && ct.Sub[1].Ivalue.Vtype == "int" && -2147483648 < ct.Sub[1].Ivalue.Ivalue.(int) && ct.Sub[1].Ivalue.Ivalue.(int) < 2147483648 {
					if ct.Sub[0].Vtype == "innercall" || ct.Sub[0].Vtype == "outercall" { // push from MA
						t0, t1 = cp.gencall(&ct.Sub[0], false)
						asm = asm + t0 + fmt.Sprintf("addr ma $%d\n", ct.Sub[1].Ivalue.Ivalue.(int))
					} else { // push operand 0
						t0, t1 = cp.genkasm(&ct.Sub[0])
						asm = asm + t0 + fmt.Sprintf("addi $%d\n", ct.Sub[1].Ivalue.Ivalue.(int))
					}
					return asm, t1
				}

			} else if ct.Value == "-" { // F - C, ? - C
				if ct.Sub[1].Vtype == "literal" && ct.Sub[1].Ivalue.Vtype == "int" && -2147483648 < ct.Sub[1].Ivalue.Ivalue.(int) && ct.Sub[1].Ivalue.Ivalue.(int) < 2147483648 {
					if ct.Sub[0].Vtype == "innercall" || ct.Sub[0].Vtype == "outercall" { // push from MA
						t0, t1 = cp.gencall(&ct.Sub[0], false)
						asm = asm + t0 + fmt.Sprintf("addr ma $%d\n", -ct.Sub[1].Ivalue.Ivalue.(int))
					} else { // push operand 0
						t0, t1 = cp.genkasm(&ct.Sub[0])
						asm = asm + t0 + fmt.Sprintf("addi $%d\n", -ct.Sub[1].Ivalue.Ivalue.(int))
					}
					return asm, t1
				}

			} else if ct.Value == "*" { // ? * C, C * ?
				if ct.Sub[0].Vtype == "literal" && ct.Sub[0].Ivalue.Vtype == "int" && -2147483648 < ct.Sub[0].Ivalue.Ivalue.(int) && ct.Sub[0].Ivalue.Ivalue.(int) < 2147483648 {
					t0, t1 = cp.genkasm(&ct.Sub[1])
					asm = asm + t0 + fmt.Sprintf("muli $%d\n", ct.Sub[0].Ivalue.Ivalue.(int))
					return asm, t1

				} else if ct.Sub[1].Vtype == "literal" && ct.Sub[1].Ivalue.Vtype == "int" && -2147483648 < ct.Sub[1].Ivalue.Ivalue.(int) && ct.Sub[1].Ivalue.Ivalue.(int) < 2147483648 {
					t0, t1 = cp.genkasm(&ct.Sub[0])
					asm = asm + t0 + fmt.Sprintf("muli $%d\n", ct.Sub[1].Ivalue.Ivalue.(int))
					return asm, t1
				}
			}
		}

		t0, t1 = cp.genkasm(&ct.Sub[0])
		asm = asm + t0
		if t1 != nil {
			return asm, t1
		}
		t0, t1 = cp.genkasm(&ct.Sub[1])
		asm = asm + t0
		if t1 != nil {
			return asm, t1
		}

		switch ct.Value {
		case "+":
			asm = asm + "add\n"
		case "-":
			asm = asm + "sub\n"
		case "*":
			asm = asm + "mul\n"
		case "/":
			asm = asm + "div\n"
		case "//":
			asm = asm + "divs\n"
		case "%":
			asm = asm + "divr\n"
		case "**":
			asm = asm + "pow\n"
		case ">":
			asm = asm + "grt\n"
		case ">=":
			asm = asm + "grte\n"
		case "<":
			asm = asm + "sml\n"
		case "<=":
			asm = asm + "smle\n"
		case "==":
			asm = asm + "eql\n"
		case "!=":
			asm = asm + "eqln\n"
		default:
			err = fmt.Errorf("e414 : unsupported bicalc %s", ct.Value)
		}

	case "innercall", "outercall":
		return cp.gencall(ct, true)

	case "assign":
		var vasm string
		if cp.VarLocal[ct.Value] == 0 { // global v
			vasm = fmt.Sprintf("global &%d ;assign\n", cp.VarGlobal[ct.Value])
		} else { // local v
			vasm = fmt.Sprintf("local &%d ;assign\n", cp.VarLocal[ct.Value])
		}

		if cp.OptAsm {
			if ct.Sub[0].Vtype == "innercall" || ct.Sub[0].Vtype == "outercall" { // v = F()
				t0, t1 := cp.gencall(&ct.Sub[0], false)
				asm = asm + t0 + "store ma " + vasm
				return asm, t1

			} else if ct.Sub[0].Vtype == "literal" { // v = l
				asm = asm + fmt.Sprintf("load mb const &%d\n", ct.Sub[0].Ivalue.Vpos) + "store mb " + vasm
				return asm, nil

			} else if ct.Sub[0].Vtype == "name" { // v = u
				if cp.VarLocal[ct.Sub[0].Value] == 0 { // global u
					asm = asm + fmt.Sprintf("load mb global &%d\n", cp.VarGlobal[ct.Sub[0].Value])
				} else { // local u
					asm = asm + fmt.Sprintf("load mb local &%d\n", cp.VarLocal[ct.Sub[0].Value])
				}
				asm = asm + "store mb " + vasm
				return asm, nil

			} else if ct.Sub[0].Vtype == "bicalc" {
				if ct.Sub[0].Sub[0].Vtype == "name" && ct.Value == ct.Sub[0].Sub[0].Value && ct.Sub[0].Sub[1].Vtype == "literal" {
					if ct.Sub[0].Value == "+" && ct.Sub[0].Sub[1].Ivalue.Vtype == "int" && ct.Sub[0].Sub[1].Ivalue.Ivalue.(int) == 1 { // v = v + 1
						asm = asm + "inc " + vasm
						return asm, nil
					} else if ct.Sub[0].Value == "-" && ct.Sub[0].Sub[1].Ivalue.Vtype == "int" && ct.Sub[0].Sub[1].Ivalue.Ivalue.(int) == 1 { // v = v - 1
						asm = asm + "dec " + vasm
						return asm, nil
					} else if ct.Sub[0].Value == "*" && ct.Sub[0].Sub[1].Ivalue.Vtype == "int" && ct.Sub[0].Sub[1].Ivalue.Ivalue.(int) == 2 { // v = v * 2
						asm = asm + "shm " + vasm
						return asm, nil
					} else if ct.Sub[0].Value == "/" && ct.Sub[0].Sub[1].Ivalue.Vtype == "int" && ct.Sub[0].Sub[1].Ivalue.Ivalue.(int) == 2 { // v = v / 2
						asm = asm + "shd " + vasm
						return asm, nil
					}

				} else if ct.Sub[0].Sub[0].Vtype == "literal" && ct.Sub[0].Sub[1].Vtype == "name" && ct.Value == ct.Sub[0].Sub[1].Value {
					if ct.Sub[0].Value == "+" && ct.Sub[0].Sub[0].Ivalue.Vtype == "int" && ct.Sub[0].Sub[0].Ivalue.Ivalue.(int) == 1 { // v = 1 + v
						asm = asm + "inc " + vasm
						return asm, nil
					} else if ct.Sub[0].Value == "*" && ct.Sub[0].Sub[0].Ivalue.Vtype == "int" && ct.Sub[0].Sub[0].Ivalue.Ivalue.(int) == 2 { // v = 2 * v
						asm = asm + "shm " + vasm
						return asm, nil
					}
				}
			}
		}

		t0, t1 := cp.genkasm(&ct.Sub[0])
		asm = asm + t0 + "popset " + vasm
		err = t1

	case "function":
		return cp.genfunc(ct)

	case "control":
		return cp.gencontrol(ct)

	default:
		err = fmt.Errorf("e415 : invalid CTree %s", ct.Vtype)
	}
	return asm, err
}

// phase 5 - gen kasm (call)
func (cp *Compiler) gencall(ct *CTree, push bool) (string, error) {
	if ct.Vtype != "innercall" && ct.Vtype != "outercall" {
		return "", errors.New("e416 : calling non-function")
	}
	asm := ""
	for _, r := range ct.Sub { // push args
		t0, t1 := cp.genkasm(&r)
		asm = asm + t0
		if t1 != nil {
			return asm, t1
		}
	}
	if ct.Vtype == "innercall" { // innercall
		asm = asm + fmt.Sprintf("load ma const &%d\n", cp.CountLiteral["funcname"+ct.Value].Vpos)
		asm = asm + fmt.Sprintf("call ma $%d @%d ;innercall\n", cp.InnerLocal[ct.Value]-cp.InnerParms[ct.Value], cp.InnerNum[ct.Value])
	} else { // outercall
		asm = asm + fmt.Sprintf("intr $%d $%d ;outercall\n", cp.OuterParms[ct.Value], cp.OuterNum[ct.Value])
	}
	if push { // push ma to stack
		asm = asm + "push ma\n"
	}
	return asm, nil
}

// phase 5 - gen kasm (control)
func (cp *Compiler) gencontrol(ct *CTree) (string, error) {
	var err error
	err = nil
	asm := ""

	switch ct.Value {
	case "if":
		t0, t1, t2 := cp.gencond(&ct.Sub[0])
		if t2 != nil {
			return asm, t2
		}
		asm = asm + t0

		if len(ct.Sub) == 2 { // plain if
			t0, t2 = cp.genkasm(&ct.Sub[1])
			err = t2
			asm = asm + t0 + t1

		} else { // if-else
			t0, t2 = cp.genkasm(&ct.Sub[1])
			if t2 != nil {
				return asm, t2
			}
			asm = asm + t0
			t0, t2 = cp.genkasm(&ct.Sub[2])
			err = t2
			t1 = t1 + t0

			asm = asm + fmt.Sprintf("jmp @%d\n", cp.CountLabel)
			t1 = t1 + fmt.Sprintf("label @%d\n", cp.CountLabel)
			cp.CountLabel = cp.CountLabel + 1
			asm = asm + t1
		}

	case "while":
		t0, t1, t2 := cp.gencond(&ct.Sub[0])
		if t2 != nil {
			return asm, t2
		}
		t3, t2 := cp.genkasm(&ct.Sub[1])
		err = t2

		asm = asm + fmt.Sprintf("label @%d\n", cp.CountLabel) + t0 + t3 + fmt.Sprintf("jmp @%d\n", cp.CountLabel) + t1
		cp.CountLabel = cp.CountLabel + 1

	case "for":
		t0 := "" // for count (invisible)
		t1 := "" // for value (invisible)
		if cp.VarLocal[fmt.Sprintf("for count %d", ct.Ivalue.Vpos)] == 0 {
			t0 = fmt.Sprintf("global &%d", cp.VarGlobal[fmt.Sprintf("for count %d", ct.Ivalue.Vpos)])
			t1 = fmt.Sprintf("global &%d", cp.VarGlobal[fmt.Sprintf("for value %d", ct.Ivalue.Vpos)])
		} else {
			t0 = fmt.Sprintf("local &%d", cp.VarLocal[fmt.Sprintf("for count %d", ct.Ivalue.Vpos)])
			t1 = fmt.Sprintf("local &%d", cp.VarLocal[fmt.Sprintf("for value %d", ct.Ivalue.Vpos)])
		}
		t2 := "" // var i
		t3 := "" // var r
		t4 := "" // target (i, r <- target)
		if cp.VarLocal[ct.Sub[0].Value] == 0 {
			t2 = fmt.Sprintf("global &%d", cp.VarGlobal[ct.Sub[0].Value])
		} else {
			t2 = fmt.Sprintf("local &%d", cp.VarLocal[ct.Sub[0].Value])
		}
		if cp.VarLocal[ct.Sub[1].Value] == 0 {
			t3 = fmt.Sprintf("global &%d", cp.VarGlobal[ct.Sub[1].Value])
		} else {
			t3 = fmt.Sprintf("local &%d", cp.VarLocal[ct.Sub[1].Value])
		}
		if ct.Sub[2].Vtype == "literal" {
			t4 = fmt.Sprintf("const &%d", ct.Sub[2].Ivalue.Vpos)
		} else if cp.VarLocal[ct.Sub[2].Value] == 0 {
			t4 = fmt.Sprintf("global &%d", cp.VarGlobal[ct.Sub[2].Value])
		} else {
			t4 = fmt.Sprintf("local &%d", cp.VarLocal[ct.Sub[2].Value])
		}

		t5, t6 := cp.genkasm(&ct.Sub[3]) // worker code
		err = t6
		t7 := fmt.Sprintf("jmpiff @%d\n", cp.CountLabel)
		t8 := fmt.Sprintf("label @%d\n", cp.CountLabel)
		cp.CountLabel = cp.CountLabel + 1

		asm = asm + fmt.Sprintf("load mb const &3 ;for_load\nload ma %s\n", t4)
		asm = asm + fmt.Sprintf("store mb %s\nstore ma %s\n", t0, t1)
		asm = asm + fmt.Sprintf("label @%d\nforcond mb %s ;for_check\n", cp.CountLabel, t1)
		asm = asm + t7 + fmt.Sprintf("forset mb %s ;for_assign\nstore mb %s\nstore ma %s\n", t1, t2, t3) + t5

		if cp.OptAsm {
			asm = asm + fmt.Sprintf("inc %s ;for_add\nload mb %s\n", t0, t0)
		} else {
			asm = asm + fmt.Sprintf("pushset %s ;for_add\npushset const &4\nadd\npop mb\nstore mb %s\n", t0, t0)
		}
		asm = asm + fmt.Sprintf("jmp @%d\n", cp.CountLabel) + t8
		cp.CountLabel = cp.CountLabel + 1

	case "return":
		if cp.OptAsm {
			if ct.Sub[0].Vtype == "literal" { // return l
				asm = asm + fmt.Sprintf("load mb const &%d\n", ct.Sub[0].Ivalue.Vpos) + "store mb local &3\n" + fmt.Sprintf("ret $%d ;return\n", cp.VarArgs)
				return asm, err
			} else if ct.Sub[0].Vtype == "name" { // return v
				if cp.VarLocal[ct.Sub[0].Value] == 0 { // global v
					asm = asm + fmt.Sprintf("load mb global &%d\n", cp.VarGlobal[ct.Sub[0].Value])
				} else { // local v
					asm = asm + fmt.Sprintf("load mb local &%d\n", cp.VarLocal[ct.Sub[0].Value])
				}
				asm = asm + "store mb local &3\n" + fmt.Sprintf("ret $%d ;return\n", cp.VarArgs)
				return asm, err
			}
		}

		t0, t1 := cp.genkasm(&ct.Sub[0])
		asm = asm + t0 + "popset local &3\n" + fmt.Sprintf("ret $%d ;return\n", cp.VarArgs)
		err = t1

	case "flow":
		var t0 string
		var t1 error
		for _, r := range ct.Sub {
			if r.Vtype == "innercall" || r.Vtype == "outercall" {
				t0, t1 = cp.gencall(&r, false)
			} else {
				t0, t1 = cp.genkasm(&r)
			}
			asm = asm + t0
			if t1 != nil {
				return asm, t1
			}
		}

	default:
		err = fmt.Errorf("e417 : invalid CTree %s", ct.Vtype)
	}
	return asm, err
}

// phase 5 - gen kasm (condition)
func (cp *Compiler) gencond(ct *CTree) (string, string, error) {
	var err error
	err = nil
	asm0 := ""
	asm1 := ""
	opcode := map[string]int{"==": 1, "!=": 2, "<": 3, ">": 4, "<=": 5, ">=": 6}

	if cp.OptAsm && ct.Vtype == "bicalc" && opcode[ct.Value] != 0 {
		if ct.Sub[0].Vtype == "name" && ct.Sub[1].Vtype == "name" { // v == u
			if cp.VarLocal[ct.Sub[0].Value] == 0 { // global v
				asm0 = asm0 + fmt.Sprintf("load ma global &%d\n", cp.VarGlobal[ct.Sub[0].Value])
			} else { // local v
				asm0 = asm0 + fmt.Sprintf("load ma local &%d\n", cp.VarLocal[ct.Sub[0].Value])
			}
			if cp.VarLocal[ct.Sub[1].Value] == 0 { // global u
				asm0 = asm0 + fmt.Sprintf("load mb global &%d\n", cp.VarGlobal[ct.Sub[1].Value])
			} else { // local u
				asm0 = asm0 + fmt.Sprintf("load mb local &%d\n", cp.VarLocal[ct.Sub[1].Value])
			}
			asm0 = asm0 + fmt.Sprintf("jmpi $%d @%d ;condition\n", opcode[ct.Value], cp.CountLabel)
			asm1 = fmt.Sprintf("label @%d\n", cp.CountLabel)
			cp.CountLabel = cp.CountLabel + 1
			return asm0, asm1, nil

		} else if ct.Sub[0].Vtype == "name" && ct.Sub[1].Vtype == "literal" { // v == l
			if cp.VarLocal[ct.Sub[0].Value] == 0 { // global v
				asm0 = asm0 + fmt.Sprintf("load ma global &%d\n", cp.VarGlobal[ct.Sub[0].Value])
			} else { // local v
				asm0 = asm0 + fmt.Sprintf("load ma local &%d\n", cp.VarLocal[ct.Sub[0].Value])
			}
			asm0 = asm0 + fmt.Sprintf("load mb const &%d\n", ct.Sub[1].Ivalue.Vpos)
			asm0 = asm0 + fmt.Sprintf("jmpi $%d @%d ;condition\n", opcode[ct.Value], cp.CountLabel)
			asm1 = fmt.Sprintf("label @%d\n", cp.CountLabel)
			cp.CountLabel = cp.CountLabel + 1
			return asm0, asm1, nil

		} else if ct.Sub[0].Vtype == "literal" && ct.Sub[1].Vtype == "name" { // l == v
			asm0 = asm0 + fmt.Sprintf("load ma const &%d\n", ct.Sub[0].Ivalue.Vpos)
			if cp.VarLocal[ct.Sub[1].Value] == 0 { // global v
				asm0 = asm0 + fmt.Sprintf("load mb global &%d\n", cp.VarGlobal[ct.Sub[1].Value])
			} else { // local v
				asm0 = asm0 + fmt.Sprintf("load mb local &%d\n", cp.VarLocal[ct.Sub[1].Value])
			}
			asm0 = asm0 + fmt.Sprintf("jmpi $%d @%d ;condition\n", opcode[ct.Value], cp.CountLabel)
			asm1 = fmt.Sprintf("label @%d\n", cp.CountLabel)
			cp.CountLabel = cp.CountLabel + 1
			return asm0, asm1, nil

		} else if ct.Sub[0].Vtype == "literal" && ct.Sub[1].Vtype == "literal" { // l == m
			asm0 = asm0 + fmt.Sprintf("load ma const &%d\n", ct.Sub[0].Ivalue.Vpos)
			asm0 = asm0 + fmt.Sprintf("load mb const &%d\n", ct.Sub[1].Ivalue.Vpos)
			asm0 = asm0 + fmt.Sprintf("jmpi $%d @%d ;condition\n", opcode[ct.Value], cp.CountLabel)
			asm1 = fmt.Sprintf("label @%d\n", cp.CountLabel)
			cp.CountLabel = cp.CountLabel + 1
			return asm0, asm1, nil
		}
	}

	asm0, err = cp.genkasm(ct)
	asm0 = asm0 + fmt.Sprintf("jmpiff @%d ;condition\n", cp.CountLabel)
	asm1 = fmt.Sprintf("label @%d\n", cp.CountLabel)
	cp.CountLabel = cp.CountLabel + 1
	return asm0, asm1, err
}

// phase 6 - gen kasm (function)
func (cp *Compiler) genfunc(ct *CTree) (string, error) {
	var err error
	err = nil
	asm := ""
	if ct.Vtype != "function" {
		return asm, errors.New("e418 : compiling non-function")
	}
	t0, t1 := cp.genkasm(&ct.Sub[len(ct.Sub)-1])
	err = t1
	asm = asm + fmt.Sprintf("label @%d ;function_start\n", cp.InnerNum[ct.Value]) + t0 + fmt.Sprintf("ret $%d ;function_end\n", cp.InnerParms[ct.Value])
	return asm, err
}

// phase 7 - gen kasm (rodata / data)
func (cp *Compiler) gendata() (string, string) {
	rodata := make([]string, len(cp.CountLiteral))
	for _, r := range cp.CountLiteral {
		switch x := r.Ivalue.(type) {
		case bool:
			if x {
				rodata[r.Vpos] = "data bool t"
			} else {
				rodata[r.Vpos] = "data bool f"
			}
		case int:
			rodata[r.Vpos] = fmt.Sprintf("data int %d", x)
		case float64:
			rodata[r.Vpos] = fmt.Sprintf("data float %f", x)
		case string:
			rodata[r.Vpos] = fmt.Sprintf("data string %s", kio.Bprint([]byte(x)))
		case []byte:
			rodata[r.Vpos] = fmt.Sprintf("data bytes %s", kio.Bprint(x))
		case nil:
			rodata[r.Vpos] = "data none"
		default:
			rodata[r.Vpos] = "data none"
		}
		rodata[r.Vpos] = rodata[r.Vpos] + fmt.Sprintf(" ;const_%d", r.Vpos)
	}
	data := make([]string, len(cp.VarGlobal))
	for i := 0; i < len(cp.VarGlobal); i++ {
		data[i] = fmt.Sprintf("data int 0 ;global_%d", i+len(cp.CountLiteral))
	}
	return strings.Join(rodata, "\n"), strings.Join(data, "\n")
}

// init compiler with basic settings
func (cp *Compiler) Init() {
	cp.OptConst = true
	cp.OptAsm = true

	cp.OuterNum = map[string]int{"test.input": 16, "test.print": 17, "test.read": 18, "test.write": 19, "test.time": 20, "test.sleep": 21}
	cp.OuterParms = map[string]int{"test.input": 1, "test.print": 1, "test.read": 2, "test.write": 2, "test.time": 0, "test.sleep": 1}
	cp.InnerNum = make(map[string]int)
	cp.InnerParms = make(map[string]int)
	cp.InnerLocal = make(map[string]int)

	cp.VarGlobal = make(map[string]int)
	cp.VarLocal = make(map[string]int)
	cp.VarArgs = 0
	cp.CountLocal = 0

	cp.CountLabel = 1
	cp.CountFor = 1
	cp.CountLiteral = make(map[string]CLiteral)
	temp := []string{"none", "None", "bool", "True", "bool", "False", "int", "0", "int", "1", "float", "0.0", "string", "", "bytes", ""}
	for i := 0; i < len(temp)/2; i++ {
		var cl CLiteral
		cl.Set(temp[2*i], temp[2*i+1], i)
		cp.CountLiteral[temp[2*i]+temp[2*i+1]] = cl
	}

	cp.il_tree = make([]CTree, 0)
	cp.il_main = ""
	cp.il_func = make([]string, 0)
}

// compile AST -> kasm
func (cp *Compiler) Compile(mainflow *AST, functions []AST) (asm string, err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("%s at Compiler.Compile exit", ferr)
		}
	}()
	asm = ""
	err = nil
	cp.il_tree = make([]CTree, len(functions)+1)
	cp.il_main = ""
	cp.il_func = make([]string, len(functions))

	// phase 0 - set ctree
	err = cp.il_tree[0].Set(*mainflow, cp.OptConst)
	if err != nil {
		return asm, err
	}
	for i, r := range functions {
		err = cp.il_tree[i+1].Set(r, cp.OptConst)
		if err != nil {
			return asm, err
		}
	}

	// phase 1 - update const
	for i := 0; i < len(cp.il_tree); i++ {
		cp.readconst(&cp.il_tree[i])
	}

	// phase 2 - update inner functions
	for i := 1; i < len(cp.il_tree); i++ {
		err = cp.readfunc(&cp.il_tree[i])
		if err != nil {
			return asm, fmt.Errorf("%s at Compiler.Compile p2[%d]", err, i)
		}
	}

	// phase 3 - set global variable
	err = cp.readglobal(&cp.il_tree[0])
	if err != nil {
		return asm, fmt.Errorf("%s at Compiler.Compile p3[0]", err)
	}

	// phase 4 - set inner local num
	for i := 1; i < len(cp.il_tree); i++ {
		cp.VarLocal = make(map[string]int)
		cp.VarArgs = cp.InnerParms[cp.il_tree[i].Value]
		cp.CountLocal = 0
		cp.readlocal(&cp.il_tree[i])
		cp.InnerLocal[cp.il_tree[i].Value] = len(cp.VarLocal)
	}

	// phase 5 - compile mainflow
	cp.VarLocal = make(map[string]int)
	cp.VarArgs = 0
	cp.CountLocal = 0
	cp.il_main, err = cp.genkasm(&cp.il_tree[0])
	if err != nil {
		return asm, fmt.Errorf("%s at Compiler.Compile p5[0]", err)
	}

	// repeat for inner functions
	for i := 1; i < len(cp.il_tree); i++ {
		// phase 4 - set local variable
		cp.VarLocal = make(map[string]int)
		cp.VarArgs = cp.InnerParms[cp.il_tree[i].Value]
		cp.CountLocal = 0
		err = cp.readlocal(&cp.il_tree[i])
		if err != nil {
			return asm, fmt.Errorf("%s at Compiler.Compile p4[%d]", err, i)
		}

		// phase 6 - compile inner functions
		cp.il_func[i-1], err = cp.genfunc(&cp.il_tree[i])
		if err != nil {
			return asm, fmt.Errorf("%s at Compiler.Compile p6[%d]", err, i)
		}
	}

	// phase 7 - compile const & global
	rodata, data := cp.gendata()

	asm = fmt.Sprintf(".rodata\n%s\n\n.data\n%s\n\n.text\nnop ;program_start\n%shlt ;program_end\n", rodata, data, cp.il_main)
	for _, r := range cp.il_func {
		asm = asm + "\n" + r
	}
	return asm, err
}

// ctree literal value (part 4.5 ~ 4.6)
type CLiteral struct {
	Vtype  string      // none bool int float string bytes
	Vpos   int         // rodata position
	Ivalue interface{} // bool/int/float/string/bytes
}

// set CLiteral value
func (cl *CLiteral) Set(vtype string, value string, pos int) (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("%s at CLiteral.Set exit", ferr)
		}
	}()
	err = nil
	cl.Vtype = vtype
	cl.Vpos = pos
	cl.Ivalue = nil

	switch vtype {
	case "none":
		cl.Ivalue = nil
	case "bool":
		if value == "True" {
			cl.Ivalue = true
		} else {
			cl.Ivalue = false
		}
	case "int":
		cl.Ivalue, err = strconv.Atoi(value)
	case "float":
		cl.Ivalue, err = strconv.ParseFloat(value, 64)
	case "string":
		cl.Ivalue = value
	case "bytes":
		cl.Ivalue, err = kio.Bread(value)
	case "fornum", "funcname":
		cl.Ivalue = value
	default:
		err = fmt.Errorf("e450 : invalid type %s", vtype)
	}
	return err
}

// optimized compile tree (part 4.7 ~ 4.9)
type CTree struct {
	Value  string    // string formatted value
	Vtype  string    // literal name bicalc innercall outercall assign function control
	Sub    []CTree   // optimized sub component
	Ivalue *CLiteral // immediate value (literal only)
}

func (ct *CTree) add(a interface{}, b interface{}) (*AST, error) {
	tp0 := "e470"
	tp1 := "e470"
	var ret *AST

	switch a := a.(type) {
	case bool:
		tp0 = "bool"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
			ret = ct.setast(a || b)
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case int:
		tp0 = "int"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(a + b)
		case float64:
			tp1 = "float"
			ret = ct.setast(float64(a) + b)
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case float64:
		tp0 = "float"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(a + float64(b))
		case float64:
			tp1 = "float"
			ret = ct.setast(a + b)
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case string:
		tp0 = "string"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
			ret = ct.setast(a + b)
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case []byte:
		tp0 = "bytes"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
			ret = ct.setast(append(a, b...))
		case nil:
			tp1 = "none"
		}

	case nil:
		tp0 = "none"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("e471 : wrong type ADD(%s, %s)", tp0, tp1)
	} else {
		return ret, nil
	}
}

func (ct *CTree) sub(a interface{}, b interface{}) (*AST, error) {
	tp0 := "e470"
	tp1 := "e470"
	var ret *AST

	switch a := a.(type) {
	case bool:
		tp0 = "bool"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case int:
		tp0 = "int"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(a - b)
		case float64:
			tp1 = "float"
			ret = ct.setast(float64(a) - b)
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case float64:
		tp0 = "float"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(a - float64(b))
		case float64:
			tp1 = "float"
			ret = ct.setast(a - b)
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case string:
		tp0 = "string"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case []byte:
		tp0 = "bytes"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case nil:
		tp0 = "none"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("e472 : wrong type SUB(%s, %s)", tp0, tp1)
	} else {
		return ret, nil
	}
}

func (ct *CTree) mul(a interface{}, b interface{}) (*AST, error) {
	tp0 := "e470"
	tp1 := "e470"
	var ret *AST

	switch a := a.(type) {
	case bool:
		tp0 = "bool"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
			ret = ct.setast(a && b)
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case int:
		tp0 = "int"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(a * b)
		case float64:
			tp1 = "float"
			ret = ct.setast(float64(a) * b)
		case string:
			tp1 = "string"
			if a < 0 {
				tp0 = "int_n"
			} else {
				ret = ct.setast(strings.Repeat(b, a))
			}
		case []byte:
			tp1 = "bytes"
			if a < 0 {
				tp0 = "int_n"
			} else {
				ret = ct.setast(slices.Repeat(b, a))
			}
		case nil:
			tp1 = "none"
		}

	case float64:
		tp0 = "float"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(a * float64(b))
		case float64:
			tp1 = "float"
			ret = ct.setast(a * b)
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case string:
		tp0 = "string"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			if b < 0 {
				tp0 = "int_n"
			} else {
				ret = ct.setast(strings.Repeat(a, b))
			}
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case []byte:
		tp0 = "bytes"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			if b < 0 {
				tp0 = "int_n"
			} else {
				ret = ct.setast(slices.Repeat(a, b))
			}
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case nil:
		tp0 = "none"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("e473 : wrong type MUL(%s, %s)", tp0, tp1)
	} else {
		return ret, nil
	}
}

func (ct *CTree) div(a interface{}, b interface{}) (*AST, error) {
	tp0 := "e470"
	tp1 := "e470"
	var ret *AST

	switch a := a.(type) {
	case bool:
		tp0 = "bool"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case int:
		tp0 = "int"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			if b == 0 {
				tp1 = "int_0"
			} else {
				ret = ct.setast(float64(a) / float64(b))
			}
		case float64:
			tp1 = "float"
			if b == 0.0 {
				tp1 = "float_0"
			} else {
				ret = ct.setast(float64(a) / b)
			}
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case float64:
		tp0 = "float"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			if b == 0 {
				tp1 = "int_0"
			} else {
				ret = ct.setast(a / float64(b))
			}
		case float64:
			tp1 = "float"
			if b == 0.0 {
				tp1 = "float_0"
			} else {
				ret = ct.setast(a / b)
			}
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case string:
		tp0 = "string"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case []byte:
		tp0 = "bytes"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case nil:
		tp0 = "none"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("e474 : wrong type DIV(%s, %s)", tp0, tp1)
	} else {
		return ret, nil
	}
}

func (ct *CTree) divs(a interface{}, b interface{}) (*AST, error) {
	tp0 := "e470"
	tp1 := "e470"
	var ret *AST

	switch a := a.(type) {
	case bool:
		tp0 = "bool"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case int:
		tp0 = "int"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			if b == 0 {
				tp1 = "int_0"
			} else if a > 0 && b > 0 {
				ret = ct.setast(a / b)
			} else {
				ret = ct.setast(int(math.Floor(float64(a) / float64(b))))
			}
		case float64:
			tp1 = "float"
			if b == 0.0 {
				tp1 = "float_0"
			} else {
				ret = ct.setast(math.Floor(float64(a) / b))
			}
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case float64:
		tp0 = "float"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			if b == 0 {
				tp1 = "int_0"
			} else {
				ret = ct.setast(math.Floor(a / float64(b)))
			}
		case float64:
			tp1 = "float"
			if b == 0.0 {
				tp1 = "float_0"
			} else {
				ret = ct.setast(math.Floor(a / b))
			}
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case string:
		tp0 = "string"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case []byte:
		tp0 = "bytes"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case nil:
		tp0 = "none"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("e475 : wrong type DIVS(%s, %s)", tp0, tp1)
	} else {
		return ret, nil
	}
}

func (ct *CTree) divr(a interface{}, b interface{}) (*AST, error) {
	tp0 := "e470"
	tp1 := "e470"
	var ret *AST

	switch a := a.(type) {
	case bool:
		tp0 = "bool"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case int:
		tp0 = "int"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			if b == 0 {
				tp1 = "int_0"
			} else if a > 0 && b > 0 {
				ret = ct.setast(a % b)
			} else {
				ret = ct.setast(int(float64(a) - float64(b)*math.Floor(float64(a)/float64(b))))
			}
		case float64:
			tp1 = "float"
			if b == 0.0 {
				tp1 = "float_0"
			} else {
				ret = ct.setast(float64(a) - b*math.Floor(float64(a)/b))
			}
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case float64:
		tp0 = "float"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			if b == 0 {
				tp1 = "int_0"
			} else {
				ret = ct.setast(a - float64(b)*math.Floor(a/float64(b)))
			}
		case float64:
			tp1 = "float"
			if b == 0.0 {
				tp1 = "float_0"
			} else {
				ret = ct.setast(a - b*math.Floor(a/b))
			}
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case string:
		tp0 = "string"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case []byte:
		tp0 = "bytes"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case nil:
		tp0 = "none"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("e476 : wrong type DIVR(%s, %s)", tp0, tp1)
	} else {
		return ret, nil
	}
}

func (ct *CTree) pow(a interface{}, b interface{}) (*AST, error) {
	tp0 := "e470"
	tp1 := "e470"
	var ret *AST

	switch a := a.(type) {
	case bool:
		tp0 = "bool"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case int:
		tp0 = "int"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			if b >= 0 {
				ret = ct.setast(int(math.Pow(float64(a), float64(b))))
			} else {
				ret = ct.setast(math.Pow(float64(a), float64(b)))
			}
		case float64:
			tp1 = "float"
			ret = ct.setast(math.Pow(float64(a), b))
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case float64:
		tp0 = "float"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(math.Pow(a, float64(b)))
		case float64:
			tp1 = "float"
			ret = ct.setast(math.Pow(a, b))
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case string:
		tp0 = "string"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case []byte:
		tp0 = "bytes"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case nil:
		tp0 = "none"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("e477 : wrong type POW(%s, %s)", tp0, tp1)
	} else {
		return ret, nil
	}
}

func (ct *CTree) eql(a interface{}, b interface{}) (*AST, error) {
	tp0 := "e470"
	tp1 := "e470"
	var ret *AST

	switch a := a.(type) {
	case bool:
		tp0 = "bool"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
			ret = ct.setast(a == b)
		case int:
			tp1 = "int"
			ret = ct.setast(false)
		case float64:
			tp1 = "float"
			ret = ct.setast(false)
		case string:
			tp1 = "string"
			ret = ct.setast(false)
		case []byte:
			tp1 = "bytes"
			ret = ct.setast(false)
		case nil:
			tp1 = "none"
			ret = ct.setast(false)
		}

	case int:
		tp0 = "int"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
			ret = ct.setast(false)
		case int:
			tp1 = "int"
			ret = ct.setast(a == b)
		case float64:
			tp1 = "float"
			ret = ct.setast(float64(a) == b)
		case string:
			tp1 = "string"
			ret = ct.setast(false)
		case []byte:
			tp1 = "bytes"
			ret = ct.setast(false)
		case nil:
			tp1 = "none"
			ret = ct.setast(false)
		}

	case float64:
		tp0 = "float"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
			ret = ct.setast(false)
		case int:
			tp1 = "int"
			ret = ct.setast(a == float64(b))
		case float64:
			tp1 = "float"
			ret = ct.setast(a == b)
		case string:
			tp1 = "string"
			ret = ct.setast(false)
		case []byte:
			tp1 = "bytes"
			ret = ct.setast(false)
		case nil:
			tp1 = "none"
			ret = ct.setast(false)
		}

	case string:
		tp0 = "string"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
			ret = ct.setast(false)
		case int:
			tp1 = "int"
			ret = ct.setast(false)
		case float64:
			tp1 = "float"
			ret = ct.setast(false)
		case string:
			tp1 = "string"
			ret = ct.setast(a == b)
		case []byte:
			tp1 = "bytes"
			ret = ct.setast(false)
		case nil:
			tp1 = "none"
			ret = ct.setast(false)
		}

	case []byte:
		tp0 = "bytes"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
			ret = ct.setast(false)
		case int:
			tp1 = "int"
			ret = ct.setast(false)
		case float64:
			tp1 = "float"
			ret = ct.setast(false)
		case string:
			tp1 = "string"
			ret = ct.setast(false)
		case []byte:
			tp1 = "bytes"
			ret = ct.setast(kio.Bequal(a, b))
		case nil:
			tp1 = "none"
			ret = ct.setast(false)
		}

	case nil:
		tp0 = "none"
		switch b.(type) {
		case bool:
			tp1 = "bool"
			ret = ct.setast(false)
		case int:
			tp1 = "int"
			ret = ct.setast(false)
		case float64:
			tp1 = "float"
			ret = ct.setast(false)
		case string:
			tp1 = "string"
			ret = ct.setast(false)
		case []byte:
			tp1 = "bytes"
			ret = ct.setast(false)
		case nil:
			tp1 = "none"
			ret = ct.setast(true)
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("e478 : wrong type EQL(%s, %s)", tp0, tp1)
	} else {
		return ret, nil
	}
}

func (ct *CTree) sml(a interface{}, b interface{}) (*AST, error) {
	tp0 := "e470"
	tp1 := "e470"
	var ret *AST

	switch a := a.(type) {
	case bool:
		tp0 = "bool"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case int:
		tp0 = "int"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(a < b)
		case float64:
			tp1 = "float"
			ret = ct.setast(float64(a) < b)
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case float64:
		tp0 = "float"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(a < float64(b))
		case float64:
			tp1 = "float"
			ret = ct.setast(a < b)
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case string:
		tp0 = "string"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
			ret = ct.setast(a < b)
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case []byte:
		tp0 = "bytes"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
			ret = ct.setast(bytes.Compare(a, b) == -1)
		case nil:
			tp1 = "none"
		}

	case nil:
		tp0 = "none"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("e479 : wrong type SML(%s, %s)", tp0, tp1)
	} else {
		return ret, nil
	}
}

func (ct *CTree) smle(a interface{}, b interface{}) (*AST, error) {
	tp0 := "e470"
	tp1 := "e470"
	var ret *AST

	switch a := a.(type) {
	case bool:
		tp0 = "bool"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case int:
		tp0 = "int"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(a <= b)
		case float64:
			tp1 = "float"
			ret = ct.setast(float64(a) <= b)
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case float64:
		tp0 = "float"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
			ret = ct.setast(a <= float64(b))
		case float64:
			tp1 = "float"
			ret = ct.setast(a <= b)
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case string:
		tp0 = "string"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
			ret = ct.setast(a <= b)
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}

	case []byte:
		tp0 = "bytes"
		switch b := b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
			ret = ct.setast(bytes.Compare(a, b) != 1)
		case nil:
			tp1 = "none"
		}

	case nil:
		tp0 = "none"
		switch b.(type) {
		case bool:
			tp1 = "bool"
		case int:
			tp1 = "int"
		case float64:
			tp1 = "float"
		case string:
			tp1 = "string"
		case []byte:
			tp1 = "bytes"
		case nil:
			tp1 = "none"
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("e480 : wrong type SMLE(%s, %s)", tp0, tp1)
	} else {
		return ret, nil
	}
}

func (ct *CTree) setast(v interface{}) (ret *AST) {
	var temp AST
	ret = &temp
	switch v := v.(type) {
	case bool:
		temp.Vtype = "bool"
		if v {
			temp.Value = "True"
		} else {
			temp.Value = "False"
		}
	case int:
		temp.Vtype = "int"
		temp.Value = fmt.Sprintf("%d", v)
	case float64:
		temp.Vtype = "float"
		temp.Value = fmt.Sprintf("%f", v)
	case string:
		temp.Vtype = "string"
		temp.Value = v
	case []byte:
		temp.Vtype = "bytes"
		temp.Value = kio.Bprint(v)
	case nil:
		temp.Vtype = "none"
		temp.Value = "None"
	default:
		ret = nil
	}
	return ret
}

func (ct *CTree) optimize() error {
	if ct.Sub[0].Vtype != "literal" || ct.Sub[1].Vtype != "literal" {
		return nil
	}
	var ast *AST
	var err error
	a := ct.Sub[0].Ivalue.Ivalue
	b := ct.Sub[1].Ivalue.Ivalue

	switch ct.Value {
	case "+":
		ast, err = ct.add(a, b)
	case "-":
		ast, err = ct.sub(a, b)
	case "*":
		ast, err = ct.mul(a, b)
	case "/":
		ast, err = ct.div(a, b)
	case "//":
		ast, err = ct.divs(a, b)
	case "%":
		ast, err = ct.divr(a, b)
	case "**":
		ast, err = ct.pow(a, b)
	case "==":
		ast, err = ct.eql(a, b)
	case "!=":
		ast, err = ct.eql(a, b)
		if ast.Value == "True" {
			ast.Value = "False"
		} else {
			ast.Value = "True"
		}
	case "<":
		ast, err = ct.sml(a, b)
	case "<=":
		ast, err = ct.smle(a, b)
	case ">":
		ast, err = ct.sml(b, a)
	case ">=":
		ast, err = ct.smle(b, a)
	default:
		return fmt.Errorf("e482 : unsupported bicalc %s", ct.Value)
	}

	if err != nil {
		return err
	}
	return ct.Set(*ast, true)
}

// set CTree, fold const
func (ct *CTree) Set(ast AST, opt bool) error {
	var err error
	err = nil

	switch ast.Vtype {
	case "none", "bool", "int", "float", "string", "bytes":
		ct.Value = ast.Value
		ct.Vtype = "literal"
		ct.Sub = nil
		var temp CLiteral
		err = temp.Set(ast.Vtype, ast.Value, -1) // set CLpos after
		ct.Ivalue = &temp

	case "bicalc":
		ct.Value = ast.Value
		ct.Vtype = ast.Vtype
		ct.Sub = make([]CTree, 2)
		ct.Ivalue = nil
		err = ct.Sub[0].Set(ast.Sub[0], opt)
		if err != nil {
			return err
		}
		err = ct.Sub[1].Set(ast.Sub[1], opt)
		if err != nil {
			return err
		}
		if opt {
			err = ct.optimize()
		}

	default:
		ct.Value = ast.Value
		ct.Vtype = ast.Vtype
		ct.Sub = make([]CTree, len(ast.Sub))
		ct.Ivalue = nil
		for i, r := range ast.Sub {
			err = ct.Sub[i].Set(r, opt)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// kasm to kelf (part 5)
type Assembler struct {
	Icon []byte // ksc5 icon
	Info string // program name & info
	ABIf int    // ABI format

	public  string // public key (SHA3-512 ro + d + t)
	private string // private key (SHA3-512 ro + d + t)
}

// encode little-endian signed int 16 / 32
func (as *Assembler) enc(data int, fullword bool) ([]byte, error) {
	var out []byte
	if fullword { // i32
		if data < math.MinInt32 || data > math.MaxInt32 {
			return nil, errors.New("e500 : overflow i32")
		}
		out = make([]byte, 4)
		out[0] = byte(data & 0xff)
		out[1] = byte((data >> 8) & 0xff)
		out[2] = byte((data >> 16) & 0xff)
		out[3] = byte((data >> 24) & 0xff)
	} else { // i16
		if data < math.MinInt16 || data > math.MaxInt16 {
			return nil, errors.New("e501 : overflow i16")
		}
		out = make([]byte, 2)
		out[0] = byte(data & 0xff)
		out[1] = byte((data >> 8) & 0xff)
	}
	return out, nil
}

// split tokens from kasm
func (as *Assembler) tokenize(kasm string) [][]string {
	out := make([][]string, 0)
	for _, r := range strings.Split(strings.ReplaceAll(kasm, "\r", "\n"), "\n") {
		temp := make([]string, 0)
		flag := true
		for _, l := range strings.Split(r, " ") {
			if l != "" {
				if l[0] == ';' {
					flag = false
				} else if flag {
					temp = append(temp, l)
				}
			}
		}
		if len(temp) != 0 {
			out = append(out, temp)
		}
	}
	return out
}

// kasm code : data *
func (as *Assembler) asmdata(tokens [][]string) ([]byte, error) {
	out := make([][]byte, len(tokens))
	for i, r := range tokens {
		if r[0] != "data" {
			return nil, fmt.Errorf("e502 : invalid opcode %s", r[0])
		}

		switch r[1] {
		case "none":
			out[i] = []byte{78}

		case "bool":
			if r[2] == "t" {
				out[i] = []byte{66, 0}
			} else {
				out[i] = []byte{66, 1}
			}

		case "int":
			tgt, err := strconv.Atoi(r[2])
			temp := bytes.NewBuffer(nil)
			binary.Write(temp, binary.LittleEndian, int64(tgt))
			out[i] = append([]byte{73}, temp.Bytes()...)
			if err != nil {
				return nil, fmt.Errorf("%s at Assembler.asmdata int", err)
			}

		case "float":
			tgt, err := strconv.ParseFloat(r[2], 64)
			temp := make([]byte, 8)
			binary.LittleEndian.PutUint64(temp, math.Float64bits(tgt))
			out[i] = append([]byte{70}, temp...)
			if err != nil {
				return nil, fmt.Errorf("%s at Assembler.asmdata float", err)
			}

		case "string":
			if len(r) == 2 {
				out[i] = []byte{83, 0, 0, 0, 0}
			} else {
				t0, e0 := kio.Bread(r[2])
				t1, e1 := as.enc(len(t0), true)
				out[i] = append(append([]byte{83}, t1...), t0...)
				if e0 != nil || e1 != nil {
					return nil, fmt.Errorf("(%s %s) at Assembler.asmdata string", e0, e1)
				}
			}

		case "bytes":
			if len(r) == 2 {
				out[i] = []byte{67, 0, 0, 0, 0}
			} else {
				t0, e0 := kio.Bread(r[2])
				t1, e1 := as.enc(len(t0), true)
				out[i] = append(append([]byte{67}, t1...), t0...)
				if e0 != nil || e1 != nil {
					return nil, fmt.Errorf("(%s %s) at Assembler.asmdata bytes", e0, e1)
				}
			}

		default:
			return nil, fmt.Errorf("e503 : invalid datatype %s", r[1])
		}
	}
	return bytes.Join(out, nil), nil
}

// kasm code : i16 / i32
func (as *Assembler) asmint(itype string, icode string, fullword bool) ([]byte, error) {
	if icode[0] != itype[0] {
		return nil, errors.New("e504 : invalid int syntax")
	}
	t0, e0 := strconv.Atoi(icode[1:])
	t1, e1 := as.enc(t0, fullword)
	if e0 == nil && e1 == nil {
		return t1, nil
	} else {
		return nil, fmt.Errorf("(%s %s) at Assembler.asmint", e0, e1)
	}
}

// kasm code : seg addr (i16 i32)
func (as *Assembler) asmvar(seg string, addr string) ([]byte, error) {
	var out []byte
	if seg == "const" {
		out = []byte{99, 0}
	} else if seg == "global" {
		out = []byte{103, 0}
	} else if seg == "local" {
		out = []byte{108, 0}
	} else {
		return nil, errors.New("e505 : invalid segment token")
	}
	temp, err := as.asmint("&", addr, true)
	if err == nil {
		return append(out, temp...), nil
	} else {
		return nil, err
	}
}

// kasm code : reg
func (as *Assembler) asmreg(reg string) ([]byte, error) {
	if reg == "ma" {
		return []byte{97}, nil
	} else if reg == "mb" {
		return []byte{98}, nil
	} else {
		return nil, fmt.Errorf("e506 : invalid register %s", reg)
	}
}

// kasm code : jmp pos
func (as *Assembler) asmjmp(pos string, table map[string]int) ([]byte, error) {
	t0, e0 := table[pos]
	t1, e1 := as.enc(t0, true)
	if e0 {
		if e1 == nil {
			return t1, nil
		} else {
			return nil, fmt.Errorf("%s at Assembler.asmjmp", e1)
		}
	} else {
		return nil, fmt.Errorf("e507 : label not exist %s", pos)
	}
}

// convert tokens to binary code (rodata, data, text)
func (as *Assembler) assemble(rodata [][]string, data [][]string, text [][]string) ([]byte, []byte, []byte, error) {
	var rodata_out []byte
	var data_out []byte
	text_out := make([][]byte, len(text))
	label := make(map[string]int)
	var err error

	rodata_out, err = as.asmdata(rodata) // assemble rodata
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s at Assembler.assemble rodata", err)
	}

	data_out, err = as.asmdata(data) // assemble data
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s at Assembler.assemble data", err)
	}

	for i, r := range text { // labeling text
		if r[0] == "label" {
			if r[1][0] == '@' {
				label[r[1]] = i
			} else {
				return nil, nil, nil, fmt.Errorf("e508 : invalid label ptr %s at text", r[1])
			}
		}
	}

	for i, r := range text { // assemble text
		switch r[0] {
		case "hlt": // 0x00
			text_out[i] = []byte{0, 0, 0, 0, 0, 0, 0, 0}

		case "nop", "label": // 0x01
			text_out[i] = []byte{1, 0, 0, 0, 0, 0, 0, 0}

		case "intr": // 0x10 i16 i32
			text_out[i] = []byte{16, 0}
			fr, er := as.asmint("$", r[1], false)
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble intr", er)
			} else {
				fr, er = as.asmint("$", r[2], true)
				text_out[i] = append(text_out[i], fr...)
				if er != nil {
					err = fmt.Errorf("%s at Assembler.assemble intr", er)
				}
			}

		case "call": // 0x11 r i16 i32
			text_out[i] = []byte{17}
			fr, er := as.asmreg(r[1])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble call", er)
			} else {
				fr, er = as.asmint("$", r[2], false)
				text_out[i] = append(text_out[i], fr...)
				if er != nil {
					err = fmt.Errorf("%s at Assembler.assemble call", er)
				} else {
					fr, er = as.asmjmp(r[3], label)
					text_out[i] = append(text_out[i], fr...)
					if er != nil {
						err = fmt.Errorf("%s at Assembler.assemble call", er)
					}
				}
			}

		case "ret": // 0x12 i16
			text_out[i] = []byte{18, 0}
			fr, er := as.asmint("$", r[1], false)
			text_out[i] = append(append(text_out[i], fr...), 0, 0, 0, 0)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble ret", er)
			}

		case "jmp": // 0x13 i32
			text_out[i] = []byte{19, 0, 0, 0}
			fr, er := as.asmjmp(r[1], label)
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble jmp", er)
			}

		case "jmpiff": // 0x14 i32
			text_out[i] = []byte{20, 0, 0, 0}
			fr, er := as.asmjmp(r[1], label)
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble jmpiff", er)
			}

		case "forcond": // 0x15 r i16 i32
			text_out[i] = []byte{21}
			fr, er := as.asmreg(r[1])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble forcond", er)
			} else {
				fr, er = as.asmvar(r[2], r[3])
				text_out[i] = append(text_out[i], fr...)
				if er != nil {
					err = fmt.Errorf("%s at Assembler.assemble forcond", er)
				}
			}

		case "forset": // 0x16 r i16 i32
			text_out[i] = []byte{22}
			fr, er := as.asmreg(r[1])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble forset", er)
			} else {
				fr, er = as.asmvar(r[2], r[3])
				text_out[i] = append(text_out[i], fr...)
				if er != nil {
					err = fmt.Errorf("%s at Assembler.assemble forset", er)
				}
			}

		case "load": // 0x20 r i16 i32
			text_out[i] = []byte{32}
			fr, er := as.asmreg(r[1])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble load", er)
			} else {
				fr, er = as.asmvar(r[2], r[3])
				text_out[i] = append(text_out[i], fr...)
				if er != nil {
					err = fmt.Errorf("%s at Assembler.assemble load", er)
				}
			}

		case "store": // 0x21 r i16 i32
			text_out[i] = []byte{33}
			fr, er := as.asmreg(r[1])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble store", er)
			} else {
				fr, er = as.asmvar(r[2], r[3])
				text_out[i] = append(text_out[i], fr...)
				if er != nil {
					err = fmt.Errorf("%s at Assembler.assemble store", er)
				}
			}

		case "push": // 0x22 r
			text_out[i] = []byte{34}
			fr, er := as.asmreg(r[1])
			text_out[i] = append(append(text_out[i], fr...), 0, 0, 0, 0, 0, 0)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble push", er)
			}

		case "pop": // 0x23 r
			text_out[i] = []byte{35}
			fr, er := as.asmreg(r[1])
			text_out[i] = append(append(text_out[i], fr...), 0, 0, 0, 0, 0, 0)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble pop", er)
			}

		case "pushset": // 0x24 i16 i32
			text_out[i] = []byte{36, 0}
			fr, er := as.asmvar(r[1], r[2])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble pushset", er)
			}

		case "popset": // 0x25 i16 i32
			text_out[i] = []byte{37, 0}
			fr, er := as.asmvar(r[1], r[2])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble popset", er)
			}

		case "add": // 0x30
			text_out[i] = []byte{48, 0, 0, 0, 0, 0, 0, 0}

		case "sub": // 0x31
			text_out[i] = []byte{49, 0, 0, 0, 0, 0, 0, 0}

		case "mul": // 0x32
			text_out[i] = []byte{50, 0, 0, 0, 0, 0, 0, 0}

		case "div": // 0x33
			text_out[i] = []byte{51, 0, 0, 0, 0, 0, 0, 0}

		case "divs": // 0x40
			text_out[i] = []byte{64, 0, 0, 0, 0, 0, 0, 0}

		case "divr": // 0x41
			text_out[i] = []byte{65, 0, 0, 0, 0, 0, 0, 0}

		case "pow": // 0x42
			text_out[i] = []byte{66, 0, 0, 0, 0, 0, 0, 0}

		case "eql": // 0x50
			text_out[i] = []byte{80, 0, 0, 0, 0, 0, 0, 0}

		case "eqln": // 0x51
			text_out[i] = []byte{81, 0, 0, 0, 0, 0, 0, 0}

		case "sml": // 0x52
			text_out[i] = []byte{82, 0, 0, 0, 0, 0, 0, 0}

		case "grt": // 0x53
			text_out[i] = []byte{83, 0, 0, 0, 0, 0, 0, 0}

		case "smle": // 0x54
			text_out[i] = []byte{84, 0, 0, 0, 0, 0, 0, 0}

		case "grte": // 0x55
			text_out[i] = []byte{85, 0, 0, 0, 0, 0, 0, 0}

		case "inc": // 0x60 i16 i32
			text_out[i] = []byte{96, 0}
			fr, er := as.asmvar(r[1], r[2])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble inc", er)
			}

		case "dec": // 0x61 i16 i32
			text_out[i] = []byte{97, 0}
			fr, er := as.asmvar(r[1], r[2])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble dec", er)
			}

		case "shm": // 0x62 i16 i32
			text_out[i] = []byte{98, 0}
			fr, er := as.asmvar(r[1], r[2])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble shm", er)
			}

		case "shd": // 0x63 i16 i32
			text_out[i] = []byte{99, 0}
			fr, er := as.asmvar(r[1], r[2])
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble shd", er)
			}

		case "addi": // 0x70 i32
			text_out[i] = []byte{112, 0, 0, 0}
			fr, er := as.asmint("$", r[1], true)
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble addi", er)
			}

		case "muli": // 0x71 i32
			text_out[i] = []byte{113, 0, 0, 0}
			fr, er := as.asmint("$", r[1], true)
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble muli", er)
			}

		case "addr": // 0x72 r i32
			text_out[i] = []byte{114}
			fr, er := as.asmreg(r[1])
			text_out[i] = append(append(text_out[i], fr...), 0, 0)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble addr", er)
			} else {
				fr, er := as.asmint("$", r[2], true)
				text_out[i] = append(text_out[i], fr...)
				if er != nil {
					err = fmt.Errorf("%s at Assembler.assemble addr", er)
				}
			}

		case "jmpi": // 0x73 i16 i32
			text_out[i] = []byte{115, 0}
			fr, er := as.asmint("$", r[1], false)
			text_out[i] = append(text_out[i], fr...)
			if er != nil {
				err = fmt.Errorf("%s at Assembler.assemble jmpi", er)
			} else {
				fr, er := as.asmjmp(r[2], label)
				text_out[i] = append(text_out[i], fr...)
				if er != nil {
					err = fmt.Errorf("%s at Assembler.assemble jmpi", er)
				}
			}

		default:
			err = fmt.Errorf("e509 : invalid opcode %s", r[0])
		}
		if err != nil {
			return rodata_out, data_out, bytes.Join(text_out, nil), err
		}
	}

	temp := bytes.Join(text_out, nil)
	if len(temp) != len(text_out)*8 {
		return rodata_out, data_out, temp, errors.New("e510 : assemble text fail")
	}
	return rodata_out, data_out, temp, nil
}

// set icon, keys (empty string to pass)
func (as *Assembler) SetKey(iconpath string, public string, private string) {
	if iconpath != "" {
		f, err := kio.Open(iconpath, "r")
		if err == nil {
			defer f.Close()
			as.Icon, _ = kio.Read(f, -1)
		} else {
			as.Icon = nil
		}
	}
	if public == "" || private == "" {
		as.public = ""
		as.private = ""
	} else {
		as.public = public
		as.private = private
	}
}

// generate kelf format exe file
func (as *Assembler) GenExe(asm string) (exe []byte, err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			exe = nil
			err = fmt.Errorf("%s at Assembler.GenExe exit", ferr)
		}
	}()
	asm = strings.ToLower(asm)
	pos0 := strings.Index(asm, ".rodata")
	if pos0 == -1 {
		return nil, errors.New("e511 : section .rodata not exist")
	}
	pos1 := strings.Index(asm, ".data")
	if pos1 == -1 {
		return nil, errors.New("e512 : section .data not exist")
	}
	pos2 := strings.Index(asm, ".text")
	if pos2 == -1 {
		return nil, errors.New("e513 : section .text not exist")
	}

	// assemble kasm, chunk(rodata, data, text)
	tk0 := as.tokenize(asm[pos0+7 : pos1])
	tk1 := as.tokenize(asm[pos1+5 : pos2])
	tk2 := as.tokenize(asm[pos2+5:])
	rodata, data, text, err := as.assemble(tk0, tk1, tk2)
	chunk := append(append(rodata, data...), text...)
	if err != nil {
		return chunk, fmt.Errorf("%s at Assembler.GenExe assemble", err)
	}

	// generate header (info, abi, sign, public)
	worker0 := kdb.Initkdb()
	worker0.Read("info = \"\"\nabi = 0\nsign = ''\npublic = \"\"")
	worker0.Fix("info", as.Info)
	worker0.Fix("abi", as.ABIf)
	if as.public != "" {
		hworker := sha3.New512()
		hworker.Write(chunk)
		hvalue := hworker.Sum(nil)
		henc, err := ksign.Sign(as.private, hvalue)
		if err != nil {
			return chunk, fmt.Errorf("%s at Assembler.GenExe sign", err)
		}
		worker0.Fix("sign", henc)
		worker0.Fix("public", as.public)
	}
	temp, _ := worker0.Write()
	header := []byte(temp)

	// generate kelf binary
	var out []byte
	as.Icon = append(as.Icon, make([]byte, (51200000-len(as.Icon))%512)...)
	worker1 := ksc.Initksc()
	worker1.Prehead = as.Icon
	worker1.Subtype = []byte("KELF")
	worker1.Reserved = append(ksc.Crc32hash(header), []byte("v5.3")...)
	out, _ = worker1.Writeb()
	out = worker1.Linkb(out, header)
	out = worker1.Linkb(out, rodata)
	out = worker1.Linkb(out, data)
	out = worker1.Linkb(out, text)
	out, _ = worker1.Addb(out, "")
	return out, nil
}
