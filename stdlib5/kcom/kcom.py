# test581 : kcom (py)

import socket
import random

def encode(num, length): # little endian encoding
    temp = [0] * length
    for i in range(0, length):
        temp[i] = num % 256
        num = num // 256
    return bytes(temp)

def decode(data): # little endian decoding
    temp = 0
    for i in range( 0, len(data) ):
        if data[i] != 0:
            temp = temp + data[i] * 256 ** i
    return temp

# port int + key bytes -> address str
def pack(port, key):
    if 48 % len(key) != 0:
        raise Exception("invalid keylen")
    temp = [""] * len(key)
    num = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"]
    for i in range( 0, len(key) ):
        temp[i] = num[key[i] // 16] + num[key[i] % 16]
    return f"{port}.{'.'.join(temp)}"

# address str -> port int + key 48B
def unpack(address):
    address = address.lower()
    temp = address.split(".")
    port = int( temp[0] )
    temp = temp[1:]
    if 48 % len(temp) != 0:
        raise Exception("invalid keylen")
    key = [0] * len(temp)
    num = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"]
    for i in range( 0, len(temp) ):
        j = 0
        while temp[i][0] != num[j]:
            j = j + 1
        k = 0
        while temp[i][1] != num[k]:
            k = k + 1
        key[i] = 16 * j + k
    return port, bytes(key) * ( 48 // len(key) )

# transmit data
class server:
    def __init__(self):
        self.ipv6 = True
        self.port = 13600
        self.close = 150
        self.msg = ""

    def send(self, data):
        # 서버 열기
        if self.ipv6:
            iprange = socket.AF_INET6
            ipaddr = '::1'
            protocol = socket.SOCK_STREAM
        else:
            iprange = socket.AF_INET
            ipaddr = '127.0.0.1'
            protocol = socket.SOCK_STREAM
        svrsc = socket.socket(iprange, protocol)
        svrsc.bind( (ipaddr, self.port) )
        svrsc.listen(1)
        if self.close > 0:
            svrsc.settimeout(self.close)

        # KCOM5 + rand 3B echo
        clisc, addr = svrsc.accept()
        mnum = clisc.recv(8)
        if mnum[0:5] == b"KCOM5":
            clisc.sendall(mnum)
        else:
            raise Exception("invalid connection")

        # data transmit
        localmsg = bytes(self.msg, encoding="utf-8")
        clisc.sendall(encode(len(localmsg), 4) + localmsg)
        clisc.sendall(encode(len(data), 8) + data)

        # socket close
        clisc.close()
        svrsc.close()

# recieve data
class client:
    def __init__(self):
        self.ipv6 = True
        self.port = 13600
        self.close = 150
        self.msg = ""

    def recieve(self):
        # 서버 접속
        if self.ipv6:
            iprange = socket.AF_INET6
            ipaddr = '::1' # linux에서는 "localhost"
            protocol = socket.SOCK_STREAM
        else:
            iprange = socket.AF_INET
            ipaddr = '127.0.0.1'
            protocol = socket.SOCK_STREAM
        clisc = socket.socket(iprange, protocol)
        if self.close > 0:
            clisc.settimeout(self.close)
        clisc.connect( (ipaddr, self.port) )

        # KCOM5 + rand 3B echo
        mnum = b"KCOM5" + random.randbytes(3)
        clisc.send(mnum)
        chk = clisc.recv(8)
        if mnum != chk:
            raise Exception("invalid connection")

        # data recieve
        mlen = decode( clisc.recv(4) )
        localmsg = clisc.recv(mlen)
        dlen = decode( clisc.recv(8) )
        data = clisc.recv(dlen)

        # socket close
        clisc.close()
        self.msg = str(localmsg, encoding="utf-8")
        return data
