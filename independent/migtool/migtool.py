# test710 : independent.migtool

import time

import tkinter
import tkinter.messagebox
import tkinter.filedialog
import reader

import kobj

class mainclass(reader.toolbox):
    def __init__(self, iswin):
        super().__init__("MiGtool", iswin, 2)
        self.menus = ["    "] * 12

    def custom1(self, x):
        if x == 2:
            global tgt
            global data
            with open(tgt, "w", encoding="utf-8") as f:
                temp = f"###{self.current[0]},{self.current[1]}\r\n" + data if "\r" in data else f"###{self.current[0]},{self.current[1]}\n" + data
                f.write(temp)

def readtxt(raw): # read manual txt, returns (pos, big, middle, small)
    if raw[0:3] == "###": # find last readpos
        idx = raw.find("\n")
        temp = raw[3:idx]
        raw = raw[idx+1:]
        try:
            idx = temp.find(",")
            pos = [ int( temp[:idx] ), int( temp[idx+1:] ) ]
        except:
            pos = [0, 0]
    else:
        pos = [0, 0]

    raw, big, middle, small = raw + "\n", [ ], [ ], [ ]
    while "<" in raw or "{" in raw:
        if raw[0] == "<": # unit title
            idx = raw.find(">")
            big.append( raw[1:idx] )
            middle.append( [ ] )
            small.append( [ ] )
            raw = raw[idx+1:]
        elif raw[0] == "{": # name, content
            idx = raw.find("}")
            middle[-1].append( raw[1:idx] )
            raw = raw[idx+1:]
            idx = raw.find(")")
            small[-1].append( raw[raw.find("(")+1:idx] )
            raw = raw[idx+1:]
        elif raw[0] == "#": # comment
            raw = raw[raw.find("\n")+1:]
        else: # whitespace
            raw = raw[1:]
    return pos, big, middle, small

iswin = True # windows / linux
args = kobj.repath()
try:
    tgt = tkinter.filedialog.askopenfile( title="메뉴얼 파일 선택", filetypes=( ('Manual Text', '*.txt'), ) ).name if len(args) == 1 else args[1]
    with open(tgt, "r", encoding="utf-8") as f:
        data = f.read()
    a, b, c, d = readtxt(data)
    if data[0:3] == "###":
        data = data[data.find("\n")+1:]
except Exception as e:
    a, b, c, d, data = [0, 0], [""], [ [""] ], [ [""] ], ""
    tkinter.messagebox.showerror(title="Invalid Manual", message=f" Error occurred while reading manual text. \n {e} ")

if data != "":
    worker = mainclass(iswin)
    worker.big, worker.middle, worker.small = b, c, d
    worker.entry()
    worker.current = a
    worker.guiloop()
time.sleep(0.5)
