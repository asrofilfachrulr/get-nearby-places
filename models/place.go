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
*	  Also some places discovered (not did much reasearches) have really short
*	  distance to its border from its provided location (about < 1km)
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
		Name     string  // for debuggin purposes
		Places   []Place // all places in a city
	}
)

func GeneratePlaces(data BatchData) BatchPlace {
	// for formatting region name
	caser := cases.Title(language.Indonesian)

	// per item in slice contains of all places in a city with location of the city inserted to grouping them and ease for search nearby in future. Further info look GetNearbyPlaces function at line 227
	cplaces := []CityPlace{}

	// iterate per city
	for i := 0; i < len(data.Cities); i++ {
		// initiate CityPlace with current location (name attr for debugging purposes)
		cp := CityPlace{
			Location: Location{
				Latitude:  data.Cities[i].CoreInfo.Location.Latitude,
				Longitude: data.Cities[i].CoreInfo.Location.Longitude,
			},
			Name: data.Cities[i].CoreInfo.Name,
		}

		// container of all generated place (not caring about level) on a city
		places := []Place{}

		// id for every place in this city
		buildingId := 1

		// generate places per category due each category has different rule (counts)
		for _, category := range MapCategories["city"] {
			// id for place for current category
			categoryPlaceId := 1

			// generating by category place count
			for c := 0; c < int(category.Count); c++ {
				// creating new location around the "city location" by shifting diagonally randomly in such range that has been declared in constants on line 45-47
				l := RandShiftLoc(data.Cities[i].CoreInfo.Location, rand.Float64()*RADIUS_DISTRIBUTION_CITY)

				// region name with title case format
				regionName := caser.String(data.Cities[i].CoreInfo.Name)

				// the generated place name
				placeName := fmt.Sprintf("%s %s %d", category.Name, regionName, categoryPlaceId)

				p := Place{
					ID:         uint8(buildingId),
					CityName:   regionName,
					CategoryID: category.ID,
					Longitude:  l.Longitude,
					Latitude:   l.Latitude,
					Name:       placeName,
				}

				// append to current city slice of places
				places = append(places, p)

				// increment for next generated places
				categoryPlaceId += 1 // will be reset to 1 in every new category
				buildingId += 1      // incremented continously as a city's place "counter"
			}
		}

		// iterate through all districts on current city
		for j := 0; j < len(data.Cities[i].Districts); j++ {
			// id for every place in this district
			buildingId := 1

			// generate places per category due each category has different rule (counts)
			for _, category := range MapCategories["district"] {
				// id for place for current category
				categoryPlaceId := 1

				// generating by category place count
				for z := 0; z < int(category.Count); z++ {
					// creating new location around the "district location" by shifting diagonally randomly in such range that has been declared in constants on line 45-47
					l := RandShiftLoc(data.Cities[i].Districts[j].CoreInfo.Location, rand.Float64()*RADIUS_DISTRIBUTION_DISTRICT)

					// name of regions and generated place
					prefix := "Kecamatan "
					cityName := caser.String(data.Cities[i].CoreInfo.Name)
					districtName := caser.String(data.Cities[i].Districts[j].CoreInfo.Name)
					placeName := fmt.Sprintf("%s %s %d", category.Name, districtName, categoryPlaceId)

					p := Place{
						ID:           uint8(buildingId),
						CityName:     cityName,
						DistrictName: prefix + districtName,
						CategoryID:   category.ID,
						Longitude:    l.Longitude,
						Latitude:     l.Latitude,
						Name:         placeName,
					}

					// append to current city slice of places
					places = append(places, p)

					// increment for next generated places
					categoryPlaceId += 1 // will be reset to 1 in every new category
					buildingId += 1      // incremented continously as district's place "counter"
				}
			}

			// iterate through all village in current district in current city
			for k := 0; k < len(data.Cities[i].Districts[j].Villages); k++ {
				// id for every place in this village
				buildingId := 1

				// generate places per category due each category has different rule (counts)
				for _, category := range MapCategories["village"] {
					// id for place for current category
					categoryPlaceId := 1

					// generating by category place count
					for z := 0; z < int(category.Count); z++ {
						// creating new location around the "village location" by shifting diagonally randomly in such range that has been declared in constants on line 45-47
						l := RandShiftLoc(data.Cities[i].Districts[j].Villages[k].CoreInfo.Location, rand.Float64()*RADIUS_DISTRIBUTION_VILLAGE)

						// name of regions and generated place
						cityName := caser.String(data.Cities[i].CoreInfo.Name)

						districtName := fmt.Sprintf("%s %s", "Kecamatan", caser.String(data.Cities[i].Districts[j].CoreInfo.Name))

						vilName := data.Cities[i].Districts[j].Villages[k].CoreInfo.Name

						vilName = caser.String(vilName)

						placeName := fmt.Sprintf("%s %s %d", category.Name, vilName, categoryPlaceId)

						// SD not begin with level name (SD Cingcin not SD Desa Cingcin)
						if category.ID == 9 {
							level := caser.String(data.Cities[i].Districts[j].Villages[k].CoreInfo.Level)
							placeName = strings.Replace(placeName, level+" ", "", 1)
						}

						p := Place{
							ID:           uint8(buildingId),
							CityName:     cityName,
							DistrictName: districtName,
							VillageName:  vilName,
							CategoryID:   category.ID,
							Longitude:    l.Longitude,
							Latitude:     l.Latitude,
							Name:         placeName,
						}

						// append to current city slice of places
						places = append(places, p)

						// increment for next generated places
						categoryPlaceId += 1 // will be reset to 1 in every new category
						buildingId += 1      // incremented continously as district's place "counter"
					}
				}
			}
		}
		// assign all places in this city to Places of city place
		cp.Places = places
		// append to slice of city places
		cplaces = append(cplaces, cp)
	}

	// assign to super struct
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
					// set new max
					if p.Dist > max {
						maxIndex = i
						max = p.Dist
					}
				}
			}
		}
	}

	// debuggin show five closest cities against given location
	log.Println("5 closest cities/regeencies to given location:")
	// scan all places only in 5 closest cities
	for _, city := range closestCities {
		// print the city's name
		fmt.Println(city.CityPlace.Name)

		// iterate through all places in the city
		for _, place := range city.CityPlace.Places {
			_, distKm := haversine.Distance(pinned, haversine.Coord{
				Lat: place.Latitude,
				Lon: place.Longitude,
			})

			// only append if the distance is equal or less than 5 km
			if distKm <= 5 {
				// filtering place by category if exists
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
