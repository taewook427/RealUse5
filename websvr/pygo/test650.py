import time
import ctypes

import kobj

ext = ctypes.CDLL("./ex.dll")

s0, s1 = kobj.call("b", "")
ext.exf0.argtypes, ext.exf0.restype = s0, s1
s0, s1 = kobj.call("bi", "b")
ext.exf1.argtypes, ext.exf1.restype = s0, s1
s0, s1 = kobj.call("bi", "b")
ext.exf2.argtypes, ext.exf2.restype = s0, s1
s0, s1 = kobj.call("bi", "")
ext.exf3.argtypes, ext.exf3.restype = s0, s1
s0, s1 = kobj.call("", "b")
ext.exf4.argtypes, ext.exf4.restype = s0, s1
s0, s1 = kobj.call("", "")
ext.exf5.argtypes, ext.exf5.restype = s0, s1
s0, s1 = kobj.call("", "f")
ext.exf6.argtypes, ext.exf6.restype = s0, s1

def test0():
    v0 = bytearray(1024 * 1024 * 980) # 980 MiB
    print("v0 생성")
    time.sleep(3)

    for i in range( 0, len(v0) ):
        v0[i] = i % 4
    v0 = bytes(v0)
    print("v0 바이트화")
    time.sleep(3)

    p0, p1 = kobj.send(v0)
    print("ctypes 포인터 생성")
    time.sleep(3)

    print("dll 진입")
    o0 = ext.exf1(p0, p1)
    print("dll 탈출")
    time.sleep(3)

    v1 = kobj.recv( o0, len(v0) )
    print(f"v1 생성 : {v1[0:16]}")
    time.sleep(3)

    ext.exf0(o0)
    print("포인터 해제")
    time.sleep(3)

    print("함수 탈출")

def test1():
    v0 = bytearray(1024 * 1024 * 980) # 980 MiB
    for i in range( 0, len(v0) ):
        v0[i] = i % 8
    v0 = bytes(v0)
    time.sleep(1.5)

    t0 = time.time()
    p0, p1 = kobj.send(v0)
    o0 = ext.exf2(p0, p1)
    v1 = kobj.recv( o0, len(v0) )
    ext.exf0(o0)
    print(v1[0:16])
    print(time.time() - t0)
    time.sleep(1.5)

def test2():
    p0, p1 = kobj.send(b"Hello, world!")
    ext.exf3(p0, p1)
    time.sleep(3)
    o0 = ext.exf4()
    print( kobj.recvauto(o0) )
    ext.exf0(o0)

def test3():
    t = ext.exf6()
    ext.exf5()
    print("작업 시작")
    while t < 1:
        print(t)
        t = ext.exf6()
        time.sleep(0.3)
    print(t)
