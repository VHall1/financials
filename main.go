package main

import (
	"context"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	c := NewMonzoClient(os.Getenv("MONZO_TOKEN"))
	transactions, err := c.ListTransactions(ctx, os.Getenv("MONZO_ACCOUNT_ID"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Fetched %d transactions", len(transactions))
	for _, t := range transactions {
		log.Printf("%+v", t)
	}
}
