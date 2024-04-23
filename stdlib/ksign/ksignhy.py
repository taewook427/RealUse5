# test640 : stdlib5.ksign hy

import ctypes
import hashlib

import kobj

class toolbox:
    def __init__(self): # "windows" / "linux"
        myos = "windows"
        if myos == "windows":
            dll = "./ksign5hy.dll"
            hvalue = b'~i\xbfB\x95Nfv\xf8c\xcf\xc5\xdf\xfd\xe9\x18\xdf\x10\xfa\xf8\xf3\x99\r\x17bT\x00\xee\xf3Q_y\x82\xe1S=\x8fN\r\\\xecG]#+@\xb4\x80\xdb\xaf\xb9\x05E\xd1\x00\xf2\x07-h\xb8\xab\x9a\xfaS'
        else:
            dll = "./ksign5hy.so"
            hvalue = b'\x87s\xd4\r\x07\xe7?\xec\xb2\xe6x\xa1?o\n\xa5\x93`\x9e\x1dl|1\xb6m\x83\xa4E\x84\xf3\x0cY\xd6\xb5\xed\x8dF\x87\x8a\xf8\xfe\x9f!k\xd5\xc03\x06P\x874%\xba\x84\x7f\x90d\x98\x8aC\xf9\xb6>\x06'
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
