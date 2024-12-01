// test718 : common.packager

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"stdlib5/cliget"
	"stdlib5/kcom"
	"stdlib5/kdb"
	"stdlib5/kio"
	"stdlib5/kobj"
	"stdlib5/kpkg"
	"stdlib5/ksc"
	"strconv"
	"strings"
)

// sign key
type signdata struct {
	name   string // sign name
	phash  []byte // public key hash
	public string // public key
}

// package data
type pkgdata struct {
	name     string  // pkg name
	ver      float64 // pkg version (N.MLL)
	text     string  // pkg explanation
	release  string  // release date
	download string  // download date
	path     string  // local install path
	num      int     // chunk num
}

// main struct
type mainclass struct {
	config_path [3]string // export, desktop, local
	config_url  [3]string // download, info, help
	config_dev  [2]bool   // iswin, devmode

	sign_db    []signdata // signkey DB
	install_db []pkgdata  // installed package DB
	package_db []pkgdata  // whole web package DB
}

// get status info
func (tbox *mainclass) get_status() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil

	f, er := kio.Open("../../_ST5_CONFIG.txt", "r")
	if er != nil {
		return er
	}
	d, _ := kio.Read(f, -1)
	f.Close()
	worker := kdb.Initkdb()
	worker.Read(string(d))
	tv, _ := worker.Get("path.export")
	tbox.config_path[0] = tv.Dat6
	tv, _ = worker.Get("path.desktop")
	tbox.config_path[1] = tv.Dat6
	tv, _ = worker.Get("path.local")
	tbox.config_path[2] = tv.Dat6
	tv, _ = worker.Get("url.download")
	tbox.config_url[0] = tv.Dat6
	tv, _ = worker.Get("url.info")
	tbox.config_url[1] = tv.Dat6
	tv, _ = worker.Get("url.help")
	tbox.config_url[2] = tv.Dat6
	tv, _ = worker.Get("dev.os")
	tbox.config_dev[0] = tv.Dat6 == "windows"
	tv, _ = worker.Get("dev.activate")
	tbox.config_dev[1] = tv.Dat1

	f, er = kio.Open("../../_ST5_SIGN.txt", "r")
	if er != nil {
		return er
	}
	d, _ = kio.Read(f, -1)
	f.Close()
	worker = kdb.Initkdb()
	worker.Read(string(d))
	num := 0
	for {
		if _, ext := worker.Name[fmt.Sprintf("%d.name", num)]; ext {
			num = num + 1
		} else {
			break
		}
	}
	tbox.sign_db = make([]signdata, 0)
	for i := 0; i < num; i++ {
		var temp signdata
		tv, _ = worker.Get(fmt.Sprintf("%d.name", i))
		temp.name = tv.Dat6
		tv, _ = worker.Get(fmt.Sprintf("%d.phash", i))
		temp.phash = tv.Dat5
		tv, _ = worker.Get(fmt.Sprintf("%d.public", i))
		temp.public = tv.Dat6
		if kio.Bequal(ksc.Crc32hash([]byte(temp.public)), temp.phash) {
			tbox.sign_db = append(tbox.sign_db, temp)
		} else {
			fmt.Printf("Warning : Sign %s with phash %s does not match to public key.\n", temp.name, kio.Bprint(temp.phash))
		}
	}

	fmt.Printf("Config : path\n<입출력 폴더> : %s\n<바탕화면> : %s\n<로컬 홈> : %s\n\n", tbox.config_path[0], tbox.config_path[1], tbox.config_path[2])
	fmt.Printf("Config : url\n<다운로드 저장소> : %s\n<패키지 저장소> : %s\n<도움말 페이지> : %s\n\n", tbox.config_url[0], tbox.config_url[1], tbox.config_url[2])
	fmt.Printf("Config : dev\n<윈도우 OS> : %t\n<개발자 모드> : %t\n\n", tbox.config_dev[0], tbox.config_dev[1])
	for _, r := range tbox.sign_db {
		fmt.Printf("Sign : %s (%s)\n%s\n\n", r.name, kio.Bprint(r.phash), r.public)
	}
	return err
}

// get_install_sub (path ~/)
func (tbox *mainclass) sub0(path string) {
	dirs, _ := os.ReadDir(path)
	for _, r := range dirs {
		name := r.Name()
		if name[len(name)-1] == '/' {
			name = name[:len(name)-1]
		}
		worker := kdb.Initkdb()
		f, err := kio.Open(path+name+"/_ST5_VERSION.txt", "r")
		if err == nil {
			d, _ := kio.Read(f, -1)
			f.Close()
			worker.Read(string(d))
		} else {
			worker.Read("name = \"None\"; version = 0.0; text = \"Invalid Package\"; release = \"00000000\"; download = \"00000000\"")
		}
		var temp pkgdata
		temp.name = name
		tv, _ := worker.Get("version")
		temp.ver = tv.Dat3
		tv, _ = worker.Get("text")
		temp.text = tv.Dat6
		tv, _ = worker.Get("release")
		temp.release = tv.Dat6
		tv, _ = worker.Get("download")
		temp.download = tv.Dat6
		temp.path = path + name + "/"
		temp.num = 1
		tbox.install_db = append(tbox.install_db, temp)
	}
}

// get installed package data
func (tbox *mainclass) get_install() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil

	tbox.install_db = make([]pkgdata, 0)
	tbox.sub0("../../_ST5_EXTENSION/")
	tbox.sub0("../../_ST5_COMMON/")
	for i, r := range tbox.install_db {
		fmt.Printf("[%03d] Package : %s ver%.3f\nInfo : %s\nRelease : %s, Download : %s\n\n", i, r.name, r.ver, r.text, r.release, r.download)
	}

	if num, er := strconv.Atoi(kio.Input("Select package to delete (x/N) : ")); er == nil {
		ans := kio.Input(fmt.Sprintf("Are you sure to delete package %s at %s? (y/n) : ", tbox.install_db[num].name, tbox.install_db[num].path))
		if ans == "y" || ans == "Y" || ans == "yes" {
			os.RemoveAll(tbox.install_db[num].path)
			fmt.Println("msg : deleted")
		} else {
			fmt.Println("msg : cancelled")
		}
	}
	return err
}

// get_package_sub (domain e/c, path ~/)
func (tbox *mainclass) sub1(domain string, path string) {
	worker := kdb.Initkdb()
	data, err := kcom.Gettxt(tbox.config_url[1], domain)
	if err == nil {
		worker.Read(data)
	} else {
		fmt.Printf("err : web fail (%s) should check status (mode 1) first\n", err)
		worker.Read("0.name = \"None\"; 0.version = 0.0; 0.devonly = False; 0.release = \"00000000\"; 0.text = \"Web access fail.\"; 0.num = 1;")
	}
	num := 0
	for {
		if _, ext := worker.Name[fmt.Sprintf("%d.name", num)]; ext {
			num = num + 1
		} else {
			break
		}
	}
	for i := 0; i < num; i++ {
		tv, _ := worker.Get(fmt.Sprintf("%d.devonly", i))
		if tbox.config_dev[1] || !tv.Dat1 {
			var temp pkgdata
			tv, _ = worker.Get(fmt.Sprintf("%d.name", i))
			temp.name = tv.Dat6
			tv, _ = worker.Get(fmt.Sprintf("%d.version", i))
			temp.ver = tv.Dat3
			tv, _ = worker.Get(fmt.Sprintf("%d.text", i))
			temp.text = tv.Dat6
			tv, _ = worker.Get(fmt.Sprintf("%d.release", i))
			temp.release = tv.Dat6
			tv, _ = worker.Get(fmt.Sprintf("%d.num", i))
			temp.num = tv.Dat2
			temp.download = "00000000"
			temp.path = path + temp.name + "/"
			tbox.package_db = append(tbox.package_db, temp)
		}
	}
}

// get whole package data
func (tbox *mainclass) get_package() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil

	tbox.package_db = make([]pkgdata, 0)
	tbox.sub1("extension", "../../_ST5_EXTENSION/")
	tbox.sub1("common", "../../_ST5_COMMON/")
	for i, r := range tbox.package_db {
		fmt.Printf("[%03d] Package : %s ver%.3f\nInfo : %s\nRelease : %s, Path : %s\n\n", i, r.name, r.ver, r.text, r.release, r.path)
	}

	if num, er := strconv.Atoi(kio.Input("Select package to download (x/N) : ")); er == nil {
		os.RemoveAll("./_ST5_DATA/")
		os.Mkdir("./_ST5_DATA/", os.ModePerm)
		count := 0.0
		tgt := tbox.package_db[num]
		kcom.Download(tbox.config_url[0], tgt.name, tgt.num, "./_ST5_DATA/download.webp", &count)
		fmt.Printf("msg : complete download %s at ./_ST5_DATA/download.webp\n", tgt.name)
	}
	return err
}

// get whole package data
func (tbox *mainclass) get_update() (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil
	if len(tbox.install_db) == 0 || len(tbox.package_db) == 0 {
		fmt.Println("msg : should check installed/web package (mode 2, 3) first")
	}
	db0 := make(map[string]float64)
	db1 := make(map[string]float64)
	for _, r := range tbox.install_db {
		db0[r.name] = r.ver
	}
	for _, r := range tbox.package_db {
		db1[r.name] = r.ver
	}
	for i, r := range db0 {
		if l, ext := db1[i]; ext {
			if l-r > 0.00001 {
				fmt.Printf("Package : %s ver%.3f -> ver%.3f\n", i, r, l)
			}
		}
	}
	return err
}

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

// div pack/unpack
func divf(pack bool) (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil

	os.RemoveAll("./_ST5_DATA/")
	os.Mkdir("./_ST5_DATA/", os.ModePerm)
	var sel cliget.OptSel
	if pack {
		sel.Init([]string{"< Select File to div pack >", "target file", "div size (MiB)"}, []string{"bool", "file", "int"}, []string{"T", "webp", "+"}, *getpath(), nil)
		sel.StrRes[2] = "24"
	} else {
		sel.Init([]string{"< Select File to div unpack >", "target file"}, []string{"bool", "file"}, []string{"T", "0"}, *getpath(), nil)
	}
	sel.GetOpt()

	if pack {
		tgt := sel.StrRes[1]
		dsize, _ := strconv.Atoi(sel.StrRes[2])
		if dsize < 1 {
			dsize = 1
		}
		dsize = dsize * 1048576
		fsize := kio.Size(tgt)
		if fsize < 0 {
			return errors.New("invalid file selection to div pack")
		}
		name := tgt[strings.LastIndex(tgt, "/")+1 : strings.LastIndex(tgt, ".")]
		f, _ := kio.Open(tgt, "r")
		defer f.Close()
		for i := 0; i < fsize/dsize; i++ {
			t, _ := kio.Open(fmt.Sprintf("./_ST5_DATA/%s.%d", name, i), "w")
			temp, _ := kio.Read(f, dsize)
			kio.Write(t, temp)
			t.Close()
		}
		if fsize%dsize != 0 {
			t, _ := kio.Open(fmt.Sprintf("./_ST5_DATA/%s.%d", name, fsize/dsize), "w")
			temp, _ := kio.Read(f, -1)
			kio.Write(t, temp)
			t.Close()
		}
		if fsize == 0 {
			t, _ := kio.Open(fmt.Sprintf("./_ST5_DATA/%s.0", name), "w")
			t.Close()
		}
		fmt.Println("msg : div pack at ./_ST5_DATA/ complete!")

	} else {
		tgt := sel.StrRes[1]
		upper := tgt[:strings.LastIndex(tgt, ".")]
		name := upper[strings.LastIndex(upper, "/")+1:]
		num := 0
		for {
			if _, er := os.Stat(fmt.Sprintf("%s.%d", upper, num)); er == nil {
				num = num + 1
			} else {
				break
			}
		}
		f, _ := kio.Open("./_ST5_DATA/"+name+".webp", "w")
		defer f.Close()
		for i := 0; i < num; i++ {
			t, _ := kio.Open(fmt.Sprintf("%s.%d", upper, i), "r")
			temp, _ := kio.Read(t, -1)
			kio.Write(f, temp)
			t.Close()
		}
		fmt.Println("msg : div unpack at ./_ST5_DATA/ complete!")
	}
	return err
}

// kpkg pack/unpack
func (tbox *mainclass) pkgf(pack bool) (err error) {
	defer func() {
		if ferr := recover(); ferr != nil {
			err = fmt.Errorf("critical : %s", ferr)
		}
	}()
	err = nil

	os.RemoveAll("./_ST5_DATA/")
	os.Mkdir("./_ST5_DATA/", os.ModePerm)
	defer os.RemoveAll("./temp675/")
	var sel cliget.OptSel
	if pack {
		sel.Init([]string{"< Select Folder & OSnum >", "tgt 0", "osnum 0", "tgt 1", "osnum 1", "tgt 2", "osnum 2", "tgt 3", "osnum 3"},
			[]string{"bool", "folder", "int", "folder", "int", "folder", "int", "folder", "int"}, []string{"T", "NR", "0+", "NR", "0+", "NR", "0+", "NR", "0+"}, *getpath(), nil)
		sel.StrRes = []string{"", "", "-1", "", "-1", "", "-1", "", "-1"}
	} else {
		sel.Init([]string{"< Select Package File >", "package file"}, []string{"bool", "file"}, []string{"T", "webp"}, *getpath(), nil)
	}
	sel.GetOpt()

	if pack {
		tgt_db := make([]string, 0)
		osnum_db := make([]int, 0)
		for i := 0; i < 4; i++ {
			if temp, err := strconv.Atoi(sel.StrRes[2*i+2]); err == nil && sel.StrRes[2*i+2] != "-1" {
				tgt_db = append(tgt_db, sel.StrRes[2*i+1])
				osnum_db = append(osnum_db, temp)
			}
		}
		osnum := 0
		if tbox.config_dev[0] {
			osnum = 1
		} else {
			osnum = 2
		}
		worker := kpkg.Initkpkg(osnum)
		sel.Init([]string{"< Set Package Field >", "name", "version", "text", "signkey"},
			[]string{"bool", "string", "float", "string", "keyfile"}, []string{"T", "1+", "0+", "1+", "0+"}, *getpath(), nil)
		sel.GetOpt()
		worker.Name = sel.StrRes[1]
		worker.Version, _ = strconv.ParseFloat(sel.StrRes[2], 64)
		worker.Text = sel.StrRes[3]
		worker.Private = ""
		worker.Public = ""
		if len(sel.ByteRes[4]) == 0 {
			fmt.Println("msg : no sign")
		} else {
			temp := kdb.Initkdb()
			er := temp.Read(string(sel.ByteRes[4]))
			if er == nil {
				tv, _ := temp.Get("private")
				worker.Private = tv.Dat6
				tv, _ = temp.Get("public")
				worker.Public = tv.Dat6
				tv, _ = temp.Get("phash")
				fmt.Printf("msg : sign phash %s\n", kio.Bprint(tv.Dat5))
			} else {
				fmt.Println("err : cannot read signkey")
			}
		}
		err = worker.Pack(osnum_db, tgt_db, fmt.Sprintf("./_ST5_DATA/%s.webp", worker.Name))
		fmt.Println("msg : kpkg pack at ./_ST5_DATA/ complete!")

	} else {
		osnum := 0
		if tbox.config_dev[0] {
			osnum = 1
		} else {
			osnum = 2
		}
		worker := kpkg.Initkpkg(osnum)
		path, er := worker.Unpack(sel.StrRes[1])
		if path[len(path)-1] != '/' {
			path = path + "/"
		}
		if !tbox.config_dev[0] {
			temp, _ := os.ReadDir(path)
			for _, r := range temp {
				if !r.IsDir() {
					os.Chmod(path+r.Name(), 0o755)
				}
			}
		}
		os.Rename(path, "./_ST5_DATA/"+worker.Name)
		err = er
		if worker.Public == "" {
			fmt.Println("warning : This package does not have ksign data!")
		} else {
			flag := true
			for _, r := range tbox.sign_db {
				if r.public == worker.Public {
					flag = false
					break
				}
			}
			if flag {
				fmt.Println("warning : Starter5 signDB does not contains the sign of this package!")
			}
		}
		fmt.Println("msg : kpkg unpack at ./_ST5_DATA/ complete!")
	}
	return err
}

func main() {
	kobj.Repath()
	flag := true
	var worker mainclass
	for flag {
		fmt.Printf("\n%24s   (0) %20s\n(1) %20s   (2) %20s\n(3) %20s   (4) %20s\n(5) %20s   (6) %20s\n(7) %20s   (8) %20s\n",
			"< Mode Selection >", "Exit", "View Status", "View Installed", "View Web Package", "View Update Queue",
			"Pack div", "Unpack div", "Pack kpkg", "Unpack kpkg")
		mode := kio.Input(">>> ")
		switch mode {
		case "0", "0\r", "0\n", "0\r\n": // exit
			flag = false
		case "1": // view status
			fmt.Println("===== Status =====")
			if err := worker.get_status(); err != nil {
				fmt.Println(err)
			}
		case "2": // view installed
			fmt.Println("===== Installed =====")
			if err := worker.get_install(); err != nil {
				fmt.Println(err)
			}
		case "3": // view web package
			fmt.Println("===== All Packages =====")
			if err := worker.get_package(); err != nil {
				fmt.Println(err)
			}
		case "4": // view update queue
			fmt.Println("===== Update Queue =====")
			if err := worker.get_update(); err != nil {
				fmt.Println(err)
			}
		case "5": // div pack
			fmt.Println("===== Pack div =====")
			if err := divf(true); err != nil {
				fmt.Println(err)
			}
		case "6": // div unpack
			fmt.Println("===== Unpack div =====")
			if err := divf(false); err != nil {
				fmt.Println(err)
			}
		case "7": // kpkg pack
			fmt.Println("===== Pack kpkg =====")
			if err := worker.pkgf(true); err != nil {
				fmt.Println(err)
			}
		case "8": // kpkg unpack
			fmt.Println("===== Unpack kpkg =====")
			if err := worker.pkgf(false); err != nil {
				fmt.Println(err)
			}
		}
	}
	kio.Input("press ENTER to exit... ")
}
