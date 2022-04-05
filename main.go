package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/asrofilfachrulr/get-nearby-places/models"
	"github.com/go-playground/validator/v10"
)

func main() {
	// village
	var v models.DataVillageReq
	resp, _ := http.Get("https://satudata.jabarprov.go.id/api-backend/bigdata/diskominfo/od_kode_wilayah_dan_nama_wilayah_desa_kelurahan?limit=5957")

	validate := validator.New()

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &v)

	districtVillageMap := map[string][]models.Village{}

	for i := 0; i < len(v.Data); i++ {
		if err := validate.Struct(&v.Data[i]); err != nil {
			continue
		}
		kodeString := strings.Split(v.Data[i].Code, ".")
		kodeKecamatan := fmt.Sprintf("%s.%s.%s", kodeString[0], kodeString[1], kodeString[2])
		village := models.Village{
			CoreInfo: models.CoreInfo{
				Location: models.Location{
					Latitude:  v.Data[i].Latitude,
					Longitude: v.Data[i].Longitude,
				},
				Name:  v.Data[i].Nama,
				Level: "KELURAHAN",
				Code:  v.Data[i].Code,
			},
		}

		if _, ok := districtVillageMap[kodeKecamatan]; !ok {
			districtVillageMap[kodeKecamatan] = []models.Village{village}
		} else {
			districtVillageMap[kodeKecamatan] = append(districtVillageMap[kodeKecamatan], village)
		}
	}

	// district
	var d models.DataDistrictReq
	resp, _ = http.Get("https://satudata.jabarprov.go.id/api-backend/bigdata/diskominfo/od_16357_kode_wilayah_dan_nama_wilayah_kecamatan?limit=627")

	//We Read the response body on the line below.
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &d)

	cityDistrictMap := map[string][]models.District{}
	for i := 0; i < len(d.Data); i++ {
		if err := validate.Struct(&d.Data[i]); err != nil {
			continue
		}
		kodeString := strings.Split(d.Data[i].Code, ".")
		kodeKota := fmt.Sprintf("%s.%s", kodeString[0], kodeString[1])
		kodeKecamatan := fmt.Sprintf("%s.%s.%s", kodeString[0], kodeString[1], kodeString[2])
		district := models.District{
			CoreInfo: models.CoreInfo{
				Location: models.Location{
					Latitude:  d.Data[i].Latitude,
					Longitude: d.Data[i].Longitude,
				},
				Level: "KECAMATAN",
				Name:  d.Data[i].Nama,
				Code:  d.Data[i].Code,
			},
			Villages: districtVillageMap[kodeKecamatan],
		}
		if _, ok := cityDistrictMap[kodeKota]; !ok {
			cityDistrictMap[kodeKota] = []models.District{district}
		} else {
			cityDistrictMap[kodeKota] = append(cityDistrictMap[kodeKota], district)
		}
	}

	// city
	var c models.DataCityReq
	resp, _ = http.Get("https://satudata.jabarprov.go.id/api-backend/bigdata/diskominfo/od_kode_wilayah_dan_nama_wilayah_kota_kabupaten?limit=27")

	//We Read the response body on the line below.
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &c)

	cities := []models.City{}
	for i := 0; i < len(c.Data); i++ {
		if err := validate.Struct(&c.Data[i]); err != nil {
			continue
		}
		kodeKota := fmt.Sprintf("%.2f", c.Data[i].Code)
		cities = append(cities, models.City{
			CoreInfo: models.CoreInfo{
				Location: models.Location{
					Longitude: c.Data[i].Longitude,
					Latitude:  c.Data[i].Latitude,
				},
				Level: "KOTA",
				Name:  c.Data[i].Nama,
				Code:  kodeKota,
			},
			Districts: cityDistrictMap[kodeKota],
		})
	}

	fmt.Printf("%v\n\n", cities[0].CoreInfo)

	for _, d := range cities[0].Districts {
		fmt.Printf("\t%v\n\n", d.CoreInfo)
		for _, v := range d.Villages {
			fmt.Printf("\t\t%v\n\n", v.CoreInfo)
		}
	}

	jabarData := &models.JabarData{
		Cities: cities,
	}

	fmt.Println(jabarData)
}
