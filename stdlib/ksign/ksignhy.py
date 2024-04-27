# test640 : stdlib5.ksign hy

import ctypes
import hashlib

import kobj

class toolbox:
    def __init__(self): # "windows" / "linux"
        myos = "windows"
        if myos == "windows":
            dll = "./ksign5hy.dll"
            hvalue = b'\xa2G<:\xd5=\xacT\x8e\xbb\x96f\xc6\xcaSDq3\x18M4\x0f\x14KB>\xfeb\xa2\x9fT\xfd\xcc\x16\xfe\xb5\x89\xaf\xb7\xf7\xd9I7\xe7\x18\xda\xde\xc3\x16\xcf\xc0xw\\\xc6\x1b\xee\xf3\x9b\xae\xbc\xb4\xe8X'
        else:
            dll = "./ksign5hy.so"
            hvalue = b'\x99_\x9c\x95\xfd"\'\x14\x84\n\xb9\xc2\x90\r\xad\x1c\xa2 \xa5\xe2\x8d+\xf8\x9b\x0fub\xc5\xbe\x9f\'\x1d}zh\x94dX\x82A\x8f\x1dI\xe6\xc5o\xc3Q\x8aM[\xd3\x94\xbe\xaf\x8c\xe9\xf2\x97\x1d\xff\xb7\xf1\x83'
        with open(dll, "rb") as f:
            hcheck = hashlib.sha3_512( f.read() ).digest()
        if hcheck != hvalue:
            raise Exception("wrong FFI")
        
        if myos == "windows":
            self.ext = ctypes.CDLL(dll)
        else:
            self.ext = ctypes.cdll.LoadLibrary(dll)
        args, rets = kobj.call("b", "")
        self.ext.func0.argtypes, self.ext.func0.restype = args, rets # free
        args, rets = kobj.call("bi", "b")
        self.ext.func1.argtypes, self.ext.func1.restype = args, rets # khash
        args, rets = kobj.call("bi", "b")
        self.ext.func2.argtypes, self.ext.func2.restype = args, rets # kinfo
        args, rets = kobj.call("i", "b")
        self.ext.func3.argtypes, self.ext.func3.restype = args, rets # genkey
        args, rets = kobj.call("bibi", "b")
        self.ext.func4.argtypes, self.ext.func4.restype = args, rets # sign
        args, rets = kobj.call("bibibi", "b")
        self.ext.func5.argtypes, self.ext.func5.restype = args, rets # verify

    # file/folder -> 64B khash value (general path)
    def khash(self, path):
        p0, p1 = kobj.send( bytes(path, encoding="utf-8") )
        o0 = self.ext.func1(p0, p1)
        v0 = kobj.unpack( kobj.recvauto(o0) )
        self.ext.func0(o0)
        if v0[1] == b"":
            return v0[0]
        else:
            raise Exception( str(v0[1], encoding="utf-8") )
        
    # get folder/file num, size info (general path) (size, file, folder)
    def kinfo(self, path):
        p0, p1 = kobj.send( bytes(path, encoding="utf-8") )
        o0 = self.ext.func2(p0, p1)
        v0 = kobj.unpack( kobj.recvauto(o0) )
        self.ext.func0(o0)
        if v0[3] == b"":
            return kobj.decode( v0[0] ), kobj.decode( v0[1] ), kobj.decode( v0[2] )
        else:
            raise Exception( str(v0[3], encoding="utf-8") )
        
    # gen N bit public, private key 2048/4096/8192
    def genkey(self, n):
        o0 = self.ext.func3(n)
        v0 = kobj.unpack( kobj.recvauto(o0) )
        self.ext.func0(o0)
        if v0[2] == b"":
            return str(v0[0], encoding="utf-8"), str(v0[1], encoding="utf-8")
        else:
            raise Exception( str(v0[2], encoding="utf-8") )

    # private S + plain nB -> enc B
    def sign(self, private, plain):
        p0, p1 = kobj.send( bytes(private, encoding="utf-8") )
        p2, p3 = kobj.send(plain)
        o0 = self.ext.func4(p0, p1, p2, p3)
        v0 = kobj.unpack( kobj.recvauto(o0) )
        self.ext.func0(o0)
        if v0[1] == b"":
            return v0[0]
        else:
            raise Exception( str(v0[1], encoding="utf-8") )
        
    # public S + enc B + plain nB -> T/F (True is ok)
    def verify(self, public, enc, plain):
        p0, p1 = kobj.send( bytes(public, encoding="utf-8") )
        p2, p3 = kobj.send(enc)
        p4, p5 = kobj.send(plain)
        o0 = self.ext.func5(p0, p1, p2, p3, p4, p5)
        v0 = kobj.unpack( kobj.recvauto(o0) )
        self.ext.func0(o0)
        if v0[0] == b"P":
            return True
        elif v0[0] == b"F":
            return False
        else:
            raise Exception( str(v0[0], encoding="utf-8") )

# format : PKIX (public), PKCS1 (private), PEM, PSS
