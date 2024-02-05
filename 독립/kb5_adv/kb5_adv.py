# test607 : kboom5 adv

import io
import os
import shutil
import time
import random
import zlib

import tkinter
import tkinter.ttk
import tkinter.messagebox
from tkinter import filedialog
from PIL import Image, ImageTk

import threading as thr
import multiprocessing as mp

import kdb
import kaes
import kcom
import simen

# ===== kcom key transfer =====

# skey 4B, sport int -(fpath S, fkey 48B)-> (errmsg S)
def sendkey(skey, sport, fpath, fkey):
    fpath = os.path.abspath(fpath).replace("\\", "/")
    skey = skey * ( 48 // len(skey) )
    tbox0 = kdb.toolbox()
    tbox1 = simen.toolbox()
    tbox2 = kcom.server()
    
    tbox0.readstr("path = 0\nkey = 0")
    tbox0.fixdata("path", fpath)
    tbox0.fixdata("key", fkey)
    data = bytes(tbox0.writestr(), encoding="utf-8") # 평문 데이터
    crcv = str( zlib.crc32(data) ) # 평문 데이터의 crc32 값
    
    tbox1.setkey(skey)
    enc = tbox1.encrypt(data) # 전송할 암호문
    tbox2.close = 90 # 시간제한 90s
    tbox2.port = sport
    tbox2.msg = crcv

    try:
        tbox2.send(enc)
        errmsg = ""
    except Exception as e:
        errmsg = str(e)
    return errmsg

# skey 4B, sport int -> (errmsg S, fpath S, fkey 48B)
def recievekey(skey, sport):
    tbox0 = kdb.toolbox()
    tbox1 = simen.toolbox()
    tbox2 = kcom.client()

    skey = skey * ( 48 // len(skey) )
    tbox2.close = 20 # 시간제한 20s
    tbox2.port = sport

    try:
        enc = tbox2.recieve()
        crcv = tbox2.msg

        tbox1.setkey(skey)
        data = tbox1.decrypt(enc)
        if zlib.crc32(data) != int(crcv):
            raise Exception("wrong crc32 value")

        tbox0.readstr( str(data, encoding="utf-8") )
        fpath = tbox0.getdata("path")[3]
        fkey = tbox0.getdata("key")[3]
        
        errmsg = ""
    except Exception as e:
        errmsg = str(e)

    if errmsg == "":
        return "", fpath, fkey
    else:
        return errmsg, "", b""
    
# enc kf path S + key 48B -> plain kf B
def readkey(path, key):
    tbox = kaes.funcbytes()
    with open(path, "rb") as f:
        enc = f.read()
    return tbox.de(key, enc)

class mainclass():
    def __init__(self):
        # 공통 변수 설정
        self.password = b"0000"
        self.keyfile = kaes.genkf("기본키파일")
        self.keyfilepath = "기본키파일"
        self.hint = bytes("초기화 비밀번호 : 0000", encoding="utf-8")
        self.message = "Kboom5 사용을 환영합니다."
        self.folders = [ ]
        self.boomnum = 104857600
        self.filekeys = dict()
        self.filenames = [ ]

        self.loginsuccess = False

    def finit(self): # 최초 초기화
        self.folders = [ ]
        for i in os.listdir("./"):
            if os.path.isdir(i) and len(i) > 3:
                if i[-3:] == "_en":
                    if os.path.exists(i[0:-3] + "_de"):
                        self.folders.append( i[0:-3] )
        if os.path.exists("settings.webp"):
            with open("settings.webp", "rb") as f:
                temp = f.read()
        else:
            temp = self.pkgset()
            with open("settings.webp", "wb") as f:
                f.write(temp)
        return temp
    
    def pkgset(self): # 세팅 파일 패키징
        tbox0 = kdb.toolbox()
        tbox1 = kaes.genbytes()
        tbox1.msg = "Kboom5 settings file"

        temp = ["msg = 0", "boom = 0", "num = 0"]
        for i in range( 0, len(self.filekeys) ):
            temp.append(f"{i}.name = 0")
            temp.append(f"{i}.key = 0")
        tbox0.readstr( "\n".join(temp) )

        tbox0.fixdata("msg", self.message)
        tbox0.fixdata("boom", self.boomnum)
        tbox0.fixdata( "num", len(self.filekeys) )
        count = 0
        for i in self.filekeys:
            tbox0.fixdata(f"{count}.name", i)
            tbox0.fixdata( f"{count}.key", self.filekeys[i] )
            count = count + 1
        temp = tbox0.writestr()

        return tbox1.en( self.password, self.keyfile, self.hint, bytes(temp, encoding="utf-8") )

    def readpkg(self, enc, stpoint): # 세팅 파일 읽기 (0 : 성공, 1 : 실패)
        tbox0 = kdb.toolbox()
        tbox1 = kaes.genbytes()
        self.filekeys = dict()
        self.filenames = [ ]

        try:
            temp = str(tbox1.de(self.password, self.keyfile, enc, stpoint), encoding="utf-8")
            tbox0.readstr(temp)
            self.message = tbox0.getdata("msg")[3]
            self.boomnum = tbox0.getdata("boom")[3]
            num = tbox0.getdata("num")[3]
            for i in range(0, num):
                tnm = tbox0.getdata(f"{i}.name")[3]
                tky = tbox0.getdata(f"{i}.key")[3]
                self.filekeys[tnm] = tky
            return 0
        except:
            return 1

    def login(self, enc): # 로그인
        # 힌트 뷰어
        tbox = kaes.genbytes()
        try:
            self.hint, msg, stpoint = tbox.view(enc)
        except:
            self.hint = bytes("올바른 설정파일 아님", encoding="utf-8")
            stpoint = 0

        # 로그인 윈도우
        win = tkinter.Tk()
        win.title('KB5 login')
        win.geometry("300x200+100+50")
        win.resizable(False, False)

        # 프로그램 검증
        v0 = kaes.genbytes()
        v1 = kaes.genfile()
        v2 = kaes.funcbytes()
        v3 = kaes.funcfile()
        if not (v0.valid and v1.valid and v2.valid and v3.valid):
            tkinter.messagebox.showinfo(title='심각한 보안 경고', message=' 네이티브 가속 공유 라이브러리 파일이 \n 유효하지 않습니다. (dll/so) ')

        def resetkf():
            time.sleep(0.1)
            nonlocal strvar1
            try:
                temp = filedialog.askopenfile(title='파일 선택').name
                with open(temp, "rb") as f:
                    self.keyfile = f.read()
                self.keyfilepath = temp
                strvar1.set(temp)
            except:
                self.keyfile = kaes.genkf("기본키파일")
                self.keyfilepath = "기본키파일"
                strvar1.set("기본키파일")

        def resetkfremote():
            time.sleep(0.1)
            nonlocal strvar1
            sport, skey = kcom.unpack( ent1b.get() )
            err, fpath, fkey = recievekey(skey, sport)
            if err != "":
                self.keyfile = kaes.genkf("기본키파일")
                self.keyfilepath = "기본키파일"
                tkinter.messagebox.showinfo(title='다운로드 실패', message=f' 키 파일 다운로드 실패 \n {err} ')
            else:
                self.keyfile = readkey(fpath, fkey)
                self.keyfilepath = fpath
            strvar1.set(self.keyfilepath)

        def gologin():
            time.sleep(0.1)
            self.password = bytes(ent3.get(), encoding="utf-8")
            if self.readpkg(enc, stpoint) == 0:
                self.loginsuccess = True
                win.destroy()
            else:
                tkinter.messagebox.showinfo(title='wrong PWKF', message=' 비밀번호 또는 키 파일 데이터가 \n 일치하지 않습니다. ')

        but0 = tkinter.Button(win, font=("맑은 고딕", 12), text=". . .", command=resetkf)
        but0.place(x=5, y=5) # kf reset
        but0b = tkinter.Button(win, font=("맑은 고딕", 12), text=". . .", command=resetkfremote)
        but0b.place(x=5, y=55) # kf reset remote
        strvar1 = tkinter.StringVar()
        strvar1.set(self.keyfilepath)
        ent1 = tkinter.Entry(win, textvariable=strvar1, font=("맑은 고딕", 14), width=23, state="readonly")
        ent1.place(x=55, y=10) # kf path
        ent1b = tkinter.Entry(win, font=("맑은 고딕", 14), width=23)
        ent1b.place(x=55, y=60) # kf path remote
        lbl2 = tkinter.Label( win, font=("맑은 고딕", 14), text=str(self.hint, encoding="utf-8") )
        lbl2.place(x=5, y=105) # hint
        ent3 = tkinter.Entry(win, font=("맑은 고딕", 14), width=23, show="*")
        ent3.place(x=5, y=145) # pw input
        but4 = tkinter.Button(win, font=("맑은 고딕", 12), text=" Go ", command=gologin)
        but4.place(x=245, y=140)

        win.mainloop()

    def mainfunc(self): # GUI
        def sendfunc(): # fr0 보내기
            time.sleep(0.1)
            port = random.randrange(10000, 40000)
            key = kaes.genrandom(4)
            realkey = key * 12
            addr = kcom.pack(port, key)

            name = self.filenames[ listbox.curselection()[0] ]
            fkey = self.filekeys[name]
            pos = name.find("/")
            name = name[0:pos] + "_en/" + name[pos + 1:]
            ret = [""]

            def transmit(skey, sport, fname, fkey, resret): # 전송 인라인
                err = sendkey(skey, sport, fname, fkey)
                resret[0] = err

            x = thr.Thread( target=transmit, args=(key, port, name, fkey, ret) )
            x.start()
            tkinter.messagebox.showinfo(title='전송 시작', message=f' 90초 안에 다음 주소를 \n 수신 프로그램에 입력하세요. \n {addr} ')
            x.join()
            if ret[0] == "":
                ret[0] = "transfered successfully"
            tkinter.messagebox.showinfo(title='전송 결과', message=" " + ret[0] + " ")
            regen()

        def impfunc(): # fr0 가져오기
            time.sleep(0.1)
            res0 = self.fimp()
            res1 = self.fclear()
            temp = self.pkgset()
            with open("settings.webp", "wb") as f:
                f.write(temp)
            tkinter.messagebox.showinfo(title='import files', message=f' 파일 {res0} 개를 가져왔습니다. \n 드라이브 삭제 : {res1[0]} 개, 키 삭제 : {res1[1]} 개 ')
            regen()

        def expfunc(): # fr0 내보내기
            time.sleep(0.1)
            sel = listbox.curselection()
            if len(sel) != 0:
                self.fexp(sel)
            temp = self.pkgset()
            with open("settings.webp", "wb") as f:
                f.write(temp)
            tkinter.messagebox.showinfo(title='export files', message=f' 파일 {len(sel)} 개를 내보냈습니다. ')
            regen()

        def boomfunc(): # fr0 붐
            time.sleep(0.1)
            res = self.fclear()
            if tkinter.messagebox.askokcancel(title='export files', message=f' 드라이브 삭제 : {res[0]} 개, 키 삭제 : {res[1]} 개 \n boom 진행 시 모든 키 파일이 삭제되고 \n 저장소가 {self.boomnum} 바이트만큼 \n 덮어씌워집니다. 계속하시겠습니까? '):
                try:
                    res = self.fboom()
                except Exception as e:
                    res = str(e)
                tkinter.messagebox.showinfo(title='BOOM complete', message=f'{res}')
            regen()

        def viewfunc(): # fr1 view
            nonlocal cnvimg
            try:
                temp = self.filenames[ listbox.curselection()[0] ]
                pos = temp.find("/")
                tgt = temp[0:pos] + "_en/" + temp[pos + 1:]
                with open(tgt, "rb") as f:
                    tdata = f.read()
                tbox = kaes.funcbytes()
                idata = tbox.de(self.filekeys[temp], tdata)

                pimg = Image.open( io.BytesIO(idata) )
                iw, ih = pimg.size
                ratio = min(420 / iw, 420 / ih)
                sw, sh = int(iw * ratio), int(ih * ratio)
                rimg = pimg.resize( (sw, sh), Image.LANCZOS )
                cnvimg = ImageTk.PhotoImage(rimg)
                canvas.create_image(5, 5, anchor=tkinter.NW, image=cnvimg)
            except:
                pass

        def refig0(): # 메세지/붐 재설정 인라인
            self.message = tbox0.get('1.0', tkinter.END)[0:-1]
            self.boomnum = int( ent1.get() )
            multi = 1
            temp = cbox1.get().lower()[0]
            if temp == "b":
                multi = 1
            elif temp == "k":
                multi = 1024
            elif temp == "m":
                multi = 1048576
            elif temp == "g":
                multi = 1073741824
            self.boomnum = self.boomnum * multi

        def refig1(): # 비밀번호/힌트 재설정 인라인
            if chkvar.get() != 0:
                self.password = bytes(ent3.get(), encoding="utf-8")
            self.hint = bytes(ent4.get(), encoding="utf-8")

        def refig2(mode): # 패키지 저장
            temp = self.pkgset()
            with open("settings.webp", "wb") as f:
                f.write(temp)
            if mode == 0: # 메세지
                tkinter.messagebox.showinfo(title='config update', message=' 메세지/boom 설정이 \n 업데이트되었습니다. ')
            elif mode == 1: # 비밀번호
                tkinter.messagebox.showinfo(title='pwhint update', message=' 메세지/boom/비밀번호/힌트 \n 설정이 업데이트되었습니다. ')
            else: # 키파일
                tkinter.messagebox.showinfo(title='keyfile update', message=' 메세지/boom/비밀번호/힌트/키파일 \n 설정이 업데이트되었습니다. ')
            regen()

        def resetcfg(): # 설정 재설정
            time.sleep(0.1)
            refig0()
            refig2(0)

        def resetpw(): # 비밀번호 재설정
            time.sleep(0.1)
            refig0()
            refig1()
            refig2(1)

        def resetkf(): # 키파일 재설정
            time.sleep(0.1)
            try:
                temp = filedialog.askopenfile(title='파일 선택').name
                with open(temp, "rb") as f:
                    self.keyfile = f.read()
                self.keyfilepath = temp
            except:
                self.keyfile = kaes.genkf("기본키파일")
                self.keyfilepath = "기본키파일"
            strvar2.set(self.keyfilepath)
            refig0()
            refig1()
            refig2(2)

        def resetkfremote(): # 키파일 재설정 원격
            time.sleep(0.1)
            sport, skey = kcom.unpack( ent2b.get() )
            err, fpath, fkey = recievekey(skey, sport)
            if err != "":
                self.keyfile = kaes.genkf("기본키파일")
                self.keyfilepath = "기본키파일"
                tkinter.messagebox.showinfo(title='다운로드 실패', message=f' 키 파일 다운로드 실패 \n {err} ')
            else:
                self.keyfile = readkey(fpath, fkey)
                self.keyfilepath = fpath
            strvar2.set(self.keyfilepath)
            refig0()
            refig1()
            refig2(2)

        def pwshow(): # 비밀번호 보기/가리기
            time.sleep(0.1)
            if chkvar.get() == 0:
                strvar3.set( "*" * len( str(self.password, encoding="utf-8") ) )
                ent3.config(state="readonly")
            else:
                strvar3.set( str(self.password, encoding="utf-8") )
                ent3.config(state="normal")

        def regen(): # status/config 다시보기
            # fr0
            self.filenames = [ ]
            for i in self.folders:
                for j in os.listdir(i + "_en"):
                    if i + "/" + j in self.filekeys:
                        self.filenames.append(i + "/" + j)
            listbox.delete( 0, listbox.size() )
            for i in self.filenames:
                listbox.insert(listbox.size(), i)

            # fr1
            viewfunc()

            # fr2
            tbox0.delete(1.0, tkinter.END)
            tbox0.insert(1.0, self.message)
            strvar1.set( str(self.boomnum) )
            cbox1.set("Bytes")
            strvar2.set(self.keyfilepath)
            if chkvar.get() == 0:
                strvar3.set( "*" * len( str(self.password, encoding="utf-8") ) )
            else:
                strvar3.set( str(self.password, encoding="utf-8") )
            strvar4.set( str(self.hint, encoding="utf-8") )
            strvar5.set(f"kf size : {len(self.keyfile)}, folders : {len(self.folders)}\nfilekeys : {len(self.filekeys)}, filenames : {len(self.filenames)}")

            win.update()

        # 메인 윈도우
        win = tkinter.Tk()
        win.title('KB5_adv')
        win.geometry("450x480+200+100")
        win.resizable(False, False)

        # 프레임
        notebook = tkinter.ttk.Notebook(win, width=440, height=440)
        notebook.place(x=5, y=5)
        fr0 = tkinter.Frame(win)
        notebook.add(fr0, text="  status  ")
        fr1 = tkinter.Frame(win)
        notebook.add(fr1, text="  view  ")
        fr2 = tkinter.Frame(win)
        notebook.add(fr2, text="  config  ")

        # 사진 보기 연결
        def clickevent(event):
            if notebook.index( notebook.select() ) == 1:
                viewfunc()
        notebook.bind("<<NotebookTabChanged>>", clickevent)

        sendbut = tkinter.Button(fr0, text=" send ", font=("Consolas", 14), command=sendfunc)
        sendbut.place(x=5, y=5) # fr0 보내기
        impbut = tkinter.Button(fr0, text="import", font=("Consolas", 14), command=impfunc)
        impbut.place(x=115, y=5) # fr0 가져오기
        expbut = tkinter.Button(fr0, text="export", font=("Consolas", 14), command=expfunc)
        expbut.place(x=225, y=5) # fr0 내보내기
        boombut = tkinter.Button(fr0, text=" boom ", font=("Consolas", 14), command=boomfunc)
        boombut.place(x=335, y=5) # fr0 붐

        lstfr = tkinter.Frame(fr0)
        lstfr.place(x=5,y=50)
        listbox = tkinter.Listbox(lstfr, width=40,  height=14, font = ("맑은 고딕", 14), selectmode = 'extended')
        listbox.pack(side="left", fill="y")
        scbar = tkinter.Scrollbar(lstfr, orient="vertical")
        scbar.config(command=listbox.yview)
        scbar.pack(side="right", fill="y")
        listbox.config(yscrollcommand=scbar.set) # files listbox

        canvas = tkinter.Canvas(fr1, width=430, height=430)
        canvas.place(x=5, y=5) # fr1 사진창
        cnvimg = None

        tbox0 = tkinter.Text( fr2, width=41, height=3, font=("맑은 고딕", 14) )
        tbox0.place(x=10, y=5) # msg
        lbl1 = tkinter.Label(fr2, font=("맑은 고딕", 14), text="boom")
        lbl1.place(x=5, y=95)
        strvar1 = tkinter.StringVar()
        strvar1.set("")
        ent1 = tkinter.Entry(fr2, font=("맑은 고딕", 14), textvariable=strvar1, width=25)
        ent1.place(x=80, y=100)
        cbox1 = tkinter.ttk.Combobox(fr2, values=["Bytes", "KiB", "MiB", "GiB"], font=("맑은 고딕", 14), width=6, height=4)
        cbox1.place(x=350, y=100)
        cbox1.set("Bytes") # boomnum

        but2 = tkinter.Button(fr2, font=("맑은 고딕", 12), text=". . .", command=resetkf)
        but2.place(x=5, y=145) # kf reset
        strvar2 = tkinter.StringVar()
        strvar2.set("")
        ent2 = tkinter.Entry(fr2, textvariable=strvar2, font=("맑은 고딕", 14), width=36, state="readonly")
        ent2.place(x=60, y=150) # kf path
        but2b = tkinter.Button(fr2, font=("맑은 고딕", 12), text=". . .", command=resetkfremote)
        but2b.place(x=5, y=190) # kf reset remote
        ent2b = tkinter.Entry(fr2, font=("맑은 고딕", 14), width=36)
        ent2b.place(x=60, y=195) # kf path remote
        chkvar = tkinter.IntVar()
        chkbut = tkinter.Checkbutton(fr2, text="PW 보기", font=("맑은 고딕", 14), variable=chkvar, command=pwshow)
        chkbut.place(x=5, y=235) # pw show
        strvar3 = tkinter.StringVar()
        strvar3.set("")
        ent3 = tkinter.Entry(fr2, textvariable=strvar3, font=("맑은 고딕", 14), state="readonly", width=30)
        ent3.place(x=120, y=240) # pw input
        strvar4 = tkinter.StringVar()
        strvar4.set("")
        lbl4 = tkinter.Label(fr2, font=("맑은 고딕", 14), text="hint")
        lbl4.place(x=5, y=280)
        ent4 = tkinter.Entry(fr2, textvariable=strvar4, font=("맑은 고딕", 14), width=36)
        ent4.place(x=65, y=285) # hint input

        strvar5 = tkinter.StringVar()
        strvar5.set("kf size : 0, folders : 0\nfilekeys : 0, filenames : 0")
        lbl5 = tkinter.Label(fr2, font=("맑은 고딕", 14), textvariable=strvar5)
        lbl5.place(x=5, y=325) # debug info
        but6 = tkinter.Button(fr2, font=("맑은 고딕", 14), text=" 메세지 재설정 ", command=resetcfg)
        but6.place(x=5, y=390) # reset cfg
        but7 = tkinter.Button(fr2, font=("맑은 고딕", 14), text="비밀번호 재설정", command=resetpw)
        but7.place(x=225, y=390) # reset pw

        self.fclear()
        regen()
        win.mainloop()

    def fclear(self): # 클리어
        todel0 = [ ] # 드라이브에서 삭제
        todel1 = [ ] # 키 저장소에서 삭제
        for i in self.folders:
            for j in os.listdir(i + "_en"):
                if f"{i}/{j}" not in self.filekeys:
                    todel0.append(f"./{i}_en/{j}")
            for j in os.listdir(i + "_de"):
                todel0.append(f"./{i}_de/{j}")
        for i in self.filekeys:
            pos = i.find("/")
            if not os.path.exists( i[0:pos] + "_en/" + i[pos + 1:] ):
                todel1.append(i)
        for i in todel0:
            os.remove(i)
        for i in todel1:
            del self.filekeys[i]
        return [ len(todel0), len(todel1) ]
    
    def fimp(self): # 가져오기
        toadd = [ ]
        for i in self.folders:
            for j in os.listdir(i + "_de"):
                if f"{i}/{j}" not in self.filekeys:
                    toadd.append(f"{i}/{j}")
        tbox = kaes.funcfile()
        for i in toadd:
            key = kaes.genrandom(48)
            pos = i.find("/")
            fo = i[0:pos]
            fn = i[pos + 1:]
            tbox.en(key, fo+"_de/"+fn, fo+"_en/"+fn)
            self.filekeys[i] = key
        return len(toadd)
    
    def fexp(self, nums): # 내보내기
        tbox = kaes.funcfile()
        for i in nums:
            name = self.filenames[i]
            pos = name.find("/")
            fo = name[0:pos]
            fn = name[pos + 1:]
            tbox.de(self.filekeys[name], f"{fo}_en/{fn}", f"{fo}_de/{fn}")

    def fboom(self): # 붐
        count0 = 0 # 현재 폴더 지우기 성공수
        count1 = 0 # 실패수
        count2 = self.boomnum # 남은 덮어쓰기 크기
        for i in os.listdir("./"):
            try:
                if os.path.isdir(i):
                    shutil.rmtree(i)
                else:
                    os.remove(i)
                count0 = count0 + 1
            except:
                count1 = count1 + 1
        flag = True
        while flag:
            size = count2 // 2
            try:
                with open(f"{size}.bin", "wb") as f:
                    f.write(b"\xff" * size)
                count2 = count2 - size
                if count2 < 512:
                    flag = False
            except:
                flag = False
        try:
            with open(f"final.bin", "wb") as f:
                f.write(b"\x00" * count2)
            count2 = 0
        except:
            pass
        return f" 삭제 실패 : {count1} / {count0 + count1} \n 덮어쓰기 실패 : {count2} / {self.boomnum} "

if __name__ == "__main__":
    mp.freeze_support()
    classloader = mainclass()
    encbdata = classloader.finit()
    classloader.login(encbdata)
    time.sleep(0.1)
    if classloader.loginsuccess:
        classloader.mainfunc()
    time.sleep(0.5)
