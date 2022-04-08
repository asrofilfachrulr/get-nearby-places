package models

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/umahmood/haversine"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Struct for response [GET] /search
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

// Struct for request [GET] /search
type WebQuery struct {
	Latitude   float64
	Longitude  float64
	CategoryId uint8
}

/*
*	Assumptions about distribute random places:
*	- Due unknown border locations, assume below constants
*	  are a safe radius for distributing random places which
*	  those places are still belong to such region we referred to,
*	  Eventhough in real life, the provided location isn't centralized
*	  against the shape as overall and has various distance to its border
*	  Also some places discovered (not did much reasearch) have really short
*	  distance to its border from its provided location (< 1km)
 */
const (
	// in meter
	RADIUS_DISTRIBUTION_CITY     float64 = 5000
	RADIUS_DISTRIBUTION_DISTRICT float64 = 3000
	RADIUS_DISTRIBUTION_VILLAGE  float64 = 1500
)

// for faster used algorithm (linear search), clustering places to theirs upper region level location
// (village and district embed to their city)

type (
	BatchPlace struct {
		CityPlaces []CityPlace
	}
	CityPlace struct {
		Location Location
		Name     string
		Places   []Place // all places
	}
)

func GeneratePlaces(data BatchData) BatchPlace {
	// bplaces := BatchPlace{}

	// batchPlaces := []Place{}
	caser := cases.Title(language.Indonesian)

	cplaces := []CityPlace{}
	for i := 0; i < len(data.Cities); i++ {
		cp := CityPlace{
			Location: Location{
				Latitude:  data.Cities[i].CoreInfo.Location.Latitude,
				Longitude: data.Cities[i].CoreInfo.Location.Longitude,
			},
			Name: data.Cities[i].CoreInfo.Name,
		}
		places := []Place{}

		buildingId := 1
		for _, category := range MapCategories["city"] {
			placeId := 1
			for c := 0; c < int(category.Count); c++ {

				l := RandShiftLoc(data.Cities[i].CoreInfo.Location, rand.Float64()*RADIUS_DISTRIBUTION_CITY)

				regionName := caser.String(data.Cities[i].CoreInfo.Name)

				placeName := fmt.Sprintf("%s %s %d", category.Name, caser.String(regionName), placeId)

				// batchPlaces = append(batchPlaces)

				p := Place{
					ID:         uint8(buildingId),
					CityName:   caser.String(data.Cities[i].CoreInfo.Name),
					CategoryID: category.ID,
					Longitude:  l.Longitude,
					Latitude:   l.Latitude,
					Name:       placeName,
				}

				// batchPlaces = append(batchPlaces, p)

				places = append(places, p)

				placeId += 1
				buildingId += 1
			}
		}

		for j := 0; j < len(data.Cities[i].Districts); j++ {

			buildingId := 1
			for _, category := range MapCategories["district"] {
				placeId := 1
				for z := 0; z < int(category.Count); z++ {
					l := RandShiftLoc(data.Cities[i].CoreInfo.Location, rand.Float64()*RADIUS_DISTRIBUTION_DISTRICT)

					p := Place{
						ID:           uint8(buildingId),
						CityName:     caser.String(data.Cities[i].CoreInfo.Name),
						DistrictName: "Kecamatan " + caser.String(data.Cities[i].Districts[j].CoreInfo.Name),
						CategoryID:   category.ID,
						Longitude:    l.Longitude,
						Latitude:     l.Latitude,
						Name:         fmt.Sprintf("%s %s %d", category.Name, caser.String(data.Cities[i].Districts[j].CoreInfo.Name), placeId),
					}

					// batchPlaces = append(batchPlaces, p)
					places = append(places, p)

					placeId += 1
					buildingId += 1
				}
			}

			for k := 0; k < len(data.Cities[i].Districts[j].Villages); k++ {

				buildingId := 1
				for _, category := range MapCategories["village"] {
					placeId := 1
					for z := 0; z < int(category.Count); z++ {
						l := RandShiftLoc(data.Cities[i].Districts[j].Villages[k].CoreInfo.Location, rand.Float64()*RADIUS_DISTRIBUTION_VILLAGE)

						districtName := fmt.Sprintf("%s %s", "Kecamatan", caser.String(data.Cities[i].Districts[j].CoreInfo.Name))

						regionName := data.Cities[i].Districts[j].Villages[k].CoreInfo.Name
						vilName := caser.String(regionName)

						if category.ID != 3 {
							regionName = strings.Replace(regionName, data.Cities[i].Districts[j].Villages[k].CoreInfo.Level+" ", "", 1)
						}
						placeName := fmt.Sprintf("%s %s %d", category.Name, caser.String(regionName), placeId)

						p := Place{
							ID:           uint8(buildingId),
							CityName:     caser.String(data.Cities[i].CoreInfo.Name),
							DistrictName: districtName,
							VillageName:  vilName,
							CategoryID:   category.ID,
							Longitude:    l.Longitude,
							Latitude:     l.Latitude,
							Name:         placeName,
						}

						// batchPlaces = append(batchPlaces, p)
						places = append(places, p)

						placeId += 1
						buildingId += 1
					}
				}
			}
		}
		cp.Places = places
		cplaces = append(cplaces, cp)
	}

	bp := BatchPlace{
		CityPlaces: cplaces,
	}
	return bp
}

// assumption: the pinned location cant be surrounded by more than five cities.
// only scan places in closest cities.
func GetNearbyPlaces(q WebQuery, bp BatchPlace) ([]Place, error) {
	places := []Place{}

	pinned := haversine.Coord{Lat: q.Latitude, Lon: q.Longitude}

	// find 5 closest city to the pinned
	type CpDistance struct {
		CityPlace CityPlace
		Dist      float64
	}
	closestCities := []CpDistance{}

	// for fast replacement
	maxIndex := 0
	max := 0.0

	// iterate to all cities
	for _, cplace := range bp.CityPlaces {
		// find the distance from the city to the pinned
		_, distKm := haversine.Distance(pinned, haversine.Coord{
			Lat: cplace.Location.Latitude,
			Lon: cplace.Location.Longitude,
		})

		// directly append if slice isnt full yet
		if len(closestCities) < 5 {
			closestCities = append(closestCities, CpDistance{
				CityPlace: cplace,
				Dist:      distKm,
			})

			// set for future fast compare current distance against the max distance has been registered
			if max != 0.0 {
				if distKm > max {
					max = distKm
					maxIndex = len(closestCities) - 1
				}
			} else {
				max = distKm
			}
		} else {
			// find the max dist registered , if the current city closer, the one in slice will be replaced
			if distKm < max {
				closestCities[maxIndex] = CpDistance{
					CityPlace: cplace,
					Dist:      distKm,
				}

				max = distKm

				// looking for new max and max index
				for i, p := range closestCities {
					if maxIndex == i {
						continue
					}
					if p.Dist > max {
						maxIndex = i
						max = p.Dist
					}
				}
			}
		}
	}

	log.Println("5 closest cities/regeencies to given location:")
	// scan all places only in 5 closest cities
	for _, city := range closestCities {
		fmt.Println(city.CityPlace.Name)
		for _, place := range city.CityPlace.Places {
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
	}
	fmt.Println()

	return places, nil
}
