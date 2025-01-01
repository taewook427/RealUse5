# test741 : extension.kdm5

import os
import shutil
import time
import random
import zipfile
import multiprocessing

import tkinter as tk
import tkinter.ttk as tkt
import tkinter.filedialog as tkf
import tkinter.messagebox as tkm
import blocksel

import kaes
import kcom
import kdb
import kobj
import ksc
import kpic
import kzip
import picdt

class zipre:
    def __init__(self, iswin, dsk):
        self.mwin, self.comp, self.parms, self.desktop, self.files, self.webp = None, dict(), {"color":"PaleTurquoise1"}, dsk, [ ], True
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "470x300+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["width"], self.parms["height"], self.parms["ent.1"] = 42, 10, 20
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "560x400+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["width"], self.parms["height"], self.parms["ent.1"] = 40, 8, 16
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def sel_file(): # select files
            time.sleep(0.1)
            self.files = [ x.name.replace("\\", "/") for x in tkf.askopenfiles(title="Select Files", initialdir=self.desktop) ]
            self.comp["list0"].delete( 0, self.comp["list0"].size() )
            for i in self.files:
                self.comp["list0"].insert(self.comp["list0"].size(), i)

        def mode(): # change mode
            time.sleep(0.1)
            self.webp = not self.webp
            t = "webp" if self.webp else " png "
            self.comp["strvar1"].set(t)

        def pack(): # pack to zip
            time.sleep(0.1)
            name = self.comp["entry1"].get()
            if name == "":
                name = "result"
            if self.webp:
                pic, name = picdt.toolbox().data7, "./_ST5_DATA/" + name + ".webp"
            else:
                pic, name = picdt.toolbox().data6, "./_ST5_DATA/" + name + ".png"
            try:
                zipdata(self.files, pic, name)
                tkm.showinfo(title="Release Success", message=f" ZipFile generated at {name} ")
            except Exception as e:
                tkm.showerror(title="Release Fail", message=f" {e} ")

        self.mwin = tk.Tk() # main window
        self.mwin.title("KDM5 zipre")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # selection viewer
        self.comp["frame0"].pack(padx=5, pady=5)
        self.comp["list0"] = tk.Listbox( self.comp["frame0"], width=self.parms["width"],  height=self.parms["height"], font=self.parms["font.1"] )
        self.comp["list0"].pack(side="left", fill="y")
        self.comp["scroll0"] = tk.Scrollbar(self.comp["frame0"], orient="vertical")
        self.comp["scroll0"].config(command=self.comp["list0"].yview)
        self.comp["scroll0"].pack(side="right", fill="y")
        self.comp["list0"].config(yscrollcommand=self.comp["scroll0"].set)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # sel, name, mode, pack
        self.comp["frame1"].pack(fill="x", padx=5, pady=5)
        self.comp["button1a"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], text=" . . . ", command=sel_file)
        self.comp["button1a"].grid(row=0, column=0, padx=5)
        self.comp["label1"] = tk.Label( self.comp["frame1"], font=self.parms["font.0"], text="Name", bg=self.parms["color"] )
        self.comp["label1"].grid(row=0, column=1)
        self.comp["entry1"] = tk.Entry( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.1"] )
        self.comp["entry1"].grid(row=0, column=2, padx=5)
        self.comp["strvar1"] = tk.StringVar()
        self.comp["strvar1"].set("webp")
        self.comp["button1b"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], textvariable=self.comp["strvar1"], command=mode)
        self.comp["button1b"].grid(row=0, column=3, padx=5)
        self.comp["button1c"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], text="Pack", command=pack)
        self.comp["button1c"].grid(row=0, column=4, padx=5)

class divfile:
    def __init__(self, iswin, dsk):
        self.mwin, self.comp, self.parms, self.desktop, self.head, self.files = None, dict(), {"color":"PaleTurquoise1"}, dsk, "", [ ]
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "470x150+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["ent.0"], self.parms["ent.1"] = 37, 8
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "600x200+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["ent.1"] = 37, 6
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def sel_file(): # select file
            time.sleep(0.1)
            temp = tkf.askopenfile(title="Select Files", filetypes=[ ("0 Files", "*.0"), ("Webp Files", "*.webp"), ("All Files", "*.*") ], initialdir=self.desktop).name.replace("\\", "/")
            if temp[-2:] == ".0":
                name, num, self.files = temp[:-2], 0, [ ]
                while os.path.exists(f"{name}.{num}"):
                    self.files.append(f"{name}.{num}")
                    num = num + 1
                self.comp["strvar0b"].set(f"Series of {num} files")
            else:
                self.files = [temp]
                size = os.path.getsize(temp)
                self.comp["strvar0b"].set(f"{size/1048576:.2f} MiB ({size} B)")
            self.comp["strvar0a"].set(temp)

        def pack(): # pack file to div
            time.sleep(0.1)
            try:
                num = int(float( self.comp["entry1a"].get() ) * 1048576)
            except:
                num = -1
            try:
                div = int( self.comp["entry1b"].get() )
            except:
                div = -1
            if num > 0:
                size = num
            elif div > 0:
                temp = os.path.getsize( self.files[0] )
                size = temp // div if temp % div == 0 else int(temp / div) + 1
            else:
                size = 25165824
            temp = self.files[0]
            name, chunk, num = temp[temp.rfind("/")+1:], os.path.getsize(temp), 0
            with open(temp, "rb") as f:
                while size * num < chunk:
                    with open(f"./_ST5_DATA/{name}.{num}", "wb") as t:
                        t.write( f.read(size) )
                    num = num + 1

        def unpack(): # unpack file from div
            time.sleep(0.1)
            temp = self.files[0]
            name = temp[temp.rfind("/")+1:temp.rfind(".")]
            with open(f"./_ST5_DATA/{name}", "wb") as f:
                for i in self.files:
                    with open(i, "rb") as t:
                        f.write( t.read() )

        self.mwin = tk.Tk() # main window
        self.mwin.title("KDM5 div")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # path, info
        self.comp["frame0"].pack(padx=5, pady=5)
        self.comp["button0"] = tk.Button(self.comp["frame0"], bg=self.parms["color"], font=self.parms["font.0"], text=" . . . ", command=sel_file)
        self.comp["button0"].grid(row=0, column=0, padx=5, pady=5)
        self.comp["strvar0a"] = tk.StringVar()
        self.comp["strvar0a"].set("No Selection")
        self.comp["entry0a"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], width=self.parms["ent.0"], textvariable=self.comp["strvar0a"], state="readonly")
        self.comp["entry0a"].grid(row=0, column=1, padx=5, pady=5)
        self.comp["label0"] = tk.Label( self.comp["frame0"], font=self.parms["font.0"], text="Info", bg=self.parms["color"] )
        self.comp["label0"].grid(row=1, column=0, padx=5, pady=5)
        self.comp["strvar0b"] = tk.StringVar()
        self.comp["strvar0b"].set("No Information")
        self.comp["entry0b"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], width=self.parms["ent.0"], textvariable=self.comp["strvar0b"], state="readonly")
        self.comp["entry0b"].grid(row=1, column=1, padx=5, pady=5)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # size, split, pack, unpack
        self.comp["frame1"].pack(padx=5, pady=5)
        self.comp["label1a"] = tk.Label( self.comp["frame1"], font=self.parms["font.0"], text="Size(MiB)", bg=self.parms["color"] )
        self.comp["label1a"].grid(row=0, column=0)
        self.comp["entry1a"] = tk.Entry( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.1"] )
        self.comp["entry1a"].grid(row=0, column=1, padx=5)
        self.comp["label1b"] = tk.Label( self.comp["frame1"], font=self.parms["font.0"], text="AutoSplit", bg=self.parms["color"] )
        self.comp["label1b"].grid(row=0, column=2)
        self.comp["entry1b"] = tk.Entry( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.1"] )
        self.comp["entry1b"].grid(row=0, column=3, padx=5)
        self.comp["button1a"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], text="Pack", command=pack)
        self.comp["button1a"].grid(row=0, column=4)
        self.comp["button1b"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], text="Unpack", command=unpack)
        self.comp["button1b"].grid(row=0, column=5)

class kzipui:
    def __init__(self, iswin, dsk):
        self.mwin, self.comp, self.parms, self.desktop, self.paths = None, dict(), {"color":"PaleTurquoise1"}, dsk, [ ]
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "370x290+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["width"], self.parms["height"], self.parms["ent.1"], self.parms["combo.1"] = 32, 8, 20, 6
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "520x380+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["width"], self.parms["height"], self.parms["ent.1"], self.parms["combo.1"] = 36, 6, 24, 6
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def sel_file(): # select files
            time.sleep(0.1)
            for i in tkf.askopenfiles(title="Select Files", filetypes=[ ("Webp Files", "*.webp"), ("Png Files", "*.png"), ("All Files", "*.*") ], initialdir=self.desktop):
                self.paths.append( i.name.replace("\\", "/") )
            self.render()

        def sel_dir(): # select folder
            time.sleep(0.1)
            temp = tkf.askdirectory(title="Select Folder", initialdir=self.desktop).replace("\\", "/")
            if temp[-1] != "/":
                temp = temp + "/"
            self.paths.append(temp)
            self.render()

        def clear(): # clear selection
            time.sleep(0.1)
            self.paths = [ ]
            self.render()

        def pack(): # pack to kzip
            time.sleep(0.1)
            name = self.comp["entry1"].get()
            if name == "":
                name = "result"
            mode = self.comp["combo1"].get()
            name = name + "." + mode if mode == "webp" or mode == "png" else name + ".bin"
            kzip.dozip(self.paths, mode, "./_ST5_DATA/" + name)

        def unpack(): # unpack from kzip
            time.sleep(0.1)
            for i in range( 0, len(self.paths) ):
                try:
                    if self.paths[i][-1] == "/":
                        self.paths[i] = "Cannot unpack folder."
                    else:
                        kzip.unzip(self.paths[i], f"./_ST5_DATA/{random.randrange(1000, 10000)}/", True)
                        self.paths[i] = "Unpacked successfully."
                except Exception as e:
                    self.paths[i] = f"Error : {e}"
                self.render()

        self.mwin = tk.Tk() # main window
        self.mwin.title("KDM5 kzip")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # selection viewer
        self.comp["frame0"].pack(padx=5, pady=5)
        self.comp["list0"] = tk.Listbox( self.comp["frame0"], width=self.parms["width"],  height=self.parms["height"], font=self.parms["font.1"] )
        self.comp["list0"].pack(side="left", fill="y")
        self.comp["scroll0"] = tk.Scrollbar(self.comp["frame0"], orient="vertical")
        self.comp["scroll0"].config(command=self.comp["list0"].yview)
        self.comp["scroll0"].pack(side="right", fill="y")
        self.comp["list0"].config(yscrollcommand=self.comp["scroll0"].set)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # name, mode
        self.comp["frame1"].pack(fill="x", padx=5, pady=5)
        self.comp["label1"] = tk.Label( self.comp["frame1"], font=self.parms["font.0"], text="Name", bg=self.parms["color"] )
        self.comp["label1"].grid(row=0, column=1)
        self.comp["entry1"] = tk.Entry( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.1"] )
        self.comp["entry1"].grid(row=0, column=2, padx=5)
        self.comp["combo1"] = tkt.Combobox( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["combo.1"], values=["webp", "png", "None"] )
        self.comp["combo1"].set("webp")
        self.comp["combo1"].grid(row=0, column=3, padx=5)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # sel_file, sel_dir, clear, pack, unpack
        self.comp["frame2"].pack(fill="x", padx=10, pady=5)
        self.comp["button2a"] = tk.Button(self.comp["frame2"], bg=self.parms["color"], font=self.parms["font.0"], text=" SelFile ", command=sel_file)
        self.comp["button2a"].grid(row=0, column=0)
        self.comp["button2b"] = tk.Button(self.comp["frame2"], bg=self.parms["color"], font=self.parms["font.0"], text=" SelDir ", command=sel_dir)
        self.comp["button2b"].grid(row=0, column=1)
        self.comp["button2c"] = tk.Button(self.comp["frame2"], bg=self.parms["color"], font=self.parms["font.0"], text=" Clear ", command=clear)
        self.comp["button2c"].grid(row=0, column=2)
        self.comp["button2d"] = tk.Button(self.comp["frame2"], bg=self.parms["color"], font=self.parms["font.0"], text=" Pack ", command=pack)
        self.comp["button2d"].grid(row=0, column=3)
        self.comp["button2e"] = tk.Button(self.comp["frame2"], bg=self.parms["color"], font=self.parms["font.0"], text=" Unpack ", command=unpack)
        self.comp["button2e"].grid(row=0, column=4)

    def render(self):
        self.comp["list0"].delete( 0, self.comp["list0"].size() )
        for i in self.paths:
            self.comp["list0"].insert(self.comp["list0"].size(), i)
        self.mwin.update()

class kpicui:
    def __init__(self, iswin, dsk):
        self.mwin, self.comp, self.parms, self.desktop, self.path, self.high = None, dict(), {"color":"PaleTurquoise1", "activate":"lawn green"}, dsk, "", True
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "450x200+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["ent.0"], self.parms["combo.1"] = 36, 8
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "570x250+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["combo.1"] = 36, 8
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def sel_file(): # select file
            time.sleep(0.1)
            self.path = tkf.askopenfile(title="Select File", initialdir=self.desktop).name.replace("\\", "/")
            size = os.path.getsize(self.path)
            div = {"default":960000, "little":1382400, "standard":10140000, "big":49766400}[ self.comp["combo1a"].get() ]
            if not self.high:
                div = div // 2
            self.comp["strvar0a"].set(self.path)
            self.comp["strvar0b"].set("")
            self.comp["strvar0c"].set(f"{size/1048576:.2f} MiB ({size} B), expecting {size // div + 1} pictures")

        def sel_dir(): # select folder
            time.sleep(0.1)
            self.path = tkf.askdirectory(title="Select Folder", initialdir=self.desktop).replace("\\", "/")
            if self.path[-1] != "/":
                self.path = self.path + "/"
            wk = kpic.toolbox()
            wk.target = self.path
            _, num, style = wk.detect()
            self.comp["strvar0a"].set("")
            self.comp["strvar0b"].set(self.path)
            if num == 0:
                self.comp["strvar0c"].set(f"not detected [hidden mode]")
            else:
                self.comp["strvar0c"].set(f"{num} pictures detected [{style}]")

        def mode(): # high mode
            time.sleep(0.1)
            self.high = not self.high
            t = "activate" if self.high else "color"
            self.comp["button1a"].configure( bg=self.parms[t] )

        def pack(): # pack to kpic
            time.sleep(0.1)
            if self.path[-1] == "/":
                tkm.showerror(title="Pack Fail", message=" KPIC can pack file only. ")
                return
            wk = kpic.toolbox()
            mode = self.comp["combo1a"].get()
            if mode == "little":
                wk.setmold("./little.webp", -1, -1)
            elif mode == "standard":
                wk.setmold("./standard.webp", -1, -1)
            elif mode == "big":
                wk.setmold("./big.webp", -1, -1)
            else:
                wk.setmold("", 800, 800)
            wk.style, wk.target, wk.export = self.comp["combo1b"].get(), self.path, "./_ST5_DATA/"
            try:
                if self.high:
                    wk.pack(2)
                else:
                    wk.pack(4)
                tkm.showinfo(title="Pack Success", message=f" Style : {wk.style}, Size : {mode} \n High integration : {self.high} ")
            except Exception as e:
                tkm.showerror(title="Pack Fail", message=f" {e} ")

        def unpack(): # unpack from kpic
            time.sleep(0.1)
            if self.path[-1] != "/":
                tkm.showerror(title="Unpack Fail", message=" KPIC can detect folder only. ")
                return
            wk = kpic.toolbox()
            wk.target = self.path
            name, num, style = wk.detect()
            if num == 0:
                mode = 2 if self.high else 4
                try:
                    name, num = wk.restore( [ self.path+x for x in os.listdir(self.path) ], mode )
                    tkm.showinfo(title="Restore Success", message=f" Name : {name}, Num : {num} \n High integration : {self.high} ")
                except Exception as e:
                    tkm.showerror(title="Restore Fail", message=f" {e} ")
            else:
                try:
                    wk.style, wk.export = style, "./_ST5_DATA/result.bin"
                    wk.unpack(name, num)
                    tkm.showinfo(title="Unpack Success", message=" result generated at ./_ST5_DATA/result.bin ")
                except Exception as e:
                    tkm.showerror(title="Unpack Fail", message=f" {e} ")

        self.mwin = tk.Tk() # main window
        self.mwin.title("KDM5 kpic")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # path_file, path_dir, info
        self.comp["frame0"].pack(padx=5, pady=5)
        self.comp["button0a"] = tk.Button(self.comp["frame0"], bg=self.parms["color"], font=self.parms["font.0"], text=" . . . ", command=sel_file)
        self.comp["button0a"].grid(row=0, column=0, padx=5, pady=5)
        self.comp["strvar0a"] = tk.StringVar()
        self.comp["strvar0a"].set("No File Selection")
        self.comp["entry0a"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], width=self.parms["ent.0"], textvariable=self.comp["strvar0a"], state="readonly")
        self.comp["entry0a"].grid(row=0, column=1, padx=5, pady=5)
        self.comp["button0b"] = tk.Button(self.comp["frame0"], bg=self.parms["color"], font=self.parms["font.0"], text=" . . . ", command=sel_dir)
        self.comp["button0b"].grid(row=1, column=0, padx=5, pady=5)
        self.comp["strvar0b"] = tk.StringVar()
        self.comp["strvar0b"].set("No Folder Selection")
        self.comp["entry0b"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], width=self.parms["ent.0"], textvariable=self.comp["strvar0b"], state="readonly")
        self.comp["entry0b"].grid(row=1, column=1, padx=5, pady=5)
        self.comp["label0"] = tk.Label( self.comp["frame0"], font=self.parms["font.0"], text="Info", bg=self.parms["color"] )
        self.comp["label0"].grid(row=2, column=0)
        self.comp["strvar0c"] = tk.StringVar()
        self.comp["strvar0c"].set("No Information")
        self.comp["entry0c"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], width=self.parms["ent.0"], textvariable=self.comp["strvar0c"], state="readonly")
        self.comp["entry0c"].grid(row=2, column=1, padx=5, pady=5)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # size, style, mode, pack, unpack
        self.comp["frame1"].pack(padx=5, pady=5)
        self.comp["combo1a"] = tkt.Combobox( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["combo.1"], values=["little", "standard", "big", "default"] )
        self.comp["combo1a"].set("default")
        self.comp["combo1a"].grid(row=0, column=0, padx=5)
        self.comp["combo1b"] = tkt.Combobox( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["combo.1"], values=["webp", "png", "bmp"] )
        self.comp["combo1b"].set("webp")
        self.comp["combo1b"].grid(row=0, column=1, padx=5)
        self.comp["button1a"] = tk.Button(self.comp["frame1"], bg=self.parms["activate"], font=self.parms["font.0"], text=" High ", command=mode)
        self.comp["button1a"].grid(row=0, column=2)
        self.comp["button1b"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], text=" Pack ", command=pack)
        self.comp["button1b"].grid(row=0, column=3)
        self.comp["button1c"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], text=" Unpack ", command=unpack)
        self.comp["button1c"].grid(row=0, column=4)

class kaesui:
    def __init__(self, iswin, dsk):
        self.mwin, self.comp, self.parms, self.desktop, self.aes = None, dict(), {"color":"PaleTurquoise1", "activate":"lawn green"}, dsk, kaes.funcmode()
        self.sign, self.kf, self.pw, self.hint, self.paths, self.delori, self.pmode = ("", ""), self.aes.basickey(), b"", "", [ ], False, 0
        self.getphash = lambda x: ksc.crc32hash( bytes(x, encoding="utf-8") ).hex()
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "450x400+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["width"], self.parms["height"], self.parms["ent.2"], self.parms["ent.3"] = 40, 6, 12, 36
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "580x480+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["width"], self.parms["height"], self.parms["ent.2"], self.parms["ent.3"] = 40, 4, 12, 36

    def entry(self, isenc):
        def del_mode():
            time.sleep(0.1)
            self.delori = not self.delori
            t = "activate" if self.delori else "color"
            self.comp["button4e"].configure( bg=self.parms[t] )

        def p_mode():
            time.sleep(0.1)
            self.pmode = (self.pmode + 1) % 3
            self.comp["strvar4"].set( {0:" webp ", 1:" png ", 2:" None "}[self.pmode] )

        def clear_path():
            time.sleep(0.1)
            self.paths = [ ]
            self.render()

        def sign_direct():
            path, data = self.getfile("")
            if path == "":
                path, data = "No Sign", b""
            self.getsign( str(data, encoding="utf-8") )
            self.comp["strvar1"].set(path)

        def sign_comm():
            path, data = self.getfile( self.comp["entry1b"].get() )
            if path == "":
                path, data = "No Sign", b""
            self.getsign( str(data, encoding="utf-8") )
            self.comp["strvar1"].set(path)

        def kf_direct():
            path, self.kf = self.getfile("")
            if path == "":
                path = "bkf"
            self.comp["strvar2"].set(path)

        def kf_comm():
            path, self.kf = self.getfile( self.comp["entry2b"].get() )
            if path == "":
                path = "bkf"
            self.comp["strvar2"].set(path)

        def sel_file():
            time.sleep(0.1)
            for i in tkf.askopenfiles(title="Select Files", filetypes=[ ("Webp Files", "*.webp"), ("Png Files", "*.png"), ("K Files", "*.k"), ("All Files", "*.*") ], initialdir=self.desktop):
                self.paths.append( i.name.replace("\\", "/") )
            self.render()

        def sel_dir():
            time.sleep(0.1)
            temp = tkf.askdirectory(title="Select Folder", initialdir=self.desktop).replace("\\", "/")
            if temp[-1] != "/":
                temp = temp + "/"
            self.paths.append(temp)
            self.render()

        self.mwin = tk.Tk() # main window
        self.mwin.title("KDM5 kaes")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # selection viewer
        self.comp["frame0"].pack(padx=5, pady=5)
        self.comp["list0"] = tk.Listbox( self.comp["frame0"], width=self.parms["width"],  height=self.parms["height"], font=self.parms["font.1"] )
        self.comp["list0"].pack(side="left", fill="y")
        self.comp["scroll0"] = tk.Scrollbar(self.comp["frame0"], orient="vertical")
        self.comp["scroll0"].config(command=self.comp["list0"].yview)
        self.comp["scroll0"].pack(side="right", fill="y")
        self.comp["list0"].config(yscrollcommand=self.comp["scroll0"].set)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # sign
        self.comp["frame1"].pack(padx=5, pady=5)
        self.comp["label1"] = tk.Label( self.comp["frame1"], font=self.parms["font.0"], text="Sign", bg=self.parms["color"] )
        self.comp["label1"].grid(row=0, column=0, pady=5)
        self.comp["strvar1"] = tk.StringVar()
        self.comp["strvar1"].set("No Sign")
        if isenc:
            self.comp["button1a"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], text=" . . . ", command=sign_direct)
            self.comp["button1a"].grid(row=0, column=1, padx=5, pady=5)
            self.comp["entry1a"] = tk.Entry(self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.2"], textvariable=self.comp["strvar1"], state="readonly")
            self.comp["entry1a"].grid(row=0, column=2, padx=5, pady=5)
            self.comp["button1b"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], text=" . . . ", command=sign_comm)
            self.comp["button1b"].grid(row=0, column=3, padx=5, pady=5)
            self.comp["entry1b"] = tk.Entry( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.2"] )
            self.comp["entry1b"].grid(row=0, column=4, padx=5, pady=5)
        else:
            self.comp["entry1"] = tk.Entry(self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.3"], textvariable=self.comp["strvar1"], state="readonly")
            self.comp["entry1"].grid(row=0, column=1, padx=5, pady=5)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # kf
        self.comp["frame2"].pack(padx=5, pady=5)
        self.comp["strvar2"] = tk.StringVar()
        self.comp["strvar2"].set("bkf")
        if isenc:
            self.comp["label2"] = tk.Label( self.comp["frame1"], font=self.parms["font.0"], text="KF", bg=self.parms["color"] )
            self.comp["label2"].grid(row=1, column=0)
            self.comp["button2a"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], text=" . . . ", command=kf_direct)
            self.comp["button2a"].grid(row=1, column=1, padx=5)
            self.comp["entry2a"] = tk.Entry(self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.2"], textvariable=self.comp["strvar2"], state="readonly")
            self.comp["entry2a"].grid(row=1, column=2, padx=5)
            self.comp["button2b"] = tk.Button(self.comp["frame1"], bg=self.parms["color"], font=self.parms["font.0"], text=" . . . ", command=kf_comm)
            self.comp["button2b"].grid(row=1, column=3, padx=5)
            self.comp["entry2b"] = tk.Entry( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.2"] )
            self.comp["entry2b"].grid(row=1, column=4, padx=5)
        else:
            self.comp["label2"] = tk.Label( self.comp["frame2"], font=self.parms["font.0"], text="KF", bg=self.parms["color"] )
            self.comp["label2"].grid(row=0, column=0)
            self.comp["button2a"] = tk.Button(self.comp["frame2"], bg=self.parms["color"], font=self.parms["font.0"], text=" . . . ", command=kf_direct)
            self.comp["button2a"].grid(row=0, column=1, padx=5)
            self.comp["entry2a"] = tk.Entry(self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"], textvariable=self.comp["strvar2"], state="readonly")
            self.comp["entry2a"].grid(row=0, column=2, padx=5)
            self.comp["button2b"] = tk.Button(self.comp["frame2"], bg=self.parms["color"], font=self.parms["font.0"], text=" . . . ", command=kf_comm)
            self.comp["button2b"].grid(row=0, column=3, padx=5)
            self.comp["entry2b"] = tk.Entry( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"] )
            self.comp["entry2b"].grid(row=0, column=4, padx=5)

        self.comp["frame3"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # pw, hint
        self.comp["frame3"].pack(padx=5, pady=5)
        self.comp["strvar3"] = tk.StringVar()
        self.comp["strvar3"].set("No Hint")
        self.comp["label3a"] = tk.Label( self.comp["frame3"], font=self.parms["font.0"], text="PW", bg=self.parms["color"] )
        self.comp["label3a"].grid(row=0, column=0)
        self.comp["entry3a"] = tk.Entry(self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.3"], show="*")
        self.comp["entry3a"].grid(row=0, column=1, padx=5, pady=5)
        self.comp["label3b"] = tk.Label( self.comp["frame3"], font=self.parms["font.0"], text="Hint", bg=self.parms["color"] )
        self.comp["label3b"].grid(row=1, column=0)
        if isenc:
            self.comp["entry3b"] = tk.Entry( self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.3"] )
        else:
            self.comp["entry3b"] = tk.Entry(self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.3"], textvariable=self.comp["strvar3"], state="readonly")
        self.comp["entry3b"].grid(row=1, column=1, padx=5, pady=5)

        self.comp["frame4"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # sel_file, sel_dir, clear, mode01, submit
        self.comp["frame4"].pack(padx=5, pady=5)
        self.comp["strvar4"] = tk.StringVar()
        self.comp["strvar4"].set(" webp ")
        self.comp["button4a"] = tk.Button(self.comp["frame4"], bg=self.parms["color"], font=self.parms["font.0"], text=" SelFile ", command=sel_file)
        self.comp["button4a"].grid(row=0, column=0)
        self.comp["button4b"] = tk.Button(self.comp["frame4"], bg=self.parms["color"], font=self.parms["font.0"], text=" SelDir ", command=sel_dir)
        self.comp["button4b"].grid(row=0, column=1)
        self.comp["button4c"] = tk.Button(self.comp["frame4"], bg=self.parms["color"], font=self.parms["font.0"], text=" Clear ", command=clear_path)
        self.comp["button4c"].grid(row=0, column=2)
        self.comp["button4d"] = tk.Button(self.comp["frame4"], bg=self.parms["color"], font=self.parms["font.0"], textvariable=self.comp["strvar4"], command=p_mode)
        self.comp["button4d"].grid(row=0, column=3)
        self.comp["button4e"] = tk.Button(self.comp["frame4"], bg=self.parms["color"], font=self.parms["font.0"], text=" DEL ", command=del_mode)
        self.comp["button4e"].grid(row=0, column=4)
        temp = " Encrypt " if isenc else " Decrypt "
        self.comp["button4f"] = tk.Button(self.comp["frame4"], bg=self.parms["color"], font=self.parms["font.0"], text=temp, command=self.submit)
        self.comp["button4f"].grid(row=0, column=5)

    def render(self):
        self.comp["list0"].delete( 0, self.comp["list0"].size() )
        for i in self.paths:
            self.comp["list0"].insert(self.comp["list0"].size(), i)
        self.mwin.update()

    def getfile(self, code):
        time.sleep(0.1)
        try:
            if code == "":
                path = tkf.askopenfile( title="Select File", filetypes=[ ("Jpg Files", "*.jpg"), ("Png Files", "*.png"), ("Webp Files", "*.webp"), ("Txt Files", "*.txt"), ("All Files", "*.*") ] ).name
                with open(path, "rb") as f:
                    data = f.read()
            else:
                sport, skey = kcom.unpack(code)
                node, wk = kcom.node(), kaes.funcmode()
                node.port = sport
                data = node.recieve(skey)
                path = str(data[48:], encoding="utf-8")
                with open(path, "rb") as f:
                    tgt = f.read()
                wk.before = tgt
                wk.decrypt( data[0:48] )
                data = wk.after
                wk.before, wk.after = "", ""
        except:
            path, data = "", self.aes.basickey()
        return path, data

    def getsign(self, txt):
        try:
            wk  = kdb.toolbox()
            wk.read(txt)
            self.sign = ( wk.get("public")[3], wk.get("private")[3] )
        except:
            self.sign = ("", "")

    def delorigin(self, path):
        if self.delori:
            if path[-1] == "/":
                shutil.rmtree(path)
            else:
                os.remove(path)

    def submit(self):
        time.sleep(0.1)
        self.pw, self.hint = bytes(self.comp["entry3a"].get(), encoding="utf-8"), self.comp["entry3b"].get()
        a, b, c = len(self.pw), len(str(self.pw, encoding='utf-8')), len(self.kf)
        temp = f" {len(self.paths)} Objects \n PW {a} B (len {b}), KF {c} B \n phash : {self.getphash(self.sign[0])} \n Delete original : {self.delori} "
        return tkm.askokcancel(title="Submit Check", message=temp)

class kaesen(kaesui):
    def __init__(self, iswin, dsk):
        super().__init__(iswin, dsk)
        self.entry(True)
        self.mwin.mainloop()

    def submit(self):
        if not super().submit():
            return
        wk = kaes.allmode()
        for i in range( 0, len(self.paths) ):
            wk.hint, wk.signkey, original = self.hint, self.sign, self.paths[i]
            try:
                if self.paths[i][-1] == "/":
                    wk.msg, temp = "KDM5.kaes.endir", self.paths[i][:-1]
                    temp = temp[:temp.rfind("/")+1] + "temp741.bin"
                    kzip.dozip([ self.paths[i] ], "", temp)
                    temp2 = wk.encrypt(self.pw, self.kf, temp, self.pmode)
                    os.remove(temp)
                    self.paths[i] = f"Encrypt success (dir) : {temp2[temp2.rfind("/")+1:]}"
                else:
                    wk.msg = "KDM5.kaes.enfile"
                    temp = wk.encrypt(self.pw, self.kf, self.paths[i], self.pmode)
                    self.paths[i] = f"Encrypt success (file) : {temp[temp.rfind("/")+1:]}"
                self.delorigin(original)
            except Exception as e:
                self.paths[i] = f"Error : {e}"
            self.render()
        tkm.showinfo(title="Encrypt Complete", message=f" Encrypted {len(self.paths)} data. ")

class kaesde(kaesui):
    def __init__(self, iswin, dsk):
        def select(event):
            time.sleep(0.1)
            self.viewer.view( self.paths[ self.comp["list0"].curselection()[0] ] )
            self.comp["strvar1"].set( self.chksigndata( self.viewer.signkey[0] ) )
            self.comp["strvar3"].set(self.viewer.hint)
        super().__init__(iswin, dsk)
        self.entry(False)
        self.getsigndata()
        self.viewer = kaes.allmode()
        self.comp["list0"].bind("<ButtonRelease-1>", select)
        self.mwin.mainloop()

    def getsigndata(self):
        wk = kdb.toolbox()
        with open("../../_ST5_SIGN.txt", "r", encoding="utf-8") as f:
            wk.read( f.read() )
        self.signdata, num = [ ], 0
        while f"{num}.public" in wk.name:
            self.signdata.append( wk.get(f"{num}.public")[3] )
            num = num + 1

    def chksigndata(self, pub):
        if pub == "":
            return "No Sign [00000000]"
        else:
            for i in self.signdata:
                if i == pub:
                    return f"Valid [{self.getphash(i)}]"
            return f"Untrusted [{self.getphash(pub)}]"

    def submit(self):
        if not super().submit():
            return
        for i in range( 0, len(self.paths) ):
            self.viewer.msg, self.viewer.signkey, original = "", ["", ""], self.paths[i]
            try:
                if self.paths[i][-1] == "/":
                    raise Exception("cannot decrypt folder")
                self.viewer.view( self.paths[i] )
                if self.viewer.msg == "KDM5.kaes.endir":
                    temp = self.paths[i][:self.paths[i].rfind("/")+1]
                    temp1 = temp + "temp741/"
                    temp2 = self.viewer.decrypt( self.pw, self.kf, self.paths[i] )
                    kzip.unzip(temp2, temp1, True)
                    for j in os.listdir(temp1):
                        shutil.move(temp1+j, temp+j)
                        self.paths[i] = f"Decrypt success (dir) : {j} {self.chksigndata(self.viewer.signkey[0])}"
                    shutil.rmtree(temp1)
                    os.remove(temp2)
                else:
                    temp = self.viewer.decrypt( self.pw, self.kf, self.paths[i] )
                    if self.viewer.msg == "KDM5.kaes.enfile":
                        self.paths[i] = f"Decrypt success (file) : {temp[temp.rfind("/")+1:]} {self.chksigndata(self.viewer.signkey[0])}"
                    else:
                        self.paths[i] = f"Decrypt success (unknown) : {temp[temp.rfind("/")+1:]} {self.chksigndata(self.viewer.signkey[0])}"
                self.delorigin(original)
            except Exception as e:
                self.paths[i] = f"Error : {e}"
            self.render()
        tkm.showinfo(title="Decrypt Complete", message=f" Decrypted {len(self.paths)} data. ")
    
class mainclass(blocksel.toolbox):
    def __init__(self, iswin):
        super().__init__("KDM5", iswin)
        path, self.selection = "../../_ST5_COMMON/iconpack/kdm/", -1
        self.txts, self.curpos, self.upos, self.umsg = ["ZIPrelease", "FileDiv", "KZIP_ende", "KPIC_ende", "KAES_enc", "KAES_dec"], 1, 0, ["Select Mode"]
        self.pics = [path+"zipre.png", path+"div_ende.png", path+"kzip_ende.png", path+"kpic_ende.png", path+"kaes_en.png", path+"kaes_de.png"]

    def custom0(self, x):
        self.selection = x
        self.mwin.destroy()

def zipdata(files, pic, path): # zip release files to path
    if not os.path.exists("./temp741/"):
        os.mkdir("./temp741/")
    temp = zipfile.ZipFile("./temp741.zip", "w")
    for i in files:
        name = i[i.rfind("/")+1:]
        shutil.copy(i, "./temp741/"+name)
        temp.write("./temp741/"+name, compress_type=zipfile.ZIP_DEFLATED)
    temp.close() # temp zip file

    with open('./temp741.zip','rb') as f:
        binary = f.read()
    cheadloc = kobj.decode( binary[-6:-2] ) # centeral header
    num = kobj.decode( binary[-12:-10] ) # file num
    fdata = binary[0:cheadloc] # file data    
    ehead = binary[-22:-6] + kobj.encode(len(pic) + cheadloc, 4) + binary[-2:] # end header
    chead = b'' # centeral header
    temp = cheadloc
    for i in range(0, num):
        seta = binary[temp:temp+28]
        flen = kobj.decode( binary[temp+28:temp+30] ) # file name len
        setb = binary[temp+30:temp+42]
        start = kobj.encode(kobj.decode( binary[temp+42:temp+46] ) + len(pic), 4)
        fname = binary[temp+46:temp+46+flen]
        chead = chead + seta + binary[temp+28:temp+30] + setb + start + fname
        temp = temp + 46 + flen

    with open(path, 'wb') as f:
        f.write(pic)
        f.write(fdata)
        f.write(chead)
        f.write(ehead)
    os.remove("./temp741.zip")
    shutil.rmtree("./temp741/")

if __name__ == "__main__":
    multiprocessing.freeze_support()
    kobj.repath()
    cfg = kdb.toolbox()
    with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
        cfg.read( f.read() )
    iswin, desktop = (cfg.get("dev.os")[3] == "windows"), cfg.get("path.desktop")[3]
    if not os.path.exists("./_ST5_DATA/"):
        os.mkdir("./_ST5_DATA/")
    if os.path.exists("../../_ST5_COMMON/iconpack/"):
        worker = mainclass(iswin)
        worker.entry()
        worker.guiloop()
        if worker.selection == 0:
            t = zipre(iswin, desktop)
        elif worker.selection == 1:
            t = divfile(iswin, desktop)
        elif worker.selection == 2:
            t = kzipui(iswin, desktop)
        elif worker.selection == 3:
            t = kpicui(iswin, desktop)
        elif worker.selection == 4:
            t = kaesen(iswin, desktop)
        elif worker.selection == 5:
            t = kaesde(iswin, desktop)
    else:
        tkm.showinfo(title="No Package", message=" KDM5 requires package >>> common.iconpack <<<. \n Install dependent package and start again. ")
    time.sleep(0.5)