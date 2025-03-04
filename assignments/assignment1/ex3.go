package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

func randomNumbers() int {
	return rand.Intn(1000000000)
}

func isPrime(v int) bool {
	sq := int(math.Sqrt(float64(v))) + 1
	var i int
	for i = 2; i < sq; i++ {
		if v%i == 0 {
			return false
		}
	}
	return true
}

func intGenerator(wg *sync.WaitGroup, stop <-chan bool) <-chan int64 {
	intStream := make(chan int64)
	go func() {
		defer func() { wg.Done() }()
		defer close(intStream)
		defer fmt.Println("\nEnd of intGenerator...")

		var i int64
		for {
			select {
			case <-stop:
				return
			case intStream <- i:
				i++
			}
		}
	}()
	return intStream
}

func mersenneGenerator(wg *sync.WaitGroup, stop <-chan bool, input <-chan int64) <-chan int64 {
	output := make(chan int64)
	go func() {
		defer func() { wg.Done() }()
		defer close(output)
		defer fmt.Println("\nEnd of mersenneGenerator...")

		for n := range input {
			select {
			case <-stop:
				return
			case output <- int64(math.Pow(2, float64(n)) - 1):
			}
		}
	}()
	return output
}

func takeN(wg *sync.WaitGroup, stop <-chan bool, inputIntstream <-chan int64, n int) <-chan int64 {
	outputIntStream := make(chan int64)
	go func() {
		defer func() { wg.Done() }()
		defer close(outputIntStream)
		defer fmt.Println("\nEnd of takeN...")

		for i := 0; i < n; i++ {
			select {
			case <-stop:
				return
			case val, ok := <-inputIntstream:
				if !ok {
					return
				}
				outputIntStream <- val
			}
		}
	}()
	return outputIntStream
}

func main() {
	stop := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(1)
	intChannel := intGenerator(&wg, stop)

	wg.Add(1)
	mersenneChannel := mersenneGenerator(&wg, stop, intChannel)

	wg.Add(1)
	myChannel := takeN(&wg, stop, mersenneChannel, 15)

	for mersenne := range myChannel {
		fmt.Printf("%d ", mersenne)
	}

	close(stop)
	wg.Wait() // Wait for all goroutines to complete
}
