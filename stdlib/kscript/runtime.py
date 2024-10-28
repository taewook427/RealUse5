# test692 : stdlib5.kscript runtime py

import time
import struct
import hashlib
import collections

import kdb
import ksc
import ksign

# arithmetic logic calculation : get 2 operand, returns result & error

def add(a, b):
    cond0 = { (int, int), (int, float), (float, int), (float, float), (str, str), (bytes, bytes) }
    cond1 = { (bool, bool) }
    vtype = ( type(a), type(b) )
    if vtype in cond0:
        return a + b, ""
    elif vtype in cond1:
        return a or b, ""
    else:
        return None, f"e800 : cannot add {vtype}"

def sub(a, b):
    cond0 = { (int, int), (int, float), (float, int), (float, float) }
    vtype = ( type(a), type(b) )
    if vtype in cond0:
        return a - b, ""
    else:
        return None, f"e801 : cannot sub {vtype}"

def mul(a, b):
    cond0 = { (int, int), (int, float), (float, int), (float, float) }
    cond1 = { (int, str), (int, bytes), (str, int), (bytes, int) }
    cond2 = { (bool, bool) }
    vtype = ( type(a), type(b) )
    if vtype in cond0:
        return a * b, ""
    elif vtype in cond1:
        if type(a) == int and a < 0:
            return None, "e802 : cannot mul (int_n, str|bytes)"
        elif type(b) == int and b < 0:
            return None, "e803 : cannot mul (str|bytes, int_n)"
        else:
            return a * b, ""
    elif vtype in cond2:
        return a and b, ""
    else:
        return None, f"e804 : cannot mul {vtype}"

def div(a, b):
    cond0 = { (int, int), (int, float), (float, int), (float, float) }
    vtype = ( type(a), type(b) )
    if vtype in cond0:
        return a / b, ""
    else:
        return None, f"e805 : cannot div {vtype}"

def divs(a, b):
    cond0 = { (int, int), (int, float), (float, int), (float, float) }
    vtype = ( type(a), type(b) )
    if vtype in cond0:
        return a // b, ""
    else:
        return None, f"e806 : cannot divs {vtype}"

def divr(a, b):
    cond0 = { (int, int), (int, float), (float, int), (float, float) }
    vtype = ( type(a), type(b) )
    if vtype in cond0:
        return a % b, ""
    else:
        return None, f"e807 : cannot divr {vtype}"

def pow(a, b):
    cond0 = { (int, int), (int, float), (float, int), (float, float) }
    vtype = ( type(a), type(b) )
    if vtype in cond0:
        return a ** b, ""
    else:
        return None, f"e808 : cannot pow {vtype}"

def eql(a, b):
    cond0 = { (bool, bool), (int, int), (int, float), (float, int), (float, float), (str, str), (bytes, bytes) }
    vtype = ( type(a), type(b) )
    if vtype in cond0:
        return a == b, ""
    elif a == None and b == None:
        return True, ""
    else:
        return False, ""

def sml(a, b):
    cond0 = { (int, int), (int, float), (float, int), (float, float), (str, str), (bytes, bytes) }
    vtype = ( type(a), type(b) )
    if vtype in cond0:
        return a < b, ""
    else:
        return None, f"e809 : cannot sml {vtype}"

def smle(a, b):
    cond0 = { (int, int), (int, float), (float, int), (float, float), (str, str), (bytes, bytes) }
    vtype = ( type(a), type(b) )
    if vtype in cond0:
        return a <= b, ""
    else:
        return None, f"e810 : cannot smle {vtype}"

def inc(a):
    vtype = type(a)
    if vtype == int or vtype == float:
        return a + 1, ""
    else:
        return None, f"e811 : cannot inc {vtype}"

def dec(a):
    vtype = type(a)
    if vtype == int or vtype == float:
        return a - 1, ""
    else:
        return None, f"e812 : cannot dec {vtype}"

def shm(a):
    vtype = type(a)
    if vtype == int or vtype == float or vtype == str or vtype == bytes:
        return a * 2, ""
    else:
        return None, f"e813 : cannot shm {vtype}"

def shd(a):
    vtype = type(a)
    if vtype == int or vtype == float:
        return a / 2, ""
    else:
        return None, f"e814 : cannot shd {vtype}"

def addi(a, i):
    vtype = type(a)
    if vtype == int or vtype == float:
        return a + i, ""
    else:
        return None, f"e815 : cannot addi {vtype}"

def muli(a, i):
    vtype = type(a)
    if vtype == int or vtype == float:
        return a * i, ""
    elif vtype == str or vtype == bytes:
        if i < 0:
            return None, "e816 : cannot muli int_n"
        else:
            return a * i, ""
    else:
        return None, f"e817 : cannot muli {vtype}"

def tostr(a): # get string of var
    if type(a) == bytes:
        return a.hex()
    else:
        return str(a)

def decode(data): # decode little-endian signed int 16 / 32 / 64
    if len(data) == 2:
        return struct.unpack("<h", data)[0]
    elif len(data) == 4:
        return struct.unpack("<i", data)[0]
    elif len(data) == 8:
        return struct.unpack("<q", data)[0]
    else:
        return 0

# kscript virtual machine (part 8)
class kvm:
    def __init__(self):
        self.callmem = [ ] # outercall args
        self.errmsg = "" # error message
        self.maxstk = 1048576 # stack max size
        self.info, self.abif, self.sign, self.public = "", -1, b"", "" # kelf info

        # internal data
        self.rodata, self.data, self.text = b"", b"", b""
        self.pc, self.sp, self.ma, self.mb = 0, 0, 0, 0
        self.order, self.stack = [ ], collections.deque()

    def readelf(self, path): # read file, set value
        with open(path, "rb") as f:
            data = f.read()
        worker = ksc.toolbox()
        worker.predetect = True
        worker.readb(data)
        if worker.subtype != b"KELF":
            raise Exception("e818 : not kelf file")

        header = data[ worker.chunkpos[0] + 8 : worker.chunkpos[0] + 8 + worker.chunksize[0] ]
        self.rodata = data[ worker.chunkpos[1] + 8 : worker.chunkpos[1] + 8 + worker.chunksize[1] ]
        self.data = data[ worker.chunkpos[2] + 8 : worker.chunkpos[2] + 8 + worker.chunksize[2] ]
        self.text = data[ worker.chunkpos[3] + 8 : worker.chunkpos[3] + 8 + worker.chunksize[3] ]
        if ksc.crc32hash(header) != worker.reserved[0:4]:
            raise Exception("e819 : invalid header crc32")
        
        worker = kdb.toolbox()
        worker.read( str(header, encoding="utf-8") )
        self.info = worker.get("info")[3]
        self.abif = worker.get("abi")[3]
        self.sign = worker.get("sign")[3]
        self.public = worker.get("public")[3]

    def load_data(self, code): # load rodata/data
        pos = 0
        while pos < len(code):
            if code[pos] == 78: # none
                self.stack.append(None)
                pos = pos + 1
            elif code[pos] == 66: # bool
                if code[pos + 1] == 0:
                    self.stack.append(True)
                else:
                    self.stack.append(False)
                pos = pos + 2
            elif code[pos] == 73: # int
                self.stack.append( decode( code[pos + 1:pos + 9] ) )
                pos = pos + 9
            elif code[pos] == 70: # float
                self.stack.append( struct.unpack( "<d", code[pos + 1:pos + 9] )[0] )
                pos = pos + 9
            elif code[pos] == 83: # string
                length = decode( code[pos + 1:pos + 5] )
                pos = pos + 5
                self.stack.append( str(code[pos:pos + length], encoding="utf-8") )
                pos = pos + length
            elif code[pos] == 67: # bytes
                length = decode( code[pos + 1:pos + 5] )
                pos = pos + 5
                self.stack.append( code[pos:pos + length] )
                pos = pos + length
            else:
                raise Exception("e820 : decode fail DATA")

    def load_text(self, code): # load text
        if len(code) % 8 != 0:
            raise Exception("e821 : invalid code length")
        opcond = {0: 0, 1: 0,
		16: 2, 17: 7, 18: 2, 19: 1, 20: 1, 21: 6, 22: 6,
		32: 6, 33: 6, 34: 4, 35: 4, 36: 2, 37: 2,
		48: 0, 49: 0, 50: 0, 51: 0,
		64: 0, 65: 0, 66: 0,
		80: 0, 81: 0, 82: 0, 83: 0, 84: 0, 85: 0,
		96: 2, 97: 2, 98: 2, 99: 2,
		112: 0, 113: 0, 114: 4, 115: 3}
        for i in range(0, len(code) // 8):
            temp = 8 * i
            opcode, reg, i16, i32 = code[temp], code[temp + 1], decode( code[temp + 2:temp + 4] ), decode( code[temp + 4:temp + 8] )
            if opcode not in opcond:
                raise Exception("e822 : invalid opcode")
            cond = opcond[opcode]
            if cond // 4 == 1 and reg != 97 and reg != 98:
                raise Exception("e823 : invalid register")
            if (cond % 4) // 2 == 1 and i16 < 0:
                raise Exception("e824 : invalid int16")
            if cond % 2 == 1 and i32 < 0:
                raise Exception("e825 : invalid int32")
            self.order.append( [opcode, reg, i16, i32] )

    def getpos(self, i16, i32): # seg + addr -> pos
        if i16 == 108:
            return self.sp + i32
        else:
            return i32

    def forcalc(self, reg, i16, i32, iscond): # for (i, r <- v) calc
        pos = self.getpos(i16, i32)
        if pos < 0 or pos >= len(self.stack):
            return None, "e826 : ram access fail"
        if reg == 97:
            idx = self.ma
        else:
            idx = self.mb
        value = self.stack[pos]
        vtype = type(value)
        if vtype == int:
            if value > 0:
                if iscond:
                    return idx < value, ""
                else:
                    return idx, ""
            elif value < 0:
                if iscond:
                    return -idx > value, ""
                else:
                    return -idx, ""
            else:
                if iscond:
                    return False, ""
                else:
                    return None, ""
        elif vtype == str:
            if iscond:
                return idx < len(value), ""
            else:
                return value[idx], ""
        elif vtype == bytes:
            if iscond:
                return idx < len(value), ""
            else:
                return bytes( [ value[idx] ] ), ""
        else:
            return False, "e827 : invalid for type"

    def cycle(self): # fetch - decode - execute - interupt (0 normal, 1 hlt, 2 err, -1 c_err)
        # fetch
        if self.pc >= len(self.order):
            self.errmsg = "e828 : order access fail"
            return -1
        op, reg, i16, i32 = self.order[self.pc]
        unib, dnib, self.pc = op >> 4, op & 0x0f, self.pc + 1

        # decode & execute
        if unib == 6:
            pos = self.getpos(i16, i32)
            if pos < 0 or pos >= len(self.stack):
                self.errmsg = "e839 : stack access fail"
                return -1
            if i16 == 99:
                self.errmsg = "e840 : cannot write to const"
                return -1
            err = ""
            if dnib == 0: # inc
                self.stack[pos], err = inc( self.stack[pos] )
            elif dnib == 1: # dec
                self.stack[pos], err = dec( self.stack[pos] )
            elif dnib == 2: # shm
                self.stack[pos], err = shm( self.stack[pos] )
            elif dnib == 3: # shd
                self.stack[pos], err = shd( self.stack[pos] )
            if err != "":
                self.errmsg = err
                return 2

        elif unib == 7:
            if dnib == 0: # addi
                temp, err = addi(self.stack.pop(), i32)
                self.stack.append(temp)
                if err != "":
                    self.errmsg = err
                    return 2

            elif dnib == 1: # muli
                temp, err = muli(self.stack.pop(), i32)
                self.stack.append(temp)
                if err != "":
                    self.errmsg = err
                    return 2

            elif dnib == 2: # addr
                if reg == 97:
                    temp, err = addi(self.ma, i32)
                else:
                    temp, err = addi(self.mb, i32)
                self.stack.append(temp)
                if err != "":
                    self.errmsg = err
                    return 2

            elif dnib == 3: # jmpi
                temp, err = False, ""
                if i16 == 1:
                    temp, err = eql(self.ma, self.mb)
                elif i16 == 2:
                    temp, err = eql(self.ma, self.mb)
                    temp = not temp
                elif i16 == 3:
                    temp, err = sml(self.ma, self.mb)
                elif i16 == 4:
                    temp, err = sml(self.mb, self.ma)
                elif i16 == 5:
                    temp, err = smle(self.ma, self.mb)
                elif i16 == 6:
                    temp, err = smle(self.mb, self.ma)
                else:
                    self.errmsg = "e841 : decode fail JMPI"
                    return -1
                if type(temp) != bool or not temp:
                    self.pc = i32
                if err != "":
                    self.errmsg = err
                    return 2

        elif unib == 1:
            if dnib == 0: # intr
                if i16 > len(self.stack):
                    self.errmsg = "e829 : stack access fail"
                    return -1
                self.callmem = [0] * i16
                for i in range(0, i16):
                    self.callmem[i16 - i - 1] = self.stack.pop()
                return i32

            elif dnib == 1: # call
                if len(self.stack) > self.maxstk - 4:
                    self.errmsg = "e830 : stack overflow"
                    return -1
                temp = len(self.stack)
                if reg == 97:
                    self.stack.append(self.ma)
                else:
                    self.stack.append(self.mb)
                self.stack.append(self.pc)
                self.stack.append(self.sp)
                for i in range(0, i16 + 1):
                    self.stack.append(None)
                self.pc, self.sp = i32, temp

            elif dnib == 2: # ret
                temp = self.sp
                self.pc, self.sp, self.ma = self.stack[temp + 1], self.stack[temp + 2], self.stack[temp + 3]
                if temp < i16:
                    self.errmsg = "e831 : stack access fail"
                    return -1
                for i in range(0, len(self.stack) - temp + i16):
                    self.stack.pop()

            elif dnib == 3: # jmp
                self.pc = i32

            elif dnib == 4: # jmpiff
                temp = self.stack.pop()
                if type(temp) == bool:
                    if not temp:
                        self.pc = i32
                else:
                    self.pc = i32
                    self.errmsg = "e832 : invalid condition"
                    return 2

            elif dnib == 5: # forcond
                temp, err = self.forcalc(reg, i16, i32, True)
                self.stack.append(temp)
                if err != "":
                    self.errmsg = err
                    return 2

            elif dnib == 6: # forset
                self.ma, _ = self.forcalc(reg, i16, i32, False)

        elif unib == 2:
            if dnib == 0: # load
                pos = self.getpos(i16, i32)
                if pos < 0 or pos >= len(self.stack):
                    self.errmsg = "e833 : stack access fail"
                    return -1
                if reg == 97:
                    self.ma = self.stack[pos]
                else:
                    self.mb = self.stack[pos]

            elif dnib == 1: # store
                pos = self.getpos(i16, i32)
                if pos < 0 or pos >= len(self.stack):
                    self.errmsg = "e834 : stack access fail"
                    return -1
                if i16 == 99:
                    self.errmsg = "e835 : cannot write to const"
                    return -1
                if reg == 97:
                    self.stack[pos] = self.ma
                else:
                    self.stack[pos] = self.mb

            elif dnib == 2: # push
                if reg == 97:
                    self.stack.append(self.ma)
                else:
                    self.stack.append(self.mb)

            elif dnib == 3: # pop
                if reg == 97:
                    self.ma = self.stack.pop()
                else:
                    self.mb = self.stack.pop()

            elif dnib == 4: # pushset
                pos = self.getpos(i16, i32)
                if pos < 0 or pos >= len(self.stack):
                    self.errmsg = "e836 : stack access fail"
                    return -1
                self.stack.append( self.stack[pos] )

            elif dnib == 5: # popset
                pos = self.getpos(i16, i32)
                if pos < 0 or pos >= len(self.stack):
                    self.errmsg = "e837 : stack access fail"
                    return -1
                if i16 == 99:
                    self.errmsg = "e838 : cannot write to const"
                    return -1
                temp = self.stack.pop()
                self.stack[pos] = temp

        elif unib == 3:
            b = self.stack.pop()
            a = self.stack.pop()
            temp, err = None, ""
            if dnib == 0: # add
                temp, err = add(a, b)
            elif dnib == 1: # sub
                temp, err = sub(a, b)
            elif dnib == 2: # mul
                temp, err = mul(a, b)
            elif dnib == 3: # div
                temp, err = div(a, b)
            self.stack.append(temp)
            if err != "":
                self.errmsg = err
                return 2

        elif unib == 4:
            b = self.stack.pop()
            a = self.stack.pop()
            temp, err = None, ""
            if dnib == 0: # divs
                temp, err = divs(a, b)
            elif dnib == 1: # divr
                temp, err = divr(a, b)
            elif dnib == 2: # pow
                temp, err = pow(a, b)
            self.stack.append(temp)
            if err != "":
                self.errmsg = err
                return 2

        elif unib == 5:
            b = self.stack.pop()
            a = self.stack.pop()
            temp, err = None, ""
            if dnib == 0: # eql
                temp, err = eql(a, b)
            elif dnib == 1: # eqln
                temp, err = eql(a, b)
                temp = not temp
            elif dnib == 2: # sml
                temp, err = sml(a, b)
            elif dnib == 3: # grt
                temp, err = sml(b, a)
            elif dnib == 4: # smle
                temp, err = smle(a, b)
            elif dnib == 5: # grte
                temp, err = smle(b, a)
            self.stack.append(temp)
            if err != "":
                self.errmsg = err
                return 2

        elif unib == 0:
            if dnib == 0: # hlt
                return 1

        # interupt
        return 0

    # view kelf file, update internal field -> (info, abif, public)
    def view(self, path):
        self.readelf(path)
        return self.info, self.abif, self.public

    # load kelf file (should done view first), check sign
    def load(self, sign):
        if len(self.text) == 0:
            raise Exception("e842 : should done kvm.view() first")
        if sign and self.public != "":
            hvalue = hashlib.sha3_512(self.rodata + self.data + self.text).digest()
            worker = ksign.toolbox()
            if not worker.verify(self.public, self.sign, hvalue):
                raise Exception("e843 : ksign verify fail")
        self.load_data(self.rodata)
        self.load_data(self.data)
        self.load_text(self.text)
    
    # run code, exit with intr code, !!! set outercall return at kvm.ma !!!
    def run(self):
        self.callmem = [ ]
        self.errmsg = ""
        try:
            intr = self.cycle()
            return intr
        except Exception as e:
            self.errmsg = str(e)
            return -1

# test io support (input, print, read, write, time, sleep)
def testio(mode, v):
    out = None
    if mode == 16: # input(v)
        out = input( tostr( v[0] ) )
    elif mode == 17: # print(v)
        print(tostr( v[0] ), end="")
    elif mode == 18: # read(s, i)
        if type( v[0] ) == str and type( v[1] ) == int:
            try:
                with open(v[0], "rb") as f:
                    if v[1] < 0:
                        out = f.read()
                    else:
                        out = f.read( v[1] )
            except:
                out = b""
    elif mode == 19: # write(s, s|b)
        try:
            if type( v[0] ) == str and type( v[1] ) == str:
                with open(v[0], "w", encoding="utf-8") as f:
                    f.write( v[1] )
            elif type( v[0] ) == str and type( v[1] ) == bytes:
                with open(v[0], "wb") as f:
                    f.write( v[1] )
        except:
            pass
    elif mode == 20: # time()
        out = time.time()
    elif mode == 21: # sleep(i|f)
        if type( v[0] ) == int or type( v[0] ) == float:
            time.sleep( v[0] )
    return out
