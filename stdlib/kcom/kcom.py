# test625 : stdlib5.kcom

import random
import zlib

import socket
import requests
import bs4

import kobj

# get domain txt from (http~ *.html)
def gettxt(url, domain):
    html = requests.get(url).text
    dom = bs4.BeautifulSoup(html, "html.parser")
    data = dom.find("p", id=domain)
    if data == None:
        raise Exception(f"No Domain : {domain}")
    else:
        return data.text
    
# download binary name -> path from (http~ */) + (*.num)
def download(url, name, num, path, proc):
    proc[0] = 0.0
    if url[-1] != "/":
        url = url + "/"
    with open(path, "wb") as f:
        for i in range(0, num):
            proc[0] = i / num
            html = requests.get(url + name + f".{i}")
            if html.content == None:
                raise Exception(f"No Data : {name}.{i}")
            else:
                f.write(html.content)
    proc[0] = 2.0

# (port int, key 4B) -> address str
def pack(port, key):
    out = [""] * 6
    num = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"]
    port = kobj.encode(port, 2)
    for i in range(0, 2):
        out[i] = num[port[i] // 16] + num[port[i] % 16]
    for i in range(0, 4):
        out[i + 2] = num[key[i] // 16] + num[key[i] % 16]
    return f"{out[0]}{out[1]}.{out[2]}{out[3]}.{out[4]}{out[5]}"

# address str -> (port int, key 4B)
def unpack(address):
    address = address.lower().replace(".", "")
    num = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"]
    temp = [0] * 6
    for i in range(0, 6):
        j = 0
        while num[j] != address[2 * i]:
            j = j + 1
        k = 0
        while num[k] != address[2 * i + 1]:
            k = k + 1
        temp[i] = 16 * j + k
    temp = bytes(temp)
    port = kobj.decode( temp[0:2] )
    key = temp[2:6]
    return port, key

class node:
    def __init__(self):
        self.ipv6 = True
        self.port = 13600
        self.close = 150

    # send simple data (data B, key 4B)
    def send(self, data, key):
        temp = [0] * len(data)
        for i in range( 0, len(data) ):
            temp[i] = ( data[i] + key[i % 4] ) % 256
        msg = kobj.encode(len(data), 8) + kobj.encode(zlib.crc32(data), 4) + bytes(temp)

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

        clisc, addr = svrsc.accept()
        mnum = clisc.recv(8)
        if mnum[0:5] == b"kcom5":
            clisc.sendall(mnum)
        else:
            raise Exception("invalid connection")
        
        clisc.sendall(msg)
        clisc.close()
        svrsc.close()

    # recieve data (key 4B)
    def recieve(self, key):
        if self.ipv6:
            iprange = socket.AF_INET6
            ipaddr = '::1'
            protocol = socket.SOCK_STREAM
        else:
            iprange = socket.AF_INET
            ipaddr = '127.0.0.1'
            protocol = socket.SOCK_STREAM
        clisc = socket.socket(iprange, protocol)
        if self.close > 0:
            clisc.settimeout(self.close)
        clisc.connect( (ipaddr, self.port) )

        mnum = b"kcom5" + random.randbytes(3)
        clisc.send(mnum)
        chk = clisc.recv(8)
        if mnum != chk:
            raise Exception("invalid connection")

        mlen = kobj.decode( clisc.recv(8) )
        crcv = kobj.decode( clisc.recv(4) )
        data = clisc.recv(mlen)
        temp = [0] * mlen
        for i in range(0, mlen):
            temp[i] = ( data[i] - key[i % 4] ) % 256
        data = bytes(temp)

        if zlib.crc32(data) != crcv:
            raise Exception("invalid key")
        clisc.close()
        return data
