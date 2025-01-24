package main

import (
	"fmt"
	"math"
)

func FloorAndCeil(value float32) (int, int) {
	floorValue := int(math.Floor(float64(value)))
	ceilValue := int(math.Ceil(float64(value)))
	return floorValue, ceilValue
}

func main() {

	var number float32 = 2.5
	floor, ceil := FloorAndCeil(number)

	fmt.Printf("Pour le nombre %f :\n", number)
	fmt.Printf("Arrondissement inférieur (floor) : %d\n", floor)
	fmt.Printf("Arrondissement supérieur (ceil) : %d\n", ceil)

	number = -1.5
	floor, ceil = FloorAndCeil(number)

	fmt.Printf("\nPour le nombre %f :\n", number)
	fmt.Printf("Arrondissement inférieur (floor) : %d\n", floor)
	fmt.Printf("Arrondissement supérieur (ceil) : %d\n", ceil)
}
