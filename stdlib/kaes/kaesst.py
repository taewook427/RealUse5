# test634 : stdlib5.kaes st

import io
import os
import multiprocessing as mp

import hashlib
import secrets
import scrypt
from Cryptodome.Cipher import AES

import kdb
import ksc
import ksign
import picdt

def hash3512(data): # sha3-512 hash, data bytes
    return hashlib.sha3_512(data).digest()

def aescalc(data, key, iv, isenc, ispad): # aes-256 ende, key 32B iv 16B, isenc T:encrypt F:decrypt, ispad T:dopad F:nopad
    if len(key) != 32 or len(iv) != 16:
        raise Exception("invalid keyiv")
    module = AES.new(key, AES.MODE_CBC, iv)
    pad = lambda x : x + bytes(chr(16 - len(x) % 16),'utf-8') * (16 - len(x) % 16)
    unpad = lambda x : x[:-x[-1]]
    if isenc:
        if ispad:
            data = pad(data)
        data = module.encrypt(data)
    else:
        data = module.decrypt(data)
        if ispad:
            data = unpad(data)
    return data

def aesenc(key, iv, reader, writer, header, dsize): # generic encrypt, key bytes[40] / iv bytes[40] / rw BFIO / header bytes / dsize int
    if len(key) != 40 or len(iv) != 40:
        raise Exception("invalid keyiv")
    writer.write(header)
    p = mp.Pool(40)
    orders = [0] * 40
    num0 = dsize // 20971520
    num1 = dsize % 20971520
    num2 = num1 // 524288
    num3 = num1 % 524288
    inbuf = b""
    exbuf = b""
    tmbuf = b""

    if num0 > 0:
        inbuf = reader.read(20971520)

        for i in range(0, num0 - 1):
            for j in range(0, 40):
                temp = 524288 * j
                orders[j] = p.apply_async( aescalc, (inbuf[temp:temp + 524288], key[j], iv[j], True, False) )
            writer.write(exbuf)
            tmbuf, exbuf = reader.read(20971520), [b""] * 40
            for j in range(0, 40):
                exbuf[j] = orders[j].get()
                iv[j] = exbuf[j][-16:]
            exbuf = b"".join(exbuf)
            inbuf, tmbuf = tmbuf, b""

        for j in range(0, 40):
            temp = 524288 * j
            orders[j] = p.apply_async( aescalc, (inbuf[temp:temp + 524288], key[j], iv[j], True, False) )
        writer.write(exbuf)
        tmbuf, exbuf = reader.read(524288 * num2), [b""] * 40
        for j in range(0, 40):
            exbuf[j] = orders[j].get()
            iv[j] = exbuf[j][-16:]
        exbuf = b"".join(exbuf)
        inbuf, tmbuf = tmbuf, b""

    else:
        inbuf = reader.read(524288 * num2)

    for j in range(0, num2):
        temp = 524288 * j
        orders[j] = p.apply_async( aescalc, (inbuf[temp:temp + 524288], key[j], iv[j], True, False) )
    writer.write(exbuf)
    tmbuf, exbuf = reader.read(num3), [b""] * num2
    for j in range(0, num2):
        exbuf[j] = orders[j].get()
        iv[j] = exbuf[j][-16:]
    exbuf = b"".join(exbuf)
    inbuf, tmbuf = tmbuf, b""

    p.close()
    p.join()
    writer.write(exbuf)
    writer.write( aescalc(inbuf, key[num2], iv[num2], True, True) )

def aesdec(key, iv, reader, writer, stpoint, dsize): # generic decrypt, key bytes[40] / iv bytes[40] / rw BFIO / stpoint int / dsize int
    if len(key) != 40 or len(iv) != 40:
        raise Exception("invalid keyiv")
    reader.seek(stpoint)
    p = mp.Pool(40)
    orders = [0] * 40
    num0 = dsize // 20971520
    num1 = dsize % 20971520
    if num1 == 0:
        num0 = num0 - 1
        num1 = 20971520
    num2 = num1 // 524288
    num3 = num1 % 524288
    if num3 == 0:
        num2 = num2 - 1
        num3 = 524288
    inbuf = b""
    exbuf = b""
    tmbuf = b""

    if num0 > 0:
        inbuf = reader.read(20971520)

        for i in range(0, num0 - 1):
            for j in range(0, 40):
                temp = 524288 * j
                orders[j] = p.apply_async( aescalc, (inbuf[temp:temp + 524288], key[j], iv[j], False, False) )
            writer.write(exbuf)
            tmbuf, exbuf = reader.read(20971520), [b""] * 40
            for j in range(0, 40):
                temp = 524288 * (j + 1)
                exbuf[j] = orders[j].get()
                iv[j] = inbuf[temp - 16:temp]
            exbuf = b"".join(exbuf)
            inbuf, tmbuf = tmbuf, b""

        for j in range(0, 40):
            temp = 524288 * j
            orders[j] = p.apply_async( aescalc, (inbuf[temp:temp + 524288], key[j], iv[j], False, False) )
        writer.write(exbuf)
        tmbuf, exbuf = reader.read(524288 * num2), [b""] * 40
        for j in range(0, 40):
            temp = 524288 * (j + 1)
            exbuf[j] = orders[j].get()
            iv[j] = inbuf[temp - 16:temp]
        exbuf = b"".join(exbuf)
        inbuf, tmbuf = tmbuf, b""

    else:
        inbuf = reader.read(524288 * num2)

    for j in range(0, num2):
        temp = 524288 * j
        orders[j] = p.apply_async( aescalc, (inbuf[temp:temp + 524288], key[j], iv[j], False, False) )
    writer.write(exbuf)
    tmbuf, exbuf = reader.read(num3), [b""] * num2
    for j in range(0, num2):
        temp = 524288 * (j + 1)
        exbuf[j] = orders[j].get()
        iv[j] = inbuf[temp - 16:temp]
    exbuf = b"".join(exbuf)
    inbuf, tmbuf = tmbuf, b""

    p.close()
    p.join()
    writer.write(exbuf)
    writer.write( aescalc(inbuf, key[num2], iv[num2], False, True) )

def genpm(pw, kf, salt): # generate pwhash 128B, mkey 96B
    pwh = scrypt.hash(pw + pw + kf + pw + kf, salt, 524288, 8, 1, 128)
    mkey = scrypt.hash(kf + pw + kf + kf + pw, salt, 16384, 8, 1, 96)
    return pwh, mkey

def gensize(n): # generate size when n bytes is pad encrypted
    return n + 16 - n % 16

def genpath(size): # generate random word with length 2*size
    temp = genrand(size)
    word = ""
    for i in temp:
        word = word + str( hex(i + 16) )[2:4]
    return word

class allmode:
    def __init__(self):
        self.hint = "" # pwkf hint str
        self.msg = "" # program msg str
        self.signkey = ["", ""] # (pub, pri) ksign rsa key
        self.proc = 0.0 # -1 : not started, 0~1 : working, 2 : end

        self.salt = b""
        self.pwhash = b""
        self.ckeydata = b""
        self.tkeydata = b""
        self.encname = b""
        self.dinfo = [0, 0] # stpoint, encdsize

    # encryption, pwkf bytes / data bytes / pmode 0:webp 1:png 2:none -> result bytes / data str -> encpath str
    def encrypt(self, pw, kf, data, pmode):
        self.proc = -1.0
        ckey = genrand(1920) # (iv16 key32) * 40
        tkey = genrand(48) # iv16 key32
        if type(data) == bytes:
            oldpath = "NewData.bin"
            newpath = ""
            dsize = len(data)
            reader = io.BytesIO(data)
            writer = io.BytesIO()
        elif type(data) == str:
            data = os.path.abspath(data).replace("\\", "/")
            oldpath = data[data.rfind("/") + 1:]
            newpath = data[ 0:data.rfind("/") ] + "/" + genpath(2)
            if pmode == 0:
                newpath = newpath + ".webp"
            elif pmode == 1:
                newpath = newpath + ".png"
            else:
                newpath = newpath + ".k"
            dsize = os.path.getsize(data)
            reader = open(data, "rb")
            writer = open(newpath, "wb")
        else:
            raise Exception("invalid datatype")

        self.salt = genrand(40)
        self.pwhash, mkey = genpm(pw, kf, self.salt)
        self.ckeydata = aescalc(ckey, mkey[16:48], mkey[0:16], True, False)
        self.tkeydata = aescalc(tkey, mkey[64:96], mkey[48:64], True, False)
        self.encname = aescalc(bytes(oldpath, encoding="utf-8"), tkey[16:48], tkey[0:16], True, True)

        worker = kdb.toolbox()
        worker.read("salt = 0\npwhash = 0\nckeydata = 0\ntkeydata = 0\nencname = 0\nhint = 0\nmsg = 0")
        worker.fix("salt", self.salt)
        worker.fix("pwhash", self.pwhash)
        worker.fix("ckeydata", self.ckeydata)
        worker.fix("tkeydata", self.tkeydata)
        worker.fix("encname", self.encname)
        worker.fix("hint", self.hint)
        worker.fix("msg", self.msg)
        encheader = bytes(worker.write(), encoding="utf-8") # c0 : encheader

        if self.signkey[0] != "" and self.signkey[1] != "":
            w0, w1 = kdb.toolbox(), ksign.toolbox()
            w0.read("publickey = 0\nsigndata = 0")
            w0.fix( "publickey", self.signkey[0] )
            w0.fix( "signdata", w1.sign( self.signkey[1], hash3512(encheader) ) )
            signheader = bytes(w0.write(), encoding="utf-8")
        else:
            signheader = b"" # c1 : signheader

        worker, temp = ksc.toolbox(), picdt.toolbox()
        if pmode == 0:
            worker.prehead = temp.data5 + b"\x00" * (512 - len(temp.data5) % 512)
        elif pmode == 1:
            worker.prehead = temp.data4 + b"\x00" * (512 - len(temp.data4) % 512)
        else:
            worker.prehead = b""
        worker.subtype = b"KAES"
        worker.reserved = ksc.crc32hash(encheader) + ksc.crc32hash(signheader)
        header = worker.writeb() # real header to write
        header = worker.linkb(header, encheader)
        header = worker.linkb(header, signheader)
        header = header + ksc.encode(gensize(dsize), 8)

        keys, ivs = [0] * 40, [0] * 40
        for i in range(0, 40):
            temp = 48 * i
            keys[i], ivs[i] = ckey[temp + 16:temp + 48], ckey[temp:temp + 16]
        aesenc(keys, ivs, reader, writer, header, dsize)
        writer.write(b"\xff" * 8)

        self.proc = 2.0
        if type(data) == bytes:
            return writer.getvalue()
        else:
            reader.close()
            writer.close()
            return newpath

    # view encfile, data bytes/str
    def view(self, data):
        self.proc = -1.0
        worker = ksc.toolbox()
        worker.predetect = True
        if type(data) == bytes:
            worker.readb(data)
        elif type(data) == str:
            worker.path = data
            worker.readf()
        else:
            raise Exception("invalid datatype")
        if worker.subtype != b"KAES":
            raise Exception("invalidKAES5")
        
        if type(data) == bytes:
            encheader = data[worker.chunkpos[0] + 8:worker.chunkpos[0] + worker.chunksize[0] + 8]
            signheader = data[worker.chunkpos[1] + 8:worker.chunkpos[1] + worker.chunksize[1] + 8]
        else:
            with open(data, "rb") as f:
                f.seek(worker.chunkpos[0] + 8)
                encheader = f.read( worker.chunksize[0] )
                f.seek(worker.chunkpos[1] + 8)
                signheader = f.read( worker.chunksize[1] )
        self.dinfo = [ worker.chunkpos[2] + 8, worker.chunksize[2] ]
        if worker.reserved != ksc.crc32hash(encheader) + ksc.crc32hash(signheader):
            raise Exception("invalidCRC32")

        self.signkey = ["", ""]
        if len(signheader) != 0:
            worker, temp = kdb.toolbox(), ksign.toolbox()
            worker.read( str(signheader, encoding="utf-8") )
            self.signkey = [worker.get("publickey")[3], ""]
            signdata = worker.get("signdata")[3]
            if not temp.verify( self.signkey[0], signdata, hash3512(encheader) ):
                raise Exception("invalidRSAsign")
            
        worker = kdb.toolbox()
        worker.read( str(encheader, encoding="utf-8") )
        self.salt = worker.get("salt")[3]
        self.pwhash = worker.get("pwhash")[3]
        self.ckeydata = worker.get("ckeydata")[3]
        self.tkeydata = worker.get("tkeydata")[3]
        self.encname = worker.get("encname")[3]
        self.hint = worker.get("hint")[3]
        self.msg = worker.get("msg")[3]

    # decryption, pwkf bytes / data bytes -> result bytes / data str -> decpath str
    def decrypt(self, pw, kf, data):
        self.proc = -1.0
        if len(self.salt) == 0:
            raise Exception("should done view() first")
        pwhcmp, mkey = genpm(pw, kf, self.salt)
        if pwhcmp != self.pwhash:
            raise Exception("invalidPWKF")
        ckey = aescalc(self.ckeydata, mkey[16:48], mkey[0:16], False, False)
        tkey = aescalc(self.tkeydata, mkey[64:96], mkey[48:64], False, False)
        oldpath = str(aescalc(self.encname, tkey[16:48], tkey[0:16], False, True), encoding="utf-8")
        if type(data) == bytes:
            newpath = ""
            reader = io.BytesIO(data)
            writer = io.BytesIO()
        elif type(data) == str:
            data = os.path.abspath(data).replace("\\", "/")
            newpath = data[0:data.rfind("/") + 1] + oldpath
            reader = open(data, "rb")
            writer = open(newpath, "wb")
        else:
            raise Exception("invalid datatype")
        
        keys, ivs = [0] * 40, [0] * 40
        for i in range(0, 40):
            temp = 48 * i
            keys[i], ivs[i] = ckey[temp + 16:temp + 48], ckey[temp:temp + 16]
        aesdec( keys, ivs, reader, writer, self.dinfo[0], self.dinfo[1] )
        self.proc = 2.0
        if type(data) == bytes:
            return writer.getvalue()
        else:
            reader.close()
            writer.close()
            return newpath
        
class funcmode:
    def __init__(self):
        self.before = None # reader tgt, bytes/str
        self.after = None # writer tgt, bytes/str
        self.proc = 0.0

    # encrypt with akey 48B
    def encrypt(self, akey):
        self.proc = -1.0
        if len(akey) != 48:
            raise Exception("invalidAKEY")
        ckey = genrand(1920)
        ckeydata = aescalc(ckey, akey[16:48], akey[0:16], True, False)
        keys, ivs = [0] * 40, [0] * 40
        for i in range(0, 40):
            temp = 48 * i
            keys[i], ivs[i] = ckey[temp + 16:temp + 48], ckey[temp:temp + 16]
        if type(self.before) == bytes:
            dsize = len(self.before)
            reader = io.BytesIO(self.before)
            writer = io.BytesIO()
        elif type(self.before) == str:
            dsize = os.path.getsize(self.before)
            reader = open(self.before, "rb")
            writer = open(self.after, "wb")
        else:
            raise Exception("invalid datatype")
        aesenc(keys, ivs, reader, writer, ckeydata, dsize)
        self.proc = 2.0
        if type(self.before) == bytes:
            self.after = writer.getvalue()

    # decrypt with akey 48B
    def decrypt(self, akey):
        self.proc = -1.0
        if len(akey) != 48:
            raise Exception("invalidAKEY")
        if type(self.before) == bytes:
            dsize = len(self.before) - 1920
            reader = io.BytesIO(self.before)
            writer = io.BytesIO()
        elif type(self.before) == str:
            dsize = os.path.getsize(self.before) - 1920
            reader = open(self.before, "rb")
            writer = open(self.after, "wb")
        else:
            raise Exception("invalid datatype")
        ckeydata = reader.read(1920)
        ckey = aescalc(ckeydata, akey[16:48], akey[0:16], False, False)
        keys, ivs = [0] * 40, [0] * 40
        for i in range(0, 40):
            temp = 48 * i
            keys[i], ivs[i] = ckey[temp + 16:temp + 48], ckey[temp:temp + 16]
        aesdec(keys, ivs, reader, writer, 1920, dsize)
        self.proc = 2.0
        if type(self.before) == bytes:
            self.after = writer.getvalue()

# generate secure random nB
def genrand(size):
    return secrets.token_bytes(size)

# returns gen5kaes basic keyfile
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
