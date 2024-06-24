package main

import (
	"fmt"
	"os"
	"stdlib5/kaes"
	"stdlib5/kio"
	"stdlib5/kpic"
	"stdlib5/ksign"
	"stdlib5/kzip"
	"stdlib5/legsup"
	"time"
)

type g0 struct {
}

func (tbox *g0) read(path string) ([]byte, float64) {
	t0 := time.Now()
	f, _ := kio.Open(path, "r")
	defer f.Close()
	data, _ := kio.Read(f, -1)
	tpass := time.Since(t0).Seconds()
	time.Sleep(2 * time.Second)
	return data, tpass
}

func (tbox *g0) write(data []byte) float64 {
	t0 := time.Now()
	f, _ := kio.Open("./temp", "w")
	defer f.Close()
	kio.Write(f, data)
	tpass := time.Since(t0).Seconds()
	time.Sleep(2 * time.Second)
	return tpass
}

func (tbox *g0) calc(rept int) (float64, float64, float64) {
	var r0, r1, r2, r3 int
	t0 := time.Now()
	for i := 0; i < rept; i++ {
		r0 = r2 + r1
		r1 = r1 + 17
		r2 = r0 + r3
		r3 = r3 + 86
	}
	tpass0 := time.Since(t0).Seconds()
	time.Sleep(1 * time.Second)

	var v0, v1, v2, v3 float64
	t1 := time.Now()
	for i := 0; i < rept; i++ {
		v0 = v2 + v1
		v1 = v1 + 18.3617
		v2 = v0 + v3
		v3 = v3 + 24.9287
	}
	tpass1 := time.Since(t1).Seconds()
	time.Sleep(1 * time.Second)

	return tpass0, tpass1, float64(r0) + v0
}

func (tbox *g0) test() {
	data, t0 := tbox.read("../setB.bin")
	fmt.Printf("%f MiB/s\n", 1024/t0)
	t1 := tbox.write(data)
	fmt.Printf("%f MiB/s\n", 1024/t1)
	t0, t1, _ = tbox.calc(10000000)
	fmt.Printf("%f M/s %f M/s\n", 40/t0, 40/t1)
}

type g1 struct {
	w legsup.G1enc
}

func (tbox *g1) en() float64 {
	time.Sleep(2 * time.Second)
	tbox.w.Path = "../tcond/setA.bin"
	tbox.w.Pw = "0000"
	t0 := time.Now()
	fmt.Println(tbox.w.Encrypt())
	return time.Since(t0).Seconds()
}

func (tbox *g1) de() float64 {
	time.Sleep(2 * time.Second)
	tbox.w.Path = "../tcond/setA.bin.k"
	fmt.Println(tbox.w.View())
	tbox.w.Pw = "0000"
	t0 := time.Now()
	fmt.Println(tbox.w.Decrypt())
	return time.Since(t0).Seconds()
}

func (tbox *g1) test() {
	tbox.w.Init()
	t0 := tbox.en()
	t1 := tbox.de()
	fmt.Printf("%f MiB/s %f MiB/s\n", 300/t0, 300/t1)
}

type g2 struct {
	w legsup.G2enc
}

func (tbox *g2) en() float64 {
	time.Sleep(2 * time.Second)
	tbox.w.Path = "../tcond/setA.bin"
	tbox.w.Pw = "0000"
	t0 := time.Now()
	fmt.Println(tbox.w.Encrypt())
	return time.Since(t0).Seconds()
}

func (tbox *g2) de() float64 {
	time.Sleep(2 * time.Second)
	tbox.w.Path = "../tcond/setA.bin.k"
	fmt.Println(tbox.w.View())
	tbox.w.Pw = "0000"
	t0 := time.Now()
	fmt.Println(tbox.w.Decrypt())
	return time.Since(t0).Seconds()
}

func (tbox *g2) test() {
	tbox.w.Init()
	tbox.w.Hidename = false
	t0 := tbox.en()
	t1 := tbox.de()
	fmt.Printf("%f MiB/s %f MiB/s\n", 300/t0, 300/t1)
}

type g3 struct {
	w0 legsup.G3kzip
	w1 legsup.G3kaesall
	w2 legsup.G3kaesfunc
	w3 legsup.G3kpic
}

func (tbox *g3) kzip(path string) (float64, float64) {
	tbox.w0.Init()
	time.Sleep(2 * time.Second)
	t0 := time.Now()
	fmt.Println(tbox.w0.Packd(path, "../tcond/temp"))
	t1 := time.Since(t0).Seconds()

	fmt.Println(tbox.w0.View("../tcond/temp"))
	time.Sleep(2 * time.Second)
	t0 = time.Now()
	fmt.Println(tbox.w0.Unpack("../tcond/temp"))
	t2 := time.Since(t0).Seconds()
	return t1, t2
}

func (tbox *g3) kaes(path string) (float64, float64, float64, float64) {
	tbox.w1.Init(0, 0)
	tbox.w1.Hidename = false
	time.Sleep(2 * time.Second)
	t0 := time.Now()
	fmt.Println(tbox.w1.Encrypt(path, "0000", legsup.G3kf()))
	t1 := time.Since(t0).Seconds()
	npath := tbox.w1.Respath

	fmt.Println(tbox.w1.View(npath))
	time.Sleep(2 * time.Second)
	t0 = time.Now()
	fmt.Println(tbox.w1.Decrypt(npath, "0000", legsup.G3kf()))
	t2 := time.Since(t0).Seconds()

	time.Sleep(2 * time.Second)
	t0 = time.Now()
	fmt.Println(tbox.w2.Encrypt(path, path+".k", make([]byte, 32)))
	t3 := time.Since(t0).Seconds()

	time.Sleep(2 * time.Second)
	t0 = time.Now()
	fmt.Println(tbox.w2.Decrypt(path+".k", path, make([]byte, 32)))
	t4 := time.Since(t0).Seconds()

	return t1, t2, t3, t4
}

func (tbox *g3) kpng(path string) (float64, float64) {
	os.Mkdir("../tcond/tempp", os.ModePerm)
	tbox.w3.Init("", 2600, 2600)
	tbox.w3.Pcover = true
	time.Sleep(2 * time.Second)
	t0 := time.Now()
	fmt.Println(tbox.w3.Pack(path, "../tcond/tempp"))
	t1 := time.Since(t0).Seconds()

	name, num, _ := tbox.w3.Detect("../tcond/tempp")
	time.Sleep(2 * time.Second)
	t0 = time.Now()
	tbox.w3.Unpack("../tcond/temp", "../tcond/tempp", name, num)
	t2 := time.Since(t0).Seconds()

	return t1, t2
}

func (tbox *g3) test() {
	t0, t1 := tbox.kzip("../setE")
	fmt.Printf("%f MiB/s %f MiB/s\n", 10240/t0, 10240/t1)

	t0, t1, t2, t3 := tbox.kaes("../tcond/setC.bin")
	fmt.Printf("%f MiB/s %f MiB/s %f MiB/s %f MiB/s\n", 3072/t0, 3072/t1, 3072/t2, 3072/t3)

	t0, t1 = tbox.kpng("../tcond/setB.bin")
	fmt.Printf("%f MiB/s %f MiB/s\n", 1024/t0, 1024/t1)
}

type g4 struct {
	w0 legsup.G4enc
	w1 legsup.G4kaesall
	w2 legsup.G4kaesfunc
	pw []byte
	kf []byte
}

func (tbox *g4) kenc(path string) (float64, float64) {
	time.Sleep(2 * time.Second)
	t0 := time.Now()
	npath, _ := tbox.w0.Encrypt([]string{path}, []byte("0000"))
	t1 := time.Since(t0).Seconds()

	time.Sleep(12 * time.Second)
	t0 = time.Now()
	fmt.Println(tbox.w0.Decrypt(npath, []byte("0000")))
	t2 := time.Since(t0).Seconds()

	return t1, t2
}

func (tbox *g4) kaesF(path string) (float64, float64, float64, float64) {
	tbox.pw = []byte("0000")
	tbox.kf = legsup.G4kf()
	time.Sleep(12 * time.Second)
	t0 := time.Now()
	npath, _ := tbox.w1.EnFile(tbox.pw, tbox.kf, path)
	t1 := time.Since(t0).Seconds()

	fmt.Println(tbox.w1.ViewFile(npath))
	time.Sleep(12 * time.Second)
	t0 = time.Now()
	fmt.Println(tbox.w1.DeFile(tbox.pw, tbox.kf, npath))
	t2 := time.Since(t0).Seconds()

	time.Sleep(12 * time.Second)
	tbox.w2.Inbuf.OpenF(path, true)
	tbox.w2.Exbuf.OpenF(path+".k", false)
	t0 = time.Now()
	fmt.Println(tbox.w2.Encrypt(make([]byte, 48)))
	t3 := time.Since(t0).Seconds()
	tbox.w2.Inbuf.CloseF()
	tbox.w2.Exbuf.CloseF()

	time.Sleep(12 * time.Second)
	tbox.w2.Inbuf.OpenF(path+".k", true)
	tbox.w2.Exbuf.OpenF(path, false)
	t0 = time.Now()
	fmt.Println(tbox.w2.Decrypt(make([]byte, 48)))
	t4 := time.Since(t0).Seconds()
	tbox.w2.Inbuf.CloseF()
	tbox.w2.Exbuf.CloseF()

	return t1, t2, t3, t4
}

func (tbox *g4) kaesB(data []byte) (float64, float64, float64, float64) {
	tbox.pw = []byte("0000")
	tbox.kf = legsup.G4kf()
	time.Sleep(12 * time.Second)
	t0 := time.Now()
	ndata, _ := tbox.w1.EnBin(tbox.pw, tbox.kf, data)
	t1 := time.Since(t0).Seconds()
	data = nil

	fmt.Println(tbox.w1.ViewBin(ndata))
	time.Sleep(12 * time.Second)
	t0 = time.Now()
	data, _ = tbox.w1.DeBin(tbox.pw, tbox.kf, ndata)
	t2 := time.Since(t0).Seconds()
	ndata = nil

	time.Sleep(12 * time.Second)
	tbox.w2.Inbuf.OpenB(data, true)
	tbox.w2.Exbuf.OpenB(make([]byte, 0, len(data)+1000000), false)
	t0 = time.Now()
	fmt.Println(tbox.w2.Encrypt(make([]byte, 48)))
	t3 := time.Since(t0).Seconds()
	tbox.w2.Inbuf.CloseB()
	ndata = tbox.w2.Exbuf.CloseB()
	data = nil

	time.Sleep(12 * time.Second)
	tbox.w2.Inbuf.OpenB(ndata, true)
	tbox.w2.Exbuf.OpenB(make([]byte, 0, len(ndata)+1000000), false)
	t0 = time.Now()
	fmt.Println(tbox.w2.Decrypt(make([]byte, 48)))
	t4 := time.Since(t0).Seconds()
	tbox.w2.Inbuf.CloseB()
	tbox.w2.Exbuf.CloseB()

	return t1, t2, t3, t4
}

func (tbox *g4) test() {
	t0, t1 := tbox.kenc("../tcond/setC.bin")
	fmt.Printf("%f MiB/s %f MiB/s\n", 3072/t0, 3072/t1)

	t0, t1, t2, t3 := tbox.kaesF("../tcond/setD.bin")
	fmt.Printf("%f MiB/s %f MiB/s %f MiB/s %f MiB/s\n", 10240/t0, 10240/t1, 10240/t2, 10240/t3)

	f, _ := kio.Open("../tcond/setC.bin", "r")
	temp, _ := kio.Read(f, -1)
	f.Close()
	t0, t1, t2, t3 = tbox.kaesB(temp)
	fmt.Printf("%f MiB/s %f MiB/s %f MiB/s %f MiB/s\n", 3072/t0, 3072/t1, 3072/t2, 3072/t3)
}

type g5 struct {
	w0  kaes.Allmode
	w1  kaes.Funcmode
	pub string
	pri string
}

func (tbox *g5) kzip(path string) (float64, float64) {
	os.Mkdir("../tcond/temp/", os.ModePerm)
	time.Sleep(10 * time.Second)
	t0 := time.Now()
	fmt.Println(kzip.Dozip([]string{path}, "webp", "../tcond/temp.webp"))
	t1 := time.Since(t0).Seconds()

	time.Sleep(10 * time.Second)
	t0 = time.Now()
	fmt.Println(kzip.Unzip("../tcond/temp.webp", "../tcond/temp/", true))
	t2 := time.Since(t0).Seconds()

	os.Remove("../tcond/temp.webp")
	os.RemoveAll("../tcond/temp/")
	return t1, t2
}

func (tbox *g5) rgen(size int) float64 {
	time.Sleep(4 * time.Second)
	t0 := time.Now()
	kaes.Genrand(size)
	t1 := time.Since(t0).Seconds()
	return t1
}

func (tbox *g5) kaesF(path string) (float64, float64, float64, float64) {
	tbox.w0.Signkey[0] = tbox.pub
	tbox.w0.Signkey[1] = tbox.pri
	pw := []byte("0000")
	kf := kaes.Basickey()
	akey := make([]byte, 48)

	time.Sleep(14 * time.Second)
	t0 := time.Now()
	npath, _ := tbox.w0.EnFile(pw, kf, path, 0)
	t1 := time.Since(t0).Seconds()

	tbox.w0.ViewFile(npath)
	time.Sleep(14 * time.Second)
	t0 = time.Now()
	tbox.w0.DeFile(pw, kf, npath)
	t2 := time.Since(t0).Seconds()
	os.Remove(npath)

	tbox.w1.Before.Open(path, true)
	tbox.w1.After.Open(path+".k", false)
	time.Sleep(14 * time.Second)
	t0 = time.Now()
	tbox.w1.Encrypt(akey)
	t3 := time.Since(t0).Seconds()
	tbox.w1.Before.Close()
	tbox.w1.After.Close()
	akey = make([]byte, 48)

	tbox.w1.Before.Open(path+".k", true)
	tbox.w1.After.Open(path, false)
	time.Sleep(14 * time.Second)
	t0 = time.Now()
	tbox.w1.Decrypt(akey)
	t4 := time.Since(t0).Seconds()
	tbox.w1.Before.Close()
	tbox.w1.After.Close()

	os.Remove(path + ".k")
	return t1, t2, t3, t4
}

func (tbox *g5) kaesB(data []byte) (float64, float64, float64, float64) {
	tbox.w0.Signkey[0] = tbox.pub
	tbox.w0.Signkey[1] = tbox.pri
	pw := []byte("0000")
	kf := kaes.Basickey()
	akey := make([]byte, 48)

	time.Sleep(14 * time.Second)
	t0 := time.Now()
	ndata, _ := tbox.w0.EnBin(pw, kf, data, 0)
	t1 := time.Since(t0).Seconds()
	time.Sleep(4 * time.Second)
	data = nil

	tbox.w0.ViewBin(ndata)
	time.Sleep(14 * time.Second)
	t0 = time.Now()
	data, _ = tbox.w0.DeBin(pw, kf, ndata)
	t2 := time.Since(t0).Seconds()
	time.Sleep(4 * time.Second)
	ndata = nil

	tbox.w1.Before.Open(data, true)
	tbox.w1.After.Open(make([]byte, 0, len(data)+100000), false)
	data = nil
	time.Sleep(14 * time.Second)
	t0 = time.Now()
	tbox.w1.Encrypt(akey)
	t3 := time.Since(t0).Seconds()
	tbox.w1.Before.Close()
	data = tbox.w1.After.Close()
	akey = make([]byte, 48)

	tbox.w1.Before.Open(data, true)
	tbox.w1.After.Open(make([]byte, 0, len(data)+100000), false)
	data = nil
	time.Sleep(14 * time.Second)
	t0 = time.Now()
	tbox.w1.Decrypt(akey)
	t4 := time.Since(t0).Seconds()
	tbox.w1.Before.Close()
	tbox.w1.After.Close()

	return t1, t2, t3, t4
}

func (tbox *g5) kpic(path string, mode string) (float64, float64) {
	w, _ := kpic.Initpic("", 2600, 2600)
	w.Style = mode
	os.Mkdir("../tcond/temp/", os.ModePerm)
	w.Target = path
	w.Export = "../tcond/temp/"

	time.Sleep(14 * time.Second)
	t0 := time.Now()
	w.Pack(2)
	t1 := time.Since(t0).Seconds()

	w.Target = "../tcond/temp/"
	w.Export = "../tcond/tdata"
	name, num, _ := w.Detect()
	time.Sleep(14 * time.Second)
	t0 = time.Now()
	w.Unpack(name, num)
	t2 := time.Since(t0).Seconds()

	os.Remove("../tcond/tdata")
	os.RemoveAll("../tcond/temp/")
	return t1, t2
}

func (tbox *g5) test() {
	tbox.pub, tbox.pri, _ = ksign.Genkey(2048)
	t0, t1 := tbox.kzip("../tcond/setE")
	fmt.Printf("%f MiB/s %f MiB/s\n", 10240/t0, 10240/t1)

	time.Sleep(4 * time.Second)
	t0 = tbox.rgen(3 * 1024 * 1024 * 1024)
	fmt.Printf("%f MiB/s\n", 3072/t0)

	time.Sleep(4 * time.Second)
	t0, t1, t2, t3 := tbox.kaesF("../tcond/setD.bin")
	fmt.Printf("%f MiB/s %f MiB/s %f MiB/s %f MiB/s\n", 10240/t0, 10240/t1, 10240/t2, 10240/t3)

	f, _ := kio.Open("../tcond/setC.bin", "r")
	temp, _ := kio.Read(f, -1)
	f.Close()
	time.Sleep(4 * time.Second)
	t0, t1, t2, t3 = tbox.kaesB(temp)
	fmt.Printf("%f MiB/s %f MiB/s %f MiB/s %f MiB/s\n", 3072/t0, 3072/t1, 3072/t2, 3072/t3)

	t0, t1 = tbox.kpic("../tcond/setB.bin", "png")
	fmt.Printf("%f MiB/s %f MiB/s\n", 1024/t0, 1024/t1)
	t0, t1 = tbox.kpic("../tcond/setB.bin", "webp")
	fmt.Printf("%f MiB/s %f MiB/s\n", 1024/t0, 1024/t1)
}

func main() {
	var k g5
	k.test()
}
