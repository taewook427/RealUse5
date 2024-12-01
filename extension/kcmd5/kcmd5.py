# test719 : extension.kcmd5

import os
import time

import tkinter.messagebox
import linesel

import kobj
import kdb

class mainclass(linesel.toolbox):
    def __init__(self):
        self.cmd, worker = [ ], kdb.toolbox()
        with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
            worker.read( f.read() )
        super().__init__("kcmd5", worker.get("dev.os")[3] == "windows")
        for i in os.listdir("./_ST5_DATA/"):
            names, orders = readopt("./_ST5_DATA/" + i)
            self.infos.append(i)
            self.options.append(names)
            self.cmd.append(orders)
        if len(self.infos) == 0:
            tkinter.messagebox.showerror(title="No Orders", message=" Check orders in ./_ST5_DATA/ and try again. ")
        self.curpos = 0

    def custom0(self, x, y):
        res = exe( self.cmd[x][y] )
        tkinter.messagebox.showinfo(title="cmd result", message=f" {res} ")

def readopt(path): # read kdb optionsD
    try:
        worker, names, orders = kdb.toolbox(), [ ], [ ]
        with open(path, "r", encoding="utf-8") as f:
            worker.read( f.read() )
        i = 0
        while f"{i}.name" in worker.name:
            names.append( worker.get(f"{i}.name")[3] )
            temp, j = [ ], 0
            while f"{i}.{j}" in worker.name:
                temp.append( worker.get(f"{i}.{j}")[3] )
                j = j + 1
            orders.append(temp)
            i = i + 1
        return names, orders
    except Exception as e:
        return [f"read error : {e}"], [ [ ] ]

def exe(orders): # execute order, returns exit num
    out = [ ]
    for i in orders:
        time.sleep(0.1)
        try:
            temp = os.system(i)
        except:
            temp = 1
        out.append( str(temp) )
    return " ".join(out)

kobj.repath()
if not os.path.exists("./_ST5_DATA/"):
    os.mkdir("./_ST5_DATA/")
worker = mainclass()
worker.entry()
worker.guiloop()
time.sleep(0.5)

# 항상 관리자 권한으로 실행 (윈도우)
# Resource hacker - Manifest - requestedExecutionLevel level
# "asInvoker" -> "requireAdministrator"
