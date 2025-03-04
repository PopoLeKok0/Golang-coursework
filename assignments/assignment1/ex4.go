package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// checks if it is a prime number
func isPrime(v int64) bool {
	sq := int64(math.Sqrt(float64(v))) + 1
	var i int64
	for i = 2; i < sq; i++ {
		if v%i == 0 {
			return false
		}
	}
	return true
}

// returns a prime number
func getPrime(maxValue int64) int64 {
	for {
		n := rand.Int63n(maxValue)
		if isPrime(n) {
			return n
		}
	}
}

// a special prime is a prime number that ends with the pattern sequence
// after nTrials the function returns with a false error code
func getSpecialPrime(pattern int64, maxValue int64, nTrials int) (int64, bool) {
	var div int64
	for div = 10; pattern/div != 0; div *= 10 {
	}

	for i := 0; i < nTrials; i++ {
		n := getPrime(maxValue)
		if n%div == pattern {
			return n, true // special prime found
		}
	}

	return 0, false // we failed to find a special prime
}

// creates a goroutine that continuously searches for special primes
// and sends them through the returned channel
func getSpecialPrimeStream(wg *sync.WaitGroup, stop <-chan struct{}, pattern int64, maxValue int64, nTrials int) <-chan int64 {
	ch := make(chan int64, 1) // Add buffer to prevent blocking

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch)

		for {
			select {
			case <-stop:
				return
			default:
				n, ok := getSpecialPrime(pattern, maxValue, nTrials)
				if ok {
					select {
					case ch <- n:
						// Successfully sent the prime
					case <-stop:
						return
					}
				}
			}
		}
	}()

	return ch
}

// This is the version with threads
func main() {
	fmt.Println("Solution with one channel per thread.")
	var sp []int64
	var pattern int64 = 1111       // the suffix pattern we are looking for
	var maxV int64 = 1000000000000 // maximum value for the prime number
	var trials int = 3             // number of trials for each attempts
	var nPrimes int = 4            // number of special primes we want

	var wg sync.WaitGroup
	stop := make(chan struct{})

	start := time.Now()

	// creates the threads and returns the channel
	ch1 := getSpecialPrimeStream(&wg, stop, pattern, maxV, trials)
	ch2 := getSpecialPrimeStream(&wg, stop, pattern, maxV, trials)
	ch3 := getSpecialPrimeStream(&wg, stop, pattern, maxV, trials)
	ch4 := getSpecialPrimeStream(&wg, stop, pattern, maxV, trials)
	ch5 := getSpecialPrimeStream(&wg, stop, pattern, maxV, trials)
	ch6 := getSpecialPrimeStream(&wg, stop, pattern, maxV, trials)
	ch7 := getSpecialPrimeStream(&wg, stop, pattern, maxV, trials)
	ch8 := getSpecialPrimeStream(&wg, stop, pattern, maxV, trials)

	// monitor channels with select statement
	for len(sp) < nPrimes {
		select {
		case n := <-ch1:
			if n != 0 {
				sp = append(sp, n)
			}
		case n := <-ch2:
			if n != 0 {
				sp = append(sp, n)
			}
		case n := <-ch3:
			if n != 0 {
				sp = append(sp, n)
			}
		case n := <-ch4:
			if n != 0 {
				sp = append(sp, n)
			}
		case n := <-ch5:
			if n != 0 {
				sp = append(sp, n)
			}
		case n := <-ch6:
			if n != 0 {
				sp = append(sp, n)
			}
		case n := <-ch7:
			if n != 0 {
				sp = append(sp, n)
			}
		case n := <-ch8:
			if n != 0 {
				sp = append(sp, n)
			}
		}
	}

	// stop all goroutines
	close(stop)

	// Wait for all goroutines to finish
	wg.Wait()

	end := time.Now()
	fmt.Println("Special prime numbers are: ", sp)
	fmt.Println("End of program.", end.Sub(start))
}
