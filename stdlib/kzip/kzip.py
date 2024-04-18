# test642 : stdlib5.kzip

import os
import shutil

import picdt
import ksc

def delpath(path): # delete file/folder if exists
    if os.path.exists(path):
        if os.path.isfile(path):
            os.remove(path)
        else:
            shutil.rmtree(path)

def abspath(tgts): # return stdpath of folder/file list
    temp = [os.path.abspath(x).replace("\\", "/") for x in tgts]
    for i in range( 0, len(temp) ):
        if os.path.isdir( temp[i] ) and temp[i][-1] != "/":
            temp[i] = temp[i] + "/"
    return temp

def getlist(tgts): # get folder/file names from tgts
    lstd, lstf, absf = [ ], [ ], [ ]
    tgts = abspath(tgts)
    count = 0
    for i in tgts:
        if i[-1] == "/":
            if i.count("/") == 1:
                temp = f"LargeVolume{count}/"
                count = count + 1
            else:
                temp = i[0:-1]
                temp = temp[temp.rfind("/") + 1:] + "/"
            getsub(i, temp, lstd, lstf, absf)
        else:
            lstf.append( i[i.rfind("/") + 1:] )
            absf.append(i)
    return lstd, lstf, absf

def getsub(path, prefix, lstd, lstf, absf): # path : stdpath(~/) of root, prefix : ("root/")
    lstd.append(prefix)
    temp = [ x for x in os.listdir(path) ]
    for i in range( 0, len(temp) ):
        if os.path.isdir( path + temp[i] ) and temp[i][-1] != "/":
            temp[i] = temp[i] + "/"
    for i in temp:
        if i[-1] == "/":
            getsub(path + i, prefix + i, lstd, lstf, absf)
        else:
            lstf.append(prefix + i)
            absf.append(path + i)

# zip. tgts : folder/file path, mode ("webp", "png", ""), path ("" -> "./temp570.webp")
def dozip(tgts, mode, path):
    lstd, lstf, absf = getlist(tgts)
    header = bytes("\n".join(lstd), encoding="utf-8")

    writer = ksc.toolbox()
    if mode == "webp":
        writer.prehead = picdt.toolbox().data1
        writer.prehead = writer.prehead + b"\x00" * ( ( 0 - len(writer.prehead) ) % 512 )
    elif mode == "png":
        writer.prehead = picdt.toolbox().data0
        writer.prehead = writer.prehead + b"\x00" * ( ( 0 - len(writer.prehead) ) % 512 )
    else:
        writer.prehead = b""
    writer.subtype = b"kzip"
    writer.reserved = ksc.crc32hash(header) + b"\x00\x00\x00\x00"
    if path == "":
        writer.path = "./temp570.webp"
    else:
        writer.path = path
    delpath(writer.path)

    writer.writef()
    writer.linkf(header)
    for i in range( 0, len(lstf) ):
        writer.linkf( bytes(lstf[i], encoding="utf-8") )
        writer.addf( absf[i] )
    writer.addf("")

# unzip. kzip file path, export folder ("" -> "./temp570/"), check subtype/crc error
def unzip(path, export, chkerr):
    if export == "":
        export = "./temp570/"
    export = export.replace("\\", "/")
    if export[-1] != "/":
        export = export + "/"
    delpath(export)
    os.mkdir(export)

    reader = ksc.toolbox()
    reader.path = path
    reader.predetect = True
    reader.readf()

    with open(path, "rb") as f:
        f.seek(reader.chunkpos[0] + 8)
        header = f.read( reader.chunksize[0] )
        if chkerr:
            if reader.subtype != b"kzip":
                raise Exception("InvalidKZIP")
            if reader.reserved[0:4] != ksc.crc32hash(header):
                raise Exception("InvalidCRC")
            if len(reader.chunkpos) % 2 != 1:
                raise Exception("InvalidChunk")

        if header == b"":
            header = [ ]
        else:
            header = str(header, encoding="utf-8").split("\n")
        for i in header:
            os.mkdir(export + i)

        for i in range(0, (len(reader.chunkpos) - 1) // 2):
            pos = 2 * i + 1
            f.seek(reader.chunkpos[pos] + 8)
            name = str(f.read( reader.chunksize[pos] ), encoding="utf-8")
            size = reader.chunksize[pos + 1]
            f.seek(reader.chunkpos[pos + 1] + 8)
            with open(export + name, "wb") as t:
                for j in range(0, size // 10485760):
                    t.write( f.read(10485760) )
                if size % 10485760 != 0:
                    t.write( f.read(size % 10485760) )
