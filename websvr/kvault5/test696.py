import kaes
import ksc
import kdb

import scrypt
from Cryptodome.Cipher import AES

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

def genpm(pw, kf, salt): # generate pwhash 192B, mkey 144B
    pwh = scrypt.hash(pw + pw + kf + pw + kf, salt, 524288, 8, 1, 192)
    mkey = scrypt.hash(kf + pw + kf + kf + pw, salt, 16384, 8, 1, 144)
    return pwh, mkey

# ===== access data =====
path = "E:/t/"
pw = b"0000"
kf = kaes.basickey()
blocknum = 2
# ===== ===== =====

# block A unpack
t = ksc.toolbox()
t.predetect = True
t.path = path + "0a.webp"
t.readf()
with open(path + "0a.webp", "rb") as f:
    f.seek( t.chunkpos[0] + 8 )
    ptA = f.read( t.chunksize[0] )
    f.seek( t.chunkpos[1] + 8 )
    ptB = f.read( t.chunksize[1] )
    f.seek( t.chunkpos[2] + 8 )
    ptC = f.read( t.chunksize[2] )
    f.seek( t.chunkpos[3] + 8 )
    ptD = f.read( t.chunksize[3] )

# block A part A
t = kdb.toolbox()
t.read( str(ptA, encoding="utf-8") )
salt = t.get("salt")[3]
pwhash = t.get("pwhash")[3]
fsyskdt = t.get("fsyskdt")[3]
fkeykdt = t.get("fkeykdt")[3]
fphykdt = t.get("fphykdt")[3]

# get masterkey
pwh, mkey = genpm(pw, kf, salt)
if pwh != pwhash:
    print("wrong pwkf")

# get session key, decrypt part BCD
t = kaes.funcmode()
fsyskey = aescalc(fsyskdt, mkey[16:48], mkey[0:16], False, False)
t.before = ptB
t.decrypt(fsyskey)
fsys = t.after
fkeykey = aescalc(fkeykdt, mkey[64:96], mkey[48:64], False, False)
t.before = ptC
t.decrypt(fkeykey)
fkey = t.after
fphykey = aescalc(fphykdt, mkey[112:144], mkey[96:112], False, False)
t.before = ptD
t.decrypt(fphykey)
fphy = t.after
with open(path + "_fsys.txt", "wb") as f:
    f.write(fsys)
with open(path + "_fkey.bin", "wb") as f:
    f.write(fkey)
with open(path + "_fphy.bin", "wb") as f:
    f.write(fphy)

# decrypt & rewrite blocks
fphy = fphy[8:]
for i in range(0, blocknum):
    key = b""
    pos = 16 * (int(i / 1) % 256)
    key = key + fphy[pos:pos + 16]
    pos = 16 * (int(i / 256) % 256)
    key = key + fphy[pos:pos + 16]
    pos = 16 * (int(i / 65536) % 256)
    key = key + fphy[pos:pos + 16]

    p0 = path + f"{int(i/256)}/{i%256}c.kv5"
    p1 = path + f"{int(i/256)}/_{i%256}c.kv5"
    with open(p0, "rb") as f:
        with open(p1, "wb") as t:
            t.write( aescalc(f.read(), key[16:48], key[0:16], False, False) )
