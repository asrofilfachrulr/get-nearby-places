package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

const (
	VILLAGE_URL  string = "https://satudata.jabarprov.go.id/api-backend/bigdata/diskominfo/od_kode_wilayah_dan_nama_wilayah_desa_kelurahan?limit=5957"
	DISTRICT_URL string = "https://satudata.jabarprov.go.id/api-backend/bigdata/diskominfo/od_16357_kode_wilayah_dan_nama_wilayah_kecamatan?limit=627"
	CITY_URL     string = "https://satudata.jabarprov.go.id/api-backend/bigdata/diskominfo/od_kode_wilayah_dan_nama_wilayah_kota_kabupaten?limit=27"
)

// common attributes / structures
type (
	Location struct {
		Latitude  float64
		Longitude float64
	}
	CoreInfo struct {
		Location Location
		Level    string
		Name     string
		Code     string
	}
)

// structuring unmarshalled json from API call
type (
	Village struct {
		CoreInfo CoreInfo
	}
	District struct {
		CoreInfo CoreInfo
		Villages []Village
	}
	City struct {
		CoreInfo  CoreInfo
		Districts []District
	}
	BatchData struct {
		Cities    []City
		NCity     uint
		NDistrict uint
		NVillage  uint
	}
)

var validate = validator.New()

// The code using golang feature: generic. Please use minimum golang version 1.18

// Using generic, check if key exists of a region structs, append to value which is a slice
// If doesn't extist, initialize the new one
func MapAppend[R Regionable](dataMap *map[string][]R, k string, data R) {
	if _, ok := (*dataMap)[k]; !ok {
		(*dataMap)[k] = []R{data}
	} else {
		(*dataMap)[k] = append((*dataMap)[k], data)
	}
}

// Using generic, request to the given URL and unmarshalling to data which R type.
func RequestThenUnmarshall[R Unmarshallable](url string, data *R) error {
	resp, _ := http.Get(url)

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return json.Unmarshal(body, data)
}

// Load all data and seed to defined structs
func LoadAll() BatchData {
	//---------------village data---------------//
	var v DataVillageReq

	if err := RequestThenUnmarshall(VILLAGE_URL, &v); err != nil {
		log.Fatalln(err)
	}

	// relation between district and its villages, mapped by district code
	districtVillageMap := map[string][]Village{}

	for i := 0; i < len(v.Data); i++ {
		if err := validate.Struct(&v.Data[i]); err != nil {
			continue
		}

		// split village code
		codeSegments := strings.Split(v.Data[i].Code, ".")

		districtCode := fmt.Sprintf("%s.%s.%s", codeSegments[0], codeSegments[1], codeSegments[2])

		// level name determined by the begin number of the code last segment. [1] for KELURAHAN and [2] for DESA
		level := "KELURAHAN"
		if string(codeSegments[3][0]) == "2" {
			level = "DESA"
		}

		village := Village{
			CoreInfo: CoreInfo{
				Location: Location{
					Latitude:  v.Data[i].Latitude,
					Longitude: v.Data[i].Longitude,
				},
				Name:  strings.Join([]string{level, v.Data[i].Name}, " "),
				Level: level,
				Code:  v.Data[i].Code,
			},
		}
		MapAppend(&districtVillageMap, districtCode, village)
	}

	//---------------district data---------------//
	var d DataDistrictReq

	if err := RequestThenUnmarshall(DISTRICT_URL, &d); err != nil {
		log.Fatalln(err)
	}

	// relation between city and all its districts, mapped by city code
	cityDistrictMap := map[string][]District{}

	for i := 0; i < len(d.Data); i++ {
		if err := validate.Struct(&d.Data[i]); err != nil {
			continue
		}

		// split district code to set city code and district code
		codeSegments := strings.Split(d.Data[i].Code, ".")

		cityCode := fmt.Sprintf("%s.%s", codeSegments[0], codeSegments[1])
		districtCode := fmt.Sprintf("%s.%s.%s", codeSegments[0], codeSegments[1], codeSegments[2])

		// level name, true for all data. not used as prefix to name a district
		level := "KECAMATAN"

		district := District{
			CoreInfo: CoreInfo{
				Location: Location{
					Latitude:  d.Data[i].Latitude,
					Longitude: d.Data[i].Longitude,
				},
				Level: level,
				Name:  d.Data[i].Name,
				Code:  d.Data[i].Code,
			},
			// relate with slice of villages which has same district code
			Villages: districtVillageMap[districtCode],
		}

		MapAppend(&cityDistrictMap, cityCode, district)
	}

	//--------------- city data ---------------//
	var c DataCityReq

	if err := RequestThenUnmarshall(CITY_URL, &c); err != nil {
		log.Fatalln(err)
	}

	cities := []City{}

	for i := 0; i < len(c.Data); i++ {
		if err := validate.Struct(&c.Data[i]); err != nil {
			continue
		}

		// due city code from API has float type (uniquely against other codes), format to string to achieve data uniformity
		cityCode := fmt.Sprintf("%.2f", c.Data[i].Code)

		// get city level name from its name
		cityLevel := strings.Split(c.Data[i].Name, " ")[0]

		if cityLevel == "KAB." {
			// remove the abbreviation for city level
			cityLevel = "KABUPATEN"
		}

		// remove the abbreviation on city name for city name
		cityName := strings.Replace(c.Data[i].Name, "KAB.", "KABUPATEN", 1)

		cities = append(cities, City{
			CoreInfo: CoreInfo{
				Location: Location{
					Longitude: c.Data[i].Longitude,
					Latitude:  c.Data[i].Latitude,
				},
				Level: cityLevel,
				Name:  cityName,
				Code:  cityCode,
			},
			// relate with slice of districts which has same city code
			Districts: cityDistrictMap[cityCode],
		})

	}

	// insert cities to super struct
	return BatchData{
		Cities: cities,
	}
}
