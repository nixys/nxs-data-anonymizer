package main

import (
	"fmt"
	"os"
	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: faker-go [name|email|phone|date|salary|description]")
		os.Exit(1)
	}

	gofakeit.Seed(0)
	
	switch os.Args[1] {
	case "name":
		fmt.Print(gofakeit.Name())
	case "email":
		fmt.Print(gofakeit.Email())
	case "phone":
		fmt.Printf("+33%d", gofakeit.Number(600000000, 699999999))
	case "date":
		fmt.Print(gofakeit.Date().Format("2006-01-02"))
	case "salary":
		fmt.Printf("%.2f", gofakeit.Price(25000, 75000))
	case "description":
		fmt.Print(gofakeit.Sentence(15))
	default:
		fmt.Printf("Unknown type: %s\n", os.Args[1])
		os.Exit(1)
	}
}