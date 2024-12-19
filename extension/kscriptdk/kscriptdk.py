# test736 : extension.kscriptdk
# 윈도우 빌드 추가 : --hidden-import "pynput.keyboard._win32" --hidden-import "pynput.mouse._win32"
# 리눅스 빌드 추가 : --hidden-import "pynput.keyboard._xorg" --hidden-import "pynput.mouse._xorg"

import os
import time
import ctypes

import tkinter as tk
import tkinter.filedialog as tkf
import tkinter.messagebox as tkm

import pyautogui
import kobj
import kdb
import kcom
import kaes
import runtime
import kscript_lib
import kscript_macro

class ks_cplr:
    def __init__(self, iswin, desktop):
        self.mwin, self.comp, self.parms, self.desktop = None, dict(), {"color":"DarkSeaGreen1", "activate":"cyan"}, desktop
        self.source = ["", "", ""] # txtpath, txtdata, exitmode
        self.compile = ["", "", "", -1] # info, iconpath, signdata, stack
        self.status = [True, True, False, False, False, True, True, True, True, True, True] # optconst, optasm, viewtoken, viewast, viewasm, testfunc, stdlib, stdio, osfs, macro, errhlt
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "800x600+200+100", ("맑은 고딕", 10), ("Consolas", 12)
            self.parms["ent.0"], self.parms["w"], self.parms["h"], self.parms["ent.2"], self.parms["ent.5"] = 52, 70, 14, 28, 12
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "1080x770+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["w"], self.parms["h"], self.parms["ent.2"], self.parms["ent.5"] = 63, 81, 12, 33, 12
        self.entry()
        self.guiloop()

    def entry(self):
        def sel_source(): # select source code
            time.sleep(0.1)
            self.source[0] = tkf.askopenfile(title="Select Code File", filetypes=[ ("Text Files", "*.txt"), ("All Files", "*.*") ], initialdir=self.desktop).name
            with open(self.source[0], "r", encoding="utf-8") as f:
                self.source[1] = f.read()
            self.comp["text1"].delete("1.0", tk.END)
            self.comp["text1"].insert( tk.END, self.source[1] )
            self.comp["strvar0a"].set( self.source[0] )

        def sel_icon(): # select icon file
            time.sleep(0.1)
            try:
                self.compile[1] = tkf.askopenfile(title="Select Code File", filetypes=[ ("WEBP Files", "*.webp"), ("PNG Files", "*.png"), ("JPG Files", "*.jpg"), ("All Files", "*.*") ], initialdir=self.desktop).name
                self.comp["strvar2"].set( self.compile[1] )
            except Exception as e:
                self.compile[1] = ""
                self.comp["strvar2"].set(f"{e}")

        def sel_sign_direct(): # select sign direct
            time.sleep(0.1)
            try:
                path = tkf.askopenfile(title="Select Sign Private", filetypes=[ ("Text Files", "*.txt"), ("All Files", "*.*") ], initialdir=self.desktop).name
                with open(path, "r", encoding="utf-8") as f:
                    data = f.read()
                self.comp["strvar3"].set(path)
                self.compile[2] = data
            except Exception as e:
                self.comp["strvar3"].set(f"{e}")
                self.compile[2] = ""

        def sel_sign_comm(): # select sign comm
            time.sleep(0.1)
            try:
                sport, skey = kcom.unpack( self.comp["entry3b"].get() )
                node, wk = kcom.node(), kaes.funcmode()
                node.port = sport
                data = node.recieve(skey)
                path = str(data[48:], encoding="utf-8")
                with open(path, "rb") as f:
                    tgt = f.read()
                wk.before = tgt
                wk.decrypt( data[0:48] )
                data = str(wk.after, encoding="utf-8")
                wk.before, wk.after = "", ""
                self.comp["strvar3"].set(path)
                self.compile[2] = data
            except Exception as e:
                self.comp["strvar3"].set(f"{e}")
                self.compile[2] = ""

        def click_but(x): # option click (x 0~10)
            time.sleep(0.1)
            self.status[x] = not self.status[x]
            if self.status[x]:
                self.comp[f"button4.{x}"].configure( bg=self.parms["activate"] )
            else:
                self.comp[f"button4.{x}"].configure( bg=self.parms["color"] )

        def click_save(): # save source
            time.sleep(0.1)
            temp = self.comp["text1"].get("1.0", tk.END)
            with open(self.source[0], "w", encoding="utf-8") as f:
                f.write( temp[:len(temp)-1] )
            tkm.showinfo(title="Text Saved", message=f" Source code saved at \n {self.source[0]} ")

        def click_build(): # build source
            time.sleep(0.1)
            self.source[1], self.source[2], self.compile[0] = self.comp["text1"].get("1.0", tk.END), "build", self.comp["entry2a"].get()
            self.mwin.destroy()
            self.mwin = None

        def click_run(): # run source
            time.sleep(0.1)
            self.source[1], self.source[2], self.compile[0], temp = self.comp["text1"].get("1.0", tk.END), "run", self.comp["entry2a"].get(), self.comp["entry5"].get()
            self.compile[3] = -1 if temp == "" else int(temp)
            self.mwin.destroy()
            self.mwin = None

        def gui_exit(): # close the screen
            time.sleep(0.1)
            self.source[2] = ""
            self.mwin.destroy()
            self.mwin = None

        self.mwin = tk.Tk() # main window
        self.mwin.title("KscriptDK Compiler")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )
        self.mwin.protocol('WM_DELETE_WINDOW', gui_exit)

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # source, mptr
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["button0"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], bg=self.parms["color"], text=" . . . ", command=sel_source)
        self.comp["button0"].grid(row=0, column=0)
        self.comp["strvar0a"] = tk.StringVar()
        self.comp["strvar0a"].set("Select Kscript Source")
        self.comp["entry0"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], width=self.parms["ent.0"], textvariable=self.comp["strvar0a"], state="readonly")
        self.comp["entry0"].grid(row=0, column=1, padx=5)
        self.comp["strvar0b"] = tk.StringVar()
        self.comp["strvar0b"].set("x00000 y00000")
        self.comp["label0"] = tk.Label( self.comp["frame0"], font=self.parms["font.1"], bg=self.parms["color"], textvariable=self.comp["strvar0b"] )
        self.comp["label0"].grid(row=0, column=2)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # code text
        self.comp["frame1"].pack(fill="x", padx=5, pady=5)
        self.comp["text1"] = tk.Text( self.comp["frame1"], width=self.parms["w"],  height=self.parms["h"], font=self.parms["font.1"] )
        self.comp["text1"].pack(side="left", fill="y")
        self.comp["scroll1"] = tk.Scrollbar(self.comp["frame1"], orient="vertical")
        self.comp["scroll1"].config(command=self.comp["text1"].yview)
        self.comp["scroll1"].pack(side="right", fill="y")
        self.comp["text1"].config(yscrollcommand=self.comp["scroll1"].set)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # info, icon
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["label2a"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="Info")
        self.comp["label2a"].grid(row=0, column=0)
        self.comp["entry2a"] = tk.Entry( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"] )
        self.comp["entry2a"].grid(row=0, column=1, padx=5)
        self.comp["label2b"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="Icon")
        self.comp["label2b"].grid(row=0, column=2, padx=5)
        self.comp["button2"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text=" . . . ", command=sel_icon)
        self.comp["button2"].grid(row=0, column=3)
        self.comp["strvar2"] = tk.StringVar()
        self.comp["strvar2"].set("Select Binary Icon")
        self.comp["entry2b"] = tk.Entry(self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"], textvariable=self.comp["strvar2"], state="readonly")
        self.comp["entry2b"].grid(row=0, column=4, padx=5)

        self.comp["frame3"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # sign
        self.comp["frame3"].pack(fill="x", padx=5, pady=5)
        self.comp["label3"] = tk.Label(self.comp["frame3"], font=self.parms["font.0"], bg=self.parms["color"], text="Sign")
        self.comp["label3"].grid(row=0, column=0)
        self.comp["button3a"] = tk.Button(self.comp["frame3"], font=self.parms["font.0"], bg=self.parms["color"], text=" . . . ", command=sel_sign_direct)
        self.comp["button3a"].grid(row=0, column=1)
        self.comp["strvar3"] = tk.StringVar()
        self.comp["strvar3"].set("")
        self.comp["entry3a"] = tk.Entry(self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.2"], textvariable=self.comp["strvar3"], state="readonly")
        self.comp["entry3a"].grid(row=0, column=2, padx=5)
        self.comp["button3b"] = tk.Button(self.comp["frame3"], font=self.parms["font.0"], bg=self.parms["color"], text=" . . . ", command=sel_sign_comm)
        self.comp["button3b"].grid(row=0, column=3)
        self.comp["entry3b"] = tk.Entry( self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.2"] )
        self.comp["entry3b"].grid(row=0, column=4, padx=5)

        optname = ["OptConst", "OptAsm", "ViewTKN", "ViewAST", "ViewASM", "testfunc", " stdlib ", " stdio ", "  osfs  ", " macro "]
        self.comp["frame4"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # options
        self.comp["frame4"].pack(fill="x", padx=5, pady=5)
        for i in range(0, 10):
            self.comp[f"button4.{i}"] = tk.Button( self.comp["frame4"], font=self.parms["font.0"], text=optname[i], command=lambda i=i: click_but(i) )
            self.comp[f"button4.{i}"].grid(row=0, column=i)
            if self.status[i]:
                self.comp[f"button4.{i}"].configure( bg=self.parms["activate"] )
            else:
                self.comp[f"button4.{i}"].configure( bg=self.parms["color"] )

        self.comp["frame5"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # stack, errhlt, save, build, run
        self.comp["frame5"].pack(fill="x", padx=5, pady=5)
        self.comp["label5"] = tk.Label(self.comp["frame5"], font=self.parms["font.0"], bg=self.parms["color"], text="Custom Stack")
        self.comp["label5"].grid(row=0, column=0)
        self.comp["entry5"] = tk.Entry( self.comp["frame5"], font=self.parms["font.1"], width=self.parms["ent.5"] )
        self.comp["entry5"].grid(row=0, column=1, padx=5)
        self.comp["button4.10"] = tk.Button( self.comp["frame5"], font=self.parms["font.0"], bg=self.parms["activate"], text="  ErrHlt  ", command=lambda: click_but(10) )
        self.comp["button4.10"].grid(row=0, column=2, padx=10)
        self.comp["button5b"] = tk.Button(self.comp["frame5"], font=self.parms["font.0"], bg=self.parms["color"], text="     Save     ", command=click_save)
        self.comp["button5b"].grid(row=0, column=3)
        self.comp["button5c"] = tk.Button(self.comp["frame5"], font=self.parms["font.0"], bg=self.parms["color"], text="     Build     ", command=click_build)
        self.comp["button5c"].grid(row=0, column=4)
        self.comp["button5d"] = tk.Button(self.comp["frame5"], font=self.parms["font.0"], bg=self.parms["color"], text="     Run     ", command=click_run)
        self.comp["button5d"].grid(row=0, column=5)

    def guiloop(self):
        while self.mwin != None:
            temp = pyautogui.position()
            self.comp["strvar0b"].set(f"x{temp[0]} y{temp[1]}")
            self.mwin.update()
            time.sleep(0.05)

class ks_dbg:
    def __init__(self, iswin, rt):
        self.mwin, self.comp, self.parms, self.stack, self.rt = None, dict(), {"color":"DarkSeaGreen1", "activate":"cyan", "run":"hot pink"}, [0], rt
        self.status = [False, False, True, False] # lock, stack value, reg value, io waiting
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "800x600+200+100", ("맑은 고딕", 10), ("Consolas", 12)
            self.parms["e"], self.parms["w"], self.parms["h"], self.parms["ent.0"], self.parms["ent.2a"], self.parms["ent.2b"] = 30, 30, 18, 22, 8, 18
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "1080x770+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["e"], self.parms["w"], self.parms["h"], self.parms["ent.0"], self.parms["ent.2a"], self.parms["ent.2b"] = 37, 37, 14, 25, 8, 22
        self.entry()

    def entry(self):
        def click_clear(): # clear log
            time.sleep(0.1)
            self.comp["list_r"].delete( 0, self.comp["list_r"].size() )

        def click_run(): # run 1 step
            if not self.status[3]:
                self.rt.run()
                self.setreg(self.rt.vm.pc, self.rt.vm.sp, self.rt.vm.ma, self.rt.vm.mb)
                self.setstk(self.rt.vm.stack)
            time.sleep(0.02)

        def click_status(x): # change status (x 0~2)
            time.sleep(0.1)
            self.status[x] = not self.status[x]
            if self.status[x]:
                self.comp[f"button1.{x}"].configure( bg=self.parms["activate"] )
            else:
                self.comp[f"button1.{x}"].configure( bg=self.parms["color"] )
            if x == 0:
                if self.status[x]:
                    self.mwin.wm_attributes("-topmost", 1)
                else:
                    self.mwin.wm_attributes("-topmost", 0)
            elif x == 1:
                self.setstk(self.stack)

        def press_enter(event): # data input
            if self.status[3]:
                self.rt.vm.ma, self.status[3] = self.comp["entry0"].get(), False
                self.stdout(self.rt.vm.ma + "\n")
                click_run()
        def click_enter():
            press_enter(None)

        def sel_stack(event): # select stack line
            time.sleep(0.1)
            self.selstk()

        self.mwin = tk.Tk() # main window, left/right frame
        self.mwin.title("KscriptDK Debugger")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )
        self.comp["l"] = tk.Frame( self.mwin, bg=self.parms["color"] )
        self.comp["l"].pack(side="left", fill=tk.BOTH, expand=True)
        self.comp["r"] = tk.Frame( self.mwin, bg=self.parms["color"] )
        self.comp["r"].pack(side="right", fill=tk.BOTH, expand=True)

        self.comp["strvar_l"] = tk.StringVar() # msg, log
        self.comp["strvar_l"].set("No Message")
        self.comp["entry_l"] = tk.Entry(self.comp["l"], font=self.parms["font.1"], textvariable=self.comp["strvar_l"], width=self.parms["e"], state="readonly")
        self.comp["entry_l"].pack(side="top", padx=5, pady=5)
        self.comp["frame_l"] = tk.Frame( self.comp["l"], bg=self.parms["color"] )
        self.comp["frame_l"].pack(side="top", padx=5, pady=5)
        self.comp["list_l"] = tk.Listbox( self.comp["frame_l"], width=self.parms["w"],  height=self.parms["h"], font=self.parms["font.1"] )
        self.comp["list_l"].pack(side="left", fill="y")
        self.comp["scroll_l"] = tk.Scrollbar(self.comp["frame_l"], orient="vertical")
        self.comp["scroll_l"].config(command=self.comp["list_l"].yview)
        self.comp["scroll_l"].pack(side="right", fill="y")
        self.comp["list_l"].config(yscrollcommand=self.comp["scroll_l"].set)

        self.comp["strvar_r"] = tk.StringVar() # value, stack
        self.comp["strvar_r"].set("No Data")
        self.comp["entry_r"] = tk.Entry(self.comp["r"], font=self.parms["font.1"], textvariable=self.comp["strvar_r"], width=self.parms["e"], state="readonly")
        self.comp["entry_r"].pack(side="top", padx=5, pady=5)
        self.comp["frame_r"] = tk.Frame( self.comp["r"], bg=self.parms["color"] )
        self.comp["frame_r"].pack(side="top", padx=5, pady=5)
        self.comp["list_r"] = tk.Listbox( self.comp["frame_r"], width=self.parms["w"],  height=self.parms["h"], font=self.parms["font.1"] )
        self.comp["list_r"].pack(side="left", fill="y")
        self.comp["scroll_r"] = tk.Scrollbar(self.comp["frame_r"], orient="vertical")
        self.comp["scroll_r"].config(command=self.comp["list_r"].yview)
        self.comp["scroll_r"].pack(side="right", fill="y")
        self.comp["list_r"].config(yscrollcommand=self.comp["scroll_r"].set)
        self.comp["list_r"].bind("<ButtonRelease-1>", sel_stack)

        self.comp["frame0"] = tk.Frame( self.comp["l"], bg=self.parms["color"] ) # clear, input
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["button0a"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], text="Clear", bg=self.parms["color"], command=click_clear)
        self.comp["button0a"].grid(row=0, column=0)
        self.comp["button0b"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], text="Enter", bg=self.parms["color"], command=click_enter)
        self.comp["button0b"].grid(row=0, column=1)
        self.comp["entry0"] = tk.Entry( self.comp["frame0"], font=self.parms["font.1"], width=self.parms["ent.0"] )
        self.comp["entry0"].grid(row=0, column=2, padx=5)
        self.comp["entry0"].bind("<Return>", press_enter)

        self.comp["frame1"] = tk.Frame( self.comp["l"], bg=self.parms["color"] ) # run, lock, stack, reg, pos
        self.comp["frame1"].pack(fill="x", padx=5, pady=5)
        self.comp["button1"] = tk.Button(self.comp["frame1"], font=self.parms["font.0"], text="Run", bg=self.parms["run"], command=click_run)
        self.comp["button1"].grid(row=0, column=0)
        self.comp["button1.0"] = tk.Button( self.comp["frame1"], font=self.parms["font.0"], text="Lock", bg=self.parms["color"], command=lambda:click_status(0) )
        self.comp["button1.0"].grid(row=0, column=1)
        self.comp["button1.1"] = tk.Button( self.comp["frame1"], font=self.parms["font.0"], text="STK", bg=self.parms["color"], command=lambda:click_status(1) )
        self.comp["button1.1"].grid(row=0, column=2)
        self.comp["button1.2"] = tk.Button( self.comp["frame1"], font=self.parms["font.0"], text="REG", bg=self.parms["activate"], command=lambda:click_status(2) )
        self.comp["button1.2"].grid(row=0, column=3)
        self.comp["strvar1"] = tk.StringVar()
        self.comp["strvar1"].set("STK 00000000")
        self.comp["label1"] = tk.Label(self.comp["frame1"], font=self.parms["font.1"], bg=self.parms["color"], textvariable=self.comp["strvar1"])
        self.comp["label1"].grid(row=0, column=4, padx=5)

        self.comp["frame2"] = tk.Frame( self.comp["r"], bg=self.parms["color"] ) # pc, sp, ma, mb
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["label2a"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="PC")
        self.comp["label2a"].grid(row=0, column=0, padx=5, pady=5)
        self.comp["strvar2a"] = tk.StringVar()
        self.comp["strvar2a"].set("NO PC")
        self.comp["entry2a"] = tk.Entry(self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2a"], textvariable=self.comp["strvar2a"], state="readonly")
        self.comp["entry2a"].grid(row=0, column=1, padx=5, pady=5)
        self.comp["label2b"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="MA")
        self.comp["label2b"].grid(row=0, column=2, padx=5, pady=5)
        self.comp["strvar2b"] = tk.StringVar()
        self.comp["strvar2b"].set("NO MA")
        self.comp["entry2b"] = tk.Entry(self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2b"], textvariable=self.comp["strvar2b"], state="readonly")
        self.comp["entry2b"].grid(row=0, column=3, padx=5, pady=5)
        self.comp["label2c"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="SP")
        self.comp["label2c"].grid(row=1, column=0, padx=5, pady=5)
        self.comp["strvar2c"] = tk.StringVar()
        self.comp["strvar2c"].set("NO SP")
        self.comp["entry2c"] = tk.Entry(self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2a"], textvariable=self.comp["strvar2c"], state="readonly")
        self.comp["entry2c"].grid(row=1, column=1, padx=5, pady=5)
        self.comp["label2d"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="MB")
        self.comp["label2d"].grid(row=1, column=2, padx=5, pady=5)
        self.comp["strvar2d"] = tk.StringVar()
        self.comp["strvar2d"].set("NO MB")
        self.comp["entry2d"] = tk.Entry(self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2b"], textvariable=self.comp["strvar2d"], state="readonly")
        self.comp["entry2d"].grid(row=1, column=3, padx=5, pady=5)

    def selstk(self): # stack selection
        pos = self.comp["list_r"].curselection()
        pos = 0 if len(pos) == 0 else pos[0]
        self.comp["strvar1"].set(f"STK {pos}")
        if self.status[1]:
            self.comp["strvar_r"].set( runtime.tostr( self.stack[pos] ) )
        else:
            self.comp["strvar_r"].set("")
        self.comp["list_r"].see(self.comp["list_r"].size() - 1)

    def setstk(self, newstk): # update stack
        self.comp["list_r"].delete( 0, self.comp["list_r"].size() )
        if self.status[1]:
            for i in range( 0, len(newstk) ):
                self.comp["list_r"].insert( i, runtime.tostr( newstk[i] ) )
        else:
            for i in range( 0, len(newstk) ):
                self.comp["list_r"].insert(i, f"STK {i}")
        self.stack = newstk
        self.selstk()

    def setreg(self, pc, sp, ma, mb): # update reg
        if self.status[2]:
            self.comp["strvar2a"].set( runtime.tostr(pc) )
            self.comp["strvar2b"].set( runtime.tostr(ma) )
            self.comp["strvar2c"].set( runtime.tostr(sp) )
            self.comp["strvar2d"].set( runtime.tostr(mb) )
        else:
            self.comp["strvar2a"].set("")
            self.comp["strvar2b"].set("")
            self.comp["strvar2c"].set("")
            self.comp["strvar2d"].set("")

    def testio(self, mode, v): # runtime.testio impl
        if mode == 16:
            self.stdout( runtime.tostr(v[0]) )
            self.stdin()
        elif mode == 17:
            self.stdout( runtime.tostr(v[0]) )
            return None
        else:
            return runtime.testio(mode, v)

    def stdin(self):
        self.status[3] = True
        return ""

    def stdout(self, word):
        idx = self.comp["list_l"].size() - 1
        if idx >= 0:
            word = self.comp["list_l"].get(idx) + word
            self.comp["list_l"].delete(idx)
        word = word.split("\n")
        for i in word:
            self.comp["list_l"].insert(self.comp["list_l"].size(), i)
        self.comp["list_l"].see(self.comp["list_l"].size()-1)
        self.mwin.update()

    def stderr(self, word):
        self.stdout("[err] " + word)

    def stdmsg(self, word):
        self.comp["strvar_l"].set(word)
        self.mwin.update()

class ks_run:
    def __init__(self, iswin, stack, libmode):
        self.vm, self.lb, self.tl, self.eh, self.ui = runtime.kvm(), kscript_lib.lib(libmode[1], libmode[2], libmode[3], iswin), libmode[0], libmode[5], None
        self.mc = kscript_macro.lib(True, True, True, True) if libmode[4] else kscript_macro.lib(False, False, False, False)
        self.mc.drvpath = os.path.abspath("../../_ST5_COMMON/edgewd/msedgedriver.exe") if iswin else os.path.abspath("../../_ST5_COMMON/edgewd/msedgedriver")
        for i in os.listdir("../../_ST5_COMMON/videopack/"):
            if "wkhtmltopdf" in i:
                self.mc.kitpath = os.path.abspath("../../_ST5_COMMON/videopack/" + i)
        cfg = kdb.toolbox()
        with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
            cfg.read( f.read() )
        self.lb.io, self.lb.p_desktop, self.lb.p_local, self.lb.p_starter = self.ui, cfg.get("path.desktop")[3], cfg.get("path.local")[3], cfg.get("path.export")[3]
        self.lb.p_base = os.path.abspath("./").replace("\\", "/")
        if self.lb.p_base[-1] != "/":
            self.lb.p_base = self.lb.p_base + "/"
        if stack > 1:
            self.vm.maxstk = stack

    def load(self, path):
        try:
            info, _, _ = self.vm.view(path)
            self.ui.stdmsg(info)
            self.vm.load(False)
            self.ui.setreg(self.vm.pc, self.vm.sp, self.vm.ma, self.vm.mb)
            self.ui.setstk(self.vm.stack)
        except Exception as e:
            tkm.showerror(title="Load Fail", message=f" {e} ")
    
    def run(self):
        intr = self.vm.run()
        if intr >= 500: # kscript macro
            self.vm.ma, self.vm.errmsg = self.mc.run(self.vm.callmem, intr)
            if self.vm.errmsg != None and self.rh:
                self.ui.stdmsg(f"[err] {self.vm.errmsg}")
                self.vm = None
            elif self.vm.errmsg != None and not self.eh:
                self.ui.stdout(f"[err] {self.vm.errmsg}")
        elif intr >= 32: # kscript lib
            self.vm.ma, self.vm.errmsg = self.lb.run(self.vm.callmem, intr)
            if self.vm.errmsg != None and self.eh:
                self.ui.stdmsg(f"[err] {self.vm.errmsg}")
                self.vm = None
            elif self.vm.errmsg != None and not self.eh:
                self.ui.stdout(f"[err] {self.vm.errmsg}")
        elif intr >= 16: # testio
            if self.tl:
                self.vm.ma = self.ui.testio(intr, self.vm.callmem)
            elif self.eh:
                self.ui.stdmsg("[err] not supported func")
                self.vm = None
            else:
                self.ui.stdout("[err] not supported func")
        elif intr >= 0: # kscript
            if intr == 1:
                self.ui.stdmsg("[msg] program end")
                self.vm = None
            elif intr == 2:
                if self.eh:
                    self.ui.stdmsg(f"[err] {self.vm.errmsg}")
                    self.vm = None
                else:
                    self.ui.stdout(f"[err] {self.vm.errmsg}")
        else: # critical
            self.ui.stdmsg(f"[critical] {self.vm.errmsg}")
            self.vm = None

def cplr(iswin, worker): # compile binary
    dll = ctypes.CDLL("../../_ST5_COMMON/kscriptc/kscriptc.dll") if iswin else ctypes.cdll.LoadLibrary("../../_ST5_COMMON/kscriptc/kscriptc.so")
    dll.func0.argtypes, dll.func0.restype = kobj.call("b", "") # free
    dll.func1.argtypes, dll.func1.restype = kobj.call("iiiii", "") # init c
    dll.func2.argtypes, dll.func2.restype = kobj.call("bii", "") # set data
    dll.func3.argtypes, dll.func3.restype = kobj.call("", "b") # get data
    dll.func4.argtypes, dll.func4.restype = kobj.call("", "f") # get tpass
    dll.func5.argtypes, dll.func5.restype = kobj.call("bi", "") # addpkg
    dll.func6.argtypes, dll.func6.restype = kobj.call("", "b") # compile
    temp = [0 if x else 1 for x in worker.status]
    dll.func1( temp[2], temp[3], temp[4], temp[0], temp[1] )
    o0, o1 = kobj.send( bytes(worker.source[1], encoding="utf-8") )
    dll.func2(o0, o1, 0)
    o0, o1 = kobj.send( bytes(worker.compile[0], encoding="utf-8") )
    dll.func2(o0, o1, 1)
    o0, o1 = kobj.send( bytes(worker.compile[1], encoding="utf-8") )
    dll.func2(o0, o1, 2)
    o0, o1 = kobj.send( bytes(worker.compile[2], encoding="utf-8") )
    dll.func2(o0, o1, 3)
    temp = ["testfunc.txt", "stdlib.txt", "stdio.txt", "osfs.txt", "macro.txt"]
    for i in range(0, 5):
        if worker.status[i + 5]:
            o0, o1 = kobj.send( bytes("../../_ST5_COMMON/kscriptc/_ST5_DATA/" + temp[i], encoding="utf-8") )
            dll.func5(o0, o1)
    t0 = dll.func6()
    err, tpass = str(kobj.recvauto(t0), encoding="utf-8"), dll.func4()
    t1, t2, t3, t4 = dll.func3(4), dll.func3(5), dll.func3(6), dll.func3(-1)
    tkn, ast, asm, bin = kobj.recvauto(t1), kobj.recvauto(t2), kobj.recvauto(t3), kobj.recvauto(t4)
    dll.func0(t0)
    dll.func0(t1)
    dll.func0(t2)
    dll.func0(t3)
    dll.func0(t4)
    if err == "":
        with open("./_ST5_DATA/result.webp", "wb") as f:
            f.write(bin)
        if worker.status[2]:
            with open("./_ST5_DATA/token.txt", "wb") as f:
                f.write(tkn)
        if worker.status[3]:
            with open("./_ST5_DATA/ast.txt", "wb") as f:
                f.write(ast)
        if worker.status[4]:
            with open("./_ST5_DATA/asm.txt", "wb") as f:
                f.write(asm)
        m0, m1 = "Compile Success", f"Binary result generated at ./_ST5_DATA/ with {tpass:.6f}s of compile time."
    else:
        m0, m1 = "Compile Error", f"{err}"
    m1 = " " + " \n ".join( [ m1[x:x+28] for x in range(0, len(m1), 28) ] ) + " "
    root = tk.Tk()
    root.title("KscriptDK")
    lbl = tk.Label(root, font=("Consolas", 14), text=m0) if iswin else tk.Label(root, font=("Consolas", 10), text=m0)
    lbl.pack(padx=5, pady=5)
    text = tk.Text(root, font=("Consolas", 14), width=30,  height=4) if iswin else tk.Text(root, font=("Consolas", 10), width=30,  height=3)
    text.pack(padx=5, pady=5)
    text.insert("1.0", m1)
    root.mainloop()
    return err

def run(iswin, worker): # run binary
    rt = ks_run( iswin, worker.compile[3], worker.status[5:] )
    ui = ks_dbg(iswin, rt)
    rt.ui = ui
    rt.load("./_ST5_DATA/result.webp")
    ui.mwin.mainloop()

kobj.repath()
cfg = kdb.toolbox()
with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
    cfg.read( f.read() )
iswin, flag = (cfg.get("dev.os")[3] == "windows"), True
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
    worker = ks_cplr( iswin, cfg.get("path.desktop")[3] )
    if worker.source[2] == "build":
        cplr(iswin, worker)
    elif worker.source[2] == "run":
        if cplr(iswin, worker) == "":
            run(iswin, worker)
time.sleep(0.5)
