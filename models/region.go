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

func MapAppend[R Regionable](dataMap *map[string][]R, k string, data R) {
	if _, ok := (*dataMap)[k]; !ok {
		(*dataMap)[k] = []R{data}
	} else {
		(*dataMap)[k] = append((*dataMap)[k], data)
	}
}

func RequestThenUnmarshall[R Unmarshallable](url string, data *R) error {
	resp, _ := http.Get(url)

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return json.Unmarshal(body, data)
}

func LoadAll() BatchData {
	// village
	var v DataVillageReq

	if err := RequestThenUnmarshall(VILLAGE_URL, &v); err != nil {
		log.Fatalln(err)
	}

	districtVillageMap := map[string][]Village{}

	for i := 0; i < len(v.Data); i++ {
		if err := validate.Struct(&v.Data[i]); err != nil {
			continue
		}
		codeString := strings.Split(v.Data[i].Code, ".")

		districtCode := fmt.Sprintf("%s.%s.%s", codeString[0], codeString[1], codeString[2])

		village := Village{
			CoreInfo: CoreInfo{
				Location: Location{
					Latitude:  v.Data[i].Latitude,
					Longitude: v.Data[i].Longitude,
				},
				Name:  strings.Join([]string{"KELURAHAN", v.Data[i].Name}, " "),
				Level: "KELURAHAN",
				Code:  v.Data[i].Code,
			},
		}
		MapAppend(&districtVillageMap, districtCode, village)
	}

	// district
	var d DataDistrictReq
	if err := RequestThenUnmarshall(DISTRICT_URL, &d); err != nil {
		log.Fatalln(err)
	}

	cityDistrictMap := map[string][]District{}
	for i := 0; i < len(d.Data); i++ {
		if err := validate.Struct(&d.Data[i]); err != nil {
			continue
		}

		codeString := strings.Split(d.Data[i].Code, ".")
		cityCode := fmt.Sprintf("%s.%s", codeString[0], codeString[1])
		districtCode := fmt.Sprintf("%s.%s.%s", codeString[0], codeString[1], codeString[2])

		district := District{
			CoreInfo: CoreInfo{
				Location: Location{
					Latitude:  d.Data[i].Latitude,
					Longitude: d.Data[i].Longitude,
				},
				Level: "KECAMATAN",
				Name:  d.Data[i].Name,
				Code:  d.Data[i].Code,
			},
			Villages: districtVillageMap[districtCode],
		}

		MapAppend(&cityDistrictMap, cityCode, district)
	}

	// city
	var c DataCityReq

	if err := RequestThenUnmarshall(CITY_URL, &c); err != nil {
		log.Fatalln(err)
	}

	cities := []City{}
	for i := 0; i < len(c.Data); i++ {
		if err := validate.Struct(&c.Data[i]); err != nil {
			continue
		}
		cityCode := fmt.Sprintf("%.2f", c.Data[i].Code)

		cityLevel := strings.Split(c.Data[i].Name, " ")[0]

		if cityLevel == "KAB." {
			cityLevel = "KABUPATEN"
			vilLevel := "DESA"

			districts := cityDistrictMap[cityCode]

			for i := 0; i < len(districts); i++ {
				for j := 0; j < len(districts[i].Villages); j++ {
					districts[i].Villages[j].CoreInfo.Level = vilLevel
					districts[i].Villages[j].CoreInfo.Name = strings.Replace(districts[i].Villages[j].CoreInfo.Name, "KELURAHAN", vilLevel, 1)
				}
			}
		}

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
			Districts: cityDistrictMap[cityCode],
		})

	}

	return BatchData{
		Cities: cities,
	}
}
