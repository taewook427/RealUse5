# test709 : independent.kboom
# 리눅스 빌드 추가 : --hidden-import='PIL._tkinter_finder'

import os
import shutil
import time
import base64
import random
import threading

import tkinter
import tkinter.messagebox
import tkinter.filedialog
import reader

import kaes
import kcom
import kobj

class iomgr: # settings.webp IO layer
    def __init__(self):
        self.aesA, self.aesF, self.buffer, self.kfpath = kaes.allmode(), kaes.funcmode(), b"", "Basic Keyfile"
        self.pw, self.kf, self.hint, self.folder, self.file = b"0000", self.aesA.basickey(), "init pw : 0000, kf : bkf", [ ], dict()

    def boot(self): # when settings.webp is available
        with open("./settings.webp", "rb") as f:
            self.buffer = f.read()
        self.aesA.view(self.buffer)
        self.hint, self.folder, temp = self.aesA.hint, [ ], [ x.replace("\\", "/") for x in os.listdir("./") ]
        for i in temp:
            if os.path.isdir(i):
                if i[-1] != "/":
                    i = i + "/"
                self.folder.append(i)

    def login(self, pw, kf): # login to vault & update value
        self.pw, self.kf, self.file = pw, kf, dict()
        data = str(self.aesA.decrypt(self.pw, self.kf, self.buffer), encoding="utf-8").split("\n")
        self.buffer = b""
        for i in range(0, len(data) // 2):
            name, key = data[2 * i], base64.b64decode( data[2 * i + 1] )
            if os.path.exists("./" + name):
                self.file[name] = key

    def save(self): # save settings.webp
        data = [ ]
        for i in self.file:
            data.append(i)
            data.append( str(base64.b64encode( self.file[i] ), encoding="utf-8") )
        data, self.aesA.hint, self.aesA.msg, self.aesA.signkey = "\n".join(data), self.hint, "independent.kboom5", ["", ""]
        with open("./settings.webp", "wb") as f:
            f.write( self.aesA.encrypt(self.pw, self.kf, bytes(data, encoding="utf-8"), 0) )

    def send(self, name, portkey, port): # send kcom by folder/file
        key, name, node = self.file[name], os.path.abspath(name).replace("\\", "/"), kcom.node()
        node.port, node.close = port, 45
        node.send(key + bytes(name, encoding="utf-8"), portkey)

    def imfile(self, files, folder): # import files to folder(~/)
        if folder not in self.folder:
            raise Exception("invalid folder name")
        names = [ ]
        for i in files:
            i = i.replace("\\", "/")
            name, key = i[i.rfind("/")+1:], self.aesF.genrand(48)
            self.aesF.before, self.aesF.after = i, "./" + folder + name
            self.aesF.encrypt(key)
            self.aesF.before, self.aesF.after = "", ""
            self.file[folder + name] = key
            names.append(name)
        self.save()
        return names
    
    def exfile(self, exdir, names): # export files to exdir
        exdir = os.path.abspath(exdir).replace("\\", "/")
        if exdir[-1] != "/":
            exdir = exdir + "/"
        for i in names:
            key, self.aesF.before, self.aesF.after = self.file[i], "./" + i, exdir + i[i.find("/")+1:]
            self.aesF.decrypt(key)
            self.aesF.before, self.aesF.after = "", ""

    def imbin(self, data, name): # import binary
        key, self.aesF.before = self.aesF.genrand(48), data
        self.aesF.encrypt(key)
        data = self.aesF.after
        self.aesF.before, self.aesF.after = "", ""
        with open("./" + name, "wb") as f:
            f.write(data)
        self.file[name] = key
        self.save()

    def exbin(self, name): # export binary
        with open("./" + name, "rb") as f:
            data, key = f.read(), self.file[name]
        self.aesF.before = data
        self.aesF.decrypt(key)
        data = self.aesF.after
        self.aesF.before, self.aesF.after = "", ""
        return data

    def delete(self, name): # delete file
        os.remove("./" + name)
        del self.file[name]
        self.save()

class selector: # get pw, kf, hint
    def __init__(self, iswin, iolayer):
        self.parms, self.io, self.exit = { }, iolayer, False
        if iswin:
            # windows
            self.parms["mwin.size"], self.parms["but.font"], self.parms["entry.font"], self.parms["entry.width"] = "300x200+100+50", ("맑은 고딕", 12), ("맑은 고딕", 14), 23
            self.parms["entry.x"], self.parms["y1"], self.parms["y2"], self.parms["y3"] = 55, 50, 105, 150
        else:
            # linux
            self.parms["mwin.size"], self.parms["but.font"], self.parms["entry.font"], self.parms["entry.width"] = "600x400+200+100", ("맑은 고딕", 12), ("맑은 고딕", 14), 23
            self.parms["entry.x"], self.parms["y1"], self.parms["y2"], self.parms["y3"] = 95, 95, 185, 260

    def getparms(self, ismain, topgui): # set pw, kf, hint value
        def func0(): # reset kf (direct)
            time.sleep(0.1)
            try:
                temp = tkinter.filedialog.askopenfile(title='키 파일 선택').name
                with open(temp, "rb") as f:
                    self.io.kf = f.read()
                self.io.kfpath = temp
            except:
                self.io.kf, self.io.kfpath = self.io.aesA.basickey(), "Basic Keyfile"
            nonlocal strvar1
            strvar1.set(self.io.kfpath)

        def func1(): # reset kf (comm)
            time.sleep(0.1)
            sport, skey = kcom.unpack( ent1b.get() )
            node = kcom.node()
            node.port = sport
            data = node.recieve(skey)
            self.io.kfpath = str(data[48:], encoding="utf-8")
            with open(self.io.kfpath, "rb") as f:
                tgt = f.read()
            self.io.aesF.before = tgt
            self.io.aesF.decrypt( data[0:48] )
            self.io.kf = self.io.aesF.after
            self.io.aesF.before, self.io.aesF.after = "", ""
            nonlocal strvar1
            strvar1.set(self.io.kfpath)
            tkinter.messagebox.showinfo(title='Keyfile Recieved', message=f' path : {self.io.kfpath} \n size : {len(self.io.kf)} B ')

        def func2(): # select
            time.sleep(0.1)
            self.io.pw = bytes(ent3.get(), encoding="utf-8")
            if not ismain:
                self.io.hint = ent2.get()
                self.io.save()
                tkinter.messagebox.showinfo(title='New PW Saved', message=f' PW : {self.io.pw} \n KFlen : {len(self.io.kf)} \n hint : {self.io.hint} ')
            win.destroy()
            self.exit = True

        self.exit = False
        win = tkinter.Tk() if ismain else tkinter.Toplevel(topgui)
        win.title('KB5 password')
        win.geometry( self.parms["mwin.size"] )
        win.resizable(False, False)
        but0 = tkinter.Button(win, font=self.parms["but.font"], text=". . .", command=func0)
        but0.place(x=5, y=5) # kf reset
        but0b = tkinter.Button(win, font=self.parms["but.font"], text=". . .", command=func1)
        but0b.place(x=5, y=self.parms["y1"]) # kf reset remote
        strvar1 = tkinter.StringVar()
        strvar1.set(self.io.kfpath)
        ent1 = tkinter.Entry(win, textvariable=strvar1, font=self.parms["entry.font"], width=self.parms["entry.width"], state="readonly")
        ent1.place(x=self.parms["entry.x"], y=10) # kf path
        ent1b = tkinter.Entry(win, font=self.parms["entry.font"], width=self.parms["entry.width"])
        ent1b.place(x=self.parms["entry.x"], y=self.parms["y1"]+5) # kf path remote
        if ismain:
            lbl2 = tkinter.Label(win, font=self.parms["entry.font"], text=self.io.hint)
            lbl2.place(x=self.parms["entry.x"], y=self.parms["y2"]) # hint lbl
        else:
            ent2 = tkinter.Entry(win, font=self.parms["entry.font"], width=self.parms["entry.width"])
            ent2.place(x=self.parms["entry.x"], y=self.parms["y2"]+5) # hint entry
        lbl2b = tkinter.Label(win, font=self.parms["entry.font"], text="hint")
        lbl2b.place(x=5, y=self.parms["y2"]) # hint lbl
        but4 = tkinter.Button(win, font=self.parms["but.font"], text=" Go ", command=func2)
        but4.place(x=5, y=self.parms["y3"])
        ent3 = tkinter.Entry(win, font=self.parms["entry.font"], width=self.parms["entry.width"], show="*")
        ent3.place(x=self.parms["entry.x"], y=self.parms["y3"]+5) # pw input
        if ismain:
            win.mainloop()

class mainclass(reader.toolbox):
    def __init__(self, iswin):
        self.io = iomgr()
        if not os.path.exists("./settings.webp"):
            self.io.save()
        self.io.boot()
        self.sel = selector(iswin, self.io)

        flag = True
        while flag:
            self.sel.getparms(True, None)
            if self.sel.exit:
                pw, kf = self.sel.io.pw, self.sel.io.kf
                try:
                    self.io.login(pw, kf)
                    flag = False
                except Exception as e:
                    tkinter.messagebox.showinfo(title='Login Fail', message=f' Error occurred while login. \n {e} ')
            else:
                flag = False
        if self.sel.exit:
            self.start_gui()

    def start_gui(self):
        super().__init__("KB5 mainpage", iswin, 0)
        self.menus = ["PWreset", "      ", "      ", "      ", "      ", "boom!", "Import", "Export", "Delete", "View", "Send", "Edit"]
        for i in self.io.folder: # get txt/pic from folder
            self.big.append(i)
            temp0, temp1 = [ ], [ ]
            for j in os.listdir("./" + i): # check if txt/pic
                temp0.append(j)
                if len(j) > 4 and j[-4:] == ".txt":
                    temp1.append( str(self.io.exbin(i + j), encoding="utf-8") )
                else:
                    temp1.append(i + j)
            self.middle.append(temp0)
            self.small.append(temp1)
        self.current = [0, 0]
        self.entry()
        self.guiloop()

    def custom0(self, x):
        if x == 0: # pw reset
            if tkinter.messagebox.askokcancel(title='Change Password', message=f' Are you sure to change password? '):
                self.sel.getparms(False, self.mwin)
                
        elif x == 5: # boom
            if tkinter.messagebox.askokcancel(title='KB5 boom', message=f' Are you sure to delete all keyfiles? '):
                for i in self.io.folder:
                    shutil.rmtree("./" + i)
                with open("./settings.webp", "wb") as f:
                    f.write(b"cafebabe" * 10485760)
                self.mwin.destroy()

        elif x == 6: # import
            if tkinter.messagebox.askokcancel(title="Import Check", message=" This work will add files to the vault. "):
                i, files = self.current[0], [ ]
                for j in tkinter.filedialog.askopenfiles(title='가져올 파일들 선택'):
                    files.append(j.name)
                names = self.io.imfile( files, self.big[i] )
                self.middle[i] = self.middle[i] + names
                self.small[i] = self.small[i] + [self.big[i] + x for x in names]
                tkinter.messagebox.showinfo(title='Import Complete', message=f' {len(names)} files are added to vault. ')
                self.render(True, False)

        elif x == 7: # export
            if tkinter.messagebox.askokcancel(title="Export Check", message=" This work will make files from the vault. "):
                i, path = self.current[0], tkinter.filedialog.askdirectory(title="내보낼 폴더 선택")
                self.io.exfile( path, [ self.big[i] + x for x in self.middle[i] ] )
                tkinter.messagebox.showinfo(title='Export Complete', message=f' Files generated at {path} ')
        
        elif x == 8: # delete
            i, j = self.current
            name = self.big[i] + self.middle[i][j]
            if tkinter.messagebox.askokcancel(title='Keyfile Delete', message=f' Are you sure to delete following file? \n {name} '):
                self.io.delete(name)
                self.middle[i] = self.middle[i][0:j] + self.middle[i][j+1:]
                self.small[i] = self.small[i][0:j] + self.small[i][j+1:]
            self.render(True, False)

        elif x == 9: # view
            i = self.current[0]
            for j in range( 0, len( self.middle[i] ) ):
                name = self.big[i] + self.middle[i][j]
                if "." in name and name[name.rfind(".")+1:] in ["bmp", "png", "jpg", "jpeg", "gif", "webp"]:
                    self.small[i][j] = self.io.exbin(name)

        elif x == 10: # send
            i, j = self.current
            port, key = random.randrange(10000, 30000), self.io.aesA.genrand(4)
            addr = kcom.pack(port, key)
            temp = tkinter.Toplevel(self.mwin)
            temp.title('KB5 sending')
            temp.geometry("450x300+100+50")
            temp.resizable(False, False)
            lbl = tkinter.Label(temp, font=("Consolas", 14), text=f"sending file\n{self.big[i]}\n{self.middle[i][j]}")
            lbl.place(x=5, y=5)
            strvar = tkinter.StringVar()
            strvar.set(addr)
            ent = tkinter.Entry(temp, textvariable=strvar, font=("Consolas", 14), width=20, state="readonly")
            ent.place(x=5, y=220)
            self.mwin.update()
            thr = threading.Thread(target=self.io.send, args=(self.big[i] + self.middle[i][j], key, port))
            thr.start()

        elif x == 11: # edit
            self.editable = [False, False, True]
            self.render(True, False)

    def custom1(self, x):
        self.editable = [False, False, False]
        if x == 0: # go left
            i = self.current[0] + 1
            for j in range( 0, len( self.small[i] ) ):
                if type(self.small[i][j]) == bytes:
                    self.small[i][j] = self.big[i] + self.middle[i][j]

        elif x == 0: # go right
            i = self.current[0] - 1
            for j in range( 0, len( self.small[i] ) ):
                if type(self.small[i][j]) == bytes:
                    self.small[i][j] = self.big[i] + self.middle[i][j]

        elif x == 2: # renew
            i, j = self.current
            if len( self.middle[i][j] ) > 4 and self.middle[i][j][-4:] == ".txt":
                self.io.imbin( bytes(self.small[i][j], encoding="utf-8"), self.big[i] + self.middle[i][j] )

if __name__ == "__main__":
    iswin = True # iswin (windows T, linux F)
    kobj.repath()
    worker = mainclass(iswin)
    time.sleep(0.5)
