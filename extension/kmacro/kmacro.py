# test735 : extension.kmacro
# 윈도우 빌드 추가 : --hidden-import "pynput.keyboard._win32" --hidden-import "pynput.mouse._win32"
# 리눅스 빌드 추가 : --hidden-import "pynput.keyboard._xorg" --hidden-import "pynput.mouse._xorg" --hidden-import='PIL._tkinter_finder'

import os
import time
import ctypes
import base64
import requests
import threading
import subprocess

import tkinter as tk
import tkinter.ttk as tkt
import tkinter.filedialog as tkf
import tkinter.messagebox as tkm
from PIL import Image, ImageTk
import blocksel

import pyautogui
import keyboard
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.edge.service import Service

import kobj
import kdb
import runtime
import kscript_lib
import kscript_macro

class capture:
    def __init__(self, iswin, ff):
        self.mwin, self.comp, self.parms = None, dict(), {"color":"wheat1", "activate":"green yellow"}
        self.path_ff, self.pic_path, self.pic_num, self.status_lock, self.status_cut, self.key = ff, "", 0, False, False, None
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "280x300+200+100", ("맑은 고딕", 10), ("Consolas", 12)
            self.parms["combo.1"], self.parms["ent.2"] = 6, 5
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "360x340+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["combo.1"], self.parms["ent.2"] = 6, 5
        self.entry()
        thr = threading.Thread(target=self.getkey)
        thr.start()
        self.guiloop()

    def entry(self):
        def gui_exit(): # close the screen
            time.sleep(0.1)
            self.mwin.destroy()
            self.mwin = None

        def click_lock(): # lock button clicked
            time.sleep(0.1)
            self.status_lock = not self.status_lock
            if self.status_lock:
                self.comp["button1a"].configure( bg=self.parms["activate"] )
                self.mwin.wm_attributes("-topmost", 1)
            else:
                self.comp["button1a"].configure( bg=self.parms["color"] )
                self.mwin.wm_attributes("-topmost", 0)

        def click_cut(): # cut button clicked
            time.sleep(0.1)
            self.status_cut = not self.status_cut
            if self.status_cut:
                self.comp["button1b"].configure( bg=self.parms["activate"] )
            else:
                self.comp["button1b"].configure( bg=self.parms["color"] )

        def cap_screen(): # capture screen
            self.shot()
            time.sleep(0.05)

        self.mwin = tk.Tk() # main window
        self.mwin.title("Kmacro Capture")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )
        self.mwin.protocol('WM_DELETE_WINDOW', gui_exit)

        self.comp["photo"] = None # ImageTk photo
        self.comp["canvas"] = tk.Canvas(self.mwin, width=256, height=144)
        self.comp["canvas"].pack(padx=10, pady=10)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # lock, screen, key
        self.comp["frame1"].pack(fill="x", padx=5, pady=5)
        self.comp["button1a"] = tk.Button(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Lock", command=click_lock)
        self.comp["button1a"].grid(row=0, column=0)
        self.comp["button1b"] = tk.Button(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Cut", command=click_cut)
        self.comp["button1b"].grid(row=0, column=1)
        self.comp["combo1"] = tkt.Combobox( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["combo.1"], values=["f2", "f8", "f11", "ctrl", "shift", "PrtSc", "Manual"] )
        self.comp["combo1"].set("f2")
        self.comp["combo1"].grid(row=0, column=2, padx=5)
        self.comp["label1"] = tk.Label(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="x0y0 x1y1")
        self.comp["label1"].grid(row=0, column=3)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # cut cord
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["entry2a"] = tk.Entry( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"] )
        self.comp["entry2a"].grid(row=0, column=0, padx=5)
        self.comp["entry2b"] = tk.Entry( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"] )
        self.comp["entry2b"].grid(row=0, column=1, padx=5)
        self.comp["entry2c"] = tk.Entry( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"] )
        self.comp["entry2c"].grid(row=0, column=2, padx=5)
        self.comp["entry2d"] = tk.Entry( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"] )
        self.comp["entry2d"].grid(row=0, column=3, padx=5)

        self.comp["frame3"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # info, save
        self.comp["frame3"].pack(fill="x", padx=5, pady=5)
        self.comp["strvar3"] = tk.StringVar()
        self.comp["strvar3"].set("p0000 x00000 y00000")
        self.comp["label3"] = tk.Label(self.comp["frame3"], font=self.parms["font.1"], bg=self.parms["color"], textvariable=self.comp["strvar3"])
        self.comp["label3"].pack(side="left")
        self.comp["button3"] = tk.Button(self.comp["frame3"], font=self.parms["font.0"], bg=self.parms["color"], text="Capture", command=cap_screen)
        self.comp["button3"].pack(side="right")

    def guiloop(self):
        while self.mwin != None:
            temp = pyautogui.position()
            self.comp["strvar3"].set(f"p{self.pic_num:04} x{temp[0]} y{temp[1]}")
            self.mwin.update()
            temp = self.comp["combo1"].get()
            if temp == "PrtSc":
                temp = "print screen"
            if self.key == temp:
                self.shot()
                self.key = None
            time.sleep(0.03)

    def getkey(self):
        while self.mwin != None:
            self.key = keyboard.read_key()
            time.sleep(0.05)

    def shot(self):
        path = f"./_ST5_DATA/{self.pic_num:04d}.png"
        if self.status_cut:
            x0, y0, x1, y1 = int( self.comp["entry2a"].get() ), int( self.comp["entry2b"].get() ), int( self.comp["entry2c"].get() ), int( self.comp["entry2d"].get() )
            pyautogui.screenshot( path, region=(x0, y0, x1-x0, y1-y0) )
        else:
            pyautogui.screenshot(path)
        self.comp["photo"] = ImageTk.PhotoImage( Image.open(path).resize( (256, 144) ) )
        self.comp["canvas"].create_image( 128, 72, image=self.comp["photo"] )
        self.pic_num = self.pic_num + 1

class web_dwn:
    def __init__(self, iswin, ff):
        self.mwin, self.comp, self.parms, self.status_lock = None, dict(), {"color":"wheat1", "activate":"green yellow"}, False
        self.path_ff = os.path.abspath(ff + "msedgedriver.exe") if iswin else os.path.abspath(ff + "msedgedriver")
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "350x150+200+100", ("맑은 고딕", 10), ("Consolas", 12)
            self.parms["ent.0"], self.parms["combo.1"], self.parms["ent.2"] = 25, 8, 18
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "400x180+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["combo.1"], self.parms["ent.2"] = 25, 8, 16
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def click_lock(): # lock button clicked
            time.sleep(0.1)
            self.status_lock = not self.status_lock
            if self.status_lock:
                self.comp["button1"].configure( bg=self.parms["activate"] )
                self.mwin.wm_attributes("-topmost", 1)
            else:
                self.comp["button1"].configure( bg=self.parms["color"] )
                self.mwin.wm_attributes("-topmost", 0)

        def save_url(): # save link
            time.sleep(0.1)
            drv, tp, nm = getdrv(self.path_ff), self.comp["combo1"].get(), self.comp['entry2'].get()
            drv.get( self.comp["entry0"].get() )
            drv.implicitly_wait(90)
            time.sleep(0.5)
            if tp == "PDF":
                with open(f"./_ST5_DATA/{nm}.pdf", 'wb') as f:
                    f.write( base64.b64decode( drv.print_page() ) )
            elif tp == "Image":
                drv.save_screenshot(f"./_ST5_DATA/{nm}.png")
            elif tp == "Source":
                temp = drv.page_source
                with open(f"./_ST5_DATA/{nm}.html", 'w', encoding="utf-8") as f:
                    f.write(temp)
                num, scripts = 0, drv.find_elements(By.TAG_NAME, "script")
                for i in scripts:
                    src = i.get_attribute("src")
                    if src != None:
                        with open(f"./_ST5_DATA/{nm}_{num}.js", 'w', encoding="utf-8") as f:
                            f.write(requests.get(src).text)
                    num = num + 1
                num, links = 0, drv.find_elements(By.TAG_NAME, "link")
                for i in links:
                    href = i.get_attribute("href")
                    if href != None and 'stylesheet' in i.get_attribute('rel'):
                        with open(f"./_ST5_DATA/{nm}_{num}.css", 'w', encoding="utf-8") as f:
                            f.write(requests.get(href).text)
                    num = num + 1

        self.mwin = tk.Tk() # main window
        self.mwin.title("Kmacro Web Download")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # url
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["label0"] = tk.Label(self.comp["frame0"], font=self.parms["font.0"], bg=self.parms["color"], text="URL")
        self.comp["label0"].grid(row=0, column=0)
        self.comp["entry0"] = tk.Entry( self.comp["frame0"], font=self.parms["font.1"], width=self.parms["ent.0"] )
        self.comp["entry0"].grid(row=0, column=1, padx=5)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # type, lock
        self.comp["frame1"].pack(fill="x", padx=5)
        self.comp["label1"] = tk.Label(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Type")
        self.comp["label1"].grid(row=0, column=0)
        self.comp["combo1"] = tkt.Combobox( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["combo.1"], values=["PDF", "Image", "Source"] )
        self.comp["combo1"].set("PDF")
        self.comp["combo1"].grid(row=0, column=1, padx=5)
        self.comp["button1"] = tk.Button(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Lock", command=click_lock)
        self.comp["button1"].grid(row=0, column=2, padx=5)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # save
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["label2"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="Name")
        self.comp["label2"].grid(row=0, column=0)
        self.comp["entry2"] = tk.Entry( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"] )
        self.comp["entry2"].grid(row=0, column=1, padx=5)
        self.comp["button2"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text=" Save ", bg=self.parms["color"], command=save_url)
        self.comp["button2"].grid(row=0, column=2, padx=5)

class yt_dwn:
    def __init__(self, iswin, ff):
        self.mwin, self.comp, self.parms = None, dict(), {"color":"wheat1"}
        for i in os.listdir(ff):
            if "yt-dlp" in i:
                self.path_ff = os.path.abspath(ff + i)
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "400x250+200+100", ("맑은 고딕", 10), ("Consolas", 12)
            self.parms["ent.0"], self.parms["w"], self.parms["h"], self.parms["ent.1"] = 30, 35, 6, 5
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "500x320+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["w"], self.parms["h"], self.parms["ent.1"] = 32, 35, 5, 5
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def get_info(): # get video info
            time.sleep(0.1)
            order = f"{self.path_ff} -F {self.comp["entry0"].get()}"
            result = subprocess.run(order, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True, text=True, encoding="utf-8")
            txt, out = (result.stdout + result.stderr).split("\n"), [ ]
            for i in txt:
                i = [ x for x in filter( lambda x: x != "", i.split(" ") ) ]
                try:
                    t, temp = int(i[0]), f"{i[0]} {i[1]} {i[2]}"
                    for j in i:
                        if "k" in j:
                            temp = temp + " " + j
                            break
                    out.append(temp)
                except:
                    pass
            self.comp["text0"].delete("1.0", tk.END)
            self.comp["text0"].insert( tk.END, "\n".join(out) )

        def get_vod(): # get video file
            time.sleep(0.1)
            order, cv, ca = f"{self.path_ff} -P ./_ST5_DATA", self.comp["entry1a"].get(), self.comp["entry1b"].get()
            if cv != "" and ca != "":
                order = order + f" -f {cv}+{ca}"
            order = order + f' {self.comp["entry0"].get()}'
            result = subprocess.run(order, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True, text=True, encoding="utf-8")
            with open("./_ST5_DATA/log.txt", "w", encoding="utf-8") as f:
                f.write(f"[order]\n{order}\n{time.strftime("%Y.%m.%d;%H:%M:%S",time.localtime(time.time()))}\n")
                f.write(f"[stdout]\n{result.stdout}\n")
                f.write(f"[stderr]\n{result.stderr}\n")

        self.mwin = tk.Tk() # main window
        self.mwin.title("Kmacro YT Download")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # url
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["label0"] = tk.Label(self.comp["frame0"], font=self.parms["font.0"], bg=self.parms["color"], text="URL")
        self.comp["label0"].grid(row=0, column=0)
        self.comp["entry0"] = tk.Entry( self.comp["frame0"], font=self.parms["font.1"], width=self.parms["ent.0"] )
        self.comp["entry0"].grid(row=0, column=1, padx=5)
        self.comp["text0"] = tk.Text( self.mwin, width=self.parms["w"],  height=self.parms["h"], font=self.parms["font.1"] )
        self.comp["text0"].pack(side="top")

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # video, audio, info, save
        self.comp["frame1"].pack(fill="x", padx=5, pady=5)
        self.comp["label1a"] = tk.Label(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Video")
        self.comp["label1a"].grid(row=0, column=0)
        self.comp["entry1a"] = tk.Entry( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.1"] )
        self.comp["entry1a"].grid(row=0, column=1, padx=5)
        self.comp["label1b"] = tk.Label(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Audio")
        self.comp["label1b"].grid(row=0, column=2)
        self.comp["entry1b"] = tk.Entry( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.1"] )
        self.comp["entry1b"].grid(row=0, column=3, padx=5)
        self.comp["button1a"] = tk.Button(self.comp["frame1"], font=self.parms["font.0"], text=" Info ", bg=self.parms["color"], command=get_info)
        self.comp["button1a"].grid(row=0, column=4, padx=5)
        self.comp["button1b"] = tk.Button(self.comp["frame1"], font=self.parms["font.0"], text=" Save ", bg=self.parms["color"], command=get_vod)
        self.comp["button1b"].grid(row=0, column=5, padx=5)

class yt_info:
    def __init__(self, iswin, ff):
        self.mwin, self.comp, self.parms = None, dict(), {"color":"wheat1"}
        for i in os.listdir(ff):
            if "yt-dlp" in i:
                self.path_ff = os.path.abspath(ff + i)
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "400x250+200+100", ("맑은 고딕", 10), ("Consolas", 12)
            self.parms["ent.0"], self.parms["w"], self.parms["h"] = 25, 35, 8
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "500x320+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["w"], self.parms["h"] = 25, 35, 6
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def get_info(): # get video info
            time.sleep(0.1)
            order = f"{self.path_ff} -F {self.comp["entry0"].get()}"
            result = subprocess.run(order, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True, text=True, encoding="utf-8")
            txt = f"[order]\n{order}\n{time.strftime("%Y.%m.%d;%H:%M:%S",time.localtime(time.time()))}\n"
            txt = txt + f"[stdout]\n{result.stdout}\n[stderr]\n{result.stderr}\n"
            self.comp["text0"].delete("1.0", tk.END)
            self.comp["text0"].insert(tk.END, txt)

        self.mwin = tk.Tk() # main window
        self.mwin.title("Kmacro YT Info")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # url
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["label0"] = tk.Label(self.comp["frame0"], font=self.parms["font.0"], bg=self.parms["color"], text="URL")
        self.comp["label0"].grid(row=0, column=0)
        self.comp["entry0"] = tk.Entry( self.comp["frame0"], font=self.parms["font.1"], width=self.parms["ent.0"] )
        self.comp["entry0"].grid(row=0, column=1, padx=5)
        self.comp["button0"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], text=" Info ", bg=self.parms["color"], command=get_info)
        self.comp["button0"].grid(row=0, column=2, padx=5)
        self.comp["text0"] = tk.Text( self.mwin, width=self.parms["w"],  height=self.parms["h"], font=self.parms["font.1"] )
        self.comp["text0"].pack(side="top")

def ks_sel(iswin, mode, infof): # select kscript code/bin -> (path, stack, errhlt)
    def sel_file(): # select file
        time.sleep(0.1)
        nonlocal path
        nonlocal mode
        nonlocal comp
        if mode == "txt":
            path = tkf.askopenfile(title="Select Kscript", filetypes=[ ("Source Code", "*.txt"), ("All Files", "*.*") ], initialdir="./_ST5_DATA/").name
        elif mode == "webp":
            path = tkf.askopenfile(title="Select Kscript", filetypes=[ ("Binary", "*.webp"), ("All Files", "*.*") ], initialdir="./_ST5_DATA/").name
        comp["strvar0"].set(path)
        comp["strvar2"].set( infof(path) )

    def run_file(): # exit & return
        time.sleep(0.1)
        nonlocal mwin
        nonlocal stack
        nonlocal comp
        try:
            stack = int( comp["entry1"].get() )
        except:
            stack = -1
        mwin.destroy()

    def click_hlt(): # errhlt click
        time.sleep(0.1)
        nonlocal errhlt
        nonlocal comp
        nonlocal parms
        errhlt = not errhlt
        if errhlt:
            comp["button1a"].configure( bg=parms["activate"] )
        else:
            comp["button1a"].configure( bg=parms["color"] )

    comp, parms, path, stack, errhlt = dict(), {"color":"wheat1", "activate":"green yellow"}, "", -1, True
    if iswin: # windows
        parms["mwin.size"], parms["font.0"], parms["font.1"] = "360x150+200+100", ("맑은 고딕", 10), ("Consolas", 12)
        parms["ent.0"], parms["ent.1"], parms["ent.2"] = 25, 8, 30
    else: # linux
        parms["mwin.size"], parms["font.0"], parms["font.1"] = "450x180+200+100", ("맑은 고딕", 8), ("Consolas", 10)
        parms["ent.0"], parms["ent.1"], parms["ent.2"] = 27, 8, 32

    mwin = tk.Tk() # main window
    mwin.title("Kmacro Kscript")
    mwin.geometry( parms["mwin.size"] )
    mwin.resizable(False, False)
    mwin.configure( bg=parms["color"] )

    comp["frame0"] = tk.Frame( mwin, bg=parms["color"] ) # file selector
    comp["frame0"].pack(fill="x", padx=5, pady=5)
    comp["button0"] = tk.Button(comp["frame0"], bg=parms["color"], font=parms["font.0"], text=" . . . ", command=sel_file)
    comp["button0"].grid(row=0, column=0)
    comp["strvar0"] = tk.StringVar()
    comp["strvar0"].set("Select Kscript File")
    comp["entry0"] = tk.Entry(comp["frame0"], font=parms["font.1"], textvariable=comp["strvar0"], width=parms["ent.0"], state="readonly")
    comp["entry0"].grid(row=0, column=1, padx=5)

    comp["strvar2"] = tk.StringVar()
    comp["strvar2"].set("No Info")
    comp["entry2"] = tk.Entry(mwin, font=parms["font.1"], textvariable=comp["strvar2"], width=parms["ent.2"], state="readonly")
    comp["entry2"].pack(padx=5, pady=5)

    comp["frame1"] = tk.Frame( mwin, bg=parms["color"] ) # stack, errhlt, run
    comp["frame1"].pack(fill="x", padx=5, pady=5)
    comp["label1"] = tk.Label( comp["frame1"], font=parms["font.0"], text="Custom Stack", bg=parms["color"] )
    comp["label1"].grid(row=0, column=0)
    comp["entry1"] = tk.Entry( comp["frame1"], font=parms["font.1"], width=parms["ent.1"] )
    comp["entry1"].grid(row=0, column=1, padx=5)
    comp["button1a"] = tk.Button(comp["frame1"], font=parms["font.0"], text="ErrHlt", bg=parms["activate"], command=click_hlt)
    comp["button1a"].grid(row=0, column=2, padx=5)
    comp["button1b"] = tk.Button(comp["frame1"], font=parms["font.0"], text="Run", bg=parms["color"], command=run_file)
    comp["button1b"].grid(row=0, column=3, padx=5)

    mwin.mainloop()
    return path, stack, errhlt

class ks_console:
    def __init__(self, iswin):
        self.mwin, self.comp, self.parms, self.status_run, self.status_lock, self.io = None, dict(), {"color":"wheat1", "activate":"green yellow"}, False, False, ("none", "")
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "600x400+200+100", ("맑은 고딕", 10), ("Consolas", 12)
            self.parms["ent.0"], self.parms["w"], self.parms["h"], self.parms["ent.2"] = 43, 49, 11, 41
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "700x500+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["w"], self.parms["h"], self.parms["ent.2"] = 40, 49, 9, 38

    def start(self):
        def gui_exit(): # close the screen
            time.sleep(0.1)
            self.mwin.destroy()
            self.mwin = None

        def click_run(): # run on/off
            self.status_run = not self.status_run
            if self.status_run:
                self.comp["button0a"].configure( bg=self.parms["activate"] )
            else:
                self.comp["button0a"].configure( bg=self.parms["color"] )
            time.sleep(0.1)

        def click_lock(): # lock on/off
            time.sleep(0.1)
            self.status_lock = not self.status_lock
            if self.status_lock:
                self.comp["button0b"].configure( bg=self.parms["activate"] )
                self.mwin.wm_attributes("-topmost", 1)
            else:
                self.comp["button0b"].configure( bg=self.parms["color"] )
                self.mwin.wm_attributes("-topmost", 0)

        def click_clear(): # clear log
            time.sleep(0.1)
            self.comp["list1"].delete( 0, self.comp["list1"].size() )

        def press_enter(event): # data input
            temp = self.comp["entry2"].get()
            self.setout(temp+"\n")
            self.io = ("ret", temp)
            time.sleep(0.1)
        def click_enter():
            press_enter(None)

        self.mwin = tk.Tk() # main window
        self.mwin.title("Kmacro Kscript")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )
        self.mwin.protocol('WM_DELETE_WINDOW', gui_exit)

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # run, lock, msg
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["button0a"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], text="Run", bg=self.parms["color"], command=click_run)
        self.comp["button0a"].grid(row=0, column=0)
        self.comp["button0b"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], text="Lock", bg=self.parms["color"], command=click_lock)
        self.comp["button0b"].grid(row=0, column=1)
        self.comp["strvar0"] = tk.StringVar()
        self.comp["strvar0"].set("No Data")
        self.comp["entry0"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], textvariable=self.comp["strvar0"], width=self.parms["ent.0"], state="readonly")
        self.comp["entry0"].grid(row=0, column=2, padx=5)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # console log
        self.comp["frame1"].pack(padx=5, pady=5)
        self.comp["list1"] = tk.Listbox( self.comp["frame1"], width=self.parms["w"],  height=self.parms["h"], font=self.parms["font.1"] )
        self.comp["list1"].pack(side="left", fill="y")
        self.comp["scroll1"] = tk.Scrollbar(self.comp["frame1"], orient="vertical")
        self.comp["scroll1"].config(command=self.comp["list1"].yview)
        self.comp["scroll1"].pack(side="right", fill="y")
        self.comp["list1"].config(yscrollcommand=self.comp["scroll1"].set)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # clear, input
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["button2a"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text="Clear", bg=self.parms["color"], command=click_clear)
        self.comp["button2a"].grid(row=0, column=0)
        self.comp["button2b"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text="Enter", bg=self.parms["color"], command=click_enter)
        self.comp["button2b"].grid(row=0, column=1)
        self.comp["entry2"] = tk.Entry( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"] )
        self.comp["entry2"].grid(row=0, column=2, padx=5)
        self.comp["entry2"].bind("<Return>", press_enter)

    def guiloop(self): # io status "none" "print" "msg"
        while self.mwin != None:
            if self.io[0] == "print":
                self.setout( self.io[1] )
                self.io = ("none", "")
            elif self.io[0] == "msg":
                self.setmsg( self.io[1] )
                self.io = ("none", "")
            else:
                time.sleep(0.01)
                self.mwin.update()

    def setmsg(self, text): # set msgbox
        self.comp["strvar0"].set(text)
        self.mwin.update()

    def setout(self, text): # set output listbox
        idx = self.comp["list1"].size() - 1
        if idx >= 0:
            text = self.comp["list1"].get(idx) + text
            self.comp["list1"].delete(idx)
        text = text.split("\n")
        for i in text:
            self.comp["list1"].insert(self.comp["list1"].size(), i)
        self.comp["list1"].see(self.comp["list1"].size()-1)

    def testio(self, mode, v): # runtime.testio impl
        if mode == 16:
            self.stdout( runtime.tostr(v[0]) )
            return self.stdin()
        elif mode == 17:
            self.stdout( runtime.tostr(v[0]) )
            return None
        else:
            return runtime.testio(mode, v)
        
    def stdin(self):
        while self.io[0] != "ret":
            time.sleep(0.01)
        temp, self.io = self.io[1], ("none", "")
        return temp

    def stdout(self, word):
        self.io = ("print", word)
        while self.io[0] != "none":
            time.sleep(0.01)

    def stderr(self, word):
        self.io = ("print", "[err] " + word)
        while self.io[0] != "none":
            time.sleep(0.01)

    def stdmsg(self, word):
        self.io = ("msg", word)
        while self.io[0] != "none":
            time.sleep(0.01)

class mainclass(blocksel.toolbox):
    def __init__(self, iswin):
        super().__init__("Kmacro", iswin)
        path, self.selection = "../../_ST5_COMMON/iconpack/kmacro/", -1
        self.txts, self.curpos, self.upos, self.umsg = ["Capture", "KSrun_txt", "KSrun_bin", "WebDwn", "YTubeDwn", "YTubeInfo"], 1, 0, ["Select Mode"]
        self.pics = [path+"capture.png", path+"ks_txt.png", path+"ks_exe.png", path+"web_dwn.png", path+"yt_dwn.png", path+"yt_info.png"]

    def custom0(self, x):
        self.selection = x
        self.mwin.destroy()

def getdrv(path): # get driver
    options = webdriver.EdgeOptions()
    options.add_argument('headless')
    options.add_argument("disable-gpu")
    options.add_argument('window-size=1920x1080')
    options.add_argument("user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/177.0.0.0 Safari/537.36 Edg/177.0.1938.62")
    service = Service(path)
    service.creation_flags = 0x08000000
    driver = webdriver.Edge(options=options,service=service)
    return driver

def getinfo(path): # get program info
    try:
        if "." in path and path[path.rfind("."):] == ".txt":
            with open(path, "r", encoding="utf-8") as f:
                lines = f.readlines()
            info = "No Info"
            for i in lines:
                temp = i.replace(" ", "")
                if len(temp) > 3 and temp[0] == "#":
                    info = i[i.find("#")+1:].strip()
                    break
        else:
            temp = runtime.kvm()
            info, _, _ = temp.view(path)
        return info
    except Exception as e:
        return f"{e}"

def compile(iswin, path): # compile kscript
    dll = ctypes.CDLL("../../_ST5_COMMON/kscriptc/kscriptc.dll") if iswin else ctypes.cdll.LoadLibrary("../../_ST5_COMMON/kscriptc/kscriptc.so")
    dll.func0.argtypes, dll.func0.restype = kobj.call("b", "") # free
    dll.func1.argtypes, dll.func1.restype = kobj.call("iiiii", "") # init c
    dll.func2.argtypes, dll.func2.restype = kobj.call("bii", "") # set data
    dll.func3.argtypes, dll.func3.restype = kobj.call("", "b") # get data
    dll.func4.argtypes, dll.func4.restype = kobj.call("", "f") # get tpass
    dll.func5.argtypes, dll.func5.restype = kobj.call("bi", "") # addpkg
    dll.func6.argtypes, dll.func6.restype = kobj.call("", "b") # compile
    dll.func1(1, 1, 1, 0, 0)
    with open(path, "r", encoding="utf-8") as f:
        o0, o1 = kobj.send( bytes(f.read(), encoding="utf-8") )
        dll.func2(o0, o1, 0)
        o0, o1 = kobj.send( bytes("실시간 컴파일된 바이너리", encoding="utf-8") )
        dll.func2(o0, o1, 1)
    for i in os.listdir("../../_ST5_COMMON/kscriptc/_ST5_DATA/"):
        o0, o1 = kobj.send( bytes("../../_ST5_COMMON/kscriptc/_ST5_DATA/"+i, encoding="utf-8") )
        dll.func5(o0, o1)
    t = dll.func6()
    err = str(kobj.recvauto(t), encoding="utf-8")
    dll.func0(t)
    t = dll.func3(-1)
    bin = kobj.recvauto(t)
    dll.func0(t)
    if err == "":
        path = f"./_ST5_DATA/{int(time.time())}.bin"
        with open(path, "wb") as f:
            f.write(bin)
        return path
    else:
        tkm.showerror(title="Compile Error", message=f" {err} ")
        return ""

def run(iswin, path, stack, errhlt): # run kscript binary
    rt, vm, cfg, num, public, phash = ks_console(iswin), runtime.kvm(), kdb.toolbox(), 0, [ ], [ ]
    rt.start()
    rt.setmsg(path)
    try:
        info, abi, pub = vm.view(path)
        if abi > 79:
            raise Exception("invalid ABI")
        with open("../../_ST5_SIGN.txt", "r", encoding="utf-8") as f:
            cfg.read( f.read() )
        while f"{num}.name" in cfg.name:
            public.append( cfg.get(f"{num}.public")[3] )
            phash.append( cfg.get(f"{num}.phash")[3] )
            num = num + 1
    except Exception as e:
        rt.setmsg(f"[load error] {e}")
        rt.mwin.mainloop()
        return
    try:
        vm.load(True)
        if pub == "":
            raise Exception("no sign exists")
        if pub not in public:
            raise Exception("untrusted sign")
        info = info + " phash " + phash[ public.index(pub) ].hex()
    except Exception as e:
        tkm.showinfo(title="Sign Warning", message=f" {e} ")
    rt.setmsg(info)
    if stack > 1:
        vm.maxstk = stack
    vm.runone, mc, lb = False, kscript_macro.lib(True, True, True, True), kscript_lib.lib(True, True, True, iswin)
    mc.drvpath = os.path.abspath("../../_ST5_COMMON/edgewd/msedgedriver.exe") if iswin else os.path.abspath("../../_ST5_COMMON/edgewd/msedgedriver")
    for i in os.listdir("../../_ST5_COMMON/videopack/"):
        if "wkhtmltopdf" in i:
            mc.kitpath = os.path.abspath("../../_ST5_COMMON/videopack/" + i)
    cfg = kdb.toolbox()
    with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
        cfg.read( f.read() )
    lb.io, lb.p_desktop, lb.p_local, lb.p_starter = rt, cfg.get("path.desktop")[3], cfg.get("path.local")[3], cfg.get("path.export")[3]
    lb.p_base = os.path.abspath("./").replace("\\", "/")
    if lb.p_base[-1] != "/":
        lb.p_base = lb.p_base + "/"
    thr = threading.Thread( target=run_sub, args=(rt, vm, mc, lb, errhlt) )
    thr.start()
    rt.guiloop()

def run_sub(rt, vm, mc, lb, errhlt):
    while rt.mwin != None:
        if rt.status_run:
            intr = vm.run()
            if intr >= 500: # kscript macro
                vm.ma, vm.errmsg = mc.run(vm.callmem, intr)
                if vm.errmsg != None and errhlt:
                    rt.stdmsg(f"[err] {vm.errmsg}")
                    break
                elif vm.errmsg != None and not errhlt:
                    rt.stdout(f"[err] {vm.errmsg}")
            elif intr >= 32: # kscript lib
                vm.ma, vm.errmsg = lb.run(vm.callmem, intr)
                if vm.errmsg != None and errhlt:
                    rt.stdmsg(f"[err] {vm.errmsg}")
                    break
                elif vm.errmsg != None and not errhlt:
                    rt.stdout(f"[err] {vm.errmsg}")
            elif intr >= 16: # testio
                vm.ma = rt.testio(intr, vm.callmem)
            elif intr >= 0: # kscript
                if intr == 1:
                    rt.stdmsg("[msg] program end")
                    break
                elif intr == 2:
                    if errhlt:
                        rt.stdmsg(f"[err] {vm.errmsg}")
                        break
                    else:
                        rt.stdout(f"[err] {vm.errmsg}")
            else: # critical
                rt.stdmsg(f"[critical] {vm.errmsg}")
        else:
            time.sleep(0.1)

kobj.repath()
cfg = kdb.toolbox()
with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
    cfg.read( f.read() )
iswin, flag = (cfg.get("dev.os")[3] == "windows"), True
if not os.path.exists("../../_ST5_COMMON/iconpack/"):
    flag = False
    tkm.showinfo(title="No Package", message=" Convs requires package >>> common.iconpack <<<. \n Install dependent package and start again. ")
if not os.path.exists("../../_ST5_COMMON/videopack/"):
    flag = False
    tkm.showinfo(title="No Package", message=" Convs requires package >>> common.videopack <<<. \n Install dependent package and start again. ")
if not os.path.exists("../../_ST5_COMMON/edgewd/"):
    flag = False
    tkm.showinfo(title="No Package", message=" Convs requires package >>> common.edgewd <<<. \n Install dependent package and start again. ")
if not os.path.exists("../../_ST5_COMMON/kscriptc/"):
    flag = False
    tkm.showinfo(title="No Package", message=" Convs requires package >>> common.kscriptc <<<. \n Install dependent package and start again. ")
if not os.path.exists("./_ST5_DATA/"):
    os.mkdir("./_ST5_DATA/")
if flag:
    worker = mainclass(iswin)
    worker.entry()
    worker.guiloop()
    if worker.selection == 0:
        t = capture(iswin, "../../_ST5_COMMON/videopack/")
    elif worker.selection == 1:
        path, stack, errhlt = ks_sel(iswin, "txt", getinfo)
        path = compile(iswin, path)
        if path != "":
            run(iswin, path, stack, errhlt)
            os.remove(path)
    elif worker.selection == 2:
        path, stack, errhlt = ks_sel(iswin, "webp", getinfo)
        if path != "":
            run(iswin, path, stack, errhlt)
    elif worker.selection == 3:
        t = web_dwn(iswin, "../../_ST5_COMMON/edgewd/")
    elif worker.selection == 4:
        t = yt_dwn(iswin, "../../_ST5_COMMON/videopack/")
    elif worker.selection == 5:
        t = yt_info(iswin, "../../_ST5_COMMON/videopack/")
time.sleep(0.5)
