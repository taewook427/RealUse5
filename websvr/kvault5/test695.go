package main

import (
	"fmt"
	"stdlib5/kaes"
	"stdlib5/kio"
	"stdlib5/kvault"
	"strings"
	"time"
)

// basic test
func test0(remote string) {
	a, b, c := kvault.Test_Basic(remote)
	fmt.Printf("cluster gen : %f s\nplain read : %f MiB/s\nplain write : %f MiB/s\n", a, 4096/b, 4096/c)
}

// file io test
func test1(remote string) {
	a, b, c := kvault.Test_IO(remote)
	fmt.Printf("cluster login : %f s\nvfile read : %f MiB/s\nvfile write : %f MiB/s\n", a, 4096/b, 4096/c)
}

// multi file test
func test2(remote string) {
	a, b, c := kvault.Test_Multi(remote)
	fmt.Printf("complex gen : %f s\ncomplex read : %f s\ncomplex write : %f s\n", a, b, c)
}

// read eval print loop
func repl() {
	var k kvault.Shell
	flag := true
	for flag {
		temp := strings.Split(kio.Input(">>> "), " ")
		switch temp[0] {
		case "stop":
			flag = false
		case "status":
			fmt.Printf("%s  dir %d  file %d\n", k.CurPath, k.CurNum[0], k.CurNum[1])
			for i := 0; i < k.CurNum[0]+k.CurNum[1]; i++ {
				fmt.Printf("%s, (%s %t), %d B\n", k.CurName[i], k.CurTime[i], k.CurLock[i], k.CurSize[i])
			}
		case "cmd":
			fmt.Println(k.Command(temp[1], temp[2:]))
		case "load":
			switch temp[1] {
			case "flag":
				fmt.Printf("working %t, readonly %t, viewsize %t\n", k.FlagWk, k.FlagRo, k.FlagSz)
			case "err":
				fmt.Printf("asyncerr : %s\n", k.AsyncErr)
			case "str":
				for i, r := range k.IOstr {
					fmt.Printf("IOstr[%d] : %s\n", i, r)
				}
			case "byte":
				for i, r := range k.IObyte {
					if len(r) < 1000 {
						fmt.Printf("IObyte[%d] : %s\n", i, string(r))
					} else {
						fmt.Printf("IObyte[%d] : <size %d>\n", i, len(r))
					}
				}
			default:
				fmt.Println("flag err str byte")
			}
		case "store":
			switch temp[1] {
			case "viewsize":
				k.FlagSz = (temp[2] == "true")
			case "curpath":
				k.CurPath = temp[2]
			case "str":
				k.IOstr = append(k.IOstr, temp[2])
			case "byte":
				if temp[2] == "*" {
					k.IObyte = append(k.IObyte, kaes.Basickey())
				} else {
					k.IObyte = append(k.IObyte, []byte(temp[2]))
				}
			default:
				fmt.Println("viewsize curpath str byte")
			}
		default:
			fmt.Println("stop status cmd load store")
		}
	}
}

// mode 0 1 2 else
func main() {
	remote := "E:/t/" // ext drive ~/
	switch kio.Input("mode (0 1 2 else) >>> ") {
	case "0":
		test0(remote)
	case "1":
		test1(remote)
	case "2":
		test2(remote)
	default:
		repl()
	}
	time.Sleep(5 * time.Second)
	kio.Input("press ENTER to exit... ")
}
