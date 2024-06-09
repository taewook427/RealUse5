# test636 : stdlib5.kaes hy

import io
import os
import time
import ctypes
import threading as thr

import hashlib
import secrets

import kobj
import kdb
import ksc
import ksign
import picdt

class premode:
    def __init__(self): # "windows" / "linux"
        myos = "windows"
        if myos == "windows":
            dll = "./kaes5hy.dll"
            hv = b'\x87e\xa9Xy\x8c\xb1OL\x86\x9blQW\xc9\xc0OZ\x95\x06P\xee(\x10\xd1V\x7f}\xa6\x92fY\xc9\xe2Z\x05S\x12(\x9f\xe9V\x8e\x89\xeel\xe1Yl\x1eod>\\\xe2|\x90\x06\xf7\xbe\x88"\x0b\xac'
        else:
            dll = "./kaes5hy.so"
            hv = b'\xdf\x94E]\x8br\xee|\x15pU\xeeUUY?\x0c\xc0\xf8\xb6\x06S\xbe\x8b<\x85\xebz\x9f\xf8k\xa9@\xab\x01O\xf6\x7f8\x168\x17\xeas\xca\xefT]\xea\xde\xf9}G\x06\x86@\x02Q\xc2NV\xb9\xf4^'
        with open(dll, "rb") as f:
            hc = hashlib.sha3_512( f.read() ).digest()
        if hc != hv:
            raise Exception("wrong FFI")
        
        if myos == "windows":
            self.ext = ctypes.CDLL(dll)
        else:
            self.ext = ctypes.cdll.LoadLibrary(dll)
        self.ext.func0.argtypes, self.ext.func0.restype = kobj.call("b", "") # free
        self.ext.func1.argtypes, self.ext.func1.restype = kobj.call("", "b") # get basickey
        self.ext.func2.argtypes, self.ext.func2.restype = kobj.call("i", "f") # set/get proc
        self.ext.func3.argtypes, self.ext.func3.restype = kobj.call("bibibi", "b") # get pm
        self.ext.func4.argtypes, self.ext.func4.restype = kobj.call("bibii", "b") # aes calc
        self.ext.func5.argtypes, self.ext.func5.restype = kobj.call("bibii", "b") # aes chunk
        self.ext.func6.argtypes, self.ext.func6.restype = kobj.call("bi", "b") # all-mode file encrypt
        self.ext.func7.argtypes, self.ext.func7.restype = kobj.call("bi", "b") # all-mode file view
        self.ext.func8.argtypes, self.ext.func8.restype = kobj.call("bi", "b") # all-mode file decrypt
        self.ext.func9.argtypes, self.ext.func9.restype = kobj.call("bi", "b") # func-mode file ende

        self.proc = 0.0 # -1 : not started, 0~1 : working, 2 : end

    def hash3512(self, data): # sha3-512 hash, data bytes
        return hashlib.sha3_512(data).digest()
    
    def aescalc(self, data, key, iv, isenc, ispad): # aes-256 ende, key 32B iv 16B, isenc T:encrypt F:decrypt, ispad T:dopad F:nopad
        if len(key) != 32 or len(iv) != 16:
            raise Exception("invalid keyiv")
        p0, p1 = kobj.send(data)
        p2, p3 = kobj.send(iv + key)
        p4 = 0 if isenc else 2
        if not ispad:
            p4 = p4 + 1

        r0 = self.ext.func4(p0, p1, p2, p3, p4)
        data = kobj.recvauto(r0)
        self.ext.func0(r0)
        return data
    
    def aesenc(self, ckey, reader, writer, header, dsize): # Bmode encrypt, ckey 1920B / rw Bio / header bytes / dsize int
        if len(ckey) != 1920:
            raise Exception("invalid keyiv")
        writer.write(header)
        num0 = dsize // 104857600
        num1 = dsize % 104857600
        num2 = num1 // 524288
        num3 = num1 % 524288

        self.proc = 0.0
        for i in range(0, num0):
            p0, p1 = kobj.send(ckey)
            p2, p3 = kobj.send( reader.read(104857600) )

            r0 = self.ext.func5(p0, p1, p2, p3, 0)
            temp = kobj.recv(r0, 104857600 + 1920)
            self.ext.func0(r0)

            ckey = temp[0:1920]
            writer.write( temp[1920:104857600 + 1920] )
            self.proc = i / num0

        if num2 != 0:
            p0, p1 = kobj.send(ckey)
            p2, p3 = kobj.send( reader.read(524288 * num2) )

            r0 = self.ext.func5(p0, p1, p2, p3, 0)
            temp = kobj.recv(r0, 524288 * num2 + 1920)
            self.ext.func0(r0)

            ckey = temp[0:1920]
            writer.write( temp[1920:524288 * num2 + 1920] )

        ti = 48 * (num2 % 40)
        temp = self.aescalc(reader.read(num3), ckey[ti + 16:ti + 48], ckey[ti:ti + 16], True, True)
        writer.write(temp)

    def aesdec(self, ckey, reader, writer, stpoint, dsize): # Bmode decrypt, ckey 1920B / rw Bio / stpoint int / dsize int
        if len(ckey) != 1920:
            raise Exception("invalid keyiv")
        reader.seek(stpoint)
        num0 = dsize // 104857600
        num1 = dsize % 104857600
        if num1 == 0:
            num0 = num0 - 1
            num1 = 104857600
        num2 = num1 // 524288
        num3 = num1 % 524288
        if num3 == 0:
            num2 = num2 - 1
            num3 = 524288

        self.proc = 0.0
        for i in range(0, num0):
            p0, p1 = kobj.send(ckey)
            p2, p3 = kobj.send( reader.read(104857600) )

            r0 = self.ext.func5(p0, p1, p2, p3, 1)
            temp = kobj.recv(r0, 104857600 + 1920)
            self.ext.func0(r0)

            ckey = temp[0:1920]
            writer.write( temp[1920:104857600 + 1920] )
            self.proc = i / num0

        if num2 != 0:
            p0, p1 = kobj.send(ckey)
            p2, p3 = kobj.send( reader.read(524288 * num2) )

            r0 = self.ext.func5(p0, p1, p2, p3, 1)
            temp = kobj.recv(r0, 524288 * num2 + 1920)
            self.ext.func0(r0)

            ckey = temp[0:1920]
            writer.write( temp[1920:524288 * num2 + 1920] )

        ti = 48 * (num2 % 40)
        temp = self.aescalc(reader.read(num3), ckey[ti + 16:ti + 48], ckey[ti:ti + 16], False, True)
        writer.write(temp)

    def genpm(self, pw, kf, salt): # generate pwhash 128B, mkey 96B
        p0, p1 = kobj.send(pw)
        p2, p3 = kobj.send(kf)
        p4, p5 = kobj.send(salt)

        r0 = self.ext.func3(p0, p1, p2, p3, p4, p5)
        temp = kobj.recv(r0, 224)
        self.ext.func0(r0)
        return temp[0:128], temp[128:224]
    
    def gensize(self, n): # generate size when n bytes is pad encrypted
        return n + 16 - n % 16

    def genpath(self, size): # generate random word with length 2*size
        temp = self.genrand(size)
        word = ""
        for i in temp:
            word = word + str( hex(i + 16) )[2:4]
        return word
    
    def genproc(self): # update proc until reach 1+
        while True:
            time.sleep(0.1)
            self.proc = self.ext.func2(1)
            if self.proc > 1.1:
                break
    
    # generate secure random nB
    def genrand(self, size):
        return secrets.token_bytes(size)
    
    # returns gen5kaes basic keyfile
    def basickey(self):
        r0 = self.ext.func1()
        data = kobj.recvauto(r0)
        self.ext.func0(r0)
        return data

class allmode(premode):
    def __init__(self):
        super().__init__()
        self.hint = "" # pwkf hint str
        self.msg = "" # program msg str
        self.signkey = ["", ""] # (pub, pri) ksign rsa key

        self.salt = b""
        self.pwhash = b""
        self.ckeydata = b""
        self.tkeydata = b""
        self.encname = b""
        self.dinfo = [0, 0] # stpoint, encdsize

    def mkhead(self, pmode, dsize): # make real header
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
            w0.fix( "signdata", w1.sign( self.signkey[1], self.hash3512(encheader) ) )
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
        header = header + ksc.encode(self.gensize(dsize), 8)
        return header

    def rdhead(self, encheader, signheader): # read header, update value
        self.signkey = ["", ""]
        if len(signheader) != 0:
            worker, temp = kdb.toolbox(), ksign.toolbox()
            worker.read( str(signheader, encoding="utf-8") )
            self.signkey = [worker.get("publickey")[3], ""]
            signdata = worker.get("signdata")[3]
            if not temp.verify( self.signkey[0], signdata, self.hash3512(encheader) ):
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

    # encryption, pwkf bytes / data bytes / pmode 0:webp 1:png 2:none -> result bytes / data str -> encpath str
    def encrypt(self, pw, kf, data, pmode):
        self.proc = -1.0
        out = None # return value

        if type(data) == str:
            temp = [0] * 8
            temp[0] = bytes(self.hint, encoding="utf-8")
            temp[1] = bytes(self.msg, encoding="utf-8")
            temp[2] = bytes(self.signkey[0], encoding="utf-8")
            temp[3] = bytes(self.signkey[1], encoding="utf-8")
            temp[4] = pw
            temp[5] = kf
            temp[6] = bytes(data, encoding="utf-8")
            temp[7] = kobj.encode(pmode, 8)
            p0, p1 = kobj.send( kobj.pack(temp) )

            self.ext.func2(0)
            t = thr.Thread(target=self.genproc)
            t.start()

            r0 = self.ext.func6(p0, p1)
            data = kobj.unpack( kobj.recvauto(r0) )
            self.ext.func0(r0)
            if len( data[0] ) != 0:
                raise Exception( str(data[0], encoding="utf-8") )
            out = str(data[1], encoding="utf-8")

        elif type(data) == bytes:
            ckey = self.genrand(1920) # (iv16 key32) * 40
            tkey = self.genrand(48) # iv16 key32
            dsize = len(data)
            reader = io.BytesIO(data)
            writer = io.BytesIO()

            self.salt = self.genrand(40)
            self.pwhash, mkey = self.genpm(pw, kf, self.salt)
            self.ckeydata = self.aescalc(ckey, mkey[16:48], mkey[0:16], True, False)
            self.tkeydata = self.aescalc(tkey, mkey[64:96], mkey[48:64], True, False)
            self.encname = self.aescalc(b"NewData.bin", tkey[16:48], tkey[0:16], True, True)

            header = self.mkhead(pmode, dsize) # real header to write
            self.aesenc(ckey, reader, writer, header, dsize)
            writer.write(b"\xff" * 8)
            out = writer.getvalue()

        else:
            raise Exception("invalid datatype")
        self.proc = 2.0
        return out

    # view encfile, data bytes/str
    def view(self, data):
        self.proc = -1.0

        if type(data) == str:
            p0, p1 = kobj.send( bytes(data, encoding="utf-8") )

            r0 = self.ext.func7(p0, p1)
            data = kobj.unpack( kobj.recvauto(r0) )
            self.ext.func0(r0)
            if len( data[0] ) != 0:
                raise Exception( str(data[0], encoding="utf-8") )

            self.hint = str(data[1], encoding="utf-8")
            self.msg = str(data[2], encoding="utf-8")
            self.signkey[0] = str(data[3], encoding="utf-8")
            self.signkey[1] = ""

        elif type(data) == bytes:
            worker = ksc.toolbox()
            worker.predetect = True
            worker.readb(data)
            if worker.subtype != b"KAES":
                raise Exception("invalidKAES5")

            encheader = data[worker.chunkpos[0] + 8:worker.chunkpos[0] + worker.chunksize[0] + 8]
            signheader = data[worker.chunkpos[1] + 8:worker.chunkpos[1] + worker.chunksize[1] + 8]
            self.dinfo = [ worker.chunkpos[2] + 8, worker.chunksize[2] ]
            if worker.reserved != ksc.crc32hash(encheader) + ksc.crc32hash(signheader):
                raise Exception("invalidCRC32")

            self.rdhead(encheader, signheader)

        else:
            raise Exception("invalid datatype")

    # decryption, pwkf bytes / data bytes -> result bytes / data str -> decpath str
    def decrypt(self, pw, kf, data):
        self.proc = -1.0
        out = None # return value

        if type(data) == str:
            temp = [0] * 3
            temp[0] = pw
            temp[1] = kf
            temp[2] = bytes(data, encoding="utf-8")
            p0, p1 = kobj.send( kobj.pack(temp) )

            self.ext.func2(0)
            t = thr.Thread(target=self.genproc)
            t.start()

            r0 = self.ext.func8(p0, p1)
            data = kobj.unpack( kobj.recvauto(r0) )
            self.ext.func0(r0)
            if len( data[0] ) != 0:
                raise Exception( str(data[0], encoding="utf-8") )
            out = str(data[1], encoding="utf-8")

        elif type(data) == bytes:
            if len(self.salt) == 0:
                raise Exception("should done view() first")
            pwhcmp, mkey = self.genpm(pw, kf, self.salt)
            if pwhcmp != self.pwhash:
                raise Exception("invalidPWKF")
            ckey = self.aescalc(self.ckeydata, mkey[16:48], mkey[0:16], False, False)

            reader = io.BytesIO(data)
            writer = io.BytesIO()
            self.aesdec( ckey, reader, writer, self.dinfo[0], self.dinfo[1] )
            out = writer.getvalue()

        else:
            raise Exception("invalid datatype")
        self.proc = 2.0
        return out

class funcmode(premode):
    def __init__(self):
        super().__init__()
        self.before = None # reader tgt, bytes/str
        self.after = None # writer tgt, bytes/str

    # encrypt with akey 48B
    def encrypt(self, akey):
        self.proc = -1.0
        if len(akey) != 48:
            raise Exception("invalidAKEY")
        
        if type(self.before) == str:
            temp = [0] * 4
            temp[0] = bytes(self.before, encoding="utf-8")
            temp[1] = bytes(self.after, encoding="utf-8")
            temp[2] = akey
            temp[3] = kobj.encode(0, 8)
            p0, p1 = kobj.send( kobj.pack(temp) )

            self.ext.func2(0)
            t = thr.Thread(target=self.genproc)
            t.start()

            r0 = self.ext.func9(p0, p1)
            data = kobj.recvauto(r0)
            self.ext.func0(r0)
            if len(data) != 0:
                raise Exception( str(data, encoding="utf-8") )

        elif type(self.before) == bytes:
            ckey = self.genrand(1920)
            ckeydata = self.aescalc(ckey, akey[16:48], akey[0:16], True, False)
            dsize = len(self.before)
            reader = io.BytesIO(self.before)
            writer = io.BytesIO()
            self.aesenc(ckey, reader, writer, ckeydata, dsize)
            self.after = writer.getvalue()

        else:
            raise Exception("invalid datatype")
        self.proc = 2.0

    # decrypt with akey 48B
    def decrypt(self, akey):
        self.proc = -1.0
        if len(akey) != 48:
            raise Exception("invalidAKEY")
        
        if type(self.before) == str:
            temp = [0] * 4
            temp[0] = bytes(self.before, encoding="utf-8")
            temp[1] = bytes(self.after, encoding="utf-8")
            temp[2] = akey
            temp[3] = kobj.encode(1, 8)
            p0, p1 = kobj.send( kobj.pack(temp) )

            self.ext.func2(0)
            t = thr.Thread(target=self.genproc)
            t.start()

            r0 = self.ext.func9(p0, p1)
            data = kobj.recvauto(r0)
            self.ext.func0(r0)
            if len(data) != 0:
                raise Exception( str(data, encoding="utf-8") )
            
        elif type(self.before) == bytes:
            dsize = len(self.before) - 1920
            reader = io.BytesIO(self.before)
            writer = io.BytesIO()
            ckeydata = reader.read(1920)
            ckey = self.aescalc(ckeydata, akey[16:48], akey[0:16], False, False)
            self.aesdec(ckey, reader, writer, 1920, dsize)
            self.after = writer.getvalue()

        else:
            raise Exception("invalid datatype")
        self.proc = 2.0
