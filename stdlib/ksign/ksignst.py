# test638 : stdlib5.ksign st

import os
import threading
import hashlib

from Cryptodome.PublicKey import RSA
from Cryptodome.Signature import pss
from Cryptodome.Hash import SHA3_512

def hashf(path, res): # file pash (abspath)
    with open(path, "rb") as f:
        temp = f.read()
    if temp == b"":
        res[0] = b"\x00" * 64
    else:
        res[0] = hashlib.sha3_512(temp).digest()

def hashd(path, res): # folder hash (abspath/)
    ps = [ path + x for x in os.listdir(path) ]
    for i in range( 0, len(ps) ):
        if os.path.isdir( ps[i] ) and ps[i][-1] != "/":
            ps[i] = ps[i] + "/"
        ps[i] = bytes(ps[i], encoding="utf-8")
    ps.sort()
    if ps == [ ]:
        res[0] = b"\x00" * 64
    else:
        temp = [b""] * len(ps)
        thr = [0] * len(ps)
        ret = [ [b""] for x in range( 0, len(ps) ) ]
        for i in range( 0, len(ps) ):
            ps[i] = str(ps[i], encoding="utf-8")
            if ps[i][-1] == "/":
                thr[i] = threading.Thread( target=hashd, args=( ps[i], ret[i] ) )
            else:
                thr[i] = threading.Thread( target=hashf, args=( ps[i], ret[i] ) )
            thr[i].start()
        for i in range( 0, len(ps) ):
            thr[i].join()
            temp[i] = ret[i][0]
        res[0] = hashlib.sha3_512( b"".join(temp) ).digest()

def infod(path): # get folder info (abspath/)
    size, file, folder = 0, 0, 1
    temp = [ path + x for x in os.listdir(path) ]
    for i in range( 0, len(temp) ):
        if os.path.isdir( temp[i] ) and temp[i][-1] != "/":
            temp[i] = temp[i] + "/"
        if temp[i][-1] == "/":
            t0, t1, t2 = infod( temp[i] )
            size, file, folder = size + t0, file + t1, folder + t2
        else:
            size, file = size + os.path.getsize( temp[i] ), file + 1
    return size, file, folder

# file/folder -> 64B khash value (general path)
def khash(path):
    ret = [b""]
    path = os.path.abspath(path).replace("\\", "/")
    if os.path.isdir(path):
        if path[-1] != "/":
            path = path + "/"
        hashd(path, ret)
    else:
        hashf(path, ret)
    return ret[0]

# get folder/file num, size info (general path) (size, file, folder)
def kinfo(path):
    path = os.path.abspath(path).replace("\\", "/")
    if os.path.isdir(path):
        if path[-1] != "/":
            path = path + "/"
        return infod(path)
    else:
        return os.path.getsize(path), 1, 0

# gen N bit public, private key 2048/4096/8192
def genkey(n):
    temp = RSA.generate(n)
    private = str(temp.export_key(format='PEM', pkcs=1), encoding="utf-8")
    public = str(temp.publickey().export_key(format='PEM'), encoding="utf-8")
    return public, private

# private S + plain nB -> enc B
def sign(private, plain):
    key = RSA.import_key( bytes(private, encoding="utf-8") )
    temp = pss.new(key).sign( SHA3_512.new(plain) )
    return temp

# public S + enc B + plain nB -> T/F (True is ok)
def verify(public, enc, plain): # !!!cannot verify sign gen by golang!!!
    key = RSA.import_key( bytes(public, encoding="utf-8") )
    v = pss.new(key)
    h = SHA3_512.new(plain)
    try:
        v.verify(h, enc)
        return True
    except:
        return False

# format : PKIX (public), PKCS1 (private), PEM, PSS
