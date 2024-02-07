# test602 : viewer

import os
import time
import clipboard

import ctypes
import hashlib

import tkinter
import tkinter.messagebox
import tkinter.filedialog

import kdb

# little endian encoding
def encode(num, length):
    temp = [0] * length
    for i in range(0, length):
        temp[i] = num % 256
        num = num // 256
    return bytes(temp)

# little endian decoding
def decode(data):
    temp = 0
    for i in range( 0, len(data) ):
        if data[i] != 0:
            temp = temp + data[i] * 256 ** i
    return temp

class mainclass:
    def __init__(self): # 초기화
        self.valid = False # settings.txt 여부
        self.viewsize = True # 크기 보기?
        self.start = "./" # 시작점
        self.tpname = [ ] # 바로가기 이름들
        self.tppath = [ ] # 바로가기 경로들
        self.divsize = 10485760 # div 크기

        try:
            tbox = kdb.toolbox()
            tbox.readfile("./settings.txt")
            self.viewsize = tbox.getdata("viewsize")[3]
            self.start = os.path.abspath( tbox.getdata("start")[3] ).replace("\\", "/")
            num = tbox.getdata("num")[3]
            for i in range(0, num):
                self.tpname.append( tbox.getdata(f"{i}.name")[3] )
                self.tppath.append( os.path.abspath( tbox.getdata(f"{i}.path")[3] ).replace("\\", "/") )
            self.divsize = tbox.getdata("divsize")[3]
            self.valid = True
        except:
            self.start = os.path.abspath(self.start).replace("\\", "/")
            la = len(self.tpname)
            lb = len(self.tppath)
            if la != lb:
                num = min(la, lb)
                self.tpname = self.tpname[0:num]
                self.tppath = self.tppath[0:num]

        if self.start[-1] != "/":
            self.start = self.start + "/"
        self.current = self.start # 현재 폴더

        myos = "windows" # "windows"/"linux"
        if myos == "windows":
            dllname = "./kviewer5.dll"
            dllhash = b'x\x074u\xe8\xc4X\xa3\xb8C\xf87Ad\xd4\xc2\xda\x8ab\x01\xa8\xcf\x8d\x85\xa7\xa4\xae\x18\xfd)n7\xca\xfa\x94\x7fLb\x80;R/4\xaf\xab\xe6>\x10n\x9a \x88\x93\xcc\xcd\x9be\xfe+\xcb\x03\xb2\xeb\x06'
        else:
            dllname = "./kviewer5.so"
            dllhash = b""
        with open(dllname, "rb") as f:
            temp = f.read()
        temp = hashlib.sha3_512(temp).digest()
        #print(temp)
        if temp != dllhash:
            raise Exception("wrong ffi")
        if myos == "windows":
            self.ext = ctypes.CDLL(dllname)
        else:
            self.ext = ctypes.cdll.LoadLibrary(dllname)
        self.ext.ex0.argtype = (ctypes.POINTER(ctypes.c_char), ctypes.c_int)
        self.ext.ex0.restype = ctypes.POINTER(ctypes.c_char)
        self.ext.ex1.argtype = (ctypes.POINTER(ctypes.c_char), ctypes.c_int)
        self.ext.ex1.restype = ctypes.POINTER(ctypes.c_char)
        self.ext.ex2.argtype = (ctypes.POINTER(ctypes.c_char), ctypes.c_int, ctypes.c_int)
        self.ext.ex2.restype = ctypes.POINTER(ctypes.c_char)
        self.ext.ex3.argtype = (ctypes.POINTER(ctypes.c_char), ctypes.c_int, ctypes.c_int, ctypes.c_int)
        self.ext.ex3.restype = ctypes.POINTER(ctypes.c_char)
        self.ext.ex4.argtype = (ctypes.POINTER(ctypes.c_char),)

    def convsize(self, size): # B/KiB/MiB/GiB로 변환
        if size < 1024:
            return f"{size} B"
        elif size < 1048576:
            size = size / 1024
            return f"{size:.1f} KiB"
        elif size < 1073741824:
            size = size / 1048576
            return f"{size:.1f} MiB"
        else:
            size = size / 1073741824
            return f"{size:.1f} GiB"

    def fstruct(self, path, show, struct): # 절대경로 받아 formated size 보기
        a, b, c, d = self.getinfo(path)
        c = self.convsize(c)
        l = 0
        for i in show:
            if ord(i) > 256:
                l = l + 2
            else:
                l = l + 1
        if l + 12 > struct:
            return show + "    " + c
        else:
            return show + " " * (struct - 8 - l) + c

    def getinfo(self, path): # 폴더/파일 정보 -> 이하 폴더수 N, 이하 파일수 N, 크기 N, 수정시각 S
        if os.path.isdir(path):
            path = bytes(path, encoding="utf-8")
            dl = len(path)
            dp = ctypes.c_char_p(path)
            dr = self.ext.ex0(dp, dl)
            temp = [0] * 24
            for i in range(0, 24):
                temp[i] = dr[i][0]
            self.ext.ex4(dr)
            out = [0, 0, 0, 0]
            out[0] = decode( temp[0:8] )
            out[1] = decode( temp[8:16] )
            out[2] = decode( temp[16:24] )
        else:
            out = [0, 1, 0, 0]
            out[2] = os.path.getsize(path)
        out[3] = time.strftime( '%Y-%m-%d_%H:%M:%S', time.localtime( os.path.getmtime(path) ) )
        return out[0], out[1], out[2], out[3]

    def khashf(self, path): # 폴더/파일 khash -> hashv 64B -> resb S (len 128)
        path = bytes(path, encoding="utf-8")
        dl = len(path)
        dp = ctypes.c_char_p(path)
        dr = self.ext.ex1(dp, dl)
        temp = [0] * 64
        for i in range(0, 64):
            temp[i] = dr[i][0]
        self.ext.ex4(dr)
        res = [""] * 128
        symbols = ["0","1","2","3","4","5","6","7","8","9","a","b","c","d","e","f"]
        for i in range(0, 64):
            res[2 * i] = symbols[temp[i] // 16]
            res[2 * i + 1] = symbols[temp[i] % 16]
        return "".join(res)

    def kzipf(self, mode, path): # kzip 제어 -> errmsg S
        if mode:
            mode = 0
        else:
            mode = 1
        path = bytes(path, encoding="utf-8")
        dl = len(path)
        dp = ctypes.c_char_p(path)
        dr = self.ext.ex2(dp, dl, mode)
        temp = [0] * 4
        for i in range(0, 4):
            temp[i] = dr[i][0]
        tsize = decode(temp)
        temp = [0] * tsize
        for i in range(0, tsize):
            temp[i] = dr[i + 4][0]
        self.ext.ex4(dr)
        return str(bytes(temp), encoding="utf-8")

    def divf(self, mode, path): # 파일 단순 나눔 제어 -> errmsg S
        if mode:
            mode = 0
        else:
            mode = 1
        path = bytes(path, encoding="utf-8")
        dl = len(path)
        dp = ctypes.c_char_p(path)
        dr = self.ext.ex3(dp, dl, mode, self.divsize)
        temp = [0] * 4
        for i in range(0, 4):
            temp[i] = dr[i][0]
        tsize = decode(temp)
        temp = [0] * tsize
        for i in range(0, tsize):
            temp[i] = dr[i + 4][0]
        self.ext.ex4(dr)
        return str(bytes(temp), encoding="utf-8")

    def startf(self, path, mode): # 폴더/파일 실행
        # windows/linux difference
        if mode == 0:
            os.startfile(path) # 폴더 열기
            #os.system(f"open {path}")
        else:
            os.startfile(path) # 파일 열기
            #os.system(f"open {path}")

    def mainfunc(self): # GUI
        def mf0(): # 바로가기 선택창 열기
            def gof(event):
                time.sleep(0.1)
                nonlocal subpast
                temp = sublist.curselection()[0]
                if temp != subpast:
                    subpast = temp
                else:
                    self.current = self.tppath[temp]
                    subwin.destroy()
                if self.current[-1] != "/":
                    self.current = self.current + "/"
                regen()
            
            # 하위 윈도우, 리스트박스
            time.sleep(0.1)
            subwin = tkinter.Toplevel(win)
            subwin.title('kviewer5')
            subwin.geometry("250x300+200+100") # lxgp= 330x430+200+100
            subwin.resizable(False, False)
            sublist = tkinter.Listbox(subwin, font = ("맑은 고딕", 14), width=23,  height=11) # lxgp= w=14 h=7
            sublist.place(x=5, y=5)
            for i in self.tpname:
                sublist.insert(sublist.size(), i)
            subpast = -1
            sublist.bind('<ButtonRelease-1>', gof)

        def mf1(): # 폴더 수동 이동
            def gof():
                time.sleep(0.1)
                temp = subent.get()
                if temp[0] == '"' and temp[-1] == '"':
                    temp = temp[1:-1]
                temp = os.path.abspath(temp).replace("\\", "/")
                if os.path.exists(temp):
                    if os.path.isdir(temp):
                        self.current = temp
                        subwin.destroy()
                if self.current[-1] != "/":
                    self.current = self.current + "/"
                regen()

            # 하위 윈도우, 문구, 입력창, 버튼
            time.sleep(0.1)
            subwin = tkinter.Toplevel(win)
            subwin.title('kviewer5')
            subwin.geometry("300x150+200+100") # lxgp= 450x200+200+100
            subwin.resizable(False, False)
            sublbl = tkinter.Label(subwin, font=("맑은 고딕", 12), text="이동할 폴더 경로를 입력하세요")
            sublbl.place(x=5, y=5)
            subent = tkinter.Entry(subwin, font=("맑은 고딕", 14), width=20) # lxgp= w=14
            subent.place(x=5, y=90) # lxgp= x=5 y=90
            subbut = tkinter.Button(subwin, font=("맑은 고딕", 14), text="이동", command=gof) # lxgp= font=14
            subbut.place(x=240, y=85) # lxgp= x=330 y=85
        
        def mf2(): # 현재 폴더 열기
            time.sleep(0.1)
            self.startf(self.current, 0)

        def mf3(): # 경로 복사
            time.sleep(0.1)
            temp = listbox.curselection()[0]
            tocopy = self.current + curshow[temp]
            try:
                clipboard.copy(tocopy)
            except:
                tkinter.messagebox.showinfo("클립보드 복사오류", " 현재 시스템에서 복사-붙여넣기 할 수 없습니다. 다음 두 명령어 중 하나를 실행하십시오. \n sudo apt-get install xsel \n sudo apt-get install xclip ")
                    

        def mf4(): # 정보 보기
            time.sleep(0.1)
            temp = listbox.curselection()[0]
            numd, numf, nums, ledt = self.getinfo( self.current + curshow[temp] )
            tkinter.messagebox.showinfo("선택 항목 정보", f" 포함 폴더 수 : {numd}, 포함 파일 수 : {numf} \n 크기 : {self.convsize(nums)}, 상세크기 : {nums} B \n 수정시각 : {ledt} ")

        def mf5(): # kzip pack
            time.sleep(0.1)
            tpath = tkinter.filedialog.askdirectory(title='폴더 선택')
            msg = self.kzipf(True, tpath)
            tkinter.messagebox.showinfo("kzip pack 결과", f" {msg} ")
            regen()

        def mf6(): # kzip unpack
            time.sleep(0.1)
            tpath = tkinter.filedialog.askopenfile(title='파일 선택', filetypes=(('webp files', '*.webp'),('png files', '*.png'),('all files', '*.*'))).name
            msg = self.kzipf(False, tpath)
            tkinter.messagebox.showinfo("kzip unpack 결과", f" {msg} ")
            regen()

        def mf7(): # div pack
            time.sleep(0.1)
            tpath = tkinter.filedialog.askopenfile(title='파일 선택').name
            msg = self.divf(True, tpath)
            tkinter.messagebox.showinfo("div pack 결과", f" {msg} ")
            regen()

        def mf8(): # div unpack
            time.sleep(0.1)
            tpath = tkinter.filedialog.askopenfile(title='파일 선택', filetypes=(('starting files', '*.0'),)).name
            msg = self.divf(False, tpath)
            tkinter.messagebox.showinfo("div unpack 결과", f" {msg} ")
            regen()

        def mf9(): # khash
            time.sleep(0.1)
            temp = listbox.curselection()[0]
            resb = self.khashf( self.current + curshow[temp] )
            tkinter.messagebox.showinfo("KHASH 결과", f" {resb[0:32]} \n {resb[32:64]} \n {resb[64:96]} \n {resb[96:128]} ")

        def mfx(): # toggle viewsize
            time.sleep(0.1)
            self.viewsize = not self.viewsize
            regen()

        def clickf(event): # 클릭 시 하위 항목으로
            time.sleep(0.1)
            nonlocal past
            temp = listbox.curselection()[0]
            if temp == past:
                if curshow[temp][-1] == "/":
                    self.current = self.current + curshow[temp]
                    regen()
                else:
                    self.startf(self.current + curshow[temp], 1)
            else:
                past = temp

        def refunc(): # 상위폴더로
            time.sleep(0.1)
            if self.current[-1] == "/":
                self.current = self.current[0:-1]
            if self.current.count("/") != 0:
                pos = self.current.rfind("/")
                self.current = self.current[0:pos]
            self.current = self.current + "/"
            regen()

        def regen(): # 상태 업데이트
            strvar.set(self.current)
            nonlocal listbox
            nonlocal curshow
            folder = [ ]
            file = [ ]
            curshow = [ ]
            for i in os.listdir(self.current):
                i = i.replace("\\", "/")
                if i[-1] == "/":
                    i = i[0:-1]
                if os.path.isdir(self.current + i):
                    folder.append(i + "/")
                else:
                    file.append(i)
            listbox.delete( 0, listbox.size() )
            for i in folder:
                if self.viewsize:
                    listbox.insert( listbox.size(), self.fstruct(self.current + i, i, 50) ) # lxgp= 40
                else:
                    listbox.insert(listbox.size(), i)
                curshow.append(i)
            for i in file:
                if self.viewsize:
                    listbox.insert( listbox.size(), self.fstruct(self.current + i, i, 50) ) # lxgp= 40
                else:
                    listbox.insert(listbox.size(), i)
                curshow.append(i)
            win.update()

        # 메인 윈도우
        win = tkinter.Tk()
        win.title('kviewer5')
        win.geometry("600x450+200+100") # lxgp= 850x650+200+100
        win.resizable(False, False)

        mbar = tkinter.Menu(win) # 메뉴 바

        menu0 = tkinter.Menu(mbar, tearoff=0) # teleport
        menu0.add_command(label="바로 가기", font = ("맑은 고딕", 14), command=mf0)
        menu0.add_command(label="수동 지정", font = ("맑은 고딕", 14), command=mf1)
        mbar.add_cascade(label="  Teleport  ", menu=menu0)

        menu1 = tkinter.Menu(mbar, tearoff=0) # general
        menu1.add_command(label="현재 폴더 열기", font = ("맑은 고딕", 14), command=mf2)
        menu1.add_command(label="선택 경로 복사", font = ("맑은 고딕", 14), command=mf3)
        menu1.add_separator()
        menu1.add_command(label="세부 정보 확인", font = ("맑은 고딕", 14), command=mf4)
        mbar.add_cascade(label="  General  ", menu=menu1)

        menu2 = tkinter.Menu(mbar, tearoff=0) # files
        menu2.add_command(label="KZIP pack", font = ("맑은 고딕", 14), command=mf5)
        menu2.add_command(label="KZIP unpack", font = ("맑은 고딕", 14), command=mf6)
        menu2.add_separator()
        menu2.add_command(label="DIV pack", font = ("맑은 고딕", 14), command=mf7)
        menu2.add_command(label="DIV unpack", font = ("맑은 고딕", 14), command=mf8)
        menu2.add_separator()
        menu2.add_command(label="KHASH", font = ("맑은 고딕", 14), command=mf9)
        mbar.add_cascade(label="  Files  ", menu=menu2)

        menu3 = tkinter.Menu(mbar, tearoff=0) # viewsize
        menu3.add_command(label="Toggle", font = ("맑은 고딕", 14), command=mfx)
        mbar.add_cascade(label="  Viewsize  ", menu=menu3)

        win.config(menu=mbar)

        but = tkinter.Button(win, font=("맑은 고딕", 12), text=" < ", command=refunc)
        but.place(x=5, y=5) # 상위폴더로 버튼
        strvar = tkinter.StringVar()
        strvar.set(self.current)
        ent = tkinter.Entry(win, textvariable=strvar, font=("Consolas", 14), width=50, state="readonly") # lxgp= w=40
        ent.place(x=65, y=10) # 현재 폴더 표시 lxgp= x=75 y=10

        lstfr = tkinter.Frame(win)
        lstfr.place(x=5, y=55) # lxgp= x=5 y=75
        listbox = tkinter.Listbox(lstfr, width=55,  height=16, font = ("Consolas", 14), selectmode = 'extended') # lxgp= w=43 h=10
        listbox.pack(side="left", fill="y")
        scbar = tkinter.Scrollbar(lstfr, orient="vertical")
        scbar.config(command=listbox.yview)
        scbar.pack(side="right", fill="y")
        listbox.config(yscrollcommand=scbar.set) # 파일 리스트

        past = -1 # 직전 선택
        curshow = [ ] # 현재 보여지는 리스트
        listbox.bind('<ButtonRelease-1>', clickf)
        regen()

        win.mainloop()

classloader = mainclass()
classloader.mainfunc()
time.sleep(0.5)
