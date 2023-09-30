# test516 : 선택도우미 / win & linux

import time

import tkinter
import tkinter.messagebox
from tkinter import filedialog

import clipboard

style = 0 # 0 : linux, 1 : windows
mode = 0 # 0 : file, 1 : folder

def func0():
    time.sleep(0.1)
    global show0
    global style
    if style == 0:
        style = 1
        show0.set("style : windows")
    else:
        style = 0
        show0.set("style : linux")

def func1():
    time.sleep(0.1)
    global show1
    global mode
    if mode == 0:
        mode = 1
        show1.set("mode : folder")
    else:
        mode = 0
        show1.set("mode : file")

def func2():
    time.sleep(0.1)
    global show3
    if mode == 0:
        ftype = (
            ('all files', '*.*'),('k files', '*.k'),('png files', '*.png'),('jpg files', '*.jpg'),
            ('webp files', '*.webp'),('txt files', '*.txt'),('mp4 files', '*.mp4'),('exe files', '*.exe')
            )
        text = filedialog.askopenfile(title='파일 선택', filetypes = ftype).name
    else:
        text = filedialog.askdirectory(title='폴더 선택')
    if style == 0:
        text = text.replace("\\", "/")
    else:
        text = text.replace("/", "\\")
    show3.set(text)
    try:
        clipboard.copy(text)
    except:
        tkinter.messagebox.showinfo("클립보드 복사오류", f" 현재 시스템에서 복사-붙여넣기 할 수 없습니다. 다음 두 명령어 중 하나를 실행하십시오. \n sudo apt-get install xsel \n sudo apt-get install xclip ")

win = tkinter.Tk()
win.title("선택 도우미")
win.geometry("300x200+100+50")
show0 = tkinter.StringVar()
show0.set("style : linux")
but0 = tkinter.Button(win, textvariable = show0, font = ("맑은 고딕", 14), command = func0)
but0.place(x = 5, y = 5)
show1 = tkinter.StringVar()
show1.set("mode : file")
but1 = tkinter.Button(win, textvariable = show1, font = ("맑은 고딕", 14), command = func1)
but1.place(x = 5, y = 55)
but2 = tkinter.Button(win, text = "browse", font = ("맑은 고딕", 14), command = func2)
but2.place(x = 5, y = 105)
show3 = tkinter.StringVar()
in3 = tkinter.Entry(win, width = 28, textvariable = show3, font = ("맑은 고딕", 14), state = "readonly")
in3.place(x = 5, y = 160)

win.mainloop()
