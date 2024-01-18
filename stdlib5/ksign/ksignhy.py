# test566 : ksign hybrid (py)

import ctypes
import hashlib

class toolbox:
    def __init__(self):
        # windows
        dll = "./ksignhy.dll"
        hvalue = b'\x16B6k\x18\xb4x\x8e\xb3j\xd4x\xfa\x12\x16\x81\r\xa3\xe1\x9b\xde\x05\x86\x03E\xb4;,\x04\x95\x9fL\x96H\xcc\xc9}\xbf\xf4v\x8f\xa5\xf1\xaa\r\x99\x05h/\x83I\x02\xbe+\r?\xb5k\xdb\x1eH\xc6\xc1\x19'
        # linux
        #dll = "./ksignhy.so"
        #hvalue = b'\xe91#Xa\x94\x84\xae\x9c(\xca\xce$\x98R\xfa\xd2\xcc-\xd9\xff\xe1\xf7\xeb\x8f\x9b\xdb\x12\x06tp\x86\xadp\xda#\x1d\x8dm\x89\n\x91\x16}\xe7G\xb8\xef\x0e\xae\r\x0f_`\xb5\xd6\x80(Re\xcd\x12\x08I'
        with open(dll, "rb") as f:
            temp = f.read()
        if hashlib.sha3_512(temp).digest() != hvalue:
            raise Exception("wrong ffi")

        # windows
        self.ext = ctypes.CDLL(dll)
        # linux
        #self.ext = ctypes.cdll.LoadLibrary(dll)
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
