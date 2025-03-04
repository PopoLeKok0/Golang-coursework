package main

import (
	"fmt"
	"sync"
)

// greatest common divisor
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	x := []int{100, 15, 18, 10, 25, 90, 22, 60}
	y := []int{60, 5, 72, 5, 250, 108, 44, 3600}
	var z [8]int

	var wg sync.WaitGroup
	wg.Add(len(z))

	// parallel loop
	for i := range z {
		go func(i int) {
			defer wg.Done()
			z[i] = gcd(x[i], y[i])
		}(i)
	}

	wg.Wait()

	// print result
	for _, v := range z {
		fmt.Println(v)
	}
}
