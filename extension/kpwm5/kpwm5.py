# test738 : extension.kpwm5

import os
import time
import base64

import tkinter as tk
import tkinter.filedialog as tkf
import tkinter.messagebox as tkm
import reader

import kobj
import kdb
import kcom
import kaes

def readmd(text, is64): # read txt, returns chap, title, content, remain
    chap, title, cont = "", [ ], [ ]
    idx = text.find(">")
    chap = str(base64.b64decode( text[text.find("<")+1:idx] ), encoding="utf-8") if is64 else text[text.find("<")+1:idx]
    text = text[idx+1:]
    while len(text) != 0:
        if text[0] == "<":
            break
        elif text[0] == "{":
            idx = text.find("}")
            if is64:
                title.append( str(base64.b64decode( text[1:idx] ), encoding="utf-8") )
            else:
                title.append( text[1:idx] )
            text = text[idx+1:]
            idx = text.find(")")
            if is64:
                cont.append( str(base64.b64decode( text[text.find("(")+1:idx] ), encoding="utf-8") )
            else:
                cont.append( text[text.find("(")+1:idx] )
            text = text[idx+1:]
        elif text[0] == "#":
            text = text[text.find("\n")+1:]
        else:
            text = text[1:]
    remain = text if "<" in text else ""
    return chap, title, cont, remain

def writemd(chap, title, cont, is64): # write md with chap, title, content
    title2, cont2 = [ ], [ ]
    chap = base64.b64encode( bytes(chap, encoding="utf-8") ).decode("ascii") if is64 else chap.replace("<", " ").replace(">", " ")
    temp = [f"<{chap}>", ""]
    for i in range( 0, len(title) ):
        if is64:
            title2.append( base64.b64encode( bytes(title[i], encoding="utf-8") ).decode("ascii") )
            cont2.append( base64.b64encode( bytes(cont[i], encoding="utf-8") ).decode("ascii") )
        else:
            title2.append( title[i].replace("{", " ").replace("}", " ") )
            cont2.append( cont[i].replace("(", " ").replace(")", " ") )
        temp.append(f"{{{title2[i]}}} ({cont2[i]})")
    return "\n".join(temp) + "\n"

class pwkf:
    def __init__(self, iswin, aes):
        self.mwin, self.comp, self.parms, self.path, self.hint, self.pw = None, dict(), {"color":"thistle1"}, "", "Not Selected", b""
        self.aes, self.data, self.kf = aes, b"", aes.basickey()
        if iswin: # windows
            self.parms["ms"], self.parms["f0"], self.parms["f1"], self.parms["e0"], self.parms["e1"], self.parms["e2"] = "450x200+200+100", ("맑은 고딕", 12), ("Consolas", 14), 36, 14, 34
        else: # linux
            self.parms["ms"], self.parms["f0"], self.parms["f1"], self.parms["e0"], self.parms["e1"], self.parms["e2"] = "500x240+200+100", ("맑은 고딕", 8), ("Consolas", 10), 30, 11, 28

    def entry(self, isboot, master, callback):
        def sel_data(): # select data file
            time.sleep(0.1)
            if isboot:
                try:
                    temp = tkf.askopenfile(title="Select Data", filetypes=[ ("Link", "*.txt"), ("Data", "*.webp"), ("All Files", "*.*") ], initialdir="./_ST5_DATA/").name
                    if "." in temp and temp[temp.rfind("."):] == ".txt":
                        with open(temp, "r", encoding="utf-8") as f:
                            temp = os.path.abspath( f.read().replace("\n", "") )
                    self.path, self.hint = temp, self.isvalid(temp)
                except Exception as e:
                    self.path, self.hint = "", f"err : {e}"
                self.comp["strvar0a"].set(self.path)
                self.comp["strvar0b"].set(self.hint)

        def sel_key_direct(): # select kf direct
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
            if not isboot:
                self.hint = self.comp["entry0b"].get()
            self.pw = bytes(self.comp["entry2"].get(), encoding="utf-8")
            msg = f" Are you sure to submit current status? \n PW {len(self.pw)} B (len {len(str(self.pw, encoding='utf-8'))}), KF {len(self.kf)} B "
            if tkm.askokcancel(title="Submit PWKF", message=msg):
                if isboot:
                    self.login()
                else:
                    callback()

        self.mwin = tk.Tk() if isboot else tk.Toplevel(master) # main window
        self.mwin.title("PWKF")
        self.mwin.geometry( self.parms["ms"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # path, hint
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["button0"] = tk.Button(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text=" . . . ", command=sel_data)
        self.comp["button0"].grid(row=0, column=0, padx=5, pady=5)
        self.comp["strvar0a"] = tk.StringVar()
        self.comp["strvar0a"].set(self.path)
        self.comp["entry0a"] = tk.Entry(self.comp["frame0"], font=self.parms["f1"], width=self.parms["e0"], textvariable=self.comp["strvar0a"], state="readonly")
        self.comp["entry0a"].grid(row=0, column=1, padx=5, pady=5)
        self.comp["label0"] = tk.Label(self.comp["frame0"], font=self.parms["f0"], bg=self.parms["color"], text="Hint")
        self.comp["label0"].grid(row=1, column=0, padx=5, pady=5)
        if isboot:
            self.comp["strvar0b"] = tk.StringVar()
            self.comp["strvar0b"].set(self.hint)
            self.comp["entry0b"] = tk.Entry(self.comp["frame0"], font=self.parms["f1"], width=self.parms["e0"], textvariable=self.comp["strvar0b"], state="readonly")
        else:
            self.comp["entry0b"] = tk.Entry( self.comp["frame0"], font=self.parms["f1"], width=self.parms["e0"] )
        self.comp["entry0b"].grid(row=1, column=1, padx=5, pady=5)

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
        self.comp["button2"] = tk.Button(self.comp["frame2"], font=self.parms["f0"], bg=self.parms["color"], text=" Go ", command=submit)
        self.comp["button2"].grid(row=0, column=2)

    def isvalid(self, path): # check if valid path
        with open(path, "rb") as f:
            self.data = f.read()
        self.aes.view(self.data)
        return self.aes.hint
    
    def login(self): # try login when isboot
        if self.path == "":
            time.sleep(0.1)
            self.path, self.data = f"./_ST5_DATA/{int(time.time())}.webp", b"<TmV3Tm90ZQ==> {TmV3TGluZQ==} (TmV3Q29udGVudA==)"
            self.pw, self.kf, self.hint = b"0000", self.aes.basickey(), "init pw : 0000, kf : bkf"
            tkm.showinfo(title="Temp DataChunk", message=f" New datachunk will be generated at {self.path} ")
            self.mwin.destroy()
        else:
            try:
                self.data = self.aes.decrypt(self.pw, self.kf, self.data)
                self.mwin.destroy()
            except Exception as e:
                tkm.showerror(title="Login Fail", message=f" {e} ")
    
class mainclass(reader.toolbox):
    def __init__(self, iswin, data):
        super().__init__("KPWM5", iswin, 1)
        self.hint, self.pw, self.kf, self.path, self.iswin = data.hint, data.pw, data.kf, data.path, iswin
        self.current, self.editable, self.aes, self.ui = [0, 0], [False, False, False], data.aes, None
        self.menus = [" PWKF ", " Sort ", "Import", "Export", "AddPage", "DelPage",
                    " Lock ", " Next ", " Info ", " Save ", "AddLine", "DelLine"]
        temp = str(data.data, encoding="utf-8")
        while "<" in temp:
            a, b, c, temp = readmd(temp, True)
            self.big.append(a)
            self.middle.append(b)
            self.small.append(c)
        if len(self.big) == 0:
            self.big.append("NewPage")
            self.middle.append( ["NewLine"] )
            self.small.append( ["NewContent"] )
        self.entry()
        self.guiloop()

    def custom0(self, x):
        if x == 0: # pwkf
            self.ui = pwkf(self.iswin, self.aes)
            self.ui.path, self.ui.hint, self.ui.kf = self.path, self.hint, self.kf
            self.ui.entry(False, self.mwin, self.setpwkf)
            self.ui.comp["strvar1"].set(f"kf {len(self.kf)} B")
        elif x == 1: # sort
            temp = [ ]
            for i in range( 0, len(self.big) ):
                buf = [ ]
                for j in range( 0, len( self.middle[i] ) ):
                    buf.append( ( self.middle[i][j], self.small[i][j] ) )
                buf.sort(key=lambda x:x[0])
                temp.append( ( self.big[i], [x[0] for x in buf], [x[1] for x in buf] ) )
            temp.sort(key=lambda x:x[0])
            self.big, self.middle, self.small = [ ], [ ], [ ]
            for i in temp:
                self.big.append( i[0] )
                self.middle.append( i[1] )
                self.small.append( i[2] )
        elif x == 2: # import
            if tkm.askokcancel(title="Import Check", message=" This work will reset your current data. "):
                path = tkf.askopenfile( title="", filetypes=[ ("Text Files", "*.txt"), ("All Files", "*.*") ] ).name
                with open(path, "r", encoding="utf-8") as f:
                    temp, self.big, self.middle, self.small = f.read(), [ ], [ ], [ ]
                while "<" in temp:
                    a, b, c, temp = readmd(temp, False)
                    self.big.append(a)
                    self.middle.append(b)
                    self.small.append(c)
                if len(self.big) == 0:
                    self.big.append("NewPage")
                    self.middle.append( ["NewLine"] )
                    self.small.append( ["NewContent"] )
                tkm.showinfo(title="Import Complete", message=f" Imported {len(self.big)} pages. ")
        elif x == 3: # export
            temp = [ ]
            for i in range( 0, len(self.big) ):
                temp.append( writemd(self.big[i], self.middle[i], self.small[i], False) )
            temp = "\n\n".join(temp)
            if tkm.askokcancel(title="Export Check", message=" This work will write your data in plain text. "):
                path = tkf.askdirectory(title="Select Export Dir", initialdir="./_ST5_DATA/").replace("\\", "/")
                if path[-1] != "/":
                    path = path + "/"
                with open(path + "export.txt", "w", encoding="utf-8") as f:
                    f.write(temp)
                tkm.showinfo(title="Export Complete", message=f" Exported at {path} (len {len(temp)}). ")
        elif x == 4: # addpage
            self.big.insert(self.current[0], "NewPage")
            self.middle.insert( self.current[0], ["NewLine"] )
            self.small.insert( self.current[0], ["NewContent"] )
        elif x == 5: # delpage
            num = self.current[0]
            if tkm.askokcancel(title="Del Check", message=f" Are you sure to delete page {num} ({self.big[num]})? "):
                del self.big[num]
                del self.middle[num]
                del self.small[num]
                self.current[0] = 0 if num == 0 else num - 1
        elif x == 6: # lock
            self.editable = [ not self.editable[0] ] * 3
        elif x == 7: # next
            if self.current[1] + 1 < len( self.middle[ self.current[0] ] ):
                self.current[1] = self.current[1] + 1
            elif self.current[0] + 1 < len(self.big):
                self.current[0], self.current[1] = self.current[0] + 1, 0
        elif x == 8: # info
            if tkm.askokcancel(title="Info Check", message=f" This work will show password of datachunk. "):
                msg = f" PW : {str(self.pw, encoding='utf-8')}, KF : {len(self.kf)} B \n Hint : {self.hint} \n Path : {self.path} \n PageLen : {len(self.big)}, Edit : {self.editable[0]} "
                tkm.showinfo(title=f"Pos ({self.current[0]}, {self.current[1]})", message=msg)
        elif x == 9: # save
            self.save()
        elif x == 10: # addline
            n0, n1 = self.current
            self.middle[n0].insert(n1, "NewLine")
            self.small[n0].insert(n1, "NewContent")
        elif x == 11: # delline
            n0, n1 = self.current
            if tkm.askokcancel(title="Del Check", message=f" Are you sure to delete line {n0}.{n1} ({self.middle[n0][n1]})? "):
                del self.middle[n0][n1]
                del self.small[n0][n1]
                self.current[1] = 0 if n1 == 0 else n1 - 1
        self.render(True, False)
    
    def save(self):
        temp = [ ]
        for i in range( 0, len(self.big) ):
            temp.append( writemd(self.big[i], self.middle[i], self.small[i], True) )
        temp = bytes("\n\n".join(temp), encoding="utf-8")
        self.aes.hint, self.aes.msg, self.aes.signkey = self.hint, "KPWM5 datachunk", ["", ""]
        temp = self.aes.encrypt(self.pw, self.kf, temp, 0)
        with open(self.path, "wb") as f:
            f.write(temp)
        tkm.showinfo(title="Data Saved", message=f" Data size {len(temp)} B at {self.path} ")
    
    def setpwkf(self):
        self.hint, self.pw, self.kf = self.ui.hint, self.ui.pw, self.ui.kf
        self.ui = None
        self.save()

kobj.repath()
cfg = kdb.toolbox()
with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
    cfg.read( f.read() )
iswin = (cfg.get("dev.os")[3] == "windows")
if not os.path.exists("./_ST5_DATA/"):
    os.mkdir("./_ST5_DATA/")
worker = pwkf( iswin, kaes.allmode() )
worker.entry(True, None, None)
worker.mwin.mainloop()
if worker.data != b"":
    worker = mainclass(iswin, worker)
time.sleep(0.5)
