# test618 : andcom5

import time
import os
import shutil
import random

import tkinter
import tkinter.filedialog

import py7zr
import zipfile

def zip7z(pw, path): # notin618_temp/ -> path(7z)
    with py7zr.SevenZipFile(path, "w", password=pw) as f:
        f.writeall("notin618_temp")

def unzip7z(tgt, pw, path): # tgt(7z) -> path(/)
    with py7zr.SevenZipFile(tgt, mode="r", password=pw) as f:
        f.extractall(path)

def zippk(path): # cur dir files -> notin618_temp/path(zip)
    names = [ x for x in os.listdir("./") if ("notin618_" not in x) and ("st5adv_" not in x) ]
    with zipfile.ZipFile(f"notin618_temp/{path}", "w") as f:
        for i in names:
            f.write(i, compress_type=zipfile.ZIP_DEFLATED)
    for i in names:
        os.remove(i)

def unzippk(tgt, path): # tgt(zip) -> path(/)
    with zipfile.ZipFile(tgt, "r") as f:
        f.extractall(path)

def dclear(mode): # clear notin618_temp/, mknew? mode
    if os.path.exists("./notin618_temp"):
        shutil.rmtree("./notin618_temp")
    if mode:
        os.mkdir("./notin618_temp")

def fclear(): # clear cur dir files
    names = [ x for x in os.listdir("./") if ("notin618_" not in x) and ("st5adv_" not in x) ]
    for i in names:
        os.remove(i)

def getf(paths, tof): # tgt files -> tof/tgt
    paths = [os.path.abspath(x).replace("\\", "/") for x in paths]
    for i in paths:
        shutil.copyfile( i, tof + i[i.rfind("/"):] )

def getd(): # get desktop dir
    # windows
    desk = os.path.expanduser("~/Desktop").replace("\\", "/")
    # linux
    #desk = os.path.expanduser("~/Desktop")
    #if not os.path.exists(desk):
    #desk = os.path.expanduser("~/바탕화면")
    if desk[-1] != "/":
        desk = desk + "/"
    return desk

def delf(paths): # del tgt files
    for i in paths:
        os.remove(i)

def mono(enfiles, defile, pw0, pw1, mode): # only 7z en/de
    res = ""
    desk = getd()
    try:
        if enfiles != [ ]:
            if pw0 == pw1:
                dclear(True)
                getf(enfiles, "notin618_temp")
                temp = desk + str( random.randrange(1000,10000) ) + '.7z'
                zip7z(pw0, temp)
                dclear(False)
                if mode:
                    delf(enfiles)
                res = temp
            else:
                res = "PW not match"

        elif defile != "":
            temp = desk + str( random.randrange(1000,10000) )
            os.mkdir(temp)
            unzip7z(defile, pw0, temp)
            if mode:
                delf( [defile] )
            res = temp

        else:
            res = "Done Nothing"
    except Exception as e:
        res = str(e)
    return res

def dual(enfiles, defile, pw0, pw1, mode): # both zip - 7z en/de
    res = ""
    desk = getd()
    try:
        if enfiles != [ ]:
            if pw0 == pw1:
                dclear(True)
                fclear()
                getf(enfiles, ".")
                temp = str( random.randrange(1000,10000) ) + '.zip'
                zippk(temp)
                temp = desk + str( random.randrange(1000,10000) ) + '.7z'
                zip7z(pw0, temp)
                dclear(False)
                fclear()
                if mode:
                    delf(enfiles)
                res = temp
            else:
                res = "PW not match"

        elif defile != "":
            dclear(False)
            fclear()
            unzip7z(defile, pw0, "notin618_temp")
            temp = os.listdir("./notin618_temp")
            if temp == ["notin618_temp"]:
                temp = os.listdir("notin618_temp/notin618_temp")
                zname = temp[0]
                shutil.move(f"notin618_temp/notin618_temp/{zname}", zname)
            else:
                zname = temp[0]
                shutil.move(f"notin618_temp/{zname}", zname)
            temp = desk + str( random.randrange(1000,10000) )
            os.mkdir(temp)
            unzippk(zname, temp)
            dclear(True)
            fclear()
            if mode:
                delf( [defile] )
            res = temp

        else:
            res = "Done Nothing"
    except Exception as e:
        res = str(e)
    return res

def maingui():
    win = tkinter.Tk()
    win.title('AndCom5')
    win.geometry("450x350+100+50") # lxgp= 750x600+200+100
    win.resizable(False, False)

    enfiles = [ ]
    defile = [ ]
    dmode = False
    zmode = False

    def func0(): # get enfiles
        time.sleep(0.1)
        nonlocal enfiles
        enfiles = [ ]
        files = tkinter.filedialog.askopenfiles(title='파일들 선택')
        for i in files:
            enfiles.append(i.name)
        if enfiles != [ ]:
            ens.set(f"({len(enfiles)}) {enfiles[0]}")
        else:
            ens.set("(0) - ")
        nonlocal defile
        defile = ""
        des.set("(0) - ")
        win.update()

    def func1(): # get defile
        time.sleep(0.1)
        try:
            file = tkinter.filedialog.askopenfile( title="파일 선택", filetypes=( ('7z files', '*.7z'), ('all files', '*.*') ) ).name
        except:
            file = ""
        nonlocal enfiles
        enfiles = [ ]
        ens.set("(0) - ")
        nonlocal defile
        defile = file
        if defile != "":
            des.set(f"(1) {defile}")
        else:
            des.set("(0) - ")
        win.update()

    def func2(): # toggle del mode
        time.sleep(0.1)
        nonlocal dmode
        if chkv2.get() == 0:
            dmode = False
        else:
            dmode = True

    def func3(): # toggle 7z mode
        time.sleep(0.1)
        nonlocal zmode
        if chkv3.get() == 0:
            zmode = False
        else:
            zmode = True

    def func4(): # execute
        time.sleep(0.1)
        if zmode:
            res = mono(enfiles, defile, in5.get(), in6.get(), dmode)
        else:
            res = dual(enfiles, defile, in5.get(), in6.get(), dmode)
        if "/" in res:
            res = res[res.rfind("/") + 1:]
        status.set(res)
        win.update()

    but0 = tkinter.Button(win, text="EN", font=("Consolas", 20), command=func0)
    but0.place(x=5, y=5)
    ens = tkinter.StringVar()
    ens.set("(0) - ")
    lbl0 = tkinter.Label(win, font=("Consolas", 16), textvariable=ens)
    lbl0.place(x=60, y=15) # lxgp= x=90 y=15

    but1 = tkinter.Button(win, text="DE", font=("Consolas", 20), command=func1)
    but1.place(x=5, y=65) # lxgp= x=5 y=105
    des = tkinter.StringVar()
    des.set("(0) - ")
    lbl1 = tkinter.Label(win, font=("Consolas", 16), textvariable=des)
    lbl1.place(x=60, y=75) # lxgp= x=90 y=115

    chkv2 = tkinter.IntVar()
    chkb2 = tkinter.Checkbutton(win, text="원본 삭제", font=("Consolas", 16), variable=chkv2, command=func2)
    chkb2.place(x=5, y=230) # lxgp= x=5 y=400

    chkv3 = tkinter.IntVar()
    chkb3 = tkinter.Checkbutton(win, text="7z 형식만", font=("Consolas", 16), variable=chkv3, command=func3)
    chkb3.place(x=140, y=230) # lxgp= x=260 y=400

    but4 = tkinter.Button(win, text="  G O  ", font=("Consolas", 20), command=func4)
    but4.place(x=290,y=220) # lxgp= x=520 y=400
    status = tkinter.StringVar()
    status.set("idle")
    lbl4 = tkinter.Label(win, font=("Consolas", 16), textvariable=status)
    lbl4.place(x=5, y=300) # lxgp= x=5 y=500

    lbl5 = tkinter.Label(win, font=("Consolas", 16), text="pw")
    lbl5.place(x=5, y=140) # lxgp= x=5 y=220
    in5 = tkinter.Entry(width=28, font=("Consolas", 16), show="*")
    in5.place(x=80, y=140) # lxgp= x=80 y=220

    lbl6 = tkinter.Label(win, font=("Consolas", 16), text="pw")
    lbl6.place(x=5, y=175) # lxgp= x=5 y=295
    in6 = tkinter.Entry(width=28, font=("Consolas", 16), show="*")
    in6.place(x=80, y=175) # lxgp= x=80 y=295

    win.mainloop()

dclear(True)
fclear()
maingui()
dclear(True)
fclear()
time.sleep(0.5)
