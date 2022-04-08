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

	// need to be seperated due different JSON struct tag
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
