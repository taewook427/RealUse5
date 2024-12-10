# test732 : extension.kdoc

import os
import time

import tkinter.filedialog
import tkinter.messagebox
import reader

import kobj
import kdb

def readmd(text): # read txt, returns chap, title, content, remain
    chap, title, cont = "", [ ], [ ]
    idx = text.find(">")
    chap = text[text.find("<")+1:idx]
    text = text[idx+1:]
    while len(text) != 0:
        if text[0] == "<":
            break
        elif text[0] == "{":
            idx = text.find("}")
            title.append( text[1:idx] )
            text = text[idx+1:]
            idx = text.find(")")
            cont.append( text[text.find("(")+1:idx] )
            text = text[idx+1:]
        elif text[0] == "#":
            text = text[text.find("\n")+1:]
        else:
            text = text[1:]
    remain = text if "<" in text else ""
    return chap, title, cont, remain

def writemd(chap, title, cont): # write md with chap, title, content
    chap = chap.replace("<", " ").replace(">", " ")
    temp = [f"<{chap}>", ""]
    for i in range( 0, len(title) ):
        title[i] = title[i].replace("{", " ").replace("}", " ")
        cont[i] = cont[i].replace("(", " ").replace(")", " ")
        temp.append(f"{{{title[i]}}} ({cont[i]})")
    return "\n".join(temp) + "\n"

class page: # class for sort
    def __init__(self, chap, title, cont):
        self.chap, self.title, self.cont = chap, title, cont

class mainclass(reader.toolbox):
    def __init__(self):
        worker = kdb.toolbox()
        with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
            worker.read( f.read() )
        super().__init__("Kdoc", worker.get("dev.os")[3] == "windows", 1)
        self.path = tkinter.filedialog.askopenfile(filetypes=( ("Text File", "*.txt"), ("All File", "*.*") ), initialdir="./_ST5_DATA/").name
        with open(self.path, "r", encoding="utf-8") as f:
            temp = f.read()
        self.menus = ["Save", "Add Page", "Del Page", "Sort Page", "Move Page Left", "Move Page Right",
                    "Lock", "Add Line", "Del Line", "Sort Line", "Move Line Up", "Move Line Down"]
        self.cache, self.current, self.editable = [ ], [0, 0], [False, False, False]
        while "<" in temp:
            a, b, c, temp = readmd(temp)
            self.cache.append( writemd(a, b, c) )
            self.big.append(a)
            self.middle.append(b)
            self.small.append(c)
        if len(self.cache) == 0:
            self.cache.append("<page> {line} (data)")
            self.big.append("page")
            self.middle.append( ["line"] )
            self.small.append( ["data"] )
        self.entry()
        self.guiloop()

    def custom0(self, x):
        if x == 0: # save
            temp = "\n\n".join(self.cache)
            with open(self.path, "w", encoding="utf-8") as f:
                f.write(temp)
            tkinter.messagebox.showinfo(title="Text Save", message=f" Save markdown (len {len(temp)}) at \n {self.path} ")
        elif x == 1: # add page
            num = self.current[0]
            self.cache.insert(num, "<page> {line} (data)")
            self.big.insert(num, "page")
            self.middle.insert( num, ["line"] )
            self.small.insert( num, ["data"] )
        elif x == 2: # del page
            num = self.current[0]
            if tkinter.messagebox.askokcancel(title="Del Check", message=f" Are you sure to delete page {num} ({self.big[num]})? "):
                del self.cache[num]
                del self.big[num]
                del self.middle[num]
                del self.small[num]
                self.current[0] = 0 if num == 0 else num - 1
        elif x == 3: # sort page
            temp = [0] * len(self.cache)
            for i in range( 0, len(temp) ):
                temp[i] = page( self.big[i], self.middle[i], self.small[i] )
            temp = sorted(temp, key=lambda x:x.chap)
            for i in range( 0, len(temp) ):
                self.cache[i] = writemd(temp[i].chap, temp[i].title, temp[i].cont)
                self.big[i] = temp[i].chap
                self.middle[i] = temp[i].title
                self.small[i] = temp[i].cont
        elif x == 4: # move page left
            num = self.current[0]
            if num != 0:
                self.cache[num-1], self.cache[num] = self.cache[num], self.cache[num-1]
                self.big[num-1], self.big[num] = self.big[num], self.big[num-1]
                self.middle[num-1], self.middle[num] = self.middle[num], self.middle[num-1]
                self.small[num-1], self.small[num] = self.small[num], self.small[num-1]
        elif x == 5: # move page right
            num = self.current[0]
            if num != len(self.cache) - 1:
                self.cache[num+1], self.cache[num] = self.cache[num], self.cache[num+1]
                self.big[num+1], self.big[num] = self.big[num], self.big[num+1]
                self.middle[num+1], self.middle[num] = self.middle[num], self.middle[num+1]
                self.small[num+1], self.small[num] = self.small[num], self.small[num+1]
        elif x == 6: # lock
            self.editable = [ not self.editable[0] ] * 3
        elif x == 7: # add line
            n0, n1 = self.current
            self.middle[n0].insert(n1, "line")
            self.small[n0].insert(n1, "data")
            self.cache[n0] = writemd( self.big[n0], self.middle[n0], self.small[n0] )
        elif x == 8: # del line
            n0, n1 = self.current
            if tkinter.messagebox.askokcancel(title="Del Check", message=f" Are you sure to delete line {n0}.{n1} ({self.middle[n0][n1]})? "):
                del self.middle[n0][n1]
                del self.small[n0][n1]
                self.cache[n0] = writemd( self.big[n0], self.middle[n0], self.small[n0] )
                self.current[1] = 0 if n1 == 0 else n1 - 1
        elif x == 9: # sort line
            n0, n1 = self.current
            temp = [0] * len( self.middle[n0] )
            for i in range( 0, len(temp) ):
                temp[i] = page( "", self.middle[n0][i], self.small[n0][i] )
            temp = sorted(temp, key=lambda x:x.title)
            for i in range( 0, len(temp) ):
                self.middle[n0][i] = temp[i].title
                self.small[n0][i] = temp[i].cont
            self.cache[n0] = writemd( self.big[n0], self.middle[n0], self.small[n0] )
        elif x == 10: # move line up
            n0, n1 = self.current
            if n1 != 0:
                self.middle[n0][n1-1], self.middle[n0][n1] = self.middle[n0][n1], self.middle[n0][n1-1]
                self.small[n0][n1-1], self.small[n0][n1] = self.small[n0][n1], self.small[n0][n1-1]
                self.cache[n0] = writemd( self.big[n0], self.middle[n0], self.small[n0] )
        elif x == 11: # move line down
            n0, n1 = self.current
            if n1 != len( self.middle[n0] ) - 1:
                self.middle[n0][n1+1], self.middle[n0][n1] = self.middle[n0][n1], self.middle[n0][n1+1]
                self.small[n0][n1+1], self.small[n0][n1] = self.small[n0][n1], self.small[n0][n1+1]
                self.cache[n0] = writemd( self.big[n0], self.middle[n0], self.small[n0] )
        self.render(True, False)
    
    def custom1(self, x):
        if x == 2 and self.editable[0]:
            num = self.current[0]
            self.cache[num] = writemd( self.big[num], self.middle[num], self.small[num] )

kobj.repath()
if not os.path.exists("./_ST5_DATA/"):
    os.mkdir("./_ST5_DATA/")
    with open("./_ST5_DATA/example.txt", "w", encoding="utf-8") as f:
        f.write("<My Note>\n  # example text\n\n{First Record} (only text file)\n{Second Record} (we\ncan\nread\nit)  # multi-line\n")
worker = mainclass()
time.sleep(0.5)
