# test574 : kpic (py)

import zlib # ksc5

import random
import os
import shutil
from PIL import Image
import multiprocessing as mp

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

# pack file -> bmp
def inf0(head, body, data, path, mode):
    buffer = [0] * len(body)
    if mode == 2:
        for i in range( 0, len(data) ):
            temp = 2 * i
            buffer[temp] = 16 * (body[temp] // 16) + (data[i] // 16)
            temp = temp + 1
            buffer[temp] = 16 * (body[temp] // 16) + (data[i] % 16)
        with open(path, "wb") as f:
            f.write(head)
            f.write( bytes(buffer) )
            
    else:
        for i in range( 0, len(data) ):
            temp = 4 * i
            buffer[temp] = 4 * (body[temp] // 4) + (data[i] // 64)
            temp = temp + 1
            buffer[temp] = 4 * (body[temp] // 4) + (data[i] // 16) % 4
            temp = temp + 1
            buffer[temp] = 4 * (body[temp] // 4) + (data[i] // 4) % 4
            temp = temp + 1
            buffer[temp] = 4 * (body[temp] // 4) + (data[i] % 4)
        with open(path, "wb") as f:
            f.write(head)
            f.write( bytes(buffer) )

# pack bmp -> png/webp
def inf1(before, after):
    temp = Image.open(before)
    temp.save(after, lossless=True)

# unpack png/webp -> bmp
def inf2(before, after):
    temp = Image.open(before)
    temp.save(after)

# unpack bmp -> file, return bytes
def inf3(path, mode):
    with open(path, "rb") as f:
        f.read(10)
        bmphsize = decode( f.read(4) )
    with open(path, "rb") as f:
        f.read(bmphsize)
        raw = f.read()
    if mode == 2:
        data = [0] * (len(raw) // 2)
        for i in range( 0, len(data) ):
            temp = 2 * i
            data[i] = 16 * (raw[temp] % 16) + (raw[temp + 1] % 16)
    else:
        data = [0] * (len(raw) // 4)
        for i in range( 0, len(data) ):
            temp = 4 * i
            data[i] = 64 * (raw[temp] % 4) + 16 * (raw[temp + 1] % 4) + 4 * (raw[temp + 2] % 4) + (raw[temp + 3] % 4)
    return bytes( data[8:] )

# move bmp if bmp is valid, return max num
def inf4(path, after, mode):
    with open(path, "rb") as f:
        f.read(10)
        bmphsize = decode( f.read(4) )
    with open(path, "rb") as f:
        f.read(bmphsize)
        raw = f.read(64)
    data = [0] * 8
    if mode == 2:
        for i in range(0, 8):
            temp = 2 * i
            data[i] = 16 * (raw[temp] % 16) + (raw[temp + 1] % 16)
    else:
        for i in range(0, 8):
            temp = 4 * i
            data[i] = 64 * (raw[temp] % 4) + 16 * (raw[temp + 1] % 4) + 4 * (raw[temp + 2] % 4) + (raw[temp + 3] % 4)
    test = bytes(data)
    for i in range(0, 4):
        if 65 > data[i] or data[i] > 64 + 26:
            return 0
    name = str(test[0:4], encoding="utf-8")
    num0 = decode( test[4:6] )
    num1 = decode( test[6:8] )
    shutil.move(path, after + f"{name}{num0}.bmp")
    return num1

class toolbox:
    def __init__(self):
        self.moldsize = [0, 0] # 주형 사진 크기 (가로 * 세로)
        self.moldpath = "" # 주형 사진 경로
        self.temppath = "./temp574/" # 임시 폴더 경로
        self.target = "./" # 패킹할 파일/언패킹할 폴더
        self.export = "./" # 패킹 결과/언패킹 결과로 출력할 폴더/파일
        self.style = "webp" # 패킹 결과/언패킹 대상의 스타일, "png"/"webp"/"bmp"

    # 주형 사진을 설정
    def setmold(self, path):
        # 주형 사진 경로 -> 내부 주형 사진 크기와 경로 설정
        temp = Image.open(path)
        w, h = temp.size
        self.moldsize = [w, h]
        path = os.path.abspath(path)
        path = path.replace("\\", "/")
        self.moldpath = path

    # 주형 사진을 임시 폴더 내부에 mold.bmp로 변형, (가로 * 세로)
    def convmold(self, row, col):
        # 주형 사진 경로와 크기를 바꿈
        if row % 4 != 0 or col % 4 != 0:
            raise Exception("Mold should be 4n*4m")
        self.moldsize = [row, col]
        self.temppath = os.path.abspath(self.temppath)
        self.temppath = self.temppath.replace("\\", "/")
        if self.temppath[-1] != "/":
            self.temppath = self.temppath + "/"
        self.clear(True)
        temp = Image.open(self.moldpath)
        temp = temp.resize( (row, col) )
        temp.save(self.temppath + "mold.bmp")

    # 임시 폴더 삭제, mknew가 T면 초기화 후 새로 생성
    def clear(self, mknew):
        if os.path.exists(self.temppath):
            shutil.rmtree(self.temppath)
        if mknew:
            os.mkdir(self.temppath)

    # 파일을 패키징, mode는 2/4, 일련번호와 개수 반환
    def pack(self, mode):
        # 먼저 setmold, convmold를 해야 함
        fsize = os.path.getsize(self.target) # 파일 사이즈
        csize = self.moldsize[0] * self.moldsize[1] * 3 // mode - 8 # 한 사진당 실질용량
        with open(self.temppath + "mold.bmp", "rb") as f:
            f.read(10)
            bmphsize = decode( f.read(4) )
            f.read(14)
            colorsize = decode( f.read(2) )
            if colorsize != 24:
                raise Exception(f"BMPbpp should be 24")
        with open(self.temppath + "mold.bmp", "rb") as f:
            bmph = f.read(bmphsize) # bmp 헤더
            bmpd = f.read() # bmp 데이터
            
        data0 = ["A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"]
        data1 = ["A", "E", "I", "O", "U"]
        data2 = ["B", "C", "D", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "V", "W", "X", "Y", "Z"]
        name = data0[ random.randrange(0,26) ] + data0[ random.randrange(0,26) ] + data0[ random.randrange(0,26) ]
        if mode == 2:
            name = data2[ random.randrange(0, 21) ] + name
        else:
            name = data1[ random.randrange(0, 5) ] + name # 사진 일련번호
        if fsize % csize == 0:
            num = fsize // csize
            add = 0 # 추가로 패딩할 크기
        else:
            num = fsize // csize + 1 # 사진 개수
            add = csize - fsize % csize
        count0 = (num - 1) // 32 # 32x 반복수
        count1 = (num - 1) % 32 + 1 # 1~32 마지막수
        h0 = bytes(name, encoding="utf-8") # 일련번호 헤더
        h1 = encode(num, 2) # 전체개수 헤더

        with open(self.target, "rb") as f:
            current = 0
            p = mp.Pool(32)
            for i in range(0, count0):
                for j in range(0, 32):
                    temp = f.read(csize)
                    temp = h0 + encode(current, 2) + h1 + temp
                    p.apply_async( inf0, (bmph, bmpd, temp, self.temppath + name + f"{current}.bmp", mode) )
                    current = current + 1
            p.close()
            p.join()

            if count1 != 1:
                p = mp.Pool(count1 - 1)
                for i in range(0, count1 - 1):
                    temp = f.read(csize)
                    temp = h0 + encode(current, 2) + h1 + temp
                    p.apply_async( inf0, (bmph, bmpd, temp, self.temppath + name + f"{current}.bmp", mode) )
                    current = current + 1
                p.close()
                p.join()

            temp = f.read(csize - add)
            temp = h0 + encode(current, 2) + h1 + temp + random.randbytes(add)
            inf0(bmph, bmpd, temp, self.temppath + name + f"{current}.bmp", mode)

        current = 0
        p = mp.Pool(32)
        for i in range(0, num // 32):
            for j in range(0, 32):
                p.apply_async( inf1, (self.temppath + name + f"{current}.bmp", self.temppath + name + f"{current}.{self.style}") )
                current = current + 1
        p.close()
        p.join()
        
        if num % 32 != 0:
            p = mp.Pool(num % 32)
            for i in range(0, num % 32):
                p.apply_async( inf1, (self.temppath + name + f"{current}.bmp", self.temppath + name + f"{current}.{self.style}") )
                current = current + 1
            p.close()
            p.join()

        self.export = os.path.abspath(self.export)
        self.export = self.export.replace("\\", "/")
        if self.export[-1] != "/":
            self.export = self.export + "/"
        for i in range(0, num):
            shutil.move(self.temppath + name + f"{i}.{self.style}", self.export + name + f"{i}.{self.style}")
            
        return name, num

    # 폴더경로, 일련번호, 개수를 입력받아 파일로 언패킹
    def unpack(self, name, num):
        self.target = os.path.abspath(self.target)
        self.target = self.target.replace("\\", "/")
        if self.target[-1] != "/":
            self.target = self.target + "/"
        self.temppath = os.path.abspath(self.temppath)
        self.temppath = self.temppath.replace("\\", "/")
        if self.temppath[-1] != "/":
            self.temppath = self.temppath + "/"
        self.clear(True)
        if name[0] in ["A", "E", "I", "O", "U"]:
            mode = 4
        else:
            mode = 2

        current = 0
        p = mp.Pool(32)
        for i in range(0, num // 32):
            for j in range(0, 32):
                p.apply_async( inf2, (self.target + name + f"{current}.{self.style}", self.temppath + name + f"{current}.bmp") )
                current = current + 1
        p.close()
        p.join()
        
        if num % 32 != 0:
            p = mp.Pool(num % 32)
            for i in range(0, num % 32):
                p.apply_async( inf2, (self.target + name + f"{current}.{self.style}", self.temppath + name + f"{current}.bmp") )
                current = current + 1
            p.close()
            p.join()

        with open(self.export, "wb") as f:
            current = 0
            wrbuf = [0] * 32
            for i in range(0, num // 32):
                p = mp.Pool(32)
                for j in range(0, 32):
                    wrbuf[j] = p.apply_async( inf3, (self.temppath + name + f"{current}.bmp", mode) )
                    current = current + 1
                p.close()
                p.join()
                for j in range(0, 32):
                    f.write( wrbuf[j].get() )

            if num % 32 != 0:
                wrbuf = [0] * (num % 32)
                p = mp.Pool(num % 32)
                for i in range(0, num % 32):
                    wrbuf[i] = p.apply_async( inf3, (self.temppath + name + f"{current}.bmp", mode) )
                    current = current + 1
                p.close()
                p.join()
                for i in range(0, num % 32):
                    f.write( wrbuf[i].get() )

    # 폴더를 입력받아 그 안의 kpic의 일련번호와 개수 반환
    def detect(self):
        self.target = os.path.abspath(self.target)
        self.target = self.target.replace("\\", "/")
        if self.target[-1] != "/":
            self.target = self.target + "/"
        flist = os.listdir(self.target)
        plist = [ ] # 사진 리스트
        alp = [ chr(x) for x in range(65, 65 + 26) ] # 대문자 알파벳
        numb = [ str(x) for x in range(0, 10) ] # 0~9 숫자
        for i in flist:
            if "." in i and len(i) > 8:
                forward = i[ 0:i.rfind(".") ] # . 전
                backward = i[i.rfind(".") + 1:].lower() # . 후
                flag = True
                for j in range(0, 4):
                    if forward[j] not in alp:
                        flag = False
                for j in range( 4, len(forward) ):
                    if forward[j] not in numb:
                        flag = False
                if flag and backward in ["bmp", "png", "webp"]:
                    plist.append(i)
        if plist == [ ]:
            return "", 0
        else:
            plist.sort()
            temp = plist[0]
            name = temp[0:4]
            num = 0
            style = temp[temp.find(".") + 1:]
            while f"{name}{num}.{style}" in plist:
                num = num + 1
            return name, num, backward

    # 파일들을 입력받아 내부 기록을 이용하여 올바른 kpic 파일로 복구, 전체 개수 반환
    def restore(self, files, mode):
        self.export = os.path.abspath(self.export)
        self.export = self.export.replace("\\", "/")
        if self.export[-1] != "/":
            self.export = self.export + "/"
        self.temppath = os.path.abspath(self.temppath)
        self.temppath = self.temppath.replace("\\", "/")
        if self.temppath[-1] != "/":
            self.temppath = self.temppath + "/"
        self.clear(True)

        current = 0
        p = mp.Pool(32)
        for i in range(0, len(files) // 32):
            for j in range(0, 32):
                p.apply_async( inf2, (files[current], self.temppath + str(current) + ".bmp") )
                current = current + 1
        p.close()
        p.join()

        if len(files) % 32 != 0:
            p = mp.Pool(len(files) % 32)
            for i in range(0, len(files) % 32):
                p.apply_async( inf2, (files[current], self.temppath + str(current) + ".bmp") )
                current = current + 1
            p.close()
            p.join()

        maxs = [0] * len(files)
        for i in range( 0, len(files) ):
            maxs[i] = inf4(self.temppath + str(i) + ".bmp", self.export, mode)
        return max(maxs)

# ========== ========== kzip5 end ========== ==========
