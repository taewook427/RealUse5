# test714 : independent.starter5

import os
import time

import tkinter
import tkinter.messagebox
import blocksel

import kdb
import kobj

class procmgr: # process manager
    def __init__(self):
        worker = kdb.toolbox()
        with open("./_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
            worker.read( f.read() )
        self.iswin = worker.get("dev.os")[3] == "windows"
        self.name, self.path, self.exe, self.icon, self.sub, self.info = [ ], [ ], [ ], [ ], dict(), dict()
        self.setup("./_ST5_EXTENSION/")
        self.setup("./_ST5_COMMON/")

    def setup(self, path): # update name, path, exe, icon, sub
        for i in [ x[:-1] if x[-1] == "/" else x for x in os.listdir(path) ]:
            self.name.append(i)
            self.path.append(path + i + "/")
            exe, icon = None, None
            for j in [ x[:-1] if x[-1] == "/" else x for x in os.listdir(path + i) ]:
                if os.path.isdir(path + i + "/" + j):
                    j = j + "/"
                if "_ST5_EXE" in j:
                    exe = path + i + "/" + j
                elif "_ST5_ICON" in j:
                    icon = path + i + "/" + j
                self.sub[path + i + "/" + j] = True
            self.exe.append(exe)
            self.icon.append(icon)

    def openfile(self, path): # execute file (cwd matched)
        cwd, path = os.path.abspath("./"), os.path.abspath(path).replace("\\", "/")
        os.chdir( path[ 0:path.rfind("/") ] )
        if self.iswin:
            os.startfile( path.replace("/", "\\") )
        else:
            os.system(f"gnome-terminal -- ./{path[path.rfind("/")+1:]}")
        os.chdir(cwd)

    def openfolder(self, path): # open folder
        if self.iswin:
            os.startfile( path.replace("/", "\\") )
        else:
            os.system(f"open {path}")

    def openinfo(self, path): # get version info
        if path in self.info:
            return self.info[path]
        else:
            worker = kdb.toolbox()
            with open(path, "r", encoding="utf-8") as f:
                worker.read( f.read() )
            temp = [ worker.get("name")[3], worker.get("version")[3], worker.get("text")[3], worker.get("release")[3], worker.get("download")[3] ]
            self.info[path] = temp
            return temp
        
class mainclass(blocksel.toolbox):
    def __init__(self):
        self.proc = procmgr()
        super().__init__("Starter5", self.proc.iswin)
        addnum = 6 if len(self.proc.name) == 0 else 5 - (len(self.proc.name) - 1) % 6
        for i in range( 0, len(self.proc.name) ):
            self.txts.append( self.proc.name[i] )
            if self.proc.icon[i] == None:
                self.pics.append("./_ST5_ICON.png")
            else:
                self.pics.append( self.proc.icon[i] )
        for i in range(0, addnum):
            self.txts.append("")
            self.pics.append("./_ST5_ICON.png")
        self.curpos, self.upos, self.umsg = 1, 0, ["Mode 0 : Package Execute", "Mode 1 : Open Folder", "Mode 2 : Open Data Storage", "Mode 3 : Show Info"]

    def custom0(self, x):
        if x < len(self.proc.name):
            if self.upos == 0:
                if self.proc.exe[x] == None:
                    tkinter.messagebox.showinfo(title="Cannot Execute", message=f" No executable file in {self.proc.name[x]}. ")
                else:
                    self.proc.openfile( self.proc.exe[x] )
            elif self.upos == 1:
                self.proc.openfolder( self.proc.path[x] )
            elif self.upos == 2:
                if self.proc.path[x] + "_ST5_DATA/" in self.proc.sub:
                    self.proc.openfolder(self.proc.path[x] + "_ST5_DATA/")
                else:
                    tkinter.messagebox.showinfo(title="Cannot Open", message=f" No data storage in {self.proc.name[x]}. ")
            elif self.upos == 3:
                a, b, c, d, e = self.proc.openinfo(self.proc.path[x] + "_ST5_VERSION.txt")
                tkinter.messagebox.showinfo(title=f"Info of {self.proc.name[x]}", message=f" Package : {a} ver{b} \n Info : {c} \n Release : {d}, Download : {e} ")

kobj.repath()
try:
    worker = mainclass()
    flag = True
except Exception as e:
    tkinter.messagebox.showerror(title="Setup Fail", message=f" Error occurred while setup. \n {e} ")
    flag = False
if flag:
    worker.entry()
    worker.guiloop()
time.sleep(0.5)
