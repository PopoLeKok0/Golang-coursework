// Le programme ci-dessous permet de lire des nombres pairs envoyés sur un channel.
// ajouter dans la fonction main, la boucle de lecture au channel
// permettant d'afficher les nombres a la console

package main

import "fmt"

func main() {

	var st = []int{3, 6, 9, 15, 18, 34, 22, 5, 77, 99, 44}

	ch := sendPair(st)
	for {
		valeur, ok := <-ch
		if !ok {
			break
		}
		fmt.Printf("%d ", valeur)
	}
}

func sendPair(intArr []int) chan int {

	ch := make(chan int)

	go func() {

		for _, k := range intArr {

			if k%2 == 0 {
				ch <- k
			}
		}

		close(ch)
	}()

	return ch
}

// La sortie produite est:
// 6 18 34 22 44
