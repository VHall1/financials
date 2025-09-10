package main

import (
	"log"
	"os"
)

func main() {
	srv := NewHTTPServer(":8080", NewMonzoClient(os.Getenv("MONZO_TOKEN")))
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
