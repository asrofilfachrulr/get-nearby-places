package main

import (
	"encoding/json"
	"fmt"

	m "github.com/asrofilfachrulr/get-nearby-places/models"
)

func main() {
	batchData := m.LoadAll()
	batchPlaces := m.GeneratePlaces(batchData)

	for _, v := range [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		b, _ := json.MarshalIndent(batchPlaces[v], "", " ")
		fmt.Printf("b: %v\n", string(b))

	}

}
