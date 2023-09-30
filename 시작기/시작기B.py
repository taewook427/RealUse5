# test518 : gen5linux 시작기

import time
import os
import subprocess

import tkinter
import tkinter.messagebox
from tkinter import filedialog

import clipboard

class mainclass:
    
    def __init__(self):
        temp = os.listdir("./")
        self.helper = "" # 선택 도우미 경로
        for i in temp:
            if "helper5" in i:
                self.helper = i
        temp = os.listdir("./extension")
        self.kmain = [ ] # program path
        self.kicon = [ ] # program icon
        self.kname = [ ] # program name
        for i in temp:
            names = os.listdir(f"./extension/{i}")
            nm0 = ""
            nm1 = ""
            for j in names:
                if "kmain5" in j:
                    nm0 = j
                elif "kicon5" in j:
                    nm1 = j
            if nm0 != "":
                self.kmain.append(nm0)
                self.kname.append(i)
                if nm1 == "":
                    self.kicon.append(f"./fill.png")
                else:
                    self.kicon.append(f"./extension/{i}/{nm1}")
        if len(self.kname) == 0:
            add = 6
        else:
            add = 5 - (len(self.kname) - 1) % 6
        self.kmain = self.kmain + [""] * add
        self.kicon = self.kicon + [f"./fill.png"] * add
        self.kname = self.kname + [""] * add

    def mainfunc(self):
        win = tkinter.Tk()
        win.title("5세대 시작기")
        win.geometry("600x500+100+50")

        current = 0 # 현재 페이지
        pmax = len(self.kname) // 6 # 페이지 수

        def func0():
            time.sleep(0.1)
            try:
                clipboard.copy("./" + self.helper)
            except:
                tkinter.messagebox.showinfo("클립보드 복사오류", f" 현재 시스템에서 복사-붙여넣기 할 수 없습니다. 다음 두 명령어 중 하나를 실행하십시오. \n sudo apt-get install xsel \n sudo apt-get install xclip ")
            subprocess.run(['xdg-open', './'])

        def func1():
            time.sleep(0.1)
            nonlocal win
            win.destroy()
            self.divfile()

        def func2():
            time.sleep(0.1)
            nonlocal win
            win.destroy()
            self.getfile()

        def func3():
            time.sleep(0.1)
            nonlocal current
            nonlocal pmax
            if current != 0:
                current = current - 1
                nonlocal show5
                nonlocal photo5
                nonlocal but5
                show5.set( self.kname[6 * current][0:15] )
                photo5 = tkinter.PhotoImage( file = self.kicon[6 * current] )
                but5.config(image = photo5)
                nonlocal show6
                nonlocal photo6
                nonlocal but6
                show6.set( self.kname[6 * current + 1][0:15] )
                photo6 = tkinter.PhotoImage( file = self.kicon[6 * current + 1] )
                but6.config(image = photo6)
                nonlocal show7
                nonlocal photo7
                nonlocal but7
                show7.set( self.kname[6 * current + 2][0:15] )
                photo7 = tkinter.PhotoImage( file = self.kicon[6 * current + 2] )
                but7.config(image = photo7)
                nonlocal show8
                nonlocal photo8
                nonlocal but8
                show8.set( self.kname[6 * current + 3][0:15] )
                photo8 = tkinter.PhotoImage( file = self.kicon[6 * current + 3] )
                but8.config(image = photo8)
                nonlocal show9
                nonlocal photo9
                nonlocal but9
                show9.set( self.kname[6 * current + 4][0:15] )
                photo9 = tkinter.PhotoImage( file = self.kicon[6 * current + 4] )
                but9.config(image = photo9)
                nonlocal show10
                nonlocal photo10
                nonlocal but10
                show10.set( self.kname[6 * current + 5][0:15] )
                photo10 = tkinter.PhotoImage( file = self.kicon[6 * current + 5] )
                but10.config(image = photo10)
                show11.set(f"( {current + 1} / {pmax} )")

        def func4():
            time.sleep(0.1)
            nonlocal current
            nonlocal pmax
            if current + 1 < pmax:
                current = current + 1
                nonlocal show5
                nonlocal photo5
                nonlocal but5
                show5.set( self.kname[6 * current][0:15] )
                photo5 = tkinter.PhotoImage( file = self.kicon[6 * current] )
                but5.config(image = photo5)
                nonlocal show6
                nonlocal photo6
                nonlocal but6
                show6.set( self.kname[6 * current + 1][0:15] )
                photo6 = tkinter.PhotoImage( file = self.kicon[6 * current + 1] )
                but6.config(image = photo6)
                nonlocal show7
                nonlocal photo7
                nonlocal but7
                show7.set( self.kname[6 * current + 2][0:15] )
                photo7 = tkinter.PhotoImage( file = self.kicon[6 * current + 2] )
                but7.config(image = photo7)
                nonlocal show8
                nonlocal photo8
                nonlocal but8
                show8.set( self.kname[6 * current + 3][0:15] )
                photo8 = tkinter.PhotoImage( file = self.kicon[6 * current + 3] )
                but8.config(image = photo8)
                nonlocal show9
                nonlocal photo9
                nonlocal but9
                show9.set( self.kname[6 * current + 4][0:15] )
                photo9 = tkinter.PhotoImage( file = self.kicon[6 * current + 4] )
                but9.config(image = photo9)
                nonlocal show10
                nonlocal photo10
                nonlocal but10
                show10.set( self.kname[6 * current + 5][0:15] )
                photo10 = tkinter.PhotoImage( file = self.kicon[6 * current + 5] )
                but10.config(image = photo10)
                show11.set(f"( {current + 1} / {pmax} )")

        but0 = tkinter.Button(win, text = "선택 도우미", font = ("맑은 고딕", 14), command = func0)
        but0.place(x = 5, y = 5)
        but1 = tkinter.Button(win, text = "조각 나누기", font = ("맑은 고딕", 14), command = func1)
        but1.place(x = 150, y = 5)
        but2 = tkinter.Button(win, text = "조각 모으기", font = ("맑은 고딕", 14), command = func2)
        but2.place(x = 295, y = 5)
        but3 = tkinter.Button(win, text = "\n\n\n\n\n\n\n < \n\n\n\n\n\n\n", font = ("맑은 고딕", 14), command = func3)
        but3.place(x = 5, y = 55)
        but4 = tkinter.Button(win, text = "\n\n\n\n\n\n\n > \n\n\n\n\n\n\n", font = ("맑은 고딕", 14), command = func4)
        but4.place(x = 550, y = 55)

        def func5():
            time.sleep(0.1)
            nonlocal current
            count = 6 * current
            if self.kname[count] != "":
                nonlocal mode12
                tgt = "./extension/" + self.kname[count]
                if mode12 == 0:
                    try:
                        clipboard.copy( "./" + self.kmain[count] )
                    except:
                        tkinter.messagebox.showinfo("클립보드 복사오류", f" 현재 시스템에서 복사-붙여넣기 할 수 없습니다. 다음 두 명령어 중 하나를 실행하십시오. \n sudo apt-get install xsel \n sudo apt-get install xclip ")
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )
                else:
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )

        def func6():
            time.sleep(0.1)
            nonlocal current
            count = 6 * current + 1
            if self.kname[count] != "":
                nonlocal mode12
                tgt = "./extension/" + self.kname[count]
                if mode12 == 0:
                    try:
                        clipboard.copy( "./" + self.kmain[count] )
                    except:
                        tkinter.messagebox.showinfo("클립보드 복사오류", f" 현재 시스템에서 복사-붙여넣기 할 수 없습니다. 다음 두 명령어 중 하나를 실행하십시오. \n sudo apt-get install xsel \n sudo apt-get install xclip ")
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )
                else:
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )

        def func7():
            time.sleep(0.1)
            nonlocal current
            count = 6 * current + 2
            if self.kname[count] != "":
                nonlocal mode12
                tgt = "./extension/" + self.kname[count]
                if mode12 == 0:
                    try:
                        clipboard.copy( "./" + self.kmain[count] )
                    except:
                        tkinter.messagebox.showinfo("클립보드 복사오류", f" 현재 시스템에서 복사-붙여넣기 할 수 없습니다. 다음 두 명령어 중 하나를 실행하십시오. \n sudo apt-get install xsel \n sudo apt-get install xclip ")
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )
                else:
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )

        def func8():
            time.sleep(0.1)
            nonlocal current
            count = 6 * current + 3
            if self.kname[count] != "":
                nonlocal mode12
                tgt = "./extension/" + self.kname[count]
                if mode12 == 0:
                    try:
                        clipboard.copy( "./" + self.kmain[count] )
                    except:
                        tkinter.messagebox.showinfo("클립보드 복사오류", f" 현재 시스템에서 복사-붙여넣기 할 수 없습니다. 다음 두 명령어 중 하나를 실행하십시오. \n sudo apt-get install xsel \n sudo apt-get install xclip ")
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )
                else:
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )

        def func9():
            time.sleep(0.1)
            nonlocal current
            count = 6 * current + 4
            if self.kname[count] != "":
                nonlocal mode12
                tgt = "./extension/" + self.kname[count]
                if mode12 == 0:
                    try:
                        clipboard.copy( "./" + self.kmain[count] )
                    except:
                        tkinter.messagebox.showinfo("클립보드 복사오류", f" 현재 시스템에서 복사-붙여넣기 할 수 없습니다. 다음 두 명령어 중 하나를 실행하십시오. \n sudo apt-get install xsel \n sudo apt-get install xclip ")
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )
                else:
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )

        def func10():
            time.sleep(0.1)
            nonlocal current
            count = 6 * current + 5
            if self.kname[count] != "":
                nonlocal mode12
                tgt = "./extension/" + self.kname[count]
                if mode12 == 0:
                    try:
                        clipboard.copy( "./" + self.kmain[count] )
                    except:
                        tkinter.messagebox.showinfo("클립보드 복사오류", f" 현재 시스템에서 복사-붙여넣기 할 수 없습니다. 다음 두 명령어 중 하나를 실행하십시오. \n sudo apt-get install xsel \n sudo apt-get install xclip ")
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )
                else:
                    tgt = tgt.replace("\\", "/")
                    subprocess.run( ['xdg-open', tgt] )

        show5 = tkinter.StringVar()
        show5.set( self.kname[6 * current][0:15] )
        show6 = tkinter.StringVar()
        show6.set( self.kname[6 * current + 1][0:15] )
        show7 = tkinter.StringVar()
        show7.set( self.kname[6 * current + 2][0:15] )
        show8 = tkinter.StringVar()
        show8.set( self.kname[6 * current + 3][0:15] )
        show9 = tkinter.StringVar()
        show9.set( self.kname[6 * current + 4][0:15] )
        show10 = tkinter.StringVar()
        show10.set( self.kname[6 * current + 5][0:15] )

        photo5 = tkinter.PhotoImage( file = self.kicon[6 * current] )
        photo6 = tkinter.PhotoImage( file = self.kicon[6 * current + 1] )
        photo7 = tkinter.PhotoImage( file = self.kicon[6 * current + 2] )
        photo8 = tkinter.PhotoImage( file = self.kicon[6 * current + 3] )
        photo9 = tkinter.PhotoImage( file = self.kicon[6 * current + 4] )
        photo10 = tkinter.PhotoImage( file = self.kicon[6 * current + 5] )

        but5 = tkinter.Button(win, image = photo5, font = ("맑은 고딕", 14), command = func5)
        but5.place(x = 55, y = 55)
        lbl5 = tkinter.Label( win, textvariable = show5, font = ("맑은 고딕", 14) )
        lbl5.place(x = 55, y = 215)

        but6 = tkinter.Button(win, image = photo6, font = ("맑은 고딕", 14), command = func6)
        but6.place(x = 220, y = 55)
        lbl6 = tkinter.Label( win, textvariable = show6, font = ("맑은 고딕", 14) )
        lbl6.place(x = 220, y = 215)

        but7 = tkinter.Button(win, image = photo7, font = ("맑은 고딕", 14), command = func7)
        but7.place(x = 385, y = 55)
        lbl7 = tkinter.Label( win, textvariable = show7, font = ("맑은 고딕", 14) )
        lbl7.place(x = 385, y = 215)

        but8 = tkinter.Button(win, image = photo8, font = ("맑은 고딕", 14), command = func8)
        but8.place(x = 55, y = 255)
        lbl8 = tkinter.Label( win, textvariable = show8, font = ("맑은 고딕", 14) )
        lbl8.place(x = 55, y = 415)

        but9 = tkinter.Button(win, image = photo9, font = ("맑은 고딕", 14), command = func9)
        but9.place(x = 220, y = 255)
        lbl9 = tkinter.Label( win, textvariable = show9, font = ("맑은 고딕", 14) )
        lbl9.place(x = 220, y = 415)

        but10 = tkinter.Button(win, image = photo10, font = ("맑은 고딕", 14), command = func10)
        but10.place(x = 385, y = 255)
        lbl10 = tkinter.Label( win, textvariable = show10, font = ("맑은 고딕", 14) )
        lbl10.place(x = 385, y = 415)

        show11 = tkinter.StringVar()
        show11.set(f"( {current + 1} / {pmax} )")
        lbl11 = tkinter.Label( win, textvariable = show11, font = ("맑은 고딕", 14) )
        lbl11.place(x = 265, y = 460)

        def func12():
            time.sleep(0.1)
            nonlocal mode12
            nonlocal show12
            if mode12 == 0:
                mode12 = 1
                show12.set("모드 : 폴더 열기")
            else:
                mode12 = 0
                show12.set("모드 : 파일 실행")

        show12 = tkinter.StringVar()
        show12.set("모드 : 파일 실행")
        mode12 = 0 # 0 : 파일실행, 1 : 폴더열기
            
        but12 = tkinter.Button(win, textvariable = show12, font = ("맑은 고딕", 14), command = func12)
        but12.place(x = 435, y = 5)

        win.mainloop()

    def divfile(self):
        win = tkinter.Tk()
        win.title("5세대 시작기")
        win.geometry("300x200+100+50")

        def func0():
            time.sleep(0.1)
            nonlocal show1
            try:
                text = filedialog.askopenfile( title='ZIP 파일 선택', filetypes = ( ("zip files", "*.zip"), ("all files", "*.*") ) ).name
            except:
                text = ""
            text = text.replace("\\", "/")
            show1.set(text)

        but0 = tkinter.Button(win, text = ". . .", font = ("맑은 고딕", 14), command = func0)
        but0.place(x = 5, y = 5)

        show1 = tkinter.StringVar()
        show1.set("")
        in1 = tkinter.Entry(win, textvariable = show1, width = 23, font = ("맑은 고딕", 14), state = "readonly")
        in1.place(x = 55, y = 5)

        lbl2 = tkinter.Label( win, text = "분할 크기 (MiB) : ", font = ("맑은 고딕", 14) )
        lbl2.place(x = 5, y = 65)

        in3 = tkinter.Entry( win, width = 10, font = ("맑은 고딕", 14) )
        in3.place(x = 165, y = 70)

        def func4():
            time.sleep(0.1)
            nonlocal show1
            nonlocal in3
            temp = show1.get()
            path = temp[0:temp.rfind('/') + 1] # 폴더 경로 (~/)
            name = temp[temp.rfind('/') + 1:] # 이름
            div = int( in3.get() ) * 1024 * 1024 # 분할 크기

            with open(temp, 'rb') as f:
                data = f.read()
            size = len(data)
            count = 0

            for i in range(0, size // div):
                with open(path + name + f".{count}", "wb") as f:
                    f.write( data[count * div : count * div + div] )
                count = count + 1
            if size % div != 0:
                with open(path + name + f".{count}", "wb") as f:
                    f.write( data[count * div:] )
                count = count + 1

            tkinter.messagebox.showinfo("나누기 완료", f" {name}.0 ~ {name}.{count - 1} \n {count}개 파일 생성되었습니다. ")

        but4 = tkinter.Button(win, text = "나누기", font = ("맑은 고딕", 14), command = func4)
        but4.place(x = 5, y = 125)

        win.mainloop()

    def getfile(self):
        win = tkinter.Tk()
        win.title("5세대 시작기")
        win.geometry("300x200+100+50")

        def func0():
            time.sleep(0.1)
            nonlocal path
            nonlocal name
            nonlocal num
            try:
                text = filedialog.askopenfile( title='시작 파일 선택', filetypes = ( ("start files", "*.0"), ("all files", "*.*") ) ).name
            except:
                text = ""
            text = text.replace("\\", "/")
            if text != "":
                path = text[0:text.rfind("/") + 1]
                name = text[text.rfind("/") + 1:-2]
                num = 0
                while os.path.isfile(f"{path}{name}.{num}"):
                    num = num + 1
                nonlocal show1
                show1.set(f"조각 이름 : {name}\n\n조각 개수 : {num}")

        but0 = tkinter.Button(win, text = ". . .", font = ("맑은 고딕", 14), command = func0)
        but0.place(x = 5, y = 5)
        path = "" # 폴더 경로 (~/)
        name = "" # 결과 파일 이름
        num = 0 # 타겟 파일 개수

        show1 = tkinter.StringVar()
        show1.set("조각 이름 : \n\n조각 개수 : ")
        lbl1 = tkinter.Label( win, textvariable = show1, font = ("맑은 고딕", 14) )
        lbl1.place(x = 5, y = 55)

        def func2():
            time.sleep(0.1)
            nonlocal path
            nonlocal name
            nonlocal num
            if path != "":
                with open(path + name, "wb") as f:
                    for i in range(0, num):
                        with open(path + name + f".{i}", "rb") as t:
                            f.write( t.read() )
            tkinter.messagebox.showinfo("모으기 완료", f" {num}개 파일 모음 \n {name} 파일 생성되었습니다. ")

        but2 = tkinter.Button(win, text = "모으기", font = ("맑은 고딕", 14), command = func2)
        but2.place(x = 5, y = 145)

        win.mainloop()

# 모든 아이콘 사진은 150 * 150 픽셀 사이즈로
k = mainclass()
k.mainfunc()
time.sleep(0.5)
