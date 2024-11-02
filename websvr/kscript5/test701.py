# ./pyvm.exe program -sign -info -o

import kobj
import kvm

order = kobj.repath()
program = ""
sign = False
info = False
opt = False
for i in order[1:]:
    if i[0] == "-":
        if i.lower() == "-sign":
            sign = True
        elif i.lower() == "-info":
            info = True
        elif i.lower() == "-o":
            opt = True
    else:
        program = i

def f(p, s, i, o):
    k = kvm.kvm()
    ix = 0
    if program == "":
        return "err : no program"
    else:
        try:
            a, b, c = k.view(p)
            if i:
                return f"info : {a}\nabi : {b}\npublic : {c}"
            k.load(s)
        except Exception as e:
            return f"err : {e}"
    k.runone = not o
    while True:
        ix = k.run()
        if ix != 0:
            if ix < 8:
                return f"exit code : {ix} with {k.errmsg}"
            else:
                k.ma = kvm.testio(ix, k.callmem)

print( "\n" + f(program, sign, info, opt) )
input("press ENTER to exit... ")
