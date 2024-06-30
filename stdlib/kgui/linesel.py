# test666 : stdlib5.kgui linesel

import time

import tkinter as tk

class toolbox:
    def __init__(self, title, iswin):
        self.title = title
        self.parms = {"color.0" : "plum1", "color.1" : "deep pink"}
        if iswin:
            # windows
            self.parms["mwin.size"] = "320x400+100+50"
            self.parms["frame.x"], self.parms["frame.y"] = 10, 10
            self.parms["lstbox.w"], self.parms["lstbox.h"], self.parms["lstbox.font"] = 25, 14, ('Consolas', 15)
            self.parms["lbut.x"], self.parms["lbut.y"], self.parms["rbut.x"], self.parms["rbut.y"] = 10, 355, 280, 355
            self.parms["clbl.x"], self.parms["clbl.y"], self.parms["clbl.font"] = 40, 355, ('Consolas', 14)
        else:
            # linux
            self.parms["mwin.size"] = "480x600+150+75"
            self.parms["frame.x"], self.parms["frame.y"] = 10, 10
            self.parms["lstbox.w"], self.parms["lstbox.h"], self.parms["lstbox.font"] = 22, 8, ('Consolas', 15)
            self.parms["lbut.x"], self.parms["lbut.y"], self.parms["rbut.x"], self.parms["rbut.y"] = 10, 525, 425, 525
            self.parms["clbl.x"], self.parms["clbl.y"], self.parms["clbl.font"] = 70, 525, ('Consolas', 14)

        # user config range
        self.infos = [ ] # window infos, str[]
        self.options = [ ] # select options, str[][]
        
        # gui memory
        self.mwin = None # main window, obj
        self.curpos = -1 # current infos position, int
        self.comp = dict() # other gui components

    # init window, start working
    def entry(self):
        def lbut_click(): # work when lbut clicked
            time.sleep(0.1)
            if self.curpos > 0:
                self.curpos = self.curpos - 1
                self.render()

        def rbut_click(): # work when rbut clicked
            time.sleep(0.1)
            if len(self.infos) > self.curpos + 1:
                self.curpos = self.curpos + 1
                self.render()

        def lstbox_click(event): # work when lstbox clicked
            time.sleep(0.1)
            temp = self.comp["lstbox"].curselection()[0]
            if self.comp["lclick"] == temp:
                self.custom0(self.curpos, temp) # launch user func
            else:
                self.comp["lclick"] = temp

        if len(self.infos) != len(self.options):
            raise Exception("invalid GUI orders")

        self.mwin = tk.Tk()
        self.mwin.title(self.title)
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color.0"] )

        self.comp["frame"] = tk.Frame(self.mwin)
        self.comp["frame"].place( x=self.parms["frame.x"], y=self.parms["frame.y"] )
        self.comp["lstbox"] = tk.Listbox(
            self.comp["frame"], width=self.parms["lstbox.w"],  height=self.parms["lstbox.h"], font=self.parms["lstbox.font"],
            bg=self.parms["color.0"], fg=self.parms["color.1"], selectbackground=self.parms["color.1"], selectforeground=self.parms["color.0"] )
        self.comp["lstbox"].pack(side="left", fill="y") # options selection box
        self.comp["scroll"] = tk.Scrollbar(self.comp["frame"], orient="vertical")
        self.comp["scroll"].config(command=self.comp["lstbox"].yview)
        self.comp["scroll"].pack(side="right", fill="y")
        self.comp["lstbox"].config(yscrollcommand=self.comp["scroll"].set)
        self.comp["lstbox"].bind('<ButtonRelease-1>', lstbox_click)
        self.comp["lclick"] = -1 # pos of last clicked

        self.comp["lbut"] = tk.Button( self.mwin, font=self.parms["clbl.font"], text="<", command=lbut_click,
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["lbut"].place( x=self.parms["lbut.x"], y=self.parms["lbut.y"] ) # go to left stage
        self.comp["rbut"] = tk.Button( self.mwin, font=self.parms["clbl.font"], text=">", command=rbut_click,
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["rbut"].place( x=self.parms["rbut.x"], y=self.parms["rbut.y"] ) # go to right stage
        self.comp["clbl_strvar"] = tk.StringVar() # current stage
        self.comp["clbl_strvar"].set("None")
        self.comp["clbl"] = tk.Label( self.mwin, font=self.parms["clbl.font"], textvariable=self.comp["clbl_strvar"],
            bg=self.parms["color.0"], fg=self.parms["color.1"] )
        self.comp["clbl"].place( x=self.parms["clbl.x"], y=self.parms["clbl.y"] )

    # fill lstbox, center entry by self.curpos
    def render(self):
        self.comp["lclick"] = -1
        self.comp["lstbox"].delete( 0, self.comp["lstbox"].size() )
        if self.curpos >= 0:
            self.comp["clbl_strvar"].set( self.infos[self.curpos] )
            self.mwin.update()
            time.sleep(0.5)
            for i in self.options[self.curpos]:
                self.comp["lstbox"].insert(self.comp["lstbox"].size(), i)
                self.mwin.update()
                time.sleep(0.1)

    # go inside of GUI mainloop
    def guiloop(self):
        self.render()
        self.mwin.mainloop()
        time.sleep(0.5)
    
    # user custom func, options[x][y]
    def custom0(self, x, y):
        pass
