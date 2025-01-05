# test745 : independent.installer user

import os
import shutil
import time
import threading

import webbrowser
from win32com.client import Dispatch # pkg pywin32

import tkinter as tk
import tkinter.ttk as tkt
import tkinter.filedialog as tkf
import tkinter.messagebox as tkm

import ksc
import kdb
import kobj
import kcom
import kpkg

class downloader:
    def __init__(self):
        self.mwin, self.font0, self.font1, self.bg, self.act = None, ("맑은 고딕", 12), ("Consolas", 14), "gray90", "lawn green"
        self.desktop, self.home, self.startup, self.devmode, self.com_pkg, self.ind_pkg, self.check = "", "", "", False, [ ], [ ], [ ]

    def setup(self):
        def select(n):
            time.sleep(0.1)
            if n == 0:
                self.desktop = abspath( tkf.askdirectory(title="폴더 선택") )
                svar1a.set(self.desktop)
            elif n == 1:
                self.home = abspath( tkf.askdirectory(title="폴더 선택") )
                svar1b.set(self.home)
            elif n == 2:
                self.startup = abspath( tkf.askdirectory(title="폴더 선택") )
                svar1c.set(self.startup)

        def switch():
            time.sleep(0.1)
            self.devmode = not self.devmode
            t = self.act if self.devmode else self.bg
            but2a.configure(bg=t)

        def submit():
            time.sleep(0.1)
            if os.path.isdir(self.desktop) and os.path.isdir(self.home) and os.path.isdir(self.startup):
                self.mwin.destroy()
                self.checkup()
            else:
                tkm.showerror(title="경로 설정 오류", message=" 경로가 잘못되었습니다. \n 바탕화면, 홈 폴더, 시작 프로그램 폴더 경로를 설정하세요. ")

        self.mwin = tk.Tk()
        self.mwin.title("다운로더 - 경로 설정")
        self.mwin.geometry("500x240+400+200")
        self.mwin.resizable(False, False)

        lbl0 = tk.Label(self.mwin, font=self.font0, text="RealUse5 윈도우 사용자용 설치기입니다.\n바탕화면, 홈 폴더, 시작 프로그램 폴더 경로를 설정하세요.")
        lbl0.pack(padx=5, pady=5)
        self.home = abspath(os.path.expanduser("~"))
        self.desktop, self.startup = self.home + "Desktop/", self.home + "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup/"
        fr1 = tk.Frame(self.mwin)
        fr1.pack(fill="x", padx=5, pady=5)

        but1a = tk.Button(fr1, font=self.font0, text=" . . . ", command=lambda:select(0))
        but1a.grid(row=0, column=0)
        svar1a = tk.StringVar()
        svar1a.set(self.desktop)
        ent1a = tk.Entry(fr1, font=self.font1, textvariable=svar1a, width=42, state="readonly")
        ent1a.grid(row=0, column=1, padx=5)
        but1b = tk.Button(fr1, font=self.font0, text=" . . . ", command=lambda:select(1))
        but1b.grid(row=1, column=0)
        svar1b = tk.StringVar()
        svar1b.set(self.home)
        ent1b = tk.Entry(fr1, font=self.font1, textvariable=svar1b, width=42, state="readonly")
        ent1b.grid(row=1, column=1, padx=5)
        but1c = tk.Button(fr1, font=self.font0, text=" . . . ", command=lambda:select(2))
        but1c.grid(row=2, column=0)
        svar1c = tk.StringVar()
        svar1c.set(self.startup)
        ent1c = tk.Entry(fr1, font=self.font1, textvariable=svar1c, width=42, state="readonly")
        ent1c.grid(row=2, column=1, padx=5)

        fr2 = tk.Frame(self.mwin)
        fr2.pack(fill="x", padx=5, pady=5)
        but2a = tk.Button(fr2, font=self.font0, text=" 개발자 모드 ", command=switch)
        but2a.pack(side="left")
        but2b = tk.Button(fr2, font=self.font0, text="   Help   ", command=openhelp)
        but2b.pack(side="left")
        but2c = tk.Button(fr2, font=self.font0, text="   Next   ", command=submit)
        but2c.pack(side="right")
        self.mwin.mainloop()

    def checkup(self):
        try:
            self.com_pkg, self.ind_pkg = get_pkg(self.devmode)
            self.check = [True, True, True, False, True] if self.devmode else [True, False, False, False, True]
        except Exception as e:
            tkm.showerror(title="정보 가져오기 실패", message=f" {e} ")

        def switch(n):
            time.sleep(0.1)
            self.check[n] = not self.check[n]
            t = self.act if self.check[n] else self.bg
            but[n].configure(bg=t)

        def submit():
            time.sleep(0.1)
            if self.check[4] and not self.check[0]:
                tkm.showerror(title="라이브러리 설치 경고", message=" common 라이브러리는 Starter5를 설치한 뒤, \n 추가로 다운로드 받아야 합니다. ")
            else:
                self.mwin.destroy()
                self.working()
        
        self.mwin = tk.Tk()
        self.mwin.title("다운로더 - 구성 설정")
        self.mwin.geometry("500x300+400+200")
        self.mwin.resizable(False, False)

        lbl0 = tk.Label(self.mwin, font=self.font0, text="설치할 프로그램을 설정하세요.\n초록색으로 활성화된 패키지만 설치됩니다.")
        lbl0.pack(padx=5, pady=5)
        fr1 = tk.Frame(self.mwin)
        fr1.pack(fill="x", padx=5, pady=5)
        but = [0, 0, 0, 0, 0]

        but[0] = tk.Button(fr1, font=self.font0, text=" 설치 ", command=lambda:switch(0))
        but[0].grid(row=0, column=0)
        svar0 = tk.StringVar()
        svar0.set(f"{self.ind_pkg[0][0]} : {self.ind_pkg[0][1]}")
        ent0 = tk.Entry(fr1, font=self.font1, textvariable=svar0, width=42, state="readonly")
        ent0.grid(row=0, column=1, padx=5)
        but[1] = tk.Button(fr1, font=self.font0, text=" 설치 ", command=lambda:switch(1))
        but[1].grid(row=1, column=0)
        svar1 = tk.StringVar()
        svar1.set(f"{self.ind_pkg[1][0]} : {self.ind_pkg[1][1]}")
        ent1 = tk.Entry(fr1, font=self.font1, textvariable=svar1, width=42, state="readonly")
        ent1.grid(row=1, column=1, padx=5)
        but[2] = tk.Button(fr1, font=self.font0, text=" 설치 ", command=lambda:switch(2))
        but[2].grid(row=2, column=0)
        svar2 = tk.StringVar()
        svar2.set(f"{self.ind_pkg[2][0]} : {self.ind_pkg[2][1]}")
        ent2 = tk.Entry(fr1, font=self.font1, textvariable=svar2, width=42, state="readonly")
        ent2.grid(row=2, column=1, padx=5)
        but[3] = tk.Button(fr1, font=self.font0, text=" 설치 ", command=lambda:switch(3))
        but[3].grid(row=3, column=0)
        svar3 = tk.StringVar()
        svar3.set(f"{self.ind_pkg[3][0]} : {self.ind_pkg[3][1]}")
        ent3 = tk.Entry(fr1, font=self.font1, textvariable=svar3, width=42, state="readonly")
        ent3.grid(row=3, column=1, padx=5)

        but[4] = tk.Button(fr1, font=self.font0, text=" 설치 ", command=lambda:switch(4))
        but[4].grid(row=4, column=0)
        svar4 = tk.StringVar()
        svar4.set("common : Starter5 공용 라이브러리")
        ent4 = tk.Entry(fr1, font=self.font1, textvariable=svar4, width=42, state="readonly")
        ent4.grid(row=4, column=1, padx=5)
        for i in range(0, 5):
            t = self.act if self.check[i] else self.bg
            but[i].configure(bg=t)

        fr2 = tk.Frame(self.mwin)
        fr2.pack(fill="x", padx=5, pady=5)
        but2a = tk.Button(fr2, font=self.font0, text="   Help   ", command=openhelp)
        but2a.pack(side="left")
        but2b = tk.Button(fr2, font=self.font0, text="   Next   ", command=submit)
        but2b.pack(side="right")
        self.mwin.mainloop()

    def working(self):
        def write(text):
            lst.insert(lst.size(), text)
            lst.see(lst.size() - 1)
            self.mwin.update()
            time.sleep(0.1)

        def addproc():
            nonlocal curnum
            curnum = curnum + 1
            proc["value"] = curnum
            self.mwin.update()
            time.sleep(0.1)

        def install(name, mode):
            wk, mode = kpkg.toolbox(), mode.split("/")
            wk.osnum = 1
            if mode[0] == "standalone":
                path = abspath( wk.unpack("./temp.webp") )
                for i in os.listdir(path):
                    if i != "_ST5_VERSION.txt":
                        shutil.move(path + i, self.desktop + i)
            elif mode[0] == "startup":
                path = abspath( wk.unpack("./temp.webp") )
                for i in os.listdir(path):
                    if i != "_ST5_VERSION.txt":
                        shutil.move(path + i, self.startup + i)
            elif mode[0] == "install":
                path = abspath( wk.unpack("./temp.webp") )
                shutil.move(path, self.home + name)
                path = abspath(self.home + name)
                for i in os.listdir(path):
                    if mode[1] in i:
                        mkshort(path + i, self.desktop + name + ".lnk")
            else:
                path = abspath( wk.unpack("./temp.webp") )
                shutil.move(path, self.desktop + name)
            try:
                os.remove("./temp.webp")
                shutil.rmtree("./temp674/")
            except:
                time.sleep(0.1)
            return ksc.crc32hash( bytes(wk.public, encoding="utf-8") ).hex()

        def addcfg():
            wk, path = kdb.toolbox(), self.home + "Starter5/_ST5_CONFIG.txt"
            with open(path, "r", encoding="utf-8") as f:
                wk.read( f.read() )
            wk.fix("path.export", self.home + "Starter5/_ST5_DATA/")
            wk.fix("path.desktop", self.desktop)
            wk.fix("path.local", self.home)
            wk.fix("url.download", url_dwn0)
            wk.fix("url.info", url_dwn0+"list.html")
            wk.fix("url.help", url_help)
            wk.fix("dev.os", "windows")
            wk.fix("dev.activate", self.devmode)
            with open(path, "w", encoding="utf-8") as f:
                f.write( wk.write() )

        def unpack(name):
            wk = kpkg.toolbox()
            wk.osnum = 1
            path = wk.unpack("./temp.webp")
            shutil.move(path, self.home + "Starter5/_ST5_COMMON/" + name)
            try:
                os.remove("./temp.webp")
                shutil.rmtree("./temp674/")
            except:
                time.sleep(0.1)
            return ksc.crc32hash( bytes(wk.public, encoding="utf-8") ).hex()

        maxnum, curnum = 0, 1
        for i in range(0, 4):
            if self.check[i]:
                maxnum = maxnum + 2
        if self.check[3]:
            maxnum = maxnum + 1
        if self.check[4]:
            maxnum = maxnum + 2 * len(self.com_pkg)

        self.mwin = tk.Tk()
        self.mwin.title("다운로더 - 설치 진행중")
        self.mwin.geometry("500x300+400+200")
        self.mwin.resizable(False, False)

        lbl = tk.Label(self.mwin, font=self.font0, text="다운로드와 설치가 진행 중입니다.\n설치가 완료될 때까지 창을 닫지 마세요.")
        lbl.pack(padx=5, pady=5)
        fr = tk.Frame(self.mwin)
        fr.pack(padx=5, pady=5)
        lst = tk.Listbox(fr, width=42,  height=7, font=self.font1)
        lst.pack(side="left", fill="y")
        scr = tk.Scrollbar(fr, orient="vertical")
        scr.config(command=lst.yview)
        scr.pack(side="right", fill="y")
        lst.config(yscrollcommand=scr.set)
        proc = tkt.Progressbar(self.mwin, length=470, maximum=maxnum+1)
        proc.pack(padx=15, pady=15)

        for i in range(0, 4):
            if self.check[i]:
                write(f"다운로드 시작 : {self.ind_pkg[i][0]}")
                try:
                    if i == 0:
                        kcom.download( url_dwn0, self.ind_pkg[i][0], self.ind_pkg[i][2], "./temp.webp", [0] )
                    else:
                        kcom.download( url_dwn1, self.ind_pkg[i][0], self.ind_pkg[i][2], "./temp.webp", [0] )
                    addproc()
                    write("다운로드 완료")
                    try:
                        write(f"설치 시작 : {self.ind_pkg[i][0]}")
                        ph = install( self.ind_pkg[i][0], self.ind_pkg[i][3] )
                        addproc()
                        write(f"설치 완료 ({ph})")
                    except Exception as e:
                        addproc()
                        write(f"설치 실패 : {e}")
                except Exception as e:
                    addproc()
                    write(f"다운로드 실패 : {e}")

        if self.check[0]:
            try:
                addcfg()
                addproc()
                write("개인설정 초기화 완료")
            except Exception as e:
                addproc()
                write(f"개인설정 초기화 실패 : {e}")

        if self.check[4]:
            for i in self.com_pkg:
                write(f"다운로드 시작 : {i[0]}")
                try:
                    kcom.download( url_dwn0, i[0], i[2], "./temp.webp", [0] )
                    addproc()
                    write("다운로드 완료")
                    try:
                        write(f"설치 시작 : {i[0]}")
                        ph = unpack( i[0] )
                        addproc()
                        write(f"설치 완료 ({ph})")
                    except Exception as e:
                        addproc()
                        write(f"설치 실패 : {e}")
                except Exception as e:
                    addproc()
                    write(f"다운로드 실패 : {e}")

        tkm.showinfo(title="다운로더 - 설치 완료", message=" 설치를 완료했습니다. \n 창을 닫으셔도 됩니다. ")
        self.mwin.mainloop()

def get_pkg(devmode):
    wk0, wk1, num0, num1, out0, out1 = kdb.toolbox(), kdb.toolbox(), 0, 0, [ ], [ ]
    wk0.read( kcom.gettxt(url_dwn0+"list.html", "common") )
    wk1.read( kcom.gettxt(url_dwn0+"list.html", "independent") )
    while f"{num0}.name" in wk0.name:
        num0 = num0 + 1
    while f"{num1}.name" in wk1.name:
        num1 = num1 + 1
    for i in range(0, num0):
        if devmode or not wk0.get(f"{i}.devonly")[3]:
            out0.append( ( wk0.get(f"{i}.name")[3], wk0.get(f"{i}.text")[3], wk0.get(f"{i}.num")[3] ) )
    for i in range(0, num1):
        out1.append( ( wk1.get(f"{i}.name")[3], wk1.get(f"{i}.text")[3], wk1.get(f"{i}.num")[3], wk1.get(f"{i}.type")[3] ) )

    wk0, wk1, num0, num1 = kdb.toolbox(), kdb.toolbox(), 0, 0
    wk0.read( kcom.gettxt(url_dwn1+"list.html", "common") )
    wk1.read( kcom.gettxt(url_dwn1+"list.html", "independent") )
    while f"{num0}.name" in wk0.name:
        num0 = num0 + 1
    while f"{num1}.name" in wk1.name:
        num1 = num1 + 1
    for i in range(0, num0):
        if devmode or not wk0.get(f"{i}.devonly")[3]:
            out0.append( ( wk0.get(f"{i}.name")[3], wk0.get(f"{i}.text")[3], wk0.get(f"{i}.num")[3] ) )
    for i in range(0, num1):
        out1.append( ( wk1.get(f"{i}.name")[3], wk1.get(f"{i}.text")[3], wk1.get(f"{i}.num")[3], wk1.get(f"{i}.type")[3] ) )
    return out0, out1

def abspath(path):
    path = os.path.abspath(path).replace("\\", "/")
    if os.path.isdir(path) and path[-1] != "/":
        path = path + "/"
    return path

def mkshort(src, dst):
    shell = Dispatch('WScript.Shell')
    shortcut = shell.CreateShortCut( dst.replace("/", "\\") )
    shortcut.Targetpath = src.replace("/", "\\")
    shortcut.WorkingDirectory = src[:src.rfind("/")].replace("/", "\\")
    shortcut.save()

def openhelp():
    time.sleep(0.1)
    webbrowser.open(url_help)

if __name__ == "__main__":
    kobj.repath() # !!! base info section !!!
    url_dwn0, url_dwn1 = "https://taewook427.github.io/RealUse5_Sub0/data/", "https://taewook427.github.io/RealUse5_Sub1/data/"
    url_help = "https://taewook427.github.io/RealUse5/websvr/helpurl/helpurl.html"
    k = downloader()
    try:
        if not os.path.exists("./ksign5hy.dll"):
            raise Exception("DLL 파일을 로드하지 못했습니다\n모든 파일의 압축을 풀어주세요")
        k.setup()
    except Exception as e:
        tkm.showerror(title="Critical Error", message=f" {e} ")
    time.sleep(0.5)
