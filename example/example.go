package main

import (
	"github.com/dcrichards/go-fhrs/fhrs"
	"log"
)

func main() {
	client, err := fhrs.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	
	establishment, err := client.Establishments.GetByID("82940")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n%+v\n", establishment)
}
