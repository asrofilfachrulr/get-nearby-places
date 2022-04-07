package models

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/umahmood/haversine"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

type WebQuery struct {
	Latitude   float64
	Longitude  float64
	CategoryId uint8
}

func GeneratePlaces(data BatchData) []Place {
	batchPlaces := []Place{}
	caser := cases.Title(language.Indonesian)

	for i := 0; i < len(data.Cities); i++ {

		buildingId := 1
		for _, category := range MapCategories["city"] {
			placeId := 1
			for c := 0; c < int(category.Count); c++ {

				l := RandShiftLoc(data.Cities[i].CoreInfo.Location, float64(rand.Intn(5000)))

				regionName := caser.String(data.Cities[i].CoreInfo.Name)

				placeName := fmt.Sprintf("%s %s %d", category.Name, caser.String(regionName), placeId)

				batchPlaces = append(batchPlaces, Place{
					ID:         uint8(buildingId),
					CityName:   caser.String(data.Cities[i].CoreInfo.Name),
					CategoryID: category.ID,
					Longitude:  l.Longitude,
					Latitude:   l.Latitude,
					Name:       placeName,
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
					l := RandShiftLoc(data.Cities[i].CoreInfo.Location, float64(rand.Intn(3000)))

					batchPlaces = append(batchPlaces, Place{
						ID:           uint8(buildingId),
						CityName:     caser.String(data.Cities[i].CoreInfo.Name),
						DistrictName: "Kecamatan " + caser.String(data.Cities[i].Districts[j].CoreInfo.Name),
						CategoryID:   category.ID,
						Longitude:    l.Longitude,
						Latitude:     l.Latitude,
						Name:         fmt.Sprintf("%s %s %d", category.Name, caser.String(data.Cities[i].Districts[j].CoreInfo.Name), placeId),
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
						l := RandShiftLoc(data.Cities[i].CoreInfo.Location, float64(rand.Intn(1500)))

						districtName := fmt.Sprintf("%s %s", "Kecamatan", caser.String(data.Cities[i].Districts[j].CoreInfo.Name))

						regionName := data.Cities[i].Districts[j].Villages[k].CoreInfo.Name
						vilName := caser.String(regionName)

						if category.ID != 3 {
							regionName = strings.Replace(regionName, data.Cities[i].Districts[j].Villages[k].CoreInfo.Level+" ", "", 1)
						}
						placeName := fmt.Sprintf("%s %s %d", category.Name, caser.String(regionName), placeId)

						batchPlaces = append(batchPlaces, Place{
							ID:           uint8(buildingId),
							CityName:     caser.String(data.Cities[i].CoreInfo.Name),
							DistrictName: districtName,
							VillageName:  vilName,
							CategoryID:   category.ID,
							Longitude:    l.Longitude,
							Latitude:     l.Latitude,
							Name:         placeName,
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

func GetNearbyPlaces(q WebQuery, BatchPlaces []Place) ([]Place, error) {
	places := []Place{}

	pinned := haversine.Coord{Lat: q.Latitude, Lon: q.Longitude}

	for _, place := range BatchPlaces {
		_, distKm := haversine.Distance(pinned, haversine.Coord{
			Lat: place.Latitude,
			Lon: place.Longitude,
		})

		if distKm <= 5 {
			if q.CategoryId != 0 {
				if place.CategoryID == q.CategoryId {
					places = append(places, place)
				}
			} else {
				places = append(places, place)
			}
		}
	}

	return places, nil
}
