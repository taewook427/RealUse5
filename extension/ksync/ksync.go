// test730 : extension.ksync

package main

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"stdlib5/kio"
	"stdlib5/kobj"
	"time"
)

// change file time (path ~/, t set time, m T : rand)
func func1(p string, t time.Time, m bool) {
	fs, _ := os.ReadDir(p)
	var nm string
	var ts time.Time
	for _, r := range fs {
		if m {
			ts = time.Unix(time.Now().Unix()-rand.Int64N(94608000), 0)
		} else {
			ts = t
		}
		nm = p + r.Name()
		if r.IsDir() {
			if nm[len(nm)-1] != '/' {
				nm = nm + "/"
			}
			func1(nm, ts, m)
		} else {
			os.Chtimes(nm, ts, ts)
			fmt.Printf("[change file time] %s\n", nm)
		}
	}
	os.Chtimes(p, t, t)
	fmt.Printf("[change dir time] %s\n", p)
	time.Sleep(50 * time.Millisecond)
}

// sync cluster (src -> dst), src & dst : ~/, mode : "n", "s", "t", "c"
func func2(src string, dst string, mode string) {
	src_fs, _ := os.ReadDir(src)
	src_dir := make([]string, 0)
	src_file := make([]string, 0)
	for _, r := range src_fs {
		nm := r.Name()
		if r.IsDir() {
			if nm[len(nm)-1] != '/' {
				nm = nm + "/"
				src_dir = append(src_dir, nm)
			}
		} else {
			src_file = append(src_file, nm)
		}
	}

	dst_fs, _ := os.ReadDir(dst)
	dst_dir := make([]string, 0)
	dst_file := make([]string, 0)
	for _, r := range dst_fs {
		nm := r.Name()
		if r.IsDir() {
			if nm[len(nm)-1] != '/' {
				nm = nm + "/"
				dst_dir = append(dst_dir, nm)
			}
		} else {
			dst_file = append(dst_file, nm)
		}
	}

	for _, r := range dst_file { // delete dst-only file
		if !slices.Contains(src_file, r) {
			os.Remove(dst + r)
			fmt.Printf("[delete file] _ -> %s\n", dst+r)
		}
	}
	for _, r := range src_file { // check & copy src_file
		if slices.Contains(dst_file, r) {
			if checkfile(src+r, dst+r, mode) {
				fmt.Printf("[same file] %s -> %s\n", src+r, dst+r)
			} else {
				fmt.Print("[update] ")
				copyfile(src+r, dst+r)
			}
		} else {
			copyfile(src+r, dst+r)
		}
	}

	for _, r := range dst_dir { // delete dst-only dir
		if !slices.Contains(src_dir, r) {
			os.RemoveAll(dst + r)
			fmt.Printf("[delete dir] _ -> %s\n", dst+r)
		}
	}
	for _, r := range src_dir { // check & copy src_dir
		if slices.Contains(dst_dir, r) {
			func2(src+r, dst+r, mode)
		} else {
			copydir(src+r, dst+r)
		}
	}
	fmt.Printf("[update dir] %s -> %s\n", src, dst)
}

// check if file is same, mode : "n", "s", "t", "c", no status print
func checkfile(src string, dst string, mode string) bool {
	res := false
	switch mode {
	case "n": // by name
		res = true
	case "s": // by size
		res = kio.Size(src) == kio.Size(dst)
	case "t": // by time
		ta, _ := os.Stat(src)
		tb, _ := os.Stat(dst)
		res = ta.ModTime().Unix() == tb.ModTime().Unix()
	case "c": // by content
		sz := kio.Size(src)
		if sz == kio.Size(dst) {
			res = true
			f, _ := kio.Open(src, "r")
			t, _ := kio.Open(dst, "r")
			defer f.Close()
			defer t.Close()
			for i := 0; i < sz/104857600; i++ {
				d0, _ := kio.Read(f, 104857600)
				d1, _ := kio.Read(t, 104857600)
				if !bytes.Equal(d0, d1) {
					res = false
					break
				}
			}
			if res && (sz%104857600 != 0) {
				d0, _ := kio.Read(f, sz%104857600)
				d1, _ := kio.Read(t, sz%104857600)
				if !bytes.Equal(d0, d1) {
					res = false
				}
			}
		}
	}
	return res
}

// copy file src -> dst
func copyfile(src string, dst string) {
	defer time.Sleep(10 * time.Millisecond)
	sz := kio.Size(src)
	f, _ := kio.Open(src, "r")
	t, _ := kio.Open(dst, "w")
	var d []byte
	defer f.Close()
	defer t.Close()
	for i := 0; i < sz/104857600; i++ {
		d, _ = kio.Read(f, 104857600)
		kio.Write(t, d)
	}
	if sz%104857600 != 0 {
		d, _ = kio.Read(f, sz%104857600)
		kio.Write(t, d)
	}
	fmt.Printf("[copy file] %s -> %s\n", src, dst)
}

// copy dir src -> dst (path ~/)
func copydir(src string, dst string) {
	os.Mkdir(dst, os.ModePerm)
	fs, _ := os.ReadDir(src)
	for _, r := range fs {
		nm := r.Name()
		if r.IsDir() {
			if nm[len(nm)-1] != '/' {
				nm = nm + "/"
			}
			copydir(src+nm, dst+nm)
		} else {
			copyfile(src+nm, dst+nm)
		}
	}
	fmt.Printf("[copy dir] %s -> %s\n", src, dst)
}

// patch cluster (src -> dst), src & dst : ~/
func func3(src string, dst string) {
	src_fs, _ := os.ReadDir(src)
	src_dir := make([]string, 0)
	src_file := make([]string, 0)
	for _, r := range src_fs {
		nm := r.Name()
		if r.IsDir() {
			if nm[len(nm)-1] != '/' {
				nm = nm + "/"
				src_dir = append(src_dir, nm)
			}
		} else {
			src_file = append(src_file, nm)
		}
	}

	dst_fs, _ := os.ReadDir(dst)
	dst_dir := make([]string, 0)
	dst_file := make([]string, 0)
	for _, r := range dst_fs {
		nm := r.Name()
		if r.IsDir() {
			if nm[len(nm)-1] != '/' {
				nm = nm + "/"
				dst_dir = append(dst_dir, nm)
			}
		} else {
			dst_file = append(dst_file, nm)
		}
	}

	for _, r := range dst_file { // delete dst-only file
		if !slices.Contains(src_file, r) {
			os.Remove(dst + r)
			fmt.Printf("[delete file] _ -> %s\n", dst+r)
		}
	}
	for _, r := range src_file { // check & copy src_file
		if slices.Contains(dst_file, r) {
			patchfile(src+r, dst+r)
		} else {
			copyfile(src+r, dst+r)
		}
	}

	for _, r := range dst_dir { // delete dst-only dir
		if !slices.Contains(src_dir, r) {
			os.RemoveAll(dst + r)
			fmt.Printf("[delete dir] _ -> %s\n", dst+r)
		}
	}
	for _, r := range src_dir { // check & copy src_dir
		if slices.Contains(dst_dir, r) {
			func3(src+r, dst+r)
		} else {
			copydir(src+r, dst+r)
		}
	}
	fmt.Printf("[patch dir] %s -> %s\n", src, dst)
}

// patch & update file (src -> dst)
func patchfile(src string, dst string) {
	defer time.Sleep(10 * time.Millisecond)
	flag := false
	sz := kio.Size(src)
	if sz != kio.Size(dst) {
		flag = true
		os.Truncate(dst, int64(sz))
	}
	f, _ := os.Open(src)
	t, _ := os.OpenFile(dst, os.O_RDWR, 0644)
	defer f.Close()
	defer t.Close()
	for i := 0; i < sz/1048576; i++ {
		d0, _ := kio.Read(f, 1048576)
		d1, _ := kio.Read(t, 1048576)
		if !bytes.Equal(d0, d1) {
			flag = true
			t.Seek(-1048576, 1)
			kio.Write(t, d0)
		}
	}
	if sz%1048576 != 0 {
		d0, _ := kio.Read(f, sz%1048576)
		d1, _ := kio.Read(t, sz%1048576)
		if !bytes.Equal(d0, d1) {
			flag = true
			t.Seek(-int64(sz%1048576), 1)
			kio.Write(t, d0)
		}
	}
	if flag {
		fmt.Printf("[patch file] %s -> %s\n", src, dst)
	} else {
		fmt.Printf("[same file] %s -> %s\n", src, dst)
	}
}

// get dir path
func func4(msg string) string {
	temp := kio.Input(msg)
	if temp[0] == '"' {
		temp = temp[1 : len(temp)-1]
	}
	temp = kio.Abs(temp)
	fmt.Printf("selected : %s\n", temp)
	return temp
}

// get time
func func5(msg string) time.Time {
	temp := kio.Input(msg)
	t, e := time.Parse("2006.01.02", temp)
	if e == nil {
		fmt.Printf("selected : %d\n", t.Unix())
		return t
	} else {
		t = time.Now()
		fmt.Printf("error %s : %d\n", e, t.Unix())
		return t
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("critical : %s\n", err)
		}
		kio.Input("\npress ENTER to exit... ")
	}()
	kobj.Repath()
	fmt.Printf("%24s   [0] %20s   [1] %20s\n[2] %20s   [3] %20s   [4] %20s\n[5] %20s   [6] %20s   [7] %20s\n",
		"< Select Mode >", "Cluster Sync (name)", "Cluster Sync (size)",
		"Cluster Sync (time)", "Cluster Sync (data)", "Sync by Patch",
		"Change Time (Now)", "Change Time (Rand)", "Change Time (Preset)")
	switch kio.Input(">>> ") {
	case "0":
		fmt.Println("===== Cluster Sync by Name =====")
		func2(func4("src path : "), func4("dst path : "), "n")
	case "1":
		fmt.Println("===== Cluster Sync by Size =====")
		func2(func4("src path : "), func4("dst path : "), "s")
	case "2":
		fmt.Println("===== Cluster Sync by ModTime =====")
		func2(func4("src path : "), func4("dst path : "), "t")
	case "3":
		fmt.Println("===== Cluster Sync by Content =====")
		func2(func4("src path : "), func4("dst path : "), "c")
	case "4":
		fmt.Println("===== Cluster Patch =====")
		func3(func4("src path : "), func4("dst path : "))
	case "5":
		fmt.Println("===== Change Time to Now =====")
		func1(func4("tgt path : "), time.Now(), false)
	case "6":
		fmt.Println("===== Change Time to Rand 3yr =====")
		func1(func4("tgt path : "), time.Now(), true)
	case "7":
		fmt.Println("===== Change Time to Preset Time =====")
		func1(func4("tgt path : "), func5("setting time (yyyy.mm.dd) : "), false)
	default:
		fmt.Println("invalid option")
	}
}
