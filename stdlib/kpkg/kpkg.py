# test674 : stdlib5.kpkg

import os
import shutil
import time
import hashlib

import kdb
import ksc
import ksignhy
import kzip

class toolbox: # use ./temp674/
    def __init__(self):
        self.public = "" # RSA public key, empty if not using
        self.private = "" # RSA private key, empty if not using
        self.osnum = 0 # any : 0,  windows : 1, linux mint : 2
        self.name = "" # package name
        self.version = 0.0 # package version
        self.text = "" # package explanation
        self.rel_date = "" # package release date
        self.dwn_date = "" # package download date

    # pack package by osnum/dirpaths, gen file at respath
    def pack(self, osnums, dirpaths, respath):
        # setting required : public, private, name, version, text
        if len(osnums) != len(dirpaths):
            raise Exception("invalid package order")
        for i in dirpaths:
            if not os.path.isdir(i):
                raise Exception("invalid package order")
        temp = time.localtime( time.time() )
        self.rel_date = time.strftime('%Y%m%d' , temp)

        # chunk0
        worker = kdb.toolbox()
        temp = ""
        for i in range( 0, len(osnums) ):
            temp = temp + f"pkg{i} = {osnums[i]}\n"
        worker.read(temp + "name = 0\nversion = 0\ntext = 0\nrelease = 0")
        worker.fix("name", self.name)
        worker.fix("version", self.version)
        worker.fix("text", self.text)
        worker.fix("release", self.rel_date)
        chunk0 = bytes(worker.write(), encoding="utf-8") # chunk0 : basic info

        # init temp dir, pack with kzip
        dirinit(True)
        for i in range( 0, len(osnums) ):
            kzip.dozip([ dirpaths[i] ], "webp", f"./temp674/{i}.webp")

        if self.private == "":
            # not using sign
            chunk1 = b"" # chunk1 : sign info
        else:
            # using sign
            worker = kdb.toolbox()
            temp = ""
            for i in range( 0, len(osnums) ):
                temp = temp + f"pkg{i} = 0\n"
            temp = temp + "public = 0"
            worker.read(temp)
            signtool = ksignhy.toolbox()
            for i in range( 0, len(osnums) ):
                with open(f"./temp674/{i}.webp", "rb") as f:
                    hvalue = hashlib.sha3_512( chunk0 + f.read() ).digest()
                worker.fix( f"pkg{i}", signtool.sign(self.private, hvalue) )
            worker.fix("public", self.public)
            chunk1 = bytes(worker.write(), encoding="utf-8") # chunk1 : sign info

        worker = ksc.toolbox() # make basic structure
        worker.path = respath
        worker.prehead = genwebp()
        worker.subtype = b"KPKG"
        worker.reserved = ksc.crc32hash(chunk0) + ksc.crc32hash(chunk1)
        worker.writef()
        worker.linkf(chunk0)
        worker.linkf(chunk1)

        # add kzip files
        for i in range( 0, len(osnums) ):
            worker.addf(f"./temp674/{i}.webp")
        dirinit(True)

    # unpack package by os, update internal data, gen folder at temp674/temp/, returns dir path
    def unpack(self, path):
        # setting required : osnum
        dirinit(True)
        worker = ksc.toolbox()
        worker.predetect = True
        worker.path = path
        worker.readf()
        if worker.subtype != b"KPKG":
            raise Exception("invalid package")
        with open(path, "rb") as f:
            f.seek(worker.chunkpos[0] + 8)
            chunk0 = f.read( worker.chunksize[0] ) # basic data
            f.seek(worker.chunkpos[1] + 8)
            chunk1 = f.read( worker.chunksize[1] ) # sign data
        if worker.reserved != ksc.crc32hash(chunk0) + ksc.crc32hash(chunk1):
            raise Exception("invalid package")

        c0data = kdb.toolbox() # get info chunk
        c0data.read( str(chunk0, encoding="utf-8") )
        c1valid = False if len(chunk1) == 0 else True
        if c1valid:
            c1data = kdb.toolbox()
            c1data.read( str(chunk1, encoding="utf-8") )

        temp = time.localtime( time.time() ) # init install info
        self.name = c0data.get("name")[3]
        self.version = c0data.get("version")[3]
        self.text= c0data.get("text")[3]
        self.rel_date = c0data.get("release")[3]
        self.dwn_date = time.strftime('%Y%m%d' , temp)

        osnums = [ ] # get os numbers
        temp = 0
        while f"pkg{temp}" in c0data.name:
            osnums.append( c0data.get(f"pkg{temp}")[3] )
            temp = temp + 1
        if (0 not in osnums) and (self.osnum not in osnums):
            raise Exception("OS not support PKG")
        if 0 in osnums:
            pos = osnums.index(0)
        else:
            pos = osnums.index(self.osnum)

        # extract chunk corresponds osnum
        with open(path, "rb") as f:
            with open("./temp674/temp.webp", "wb") as t:
                f.seek(worker.chunkpos[pos + 2] + 8)
                size = worker.chunksize[pos + 2]
                for i in range(0, size // 10485760):
                    t.write( f.read(10485760) )
                t.write( f.read(size % 10485760) )

        # check sign
        self.public = ""
        if c1valid:
            with open("./temp674/temp.webp", "rb") as f:
                hvalue = hashlib.sha3_512( chunk0 + f.read() ).digest()
            enc = c1data.get(f"pkg{pos}")[3]
            self.public = c1data.get("public")[3]
            worker = ksignhy.toolbox()
            if not worker.verify(self.public, enc, hvalue):
                raise Exception("invalid sign")

        os.mkdir("./temp674/temp") # unpack package
        kzip.unzip("./temp674/temp.webp", "./temp674/temp", True)
        pkgpath = ( "./temp674/temp/" + os.listdir("./temp674/temp")[0] )
        pkgpath = os.path.abspath(pkgpath).replace("\\", "/")
        if pkgpath[-1] != "/":
            pkgpath = pkgpath + "/"
        with open(f"{pkgpath}_ST5_VERSION.txt", "w", encoding="utf-8") as f:
            worker = kdb.toolbox()
            worker.read("name = 0\nversion = 0\ntext = 0\nrelease = 0\ndownload = 0")
            worker.fix("name", self.name)
            worker.fix("version", self.version)
            worker.fix("text", self.text)
            worker.fix("release", self.rel_date)
            worker.fix("download", self.dwn_date)
            f.write( worker.write() )
        return pkgpath

def dirinit(mode): # init dir True : 674, False : 675
    name = "./temp674" if mode else "./temp675"
    if os.path.exists(name):
        shutil.rmtree(name)
    os.mkdir(name)

def genwebp(): # make kpkg base picture
    data0 = b""
    data0 = data0 + b"\x52\x49\x46\x46\x96\x02\x00\x00\x57\x45\x42\x50\x56\x50\x38\x20\x8a\x02\x00\x00\x50\x0f\x00\x9d\x01\x2a\x40\x00\x40\x00\x3e\x6d\x2e\x93\x46\xa4\x22\xa1\xa1\x24"
    data0 = data0 + b"\x0e\xd8\x80\x0d\x89\x6a\x00\xc0\xd4\x64\x41\x47\xb2\x7d\x80\x7e\x99\xf4\xd1\xef\xc4\x54\xde\xf6\xda\xbe\xaa\xf6\xd4\x79\x80\xfd\x6e\xfd\x80\xf7\x7c\xfe\xe5\xea"
    data0 = data0 + b"\x03\xa1\x9b\xd6\x03\xd0\x03\xf6\x33\xd3\x2f\xd9\x5b\xfc\xc5\x7d\x1f\x7a\x48\x60\xcc\x35\x8f\x03\xd0\x87\x3b\x3f\x48\x11\x21\xc2\x91\x35\xa2\x51\xd4\xe0\x39\x06"
    data0 = data0 + b"\xb4\x5d\xad\x6a\x19\x4b\x7c\x03\x68\xbe\xf8\x64\x31\xe1\x61\x3a\xff\xf9\x85\x0e\x20\x00\x5a\x8b\x01\x8b\xf2\x81\x25\x30\x60\x00\x44\x46\x4c\x57\x4c\x06\x33\x55"
    data0 = data0 + b"\xe2\x33\x87\xef\x4a\xe2\xfe\x0c\xfd\xcd\x04\xb2\x22\xef\xc0\xb5\x7f\xf0\xa7\x0b\x36\xcb\xb0\x88\x8b\xe3\x09\x14\x21\x1f\x76\xe0\x5c\xd9\x8f\xd1\xff\xfd\x6c\xfb"
    data0 = data0 + b"\xdd\xd9\x92\xaf\xad\xd1\x9a\x53\x77\x0f\xdf\xff\xbf\xb3\x79\xcb\x30\x2a\xbe\xbc\xdf\x4e\x35\x08\xa9\xe2\x80\xcd\xcc\x8c\xd6\x0f\x3a\x91\x1b\x0a\xef\x85\x55\x07"
    data0 = data0 + b"\x08\x4d\x6c\x67\x58\x7e\xa9\xc9\xa0\xb9\x61\xf7\x11\x68\xdf\x5c\x80\xa6\x19\xfe\x86\x35\x9a\xde\x4f\x75\xee\x6d\xa1\xcd\x9f\x23\x79\xb6\x2d\x10\x05\x1d\x14\x98"
    data0 = data0 + b"\xc7\x06\xfc\x33\x59\x24\x79\x13\x25\x60\x04\x4c\x18\xbb\x22\xc5\xa4\xaf\x86\xe2\x88\x03\x5a\xe9\xa4\xc0\xd6\x7d\x86\x45\xa3\xf6\x39\x34\x67\x9e\xce\x7e\x72\x28"
    data0 = data0 + b"\x51\xc0\x7d\xa0\x78\x17\xcb\x29\xc5\x5d\xf4\x5b\x06\xaf\x5a\xae\xe4\x4a\x98\xdd\xe5\x91\xca\xd2\x6b\x7b\x14\x2d\xe5\xa2\xaf\x5d\xa3\x90\xff\x69\xef\x07\x3e\x0c"
    data0 = data0 + b"\x95\xd3\x23\x6c\xef\x45\xf0\x30\x34\x48\xcc\xaf\x29\x19\x12\x3f\x6d\x61\x89\x22\xcc\xb4\x22\x34\x2b\x5f\xdb\xc1\xf6\xe6\x1d\x17\x73\xfd\x38\x84\xe9\xfc\x47\xd0"
    data0 = data0 + b"\x8c\x02\xfd\x30\xca\x65\xc8\xb6\xfa\x21\xa2\xf3\x0e\x1f\x69\x3c\x40\x7f\xcd\x9e\x49\x61\xfa\xbd\x1d\x26\xfb\x62\xe8\x87\x3f\x2b\x8c\x96\x72\x2a\x9f\xc3\xd5\x2e"
    data0 = data0 + b"\xea\x40\xb1\x9f\x32\xb4\x8a\x20\xfe\xea\xf5\xf5\x53\xd2\xdc\xe1\xe9\x98\xe8\x4f\xd2\xcc\x2a\x9f\xc2\x55\x44\xc1\x60\x06\xf1\x59\x0c\x2e\xba\xbe\x7c\x5f\x38\xa3"
    data0 = data0 + b"\x79\x76\xac\xfc\xd9\x19\x0f\xc7\xd8\xff\xfa\xb5\x23\x5c\x2d\xe5\xae\x97\xbc\x1a\xe3\xd2\xcd\x01\x79\xec\xbc\x09\x8e\xc6\x6d\x8d\xc8\x2b\x94\xbc\x80\x13\xb8\x7f"
    data0 = data0 + b"\x1b\xd2\xde\x01\x27\xcb\xc2\xfd\xdf\xc5\xce\x12\x8c\x82\x95\x33\x37\x56\x8a\x12\xf0\xcd\x98\x6c\x18\xe2\xd0\x9d\x0c\x67\xf9\x83\x49\xb9\xb3\x20\xca\xe8\x03\x56"
    data0 = data0 + b"\x6c\x21\x64\xaf\x25\xd2\xc9\x38\x00\x65\x0a\xc1\x90\x7e\x8f\x6b\x26\xef\x65\xfc\x7f\x30\xfb\xca\xde\xc5\xb1\x06\x23\x74\xb8\xb6\x80\xed\x73\x0d\x55\x16\x79\x56"
    data0 = data0 + b"\x93\x9f\x64\x2b\x98\xa4\xa0\xca\x3e\x16\xcd\x7e\x99\xf2\xd3\x31\x5e\x25\x4a\xff\x56\x27\x23\x38\x89\x92\xd8\x62\x04\x34\x37\x19\x15\x7b\x75\x13\x26\xfd\x65\xfe"
    data0 = data0 + b"\x39\xec\xac\x77\x11\x5a\x24\xbb\x69\x3d\x7f\x24\x7f\xa8\xee\x90\xc7\x22\xa2\x38\xfc\x77\xd9\x4e\xcf\xe1\x83\x40\x00\x00"
    return data0 + ( 1024 - len(data0) ) * b"\x00"
