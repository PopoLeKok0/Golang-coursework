package main

import (
	"fmt"
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

	done := make(chan bool, len(z))

	// Parallel loop
	for i := range z {
		go func(i int) {
			z[i] = gcd(x[i], y[i])
			done <- true
		}(i)
	}

	for i := 0; i < len(z); i++ {
		<-done
	}

	// Print result
	for _, v := range z {
		fmt.Println(v)
	}
}
