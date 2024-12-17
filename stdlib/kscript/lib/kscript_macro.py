# test707 : stdlib5.kscript macro

import os
import time
import base64

import webbrowser
import pyautogui
import keyboard
import pypdf
import pdfkit

from selenium import webdriver
from selenium.webdriver.edge.service import Service

def getdrv(path): # get driver
    #옵션설정
    options = webdriver.EdgeOptions()
    options.add_argument('headless') #헤드레스만
    options.add_argument("disable-gpu") #헤드레스만
    options.add_argument('window-size=1920x1080')
    options.add_argument("user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/177.0.0.0 Safari/537.36 Edg/177.0.1938.62")

    #드라이버 설정
    service = Service(path)
    service.creation_flags = 0x08000000
    driver = webdriver.Edge(options=options,service=service)

    return driver

class pdf: # pdf reader/writer
    def __init__(self, path, mode):
        self.path, self.mode, self.pdf = os.path.abspath(path).replace("\\", "/"), mode, None
        if mode == "r":
            self.pdf = pypdf.PdfReader(self.path)
        elif mode == "w":
            self.pdf = pypdf.PdfWriter()

    def close(self): # close & save pdf
        if self.mode == "r":
            self.pdf.close()
        elif self.mode == "w":
            self.pdf.write(self.path)

    def getlen(self): # get len(pdf.pages)
        return len(self.pdf.pages)

    def getpage(self, num): # get pdf.page[num]
        return self.pdf.pages[num]

    def addpage(self, page): # add pdf.page
        if self.mode == "w":
            self.pdf.add_page(page)

    def addpdf(self, pdf): # add pdf file
        if self.mode == "w":
            self.pdf.append(pdf)

# 4-module worker, (pdf, web, mouse, key), !!! need edgedriver, mkhtmltopdf !!!
class lib:
    def __init__(self, u_pdf, u_web, u_mouse, u_key):
        self.u_pdf, self.u_web, self.u_mouse, self.u_key = u_pdf, u_web, u_mouse, u_key
        self.phnd = dict() # pdf handle (pdf[fullpath])
        self.drvpath = "msedgedriver.exe" # ! manual set drvpath !
        self.kitpath = "wkhtmltopdf.exe" # ! manual set wk !

    def pdf(self, pa, cd): # pdf func (500~510)
        if not self.u_pdf:
            return None, "not supported func"
        vo, er = None, None
        if cd == 500: # pdf.open
            if not self.check( pa, ["s", "s"] ):
                er = "type error"
            elif pa[1] != "r" and pa[1] != "w":
                er = "invalid option"
            else:
                path = os.path.abspath( pa[0] ).replace("\\", "/")
                if path in self.phnd:
                    self.phnd[path].close()
                self.phnd[path] = pdf( path, pa[1] )

        elif cd == 501: # pdf.close
            if not self.check( pa, ["s"] ):
                er = "type error"
            else:
                path = os.path.abspath( pa[0] ).replace("\\", "/")
                if path in self.phnd:
                    self.phnd[path].close()
                    del self.phnd[path]

        elif cd == 502: # pdf.len
            if not self.check( pa, ["s"] ):
                er = "type error"
            else:
                path = os.path.abspath( pa[0] ).replace("\\", "/")
                vo = self.phnd[path].getlen() if path in self.phnd else -1

        elif cd == 503: # pdf.addpage
            if not self.check( pa, ["s", "s", "i"] ):
                er = "type error"
            else:
                src = os.path.abspath( pa[0] ).replace("\\", "/")
                dst = os.path.abspath( pa[1] ).replace("\\", "/")
                if src not in self.phnd or dst not in self.phnd:
                    er = "not opened pdf"
                else:
                    self.phnd[dst].addpage( self.phnd[src].getpage( pa[2] ) )

        elif cd == 504: # pdf.addpdf
            if not self.check( pa, ["s", "s"] ):
                er = "type error"
            else:
                path = os.path.abspath( pa[1] ).replace("\\", "/")
                self.phnd[path].addpdf( pa[0] )

        elif cd == 505: # pdf.mkpdf
            if not self.check( pa, ["s", "s", "b"] ):
                er = "type error"
            else:
                cfg = pdfkit.configuration(wkhtmltopdf=self.kitpath)
                if pa[2]:
                    pt0 = '<!DOCTYPE html> <html> <head> <meta charset="utf-8"> <style> p { font-size: 24px; } </style> </head> <body> '
                    pt1 = ' </body> </html>'
                    pt2 = [ "<p>" + x + "</p>" for x in pa[1].split("\n") ]
                    pdfkit.from_string(pt0+" ".join(pt2)+pt1, pa[0], configuration=cfg)
                else:
                    pdfkit.from_string(pa[1], pa[0], configuration=cfg)

        else:
            er = "not supported func"
        return vo, er
    
    def web(self, pa, cd): # web func (510~520)
        if not self.u_web:
            return None, "not supported func"
        vo, er = None, None
        if cd == 510: # web.open
            webbrowser.open( pa[0] )

        elif cd == 511: # web.save
            if not self.check( pa, ["s", "s"] ):
                er = "type error"
            else:
                driver = getdrv(self.drvpath)
                driver.get( pa[0] )
                driver.implicitly_wait(90)
                time.sleep(0.5)
                with open(pa[1], 'wb') as f:
                    f.write( base64.b64decode( driver.print_page() ) )

        else:
            er = "not supported func"
        return vo, er
    
    def mouse(self, pa, cd): # mouse func (520~530)
        if not self.u_mouse:
            return None, "not supported func"
        vo, er = None, None
        if cd == 520: # mouse.delay
            if not self.check( pa, ["i|f"] ):
                er = "type error"
            elif pa[0] < 0:
                er = "invalid option"
            else:
                pyautogui.PAUSE = pa[0]

        elif cd == 521: # mouse.shot
            if not self.check( pa, ["s", "n|i", "n|i", "n|i", "n|i"] ):
                er = "type error"
            else:
                xs, ys = pyautogui.size()
                x_st = 0 if pa[1] == None else pa[1]
                y_st = 0 if pa[2] == None else pa[2]
                x_ed = xs if pa[3] == None else pa[3]
                y_ed = ys if pa[4] == None else pa[4]
                pyautogui.screenshot( pa[0], region=(x_st, y_st, x_ed-x_st, y_ed-y_st) )

        elif cd == 522: # mouse.size
            xs, ys = pyautogui.size()
            vo = xs if pa[0] else ys
            
        elif cd == 523: # mouse.pos
            xs, ys = pyautogui.position()
            vo = xs if pa[0] else ys

        elif cd == 524: # mouse.move
            if not self.check( pa, ["i|f", "i|f", "b", "n|i|f"] ):
                er = "type error"
            else:
                x, y = int( pa[0] ), int( pa[1] )
                t = 0.01 if pa[3] == None else pa[3]
                if pa[2]: # drag
                    pyautogui.drag(x, y, duration=t)
                else: # non-drag
                    pyautogui.move(x, y, duration=t)

        elif cd == 525: # mouse.moveto
            if not self.check( pa, ["i|f", "i|f", "b", "n|i|f"] ):
                er = "type error"
            else:
                x, y = int( pa[0] ), int( pa[1] )
                t = 0.01 if pa[3] == None else pa[3]
                if pa[2]: # drag
                    pyautogui.dragTo(x, y, duration=t)
                else: # non-drag
                    pyautogui.moveTo(x, y, duration=t)

        elif cd == 526: # mouse.click
            if not self.check( pa, ["i", "n|i|f", "b"] ):
                er = "type error"
            elif pa[0] < 0:
                er = "invalid option"
            else:
                b = "left" if pa[2] else "right"
                t = 0.01 if pa[1] == None else pa[1]
                pyautogui.click( button=b, clicks=pa[0], interval=t )

        elif cd == 527: # mouse.scroll
            if not self.check( pa, ["i|f"] ):
                er = "type error"
            else:
                pyautogui.scroll( int( pa[0] ) )

        else:
            er = "not supported func"
        return vo, er
    
    def key(self, pa, cd): # key func (530~540)
        if not self.u_key:
            return None, "not supported func"
        vo, er = None, None
        if cd == 530: # key.write
            if not self.check( pa, ["s", "n|i|f"] ):
                er = "type error"
            else:
                if pa[1] == None:
                    keyboard.write( pa[0] )
                else:
                    keyboard.write( pa[0], delay=pa[1] )

        elif cd == 531: # key.press
            if not self.check( pa, ["s", "n|i|f"] ):
                er = "type error"
            else:
                if pa[1] == None:
                    keyboard.press_and_release( pa[0] )
                else:
                    keyboard.press( pa[0] )
                    time.sleep( pa[1] )
                    keyboard.release( pa[0] )

        elif cd == 532: # key.set
            if not self.check( pa, ["s", "b"] ):
                er = "type error"
            else:
                if pa[1] and not keyboard.is_pressed( pa[0] ):
                    keyboard.press( pa[0] )
                elif not pa[1] and keyboard.is_pressed( pa[0] ):
                    keyboard.release( pa[0] )

        elif cd == 533: # key.status
            vo = keyboard.is_pressed( pa[0] )

        elif cd == 534: # key.event
            vo = keyboard.read_key()

        elif cd == 535: # key.wait
            keyboard.wait( pa[0] )
            
        else:
            er = "not supported func"
        return vo, er

    def check(self, pa, tp): # check parms types (nbifsc | a)
        if len(pa) != len(tp):
            return False
        pos = 0
        for i in pa:
            flag = False
            for j in tp[pos].split("|"):
                if j == "a":
                    flag = True
                elif j == "n" and i == None:
                    flag = True
                elif j == "b" and type(i) == bool:
                    flag = True
                elif j == "i" and type(i) == int:
                    flag = True
                elif j == "f" and type(i) == float:
                    flag = True
                elif j == "s" and type(i) == str:
                    flag = True
                elif j == "c" and type(i) == bytes:
                    flag = True
                if flag:
                    break
            if not flag:
                return False
            pos = pos + 1
        return True

    # close all pdf handles
    def exit(self):
        for i in self.phnd:
            try:
                self.phnd[i].close()
            except:
                pass
        self.__init__(False, False, False, False)

    # outer func works
    def run(self, parms, icode):
        vout, err = None, None
        try:
            if icode < 500:
                err = "not supported func"
            elif icode < 510: # pdf
                vout, err = self.pdf(parms, icode)
            elif icode < 520: # web
                vout, err = self.web(parms, icode)
            elif icode < 530: # mouse
                vout, err = self.mouse(parms, icode)
            elif icode < 540: # key
                vout, err = self.key(parms, icode)
            else:
                err = "not supported func"
        except Exception as e:
            vout, err = None, f"critical : {e}"
        if err == "":
            err = None
        return vout, err
