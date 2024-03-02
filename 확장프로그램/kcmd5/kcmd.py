# test620 : kcmd

import time
import os

import tkinter
import tkinter.messagebox
import tkinter.ttk

import kdb

def getdata():
    tbox = kdb.toolbox()
    tbox.readfile("./data620.txt")
    out = [0] * tbox.getdata("num")[3]
    for i in range( 0, len(out) ):
        temp = [ "", "", [0] * tbox.getdata(f"{i}.num")[3] ]
        temp[0] = tbox.getdata(f"{i}.name")[3]
        temp[1] = tbox.getdata(f"{i}.info")[3]
        for j in range( 0, len( temp[2] ) ):
            temp[2][j] = tbox.getdata(f"{i}.{j}")[3]
        out[i] = temp
    return out

def execute(orders):
    for i in range( 0, len(orders) ):
        time.sleep(0.1)
        try:
            orders[i] = os.system( orders[i] )
        except:
            orders[i] = 1
    return str(orders)

def main():
    win = tkinter.Tk()
    win.title('Kcmd5')
    win.geometry("270x330+200+100")
    win.resizable(False, False)
    win.configure(bg=c0)

    frame = tkinter.Frame(win)
    frame.place(x=10,y=10)
    listbox = tkinter.Listbox(frame, width=20,  height=10, font=('Consolas', 15), bg=c0, fg=c1, selectbackground=c1, selectforeground=c0)
    listbox.pack(side="left", fill="y")
    scrollbar0 = tkinter.Scrollbar(frame, orient="vertical")
    scrollbar0.config(command=listbox.yview)
    scrollbar0.pack(side="right", fill="y")
    listbox.config(yscrollcommand=scrollbar0.set)
    win.update()

    time.sleep(0.2)
    for i in base:
        listbox.insert( listbox.size(), i[0] )
        time.sleep(0.05)
        win.update()

    status = tkinter.StringVar()
    status.set('====================\n쉽고 빠른 명령어 실행\n====================')
    label0 = tkinter.Label(win, textvariable=status, font=('Consolas', 15), bg=c0, fg=c1)
    label0.place(x=10, y=255)

    past = -1
    def click(event):
        time.sleep(0.1)
        temp = listbox.curselection()[0]
        nonlocal past
        if past == temp:
            tkinter.messagebox.showinfo( "실행 결과", execute( base[past][2] ) )
        else:
            past = temp
            status.set( base[past][1] )
        win.update()
    listbox.bind('<ButtonRelease-1>',click)

    win.mainloop()

c0, c1, base = "RoyalBlue1", "turquoise1", getdata()
main()
time.sleep(0.5)