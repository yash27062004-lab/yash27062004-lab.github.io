package main

import (
	"log"
	"net/http"
)

func main() {
	const address = "127.0.0.1:8080"

	files := http.FileServer(http.Dir("dist"))

	log.Printf("Website running at http://%s", address)

	if err := http.ListenAndServe(address, files); err != nil {
		log.Fatal(err)
	}
}
