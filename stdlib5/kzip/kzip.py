# test570 : kzip (py)

import zlib # ksc5

import os
import shutil
import picdt

# ========== ========== ksc5 start ========== ==========

# little endian encoding
def encode(num, length):
    temp = [0] * length
    for i in range(0, length):
        temp[i] = num % 256
        num = num // 256
    return bytes(temp)

# little endian decoding
def decode(data):
    temp = 0
    for i in range( 0, len(data) ):
        if data[i] != 0:
            temp = temp + data[i] * 256 ** i
    return temp

# crc32
def crc32(data):
    return encode(zlib.crc32(data), 4)

# find 1024nB + KSC5 pos
def findpos(data):
    temp = 0
    while len(data) >= temp + 4:
        if data[temp:temp + 4] == b"KSC5":
            return temp
        else:
            temp = temp + 1024
    return -1

# ========== ========== ksc5 end ========== ==========

# ========== ========== kzip5 start ========== ==========

def initfile(path): # file 존재시 삭제
    if os.path.exists(path):
        os.remove(path)

def initfolder(path): # folder 존재시 삭제
    if os.path.exists(path):
        shutil.rmtree(path)

def getlist(root, path, lstf, lstd): # root : 표준절대경로/, path : 표준절대경로/에 폴더,파일 표준절대경로 추가
    temp = [ path + x for x in os.listdir(root + path) ]
    for i in range( 0, len(temp) ):
        if os.path.isdir( root + temp[i] ) and temp[i][-1] != "/":
            temp[i] = temp[i] + "/"
    for i in temp:
        if i[-1] == "/":
            lstd.append(i)
            getlist(root, i, lstf, lstd)
        else:
            lstf.append(i)

class toolbox:
    def __init__(self): # ! 실행 전 파라미터 4개를 맞춰야 함 !
        self.noerr = False # crc 오류 무시 여부
        self.folder = "" # 폴더 절대경로
        self.file = [ ] # 파일들 절대경로
        self.export = "./temp570" # 결과 출력 위치

    # 파일/폴더 표준절대경로화
    def abs(self):
        self.folder = os.path.abspath(self.folder)
        self.folder = self.folder.replace("\\", "/")
        if self.folder[-1] != "/":
            self.folder = self.folder + "/"
        temp = [os.path.abspath(x) for x in self.file]
        self.file = [x.replace("\\", "/") for x in temp]

    # 파일을 패키징, mode = "png"/"webp"/"nah"
    def zipfile(self, mode):
        if mode == "png":
            fakeh = picdt.toolbox().data0
            fakeh = fakeh + b"\x00" * ( ( 0 - len(fakeh) ) % 1024 )
        elif mode == "webp":
            fakeh = picdt.toolbox().data1
            fakeh = fakeh + b"\x00" * ( ( 0 - len(fakeh) ) % 1024 )
        else:
            fakeh = b"" # prehead + padding
        commonh = b"KSC5" # common head
        subtypeh = b"KZIP" # subtype head
        h0 = f"folder = 0; file = {len(self.file)}"
        ht = [x.replace("\\", "/") for x in self.file]
        h1 = [x[x.rfind("/") + 1:] for x in ht]
        mainh = bytes(h0, encoding="utf-8") # main header
        hsize = encode(len(mainh), 4) # main head size
        res = crc32(mainh) # reserved

        initfile(self.export)
        with open(self.export, "wb") as f:
            f.write(fakeh)
            f.write(commonh)
            f.write(subtypeh)
            f.write(res)
            f.write(hsize)
            f.write(mainh)

            for i in range( 0, len(self.file) ):
                nmb = bytes(h1[i], encoding="utf-8") # name B
                nms = encode(len(nmb), 8) # name len
                f.write(nms + nmb)

                fsize = os.path.getsize( self.file[i] )
                fls = encode(fsize, 8) # file len
                f.write(fls)
                with open(self.file[i], "rb") as t:
                    for j in range(0, fsize // 10485760):
                        flb = t.read(10485760)
                        f.write(flb)
                    flb = t.read(fsize % 10485760) # file B
                    f.write(flb)

    # 폴더를 패키징, mode = "png"/"webp"/"nah"
    def zipfolder(self, mode):
        if mode == "png":
            fakeh = picdt.toolbox().data0
            fakeh = fakeh + b"\x00" * ( ( 0 - len(fakeh) ) % 1024 )
        elif mode == "webp":
            fakeh = picdt.toolbox().data1
            fakeh = fakeh + b"\x00" * ( ( 0 - len(fakeh) ) % 1024 )
        else:
            fakeh = b"" # prehead + padding
        commonh = b"KSC5" # common head
        subtypeh = b"KZIP" # subtype head

        root = self.folder.replace("\\", "/")
        if root[-1] == "/": # ~/에서 /제거
            root = root[0:-1]
        fr0 = root[0:root.rfind("/") + 1] # root의 상위 폴더 절대경로
        fr1 = root[root.rfind("/") + 1:] + "/" # 패키징할 root
        wrfile = [ ]
        wrfolder = [fr1]
        getlist(fr0, fr1, wrfile, wrfolder)

        h0 = f"folder = {len(wrfolder)}; file = {len(wrfile)}"
        h1 = "\n".join( [h0] + wrfolder )
        mainh = bytes(h1, encoding="utf-8") # main header
        hsize = encode(len(mainh), 4) # main head size
        res = crc32(mainh) # reserved

        initfile(self.export)
        with open(self.export, "wb") as f:
            f.write(fakeh)
            f.write(commonh)
            f.write(subtypeh)
            f.write(res)
            f.write(hsize)
            f.write(mainh)

            for i in range( 0, len(wrfile) ):
                nmb = bytes(wrfile[i], encoding="utf-8") # name B
                nms = encode(len(nmb), 8) # name len
                f.write(nms + nmb)

                fsize = os.path.getsize( fr0 + wrfile[i] )
                fls = encode(fsize, 8) # file len
                f.write(fls)
                with open(fr0 + wrfile[i], "rb") as t:
                    for j in range(0, fsize // 10485760):
                        flb = t.read(10485760)
                        f.write(flb)
                    flb = t.read(fsize % 10485760) # file B
                    f.write(flb)

    # 패키징 해제, path는 kzip파일 경로
    def unzip(self, path):
        with open(path, "rb") as f:
            if os.path.getsize(path) > 16384:
                temp = f.read(16384)
            else:
                temp = f.read()
        pos = findpos(temp)
        if pos == -1:
            raise Exception("Not Valid KSC5")

        expath = os.path.abspath(self.export)
        expath = expath.replace("\\", "/")
        if expath[-1] != "/":
            expath = expath + "/"
        initfolder(expath)
        os.mkdir(expath)

        with open(path, "rb") as f:
            f.read(pos + 4)
            subtypeh = f.read(4) # subtype head
            res = f.read(4) # reserved
            hsize = decode( f.read(4) )
            mainh = f.read(hsize) # main header

            if not self.noerr:
                if subtypeh != b"KZIP":
                    raise Exception("Not Valid KZIP")
                if crc32(mainh) != res:
                    raise Exception("Broken Header")

            doc = str(mainh, encoding="utf-8").split("\n") # 헤더 분리
            if doc[-1] == "":
                doc = doc[0:-1]
            infoline = doc[0] # file folder num
            if len(doc) == 1:
                doc = [ ]
            else:
                doc = doc[1:] # folders
                
            foldernum = 0
            filenum = 0
            for i in infoline.split(";"):
                if "folder" in i:
                    i = i.replace(" ", "")
                    foldernum = int( i[i.find("=") + 1:] )
                elif "file" in i:
                    i = i.replace(" ", "")
                    filenum = int( i[i.find("=") + 1:] )
            for i in doc:
                os.mkdir(expath + i)

            for i in range(0, filenum):
                nms = decode( f.read(8) )
                nm = str(f.read(nms), encoding="utf-8")
                fls = decode( f.read(8) )
                with open(expath + nm, "wb") as t:
                    for j in range(0, fls // 10485760):
                        flb = f.read(10485760)
                        t.write(flb)
                    flb = f.read(fls % 10485760)
                    t.write(flb)

# ========== ========== kzip5 end ========== ==========
