# test646 : stdlib5.kpic

import os
import shutil
import io
import random

from PIL import Image
import multiprocessing as mp

import kobj
import picdt

def inf0(head, body, data, path, zmode, pmode): # pack bytes -> pic
    temp = bytearray( len(body) )
    div = 16 if zmode == 2 else 4
    for i in range( 0, len(data) ):
        offset, adjust, reg = zmode * i, zmode - 1, data[i]
        while adjust >= 0:
            temp[offset + adjust] = div * (body[offset + adjust] // div) + reg % div
            adjust, reg = adjust - 1, reg // div

    bmpv = head + bytes(temp)
    temp = None
    if pmode == "bmp":
        with open(path, "wb") as f:
            f.write(bmpv)
    else:
        tgt = Image.open( io.BytesIO(bmpv) )
        tgt.save(path, format=pmode, lossless=True)
        tgt, bmpv = None, None

def inf1(path, zmode, pmode): # unpack pic -> bytes
    raw = inf3(path, pmode)
    temp = bytearray(len(raw) // zmode)
    div = 16 if zmode == 2 else 4
    for i in range( 0, len(temp) ):
        offset, adjust, reg = zmode * i, 0, div ** (zmode - 1)
        while adjust < zmode:
            temp[i] = temp[i] + reg * (raw[offset + adjust] % div)
            adjust, reg = adjust + 1, reg // div
    return bytes(temp)[8:]

def inf2(path, zmode, pmode): # examine pic by internal data, returns (name, curnum, maxnum)
    raw = inf3(path, pmode)
    temp = bytearray(8)
    div = 16 if zmode == 2 else 4
    for i in range( 0, len(temp) ):
        offset, adjust, reg = zmode * i, 0, div ** (zmode - 1)
        while adjust < zmode:
            temp[i] = temp[i] + reg * (raw[offset + adjust] % div)
            adjust, reg = adjust + 1, reg // div
    temp = bytes(temp)
    
    for i in temp[0:4]:
        if i < 65 or i > 64 + 26:
            return ("", -1, -1)
    name, num0, num1 = str(temp[0:4], encoding="utf-8"), kobj.decode( temp[4:6] ), kobj.decode( temp[6:8] )
    if num0 < num1:
        return (name, num0, num1)
    else:
        return ("", -1, -1)
        
def inf3(path, pmode): # read pic data, returns bmpv
    if pmode == "bmp":
        with open(path, "rb") as f:
            raw = f.read()
    else:
        tgt, bmpv = Image.open(path), io.BytesIO()
        tgt.save(bmpv, format="bmp")
        raw = bmpv.getvalue()
        tgt, bmpv = None, None
    hsize = kobj.decode( raw[10:14] )
    return raw[hsize:]

class toolbox:
    def __init__(self):
        self.moldsize = [0, 0]
        self.moldhead = b""
        self.moldbody = b""

        self.target = "" # conv target file/folder
        self.export = "" # conv result file/folder
        self.style = "webp" # conv style, "webp"/"png"/"bmp"

        self.proc = -1.0 # progress, -1 : not started, 0~1 : working, 2 : finished

    # set mold pic, (row, col) : -1/4n
    def setmold(self, path, row, col):
        if path == "":
            raw = Image.open( io.BytesIO(picdt.toolbox().data3) )
        else:
            raw = Image.open(path)
        if row < 0 or col < 0:
            row, col = raw.size
        else:
            raw = raw.resize( (row, col) )
        if row % 4 != 0 or col % 4 != 0:
            raise Exception("mold should be 4N*4M size")
        self.moldsize = [row, col]

        temp = io.BytesIO()
        raw.save(temp, format='bmp')
        bmpv = temp.getvalue()
        hsize, csize = kobj.decode( bmpv[10:14] ), kobj.decode( bmpv[28:30] )
        if csize != 24:
            raise Exception("color should be 24 bit")
        self.moldhead, self.moldbody = bmpv[0:hsize], bmpv[hsize:]

    # detect kpic info by target path -> (name, num, style)
    def detect(self):
        self.target = os.path.abspath(self.target).replace("\\", "/")
        if self.target[-1] != "/":
            self.target = self.target + "/"
        flist, plist, chars, nums = os.listdir(self.target), [ ], [ chr(x) for x in range(65, 65 + 26) ], [ str(x) for x in range(0, 10) ]

        for i in flist:
            if "." in i and len(i) > 8:
                fr0, fr1, flag = i[ 0:i.rfind(".") ], i[i.rfind(".") + 1:], True
                for j in range(0, 4):
                    if fr0[j] not in chars:
                        flag = False
                        break
                for j in range( 4, len(fr0) ):
                    if fr0[j] not in nums:
                        flag = False
                        break
                if flag and fr1 in ["webp", "png", "bmp"]:
                    plist.append(i)

        if plist == [ ]:
            return "", 0, ""
        else:
            fr = plist[0]
            name, num, style = fr[0:4], 0, fr[fr.find(".") + 1:]
            while f"{name}{num}.{style}" in plist:
                num = num + 1
            return name, num, style

    # pack file to pic, zmode = 2/4,  ! do setmold() first !
    def pack(self, zmode):
        self.proc = -1.0
        if os.path.exists(self.export):
            shutil.rmtree(self.export)
        os.mkdir(self.export) # !!! automatically clear export folder !!!
        self.export = os.path.abspath(self.export).replace("\\", "/")
        if self.export[-1] != "/":
            self.export = self.export + "/"
        fsize = os.path.getsize(self.target)
        csize = self.moldsize[0] * self.moldsize[1] * 3 // zmode - 8

        chars0 = [ chr(x) for x in range(65, 65 + 26) ]
        chars1 = ["A", "E", "I", "O", "U"]
        chars2 = [x for x in chars0 if x not in chars1]
        name = "".join( [ chars0[ random.randrange(0, 26) ] for x in range(0, 3) ] )
        name = chars2[ random.randrange(0, 21) ] + name if zmode == 2 else chars1[ random.randrange(0, 5) ] + name
        if fsize % csize == 0:
            num0, num1 = fsize // csize, 0 # pic num, added bytes length
        else:
            num0, num1 = fsize // csize + 1, csize - fsize % csize
        num2, num3 = (num0 - 1) // 32, (num0 - 1) % 32 + 1 # x32 repeating, 1~32 last cycle
        head0, head1 = bytes(name, encoding="utf-8"), kobj.encode(num0, 2) # name, maxnum header

        self.proc = 0.0
        current = 0
        with open(self.target, "rb") as f:
            p, q = mp.Pool(32), [0] * 32
            for i in range(0, num2):
                for j in range(0, 32):
                    data = head0 + kobj.encode(current, 2) + head1 + f.read(csize)
                    q[j] = p.apply_async( inf0, (self.moldhead, self.moldbody, data, f"{self.export}{name}{current}.{self.style}", zmode, self.style) )
                    current, self.proc = current + 1, current / num0
                q[0].wait()
            p.close()
            p.join()

            if num3 != 1:
                p, q = mp.Pool(num3 - 1), [0] * (num3 - 1)
                for j in range(0, num3 - 1):
                    data = head0 + kobj.encode(current, 2) + head1 + f.read(csize)
                    q[j] = p.apply_async( inf0, (self.moldhead, self.moldbody, data, f"{self.export}{name}{current}.{self.style}", zmode, self.style) )
                    current, self.proc = current + 1, current / num0
                q[0].wait()
                p.close()
                p.join()

            data = head0 + kobj.encode(current, 2) + head1 + f.read(csize - num1) + random.randbytes(num1)
            inf0(self.moldhead, self.moldbody, data, f"{self.export}{name}{current}.{self.style}", zmode, self.style)
            self.proc = 2.0
        return name, num0

    # unpack pic to file with name, num
    def unpack(self, name, num):
        self.proc = -1.0
        self.target = os.path.abspath(self.target).replace("\\", "/")
        if self.target[-1] != "/":
            self.target = self.target + "/"
        zmode = 4 if name[0] in ["A", "E", "I", "O", "U"] else 2

        self.proc = 0.0
        current = 0
        with open(self.export, "wb") as f:
            p, q = mp.Pool(32), [0] * 32
            for i in range(0, num // 32):
                for j in range(0, 32):
                    q[j] = p.apply_async( inf1, (f"{self.target}{name}{current}.{self.style}", zmode, self.style) )
                    current, self.proc = current + 1, current / num
                for j in range(0, 32):
                    f.write( q[j].get() )
            p.close()
            p.join()

            if num % 32 != 0:
                p, q = mp.Pool(num % 32), [0] * (num % 32)
                for j in range(0, num % 32):
                    q[j] = p.apply_async( inf1, (f"{self.target}{name}{current}.{self.style}", zmode, self.style) )
                    current, self.proc = current + 1, current / num
                for j in range(0, num % 32):
                    f.write( q[j].get() )
                p.close()
                p.join()
        self.proc = 2.0

    # restore and change name by internal data
    def restore(self, files, zmode):
        self.proc = -1.0
        files = [os.path.abspath(x).replace("\\", "/") for x in files]
        size = len(files)
        retv = [0] * size

        self.proc = 0.0
        current = 0
        p = mp.Pool(32)
        for i in range(0, size // 32):
            for j in range(0, 32):
                retv[current] = p.apply_async( inf2, (files[32 * i + j], zmode, self.style) )
                current, self.proc = current + 1, current / size
            for j in range(0, 32):
                retv[32 * i + j] = retv[32 * i + j].get()
        p.close()
        p.join()

        if size % 32 != 0:
            p = mp.Pool(size % 32)
            i = 32 * (size // 32)
            for j in range(0, size % 32):
                retv[current] = p.apply_async( inf2, (files[32 * i + j], zmode, self.style) )
                current, self.proc = current + 1, current / size
            for j in range(0, size % 32):
                retv[i + j] = retv[i + j].get()
            p.close()
            p.join()

        name, maxnum = "", 0
        for i in range(0, size):
            if retv[i][1] >= 0 and retv[i][2] > 0:
                path = files[i]
                os.rename(path, path[0:path.rfind("/") + 1] + f"{retv[i][0]}{retv[i][1]}.{self.style}")
                name = retv[i][0]
                if maxnum < retv[i][2]:
                    maxnum = retv[i][2]
        self.proc = 2.0
        return name, maxnum
