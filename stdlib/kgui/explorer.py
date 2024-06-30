# test668 : stdlib5.kgui explorer

import io
import time

import tkinter as tk
import tkinter.ttk as tkt
import tkinter.messagebox as tkm

from PIL import Image, ImageTk

class toolbox:
    def __init__(self, title, iswin):
        self.title = title
        self.parms = {"color.0" : "light cyan", "color.1" : "ghost white", "color.2" : "gray95", "color.3" : "lawn green"}
        if iswin:
            # windows
            self.parms["mwin.size"] = "720x640+100+50"
            self.parms["font.0"], self.parms["font.1"], self.parms["font.2"], self.parms["font.3"] = ("맑은 고딕", 14), ("맑은 고딕", 18), ('Consolas', 14), ("맑은 고딕", 13)
            self.parms["ubut0.x"], self.parms["ubut0.y"], self.parms["ubut1.x"], self.parms["ubut1.y"] = 10, 10, 60, 10
            self.parms["uent0.x"], self.parms["uent0.y"], self.parms["uent0.w"] = 110, 15, 28
            self.parms["ubut2.x"], self.parms["ubut2.y"], self.parms["ubut3.x"], self.parms["ubut3.y"] = 505, 10, 560, 10
            self.parms["uent1.x"], self.parms["uent1.y"], self.parms["uent1.w"] = 615, 15, 7
            self.parms["nb.x"], self.parms["nb.y"], self.parms["nb.w"], self.parms["nb.h"] = 10, 80, 695, 520
            self.parms["array.x"], self.parms["array.y"], self.parms["array.movx"], self.parms["array.movy"] = 10, 10, 170, 170
            self.parms["array.lbl_h"], self.parms["array.lbl_w"] = 30, 17
            self.parms["sbut0.x"], self.parms["sbut0.y"], self.parms["sent0.x"], self.parms["sent0.y"], self.parms["sent0.w"] = 10, 10, 65, 15, 20
            self.parms["slst0.x"], self.parms["slst0.y"], self.parms["slst0.w"], self.parms["slst0.h"], self.parms["sadd"] = 10, 60, 30, 18, 340
            self.parms["marr.x0"], self.parms["marr.x1"], self.parms["marr.x2"], self.parms["marr.y"] = 20, 500, 600, 10
            self.parms["tbox.x"], self.parms["tbox.y"], self.parms["tbox.w"], self.parms["tbox.h"] = 20, 60, 64, 20
            self.parms["pimg.x"], self.parms["pimg.y"], self.parms["pimg.w"], self.parms["pimg.h"] = 15, 60, 660, 450
            self.parms["bbox.x"], self.parms["bbox.y"], self.parms["bbox.w"], self.parms["bbox.h"] = 20, 60, 64, 20
        else:
            # linux
            self.parms["mwin.size"] = "900x760+150+75"
            self.parms["font.0"], self.parms["font.1"], self.parms["font.2"], self.parms["font.3"] = ("맑은 고딕", 12), ("맑은 고딕", 14), ('Consolas', 10), ("맑은 고딕", 10)
            self.parms["ubut0.x"], self.parms["ubut0.y"], self.parms["ubut1.x"], self.parms["ubut1.y"] = 10, 10, 90, 10
            self.parms["uent0.x"], self.parms["uent0.y"], self.parms["uent0.w"] = 160, 10, 18
            self.parms["ubut2.x"], self.parms["ubut2.y"], self.parms["ubut3.x"], self.parms["ubut3.y"] = 580, 10, 660, 10
            self.parms["uent1.x"], self.parms["uent1.y"], self.parms["uent1.w"] = 740, 10, 7
            self.parms["nb.x"], self.parms["nb.y"], self.parms["nb.w"], self.parms["nb.h"] = 10, 80, 880, 620
            self.parms["array.x"], self.parms["array.y"], self.parms["array.movx"], self.parms["array.movy"] = 10, 10, 220, 200
            self.parms["array.lbl_h"], self.parms["array.lbl_w"] = 45, 12
            self.parms["sbut0.x"], self.parms["sbut0.y"], self.parms["sent0.x"], self.parms["sent0.y"], self.parms["sent0.w"] = 10, 10, 90, 10, 15
            self.parms["slst0.x"], self.parms["slst0.y"], self.parms["slst0.w"], self.parms["slst0.h"], self.parms["sadd"] = 10, 80, 30, 13, 440
            self.parms["marr.x0"], self.parms["marr.x1"], self.parms["marr.x2"], self.parms["marr.y"] = 15, 590, 740, 10
            self.parms["tbox.x"], self.parms["tbox.y"], self.parms["tbox.w"], self.parms["tbox.h"] = 15, 90, 64, 13
            self.parms["pimg.x"], self.parms["pimg.y"], self.parms["pimg.w"], self.parms["pimg.h"] = 15, 90, 850, 520
            self.parms["bbox.x"], self.parms["bbox.y"], self.parms["bbox.w"], self.parms["bbox.h"] = 15, 90, 64, 13

        # user config range
        self.paths = [ ] # paths of tabs, str[5]
        self.names = [ ] # names of view tab, str[], ~/ if dir
        self.sizes = [ ] # sizes of view tab, int[], -1 if no info
        self.locked = [ ] # T:encdir / F:else, bool[]
        self.selected = [ ] # T:selected / F:not, bool[]
        self.search = [ ] # search result, str[]
        self.log0 = "" # current status, str
        self.log1 = [ ] # full log data, str[]
        self.txtdata = "" # txt edit data, str
        self.picdata = b"" # pic view data, bytes
        self.bindata = b"" # bin edit data, bytes

        # gui memory
        self.icons = dict() # icon pics (none, file, dir, .ext), dict(bytes)[str]
        self.tabpos = -1 # current tab number, int 0~4
        self.viewpos = 0 # view tab position, int 1+
        self.mwin = None # main window
        self.mbar = None # menu bar
        self.comp = dict() # other gui components

    # init window, start working
    def entry(self):
        self.resize()
        def exit_click(): # close button clicked
            time.sleep(0.1)
            if tkm.askyesno("Closing GUI", " Do you really want to close the screen? "):
                self.mwin.destroy()
                self.mwin = None

        def ubut0_click(): # upper "ok" button
            time.sleep(0.1)
            self.custom1(0)

        def ubut1_click(): # upper "back" button
            time.sleep(0.1)
            self.custom1(1)

        def ubut2_click(): # upper "go up" button
            time.sleep(0.1)
            if self.viewpos > 1:
                self.viewpos = self.viewpos - 1
            self.render(True, True, True, False, False, False, False)

        def ubut3_click(): # upper "go down" button
            time.sleep(0.1)
            if self.viewpos < (len(self.names) - 1) // 12 + 1:
                self.viewpos = self.viewpos + 1
            self.render(True, True, True, False, False, False, False)

        def sbut0_click(): # search button 0
           time.sleep(0.1)
           self.custom5( 0, self.comp["sent0"].get() )

        def sbut1_click(): # search button 1
           time.sleep(0.1)
           self.custom5( 1, self.comp["sent0"].get() )

        def tbut0_click(): # txt button 0
           time.sleep(0.1)
           self.render(False, False, False, False, True, False, False)

        def tbut1_click(): # txt button 1
           time.sleep(0.1)
           self.custom6( self.comp["tbox"].get('1.0', tk.END)[0:-1] )

        def pbut0_click(): # pic button 0
           time.sleep(0.1)
           self.custom7(0)

        def pbut1_click(): # pic button 1
           time.sleep(0.1)
           self.custom7(1)

        def bbut0_click(): # bin button 0
            time.sleep(0.1)
            try:
                worker = bitgui(b"")
                worker.TtoB( self.comp["bbox"].get('1.0', tk.END)[0:-1] )
                if self.parms["bbox.w"] > 63:
                    temp = worker.BtoT(16)
                else:
                    temp = worker.BtoT(8)
                self.comp["bbox"].delete('1.0', tk.END)
                self.comp["bbox"].insert('1.0', temp)
                self.mwin.update()
            except Exception as e:
                tkm.showerror("binary error", f" {e} ")

        def bbut1_click(): # bin button 1
            time.sleep(0.1)
            self.render(False, False, False, False, False, False, True)

        def bbut2_click(): # bin button 2
            time.sleep(0.1)
            try:
                worker = bitgui(b"")
                worker.TtoB( self.comp["bbox"].get('1.0', tk.END)[0:-1] )
                self.custom8(worker.data)
            except Exception as e:
                tkm.showerror("binary error", f" {e} ")

        def vbut_click(num, tp): # general view button click
            time.sleep(0.1)
            num = 12 * (self.viewpos - 1) + num
            if num < len(self.names):
                if tp == 0:
                    self.custom2(num)
                elif tp == 1:
                    self.custom3(num)
                else:
                    self.custom4(num)

        def vf0a(): vbut_click(0, 0)
        def vf1a(): vbut_click(1, 0)
        def vf2a(): vbut_click(2, 0)
        def vf3a(): vbut_click(3, 0)
        def vf4a(): vbut_click(4, 0)
        def vf5a(): vbut_click(5, 0)
        def vf6a(): vbut_click(6, 0)
        def vf7a(): vbut_click(7, 0)
        def vf8a(): vbut_click(8, 0)
        def vf9a(): vbut_click(9, 0)
        def vf10a(): vbut_click(10, 0)
        def vf11a(): vbut_click(11, 0)
        def vf0b(): vbut_click(0, 1)
        def vf1b(): vbut_click(1, 1)
        def vf2b(): vbut_click(2, 1)
        def vf3b(): vbut_click(3, 1)
        def vf4b(): vbut_click(4, 1)
        def vf5b(): vbut_click(5, 1)
        def vf6b(): vbut_click(6, 1)
        def vf7b(): vbut_click(7, 1)
        def vf8b(): vbut_click(8, 1)
        def vf9b(): vbut_click(9, 1)
        def vf10b(): vbut_click(10, 1)
        def vf11b(): vbut_click(11, 1)
        def vf0c(): vbut_click(0, 2)
        def vf1c(): vbut_click(1, 2)
        def vf2c(): vbut_click(2, 2)
        def vf3c(): vbut_click(3, 2)
        def vf4c(): vbut_click(4, 2)
        def vf5c(): vbut_click(5, 2)
        def vf6c(): vbut_click(6, 2)
        def vf7c(): vbut_click(7, 2)
        def vf8c(): vbut_click(8, 2)
        def vf9c(): vbut_click(9, 2)
        def vf10c(): vbut_click(10, 2)
        def vf11c(): vbut_click(11, 2)

        def tab_click(event): # notebook tab click
            time.sleep(0.1)
            self.tabpos = self.comp["nb"].index("current")
            self.render(True, False, False, False, False, False, False)

        self.mwin = tk.Tk() # main window
        self.mwin.title(self.title)
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color.0"] )
        self.mwin.protocol('WM_DELETE_WINDOW', exit_click)
        self.mbar = tk.Menu(self.mwin) # menu bar
        self.menubuilder()
        self.mwin.config(menu=self.mbar)

        self.comp["ubut0"] = tk.Button(self.mwin, font=self.parms["font.0"], text=" ◉ ", command=ubut0_click) # upper button 0 (ok)
        self.comp["ubut0"].place( x=self.parms["ubut0.x"], y=self.parms["ubut0.y"] )
        self.comp["ubut1"] = tk.Button(self.mwin, font=self.parms["font.0"], text=" < ", command=ubut1_click) # upper button 1 (back)
        self.comp["ubut1"].place( x=self.parms["ubut1.x"], y=self.parms["ubut1.y"] )
        self.comp["uent0_strvar"] = tk.StringVar() # upper entry 0 (path)
        self.comp["uent0_strvar"].set("None")
        self.comp["uent0"] = tk.Entry(self.mwin, width=self.parms["uent0.w"], font=self.parms["font.1"], textvariable=self.comp["uent0_strvar"], state="readonly")
        self.comp["uent0"].place( x=self.parms["uent0.x"], y=self.parms["uent0.y"] )
        self.comp["ubut2"] = tk.Button(self.mwin, font=self.parms["font.0"], text=" ↑ ", command=ubut2_click) # upper button 2 (go up)
        self.comp["ubut2"].place( x=self.parms["ubut2.x"], y=self.parms["ubut2.y"] )
        self.comp["ubut3"] = tk.Button(self.mwin, font=self.parms["font.0"], text=" ↓ ", command=ubut3_click) # upper button 3 (go down)
        self.comp["ubut3"].place( x=self.parms["ubut3.x"], y=self.parms["ubut3.y"] )
        self.comp["uent1_strvar"] = tk.StringVar() # upper entry 1 (curpos)
        self.comp["uent1_strvar"].set("None")
        self.comp["uent1"] = tk.Entry(self.mwin, width=self.parms["uent1.w"], font=self.parms["font.1"], textvariable=self.comp["uent1_strvar"], state="readonly")
        self.comp["uent1"].place( x=self.parms["uent1.x"], y=self.parms["uent1.y"] )

        self.comp["nb"] = tkt.Notebook( self.mwin, width=self.parms["nb.w"], height=self.parms["nb.h"] ) # gui notebook
        self.comp["nb"].place( x=self.parms["nb.x"], y=self.parms["nb.y"] )
        self.comp["nb"].bind('<ButtonRelease-1>', tab_click)
        self.comp["nbfr0"] = tk.Frame( self.mwin, bg=self.parms["color.0"] )
        self.comp["nbfr1"] = tk.Frame( self.mwin, bg=self.parms["color.0"] )
        self.comp["nbfr2"] = tk.Frame( self.mwin, bg=self.parms["color.0"] )
        self.comp["nbfr3"] = tk.Frame( self.mwin, bg=self.parms["color.0"] )
        self.comp["nbfr4"] = tk.Frame( self.mwin, bg=self.parms["color.0"] )
        self.comp["nb"].add(self.comp["nbfr0"], text="    view    ") # notebook frames
        self.comp["nb"].add(self.comp["nbfr1"], text="   search   ")
        self.comp["nb"].add(self.comp["nbfr2"], text="  txt edit  ")
        self.comp["nb"].add(self.comp["nbfr3"], text="  pic view  ")
        self.comp["nb"].add(self.comp["nbfr4"], text="  bin edit  ")

        # 12 button data
        for i in range(0, 12):
            self.comp[f"vent{i}a_strvar"], self.comp[f"vent{i}b_strvar"], self.comp[f"vbut{i}_photo"] = tk.StringVar(), tk.StringVar(), ImageTk.PhotoImage( self.icons["dir"] )
            self.comp[f"vent{i}a_strvar"].set("None")
            self.comp[f"vent{i}b_strvar"].set("None")

        # 12 button array
        orders = [vf0a, vf1a, vf2a, vf3a, vf4a, vf5a, vf6a, vf7a, vf8a, vf9a, vf10a, vf11a,
                vf0b, vf1b, vf2b, vf3b, vf4b, vf5b, vf6b, vf7b, vf8b, vf9b, vf10b, vf11b,
                vf0c, vf1c, vf2c, vf3c, vf4c, vf5c, vf6c, vf7c, vf8c, vf9c, vf10c, vf11c]
        win = self.comp["nbfr0"]
        x, y, mx, my, h, w = self.parms["array.x"], self.parms["array.y"], self.parms["array.movx"], self.parms["array.movy"], self.parms["array.lbl_h"], self.parms["array.lbl_w"]
        for i in range(0, 3):
            for j in range(0, 4):
                num = 4 * i + j
                self.comp[f"vbut{num}a"] = tk.Button(win, image=self.comp[f"vbut{num}_photo"], borderwidth=0, command=orders[num],
                    bg=self.parms["color.1"] ) # icon button
                self.comp[f"vbut{num}a"].place(x=j*mx + x, y=i*my + y)
                self.comp[f"vbut{num}b"] = tk.Button(win, font=self.parms["font.2"], text="LOCK", command=orders[num + 12],
                    bg=self.parms["color.3"] ) # lock sign
                self.comp[f"vbut{num}b"].place(x=j*mx + x+105, y=i*my + y)
                self.comp[f"vbut{num}c"] = tk.Button(win, font=self.parms["font.2"], text=" SEL", command=orders[num + 24],
                    bg=self.parms["color.3"] ) # checked sign
                self.comp[f"vbut{num}c"].place(x=j*mx + x+105, y=i*my + y+50)
                self.comp[f"vent{num}a"] = tk.Entry(win, font=self.parms["font.3"], textvariable=self.comp[f"vent{num}a_strvar"], width=w, state="readonly")
                self.comp[f"vent{num}a"].place(x=j*mx + x, y=i*my + y+105) # name
                self.comp[f"vent{num}b"] = tk.Entry(win, font=self.parms["font.3"], textvariable=self.comp[f"vent{num}b_strvar"], width=w, state="readonly")
                self.comp[f"vent{num}b"].place(x=j*mx + x, y=i*my + y+105+h) # size

        self.comp["sbut0"] = tk.Button(self.comp["nbfr1"], font=self.parms["font.0"], text=" ⌕ ", command=sbut0_click)
        self.comp["sbut0"].place( x=self.parms["sbut0.x"], y=self.parms["sbut0.y"] ) # search button 0
        self.comp["sent0"] = tk.Entry( self.comp["nbfr1"], font=self.parms["font.1"], width=self.parms["sent0.w"] )
        self.comp["sent0"].place( x=self.parms["sent0.x"], y=self.parms["sent0.y"] ) # search entry 0
        self.comp["sfr0"] = tk.Frame( self.comp["nbfr1"] ) # search listbox frame 0
        self.comp["sfr0"].place( x=self.parms["slst0.x"], y=self.parms["slst0.y"] )
        self.comp["slst0"] = tk.Listbox( self.comp["sfr0"], width=self.parms["slst0.w"], height=self.parms["slst0.h"], font=self.parms["font.2"] )
        self.comp["slst0"].pack(side="left", fill="y") # search listbox 0
        self.comp["slst0_scroll"] = tk.Scrollbar(self.comp["sfr0"], orient="vertical")
        self.comp["slst0_scroll"].config(command=self.comp["slst0"].yview)
        self.comp["slst0_scroll"].pack(side="right", fill="y")
        self.comp["slst0"].config(yscrollcommand=self.comp["slst0_scroll"].set)

        self.comp["sbut1"] = tk.Button(self.comp["nbfr1"], font=self.parms["font.0"], text=" ↻ ", command=sbut1_click,
            bg=self.parms["color.3"] )
        self.comp["sbut1"].place( x=self.parms["sbut0.x"]+self.parms["sadd"], y=self.parms["sbut0.y"] ) # search button 1
        self.comp["sent1_strvar"] = tk.StringVar() # search entry 1
        self.comp["sent1_strvar"].set("None")
        self.comp["sent1"] = tk.Entry(self.comp["nbfr1"], font=self.parms["font.1"], textvariable=self.comp["sent1_strvar"], width=self.parms["sent0.w"], state="readonly")
        self.comp["sent1"].place( x=self.parms["sent0.x"]+self.parms["sadd"], y=self.parms["sent0.y"] )
        self.comp["sfr1"] = tk.Frame( self.comp["nbfr1"] ) # search listbox frame 1
        self.comp["sfr1"].place( x=self.parms["slst0.x"]+self.parms["sadd"], y=self.parms["slst0.y"] )
        self.comp["slst1"] = tk.Listbox( self.comp["sfr1"], width=self.parms["slst0.w"], height=self.parms["slst0.h"], font=self.parms["font.2"] )
        self.comp["slst1"].pack(side="left", fill="y") # search listbox 1
        self.comp["slst1_scroll"] = tk.Scrollbar(self.comp["sfr1"], orient="vertical")
        self.comp["slst1_scroll"].config(command=self.comp["slst1"].yview)
        self.comp["slst1_scroll"].pack(side="right", fill="y")
        self.comp["slst1"].config(yscrollcommand=self.comp["slst1_scroll"].set)

        self.comp["tbut0"] = tk.Button(self.comp["nbfr2"], font=self.parms["font.0"], text=" reload ", command=tbut0_click)
        self.comp["tbut0"].place( x=self.parms["marr.x1"], y=self.parms["marr.y"] ) # txt button 0
        self.comp["tbut1"] = tk.Button(self.comp["nbfr2"], font=self.parms["font.0"], text="  save  ", command=tbut1_click)
        self.comp["tbut1"].place( x=self.parms["marr.x2"], y=self.parms["marr.y"] ) # txt button 1
        self.comp["tfr"] = tk.Frame( self.comp["nbfr2"] ) # txt textbox frame
        self.comp["tfr"].place( x=self.parms["tbox.x"], y=self.parms["tbox.y"] )
        self.comp["tbox"] = tk.Text( self.comp["tfr"], width=self.parms["tbox.w"], height=self.parms["tbox.h"], font=self.parms["font.2"] )
        self.comp["tbox"].pack(side="left", fill="y") # txt textbox
        self.comp["tbox_scroll"] = tk.Scrollbar(self.comp["tfr"], orient="vertical")
        self.comp["tbox_scroll"].config(command=self.comp["tbox"].yview)
        self.comp["tbox_scroll"].pack(side="right", fill="y")
        self.comp["tbox"].config(yscrollcommand=self.comp["tbox_scroll"].set)

        self.comp["pbut0"] = tk.Button(self.comp["nbfr3"], font=self.parms["font.0"], text="   <-   ", command=pbut0_click)
        self.comp["pbut0"].place( x=self.parms["marr.x1"], y=self.parms["marr.y"] ) # pic button 0
        self.comp["pbut1"] = tk.Button(self.comp["nbfr3"], font=self.parms["font.0"], text="   ->   ", command=pbut1_click)
        self.comp["pbut1"].place( x=self.parms["marr.x2"], y=self.parms["marr.y"] ) # pic button 1
        self.comp["pimg"] = tk.Canvas( self.comp["nbfr3"], width=self.parms["pimg.w"], height=self.parms["pimg.h"], bg=self.parms["color.1"] )
        self.comp["pimg"].place( x=self.parms["pimg.x"], y=self.parms["pimg.y"] ) # pic image

        self.comp["bbut0"] = tk.Button(self.comp["nbfr4"], font=self.parms["font.0"], text="  sort  ", command=bbut0_click)
        self.comp["bbut0"].place( x=self.parms["marr.x0"], y=self.parms["marr.y"] ) # bin button 0
        self.comp["bbut1"] = tk.Button(self.comp["nbfr4"], font=self.parms["font.0"], text=" reload ", command=bbut1_click)
        self.comp["bbut1"].place( x=self.parms["marr.x1"], y=self.parms["marr.y"] ) # bin button 1
        self.comp["bbut2"] = tk.Button(self.comp["nbfr4"], font=self.parms["font.0"], text="  save  ", command=bbut2_click)
        self.comp["bbut2"].place( x=self.parms["marr.x2"], y=self.parms["marr.y"] ) # bin button 2
        self.comp["bfr"] = tk.Frame( self.comp["nbfr4"] ) # bin textbox frame
        self.comp["bfr"].place( x=self.parms["bbox.x"], y=self.parms["bbox.y"] )
        self.comp["bbox"] = tk.Text( self.comp["bfr"], width=self.parms["bbox.w"], height=self.parms["bbox.h"], font=self.parms["font.2"] )
        self.comp["bbox"].pack(side="left", fill="y") # bin textbox
        self.comp["bbox_scroll"] = tk.Scrollbar(self.comp["bfr"], orient="vertical")
        self.comp["bbox_scroll"].config(command=self.comp["bbox"].yview)
        self.comp["bbox_scroll"].pack(side="right", fill="y")
        self.comp["bbox"].config(yscrollcommand=self.comp["bbox_scroll"].set)

    # render gui, each bool parm determines render or not, view0:picture/txts view1:chked
    def render(self, upper, view0, view1, sch, txtd, picd, bind):
        self.check()
        if upper:
            self.comp["uent0_strvar"].set( self.paths[self.tabpos] )
            self.comp["uent1_strvar"].set(f"{self.viewpos} / {(len(self.names)-1)//12+1}")
        if view0:
            for i in range(0, 12):
                num = 12 * (self.viewpos - 1) + i
                if num < len(self.names):
                    self.comp[f"vent{i}a_strvar"].set( self.names[num] )
                    self.comp[f"vent{i}b_strvar"].set( self.sizeasume( self.sizes[num] ) )
                    self.comp[f"vbut{i}_photo"] = ImageTk.PhotoImage( self.picasume( self.names[num] ) )
                    self.comp[f"vbut{i}a"].configure( image=self.comp[f"vbut{i}_photo"] )
                else:
                    self.comp[f"vent{i}a_strvar"].set("")
                    self.comp[f"vent{i}b_strvar"].set("")
                    self.comp[f"vbut{i}_photo"] = ImageTk.PhotoImage( self.picasume("") )
                    self.comp[f"vbut{i}a"].configure( image=self.comp[f"vbut{i}_photo"] )
        if view1:
            for i in range(0, 12):
                num = 12 * (self.viewpos - 1) + i
                if num < len(self.names):
                    if self.locked[num]:
                        self.comp[f"vbut{i}b"].configure( bg=self.parms["color.3"] )
                    else:
                        self.comp[f"vbut{i}b"].configure( bg=self.parms["color.2"] )
                    if self.selected[num]:
                        self.comp[f"vbut{i}c"].configure( bg=self.parms["color.3"] )
                    else:
                        self.comp[f"vbut{i}c"].configure( bg=self.parms["color.2"] )
                else:
                    self.comp[f"vbut{i}b"].configure( bg=self.parms["color.2"] )
                    self.comp[f"vbut{i}c"].configure( bg=self.parms["color.2"] )
        if sch:
            self.comp["slst0"].delete( 0, self.comp["slst0"].size() )
            for i in self.search:
                self.comp["slst0"].insert(self.comp["slst0"].size(), i)
            self.comp["sbut1"].configure( bg=self.parms["color.2"] )
            self.comp["sent1_strvar"].set(self.log0)
            self.comp["slst1"].delete( 0, self.comp["slst1"].size() )
            for i in self.log1:
                self.comp["slst1"].insert(self.comp["slst1"].size(), i)
        if txtd:
            self.comp["tbox"].delete('1.0', tk.END)
            self.comp["tbox"].insert('1.0', self.txtdata)
        if picd:
            try:
                pt = Image.open( io.BytesIO(self.picdata) )
                if self.parms["pimg.w"] / self.parms["pimg.h"] > pt.size[0] / pt.size[1]:
                    px, py = int( pt.size[0] / pt.size[1] * self.parms["pimg.h"] ), self.parms["pimg.h"]
                    px = px - px % 2
                else:
                    px, py = self.parms["pimg.w"], int( self.parms["pimg.w"] / pt.size[0] * pt.size[1] )
                    py = py - py % 2
                pt = pt.resize( (px, py) )
                self.comp["pimg"].destroy()
                self.comp["pimg"] = tk.Canvas( self.comp["nbfr3"], width=self.parms["pimg.w"], height=self.parms["pimg.h"], bg=self.parms["color.1"] )
                self.comp["pimg"].place( x=self.parms["pimg.x"], y=self.parms["pimg.y"] )
                self.comp["pimg_image"] = ImageTk.PhotoImage(pt)
                self.comp["pimg"].create_image( px//2, py//2, image=self.comp["pimg_image"] )
            except:
                pass
        if bind:
            worker = bitgui(self.bindata)
            self.comp["bbox"].delete('1.0', tk.END)
            if self.parms["bbox.w"] > 63:
                self.comp["bbox"].insert( '1.0', worker.BtoT(16) )
            else:
                self.comp["bbox"].insert( '1.0', worker.BtoT(8) )
        self.mwin.update()

    # go inside of GUI mainloop
    def guiloop(self):
        self.render(True, True, True, True, True, True, True)
        while self.mwin != None:
            time.sleep(0.1)
            if self.custom9():
                try:
                    self.comp["sbut1"].configure( bg=self.parms["color.3"] )
                    self.mwin.update()
                except:
                    pass
            self.mwin.update()
        time.sleep(0.5)

    def check(self): # check whether render condition is valid
        cond = True
        if len(self.paths) != 5:
            cond = False
        if len(self.names) != len(self.sizes) or len(self.names) != len(self.selected) or len(self.names) != len(self.locked):
            cond = False
        if self.tabpos < 0 or self.tabpos > 4:
            cond = False
        if (self.viewpos < 1 or self.viewpos > (len(self.names) - 1) // 12 + 1) and len(self.names) != 0:
            cond = False
        if not cond:
            raise Exception("invalid GUI order")

    def resize(self): # resize icons image to 100x100
        for i in self.icons:
            self.icons[i] = Image.open( io.BytesIO( self.icons[i] ) ).resize( (100, 100) )

    def sizeasume(self, n): # gen txt from size num, NoInfo if -1
        if n < 0:
            return "     ? B  "
        elif n < 1024:
            return f"{n:>6} B  "
        elif n < 1048576:
            return f"{n/1024:6.1f} KiB"
        elif n < 1073741824:
            return f"{n/1048576:6.1f} MiB"
        else:
            return f"{n/1073741824:6.1f} GiB"
        
    def picasume(self, nm): # return pic bytes match with name
        if len(nm) == 0:
            return self.icons["none"]
        if nm[-1] == "/":
            return self.icons["dir"]
        if "." in nm:
            tmp = nm[nm.rfind("."):].lower()
            if tmp in self.icons:
                return self.icons[tmp]
            else:
                return self.icons["file"]
        else:
            return self.icons["file"]

    # user custom func, build menu bar
    def menubuilder(self):
        pass

    # user custom func, menu option[x]
    def custom0(self, x):
        pass

    # user custom func, upper button x:0/1
    def custom1(self, x):
        pass

    # user custom func, view button A (pic) x:self.names[x]
    def custom2(self, x):
        pass

    # user custom func, view button B (lock) x:self.names[x]
    def custom3(self, x):
        pass

    # user custom func, view button C (select) x:self.names[x]
    def custom4(self, x):
        pass

    # user custom func, search button x:0/1, y:entry input str
    def custom5(self, x, y):
        pass

    # user custom func, txt "save" button, x:box input str
    def custom6(self, x):
        pass

    # user custom func, pic button x:0/1
    def custom7(self, x):
        pass
    
    # user custom func, bin "save" button, x:box input bytes
    def custom8(self, x):
        pass

    # user custom func, returns whether !!! update is required !!! (->bool)
    def custom9(self):
        return True

class bitgui:
    def __init__(self, data):
        if len(data) < 2147483648:
            self.data = data # bytes under 2GiB
        else:
            self.data = b""
        self.table = "00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F "

    def BtoT(self, wnum): # gen gui text, wnum = 8 or 16
        buf = [ ]
        if wnum == 16:
            buf.append("HEX_EDIT  " + self.table + "__ __ ")
            buf.append(" " * 64)
            for i in range(0, len(self.data) // 16):
                temp = [ f"{hex(16*i)[2:]:0>8}  ".upper() ]
                for j in range(0, 16):
                    temp.append( f"{hex(self.data[16*i+j])[2:]:0>2} ".lower() )
                temp.append("      ")
                buf.append( "".join(temp) )
            if len(self.data) % 16 != 0:
                i = len(self.data) // 16
                temp = [ f"{hex(16*i)[2:]:0>8}  ".upper() ]
                for j in range(0, len(self.data) % 16):
                    temp.append( f"{hex(self.data[16*i+j])[2:]:0>2} ".lower() )
                buf.append( "".join(temp) )

        else:
            buf.append("HEX_EDIT  " + self.table[0:24] + "__ __ ")
            buf.append(" " * 40)
            for i in range(0, len(self.data) // 8):
                temp = [ f"{hex(8*i)[2:]:0>8}  ".upper() ]
                for j in range(0, 8):
                    temp.append( f"{hex(self.data[8*i+j])[2:]:0>2} ".lower() )
                temp.append("      ")
                buf.append( "".join(temp) )
            if len(self.data) % 8 != 0:
                i = len(self.data) // 8
                temp = [ f"{hex(8*i)[2:]:0>8}  ".upper() ]
                for j in range(0, len(self.data) % 8):
                    temp.append( f"{hex(self.data[8*i+j])[2:]:0>2} ".lower() )
                buf.append( "".join(temp) )
        return "\n".join(buf)

    def TtoB(self, text): # gen bytes data
        buf = [ ]
        text = filter( lambda x: len(x) > 10, text.split("\n")[2:] )
        temp = [x[8:] for x in text]
        for i in temp:
            buf.append( bytes.fromhex( i.replace(" ", "") ) )
        self.data = b"".join(buf)
