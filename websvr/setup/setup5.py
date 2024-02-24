# test617 : gen5 AutoSetup (win)

import time
import os
import shutil

import zlib
import webbrowser

import tkinter
import tkinter.messagebox
import tkinter.filedialog
import tkinter.ttk

import winshell
from win32com.client import Dispatch

import kweb
import kzip
import kdb

class mainclass:

    def __init__(self):
        self.name = [ ] # 표시 이름
        self.info = [ ] # 표시 정보
        self.chunkname = [ ] # 파일 이름
        self.chunknum = [ ] # 청크 수
        self.chunkorder = [ ] # 설치 명령

        # 다운로드 정보
        self.url = "https://taewook427.github.io/project_gen5/websvr/datasvr0.html"
        self.bin = "binary"
        self.val = "value"

        # 도움말 정보
        self.help0 = "https://taewook427.github.io/project_gen5/websvr/userguide/test593.html"
        self.help1 = "https://taewook427.github.io/project_gen5/websvr/userguide/test594.html"
        self.help2 = "https://taewook427.github.io/project_gen5/websvr/userguide/test595.html"

        self.msg = "" # 결과 메세지

        self.fp0 = "" # 바탕화면
        self.fp1 = "" # 홈 폴더
        self.fp2 = "" # 시작 프로그램 폴더
        self.install = [ ] # 설치 여부
        self.status = 0 # 0 초기, 1 경로 설정, 2 설치 여부, 3 실제 설치, 4 완료

        self.num = 0 # 설치 번호
        self.flag = 0 # 설치된 숫자

    def dclear(self, mknew): # temp617 폴더 초기화, mknew는 생성 여부
        if os.path.exists("./temp617"):
            shutil.rmtree("./temp617")
        if mknew:
            os.mkdir("./temp617")

    def fclear(self, path): # path 파일 삭제
        if os.path.exists(path):
            os.remove(path)

    def hash(self): # ./data.webp CRC32 value
        with open("./data.webp", "rb") as f:
            data = f.read()
        temp = hex( zlib.crc32(data) )
        temp = temp[2:]
        temp = "0" * ( 8 - len(temp) ) + temp
        return temp

    def getlist(self): # 다운로드 리스트 구하기
        try:
            temp = kweb.gettxt(self.url, self.val)
            tbox = kdb.toolbox()
            tbox.readstr(temp)
            num = tbox.getdata("setup5.num")[3]
            for i in range(0, num):
                self.name.append( tbox.getdata(f"setup5.{i}.name")[3] )
                self.info.append( tbox.getdata(f"setup5.{i}.info")[3] )
            
            temp = kweb.gettxt(self.url, self.bin)
            tbox = kdb.toolbox()
            tbox.readstr(temp)
            for i in range(0, num):
                self.chunkname.append( tbox.getdata(f"setup5.{i}.name")[3] )
                self.chunknum.append( tbox.getdata(f"setup5.{i}.num")[3] )
                self.chunkorder.append( tbox.getdata(f"setup5.{i}.order")[3] )

            self.msg = "온라인 연결 성공"
        except Exception as e:
            self.msg = str(e)
            self.name = [ ]
            self.info = [ ]
            self.chunkname = [ ]
            self.chunknum = [ ]
            self.chunkorder = [ ]

        self.fp1 = os.path.abspath( os.path.expanduser('~') ).replace("\\", "/")
        if self.fp1[-1] != "/":
            self.fp1 = self.fp1 + "/"
        self.fp0 = self.fp1 + "Desktop/"
        self.fp2 = self.fp1 + "/AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup/"
        for i in self.chunkorder:
            if i[0] == "y":
                self.install.append(True)
            else:
                self.install.append(False)

    def download(self, num): # 패키지 다운로드
        self.fclear("./data.webp")
        try:
            kweb.download(self.url[0:self.url.rfind("/") + 1], self.chunkname[num], self.chunknum[num], "./data.webp")
            res = "성공"
        except Exception as e:
            res = str(e)
        return res

    def unpack(self): # 패키지 풀기
        self.dclear(True)
        tbox = kzip.toolbox()
        tbox.export = "./temp617"
        tbox.unzip("./data.webp")

    def move(self, num): # 패키지 설치
        temp = self.chunkorder[num]
        if temp[1] == "d":
            isfolder = True
        else:
            isfolder = False
        if temp[2] == "d":
            unpackpath = self.fp0
        elif temp[2] == "h":
            unpackpath = self.fp1
        else:
            unpackpath = self.fp2
        istp = [0, 0, 0]
        if temp[3] == "y" and isfolder:
            istp[0] = True
            td = temp.split("/")
            istp[1] = td[1]
            istp[2] = td[2]
        else:
            istp[0] = False
        names = os.listdir("./temp617")
        for i in range( 0, len(names) ):
            if names[i][-1] == "/":
                names[i] = names[i][0:-1]
        for i in names:
            shutil.move(f"./temp617/{i}", f"{unpackpath}{i}")
        if istp[0]:
            self.mktp(f"{unpackpath}{names[0]}/{istp[1]}", f"{self.fp0}{istp[2]}")
        self.dclear(False)
        self.fclear("./data.webp")

    def mktp(self, tgt, path): # 바로가기 생성
        tgt = os.path.abspath(tgt).replace("/", "\\")
        path = os.path.abspath(path).replace("/", "\\")
        shell = Dispatch('WScript.Shell')
        shortcut = shell.CreateShortCut(path)
        shortcut.Targetpath = tgt
        shortcut.WorkingDirectory = tgt[ 0:tgt.rfind("\\") ]
        shortcut.save()

    def guif1(self):
        def bf2():
            time.sleep(0.1)
            webbrowser.open(self.help0)
        def bf3():
            time.sleep(0.1)
            webbrowser.open(self.help1)
        def bf4():
            time.sleep(0.1)
            webbrowser.open(self.help2)
        def bf5():
            time.sleep(0.1)
            webbrowser.open(self.url)

        def bf6():
            time.sleep(0.1)
            temp = tkinter.filedialog.askdirectory(title="바탕화면 설정")
            if temp != None:
                temp = os.path.abspath(temp).replace("\\", "/")
                if temp[-1] != "/":
                    temp = temp + "/"
                self.fp0 = temp
                var6.set(temp)
                win.update()
        def bf7():
            time.sleep(0.1)
            temp = tkinter.filedialog.askdirectory(title="홈 폴더 설정")
            if temp != None:
                temp = os.path.abspath(temp).replace("\\", "/")
                if temp[-1] != "/":
                    temp = temp + "/"
                self.fp1 = temp
                var7.set(temp)
                win.update()
        def bf8():
            time.sleep(0.1)
            temp = tkinter.filedialog.askdirectory(title="시작 프로그램 폴더")
            if temp != None:
                temp = os.path.abspath(temp).replace("\\", "/")
                if temp[-1] != "/":
                    temp = temp + "/"
                self.fp2 = temp
                var8.set(temp)
                win.update()

        def gof():
            time.sleep(0.1)
            if os.path.exists(self.fp0) and os.path.exists(self.fp1) and os.path.exists(self.fp2):
                self.status = 1
                win.destroy()
                time.sleep(0.5)
            else:
                tkinter.messagebox.showerror("경로 설정 실패", " 바탕화면, 홈 폴더, 시작 프로그램 폴더 \n 세 폴더의 경로 중 잘못된 경로가 존재합니다. ")

        win = tkinter.Tk()
        win.title("Setup5")
        win.geometry("600x320+100+50")
        win.configure(bg=cl0)

        lbl0 = tkinter.Label(win, font=("Consolas", 14), text="STEP 1 : 기본 경로 설정")
        lbl0.place(x=5, y=5)
        lbl0.configure(bg=cl0, fg=cl1)
        lbl1 = tkinter.Label(win, font=("Consolas", 14), text="웹 도움말 버튼")
        lbl1.place(x=455, y=5)
        lbl1.configure(bg=cl0, fg=cl1)
        but2 = tkinter.Button(win, font=("Consolas", 14), text="독립 프로그램", command=bf2)
        but2.place(x=5, y=45)
        but2.configure(bg=cl2, fg=cl1)
        but3 = tkinter.Button(win, font=("Consolas", 14), text="실행 매니저", command=bf3)
        but3.place(x=155, y=45)
        but3.configure(bg=cl2, fg=cl1)
        but4 = tkinter.Button(win, font=("Consolas", 14), text="확장 프로그램", command=bf4)
        but4.place(x=305, y=45)
        but4.configure(bg=cl2, fg=cl1)
        but5 = tkinter.Button(win, font=("Consolas", 14), text="데이터 서버", command=bf5)
        but5.place(x=455, y=45)
        but5.configure(bg=cl2, fg=cl1)

        lbl6 = tkinter.Label(win, font=("Consolas", 14), text="바탕 화면")
        lbl6.place(x=5, y=105)
        lbl6.configure(bg=cl0, fg=cl1)
        but6 = tkinter.Button(win, font=("Consolas", 14), text="...", command=bf6)
        but6.place(x=155, y=105)
        but6.configure(bg=cl2, fg=cl1)
        var6 = tkinter.StringVar()
        var6.set(self.fp0)
        ent6 = tkinter.Entry(win, font=("Consolas", 14), textvariable=var6, width=37, state="readonly")
        ent6.place(x=210, y=110)
        ent6.configure(readonlybackground=cl2, fg=cl1)

        lbl7 = tkinter.Label(win, font=("Consolas", 14), text="홈 폴더")
        lbl7.place(x=5, y=155)
        lbl7.configure(bg=cl0, fg=cl1)
        but7 = tkinter.Button(win, font=("Consolas", 14), text="...", command=bf7)
        but7.place(x=155, y=155)
        but7.configure(bg=cl2, fg=cl1)
        var7 = tkinter.StringVar()
        var7.set(self.fp1)
        ent7 = tkinter.Entry(win, font=("Consolas", 14), textvariable=var7, width=37, state="readonly")
        ent7.place(x=210, y=160)
        ent7.configure(readonlybackground=cl2, fg=cl1)

        lbl8 = tkinter.Label(win, font=("Consolas", 14), text="시작 프로그램")
        lbl8.place(x=5, y=205)
        lbl8.configure(bg=cl0, fg=cl1)
        but8 = tkinter.Button(win, font=("Consolas", 14), text="...", command=bf8)
        but8.place(x=155, y=205)
        but8.configure(bg=cl2, fg=cl1)
        var8 = tkinter.StringVar()
        var8.set(self.fp2)
        ent8 = tkinter.Entry(win, font=("Consolas", 14), textvariable=var8, width=37, state="readonly")
        ent8.place(x=210, y=210)
        ent8.configure(readonlybackground=cl2, fg=cl1)

        var9 = tkinter.StringVar()
        var9.set(self.msg)
        ent9 = tkinter.Entry(win, font=("Consolas", 14), textvariable=var9, width=47, state="readonly")
        ent9.place(x=105, y=280)
        ent9.configure(readonlybackground=cl2, fg=cl1)
        gobut = tkinter.Button(win, font=("Consolas", 14), text="진행하기", command=gof)
        gobut.place(x=5, y=275)
        gobut.configure(bg=cl2, fg=cl1)

        win.mainloop()

    def guif2(self):
        def bf3():
            time.sleep(0.1)
            temp = listbox.curselection()[0]
            self.install[temp] = not self.install[temp]
            draw()

        def bf4():
            time.sleep(0.1)
            self.status = 2
            win.destroy()

        def draw():
            temp = listbox.curselection()[0]
            var2.set(f"정보 :\n{self.info[temp]}\n설치 : {self.install[temp]}\n\nchunkname :\n{self.chunkname[temp]}\nchunknum : {self.chunknum[temp]}\nchunkorder : {self.chunkorder[temp][0:4]}")
            win.update()

        def clickf(event):
            time.sleep(0.1)
            draw()

        win = tkinter.Tk()
        win.title("Setup5")
        win.geometry("500x350+100+50")
        win.configure(bg=cl0)

        lbl0 = tkinter.Label(win, font=("Consolas", 14), text="STEP 2 : 설치 여부 설정")
        lbl0.place(x=5, y=5)
        lbl0.configure(bg=cl0, fg=cl1)
        listbox = tkinter.Listbox(win, width=22,  height=10, font=("Consolas", 14), selectmode='extended')
        listbox.place(x=5, y=55)
        listbox.configure(bg=cl2, fg=cl1, selectbackground=cl1, selectforeground=cl2)
        var2 = tkinter.StringVar()
        var2.set("정보 : 없음\n설치 : False\n\nchunkname : None\nchunknum : 0\nchunkorder : nfdn")
        lbl2 = tkinter.Label(win, font=("Consolas", 14), textvariable=var2)
        lbl2.place(x=255, y=55)
        lbl2.configure(bg=cl0, fg=cl1)

        but3 = tkinter.Button(win, font=("Consolas", 14), text="설치 여부 토글", command=bf3)
        but3.place(x=5, y=305)
        but3.configure(bg=cl2, fg=cl1)
        but4 = tkinter.Button(win, font=("Consolas", 14), text="진행하기", command=bf4)
        but4.place(x=255, y=305)
        but4.configure(bg=cl2, fg=cl1)
        for i in self.name:
            listbox.insert(listbox.size(), i)
        listbox.bind('<ButtonRelease-1>', clickf)

        win.mainloop()

    def guif3(self):
        def exef():
            time.sleep(0.1)
            try:
                var1.set(f"설치 번호 : {self.num}\n데이터 다운로드 : 진행 중\nCRC32 : None\n데이터 언패킹 : 대기\n프로그램 설치 : 대기")
                win.update()
                dwnv = self.download(self.num)
                var1.set(f"설치 번호 : {self.num}\n데이터 다운로드 : {dwnv}\nCRC32 : 진행 중\n데이터 언패킹 : 진행 중\n프로그램 설치 : 대기")
                win.update()
                crcv = self.hash()
                self.unpack()
                var1.set(f"설치 번호 : {self.num}\n데이터 다운로드 : {dwnv}\nCRC32 : {crcv}\n데이터 언패킹 : 완료\n프로그램 설치 : 진행 중")
                self.move(self.num)
                var1.set(f"설치 번호 : {self.num}\n데이터 다운로드 : {dwnv}\nCRC32 : {crcv}\n데이터 언패킹 : 완료\n프로그램 설치 : 완료")
                res = "installed successfully"
            except Exception as e:
                res = str(e)
            self.flag = self.flag + 1
            time.sleep(0.5)
            tkinter.messagebox.showinfo("설치 결과", f" {res} ")
            win.destroy()

        win = tkinter.Tk()
        win.title("Setup5")
        win.geometry("300x250+100+50")
        win.configure(bg=cl0)

        lbl0 = tkinter.Label(win, font=("Consolas", 14), text=f"STEP 3 : {self.name[self.num]} 설치")
        lbl0.place(x=5, y=5)
        lbl0.configure(bg=cl0, fg=cl1)
        var1 = tkinter.StringVar()
        var1.set(f"설치 번호 : {self.num}\n데이터 다운로드 : 대기\nCRC32 : None\n데이터 언패킹 : 대기\n프로그램 설치 : 대기")
        lbl1 = tkinter.Label(win, font=("Consolas", 14), textvariable=var1)
        lbl1.place(x=5, y=55)
        lbl1.configure(bg=cl0, fg=cl1)

        exebut = tkinter.Button(win, font=("Consolas", 14), text="Install", command=exef)
        exebut.place(x=105, y=205)
        exebut.configure(bg=cl2, fg=cl1)

        win.mainloop()

    def guif4(self):
        win = tkinter.Tk()
        win.title("Setup5")
        win.geometry("300x180+100+50")
        win.configure(bg=cl0)

        lbl0 = tkinter.Label(win, font=("Consolas", 14), text="STEP 4 : 작업 마무리")
        lbl0.place(x=5, y=5)
        lbl0.configure(bg=cl0, fg=cl1)
        lbl1 = tkinter.Label(win, font=("Consolas", 14), text=f"작업이 완료되었습니다.\n창을 닫아도 됩니다.\n\n설치된 프로그램 수 : {self.flag}")
        lbl1.place(x=55, y=55)
        lbl1.configure(bg=cl0, fg=cl1)

        win.mainloop()

cl0 = "RoyalBlue1" # 로얄 블루
cl1 = "gold2" # 골드
cl2 = "RoyalBlue2" # 보라
worker = mainclass()
worker.getlist()
worker.guif1()
if worker.status == 1:
    worker.guif2()
    while worker.status == 2:
        if worker.num < len(worker.name):
            if worker.install[worker.num]:
                worker.guif3()
            worker.num = worker.num + 1
        else:
            worker.status = 3
worker.guif4()
time.sleep(0.5)
