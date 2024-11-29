# test711 : independent.powercut

import time
import os

import tkinter
import tkinter.messagebox

import kobj

def wait(day, hour): # wait N days & time N hour
    for i in range(0, day):
        for j in range(0, 24):
            time.sleep(3600)
    while hour != time.localtime( time.time() ).tm_hour:
        time.sleep(300)

def genlog(path, word): # add log to path
    t = int( time.time() )
    temp = time.strftime( "%Y.%m.%d;%H:%M:%S", time.localtime(t) )
    temp = bytes(f"{temp}/{t}#{word}\n", encoding="utf-8")
    if os.path.exists(path):
        with open(path, "ab") as f:
            f.write(temp)
    else:
        with open(path, "wb") as f:
            f.write(temp)

def askoff(): # ask user to shutdown
    def f0():
        time.sleep(0.1)
        win.destroy()
    def f1():
        time.sleep(0.1)
        if tkinter.messagebox.askokcancel(title="재부팅 확인", message=" 이 컴퓨터가 재부팅됩니다. "):
            win.destroy()
            global flag
            flag = True
    global flag
    flag = False
    win = tkinter.Tk()
    win.title("POWERCUT")
    win.configure(bg="LightSteelBlue1")
    lbl = tkinter.Label(win, text=f"\n이 컴퓨터를 {day}일 동안 사용하였습니다.\n재시작을 권장합니다.\n\n", font=("맑은 고딕", 14), background="LightSteelBlue1")
    lbl.pack()
    but0 = tkinter.Button(win, text="나중에 다시 알림", font=("맑은 고딕", 14), background="LightSteelBlue1", command=f0)
    but0.pack(side="left")
    but1 = tkinter.Button(win, text="    지금 재부팅    ", font=("맑은 고딕", 14), background="LightSteelBlue1", command=f1)
    but1.pack(side="right")
    win.mainloop()

path = kobj.repath()[0].replace("\\", "/")
path = path[path.rfind("/")+1:]
pos0, pos1, pos2 = path.find("_"), path.rfind("_"), path.find(".")
try:
    day = int( path[pos0+1:pos1] ) # _N_
    hour = path[pos1+1:pos2] # _Nf.
    if hour[-1] == "f":
        hour, force = int( hour[:-1] ), True
    else:
        hour, force = int(hour), False
except:
    day, hour, force = 3, 3, False
day, hour = max(0, day), min(max(0, hour), 23)

desktop = os.path.join(os.path.expanduser("~"),"Desktop")
alert, log = os.path.join(desktop, "PC5_alert.txt"), os.path.join(desktop, "powerlog.txt")
time.sleep(180) ###
if os.path.exists(alert):
    os.remove(alert)
genlog(log, "POWERCUT_ON")
wait(day, hour) ###

if force:
    genlog(alert, "ALERT")
    time.sleep(900)
    os.remove(alert)
    genlog(log, "POWERCUT_OFF")
    os.system('shutdown -f -r -t 5')
else:
    flag = False
    while not flag:
        askoff()
        if flag:
            break
        wait(day, hour) ###
    genlog(log, "POWERCUT_OFF")
    os.system('shutdown -f -r -t 5')
