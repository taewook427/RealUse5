# test733 : extension.kmap

import os
import time
import math
import threading

import tkinter as tk
import tkinter.filedialog as tkf
import tkinter.messagebox as tkm

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

class camera:
    def __init__(self):
        self.pos, self.angle, self.mul, self.fov = [0, 0, 0], [0, 0], 1.0, math.pi / 2

    def precalc(self): # calc some parms
        self.sin_theta, self.cos_theta, self.sin_phi, self.cos_phi = math.sin(-self.angle[0]), math.cos(-self.angle[0]), math.sin(-self.angle[1]), math.cos(-self.angle[1])
        self.rxd, self.ryd = math.tan(0.5 * self.fov), math.tan(0.333333 * self.fov)
        self.xd, self.yd = -360 / self.rxd, -240 / self.ryd

    def relpos(self, node): # convert to relative pos
        x, y, z = (node.pos[0] - self.pos[0]) / self.mul, (node.pos[1] - self.pos[1]) / self.mul, (node.pos[2] - self.pos[2]) / self.mul
        x, z = self.cos_theta * x - self.sin_theta * z, self.sin_theta * x + self.cos_theta * z
        x, y = self.cos_phi * x - self.sin_phi * y, self.sin_phi * x + self.cos_phi * y
        return x, y, z

    def drawnode(self, node): # x, y, sz, far
        x, y, z = self.relpos(node)
        if x <= 0:
            return 0, 0, -1, 0
        mx, my, far = z / x, y / x, math.sqrt(x * x + y * y + z * z)
        x2d, y2d = self.xd * mx + 360, self.yd * my + 240
        if -self.rxd < mx < self.rxd and -self.ryd < my < self.ryd:
            sz = node.size / far * 240
            return x2d, y2d, max(min(sz, 239), 1), far
        else:
            return 0, 0, -1, 0
        
    def drawline(self, line): # x0, y0, x1, y1, sz, far
        x0, y0, z0 = self.relpos(line.node[0])
        x1, y1, z1 = self.relpos(line.node[1])
        dx, dy, dz = x1 - x0, y1 - y0, z1 - z0
        d2 = math.sqrt(dx * dx + dy * dy + dz * dz)
        num = max(int(d2 * 5), 1)
        start, end, step = None, None, 1 / num
        for i in range(0, num):
            temp = i * step
            x, y, z = x0 + dx * temp, y0 + dy * temp, z0 + dz * temp
            if x > 0:
                mx, my = z / x, y / x
                x2d, y2d = self.xd * mx + 360, self.yd * my + 240
                if -self.rxd < mx < self.rxd and -self.ryd < my < self.ryd:
                    if start == None:
                        start = (x2d, y2d)
                    else:
                        end = (x2d, y2d)
        if start == None or end == None or d2 == 0:
            return 0, 0, 0, 0, -1, 0
        else:
            tx, ty, tz = y0 * dz - z0 * dy, z0 * dx - x0 * dz, x0 * dy - y0 * dx
            far = math.sqrt(tx * tx + ty * ty + tz * tz) / d2
            if far == 0:
                return 0, 0, 0, 0, -1, 0
            else:
                return start[0], start[1], end[0], end[1], max(min(line.size / far * 240, 239), 1), far

class node:
    def __init__( self, idx="object", pos=(0, 0, 0) ):
        self.id, self.lbl, self.pos, self.size, self.color = idx, "", pos, 1, "black"

class line:
    def __init__( self, idx="object", node=(None, None) ):
        self.id, self.lbl, self.node, self.size, self.color, self.len = idx, "", node, 1, "black", ""

class mainclass:
    def __init__(self, path):
        self.cam, self.dots, self.lines, self.temp, self.viewlbl, self.viewcord, self.loop = camera(), [ ], [ ], dict(), False, False, True
        with open(path, "r", encoding="utf-8") as f:
            self.setup( f.read() )
        self.temp = None
        self.entry()
        self.mainloop()
    
    def setup(self, text):
        while "<" in text:
            a, b, c, text = readmd(text)
            if a == "node":
                for i in range( 0, len(b) ):
                    worker = kdb.toolbox()
                    worker.read( c[i] )
                    n = node( idx=b[i], pos=( worker.get("x")[3], worker.get("y")[3], worker.get("z")[3] ) )
                    if "lbl" in worker.name:
                        n.lbl = worker.get("lbl")[3]
                    if "color" in worker.name:
                        n.color = worker.get("color")[3]
                    if "r" in worker.name:
                        n.size = worker.get("r")[3]
                    self.dots.append(n)
                    self.temp[ b[i] ] = n

            elif a == "line":
                for i in range( 0, len(b) ):
                    worker = kdb.toolbox()
                    worker.read( c[i] )
                    l = line( idx=b[i], node=( self.temp[ worker.get("start")[3] ], self.temp[ worker.get("end")[3] ] ) )
                    if "lbl" in worker.name:
                        l.lbl = worker.get("lbl")[3]
                    if "color" in worker.name:
                        l.color = worker.get("color")[3]
                    if "r" in worker.name:
                        l.size = worker.get("r")[3]
                    x, y, z = l.node[0].pos[0] - l.node[1].pos[0], l.node[0].pos[1] - l.node[1].pos[1], l.node[0].pos[2] - l.node[1].pos[2]
                    l.len = f"({math.sqrt(x*x+y*y+z*z):.1f})"
                    self.lines.append(l)

            elif a == "config":
                 for i in range( 0, len(b) ):
                    if b[i] == "baseline":
                        num = max(float( c[i] ), 10)
                        l0 = line( node=( node( pos=(-num, 0, 0) ), node( pos=(num, 0, 0) ) ) )
                        l1 = line( node=( node( pos=(0, -num, 0) ), node( pos=(0, num, 0) ) ) )
                        l2 = line( node=( node( pos=(0, 0, -num) ), node( pos=(0, 0, num) ) ) )
                        l0.color, l1.color, l2.color = "black", "blue", "red"
                        self.lines.append(l0)
                        self.lines.append(l1)
                        self.lines.append(l2)
                    elif b[i] == "viewlbl":
                        self.viewlbl = True if "t" in c[i].lower() else False
                    elif b[i] == "viewcord":
                        self.viewcord = True if "t" in c[i].lower() else False
                    elif b[i] == "pos":
                        self.cam.pos = [ float(x) for x in c[i].replace(" ", "").split(",") ]
                    elif b[i] == "mul":
                        self.cam.mul = max(float( c[i] ), 0.01)
                    elif b[i] == "fov":
                        self.cam.fov = max(min(float( c[i] ) / 180 * math.pi, 3.05433), 0.5236)

    def entry(self):
        self.mwin = tk.Tk()
        self.mwin.title("Kmap")
        self.text = tk.StringVar()
        self.text.set("WASD Space Shift L")
        self.label = tk.Label(self.mwin, font=("Consolas", 14), textvariable=self.text)
        self.label.pack(pady=10)
        self.canvas = tk.Canvas(self.mwin, width=720, height=480, bg="white")
        self.canvas.pack()

        def exit_click(): # close button clicked
            time.sleep(0.1)
            if tkm.askyesno("Exit Check", " Do you want to exit Kmap? "):
                self.mwin.destroy()
                self.mwin, self.loop = None, False
        def click_start(event): # mouse click start
            self.last_x, self.last_y = event.x, event.y
        def click_end(event): # mouse click end
            self.last_x, self.last_y = None, None
        self.mwin.protocol('WM_DELETE_WINDOW', exit_click)
        self.last_x, self.last_y = None, None
        self.mwin.bind('<Key>', self.presskey)
        self.canvas.bind('<Button-1>', click_start)
        self.canvas.bind('<ButtonRelease-1>', click_end)
        self.canvas.bind('<B1-Motion>', self.mouseclick)
        self.canvas.bind('<MouseWheel>', self.mousesrcoll)

    def mainloop(self):
        while self.loop and self.mwin != None:
            t = time.time()
            self.render()
            time.sleep( max(0.03 - (time.time() - t), 0.01) )

    def genlbl(self):
        temp = f"(x {self.cam.pos[0]:.1f}, y {self.cam.pos[1]:.1f}, z {self.cam.pos[2]:.1f})\n"
        return temp + f"(θ {self.cam.angle[0]*180/math.pi:.1f}, φ {self.cam.angle[1]*180/math.pi:.1f}), x({self.cam.mul:.2f}), [{self.cam.fov*180/math.pi:.1f}]"
    
    def render(self):
        self.canvas.delete("all")
        self.cam.precalc()
        self.temp = ""
        def dotf():
            try:
                self.dot_queue = [ ]
                for i in self.dots:
                    a, b, c, d = self.cam.drawnode(i)
                    self.dot_queue.append( [a, b, c, d, i] )
                self.dot_queue.sort(key=lambda x:x[3], reverse=True)
            except Exception as e:
                self.temp = self.temp + str(e)
        def linef():
            try:
                self.line_queue = [ ]
                for i in self.lines:
                    a, b, c, d, e, f = self.cam.drawline(i)
                    self.line_queue.append( [a, b, c, d, e, f, i] )
                self.line_queue.sort(key=lambda x:x[5], reverse=True)
            except Exception as e:
                self.temp = self.temp + str(e)
        thr0 = threading.Thread(target=linef)
        thr1 = threading.Thread(target=dotf)
        thr0.start()
        thr1.start()
        thr0.join()
        for i in self.line_queue:
            if i[4] > 0:
                self.canvas.create_line(i[0], i[1], i[2], i[3], fill=i[6].color, width=i[4])
                if self.viewlbl:
                    text = f"{i[6].lbl} {i[6].len}" if self.viewcord else i[6].lbl
                    self.canvas.create_text((i[0]+i[2])/2, (i[1]+i[3])/2-10, fill=i[6].color, text=text)
        thr1.join()
        for i in self.dot_queue:
            if i[2] > 0:
                self.canvas.create_oval(i[0]-i[2], i[1]-i[2], i[0]+i[2], i[1]+i[2], fill=i[4].color)
                if self.viewlbl:
                    text = f"{i[4].lbl} ({i[4].pos[0]:.1f}, {i[4].pos[1]:.1f}, {i[4].pos[2]:.1f})" if self.viewcord else i[4].lbl
                    self.canvas.create_text(i[0], i[1]+10, fill=i[4].color, text=text)
        self.mwin.update()
        if self.temp != "":
            self.loop = tkm.askretrycancel(title="render error", message=f" {self.temp} \n Do you want to render again? ")

    def presskey(self, event): # press WASD Shift Space -+
        key = event.keysym.lower()
        if key == "w" or key == "up":
            self.cam.pos[0] = self.cam.pos[0] + self.cam.mul * math.cos(self.cam.angle[0])
            self.cam.pos[2] = self.cam.pos[2] + self.cam.mul * math.sin(self.cam.angle[0])
        elif key == "s" or key == "down":
            self.cam.pos[0] = self.cam.pos[0] - self.cam.mul * math.cos(self.cam.angle[0])
            self.cam.pos[2] = self.cam.pos[2] - self.cam.mul * math.sin(self.cam.angle[0])
        elif key == "a" or key == "left":
            self.cam.pos[0] = self.cam.pos[0] - self.cam.mul * math.sin(self.cam.angle[0])
            self.cam.pos[2] = self.cam.pos[2] + self.cam.mul * math.cos(self.cam.angle[0])
        elif key == "d" or key == "right":
            self.cam.pos[0] = self.cam.pos[0] + self.cam.mul * math.sin(self.cam.angle[0])
            self.cam.pos[2] = self.cam.pos[2] - self.cam.mul * math.cos(self.cam.angle[0])
        elif "space" in key:
            self.cam.pos[1] = self.cam.pos[1] + self.cam.mul
        elif "shift" in key:
            self.cam.pos[1] = self.cam.pos[1] - self.cam.mul
        elif key == "minus":
            self.cam.fov = max(self.cam.fov - 0.1, 0.5236)
        elif key == "equal":
            self.cam.fov = min(self.cam.fov + 0.1, 3.05433)
        elif key == "l":
            self.viewlbl, self.viewcord = not self.viewlbl, not self.viewlbl
        elif key == "bracketleft":
            self.cam.mul = self.cam.mul / 1.1
        elif key == "bracketright":
            self.cam.mul = self.cam.mul * 1.1
        self.text.set( self.genlbl() )

    def mouseclick(self, event): # mouse click & move
        if self.last_x != None and self.last_y != None:
            c0, c1 = self.cam.fov / 5, 180 / math.pi
            dx, dy = (event.x - self.last_x) * c0, (event.y - self.last_y) * c0
            self.last_x, self.last_y = event.x, event.y
            self.cam.angle[0], self.cam.angle[1] = self.cam.angle[0] * c1, self.cam.angle[1] * c1
            self.cam.angle[0], self.cam.angle[1] = (self.cam.angle[0] - dx + 180) % 360 - 180, (self.cam.angle[1] - dy + 180) % 360 - 180
            self.cam.angle[1] = max( -90, min( 90, self.cam.angle[1] ) )
            self.cam.angle[0], self.cam.angle[1] = self.cam.angle[0] / c1, self.cam.angle[1] / c1
        self.text.set( self.genlbl() )

    def mousesrcoll(self, event): # mouse scroll
        self.cam.mul = self.cam.mul * math.exp(event.delta / 360)
        self.text.set( self.genlbl() )

if __name__ == "__main__":
    kobj.repath()
    if not os.path.exists("./_ST5_DATA/"):
        os.mkdir("./_ST5_DATA/")
        with open("./_ST5_DATA/example.txt", "w", encoding="utf-8") as f:
            f.write('<node> {obj0} (lbl="home"; color="black"; r=5; x=50.0; y=0.0; z=0.0)\n')
            f.write('{obj1} (lbl="school"; color="blue"; r=2.5; x=0; y=50; z=0)\n')
            f.write('{obj2} (lbl="shop"; r=2.5; x=0; y=0; z=50)\n')
            f.write('<line> {obj3} (lbl="brick"; color="red"; r=1.5; start="obj0"; end="obj1")\n')
            f.write('{obj4} (lbl="dirt"; r=1.5; start="obj1"; end="obj2")\n')
            f.write('{obj5} (lbl="stone"; color="red"; r=1.5; start="obj2"; end="obj0")\n')
            f.write('<config> {baseline} (500) {viewlbl} (True) {viewcord} (True) {pos} (1, 1, 1) {mul} (1.0) {fov} (90.0)\n')
    worker = mainclass(tkf.askopenfile(filetypes=( ("Text File", "*.txt"), ("All File", "*.*") ), initialdir="./_ST5_DATA/").name)
    time.sleep(0.5)
