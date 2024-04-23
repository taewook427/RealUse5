# test628 : stdlib5.kobj

import os
import sys
import ctypes

# re-alloc current path, returns cmd args
def repath():
    temp = sys.argv
    path = os.path.abspath( temp[0] )
    path = path.replace("\\", "/")
    os.chdir( path[ 0:path.rfind("/") ] )
    return temp

# little endian encoding, I -> B
def encode(num, length):
    temp = [0] * length
    for i in range(0, length):
        temp[i] = num % 256
        num = num // 256
    return bytes(temp)

# little endian decoding, B -> I
def decode(data):
    temp = 0
    for i in range( 0, len(data) ):
        if data[i] != 0:
            temp = temp + data[i] * 256 ** i
    return temp

# package series of bytes, B[] -> B{1B len + (8B size + nB data) * n}
def pack(series):
    temp = [b""] * (2 * len(series) + 1)
    temp[0] = bytes( [ len(series) ] )
    for i in range( 0, len(series) ):
        temp[2 * i + 1] = encode(len( series[i] ), 8)
        temp[2 * i + 2] = series[i]
    return b"".join(temp)

# unpack packed B, B{1B len + (8B size + nB data) * n} -> B[]
def unpack(chunk):
    length = chunk[0]
    temp = [b""] * length
    ptr = 1
    for i in range(0, length):
        clen = decode( chunk[ptr:ptr + 8] )
        ptr = ptr + 8
        temp[i] = chunk[ptr:ptr + clen]
        ptr = ptr + clen
    return temp

# B -> (c_char_p, length int)
def send(data):
    length = len(data)
    arr = ctypes.c_char_p(data)
    return arr, length

# cptr(nB data) -> B
def recv(reader, length):
    temp = bytearray(length)
    for i in range(0, length):
        temp[i] = reader[i][0]
    return bytes(temp)

# cptr(8B len + nB data) -> B
def recvauto(reader):
    temp = bytearray(8)
    for i in range(0, 8):
        temp[i] = reader[i][0]
    length = decode( bytes(temp) )
    temp = bytearray(length)
    for i in range(0, length):
        temp[i] = reader[i + 8][0]
    return bytes(temp)

# dll/so func arg/ret types, (types, types), { int(i), float(f), charptr(b) }
def call(args, rets):
    argv, retv = [ ], None
    for i in args:
        if i == "i":
            argv.append(ctypes.c_int)
        elif i == "f":
            argv.append(ctypes.c_float)
        elif i == "b":
            argv.append( ctypes.POINTER(ctypes.c_char) )
    argv = tuple(argv)
    if rets == "i":
        retv = ctypes.c_int
    elif rets == "f":
        retv = ctypes.c_float
    elif rets == "b":
        retv = ctypes.POINTER(ctypes.c_char)
    return argv, retv
