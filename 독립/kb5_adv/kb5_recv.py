# test608 : kboom5 adv reciever test

import io
import time
import os
import zlib

import tkinter
import tkinter.messagebox
from PIL import Image, ImageTk

import multiprocessing as mp

import kdb
import kaes
import kcom
import simen

# ===== kcom key transfer =====

# skey 4B, sport int -(fpath S, fkey 48B)-> (errmsg S)
def sendkey(skey, sport, fpath, fkey):
    fpath = os.path.abspath(fpath).replace("\\", "/")
    skey = skey * ( 48 // len(skey) )
    tbox0 = kdb.toolbox()
    tbox1 = simen.toolbox()
    tbox2 = kcom.server()
    
    tbox0.readstr("path = 0\nkey = 0")
    tbox0.fixdata("path", fpath)
    tbox0.fixdata("key", fkey)
    data = bytes(tbox0.writestr(), encoding="utf-8") # 평문 데이터
    crcv = str( zlib.crc32(data) ) # 평문 데이터의 crc32 값
    
    tbox1.setkey(skey)
    enc = tbox1.encrypt(data) # 전송할 암호문
    tbox2.close = 90 # 시간제한 90s
    tbox2.port = sport
    tbox2.msg = crcv

    try:
        tbox2.send(enc)
        errmsg = ""
    except Exception as e:
        errmsg = str(e)
    return errmsg

# skey 4B, sport int -> (errmsg S, fpath S, fkey 48B)
def recievekey(skey, sport):
    tbox0 = kdb.toolbox()
    tbox1 = simen.toolbox()
    tbox2 = kcom.client()

    skey = skey * ( 48 // len(skey) )
    tbox2.close = 60 # 시간제한 60s
    tbox2.port = sport

    try:
        enc = tbox2.recieve()
        crcv = tbox2.msg

        tbox1.setkey(skey)
        data = tbox1.decrypt(enc)
        if zlib.crc32(data) != int(crcv):
            raise Exception("wrong crc32 value")

        tbox0.readstr( str(data, encoding="utf-8") )
        fpath = tbox0.getdata("path")[3]
        fkey = tbox0.getdata("key")[3]
        
        errmsg = ""
    except Exception as e:
        errmsg = str(e)

    if errmsg == "":
        return "", fpath, fkey
    else:
        return errmsg, "", b""

# enc kf path S + key 48B -> plain kf B
def readkey(path, key):
    tbox = kaes.funcbytes()
    with open(path, "rb") as f:
        enc = f.read()
    return tbox.de(key, enc)

def main():
    def getdata():
        time.sleep(0.1)
        addr = ent.get()
        port, key = kcom.unpack(addr)

        err, fp, fk = recievekey(key, port)
        if err != "":
            tkinter.messagebox.showinfo(title='수신 오류', message=f' {err} ')

        else:
            idata = readkey(fp, fk)
            with open("recv", "wb") as f:
                f.write(idata)

            try:
                pimg = Image.open( io.BytesIO(idata) )
                iw, ih = pimg.size
                ratio = min(370 / iw, 370 / ih)
                sw, sh = int(iw * ratio), int(ih * ratio)
                rimg = pimg.resize( (sw, sh), Image.LANCZOS )
                nonlocal rpimg
                rpimg = ImageTk.PhotoImage(rimg)
                canvas.create_image(5, 5, anchor=tkinter.NW, image=rpimg)
            except Exception as e:
                tkinter.messagebox.showinfo(title='사진 표시 오류', message=f' {e} ')

    win = tkinter.Tk()
    win.title("test608")
    win.geometry("400x500+200+100")
    win.resizable(False, False)

    lbl = tkinter.Label(win, font=("Consolas", 14), text="주소")
    lbl.place(x=5, y=5)
    ent = tkinter.Entry(win, font=("Consolas", 14), width=26)
    ent.place(x=50, y=10)
    but = tkinter.Button(win, font=("Consolas", 12), text="recieve", command=getdata)
    but.place(x=320, y=5)

    canvas = tkinter.Canvas(win, width=390, height=390)
    canvas.place(x=5, y=105)
    rpimg = None

    win.mainloop()

if __name__ == "__main__":
    time.sleep(0.5)
    mp.freeze_support()
    main()
    time.sleep(0.5)
