package main

import (
	"log"

	m "github.com/asrofilfachrulr/get-nearby-places/models"
	"github.com/asrofilfachrulr/get-nearby-places/router"
)

func main() {
	// API Call and load them to predefined structs in models dir
	batchData := m.LoadAll()

	// Generate random places from given data
	batchPlaces := m.GeneratePlaces(batchData)

	// for i := 0; i < len(batchPlaces)/10; i++ {
	// 	b, _ := json.MarshalIndent(batchPlaces[i], "", " ")
	// 	fmt.Println(string(b))

	// }

	r := router.SetupRouter(batchPlaces)

	log.Fatalln(r.Run(":8000"))
}
