package main

import "fmt"

func FilterPositives(numbers []int) []int {

	result := make([]int, 0, cap(numbers))
	for _, num := range numbers {
		if num >= 0 {
			result = append(result, num)
		}
	}

	return result
}

func main() {
	numbers := []int{-10, 15, -3, 7, 0, -5, 12}
	fmt.Printf("Slice original : %v (capacité : %d)\n", numbers, cap(numbers))

	filtered := FilterPositives(numbers)
	fmt.Printf("Slice filtré : %v (capacité : %d)\n", filtered, cap(filtered))
}
