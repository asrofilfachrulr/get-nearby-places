package models

import (
	"fmt"
	"math/rand"
	"strings"
)

type Place struct {
	ID           uint8   `json:"id"` // refers to building ownership in region
	CityName     string  `json:"city_name"`
	DistrictName string  `json:"district_name,omitempty"`
	VillageName  string  `json:"village_name,omitempty"`
	CategoryID   uint8   `json:"category_id"`
	Name         string  `json:"name"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
}

var directions = []Direction{North, South, West, East}

func GeneratePlaces(data BatchData) []Place {
	batchPlaces := []Place{}

	for i := 0; i < len(data.Cities); i++ {

		buildingId := 1
		for _, category := range MapCategories["city"] {
			placeId := 1
			for c := 0; c < int(category.Count); c++ {

				l, _ := ShiftLocation(data.Cities[i].CoreInfo.Location, float64(rand.Intn(5000)), directions[rand.Intn(4)])

				batchPlaces = append(batchPlaces, Place{
					ID:         uint8(buildingId),
					CityName:   strings.Title(strings.ToLower(data.Cities[i].CoreInfo.Name)),
					CategoryID: category.ID,
					Longitude:  l.Longitude,
					Latitude:   l.Latitude,
					Name:       fmt.Sprintf("%s %s %d", category.Name, strings.Title(strings.ToLower(data.Cities[i].CoreInfo.Name)), placeId),
				})

				placeId += 1
				buildingId += 1
			}
		}

		for j := 0; j < len(data.Cities[i].Districts); j++ {
			buildingId := 1
			for _, category := range MapCategories["district"] {
				placeId := 1
				for z := 0; z < int(category.Count); z++ {
					l, _ := ShiftLocation(data.Cities[i].CoreInfo.Location, float64(rand.Intn(5000)), directions[rand.Intn(4)])

					batchPlaces = append(batchPlaces, Place{
						ID:           uint8(buildingId),
						CityName:     strings.Title(strings.ToLower(data.Cities[i].CoreInfo.Name)),
						DistrictName: "Kecamatan " + strings.Title(strings.ToLower(data.Cities[i].Districts[j].CoreInfo.Name)),
						CategoryID:   category.ID,
						Longitude:    l.Longitude,
						Latitude:     l.Latitude,
						Name:         fmt.Sprintf("%s %s %d", category.Name, strings.Title(strings.ToLower(data.Cities[i].Districts[j].CoreInfo.Name)), placeId),
					})
					placeId += 1
					buildingId += 1
				}
			}

			for k := 0; k < len(data.Cities[i].Districts[j].Villages); k++ {
				buildingId := 1
				for _, category := range MapCategories["village"] {
					placeId := 1
					for z := 0; z < int(category.Count); z++ {

						l, _ := ShiftLocation(data.Cities[i].CoreInfo.Location, float64(rand.Intn(5000)), directions[rand.Intn(4)])

						batchPlaces = append(batchPlaces, Place{
							ID:           uint8(buildingId),
							CityName:     strings.Title(strings.ToLower(data.Cities[i].CoreInfo.Name)),
							DistrictName: "Kecamatan " + strings.Title(strings.ToLower(data.Cities[i].Districts[j].CoreInfo.Name)),
							VillageName:  strings.Title(strings.ToLower(data.Cities[i].Districts[j].Villages[k].CoreInfo.Name)),
							CategoryID:   category.ID,
							Longitude:    l.Longitude,
							Latitude:     l.Latitude,
							Name:         fmt.Sprintf("%s %s %d", category.Name, strings.Title(strings.ToLower(data.Cities[i].Districts[j].Villages[k].CoreInfo.Name)), placeId),
						})
						placeId += 1
						buildingId += 1
					}
				}
			}
		}
	}

	return batchPlaces
}
