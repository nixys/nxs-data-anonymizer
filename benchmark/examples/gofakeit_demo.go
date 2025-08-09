// Démonstration de gofakeit - À exécuter pour voir les données générées
package main

import (
	"fmt"
	"time"
	
	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	// Initialisation avec seed pour des résultats reproductibles
	gofakeit.Seed(time.Now().UnixNano())
	
	fmt.Println("=== DÉMONSTRATION GOFAKEIT ===")
	fmt.Println()
	
	// === NOMS ===
	fmt.Println("NOMS:")
	for i := 0; i < 5; i++ {
		fmt.Printf("  %s %s\n", gofakeit.FirstName(), gofakeit.LastName())
	}
	fmt.Printf("  Nom complet: %s\n", gofakeit.Name())
	fmt.Println()
	
	// === EMAILS ===
	fmt.Println("EMAILS:")
	for i := 0; i < 5; i++ {
		fmt.Printf("  %s\n", gofakeit.Email())
	}
	fmt.Println()
	
	// === TÉLÉPHONES ===
	fmt.Println("TÉLÉPHONES:")
	for i := 0; i < 5; i++ {
		fmt.Printf("  %s\n", gofakeit.Phone())
		fmt.Printf("  %s\n", gofakeit.PhoneFormatted())
	}
	// Téléphone français personnalisé
	fmt.Println("  Téléphones FR personnalisés:")
	for i := 0; i < 3; i++ {
		fmt.Printf("    +33%d%s\n", gofakeit.IntRange(6, 7), gofakeit.Numerify("########"))
	}
	fmt.Println()
	
	// === DATES ===
	fmt.Println("DATES:")
	fmt.Printf("  Date quelconque: %s\n", gofakeit.Date().Format("2006-01-02"))
	
	// Dates de naissance réalistes (18-80 ans)
	fmt.Println("  Dates de naissance:")
	for i := 0; i < 5; i++ {
		birthDate := gofakeit.DateRange(
			time.Now().AddDate(-80, 0, 0), 
			time.Now().AddDate(-18, 0, 0),
		)
		fmt.Printf("    %s\n", birthDate.Format("2006-01-02"))
	}
	fmt.Println()
	
	// === NOMBRES ===
	fmt.Println("NOMBRES:")
	fmt.Println("  Salaires:")
	for i := 0; i < 5; i++ {
		salary := gofakeit.Float64Range(25000, 75000)
		fmt.Printf("    %.2f€\n", salary)
	}
	fmt.Println()
	
	// === TEXTE ===
	fmt.Println("TEXTE:")
	fmt.Printf("  Phrase courte: %s\n", gofakeit.LoremIpsumSentence(6))
	fmt.Printf("  Phrase longue: %s\n", gofakeit.LoremIpsumSentence(15))
	fmt.Printf("  Paragraphe: %s\n", gofakeit.LoremIpsumParagraph(3, 8, 4, " "))
	fmt.Println()
	
	// === ADRESSES ===
	fmt.Println("ADRESSES:")
	for i := 0; i < 3; i++ {
		fmt.Printf("  %s\n", gofakeit.Address().Address)
		fmt.Printf("  %s, %s %s\n", gofakeit.Street(), gofakeit.Zip(), gofakeit.City())
	}
	fmt.Println()
	
	// === DONNÉES TECHNIQUES ===
	fmt.Println("DONNÉES TECHNIQUES:")
	fmt.Printf("  UUID: %s\n", gofakeit.UUID())
	fmt.Printf("  Username: %s\n", gofakeit.Username())
	fmt.Printf("  Password: %s\n", gofakeit.Password(true, true, true, true, false, 12))
	fmt.Println()
	
	// === COMPARAISON AVEC RANDOM ===
	fmt.Println("=== COMPARAISON RANDOM vs GOFAKEIT ===")
	fmt.Println()
	
	fmt.Println("MÉTHODE RANDOM ACTUELLE (comme Sprig natif):")
	for i := 0; i < 3; i++ {
		// Simulation de ce que fait randAlpha + randAlpha
		randomName := fmt.Sprintf("%s %s", 
			gofakeit.LetterN(8), 
			gofakeit.LetterN(10))
		randomEmail := fmt.Sprintf("%s@example.com", gofakeit.LetterN(8))
		randomPhone := fmt.Sprintf("+33%s", gofakeit.Numerify("#########"))
		
		fmt.Printf("  Nom: %s\n", randomName)
		fmt.Printf("  Email: %s\n", randomEmail) 
		fmt.Printf("  Tel: %s\n", randomPhone)
		fmt.Println()
	}
	
	fmt.Println("MÉTHODE GOFAKEIT (réaliste):")
	for i := 0; i < 3; i++ {
		fmt.Printf("  Nom: %s\n", gofakeit.Name())
		fmt.Printf("  Email: %s\n", gofakeit.Email())
		fmt.Printf("  Tel: %s\n", gofakeit.Phone())
		fmt.Println()
	}
}