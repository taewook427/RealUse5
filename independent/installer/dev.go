// test744 : independent.installer dev

package main

import (
	"fmt"
	"os"
	"stdlib5/kcom"
	"stdlib5/kdb"
	"stdlib5/kio"
	"stdlib5/kobj"
	"stdlib5/kpkg"
	"stdlib5/ksc"
	"strconv"
	"strings"
)

// unpack package
func func0(path string) {
	wk := kpkg.Initkpkg(osnum)
	pos, err := wk.Unpack(path)
	name := path[strings.LastIndex(path, "/")+1 : strings.LastIndex(path, "_")]
	if err == nil {
		os.Rename(pos, fmt.Sprintf("./%s/", name))
		os.RemoveAll("./temp675/")
		fmt.Println(name + " : phash=" + kio.Bprint(ksc.Crc32hash([]byte(wk.Public))))
	} else {
		fmt.Println(err)
	}
}

// download content
func func1() {
	ext_pkg, ext_err := get_pkg("extension")
	com_pkg, com_err := get_pkg("common")
	ind_pkg, ind_err := get_pkg("independent")
	num := 0
	var proc float64
	if ext_err != nil || com_err != nil || ind_err != nil {
		fmt.Printf("web error : %s %s %s\n", ext_err, com_err, ind_err)
	} else {
		for _, r := range ext_pkg {
			fmt.Printf("%2d %15s  ext@%.4f %s  %s\n", num, r.name, r.version, r.release, r.text)
			num = num + 1
		}
		for _, r := range com_pkg {
			fmt.Printf("%2d %15s  com@%.4f %s  %s\n", num, r.name, r.version, r.release, r.text)
			num = num + 1
		}
		for _, r := range ind_pkg {
			fmt.Printf("%2d %15s  ind@%.4f %s  %s\n", num, r.name, r.version, r.release, r.text)
			num = num + 1
		}
	}
	for {
		n, _ := strconv.Atoi(kio.Input("Download Num [-1, N] : "))
		if n < 0 {
			break
		} else if n < len(ext_pkg) {
			if err := kcom.Download(url, ext_pkg[n].name, ext_pkg[n].num, fmt.Sprintf("./%s_dwn.webp", ext_pkg[n].name), &proc); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("download : ./%s_dwn.webp\n", ext_pkg[n].name)
			}
		} else if n-len(ext_pkg) < len(com_pkg) {
			n = n - len(ext_pkg)
			if err := kcom.Download(url, com_pkg[n].name, com_pkg[n].num, fmt.Sprintf("./%s_dwn.webp", com_pkg[n].name), &proc); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("download : ./%s_dwn.webp\n", com_pkg[n].name)
			}
		} else if n-len(ext_pkg)-len(com_pkg) < len(ind_pkg) {
			n = n - len(ext_pkg) - len(com_pkg)
			if err := kcom.Download(url, ind_pkg[n].name, ind_pkg[n].num, fmt.Sprintf("./%s_dwn.webp", ind_pkg[n].name), &proc); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("download : ./%s_dwn.webp\n", ind_pkg[n].name)
			}
		} else {
			break
		}
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("critical : %s\n", err)
		}
		kio.Input("press ENTER to exit... ")
	}()
	order = kobj.Repath()
	order = append(append(order, ""), "")
	fmt.Printf("URL : %s, OSnum : %d\n", url, osnum)
	switch order[1] {
	case "unpack":
		func0(kio.Abs(order[2]))
	case "download":
		if order[2] != "" {
			url = order[2]
		}
		func1()
	default:
		fmt.Println("invalid argument : use following commands\n./~ unpack [path]\n./~ download [url]")
	}
}

// !!! base info section !!!
var order []string
var url string = "https://taewook427.github.io/RealUse5_Sub0/data/"
var osnum int = 1

type pkg_db struct {
	name    string
	version float64
	release string
	text    string
	num     int
}

func get_pkg(domain string) ([]pkg_db, error) {
	wk := kdb.Initkdb()
	data, err := kcom.Gettxt(url+"list.html", domain)
	if err != nil {
		return nil, err
	}
	err = wk.Read(data)
	if err != nil {
		return nil, err
	}
	num := 0
	for {
		if _, ext := wk.Name[fmt.Sprintf("%d.name", num)]; ext {
			num = num + 1
		} else {
			break
		}
	}
	out := make([]pkg_db, num)
	for i := 0; i < num; i++ {
		tv, _ := wk.Get(fmt.Sprintf("%d.name", i))
		out[i].name = tv.Dat6
		tv, _ = wk.Get(fmt.Sprintf("%d.version", i))
		out[i].version = tv.Dat3
		tv, _ = wk.Get(fmt.Sprintf("%d.release", i))
		out[i].release = tv.Dat6
		tv, _ = wk.Get(fmt.Sprintf("%d.text", i))
		out[i].text = tv.Dat6
		tv, _ = wk.Get(fmt.Sprintf("%d.num", i))
		out[i].num = tv.Dat2
	}
	return out, nil
}
