# test740 : extension.kvdrive frontend
# 리눅스 빌드 추가 : --hidden-import='PIL._tkinter_finder'

import os
import shutil
import time
import ctypes
import hashlib

import tkinter as tk
import tkinter.ttk as tkt
import tkinter.filedialog as tkf
import tkinter.simpledialog as tks
import tkinter.messagebox as tkm
import explorer

import kobj
import kdb
import kcom
import kaes

class pwkf:
    def __init__(self, iswin, aes, back):
        self.mwin, self.comp, self.parms, self.clu, self.acc, self.hint, self.aes, self.iomgr = None, dict(), {"color":"thistle1", "act":"cyan"}, "", "", "Not Selected", aes, back
        self.pw, self.kf, self.clutype, self.clusize, self.mode, self.unique, self.local = b"", aes.basickey(), "KV4adv", "default", "Login", True, f"./{int(time.time())}/"
        if iswin: # windows
            self.parms["ms"], self.parms["f0"], self.parms["f1"] = "540x240+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["e0"], self.parms["c0"], self.parms["e1"], self.parms["e2"] = 27, 8, 18, 40
        else: # linux
            self.parms["ms"], self.parms["f0"], self.parms["f1"] = "640x300+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["e0"], self.parms["c0"], self.parms["e1"], self.parms["e2"] = 24, 8, 16, 35

    def entry(self, isboot, master, callback):
        def sel_clu(): # select cluster
            time.sleep(0.1)
            if isboot:
                self.clu = os.path.abspath( tkf.askdirectory(title="Select Cluster", initialdir="./_ST5_DATA/") ).replace("\\", "/")
                if self.clu[-1] != "/":
                    self.clu = self.clu + "/"
                self.acc = self.clu + "0a.webp"
                self.comp["strvar0a"].set(self.clu)
                self.comp["strvar0b"].set(self.acc)
                self.isvalid()

        def sel_acc(): # select account
            time.sleep(0.1)
            if isboot:
                self.acc = os.path.abspath(tkf.askopenfile(title="", filetypes=[ ("BlockA", "*.webp"), ("All Files", "*.*") ], initialdir="./_ST5_DATA/").name).replace("\\", "/")
                self.comp["strvar0b"].set(self.acc)
                self.isvalid()

        def sel_uni(): # login unique
            time.sleep(0.1)
            if isboot:
                self.unique = not self.unique
                if self.unique:
                    self.comp["button0c"].configure( bg=self.parms["act"] )
                else:
                    self.comp["button0c"].configure( bg=self.parms["color"] )

        def sel_key_direct(): # select kf direct
            time.sleep(0.1)
            time.sleep(0.1)
            try:
                path = tkf.askopenfile( title="Select KeyFile", filetypes=[ ("Jpg Files", "*.jpg"), ("Png Files", "*.png"), ("Webp Files", "*.webp"), ("All Files", "*.*") ] ).name
                with open(path, "rb") as f:
                    data = f.read()
                self.comp["strvar1"].set(path)
                self.kf = data
            except:
                self.comp["strvar1"].set("bkf")
                self.kf = self.aes.basickey()

        def sel_key_comm(): # select kf comm
            time.sleep(0.1)
            try:
                sport, skey = kcom.unpack( self.comp["entry1b"].get() )
                node, wk = kcom.node(), kaes.funcmode()
                node.port = sport
                data = node.recieve(skey)
                path = str(data[48:], encoding="utf-8")
                with open(path, "rb") as f:
                    tgt = f.read()
                wk.before = tgt
                wk.decrypt( data[0:48] )
                self.comp["strvar1"].set(path)
                self.kf = wk.after
                wk.before, wk.after = "", ""
            except:
                self.comp["strvar1"].set("bkf")
                self.kf = self.aes.basickey()

        def submit(): # submit data
            time.sleep(0.1)
            if isboot:
                self.clutype, self.clusize, self.mode = self.comp["combo0a"].get(), self.comp["combo0b"].get(), self.comp["combo0c"].get()
            else:
                self.hint = bytes(self.comp["entry0c"].get(), encoding="utf-8")
            self.pw = bytes(self.comp["entry2"].get(), encoding="utf-8")
            if isboot:
                if self.mode == "MkNew":
                    self.mknew()
                elif self.mode == "Rebuild":
                    self.rebuild()
                elif self.mode == "Login":
                    self.login()
            else:
                callback()

        self.mwin = tk.Tk() if isboot else tk.Toplevel(master) # main window
        self.mwin.title("PWKF")
        self.mwin.geometry( self.parms["ms"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # clu, acc, hint
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["label0a"] = tk.Label(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text="Cluster")
        self.comp["label0a"].grid(row=0, column=0, padx=5, pady=5)
        self.comp["button0a"] = tk.Button(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text=" . . . ", command=sel_clu)
        self.comp["button0a"].grid(row=0, column=1)
        self.comp["strvar0a"] = tk.StringVar()
        self.comp["strvar0a"].set(self.clu)
        self.comp["entry0a"] = tk.Entry(self.comp["frame0"], font=self.parms["f1"], width=self.parms["e0"], textvariable=self.comp["strvar0a"], state="readonly")
        self.comp["entry0a"].grid(row=0, column=2, padx=5, pady=5)
        if isboot:
            self.comp["combo0a"] = tkt.Combobox( self.comp["frame0"], font=self.parms["f1"], width=self.parms["c0"], values=["KV4adv", "KV5st"] )
            self.comp["combo0a"].set(self.clutype)
        else:
            self.comp["combo0a"] = tk.Label(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text=self.clutype)
        self.comp["combo0a"].grid(row=0, column=3, padx=5, pady=5)

        self.comp["label0b"] = tk.Label(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text="Account")
        self.comp["label0b"].grid(row=1, column=0, padx=5, pady=5)
        self.comp["button0b"] = tk.Button(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text=" . . . ", command=sel_acc)
        self.comp["button0b"].grid(row=1, column=1)
        self.comp["strvar0b"] = tk.StringVar()
        self.comp["strvar0b"].set(self.acc)
        self.comp["entry0b"] = tk.Entry(self.comp["frame0"], font=self.parms["f1"], width=self.parms["e0"], textvariable=self.comp["strvar0b"], state="readonly")
        self.comp["entry0b"].grid(row=1, column=2, padx=5, pady=5)
        if isboot:
            self.comp["combo0b"] = tkt.Combobox( self.comp["frame0"], font=self.parms["f1"], width=self.parms["c0"], values=["small", "standard", "large", "default"] )
            self.comp["combo0b"].set(self.clusize)
        else:
            self.comp["combo0b"] = tk.Label(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text=self.clusize)
        self.comp["combo0b"].grid(row=1, column=3, padx=5, pady=5)

        self.comp["label0c"] = tk.Label(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text="Hint")
        self.comp["label0c"].grid(row=2, column=0, padx=5, pady=5)
        self.comp["button0c"] = tk.Button(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["act"], text=" Uni ", command=sel_uni) if self.unique else tk.Button(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text=" Uni ", command=sel_uni)
        self.comp["button0c"].grid(row=2, column=1)
        if isboot:
            self.comp["strvar0c"] = tk.StringVar()
            self.comp["strvar0c"].set(self.hint)
            self.comp["entry0c"] = tk.Entry(self.comp["frame0"], font=self.parms["f1"], width=self.parms["e0"], textvariable=self.comp["strvar0c"], state="readonly")
            self.comp["combo0c"] = tkt.Combobox( self.comp["frame0"], font=self.parms["f1"], width=self.parms["c0"], values=["MkNew", "Login", "Rebuild"] )
            self.comp["combo0c"].set(self.mode)
        else:
            self.comp["entry0c"] = tk.Entry( self.comp["frame0"], font=self.parms["f1"], width=self.parms["e0"] )
            self.comp["combo0c"] = tk.Label(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text=self.mode)
        self.comp["entry0c"].grid(row=2, column=2, padx=5, pady=5)
        self.comp["combo0c"].grid(row=2, column=3, padx=5, pady=5)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # kf
        self.comp["frame1"].pack(fill="x", padx=5, pady=5)
        self.comp["label1"] = tk.Label(self.comp["frame1"], font=self.parms["f0"], bg=self.parms["color"], text="KF")
        self.comp["label1"].grid(row=0, column=0, padx=5)
        self.comp["button1a"] = tk.Button(self.comp["frame1"], font=self.parms["f0"], bg=self.parms["color"], text=" . . . ", command=sel_key_direct)
        self.comp["button1a"].grid(row=0, column=1)
        self.comp["strvar1"] = tk.StringVar()
        self.comp["strvar1"].set("bkf")
        self.comp["entry1a"] = tk.Entry(self.comp["frame1"], font=self.parms["f1"], width=self.parms["e1"], textvariable=self.comp["strvar1"], state="readonly")
        self.comp["entry1a"].grid(row=0, column=2, padx=5)
        self.comp["button1b"] = tk.Button(self.comp["frame1"], font=self.parms["f0"], bg=self.parms["color"], text=" . . . ", command=sel_key_comm)
        self.comp["button1b"].grid(row=0, column=3)
        self.comp["entry1b"] = tk.Entry( self.comp["frame1"], font=self.parms["f1"], width=self.parms["e1"] )
        self.comp["entry1b"].grid(row=0, column=4, padx=5)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # pw & go
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["label2"] = tk.Label(self.comp["frame2"], font=self.parms["f0"], bg=self.parms["color"], text="PW")
        self.comp["label2"].grid(row=0, column=0)
        self.comp["entry2"] = tk.Entry(self.comp["frame2"], font=self.parms["f1"], width=self.parms["e2"], show="*")
        self.comp["entry2"].grid(row=0, column=1, padx=5)
        self.comp["button2"] = tk.Button(self.comp["frame2"], font=self.parms["f0"], bg=self.parms["color"], text=" Submit ", command=submit)
        self.comp["button2"].grid(row=0, column=2)

    def isvalid(self):
        self.clutype, self.mode = self.comp["combo0a"].get(), self.comp["combo0c"].get()
        if self.mode != "MkNew":
            if self.clutype == "KV5st":
                self.iomgr.ext.func1(1)
            else:
                self.iomgr.ext.func1(0)
            self.iomgr.order( "init", [ ] )
            temp = [cfg.get("path.desktop")[3], "./temp740/", self.clu, self.acc]
            if self.unique:
                temp[1] = self.local
            err = self.iomgr.order("boot", temp)
            if err == "":
                r0, r1 = self.iomgr.getter(0), self.iomgr.getter(1)
                self.hint = f"{str(r1[0],encoding='utf-8')} [{r0[1]}@{r0[0]}]"
                self.comp["strvar0c"].set(self.hint)
            else:
                self.hint = err
                self.comp["strvar0c"].set(err)

    def mknew(self):
        if self.clutype == "KV5st":
            self.iomgr.ext.func1(1)
        else:
            self.iomgr.ext.func1(0)
        self.iomgr.order( "init", [ ] )
        temp = [self.clu, tks.askstring(title="Select Name", prompt="New Cluster Name : ")]
        if self.clutype == "KV5st":
            temp.append(self.clusize)
        else:
            option = { "default":("134217728", "128"), "small":("536870912", "256"), "standard":("2147483648", "512"), "large":("17179869184", "512") }
            temp.append( option[self.clusize][0] )
            temp.append( option[self.clusize][1] )
        err = self.iomgr.order("new", temp)
        if err == "":
            tkm.showinfo(title="Cluster Generated", message=f" New cluster at {self.clu} ")
        else:
            tkm.showerror(title="MkNew Fail", message=f" {err} ")

    def rebuild(self):
        if self.clutype == "KV5st":
            self.iomgr.ext.func1(1)
        else:
            self.iomgr.ext.func1(0)
        self.iomgr.order( "init", [ ] )
        self.iomgr.setter([self.pw, self.kf], 1)
        err = self.iomgr.order( "rebuild", [self.clu] )
        if err == "":
            tkm.showinfo(title="Cluster Rebuild", message=f" Rebuild success at {self.clu} ")
        else:
            tkm.showerror(title="Rebuild Fail", message=f" {err} ")

    def login(self):
        self.iomgr.setter([self.pw, self.kf], 1)
        err = self.iomgr.order( "login", ["10"] )
        if err == "":
            self.mwin.destroy()
            self.mwin, self.pw, self.kf, self.hint = None, b"01234567", b"76543210", b"89898989"
        else:
            tkm.showerror(title="Login Fail", message=f" {err} ")

class iomgr:
    def __init__(self, iswin):
        if iswin:
            dll = "./kvdrive.dll"
            hv = b'\xf4\xdagF\x05\x056n\xd3?\xdf\xa0\xf0\xdc\xd8[Jp\xce\x80Fmj\xa6\xda\xa4\xb2\x14bqc\xa88\x8d)\t\n\xa6\x0f\xe1|\xb11\xb9qh\xf6\x8c\xf5\xef)\x12\x81\x8d\xee\x80\x8e$6\xb1NLq\x92'
        else:
            dll = "./kvdrive.so"
            hv = b'\x8e\xe2\xfel\x08\xd1\x0c:>\x88\x81\xbfX\x9b<5\xeaO]Z\xdcst<\xcc\xb8C\xe5I\x8b\xd5\x1f` \xac\\pOz\x9cW\x93V\x90]\xf4 \x9f1\xe8D\xf4E\xe9M\xc8pj\xd2\x1e7\x852%'
        with open(dll, "rb") as f:
            hc = hashlib.sha3_512( f.read() ).digest()
        if hc != hv:
            raise Exception("wrong FFI")
        if iswin:
            self.ext = ctypes.CDLL(dll)
        else:
            self.ext = ctypes.cdll.LoadLibrary(dll)
        self.ext.func0.argtypes, self.ext.func0.restype = kobj.call("b", "") # free
        self.ext.func1.argtypes, self.ext.func1.restype = kobj.call("i", "") # set flag
        self.ext.func2.argtypes, self.ext.func2.restype = kobj.call("i", "") # clear buf
        self.ext.func3.argtypes, self.ext.func3.restype = kobj.call("bii", "") # setter
        self.ext.func4.argtypes, self.ext.func4.restype = kobj.call("i", "b") # getter
        self.ext.func5.argtypes, self.ext.func5.restype = kobj.call("i", "i") # status flag
        self.ext.func6.argtypes, self.ext.func6.restype = kobj.call("bi", "b") # cmd
        self.ext.func7.argtypes, self.ext.func7.restype = kobj.call("", "b") # status self
        self.ext.func8.argtypes, self.ext.func8.restype = kobj.call("", "b") # status dir
        self.ext.func9.argtypes, self.ext.func9.restype = kobj.call("", "b") # status file
        self.ext.func1(3)

    def order(self, cmd, parms): # command order ( str, str[] )
        self.setter(parms, 3)
        o0, o1 = kobj.send( bytes(cmd, encoding="utf-8") )
        p = self.ext.func6(o0, o1)
        err = str(kobj.recvauto(p), encoding="utf-8")
        self.ext.func0(p)
        return err

    def setter(self, data, mode): # mode 0 str ( str[] ), 1 bytes ( bytes[] ), 2 curpath (str), 3 cmdbuf ( str[] )
        if mode == 0:
            self.ext.func2(1)
            for i in data:
                o0, o1 = kobj.send( bytes(i, encoding="utf-8") )
                self.ext.func3(o0, o1, 0)
        elif mode == 1:
            self.ext.func2(2)
            for i in data:
                o0, o1 = kobj.send(i)
                self.ext.func3(o0, o1, 1)
        elif mode == 2:
            o0, o1 = kobj.send( bytes(data, encoding="utf-8") )
            self.ext.func3(o0, o1, 2)
        elif mode == 3:
            self.ext.func2(0)
            for i in data:
                o0, o1 = kobj.send( bytes(i, encoding="utf-8") )
                self.ext.func3(o0, o1, 3)

    def getter(self, mode): # mode 0 str ( str[] ), 1 bytes ( bytes[] ), 2 curpath (str), 3 asyncerr (str)
        p = self.ext.func4(mode)
        data = kobj.recvauto(p)
        self.ext.func0(p)
        if mode == 0:
            return [ str(x, encoding="utf-8") for x in kobj.unpack(data) ]
        elif mode == 1:
            return kobj.unpack(data)
        elif mode == 2 or mode == 3:
            return str(data, encoding="utf-8")

    def sync(self): # get curdir data
        self.self_name, self.self_time, self.self_size, self.self_locked, self.self_subdir, self.self_subfile = "", "", -1, False, 0, 0
        self.dir_name, self.dir_time, self.dir_size, self.dir_locked = [ ], [ ], [ ], [ ]
        self.file_name, self.file_time, self.file_size, self.file_fptr = [ ], [ ], [ ], [ ]
        self.curpath = self.getter(2)
        p = self.ext.func7()
        temp = kobj.unpack( kobj.recvauto(p) )
        self.ext.func0(p)
        self.self_name, self.self_time, self.self_size = str(temp[0], encoding="utf-8"), str(temp[1], encoding="utf-8"), kobj.decode( temp[2] )
        self.self_locked, self.self_subdir, self.self_subfile == (temp[3] == b"\x00"), kobj.decode( temp[4] ), kobj.decode( temp[5] )
        p = self.ext.func8()
        temp = kobj.unpack( kobj.recvauto(p) )
        self.ext.func0(p)
        if temp[0] != b"":
            self.dir_name, self.dir_time = str(temp[0], encoding="utf-8").split("\n"), str(temp[1], encoding="utf-8").split("\n")
            self.dir_size, self.dir_locked = [ kobj.decode( temp[2][8*x:8*x+8] ) for x in range(0, len( temp[2] ) // 8) ], [ (x == 0) for x in temp[3] ]
        p = self.ext.func9()
        temp = kobj.unpack( kobj.recvauto(p) )
        self.ext.func0(p)
        if temp[0] != b"":
            self.file_name, self.file_time = str(temp[0], encoding="utf-8").split("\n"), str(temp[1], encoding="utf-8").split("\n")
            self.file_size = [ kobj.decode( temp[2][8*x:8*x+8] ) for x in range(0, len( temp[2] ) // 8) ]
            self.file_fptr = [ kobj.decode( temp[3][8*x:8*x+8] ) for x in range(0, len( temp[3] ) // 8) ]

class mainclass(explorer.toolbox):
    def __init__(self, iswin, back):
        super().__init__("KVdrive", iswin)
        for i in os.listdir("../../_ST5_COMMON/iconpack/explorer/"):
            with open(f"../../_ST5_COMMON/iconpack/explorer/{i}", "rb") as f:
                i = i[1:i.rfind(".")] if i[0] == "_" else "." + i[:i.rfind(".")]
                self.icons[i] = f.read()
        self.iomgr, self.pwkf, self.tabpos, self.viewpos, self.paths, self.log0 = back, None, 0, 1, [""] * 5, "Boot Complete"
        self.status_working, self.status_viewsize, self.status_move, self.source_move = False, False, False, [ "", [ ] ]
        self.entry()
        self.sync()
        self.guiloop()

    def menubuilder(self):
        m0 = tk.Menu( self.mbar, tearoff=0, font=self.parms["font.3"] )
        m0.add_command( label=" Session Abort ", command=lambda:self.custom0(0, 0), font=self.parms["font.3"] )
        m0.add_command( label=" Session Recover ", command=lambda:self.custom0(0, 1), font=self.parms["font.3"] )
        m0.add_separator()
        m0.add_command( label=" Clear Search ", command=lambda:self.custom0(0, 2), font=self.parms["font.3"] )
        m0.add_command( label=" Clear Log ", command=lambda:self.custom0(0, 3), font=self.parms["font.3"] )
        m0.add_command( label=" Clear View ", command=lambda:self.custom0(0, 4), font=self.parms["font.3"] )
        m0.add_command( label=" Debug Info ", command=lambda:self.custom0(0, 5), font=self.parms["font.3"] )
        self.mbar.add_cascade(label=" Session ", menu=m0)

        m1 = tk.Menu( self.mbar, tearoff=0, font=self.parms["font.3"] )
        m1.add_command( label=" Import File ", command=lambda:self.custom0(1, 0), font=self.parms["font.3"] )
        m1.add_command( label=" Import Dir ", command=lambda:self.custom0(1, 1), font=self.parms["font.3"] )
        m1.add_command( label=" Export ", command=lambda:self.custom0(1, 2), font=self.parms["font.3"] )
        m1.add_separator()
        m1.add_command( label=" View Text ", command=lambda:self.custom0(1, 3), font=self.parms["font.3"] )
        m1.add_command( label=" View Picture ", command=lambda:self.custom0(1, 4), font=self.parms["font.3"] )
        m1.add_command( label=" View Binary ", command=lambda:self.custom0(1, 5), font=self.parms["font.3"] )
        self.mbar.add_cascade(label=" File ", menu=m1)

        m2 = tk.Menu( self.mbar, tearoff=0, font=self.parms["font.3"] )
        m2.add_command( label=" Recycle ", command=lambda:self.custom0(2, 0), font=self.parms["font.3"] )
        m2.add_command( label=" Rename ", command=lambda:self.custom0(2, 1), font=self.parms["font.3"] )
        m2.add_command( label=" Move ", command=lambda:self.custom0(2, 2), font=self.parms["font.3"] )
        m2.add_separator()
        m2.add_command( label=" Cancel Move ", command=lambda:self.custom0(2, 3), font=self.parms["font.3"] )
        self.mbar.add_cascade(label=" Manage ", menu=m2)

        m3 = tk.Menu( self.mbar, tearoff=0, font=self.parms["font.3"] )
        m3.add_command( label=" Delete ", command=lambda:self.custom0(3, 0), font=self.parms["font.3"] )
        m3.add_command( label=" Deep Lock ", command=lambda:self.custom0(3, 1), font=self.parms["font.3"] )
        m3.add_command( label=" New File ", command=lambda:self.custom0(3, 2), font=self.parms["font.3"] )
        m3.add_command( label=" New Dir ", command=lambda:self.custom0(3, 3), font=self.parms["font.3"] )
        m3.add_separator()
        m3.add_command( label=" Select All ", command=lambda:self.custom0(3, 4), font=self.parms["font.3"] )
        m3.add_command( label=" Toggle Viewsize ", command=lambda:self.custom0(3, 5), font=self.parms["font.3"] )
        self.mbar.add_cascade(label=" Control ", menu=m3)

        m4 = tk.Menu( self.mbar, tearoff=0, font=self.parms["font.3"] )
        m4.add_command( label=" Reset PWKF ", command=lambda:self.custom0(4, 0), font=self.parms["font.3"] )
        m4.add_command( label=" Export Account ", command=lambda:self.custom0(4, 1), font=self.parms["font.3"] )
        m4.add_separator()
        m4.add_command( label=" Restore Name ", command=lambda:self.custom0(4, 2), font=self.parms["font.3"] )
        m4.add_command( label=" Restore Data ", command=lambda:self.custom0(4, 3), font=self.parms["font.3"] )
        m4.add_command( label=" Restore Struct ", command=lambda:self.custom0(4, 4), font=self.parms["font.3"] )
        m4.add_command( label=" Check Cluster ", command=lambda:self.custom0(4, 5), font=self.parms["font.3"] )
        self.mbar.add_cascade(label=" Advanced ", menu=m4)

    def filter_sel(self, mode): # filter selection (ex del mv ren view)
        if mode == "ex":
            file, folder = [ ], ""
            for i in range( 0, len(self.names) ):
                if self.selected[i]:
                    if self.names[i][-1] == "/":
                        folder = self.names[i]
                    else:
                        file.append( self.names[i] )
            if self.iomgr.ext.func5(3) == 1:
                folder = self.iomgr.curpath + folder
            return file, folder
        elif mode == "del" or mode == "mv":
            if self.iomgr.ext.func5(3) == 0:
                return [ self.names[x] for x in filter( lambda x: self.selected[x], range( 0, len(self.names) ) ) ]
            else:
                return [ str(x) for x in filter( lambda x: self.selected[x], range( 0, len(self.names) ) ) ]
        elif mode == "ren":
            num, name = -1, ""
            for i in range( 0, len(self.names) ):
                if self.selected[i]:
                    num, name = i, self.names[i]
            return num, name
        elif mode == "view":
            return [ self.names[x] for x in filter( lambda x: self.selected[x] and self.names[x][-1] != "/" and self.sizes[x] < maxsize, range( 0, len(self.names) ) ) ]

    def sync(self): # update status
        if self.status_move:
            self.iomgr.order( "navigate", [self.iomgr.curpath] )
            temp = self.iomgr.getter(0)[0].split("\n")
            self.paths[0], self.names, self.sizes, self.locked, self.selected = self.iomgr.curpath, temp, [-1] * len(temp), [False] * len(temp), [False] * len(temp)
        else:
            if self.status_viewsize:
                self.iomgr.ext.func1(2)
                err = self.iomgr.order( "update", ["true"] )
            else:
                self.iomgr.ext.func1(3)
                err = self.iomgr.order( "update", ["false"] )
            if err != "":
                tkm.showerror(title="Update Fail", message=f" {err} ")
            self.iomgr.sync()
            self.paths[0], self.names, self.sizes = self.iomgr.curpath, self.iomgr.dir_name + self.iomgr.file_name, self.iomgr.dir_size + self.iomgr.file_size
            self.locked, self.selected = self.iomgr.dir_locked + [False] * len(self.iomgr.file_fptr), [False] * len(self.names)
        self.render(True, True, True, False, False, False, False)

    def resetpwkf(self): # pwkf reset callback
        self.pwkf.mwin.destroy()
        self.iomgr.setter([self.pwkf.pw, self.pwkf.kf, self.pwkf.hint], 1)
        self.pwkf.pw, self.pwkf.kf, self.pwkf.hint = b"12345678", b"01234567", b"89898989"
        err, self.pwkf = self.iomgr.order("reset", [ ]), None
        if err == "":
            tkm.showinfo(title="PWKF Reset", message=" New Password updated. ")
        else:
            tkm.showerror(title="PWKF Reset Fail", message=f" {err} ")

    def extendacc(self): # account extension callback
        self.pwkf.mwin.destroy()
        name = tks.askstring(title="Select Name", prompt="New Account ID : ")
        flag = "true" if tkm.askyesno(title="Lock Check", message=" Do you want to include Locked Dir? \n (yes : include Lock, no : Unlocked only) ") else "false"
        self.iomgr.setter([self.pwkf.pw, self.pwkf.kf, self.pwkf.hint], 1)
        self.pwkf.pw, self.pwkf.kf, self.pwkf.hint = b"12345678", b"01234567", b"89898989"
        err, self.pwkf = self.iomgr.order("extend", [name, flag]), None
        name = self.iomgr.getter(0)[0]
        if err == "":
            tkm.showinfo(title="Account Generated", message=f" new account at {name} ")
        else:
            tkm.showerror(title="Extend Fail", message=f" {err} ")

    def custom0(self, x, y):
        if x == 0:
            if y == 0:
                err = self.iomgr.order( "abort", ["true", "true", "true"] )
                if err != "":
                    tkm.showerror(title="Abort Fail", message=f" {err} ")

            elif y == 1:
                self.iomgr.order( "abort", ["true", "false", "false"] )
                self.iomgr.setter("//", 2)
                self.viewpos = 1
                self.sync()

            elif y == 2:
                self.search = [ ]
                self.render(False, False, False, True, False, False, False)

            elif y == 3:
                self.iomgr.order( "log", ["true"] )
                self.log1 = [ ]
                self.render(False, False, False, True, False, False, False)

            elif y == 4:
                self.paths[2], self.paths[3], self.paths[4], self.txtdata, self.picdata, self.bindata = "", "", "", "", b"", b""
                self.picdata_name, self.picdata_bin, self.picdata_num = None, None, None
                self.render(True, False, False, False, True, True, True)

            elif y == 5:
                err0 = self.iomgr.order( "debug", ["true"] )
                temp = "\n".join( self.iomgr.getter(0) )
                err1 = self.iomgr.order( "print", ["true"] )
                temp = temp + "\n\n" + self.iomgr.getter(0)[0]
                self.paths[2], self.txtdata = "", temp
                self.render(False, False, False, False, True, False, False)
                if err0 == "" and err1 == "":
                    tkm.showinfo(title="Debug Info", message=" Debug info generated at txt section. ")
                else:
                    tkm.showerror(title="Debug Fail", message=f" {err0} \n {err1} ")

        elif x == 1:
            if y == 0:
                paths = [ ]
                for i in tkf.askopenfiles(title="Select Files", initialdir="./_ST5_DATA/"):
                    paths.append( i.name.replace("\\", "/") )
                err = self.iomgr.order("imfile", paths)
                if err == "":
                    tkm.showinfo(title="Order Start", message=" Starting import files. \n This work takes time. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")

            elif y == 1:
                path = tkf.askdirectory(title="Select Dir", initialdir="./_ST5_DATA/").replace("\\", "/")
                err = self.iomgr.order( "imdir", [path] )
                if err == "":
                    tkm.showinfo(title="Order Start", message=" Starting import folder. \n This work takes time. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")

            elif y == 2:
                file, folder = self.filter_sel("ex")
                err = self.iomgr.order("exdir", [folder]) if file == [ ] else self.iomgr.order("exfile", file)
                if err == "":
                    tkm.showinfo(title="Order Start", message=" Starting export data. \n This work takes time. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")

            elif y == 3:
                temp = self.filter_sel("view")[0]
                err0 = self.iomgr.order("exbin", [temp])
                while self.iomgr.ext.func5(0) == 0:
                    time.sleep(0.1)
                err1 = self.iomgr.getter(3)
                if err0 == "" and err1 == "":
                    self.txtdata, self.paths[2] = str(self.iomgr.getter(1)[0], encoding="utf-8"), self.iomgr.curpath + temp
                    self.render(True, False, False, False, True, False, False)
                    tkm.showinfo(title="Txt View", message=" Text generated at txt section. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err0} \n {err1} ")

            elif y == 4:
                temp = self.filter_sel("view")
                err0 = self.iomgr.order("exbin", temp)
                while self.iomgr.ext.func5(0) == 0:
                    time.sleep(0.1)
                err1 = self.iomgr.getter(3)
                if err0 == "" and err1 == "":
                    self.picdata_name, self.picdata_bin, self.picdata_num = [self.iomgr.curpath + x for x in temp], self.iomgr.getter(1), 1
                    self.picdata, self.paths[3] = self.picdata_bin[0], self.picdata_name[0]
                    self.render(True, False, False, False, False, True, False)
                    tkm.showinfo(title="Pic View", message=" Picture generated at pic section. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err0} \n {err1} ")

            elif y == 5:
                temp = self.filter_sel("view")[0]
                err0 = self.iomgr.order("exbin", [temp])
                while self.iomgr.ext.func5(0) == 0:
                    time.sleep(0.1)
                err1 = self.iomgr.getter(3)
                if err0 == "" and err1 == "":
                    self.bindata, self.paths[4] = self.iomgr.getter(1)[0], self.iomgr.curpath + temp
                    self.render(True, False, False, False, False, False, True)
                    tkm.showinfo(title="Bin View", message=" Binary generated at bin section. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err0} \n {err1} ")

        elif x == 2:
            if y == 0:
                err = self.iomgr.order( "move", ["/_BIN/"] + self.filter_sel("mv") )
                if err == "":
                    tkm.showinfo(title="Order Start", message=" Starting move to bin. \n This work takes time. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")
                
            elif y == 1:
                num, name = self.filter_sel("ren")
                temp = tks.askstring(title="Enter New Name", prompt=f" {name} ")
                if temp == "":
                    temp = "NewName"
                if name[-1] == "/" and temp[-1] != "/":
                    temp = temp + "/"
                if "\\" in temp or ":" in temp or "*" in temp or "?" in temp or '"' in temp or "|" in temp or "<" in temp or ">" in temp:
                    if not tkm.askyesno(title="Dangerous Name", message=f' Name {temp} contains letter (\\:*?"|<>). \n Are you going to change name anyway? '):
                        return
                self.iomgr.setter([temp], 0)
                err = self.iomgr.order( "rename", [name] ) if self.iomgr.ext.func5(3) == 0 else self.iomgr.order( "rename", [str(num)] )
                if err == "":
                    tkm.showinfo(title="Order Start", message=" Starting data rename. \n This work takes time. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")
                    
            elif y == 2:
                self.source_move, self.status_move = [ self.iomgr.curpath, self.filter_sel("mv") ], True
                self.sync()

            elif y == 3:
                self.iomgr.curpath = self.source_move[0]
                self.source_move, self.status_move = None, False
                self.sync()

        elif x == 3:
            if y == 0:
                if self.iomgr.curpath == "/_BIN/":
                    err = self.iomgr.order( "delete", self.filter_sel("del") )
                    if err == "":
                        tkm.showinfo(title="Order Start", message=" Starting data delete. \n This work takes time. ")
                    else:
                        tkm.showerror(title="Order Fail", message=f" {err} ")
                    
            elif y == 1:
                if self.iomgr.curpath != "/" and self.iomgr.curpath != "/_BIN/":
                    flag = "true" if tkm.askyesno(title="Lock Config", message=f" Do you want to lock {self.iomgr.curpath} ? \n (yes : Lock, no : Unlock) ") else "false"
                    err = self.iomgr.order( "dirlock", [flag, ""] ) if self.iomgr.ext.func5(3) == 0 else self.iomgr.order( "dirlock", [flag, "-1"] )
                    if err == "":
                        tkm.showinfo(title="Order Start", message=" Starting deep lock. \n This work takes time. ")
                    else:
                        tkm.showerror(title="Order Fail", message=f" {err} ")

            elif y == 2:
                self.iomgr.setter( [b"Hello, world!"], 1 )
                err = self.iomgr.order( "imbin", ["NewFile"] )
                if err == "":
                    tkm.showinfo(title="Order Start", message=" Making new file. \n This work takes time. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")

            elif y == 3:
                err = self.iomgr.order( "dirnew", ["NewDir/"] )
                if err == "":
                    tkm.showinfo(title="Order Start", message=" Making new folder. \n This work takes time. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")

            elif y == 4:
                flag = True if False in self.selected else False
                for i in range( 0, len(self.names) ):
                    if self.iomgr.curpath != "/" or (self.names[i] != "_BIN/" and self.names[i] != "_BUF"):
                        self.selected[i] = flag
                self.render(False, False, True, False, False, False, False)

            elif y == 5:
                self.status_viewsize = not self.status_viewsize
                self.sync()

        elif x == 4:
            if y == 0:
                self.pwkf = pwkf(iswin, kaes.funcmode(), None)
                self.pwkf.entry(False, self.mwin, self.resetpwkf)

            elif y == 1:
                self.pwkf = pwkf(iswin, kaes.funcmode(), None)
                self.pwkf.entry(False, self.mwin, self.extendacc)

            elif y == 2:
                err = self.iomgr.order( "restore", ["rename"] )
                if err == "":
                    tkm.showinfo(title="Order Start", message=" Check & Restoring Names. \n This work takes time. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")
                while self.iomgr.ext.func5(0) == 0:
                    time.sleep(0.1)
                self.paths[2], self.txtdata = "", self.iomgr.getter(0)[0]
                self.render(False, False, False, False, True, False, False)
                tkm.showinfo(title="Order Clear", message=" Check result at txt section. ")

            elif y == 3:
                if self.iomgr.ext.func5(3) == 0:
                    tkm.showinfo(title="Not Supported", message=" Restoring Data is not supported by KV4adv cluster. ")
                else:
                    err = self.iomgr.order( "restore", ["rewrite"] )
                    if err == "":
                        tkm.showinfo(title="Order Start", message=" Check & Restoring Data. \n This work takes time. ")
                    else:
                        tkm.showerror(title="Order Fail", message=f" {err} ")
                    while self.iomgr.ext.func5(0) == 0:
                        time.sleep(0.1)
                    self.paths[2], self.txtdata = "", self.iomgr.getter(0)[0]
                    self.render(False, False, False, False, True, False, False)
                    tkm.showinfo(title="Order Clear", message=" Check result at txt section. ")

            elif y == 4:
                err = self.iomgr.order( "restore", ["rebuild"] )
                if err == "":
                    tkm.showinfo(title="Order Start", message=" Check & Restoring Cluster. \n This work takes time. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")
                while self.iomgr.ext.func5(0) == 0:
                    time.sleep(0.1)
                self.paths[2], self.txtdata = "", self.iomgr.getter(0)[0]
                self.render(False, False, False, False, True, False, False)
                tkm.showinfo(title="Order Clear", message=" Check result at txt section. ")

            elif y == 5:
                if self.iomgr.ext.func5(3) == 0:
                    tkm.showinfo(title="Not Supported", message=" Checking Cluster is not supported by KV4adv cluster. ")
                else:
                    err = self.iomgr.order( "check", [ ] )
                    if err == "":
                        tkm.showinfo(title="Order Start", message=" Checking Data. \n This work takes time. ")
                    else:
                        tkm.showerror(title="Order Fail", message=f" {err} ")
                    while self.iomgr.ext.func5(0) == 0:
                        time.sleep(0.1)
                    self.paths[2], self.txtdata = "", self.iomgr.getter(0)[0]
                    self.render(False, False, False, False, True, False, False)
                    tkm.showinfo(title="Order Clear", message=" Check result at txt section. ")
    
    def custom1(self, x):
        if x == 0:
            if self.status_move:
                err = self.iomgr.order( "move", [self.iomgr.curpath] + self.source_move[1] )
                if err == "":
                    tkm.showinfo(title="Order Start", message=" Starting data move. \n This work takes time. ")
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")
                self.custom0(2, 3)
            else:
                self.sync()
        elif x == 1:
            if self.iomgr.curpath.count("/") > 1:
                if self.status_move:
                    self.iomgr.curpath, self.viewpos = self.iomgr.curpath[:self.iomgr.curpath[:-1].rfind("/")+1], 1
                else:
                    self.iomgr.setter(self.iomgr.curpath[:self.iomgr.curpath[:-1].rfind("/")+1], 2)
                    self.viewpos = 1
                self.sync()
    
    def custom2(self, x):
        if self.names[x][-1] == "/":
            if self.status_move:
                self.iomgr.curpath, self.viewpos = self.iomgr.curpath + self.names[x], 1
            else:
                self.iomgr.setter(self.iomgr.curpath+self.names[x], 2)
                self.viewpos = 1
            self.sync()
        else:
            self.custom4(x)
    
    def custom3(self, x):
        if self.names[x][-1] == "/" and not self.status_move:
            if self.iomgr.curpath != "/" or self.names[x] != "_BIN/":
                self.locked[x] = not self.locked[x]
                flag = "true" if self.locked[x] else "false"
                err = self.iomgr.order( "dirlock", [ flag, self.names[x] ] ) if self.iomgr.ext.func5(3) == 0 else self.iomgr.order( "dirlock", [ flag, str(x) ] )
                if err == "":
                    time.sleep(0.5)
                    self.sync()
                else:
                    tkm.showerror(title="Order Fail", message=f" {err} ")
    
    def custom4(self, x):
        if self.iomgr.curpath != "/" or (self.names[x] != "_BIN/" and self.names[x] != "_BUF"):
            self.selected[x] = not self.selected[x]
            self.render(False, False, True, False, False, False, False)
    
    def custom5(self, x, y):
        if x == 0:
            self.iomgr.order( "search", [y] )
            self.search = self.iomgr.getter(0)[0].split("\n")
        elif x == 1:
            self.iomgr.order( "log", ["false"] )
            self.log1 = self.iomgr.getter(0)[0].split("\n")
        self.render(False, False, False, True, False, False, False)
    
    def custom6(self, x):
        temp = self.paths[2]
        if self.iomgr.curpath == temp[:temp.rfind("/")+1]:
            temp = temp[temp.rfind("/")+1:]
            self.iomgr.setter([bytes(x, encoding="utf-8")], 1)
            err = self.iomgr.order( "imbin", [temp] )
            if err == "":
                tkm.showinfo(title="Order Start", message=" Saving text data. \n This work takes time. ")
            else:
                tkm.showerror(title="Save Fail", message=f" {err} ")
        else:
            tkm.showerror(title="Save Fail", message=f" File not in the current folder. ")
    
    def custom7(self, x):
        if x == 0 and self.picdata_num > 0:
            self.picdata_num = self.picdata_num - 1
        elif x == 1 and self.picdata_num + 1 < len(self.picdata_bin):
            self.picdata_num = self.picdata_num + 1
        self.paths[3], self.picdata = self.picdata_name[self.picdata_num], self.picdata_bin[self.picdata_num]
        self.render(True, False, False, False, False, True, False)
    
    def custom8(self, x):
        temp = self.paths[4]
        if self.iomgr.curpath == temp[:temp.rfind("/")+1]:
            temp = temp[temp.rfind("/")+1:]
            self.iomgr.setter([x], 1)
            err = self.iomgr.order( "imbin", [temp] )
            if err == "":
                tkm.showinfo(title="Order Start", message=" Saving binary data. \n This work takes time. ")
            else:
                tkm.showerror(title="Save Fail", message=f" {err} ")
        else:
            tkm.showerror(title="Save Fail", message=f" File not in the current folder. ")
    
    def custom9(self):
        if self.status_working:
            if self.iomgr.ext.func5(0) == 1:
                err, self.log0, self.status_working = self.iomgr.getter(3), "Idle", False
                if err != "":
                    tkm.showerror(title="Async Error", message=f" {err} ")
                return True
        else:
            if self.iomgr.ext.func5(0) == 0:
                self.log0, self.status_working = "Working", True
                return True
        if self.status_move:
            self.log0 = "Selecting Move"
        return False

cfg = kobj.repath()
try:
    maxsize = max(int( cfg[1] ), 4)
except:
    maxsize = 25165824 # default 24MiB
cfg = kdb.toolbox()
with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
    cfg.read( f.read() )
if not os.path.exists("./_ST5_DATA/"):
    os.mkdir("./_ST5_DATA/")
iswin = (cfg.get("dev.os")[3] == "windows")
if os.path.exists("../../_ST5_COMMON/iconpack/"):
    w0 = iomgr(iswin)
    w1 = pwkf(iswin, kaes.funcmode(), w0)
    w1.entry(True, None, None)
    w1.mwin.mainloop()
    if w1.mwin == None:
        w2 = mainclass(iswin, w0)
    w0.order( "exit", [ ] )
    if os.path.exists(w1.local):
        shutil.rmtree(w1.local)
else:
    tkm.showinfo(title="No Package", message=" KVdrive requires package >>> common.iconpack <<<. \n Install dependent package and start again. ")
time.sleep(0.5)
