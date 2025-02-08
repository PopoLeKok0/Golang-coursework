package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	var nGo int
	rand.Seed(42)

	fmt.Print("\nNumber of Go routines: ")
	fmt.Scanf("%d \n", &nGo)

	res := make([]chan float64, nGo)

	for i := 0; i < nGo; i++ {
		res[i] = numbers(1000)
	}

	for remaining := nGo; remaining > 0; {
		for i := 0; i < nGo; i++ {
			select {
			case num, open := <-res[i]:
				if open {
					fmt.Printf("Result %d: %f\n", i, num)
				} else {
					remaining--
				}
			default:
			}
		}
	}
}

func numbers(sz int) chan float64 {
	res := make(chan float64)
	go func() {
		defer close(res)
		num := 0.0

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)

		for i := 0; i < sz; i++ {
			num += math.Sqrt(math.Abs(rand.Float64()))
		}
		num /= float64(sz)

		res <- num
	}()
	return res
}
