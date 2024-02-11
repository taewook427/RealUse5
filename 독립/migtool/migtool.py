# test601 : migtool

import time
import os
import sys

import tkinter
import tkinter.ttk
import tkinter.filedialog
import tkinter.messagebox

import pyautogui # --hidden-import=pyautogui
import threading as thr

import kzip

class mainclass:
    def __init__(self, args): # 인자 인식 초기화
        if len(args) == 1:
            self.gettxt()
        else:
            self.finput = args[1] # 프로그램에 투입한 파일
        self.finput = os.path.abspath(self.finput).replace("\\", "/")
        if self.finput[self.finput.rfind("/") + 1:] == "data601.txt":
            with open(self.finput, "r", encoding="utf-8") as f:
                temp = f.readlines()
            self.file = temp[0][0:-1].replace("\r", "")
            self.cur = int( temp[1][0:-1].replace("\r", "") )
        else:
            self.file = self.finput # 프로그램이 읽는 파일
            self.cur = 0 # 현재 진행도 (1~len 까지)

    def gettxt(self): # txt 파일 가져오기
        self.finput = tkinter.filedialog.askopenfile(title='명세 파일 선택', filetypes=(('txt files', '*.txt'),)).name

    def analysis(self): # self.file 구분 분석
        try:
            with open(self.file, "r", encoding="utf-8") as f:
                temp = f.readlines()
            temp = [x.replace("\r", "").replace("\n", "") for x in temp]
            raw = [ ]
            for i in temp:
                if i != "":
                    raw.append(i)
            
            self.call = [ ]
            self.proc = [ ]
            self.subp = [ ]
            reg0 = ""
            reg1 = ""
            reg2 = ""

            ctrlnum = 0
            while ctrlnum < len(raw):
                temp = raw[ctrlnum]
                ctrlnum = ctrlnum + 1
                if temp[0] == "#":
                    pass
                elif temp[0] == "!":
                    reg0 = temp[1:]
                elif temp[0] == "-":
                    reg1 = temp[1:]
                    if ctrlnum == len(raw):
                        self.call.append(reg0)
                        self.proc.append(reg1)
                        self.subp.append(reg2)
                        reg2 = ""
                    elif raw[ctrlnum][0] != "/":
                        self.call.append(reg0)
                        self.proc.append(reg1)
                        self.subp.append(reg2)
                        reg2 = ""
                elif temp[0] == "/":
                    reg2 = temp[1:]
                    self.call.append(reg0)
                    self.proc.append(reg1)
                    self.subp.append(reg2)
                    reg2 = ""
                
        except Exception as e:
            self.call = ["예외 발생"] # 프로시져 묶음
            self.proc = [self.file] # 프로시져들
            self.subp = [str(e)] # 서브프로시져들
        if len(self.call) == 0:
            self.call = ["빈 명세"]
            self.proc = [self.file]
            self.subp = ["확인된 절차가 없습니다"]

    def findnum(self, num): # call의 num째 chunk의 개수
        temp = self.call[num]
        reg0 = num
        while self.call[reg0] == temp:
            reg0 = reg0 - 1
            if reg0 == -1:
                break
        reg1 = num
        while self.call[reg1] == temp:
            reg1 = reg1 + 1
            if reg1 == len(self.call):
                break
        return reg1 - reg0 - 1

    def findpos(self, num): # call의 num째의 chunk 내부에서 상대 위치
        temp = self.call[num]
        reg = num
        while self.call[reg] == temp:
            reg = reg - 1
            if reg == -1:
                break
        return num - reg - 1

    def log(self): # data601.txt로 현재 상태 기록
        temp = self.file + "\n" + str(self.cur) + "\n"
        with open(self.finput[0:self.finput.rfind("/")+1] + "data601.txt", "w", encoding="utf-8") as f:
            f.write(temp)

    def subcap(self, swin, strv): # 사각 지정 캡쳐
        def setf(event): # 설정 함수 AB
            time.sleep(0.1)
            nonlocal status
            if status == 1:
                mem[0], mem[1] = pyautogui.position()
                strv.set(f"({mem[0]}, {mem[1]}) (0, 0)\n캡처 부위 오른쪽 아래에\n마우스를 올리고 엔터")
                swin.update()
                status = 2
            elif status == 2:
                mem[2], mem[3] = pyautogui.position()
                strv.set(f"({mem[0]}, {mem[1]}) ({mem[2]}, {mem[3]})\n엔터를 눌러 촬영")
                swin.update()
                status = 3
            else:
                tpath = f"{int( time.time() )}.png"
                pyautogui.screenshot(tpath, region=(mem[0], mem[1], mem[2]-mem[0], mem[3]-mem[1])).save(tpath)
                tkinter.messagebox.showinfo("부분 스크린샷", f" {tpath}으로 \n 저장되었습니다. ")
                swin.destroy()
            
        strv.set("(0, 0) (0, 0)\n캡처 부위 왼쪽 위에\n마우스를 올리고 엔터")
        swin.update()
        status = 1 # 상태 기록
        mem = [0, 0, 0, 0] # 촬영 부위
        swin.bind("<Return>", setf)

    def mainfunc(self):
        def mf0(): # kzip pack
            time.sleep(0.1)
            tpath = tkinter.filedialog.askdirectory(title='폴더 선택')
            tpath = os.path.abspath(tpath).replace("\\", "/")
            if tpath[-1] == "/":
                tpath = tpath[0:-1]
            try:
                tbox = kzip.toolbox()
                tbox.folder = tpath
                tbox.export = tpath[0:tpath.rfind("/")+1] + "kzip5_result.webp"
                tbox.zipfolder("webp")
                msg = "successfully converted"
            except Exception as e:
                msg = str(e)
            tkinter.messagebox.showinfo("kzip pack 결과", f" {msg} ")

        def mf1(): # kzip unpack
            time.sleep(0.1)
            tpath = tkinter.filedialog.askopenfile(title='파일 선택', filetypes=(('webp files', '*.webp'),('png files', '*.png'),('all files', '*.*'))).name
            tpath = os.path.abspath(tpath).replace("\\", "/")
            try:
                tbox = kzip.toolbox()
                tbox.export = tpath[0:tpath.rfind("/")+1] + "kzip5_result/"
                tbox.unzip(tpath)
                msg = "successfully converted"
            except Exception as e:
                msg = str(e)
            tkinter.messagebox.showinfo("kzip pack 결과", f" {msg} ")

        def mf2(): # full screenshot
            time.sleep(0.1)
            time.sleep(0.5)
            tpath = f"{int( time.time() )}.png"
            pyautogui.screenshot(tpath).save(tpath)
            tkinter.messagebox.showinfo("전체 스크린샷", f" {tpath}으로 \n 저장되었습니다. ")

        def mf3(): # square screenshot
            time.sleep(0.1)
            subwin = tkinter.Toplevel(win)
            subwin.title('migtool5')
            subwin.geometry("300x150+200+100") # lxgp= 450x200+200+100
            subwin.resizable(False, False)
            substr = tkinter.StringVar()
            substr.set("")
            sublbl = tkinter.Label(subwin, font=("맑은 고딕", 14), textvariable=substr)
            sublbl.place(x=5, y=5) # lxgp= x=5 y=5
            x = thr.Thread( target=self.subcap, args=(subwin, substr) )
            x.start()

        def prevf(): # 이전 단계로
            time.sleep(0.1)
            if self.cur > 0:
                self.cur = self.cur - 1
            regen()

        def nextf(): # 다음단계로
            time.sleep(0.1)
            if self.cur < len(self.call) - 1:
                self.cur = self.cur + 1
            regen()

        def regen(): # 새로 그리기
            pbar0.config(value=self.cur+1)
            pbar1.config(maximum=self.findnum(self.cur), value=self.findpos(self.cur)+1)
            strvar1.set( self.call[self.cur] )
            strvar2.set( self.proc[self.cur] )
            strvar3.set( self.subp[self.cur] )
            win.update()
            self.log()
        
        # 메인 윈도우
        win = tkinter.Tk()
        win.title('migtool5')
        win.geometry("500x300+200+100") # lxgp= 700x500+200+100
        win.resizable(False, False)

        mbar = tkinter.Menu(win) # 메뉴 바

        menu0 = tkinter.Menu(mbar, tearoff=0) # teleport
        menu0.add_command(label="KZIP pack", font=("맑은 고딕", 14), command=mf0)
        menu0.add_command(label="KZIP unpack", font=("맑은 고딕", 14), command=mf1)
        mbar.add_cascade(label="  Tool  ", menu=menu0)

        menu1 = tkinter.Menu(mbar, tearoff=0) # general
        menu1.add_command(label="전체 스크린샷", font=("맑은 고딕", 14), command=mf2)
        menu1.add_command(label="부분 스크린샷", font=("맑은 고딕", 14), command=mf3)
        mbar.add_cascade(label="  Picture  ", menu=menu1)

        win.config(menu=mbar)

        pbar0 = tkinter.ttk.Progressbar(win, maximum=len(self.call), value=self.cur+1, length=490) # lxgp= l=690
        pbar0.place(x=5, y=5) # lxgp= x=5 y=5
        strvar0 = tkinter.StringVar()
        strvar0.set(self.finput)
        ent0 = tkinter.Entry(win, textvariable=strvar0, font=("Consolas", 14), width=36, state="readonly") # lxgp= w=36
        ent0.place(x=5, y=35) # lxgp= x=5 y=45
        
        pbar1 = tkinter.ttk.Progressbar(win, maximum=self.findnum(self.cur), value=self.findpos(self.cur)+1, length=490) # lxgp= l=690
        pbar1.place(x=5, y=95) # lxgp= x=5 y=125
        strvar1 = tkinter.StringVar()
        strvar1.set( self.call[self.cur] )
        ent1 = tkinter.Entry(win, textvariable=strvar1, font=("Consolas", 14), width=36, state="readonly") # lxgp= lxgp= w=36
        ent1.place(x=5, y=125) # lxgp= lxgp= x=5 y=165

        but0 = tkinter.Button(win, font=("Consolas", 12), text="\n < \n", command=prevf)
        but0.place(x=5, y=185) # lxgp= x=5 y=265
        but1 = tkinter.Button(win, font=("Consolas", 12), text="\n > \n", command=nextf)
        but1.place(x=445, y=185) # lxgp= x=615 y=265

        strvar2 = tkinter.StringVar()
        strvar2.set( self.proc[self.cur] )
        ent2 = tkinter.Entry(win, textvariable=strvar2, font=("Consolas", 14), width=28, state="readonly") # lxgp=
        ent2.place(x=65, y=185) # lxgp= x=85 y=270
        strvar3 = tkinter.StringVar()
        strvar3.set( self.subp[self.cur] )
        ent3 = tkinter.Entry(win, textvariable=strvar3, font=("Consolas", 14), width=28, state="readonly") # lxgp=
        ent3.place(x=65, y=225) # lxgp= x=85 y=350

        win.mainloop()

classloader = mainclass(sys.argv)
classloader.analysis()
classloader.mainfunc()
time.sleep(0.5)
