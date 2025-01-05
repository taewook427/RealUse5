# test734 : extension.convs
# 리눅스 빌드 추가 : --hidden-import='PIL._tkinter_finder'

import os
import subprocess
import time

import tkinter as tk
import tkinter.ttk as tkt
import tkinter.filedialog as tkf
import tkinter.messagebox as tkm
from PIL import Image, ImageTk, ImageOps
import blocksel

import pypdf
import pdfkit

import kobj
import kdb

def sel_file_one(mode): # select one file (img, vid, pdf)
    ft = [ ("All Files", "*.*") ]
    if mode == "img":
        ft = [ ("JPG Files", "*.jpg"), ("PNG Files", "*.png"), ("WEBP Files", "*.webp"), ("BMP Files", "*.bmp") ] + ft
    elif mode == "vid":
        ft = [ ("MP4 Files", "*.mp4"), ("WEBM Files", "*.webm"), ("AVI Files", "*.avi"), ("MKV Files", "*.mkv"), ("WMV Files", "*.wmv") ] + ft
    elif mode == "vid2":
        ft = [ ("MP4 Files", "*.mp4"), ("MP3 Files", "*.mp3"), ("GIF Files", "*.gif") ] + ft
    elif mode == "pdf":
        ft = [ ("PDF Files", "*.pdf") ] + ft
    return tkf.askopenfile(title="Select File", filetypes=ft).name

def sel_file_mutli(): # select generic files
    names, ft = [ ], [ ("JPG Files", "*.jpg"), ("PNG Files", "*.png"), ("WEBP Files", "*.webp"), ("BMP Files", "*.bmp"), ("MP4 Files", "*.mp4"), ("PDF Files", "*.pdf"), ("All Files", "*.*") ]
    for i in tkf.askopenfiles(title="Select Files", filetypes=ft):
        names.append( i.name.replace("\\", "/") )
    return names

def imgxy(x, y): # resize (x, y) -> (400, 400)
    nx, ny = 400, 400
    if x > y:
        ny = int(y * 200 / x) * 2
    else:
        nx = int(x * 200 / y) * 2
    return nx, ny

class img_resize:
    def __init__(self, iswin, ff):
        self.mwin, self.comp, self.parms, self.path_ff, self.path_image, self.img_object = None, dict(), {"color":"alice blue"}, ff, "", None
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "410x620+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["ent.0"], self.parms["ent.3"], self.parms["combo.3"], self.parms["ent.4"] = 33, 23, 7, 7
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "540x700+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["ent.3"], self.parms["combo.3"], self.parms["ent.4"] = 33, 25, 7, 7
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def sel_file(): # select image file
            time.sleep(0.1)
            self.path_image = sel_file_one("img")
            self.comp["strvar0"].set(self.path_image)
            self.img_object = Image.open(self.path_image)
            self.showimg()

        def merge_up(): # merge direction up
            time.sleep(0.1)
            tgt = Image.open( sel_file_one("img") )
            x0, y0 = self.img_object.size
            x1, y1 = tgt.size
            x, y = max(x0, x1), y0 + y1
            temp = Image.new( 'RGBA', (x, y) )
            temp.paste( tgt, (0, 0) )
            temp.paste( self.img_object, (0, y1) )
            self.img_object = temp
            self.showimg()

        def merge_down(): # merge direction down
            time.sleep(0.1)
            tgt = Image.open( sel_file_one("img") )
            x0, y0 = self.img_object.size
            x1, y1 = tgt.size
            x, y = max(x0, x1), y0 + y1
            temp = Image.new( 'RGBA', (x, y) )
            temp.paste( self.img_object, (0, 0) )
            temp.paste( tgt, (0, y0) )
            self.img_object = temp
            self.showimg()

        def merge_left(): # merge direction left
            time.sleep(0.1)
            tgt = Image.open( sel_file_one("img") )
            x0, y0 = self.img_object.size
            x1, y1 = tgt.size
            x, y = x0 + x1, max(y0, y1)
            temp = Image.new( 'RGBA', (x, y) )
            temp.paste( tgt, (0, 0) )
            temp.paste( self.img_object, (x1, 0) )
            self.img_object = temp
            self.showimg()

        def merge_right(): # merge direction right
            time.sleep(0.1)
            tgt = Image.open( sel_file_one("img") )
            x0, y0 = self.img_object.size
            x1, y1 = tgt.size
            x, y = x0 + x1, max(y0, y1)
            temp = Image.new( 'RGBA', (x, y) )
            temp.paste( self.img_object, (0, 0) )
            temp.paste( tgt, (x0, 0) )
            self.img_object = temp
            self.showimg()

        def save_img(): # save image
            time.sleep(0.1)
            x, y, path, tp, temp = self.comp["entry4a"].get(), self.comp["entry4b"].get(), self.comp["entry3"].get(), self.comp["combo3"].get(), self.img_object
            if x != "" and y != "":
                temp = temp.resize( ( int(x), int(y) ) )
            if "jpg" in tp:
                temp = temp.convert('RGB')
            if "_ll" in tp:
                tp = tp[:tp.rfind("_")]
                temp.save(f"./_ST5_DATA/{path}.{tp}", lossless=True)
            else:
                temp.save(f"./_ST5_DATA/{path}.{tp}")

        self.mwin = tk.Tk() # main window
        self.mwin.title("Convs Image Resize")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # file selector
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["button0"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], text=" . . . ", command=sel_file)
        self.comp["button0"].grid(row=0, column=0)
        self.comp["strvar0"] = tk.StringVar()
        self.comp["strvar0"].set("Select Image File")
        self.comp["entry0"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], textvariable=self.comp["strvar0"], width=self.parms["ent.0"], state="readonly")
        self.comp["entry0"].grid(row=0, column=1, padx=5)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # image info
        self.comp["frame1"].pack(fill="x", padx=5)
        self.comp["strvar1"] = tk.StringVar()
        self.comp["strvar1"].set("00000 x 00000 0000000000 (000.0 MiB)")
        self.comp["label1"] = tk.Label( self.comp["frame1"], font=self.parms["font.1"], textvariable=self.comp["strvar1"], bg=self.parms["color"] )
        self.comp["label1"].grid(row=0, column=0)

        self.comp["photo"] = None # ImageTk photo
        self.comp["canvas"] = tk.Canvas(self.mwin, width=400, height=400)
        self.comp["canvas"].pack(padx=5, pady=5)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # merge menus
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["label2"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="Merge")
        self.comp["label2"].grid(row=0, column=0)
        self.comp["button2a"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text="   ↑   ", command=merge_up)
        self.comp["button2a"].grid(row=0, column=1, padx=5)
        self.comp["button2b"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text="   ↓   ", command=merge_down)
        self.comp["button2b"].grid(row=0, column=2, padx=5)
        self.comp["button2c"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text="   ←   ", command=merge_left)
        self.comp["button2c"].grid(row=0, column=3, padx=5)
        self.comp["button2d"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text="   →   ", command=merge_right)
        self.comp["button2d"].grid(row=0, column=4, padx=5)

        self.comp["frame3"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # name & ext
        self.comp["frame3"].pack(fill="x", padx=5, pady=5)
        self.comp["label3"] = tk.Label(self.comp["frame3"], font=self.parms["font.0"], bg=self.parms["color"], text="Name")
        self.comp["label3"].grid(row=0, column=0)
        self.comp["entry3"] = tk.Entry( self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.3"] )
        self.comp["entry3"].grid(row=0, column=1, padx=5)
        self.comp["combo3"] = tkt.Combobox( self.comp["frame3"], font=self.parms["font.1"], width=self.parms["combo.3"], values=["jpg", "png", "bmp", "webp", "webp_ll", "ico", "Manual"] )
        self.comp["combo3"].set("webp")
        self.comp["combo3"].grid(row=0, column=2, padx=5)

        self.comp["frame4"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # size & save
        self.comp["frame4"].pack(fill="x", padx=5)
        self.comp["label4"] = tk.Label(self.comp["frame4"], font=self.parms["font.0"], bg=self.parms["color"], text="Size (Width, Height)")
        self.comp["label4"].grid(row=0, column=0)
        self.comp["entry4a"] = tk.Entry( self.comp["frame4"], font=self.parms["font.1"], width=self.parms["ent.4"] )
        self.comp["entry4a"].grid(row=0, column=1, padx=5)
        self.comp["entry4b"] = tk.Entry( self.comp["frame4"], font=self.parms["font.1"], width=self.parms["ent.4"] )
        self.comp["entry4b"].grid(row=0, column=2, padx=5)
        self.comp["button4"] = tk.Button(self.comp["frame4"], font=self.parms["font.0"], text=" Save ", command=save_img)
        self.comp["button4"].grid(row=0, column=3, padx=5)

    def showimg(self):
        x, y = self.img_object.size
        sz = os.path.getsize(self.path_image)
        if sz < 1024:
            self.comp["strvar1"].set(f"{x} x {y} {sz} ({sz} B)")
        elif sz < 1048576:
            self.comp["strvar1"].set(f"{x} x {y} {sz} ({sz/1024:.1f} KiB)")
        else:
            self.comp["strvar1"].set(f"{x} x {y} {sz} ({sz/1048576:.1f} MiB)")
        x, y = imgxy(x, y)
        self.comp["photo"] = ImageTk.PhotoImage( self.img_object.resize( (x, y) ) )
        self.comp["canvas"].create_image( x//2, y//2, image=self.comp["photo"] )

class img_conv:
    def __init__(self, iswin, ff):
        self.mwin, self.comp, self.parms, self.path_ff, self.path_image, self.img_object = None, dict(), {"color":"alice blue"}, ff, "", None
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "410x620+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["ent.0"], self.parms["ent.3"], self.parms["combo.4"] = 33, 5, 24
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "540x700+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["ent.3"], self.parms["combo.4"] = 33, 5, 24
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def sel_file(): # select image file
            time.sleep(0.1)
            self.path_image = sel_file_one("img")
            self.comp["strvar0"].set(self.path_image)
            self.img_object = Image.open(self.path_image)
            self.showimg()
        
        def conv_left(): # convert turn left
            time.sleep(0.1)
            self.img_object = self.img_object.rotate(90, expand=True)
            self.showimg()

        def conv_right(): # convert turn right
            time.sleep(0.1)
            self.img_object = self.img_object.rotate(270, expand=True)
            self.showimg()

        def conv_ver(): # convert mirror vertical
            time.sleep(0.1)
            self.img_object = self.img_object.transpose(Image.FLIP_TOP_BOTTOM)
            self.showimg()

        def conv_hor(): # convert mirror horizontal
            time.sleep(0.1)
            self.img_object = self.img_object.transpose(Image.FLIP_LEFT_RIGHT)
            self.showimg()

        def save_img(): # save image
            time.sleep(0.1)
            x0, y0, x1, y1 = self.comp["entry3a"].get(), self.comp["entry3b"].get(), self.comp["entry3c"].get(), self.comp["entry3d"].get()
            effect, path, temp = self.comp["combo4"].get(), self.path_image.replace("\\", "/"), self.img_object
            if x0 != "" and y0 != "" and x1 != "" and y1 != "":
                temp = temp.crop( ( int(x0), int(y0), int(x1), int(y1) ) )
            if effect == "Color Black & White":
                temp = temp.convert("L")
            elif effect == "Color Invert":
                temp = ImageOps.invert( temp.convert("RGB") )
            elif effect == "Color Red +10%":
                temp = temp.convert("RGB")
                data = [ ]
                for i in temp.getdata():
                    data.append( ( min(int(i[0]*1.1), 255), i[1], i[2] ) )
                temp.putdata(data)
            elif effect == "Color Green +10%":
                temp = temp.convert("RGB")
                data = [ ]
                for i in temp.getdata():
                    data.append( ( i[0], min(int(i[1]*1.1), 255), i[2] ) )
                temp.putdata(data)
            elif effect == "Color Blue +10%":
                temp = temp.convert("RGB")
                data = [ ]
                for i in temp.getdata():
                    data.append( ( i[0], i[1], min(int(i[2]*1.1), 255) ) )
                temp.putdata(data)
            elif effect == "Anti Aliasing":
                temp = temp.resize(temp.size, Image.LANCZOS)
            temp.save("./_ST5_DATA/" + path[path.rfind("/")+1:], lossless=True)

        self.mwin = tk.Tk() # main window
        self.mwin.title("Convs Image Convert")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # file selector
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["button0"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], text=" . . . ", command=sel_file)
        self.comp["button0"].grid(row=0, column=0)
        self.comp["strvar0"] = tk.StringVar()
        self.comp["strvar0"].set("Select Image File")
        self.comp["entry0"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], textvariable=self.comp["strvar0"], width=self.parms["ent.0"], state="readonly")
        self.comp["entry0"].grid(row=0, column=1, padx=5)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # image info
        self.comp["frame1"].pack(fill="x", padx=5)
        self.comp["strvar1"] = tk.StringVar()
        self.comp["strvar1"].set("00000 x 00000 0000000000 (000.0 MiB)")
        self.comp["label1"] = tk.Label( self.comp["frame1"], font=self.parms["font.1"], textvariable=self.comp["strvar1"], bg=self.parms["color"] )
        self.comp["label1"].grid(row=0, column=0)

        self.comp["photo"] = None # ImageTk photo
        self.comp["canvas"] = tk.Canvas(self.mwin, width=400, height=400)
        self.comp["canvas"].pack(padx=5, pady=5)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # conv menus
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["label2"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="Rotate/Mirror")
        self.comp["label2"].grid(row=0, column=0)
        self.comp["button2a"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text="  ↺  ", command=conv_left)
        self.comp["button2a"].grid(row=0, column=1, padx=5)
        self.comp["button2b"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text="  ↻  ", command=conv_right)
        self.comp["button2b"].grid(row=0, column=2, padx=5)
        self.comp["button2c"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text="  ↕  ", command=conv_ver)
        self.comp["button2c"].grid(row=0, column=3, padx=5)
        self.comp["button2d"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text="  ↔  ", command=conv_hor)
        self.comp["button2d"].grid(row=0, column=4, padx=5)

        self.comp["frame3"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # cut menus
        self.comp["frame3"].pack(fill="x", padx=5, pady=5)
        self.comp["label3"] = tk.Label(self.comp["frame3"], font=self.parms["font.0"], bg=self.parms["color"], text="Cut (x0 y0 x1 y1)")
        self.comp["label3"].grid(row=0, column=0)
        self.comp["entry3a"] = tk.Entry( self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.3"] )
        self.comp["entry3a"].grid(row=0, column=1, padx=5)
        self.comp["entry3b"] = tk.Entry( self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.3"] )
        self.comp["entry3b"].grid(row=0, column=2, padx=5)
        self.comp["entry3c"] = tk.Entry( self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.3"] )
        self.comp["entry3c"].grid(row=0, column=3, padx=5)
        self.comp["entry3d"] = tk.Entry( self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.3"] )
        self.comp["entry3d"].grid(row=0, column=4, padx=5)

        self.comp["frame4"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # effect
        self.comp["frame4"].pack(fill="x", padx=5, pady=5)
        self.comp["label4"] = tk.Label(self.comp["frame4"], font=self.parms["font.0"], bg=self.parms["color"], text="Effect")
        self.comp["label4"].grid(row=0, column=0)
        self.comp["combo4"] = tkt.Combobox( self.comp["frame4"], font=self.parms["font.1"], width=self.parms["combo.4"],
            values=["No Effect", "Color Black & White", "Color Invert", "Color Red +10%", "Color Green +10%", "Color Blue +10%", "Anti Aliasing"] )
        self.comp["combo4"].set("No Effect")
        self.comp["combo4"].grid(row=0, column=1, padx=5)
        self.comp["button4"] = tk.Button(self.comp["frame4"], font=self.parms["font.0"], text=" Save ", command=save_img)
        self.comp["button4"].grid(row=0, column=3, padx=5)

    def showimg(self):
        x, y = self.img_object.size
        sz = os.path.getsize(self.path_image)
        if sz < 1024:
            self.comp["strvar1"].set(f"{x} x {y} {sz} ({sz} B)")
        elif sz < 1048576:
            self.comp["strvar1"].set(f"{x} x {y} {sz} ({sz/1024:.1f} KiB)")
        else:
            self.comp["strvar1"].set(f"{x} x {y} {sz} ({sz/1048576:.1f} MiB)")
        x, y = imgxy(x, y)
        self.comp["photo"] = ImageTk.PhotoImage( self.img_object.resize( (x, y) ) )
        self.comp["canvas"].create_image( x//2, y//2, image=self.comp["photo"] )

class vod_resize:
    def __init__(self, iswin, ff):
        self.mwin, self.comp, self.parms, self.path_video = None, dict(), {"color":"alice blue"}, ""
        for i in os.listdir(ff):
            if "ffmpeg" in i:
                self.path_ff = os.path.abspath(ff + i)
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "550x180+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["ent.0"], self.parms["combo.1a"], self.parms["combo.1b"], self.parms["combo.1c"] = 47, 12, 8, 8
            self.parms["combo.2a"], self.parms["combo.2b"], self.parms["ent.3"], self.parms["combo.3"] = 12, 8, 29, 8
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "660x220+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["combo.1a"], self.parms["combo.1b"], self.parms["combo.1c"] = 43, 12, 8, 8
            self.parms["combo.2a"], self.parms["combo.2b"], self.parms["ent.3"], self.parms["combo.3"] = 12, 8, 27, 7
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def sel_file(): # select video
            time.sleep(0.1)
            self.path_video = sel_file_one("vid")
            self.comp["strvar0"].set(self.path_video)

        def save_vod(): # save video
            time.sleep(0.1)
            codec, bit, frame, res, aud = self.comp["combo1a"].get(), self.comp["combo1b"].get(), self.comp["combo1c"].get(), self.comp["combo2a"].get(), self.comp["combo2b"].get()
            if "vaapi" in codec:
                order = f'{self.path_ff} -hwaccel vaapi -hwaccel_device /dev/dri/renderD128 -i "{self.path_video}" -y'
            else:
                order = f'{self.path_ff} -i "{self.path_video}" -y'
            if codec == "default":
                order = order + " -c:v copy"
            else:
                order = order + f" -c:v {codec}"
            if frame != "default":
                order = order + f" -r {frame}"
            if res != "default":
                order = order + f' -vf "scale={res}"'
            if bit == "default":
                pass
            elif bit == "auto":
                if res == "default":
                    res = "1280x720"
                x, y = int( res[ :res.find("x") ] ), int( res[res.find("x")+1:] )
                f = 30 if frame == "default" else int(frame)
                order = order + f" -b:v {max(500 * int(x * y * f / 5000000), 500)}k"
            else:
                order = order + f" -b:v {bit}"
            if aud != "default":
                order = order + f" -b:a {aud}"
            order = order + f' "./_ST5_DATA/{self.comp["entry3"].get()}.{self.comp["combo3"].get()}"'
            result = subprocess.run(order, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True, text=True, encoding="utf-8")
            with open("./_ST5_DATA/log.txt", "w", encoding="utf-8") as f:
                f.write(f"[order]\n{order}\n{time.strftime("%Y.%m.%d;%H:%M:%S",time.localtime(time.time()))}\n")
                f.write(f"[stdout]\n{result.stdout}\n")
                f.write(f"[stderr]\n{result.stderr}\n")
        
        self.mwin = tk.Tk() # main window
        self.mwin.title("Convs Video Resize")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # file selector
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["button0"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], text=" . . . ", command=sel_file)
        self.comp["button0"].grid(row=0, column=0)
        self.comp["strvar0"] = tk.StringVar()
        self.comp["strvar0"].set("Select Video File")
        self.comp["entry0"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], textvariable=self.comp["strvar0"], width=self.parms["ent.0"], state="readonly")
        self.comp["entry0"].grid(row=0, column=1, padx=5)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # codec, bitrate, frame
        self.comp["frame1"].pack(fill="x", padx=5, pady=5)
        self.comp["label1a"] = tk.Label(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Codec")
        self.comp["label1a"].grid(row=0, column=0)
        self.comp["combo1a"] = tkt.Combobox( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["combo.1a"],
            values=["default", "libx264", "libx265", "h264_nvenc", "hevc_nvenc", "h264_amf", "h264_vaapi", "h264_qsv", "Manual"] )
        self.comp["combo1a"].set("default")
        self.comp["combo1a"].grid(row=0, column=1, padx=5)
        self.comp["label1b"] = tk.Label(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Bitrate")
        self.comp["label1b"].grid(row=0, column=2)
        self.comp["combo1b"] = tkt.Combobox( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["combo.1b"],
            values=["auto", "default", "1000k", "1500k", "2500k", "4000k", "6000k", "8500k", "Manual"] )
        self.comp["combo1b"].set("auto")
        self.comp["combo1b"].grid(row=0, column=3, padx=5)
        self.comp["label1c"] = tk.Label(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Frame")
        self.comp["label1c"].grid(row=0, column=4)
        self.comp["combo1c"] = tkt.Combobox( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["combo.1c"],
            values=["default", "10", "15", "30", "48", "60", "Manual"] )
        self.comp["combo1c"].set("default")
        self.comp["combo1c"].grid(row=0, column=5, padx=5)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # resolution, audio
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["label2a"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="Resolution")
        self.comp["label2a"].grid(row=0, column=0)
        self.comp["combo2a"] = tkt.Combobox( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["combo.2a"],
            values=["default", "256x144", "480x360", "640x480", "1280x720", "1920x1080", "2560x1440", "3840x2160", "Manual"] )
        self.comp["combo2a"].set("default")
        self.comp["combo2a"].grid(row=0, column=1, padx=5)
        self.comp["label2b"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="Audiorate")
        self.comp["label2b"].grid(row=0, column=2)
        self.comp["combo2b"] = tkt.Combobox( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["combo.2b"],
            values=["default", "96k", "128k", "160k", "192k", "Manual"] )
        self.comp["combo2b"].set("default")
        self.comp["combo2b"].grid(row=0, column=3, padx=5)

        self.comp["frame3"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # name, ext, save
        self.comp["frame3"].pack(fill="x", padx=5, pady=5)
        self.comp["label3"] = tk.Label(self.comp["frame3"], font=self.parms["font.0"], bg=self.parms["color"], text="Name")
        self.comp["label3"].grid(row=0, column=0)
        self.comp["entry3"] = tk.Entry( self.comp["frame3"], font=self.parms["font.1"], width=self.parms["ent.3"] )
        self.comp["entry3"].grid(row=0, column=1, padx=5)
        self.comp["combo3"] = tkt.Combobox( self.comp["frame3"], font=self.parms["font.1"], width=self.parms["combo.3"], values=["mp4", "avi", "mkv", "wmv", "Manual"] )
        self.comp["combo3"].set("mp4")
        self.comp["combo3"].grid(row=0, column=2, padx=5)
        self.comp["button3"] = tk.Button(self.comp["frame3"], font=self.parms["font.0"], text=" Save ", command=save_vod)
        self.comp["button3"].grid(row=0, column=3, padx=5)

class vod_conv:
    def __init__(self, iswin, ff):
        self.mwin, self.comp, self.parms, self.path_video = None, dict(), {"color":"alice blue"}, ""
        for i in os.listdir(ff):
            if "ffmpeg" in i:
                self.path_ff = os.path.abspath(ff + i)
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "550x180+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["ent.0"], self.parms["ent.1"], self.parms["combo.1"], self.parms["ent.2"], self.parms["combo.2"] = 47, 11, 8, 29, 8
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "660x220+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["ent.1"], self.parms["combo.1"], self.parms["ent.2"], self.parms["combo.2"] = 43, 10, 8, 27, 7
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def sel_file(): # select video
            time.sleep(0.1)
            self.path_video = sel_file_one("vid2")
            self.comp["strvar0"].set(self.path_video)

        def save_vod(): # save video
            time.sleep(0.1)
            st, ed, fr, nm = self.comp["entry1a"].get(), self.comp["entry1b"].get(), self.comp["combo1"].get(), self.comp["entry2"].get()
            mode, order = self.comp["combo2"].get(), f'{self.path_ff} -i "{self.path_video}" -y'
            if st != "":
                order = order + f" -ss {st}"
            if ed != "":
                order = order + f" -to {ed}"
            if mode == "mp4-mp4":
                order = order + f' -c:v copy -c:a copy "./_ST5_DATA/{nm}.mp4"'
            elif mode == "mp4-mp3":
                order = order + f' -vn -ar 44100 -ac 2 -ab 192 "./_ST5_DATA/{nm}.mp3"'
            elif mode == "mp3-mp4":
                order = order + f' -loop 1 -i bp.webp -vf "scale=640:360:force_original_aspect_ratio=decrease,pad=640:360:-1:-1:color=black,setsar=1,format=yuv420p" -shortest "./_ST5_DATA/{nm}.mp4"'
            elif mode == "mp4-gif":
                fr = "10" if fr == "default" else fr
                order = order + f' -vf "fps={fr},scale=640:-1:flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse" -loop 0 "./_ST5_DATA/{nm}.gif"'
            elif mode == "gif-mp4":
                order = order + f' -movflags faststart -pix_fmt yuv420p -vf "scale=trunc(iw/2)*2:trunc(ih/2)*2" "./_ST5_DATA/{nm}.mp4"'
            result = subprocess.run(order, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True, text=True, encoding="utf-8")
            with open("./_ST5_DATA/log.txt", "w", encoding="utf-8") as f:
                f.write(f"[order]\n{order}\n{time.strftime("%Y.%m.%d;%H:%M:%S",time.localtime(time.time()))}\n")
                f.write(f"[stdout]\n{result.stdout}\n")
                f.write(f"[stderr]\n{result.stderr}\n")
        
        self.mwin = tk.Tk() # main window
        self.mwin.title("Convs Video Convert")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # file selector
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["button0"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], text=" . . . ", command=sel_file)
        self.comp["button0"].grid(row=0, column=0)
        self.comp["strvar0"] = tk.StringVar()
        self.comp["strvar0"].set("Select Video File")
        self.comp["entry0"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], textvariable=self.comp["strvar0"], width=self.parms["ent.0"], state="readonly")
        self.comp["entry0"].grid(row=0, column=1, padx=5)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # time cut, frame
        self.comp["frame1"].pack(fill="x", padx=5, pady=5)
        self.comp["label1a"] = tk.Label(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Time (H:M:S.s)")
        self.comp["label1a"].grid(row=0, column=0)
        self.comp["entry1a"] = tk.Entry( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.1"] )
        self.comp["entry1a"].grid(row=0, column=1, padx=5)
        self.comp["entry1b"] = tk.Entry( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["ent.1"] )
        self.comp["entry1b"].grid(row=0, column=2, padx=5)
        self.comp["label1b"] = tk.Label(self.comp["frame1"], font=self.parms["font.0"], bg=self.parms["color"], text="Frame")
        self.comp["label1b"].grid(row=0, column=4)
        self.comp["combo1"] = tkt.Combobox( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["combo.1"],
            values=["default", "10", "15", "30", "48", "60", "Manual"] )
        self.comp["combo1"].set("default")
        self.comp["combo1"].grid(row=0, column=5, padx=5)

        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # name, ext, save
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["label2"] = tk.Label(self.comp["frame2"], font=self.parms["font.0"], bg=self.parms["color"], text="Name")
        self.comp["label2"].grid(row=0, column=0)
        self.comp["entry2"] = tk.Entry( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["ent.2"] )
        self.comp["entry2"].grid(row=0, column=1, padx=5)
        self.comp["combo2"] = tkt.Combobox( self.comp["frame2"], font=self.parms["font.1"], width=self.parms["combo.2"], values=["mp4-mp4", "mp4-mp3", "mp3-mp4", "mp4-gif", "gif-mp4"] )
        self.comp["combo2"].set("mp4-mp4")
        self.comp["combo2"].grid(row=0, column=2, padx=5)
        self.comp["button2"] = tk.Button(self.comp["frame2"], font=self.parms["font.0"], text=" Save ", command=save_vod)
        self.comp["button2"].grid(row=0, column=3, padx=5)

class pdf_conv:
    def __init__(self, iswin, ff):
        self.mwin, self.comp, self.parms = None, dict(), {"color":"alice blue"}
        self.path_ff, self.path_pdf, self.pdf_object, self.viewpos, self.pdf_reader = ff, "", [ ], 0, None
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "410x620+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["ent.0"], self.parms["w.1"], self.parms["h.1"], self.parms["ent.4"] = 33, 41, 10, 26
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "560x760+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["ent.0"], self.parms["w.1"], self.parms["h.1"], self.parms["ent.4"] = 34, 41, 8, 27
        self.entry()
        self.mwin.mainloop()
        if self.pdf_reader != None:
            self.pdf_reader.close()

    def entry(self):
        def sel_file(): # select pdf file
            time.sleep(0.1)
            self.path_pdf = sel_file_one("pdf")
            self.pdf_reader, self.pdf_object = pypdf.PdfReader(self.path_pdf), [ ]
            self.comp["list1"].delete( 0, self.comp["list1"].size() )
            for i in range( 0, len(self.pdf_reader.pages) ):
                self.pdf_object.append( self.pdf_reader.pages[i] )
                self.comp["list1"].insert(self.comp["list1"].size(), f"p{i}")
            self.viewpage()

        def sel_box(event): # select listbox
            time.sleep(0.1)
            self.viewpos = self.comp["list1"].curselection()[0]
            self.viewpage()
        
        def page_left(): # page turn left
            time.sleep(0.1)
            self.pdf_object[self.viewpos].rotate(270)
            self.viewpage()

        def page_right(): # page turn right
            time.sleep(0.1)
            self.pdf_object[self.viewpos].rotate(90)
            self.viewpage()

        def page_up(): # move page up
            time.sleep(0.1)
            if self.viewpos != 0:
                self.pdf_object[self.viewpos-1], self.pdf_object[self.viewpos] = self.pdf_object[self.viewpos], self.pdf_object[self.viewpos-1]
            self.viewpage()

        def page_down(): # move page down
            time.sleep(0.1)
            if self.viewpos != len(self.pdf_object)-1:
                self.pdf_object[self.viewpos+1], self.pdf_object[self.viewpos] = self.pdf_object[self.viewpos], self.pdf_object[self.viewpos+1]
            self.viewpage()

        def save_pdf(): # save pdf
            time.sleep(0.1)
            nm = "./_ST5_DATA/" + self.comp["entry4"].get() + ".pdf"
            with pypdf.PdfWriter(nm) as f:
                for i in self.pdf_object:
                    f.add_page(i)

        self.mwin = tk.Tk() # main window
        self.mwin.title("Convs PDF Convert")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # file selector
        self.comp["frame0"].pack(fill="x", padx=5, pady=5)
        self.comp["button0"] = tk.Button(self.comp["frame0"], font=self.parms["font.0"], text=" . . . ", command=sel_file)
        self.comp["button0"].grid(row=0, column=0)
        self.comp["strvar0"] = tk.StringVar()
        self.comp["strvar0"].set("Select PDF File")
        self.comp["entry0"] = tk.Entry(self.comp["frame0"], font=self.parms["font.1"], textvariable=self.comp["strvar0"], width=self.parms["ent.0"], state="readonly")
        self.comp["entry0"].grid(row=0, column=1, padx=5)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # page viewer
        self.comp["frame1"].pack(padx=5, pady=5)
        self.comp["list1"] = tk.Listbox( self.comp["frame1"], width=self.parms["w.1"],  height=self.parms["h.1"], font=self.parms["font.0"] )
        self.comp["list1"].pack(side="left", fill="y")
        self.comp["scroll1"] = tk.Scrollbar(self.comp["frame1"], orient="vertical")
        self.comp["scroll1"].config(command=self.comp["list1"].yview)
        self.comp["scroll1"].pack(side="right", fill="y")
        self.comp["list1"].config(yscrollcommand=self.comp["scroll1"].set)
        self.comp["list1"].bind("<ButtonRelease-1>", sel_box)

        self.comp["text1"] = tk.Text( self.mwin, width=self.parms["w.1"],  height=self.parms["h.1"], font=self.parms["font.0"] )
        self.comp["text1"].pack(side="top")
        self.comp["frame2"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # pdf page info
        self.comp["frame2"].pack(fill="x", padx=5, pady=5)
        self.comp["strvar2"] = tk.StringVar()
        self.comp["strvar2"].set("PageNum 000, Rotation 000")
        self.comp["label2"] = tk.Label( self.comp["frame2"], font=self.parms["font.1"], bg=self.parms["color"], textvariable=self.comp["strvar2"] )
        self.comp["label2"].grid(row=0, column=0)

        self.comp["frame3"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # rotation & move
        self.comp["frame3"].pack(fill="x", padx=5)
        self.comp["label3"] = tk.Label(self.comp["frame3"], font=self.parms["font.0"], bg=self.parms["color"], text="Rotate/Move")
        self.comp["label3"].grid(row=0, column=0)
        self.comp["button3a"] = tk.Button(self.comp["frame3"], font=self.parms["font.0"], text="  ↺  ", command=page_left)
        self.comp["button3a"].grid(row=0, column=1, padx=5)
        self.comp["button3b"] = tk.Button(self.comp["frame3"], font=self.parms["font.0"], text="  ↻  ", command=page_right)
        self.comp["button3b"].grid(row=0, column=2, padx=5)
        self.comp["button3c"] = tk.Button(self.comp["frame3"], font=self.parms["font.0"], text="  ↑  ", command=page_up)
        self.comp["button3c"].grid(row=0, column=3, padx=5)
        self.comp["button3d"] = tk.Button(self.comp["frame3"], font=self.parms["font.0"], text="  ↓  ", command=page_down)
        self.comp["button3d"].grid(row=0, column=4, padx=5)

        self.comp["frame4"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # save
        self.comp["frame4"].pack(fill="x", padx=5, pady=5)
        self.comp["label4"] = tk.Label(self.comp["frame4"], font=self.parms["font.0"], bg=self.parms["color"], text="Name")
        self.comp["label4"].grid(row=0, column=0)
        self.comp["entry4"] = tk.Entry( self.comp["frame4"], font=self.parms["font.1"], width=self.parms["ent.4"] )
        self.comp["entry4"].grid(row=0, column=1, padx=5)
        self.comp["button4"] = tk.Button(self.comp["frame4"], font=self.parms["font.0"], text=" Save ", command=save_pdf)
        self.comp["button4"].grid(row=0, column=2, padx=5)

    def viewpage(self):
        txt, rot = self.pdf_object[self.viewpos].extract_text(), self.pdf_object[self.viewpos].rotation
        self.comp["strvar2"].set(f"PageNum {self.viewpos}, Rotation {rot}")
        self.comp["text1"].delete("1.0", tk.END)
        self.comp["text1"].insert(tk.END, txt)

class file_conv:
    def __init__(self, iswin, ff):
        self.mwin, self.comp, self.parms, self.path_ff, self.path, self.viewpos = None, dict(), {"color":"alice blue"}, ff, [ ], 0
        for i in os.listdir(self.path_ff):
            if "wkhtmltopdf" in i:
                self.wk = os.path.abspath(self.path_ff + i)
            elif "ffmpeg" in i:
                self.ff = os.path.abspath(self.path_ff + i)
        if iswin: # windows
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "450x500+200+100", ("맑은 고딕", 12), ("Consolas", 14)
            self.parms["w"], self.parms["h"], self.parms["c"] = 40, 18, 14
        else: # linux
            self.parms["mwin.size"], self.parms["font.0"], self.parms["font.1"] = "550x600+200+100", ("맑은 고딕", 8), ("Consolas", 10)
            self.parms["w"], self.parms["h"], self.parms["c"] = 40, 13, 12
        self.entry()
        self.mwin.mainloop()

    def entry(self):
        def sel_file(): # select file
            time.sleep(0.1)
            self.path = sel_file_mutli()
            self.viewpath()

        def sel_box(event): # select listbox
            time.sleep(0.1)
            self.viewpos = self.comp["list0"].curselection()[0]

        def del_one(): # delete one selection
            time.sleep(0.1)
            del self.path[self.viewpos]
            self.viewpath()

        def del_all(): # delete all selection
            time.sleep(0.1)
            self.path = [ ]
            self.viewpath()

        def save_file(): # save file
            time.sleep(0.1)
            mode, name, log = self.comp["combo1"].get(), [ ], [ ]
            for i in self.path:
                temp = i[i.rfind("/")+1:]
                if "." in temp:
                    name.append( temp[ :temp.rfind(".") ] )
                else:
                    name.append(temp)
            if mode == "img -> webp":
                for i in range( 0, len(self.path) ):
                    try:
                        Image.open( self.path[i] ).save("./_ST5_DATA/" + name[i] + ".webp")
                        log.append(f"[webp] success")
                    except Exception as e:
                        log.append(f"[webp] err {e}")
                log = "\n".join(log)
            elif mode == "img -> webp_ll":
                for i in range( 0, len(self.path) ):
                    try:
                        Image.open( self.path[i] ).save("./_ST5_DATA/" + name[i] + ".webp", lossless=True)
                        log.append(f"[webp] success")
                    except Exception as e:
                        log.append(f"[webp] err {e}")
                log = "\n".join(log)
            elif mode == "img -> png":
                for i in range( 0, len(self.path) ):
                    try:
                        Image.open( self.path[i] ).save("./_ST5_DATA/" + name[i] + ".png")
                        log.append(f"[png] success")
                    except Exception as e:
                        log.append(f"[png] err {e}")
                log = "\n".join(log)
            elif mode == "img -> jpg":
                for i in range( 0, len(self.path) ):
                    try:
                        Image.open( self.path[i] ).convert('RGB').save("./_ST5_DATA/" + name[i] + ".jpg")
                        log.append(f"[jpg] success")
                    except Exception as e:
                        log.append(f"[jpg] err {e}")
                log = "\n".join(log)
            elif mode == "html -> pdf":
                cfg = pdfkit.configuration(wkhtmltopdf=self.wk)
                for i in range( 0, len(self.path) ):
                    try:
                        pdfkit.from_file(self.path[i], "./_ST5_DATA/"+name[i]+".pdf", configuration=cfg)
                        log.append(f"[html] success")
                    except Exception as e:
                        log.append(f"[html] err {e}")
                log = "\n".join(log)
            elif mode == "merge img":
                try:
                    imgs = [Image.open(x) for x in self.path]
                    log.append(f"[img] open {len(imgs)}")
                    if len(imgs) == 0:
                        temp = open("./_ST5_DATA/imgpdf_0.pdf", "wb")
                        temp.close()
                        log.append(f"[img] no pdf")
                    elif len(imgs) == 1:
                        imgs[0].save("./_ST5_DATA/imgpdf_1.pdf")
                        log.append(f"[img] 1 img pdf")
                    else:
                        imgs[0].save( f"./_ST5_DATA/imgpdf_{len(imgs)}.pdf", save_all=True, append_images=imgs[1:] )
                        log.append(f"[img] pdf of {len(imgs)} files")
                except Exception as e:
                    log.append(f"[img] err {e}")
                log = "\n".join(log)
            elif mode == "merge mp4":
                with open("./_ST5_DATA/list.txt", "w", encoding="utf-8") as f:
                    for i in self.path:
                        f.write(f"file '{i}'\n")
                order = f"{self.ff} -f concat -safe 0 -i ./_ST5_DATA/list.txt ./_ST5_DATA/video_{len(name)}.mp4"
                result = subprocess.run(order, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True, text=True, encoding="utf-8")
                log = f"[order]\n{order}\n{time.strftime("%Y.%m.%d;%H:%M:%S",time.localtime(time.time()))}\n[stdout]\n{result.stdout}\n[stderr]\n{result.stderr}\n"
            elif mode == "merge_cp mp4":
                with open("./_ST5_DATA/list.txt", "w", encoding="utf-8") as f:
                    for i in self.path:
                        f.write(f"file '{i}'\n")
                order = f"{self.ff} -f concat -safe 0 -i ./_ST5_DATA/list.txt -c copy ./_ST5_DATA/video_{len(name)}.mp4"
                result = subprocess.run(order, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True, text=True, encoding="utf-8")
                log = f"[order]\n{order}\n{time.strftime("%Y.%m.%d;%H:%M:%S",time.localtime(time.time()))}\n[stdout]\n{result.stdout}\n[stderr]\n{result.stderr}\n"
            elif mode == "merge pdf":
                try:
                    with pypdf.PdfWriter(f"./_ST5_DATA/pdfpdf_{len(self.path)}.pdf") as f:
                        for i in range( 0, len(self.path) ):
                            f.append( self.path[i] )
                            log.append(f"[pdf] add {name[i]}")
                except Exception as e:
                    log.append(f"[pdf] err {e}")
                log = "\n".join(log)
            with open("./_ST5_DATA/log.txt", "w", encoding="utf-8") as f:
                f.write(log)

        self.mwin = tk.Tk() # main window
        self.mwin.title("Convs File Convert")
        self.mwin.geometry( self.parms["mwin.size"] )
        self.mwin.resizable(False, False)
        self.mwin.configure( bg=self.parms["color"] )

        self.comp["frame0"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # selection viewer
        self.comp["frame0"].pack(padx=5, pady=5)
        self.comp["list0"] = tk.Listbox( self.comp["frame0"], width=self.parms["w"],  height=self.parms["h"], font=self.parms["font.1"] )
        self.comp["list0"].pack(side="left", fill="y")
        self.comp["scroll0"] = tk.Scrollbar(self.comp["frame0"], orient="vertical")
        self.comp["scroll0"].config(command=self.comp["list0"].yview)
        self.comp["scroll0"].pack(side="right", fill="y")
        self.comp["list0"].config(yscrollcommand=self.comp["scroll0"].set)
        self.comp["list0"].bind("<ButtonRelease-1>", sel_box)

        self.comp["frame1"] = tk.Frame( self.mwin, bg=self.parms["color"] ) # sel, del, clr, save
        self.comp["frame1"].pack(fill="x", padx=5, pady=5)
        self.comp["combo1"] = tkt.Combobox( self.comp["frame1"], font=self.parms["font.1"], width=self.parms["c"],
            values=["img -> webp", "img -> webp_ll", "img -> png", "img -> jpg", "html -> pdf", "merge img", "merge mp4", "merge_cp mp4", "merge pdf"] )
        self.comp["combo1"].set("merge img")
        self.comp["combo1"].grid(row=0, column=5, padx=5)
        self.comp["button1a"] = tk.Button(self.comp["frame1"], font=self.parms["font.0"], text=" . . . ", command=sel_file)
        self.comp["button1a"].grid(row=0, column=1, padx=5)
        self.comp["button1b"] = tk.Button(self.comp["frame1"], font=self.parms["font.0"], text=" DEL ", command=del_one)
        self.comp["button1b"].grid(row=0, column=2, padx=5)
        self.comp["button1c"] = tk.Button(self.comp["frame1"], font=self.parms["font.0"], text=" CLR ", command=del_all)
        self.comp["button1c"].grid(row=0, column=3, padx=5)
        self.comp["button1d"] = tk.Button(self.comp["frame1"], font=self.parms["font.0"], text=" Save ", command=save_file)
        self.comp["button1d"].grid(row=0, column=4, padx=5)

    def viewpath(self):
        self.comp["list0"].delete( 0, self.comp["list0"].size() )
        for i in self.path:
            self.comp["list0"].insert(self.comp["list0"].size(), i)

class mainclass(blocksel.toolbox):
    def __init__(self, iswin):
        super().__init__("Convs", iswin)
        path, self.selection = "../../_ST5_COMMON/iconpack/convs/", -1
        self.txts, self.curpos, self.upos, self.umsg = ["ImgResize", "ImgConv", "PDFctrl", "VidResize", "VidConv", "FileCtrl"], 1, 0, ["Select Mode"]
        self.pics = [path+"img_resize.png", path+"img_conv.png", path+"pdf_conv.png", path+"mp4_resize.png", path+"mp4_conv.png", path+"ext_conv.png"]

    def custom0(self, x):
        self.selection = x
        self.mwin.destroy()

kobj.repath()
cfg = kdb.toolbox()
with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
    cfg.read( f.read() )
iswin, flag = (cfg.get("dev.os")[3] == "windows"), True
if not os.path.exists("../../_ST5_COMMON/iconpack/"):
    flag = False
    tkm.showinfo(title="No Package", message=" Convs requires package >>> common.iconpack <<<. \n Install dependent package and start again. ")
if not os.path.exists("../../_ST5_COMMON/videopack/"):
    flag = False
    tkm.showinfo(title="No Package", message=" Convs requires package >>> common.videopack <<<. \n Install dependent package and start again. ")
if not os.path.exists("./_ST5_DATA/"):
    os.mkdir("./_ST5_DATA/")
if flag:
    worker = mainclass(iswin)
    worker.entry()
    worker.guiloop()
    if worker.selection == 0:
        t = img_resize(iswin, "../../_ST5_COMMON/videopack/")
    elif worker.selection == 1:
        t = img_conv(iswin, "../../_ST5_COMMON/videopack/")
    elif worker.selection == 2:
        t = pdf_conv(iswin, "../../_ST5_COMMON/videopack/")
    elif worker.selection == 3:
        t = vod_resize(iswin, "../../_ST5_COMMON/videopack/")
    elif worker.selection == 4:
        t = vod_conv(iswin, "../../_ST5_COMMON/videopack/")
    elif worker.selection == 5:
        t = file_conv(iswin, "../../_ST5_COMMON/videopack/")
time.sleep(0.5)
