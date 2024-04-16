# test644 : stdlib5.ksc

import zlib
import os

class toolbox:
    def __init__(self):
        self.prehead = webpbase() # prehead + pad 512nB
        self.common = b"KSC5" # common sign 4B
        self.subtype = b"\x00\x00\x00\x00" # subtype sign 4B
        self.reserved = b"\x00\x00\x00\x00\x00\x00\x00\x00" # reserved 8B

        self.headp = 512 # mainheader start point (512n)
        self.rsize = 0 # valid data size
        self.path = "" # file path

        self.predetect = False # read chunk info ahead
        self.chunkpos = [ ] # chunk (8B + nB) pos
        self.chunksize = [ ] # chunk data (nB) size

    # read by file path
    def readf(self):
        if self.path == "":
            raise Exception("NoSuchFile")
        self.prehead, self.headp, size = b"", 0, os.path.getsize(self.path)

        with open(self.path, "rb") as f:
            while size >= self.headp + 4:
                temp = f.read(4)
                if temp == self.common:
                    break
                else:
                    self.prehead = self.prehead + temp + f.read(508)
                    self.headp = self.headp + 512
            if size < self.headp + 16:
                raise Exception("InvalidKSC5")
            
            self.subtype = f.read(4)
            self.reserved = f.read(8)
            self.rsize = size - self.headp

            self.chunkpos, self.chunksize, pos = [ ], [ ], self.headp + 16
            if self.predetect:
                while size >= pos + 8:
                    temp = f.read(8)
                    if temp == b"\xff\xff\xff\xff\xff\xff\xff\xff":
                        break
                    else:
                        temp = decode(temp)
                        self.chunkpos.append(pos)
                        self.chunksize.append(temp)
                        pos = pos + 8 + temp
                        f.seek(pos)

    # init & write mainhead by file path
    def writef(self):
        self.headp = len(self.prehead)
        self.rsize = 0
        self.chunkpos, self.chunksize = [ ], [ ]
        if (self.headp % 512 != 0) or (len(self.common) != 4) or (len(self.subtype) != 4) or (len(self.reserved) != 8):
            raise Exception("InvalidKSC5")
        if self.path == "":
            raise Exception("NoSuchFile")
        
        with open(self.path, "wb") as f:
            f.write(self.prehead + self.common + self.subtype + self.reserved)

    # write one chunk by content of file, empty path : end (8x FF)
    def addf(self, path):
        if path == "":
            with open(self.path, "ab") as f:
                f.write(b"\xff\xff\xff\xff\xff\xff\xff\xff")
        else:
            size = os.path.getsize(path)
            with open(self.path, "ab") as f:
                f.write( encode(size, 8) )
                with open(path, "rb") as t:
                    for i in range(0, size // 10485760):
                        f.write( t.read(10485760) )
                    f.write( t.read(size % 10485760) )

    # write one chunk by binary data
    def linkf(self, data):
        with open(self.path, "ab") as f:
            f.write( encode(len(data), 8) )
            f.write(data)

    # read by binary data
    def readb(self, data):
        self.prehead, self.headp, size = b"", 0, len(data)

        while size >= self.headp + 4:
            if data[self.headp:self.headp + 4] == self.common:
                break
            else:
                self.prehead = self.prehead + data[self.headp:self.headp + 512]
                self.headp = self.headp + 512
        if size < self.headp + 16:
            raise Exception("InvalidKSC5")
            
        self.subtype = data[self.headp + 4:self.headp + 8]
        self.reserved = data[self.headp + 8:self.headp + 16]
        self.rsize = size - self.headp

        self.chunkpos, self.chunksize, pos = [ ], [ ], self.headp + 16
        if self.predetect:
            while size >= pos + 8:
                temp = data[pos:pos + 8]
                if temp == b"\xff\xff\xff\xff\xff\xff\xff\xff":
                    break
                else:
                    temp = decode(temp)
                    self.chunkpos.append(pos)
                    self.chunksize.append(temp)
                    pos = pos + 8 + temp

    # init & write mainhead by binary data
    def writeb(self):
        self.headp = len(self.prehead)
        self.rsize = 0
        self.chunkpos, self.chunksize = [ ], [ ]
        if (self.headp % 512 != 0) or (len(self.common) != 4) or (len(self.subtype) != 4) or (len(self.reserved) != 8):
            raise Exception("InvalidKSC5")
        
        return b"".join( [self.prehead, self.common, self.subtype, self.reserved] )

    # write one chunk by content of file, empty path : end (8x FF)
    def addb(self, stream, path):
        if path == "":
            return stream + b"\xff\xff\xff\xff\xff\xff\xff\xff"
        else:
            size = os.path.getsize(path)
            with open(path, "rb") as f:
                temp = f.read()
            return b"".join( [stream, encode(size, 8), temp] )

    # write one chunk by binary data
    def linkb(self, stream, data):
        return b"".join( [stream, encode(len(data), 8), data] )

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

# crc32
def crc32hash(data):
    return encode(zlib.crc32(data), 4)

# basic KSC5 prehead webp data
def webpbase():
    data0 = b""
    data0 = data0 + b"\x52\x49\x46\x46\xbc\x01\x00\x00\x57\x45\x42\x50\x56\x50\x38\x20\xb0\x01\x00\x00\x50\x09\x00\x9d\x01\x2a\x40\x00\x40\x00\x3e\x69\x2c\x90\x45\xa4\x22\xa1\x9a\xfa"
    data0 = data0 + b"\x34\xcc\x40\x06\x84\xb3\x80\x67\x2c\xd1\xff\xfa\x7a\x71\x96\xf7\x28\x36\x82\x88\x0e\x47\xbf\x98\xba\x75\x08\x2e\xed\xcf\x6a\xfd\x0d\xf3\xed\x21\xc2\x2e\xfc\xc0"
    data0 = data0 + b"\xc4\x87\x23\xf1\xf7\x5b\x9a\x84\xe6\x14\x3f\x5b\xc1\x1b\x75\x50\xb9\xae\xe4\x94\xb6\x49\xc0\x00\xec\x8f\xc2\xb9\x12\xdc\x9c\x2a\x54\x4d\x1f\xdb\xe8\xf0\xdd\x00"
    data0 = data0 + b"\xc3\xca\xff\x48\xc4\x4b\xf9\x50\x76\x84\x72\x68\xdf\xd4\xed\xc9\x9c\xe0\x78\xf5\xf4\x61\x6e\x63\x52\xd8\x15\x7c\xe5\x23\xd5\x3b\x97\x09\x67\x89\x6a\x8b\xe3\x18"
    data0 = data0 + b"\xb0\x8a\xa4\x41\x57\x8a\x97\x6b\xc8\x06\x2c\xd0\xd6\xce\xe6\xbb\x42\x6f\x8c\x01\x5f\xa5\x6f\x34\xcf\xbe\x81\x92\x33\x3f\x9c\xaa\x25\x16\x29\x53\x27\xa1\x48\xd8"
    data0 = data0 + b"\x8e\x58\xc3\xfd\x01\x34\x52\x36\x6c\x25\xfc\xc0\x68\xed\x2c\x0d\x95\x9a\x8a\xd3\x89\x90\xf7\x2a\xa2\x3b\xb0\xc6\xd4\xf5\x19\x51\x34\xbc\x72\x2f\x2d\xb8\x46\x32"
    data0 = data0 + b"\x8b\x59\x33\x08\xb1\x76\x37\x77\x63\x41\xa5\x0e\xee\x55\x7d\x73\x2e\x5b\x8e\x50\x99\x41\xfe\x96\xf2\x34\x36\x3a\xda\xd0\xd9\x17\x89\x13\xa0\x9b\xab\xf9\x48\xbd"
    data0 = data0 + b"\xca\xf8\x36\x33\x0e\x98\xd9\xe8\x95\x94\x2a\xb0\x92\x42\x01\x47\xe0\x53\xa7\xfd\x65\x38\x6a\x65\x6b\x39\x3c\x6c\xca\x2d\x51\x32\x33\xcd\x68\x5c\xdf\x23\xae\x58"
    data0 = data0 + b"\x4f\x0e\x7d\xbb\xdb\x98\x56\x62\xc4\x61\x94\x2a\xdd\x51\xcd\x79\xeb\x38\x99\x4a\x4a\x6b\xa6\x2d\x7b\xc5\xba\xbb\x92\x9a\xf4\x7d\xb7\xdc\x57\xc6\xfd\xe9\x66\x2a"
    data0 = data0 + b"\x8f\xf8\x04\xf3\x28\x33\x0c\xb5\x12\x8b\x39\xd0\xaa\xce\xae\xac\x90\x77\x2d\xe3\xf2\xb0\x23\x21\xa7\x4f\x26\xfb\x51\x45\x9e\x27\xc3\xb0\x6a\x6a\x0b\xd3\xfd\x3c"
    data0 = data0 + b"\x31\x65\x05\xfa\x21\x64\xa0\x36\x1e\x5d\x8b\xd4\xe9\xd7\xe1\x59\x72\x2f\x0c\x1c\x47\x3e\x1c\xbf\xe0\x15\xc0\x57\xa2\x30\x47\xc6\x4d\x02\xbd\xbf\x6f\xb6\xfc\x12"
    data0 = data0 + b"\x10\x4c\x3b\xde\xea\xdd\x7b\x00\x00\x00\x00\x00"

    data0 = data0 + b"\x00" * ( 512 - len(data0) )
    return data0