# test599 : kexec (py)

import time
import os

import tkinter
import tkinter.messagebox
import tkinter.ttk

import kdb

tbox = kdb.toolbox()
tbox.readfile("data599.txt")
title = tbox.getdata("title")[3]
length = tbox.getdata("num")[3]
path = [""] * length
name = [""] * length
for i in range(0, length):
    name[i] = tbox.getdata(f"{i}.name")[3]
    path[i] = tbox.getdata(f"{i}.path")[3]
    path[i] = path[i].replace("\\", "/")
cl0 = "plum1"
cl1 = "deep pink"

win = tkinter.Tk()
win.title(title)
win.geometry("320x360+100+50")
win.resizable(False, False)
win.configure(bg=cl0)

frame = tkinter.Frame(win)
frame.place(x=10,y=10)
listbox = tkinter.Listbox(
    frame, width=25,  height=14, font = ('Consolas', 15),
    bg=cl0, fg=cl1, selectbackground=cl1, selectforeground=cl0)
listbox.pack(side="left", fill="y")
scrollbar0 = tkinter.Scrollbar(frame, orient="vertical")
scrollbar0.config(command=listbox.yview)
scrollbar0.pack(side="right", fill="y")
listbox.config(yscrollcommand=scrollbar0.set)
win.update()

time.sleep(0.5)
for i in name:
    listbox.insert( listbox.size(),i )
    win.update()
    time.sleep(0.1)

last = -1
def click(event):
    time.sleep(0.1)
    global name
    global path
    global listbox
    global last
    temp = listbox.curselection()[0]
    if last == temp:
        time.sleep(0.1)
        current = os.getcwd()
        pmem = path[last]
        os.chdir( pmem[ 0:pmem.rfind("/") ] )
        os.startfile(pmem)
        os.chdir(current)
    else:
        last = temp
listbox.bind('<ButtonRelease-1>',click)

win.mainloop()
time.sleep(0.5)
