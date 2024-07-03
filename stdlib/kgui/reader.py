# test673 : stdlib5.kgui reader

import io
import time

import tkinter as tk
import tkinter.ttk as tkt
import tkinter.messagebox as tkm

from PIL import Image, ImageTk

class toolbox:
    def __init__(self, title, iswin, mode):
        self.title = title
        self.mode = mode # 0 : pic, 1 : sch, 2 : renew
        self.parms = {"color.0" : "dark green", "color.1" : "forest green", "color.2" : "lime green", "color.3" : "spring green", "color.4" : "mint cream"}
        if iswin:
            # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"], self.parms["font.2"] = "650x420+100+50", ("맑은 고딕", 12), ("Consolas", 14), ("맑은 고딕", 13)
            self.parms["uent.w"], self.parms["lst.x"], self.parms["lst.y"], self.parms["lst.w"], self.parms["lst.h"] = 50, 10, 120, 30, 10
            self.parms["lbut.x"], self.parms["lbut.y"], self.parms["lent.x"], self.parms["lent.y"], self.parms["lent.w"] = 10, 370, 60, 380, 26
            self.parms["tbox.x"], self.parms["tbox.y"], self.parms["tbox.w"], self.parms["tbox.h"] = 320, 120, 34, 12
            self.parms["awin.pic.size"], self.parms["awin.pic.x"], self.parms["awin.pic.y"] = "420x420+200+100", 420, 420
            self.parms["awin.sch.size"], self.parms["awin.sch.tbox.w"], self.parms["awin.sch.tbox.h"] = "320x420+200+100", 35, 4
            self.parms["awin.sch.lst.x"], self.parms["awin.sch.lst.y"], self.parms["awin.sch.lst.w"], self.parms["awin.sch.lst.h"] = 10, 150, 30, 10
        else:
            # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"], self.parms["font.2"] = "800x550+150+75", ("맑은 고딕", 8), ("Consolas", 10), ("맑은 고딕", 10)
            self.parms["uent.w"], self.parms["lst.x"], self.parms["lst.y"], self.parms["lst.w"], self.parms["lst.h"] = 42, 10, 160, 24, 8
            self.parms["lbut.x"], self.parms["lbut.y"], self.parms["lent.x"], self.parms["lent.y"], self.parms["lent.w"] = 10, 490, 80, 500, 20
            self.parms["tbox.x"], self.parms["tbox.y"], self.parms["tbox.w"], self.parms["tbox.h"] = 400, 160, 25, 9
            self.parms["awin.pic.size"], self.parms["awin.pic.x"], self.parms["awin.pic.y"] = "550x550+300+150", 550, 550
            self.parms["awin.sch.size"], self.parms["awin.sch.tbox.w"], self.parms["awin.sch.tbox.h"] = "400x550+300+150", 26, 3
            self.parms["awin.sch.lst.x"], self.parms["awin.sch.lst.y"], self.parms["awin.sch.lst.w"], self.parms["awin.sch.lst.h"] = 10, 200, 24, 8

        # user config range
        self.menus = [ ] # menu buttons, str[12]
        self.big = [ ] # seperation big, str[], size N
        self.middle = [ ] # seperation middle, str[][], size N * ~
        self.small = [ ] # content small, str/bytes[][], size N * ~
        self.editable = [False, False, False] # is seperation editable, bool[3]

        # gui memory
        self.current = [-1, -1] # current view position
        self.mwin = None # main window
        self.awin = None # additional window
        self.comp = dict() # other gui components

    # init window, start working
    def entry(self):
        if len(self.menus) != 12 or len(self.big) < 1:
            raise Exception("invalid GUI order")
        
        # menu functions
        def mf0():
            time.sleep(0.1)
            self.custom0(0)
        def mf1():
            time.sleep(0.1)
            self.custom0(1)
        def mf2():
            time.sleep(0.1)
            self.custom0(2)
        def mf3():
            time.sleep(0.1)
            self.custom0(3)
        def mf4():
            time.sleep(0.1)
            self.custom0(4)
        def mf5():
            time.sleep(0.1)
            self.custom0(5)
        def mf6():
            time.sleep(0.1)
            self.custom0(6)
        def mf7():
            time.sleep(0.1)
            self.custom0(7)
        def mf8():
            time.sleep(0.1)
            self.custom0(8)
        def mf9():
            time.sleep(0.1)
            self.custom0(9)
        def mf10():
            time.sleep(0.1)
            self.custom0(10)
        def mf11():
            time.sleep(0.1)
            self.custom0(11)

        # button functions
        def bf0(): # go left
            time.sleep(0.1)
            if self.current[0] > 0:
                self.current[0] = self.current[0] - 1
                self.current[1] = 0
                self.custom1(0)
                self.render(True, True)

        def bf1(): # go right
            time.sleep(0.1)
            if self.current[0] + 1 < len(self.big):
                self.current[0] = self.current[0] + 1
                self.current[1] = 0
                self.custom1(1)
                self.render(True, True)

        def bf2(): # regen
            time.sleep(0.1)
            if self.editable[0]:
                self.big[ self.current[0] ] = self.comp["uent"].get()
            if self.editable[1] and len( self.middle[ self.current[0] ] ) != 0:
                self.middle[ self.current[0] ][ self.current[1] ] = self.comp["lent"].get()
            if self.editable[2] and len( self.middle[ self.current[0] ] ) != 0:
                self.small[ self.current[0] ][ self.current[1] ] = self.comp["tbox"].get('1.0', tk.END)[0:-1]
            self.custom1(2)
            self.render(True, True)

        def bf3(): # mode func
            time.sleep(0.1)
            self.build()
            self.custom1(3)
            self.render(True, True)

        def lst_click(event): # listbox clicked
            time.sleep(0.1)
            temp = self.comp["lst"].curselection()
            if len(temp) > 0:
                self.current[1] = temp[0]
                self.render(True, True)

        def exit_click(): # close button clicked
            time.sleep(0.1)
            if tkm.askyesno("Closing GUI", " Do you really want to close the screen? "):
                if self.awin != None:
                    self.awin.destroy()
                    self.awin = None
                self.mwin.destroy()
                self.mwin = None

        self.mwin = tk.Tk() # main window
        self.mwin.title(self.title)
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color.3"] )
        self.mwin.protocol('WM_DELETE_WINDOW', exit_click)

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color.1"] )
        self.comp["frame0"].pack(side="top", fill="x") # menu buttons frame
        temp = [mf0, mf1, mf2, mf3, mf4, mf5, mf6, mf7, mf8, mf9, mf10, mf11]
        for i in range(0, 2):
            for j in range(0, 6):
                self.comp[f"mbut{6*i+j}"] = tk.Button( self.comp["frame0"], font=self.parms["font.0"], text=self.menus[6*i+j], command=temp[6*i+j],
                    bg=self.parms["color.1"], fg=self.parms["color.3"] ) # menu button 0~11
                self.comp[f"mbut{6*i+j}"].grid(row=i, column=j)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color.2"] )
        self.comp["frame1"].pack(side="top", fill="x") # move buttons frame
        self.comp["fr_style"] = tkt.Style() # frame1 style
        self.comp["fr_style"].theme_use('clam')
        self.comp["fr_style"].configure( "TEntry", fieldbackground=self.parms["color.2"] )

        self.comp["ubut0"] = tk.Button(self.comp["frame1"], font=self.parms["font.1"], text=" ← ", command=bf0,
            bg=self.parms["color.2"], fg=self.parms["color.0"] )
        self.comp["ubut0"].pack(side="left", fill="y") # move button (go left)
        self.comp["ubut1"] = tk.Button(self.comp["frame1"], font=self.parms["font.1"], text=" → ", command=bf1,
            bg=self.parms["color.2"], fg=self.parms["color.0"] )
        self.comp["ubut1"].pack(side="left", fill="y") # move button (go right)
        self.comp["ubut2"] = tk.Button(self.comp["frame1"], font=self.parms["font.1"], text=" ↻ ", command=bf2,
            bg=self.parms["color.2"], fg=self.parms["color.0"] )
        self.comp["ubut2"].pack(side="left", fill="y") # move button (regen)
        self.comp["uent_strvar"] = tk.StringVar()
        self.comp["uent_strvar"].set("None")
        self.comp["uent"] = tkt.Entry(self.comp["frame1"], width=self.parms["uent.w"], font=self.parms["font.1"], textvariable=self.comp["uent_strvar"], state="readonly",
             foreground=self.parms["color.0"] )
        self.comp["uent"].pack(side="right", fill="both") # upper entry

        self.comp["lst_frame"] = tk.Frame(self.mwin) # listbox
        self.comp["lst_frame"].place(x=self.parms["lst.x"], y=self.parms["lst.y"])
        self.comp["lst"] = tk.Listbox(self.comp["lst_frame"], width=self.parms["lst.w"], height=self.parms["lst.h"], font=self.parms["font.2"], selectmode="single",
            bg=self.parms["color.3"], fg=self.parms["color.0"], selectbackground=self.parms["color.0"], selectforeground=self.parms["color.3"] )
        self.comp["lst"].pack(side="left", fill="y")
        self.comp["lst_scroll"] = tk.Scrollbar(self.comp["lst_frame"], orient="vertical")
        self.comp["lst_scroll"].config(command=self.comp["lst"].yview)
        self.comp["lst_scroll"].pack(side="right", fill="y")
        self.comp["lst"].config(yscrollcommand=self.comp["lst_scroll"].set)
        self.comp["lst"].bind('<ButtonRelease-1>', lst_click)

        self.comp["lbut"] = tk.Button(self.mwin, font=self.parms["font.2"], text=" ◉ ", command=bf3,
            bg=self.parms["color.2"], fg=self.parms["color.4"] ) # lower button
        self.comp["lbut"].place( x=self.parms["lbut.x"], y=self.parms["lbut.y"] )
        self.comp["lent_strvar"] = tk.StringVar()
        self.comp["lent_strvar"].set("None")
        self.comp["lent"] = tkt.Entry(self.mwin, width=self.parms["lent.w"], font=self.parms["font.2"], textvariable=self.comp["lent_strvar"], state="readonly",
             foreground=self.parms["color.4"] )
        self.comp["lent"].place( x=self.parms["lent.x"], y=self.parms["lbut.y"] ) # lower entry

        self.comp["tbox"] = tk.Text( self.mwin, width=self.parms["tbox.w"], height=self.parms["tbox.h"], font=self.parms["font.2"],
            bg=self.parms["color.3"], fg=self.parms["color.0"] )
        self.comp["tbox"].place( x=self.parms["tbox.x"], y=self.parms["tbox.y"] ) # textbox

    # render gui, each bool parms decide render or not
    def render(self, mwin, awin):
        self.check()
        if mwin:
            self.comp["uent_strvar"].set( self.big[ self.current[0] ] ) # sep (big) renew
            self.comp["uent"].configure( textvariable=self.comp["uent_strvar"] )
            if self.editable[0]:
                self.comp["uent"].configure(state="normal")
            else:
                self.comp["uent"].configure(state="readonly")
            self.comp["lst"].delete( 0, self.comp["lst"].size() ) # sep (middle) renew, listbox
            for i in self.middle[ self.current[0] ]:
                self.comp["lst"].insert(self.comp["lst"].size(), i)
            if len( self.middle[ self.current[0] ] ) == 0: # sep (middle) renew, entry
                self.comp["lent_strvar"].set("NoData")
            else:
                self.comp["lent_strvar"].set( self.middle[ self.current[0] ][ self.current[1] ] )
            self.comp["lent"].configure( textvariable=self.comp["lent_strvar"] )
            if self.editable[1]:
                self.comp["lent"].configure(state="normal")
            else:
                self.comp["lent"].configure(state="readonly")
            self.comp["tbox"].configure(state="normal") # sep (small) renew
            self.comp["tbox"].delete('1.0', tk.END)
            if len( self.middle[ self.current[0] ] ) == 0:
                self.comp["tbox"].insert('1.0', "NoData")
            else:
                temp = self.small[ self.current[0] ][ self.current[1] ]
                if type(temp) == str:
                    self.comp["tbox"].insert('1.0', temp)
                else:
                    self.comp["tbox"].insert('1.0', "Non-String value")
            if self.editable[2]:
                self.comp["tbox"].configure(state="normal")
            else:
                self.comp["tbox"].configure(state="disabled")
            self.mwin.update()
        if awin and self.awin != None and self.mode == 0:
            self.comp["awin.pic.pimg"].destroy() # awin photo
            self.comp["awin.pic.pimg"] = tk.Canvas( self.awin, width=self.parms["awin.pic.x"], height=self.parms["awin.pic.y"], bg=self.parms["color.4"] )
            self.comp["awin.pic.pimg"].place(x=0, y=0)
            if len( self.middle[ self.current[0] ] ) > 0:
                if type( self.small[ self.current[0] ][ self.current[1] ] ) == bytes:
                    try: # generating photo object
                        pt = Image.open( io.BytesIO( self.small[ self.current[0] ][ self.current[1] ] ) )
                        if self.parms["awin.pic.x"] / self.parms["awin.pic.y"] > pt.size[0] / pt.size[1]:
                            px, py = int( pt.size[0] / pt.size[1] * self.parms["awin.pic.y"] ), self.parms["awin.pic.y"]
                            px = px - px % 2
                        else:
                            px, py = self.parms["awin.pic.x"], int( self.parms["awin.pic.x"] / pt.size[0] * pt.size[1] )
                            py = py - py % 2
                        pt = pt.resize( (px, py) )
                        self.comp["awin.pic.photo"] = ImageTk.PhotoImage(pt)
                        self.comp["awin.pic.pimg"].create_image( px//2, py//2, image=self.comp["awin.pic.photo"] )
                    except:
                        pass
            self.awin.update()

    # go inside of GUI mainloop
    def guiloop(self):
        self.render(True, False)
        self.mwin.mainloop()
        time.sleep(0.5)

    def check(self): # check render condition
        if len(self.menus) != 12:
            raise Exception("invalid GUI order")
        if len(self.big) != len(self.middle) or len(self.big) != len(self.small):
            raise Exception("invalid GUI order")
        for i in range( 0, len(self.big) ):
            if len( self.middle[i] ) != len( self.small[i] ):
                raise Exception("invalid GUI order")
        if self.current[0] < 0 or self.current[0] >= len(self.big):
            raise Exception("invalid GUI order")
        if self.current[1] < 0 or self.current[1] >= len( self.middle[ self.current[0] ] ):
            if self.current[1] != 0:
                raise Exception("invalid GUI order")
            
    def build(self): # make awin by mode
        if self.mode == 0: # picture view
            if self.awin == None:
                def exit_click():
                    time.sleep(0.1)
                    self.awin.destroy()
                    self.awin = None

                self.awin = tk.Toplevel(self.mwin)
                self.awin.title("Picture")
                self.awin.geometry( self.parms["awin.pic.size"] )
                self.awin.resizable(False, False)
                self.awin.configure( bg=self.parms["color.4"] )
                self.awin.protocol('WM_DELETE_WINDOW', exit_click)

                self.comp["awin.pic.photo"] = None # ImageTk photo
                self.comp["awin.pic.pimg"] = tk.Canvas( self.awin, width=self.parms["awin.pic.x"], height=self.parms["awin.pic.y"], bg=self.parms["color.4"] )
                self.comp["awin.pic.pimg"].place(x=0, y=0) # canvas

        elif self.mode == 1: # text search
            if self.awin == None:
                def exit_click(): # close button clicked
                    time.sleep(0.1)
                    self.awin.destroy()
                    self.awin = None

                def but_click(): # search button clicked
                    time.sleep(0.1)
                    temp = self.search( self.comp["awin.sch.tbox"].get('1.0', tk.END)[0:-1] )
                    self.comp["awin.sch.lst"].delete( 0, self.comp["awin.sch.lst"].size() )
                    for i in temp:
                        self.comp["awin.sch.lst"].insert(self.comp["awin.sch.lst"].size(), i)
                    self.awin.update()

                def lst_click(event): # listbox click & jump to text
                    time.sleep(0.1)
                    temp = self.comp["awin.sch.lst"].curselection()
                    if len(temp) > 0:
                        temp = self.comp["awin.sch.lst"].get(temp[0], temp[0] + 1)[0]
                        pos0, pos1 = int( temp[ temp.find("(") + 1:temp.find(",") ] ), int( temp[ temp.find(" ") + 1:temp.find(")") ] )
                        self.custom2(pos0, pos1)
                        self.current = [pos0, pos1]
                        self.render(True, True)

                self.awin = tk.Toplevel(self.mwin)
                self.awin.title("Search")
                self.awin.geometry( self.parms["awin.sch.size"] )
                self.awin.resizable(False, False)
                self.awin.configure( bg=self.parms["color.2"] )
                self.awin.protocol('WM_DELETE_WINDOW', exit_click)

                self.comp["awin.sch.tbox"] = tk.Text( self.awin, width=self.parms["awin.sch.tbox.w"], height=self.parms["awin.sch.tbox.h"], font=self.parms["font.2"],
                    bg=self.parms["color.2"], fg=self.parms["color.4"] )
                self.comp["awin.sch.tbox"].pack(side="top", fill="x") # textbox
                self.comp["awin.sch.tbox"].insert("1.0", "type here !!")

                self.comp["awin.sch.but"] = tk.Button(self.awin, font=self.parms["font.1"], text="Search", command=but_click,
                    bg=self.parms["color.2"], fg=self.parms["color.4"] )
                self.comp["awin.sch.but"].pack(side="top", fill="x") # button

                self.comp["awin.sch.frame"] = tk.Frame(self.awin) # listbox
                self.comp["awin.sch.frame"].place( x=self.parms["awin.sch.lst.x"], y=self.parms["awin.sch.lst.y"] )
                self.comp["awin.sch.lst"] = tk.Listbox(self.comp["awin.sch.frame"], width=self.parms["awin.sch.lst.w"], height=self.parms["awin.sch.lst.h"], font=self.parms["font.2"], selectmode="single",
                    bg=self.parms["color.2"], fg=self.parms["color.4"], selectbackground=self.parms["color.4"], selectforeground=self.parms["color.2"] )
                self.comp["awin.sch.lst"].pack(side="left", fill="y")
                self.comp["awin.sch.scroll"] = tk.Scrollbar(self.comp["awin.sch.frame"], orient="vertical")
                self.comp["awin.sch.scroll"].config(command=self.comp["awin.sch.lst"].yview)
                self.comp["awin.sch.scroll"].pack(side="right", fill="y")
                self.comp["awin.sch.lst"].config(yscrollcommand=self.comp["awin.sch.scroll"].set)
                self.comp["awin.sch.lst"].bind('<ButtonRelease-1>', lst_click)

        else: # recoder msgbox
            temp = f" Pos : {self.current}, Edit : {self.editable} \n Iter : ["
            for i in self.middle[0:-1]:
                temp = temp + f"{len(i)}, "
            temp = temp + f"{len(self.middle[-1])}] "
            tkm.showinfo("Recoder", temp)

    def search(self, words): # search from middle/small sep by words, str -> str[]
        res = [ ]
        words = [ x.lower().replace(" ", "").replace("\t", "") for x in words.split("\n") ]
        words = list( filter(lambda x: len(x) > 0, words) )
        for i in range( 0, len(self.big) ):
            for j in range( 0, len( self.middle[i] ) ):
                flag = True
                frag = self.middle[i][j].lower().replace(" ", "").replace("\t", "").replace("\n", "")
                for sch in words:
                    if sch in frag:
                        res.append(f"Title({i}, {j}) : {sch}")
                        flag = False
                        break
                if flag and type( self.small[i][j] ) == str:
                    frag = self.small[i][j].lower().replace(" ", "").replace("\t", "").replace("\n", "")
                    for sch in words:
                        if sch in frag:
                            res.append(f"Content({i}, {j}) : {sch}")
                            break
        return res

    # user custom func, menu option[x] 0~11
    def custom0(self, x):
        pass

    # user custom func, buttons clicked 0~3
    def custom1(self, x):
        pass

    # user custom func, search result move clicked, pos[x][y]
    def custom2(self, x, y):
        # custom2 is called "before" move
        pass
