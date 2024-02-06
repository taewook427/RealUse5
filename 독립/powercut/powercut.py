# test598 : powercut (PC5)
# 오전 3:25 ~ 3:35 랜덤

import time
import os

def delf(path):
    if os.path.exists(path):
        os.remove(path)

def mkf(path):
    f = open(path,'w')
    f.close()

path = os.path.join(os.path.expanduser('~'),'Desktop')
path = os.path.join(path,'alert598.txt')
delf(path)
time.sleep(900)
shutdown = True

while shutdown:
    time.sleep(250)
    local = time.localtime( time.time() )
    date = local.tm_mday
    hour = local.tm_hour
    minute = local.tm_min

    if date % 3 == 0:
        if hour == 3:
            if 19 <= minute <= 25:
                shutdown = False

mkf(path)
time.sleep(450)
delf(path)

os.system('shutdown -f -r -t 5')
