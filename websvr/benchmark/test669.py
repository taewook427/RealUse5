import time
import os
import shutil

# import mung2
# import kerbal
# import nox2

# import kaes4py as en4py
# import kaes4hy as en4hy

import kzip
import kaesst
import kaeshy
import ksign
import kpic

class g0:
    def read(self, path):
        t0 = time.time()
        with open(path, "rb") as f:
            data = f.read()
        t1 = time.time()
        print(f"G0 read {t1 - t0}")
        time.sleep(2)
        return data, t1 - t0

    def write(self, data):
        t0 = time.time()
        with open("./temp", "wb") as f:
            f.write(data)
        t1 = time.time()
        print(f"G0 write {t1 - t0}")
        time.sleep(2)
        return t1 - t0

    def calc(self):
        t0 = time.time()
        r0, r1, r2, r3 = 0, 0, 0, 0
        for i in range(0, 10000000):
            r0 = r2 + r1
            r1 = r1 + 17
            r2 = r0 + r3
            r3 = r3 + 86
        t1 = time.time()
        time.sleep(1)
        t2 = time.time()
        r0, r1, r2, r3 = 0.0, 0.0, 0.0, 0.0
        for i in range(0, 10000000):
            r0 = r2 + r1
            r1 = r1 + 18.3617
            r2 = r0 + r3
            r3 = r3 + 24.9281
        t3 = time.time()
        time.sleep(1)
        print(f"G0 calc {t1 - t0} {t3 - t2}")
        return t1 - t0, t3 - t2

    def test(self):
        data0, t0 = self.read("../setB.bin")
        print(f"{1024 / t0} MiB/s")

        t0 = self.write(data0)
        print(f"{1024 / t0} MiB/s")

        t0, t1 = self.calc()
        print(f"{40 / t0} M/s {40 / t1} M/s")

        os.remove("./temp")

class g3:
    def __init__(self):
        self.w0 = mung2.toolbox()
        self.w1 = kerbal.toolbox()
        self.w2 = nox2.toolbox()
        self.d0 = kerbal.xenon()
        self.w2.set(32, "p0.png")

    def kzip(self, path):
        path = os.path.abspath(path)
        t0 = time.time()
        self.w0.pack(path, "./temp", True)
        t1 = time.time()
        time.sleep(2)

        t2 = time.time()
        self.w0.unpack("./temp")
        t3 = time.time()
        time.sleep(2)

        print(f"G3 kzip {t1 - t0} {t3 - t2}")
        return t1 - t0, t3 - t2

    def kaes(self, path):
        t0 = time.time()
        self.w1.encryptall(path, "0000", "hint", b"01010123", 3, 32, 131072, self.d0.data1)
        t1 = time.time()
        time.sleep(2)

        _, parms = self.w1.checkall(path + ".k")
        _, ckey, iv, name = self.w1.pwall( "0000", self.d0.data1, parms[4], parms[5], [ parms[3], parms[7], parms[8] ] )
        t2 = time.time()
        self.w1.decryptall( [path + ".k", parms[0], name, 32, 131072, ckey, iv] )
        t3 = time.time()
        time.sleep(2)

        t4 = time.time()
        self.w1.encryptfunc(path, path + ".k", b"\x00" * 32, 32, 131072, b"01010123")
        t5 = time.time()
        time.sleep(2)

        _, _, _, iv, ckdt, fstart = self.w1.checkfunc(path + ".k")
        t6 = time.time()
        self.w1.decryptfunc( path + ".k", path, b"\x00" * 32, [32, 131072, iv, ckdt, fstart] )
        t7 = time.time()
        time.sleep(2)

        print(f"G3 kaes {t1 - t0} {t3 - t2} {t5 - t4} {t7 - t6}")
        return t1 - t0, t3 - t2, t5 - t4, t7 - t6

    def kpng(self, path):
        t0 = time.time()
        self.w2.pack(path, False, "png")
        t1 = time.time()
        time.sleep(2)

        os.rename("./temp270", "./tempp")
        parms = self.w2.detect("./tempp")
        t2 = time.time()
        self.w2.unpack( ["./tempp"] + parms )
        t3 = time.time()
        time.sleep(2)

        print(f"G3 kpng {t1 - t0} {t3 - t2}")
        return t1 - t0, t3 - t2

    def test(self):
        t0, t1 = self.kzip("../setE/")
        print(f"{10240 / t0} MiB/s {10240 / t1} MiB/s")

        t0, t1, t2, t3 = self.kaes("../tcond/setC.bin")
        print(f"{3072 / t0} MiB/s {3072 / t1} MiB/s {3072 / t2} MiB/s {3072 / t3} MiB/s")

        t0, t1 = self.kpng("../tcond/setB.bin")
        print(f"{1024 / t0} MiB/s {1024 / t1} MiB/s")

        os.remove("./temp")
        shutil.rmtree("./temp261/")

class g4:
    def __init__(self):
        self.w0 = en4py.toolbox()
        self.w1 = en4hy.toolbox()
        self.pw, self.kf = b"0000", en4hy.genbkf()

    def kaes(self, worker, tgt):
        t0 = time.time()
        ntgt = worker.enwhole(self.pw, self.kf, b"", tgt)
        t1 = time.time()
        del tgt
        time.sleep(12)

        worker.view(ntgt)
        t2 = time.time()
        tgt = worker.dewhole(self.pw, self.kf, ntgt)
        t3 = time.time()
        del ntgt
        time.sleep(10)

        t4 = time.time()
        ntgt = worker.enfunc(b"\x00" * 48, tgt, tgt="temp")
        t5 = time.time()
        del tgt
        time.sleep(12)

        t6 = time.time()
        tgt = worker.defunc(b"\x00" * 48, ntgt)
        t7 = time.time()
        del ntgt
        time.sleep(10)

        print(f"G4 kaes {type(tgt)} {t1 - t0} {t3 - t2} {t5 - t4} {t7 - t6}")
        return t1 - t0, t3 - t2, t5 - t4, t7 - t6

    def rgen(self, size):
        t0 = time.time()
        temp = self.w0.genrandom(size)
        t1 = time.time()
        time.sleep(4)
        del temp

        t2 = time.time()
        temp = self.w1.genrandom(size)
        t3 = time.time()
        time.sleep(4)

        print(f"G4 kaes {t1 - t0} {t3 - t2}")
        return t1 - t0, t3 - t2

    def test(self):
        print("< r-gen >")
        t0, t1 = self.rgen(3 * 1024 * 1024 * 1024)
        print(f"Rgen {3072 / t0} MiB/s {3072 / t1} MiB/s")
        time.sleep(2)
        
        print("< py-mode >")
        t0, t1, t2, t3 = self.kaes(self.w0, "../tcond/setD.bin")
        print(f"Fmode {10240 / t0} MiB/s {10240 / t1} MiB/s {10240 / t2} MiB/s {10240 / t3} MiB/s")

        with open("../tcond/setC.bin", "rb") as f:
            temp = f.read()
        time.sleep(2)
        t0, t1, t2, t3 = self.kaes(self.w0, temp)
        del temp
        print(f"Bmode {3072 / t0} MiB/s {3072 / t1} MiB/s {3072 / t2} MiB/s {3072 / t3} MiB/s")
        time.sleep(2)

        print("< hy-mode >")
        t0, t1, t2, t3 = self.kaes(self.w1, "../tcond/setD.bin")
        print(f"Fmode {10240 / t0} MiB/s {10240 / t1} MiB/s {10240 / t2} MiB/s {10240 / t3} MiB/s")

        with open("../tcond/setC.bin", "rb") as f:
            temp = f.read()
        time.sleep(2)
        t0, t1, t2, t3 = self.kaes(self.w1, temp)
        del temp
        print(f"Bmode {3072 / t0} MiB/s {3072 / t1} MiB/s {3072 / t2} MiB/s {3072 / t3} MiB/s")

class g5:
    def __init__(self):
        self.w0 = kaesst.allmode()
        self.w1 = kaesst.funcmode()
        self.w2 = kaeshy.allmode()
        self.w3 = kaeshy.funcmode()
        self.w4 = kpic.toolbox()

        self.pw = b"0000"
        self.kf = kaesst.basickey()
        self.akey = b"\x00" * 48

        self.w5 = ksign.toolbox()
        self.pub, self.pri = self.w5.genkey(2048)

    def kzip(self, path):
        t0 = time.time()
        kzip.dozip( [path], "webp", "./temp.webp" )
        t1 = time.time()
        time.sleep(10)
        
        os.mkdir("../tcond/temp/")
        t2 = time.time()
        kzip.unzip("./temp.webp", "../tcond/temp/", True)
        t3 = time.time()
        time.sleep(10)

        print(f"G5 kzip {t1 - t0} {t3 - t2}")
        os.remove("./temp.webp")
        shutil.rmtree("../tcond/temp/")
        return t1 - t0, t3 - t2

    def rgen(self, size):
        t0 = time.time()
        temp = kaesst.genrand(size)
        t1 = time.time()
        time.sleep(10)
        del temp

        t2 = time.time()
        temp = self.w2.genrand(size)
        t3 = time.time()
        time.sleep(10)
        del temp

        print(f"G5 rgen {t1 - t0} {t3 - t2}")
        return t1 - t0, t3 - t2

    def kaesF(self, wa, wb, path):
        wa.signkey = [self.pub, self.pri]
        t0 = time.time()
        npath = wa.encrypt(self.pw, self.kf, path, 0)
        t1 = time.time()
        time.sleep(12)

        wa.view(npath)
        t2 = time.time()
        path = wa.decrypt(self.pw, self.kf, npath)
        t3 = time.time()
        time.sleep(12)
        os.remove(npath)

        wb.before = path
        wb.after = path + ".k"
        t4 = time.time()
        wb.encrypt(self.akey)
        t5 = time.time()
        time.sleep(12)

        wb.before = path + ".k"
        wb.after = path
        t6 = time.time()
        wb.decrypt(self.akey)
        t7 = time.time()
        time.sleep(12)

        os.remove(path + ".k")
        print(f"G5 kaesF {t1 - t0} {t3 - t2} {t5 - t4} {t7 - t6}")
        return t1 - t0, t3 - t2, t5 - t4, t7 - t6

    def kaesB(self, wa, wb, data):
        wa.signkey = [self.pub, self.pri]
        t0 = time.time()
        ndata = wa.encrypt(self.pw, self.kf, data, 0)
        t1 = time.time()
        time.sleep(12)
        del data

        wa.view(ndata)
        t2 = time.time()
        data = wa.decrypt(self.pw, self.kf, ndata)
        t3 = time.time()
        time.sleep(12)
        del ndata

        wb.before = data
        wb.after = None
        t4 = time.time()
        wb.encrypt(self.akey)
        t5 = time.time()
        time.sleep(12)

        wb.before = wb.after
        wb.after = None
        t6 = time.time()
        wb.decrypt(self.akey)
        t7 = time.time()
        time.sleep(12)

        print(f"G5 kaesB {t1 - t0} {t3 - t2} {t5 - t4} {t7 - t6}")
        return t1 - t0, t3 - t2, t5 - t4, t7 - t6

    def kpic(self, path, style):
        self.w4.setmold("", 2600, 2600)
        self.w4.style = style
        os.mkdir("../tcond/temp/")

        self.w4.target = path
        self.w4.export = "../tcond/temp/"
        t0 = time.time()
        self.w4.pack(2)
        t1 = time.time()
        time.sleep(12)

        self.w4.target = "../tcond/temp/"
        self.w4.export = "./temp"
        name, num, style = self.w4.detect()
        t2 = time.time()
        self.w4.unpack(name, num)
        t3 = time.time()
        time.sleep(12)

        os.remove("./temp")
        shutil.rmtree("../tcond/temp/")
        print(f"G5 kpic {style} {t1 - t0} {t3 - t2}")
        return t1 - t0, t3 - t2

    def test(self):
        print("< kzip >")
        t0, t1 = self.kzip("../tcond/setE/")
        print(f"kzip {10240 / t0} MiB/s {10240 / t1} MiB/s")
        time.sleep(4)

        print("< r-gen >")
        t0, t1 = self.rgen(3 * 1024 * 1024 * 1024)
        print(f"Rgen {3072 / t0} MiB/s {3072 / t1} MiB/s")
        time.sleep(4)

        print("< py-mode >")
        t0, t1, t2, t3 = self.kaesF(self.w0, self.w1, "../tcond/setD.bin")
        print(f"Fmode {10240 / t0} MiB/s {10240 / t1} MiB/s {10240 / t2} MiB/s {10240 / t3} MiB/s")

        with open("../tcond/setC.bin", "rb") as f:
            temp = f.read()
        time.sleep(4)
        t0, t1, t2, t3 = self.kaesB(self.w0, self.w1, temp)
        del temp
        print(f"Bmode {3072 / t0} MiB/s {3072 / t1} MiB/s {3072 / t2} MiB/s {3072 / t3} MiB/s")
        time.sleep(4)

        print("< hy-mode >")
        t0, t1, t2, t3 = self.kaesF(self.w2, self.w3, "../tcond/setD.bin")
        print(f"Fmode {10240 / t0} MiB/s {10240 / t1} MiB/s {10240 / t2} MiB/s {10240 / t3} MiB/s")

        with open("../tcond/setC.bin", "rb") as f:
            temp = f.read()
        time.sleep(4)
        t0, t1, t2, t3 = self.kaesB(self.w2, self.w3, temp)
        del temp
        print(f"Bmode {3072 / t0} MiB/s {3072 / t1} MiB/s {3072 / t2} MiB/s {3072 / t3} MiB/s")
        time.sleep(4)

        print("< kpic-png >")
        t0, t1 = self.kpic("../tcond/setB.bin", "png")
        print(f"png {1024 / t0} MiB/s {1024 / t1} MiB/s")
        time.sleep(4)

        print("< kpic-webp >")
        t0, t1 = self.kpic("../tcond/setB.bin", "webp")
        print(f"webp {1024 / t0} MiB/s {1024 / t1} MiB/s")
        time.sleep(2)
