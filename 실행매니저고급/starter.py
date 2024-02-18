# test616 : ST5adv main

import os
import shutil
import time

import clipboard
import webbrowser

import tkinter
import tkinter.ttk
import tkinter.messagebox
from tkinter import filedialog
from PIL import Image, ImageTk

import threading as thr
import multiprocessing as mp

import kdb
import kzip
import kweb
import ksign
import kaes

class procmgr: # 프로세스 실행 관리
    def __init__(self):
        self.current = 0 # 현재 실행중인 프로세스

    def file(self, path): # 파일 실행(작업 폴더 맞추기)
        if guiclass.myos == "windows":
            temp = os.path.abspath("./").replace("\\", "/")
            tgt = os.path.abspath(path).replace("\\", "/")
            if tgt[-1] == "/":
                tgt = tgt[0:-1]
            pos = tgt.rfind("/")
            os.chdir( tgt[0:pos] )
            os.startfile( tgt[pos + 1:].replace("/", "\\") )
            os.chdir(temp)
        else:
            temp = os.path.abspath("./").replace("\\", "/")
            tgt = os.path.abspath(path).replace("\\", "/")
            if tgt[-1] == "/":
                tgt = tgt[0:-1]
            pos = tgt.rfind("/")
            os.chdir( tgt[0:pos] )
            os.system(f"gnome-terminal -- ./{tgt[pos+1:]}")
            os.chdir(temp)

    def folder(self, path): # 폴더 열기
        if guiclass.myos == "windows":
            os.startfile( path.replace("/", "\\") )
        else:
            os.system(f"open {path}")

class reltool: # div/kzip 배포 도구
    def __init__(self):
        self.chunksize = 10485760 # div size
        self.exfolder = "./temp616" # 확장 export 출력 폴더

    def divpack(self, path): # path -> path.num
        path = os.path.abspath(path).replace("\\", "/")
        size = os.path.getsize(path)
        num = size // self.chunksize
        with open(path, "rb") as f:
            for i in range(0, num):
                with open(f"{path}.{i}", "wb") as t:
                    temp = f.read(self.chunksize)
                    t.write(temp)
            if size % self.chunksize != 0:
                with open(f"{path}.{num}", "wb") as t:
                    temp = f.read(size % self.chunksize)
                    t.write(temp)

    def divunpack(self, path): # path.0 -> path
        path = os.path.abspath(path).replace("\\", "/")[0:-2]
        num = 0
        while os.path.exists(f"{path}.{num}"):
            num = num + 1
        with open(path, "wb") as f:
            for i in range(0, num):
                with open(f"{path}.{i}", "rb") as t:
                    temp = t.read()
                    f.write(temp)

    def kzippack(self, path, tgt): # folder path -> kzip tgt
        path = path.replace("\\", "/")
        if path[-1] == "/":
            path = path[0:-1]
        if os.path.exists(f"{path}/st5adv_versioninfo.txt"):
            d0 = time.strftime( '%Y-%m-%d_%H:%M:%S', time.localtime( time.time() ) )
            tbox = kdb.toolbox()
            tbox.readfile(f"{path}/st5adv_versioninfo.txt")
            tbox.fixdata("release", d0)
            tbox.fixdata("os", guiclass.myos)
            tbox.writefile(f"{path}/st5adv_versioninfo.txt")
        tbox = kzip.toolbox()
        tbox.export = tgt
        tbox.folder = path
        tbox.abs()
        tbox.zipfolder("webp")

    def kzipunpack(self, path, tgt): # ext folder path, kzip tgt
        tbox = kzip.toolbox()
        tbox.export = "./temp616"
        tbox.unzip(tgt)
        temp = os.listdir("./temp616")[0]
        if temp[-1] == "/":
            temp = temp[0:-1]
        if os.path.exists(f"{path}/{temp}"):
            shutil.rmtree(f"{path}/{temp}")
        shutil.move(f"./temp616/{temp}", f"{path}/{temp}")
        if os.path.exists(f"{path}/{temp}/st5adv_versioninfo.txt"):
            d0 = time.strftime( '%Y-%m-%d_%H:%M:%S', time.localtime( time.time() ) )
            tbox = kdb.toolbox()
            tbox.readfile(f"{path}/{temp}/st5adv_versioninfo.txt")
            tbox.fixdata("install", d0)
            tbox.writefile(f"{path}/{temp}/st5adv_versioninfo.txt")

    def fclear(self, mknew): # temp616 폴더 초기화, mknew는 생성 여부
        if os.path.exists("./temp616"):
            shutil.rmtree("./temp616")
        if mknew:
            os.mkdir("./temp616")

    def sign(self, tgt, pvkpwb): # 파일 서명
        tbox0 = kdb.toolbox()
        tbox1 = ksign.toolbox()
        tbox2 = kaes.genbytes()

        with open(contclass.pvk, "rb") as f:
            temp = f.read()
        h, m, s = tbox2.view(temp)
        temp = tbox2.de(pvkpwb, kaes.genkf("NoSuchPath"), temp, s)
        tbox0.readstr( str(temp, encoding="utf-8") )

        d0 = tbox0.getdata("name")[3] # 서명이름
        d1 = tbox0.getdata("date")[3] # 생성일
        d2 = tbox0.getdata("strength")[3] # 강도
        d3 = tbox0.getdata("private")[3] # 개인키
        d4 = tbox0.getdata("public")[3] # 공개키

        d5 = "rsa.name = 0\nrsa.date = 0\nrsa.strength = 0\nrsa.public = 0\nsign.explain = 0\nsign.date = 0\nsign.hash = 0\nsign.enc = 0\n"
        d6 = tbox1.khash(tgt) # 해시값
        d7 = tbox1.fm(d0, d6) # 평문
        d8 = tbox1.sign(d3, d7) # 암호문
        d9 = time.strftime( '%Y-%m-%d_%H:%M:%S', time.localtime( time.time() ) ) # 서명일
        
        tbox0 = kdb.toolbox()
        tbox0.readstr(d5)
        tbox0.fixdata("rsa.name", d0)
        tbox0.fixdata("rsa.date", d1)
        tbox0.fixdata("rsa.strength", d2)
        tbox0.fixdata("rsa.public", d4)
        tbox0.fixdata("sign.explain", "st5adv_extension")
        tbox0.fixdata("sign.date", d9)
        tbox0.fixdata("sign.hash", d6)
        tbox0.fixdata("sign.enc", d8)

        tbox0.writefile(f"{tgt}.txt")

    def verify(self, tgt, txt): # 파일 서명 검증
        tbox0 = kdb.toolbox()
        tbox0.readfile("./st5adv_publickey.txt")
        d0 = tbox0.getdata("public")[3] # st5 공개키

        tbox0 = kdb.toolbox()
        tbox0.readfile(txt)
        d1 = tbox0.getdata("rsa.public")[3] # txt의 공개키
        d2 = tbox0.getdata("sign.hash")[3] # txt의 hash
        d3 = tbox0.getdata("sign.enc")[3] # txt의 암호문
        d4 = tbox0.getdata("rsa.name")[3] # 서명 이름

        if d0 != d1:
            return False

        tbox1 = ksign.toolbox()
        if tbox1.khash(tgt) != d2:
            return False

        return tbox1.verify( d1, d3, tbox1.fm(d4, d2) )

class contmgr: # extension 관리
    def __init__(self):
        self.name = [ ] # ext 이름
        self.exe = [ ] # ext 실행파일
        self.icon = [ ] # ext 아이콘
        self.txt = [ ] # ext 정보

        self.dosign = False # 서명 여부
        self.pvk = "" # 개인키 파일 경로
        self.cmurl = "" # CM url
        self.cmexpurl = "" # content explain

        if os.path.exists("./st5adv_publickey.txt"):
            self.public = kdb.toolbox()
            self.public.readfile("./st5adv_publickey.txt")
            self.dosign = True
        else:
            self.public = None # 공개키 정보
            self.dosign = False # 서명 여부

    def search(self): # 설치된 extension 검색
        self.name = [ ]
        self.exe = [ ]
        self.icon = [ ]
        self.txt = [ ]
        for i in os.listdir("./extension"):
            if os.path.isdir(f"./extension/{i}"):
                self.name.append(i)
        for i in self.name:
            flag0 = True
            flag1 = True
            flag2 = True
            for j in os.listdir(f"./extension/{i}"):
                if "st5adv_entrypoint" in j:
                    self.exe.append(f"./extension/{i}/{j}")
                    flag0 = False
                elif "st5adv_extensionicon" in j:
                    self.icon.append(f"./extension/{i}/{j}")
                    flag1 = False
                elif "st5adv_versioninfo" in j:
                    self.txt.append(f"./extension/{i}/{j}")
                    flag2 = False
            if flag0:
                self.exe.append(None)
            if flag1:
                self.icon.append("./st5adv_extensionfill.png")
            if flag2:
                self.txt.append(None)

    def getinfo(self): # extension 별 info
        self.info = [ ] # [이름, 버전, OS, 배포일, 설치일, 추가정보]
        for i in self.txt:
            if i == None:
                self.info.append( ["Unknown", 0, "Unknown", "Unknown", "Unknown", "Unknown"] )
            else:
                temp = ["Unknown", 0, "Unknown", "Unknown", "Unknown", "Unknown"]
                try:
                    tbox = kdb.toolbox()
                    tbox.readfile(i)
                    temp[0] = tbox.getdata("name")[3]
                    temp[1] = tbox.getdata("version")[3]
                    temp[2] = tbox.getdata("os")[3]
                    temp[3] = tbox.getdata("release")[3]
                    temp[4] = tbox.getdata("install")[3]
                    temp[5] = tbox.getdata("info")[3]
                except:
                    pass
                self.info.append(temp)

    def imext(self, path, sign): # import extension, .0 -> extension/~
        relclass.fclear(True)
        relclass.divunpack(path)
        if self.dosign:
            if not relclass.verify(path[0:-2], sign):
                raise Exception("invalid RSA sign")
        relclass.kzipunpack( "./extension", path[0:-2] )
        try:
            os.remove( path[0:-2] )
        except:
            time.sleep(0.5)
            os.remove( path[0:-2] )
        relclass.fclear(False)

    def exext(self, tgt, name, pvkpw): # export extension, tgt 폴더 패키징 -> name으로 exfolder에 내보내기
        temp = f"{relclass.exfolder}{name}"
        relclass.kzippack(tgt, temp)
        if self.dosign:
            relclass.sign( temp, bytes(pvkpw, encoding="utf-8") )
        relclass.divpack(temp)
        try:
            os.remove(temp)
        except:
            time.sleep(0.5)
            os.remove(temp)

    def getlist(self): # CM online 확장 목록 구하기
        if guiclass.myos == "windows":
            webid = "winext"
        else:
            webid = "linuxext"
        temp = kweb.gettxt(self.cmurl, webid)
        tbox = kdb.toolbox()
        tbox.readstr(temp)
        d0 = tbox.getdata("num")[3] # 확장 개수 num
        d1 = [ ] # 확장별 청크 수
        d2 = [ ] # 청크 이름
        d3 = [ ] # 확장 설명
        for i in range(0, d0):
            d1.append( tbox.getdata(f"{i}.num")[3] )
            d2.append( tbox.getdata(f"{i}.name")[3] )
            d3.append( tbox.getdata(f"{i}.txt")[3] )
        return d1, d2, d3

    def download(self, num, name): # 다운로드
        url = self.cmurl[0:self.cmurl.rfind("/") + 1]
        kweb.download(url, name, num, f"{relclass.exfolder}{name}")

class maingui: # 메인 gui 화면
    def __init__(self):
        self.current = 0 # 현재 창 번호
        self.pmax = 0 # 최대 창
        self.open = True # 파일 실행 T, 폴더 열기 F
        self.msg = None # 메세지
        self.mainwin = None # 메인 화면
        self.curcount = None # 현재 창 번호 표시

        # 6개 화면 표시
        self.show0 = None
        self.show1 = None
        self.show2 = None
        self.show3 = None
        self.show4 = None
        self.show5 = None
        self.photo0 = None
        self.photo1 = None
        self.photo2 = None
        self.photo3 = None
        self.photo4 = None
        self.photo5 = None
        self.but0 = None
        self.but1 = None
        self.but2 = None
        self.but3 = None
        self.but4 = None
        self.but5 = None
        
        self.myos = ""

    def restruct(self): # 6배수 맞추기
        self.name = contclass.name.copy()
        self.icon = contclass.icon.copy()
        if len(self.name) == 0:
            add = 6
        else:
            add = 6 - len(self.name) % 6
            if add == 6:
                add = 0
        for i in range(0, add):
            self.name.append("")
            self.icon.append("./st5adv_extensionfill.png")

    def inf0(self): # mf0
        self.msg.set("패키지 개시 파일 선택")
        self.mainwin.update()
        time.sleep(0.7)
        try:
            path0 = tkinter.filedialog.askopenfile(title='파일 선택', filetypes=(('starting files', '*.0'),)).name # .0 pkg file
        except:
            path0 = None
        self.msg.set("패키지 서명 파일 선택")
        self.mainwin.update()
        time.sleep(0.7)
        try:
            path1 = tkinter.filedialog.askopenfile(title='파일 선택', filetypes=(('sign files', '*.txt'),)).name # .txt sign
        except:
            path1 = None
        try:
            contclass.imext(path0, path1)
            res = "successfully installed"
        except Exception as e:
            res = str(e)
        tkinter.messagebox.showinfo("확장 설치 결과", f" {res} ")
        contclass.search()
        contclass.getinfo()
        self.inf7()

    def inf1(self): # mf1
        def gofunc():
            time.sleep(0.1)
            name = ent0.get()
            pvkpw = ent1.get()
            path = tkinter.filedialog.askdirectory(title='폴더 선택')
            try:
                contclass.exext(path, name, pvkpw)
                res = "successfully packed"
            except Exception as e:
                res = str(e)
            tkinter.messagebox.showinfo("확장 배포 결과", f" {res} ")
            subwin.destroy()
            self.inf7()

        tbox = kaes.genfile()
        a, b, c = tbox.view(contclass.pvk)
        subwin = tkinter.Toplevel(self.mainwin)
        subwin.title('ST5adv')
        subwin.geometry("500x250+200+100") # lxgp= 700x350+300+200
        subwin.resizable(False, False)

        lbl0 = tkinter.Label(subwin, font=("맑은 고딕", 14), text="패키지 이름")
        lbl0.place(x=5, y=5)
        lbl1 = tkinter.Label( subwin, font=("맑은 고딕", 14), text=str(a, encoding="utf-8") )
        lbl1.place(x=5, y=65) # lxgp= x=5 y=75
        lbl1 = tkinter.Label(subwin, font=("맑은 고딕", 14), text="서명 비밀번호")
        lbl1.place(x=5, y=125) # lxgp= x=5 y=145
        ent0 = tkinter.Entry(subwin, font=("맑은 고딕", 14), width=30) # lxgp= w=15
        ent0.place(x=155, y=5) # lxgp= x=255 y=5
        ent1 = tkinter.Entry(subwin, font=("맑은 고딕", 14), width=30, show="*") # lxgp= w=15
        ent1.place(x=155, y=125) # lxgp= x=255 y=145
        but = tkinter.Button(subwin, font=("맑은 고딕", 14), text="확장 폴더 선택", command=gofunc)
        but.place(x=170, y=185) # lxgp= x=250 y=215

    def inf2(self): # mf2
        contclass.search()
        contclass.getinfo()

        def updt(event): # 클릭시 정보 표시
            time.sleep(0.1)
            num = listbox.curselection()[0]
            temp = contclass.info[num]
            tf = lambda x, y: [x[y * i:y * i + y] for i in range(0, len(x) // y + 1)]
            tc = [ ]
            for i in tf(f"이름 : {temp[0]}", 18):
                tc.append(i)
            for i in tf(f"버전 : {temp[1]}", 18):
                tc.append(i)
            for i in tf(f"OS : {temp[2]}", 18):
                tc.append(i)
            for i in tf(f"배포일 : {temp[3]}", 16):
                tc.append(i)
            for i in tf(f"설치일 : {temp[4]}", 16):
                tc.append(i)
            for i in tf(f"추가정보 : {temp[5]}", 10):
                tc.append(i)
            lblvar.set( "\n".join(tc) )
            subwin.update()

        def regen(): # 표시 최신화
            contclass.search()
            contclass.getinfo()
            listbox.delete( 0, listbox.size() )
            for i in contclass.name:
                listbox.insert(listbox.size(), i)
            subwin.update()

        def delfunc(): # 선택 항목 삭제
            time.sleep(0.1)
            num = listbox.curselection()[0]
            if tkinter.messagebox.askokcancel("확장 삭제", f" 정말 확장 {self.name[num]}을 삭제하시겠습니까? "):
                shutil.rmtree( "./extension/" + self.name[num] )
            regen()

        def backfunc(): # 돌아가기
            time.sleep(0.1)
            contclass.search()
            contclass.getinfo()
            subwin.destroy()
            self.inf7()

        subwin = tkinter.Toplevel(self.mainwin)
        subwin.title('ST5adv')
        subwin.geometry("500x500+200+100") # lxgp= 1000x850+300+200
        subwin.resizable(False, False)

        lstfr = tkinter.Frame(subwin)
        lstfr.place(x=5, y=55) # lxgp= x=5 y=105
        listbox = tkinter.Listbox(lstfr, width=22,  height=18, font=("Consolas", 14), selectmode='extended') # lxgp= w=22 h=12
        listbox.pack(side="left", fill="y")
        scbar = tkinter.Scrollbar(lstfr, orient="vertical")
        scbar.config(command=listbox.yview)
        scbar.pack(side="right", fill="y")
        listbox.config(yscrollcommand=scbar.set) # ext 리스트

        but0 = tkinter.Button(subwin, font=("Consolas", 14), text="delete Extension", command=delfunc)
        but0.place(x=5, y=5)
        but1 = tkinter.Button(subwin, font=("Consolas", 14), text="back to Main", command=backfunc)
        but1.place(x=255, y=5) # lxgp= x=505 y=5
        lblvar = tkinter.StringVar()
        lbl = tkinter.Label(subwin, font=("Consolas", 14), textvariable=lblvar)
        lbl.place(x=255, y=55) # lxgp= x=505 y=105

        lblvar.set(f"이름 : Unknown\n버전 : 0\nOS : Unknown\n배포일 : Unknown\n설치일 : Unknown\n추가정보 : Unknown")
        listbox.bind('<ButtonRelease-1>', updt)
        regen()

    def inf3(self): # mf3
        webbrowser.open(contclass.cmexpurl)

    def inf4(self): # mf4
        cnum, cname, cinfo = contclass.getlist()

        def dwn():
            time.sleep(0.1)
            num = listbox.curselection()[0]
            try:
                contclass.download( cnum[num], cname[num] )
                res = "download success"
            except Exception as e:
                res = str(e)
            tkinter.messagebox.showinfo("다운로드 결과", f" {res} ")

        def updt(event):
            time.sleep(0.1)
            num = listbox.curselection()[0]
            tf = lambda x, y: [x[y * i:y * i + y] for i in range(0, len(x) // y + 1)]
            tc = ["Extension", "Package", ""]
            for i in tf(f"이름 : {cname[num]}", 18):
                tc.append(i)
            for i in tf(f"설명 : {cinfo[num]}", 10):
                tc.append(i)
            lblvar.set( "\n".join(tc) )
            subwin.update()

        subwin = tkinter.Toplevel(self.mainwin)
        subwin.title('ST5adv')
        subwin.geometry("500x500+200+100") # lxgp= 1000x850+300+200
        subwin.resizable(False, False)

        lstfr = tkinter.Frame(subwin)
        lstfr.place(x=5, y=55) # lxgp= x=5 y=105
        listbox = tkinter.Listbox(lstfr, width=22,  height=18, font=("Consolas", 14), selectmode='extended') # lxgp= w=22 h=12
        listbox.pack(side="left", fill="y")
        scbar = tkinter.Scrollbar(lstfr, orient="vertical")
        scbar.config(command=listbox.yview)
        scbar.pack(side="right", fill="y")
        listbox.config(yscrollcommand=scbar.set) # ext 리스트

        but0 = tkinter.Button(subwin, font=("Consolas", 14), text="download Extension", command=dwn)
        but0.place(x=5, y=5)
        lblvar = tkinter.StringVar()
        lbl = tkinter.Label(subwin, font=("Consolas", 14), textvariable=lblvar)
        lbl.place(x=255, y=55) # lxgp= x=505 y=105

        lblvar.set(f"Extension\nPackage\n\n이름 : Unknown\n설명 : Unknown")
        listbox.bind('<ButtonRelease-1>', updt)
        for i in cname:
            listbox.insert(listbox.size(), i)
        subwin.update()

    def inf5(self): # mf5
        def getfile():
            time.sleep(0.1)
            path = tkinter.filedialog.askopenfile(title='파일 선택', filetypes=(('all files', '*.*'),)).name
            copypath(path)

        def getfolder():
            time.sleep(0.1)
            path = tkinter.filedialog.askdirectory(title='폴더 선택')
            copypath(path)

        def copypath(path):
            path = os.path.abspath(path).replace("\\", "/")
            try:
                clipboard.copy(path)
            except:
                tkinter.messagebox.showinfo("클립보드 복사오류", f" 현재 시스템에서 복사-붙여넣기 할 수 없습니다. 다음 두 명령어 중 하나를 실행하십시오. \n sudo apt-get install xsel \n sudo apt-get install xclip ")
            entvar.set(path)
            subwin.update()

        subwin = tkinter.Toplevel(self.mainwin)
        subwin.title('ST5adv')
        subwin.geometry("250x150+200+100") # lxgp= 500x300+500+300
        subwin.resizable(False, False)

        but0 = tkinter.Button(subwin, font=("Consolas", 14), text="파일 선택", command=getfile)
        but0.place(x=5, y=5)
        but1 = tkinter.Button(subwin, font=("Consolas", 14), text="폴더 선택", command=getfolder)
        but1.place(x=125, y=5) # lxgp= x=255 y=5
        entvar = tkinter.StringVar()
        ent = tkinter.Entry(subwin, font=("Consolas", 14), textvariable=entvar, width=23, state="readonly")
        ent.place(x=5, y=85) # lxgp= x=5 y=205

    def inf6(self, num): # exef N
        num = 6 * self.current + num
        if num < len(contclass.name):
            if self.open:
                procclass.file( contclass.exe[num] )
            else:
                procclass.folder( "./extension/" + self.name[num] )

    def inf7(self): # regen
        self.restruct()
        self.show0.set( self.name[6 * self.current][0:15] )
        self.show1.set( self.name[6 * self.current + 1][0:15] )
        self.show2.set( self.name[6 * self.current + 2][0:15] )
        self.show3.set( self.name[6 * self.current + 3][0:15] )
        self.show4.set( self.name[6 * self.current + 4][0:15] )
        self.show5.set( self.name[6 * self.current + 5][0:15] )
        self.photo0 = tkinter.PhotoImage( file=self.icon[6 * self.current] )
        self.photo1 = tkinter.PhotoImage( file=self.icon[6 * self.current + 1] )
        self.photo2 = tkinter.PhotoImage( file=self.icon[6 * self.current + 2] )
        self.photo3 = tkinter.PhotoImage( file=self.icon[6 * self.current + 3] )
        self.photo4 = tkinter.PhotoImage( file=self.icon[6 * self.current + 4] )
        self.photo5 = tkinter.PhotoImage( file=self.icon[6 * self.current + 5] )
        self.but0.config(image=self.photo0)
        self.but1.config(image=self.photo1)
        self.but2.config(image=self.photo2)
        self.but3.config(image=self.photo3)
        self.but4.config(image=self.photo4)
        self.but5.config(image=self.photo5)
        self.pmax = len(self.name) // 6
        self.curcount.set(f"{self.current + 1} / {self.pmax}")
        if self.open:
            self.msg.set("mode : 파일 실행")
        else:
            self.msg.set("mode : 폴더 열기")
        self.mainwin.update()

    def main(self):
        
        def mf0(): # ext import
            time.sleep(0.1)
            self.inf0()

        def mf1(): # ext export
            time.sleep(0.1)
            self.inf1()

        def mf2(): # ext manage
            time.sleep(0.1)
            self.inf2()

        def mf3(): # ext web view
            time.sleep(0.1)
            self.inf3()

        def mf4(): # ext web download
            time.sleep(0.1)
            self.inf4()

        def mf5(): # selector
            time.sleep(0.1)
            self.inf5()

        def mf6(): # st5adv folder
            time.sleep(0.1)
            procclass.folder("./")

        def goleft(): # 창 왼쪽으로
            time.sleep(0.1)
            if self.current != 0:
                self.current = self.current - 1
                self.inf7()

        def goright(): # 창 오른쪽으로
            time.sleep(0.1)
            if self.current + 1 < self.pmax:
                self.current = self.current + 1
                self.inf7()

        def toggle(): # 실행/폴더열기 토글
            time.sleep(0.1)
            if self.open:
                self.open = False
                self.msg.set("mode : 폴더 열기")
            else:
                self.open = True
                self.msg.set("mode : 파일 실행")
            self.mainwin.update()

        def exef0(): # execute 0
            time.sleep(0.1)
            self.inf6(0)

        def exef1(): # execute 1
            time.sleep(0.1)
            self.inf6(1)

        def exef2(): # execute 2
            time.sleep(0.1)
            self.inf6(2)

        def exef3(): # execute 3
            time.sleep(0.1)
            self.inf6(3)

        def exef4(): # execute 4
            time.sleep(0.1)
            self.inf6(4)

        def exef5(): # execute 5
            time.sleep(0.1)
            self.inf6(5)
        
        self.mainwin = tkinter.Tk()
        self.mainwin.title("ST5adv")
        self.mainwin.geometry("600x500+100+50") # lxgp= 800x600+300+200

        self.pmax = len(self.name) // 6
        self.msg = tkinter.StringVar()
        self.msg.set("mode : 파일 실행")
        self.curcount = tkinter.StringVar()
        self.show0 = tkinter.StringVar()
        self.show1 = tkinter.StringVar()
        self.show2 = tkinter.StringVar()
        self.show3 = tkinter.StringVar()
        self.show4 = tkinter.StringVar()
        self.show5 = tkinter.StringVar()

        mbar = tkinter.Menu(self.mainwin) # 메뉴 바

        menu0 = tkinter.Menu(mbar, tearoff=0) # local
        menu0.add_command(label="확장 가져오기", font=("맑은 고딕", 14), command=mf0)
        menu0.add_command(label="확장 내보내기", font=("맑은 고딕", 14), command=mf1)
        menu0.add_command(label="확장 관리", font = ("맑은 고딕", 14), command=mf2)
        mbar.add_cascade(label="  로컬 확장 도구  ", menu=menu0)

        menu1 = tkinter.Menu(mbar, tearoff=0) # web
        menu1.add_command(label="웹 설명서 보기", font=("맑은 고딕", 14), command=mf3)
        menu1.add_command(label="웹 확장 다운로드", font=("맑은 고딕", 14), command=mf4)
        mbar.add_cascade(label="  웹 확장 도구  ", menu=menu1)

        menu2 = tkinter.Menu(mbar, tearoff=0) # tool
        menu2.add_command(label="파일/폴더 선택기", font=("맑은 고딕", 14), command=mf5)
        menu2.add_command(label="ST5adv 폴더 열기", font=("맑은 고딕", 14), command=mf6)
        mbar.add_cascade(label="  ST5adv 도구  ", menu=menu2)

        menu3 = tkinter.Menu(mbar, tearoff=0) # toggle
        menu3.add_command(label="파일/폴더 모드 토글", font=("맑은 고딕", 14), command=toggle)
        mbar.add_cascade(label="  실행모드 설정  ", menu=menu3)

        self.mainwin.config(menu=mbar)

        show = tkinter.Entry(self.mainwin, font=("맑은 고딕", 14), textvariable=self.msg, width=58, state="readonly") # lxgp= w=35
        show.place(x=5, y=5)
        gol = tkinter.Button(self.mainwin, text="\n\n\n\n\n\n\n < \n\n\n\n\n\n\n", font=("맑은 고딕", 14), command=goleft) # lxgp= \n3
        gol.place(x=5, y=55) # lxgp= x=5 y=105
        gor = tkinter.Button(self.mainwin, text="\n\n\n\n\n\n\n > \n\n\n\n\n\n\n", font=("맑은 고딕", 14), command=goright) # lxgp= \n3
        gor.place(x=550, y=55) # lxgp= x=730 y=105

        self.but0 = tkinter.Button(self.mainwin, image=self.photo0, command=exef0)
        self.but0.place(x=55, y=55) # lxgp= x=105 y=105
        self.lbl0 = tkinter.Label(self.mainwin, font=("맑은 고딕", 14), textvariable=self.show0)
        self.lbl0.place(x=55, y=215) # lxgp= x=105 y=265

        self.but1 = tkinter.Button(self.mainwin, image=self.photo1, command=exef1)
        self.but1.place(x=220, y=55) # lxgp= x=305 y=105
        self.lbl1 = tkinter.Label(self.mainwin, font=("맑은 고딕", 14), textvariable=self.show1)
        self.lbl1.place(x=220, y=215) # lxgp= x=305 y=265

        self.but2 = tkinter.Button(self.mainwin, image=self.photo2, command=exef2)
        self.but2.place(x=385, y=55) # lxgp= x=505 y=105
        self.lbl2 = tkinter.Label(self.mainwin, font=("맑은 고딕", 14), textvariable=self.show2)
        self.lbl2.place(x=385, y=215) # lxgp= x=505 y=265

        self.but3 = tkinter.Button(self.mainwin, image=self.photo3, command=exef3)
        self.but3.place(x=55, y=255) # lxgp= x=105 y=315
        self.lbl3 = tkinter.Label(self.mainwin, font=("맑은 고딕", 14), textvariable=self.show3)
        self.lbl3.place(x=55, y=415) # lxgp= x=105 y=475

        self.but4 = tkinter.Button(self.mainwin, image=self.photo4, command=exef4)
        self.but4.place(x=220, y=255) # lxgp= x=305 y=315
        self.lbl4 = tkinter.Label(self.mainwin, font=("맑은 고딕", 14), textvariable=self.show4)
        self.lbl4.place(x=220, y=415) # lxgp= x=305 y=475

        self.but5 = tkinter.Button(self.mainwin, image=self.photo5, command=exef5)
        self.but5.place(x=385, y=255) # lxgp= x=505 y=315
        self.lbl5 = tkinter.Label(self.mainwin, font=("맑은 고딕", 14), textvariable=self.show5)
        self.lbl5.place(x=385, y=415) # lxgp= x=505 y=475

        self.curlbl = tkinter.Label(self.mainwin, font=("맑은 고딕", 14), textvariable=self.curcount)
        self.curlbl.place(x=265, y=460) # lxgp= x=335 y=535

        self.inf7()
        self.mainwin.mainloop()

def initsys(): # 시작 시 초기화
    global procclass
    procclass = procmgr()
    global relclass
    relclass = reltool()
    global contclass
    contclass = contmgr()
    global guiclass
    guiclass = maingui()

    kdbtbox = kdb.toolbox()
    kdbtbox.readfile("st5adv_settings.txt")
    relclass.chunksize = kdbtbox.getdata("chunksize")[3]
    relclass.exfolder = kdbtbox.getdata("export")[3]
    contclass.dosign = kdbtbox.getdata("dosign")[3]
    contclass.pvk = kdbtbox.getdata("pvk")[3]
    contclass.cmurl = kdbtbox.getdata("cmurl")[3]
    contclass.cmexpurl = kdbtbox.getdata("webhelp")[3]
    guiclass.myos = kdbtbox.getdata("os")[3]

    relclass.exfolder = os.path.abspath(relclass.exfolder).replace("\\", "/")
    if relclass.exfolder[-1] != "/":
        relclass.exfolder = relclass.exfolder + "/"
    relclass.fclear(True)

    contclass.search()
    contclass.getinfo()
    guiclass.restruct()

if __name__ == "__main__":
    mp.freeze_support()
    initsys()
    guiclass.main()
    while procclass.current != 0:
        time.sleep(1)
    time.sleep(0.5)
