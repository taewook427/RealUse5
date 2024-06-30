# test667 : stdlib5.kgui blocksel

import time

import tkinter as tk

class toolbox:
    def __init__(self, title, iswin):
        self.title = title
        self.parms = {"color.0" : "light sky blue", "color.1" : "midnight blue", "color.2" : "azure"}
        if iswin:
            # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "600x460+100+50", ("맑은 고딕", 14), ('Consolas', 13)
            self.parms["ulbl.x"], self.parms["ulbl.y"], self.parms["ubut.x"], self.parms["ubut.y"] = 10, 10, 550, 10
            self.parms["lbut.x"], self.parms["lbut.y"], self.parms["lbut.txt"] = 10, 60, "\n" * 7 + " < " + "\n" * 7
            self.parms["rbut.x"], self.parms["rbut.y"], self.parms["rbut.txt"] = 550, 60, "\n" * 7 + " > " + "\n" * 7
            self.parms["array.x"], self.parms["array.y"], self.parms["array.lbl_low"] = 60, 60, 160
            self.parms["array.move_x"], self.parms["array.move_y"] = 165, 200
        else:
            # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "800x540+150+75", ("맑은 고딕", 14), ('Consolas', 11)
            self.parms["ulbl.x"], self.parms["ulbl.y"], self.parms["ubut.x"], self.parms["ubut.y"] = 10, 10, 710, 10
            self.parms["lbut.x"], self.parms["lbut.y"], self.parms["lbut.txt"] = 10, 100, "\n" * 3 + " < " + "\n" * 3
            self.parms["rbut.x"], self.parms["rbut.y"], self.parms["rbut.txt"] = 710, 100, "\n" * 3 + " > " + "\n" * 3
            self.parms["array.x"], self.parms["array.y"], self.parms["array.lbl_low"] = 110, 100, 160
            self.parms["array.move_x"], self.parms["array.move_y"] = 200, 220

        # user config range
        self.pics = [ ] # selection picture path (150x150), str[]
        self.txts = [ ] # selection text (eng len 12), str[]
        self.umsg = [ ] # upper msg txt, str[]
        
        # gui memory
        self.mwin = None # main window, obj
        self.upos = -1 # upper button position, int
        self.curpos = 0 # current page position, int
        self.comp = dict() # other gui components

    # init window, start working
    def entry(self):
        def ubut_click(): # work when ubut clicked
            time.sleep(0.1)
            self.upos = (self.upos + 1) % len(self.umsg)
            self.render(True)

        def lbut_click(): # work when lbut clicked
            time.sleep(0.1)
            if self.curpos > 1:
                self.curpos = self.curpos - 1
                self.render(False)

        def rbut_click(): # work when rbut clicked
            time.sleep(0.1)
            if self.curpos < (len(self.txts) - 1) // 6 + 1:
                self.curpos = self.curpos + 1
                self.render(False)

        def but0_click(): # work when but0 clicked
            time.sleep(0.1)
            self.custom0(6 * (self.curpos - 1) + 0)

        def but1_click(): # work when but1 clicked
            time.sleep(0.1)
            self.custom0(6 * (self.curpos - 1) + 1)

        def but2_click(): # work when but2 clicked
            time.sleep(0.1)
            self.custom0(6 * (self.curpos - 1) + 2)

        def but3_click(): # work when but3 clicked
            time.sleep(0.1)
            self.custom0(6 * (self.curpos - 1) + 3)

        def but4_click(): # work when but4 clicked
            time.sleep(0.1)
            self.custom0(6 * (self.curpos - 1) + 4)

        def but5_click(): # work when but5 clicked
            time.sleep(0.1)
            self.custom0(6 * (self.curpos - 1) + 5)

        if len(self.pics) != len(self.txts) or len(self.txts) % 6 != 0:
            raise Exception("invalid GUI orders")
        
        self.mwin = tk.Tk()
        self.mwin.title(self.title)
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color.0"] )

        self.comp["ulbl_strvar"] = tk.StringVar() # upper msg label
        self.comp["ulbl_strvar"].set("None")
        self.comp["ulbl"] = tk.Label( self.mwin, font=self.parms["font.0"], textvariable=self.comp["ulbl_strvar"],
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["ulbl"].place( x=self.parms["ulbl.x"], y=self.parms["ulbl.y"] )
        self.comp["ubut"] = tk.Button( self.mwin, font=self.parms["font.0"], text=" ◉ ", command=ubut_click,
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["ubut"].place( x=self.parms["ubut.x"], y=self.parms["ubut.y"] ) # change ubut state

        self.comp["lbut"] = tk.Button( self.mwin, font=self.parms["font.0"], text=self.parms["lbut.txt"], command=lbut_click,
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["lbut"].place( x=self.parms["lbut.x"], y=self.parms["lbut.y"] ) # go to left stage
        self.comp["rbut"] = tk.Button( self.mwin, font=self.parms["font.0"], text=self.parms["rbut.txt"], command=rbut_click,
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["rbut"].place( x=self.parms["rbut.x"], y=self.parms["rbut.y"] ) # go to right stage

        # 6 button data
        self.comp["lbl0_strvar"], self.comp["photo0"] = tk.StringVar(), None
        self.comp["lbl0_strvar"].set("None")
        self.comp["lbl1_strvar"], self.comp["photo1"] = tk.StringVar(), None
        self.comp["lbl1_strvar"].set("None")
        self.comp["lbl2_strvar"], self.comp["photo2"] = tk.StringVar(), None
        self.comp["lbl2_strvar"].set("None")
        self.comp["lbl3_strvar"], self.comp["photo3"] = tk.StringVar(), None
        self.comp["lbl3_strvar"].set("None")
        self.comp["lbl4_strvar"], self.comp["photo4"] = tk.StringVar(), None
        self.comp["lbl4_strvar"].set("None")
        self.comp["lbl5_strvar"], self.comp["photo5"] = tk.StringVar(), None
        self.comp["lbl5_strvar"].set("None")

        # 6 button array
        self.comp["but0"] = tk.Button( self.mwin, image=self.comp["photo0"], borderwidth=0, command=but0_click,
            bg=self.parms["color.2"] )
        self.comp["but1"] = tk.Button( self.mwin, image=self.comp["photo1"], borderwidth=0, command=but1_click,
            bg=self.parms["color.2"] )
        self.comp["but2"] = tk.Button( self.mwin, image=self.comp["photo2"], borderwidth=0, command=but2_click,
            bg=self.parms["color.2"] )
        self.comp["but3"] = tk.Button( self.mwin, image=self.comp["photo3"], borderwidth=0, command=but3_click,
            bg=self.parms["color.2"] )
        self.comp["but4"] = tk.Button( self.mwin, image=self.comp["photo4"], borderwidth=0, command=but4_click,
            bg=self.parms["color.2"] )
        self.comp["but5"] = tk.Button( self.mwin, image=self.comp["photo5"], borderwidth=0, command=but5_click,
            bg=self.parms["color.2"] )
        
        self.comp["lbl0"] = tk.Label( self.mwin, font=self.parms["font.1"], textvariable=self.comp["lbl0_strvar"],
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["lbl1"] = tk.Label( self.mwin, font=self.parms["font.1"], textvariable=self.comp["lbl1_strvar"],
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["lbl2"] = tk.Label( self.mwin, font=self.parms["font.1"], textvariable=self.comp["lbl2_strvar"],
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["lbl3"] = tk.Label( self.mwin, font=self.parms["font.1"], textvariable=self.comp["lbl3_strvar"],
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["lbl4"] = tk.Label( self.mwin, font=self.parms["font.1"], textvariable=self.comp["lbl4_strvar"],
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["lbl5"] = tk.Label( self.mwin, font=self.parms["font.1"], textvariable=self.comp["lbl5_strvar"],
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        
        x, y, low, movx, movy = self.parms["array.x"], self.parms["array.y"], self.parms["array.lbl_low"], self.parms["array.move_x"], self.parms["array.move_y"]
        self.comp["but0"].place(x=x, y=y)
        self.comp["lbl0"].place(x=x, y=y + low)
        self.comp["but1"].place(x=x + movx, y=y)
        self.comp["lbl1"].place(x=x + movx, y=y + low)
        self.comp["but2"].place(x=x + 2 * movx, y=y)
        self.comp["lbl2"].place(x=x + 2 * movx, y=y + low)
        self.comp["but3"].place(x=x, y=y + movy)
        self.comp["lbl3"].place(x=x, y=y + low + movy)
        self.comp["but4"].place(x=x + movx, y=y + movy)
        self.comp["lbl4"].place(x=x + movx, y=y + low + movy)
        self.comp["but5"].place(x=x + 2 * movx, y=y + movy)
        self.comp["lbl5"].place(x=x + 2 * movx, y=y + low + movy)

    # fill pic, txt by self.curpos, self.upos / set uonly True to update umsg only
    def render(self, uonly):
        self.comp["ulbl_strvar"].set(f"[ {self.curpos} / {(len(self.txts) - 1) // 6 + 1} ]   {self.umsg[self.upos]}")
        if not uonly:
            num = 6 * (self.curpos - 1)
            self.comp["photo0"] = tk.PhotoImage( file=self.pics[num + 0] )
            self.comp["photo1"] = tk.PhotoImage( file=self.pics[num + 1] )
            self.comp["photo2"] = tk.PhotoImage( file=self.pics[num + 2] )
            self.comp["photo3"] = tk.PhotoImage( file=self.pics[num + 3] )
            self.comp["photo4"] = tk.PhotoImage( file=self.pics[num + 4] )
            self.comp["photo5"] = tk.PhotoImage( file=self.pics[num + 5] )
            self.comp["but0"].config( image=self.comp["photo0"] )
            self.comp["but1"].config( image=self.comp["photo1"] )
            self.comp["but2"].config( image=self.comp["photo2"] )
            self.comp["but3"].config( image=self.comp["photo3"] )
            self.comp["but4"].config( image=self.comp["photo4"] )
            self.comp["but5"].config( image=self.comp["photo5"] )
            self.comp["lbl0_strvar"].set( self.txts[num + 0] )
            self.comp["lbl1_strvar"].set( self.txts[num + 1] )
            self.comp["lbl2_strvar"].set( self.txts[num + 2] )
            self.comp["lbl3_strvar"].set( self.txts[num + 3] )
            self.comp["lbl4_strvar"].set( self.txts[num + 4] )
            self.comp["lbl5_strvar"].set( self.txts[num + 5] )
        self.mwin.update()

    # go inside of GUI mainloop
    def guiloop(self):
        self.render(False)
        self.mwin.mainloop()
        time.sleep(0.5)

    # user custom func, options[x]
    def custom0(self, x):
        pass
