import kaesst as m
#import kaeshy as m

import multiprocessing as mp
import time
import os
import shutil

def r(path):
    with open(path, "rb") as f:
        temp = f.read()
    return temp

def w(path, data):
    with open(path, "wb") as f:
        f.write(data)

def main():
    if os.path.exists("./temp583"):
        shutil.rmtree("./temp583")
    os.mkdir("./temp583")
    print("===== gen5 KAES test =====")
    
    size = [0, 524272, 9352671, 211104879, 524288, 6291456, 16777216, 36175872]
    pw = m.genrandom(48)
    kf = m.genkf("nopath")
    hint = m.genrandom(128)
    msg = "동성고 최고 귀요미 김병관"
    pre = [m.genrandom(x) for x in size]
    enc = [b""] * len(size)
    plain = [b""] * len(size)

    k0 = m.genbytes()
    k1 = m.genfile()
    k2 = m.funcbytes()
    k3 = m.funcfile()
    k0.msg = msg
    k1.msg = msg
    if k0.valid:
        print("k0 valid")
    else:
        print("k0 invalid")
    if k1.valid:
        print("k1 valid")
    else:
        print("k1 invalid")
    if k2.valid:
        print("k2 valid")
    else:
        print("k2 invalid")
    if k3.valid:
        print("k3 valid")
    else:
        print("k3 invalid")

    print()

    for i in range(0, 8):
        enc[i] = k0.en(pw, kf, hint, pre[i])
        a, b, c = k0.view(enc[i])
        print(f"test {i} : hint_{a == hint}, msg_{b == msg}")
        plain[i] = k0.de(pw, kf, enc[i], c)
        print(f"test {i} : data_{pre[i] == plain[i]}")

    print()

    for i in range(0, 8):
        w("temp583\\t", pre[i])
        enc[i] = k1.en(pw, kf, hint, "temp583\\t")
        a, b, c = k1.view(enc[i])
        print(f"test {i} : hint_{a == hint}, msg_{b == msg}")
        plain[i] = k1.de(pw, kf, enc[i], c)
        plain[i] = r(plain[i])
        print(f"test {i} : data_{pre[i] == plain[i]}")

    print()

    for i in range(0, 8):
        enc[i] = k2.en(pw, pre[i])
        plain[i] = k2.de(pw, enc[i])
        print(f"test {i} : data_{pre[i] == plain[i]}")

    print()

    for i in range(0, 8):
        w("temp583/d", pre[i])
        k3.en(pw, "temp583/d", "temp583/e")
        k3.de(pw, "temp583/e", "temp583/d")
        plain[i] = r("temp583/d")
        print(f"test {i} : data_{pre[i] == plain[i]}")

def bench():
    k0 = m.genbytes()
    k1 = m.genfile()
    k2 = m.funcbytes()
    k3 = m.funcfile()
    pw = m.genrandom(48)
    kf = m.genkf("nopath")
    hint = b""
    print()

    time.sleep(5)
    t0 = time.time()
    data = m.genrandom(2 * 1024 * 1024 * 1024) # 2048 MiB
    t1 = time.time()
    t2 = t1 - t0
    print(f"rand : time {t2} s, speed {2048 / t2} MiB/s")

    time.sleep(5)
    t0 = time.time()
    temp = k0.en(pw, kf, hint, data)
    t1 = time.time()
    t2 = t1 - t0
    print(f"k0en : time {t2} s, speed {2048 / t2} MiB/s")
    del data

    time.sleep(5)
    a, b, c = k0.view(temp)
    t0 = time.time()
    k0.de(pw, kf, temp, c)
    t1 = time.time()
    t2 = t1 - t0
    print(f"k0de : time {t2} s, speed {2048 / t2} MiB/s")
    del temp

    data = m.genrandom(4 * 1024 * 1024 * 1024) # 4096 MiB
    time.sleep(5)
    w("temp583\\b", data)
    del data
    t0 = time.time()
    temp = k1.en(pw, kf, hint, "temp583\\b")
    t1 = time.time()
    t2 = t1 - t0
    print(f"k1en : time {t2} s, speed {4096 / t2} MiB/s")

    time.sleep(5)
    a, b, c = k1.view(temp)
    t0 = time.time()
    k1.de(pw, kf, temp, c)
    t1 = time.time()
    t2 = t1 - t0
    print(f"k1de : time {t2} s, speed {4096 / t2} MiB/s")

    data = m.genrandom(1 * 1024 * 1024 * 1024) # 1024 MiB
    time.sleep(5)
    t0 = time.time()
    temp = k2.en(pw, data)
    t1 = time.time()
    t2 = t1 - t0
    print(f"k2en : time {t2} s, speed {1024 / t2} MiB/s")
    del data

    time.sleep(5)
    t0 = time.time()
    k2.de(pw, temp)
    t1 = time.time()
    t2 = t1 - t0
    print(f"k2de : time {t2} s, speed {1024 / t2} MiB/s")
    del temp

    data = m.genrandom(4 * 1024 * 1024 * 1024) # 4096 MiB
    time.sleep(5)
    w("temp583/bd", data)
    del data
    t0 = time.time()
    temp = k3.en(pw, "temp583/bd", "temp583/be")
    t1 = time.time()
    t2 = t1 - t0
    print(f"k3en : time {t2} s, speed {4096 / t2} MiB/s")

    time.sleep(5)
    t0 = time.time()
    k3.de(pw, "temp583/be", "temp583/bd")
    t1 = time.time()
    t2 = t1 - t0
    print(f"k3de : time {t2} s, speed {4096 / t2} MiB/s")
    
if __name__ == '__main__':
    mp.freeze_support()
    main()
    bench()
