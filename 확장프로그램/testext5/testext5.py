# test613 : ST5adv example user extension

import time

import ctypes

import tkinter

def load():
    myos = "windows" # "linux"
    if myos == "windows":
        path = "../testlib5/testdll.dll"
    else:
        path = "../testlib5/testso.so"
    try:
        if myos == "windows":
            ext = ctypes.CDLL(path)
        else:
            ext = ctypes.cdll.LoadLibrary(path)
        ext.primef.argtype = (ctypes.c_int,)
        ext.primef.restype = ctypes.POINTER(ctypes.c_int)
        ext.freef.argtype = (ctypes.POINTER(ctypes.c_int),)
        return ext
    except:
        return None

def main(dll):
    def gof():
        time.sleep(0.1)
        try:
            num = int( inp.get() )
            ptr = dll.primef(num)
            temp = [0] * ptr[0]
            for i in range( 0, len(temp) ):
                temp[i] = ptr[i + 1]
            dll.freef(ptr)
            out = str(temp)
        except Exception as e:
            out = str(e)
        strvar.set(out)
        
    # 메인 윈도우
    win = tkinter.Tk()
    win.title('testext5')
    win.geometry("400x250+200+100")
    win.resizable(False, False)

    # 입출력창
    strvar = tkinter.StringVar()
    if dll == None:
        strvar.set("작동에 testlib5가 필요합니다")
    else:
        strvar.set("testlib5 로드 성공")
    show = tkinter.Entry(win, font=("Consolas", 14), textvariable = strvar, width=20, state="readonly")
    show.pack()
    inp = tkinter.Entry(win, font=("Consolas", 14), width=20)
    inp.pack()
    but = tkinter.Button(win, font=("Consolas", 14), text="소인수분해", command=gof)
    but.pack()

    # 설명
    lbl = tkinter.Label(win, font=("Consolas", 14), text="ST5adv common ext / user ext\n테스트 프로그램입니다.\n32bit signed int를\n소인수분해합니다.")
    lbl.pack()

    win.mainloop()

dll = load()
main(dll)
time.sleep(0.5)
