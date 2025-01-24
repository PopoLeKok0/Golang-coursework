package main

import (
	"fmt"
	"strings"
)

type Person struct {
	lastName   string
	firstName  string
	birthDay   int
	birthMonth int
	birthYear  int
	ID         string
}

func generateID(p *Person) {

	lastNamePart := ""
	if len(p.lastName) >= 3 {
		lastNamePart = strings.Title(p.lastName[:3])
	} else {
		lastNamePart = strings.Title(p.lastName)
	}

	firstNamePart := ""
	if len(p.firstName) > 0 {
		firstNamePart = strings.ToUpper(p.firstName[:1])
	}

	birthDatePart := fmt.Sprintf("%04d%02d%02d", p.birthYear, p.birthMonth, p.birthDay)
	p.ID = fmt.Sprintf("%s%s%s", lastNamePart, firstNamePart, birthDatePart)
}

func main() {

	var lastName, firstName string
	var birthDay, birthMonth, birthYear int

	fmt.Print("Entrez le nom de famille : ")
	fmt.Scan(&lastName)

	fmt.Print("Entrez le prénom : ")
	fmt.Scan(&firstName)

	fmt.Print("Entrez le jour de naissance (1-31) : ")
	fmt.Scan(&birthDay)

	fmt.Print("Entrez le mois de naissance (1-12) : ")
	fmt.Scan(&birthMonth)

	fmt.Print("Entrez l'année de naissance : ")
	fmt.Scan(&birthYear)

	person := Person{
		lastName:   lastName,
		firstName:  firstName,
		birthDay:   birthDay,
		birthMonth: birthMonth,
		birthYear:  birthYear,
	}

	generateID(&person)

	fmt.Printf("\nPersonne créée :\n")
	fmt.Printf("Nom : %s\nPrénom : %s\nDate de naissance : %d/%d/%d\n", person.lastName, person.firstName, person.birthDay, person.birthMonth, person.birthYear)
	fmt.Printf("ID généré : %s\n", person.ID)
}
