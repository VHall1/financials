package main

import (
	"log"
	"os"
)

func main() {
	psql, err := NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	srv := NewHTTPServer(":80", psql, NewMonzoClient(os.Getenv("MONZO_TOKEN")))
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
