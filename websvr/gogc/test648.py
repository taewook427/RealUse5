"""
import time
import zlib

def read(path):
    with open(path, "rb") as f:
        data = f.read()
    return data

def checksum(data):
    crcv = zlib.crc32(data)
    res = [0] * 4
    for i in range(0, 4):
        res[i] = crcv % 256
        crcv = crcv // 256
    return res

# def read(path: str) -> bytes: (...)
# def checksum(data: bytes) -> list[int]: (...)

for i in range(0, 4):
    print( checksum( read("./big.bin") ) )
    time.sleep(2)
"""

import ctypes
import time

def send(data):
    arr = ctypes.c_char_p(data)
    return arr

def recv(ptr, length):
    temp = [0] * length
    for i in range(0, length):
        temp[i] = ptr[i][0]
    return bytes(temp)

dll = ctypes.CDLL("./ex.dll")
dll.freeptr.argtypes = (ctypes.POINTER(ctypes.c_char),)
dll.work.argtypes = (ctypes.POINTER(ctypes.c_char), ctypes.c_int)
dll.work.restype = ctypes.POINTER(ctypes.c_char)

# def send(data: bytes) -> c.char.ptr: (...)
# def recv(ptr: c,char.ptr, length: int) -> bytes: (...)
# dll: ctypes.dll; dll.freeptr: ctypes.func; dll.work: ctypes.func

inarr = [0, 1, 2, 3, 4, 5, 6, 7] * (1024 * 1024 * 120)
time.sleep(2)

inarr = bytes(inarr)
time.sleep(2)

parm0, parm1 = send(inarr), len(inarr)
time.sleep(2)

del inarr
time.sleep(2)

ptr = dll.work(parm0, parm1)
time.sleep(2)

del parm0
time.sleep(2)

resarr = recv(ptr, parm1)
print( resarr[0:32] )
time.sleep(2)

del resarr
time.sleep(2)

dll.freeptr(ptr)
time.sleep(2)
