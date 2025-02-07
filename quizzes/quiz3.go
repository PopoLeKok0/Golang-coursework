// Le programme ci-dessous démontre l'utilisation d'une boucle parallèle. Toutefois ce programme devrait attendre la fin de toutes les go routines avant de se terminer. Ajouter le mécanisme de synchronisation sync.WaitGroup afin que ce programme fonctionne correctement.
package main

import (
	"fmt"
	"sync"
)

func process(x, y int) int {
	r := x + y
	return r
}

func main() {
	arr1 := []int{40, 15, 22, 32}
	arr2 := []int{14, 21, 30, 44}
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			fmt.Println(arr1[idx], "+", arr2[idx])
			arr2[idx] = process(arr1[idx], arr2[idx])
		}(i)
	}

	wg.Wait()
	fmt.Println("arr1=", arr1)
	fmt.Println("arr2=", arr2)
}
