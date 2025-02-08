package main

import "fmt"

func main() {
	src := sendInt(10)
	dst := filterInt(src)

	for i := range dst {
		fmt.Printf("%d ", i)
	}
	fmt.Println("\nDone!")
}

func sendInt(maxNum int) chan int {
	defer fmt.Println("Sender ready!")
	ch := make(chan int)
	go func() {
		for i := 0; i < maxNum; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func filterInt(src <-chan int) chan int {
	defer fmt.Println("Filter ready!")
	dst := make(chan int)
	go func() {
		defer close(dst)
		for i := range src {
			dst <- i * 2
		}
	}()
	return dst
}
