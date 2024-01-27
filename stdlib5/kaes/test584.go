package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	kaes "example.com/kaesst"
)

func sequal(a []byte, b []byte) bool {
	if len(a) == len(b) {
		for i, r := range a {
			if b[i] != r {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func dels(todel [][]byte) {
	for i := range todel {
		todel[i] = nil
	}
}

func read(path string) []byte {
	f, _ := ioutil.ReadFile(path)
	return f
}

func write(path string, data []byte) {
	f, _ := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666) // write as f
	defer f.Close()
	f.Write(data)
}

func maintest() {
	_, err := os.ReadDir("./temp584")
	if err == nil {
		os.RemoveAll("./temp584")
	}
	os.Mkdir("./temp584", os.ModePerm)
	fmt.Println("===== gen5 KAES test =====")

	size := []int{0, 524272, 9352671, 211104879, 524288, 6291456, 16777216, 36175872}
	pw := kaes.Genrandom(48)
	kf := kaes.Genkf("nopath")
	hint := kaes.Genrandom(128)
	msg := "동성고 최고 귀요미 김병관"
	pre := make([][]byte, 8)
	enc := make([][]byte, 8)
	plain := make([][]byte, 8)
	for i := 0; i < 8; i++ {
		pre[i] = kaes.Genrandom(size[i])
	}

	k0 := kaes.Init0()
	k1 := kaes.Init1()
	k2 := kaes.Init2()
	k3 := kaes.Init3()
	k0.Msg = msg
	k1.Msg = msg
	fmt.Println("")

	for i := 0; i < 8; i++ {
		enc[i] = k0.En(pw, kf, hint, pre[i])
		a, b, c := k0.View(enc[i])
		fmt.Printf("test %d : hint_%t, msg_%t\n", i, sequal(a, hint), (b == msg))
		plain[i] = k0.De(pw, kf, enc[i], c)
		fmt.Printf("test %d : data_%t\n", i, sequal(pre[i], plain[i]))
	}
	dels(enc)
	dels(plain)
	fmt.Println("")

	for i := 0; i < 8; i++ {
		write("temp584\\t", pre[i])
		enc[i] = []byte(k1.En(pw, kf, hint, "temp584\\t"))
		a, b, c := k1.View(string(enc[i]))
		fmt.Printf("test %d : hint_%t, msg_%t\n", i, sequal(a, hint), (b == msg))
		plain[i] = []byte(k1.De(pw, kf, string(enc[i]), c))
		plain[i] = read(string(plain[i]))
		fmt.Printf("test %d : data_%t\n", i, sequal(pre[i], plain[i]))
	}
	dels(enc)
	dels(plain)
	fmt.Println("")

	for i := 0; i < 8; i++ {
		enc[i] = k2.En(pw, pre[i])
		plain[i] = k2.De(pw, enc[i])
		fmt.Printf("test %d : data_%t\n", i, sequal(pre[i], plain[i]))
	}
	dels(enc)
	dels(plain)
	fmt.Println("")

	for i := 0; i < 8; i++ {
		write("temp584/d", pre[i])
		k3.En(pw, "temp584/d", "temp584/e")
		k3.De(pw, "temp584/e", "temp584/d")
		plain[i] = read("temp584/d")
		fmt.Printf("test %d : data_%t\n", i, sequal(pre[i], plain[i]))
	}
}

func bench() {
	k0 := kaes.Init0()
	k1 := kaes.Init1()
	k2 := kaes.Init2()
	k3 := kaes.Init3()
	pw := kaes.Genrandom(48)
	kf := kaes.Genkf("nopath")
	hint := make([]byte, 0)
	fmt.Println("")

	time.Sleep(5 * time.Second)
	t0 := time.Now().UnixMicro()                   // micro seconds
	data := kaes.Genrandom(2 * 1024 * 1024 * 1024) // 2048 MiB
	t1 := time.Now().UnixMicro()
	t2 := float64(t1-t0) / 1000000
	fmt.Printf("rand : time %3f s, speed %3f MiB/s\n", t2, (2048 / t2))

	time.Sleep(5 * time.Second)
	t0 = time.Now().UnixMicro()
	temp := k0.En(pw, kf, hint, data)
	t1 = time.Now().UnixMicro()
	t2 = float64(t1-t0) / 1000000
	fmt.Printf("k0en : time %3f s, speed %3f MiB/s\n", t2, (2048 / t2))
	data = nil

	time.Sleep(5 * time.Second)
	_, _, c := k0.View(temp)
	t0 = time.Now().UnixMicro()
	k0.De(pw, kf, temp, c)
	t1 = time.Now().UnixMicro()
	t2 = float64(t1-t0) / 1000000
	fmt.Printf("k0de : time %3f s, speed %3f MiB/s\n", t2, (2048 / t2))
	temp = nil

	data = kaes.Genrandom(4 * 1024 * 1024 * 1024) // 4096 MiB
	time.Sleep(5 * time.Second)
	write("temp584\\b", data)
	data = nil
	t0 = time.Now().UnixMicro()
	temp = []byte(k1.En(pw, kf, hint, "temp584\\b"))
	t1 = time.Now().UnixMicro()
	t2 = float64(t1-t0) / 1000000
	fmt.Printf("k1en : time %3f s, speed %3f MiB/s\n", t2, (4096 / t2))

	time.Sleep(5 * time.Second)
	_, _, c = k1.View(string(temp))
	t0 = time.Now().UnixMicro()
	k1.De(pw, kf, string(temp), c)
	t1 = time.Now().UnixMicro()
	t2 = float64(t1-t0) / 1000000
	fmt.Printf("k1de : time %3f s, speed %3f MiB/s\n", t2, (4096 / t2))

	data = kaes.Genrandom(1 * 1024 * 1024 * 1024) // 1024 MiB
	time.Sleep(5 * time.Second)
	t0 = time.Now().UnixMicro()
	temp = k2.En(pw, data)
	t1 = time.Now().UnixMicro()
	t2 = float64(t1-t0) / 1000000
	fmt.Printf("k2en : time %3f s, speed %3f MiB/s\n", t2, (1024 / t2))
	data = nil

	time.Sleep(5 * time.Second)
	t0 = time.Now().UnixMicro()
	k2.De(pw, temp)
	t1 = time.Now().UnixMicro()
	t2 = float64(t1-t0) / 1000000
	fmt.Printf("k2de : time %3f s, speed %3f MiB/s\n", t2, (1024 / t2))
	temp = nil

	data = kaes.Genrandom(4 * 1024 * 1024 * 1024) // 4096 MiB
	time.Sleep(5 * time.Second)
	write("temp584/bd", data)
	data = nil
	t0 = time.Now().UnixMicro()
	k3.En(pw, "temp584/bd", "temp584/be")
	t1 = time.Now().UnixMicro()
	t2 = float64(t1-t0) / 1000000
	fmt.Printf("k3en : time %3f s, speed %3f MiB/s\n", t2, (4096 / t2))

	time.Sleep(5 * time.Second)
	t0 = time.Now().UnixMicro()
	k3.De(pw, "temp584/be", "temp584/bd")
	t1 = time.Now().UnixMicro()
	t2 = float64(t1-t0) / 1000000
	fmt.Printf("k3de : time %3f s, speed %3f MiB/s\n", t2, (4096 / t2))
}

func main() {
	maintest()
	bench()
}
