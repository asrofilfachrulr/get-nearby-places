package models

type (
	VillageRequest struct {
		Name      string `json:"kemendagri_kelurahan_nama" validate:"required"`
		Level     string
		Code      string  `json:"kemendagri_kelurahan_kode" validate:"required"`
		Latitude  float64 `json:"latitude" validate:"required"`
		Longitude float64 `json:"longitude" validate:"required"`
	}

	DistrictRequest struct {
		Name      string `json:"kemendagri_kecamatan_nama" validate:"required"`
		Level     string
		Code      string  `json:"kemendagri_kecamatan_kode" validate:"required"`
		Latitude  float64 `json:"latitude" validate:"required"`
		Longitude float64 `json:"longitude" validate:"required"`
	}

	CityRequest struct {
		Name      string `json:"kemendagri_kota_nama" validate:"required"`
		Level     string
		Code      float64 `json:"kemendagri_kota_kode" validate:"required"`
		Latitude  float64 `json:"latitude" validate:"required"`
		Longitude float64 `json:"longitude" validate:"required"`
	}

	DataVillageReq struct {
		Data []VillageRequest `json:"data" validate:"required"`
	}

	DataDistrictReq struct {
		Data []DistrictRequest `json:"data" validate:"required"`
	}
	DataCityReq struct {
		Data []CityRequest `json:"data" validate:"required"`
	}
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
	JabarData struct {
		Cities    []City
		NCity     uint
		NDistrict uint
		NVillage  uint
	}
)
