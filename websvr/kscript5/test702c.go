package main

import (
	"fmt"
	"time"
)

func int_calc(v *float64) float64 {
	t := time.Now()
	var a, b, c, d int
	for i := 0; i < 10000000; i++ {
		a = c + b
		b = b + 17
		c = a + d
		d = d + 86
	}
	*v = float64(a + b + c + d)
	return float64(time.Since(t).Microseconds()) / 1000000
}

func float_calc(v *float64) float64 {
	t := time.Now()
	var a, b, c, d float64
	for i := 0; i < 10000000; i++ {
		a = c + b
		b = b + 18.3617
		c = a + d
		d = d + 24.9287
	}
	*v = a + b + c + d
	return float64(time.Since(t).Microseconds()) / 1000000
}

func fib(n int) int {
	if n <= 2 {
		return 1
	} else {
		return fib(n-2) + fib(n-1)
	}
}

func fib_calc(n int) float64 {
	t := time.Now()
	_ = fib(n)
	return float64(time.Since(t).Microseconds()) / 1000000
}

func main() {
	var count float64
	fmt.Printf("%f M/s\n", 40/int_calc(&count))
	fmt.Printf("%f M/s\n", 40/float_calc(&count))
	fmt.Printf("%f s\n", fib_calc(38))
}
