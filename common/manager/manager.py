# test717 : common.manager

import os
import shutil
import time
import requests
import webbrowser

import tkinter.messagebox
import reader

import kobj
import kdb
import ksc
import kpkg
import kcom

class mainclass(reader.toolbox):
    def __init__(self):
        self.setup()
        super().__init__("Starter5 Manager", self.config_db[6]=="windows", 1)
        self.menus = ["Delete", "Install", "Update", "WebPage"] + ["        "] * 8
        self.get_status()
        self.get_install()
        self.get_packages()
        self.get_update()
        self.entry()
        self.current = [1, 0]
        self.guiloop()

    def setup(self): # get st5 config & sign
        worker = kdb.toolbox()
        with open("../../_ST5_CONFIG.txt", "r", encoding="utf-8") as f:
            worker.read( f.read() )
        self.config_db = [ worker.get('path.export')[3], worker.get('path.desktop')[3], worker.get('path.local')[3],
                        worker.get("url.download")[3], worker.get("url.info")[3], worker.get("url.help")[3],
                        worker.get("dev.os")[3], worker.get("dev.activate")[3] ]
        worker = kdb.toolbox()
        with open("../../_ST5_SIGN.txt", "r", encoding="utf-8") as f:
            worker.read( f.read() )
        num, self.sign_db = 0, [ ]
        while f"{num}.name" in worker.name:
            num = num + 1
        for i in range(0, num):
            a, b, c = worker.get(f"{i}.name")[3], worker.get(f"{i}.phash")[3], worker.get(f"{i}.public")[3]
            if ksc.crc32hash( bytes(c, encoding="utf-8") ) == b:
                self.sign_db.append( [a, b, c] ) # nm, ph, pub
            else:
                tkinter.messagebox.showerror(title="Invalid Sign", message=f" Sign {a} with phash {b.hex()} does not match to public key. ")

    def check_pkg(self, path, pkgr): # check if sign is valid, chmod 755
        if pkgr.public == "":
            tkinter.messagebox.showinfo(title="No Sign Warning", message=" This package does not have ksign data! ")
        else:
            pub, fl = pkgr.public, True
            for i in self.sign_db:
                if i[2] == pub:
                    fl = False
                    break
            if fl:
                tkinter.messagebox.showinfo(title="Untrusted Sign Warning", message=" Starter5 signDB does not contains the sign of this package! ")
        if self.config_db[6] != "windows":
            for i in os.listdir(path):
                if os.path.isfile(path + i):
                    os.chmod(path + i, 0o755)

    def get_status(self): # add gui config & sign data
        self.big.append("Status")
        self.middle.append( [ ] )
        self.small.append( [ ] )
        self.middle[0].append("Config : path")
        self.small[0].append(f"<입출력 폴더>\n{self.config_db[0]}\n\n<바탕화면>\n{self.config_db[1]}\n\n<로컬 홈>\n{self.config_db[2]}")
        self.middle[0].append("Config : url")
        self.small[0].append(f"<다운로드 저장소>\n{self.config_db[3]}\n\n<패키지 저장소>\n{self.config_db[4]}\n\n<도움말 페이지>\n{self.config_db[5]}")
        self.middle[0].append("Config : dev")
        self.small[0].append(f"<로컬 운영체제>\n{self.config_db[6]}\n\n<개발자 모드>\n{self.config_db[7]}")
        for i in self.sign_db:
            self.middle[0].append(f"Sign : {i[0]}")
            self.small[0].append(f"<공개키 식별값>\n{i[1].hex()}\n\n<공개키>\n{i[2]}")

    def get_install(self): # get installed package data
        self.install_db = [ ]
        self.big.append("Installed")
        self.middle.append( [ ] )
        self.small.append( [ ] )
        def sub(path):
            for i in [ x[:-1] if x[-1] == "/" else x for x in os.listdir(path) ]:
                worker = kdb.toolbox()
                if os.path.exists(path + i + "/_ST5_VERSION.txt"):
                    with open(path + i + "/_ST5_VERSION.txt", "r", encoding="utf-8") as f:
                        worker.read( f.read() )
                else:
                    worker.read('name = "None"; version = 0.0; text = "Invalid Package"; release = "00000000"; download = "00000000"')
                a, b, c, d, e, f = worker.get('name')[3], worker.get('version')[3], worker.get('text')[3], worker.get('release')[3], worker.get('download')[3], path + i + "/"
                self.install_db.append( [a, b, c, d, e, f] ) # nm ver txt rel dwn path
                self.middle[1].append(f"{i}")
                self.small[1].append(f"Package : {a} ver{b}\nInfo : {c}\nRelease : {d}\nDownload : {e}\nPath : {f}")
        sub("../../_ST5_EXTENSION/")
        sub("../../_ST5_COMMON/")

    def get_packages(self): # get whole package data
        self.big.append("All Packages")
        self.middle.append( [ ] )
        self.small.append( [ ] )
        def sub(domain, path):
            worker, temp, num = kdb.toolbox(), [ ], 0
            try:
                requests.get(self.config_db[4], timeout=5)
                worker.read( kcom.gettxt(self.config_db[4], domain) )
            except:
                worker.read(f'0.name = "None"; 0.version = 0.0; 0.devonly = False; 0.release = "00000000"; 0.text = "Web access fail."; 0.num = 1;')
            while f"{num}.name" in worker.name:
                num = num + 1
            for i in range(0, num):
                a, b, c, d = worker.get(f"{i}.name")[3], worker.get(f"{i}.version")[3], worker.get(f"{i}.text")[3], worker.get(f"{i}.release")[3]
                e, f, g = worker.get(f"{i}.devonly")[3], worker.get(f"{i}.num")[3], path + worker.get(f"{i}.name")[3] + "/"
                if self.config_db[7] or (not e):
                    temp.append( [a, b, c, d, f, g] ) # nm ver txt rel num path
                    self.middle[2].append(a)
                    self.small[2].append(f"Package : {a} ver{b}\nInfo : {c}\nRelease : {d}\nType : {domain}, dev {e}")
            return temp
        self.package_db = sub("extension", "../../_ST5_EXTENSION/") + sub("common", "../../_ST5_COMMON/")

    def get_update(self): # get update data
        db0, db1 = dict(), dict()
        for i in self.install_db:
            db0[ i[0] ] = i[1]
        for i in self.package_db:
            db1[ i[0] ] = i[1]
        self.big.append("Update Queue")
        self.middle.append( [ ] )
        self.small.append( [ ] )
        for i in db0:
            if i in db1:
                if db1[i] - db0[i] > 0.00001:
                    self.middle[3].append(i)
                    self.small[3].append(f"Current : ver{db0[i]}\nNew : ver{db1[i]}")

    def dirctrl(self, path, gen): # folder control
        if os.path.exists(path):
            shutil.rmtree(path)
        if gen:
            os.mkdir(path)

    def custom0(self, x):
        if x == 0 and self.current[0] == 1: # delete
            pos = self.current[1]
            if tkinter.messagebox.askokcancel(title="Delete Package", message=f" Are you sure to delete package {self.install_db[pos][0]} at {self.install_db[pos][5]}? "):
                shutil.rmtree(self.install_db[pos][5])
                self.middle[1], self.small[1], self.install_db = self.middle[1][:pos] + self.middle[1][pos+1:], self.small[1][:pos] + self.small[1][pos+1:], self.install_db[:pos] + self.install_db[pos+1:]
            
        elif x == 1 and self.current[0] == 2: # install
            pos, temp = self.current[1], [x[0] for x in self.install_db]
            if self.package_db[pos][0] in temp:
                tkinter.messagebox.showerror(title="Package Exists", message=f" Cannot download package because \n Package {self.package_db[pos][0]} is already installed. ")
            else:
                temp = self.package_db[pos]
                self.dirctrl("./temp625/", True)
                tkinter.messagebox.showinfo(title="Package Download", message=f" Manager will download Package {temp[0]}. \n This work takes time. ")
                kcom.download( self.config_db[3], temp[0], temp[4], "./temp625/temp.bin", [0] )
                worker = kpkg.toolbox()
                worker.osnum = 1 if self.config_db[6] == "windows" else 2
                try:
                    path = worker.unpack("./temp625/temp.bin")
                except Exception as e:
                    tkinter.messagebox.showerror(title="Install Fail", message=f" Error occurred while unpacking. \n {e}")

                os.rename( path, temp[5] )
                self.check_pkg(temp[5], worker)
                a, b, c, d, e, f = temp[0], temp[1], temp[2], temp[3], worker.dwn_date, temp[5]
                self.install_db.append( [a, b, c, d, e, f] ) # nm ver txt rel dwn path
                self.middle[1].append(f"{a}")
                self.small[1].append(f"Package : {a} ver{b}\nInfo : {c}\nRelease : {d}\nDownload : {e}\nPath : {f}")
                time.sleep(0.1)
                tkinter.messagebox.showinfo(title="Install Complete", message=f" Package {a} ver{b} is installed at {f}. ")
                self.dirctrl("./temp625/", False)
                self.dirctrl("./temp674/", False)

        elif x == 2 and self.current[0] == 3: # update
            pos, name = self.current[1], self.middle[3][ self.current[1] ]
            pos0, pos1 = [x[0] for x in self.install_db].index(name), [x[0] for x in self.package_db].index(name)
            temp = self.package_db[pos1]
            tkinter.messagebox.showinfo(title="Package Download", message=f" Manager will download Package {name}. \n This work takes time. ")
            self.dirctrl("./_ST5_DATA/", False)
            self.dirctrl("./temp625/", True)
            if os.path.exists(temp[5] + "_ST5_DATA/"):
                os.rename(temp[5] + "_ST5_DATA/", "./_ST5_DATA/")
            shutil.rmtree( temp[5] )
            kcom.download( self.config_db[3], temp[0], temp[4], "./temp625/temp.bin", [0] )
            worker = kpkg.toolbox()
            worker.osnum = 1 if self.config_db[6] == "windows" else 2
            try:
                path = worker.unpack("./temp625/temp.bin")
            except Exception as e:
                tkinter.messagebox.showerror(title="Update Fail", message=f" Error occurred while unpacking. \n {e}")

            os.rename( path, temp[5] )
            self.check_pkg(temp[5], worker)
            self.dirctrl(temp[5] + "_ST5_DATA/", False)
            if os.path.exists("./_ST5_DATA/"):
                os.rename("./_ST5_DATA/", temp[5] + "_ST5_DATA/")
            self.middle[3], self.small[3] = self.middle[3][:pos] + self.middle[3][pos+1:], self.small[3][:pos] + self.small[3][pos+1:]
            a, b, c, d, e, f = temp[0], temp[1], temp[2], temp[3], worker.dwn_date, temp[5]
            self.small[1][pos0] = f"Package : {a} ver{b}\nInfo : {c}\nRelease : {d}\nDownload : {e}\nPath : {f}"
            time.sleep(0.1)
            tkinter.messagebox.showinfo(title="Update Complete", message=f" Package {a} ver{b} is updated at {f}. ")
            self.dirctrl("./temp625/", False)
            self.dirctrl("./temp674/", False)

        elif x == 3: # help
            webbrowser.open( self.config_db[5] )
        self.current[1] = 0
        self.render(True, False)

    def custom1(self, x):
        if x == 2:
            if tkinter.messagebox.askokcancel(title="Update Current", message=" Are you sure to update whole status and get web data again? "):
                self.current, self.big, self.middle, self.small = [0, 0], [ ], [ ], [ ]
                self.setup()
                self.get_status()
                self.get_install()
                self.get_packages()
                self.get_update()

kobj.repath()
worker = mainclass()
time.sleep(0.5)
