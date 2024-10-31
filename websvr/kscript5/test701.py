# ./pyvm.exe program -sign -info

import kobj
import kvm

order = kobj.repath()
program = ""
sign = False
info = False
for i in order[1:]:
    if i[0] == "-":
        if i.lower() == "-sign":
            sign = True
        elif i.lower() == "-info":
            info = True
    else:
        program = i

def f(p, s, i):
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
    while True:
        ix = k.run()
        if ix != 0:
            if ix < 8:
                return f"exit code : {ix} with {k.errmsg}"
            else:
                k.ma = kvm.testio(ix, k.callmem)

print( "\n" + f(program, sign, info) )
input("press ENTER to exit... ")
