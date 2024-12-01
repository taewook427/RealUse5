# test720 : extension.clock5

import time

import tkinter
import tkinter.ttk
import tkinter.messagebox

import kdb

class mainclass:
    def __init__(self, iswin):
        self.iswork, self.view = False, 0 # is backend working / 0, 1, 2 view
        self.istimer, self.timersec, self.timerstart = False, 600.0, 0.0 # is timer counting / timer left / timer start
        self.isstop, self.stopsec, self.stopstart = False, 0.0, 0.0 # is stopwatch counting / stopwatch passed / stopwatch start
        self.dat, self.parms = ['월요일', '화요일', '수요일', '목요일', '금요일', '토요일', '일요일'], dict()
        if iswin:
            self.parms["win.size"], self.parms["win.w"], self.parms["win.h"] = "370x200+300+200", 350, 160
            self.parms["lbl0b.y"], self.parms["lbl1b.y"], self.parms["ent1c.x"], self.parms["lock.y"] = 55, 65, 135, 115
            self.parms["but.y0"], self.parms["but.y1"], self.parms["lst2b.font"] = 100, 95, ("맑은 고딕", 14)
            self.parms["a"], self.parms["b"], self.parms["c"], self.parms["d"], self.parms["e"], self.parms["f"] = 65, 125, 60, 185, 245, 295
        else:
            self.parms["win.size"], self.parms["win.w"], self.parms["win.h"] = "650x340+300+200", 640, 280
            self.parms["lbl0b.y"], self.parms["lbl1b.y"], self.parms["ent1c.x"], self.parms["lock.y"] = 125, 125, 230, 205
            self.parms["but.y0"], self.parms["but.y1"], self.parms["lst2b.font"] = 205, 205, ("맑은 고딕", 12)
            self.parms["a"], self.parms["b"], self.parms["c"], self.parms["d"], self.parms["e"], self.parms["f"] = 95, 185, 125, 320, 410, 500

    def gettime(self): # get time/date
        temp = time.time()
        now = time.localtime(temp)
        hour = int( time.strftime('%H' , now) )
        ap = "오전" if hour < 12 else "오후"
        hour = (hour - 1) % 12 + 1
        msec = str( temp - int(temp) )[2]
        return f"{ap} {hour:0>2}" + time.strftime(':%M:%S', now) + f".{msec}", time.strftime('%Y년 %m월 %d일 ', now) + self.dat[now.tm_wday]
    
    def tpcom(self, tf): # time float -> structed str
        ti, tf = int(tf), tf - int(tf)
        h, ti = ti // 3600, ti % 3600
        m, ti = ti // 60, ti % 60
        s = f"{ti + tf:.1f}"
        return f"{h:0>2}:{m:0>2}:{s:0>4}"
    
    def guiloop(self): # working behind
        time.sleep(0.1)
        while self.iswork:
            if self.view == 0:
                t0, t1 = self.gettime()
                self.txt0a.set(t0)
                self.txt0b.set(t1)
            elif self.view == 2:
                if self.isstop:
                    self.txt2a.set( self.tpcom(self.stopsec + time.time() - self.stopstart) )
            if self.istimer:
                left = self.timersec - (time.time() - self.timerstart)
                if left > 0:
                    self.txt1a.set( self.tpcom(left) )
                else:
                    self.istimer, self.timersec = False, 0.0
                    self.txt1a.set("00:00:00.0")
                    tkinter.messagebox.showinfo("타이머 종료", " 설정된 시간이 모두 지났습니다! ")
            self.win.update()
            time.sleep(0.1) # update 10hz

    def entry(self): # main gui
        self.win = tkinter.Tk()
        self.win.title("Clock5")
        self.win.geometry( self.parms["win.size"] )
        self.win.resizable(False, False)

        notebook = tkinter.ttk.Notebook( self.win, width=self.parms["win.w"], height=self.parms["win.h"] )
        notebook.place(x=5, y=5)
        fr0, fr1, fr2 = tkinter.Frame(self.win), tkinter.Frame(self.win), tkinter.Frame(self.win)
        notebook.add(fr0, text="    clock    ")
        notebook.add(fr1, text="    timer    ")
        notebook.add(fr2, text="  stopwatch  ")

        def clickevent(event):
            self.view = notebook.index( notebook.select() )
        notebook.bind("<<NotebookTabChanged>>", clickevent)

        self.txt0a = tkinter.StringVar()
        self.txt0a.set("오전 12:00:00.0")
        self.txt0b = tkinter.StringVar()
        self.txt0b.set("1970년 01월 01일 목요일")
        self.label0a = tkinter.Label(fr0, font=("Consolas", 30), textvariable=self.txt0a)
        self.label0a.place(x=5, y=5)
        self.label0b = tkinter.Label(fr0, font=("맑은 고딕", 15), textvariable=self.txt0b)
        self.label0b.place( x=5, y=self.parms["lbl0b.y"] )

        def lockf():
            time.sleep(0.1)
            if lockv.get() == 0:
                self.win.wm_attributes("-topmost", 0)
            else:
                self.win.wm_attributes("-topmost", 1)
        lockv = tkinter.IntVar()
        lockb = tkinter.Checkbutton(fr0, text="Always On Display", font=("Consolas", 14), variable=lockv, command=lockf)
        lockb.place( x=5, y=self.parms["lock.y"] )

        self.txt1a = tkinter.StringVar()
        self.txt1a.set("00:10:00.0")
        self.label1a = tkinter.Label(fr1, font=("Consolas", 30), textvariable=self.txt1a)
        self.label1a.place(x=5, y=5)
        label1b = tkinter.Label(fr1, font=("Consolas", 15), text="HH:MM:SS.s")
        label1b.place( x=5, y=self.parms["lbl1b.y"] )
        ent1c = tkinter.Entry(fr1, font=("Consolas", 15), width=12)
        ent1c.place( x=self.parms["ent1c.x"], y=self.parms["lbl1b.y"] )

        def setf():
            time.sleep(0.1)
            temp = ent1c.get().replace(" ", ":")
            if temp == "":
                temp = "0"
            temp = temp.split(":")
            if len(temp) == 0:
                temp = ["0", "0", "0"]
            elif len(temp) == 1:
                temp = [ "0", "0", temp[0] ]
            elif len(temp) == 2:
                temp = [ "0", temp[0], temp[1] ]
            else:
                temp = temp[0:3]
            h, m, s = int( temp[0] ), int( temp[1] ), float( temp[2] )
            self.timersec = 3600 * h + 60 * m + s
            if m >= 60 or s >= 60.0:
                tkinter.messagebox.showinfo("잘못된 시간 입력", " 분과 초는 60 이하로만 입력 가능합니다. ")
            else:
                self.txt1a.set( self.tpcom(self.timersec) )
        but1d = tkinter.Button(fr1, text=' ✔ ', font=("맑은 고딕", 15), command=setf)
        but1d.place( x=5, y=self.parms["but.y0"] )

        def togf():
            self.istimer = not self.istimer
            if self.txt1e.get() == ' ▶ ':
                self.txt1e.set(' | | ')
                self.timerstart = time.time()
            else:
                self.txt1e.set(' ▶ ')
                self.timersec = self.timersec - (time.time() - self.timerstart)
            time.sleep(0.1)
        self.txt1e = tkinter.StringVar()
        self.txt1e.set(' ▶ ')
        but1e = tkinter.Button(fr1, textvariable=self.txt1e, font=("맑은 고딕", 15), command=togf)
        but1e.place( x=self.parms["a"], y=self.parms["but.y0"] )

        def stopf():
            self.istimer = False
            time.sleep(0.2)
            self.txt1e.set(' ▶ ')
            setf()
        but1f = tkinter.Button(fr1, text=' ■ ', font=("맑은 고딕", 15), command=stopf)
        but1f.place( x=self.parms["b"], y=self.parms["but.y0"] )

        self.txt2a = tkinter.StringVar()
        self.txt2a.set("00:00:00.0")
        self.label2a = tkinter.Label(fr2, font=("Consolas", 30), textvariable=self.txt2a)
        self.label2a.place(x=5, y=5)
        list2b = tkinter.Listbox(fr2, width=17,  font=self.parms["lst2b.font"], height=3)
        list2b.place( x=5, y=self.parms["c"] )
        list2b.insert(1,'< 기록 >')

        def togglef():
            self.isstop = not self.isstop
            if self.txt2c.get() == ' ▶ ':
                self.txt2c.set(' | | ')
                self.stopstart = time.time()
            else:
                self.txt2c.set(' ▶ ')
                self.stopsec = self.stopsec + time.time() - self.stopstart
            time.sleep(0.1)
        self.txt2c = tkinter.StringVar()
        self.txt2c.set(' ▶ ')
        but2c = tkinter.Button(fr2, textvariable=self.txt2c, font=("맑은 고딕", 15), command=togglef)
        but2c.place( x=self.parms["d"], y=self.parms["but.y1"] )

        def writef():
            temp = list2b.size()
            list2b.insert(temp, f"기록{temp} : {self.txt2a.get()}")
            list2b.see(temp + 1)
            time.sleep(0.1)
        but2d = tkinter.Button(fr2, text=' ≡ ', font=("맑은 고딕", 15), command=writef)
        but2d.place( x=self.parms["e"], y=self.parms["but.y1"] )

        def resetf():
            self.isstop = False
            time.sleep(0.1)
            self.txt2c.set(' ▶ ')
            self.txt2a.set("00:00:00.0")
            self.stopsec = 0.0
            list2b.delete( 1, list2b.size() )
        but2e = tkinter.Button(fr2, text=' ■ ', font=("맑은 고딕", 15), command=resetf)
        but2e.place( x=self.parms["f"], y=self.parms["but.y1"] )

        def shutdown():
            self.iswork = False
            time.sleep(0.2)
            self.win.destroy()
        self.win.protocol('WM_DELETE_WINDOW', shutdown)
        self.iswork = True

worker = kdb.toolbox()
with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
    worker.read( f.read() )
iswin = worker.get("dev.os")[3] == "windows"
worker = mainclass(iswin)
worker.entry()
worker.guiloop()
time.sleep(0.5)