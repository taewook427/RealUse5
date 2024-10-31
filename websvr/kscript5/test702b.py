import time

def int_calc():
    t = time.time()
    a, b, c, d = 0, 0, 0, 0
    for i in range(0, 10000000):
        a = c + b
        b = b + 17
        c = a + d
        d = d + 86
    return time.time() - t

def float_calc():
    t = time.time()
    a, b, c, d = 0.0, 0.0, 0.0, 0.0
    for i in range(0, 10000000):
        a = c + b
        b = b + 18.3617
        c = a + d
        d = d + 24.9281
    return time.time() - t

def fib(n):
    if n <= 2:
        return 1
    else:
        return fib(n - 2) + fib(n - 1)

def fib_calc(n):
    t = time.time()
    fib(n)
    return time.time() - t

t = int_calc()
print(40 / t, "M/s")
time.sleep(1.5)

t = float_calc()
print(40 / t, "M/s")
time.sleep(1.5)

t = fib_calc(38)
print(t, "s")
time.sleep(1.5)
