// test723 : common.kscriptc executable

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"stdlib5/cliget"
	"stdlib5/kio"
	"stdlib5/kobj"
	"strings"
)

// get cliget.pathsel
func getpath() *cliget.PathSel {
	names := []string{"Local User", "Desktop"}
	paths := make([]string, 2)
	paths[0], _ = os.UserHomeDir()
	tp0 := filepath.Join(paths[0], "Desktop")
	tp1 := filepath.Join(paths[0], "desktop")
	tp2 := filepath.Join(paths[0], "바탕화면")

	var err error
	if _, err = os.Stat(tp0); err == nil {
		paths[1] = tp0
	} else if _, err = os.Stat(tp1); err == nil {
		paths[1] = tp1
	} else if _, err = os.Stat(tp2); err == nil {
		paths[1] = tp2
	} else {
		paths[1] = filepath.Join(paths[0], "DESKTOP")
	}

	var out cliget.PathSel
	out.Init(names, paths)
	return &out
}

func main() {
	defer func() { kio.Input("press ENTER to exit... ") }()
	kobj.Repath()
	selopt0 := []string{"program info", "source file", "icon file", "signkey", "view token", "view ast", "view asm", "optimize const", "optimize asm"}
	selopt1 := []string{"string", "file", "file", "keyfile", "bool", "bool", "bool", "bool", "bool"}
	selopt2 := []string{"*", "txt", "webp", "*", "*", "*", "*", "*", "*"}
	basicbool := [][]byte{{}, {}, {}, {}, {1}, {1}, {1}, {0}, {0}}

	pkgs, _ := os.ReadDir("./_ST5_DATA/")
	pkgsname := make([]string, 0)
	for _, r := range pkgs {
		temp := r.Name()
		if len(temp) > 4 && strings.ToLower(temp[len(temp)-4:]) == ".txt" {
			pkgsname = append(pkgsname, "./_ST5_DATA/"+temp)
			selopt0 = append(selopt0, "using "+temp[:len(temp)-4])
			selopt1 = append(selopt1, "bool")
			selopt2 = append(selopt2, "*")
			basicbool = append(basicbool, []byte{0})
		}
	}

	var sel cliget.OptSel
	sel.Init(selopt0, selopt1, selopt2, *getpath(), nil)
	sel.ByteRes = basicbool
	sel.GetOpt()

	var worker langc
	worker.init()
	f, _ := kio.Open(sel.StrRes[1], "r")
	src, _ := kio.Read(f, -1)
	f.Close()
	worker.option = [5]bool{sel.ByteRes[4][0] == 0, sel.ByteRes[5][0] == 0, sel.ByteRes[6][0] == 0, sel.ByteRes[7][0] == 0, sel.ByteRes[8][0] == 0}
	worker.data = [7]string{string(src), sel.StrRes[0], sel.StrRes[2], string(sel.ByteRes[3]), "", "", ""}
	for i, r := range pkgsname {
		if sel.ByteRes[i+9][0] == 0 {
			if err := worker.addpkg(r); err != nil {
				fmt.Printf("[header error] %s\n", err)
				return
			}
		}
	}

	exe, tpass, err := worker.compile()
	fmt.Printf("[compile time] %f\n", tpass)
	if err != nil {
		fmt.Printf("[compile error] %s\n", err)
	}
	os.RemoveAll("./result/")
	os.Mkdir("./result/", os.ModePerm)
	f, _ = kio.Open("./result/exe.webp", "w")
	kio.Write(f, exe)
	f.Close()
	if worker.option[0] {
		f, _ = kio.Open("./result/token.txt", "w")
		kio.Write(f, []byte(worker.data[4]))
		f.Close()
	}
	if worker.option[1] {
		f, _ = kio.Open("./result/ast.txt", "w")
		kio.Write(f, []byte(worker.data[5]))
		f.Close()
	}
	if worker.option[2] {
		f, _ = kio.Open("./result/kasm.txt", "w")
		kio.Write(f, []byte(worker.data[6]))
		f.Close()
	}
	fmt.Println("[result generated] ./result/")
}
