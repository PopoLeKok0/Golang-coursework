// Le programme suivant est corrigé.

package main

import "fmt"

// Boite
type Boite struct {
	poids   float64
	couleur string
}

// methodes
func (p *Boite) SetCouleur(newCouleur string) {
	p.couleur = newCouleur
}

func (p Boite) GetCouleur() string {
	return p.couleur
}

func (p Boite) GetPoids() float64 {
	return p.poids
}

func (p *Boite) doublePoids() {
	p.poids *= 2
}

func main() {

	var b = Boite{32.4, "jaune"}

	// on double le poids de la boite
	b.doublePoids()

	// on veut imprimer la couleur de la boite
	fmt.Printf("La couleur est: %s\n", b.GetCouleur())
	fmt.Printf("Le poids est: %f", b.GetPoids())
}

// La sortie produite est:
// La couleur est: jaune
// Le poids est: 64.800000
