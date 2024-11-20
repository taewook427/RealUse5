# test705 : stdlib5.kscript stdlib/stdio/osfs

import sys
import os
import shutil
import time
import math
import random

import kobj

def tostr(a): # get string of var
    return a.hex() if type(a) == bytes else str(a)

def getsize(path): # get size of file/folder
    if os.path.isfile(path):
        return os.path.getsize(path)
    else:
        if path[-1] != "/":
            path = path + "/"
        temp = [path + x for x in os.listdir(path)]
        return sum( [getsize(x) for x in temp] )

class lib_stdio: # basic stdio object
    def stdin(self): # read 1 line
        return input("")
    def stdout(self, word): # write word
        print(word, end="")
    def stderr(self, word): # write word
        sys.stderr.write(word)

# 3-module worker, (stdlib, stdio, osfs), !!! manual set path & io !!!
class lib:
    def __init__(self, u_stdlib, u_stdio, u_osfs, iswin):
        self.u_stdlib = u_stdlib # lib using flag (abi 2, iter 32~100)
        self.u_stdio = u_stdio # lib using flag (abi 4, iter 100~200)
        self.u_osfs = u_osfs # lib using flag (abi 8, iter 200~300)
        self.myos = "windows" if iswin else "linux" # working os name
        self.hmem = dict() # global memory map (mem[id])
        self.fhnd = dict() # file handle (file[fullpath])

        self.p_desktop = "!undefined_desktop" # desktop path (~/)
        self.p_local = "!undefined_local" # local user path (~/)
        self.p_starter = "!undefined_starter" # st5 path (~/)
        self.p_base = "!undefined_base" # program file position (~)
        self.io = lib_stdio() # stdio class (stdin, stdout, stderr)

    def stdlib(self, pa, cd): # stdlib func (32~100)
        if not self.u_stdlib:
            return None, "not supported func"
        vo, er = None, None
        if cd < 50: # builtin
            if cd == 32: # type
                if pa[0] == None:
                    vo = "none"
                elif type( pa[0] ) == bool:
                    vo = "bool"
                elif type( pa[0] ) == int:
                    vo = "int"
                elif type( pa[0] ) == float:
                    vo = "float"
                elif type( pa[0] ) == str:
                    vo = "str"
                elif type( pa[0] ) == bytes:
                    vo = "bytes"

            elif cd == 33: # int
                if type( pa[0] ) == int or type( pa[0] ) == float or type( pa[0] ) == str:
                    vo = int( pa[0] )
                elif type( pa[0] ) == bytes:
                    vo = kobj.decode( pa[0] )
                else:
                    er = "type error"

            elif cd == 34: # float
                if type( pa[0] ) == int or type( pa[0] ) == float or type( pa[0] ) == str:
                    vo = float( pa[0] )
                else:
                    er = "type error"

            elif cd == 35: # str
                if type( pa[0] ) == bytes:
                    vo = str(pa[0], encoding="utf-8")
                else:
                    vo = str( pa[0] )

            elif cd == 36: # bytes
                if type( pa[0] ) == int:
                    vo = kobj.encode(pa[0], 8)
                elif type( pa[0] ) == str:
                    vo = bytes(pa[0], encoding="utf-8")
                elif type( pa[0] ) == bytes:
                    vo = pa[0]
                else:
                    er = "type error"

            elif cd == 37: # hex
                if type( pa[0] ) == int:
                    vo = hex( pa[0] )[2:]
                elif type( pa[0] ) == bytes:
                    vo = pa[0].hex()
                else:
                    er = "type error"

            elif cd == 38: # chr
                if self.check( pa, ["i"] ):
                    vo = chr( pa[0] )
                else:
                    er = "type error"

            elif cd == 39: # ord
                if self.check( pa, ["s"] ):
                    vo = ord( pa[0] )
                else:
                    er = "type error"

            elif cd == 40: # len
                if self.check( pa, ["s|c"] ):
                    vo = len( pa[0] )
                else:
                    er = "type error"

            elif cd == 41: # slice
                if self.check( pa, ["s|c", "n|i", "n|i"] ):
                    st = pa[1] if pa[1] != None else 0
                    ed = pa[2] if pa[2] != None else len( pa[0] )
                    vo = pa[0][st:ed]
                else:
                    er = "type error"

        elif cd < 60: # memory
            if cd == 50: # m.malloc
                if not self.check( pa, ["s", "i"] ):
                    er = "type error"
                elif pa[1] < 0:
                    er = "invalid size"
                else:
                    if pa[0] in self.hmem:
                        er = "double alloc"
                    self.hmem[ pa[0] ] = [None] * pa[1]

            elif cd == 51: # m.realloc
                if not self.check( pa, ["s", "i"] ):
                    er = "type error"
                elif pa[1] < 0:
                    er = "invalid size"
                else:
                    tgt = self.hmem[ pa[0] ] if pa[0] in self.hmem else [ ]
                    olds, news = len(tgt), pa[1]
                    self.hmem[ pa[0] ] = tgt + [None] * (news - olds) if olds < news else tgt[0:news]

            elif cd == 52: # m.free
                if pa[0] in self.hmem:
                    del self.hmem[ pa[0] ]
                else:
                    er = "double free"

            elif cd == 53: # m.len
                vo = len( self.hmem[ pa[0] ] ) if pa[0] in self.hmem else -1

            elif cd == 54: # m.set
                if not self.check( pa, ["s", "i", "a"] ):
                    er = "type error"
                elif pa[0] not in self.hmem:
                    er = "invalid id"
                else:
                    self.hmem[ pa[0] ][ pa[1] ] = pa[2]

            elif cd == 55: # m.get
                if not self.check( pa, ["s", "i"] ):
                    er = "type error"
                elif pa[0] not in self.hmem:
                    er = "invalid id"
                else:
                    vo = self.hmem[ pa[0] ][ pa[1] ]

            elif cd == 56: # m.split
                if not self.check( pa, ["s", "s|c", "s|c"] ):
                    er = "type error"
                else:
                    self.hmem[ pa[0] ] = pa[1].split( pa[2] )

            elif cd == 57: # m.join
                if not self.check( pa, ["s", "s|c"] ):
                    er = "type error"
                elif pa[0] not in self.hmem:
                    er = "invalid id"
                else:
                    temp = [ ]
                    for i in self.hmem[ pa[0] ]:
                        if type(i) == type( pa[1] ):
                            temp.append(i)
                    vo = pa[1].join(temp)

        elif cd < 70: # string
            if cd == 60: # str.change
                if not self.check( pa, ["s", "b"] ):
                    er = "type error"
                else:
                    vo = pa[0].upper() if pa[1] else pa[0].lower()

            elif cd == 61: # str.find
                if not self.check( pa, ["s", "s", "b"] ):
                    er = "type error"
                else:
                    vo = pa[0].find( pa[1] ) if pa[2] else pa[0].rfind( pa[1] )

            elif cd == 62: # str.count
                if not self.check( pa, ["s", "s"] ):
                    er = "type error"
                else:
                    vo = pa[0].count( pa[1] )

            elif cd == 63: # str.replace
                if not self.check( pa, ["s", "s", "s", "i"] ):
                    er = "type error"
                else:
                    vo = pa[0].replace( pa[1], pa[2] ) if pa[3] < 0 else pa[0].replace( pa[1], pa[2], pa[3] )

        elif cd < 80: # time
            if cd == 70: # t.time
                vo = time.time()

            elif cd == 71: # t.stamp
                if not self.check( pa, ["i|f"] ):
                    er = "type error"
                else:
                    vo = time.strftime( "%Y.%m.%d;%H:%M:%S", time.localtime( pa[0] ) )

            elif cd == 72: # t.stampf
                if not self.check( pa, ["i|f", "s"] ):
                    er = "type error"
                else:
                    fmt = pa[1].replace("%Y", "%%Y").replace("%M", "%%M").replace("%D", "%%D").replace("%h", "%%h").replace("%m", "%%m").replace("%s", "%%s")
                    fmt = fmt.replace("%%Y", "%Y").replace("%%M", "%m").replace("%%D", "%d").replace("%%h", "%H").replace("%%m", "%M").replace("%%s", "%S")
                    vo = time.strftime( fmt, time.localtime( pa[0] ) )

            elif cd == 73: # t.sleep
                if not self.check( pa, ["i|f"] ):
                    er = "type error"
                else:
                    time.sleep( pa[0] )

        elif cd < 90: # math
            if cd == 80: # math.const
                if not self.check( pa, ["s"] ):
                    er = "type error"
                elif pa[0] == "e":
                    vo = math.e
                elif pa[0] == "pi":
                    vo = math.pi
                elif pa[0] == "phi":
                    vo = ( 1 + math.sqrt(5) ) / 2
                elif pa[0] == "sqrt2":
                    vo = math.sqrt(2)
                else:
                    er = "invalid option"

            elif cd == 81: # math.conv
                if not self.check( pa, ["s", "i|f"] ):
                    er = "type error"
                elif pa[0] == "abs":
                    vo = abs( pa[1] )
                elif pa[0] == "up":
                    vo = math.ceil( pa[1] )
                elif pa[0] == "down":
                    vo = math.floor( pa[1] )
                elif pa[0] == "round":
                    vo = round( pa[1] )
                else:
                    er = "invalid option"

            elif cd == 82: # math.log
                if not self.check( pa, ["i|f", "i|f"] ):
                    er = "type error"
                else:
                    vo = math.log( pa[1], pa[0] )

            elif cd == 83: # math.trif
                if not self.check( pa, ["s", "i|f"] ):
                    er = "type error"
                elif pa[0] == "sin":
                    vo = math.sin( pa[1] )
                elif pa[0] == "cos":
                    vo = math.cos( pa[1] )
                elif pa[0] == "tan":
                    vo = math.tan( pa[1] )
                elif pa[0] == "asin":
                    vo = math.asin( pa[1] )
                elif pa[0] == "acos":
                    vo = math.acos( pa[1] )
                elif pa[0] == "atan":
                    vo = math.atan( pa[1] )
                else:
                    er = "invalid option"

            elif cd == 84: # math.random
                vo = random.random()

            elif cd == 85: # math.randrange
                if not self.check( pa, ["i", "i"] ):
                    er = "type error"
                elif pa[0] >= pa[1]:
                    er = "invalid option"
                else:
                    vo = random.randrange( pa[0], pa[1] )

        return vo, er

    def stdio(self, pa, cd): # stdio func (100~200)
        if not self.u_stdio:
            return None, "not supported func"
        vo, er = None, None
        if cd == 100: # io.input
            self.io.stdout( tostr( pa[0] ) )
            vo = self.io.stdin()

        elif cd == 101: # io.print
            if not self.check( pa, ["a", "s"] ):
                er = "type error"
            else:
                self.io.stdout( tostr( pa[0] ) + pa[1] )

        elif cd == 102: # io.println
            self.io.stdout(tostr( pa[0] ) + "\n")

        elif cd == 103: # io.error
            if not self.check( pa, ["a", "s"] ):
                er = "type error"
            else:
                self.io.stderr( tostr( pa[0] ) + pa[1] )

        elif cd == 104: # io.open
            if not self.check( pa, ["s", "s"] ):
                er = "type error"
            else:
                path = os.path.abspath( pa[0] ).replace("\\", "/")
                if path in self.fhnd:
                    self.fhnd[path].close()
                self.fhnd[path] = open(path, pa[1]+"b")

        elif cd == 105: # io.close
            if not self.check( pa, ["s"] ):
                er = "type error"
            else:
                path = os.path.abspath( pa[0] ).replace("\\", "/")
                if path in self.fhnd:
                    self.fhnd[path].close()
                    del self.fhnd[path]

        elif cd == 106: # io.seek
            if not self.check( pa, ["s", "i", "i"] ):
                er = "type error"
            else:
                path = os.path.abspath( pa[0] ).replace("\\", "/")
                if path not in self.fhnd:
                    er = "handle not exists"
                else:
                    if pa[2] == 0:
                        self.fhnd[path].seek( pa[1] )
                    elif pa[2] == -1:
                        self.fhnd[path].seek(pa[1], 2)
                    else:
                        self.fhnd[path].seek(pa[1], 1)
                    vo = self.fhnd[path].tell()

        elif cd == 107: # io.readline
            if not self.check( pa, ["s", "i"] ):
                er = "type error"
            else:
                path = os.path.abspath( pa[0] ).replace("\\", "/")
                if path not in self.fhnd:
                    er = "handle not exists"
                else:
                    data, count = [ ], pa[1]
                    while count != 0:
                        temp = self.fhnd[path].readline()
                        if temp == b"":
                            break
                        data.append(temp)
                        count = count - 1
                    vo = str(b"".join(data), encoding="utf-8")

        elif cd == 108: # io.read
            if not self.check( pa, ["s", "i"] ):
                er = "type error"
            else:
                path = os.path.abspath( pa[0] ).replace("\\", "/")
                if path not in self.fhnd:
                    er = "handle not exists"
                else:
                    vo = self.fhnd[path].read() if pa[1] < 0 else self.fhnd[path].read( pa[1] )

        elif cd == 109: # io.write
            if not self.check( pa, ["s", "s|c"] ):
                er = "type error"
            else:
                path = os.path.abspath( pa[0] ).replace("\\", "/")
                if path not in self.fhnd:
                    er = "handle not exists"
                else:
                    self.fhnd[path].write( pa[1] ) if type( pa[1] ) == bytes else self.fhnd[path].write( bytes(pa[1], encoding="utf-8") )

        else:
            er = "not supported func"
        return vo, er

    def osfs(self, pa, cd): # osfs func (200~300)
        if not self.u_osfs:
            return None, "not supported func"
        vo, er = None, None
        if cd == 200: # os.name
            vo = self.myos

        elif cd == 201: # os.chdir
            if not self.check( pa, ["s"] ):
                er = "type error"
            else:
                os.chdir( pa[0] )

        elif cd == 202: # os.getpath
            if pa[0] == "cwd":
                vo = os.getcwd().replace("\\", "/")
                if vo[-1] != "/":
                    vo = vo + "/"
            elif pa[0] == "desktop":
                vo = self.p_desktop
            elif pa[0] == "local":
                vo = self.p_local
            elif pa[0] == "starter":
                vo = self.p_starter
            elif pa[0] == "base":
                vo = self.p_base
            else:
                er = "invalid option"

        elif cd == 203: # os.exists
            vo = os.path.exists( pa[0] )

        elif cd == 204: # os.abspath
            vo = os.path.abspath( pa[0] ).replace("\\", "/")
            if os.path.isdir(vo) and vo[-1] != "/":
                vo = vo + "/"

        elif cd == 205: # os.is
            if not self.check( pa, ["s", "b"] ):
                er = "type error"
            else:
                vo = os.path.isfile( pa[0] ) if pa[1] else os.path.isdir( pa[0] )

        elif cd == 206: # os.finfo
            if not self.check( pa, ["s", "b"] ):
                er = "type error"
            elif not os.path.exists( pa[0] ):
                vo = -1
            else:
                vo = getsize( pa[0] ) if pa[1] else os.path.getmtime( pa[0] )

        elif cd == 207: # os.listdir
            if pa[0][-1] != "/":
                pa[0] = pa[0] + "/"
            temp = os.listdir( pa[0] )
            for i in range( 0, len(temp) ):
                if os.path.isdir( pa[0] + temp[i] ) and temp[i][-1] != "/":
                    temp[i] = temp[i] + "/"
            vo = "\n".join(temp)

        elif cd == 208: # os.mkdir
            os.mkdir( pa[0] )

        elif cd == 209: # os.rename
            os.rename( pa[0], pa[1] )

        elif cd == 210: # os.move
            os.rename( pa[0], pa[1] )

        elif cd == 211: # os.remove
            if os.path.isfile( pa[0] ):
                os.remove( pa[0] )
            else:
                shutil.rmtree( pa[0] )

        else:
            er = "not supported func"
        return vo, er

    def check(self, pa, tp): # check parms types (nbifsc | a)
        if len(pa) != len(tp):
            return False
        pos = 0
        for i in pa:
            flag = False
            for j in tp[pos].split("|"):
                if j == "a":
                    flag = True
                elif j == "n" and i == None:
                    flag = True
                elif j == "b" and type(i) == bool:
                    flag = True
                elif j == "i" and type(i) == int:
                    flag = True
                elif j == "f" and type(i) == float:
                    flag = True
                elif j == "s" and type(i) == str:
                    flag = True
                elif j == "c" and type(i) == bytes:
                    flag = True
                if flag:
                    break
            if not flag:
                return False
            pos = pos + 1
        return True

    # close all file handles
    def exit(self):
        for i in self.fhnd:
            try:
                self.fhnd[i].close()
            except:
                pass
        self.__init__(False, False, False, False)

    # outer func works
    def run(self, parms, icode):
        vout, err = None, None
        try:
            if icode < 32:
                err = "not supported func"
            elif icode < 100: # stdlib
                vout, err = self.stdlib(parms, icode)
            elif icode < 200: # stdio
                vout, err = self.stdio(parms, icode)
            elif icode < 300: # osfs
                vout, err = self.osfs(parms, icode)
            else:
                err = "not supported func"
        except Exception as e:
            vout, err = None, f"critical : {e}"
        if err == "":
            err = None
        return vout, err
