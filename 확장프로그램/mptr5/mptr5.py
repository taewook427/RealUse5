# test619 : mouse pointer
# --hidden-import "pynput.keyboard._win32" --hidden-import "pynput.mouse._win32" //win
# --hidden-import "pynput.keyboard._xorg" --hidden-import "pynput.mouse._xorg" //linux
import time
import tkinter
import pyautogui
c0, c1, work = "dark green", "lawn green", True

win = tkinter.Tk()
win.title('MousePtr5')
win.geometry("600x300+200+150") # lxgp=900x400+400+300
win.resizable(False, False)
win.configure(bg=c0)

svar0, svar1, svar2 = tkinter.StringVar(), tkinter.StringVar(), tkinter.StringVar()
svar0.set("")
svar1.set("")
svar2.set("")
lbl0 = tkinter.Label(win, font=("Consolas", 20), textvariable=svar0, bg=c0, fg=c1)
lbl0.place(x=5, y=5)
lbl1 = tkinter.Label(win, font=("Consolas", 20), textvariable=svar1, bg=c0, fg=c1)
lbl1.place(x=5, y=65) # lxgp= x=5 y=85
lbl2 = tkinter.Label(win, font=("Consolas", 20), textvariable=svar2, bg=c0, fg=c1)
lbl2.place(x=5, y=125) # lxgp= x=5 y=165

def lockf():
    time.sleep(0.1)
    if lockv.get() == 0:
        win.wm_attributes("-topmost", 0)
    else:
        win.wm_attributes("-topmost", 1)
lockv = tkinter.IntVar()
lockb = tkinter.Checkbutton(win, text="Always On Display", font=("Consolas", 20), variable=lockv, command=lockf, bg=c0, fg=c1)
lockb.place(x=5, y=185) # lxgp= x=5 y=245

def shutdown():
    global work
    work = False
    win.destroy()
win.protocol('WM_DELETE_WINDOW', shutdown)

while work:
    time.sleep(0.05)
    tmp0 = pyautogui.size()
    svar0.set(f"Display : X = {tmp0[0]}, Y = {tmp0[1]}")
    tmp1 = pyautogui.position()
    svar1.set(f"MousePtr : X = {tmp1[0]}, Y = {tmp1[1]}")
    svar2.set(f"Relative : X = {tmp1[0]-tmp0[0]//2}, Y = {tmp1[1]-tmp0[1]//2}")
    win.update()
time.sleep(0.5)
