package models

// unmarshallable define type that would satisfy to unmarshall byte data from API
type Unmarshallable interface {
	DataVillageReq | DataDistrictReq | DataCityReq
}

// define region type, used in mapping
type Regionable interface {
	Village | District | City
}
