# test716 : extension.testext

import os
import time
import ctypes
import tkinter

import kobj
import kdb

class mainclass:
    def __init__(self):
        worker = kdb.toolbox()
        with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
            worker.read( f.read() )
        self.iswin = worker.get("dev.os")[3] == "windows"
        if os.path.exists("../../_ST5_COMMON/testlib/"):
            self.dll = ctypes.CDLL("../../_ST5_COMMON/testlib/testlib.dll") if self.iswin else ctypes.cdll.LoadLibrary("../../_ST5_COMMON/testlib/testlib.so")
            self.dll.func0.argtypes, self.dll.func0.restype = kobj.call("b", "") # free
            self.dll.func1.argtypes, self.dll.func1.restype = kobj.call("b", "b") # primef
        else:
            self.dll = None
        if self.iswin:
            self.parms = {"win.size" : "400x180+200+100", "lbl.font" : ("Consolas", 14), "ent.len" : 30}
        else:
            self.parms = {"win.size" : "600x400+300+150", "lbl.font" : ("Consolas", 14), "ent.len" : 30}

    def guiloop(self):
        def primef():
            time.sleep(0.1)
            if self.dll == None:
                strvar.set("common.testlib required")
            else:
                t, _ = kobj.send( kobj.encode(int( ent1.get() ), 8) )
                t = self.dll.func1(t)
                strvar.set( str(kobj.recvauto(t), encoding="utf-8") )
                self.dll.func0(t)
        
        win = tkinter.Tk()
        win.title('test extension')
        win.geometry( self.parms["win.size"] )
        win.resizable(False, False)
        lbl = tkinter.Label(win, font=self.parms["lbl.font"], text="\n테스트 프로그램입니다.\n소인수분해를 할 수 있습니다.")
        lbl.pack()
        strvar = tkinter.StringVar()
        strvar.set("program ready")
        ent0 = tkinter.Entry(win, font=self.parms["lbl.font"], textvariable = strvar, width=self.parms["ent.len"], state="readonly")
        ent0.pack()
        ent1 = tkinter.Entry(win, font=self.parms["lbl.font"], width=self.parms["ent.len"])
        ent1.pack()
        but = tkinter.Button(win, font=self.parms["lbl.font"], text=" Calculate ", command=primef)
        but.pack()
        win.mainloop()

kobj.repath()
worker = mainclass()
worker.guiloop()
time.sleep(0.5)
