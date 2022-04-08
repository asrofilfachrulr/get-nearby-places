package main

import (
	"log"
	"os"

	m "github.com/asrofilfachrulr/get-nearby-places/models"
	"github.com/asrofilfachrulr/get-nearby-places/router"
)

func main() {
	// API Call and load them to predefined structs in models dir
	batchData := m.LoadAll()

	// Generate random places from given data
	bplace := m.GeneratePlaces(batchData)

	r := router.SetupRouter(bplace)

	PORT := os.Getenv("PORT")

	log.Fatalln(r.Run(PORT))
}
