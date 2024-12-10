# test731 : extension.kviewer5
# 리눅스 빌드 추가 : --hidden-import='PIL._tkinter_finder'

import os
import time
import threading
import pyperclip

import tkinter as tk
import tkinter.messagebox as tkm
import tkinter.simpledialog as tks
import explorer

import kobj
import kdb

class submgr:
    def __init__(self, cfg, stg):
        self.iswin, self.desktop, self.auto, self.size = (cfg.get("dev.os")[3] == "windows"), cfg.get("path.desktop")[3].replace("\\", "/"), stg.get("auto")[3], stg.get("size")[3]
        if self.desktop[-1] != "/":
            self.desktop = self.desktop + "/"
        self.tpname, self.tppath, num, self.viewsize, self.lastmove, self.picname, self.picdata, self.picpos = [ ], [ ], 0, False, self.desktop, [ ], [ ], 0
        while f"{num}.name" in stg.name:
            self.tpname.append( stg.get(f"{num}.name")[3] )
            temp = stg.get(f"{num}.path")[3].replace("\\", "/")
            if temp[-1] != "/":
                temp = temp + "/"
            self.tppath.append(temp)
            num = num + 1
    
    def getdir(self, path): # path, viewsize -> name, size
        file, dir, name0, size0, name1, size1 = [ ], [ ], [ ], [ ], [ ], [ ]
        for i in os.listdir(path):
            if os.path.isdir(path + i):
                if i[-1] != "/":
                    i = i + "/"
                dir.append(i)
            else:
                file.append(i)
        file.sort()
        dir.sort()
        num = len(dir)
        temp, ret = [0] * num, [0] * num
        for i in range(0, num):
            ret[i] = [-1]
            if self.viewsize:
                temp[i] = threading.Thread( target=self.size_sub, args=( path + dir[i], ret[i] ) )
                temp[i].start()
        for i in file:
            name1.append(i)
            size1.append( os.path.getsize(path + i) )
        for i in range(0, num):
            if self.viewsize:
                temp[i].join()
            name0.append( dir[i] )
            size0.append( ret[i][0] )
        return name0 + name1, size0 + size1

    def getinfo(self, path): # get file/dir info
        size, fnum, dnum, thr0, thr1, ret0, ret1 = 0, 0, 1, [ ], [ ], [ ], [ ]
        if os.path.isdir(path):
            for i in os.listdir(path):
                if os.path.isdir(path + i):
                    if i[-1] != "/":
                        i = i + "/"
                    num = len(thr0)
                    ret0.append( [0] )
                    temp = threading.Thread( target=self.size_sub, args=( path + i, ret0[num] ) )
                    temp.start()
                    thr0.append(temp)
                    ret1.append( [0, 0] )
                    temp = threading.Thread( target=self.num_sub, args=( path + i, ret1[num] ) )
                    temp.start()
                    thr1.append(temp)
                else:
                    size = size + os.path.getsize(path + i)
                    fnum = fnum + 1
            for i in range( 0, len(thr0) ):
                thr0[i].join()
                size = size + ret0[i][0]
                thr1[i].join()
                fnum, dnum = fnum + ret1[i][0], dnum + ret1[i][1]
        else:
            size, fnum, dnum = os.path.getsize(path), 1, 0
        if size < 1024:
            st = ""
        elif size < 1048576:
            st = f"({size/1024:.1f} KiB)"
        elif size < 1073741824:
            st = f"({size/1048576:.1f} MiB)"
        else:
            st = f"({size/1073741824:.1f} GiB)"
        mt = time.strftime( "%Y.%m.%d;%H:%M:%S", time.localtime( os.path.getmtime(path) ) )
        return f" {path} \n file {fnum}, dir {dnum}, size {size} B {st} \n modtime {mt} "

    def search(self, path, name): # find file/dir having name
        name, fres, thr, ret = name.lower(), [ ], [ ], [ ]
        for i in os.listdir(path):
            if os.path.isdir(path + i):
                if i[-1] != "/":
                    i = i + "/"
                num = len(thr)
                ret.append( [ ] )
                temp = threading.Thread( target=self.name_sub, args=( path + i, name, ret[num] ) )
                temp.start()
                thr.append(temp)
            else:
                if name in i.lower():
                    fres.append(path + i)
        for i in range( 0, len(thr) ):
            thr[i].join()
            for j in ret[i]:
                fres.append(j)
        return fres

    def size_sub(self, path, ret): # get size
        size = 0
        for i in os.listdir(path):
            if os.path.isdir(path + i):
                temp = [0]
                if i[-1] != "/":
                    i = i + "/"
                self.size_sub(path + i, temp)
                size = size + temp[0]
            else:
                size = size + os.path.getsize(path + i)
        ret[0] = size

    def num_sub(self, path, ret): # get file/dir num
        file, dir = 0, 1
        for i in os.listdir(path):
            if os.path.isdir(path + i):
                temp = [0, 0]
                if i[-1] != "/":
                    i = i + "/"
                self.num_sub(path + i, temp)
                file, dir = file + temp[0], dir + temp[1]
            else:
                file = file + 1
        ret[0], ret[1] = file, dir

    def name_sub(self, path, name, ret): # find name
        if name in path.lower():
            ret.append(path)
        for i in os.listdir(path):
            if os.path.isdir(path + i):
                temp = [ ]
                if i[-1] != "/":
                    i = i + "/"
                self.name_sub(path + i, name, temp)
                for i in temp:
                    ret.append(i)
            else:
                if name in i.lower():
                    ret.append(path + i)

    def openpath(self, path): # open file/dir
        if self.iswin:
            os.startfile( path.replace("/", "\\") )
        else:
            os.system(f"open {path}")

    def execute(self, path): # execute file
        temp = os.getcwd()
        os.chdir( path[ :path.rfind("/") ] )
        if self.iswin:
            os.startfile( path.replace("/", "\\") )
        else:
            os.system(f"gnome-terminal -- ./{path[path.rfind("/")+1:]}")
        os.chdir(temp)

    def copy_sub(self, data): # copy data to clipboard
        try:
            pyperclip.copy(data)
        except:
            tkm.showinfo(title="Copy Fail", message=" Cannot use clipboard in current system. Run one of the following commands. \n sudo apt-get install xsel \n sudo apt-get install xclip ")

class mainclass(explorer.toolbox):
    def __init__(self, sub):
        self.sub = sub
        super().__init__("Kviewer5", self.sub.iswin)
        for i in os.listdir("../../_ST5_COMMON/iconpack/explorer/"):
            with open(f"../../_ST5_COMMON/iconpack/explorer/{i}", "rb") as f:
                i = i[1:i.rfind(".")] if i[0] == "_" else "." + i[:i.rfind(".")]
                self.icons[i] = f.read()
        self.tabpos, self.viewpos = 0, 1
        self.entry()
        self.paths, self.sub.viewsize = [self.sub.desktop, "", "", "", ""], (self.sub.desktop.count("/") > self.sub.auto)
        self.names, self.sizes = self.sub.getdir(self.sub.desktop)
        self.locked, self.selected = [False] * len(self.names), [False] * len(self.names)
        self.guiloop()

    def menubuilder(self):
        num = len(self.sub.tpname)
        m0 = tk.Menu( self.mbar, tearoff=0, font=self.parms["font.3"] )
        for i in range(0, num):
            m0.add_command( label=self.sub.tpname[i], command=lambda:self.custom0(0, i), font=self.parms["font.3"] )
        m0.add_separator()
        m0.add_command( label=" Manual Move ", command=lambda:self.custom0(0, num), font=self.parms["font.3"] )
        self.mbar.add_cascade(label=" Teleport ", menu=m0)

        m1 = tk.Menu( self.mbar, tearoff=0, font=self.parms["font.3"] )
        m1.add_command( label=" Open Current Dir ", command=lambda:self.custom0(1, 0), font=self.parms["font.3"] )
        m1.add_separator()
        m1.add_command( label=" View Info ", command=lambda:self.custom0(1, 1), font=self.parms["font.3"] )
        m1.add_command( label=" Copy Path (windows) ", command=lambda:self.custom0(1, 2), font=self.parms["font.3"] )
        m1.add_command( label=" Copy Path (linux) ", command=lambda:self.custom0(1, 3), font=self.parms["font.3"] )
        self.mbar.add_cascade(label=" Path ", menu=m1)

        m2 = tk.Menu( self.mbar, tearoff=0, font=self.parms["font.3"] )
        m2.add_command( label=" View Text ", command=lambda:self.custom0(2, 0), font=self.parms["font.3"] )
        m2.add_command( label=" View Pictures ", command=lambda:self.custom0(2, 1), font=self.parms["font.3"] )
        m2.add_command( label=" View Bytes ", command=lambda:self.custom0(2, 2), font=self.parms["font.3"] )
        m2.add_separator()
        m2.add_command( label=" Delete View ", command=lambda:self.custom0(2, 3), font=self.parms["font.3"] )
        self.mbar.add_cascade(label=" FileView ", menu=m2)

        m3 = tk.Menu( self.mbar, tearoff=0, font=self.parms["font.3"] )
        m3.add_command( label=" Toggle ViewSize ", command=lambda:self.custom0(3, 0), font=self.parms["font.3"] )
        m3.add_command( label=" Execute File ", command=lambda:self.custom0(3, 1), font=self.parms["font.3"] )
        m3.add_separator()
        m3.add_command( label=" Select All ", command=lambda:self.custom0(3, 2), font=self.parms["font.3"] )
        m3.add_command( label=" Select None ", command=lambda:self.custom0(3, 3), font=self.parms["font.3"] )
        self.mbar.add_cascade(label=" Functions ", menu=m3)

    def custom0(self, x, y):
        if x == 0:
            if y == len(self.sub.tpname):
                tgt = tks.askstring("Manual Move", "Enter path to move.")
                if tgt == None:
                    tgt = self.paths[0]
            else:
                tgt = self.sub.tppath[y]
            self.sub.lastmove, self.paths[0], self.sub.viewsize = self.paths[0], tgt, (tgt.count("/") > self.sub.auto)
            self.names, self.sizes = self.sub.getdir( self.paths[0] )
            self.viewpos, self.locked, self.selected = 1, [False] * len(self.names), [False] * len(self.names)
            self.render(True, True, True, False, False, False, False)

        elif x == 1:
            if y == 0:
                self.sub.openpath( self.paths[0] )
            elif y == 1:
                tkm.showinfo( title="Info", message=self.sub.getinfo( self.get_fsel() ) )
            elif y == 2:
                self.sub.copy_sub( self.get_fsel().replace("/", "\\") )
            elif y == 3:
                self.sub.copy_sub( self.get_fsel() )

        elif x == 2:
            if y == 0:
                temp = self.get_fsel()
                if temp[-1] != "/":
                    if os.path.getsize(temp) < self.sub.size:
                        self.paths[2] = temp
                        with open(temp, "rb") as f:
                            self.txtdata = str(f.read(), encoding="utf-8")
                self.render(True, False, False, False, True, False, False)
            elif y == 1:
                self.sub.picname, self.sub.picdata, self.sub.picpos = [ ], [ ], 0
                for i in range( 0, len(self.names) ):
                    if self.selected[i] and self.names[i][-1] != "/":
                        nm = self.paths[0] + self.names[i]
                        if os.path.getsize(nm) < self.sub.size:
                            with open(nm, "rb") as f:
                                self.sub.picname.append(nm)
                                self.sub.picdata.append( f.read() )
                self.paths[3], self.picdata = self.sub.picname[0], self.sub.picdata[0]
                self.render(True, False, False, False, False, True, False)
            elif y == 2:
                temp = self.get_fsel()
                if temp[-1] != "/":
                    if os.path.getsize(temp) < self.sub.size:
                        self.paths[4] = temp
                        with open(temp, "rb") as f:
                            self.bindata = f.read()
                self.render(True, False, False, False, False, False, True)
            elif y == 3:
                self.paths[2], self.paths[3], self.paths[4] = "", "", ""
                self.txtdata, self.bindata, self.picdata, self.sub.picname, self.sub.picdata, self.sub.picpos = "", b"", b"", [ ], [ ], 0
                self.render(True, False, False, False, True, True, True)

        elif x == 3:
            if y == 0:
                self.sub.viewsize = not self.sub.viewsize
                self.names, self.sizes = self.sub.getdir( self.paths[0] )
                self.render(True, True, True, False, False, False, False)
            elif y == 1:
                temp = self.get_fsel()
                if temp[-1] != "/":
                    self.sub.execute(temp)
            elif y == 2:
                self.selected = [True] * len(self.names)
                self.render(False, False, True, False, False, False, False)
            elif y == 3:
                self.selected = [False] * len(self.names)
                self.render(False, False, True, False, False, False, False)

    def custom1(self, x):
        if x == 0:
            self.paths[0], self.sub.viewsize = self.sub.lastmove, (self.sub.lastmove.count("/") > self.sub.auto)
            self.names, self.sizes = self.sub.getdir( self.paths[0] )
            self.viewpos, self.locked, self.selected = 1, [False] * len(self.names), [False] * len(self.names)
            self.render(True, True, True, False, False, False, False)
        elif x == 1:
            self.viewpos = 1
            if self.paths[0].count("/") > 1:
                self.sub.lastmove, self.paths[0] = self.paths[0], self.paths[0][:self.paths[0][:-1].rfind("/")+1]
                self.sub.viewsize = self.paths[0].count("/") > self.sub.auto
                self.names, self.sizes = self.sub.getdir( self.paths[0] )
                self.locked, self.selected = [False] * len(self.names), [False] * len(self.names)
                self.render(True, True, True, False, False, False, False)

    def custom2(self, x):
        if self.names[x][-1] == "/":
            self.viewpos, self.sub.lastmove, self.paths[0] = 1, self.paths[0], self.paths[0] + self.names[x]
            self.sub.viewsize = self.paths[0].count("/") > self.sub.auto
            self.names, self.sizes = self.sub.getdir( self.paths[0] )
            self.locked, self.selected = [False] * len(self.names), [False] * len(self.names)
            self.render(True, True, True, False, False, False, False)
        else:
            self.sub.openpath( self.paths[0] + self.names[x] )

    def custom4(self, x):
        self.selected[x] = not self.selected[x]
        self.render(False, False, True, False, False, False, False)

    def custom5(self, x, y):
        if x == 0 and len(y) > 1:
            self.paths[1] = self.paths[0]
            self.search = self.sub.search(self.paths[0], y)
            self.render(False, False, False, True, False, False, False)
        elif x == 1:
            self.paths[1], self.search = "", [ ]
            self.render(False, False, False, True, False, False, False)

    def custom6(self, x):
        if self.paths[2] != "":
            if tkm.askokcancel(title="Save Text", message=f" Are you sure to save text (len {len(x)}) at {self.paths[2]}? "):
                with open(self.paths[2], "wb") as f:
                    f.write( bytes(x, encoding="utf-8") )

    def custom7(self, x):
        if x == 0 and self.sub.picpos > 0:
            self.sub.picpos = self.sub.picpos - 1
        elif x == 1 and self.sub.picpos < len(self.sub.picname) - 1:
            self.sub.picpos = self.sub.picpos + 1
        self.paths[3], self.picdata = self.sub.picname[self.sub.picpos], self.sub.picdata[self.sub.picpos]
        self.render(True, False, False, False, False, True, False)

    def custom8(self, x):
        if self.paths[4] != "":
            if tkm.askokcancel(title="Save Bytes", message=f" Are you sure to save bytes (len {len(x)}) at {self.paths[4]}? "):
                with open(self.paths[4], "wb") as f:
                    f.write(x)

    def custom9(self):
        return False

    def get_fsel(self): # get first selected path
        tgt = ""
        for i in range( 0, len(self.selected) ):
            if self.selected[i]:
                tgt = self.paths[0] + self.names[i]
                break
        if tgt == "":
            tgt = self.paths[0]
        return tgt

if __name__ == "__main__":
    kobj.repath()
    worker0 = kdb.toolbox()
    with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
        worker0.read( f.read() )
    worker1 = kdb.toolbox()
    with open("./settings.txt", "r", encoding="utf-8") as f:
        worker1.read( f.read() )
    if os.path.exists("../../_ST5_COMMON/iconpack/"):
        worker2 = submgr(worker0, worker1)
        worker3 = mainclass(worker2)
    else:
        tkm.showinfo(title="No Package", message=" Kviewer5 requires package >>> common.iconpack <<<. \n Install dependent package and start again. ")
    time.sleep(0.5)
