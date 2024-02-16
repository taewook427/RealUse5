# test566 : ksign hybrid (py)

import ctypes
import hashlib

class toolbox:
    def __init__(self):
        myos = "windows" # "windows" / "linux"
        if myos == "windows":
            dll = "./ksignhy.dll"
            hvalue = b'\xbe\x16\xf6S\xd4\x95\xd4?\t\x9c\xcd6x\xaf\xe1\xad@\xb5e\x1b\x1cM\xeb\xa1N\xa9gA\xcd\xa2\x18\x17\xbf2\xdc\x1a\x94w\x8ep\xc6\xdc\xae\xa5\x11\xecq\x108N\xfc\xd0\xdf\xe6\xe9\x01\xf9\xe7\x10\xe0D(\xfa\xcc'
        else:
            dll = "./ksignhy.so"
            hvalue = b'\xf8|\x9a\xf4;\x94\xaa\xa2-\x94\xe6\xdbLB\xbf\x8f\x99v\xda$\xf7\xaf\x18&\xd9z\x8b\x02\xa7\x9b1\x7fbv\x1e\xf7\xf8G >\xc99\x06\x10\x96\xce)\xc8J\xd1\xb2j\x9f\xb9\n\x9f\xe1\x88L\xd2O\xed\xc8\xec'
        with open(dll, "rb") as f:
            temp = f.read()
        if hashlib.sha3_512(temp).digest() != hvalue:
            raise Exception("wrong ffi")

        if myos == "windows":
            self.ext = ctypes.CDLL(dll)
        else:
            self.ext = ctypes.cdll.LoadLibrary(dll)
        # path str B + B len N -> hash 64B
        self.ext.khashhy.argtype = (ctypes.POINTER(ctypes.c_char), ctypes.c_int)
        self.ext.khashhy.restype = ctypes.POINTER(ctypes.c_char)
        # len B N -> len 2B + public nB + len 2B + private nB
        self.ext.genkeyhy.argtype = ctypes.c_int
        self.ext.genkeyhy.restype = ctypes.POINTER(ctypes.c_char)
        # private B + len N + plain 80B-> 2B enc B
        self.ext.signhy.argtype = (ctypes.POINTER(ctypes.c_char), ctypes.c_int, ctypes.POINTER(ctypes.c_char))
        self.ext.signhy.restype = ctypes.POINTER(ctypes.c_char)
        # public B + len N + enc B + len N + plain 80B -> T(1)/F(0) (True is ok)
        self.ext.verifyhy.argtype = ( ctypes.POINTER(ctypes.c_char), ctypes.c_int, ctypes.POINTER(ctypes.c_char), ctypes.c_int, ctypes.POINTER(ctypes.c_char) )
        self.ext.verifyhy.restype = ctypes.c_int
        # free char*
        self.ext.freehy.argtype = (ctypes.POINTER(ctypes.c_char),)

    # path str -> hash B
    def khash(self, path):
        path = bytes(path.replace("\\", "/"), encoding="utf-8")
        l = len(path)
        arr = ctypes.c_char_p(path)
        ptr = self.ext.khashhy(arr, l)
        out = [0] * 64
        for i in range(0, 64):
            out[i] = ptr[i][0]
        self.ext.freehy(ptr)
        return bytes(out)

    # n int -> public str, private str
    def genkey(self, n):
        ptr = self.ext.genkeyhy(n)
        num = 0
        ra = ptr[num + 0][0]
        rb = ptr[num + 1][0]
        num = num + 2
        size = ra + 256 * rb
        out = [0] * size
        for i in range(0, size):
            out[i] = ptr[num][0]
            num = num + 1
        public = str(bytes(out), encoding="utf-8")
        ra = ptr[num + 0][0]
        rb = ptr[num + 1][0]
        num = num + 2
        size = ra + 256 * rb
        out = [0] * size
        for i in range(0, size):
            out[i] = ptr[num][0]
            num = num + 1
        private = str(bytes(out), encoding="utf-8")
        self.ext.freehy(ptr)
        return public, private

    # private str, plain B -> enc B
    def sign(self, private, plain):
        if len(plain) != 80:
            raise Exception("PlainV should be 80B")
        private = bytes(private, encoding="utf-8")
        l = len(private)
        arr0 = ctypes.c_char_p(private)
        arr1 = ctypes.c_char_p(plain)
        ptr = self.ext.signhy(arr0, l, arr1)
        num = 0
        ra = ptr[num + 0][0]
        rb = ptr[num + 1][0]
        num = num + 2
        size = ra + 256 * rb
        out = [0] * size
        for i in range(0, size):
            out[i] = ptr[num][0]
            num = num + 1
        enc = bytes(out)
        self.ext.freehy(ptr)
        return enc

    # public str, enc B, plain B -> isvalid bool
    def verify(self, public, enc, plain):
        if len(plain) != 80:
            raise Exception("PlainV should be 80B")
        public = bytes(public, encoding="utf-8")
        arr0 = ctypes.c_char_p(public)
        l0 = len(public)
        arr1 = ctypes.c_char_p(enc)
        l1 = len(enc)
        arr2 = ctypes.c_char_p(plain)
        out = self.ext.verifyhy(arr0, l0, arr1, l1, arr2)
        if out == 0:
            return False
        else:
            return True

    # name str + hashed B -> plain B
    def fm(self, name, hashed):
        bnm = bytes(name + " " * 16, encoding="utf-8")
        bnm = bnm[0:16]
        if len(hashed) != 64:
            raise Exception("HashV should be 64B")
        return bnm + hashed
