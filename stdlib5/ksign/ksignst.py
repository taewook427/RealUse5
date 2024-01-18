# test564 : ksign standard (py)

import os
import threading
import hashlib
from Cryptodome.PublicKey import RSA
from Cryptodome.Signature import pss
from Cryptodome.Hash import SHA3_512

# 파일 해시, 너무 큰 파일은 안됨. /절대경로
def hashfile(path, res):
    with open(path, "rb") as f:
        temp = f.read()
    if temp == b"":
        res[0] = b"\x00" * 64
    else:
        res[0] = hashlib.sha3_512(temp).digest()

# 폴더 해시, 너무 내부 폴더가 많으면 안됨. /절대경로(~/)
def hashfolder(path, res):
    ps = os.listdir(path)
    for i in range( 0, len(ps) ):
        ps[i] = path + ps[i]
        if os.path.isdir( ps[i] ) and ps[-1] != "/":
            ps[i] = ps[i] + "/"
    ps.sort()
    if ps == [ ]:
        res[0] = b"\x00" * 64
    else:
        temp = [b""] * len(ps)
        ret = [ [b""] for x in range( 0, len(ps) ) ]
        th = [ ]
        for i in range( 0, len(ps) ):
            if os.path.isfile( ps[i] ):
                nt = threading.Thread( target = hashfile, args = ( ps[i], ret[i] ) )
                nt.start()
                th.append(nt)
            else:
                hashfolder( ps[i], ret[i] )
        for i in th:
            i.join()
        for i in range( 0, len(ps) ):
            temp[i] = ret[i][0]
        res[0] = hashlib.sha3_512( b"".join(temp) ).digest()
        
# file/folder -> 64B hash !!!복잡한 폴더는 py와 go 결과가 다를 수 있음!!!
def khash(path):
    temp = [b""]
    path = os.path.abspath(path)
    path = path.replace("\\", "/")
    if os.path.isfile(path):
        hashfile(path, temp)
    else:
        if path[-1] != "/":
            path = path + "/"
        hashfolder(path, temp)
    return temp[0]

# gen N byte public, private key (N * 8 bit) -> 2048 : 256, 4096 : 512
def genkey(n):
    temp = RSA.generate(n * 8)
    private = str(temp.export_key(), encoding="utf-8")
    public = str(temp.publickey().export_key(), encoding="utf-8")
    return public, private

# private S + plain 80B -> enc B !!!py-go 호환 안됨!!!
def sign(private, plain):
    if len(plain) != 80:
        raise Exception("PlainV should be 80B")
    key = RSA.import_key( bytes(private, encoding="utf-8") )
    temp = pss.new(key).sign( SHA3_512.new(plain) )
    return temp

# public S + enc B + plain 80B -> T/F (True is ok)
def verify(public, enc, plain):
    if len(plain) != 80:
        raise Exception("PlainV should be 80B")
    key = RSA.import_key( bytes(public, encoding="utf-8") )
    v = pss.new(key)
    h = SHA3_512.new(plain)
    try:
        v.verify(h, enc)
        return True
    except:
        return False

# name S + hashed 64B -> plain 80B
def fm(name, hashed):
    name = bytes(name + " " * 16, encoding="utf-8")
    name = name[0:16]
    if len(hashed) != 64:
        raise Exception("HashV should be 64B")
    return name + hashed
