# test577 : kaes st (py)

import zlib # ksc5

import os
import multiprocessing as mp
import secrets
import scrypt
from Cryptodome.Cipher import AES

import kdb
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

# ========== ========== kaes5 start ========== ==========

def enshort(key, iv, data): # short encryption no padding 16nB
    func = AES.new(key, AES.MODE_CBC, iv).encrypt
    num = len(data)
    out = [b""] * (num // 16)
    for i in range(0, num // 16):
        temp = 16 * i
        out[i] = func( data[temp:temp+16] )
    return b"".join(out)

def deshort(key, iv, data): # short decryption no padding 16nB
    func = AES.new(key, AES.MODE_CBC, iv).decrypt
    num = len(data)
    out = [b""] * (num // 16)
    for i in range(0, num // 16):
        temp = 16 * i
        out[i] = func( data[temp:temp+16] )
    return b"".join(out)

def mkkey(salt, pw, kf): # generate mkey 48B
    return scrypt.hash(kf + pw + pw + kf + pw, salt, 16384, 8, 1, 48)

def svkey(salt, pw, kf): # generate key storage 256B
    return scrypt.hash(pw + pw + kf + kf + pw, salt, 524288, 8, 1, 256)

def inf0(header, key, iv, data): # 바이트 암호화 내부함수
    size = len(data)
    num0 = size // 524288 # chunk num
    num1 = size % 524288 # left size
    temp = [0] * 32 # encryption buffer
    out = [ ] # output buffer
    p = mp.Pool(32)
    count = 0 # iv, key position

    for i in range(0, num0):
        tempdata = data[524288 * i : 524288 * i + 524288] # 512kb
        temp[count] = p.apply_async( enshort, (key[count], iv[count], tempdata) )
        if count == 31:
            for j in range(0,32):
                temp[j] = temp[j].get()
                iv[j] = temp[j][-16:]
            out = out + temp
            temp = [0] * 32 # reset
            count = -1 # reset
        count = count + 1
        
    for i in range(0, num0 % 32):
        temp[i] = temp[i].get()
        iv[i] = temp[i][-16:]
    if num0 % 32 != 0:
        out = out + temp[0:num0 % 32]
    pad = lambda x : x + bytes(chr(16 - len(x) % 16),'utf-8') * (16 - len(x) % 16)
    tempdata = data[524288 * num0:]
    tempdata = pad(tempdata)
    out.append( enshort(key[count], iv[count], tempdata) )
    p.close()
    p.join()

    return header + b"".join(out)

def inf1(key, iv, data): # 바이트 복호화 내부함수
    size = len(data)
    num0 = size // 524288 # chunk num
    num1 = size % 524288 # left size
    if num1 == 0:
        num0 = num0 - 1
        num1 = 524288
    temp = [0] * 32 # decryption buffer
    out = [ ] # output buffer
    p = mp.Pool(32)
    count = 0 # iv, key position

    for i in range(0,num0):
        tempdata = data[524288 * i : 524288 * i + 524288] # 128kb
        temp[count] = p.apply_async( deshort, (key[count], iv[count], tempdata) )
        iv[count] = tempdata[-16:]
        if count == 31:
            for j in range(0,32):
                temp[j] = temp[j].get()
            out = out + temp
            temp = [0] * 32 # reset
            count = -1 # reset
        count = count + 1
        
    for i in range(0, num0 % 32):
        temp[i] = temp[i].get()
    if num0 % 32 != 0:
        out = out + temp[0:num0 % 32]
    unpad = lambda x : x[:-x[-1]]
    tempdata = data[524288 * num0:]
    out.append( unpad( deshort(key[count], iv[count], tempdata) ) )
    p.close()
    p.join()

    return b"".join(out)
    
def inf2(header, key, iv, before, after): # 파일 암호화 내부함수
    size = os.path.getsize(before)
    num0 = size // 524288 # chunk num
    num1 = size % 524288 # left size
    p = mp.Pool(32)
    order = [0] * 32 # order buffer
    write = [b''] * 32 # write buffer
    count = 0 # iv, key position

    with open(after,'wb') as f:
        with open(before,'rb') as t:
            f.write(header)

            for i in range(0, num0):
                tempdata = t.read(524288) # 512kb
                order[count] = p.apply_async( enshort, (key[count], iv[count], tempdata) )
                if count == 31:
                    for j in range(0,32):
                        write[j] = order[j].get()
                        iv[j] = write[j][-16:]
                    f.write( b''.join(write) )
                    count = -1 # reset
                    order = [0] * 32 # reset
                    write = [b''] * 32 # reset
                count = count + 1
                    
            for i in range(0, num0 % 32):
                write[i] = order[i].get()
                iv[i] = write[i][-16:]
            if num0 % 32 != 0:
                f.write( b''.join(write) )
            pad = lambda x : x + bytes(chr(16 - len(x)),'utf-8') * (16 - len(x))
            tempdata = t.read(num1)
            tempdata = tempdata[0:num1 - (num1 % 16)] + pad( tempdata[num1 - (num1 % 16):] )
            f.write( enshort(key[count], iv[count], tempdata) )
                
    p.close()
    p.join()

def inf3(stpoint, size, key, iv, before, after): # 파일 복호화 내부함수
    num0 = size // 524288 # chunk num
    num1 = size % 524288 # left size
    if num1 == 0:
        num0 = num0 - 1
        num1 = 524288
    p = mp.Pool(32)
    order = [0] * 32 # order buffer
    write = [b''] * 32 # write buffer
    count = 0 # iv, key position

    with open(before,'rb') as f:
        with open(after,'wb') as t:
            f.read(stpoint)
            
            for i in range(0, num0):
                tempdata = f.read(524288) # 512kb
                order[count] = p.apply_async( deshort, (key[count], iv[count], tempdata) )
                iv[count] = tempdata[-16:]
                if count == 31:
                    for j in range(0,32):
                        write[j] = order[j].get()
                    t.write( b''.join(write) )
                    count = -1 # reset
                    order = [0] * 32 # reset
                    write = [b''] * 32 # reset
                count = count + 1
                
            for i in range(0,num0 % 32):
                write[i] = order[i].get()
                iv[i] = write[i][-16:]
            t.write( b''.join(write) )
            unpad = lambda x : x[:-x[-1]]
            tempdata = f.read(num1)
            tempbyte = deshort(key[count], iv[count], tempdata)
            tempbyte = tempbyte[0:-16] + unpad( tempbyte[-16:] )
            t.write(tempbyte)
            
    p.close()
    p.join()

def inf4(size): # 패딩 후 사이즈 반환
    return size + 16 - size % 16

# 안전한 난수 nB 반환
def genrandom(size):
    return secrets.token_bytes(size)

# 키 파일 경로에 따라 키 파일 바이트 생성
def genkf(path):
    try:
        with open(path, 'rb') as f:
            kf = f.read()
    except:
        kf = basickey()
    return kf

class genbytes: # kdm용 일반 모드 - bytes
    def __init__(self):
        self.valid = True
        self.noerr = False
        self.mode = "webp"
        self.msg = ""

    # pw B, kf B, hint B, data B -> enc B
    def en(self, pw, kf, hint, data):
        salt = genrandom(32) # salt 32B
        pwhash = svkey(salt, pw, kf) # pwhash 256B
        mkey = mkkey(salt, pw, kf) # master key 48B
        ckey = genrandom(1536) # content key 1536B
        ckeydt = enshort(mkey[16:48], mkey[0:16], ckey) # content key data 1536B
        pw = b'0' * 64
        kf = b'0' * 64
        mkey = b'0' * 64
        
        mold = "mode = 0\nmsg = 0\nsalt = 0\npwhash = 0\nhint = 0\nckeydt = 0\n"
        kdbtbox = kdb.toolbox()
        kdbtbox.readstr(mold)
        kdbtbox.fixdata("mode", "bytes")
        kdbtbox.fixdata("msg", self.msg)
        kdbtbox.fixdata("salt", salt)
        kdbtbox.fixdata("pwhash", pwhash)
        kdbtbox.fixdata("hint", hint)
        kdbtbox.fixdata("ckeydt", ckeydt)
        mh = bytes(kdbtbox.writestr(), encoding="utf-8") # main header
        mhs = encode(len(mh), 4) # main header size

        if self.mode == "png":
            fakeh = picdt.toolbox().data4
            fakeh = fakeh + b"\x00" * ( ( 16384 - len(fakeh) ) % 1024 )
        elif self.mode == "webp":
            fakeh = picdt.toolbox().data5
            fakeh = fakeh + b"\x00" * ( ( 16384 - len(fakeh) ) % 1024 )
        else:
            fakeh = b"" # prehead + padding
        commonh = b"KSC5" # common head
        subtypeh = b"KAES" # subtype head
        res = crc32(mh) # reserved
        chunksize = encode(inf4( len(data) ), 8)

        key = [b""] * 32 # key 32B list
        iv = [b""] * 32 # iv 162B list
        for i in range(0, 32):
            iv[i] = ckey[16 * i : 16 * i + 16]
            key[i] = ckey[512 + 32 * i : 512 + 32 * i + 32]

        content = inf0(fakeh + commonh + subtypeh + res + mhs + mh + chunksize, key, iv, data)
        return content

    # pw B, kf B, data B, stpoint N -> plain B
    def de(self, pw, kf, data, stpoint):
        mhs = decode( data[stpoint:stpoint + 4] )
        mh = data[stpoint + 4:stpoint + 4 + mhs] # main header
        chunksize = decode( data[stpoint + 4 + mhs:stpoint + 12 + mhs] )
        data = data[stpoint + 12 + mhs:stpoint + 12 + mhs + chunksize]
        
        kdbtbox = kdb.toolbox()
        kdbtbox.readstr( str(mh, encoding="utf-8") )
        salt = kdbtbox.getdata("salt")[3]
        pwhash = kdbtbox.getdata("pwhash")[3]
        ckeydt = kdbtbox.getdata("ckeydt")[3]
        if svkey(salt, pw, kf) != pwhash:
            raise Exception("Not Valid PWKF")

        mkey = mkkey(salt, pw, kf) # master key 48B
        ckey = deshort(mkey[16:48], mkey[0:16], ckeydt) # content key 1536B
        pw = b'0' * 64
        kf = b'0' * 64
        mkey = b'0' * 64

        key = [b""] * 32
        iv = [b""] * 32
        for i in range(0, 32):
            iv[i] = ckey[16 * i:16 * i + 16]
            key[i] = ckey[512 + 32 * i:512 + 32 * i + 32]

        return inf1(key, iv, data)

    # data B -> hint B, msg S, stpoint N
    def view(self, data):
        if len(data) > 16384:
            pos = findpos( data[0:16384] )
        else:
            pos = findpos(data)
        if pos == -1:
            raise Exception("Not Valid KSC5")

        subtypeh = data[pos + 4:pos + 8] # subtype head
        res = data[pos + 8:pos + 12] # reserved
        mhs = decode( data[pos + 12:pos + 16] )
        mh = data[pos + 16:pos + 16 + mhs] # main header
        
        if not self.noerr:
            if subtypeh != b"KAES":
                raise Exception("Not Valid KAES")
            if crc32(mh) != res:
                raise Exception("Broken Header")

        kdbtbox = kdb.toolbox()
        kdbtbox.readstr( str(mh, encoding="utf-8") )
        hint = kdbtbox.getdata("hint")[3]
        msg = kdbtbox.getdata("msg")[3]
        return hint, msg, pos + 12

class genfile: # kdm용 일반 모드 - file
    def __init__(self):
        self.valid = True
        self.noerr = False
        self.mode = "webp"
        self.msg = ""

    # pw B, kf B, hint B, path str -> new path S
    def en(self, pw, kf, hint, path):
        path = os.path.abspath(path).replace('\\', '/') # abs path
        fopath = path[0:path.rfind('/') + 1] # folder path
        name = path[path.rfind('/') + 1:] # true name
        if self.mode == "":
            tgt = fopath + bytes.hex( genrandom(3) ) + '.ke5' # write path
        else:
            tgt = fopath + bytes.hex( genrandom(3) ) + '.' + self.mode # write path
        nmb = bytes(name, 'utf-8')
        nmb = nmb + bytes(chr(16 - len(nmb) % 16), 'utf-8') * (16 - len(nmb) % 16)

        salt = genrandom(32) # salt 32B
        pwhash = svkey(salt, pw, kf) # pwhash 256B
        mkey = mkkey(salt, pw, kf) # master key 48B
        tkey = genrandom(48) # title key 48B
        tkeydt = enshort(mkey[16:48], mkey[0:16], tkey) # title key data 48B
        namedt = enshort(tkey[16:48], tkey[0:16], nmb) # name data
        ckey = genrandom(1536) # content key 1536B
        ckeydt = enshort(mkey[16:48], mkey[0:16], ckey) # content key data 1536B
        pw = b'0' * 64
        kf = b'0' * 64
        mkey = b'0' * 64

        mold = "mode = 0\nmsg = 0\nsalt = 0\npwhash = 0\nhint = 0\ntkeydt = 0\nckeydt = 0\nnamedt = 0\n"
        kdbtbox = kdb.toolbox()
        kdbtbox.readstr(mold)
        kdbtbox.fixdata("mode", "file")
        kdbtbox.fixdata("msg", self.msg)
        kdbtbox.fixdata("salt", salt)
        kdbtbox.fixdata("pwhash", pwhash)
        kdbtbox.fixdata("hint", hint)
        kdbtbox.fixdata("tkeydt", tkeydt)
        kdbtbox.fixdata("ckeydt", ckeydt)
        kdbtbox.fixdata("namedt", namedt)
        mh = bytes(kdbtbox.writestr(), encoding="utf-8") # main header
        mhs = encode(len(mh), 4) # main header size

        if self.mode == "png":
            fakeh = picdt.toolbox().data4
            fakeh = fakeh + b"\x00" * ( ( 16384 - len(fakeh) ) % 1024 )
        elif self.mode == "webp":
            fakeh = picdt.toolbox().data5
            fakeh = fakeh + b"\x00" * ( ( 16384 - len(fakeh) ) % 1024 )
        else:
            fakeh = b"" # prehead + padding
        commonh = b"KSC5" # common head
        subtypeh = b"KAES" # subtype head
        res = crc32(mh) # reserved
        chunksize = encode(inf4( os.path.getsize(path) ), 8)

        key = [b""] * 32 # key list
        iv = [b""] * 32 # iv list
        for i in range(0, 32):
            iv[i] = ckey[16 * i:16 * i + 16]
            key[i] = ckey[512 + 32 * i:512 + 32 * i + 32]
        
        inf2(fakeh + commonh + subtypeh + res + mhs + mh + chunksize, key, iv, path, tgt)
        return tgt

    # pw B, kf B, path str, stpoint N -> new path S
    def de(self, pw, kf, path, stpoint):
        path = os.path.abspath(path).replace('\\', '/') # abs path
        fopath = path[0:path.rfind('/') + 1] # folder path
        with open(path, "rb") as f:
            f.read(stpoint)
            mhs = decode( f.read(4) )
            mh = f.read(mhs) # main header
            chunksize = decode( f.read(8) )
        stpoint = stpoint + mhs + 12

        kdbtbox = kdb.toolbox()
        kdbtbox.readstr( str(mh, encoding="utf-8") )
        salt = kdbtbox.getdata("salt")[3]
        pwhash = kdbtbox.getdata("pwhash")[3]
        tkeydt = kdbtbox.getdata("tkeydt")[3]
        ckeydt = kdbtbox.getdata("ckeydt")[3]
        namedt = kdbtbox.getdata("namedt")[3]
        if svkey(salt, pw, kf) != pwhash:
            raise Exception("Not Valid PWKF")

        mkey = mkkey(salt, pw, kf) # master key 48B
        tkey = deshort(mkey[16:48], mkey[0:16], tkeydt) # title key 48B
        ckey = deshort(mkey[16:48], mkey[0:16], ckeydt) # content key 1536B
        nmb = deshort(tkey[16:48], tkey[0:16], namedt)
        nmb = nmb[:-nmb[-1]]
        name = str(nmb, encoding="utf-8")
        pw = b'0' * 64
        kf = b'0' * 64
        mkey = b'0' * 64

        key = [b""] * 32
        iv = [b""] * 32
        for i in range(0, 32):
            iv[i] = ckey[16 * i:16 * i + 16]
            key[i] = ckey[512 + 32 * i:512 + 32 * i + 32]

        tgt = fopath + name 
        inf3(stpoint, chunksize, key, iv, path, tgt)
        return tgt

    # path str -> hint B, msg S, stpoint N
    def view(self, path):
        with open(path, "rb") as f:
            size = os.path.getsize(path)
            if size > 16384:
                pos = findpos( f.read(16384) )
            else:
                pos = findpos( f.read(size) )
        if pos == -1:
            raise Exception("Not Valid KSC5")

        with open(path, "rb") as f:
            f.read(pos + 4)
            subtypeh = f.read(4) # subtype head
            res = f.read(4) # reserved
            mhs = decode( f.read(4) )
            mh = f.read(mhs) # main header
        
        if not self.noerr:
            if subtypeh != b"KAES":
                raise Exception("Not Valid KAES")
            if crc32(mh) != res:
                raise Exception("Broken Header")

        kdbtbox = kdb.toolbox()
        kdbtbox.readstr( str(mh, encoding="utf-8") )
        hint = kdbtbox.getdata("hint")[3]
        msg = kdbtbox.getdata("msg")[3]
        return hint, msg, pos + 12

class funcbytes: # kv용 기능 모드 - bytes
    def __init__(self):
        self.valid = True

    # key 48B, data B -> enc B
    def en(self, key, data):
        ckey = genrandom(1536) # content key 1536B
        ckeydt = enshort(key[16:48], key[0:16], ckey) # content key data 1536B
        key = b'0' * 64

        key = [b""] * 32 # key 32B list
        iv = [b""] * 32 # iv 162B list
        for i in range(0, 32):
            iv[i] = ckey[16 * i:16 * i + 16]
            key[i] = ckey[512 + 32 * i:512 + 32 * i + 32]

        content = inf0(ckeydt, key, iv, data)
        return content

    # key 48B, data B -> plain B
    def de(self, key, data):
        ckeydt = data[0:1536]
        ckey = deshort(key[16:48], key[0:16], ckeydt) # content key 1536B
        key = b'0' * 64

        key = [b""] * 32 # key 32B list
        iv = [b""] * 32 # iv 162B list
        for i in range(0, 32):
            iv[i] = ckey[16 * i:16 * i + 16]
            key[i] = ckey[512 + 32 * i:512 + 32 * i + 32]

        data = data[1536:]
        content = inf1(key, iv, data)
        return content

class funcfile: # kv용 기능 모드 - file
    def __init__(self):
        self.valid = True

    # key 48B, before -> after
    def en(self, key, before, after):
        ckey = genrandom(1536) # content key 1536B
        ckeydt = enshort(key[16:48], key[0:16], ckey) # content key data 1536B
        key = b'0' * 64

        key = [b""] * 32 # key 32B list
        iv = [b""] * 32 # iv 162B list
        for i in range(0, 32):
            iv[i] = ckey[16 * i:16 * i + 16]
            key[i] = ckey[512 + 32 * i:512 + 32 * i + 32]

        inf2(ckeydt, key, iv, before, after)

    # key 48B, before -> after
    def de(self, key, before, after):
        with open(before, "rb") as f:
            ckeydt = f.read(1536)
        ckey = deshort(key[16:48], key[0:16], ckeydt) # content key 1536B
        key = b'0' * 64

        key = [b""] * 32 # key 32B list
        iv = [b""] * 32 # iv 162B list
        for i in range(0, 32):
            iv[i] = ckey[16 * i:16 * i + 16]
            key[i] = ckey[512 + 32 * i:512 + 32 * i + 32]

        size = os.path.getsize(before) - 1536
        inf3(1536, size, key, iv, before, after)

# ========== ========== kaes5 end ========== ==========

def basickey():
    data0 = b""
    data0 = data0 + b"\xea\xb7\xb8\xeb\x9e\x98\x2c\x20\xeb\x82\x98\xeb\xa5\xbc\x20\xea\xb0\x80\xeb\x91\xac\xeb\x91\x94\x20\xec\xb1\x84\x20\xec\x9d\xb4\xea\xb3\xb3\xea\xb9\x8c\xec\xa7"
    data0 = data0 + b"\x80\x20\xec\x9e\xac\xeb\xb0\x8c\xeb\x8a\x94\x20\xec\x97\xac\xec\xa0\x95\x20\xeb\xb3\xb4\xeb\x83\x88\xeb\x8b\x88\x3f\x0d\x0a\xeb\xac\xb4\xeb\x84\x88\xec\xa0\xb8"
    data0 = data0 + b"\xeb\x9d\xbc\x2e\x0d\x0a\xec\x86\x9f\xec\x95\x84\xeb\x9d\xbc\x2e\x0d\x0a\xec\x9a\xb8\xeb\xa0\xa4\xeb\x9d\xbc\x2e\x0d\x0a\xec\x98\x9b\xeb\x82\xa0\xec\x9d\x98\x20"
    data0 = data0 + b"\xea\xb0\x90\xea\xb0\x81\xeb\x93\xa4\xec\x9d\xb4\x20\xeb\x8f\x8c\xec\x95\x84\xec\x98\xa4\xeb\x8a\x94\xea\xb5\xac\xeb\x82\x98\x2e\x0d\x0a\xec\x9d\xb4\xeb\xb2\x88"
    data0 = data0 + b"\xec\x97\x94\x20\xec\x8a\xa4\xec\x8a\xa4\xeb\xa1\x9c\xec\x9d\x98\x20\xed\x9e\x98\xec\x9c\xbc\xeb\xa1\x9c\x20\xeb\x82\x98\xeb\xa5\xbc\x20\xeb\xa7\x89\xec\x9d\x84"
    data0 = data0 + b"\x20\xec\x88\x98\x20\xec\x9e\x88\xea\xb2\xa0\xeb\x8b\x88\x3f\x0d\x0a\xed\x9d\xa9\xec\x96\xb4\xec\xa7\x80\xea\xb1\xb0\xeb\x9d\xbc\x2e\x0d\x0a\xec\x9a\x94\xec\xa0"
    data0 = data0 + b"\x95\xeb\x93\xa4\xec\x9d\xb4\xec\x97\xac\x2e\x0d\x0a\xeb\xa8\xb8\xeb\xa6\xac\xeb\xa5\xbc\x20\xec\x86\x8d\xec\x9d\xb8\x20\xec\xb1\x84\x20\xec\x9d\xb4\x20\xea\xb3"
    data0 = data0 + b"\xb3\xec\x97\x90\xec\x84\x9c\x20\xec\x9d\xb4\xeb\x9f\xb0\x20\xeb\x8b\xb9\xeb\x8f\x8c\xed\x95\x9c\x20\xec\xa7\x93\xec\x9d\x84\x20\xeb\x98\x90\x20\xeb\x8b\xa4\xec"
    data0 = data0 + b"\x8b\x9c\x20\xeb\xb2\x8c\xec\x9d\xb4\xea\xb3\xa0\x20\xec\x9e\x88\xea\xb5\xac\xeb\x82\x98\x2e\x0d\x0a\xec\xa7\x91\xec\x86\x8d\x2e\x0d\x0a\xec\x97\xb4\xec\x87\xa0"
    data0 = data0 + b"\x20\xec\x9d\x91\xec\xb6\x95\x2e\x0d\x0a\xea\xb0\x9c\xeb\xb0\xa9\x2e\x0d\x0a\xec\x9d\xb4\x20\xeb\xaa\xb8\xec\x9d\x80\x20\xec\x9d\xb4\xeb\x9f\xb0\x20\xed\x9e\x98"
    data0 = data0 + b"\xeb\x8f\x84\x20\xec\x93\xb8\x20\xec\x88\x98\x20\xec\x9e\x88\xea\xb5\xac\xeb\x82\x98\x2e\x20\xed\x9d\xa5\xeb\xaf\xb8\xeb\xa1\xad\xec\xa7\x80\xeb\xa7\x8c\x20\xeb"
    data0 = data0 + b"\xaf\xb8\xec\x95\xbd\xed\x95\xb4\x2e\x0d\x0a\xec\x9d\xb4\xeb\xb2\x88\xec\x97\x94\x20\xea\xb7\xb8\xeb\x85\x80\xec\x9d\x98\x20\xeb\x8f\x84\xec\x9b\x80\x20\xec\x97"
    data0 = data0 + b"\x86\xec\x9d\xb4\x20\xeb\xa7\x89\xec\x95\x84\xeb\xb3\xb4\xeb\xa0\xa4\xeb\xac\xb4\xeb\x82\x98\x2e\x0d\x0a\xeb\xb6\x84\xec\x84\x9d\x2e\x20\xec\x95\x95\xec\xb6\x95"
    data0 = data0 + b"\x2e\x20\xec\xa0\x84\xea\xb0\x9c\x2e\x0d\x0a\xeb\x82\xb4\x20\xec\x95\x9e\xec\x97\x90\x20\xec\x84\x9c\xec\xa7\x80\x20\xeb\xa7\x90\xea\xb1\xb0\xeb\x9d\xbc\x2e\x0d"
    data0 = data0 + b"\x0a\xec\x9d\xb4\x20\xeb\xaa\xb8\xec\x9d\x80\x20\xec\xa0\x9c\xec\x95\xbd\xec\x9d\xb4\x20\xeb\x84\x88\xeb\xac\xb4\x20\xeb\xa7\x8e\xec\x95\x84\x2e\x0d\x0a\xeb\xac"
    data0 = data0 + b"\xb4\xeb\x84\x88\xec\xa0\xb8\xeb\x82\xb4\xeb\xa0\xa4\xeb\x9d\xbc\x2e\x0d\x0a\xeb\x82\x98\xec\x98\xa4\xea\xb1\xb0\xeb\x9d\xbc\x2e\x0d\x0a\xeb\x82\xa0\xeb\x9b\xb0"
    data0 = data0 + b"\xea\xb1\xb0\xeb\x9d\xbc\x2e\x0d\x0a\xec\x9d\xb4\xea\xb3\xb3\xec\x9d\x98\x20\xed\x9e\x98\xec\x9d\x84\x20\xeb\x8d\x94\x20\xeb\xa8\xbc\xec\xa0\x80\x20\xec\x95\x8c"
    data0 = data0 + b"\xec\x95\x98\xeb\x8b\xa4\xeb\xa9\xb4\x20\xec\x9a\xb0\xeb\xa6\xac\xeb\x93\xa4\xeb\x8f\x84\x20\xea\xb7\xb8\x20\xed\x9e\x98\xec\x9d\x84\x20\xec\x93\xb8\x20\xec\x88"
    data0 = data0 + b"\x98\x20\xec\x9e\x88\xec\x97\x88\xea\xb2\xa0\xec\xa7\x80\x2e\x0d\x0a\xeb\x82\x98\xeb\xa5\xbc\x20\xeb\x84\x98\xec\xa7\x80\x20\xeb\xaa\xbb\xed\x95\x98\xeb\xa9\xb4"
    data0 = data0 + b"\x20\xea\xb2\xb0\xea\xb5\xad\x20\xeb\x98\x90\x20\xeb\x8b\xa4\xec\x8b\x9c\x20\xeb\xa8\xb8\xeb\xa6\xac\xec\x97\x90\x20\xeb\xb0\x9f\xed\x9e\x90\x20\xeb\xbf\x90\xec"
    data0 = data0 + b"\x9d\xb4\xeb\x9e\x80\xeb\x8b\xa4\x2e\x0d\x0a\xeb\xac\xb8\xec\x9d\x84\x20\xec\x97\xb4\xec\x96\xb4\xec\xa3\xbc\xeb\xa7\x88\x2e\x0d\x0a\xec\x9d\xb4\xea\xb3\xb3\xec"
    data0 = data0 + b"\x97\x90\x20\xed\x95\xa8\xea\xbb\x98\x20\xea\xb0\x80\xeb\x9d\xbc\xec\x95\x89\xec\x9e\x90\xea\xb5\xac\xeb\x82\x98\x2e\x0d\x0a\xec\x9d\xb4\x20\xea\xb0\x90\xec\x98"
    data0 = data0 + b"\xa5\xec\x97\x90\xec\x84\x9c\x20\xeb\x82\x98\xea\xb0\x84\xeb\x8b\xa4\x20\xed\x95\x98\xeb\x8d\x94\xeb\x9d\xbc\xeb\x8f\x84\x20\xed\x98\xbc\xec\x9e\x90\xec\x84\x9c"
    data0 = data0 + b"\x20\xeb\xac\xb4\xec\x97\x87\xec\x9d\x84\x20\xed\x95\xa0\x20\xec\x88\x98\x20\xec\x9e\x88\xec\x9d\x84\x20\xea\xb2\x83\x20\xea\xb0\x99\xeb\x8b\x88\x3f\x0d\x0a\xec"
    data0 = data0 + b"\x9d\xb4\x20\xea\xb5\xb4\xeb\xa0\x88\xeb\xa5\xbc\x20\xeb\x81\x8a\xeb\x8a\x94\xeb\x8b\xa4\x20\xed\x95\x98\xeb\x8d\x94\xeb\x9d\xbc\xeb\x8f\x84\x20\xec\x9e\xa0\xec"
    data0 = data0 + b"\x8b\x9c\xeb\xbf\x90\xec\x9d\xb4\xec\xa7\x80\x2e\x0d\x0a\xec\x98\x85\xec\x96\xb4\xec\xa7\x80\xeb\x8a\x94\xea\xb5\xac\xeb\x82\x98\x2e\x0d\x0a\xec\x9e\xa0\xec\x9d"
    data0 = data0 + b"\xb4\x20\xec\x98\xa4\xeb\x8a\x94\xea\xb5\xac\xeb\x82\x98\x2e\x0d\x0a\xec\x9a\xb0\xeb\xa6\xac\xeb\xa5\xbc\x20\xeb\xb2\x97\xec\x96\xb4\xeb\x82\xa0\x20\xec\x88\x98"
    data0 = data0 + b"\xeb\x8a\x94\x20\xec\x97\x86\xeb\x8b\xa8\xeb\x8b\xa4\x2e\x0d\x0a\xea\xb8\xb0\xed\x9a\x8c\xeb\xa5\xbc\x20\xeb\x86\x93\xec\xb3\xa4\xea\xb5\xac\xeb\x82\x98\x2e\x0d"
    data0 = data0 + b"\x0a\xea\xb5\xbd\xed\x9e\x88\xec\xa7\x80\x20\xeb\xaa\xbb\xed\x96\x88\xea\xb5\xac\xeb\x82\x98\x2e\x0d\x0a\xed\x98\xbc\xec\x9e\x90\xec\x84\x9c\xeb\x8a\x94\x20\xed"
    data0 = data0 + b"\x9d\x90\xeb\xa6\x84\xec\x9d\x84\x20\xeb\xa9\x88\xec\xb6\x9c\x20\xec\x88\x98\x20\xec\x97\x86\xeb\x8b\xa8\xeb\x8b\xa4\x2e\x0d\x0a\xed\x8c\x8c\xeb\x8f\x84\xeb\x8a"
    data0 = data0 + b"\x94\x20\xeb\x8b\xa4\xec\x8b\x9c\x20\xec\x9d\xbc\xeb\xa0\x81\xec\x9d\xbc\x20\xea\xb2\x83\xec\x9d\xb4\xeb\x9e\x80\xeb\x8b\xa4\x2e\x0d\x0a\xeb\x84\x88\xec\x9d\x98"
    data0 = data0 + b"\x20\xeb\xaf\xb8\xec\x88\x99\xed\x95\xa8\xec\x9d\xb4\xeb\x8b\xa4\x2e\x0d\x0a\xea\xb1\xb0\xeb\x8c\x80\xed\x95\x9c\x20\xed\x9d\x90\xeb\xa6\x84\xec\x9d\x84\x20\xec"
    data0 = data0 + b"\x86\x90\xeb\xb0\x94\xeb\x8b\xa5\xec\x9c\xbc\xeb\xa1\x9c\x20\xeb\xa7\x89\xec\x9d\x84\x20\xec\x88\x98\x20\xec\x97\x86\xeb\x8b\xa8\xeb\x8b\xa4\x2e\x0d\x0a\xec\x98"
    data0 = data0 + b"\xa4\xeb\xa1\xaf\xec\x9d\xb4\x20\xeb\x84\x88\xec\x9d\x98\x20\xed\x9e\x98\xeb\xa7\x8c\xec\x9c\xbc\xeb\xa1\x9c\x20\xeb\xa7\x89\xec\x96\xb4\xeb\xb3\xb4\xeb\xa0\xa4"
    data0 = data0 + b"\xeb\xac\xb4\xeb\x82\x98\x2e\x0d\x0a\xed\x9b\x8c\xeb\xa5\xad\xed\x95\x98\xea\xb5\xac\xeb\x82\x98\x2e\x0d\x0a\xea\xb7\xb8\xeb\x9e\x98\x2c\x20\xec\x9d\xb4\x20\xec"
    data0 = data0 + b"\xa0\x95\xeb\x8f\x84\xeb\xa9\xb4\x20\xec\xa7\x80\xec\xbc\x9c\xeb\xb3\xbc\x20\xea\xb0\x80\xec\xb9\x98\xea\xb0\x80\x20\xec\x9e\x88\xea\xb2\xa0\xec\xa7\x80\x2e\x0d"
    data0 = data0 + b"\x0a\xeb\x82\xb4\xea\xb0\x80\x20\xec\x96\xb4\xeb\x96\xbb\xea\xb2\x8c\x20\xea\xb7\xb8\xeb\xa6\xac\x20\xed\x95\x9c\xec\x97\x86\xec\x9d\xb4\x20\xec\x9e\x94\xec\x9d"
    data0 = data0 + b"\xb8\xed\x95\xb4\xec\xa7\x88\x20\xec\x88\x98\x20\xec\x9e\x88\xec\x97\x88\xeb\x8a\x94\xec\xa7\x80\x20\xec\x95\x8c\xeb\xa0\xa4\xec\xa4\x84\xea\xb9\x8c\x3f\x0d\x0a"
    data0 = data0 + b"\xec\x82\xac\xeb\x9e\x8c\xeb\x93\xa4\xec\x9d\x80\x20\xeb\xaa\xa8\xeb\x91\x90\x20\xeb\xb6\x88\xec\x95\x88\xec\x9d\x84\x20\xea\xb0\x80\xec\xa7\x84\x20\xec\xb1\x84"
    data0 = data0 + b"\xeb\xa1\x9c\x20\xec\x82\xb4\xec\x95\x84\xea\xb0\x84\xeb\x8b\xa8\xeb\x8b\xa4\x2e\x0d\x0a\xec\x9d\xb4\xea\xb1\xb4\x20\xeb\xaf\xb8\xec\xa7\x80\xec\x9d\x98\x20\xec"
    data0 = data0 + b"\x98\x81\xec\x97\xad\xec\x9d\x84\x20\xeb\xa7\x9e\xeb\x8b\xa5\xeb\x9c\xa8\xeb\xa6\xb4\x20\xeb\x95\x8c\x20\xeb\x8a\x90\xeb\x81\xbc\xeb\x8a\x94\x20\xeb\x8b\xb9\xec"
    data0 = data0 + b"\x97\xb0\xed\x95\x9c\x20\xeb\x8c\x80\xea\xb0\x80\xec\x95\xbc\x2e\x0d\x0a\xed\x95\x98\xec\xa7\x80\xeb\xa7\x8c\x20\xeb\x82\x98\xeb\x8a\x94\x20\xec\x84\xb8\xec\x83"
    data0 = data0 + b"\x81\xec\x97\x90\x20\xec\x82\xb4\xec\x95\x84\xeb\x82\xa8\xea\xb8\xb0\x20\xec\x9c\x84\xed\x95\xb4\x20\xea\xb7\xb8\x20\xea\xb3\xb5\xed\x8f\xac\xeb\xa5\xbc\x20\xeb"
    data0 = data0 + b"\xb0\x9b\xec\x95\x84\xeb\x93\xa4\xec\x9d\xb4\xec\xa7\x80\x20\xec\x95\x8a\xec\x9d\x80\x20\xec\xb1\x84\x20\xec\x8a\xa4\xec\x8a\xa4\xeb\xa1\x9c\x20\xeb\xa8\xb9\xec"
    data0 = data0 + b"\x96\xb4\xeb\xb2\x84\xeb\xa0\xb8\xec\x96\xb4\x2e\x0d\x0a\xea\xb7\xb8\xea\xb2\x83\xec\x9d\xb4\x20\xeb\x82\xb4\xea\xb0\x80\x20\xec\xa0\x80\xec\xa7\x80\xeb\xa5\xb8"
    data0 = data0 + b"\x20\xec\xb5\x9c\xec\xb4\x88\xec\x9d\xb4\xec\x9e\x90\x20\xec\xb5\x9c\xec\x95\x85\xec\x9d\x98\x20\xec\x95\x85\xed\x96\x89\xec\x9d\xb4\xec\x97\x88\xec\xa7\x80\x2e"
    data0 = data0 + b"\x0d\x0a\xed\x9b\x84\xed\x9a\x8c\xed\x95\x98\xec\xa7\x84\x20\xec\x95\x8a\xec\x95\x84\x2e\x20\xec\x83\x9d\xec\xa1\xb4\xec\x9d\x84\x20\xec\x9c\x84\xed\x95\x9c\x20"
    data0 = data0 + b"\xec\x84\xa0\xed\x83\x9d\xec\x9d\xb4\xec\x97\x88\xec\x9c\xbc\xeb\x8b\x88\x2e\x0d\x0a\xeb\x84\x88\xeb\x8f\x84\x20\xeb\xa7\x88\xec\xb0\xac\xea\xb0\x80\xec\xa7\x80"
    data0 = data0 + b"\x20\xec\x95\x84\xeb\x8b\x88\xec\x97\x88\xeb\x8b\x88\x3f\x0d\x0a\xeb\x82\xb4\x20\xeb\xa8\xb8\xeb\xa6\xbf\xec\x86\x8d\xec\x9d\x84\x20\xea\xb0\x88\xea\xb8\xb0\xea"
    data0 = data0 + b"\xb0\x88\xea\xb8\xb0\x20\xeb\xb6\x84\xed\x95\xb4\xec\x8b\x9c\xed\x82\xa4\xeb\x8d\x98\x20\xeb\x84\x88\xec\x9d\x98\x20\xed\x91\x9c\xec\xa0\x95\xec\x9d\x84\x20\xeb"
    data0 = data0 + b"\xb3\xb4\xea\xb3\xa0\x20\xec\x95\x8c\x20\xec\x88\x98\x20\xec\x9e\x88\xec\x97\x88\xec\x96\xb4\x2e\x0d\x0a\xec\x96\xb4\xec\xa9\x94\x20\xec\x88\x98\x20\xec\x97\x86"
    data0 = data0 + b"\xec\x97\x88\xeb\x8b\xa4\x20\xeb\x9d\xbc\xea\xb3\xa0\x20\xeb\x84\x88\xeb\x8f\x84\x20\xeb\xa7\x90\xec\x9d\x84\x20\xed\x95\x98\xeb\xa0\xa4\xeb\x82\x98\x2e\x0d\x0a"
    data0 = data0 + b"\xec\x9d\xb4\xec\xa0\x9c\xeb\x8a\x94\x20\xeb\x8b\xa4\xec\x8b\x9c\x20\xec\x9e\x8a\xea\xb3\xa0\x20\xec\x9e\x88\xeb\x8d\x98\x20\xea\xb3\xb5\xed\x8f\xac\xeb\xa5\xbc"
    data0 = data0 + b"\x20\xeb\xa7\x88\xec\xa3\xbc\xed\x95\x98\xeb\x8a\x94\x20\xea\xb2\x8c\x20\xec\xa2\x8b\xec\x9d\x84\x20\xea\xb1\xb0\xec\x95\xbc\x2e\x0d\x0a\xea\xb5\xb4\xeb\xa0\x88"
    data0 = data0 + b"\xeb\xa5\xbc\x20\xeb\x81\x8a\xec\x96\xb4\xeb\x82\xb4\xea\xb2\xa0\xeb\x8b\xa4\xeb\xa9\xb4\x2e\x0d\x0a\xeb\x84\xa4\xea\xb0\x80\x20\xea\xb7\xb8\xeb\xa0\x87\xeb\x8b"
    data0 = data0 + b"\xa4\xeb\xa9\xb4\x20\xea\xb7\xb8\xeb\x9f\xb0\x20\xea\xb2\x83\xec\x9d\xb4\xea\xb2\xa0\xec\xa7\x80\x2e\x0d\x0a"
    return data0

